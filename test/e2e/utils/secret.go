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

package utils

import (
	"testing"

	"github.com/IBM/portieris/test/framework"
)

func CreateSecret(t *testing.T, fw *framework.Framework, manifestPath, namespace string) {
	manifest, err := fw.LoadSecretManifest(manifestPath)
	if err != nil {
		t.Fatalf("Error loading secret manifest: %v", err)
	}
	if manifest == nil {
		t.Fatalf("Error loading manifest: manifest is nil")
	}
	if err := fw.CreateSecret(namespace, manifest); err != nil {
		t.Fatalf("Error creating secret %q: %v", manifest.Name, err)
	}
}
