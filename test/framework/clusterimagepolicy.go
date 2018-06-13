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
	"log"
	"time"

	"github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// LoadClusterImagePolicyManifest takes a manifest and decodes it into a ImagePolicy object
func (f *Framework) LoadClusterImagePolicyManifest(pathToManifest string) (*v1beta1.ClusterImagePolicy, error) {
	manifest, err := openFile(pathToManifest)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file %q: %v", pathToManifest, err)
	}
	ip := v1beta1.ClusterImagePolicy{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&ip); err != nil {
		return nil, fmt.Errorf("Unable to decode file %q: %v", pathToManifest, err)
	}
	return &ip, nil
}

// CreateClusterImagePolicy creates the ClusterImagePolicy
func (f *Framework) CreateClusterImagePolicy(clusterImagePolicy *v1beta1.ClusterImagePolicy) error {
	if _, err := f.ClusterImagePolicyClient.SecurityenforcementV1beta1().ClusterImagePolicies().Create(clusterImagePolicy); err != nil {
		return fmt.Errorf("Error creating ClusterImagePolicy %q: %v", clusterImagePolicy.Name, err)
	}
	return f.WaitForClusterImagePolicy(clusterImagePolicy.Name, time.Minute)
}

// GetClusterImagePolicy retrieves the ClusterImagePolicy
func (f *Framework) GetClusterImagePolicy(name string) (*v1beta1.ClusterImagePolicy, error) {
	return f.ClusterImagePolicyClient.SecurityenforcementV1beta1().ClusterImagePolicies().Get(name, metav1.GetOptions{})
}

// ListClusterImagePolicies creates the ClusterImagePolicy
func (f *Framework) ListClusterImagePolicies() (*v1beta1.ClusterImagePolicyList, error) {
	return f.ClusterImagePolicyClient.SecurityenforcementV1beta1().ClusterImagePolicies().List(metav1.ListOptions{})
}

// WaitForClusterImagePolicy waits until the ClusterImagePolicy is created or the timeout is reached
func (f *Framework) WaitForClusterImagePolicy(name string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.GetClusterImagePolicy(name); err != nil {
			return false, err
		}
		return true, nil
	})
}

// WaitForClusterImagePolicyDefinition waits until the ClusterImagePolicy CRD is created or the timeout is reached
func (f *Framework) WaitForClusterImagePolicyDefinition(timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.GetClusterImagePolicyDefinition(); err != nil {
			return false, err
		}
		log.Print("Found ClusterImagePolicyDefinition")
		return true, nil
	})
}

// GetClusterImagePolicyDefinition retrieves the ClusterImagePolicy CRD
func (f *Framework) GetClusterImagePolicyDefinition() (*apiextensions.CustomResourceDefinition, error) {
	return f.CustomResourceDefinitionClient.Get(clusterImagePolicyCRDName, metav1.GetOptions{})
}

// DeleteClusterImagePolicy deletes the specified ClusterImagePolicy
func (f *Framework) DeleteClusterImagePolicy(name string) error {
	return f.ClusterImagePolicyClient.SecurityenforcementV1beta1().ClusterImagePolicies().Delete(name, &metav1.DeleteOptions{})
}
