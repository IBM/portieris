// Copyright 2020 Portieris Authors.
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

package simple

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
)

type regConfig struct {
	SigStore string `yaml:"sigstore"`
}
type config struct {
	DefaultDocker regConfig `yaml:"default-docker"`
}

// CreateRegistryDir write a file in a new directory containing the desired default docker configuration
func CreateRegistryDir(storeURL, storeUser, storePassword string) (string, error) {
	if storeURL == "" {
		glog.Infof("No lookaside signature store.")
		return "", nil
	}
	method := ""
	rest := ""
	if strings.HasPrefix(storeURL, "http://") {
		method = "http://"
		rest = strings.TrimPrefix(storeURL, "http://")
	}
	if strings.HasPrefix(storeURL, "https://") {
		method = "https://"
		rest = strings.TrimPrefix(storeURL, "https://")
	}

	// allow only http:// and https://
	if method == "" {
		return "", fmt.Errorf("expecting https:// or http:// URL, got: %s", storeURL)
	}
	// insert credentials as <method>user:password@<rest>
	if storeUser == "" {
		glog.Infof("Lookaside signature store at: %s", storeURL)
	} else {
		glog.Infof("Lookaside signature store at: %s%s:***@%s", method, storeUser, rest)
		storeURL = fmt.Sprintf("%s%s:%s@%s", method, storeUser, storePassword, rest)
	}

	dir, err := ioutil.TempDir("", "registry.d")
	if err != nil {
		return "", err
	}
	file, err := os.OpenFile(dir+"/default.yaml", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		os.RemoveAll(dir)
		return "", err
	}

	rConf := config{
		DefaultDocker: regConfig{
			SigStore: storeURL,
		},
	}
	bytes, err := yaml.Marshal(rConf)
	if err != nil {
		os.RemoveAll(dir)
		return "", err
	}

	_, err = file.Write(bytes)
	if err != nil {
		os.RemoveAll(dir)
		return "", err
	}
	glog.Infof("Using store config dir: %s", dir)
	return dir, nil
}

// RemoveRegistryDir .
func RemoveRegistryDir(dirName string) error {
	if dirName == "" {
		return nil
	}
	return os.RemoveAll(dirName)
}
