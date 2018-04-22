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

// LoadReplicaSetManifest takes a manifest and decodes it into a Replicaset object
func (f *Framework) LoadReplicaSetManifest(pathToManifest string) (*v1.ReplicaSet, error) {
	manifest, err := openFile(pathToManifest)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file %q: %v", pathToManifest, err)
	}
	replicaset := v1.ReplicaSet{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&replicaset); err != nil {
		return nil, fmt.Errorf("Unable to decode file %q: %v", pathToManifest, err)
	}
	return &replicaset, nil
}

// CreateReplicaSet creates a Replicaset resource and then waits for it to appear
func (f *Framework) CreateReplicaSet(namespace string, replicaset *v1.ReplicaSet) error {
	if _, err := f.KubeClient.AppsV1().ReplicaSets(namespace).Create(replicaset); err != nil {
		return err
	}
	if err := f.WaitForReplicaSet(replicaset.Name, namespace, time.Minute); err != nil {
		return err
	}
	f.WaitForReplicaSetPods(replicaset.Name, namespace, time.Second*45)
	return nil
}

// GetReplicaSet retrieves the specified deployment
func (f *Framework) GetReplicaSet(name, namespace string) (*v1.ReplicaSet, error) {
	return f.KubeClient.AppsV1().ReplicaSets(namespace).Get(name, metav1.GetOptions{})
}

// // PatchDeployment patches the specified deployment
// func (f *Framework) PatchDeployment(name, namespace, patch string) (*v1.Deployment, error) {
// 	return f.KubeClient.AppsV1().Deployments(namespace).Patch(name, types.StrategicMergePatchType, []byte(patch))
// }

// DeleteReplicaSet deletes the specified deployment
func (f *Framework) DeleteReplicaSet(name, namespace string) error {
	return f.KubeClient.AppsV1().ReplicaSets(namespace).Delete(name, &metav1.DeleteOptions{})
}

// WaitForReplicaSet waits until the specified Replicaset is created or the timeout is reached
func (f *Framework) WaitForReplicaSet(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.GetReplicaSet(name, namespace); err != nil {
			return false, nil
		}
		return true, nil
	})
}

// WaitForReplicaSetPods waits until the specified deployment's pods are created or the timeout is reached
func (f *Framework) WaitForReplicaSetPods(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		replicaset, err := f.GetReplicaSet(name, namespace)
		if err != nil {
			return false, err
		}
		if *replicaset.Spec.Replicas != replicaset.Status.AvailableReplicas {
			return false, nil
		}
		return true, nil
	})
}

// ListReplicaSet retrieves all Replicaset associated with the installed Helm release
func (f *Framework) ListReplicaSet() (*v1.ReplicaSetList, error) {
	opts := f.getHelmReleaseSelectorListOptions()
	return f.KubeClient.AppsV1().ReplicaSets(corev1.NamespaceAll).List(opts)
}
