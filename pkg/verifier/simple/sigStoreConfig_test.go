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
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestCreateRegistryFile(t *testing.T) {

	tests := []struct {
		name           string
		sigStore       string
		sigUser        string
		sigPassword    string
		expectedConfig string
	}{
		{
			name:           "url without credentials",
			sigStore:       "http://foo.com/x",
			sigUser:        "",
			sigPassword:    "",
			expectedConfig: "http://foo.com/x",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirName, fileName, err := createRegistryFile(tt.sigStore, tt.sigUser, tt.sigPassword)
			assert.NoError(t, err)
			assert.Contains(t, dirName, os.TempDir())
			assert.DirExists(t, dirName)

			assert.Contains(t, fileName, dirName)
			assert.FileExists(t, fileName)

			bytes, err := ioutil.ReadFile(fileName)
			assert.NoError(t, err)
			assert.NotZero(t, bytes)

			//fmt.Printf("content:\n%s\n", bytes)
			out := config{}
			err = yaml.Unmarshal(bytes, &out)
			assert.NoError(t, err)
			assert.NotZero(t, out)
			assert.Equal(t, out.DefaultDocker.SigStore, tt.expectedConfig)

			err = removeRegistryFile(dirName)
			assert.NoError(t, err)
		})
	}
}
