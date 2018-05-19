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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// LoadReplicationControllerManifest takes a manifest and decodes it into a Replicaset object
func (f *Framework) LoadReplicationControllerManifest(pathToManifest string) (*corev1.ReplicationController, error) {
	manifest, err := openFile(pathToManifest)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file %q: %v", pathToManifest, err)
	}
	replicationcontroller := corev1.ReplicationController{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&replicationcontroller); err != nil {
		return nil, fmt.Errorf("Unable to decode file %q: %v", pathToManifest, err)
	}
	return &replicationcontroller, nil
}

// CreateReplicationController creates a Replicaset resource and then waits for it to appear
func (f *Framework) CreateReplicationController(namespace string, replicationcontroller *corev1.ReplicationController) error {
	if _, err := f.KubeClient.CoreV1().ReplicationControllers(namespace).Create(replicationcontroller); err != nil {
		return err
	}
	if err := f.WaitForReplicationController(replicationcontroller.Name, namespace, time.Minute); err != nil {
		return err
	}
	f.WaitForReplicationControllerPods(replicationcontroller.Name, namespace, time.Second*45)
	return nil
}

// GetReplicationController retrieves the specified deployment
func (f *Framework) GetReplicationController(name, namespace string) (*corev1.ReplicationController, error) {
	return f.KubeClient.CoreV1().ReplicationControllers(namespace).Get(name, metav1.GetOptions{})
}

// // PatchDeployment patches the specified deployment
// func (f *Framework) PatchDeployment(name, namespace, patch string) (*v1.Deployment, error) {
// 	return f.KubeClient.CoreV1().Deployments(namespace).Patch(name, types.StrategicMergePatchType, []byte(patch))
// }

// DeleteReplicationController deletes the specified deployment
func (f *Framework) DeleteReplicationController(name, namespace string) error {
	return f.KubeClient.CoreV1().ReplicationControllers(namespace).Delete(name, &metav1.DeleteOptions{})
}

// WaitForReplicationController waits until the specified Replicaset is created or the timeout is reached
func (f *Framework) WaitForReplicationController(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.GetReplicationController(name, namespace); err != nil {
			return false, nil
		}
		return true, nil
	})
}

// WaitForReplicationControllerPods waits until the specified deployment's pods are created or the timeout is reached
func (f *Framework) WaitForReplicationControllerPods(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		replicationcontroller, err := f.GetReplicationController(name, namespace)
		if err != nil {
			return false, err
		}
		if *replicationcontroller.Spec.Replicas != replicationcontroller.Status.AvailableReplicas {
			return false, nil
		}
		return true, nil
	})
}

// ListReplicationController retrieves all Replicaset associated with the installed Helm release
func (f *Framework) ListReplicationController() (*corev1.ReplicationControllerList, error) {
	opts := f.getHelmReleaseSelectorListOptions()
	return f.KubeClient.CoreV1().ReplicationControllers(corev1.NamespaceAll).List(opts)
}
