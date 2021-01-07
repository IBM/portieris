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

package policy

import (
	"errors"
	"testing"

	policyclientset "github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/client/clientset/versioned"
	"github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/client/clientset/versioned/fake"
	policyv1 "github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	trueBool       = true
	falseBool      = false
	imagePolicyOne = &policyv1.ImagePolicy{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "policy-one"},
		TypeMeta:   metav1.TypeMeta{Kind: "ImagePolicy"},
		Spec: policyv1.ImagePolicySpec{
			Repositories: []policyv1.Repository{{Name: "repo-one"}},
		},
	}
	imagePolicyTwo = &policyv1.ImagePolicy{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "policy-two"},
		TypeMeta:   metav1.TypeMeta{Kind: "ImagePolicy"},
		Spec: policyv1.ImagePolicySpec{
			Repositories: []policyv1.Repository{{Name: "repo-two"}},
		},
	}

	clusterImagePolicyOne = &policyv1.ClusterImagePolicy{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "policy-one"},
		TypeMeta:   metav1.TypeMeta{Kind: "ClusterImagePolicy"},
		Spec: policyv1.ImagePolicySpec{
			Repositories: []policyv1.Repository{{Name: "repo-one"}},
		},
	}
	clusterImagePolicyTwo = &policyv1.ClusterImagePolicy{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "policy-two"},
		TypeMeta:   metav1.TypeMeta{Kind: "ClusterImagePolicy"},
		Spec: policyv1.ImagePolicySpec{
			Repositories: []policyv1.Repository{{Name: "repo-two"}},
		},
	}

	helloWorldRepositoryTrustEnabled  = policyv1.Repository{Name: "icr.io/hello/world", Policy: enabledTrustPolicy}
	helloWorldRepositoryTrustDisabled = policyv1.Repository{Name: "icr.io/hello/world", Policy: disabledTrustPolicy}
	helloEarthRepositoryTrustEnabled  = policyv1.Repository{Name: "icr.io/hello/earth", Policy: enabledTrustPolicy}
	helloEarthRepositoryTrustDisabled = policyv1.Repository{Name: "icr.io/hello/earth", Policy: disabledTrustPolicy}

	enabledTrustPolicy  = policyv1.Policy{Trust: policyv1.Trust{Enabled: &trueBool}}
	disabledTrustPolicy = policyv1.Policy{Trust: policyv1.Trust{Enabled: &falseBool}}
)

func createClusterImagePolicy(name string, repos []policyv1.Repository) *policyv1.ClusterImagePolicy {
	return &policyv1.ClusterImagePolicy{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		TypeMeta:   metav1.TypeMeta{Kind: "ClusterImagePolicy"},
		Spec: policyv1.ImagePolicySpec{
			Repositories: repos,
		},
	}
}

func createImagePolicy(name, namespace string, repos []policyv1.Repository) *policyv1.ImagePolicy {
	return &policyv1.ImagePolicy{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name},
		TypeMeta:   metav1.TypeMeta{Kind: "ImagePolicy"},
		Spec: policyv1.ImagePolicySpec{
			Repositories: repos,
		},
	}
}

func setup(policies []runtime.Object) (*Client, policyclientset.Interface) {
	clientSet := fake.NewSimpleClientset(policies...)
	return NewClient(clientSet), clientSet
}

func TestClient_GetPolicyToEnforce(t *testing.T) {

	tests := []struct {
		name      string
		namespace string
		image     string
		policies  []runtime.Object
		want      *policyv1.Policy
		wantErr   error
	}{
		{
			name:      "No Image policy, but relevant cluster policy: return cluster policy",
			image:     "icr.io/hello/world",
			namespace: "default",
			policies:  []runtime.Object{createClusterImagePolicy("policy-one", []policyv1.Repository{helloWorldRepositoryTrustEnabled})},
			want:      &enabledTrustPolicy,
		},
		{
			name:      "No Image policy, cluster policy has no repo match: return error",
			image:     "icr.io/hello/world",
			namespace: "default",
			policies:  []runtime.Object{createClusterImagePolicy("policy-one", []policyv1.Repository{helloEarthRepositoryTrustEnabled})},
			wantErr:   errors.New(`Deny "icr.io/hello/world", no matching repositories in ClusterImagePolicy and no ImagePolicies in the "default" namespace`),
		},
		{
			name:      "No Image policy or cluster policy: return error",
			image:     "icr.io/hello/world",
			namespace: "default",
			wantErr:   errors.New(`Deny "icr.io/hello/world", no image policies or cluster polices`),
		},
		{
			name:      "No Image policy, but multiple cluster policy: return cluster policy",
			image:     "icr.io/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createClusterImagePolicy("policy-one", []policyv1.Repository{helloEarthRepositoryTrustDisabled}),
				createClusterImagePolicy("policy-two", []policyv1.Repository{helloWorldRepositoryTrustEnabled}),
			},
			want: &enabledTrustPolicy,
		},
		{
			name:      "No Image policy, but multiple repos in a single cluster policy: return cluster policy",
			image:     "icr.io/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createClusterImagePolicy("policy-one", []policyv1.Repository{helloEarthRepositoryTrustDisabled, helloWorldRepositoryTrustEnabled}),
			},
			want: &enabledTrustPolicy,
		},
		{
			name:      "Relevant Image policy and cluster policy: return image policy",
			image:     "icr.io/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createClusterImagePolicy("policy-one", []policyv1.Repository{helloWorldRepositoryTrustDisabled}),
				createImagePolicy("policy-two", "default", []policyv1.Repository{helloWorldRepositoryTrustEnabled}),
			},
			want: &enabledTrustPolicy,
		},
		{
			name:      "Relevant Image policies in different namspaces: return correct image policy",
			image:     "icr.io/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createImagePolicy("policy-one", "nonDefault", []policyv1.Repository{helloWorldRepositoryTrustDisabled}),
				createImagePolicy("policy-two", "default", []policyv1.Repository{helloWorldRepositoryTrustEnabled}),
			},
			want: &enabledTrustPolicy,
		},
		{
			name:      "Image policies without relevant repository: return error",
			image:     "icr.io/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createImagePolicy("policy-one", "default", []policyv1.Repository{helloEarthRepositoryTrustDisabled}),
			},
			wantErr: errors.New(`Deny "icr.io/hello/world", no matching repositories in the ImagePolicies`),
		},
		{
			name:      "Image policies without relevant repository but matching cluster policy: return error",
			image:     "icr.io/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createImagePolicy("policy-one", "default", []policyv1.Repository{helloEarthRepositoryTrustEnabled}),
				createClusterImagePolicy("policy-one", []policyv1.Repository{helloWorldRepositoryTrustEnabled}),
			},
			wantErr: errors.New(`Deny "icr.io/hello/world", no matching repositories in the ImagePolicies`),
		},
		{
			name:      "Multiple Image policies in the same namespace: return relevant image policy",
			image:     "icr.io/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createImagePolicy("policy-one", "default", []policyv1.Repository{helloEarthRepositoryTrustDisabled}),
				createImagePolicy("policy-two", "default", []policyv1.Repository{helloWorldRepositoryTrustEnabled}),
			},
			want: &enabledTrustPolicy,
		},
		{
			name:      "Image policies with multiple repos: return relevant image policy",
			image:     "icr.io/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createImagePolicy("policy-one", "default", []policyv1.Repository{helloEarthRepositoryTrustDisabled, helloWorldRepositoryTrustEnabled}),
			},
			want: &enabledTrustPolicy,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := setup(tt.policies)
			got, err := client.GetPolicyToEnforce(tt.namespace, tt.image)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestClient_getImagePolicyList(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		policies  []runtime.Object
		wantList  *policyv1.ImagePolicyList
		wantErr   bool
	}{
		{
			name:     "returns multiple image policies from a namespace",
			policies: []runtime.Object{imagePolicyOne, imagePolicyTwo},
			wantList: &policyv1.ImagePolicyList{
				Items: []policyv1.ImagePolicy{*imagePolicyOne, *imagePolicyTwo},
			},
		},
		{
			name:     "returns no image policies from a namespace",
			wantList: &policyv1.ImagePolicyList{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := setup(tt.policies)
			got, err := client.getImagePolicyList(tt.namespace)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantList, got)
			}
		})
	}
}

func TestClient_getClusterImagePolicyList(t *testing.T) {
	tests := []struct {
		name     string
		policies []runtime.Object
		wantList *policyv1.ClusterImagePolicyList
		wantErr  bool
	}{
		{
			name:     "returns multiple cluster image policies from a namespace",
			policies: []runtime.Object{clusterImagePolicyOne, clusterImagePolicyTwo},
			wantList: &policyv1.ClusterImagePolicyList{
				Items: []policyv1.ClusterImagePolicy{*clusterImagePolicyOne, *clusterImagePolicyTwo},
			},
		},
		{
			name:     "returns no cluster image policies from a namespace",
			wantList: &policyv1.ClusterImagePolicyList{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := setup(tt.policies)
			got, err := client.getClusterImagePolicyList()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantList, got)
			}
		})
	}
}
