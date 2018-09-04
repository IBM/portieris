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
	"errors"
	"testing"

	securityenforcementclientset "github.com/IBM/portieris/pkg/apis/securityenforcement/client/clientset/versioned"
	"github.com/IBM/portieris/pkg/apis/securityenforcement/client/clientset/versioned/fake"
	securityenforcementv1beta1 "github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	trueBool       = true
	falseBool      = false
	imagePolicyOne = &securityenforcementv1beta1.ImagePolicy{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "policy-one"},
		TypeMeta:   metav1.TypeMeta{Kind: "ImagePolicy"},
		Spec: securityenforcementv1beta1.PolicySpec{
			Repositories: []securityenforcementv1beta1.Repository{{Name: "repo-one"}},
		},
	}
	imagePolicyTwo = &securityenforcementv1beta1.ImagePolicy{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "policy-two"},
		TypeMeta:   metav1.TypeMeta{Kind: "ImagePolicy"},
		Spec: securityenforcementv1beta1.PolicySpec{
			Repositories: []securityenforcementv1beta1.Repository{{Name: "repo-two"}},
		},
	}

	clusterImagePolicyOne = &securityenforcementv1beta1.ClusterImagePolicy{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "policy-one"},
		TypeMeta:   metav1.TypeMeta{Kind: "ClusterImagePolicy"},
		Spec: securityenforcementv1beta1.PolicySpec{
			Repositories: []securityenforcementv1beta1.Repository{{Name: "repo-one"}},
		},
	}
	clusterImagePolicyTwo = &securityenforcementv1beta1.ClusterImagePolicy{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "policy-two"},
		TypeMeta:   metav1.TypeMeta{Kind: "ClusterImagePolicy"},
		Spec: securityenforcementv1beta1.PolicySpec{
			Repositories: []securityenforcementv1beta1.Repository{{Name: "repo-two"}},
		},
	}

	helloWorldRepositoryTrustEnabled  = securityenforcementv1beta1.Repository{Name: "registry.bluemix.net/hello/world", Policy: enabledTrustPolicy}
	helloWorldRepositoryTrustDisabled = securityenforcementv1beta1.Repository{Name: "registry.bluemix.net/hello/world", Policy: disabledTrustPolicy}
	helloEarthRepositoryTrustEnabled  = securityenforcementv1beta1.Repository{Name: "registry.bluemix.net/hello/earth", Policy: enabledTrustPolicy}
	helloEarthRepositoryTrustDisabled = securityenforcementv1beta1.Repository{Name: "registry.bluemix.net/hello/earth", Policy: disabledTrustPolicy}

	enabledTrustPolicy  = securityenforcementv1beta1.Policy{Trust: securityenforcementv1beta1.Trust{Enabled: &trueBool}}
	disabledTrustPolicy = securityenforcementv1beta1.Policy{Trust: securityenforcementv1beta1.Trust{Enabled: &falseBool}}
)

func createClusterImagePolicy(name string, repos []securityenforcementv1beta1.Repository) *securityenforcementv1beta1.ClusterImagePolicy {
	return &securityenforcementv1beta1.ClusterImagePolicy{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		TypeMeta:   metav1.TypeMeta{Kind: "ClusterImagePolicy"},
		Spec: securityenforcementv1beta1.PolicySpec{
			Repositories: repos,
		},
	}
}

func createImagePolicy(name, namespace string, repos []securityenforcementv1beta1.Repository) *securityenforcementv1beta1.ImagePolicy {
	return &securityenforcementv1beta1.ImagePolicy{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name},
		TypeMeta:   metav1.TypeMeta{Kind: "ImagePolicy"},
		Spec: securityenforcementv1beta1.PolicySpec{
			Repositories: repos,
		},
	}
}

func setup(policies []runtime.Object) (*Client, securityenforcementclientset.Interface) {
	clientSet := fake.NewSimpleClientset(policies...)
	return NewClient(clientSet), clientSet
}

func TestClient_GetPolicyToEnforce(t *testing.T) {

	tests := []struct {
		name      string
		namespace string
		image     string
		policies  []runtime.Object
		want      *securityenforcementv1beta1.Policy
		wantErr   error
	}{
		{
			name:      "No Image policy, but relevant cluster policy: return cluster policy",
			image:     "registry.bluemix.net/hello/world",
			namespace: "default",
			policies:  []runtime.Object{createClusterImagePolicy("policy-one", []securityenforcementv1beta1.Repository{helloWorldRepositoryTrustEnabled})},
			want:      &enabledTrustPolicy,
		},
		{
			name:      "No Image policy, cluster policy has no repo match: return error",
			image:     "registry.bluemix.net/hello/world",
			namespace: "default",
			policies:  []runtime.Object{createClusterImagePolicy("policy-one", []securityenforcementv1beta1.Repository{helloEarthRepositoryTrustEnabled})},
			wantErr:   errors.New(`Deny "registry.bluemix.net/hello/world", no matching repositories in ClusterImagePolicy and no ImagePolicies in the "default" namespace`),
		},
		{
			name:      "No Image policy or cluster policy: return error",
			image:     "registry.bluemix.net/hello/world",
			namespace: "default",
			wantErr:   errors.New(`Deny "registry.bluemix.net/hello/world", no image policies or cluster polices`),
		},
		{
			name:      "No Image policy, but multiple cluster policy: return cluster policy",
			image:     "registry.bluemix.net/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createClusterImagePolicy("policy-one", []securityenforcementv1beta1.Repository{helloEarthRepositoryTrustDisabled}),
				createClusterImagePolicy("policy-two", []securityenforcementv1beta1.Repository{helloWorldRepositoryTrustEnabled}),
			},
			want: &enabledTrustPolicy,
		},
		{
			name:      "No Image policy, but multiple repos in a single cluster policy: return cluster policy",
			image:     "registry.bluemix.net/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createClusterImagePolicy("policy-one", []securityenforcementv1beta1.Repository{helloEarthRepositoryTrustDisabled, helloWorldRepositoryTrustEnabled}),
			},
			want: &enabledTrustPolicy,
		},
		{
			name:      "Relevant Image policy and cluster policy: return image policy",
			image:     "registry.bluemix.net/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createClusterImagePolicy("policy-one", []securityenforcementv1beta1.Repository{helloWorldRepositoryTrustDisabled}),
				createImagePolicy("policy-two", "default", []securityenforcementv1beta1.Repository{helloWorldRepositoryTrustEnabled}),
			},
			want: &enabledTrustPolicy,
		},
		{
			name:      "Relevant Image policies in different namspaces: return correct image policy",
			image:     "registry.bluemix.net/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createImagePolicy("policy-one", "nonDefault", []securityenforcementv1beta1.Repository{helloWorldRepositoryTrustDisabled}),
				createImagePolicy("policy-two", "default", []securityenforcementv1beta1.Repository{helloWorldRepositoryTrustEnabled}),
			},
			want: &enabledTrustPolicy,
		},
		{
			name:      "Image policies without relevant repository: return error",
			image:     "registry.bluemix.net/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createImagePolicy("policy-one", "default", []securityenforcementv1beta1.Repository{helloEarthRepositoryTrustDisabled}),
			},
			wantErr: errors.New(`Deny "registry.bluemix.net/hello/world", no matching repositories in the ImagePolicies`),
		},
		{
			name:      "Image policies without relevant repository but matching cluster policy: return error",
			image:     "registry.bluemix.net/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createImagePolicy("policy-one", "default", []securityenforcementv1beta1.Repository{helloEarthRepositoryTrustEnabled}),
				createClusterImagePolicy("policy-one", []securityenforcementv1beta1.Repository{helloWorldRepositoryTrustEnabled}),
			},
			wantErr: errors.New(`Deny "registry.bluemix.net/hello/world", no matching repositories in the ImagePolicies`),
		},
		{
			name:      "Multiple Image policies in the same namespace: return relevant image policy",
			image:     "registry.bluemix.net/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createImagePolicy("policy-one", "default", []securityenforcementv1beta1.Repository{helloEarthRepositoryTrustDisabled}),
				createImagePolicy("policy-two", "default", []securityenforcementv1beta1.Repository{helloWorldRepositoryTrustEnabled}),
			},
			want: &enabledTrustPolicy,
		},
		{
			name:      "Image policies with multiple repos: return relevant image policy",
			image:     "registry.bluemix.net/hello/world",
			namespace: "default",
			policies: []runtime.Object{
				createImagePolicy("policy-one", "default", []securityenforcementv1beta1.Repository{helloEarthRepositoryTrustDisabled, helloWorldRepositoryTrustEnabled}),
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
		wantList  *securityenforcementv1beta1.ImagePolicyList
		wantErr   bool
	}{
		{
			name:     "returns multiple image policies from a namespace",
			policies: []runtime.Object{imagePolicyOne, imagePolicyTwo},
			wantList: &securityenforcementv1beta1.ImagePolicyList{
				Items: []securityenforcementv1beta1.ImagePolicy{*imagePolicyOne, *imagePolicyTwo},
			},
		},
		{
			name:     "returns no image policies from a namespace",
			wantList: &securityenforcementv1beta1.ImagePolicyList{},
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
		wantList *securityenforcementv1beta1.ClusterImagePolicyList
		wantErr  bool
	}{
		{
			name:     "returns multiple cluster image policies from a namespace",
			policies: []runtime.Object{clusterImagePolicyOne, clusterImagePolicyTwo},
			wantList: &securityenforcementv1beta1.ClusterImagePolicyList{
				Items: []securityenforcementv1beta1.ClusterImagePolicy{*clusterImagePolicyOne, *clusterImagePolicyTwo},
			},
		},
		{
			name:     "returns no cluster image policies from a namespace",
			wantList: &securityenforcementv1beta1.ClusterImagePolicyList{},
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
