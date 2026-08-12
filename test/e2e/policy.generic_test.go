// Copyright 2018, 2023 Portieris Authors.
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func Test_JobTypesSuccess(t *testing.T) {
	utils.CheckIfTesting(t, testGeneric)

	t.Run("Policy enforced on Deployment", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on DaemonSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		utils.TestDaemonSetRunnable(t, framework, "./testdata/daemonset/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on ReplicaSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		utils.TestReplicaSetRunnable(t, framework, "./testdata/replicaset/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on ReplicationController", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		utils.TestReplicationControllerRunnable(t, framework, "./testdata/replicationcontroller/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on Pod", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		utils.TestPodRunnable(t, framework, "./testdata/pod/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on StatefulSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		utils.TestStatefulSetRunnable(t, framework, "./testdata/statefulset/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on Job", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		utils.TestJobRunnable(t, framework, "./testdata/job/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on CronJob", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		utils.TestCronJobRunnable(t, framework, "./testdata/cronjob/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
}

// TODO uncomment when issue #34 is addressed.
/* func Test_JobTypesSuccessCustomTrustServer(t *testing.T) {
	utils.CheckIfTesting(t, testGeneric)

	t.Run("Policy enforced on Deployment", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on DaemonSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestDaemonSetRunnable(t, framework, "./testdata/daemonset/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on ReplicaSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestReplicaSetRunnable(t, framework, "./testdata/replicaset/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on ReplicationController", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestReplicationControllerRunnable(t, framework, "./testdata/replicationcontroller/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on Pod", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestPodRunnable(t, framework, "./testdata/pod/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on StatefulSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestStatefulSetRunnable(t, framework, "./testdata/statefulset/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on Job", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestJobRunnable(t, framework, "./testdata/job/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on CronJob", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestCronJobRunnable(t, framework, "./testdata/cronjob/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
}

   func Test_JobTypesFailCustomTrustServer(t *testing.T) {
	utils.CheckIfTesting(t, testGeneric)

	t.Run("Policy enforced on Deployment", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on DaemonSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestDaemonSetNotRunnable(t, framework, "./testdata/daemonset/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on ReplicaSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestReplicaSetNotRunnable(t, framework, "./testdata/replicaset/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on ReplicationController", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestReplicationControllerNotRunnable(t, framework, "./testdata/replicationcontroller/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on Pod", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestPodNotRunnable(t, framework, "./testdata/pod/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on StatefulSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestStatefulSetNotRunnable(t, framework, "./testdata/statefulset/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on Job", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestJobNotRunnable(t, framework, "./testdata/job/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on CronJob", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml", "")
		utils.TestCronJobNotRunnable(t, framework, "./testdata/cronjob/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

} */

func Test_JobTypesFail(t *testing.T) {
	utils.CheckIfTesting(t, testGeneric)

	t.Run("Policy enforced on Deployment", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on DaemonSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.TestDaemonSetNotRunnable(t, framework, "./testdata/daemonset/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on ReplicaSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.TestReplicaSetNotRunnable(t, framework, "./testdata/replicaset/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on ReplicationController", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.TestReplicationControllerNotRunnable(t, framework, "./testdata/replicationcontroller/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on Pod", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.TestPodNotRunnable(t, framework, "./testdata/pod/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on StatefulSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.TestStatefulSetNotRunnable(t, framework, "./testdata/statefulset/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on Job", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.TestJobNotRunnable(t, framework, "./testdata/job/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on CronJob", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.TestCronJobNotRunnable(t, framework, "./testdata/cronjob/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

}

func Test_OperationsSucces(t *testing.T) {
	utils.CheckIfTesting(t, testGeneric)

	t.Run("Policy enforced on patch", func(t *testing.T) {
		t.Parallel()

		patchString := `{
			"spec": {
			   "template": {
				  "spec": {
					 "containers": [
						{
						   "image": "icr.io/cise/nginx:unsigned"
						}
					 ]
				  }
			   }
			}
		 }`
		// Create a namespace and policy to allow all.
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		// Start the deployment.
		deploymentName := utils.TestStartDeployNoDelete(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		// Change the policy to deny.
		utils.TestDeploymentNotRunnableOnPatch(t, framework, deploymentName, patchString, namespace.Name)

		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on replace", func(t *testing.T) {
		t.Parallel()
		// Create a namespace and policy to allow all.
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")
		utils.CreateSecret(t, framework, "./testdata/secret/sh.pubkey.yaml", namespace.Name)
		// Start the deployment.
		_ = utils.TestStartDeployNoDelete(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		// Change the policy to deny.
		utils.TestDeploymentNotRunnableOnReplace(t, framework, "./testdata/deployment/global-signed-patch-to-unsigned.yaml", namespace.Name)

		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

}

// Test_FabricatedOwnerReferenceBypass verifies that a pod with a fabricated
// ownerReference is still subject to image policy enforcement.
// A pod whose ownerReference points to a non-existent resource must be denied
// by a deny-all policy — not unconditionally allowed.
func Test_FabricatedOwnerReferenceBypass(t *testing.T) {
	utils.CheckIfTesting(t, testGeneric)

	t.Run("Unsigned pod with fabricated ownerReference is denied by policy", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml", "")

		pod, err := framework.LoadPodManifest("./testdata/pod/global-nginx-unsigned.yaml")
		if err != nil {
			t.Fatalf("failed to load pod manifest: %v", err)
		}

		// Fabricate an ownerReference pointing to a non-existent ReplicaSet.
		// This is the bypass vector: Kubernetes accepts it on CREATE without
		// verifying existence; Portieris must still enforce policy.
		pod.Name = "bypass-ownerref-test"
		pod.OwnerReferences = []metav1.OwnerReference{
			{
				APIVersion: "apps/v1",
				Kind:       "ReplicaSet",
				Name:       "does-not-exist",
				UID:        types.UID("aaaa-bbbb-cccc-dddd"),
			},
		}

		if err := framework.CreatePod(namespace.Name, pod); err == nil {
			defer framework.DeletePod(pod.Name, namespace.Name)
			t.Error("expected pod with fabricated ownerReference to be denied, but it was admitted")
			utils.DumpEvents(t, framework, namespace.Name)
			utils.DumpPolicies(t, framework, namespace.Name)
		}

		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
}
