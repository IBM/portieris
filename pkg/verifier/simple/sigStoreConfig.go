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

	"gopkg.in/yaml.v3"
)

type sigConfig struct {
	SigStore string `json:"sigstore"`
}
type config struct {
	DefaultDocker sigConfig `json:"default-docker"`
}

// CreateRegistryFile write a file in a new directory containing the desired default docker configuration
func createRegistryFile(sigStore, sigUser, sigPassword string) (string, string, error) {
	dir, err := ioutil.TempDir("", "registry.d")
	if err != nil {
		return "", "", err
	}
	file, err := os.OpenFile(dir+"/default.yaml", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		os.RemoveAll(dir)
		return "", "", err
	}

	rConf := config{
		DefaultDocker: sigConfig{
			SigStore: sigStore,
		},
	}
	bytes, err := yaml.Marshal(rConf)
	if err != nil {
		os.RemoveAll(dir)
		return "", "", err
	}

	_, err = file.Write(bytes)
	if err != nil {
		os.RemoveAll(dir)
		return "", "", err
	}
	return dir, file.Name(), nil
}

func removeRegistryFile(dirName string) error {
	return os.RemoveAll(dirName)
}
