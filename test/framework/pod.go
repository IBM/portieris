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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetPod retrieves the specified pod.
func (f *Framework) GetPod(name, namespace string) (*corev1.Pod, error) {
	return f.KubeClient.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
}

// LoadPodManifest takes a manifest and decodes it into a pod.
func (f *Framework) LoadPodManifest(pathToManifest string) (*corev1.Pod, error) {
	manifest, err := openFile(pathToManifest)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file %q: %v", pathToManifest, err)
	}
	pod := corev1.Pod{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&pod); err != nil {
		return nil, fmt.Errorf("Unable to decode file %q: %v", pathToManifest, err)
	}
	return &pod, nil
}

// CreatePod creates a pod and waits for it to show.
func (f *Framework) CreatePod(namespace string, pod *corev1.Pod) error {
	if _, err := f.KubeClient.CoreV1().Pods(namespace).Create(pod); err != nil {
		return err
	}
	return f.WaitForPod(pod.Name, namespace, 2*time.Minute)
}

// DeletePod deletes the specified pod.
func (f *Framework) DeletePod(name, namespace string) error {
	return f.KubeClient.CoreV1().Pods(namespace).Delete(name, &metav1.DeleteOptions{})
}

// WaitForPod waits until the pod deployment completes.
func (f *Framework) WaitForPod(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		pod, err := f.GetPod(name, namespace)
		if err != nil {
			return false, err
		}

		if pod.Status.Phase != corev1.PodRunning {
			return false, nil
		}
		return true, nil
	})
}

// WaitForPodDelete waits until the pod is deleted.
func (f *Framework) WaitForPodDelete(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		_, err := f.GetPod(name, namespace)
		if err == nil {
			return false, nil
		}
		return true, nil
	})
}

// DeleteRandomPod deletes the first pod that is returned in the list of pods for a specified namespace.
func (f *Framework) DeleteRandomPod(namespace string) error {
	podList, err := f.KubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	err = f.DeletePod(podList.Items[0].Name, namespace)
	if err != nil {
		return err
	}
	err = f.WaitForPodDelete(podList.Items[0].Name, namespace, time.Minute)

	return err

}
