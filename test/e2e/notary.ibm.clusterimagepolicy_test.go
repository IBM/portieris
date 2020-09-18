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

func TestNotary_ClusterImagePolicyRepositories_AllowAllDenyAll(t *testing.T) {
	utils.CheckIfTesting(t, testTrustClusterImagePolicy)
	if defaultClusterPolicy := utils.DeleteThenReturnClusterImagePolicy(t, framework, "default"); defaultClusterPolicy != nil {
		defer framework.CreateClusterImagePolicy(defaultClusterPolicy)
	}

	t.Run("Allow all images", func(t *testing.T) {
		clusterImagePolicy, namespace := utils.CreateClusterImagePolicyAndNamespace(t, framework, "./testdata/clusterimagepolicy/allow-all.yaml")
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpClusterImagePolicyTest(t, framework, clusterImagePolicy.Name, namespace.Name)
	})
	t.Run("Deny all images when no cluster image policy is present", func(t *testing.T) {
		namespace, err := framework.CreateNamespaceWithIPS("deny-all")
		if err != nil {
			t.Fatalf("error creating deny-all namespace: %v", err)
		}
		defer framework.DeleteNamespace(namespace.Name)
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
	})
}

func TestNotary_ClusterImagePolicyRepositories_BasicTrust(t *testing.T) {
	utils.CheckIfTesting(t, testTrustClusterImagePolicy)
	if defaultClusterPolicy := utils.DeleteThenReturnClusterImagePolicy(t, framework, "default"); defaultClusterPolicy != nil {
		defer framework.CreateClusterImagePolicy(defaultClusterPolicy)
	}

	t.Run("Allow signed images when trust enabled", func(t *testing.T) {
		clusterImagePolicy, namespace := utils.CreateClusterImagePolicyAndNamespace(t, framework, "./testdata/clusterimagepolicy/allow-signed.yaml")
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpClusterImagePolicyTest(t, framework, clusterImagePolicy.Name, namespace.Name)
	})
	t.Run("Deny unsigned images when trust enabled", func(t *testing.T) {
		clusterImagePolicy, namespace := utils.CreateClusterImagePolicyAndNamespace(t, framework, "./testdata/clusterimagepolicy/allow-signed.yaml")
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpClusterImagePolicyTest(t, framework, clusterImagePolicy.Name, namespace.Name)
	})
}

func TestNotary_ClusterImagePolicyRepositories_TrustPinning(t *testing.T) {
	utils.CheckIfTesting(t, testTrustClusterImagePolicy)
	if defaultClusterPolicy := utils.DeleteThenReturnClusterImagePolicy(t, framework, "default"); defaultClusterPolicy != nil {
		defer framework.CreateClusterImagePolicy(defaultClusterPolicy)
	}

	t.Run("Allow images signed by the correct single signer when trust enabled", func(t *testing.T) {
		clusterImagePolicy, namespace := utils.CreateClusterImagePolicyAndNamespace(t, framework, "./testdata/clusterimagepolicy/pinned-signer1.yaml")
		utils.CreateSecret(t, framework, "./testdata/secret/signer1.pubkey.yaml", namespace.Name)
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpClusterImagePolicyTest(t, framework, clusterImagePolicy.Name, namespace.Name)
	})
	t.Run("Allow images signed the correct multiple signers and when trust enabled", func(t *testing.T) {
		clusterImagePolicy, namespace := utils.CreateClusterImagePolicyAndNamespace(t, framework, "./testdata/clusterimagepolicy/pinned-multi.yaml")
		utils.CreateSecret(t, framework, "./testdata/secret/signer1.pubkey.yaml", namespace.Name)
		utils.CreateSecret(t, framework, "./testdata/secret/signer2.pubkey.yaml", namespace.Name)
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-multisigned.yaml", namespace.Name)
		utils.CleanUpClusterImagePolicyTest(t, framework, clusterImagePolicy.Name, namespace.Name)
	})
	t.Run("Deny images signed by the wrong signer when trust enabled", func(t *testing.T) {
		clusterImagePolicy, namespace := utils.CreateClusterImagePolicyAndNamespace(t, framework, "./testdata/clusterimagepolicy/pinned-signer2.yaml")
		utils.CreateSecret(t, framework, "./testdata/secret/signer2.pubkey.yaml", namespace.Name)
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpClusterImagePolicyTest(t, framework, clusterImagePolicy.Name, namespace.Name)
	})
	t.Run("Deny images signed by a single signer when multiple are required when trust enabled", func(t *testing.T) {
		clusterImagePolicy, namespace := utils.CreateClusterImagePolicyAndNamespace(t, framework, "./testdata/clusterimagepolicy/pinned-multi.yaml")
		utils.CreateSecret(t, framework, "./testdata/secret/signer1.pubkey.yaml", namespace.Name)
		utils.CreateSecret(t, framework, "./testdata/secret/signer2.pubkey.yaml", namespace.Name)
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpClusterImagePolicyTest(t, framework, clusterImagePolicy.Name, namespace.Name)
	})
}

func TestNotary_ClusterImagePolicyRepositories_TrustPinningMultiContainers(t *testing.T) {
	utils.CheckIfTesting(t, testTrustClusterImagePolicy)
	t.Run("Allow when both containers fulfill the policy", func(t *testing.T) {
		clusterImagePolicy, namespace := utils.CreateClusterImagePolicyAndNamespace(t, framework, "./testdata/clusterimagepolicy/allow-signed.yaml")
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-signed-signed.yaml", namespace.Name)
		utils.CleanUpClusterImagePolicyTest(t, framework, clusterImagePolicy.Name, namespace.Name)
	})
	t.Run("Deny when one container fails to fulfill the policy", func(t *testing.T) {
		clusterImagePolicy, namespace := utils.CreateClusterImagePolicyAndNamespace(t, framework, "./testdata/clusterimagepolicy/allow-signed.yaml")
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-signed-unsigned.yaml", namespace.Name)
		utils.CleanUpClusterImagePolicyTest(t, framework, clusterImagePolicy.Name, namespace.Name)
	})
}

// Temporary until we support other registries
func TestNotary_ClusterImagePolicyRepositories_ThirdPartyTrust(t *testing.T) {
	utils.CheckIfTesting(t, testTrustClusterImagePolicy)
	t.Run("Third party trust is rejected", func(t *testing.T) {
		clusterImagePolicy, namespace := utils.CreateClusterImagePolicyAndNamespace(t, framework, "./testdata/clusterimagepolicy/allow-signed-dockerhub.yaml")
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/dockerhub-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpClusterImagePolicyTest(t, framework, clusterImagePolicy.Name, namespace.Name)
	})
}
