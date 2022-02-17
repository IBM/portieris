// Copyright 2020, 2021 Portieris Authors.
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
	"os"
	"testing"

	"github.com/IBM/portieris/test/framework"
	corev1 "k8s.io/api/core/v1"
)

// CreateImagePolicyInstalledNamespaceOptionalSecrets ...
func CreateImagePolicyInstalledNamespaceOptionalSecrets(t *testing.T, fw *framework.Framework, manifestPath string, ips bool, namespaceParam string) *corev1.Namespace {
	ns := framework.MakeTestUUID()
	if namespaceParam != "" {
		ns = namespaceParam
	}
	imagePolicy, err := fw.LoadImagePolicyManifest(manifestPath)
	if err != nil {
		t.Fatalf("error loading %q ImagePolicy manifest: %v", manifestPath, err)
	}
	for idx := range imagePolicy.Spec.Repositories {
		if imagePolicy.Spec.Repositories[idx].Policy.Vulnerability.ICCRVA.Account == "ENV" {
			imagePolicy.Spec.Repositories[idx].Policy.Vulnerability.ICCRVA.Account = os.Getenv("E2E_ACCOUNT_HEADER")
			if imagePolicy.Spec.Repositories[idx].Policy.Vulnerability.ICCRVA.Account == "" {
				t.Fatalf("Unable to set Account header, did you export E2E_ACCOUNT_HEADER")
			}
			t.Log(imagePolicy.Spec.Repositories[idx].Policy.Vulnerability.ICCRVA.Account)
		}
	}

	var namespace *corev1.Namespace
	if ips {
		namespace, err = fw.CreateNamespaceWithIPS(ns)
		if err != nil {
			t.Fatalf("error creating %q namespace: %v", ns, err)
		}
	} else {
		namespace, err = fw.CreateNamespace(ns)
		if err != nil {
			t.Fatalf("error creating %q namespace: %v", ns, err)
		}
	}

	if err := fw.CreateImagePolicy(ns, imagePolicy); err != nil {
		t.Fatalf("error creating %q ImagePolicy: %v", imagePolicy.Name, err)
	}

	return namespace
}

// CreateImagePolicyInstalledNamespace ...
func CreateImagePolicyInstalledNamespace(t *testing.T, fw *framework.Framework, manifestPath string, ns string) *corev1.Namespace {
	return CreateImagePolicyInstalledNamespaceOptionalSecrets(t, fw, manifestPath, true, ns)
}

// CreateImagePolicyInstalledNamespaceNoSecrets ...
func CreateImagePolicyInstalledNamespaceNoSecrets(t *testing.T, fw *framework.Framework, manifestPath string) *corev1.Namespace {
	return CreateImagePolicyInstalledNamespaceOptionalSecrets(t, fw, manifestPath, false, "")
}

// CleanUpImagePolicyTest ...
func CleanUpImagePolicyTest(t *testing.T, fw *framework.Framework, namespace string) {
	if err := fw.DeleteNamespace(namespace); err != nil {
		t.Logf("failed to delete namespace %q: %v", namespace, err)
	}
}

// UpdateImagePolicy ...
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
