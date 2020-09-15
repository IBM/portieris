// Copyright 2018, 2020 Portieris Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Auths struct contains an embedded RegistriesStruct of name auths
type Auths struct {
	Registries RegistriesStruct `json:"auths"`
}

// RegistriesStruct is a map of registries to their credentials
type RegistriesStruct map[string]RegistryCredentials

// RegistryCredentials defines the fields stored per registry in an docker config secret
type RegistryCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Auth     string `json:"auth"`
}

// GetSecretKey obtains the "key" data from the named secret
func (w *Wrapper) GetSecretKey(namespace, secretName string) ([]byte, error) {
	// Retrieve secret
	secret, err := w.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
	if err != nil {
		glog.Error("Error: ", err)
		return nil, err
	}
	glog.Infof("Found secret %s", secretName)
	// Obtain the key data
	if key, ok := secret.Data["key"]; ok {
		return key, nil
	}
	return nil, fmt.Errorf("Secret %q in %q does not contain a \"key\" attribute", secretName, namespace)
}

// GetSecretToken retrieve the token (password field) for the given namespace/secret/registry
func (w *Wrapper) GetSecretToken(namespace, secretName, registry string) (string, string, error) {
	// glog.Infof("getSecretToken << : namespace(%s) secret(%s) registry(%s)", namespace, secretName, registry)
	var username, password string

	// Retrieve secret
	secret, err := w.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
	if err != nil {
		glog.Error("Error: ", err)
		return username, password, err
	}

	// Parse the returned data.
	auths := Auths{}
	if secretData, ok := secret.Data[".dockerconfigjson"]; ok {
		if err := json.Unmarshal(secretData, &auths); err != nil {
			glog.Errorf("Error unmarshalling .dockerconfigjson from %s: %v", secretName, err)
			return username, password, err
		}
	} else if dockerCfgData, ok := secret.Data[".dockercfg"]; ok {
		registries := RegistriesStruct{}
		if err := json.Unmarshal(dockerCfgData, &registries); err != nil {
			glog.Errorf("Error unmarshalling .dockercfg from %s: %v", secretName, err)
			return username, password, err
		}
		auths.Registries = registries
	} else {
		return username, password, fmt.Errorf("imagePullSecret %s contains neither .dockercfg nor .dockerconfigjson", secretName)
	}

	// Determine if there is a secret for the specified registry
	registries := auths.Registries
	if login, ok := registries[registry]; ok {
		username, password = w.extractRegistryCredentials(login)
	} else {
		err = fmt.Errorf("Secret %s not defined for registry: %s", secretName, registry)
	}
	// glog.Infof("getSecretToken >> : token(%s)", token)
	return username, password, err
}

func (w *Wrapper) extractRegistryCredentials(creds RegistryCredentials) (username, password string) {
	username = creds.Username
	password = creds.Password

	if creds.Auth == "" {
		return
	}

	decoder := base64.StdEncoding
	if !strings.HasSuffix(strings.TrimSpace(creds.Auth), "=") {
		// Modify the decoder to be raw if no padding is present
		decoder = decoder.WithPadding(base64.NoPadding)
	}

	base64Decoded, err := decoder.DecodeString(creds.Auth)
	if err != nil {
		glog.Warningf("Error Base64 decoding auth field, username/password fields from the registry credentials will be used instead. Error %v", err)
		return
	}

	// SplitN required here so that a colon inside the password is not treated as another delimiter
	splitted := strings.SplitN(string(base64Decoded), ":", 2)
	if len(splitted) != 2 {
		glog.Warning("Decoded auth field was not in the format username:password, the username/password fields from the registry credentials will be used instead.")
		return
	}

	username = splitted[0]
	password = splitted[1]

	return
}

// GetBasicCredentials retrieves username, password from a named secret
func (w *Wrapper) GetBasicCredentials(namespace, name string) (string, string, error) {

	if name == "" {
		return "", "", nil
	}

	// Retrieve secret
	secret, err := w.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}

	username, ok := secret.Data["username"]
	if !ok {
		return "", "", fmt.Errorf("secret: %s, does not contain username", name)
	}

	password, ok := secret.Data["password"]
	if !ok {
		return "", "", fmt.Errorf("secret: %s, does not contain password", name)
	}

	return string(username), string(password), nil
}
