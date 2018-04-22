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
	"fmt"
	"testing"
	"time"

	"github.com/IBM/portieris/test/framework"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
)

func buildDeployment(t *testing.T, fw *framework.Framework, manifestLocation, namespace string, expectCreateFail bool) *appsv1.Deployment {
	manifest, err := fw.LoadDeploymentManifest(manifestLocation)
	if err != nil {
		t.Fatalf("Error loading manifest: %v", err)
	}
	if manifest == nil {
		t.Fatalf("Error loading manifest: manifest is nil")
	}
	if err = fw.CreateDeployment(namespace, manifest); err != nil {
		if !expectCreateFail {

			t.Fatalf("Error creating %q deployment in %v: %v", manifest.Name, namespace, err)
		} else {
			return nil
		}
	}
	fw.WaitForDeploymentPods(manifest.Name, namespace, time.Minute)
	deployment, err := fw.GetDeployment(manifest.Name, namespace)
	if err != nil {
		t.Fatalf("Error getting %q deployment in %v: %v", manifest.Name, namespace, err)
	}
	return deployment
}

func patchDeployment(t *testing.T, fw *framework.Framework, deploymentName, namespace, patchString string, expectCreateFail bool) *appsv1.Deployment {
	if _, err := fw.PatchDeployment(deploymentName, namespace, patchString); err != nil {
		if !expectCreateFail {

			t.Fatalf("Error creating %q deployment in %v: %v", deploymentName, namespace, err)
		} else {
			return nil
		}
	}
	fw.WaitForDeploymentPods(deploymentName, namespace, time.Minute)
	deployment, err := fw.GetDeployment(deploymentName, namespace)
	if err != nil {
		t.Fatalf("Error getting %q deployment in %v: %v", deploymentName, namespace, err)
	}
	return deployment
}
func replaceDeployment(t *testing.T, fw *framework.Framework, namespace, manifestLocation string, expectCreateFail bool) *appsv1.Deployment {
	manifest, err := fw.LoadDeploymentManifest(manifestLocation)
	if err != nil {
		t.Fatalf("Error loading manifest: %v", err)
	}
	if manifest == nil {
		t.Fatalf("Error loading manifest: manifest is nil")
	}
	if _, err := fw.ReplaceDeployment(namespace, manifest); err != nil {
		if !expectCreateFail {

			t.Fatalf("Error creating %q deployment in %v: %v", manifest.Name, namespace, err)
		} else {
			return nil
		}
	}
	fw.WaitForDeploymentPods(manifest.Name, namespace, time.Minute)
	deployment, err := fw.GetDeployment(manifest.Name, namespace)
	if err != nil {
		t.Fatalf("Error getting %q deployment in %v: %v", manifest.Name, namespace, err)
	}
	return deployment
}

// TestDeploymentRunnable tests whether a manifest is deployable to the specified namespace
func TestDeploymentRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	deployment := buildDeployment(t, fw, manifestLocation, namespace, false)
	defer fw.DeleteDeployment(deployment.Name, deployment.Namespace)
	if !assert.Equal(t, *deployment.Spec.Replicas, deployment.Status.AvailableReplicas, "Deployment failed: available replicas did not match expected replicas") {
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
	}
}

// TestDeploymentNotRunnable tests whether a manifest is deployable to the specified namespace
func TestDeploymentNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	deployment := buildDeployment(t, fw, manifestLocation, namespace, true)
	if deployment != nil {
		defer fw.DeleteDeployment(deployment.Name, deployment.Namespace)
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
		t.Errorf("Expected deployment creation to fail")
	}
}

// TestCurrentDeployStatus checks the deployment currently has the expected number of replicas
func TestCurrentDeployStatus(t *testing.T, fw *framework.Framework, namespace, deploymentName string) {
	deployment, err := fw.GetDeployment(deploymentName, namespace)
	if err != nil {
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
		t.Errorf("Failed to check status of deployment")
	}
	if err := fw.WaitForDeploymentPods(deploymentName, namespace, 2*time.Minute); err != nil {
		t.Errorf(err.Error())

	}
	if !assert.Equal(t, *deployment.Spec.Replicas, deployment.Status.AvailableReplicas, "Deployment failed: available replicas did not match expected replicas") {
		fmt.Println(time.Now().String())
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
	}
}

// TestStartDeployNoDelete starts a deployment and only deletes on failure
func TestStartDeployNoDelete(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) (deploymentName string) {
	deployment := buildDeployment(t, fw, manifestLocation, namespace, false)
	if !assert.Equal(t, *deployment.Spec.Replicas, deployment.Status.AvailableReplicas, "Deployment failed: available replicas did not match expected replicas") {
		defer fw.DeleteDeployment(deployment.Name, deployment.Namespace)
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
	}
	return deployment.Name
}

// TestDeploymentNotRunnableOnPatch tests whether a deplomyent is not runnable after a patch
func TestDeploymentNotRunnableOnPatch(t *testing.T, fw *framework.Framework, deploymentName, patchString, namespace string) {
	deployment := patchDeployment(t, fw, deploymentName, namespace, patchString, true)
	if deployment != nil {
		defer fw.DeleteDeployment(deploymentName, deployment.Namespace)
		if !assert.Equal(t, *deployment.Spec.Replicas, deployment.Status.AvailableReplicas, "Deployment failed: available replicas did not match expected replicas") {
			DumpEvents(t, fw, namespace)
			DumpPolicies(t, fw, namespace)
		}
	}
}

// TestDeploymentNotRunnableOnReplace tests whether a deplomyent is not runnable after a replace
func TestDeploymentNotRunnableOnReplace(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	deployment := replaceDeployment(t, fw, namespace, manifestLocation, true)
	if deployment != nil {
		defer fw.DeleteDeployment(deployment.Name, deployment.Namespace)
		if !assert.Equal(t, *deployment.Spec.Replicas, deployment.Status.AvailableReplicas, "Deployment failed: available replicas did not match expected replicas") {
			DumpEvents(t, fw, namespace)
			DumpPolicies(t, fw, namespace)
		}
	}
}
