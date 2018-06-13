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

package policy

import (
	"fmt"

	securityenforcementclientset "github.com/IBM/portieris/pkg/apis/securityenforcement/client/clientset/versioned"
	securityenforcementv1beta1 "github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Interface defines the interface needed to work out which policy should be enforced
type Interface interface {
	GetPolicyToEnforce(namespace, image string) (*securityenforcementv1beta1.Policy, error)
}

// Client is responsible for working out which policy should be enforced
type Client struct {
	// policyClientSet is a securityenforcementclientset for the policy CRDs
	policyClientSet securityenforcementclientset.Interface
}

// NewClient creates a new policy client using the Security Enforcement client set it is passed
func NewClient(policyClientSet securityenforcementclientset.Interface) *Client {
	return &Client{
		policyClientSet: policyClientSet,
	}
}

// getImagePolicyList retrieves the list of image policies in the specified namespace
func (c *Client) getImagePolicyList(namespace string) (*securityenforcementv1beta1.ImagePolicyList, error) {
	policies, err := c.policyClientSet.SecurityenforcementV1beta1().ImagePolicies(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return policies, nil
}

// getClusterPolicySpec retrieves the lost of clusterwide image policies
func (c *Client) getClusterImagePolicyList() (*securityenforcementv1beta1.ClusterImagePolicyList, error) {
	policies, err := c.policyClientSet.SecurityenforcementV1beta1().ClusterImagePolicies().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return policies, nil
}

// GetPolicyToEnforce retrieves the policy that should be enforced for the specified image in the given namespace
func (c *Client) GetPolicyToEnforce(namespace, image string) (*securityenforcementv1beta1.Policy, error) {
	policyList, err := c.getImagePolicyList(namespace)
	if err != nil {
		return nil, err
	}

	if len((*policyList).Items) == 0 {
		// We don't have any image policies in the current namespace, get the list of cluster policies
		clusterPolicyList, err := c.getClusterImagePolicyList()
		if err != nil {
			return nil, err
		}

		if len((*clusterPolicyList).Items) == 0 {
			// We also don't have any cluster image policies, deny the request
			return nil, fmt.Errorf("Deny %q, no image policies or cluster polices", image)
		}

		// See if there is a match for the image
		clusterPolicy := clusterPolicyList.FindClusterImagePolicy(image)
		if clusterPolicy == nil {
			// We also don't have any cluster image policies, deny the request
			return nil, fmt.Errorf("Deny %q, no matching repositories in ClusterImagePolicy and no ImagePolicies in the %q namespace", image, namespace)
		}
		return clusterPolicy, nil
	}

	// For this image, see if there is an ImagePolicy repository that matches.
	// Get the policy if it does
	policy := policyList.FindImagePolicy(image)

	if policy == nil {
		// We also don't have any cluster image policies, deny the request
		return nil, fmt.Errorf("Deny %q, no matching repositories in the ImagePolicies", image)
	}
	return policy, nil
}
