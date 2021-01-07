/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	portieriscloudibmcomv1 "github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
)

// FakeClusterImagePolicies implements ClusterImagePolicyInterface
type FakeClusterImagePolicies struct {
	Fake *FakePortierisV1
}

var clusterimagepoliciesResource = schema.GroupVersionResource{Group: "portieris.cloud.ibm.com", Version: "v1", Resource: "clusterimagepolicies"}

var clusterimagepoliciesKind = schema.GroupVersionKind{Group: "portieris.cloud.ibm.com", Version: "v1", Kind: "ClusterImagePolicy"}

// Get takes name of the clusterImagePolicy, and returns the corresponding clusterImagePolicy object, and an error if there is any.
func (c *FakeClusterImagePolicies) Get(name string, options v1.GetOptions) (result *portieriscloudibmcomv1.ClusterImagePolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clusterimagepoliciesResource, name), &portieriscloudibmcomv1.ClusterImagePolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*portieriscloudibmcomv1.ClusterImagePolicy), err
}

// List takes label and field selectors, and returns the list of ClusterImagePolicies that match those selectors.
func (c *FakeClusterImagePolicies) List(opts v1.ListOptions) (result *portieriscloudibmcomv1.ClusterImagePolicyList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clusterimagepoliciesResource, clusterimagepoliciesKind, opts), &portieriscloudibmcomv1.ClusterImagePolicyList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &portieriscloudibmcomv1.ClusterImagePolicyList{ListMeta: obj.(*portieriscloudibmcomv1.ClusterImagePolicyList).ListMeta}
	for _, item := range obj.(*portieriscloudibmcomv1.ClusterImagePolicyList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}
