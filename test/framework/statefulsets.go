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

	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// LoadStatefulSetManifest takes a manifest and decodes it into a StatefulSet object.
func (f *Framework) LoadStatefulSetManifest(pathToManifest string) (*v1.StatefulSet, error) {
	manifest, err := openFile(pathToManifest)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file %q: %v", pathToManifest, err)
	}
	statefulset := v1.StatefulSet{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&statefulset); err != nil {
		return nil, fmt.Errorf("Unable to decode file %q: %v", pathToManifest, err)
	}
	return &statefulset, nil
}

// CreateStatefulSet creates a StatefulSet and waits for it to show.
func (f *Framework) CreateStatefulSet(namespace string, statefulset *v1.StatefulSet) error {
	if _, err := f.KubeClient.AppsV1().StatefulSets(namespace).Create(statefulset); err != nil {
		return err
	}
	if err := f.WaitForStatefulSet(statefulset.Name, namespace, time.Minute); err != nil {
		return err
	}
	f.WaitForStatefulSetPods(statefulset.Name, namespace, time.Second*45)
	return nil
}

// GetStatefulSet retrieves the specified StatefulSet.
func (f *Framework) GetStatefulSet(name, namespace string) (*v1.StatefulSet, error) {
	return f.KubeClient.AppsV1().StatefulSets(namespace).Get(name, metav1.GetOptions{})
}

// DeleteStatefulSet deletes the specified StatefulSet.
func (f *Framework) DeleteStatefulSet(name, namespace string) error {
	return f.KubeClient.AppsV1().StatefulSets(namespace).Delete(name, &metav1.DeleteOptions{})
}

// WaitForStatefulSet waits until the specified StatefulSet is created or the timeout is reached.
func (f *Framework) WaitForStatefulSet(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.GetStatefulSet(name, namespace); err != nil {
			return false, nil
		}
		return true, nil
	})
}

// WaitForStatefulSetPods waits until the specified StatefulSet's pods are created or the timeout is reached.
func (f *Framework) WaitForStatefulSetPods(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		statefulset, err := f.GetStatefulSet(name, namespace)
		if err != nil {
			return false, err
		}
		if *statefulset.Spec.Replicas != statefulset.Status.ReadyReplicas {
			return false, nil
		}
		return true, nil
	})
}

// ListStatefulSet lists the StatefulSets that are associated with the installed Helm release.
func (f *Framework) ListStatefulSet() (*v1.StatefulSetList, error) {
	opts := f.getHelmReleaseSelectorListOptions()
	return f.KubeClient.AppsV1().StatefulSets(corev1.NamespaceAll).List(opts)
}
