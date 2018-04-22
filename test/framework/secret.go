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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// LoadSecretManifest takes a manifest and decodes it into a deployment object
func (f *Framework) LoadSecretManifest(pathToManifest string) (*corev1.Secret, error) {
	manifest, err := openFile(pathToManifest)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file %q: %v", pathToManifest, err)
	}
	secret := corev1.Secret{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&secret); err != nil {
		return nil, fmt.Errorf("Unable to decode file %q: %v", pathToManifest, err)
	}
	return &secret, nil
}

// CreateSecret creates a secret resource and then waits for it to appear
func (f *Framework) CreateSecret(namespace string, secret *corev1.Secret) error {
	if _, err := f.KubeClient.CoreV1().Secrets(namespace).Create(secret); err != nil {
		return err
	}
	return f.WaitForSecret(secret.Name, namespace, time.Minute)
}

// GetSecret retrieves the specified secret
func (f *Framework) GetSecret(name, namespace string) (*corev1.Secret, error) {
	return f.KubeClient.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
}

// WaitForSecret waits until the specified deployment is created or the timeout is reached
func (f *Framework) WaitForSecret(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.GetSecret(name, namespace); err != nil {
			return false, err
		}
		return true, nil
	})
}
