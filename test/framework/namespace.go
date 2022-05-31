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
	"context"
	"encoding/json"
	"fmt"
	"time"

	pk "github.com/IBM/portieris/pkg/kubernetes"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func generateNamespace(name string) *corev1.Namespace {
	ns := &corev1.Namespace{}
	ns.Name = name
	ns.Kind = "Namespace"
	return ns
}

// IBMCloudSecretNames is the secret name that is provided to enable access to the test images.
// https://github.com/IBM/portieris/issues/34 to remove the need for this
var IBMCloudSecretNames = []string{"all-icr-io", "default-icr-io"}

// IBMTestRegistry is the default location of the test images that are used in the end-to-end tests.
var IBMTestRegistry = "de.icr.io"

// CreateNamespace creates a namespace.
func (f *Framework) CreateNamespace(name string) (*corev1.Namespace, error) {

	if _, err := f.KubeClient.CoreV1().Namespaces().Create(context.TODO(), generateNamespace(name), metav1.CreateOptions{}); err != nil {
		return nil, err
	}
	namespace, err := f.GetNamespace(name)
	if err != nil {
		return nil, err
	}
	f.WaitForNamespace(namespace.Name, time.Second*10)
	return namespace, nil
}

// GetNamespace retrieves the specified namespace.
func (f *Framework) GetNamespace(name string) (*corev1.Namespace, error) {
	return f.KubeClient.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
}

// WaitForNamespace waits until the specified namespace is created or the timeout is reached.
func (f *Framework) WaitForNamespace(name string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.GetNamespace(name); err != nil {
			return false, err
		}
		return true, nil
	})
}

// CreateNamespaceWithIPS creates a namespace, service account, and IP addresses to pull from the Global region of IBM Cloud Container Registry.
// It copies the IBM Cloud secret name, imagePullSecret, from the default namespace.
func (f *Framework) CreateNamespaceWithIPS(name string) (*corev1.Namespace, error) {
	namespace, err := f.CreateNamespace(name)
	if err != nil {
		return nil, fmt.Errorf("error creating namespace: %v", err)
	}
	var imagePullSecret *v1.Secret
	for _, secretName := range IBMCloudSecretNames {
		imagePullSecret, err = f.KubeClient.CoreV1().Secrets("default").Get(context.TODO(), secretName, metav1.GetOptions{})
	}
	if imagePullSecret == nil {
		return nil, fmt.Errorf("error getting imagePullSecret: %v", err)
	}
	imagePullSecret.Namespace = namespace.Name
	imagePullSecret.ResourceVersion = ""
	if _, err := f.KubeClient.CoreV1().Secrets(namespace.Name).Create(context.TODO(), imagePullSecret, metav1.CreateOptions{}); err != nil {
		return nil, fmt.Errorf("error creating imagePullSecret: %v", err)
	}

	// Create an invalid pull secret that is based on the valid pull secret (hard-code the password to an invalid value).
	badPullSecret := imagePullSecret.DeepCopy()
	badPullSecret.Name = "bad-" + badPullSecret.GetName()
	clientWrapper := pk.NewKubeClientsetWrapper(f.KubeClient)
	goodUser, _, err := clientWrapper.GetSecretToken(namespace.Name, imagePullSecret.GetName(), IBMTestRegistry)
	if err == nil {
		badAuths := pk.Auths{
			Registries: pk.RegistriesStruct{
				IBMTestRegistry: pk.RegistryCredentials{
					Username: goodUser,
					Password: "iamnotanapikey",
					Email:    "a@b.c",
				},
			},
		}
		badAuthData, _ := json.Marshal(badAuths)
		badPullSecret.Data[".dockerconfigjson"] = badAuthData
		if _, err := f.KubeClient.CoreV1().Secrets(namespace.Name).Create(context.TODO(), badPullSecret, metav1.CreateOptions{}); err != nil {
			return nil, fmt.Errorf("error creating bad imagePullSecret: %v", err)
		}
	}

	sa := generateServiceAccount("default")
	// Ensure the invalid imagePullSecret is before the valid one in the ServiceAccount's list.
	sa.ImagePullSecrets = []corev1.LocalObjectReference{
		{Name: badPullSecret.GetName()},
		{Name: imagePullSecret.GetName()},
	}
	if _, err := f.KubeClient.CoreV1().ServiceAccounts(namespace.Name).Update(context.TODO(), sa, metav1.UpdateOptions{}); err != nil {
		return nil, fmt.Errorf("error adding imagePullSecret to ServiceAccount: %v", err)
	}
	return namespace, nil
}

// DeleteNamespace deletes the specified namespace.
func (f *Framework) DeleteNamespace(name string) error {
	return f.KubeClient.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
}
