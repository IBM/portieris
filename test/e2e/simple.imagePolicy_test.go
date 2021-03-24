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

package e2e

import (
	"testing"

	"github.com/IBM/portieris/test/e2e/utils"
)

func TestSimple_ImagePolicyRepositories_AllowAllDenyAll(t *testing.T) {
	utils.CheckIfTesting(t, testSimpleImagePolicy)
	if defaultClusterPolicy := utils.DeleteThenReturnClusterImagePolicy(t, framework, "default"); defaultClusterPolicy != nil {
		defer framework.CreateClusterImagePolicy(defaultClusterPolicy)
	}

	t.Run("Allow all images", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/simple-accept-anything.yaml", "")
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Deny all images when no image policy is present", func(t *testing.T) {
		t.Parallel()
		namespace, err := framework.CreateNamespaceWithIPS("deny-all")
		if err != nil {
			t.Fatalf("error creating deny-all namespace: %v", err)
		}
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
}

func TestSimple_ImagePolicyRepositories_Basic(t *testing.T) {
	utils.CheckIfTesting(t, testSimpleImagePolicy)
	t.Run("Allow images signed by the correct single simple signer", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/simple-signedby1.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/simple1pubkey.yaml", namespace.Name)
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-simple.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Allow images signed the correct multiple simple signers", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/simple-signedby2.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/simple1pubkey.yaml", namespace.Name)
		utils.CreateSecret(t, framework, "./testdata/secret/simple2pubkey.yaml", namespace.Name)
		utils.TestDeploymentRunnableCheck(t, framework, "./testdata/deployment/global-nginx-simple.yaml", namespace.Name, utils.VerifySha)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Allow images signed the correct multiple simple signers with no mutation", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/simple-signedby2-nomutate.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/simple1pubkey.yaml", namespace.Name)
		utils.CreateSecret(t, framework, "./testdata/secret/simple2pubkey.yaml", namespace.Name)
		utils.TestDeploymentRunnableCheck(t, framework, "./testdata/deployment/global-nginx-simple.yaml", namespace.Name, utils.VerifyNoSha)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Allow images signed the correct multiple simple signers with explicit mutation", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/simple-signedby2-mutate.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/simple1pubkey.yaml", namespace.Name)
		utils.CreateSecret(t, framework, "./testdata/secret/simple2pubkey.yaml", namespace.Name)
		utils.TestDeploymentRunnableCheck(t, framework, "./testdata/deployment/global-nginx-simple.yaml", namespace.Name, utils.VerifySha)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Allow images signed by the correct single simple signer with a secret namespace override", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/simple-signedby1-keysecret-namespace-override.yaml", "secretnamespace")
		utils.CreateSecret(t, framework, "./testdata/secret/simple1pubkey.yaml", namespace.Name)
		utils.TestDeploymentRunnableCheck(t, framework, "./testdata/deployment/global-nginx-simple.yaml", namespace.Name, utils.VerifySha)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Deny images signed by the wrong single simple signer", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/simple-signedby1.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/simple1pubkey.yaml", namespace.Name)
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-another.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Deny images signed by a single simple signer when multiple are required", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/simple-signedby2.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/simple1pubkey.yaml", namespace.Name)
		utils.CreateSecret(t, framework, "./testdata/secret/simple2pubkey.yaml", namespace.Name)
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-another.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Allow images matched with remapIdentity policy", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/simple-remap.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/simple1pubkey.yaml", namespace.Name)
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-remapped.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Allow images without a pullSecret", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespaceNoSecrets(t, framework, "./testdata/imagepolicy/simple-accept-anything.yaml")
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-another.yaml", namespace.Name)
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/dockerhub-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Allow images which require a signature without a pullSecret", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespaceNoSecrets(t, framework, "./testdata/imagepolicy/simple-signedby1.yaml")
		utils.CreateSecret(t, framework, "./testdata/secret/simple1pubkey.yaml", namespace.Name)
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-simple.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
}
