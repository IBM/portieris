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

package e2e

import (
	"testing"

	"github.com/IBM/portieris/test/e2e/utils"
)

func TestNotary_ImagePolicyRepositories_AllowAllDenyAll(t *testing.T) {
	utils.CheckIfTesting(t, testTrustImagePolicy)
	if defaultClusterPolicy := utils.DeleteThenReturnClusterImagePolicy(t, framework, "ibmcloud-default-cluster-image-policy"); defaultClusterPolicy != nil {
		defer framework.CreateClusterImagePolicy(defaultClusterPolicy)
	}

	t.Run("Allow all images", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-all.yaml")
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

func TestNotary_ImagePolicyRepositories_BasicTrust(t *testing.T) {
	utils.CheckIfTesting(t, testTrustImagePolicy)
	t.Run("Allow signed images when trust enabled", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Deny unsigned images when trust enabled", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
}

func TestNotary_ImagePolicyRepositories_TrustPinning(t *testing.T) {
	utils.CheckIfTesting(t, testTrustImagePolicy)
	t.Run("Allow images signed by the correct single signer when trust enabled", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/pinned-signer1.yaml")
		utils.CreateSecret(t, framework, "./testdata/secret/signer1.pubkey.yaml", namespace.Name)
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Allow images signed the correct multiple signers and when trust enabled", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/pinned-multi.yaml")
		utils.CreateSecret(t, framework, "./testdata/secret/signer1.pubkey.yaml", namespace.Name)
		utils.CreateSecret(t, framework, "./testdata/secret/signer2.pubkey.yaml", namespace.Name)
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-multisigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Deny images signed by the wrong signer when trust enabled", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/pinned-signer2.yaml")
		utils.CreateSecret(t, framework, "./testdata/secret/signer2.pubkey.yaml", namespace.Name)
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Deny images signed by a single signer when multiple are required when trust enabled", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/pinned-multi.yaml")
		utils.CreateSecret(t, framework, "./testdata/secret/signer1.pubkey.yaml", namespace.Name)
		utils.CreateSecret(t, framework, "./testdata/secret/signer2.pubkey.yaml", namespace.Name)
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
}

func TestNotary_ImagePolicyRepositories_TrustPinningMultiContainers(t *testing.T) {
	utils.CheckIfTesting(t, testTrustImagePolicy)
	t.Run("Allow when both containers fulfill the policy", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-signed-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Deny when one container fails to fulfill the policy", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-signed-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
}

// Temporary until we support other registries
func TestNotary_ImagePolicyRepositories_ThirdPartyTrust(t *testing.T) {
	utils.CheckIfTesting(t, testTrustImagePolicy)
	t.Run("Third party trust is rejected", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-dockerhub.yaml")
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/dockerhub-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
}

func TestNotary_ImagePolicyRepositories_OverridesClusterPolicies(t *testing.T) {
	utils.CheckIfTesting(t, testTrustImagePolicy)
	t.Run("Unsigned image is rejected if cluster policy allows all but image policy only allows signed", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		clusterImagePolicy, redundantNamespace := utils.CreateClusterImagePolicyAndNamespace(t, framework, "./testdata/clusterimagepolicy/allow-all.yaml")
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/dockerhub-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
		utils.CleanUpClusterImagePolicyTest(t, framework, clusterImagePolicy.Name, redundantNamespace.Name)
	})
}
