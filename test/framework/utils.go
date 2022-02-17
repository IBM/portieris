// Copyright 2018, 2021 Portieris Authors.
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

package framework

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Framework) getHelmReleaseSelectorListOptions() metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: fmt.Sprintf("release=%v", f.HelmRelease),
	}
}

// MakeTestUUID is a simple wrapper to return a UUID for testing purposes
func MakeTestUUID() string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return u.String()
}

// GenerateTestAnnotation returns a unique test annotation that is used to patch resources.
func (f *Framework) GenerateTestAnnotation() string {
	return fmt.Sprintf(`
		{
			"metadata": {
				"annotations": {
					"test":"%v"
				}
			}
		}
	`, MakeTestUUID())
}

func openFile(relativePath string) (*os.File, error) {
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		return nil, fmt.Errorf("unable to work out absolute path of %q", relativePath)
	}
	return os.Open(absolutePath)
}
