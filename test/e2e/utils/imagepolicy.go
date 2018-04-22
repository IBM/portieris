// Copyright 2018 IBM
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

	uuid "github.com/satori/go.uuid"

	"github.com/IBM/portieris/test/framework"
	corev1 "k8s.io/api/core/v1"
)

func CreateImagePolicyInstalledNamespace(t *testing.T, fw *framework.Framework, manifestPath string) *corev1.Namespace {
	ns := uuid.NewV4().String()
	imagePolicy, err := fw.LoadImagePolicyManifest(manifestPath)
	if err != nil {
		t.Fatalf("error loading %q ImagePolicy manifest: %v", imagePolicy.Name, err)
	}
	namespace, err := fw.CreateNamespaceWithIPS(ns)
	if err != nil {
		t.Fatalf("error creating %q namespace: %v", ns, err)
	}
	if err := fw.CreateImagePolicy(ns, imagePolicy); err != nil {
		t.Fatalf("error creating %q ImagePolicy: %v", imagePolicy.Name, err)
	}

	return namespace
}

func CleanUpImagePolicyTest(t *testing.T, fw *framework.Framework, namespace string) {
	if err := fw.DeleteNamespace(namespace); err != nil {
		t.Logf("failed to delete namespace %q: %v", namespace, err)
	}
}

func UpdateImagePolicy(t *testing.T, fw *framework.Framework, manifestPath, namespace, oldPolicy string) {
	imagePolicy, err := fw.LoadImagePolicyManifest(manifestPath)
	if err != nil {
		t.Fatalf("error loading %q ImagePolicy manifest: %v", imagePolicy.Name, err)
	}
	if err := fw.DeleteImagePolicy(oldPolicy, namespace); err != nil {
		t.Fatalf("error updating %q ImagePolicy: %v", imagePolicy.Name, err)
	}
	if err := fw.CreateImagePolicy(namespace, imagePolicy); err != nil {
		t.Fatalf("error updating %q ImagePolicy: %v", imagePolicy.Name, err)
	}
}
