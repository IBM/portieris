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

package e2e

import (
	"testing"

	"github.com/IBM/portieris/test/e2e/utils"

	"github.com/stretchr/testify/assert"
)

func TestIBMContainerService_KubeSystem(t *testing.T) {
	utils.CheckIfTesting(t, testArmada)
	t.Run("Run anything in kube-system", func(t *testing.T) {
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/dockerhub-nginx-unsigned.yaml", "kube-system")
	})

	t.Run("Patch existing ICS deployment with annotation, verify it still works", func(t *testing.T) {
		testAnnotation := framework.GenerateTestAnnotation()
		deployment, err := framework.PatchDeployment("kubernetes-dashboard", "kube-system", testAnnotation)
		if err != nil {
			t.Fatalf("Error patching Kube Dashboard deployment: %v", err)
		}
		deployment, err = framework.GetDeployment(deployment.Name, "kube-system")
		if err != nil {
			t.Fatalf("Error refreshing Kube Dashboard deployment: %v", err)
		}
		assert.Equal(t, deployment.Status.AvailableReplicas, *deployment.Spec.Replicas, "Deployment available replicas did not match expected replicas")
		assert.Equal(t, deployment.Status.UpdatedReplicas, *deployment.Spec.Replicas, "Deployment updated replicas did not match expected replicas")
		if err = framework.DeleteDeployment(deployment.Name, deployment.Namespace); err != nil {
			t.Fatalf("Error deleting kube-dashboard deployment: %v", err)
		}
	})
}

func TestIBMContainerService_IBMSystem(t *testing.T) {
	utils.CheckIfTesting(t, testArmada)
	t.Run("Run anything in ibm-system", func(t *testing.T) {
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/dockerhub-nginx-unsigned.yaml", "ibm-system")
	})
}
