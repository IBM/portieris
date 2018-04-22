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

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
)

// LoadImagePolicyManifest takes a manifest and decodes it into a ImagePolicy object
func (f *Framework) LoadImagePolicyManifest(pathToManifest string) (*v1beta1.ImagePolicy, error) {
	manifest, err := openFile(pathToManifest)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file %q: %v", pathToManifest, err)
	}
	ip := v1beta1.ImagePolicy{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&ip); err != nil {
		return nil, fmt.Errorf("Unable to decode file %q: %v", pathToManifest, err)
	}
	return &ip, nil
}

// CreateImagePolicy creates the ImagePolicy
func (f *Framework) CreateImagePolicy(namespace string, imagePolicy *v1beta1.ImagePolicy) error {
	if _, err := f.ImagePolicyClient.SecurityenforcementV1beta1().ImagePolicies(namespace).Create(imagePolicy); err != nil {
		return fmt.Errorf("Error creating ImagePolicy %q: %v", imagePolicy.Name, err)
	}
	return f.WaitForImagePolicy(imagePolicy.Name, namespace, time.Minute)
}

// GetImagePolicy retrieves the ImagePolicy
func (f *Framework) GetImagePolicy(name, namespace string) (*v1beta1.ImagePolicy, error) {
	return f.ImagePolicyClient.SecurityenforcementV1beta1().ImagePolicies(namespace).Get(name, metav1.GetOptions{})
}

// ListImagePolicies lists all ImagePolicies in a given namespace
func (f *Framework) ListImagePolicies(namespace string) (*v1beta1.ImagePolicyList, error) {
	return f.ImagePolicyClient.SecurityenforcementV1beta1().ImagePolicies(namespace).List(metav1.ListOptions{})
}

// WaitForImagePolicy waits until the ImagePolicy is created or the timeout is reached
func (f *Framework) WaitForImagePolicy(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.GetImagePolicy(name, namespace); err != nil {
			return false, err
		}
		return true, nil
	})
}

// DeleteImagePolicy deletes the ImagePolicy
func (f *Framework) DeleteImagePolicy(name, namespace string) error {
	err := f.ImagePolicyClient.SecurityenforcementV1beta1().ImagePolicies(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// WaitForImagePolicyDefinition waits until the ImagePolicy CRD is created or the timeout is reached
func (f *Framework) WaitForImagePolicyDefinition(timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.GetImagePolicyDefinition(); err != nil {
			return false, err
		}
		return true, nil
	})
}

// GetImagePolicyDefinition retrieves the ImagePolicy CRD
func (f *Framework) GetImagePolicyDefinition() (*apiextensions.CustomResourceDefinition, error) {
	return f.CustomResourceDefinitionClient.Get(imagePolicyCRDName, metav1.GetOptions{})
}

// UpdateImagePolicy creates the ImagePolicy
func (f *Framework) UpdateImagePolicy(namespace string, imagePolicy *v1beta1.ImagePolicy) error {
	if _, err := f.ImagePolicyClient.SecurityenforcementV1beta1().ImagePolicies(namespace).Update(imagePolicy); err != nil {
		return fmt.Errorf("Error updating ImagePolicy %q: %v", imagePolicy.Name, err)
	}
	return f.WaitForImagePolicy(imagePolicy.Name, namespace, time.Minute)
}
