// Copyright 2018 Portieris Authors.
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
	"encoding/json"
	"fmt"
	"os"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Auths struct contains an embedded RegistriesStruct of name auths
type Auths struct {
	Registries RegistriesStruct `json:"auths"`
}

// RegistriesStruct is a map of registries
type RegistriesStruct map[string]struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Auth     string `json:"auth"`
}

// GetSecretToken retrieve the token (password field) for the given namespace/secret/registry
func (w *Wrapper) GetSecretToken(namespace, secretName, registry string, defaultPull bool) (string, string, error) {
	// glog.Infof("getSecretToken << : namespace(%s) secret(%s) registry(%s)", namespace, secretName, registry)
	var username, password string

	// TODO remove hard coded value after local testing
	if defaultPull {
		if os.Getenv("NAMESPACE") == "" || os.Getenv("SECRETNAME") == "" {
			errMessage := "Default namespace and/or secretname details not present to fetch token"
			glog.Errorf(errMessage)
			return username, password, fmt.Errorf(errMessage)
		}
		// fetching the default namespace and secret name for image pull secrets if nothing has been mentioned and relying on the
		// dockerconfig file of the host system
		namespace = os.Getenv("NAMESPACE")
		secretName = os.Getenv("SECRETNAME")
	}

	glog.Infof("getSecretToken << : namespace(%s) secret(%s) registry(%s)", namespace, secretName, registry)

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
		username = login.Username
		password = login.Password
	} else {
		err = fmt.Errorf("Secret not defined for registry: %s", registry)
	}
	// glog.Infof("getSecretToken >> : token(%s)", token)
	return username, password, err
}
