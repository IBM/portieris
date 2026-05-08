// Copyright 2018, 2026 Portieris Authors.
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

// Test_ZeroReplicaEnforcement verifies that admission policies are enforced
// on resources with zero replicas, ensuring they cannot bypass security checks.
func Test_ZeroReplicaEnforcement(t *testing.T) {
	utils.CheckIfTesting(t, testGeneric)

	t.Run("Policy enforced on Deployment with zero replicas - unsigned image denied", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		// This should be denied because the image is unsigned, even though replicas=0
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned-zero-replicas.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on ReplicaSet with zero replicas - unsigned image denied", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		// This should be denied because the image is unsigned, even though replicas=0
		utils.TestReplicaSetNotRunnable(t, framework, "./testdata/replicaset/global-nginx-unsigned-zero-replicas.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on StatefulSet with zero replicas - unsigned image denied", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		// This should be denied because the image is unsigned, even though replicas=0
		utils.TestStatefulSetNotRunnable(t, framework, "./testdata/statefulset/global-nginx-unsigned-zero-replicas.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on ReplicationController with zero replicas - unsigned image denied", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		// This should be denied because the image is unsigned, even though replicas=0
		utils.TestReplicationControllerNotRunnable(t, framework, "./testdata/replicationcontroller/global-nginx-unsigned-zero-replicas.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
}

// Made with Bob
