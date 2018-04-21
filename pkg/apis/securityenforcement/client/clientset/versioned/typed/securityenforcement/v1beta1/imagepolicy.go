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

package v1beta1

import (
	scheme "github.com/IBM/portieris/pkg/apis/securityenforcement/client/clientset/versioned/scheme"
	v1beta1 "github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ImagePoliciesGetter has a method to return a ImagePolicyInterface.
// A group's client should implement this interface.
type ImagePoliciesGetter interface {
	ImagePolicies(namespace string) ImagePolicyInterface
}

// ImagePolicyInterface has methods to work with ImagePolicy resources.
type ImagePolicyInterface interface {
	Create(*v1beta1.ImagePolicy) (*v1beta1.ImagePolicy, error)
	Update(*v1beta1.ImagePolicy) (*v1beta1.ImagePolicy, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1beta1.ImagePolicy, error)
	List(opts v1.ListOptions) (*v1beta1.ImagePolicyList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.ImagePolicy, err error)
	ImagePolicyExpansion
}

// imagePolicies implements ImagePolicyInterface
type imagePolicies struct {
	client rest.Interface
	ns     string
}

// newImagePolicies returns a ImagePolicies
func newImagePolicies(c *SecurityenforcementV1beta1Client, namespace string) *imagePolicies {
	return &imagePolicies{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the imagePolicy, and returns the corresponding imagePolicy object, and an error if there is any.
func (c *imagePolicies) Get(name string, options v1.GetOptions) (result *v1beta1.ImagePolicy, err error) {
	result = &v1beta1.ImagePolicy{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("imagepolicies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ImagePolicies that match those selectors.
func (c *imagePolicies) List(opts v1.ListOptions) (result *v1beta1.ImagePolicyList, err error) {
	result = &v1beta1.ImagePolicyList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("imagepolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested imagePolicies.
func (c *imagePolicies) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("imagepolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a imagePolicy and creates it.  Returns the server's representation of the imagePolicy, and an error, if there is any.
func (c *imagePolicies) Create(imagePolicy *v1beta1.ImagePolicy) (result *v1beta1.ImagePolicy, err error) {
	result = &v1beta1.ImagePolicy{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("imagepolicies").
		Body(imagePolicy).
		Do().
		Into(result)
	return
}

// Update takes the representation of a imagePolicy and updates it. Returns the server's representation of the imagePolicy, and an error, if there is any.
func (c *imagePolicies) Update(imagePolicy *v1beta1.ImagePolicy) (result *v1beta1.ImagePolicy, err error) {
	result = &v1beta1.ImagePolicy{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("imagepolicies").
		Name(imagePolicy.Name).
		Body(imagePolicy).
		Do().
		Into(result)
	return
}

// Delete takes name of the imagePolicy and deletes it. Returns an error if one occurs.
func (c *imagePolicies) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("imagepolicies").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *imagePolicies) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("imagepolicies").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched imagePolicy.
func (c *imagePolicies) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.ImagePolicy, err error) {
	result = &v1beta1.ImagePolicy{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("imagepolicies").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
