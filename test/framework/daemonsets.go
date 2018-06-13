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

package framework

import (
	"fmt"
	"time"

	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// LoadDaemonSetManifest takes a manifest and decodes it into a daemonset object
func (f *Framework) LoadDaemonSetManifest(pathToManifest string) (*v1.DaemonSet, error) {
	manifest, err := openFile(pathToManifest)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file %q: %v", pathToManifest, err)
	}
	daemonset := v1.DaemonSet{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&daemonset); err != nil {
		return nil, fmt.Errorf("Unable to decode file %q: %v", pathToManifest, err)
	}
	return &daemonset, nil
}

// CreateDaemonSet creates a daemonset resource and then waits for it to appear
func (f *Framework) CreateDaemonSet(namespace string, daemonset *v1.DaemonSet) error {
	if _, err := f.KubeClient.AppsV1().DaemonSets(namespace).Create(daemonset); err != nil {
		return err
	}
	if err := f.WaitForDaemonSet(daemonset.Name, namespace, time.Minute); err != nil {
		return err
	}
	f.WaitForDaemonSetPods(daemonset.Name, namespace, time.Second*45)
	return nil
}

// GetDaemonSets retrieves the specified deployment
func (f *Framework) GetDaemonSets(name, namespace string) (*v1.DaemonSet, error) {
	return f.KubeClient.AppsV1().DaemonSets(namespace).Get(name, metav1.GetOptions{})
}

// // PatchDeployment patches the specified deployment
// func (f *Framework) PatchDeployment(name, namespace, patch string) (*v1.Deployment, error) {
// 	return f.KubeClient.AppsV1().Deployments(namespace).Patch(name, types.StrategicMergePatchType, []byte(patch))
// }

// DeleteDaemonSet deletes the specified deployment
func (f *Framework) DeleteDaemonSet(name, namespace string) error {
	return f.KubeClient.AppsV1().DaemonSets(namespace).Delete(name, &metav1.DeleteOptions{})
}

// WaitForDaemonSet waits until the specified daemonset is created or the timeout is reached
func (f *Framework) WaitForDaemonSet(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.GetDaemonSets(name, namespace); err != nil {
			return false, nil
		}
		return true, nil
	})
}

// WaitForDaemonSetPods waits until the specified deployment's pods are created or the timeout is reached
func (f *Framework) WaitForDaemonSetPods(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		daemonset, err := f.GetDaemonSets(name, namespace)
		if err != nil {
			return false, err
		}
		if daemonset.Status.DesiredNumberScheduled != daemonset.Status.NumberReady {
			return false, nil
		}
		return true, nil
	})
}

// ListDaemonSet retrieves all daemonset associated with the installed Helm release
func (f *Framework) ListDaemonSet() (*v1.DaemonSetList, error) {
	opts := f.getHelmReleaseSelectorListOptions()
	return f.KubeClient.AppsV1().DaemonSets(corev1.NamespaceAll).List(opts)
}
