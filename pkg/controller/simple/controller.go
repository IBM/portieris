// Copyright 2020 Portieris Authors.
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

package simple

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/IBM/portieris/helpers/image"
	securityenforcementv1beta1 "github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/IBM/portieris/pkg/notary"
	"github.com/IBM/portieris/pkg/policy"
	registryclient "github.com/IBM/portieris/pkg/registry"
	"github.com/IBM/portieris/pkg/webhook"
	"github.com/IBM/portieris/types"
	"github.com/containers/image/v5/docker"
	"github.com/golang/glog"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var codec = serializer.NewCodecFactory(runtime.NewScheme())

// Controller is the notary controller
type Controller struct {
	// kubeClientsetWrapper is a standard kubernetes clientset with a wrapper for retrieving podSpec from a given object
	kubeClientsetWrapper kubernetes.WrapperInterface
	// policyClient is a securityenforcementclientset with a wrapper for retrieving the relevent policy spec
	policyClient policy.Interface
	// Container Registry client
	cr registryclient.Interface
}

// NewController creates a new controller object from the various clients passed in
func NewController(kubeWrapper kubernetes.WrapperInterface, policyClient policy.Interface, trust notary.Interface, cr registryclient.Interface) *Controller {
	return &Controller{
		kubeClientsetWrapper: kubeWrapper,
		policyClient:         policyClient,
		cr:                   cr,
	}
}

// Admit is the admissionRequest handler
func (c *Controller) Admit(admissionRequest *admissionv1beta1.AdmissionRequest) *admissionv1beta1.AdmissionResponse {
	glog.Infof("Processing Trust Admission Request for %s on %s", admissionRequest.Operation, admissionRequest.Name)

	podSpecLocation, ps, err := c.kubeClientsetWrapper.GetPodSpec(admissionRequest)
	switch err {
	case nil:
		break
	case kubernetes.ErrObjectHasParents, kubernetes.ErrObjectHasZeroReplicas:
		return &admissionv1beta1.AdmissionResponse{
			Allowed: true,
		}
	default:
		a := &webhook.AdmissionResponder{}
		a.ToAdmissionResponse(err)
		return a.Flush()
	}
	return c.mutatePodSpec(admissionRequest.Namespace, podSpecLocation, *ps)
}
func (c *Controller) mutatePodSpec(namespace, specPath string, pod corev1.PodSpec) *admissionv1beta1.AdmissionResponse {
	a := &webhook.AdmissionResponder{}
	patches := []types.JSONPatch{}

	// Iterate over each container image specified
	for _, containerType := range []string{"initContainers", "containers"} {
		var containers []corev1.Container
		switch containerType {
		case "initContainers":
			containers = pod.InitContainers
		case "containers":
			containers = pod.Containers
		default:
			a.StringToAdmissionResponse("Unhandled container type")
			return a.Flush()
		}

	containerLoop:
		for containerIndex, container := range containers {
			var policy *securityenforcementv1beta1.Policy
			img, err := image.NewReference(container.Image)
			if err != nil {
				glog.Error(err)
				a.StringToAdmissionResponse(fmt.Sprintf("Deny %q, invalid image name", container.Image))
				continue containerLoop
			}

			glog.Infof("Container Image: %s   Namespace: %s", img.String(), namespace)
			if policy, err = c.policyClient.GetPolicyToEnforce(namespace, img.String()); err != nil {
				a.StringToAdmissionResponse(err.Error())
				continue containerLoop
			} else if policy == nil || policy.Simple == nil {
				a.SetAllowed()
				continue containerLoop
			}

			// Trust is enforced
			glog.Info("Enforcing simple signatures")

			// Make sure image sure there is a ImagePullSecret defined
			// TODO: This prevents use of signed publically available images with publically available signing data
			if len(pod.ImagePullSecrets) == 0 {
				a.StringToAdmissionResponse(fmt.Sprintf("Deny %q, no ImagePullSecret defined for %s", img.String(), img.GetHostname()))
				continue containerLoop
			}

			imagePolicy, err := TransformPolicy(policy.Simple)
			if err != nil {
				glog.Error(err)
			}

		secretLoop:
			for _, secret := range pod.ImagePullSecrets {

				username, password, err := c.kubeClientsetWrapper.GetSecretToken(namespace, secret.Name, img.GetHostname())
				if err != nil {
					glog.Error(err)
					continue secretLoop
				}
				// verify secret without full verify?

				digest, err := VerifyByPolicy(img.String(), imagePolicy, username, password)
				// continue secretloop if secret fail

				if err != nil {
					switch err.(type) {
					case *docker.ErrUnauthorizedForCredentials:
						continue secretLoop
					default:
						glog.Warningf("Deny %q: %v", img.String(), err)
						continue containerLoop
					}
				}

				glog.Infof("Mutation #: %s %d  Image name: %s", containerType, containerIndex+1, img.String())
				if strings.Contains(container.Image, img.String()) {
					glog.Infof("Mutated to: %s@%s", img.String(), digest)
					patches = append(patches, types.JSONPatch{
						Op:    "replace",
						Path:  fmt.Sprintf("%s/%s/%s/image", specPath, containerType, strconv.Itoa(containerIndex)),
						Value: fmt.Sprintf("%s@%s", img.NameWithTag(), digest),
					})
				}
				a.SetAllowed()
				break secretLoop
			}
			if !a.IsAllowed() {
				a.StringToAdmissionResponse(fmt.Sprintf("Deny %q, no valid ImagePullSecret defined for %s", img.String(), img.GetHostname()))
			}
		}
	}

	if a.HasErrors() {
		return a.Flush()
	}

	// Apply patches if any
	if len(patches) > 0 {
		jsonPatch, err := json.Marshal(patches)
		if err != nil {
			a.StringToAdmissionResponse(fmt.Sprintf("Invalid Patch: %s", err.Error()))
			return a.Flush()
		}
		glog.Infof("Mutation patch: %s", string(jsonPatch))
		a.SetAllowed()
		a.SetPatch(jsonPatch)
	}

	return a.Flush()
}
