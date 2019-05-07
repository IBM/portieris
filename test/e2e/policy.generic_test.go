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

func Test_JobTypesSuccess(t *testing.T) {
	utils.CheckIfTesting(t, testGeneric)

	t.Run("Policy enforced on Deployment", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on DaemonSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestDaemonSetRunnable(t, framework, "./testdata/daemonset/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on ReplicaSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestReplicaSetRunnable(t, framework, "./testdata/replicaset/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on ReplicationController", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestReplicationControllerRunnable(t, framework, "./testdata/replicationcontroller/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on Pod", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestPodRunnable(t, framework, "./testdata/pod/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on StatefulSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestStatefulSetRunnable(t, framework, "./testdata/statefulset/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on Job", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestJobRunnable(t, framework, "./testdata/job/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on CronJob", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestCronJobRunnable(t, framework, "./testdata/cronjob/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
}

// TODO uncomment when issue #34 has been addressed
/* func Test_JobTypesSuccessCustomTrustServer(t *testing.T) {
	utils.CheckIfTesting(t, testGeneric)

	t.Run("Policy enforced on Deployment", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on DaemonSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestDaemonSetRunnable(t, framework, "./testdata/daemonset/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on ReplicaSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestReplicaSetRunnable(t, framework, "./testdata/replicaset/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on ReplicationController", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestReplicationControllerRunnable(t, framework, "./testdata/replicationcontroller/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on Pod", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestPodRunnable(t, framework, "./testdata/pod/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on StatefulSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestStatefulSetRunnable(t, framework, "./testdata/statefulset/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on Job", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestJobRunnable(t, framework, "./testdata/job/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on CronJob", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestCronJobRunnable(t, framework, "./testdata/cronjob/global-nginx-signed.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
}

   func Test_JobTypesFailCustomTrustServer(t *testing.T) {
	utils.CheckIfTesting(t, testGeneric)

	t.Run("Policy enforced on Deployment", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on DaemonSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestDaemonSetNotRunnable(t, framework, "./testdata/daemonset/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on ReplicaSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestReplicaSetNotRunnable(t, framework, "./testdata/replicaset/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on ReplicationController", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestReplicationControllerNotRunnable(t, framework, "./testdata/replicationcontroller/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on Pod", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestPodNotRunnable(t, framework, "./testdata/pod/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on StatefulSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestStatefulSetNotRunnable(t, framework, "./testdata/statefulset/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on Job", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestJobNotRunnable(t, framework, "./testdata/job/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on CronJob", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed-custom.yaml")
		utils.TestCronJobNotRunnable(t, framework, "./testdata/cronjob/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

} */

func Test_JobTypesFail(t *testing.T) {
	utils.CheckIfTesting(t, testGeneric)

	t.Run("Policy enforced on Deployment", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on DaemonSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestDaemonSetNotRunnable(t, framework, "./testdata/daemonset/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on ReplicaSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestReplicaSetNotRunnable(t, framework, "./testdata/replicaset/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on ReplicationController", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestReplicationControllerNotRunnable(t, framework, "./testdata/replicationcontroller/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on Pod", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestPodNotRunnable(t, framework, "./testdata/pod/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on StatefulSet", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestStatefulSetNotRunnable(t, framework, "./testdata/statefulset/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on Job", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestJobNotRunnable(t, framework, "./testdata/job/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Policy enforced on CronJob", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		utils.TestCronJobNotRunnable(t, framework, "./testdata/cronjob/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

}

func Test_OperationsSucces(t *testing.T) {
	utils.CheckIfTesting(t, testGeneric)

	t.Run("Policy not enforced on child resource (pod)", func(t *testing.T) {
		t.Parallel()
		// Create namespace and policy to allow all
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-all.yaml")
		//Start deployment
		deploymentName := utils.TestStartDeployNoDelete(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
		// Change policy to deny
		utils.UpdateImagePolicy(t, framework, "./testdata/imagepolicy/allow-signed.yaml", namespace.Name, "allow-all")
		// Kill Pod
		utils.KillPod(t, framework, namespace.Name)
		// Check pod comes back
		utils.TestCurrentDeployStatus(t, framework, namespace.Name, deploymentName)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

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
		// Create namespace and policy to allow all
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		//Start deployment
		deploymentName := utils.TestStartDeployNoDelete(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		// Change policy to deny
		utils.TestDeploymentNotRunnableOnPatch(t, framework, deploymentName, patchString, namespace.Name)

		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

	t.Run("Policy enforced on replace", func(t *testing.T) {
		t.Parallel()
		// Create namespace and policy to allow all
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-signed.yaml")
		//Start deployment
		_ = utils.TestStartDeployNoDelete(t, framework, "./testdata/deployment/global-nginx-signed.yaml", namespace.Name)
		// Change policy to deny
		utils.TestDeploymentNotRunnableOnReplace(t, framework, "./testdata/deployment/global-signed-patch-to-unsigned.yaml", namespace.Name)

		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

}
