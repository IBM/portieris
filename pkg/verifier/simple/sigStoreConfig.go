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
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
)

type sigConfig struct {
	SigStore string `yaml:"sigstore"`
}
type config struct {
	DefaultDocker sigConfig `yaml:"default-docker"`
}

// CreateRegistryDir write a file in a new directory containing the desired default docker configuration
func CreateRegistryDir(sigStore, sigUser, sigPassword string) (string, error) {
	if sigStore == "" {
		glog.Infof("No lookaside signature store.")
		return "", nil
	}
	glog.Infof("Lookaside signature store at: %s", sigStore)
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
		DefaultDocker: sigConfig{
			SigStore: sigStore,
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
