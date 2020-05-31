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

package v1beta1

import (
	"strings"

	"github.com/IBM/portieris/helpers/wildcard"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// TruePointer - pointer to a boolwan value of true
	TruePointer = boolPointer(true)
	// FalsePointer - pointer to a boolwan value of false
	FalsePointer = boolPointer(false)
)

func boolPointer(boolean bool) *bool {
	return &boolean
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImagePolicy is a specification for a ImagePolicy resource
type ImagePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PolicySpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImagePolicyList is a list of ImagePolicy resources
type ImagePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ImagePolicy `json:"items"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterImagePolicy is a specification for a ClusterImagePolicy resource
type ClusterImagePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PolicySpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterImagePolicyList is a list of ClusterImagePolicy resources
type ClusterImagePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ClusterImagePolicy `json:"items"`
}

// PolicySpec is the spec for a ImagePolicy or ClusterImagePolicy resource
type PolicySpec struct {
	Repositories []Repository `json:"repositories"`
}

// Repository .
type Repository struct {
	Name   string `json:"name,omitempty"` // Name may contain a * to signify one or more characters
	Policy Policy `json:"policy,omitempty"`
}

// Policy .
type Policy struct {
	Trust       Trust       `json:"trust,omitempty"`
	SimpleStore SimpleStore `json:"simpleStore,omitempty"`
	Simple      []Simple    `json:"simple,omitempty"`
	Va          VA          `json:"va,omitempty"`
}

// Trust .
type Trust struct {
	Enabled       *bool    `json:"enabled,omitempty"`
	SignerSecrets []Signer `json:"signerSecrets,omitempty"`
	TrustServer   string   `json:"trustServer,omitempty"`
}

// Signer .
type Signer struct {
	Name string `json:"name"`
}

// SimpleStore .
type SimpleStore struct {
	URL    string `json:"url"`
	Secret string `json:"secret,omitEmpty"`
}

// Simple .
type Simple struct {
	Type           string              `json:"type"`
	KeySecret      string              `json:"keySecret,omitEmpty"`
	SignedIdentity IdentityRequirement `json:"signedIdentity,omitEmpty"`
}

// IdentityRequirement .
type IdentityRequirement struct {
	Type             string `json:"type"`
	DockerReference  string `json:"dockerReference,omitEmpty"`
	DockerRepository string `json:"dockerRepository,omitEmpty"`
}

// VA .
type VA struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// FindImagePolicy - Given an ImagePolicyList, find the repository whose name
// most closely matches the image name, and returns its policy.
// If there are no matches, return a nil value.
func (apl ImagePolicyList) FindImagePolicy(image string) *Policy {
	// Variables
	bestMatchQuality := -1
	bestMatchedPolicy := Policy{}

	// Check if there are policies for the given image
Exact:
	for _, item := range apl.Items {

		// iterate over the repositories
		for _, repo := range item.Spec.Repositories {

			// get the name for the current repository
			repositoryName := repo.Name
			hasWildcard := strings.Contains(repositoryName, "*")
			// glog.Infof("repositoryName: %s", repositoryName)

			// Check if the image name matches the repository name
			match := false
			matchQuality := -1
			if !hasWildcard && repositoryName == image {
				// glog.Info("Found exact match")
				bestMatchQuality = len(image)
				bestMatchedPolicy = repo.Policy
				break Exact
			} else {
				if wildcard.CompareAnyTag(repositoryName, image) {
					match = true
					matchQuality = len(repositoryName) - strings.Count(repositoryName, "*")
				}
			}
			// glog.Infof("match: %t  matchQuality: %d", match, matchQuality)
			if match == true && matchQuality > bestMatchQuality {
				// glog.Info("Updating to this match")
				bestMatchQuality = matchQuality
				bestMatchedPolicy = repo.Policy
			}
		}
	}
	if bestMatchQuality > -1 {
		return &bestMatchedPolicy
	}
	return nil
}

// FindClusterImagePolicy - Given an ClusterImagePolicyList, find the repository whose name
// most closely matches the image name, and returns its policy.
// If there are no matches, return a nil value.
func (apl ClusterImagePolicyList) FindClusterImagePolicy(image string) *Policy {
	// Variables
	bestMatchQuality := -1
	bestMatchedPolicy := Policy{}

	// Check if there are policies for the given image
Exact:
	for _, item := range apl.Items {

		// iterate over the repositories
		for _, repo := range item.Spec.Repositories {

			// get the name for the current repository
			repositoryName := repo.Name
			hasWildcard := strings.Contains(repositoryName, "*")
			// glog.Infof("repositoryName: %s", repositoryName)

			// Check if the image name matches the repository name
			match := false
			matchQuality := -1
			if !hasWildcard && repositoryName == image {
				// glog.Info("Found exact match")
				bestMatchQuality = len(image)
				bestMatchedPolicy = repo.Policy
				break Exact
			} else {
				if wildcard.CompareAnyTag(repositoryName, image) {
					match = true
					matchQuality = len(repositoryName) - strings.Count(repositoryName, "*")
				}
			}
			// glog.Infof("match: %t  matchQuality: %d", match, matchQuality)
			if match == true && matchQuality > bestMatchQuality {
				// glog.Info("Updating to this match")
				bestMatchQuality = matchQuality
				bestMatchedPolicy = repo.Policy
			}
		}
	}
	if bestMatchQuality > -1 {
		return &bestMatchedPolicy
	}
	return nil
}
