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

// IBMCloudSecretNames secret provided to enable access to test images
// https://github.com/IBM/portieris/issues/34 to remove the need for this
var IBMCloudSecretNames = []string{"all-icr-io", "default-icr-io"}

// IBMGlobalRegistry is the default location of the test images used in these e2e tests
var IBMGlobalRegistry = "icr.io"

// CreateNamespace creates a namespace
func (f *Framework) CreateNamespace(name string) (*corev1.Namespace, error) {

	if _, err := f.KubeClient.CoreV1().Namespaces().Create(generateNamespace(name)); err != nil {
		return nil, err
	}
	namespace, err := f.GetNamespace(name)
	if err != nil {
		return nil, err
	}
	f.WaitForNamespace(namespace.Name, time.Second*10)
	return namespace, nil
}

// GetNamespace retrieves the specified namespace
func (f *Framework) GetNamespace(name string) (*corev1.Namespace, error) {
	return f.KubeClient.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
}

// WaitForNamespace waits until the specified namespace is created or the timeout is reached
func (f *Framework) WaitForNamespace(name string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.GetNamespace(name); err != nil {
			return false, err
		}
		return true, nil
	})
}

// CreateNamespaceWithIPS creates a namespace, service account and IPS to pull from the IBM Cloud Container Registry Global region
// It copies the `IBMCloudSecretName` imagePullSecret from the default namespace
func (f *Framework) CreateNamespaceWithIPS(name string) (*corev1.Namespace, error) {
	namespace, err := f.CreateNamespace(name)
	if err != nil {
		return nil, fmt.Errorf("error creating namespace: %v", err)
	}
	var imagePullSecret *v1.Secret
	for _, secretName := range IBMCloudSecretNames {
		imagePullSecret, err = f.KubeClient.CoreV1().Secrets("default").Get(secretName, metav1.GetOptions{})
		if err == nil {
			break
		}
	}
	if imagePullSecret == nil {
		return nil, fmt.Errorf("error getting imagePullSecret: %v", err)
	}
	imagePullSecret.Namespace = namespace.Name
	imagePullSecret.ResourceVersion = ""
	if _, err := f.KubeClient.CoreV1().Secrets(namespace.Name).Create(imagePullSecret); err != nil {
		return nil, fmt.Errorf("error creating imagePullSecret: %v", err)
	}

	// Create a bad pull secret based off the good one (hard-code password to a bad value)
	badPullSecret := imagePullSecret.DeepCopy()
	badPullSecret.Name = "bad-" + badPullSecret.GetName()
	clientWrapper := pk.NewKubeClientsetWrapper(f.KubeClient)
	goodUser, _, err := clientWrapper.GetSecretToken(namespace.Name, imagePullSecret.GetName(), IBMGlobalRegistry)
	if err == nil {
		badAuths := pk.Auths{
			Registries: pk.RegistriesStruct{
				IBMGlobalRegistry: pk.RegistryCredentials{
					Username: goodUser,
					Password: "iamnotanapikey",
					Email:    "a@b.c",
				},
			},
		}
		badAuthData, _ := json.Marshal(badAuths)
		badPullSecret.Data[".dockerconfigjson"] = badAuthData
		if _, err := f.KubeClient.CoreV1().Secrets(namespace.Name).Create(badPullSecret); err != nil {
			return nil, fmt.Errorf("error creating bad imagePullSecret: %v", err)
		}
	}

	sa := generateServiceAccount("default")
	// Ensure the bad imagePullSecret is before the good one in the ServiceAccount's list
	sa.ImagePullSecrets = []corev1.LocalObjectReference{
		{Name: badPullSecret.GetName()},
		{Name: imagePullSecret.GetName()},
	}
	if _, err := f.KubeClient.CoreV1().ServiceAccounts(namespace.Name).Update(sa); err != nil {
		return nil, fmt.Errorf("error adding imagePullSecret to ServiceAccount: %v", err)
	}
	return namespace, nil
}

// DeleteNamespace deletes the specified namespace
func (f *Framework) DeleteNamespace(name string) error {
	return f.KubeClient.CoreV1().Namespaces().Delete(name, &metav1.DeleteOptions{})
}
