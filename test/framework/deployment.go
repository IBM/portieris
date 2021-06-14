// Copyright 2018, 2021 Portieris Authors.
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

package framework

import (
	"fmt"
	"time"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// LoadDeploymentManifest takes a manifest and decodes it into a deployment object.
func (f *Framework) LoadDeploymentManifest(pathToManifest string) (*v1.Deployment, error) {
	manifest, err := openFile(pathToManifest)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file %q: %v", pathToManifest, err)
	}
	deployment := v1.Deployment{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&deployment); err != nil {
		return nil, fmt.Errorf("Unable to decode file %q: %v", pathToManifest, err)
	}
	return &deployment, nil
}

// CreateDeployment creates a deployment and waits for it to show.
func (f *Framework) CreateDeployment(namespace string, deployment *v1.Deployment) error {
	if _, err := f.KubeClient.AppsV1().Deployments(namespace).Create(deployment); err != nil {
		return err
	}
	if err := f.WaitForDeployment(deployment.Name, namespace, time.Minute); err != nil {
		return err
	}
	f.WaitForDeploymentPods(deployment.Name, namespace, time.Second*45)
	return nil
}

// GetDeployment retrieves the specified deployment.
func (f *Framework) GetDeployment(name, namespace string) (*v1.Deployment, error) {
	return f.KubeClient.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
}

// PatchDeployment patches the specified deployment.
func (f *Framework) PatchDeployment(name, namespace, patch string) (*v1.Deployment, error) {
	return f.KubeClient.AppsV1().Deployments(namespace).Patch(name, types.StrategicMergePatchType, []byte(patch))
}

// ReplaceDeployment replaces the specified deployment.
func (f *Framework) ReplaceDeployment(namespace string, deployment *v1.Deployment) (*v1.Deployment, error) {
	return f.KubeClient.AppsV1().Deployments(namespace).Update(deployment)
}

// DeleteDeployment deletes the specified deployment.
func (f *Framework) DeleteDeployment(name, namespace string) error {
	return f.KubeClient.AppsV1().Deployments(namespace).Delete(name, &metav1.DeleteOptions{})
}

// WaitForDeployment waits until the specified deployment is created or the timeout is reached.
func (f *Framework) WaitForDeployment(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.GetDeployment(name, namespace); err != nil {
			return false, nil
		}
		return true, nil
	})
}

// WaitForDeploymentPods waits until the specified deployment's pods are created or the timeout is reached.
func (f *Framework) WaitForDeploymentPods(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		deployment, err := f.GetDeployment(name, namespace)
		if err != nil {
			return false, err
		}
		if *deployment.Spec.Replicas != deployment.Status.AvailableReplicas {
			return false, nil
		}
		return true, nil
	})
}

// ListDeployments lists all deployments that are associated with the installed Helm release.
func (f *Framework) ListDeployments() (*v1.DeploymentList, error) {
	opts := f.getHelmReleaseSelectorListOptions()
	return f.KubeClient.AppsV1().Deployments(corev1.NamespaceAll).List(opts)
}
