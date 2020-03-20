// Copyright 2018,2020 Portieris Authors.
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

package multi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/IBM/portieris/helpers/image"
	securityenforcementv1beta1 "github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/IBM/portieris/pkg/notary"
	"github.com/IBM/portieris/pkg/policy"
	registryclient "github.com/IBM/portieris/pkg/registry"
	simpleverifier "github.com/IBM/portieris/pkg/verifier/simple"
	notaryverifier "github.com/IBM/portieris/pkg/verifier/trust"
	"github.com/IBM/portieris/pkg/webhook"
	"github.com/IBM/portieris/types"
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
	// nv notary verifier
	nv *notaryverifier.Verifier
}

// NewController creates a new controller object from the various clients passed in
func NewController(kubeWrapper kubernetes.WrapperInterface, policyClient policy.Interface, trust notary.Interface, cr registryclient.Interface) *Controller {
	nv := notaryverifier.NewVerifier(kubeWrapper, trust, cr)
	return &Controller{
		kubeClientsetWrapper: kubeWrapper,
		policyClient:         policyClient,
		nv:                   nv,
	}
}

// Admit is the admissionRequest handler
func (c *Controller) Admit(admissionRequest *admissionv1beta1.AdmissionRequest) *admissionv1beta1.AdmissionResponse {
	glog.Infof("Processing admission request for %s on %s", admissionRequest.Operation, admissionRequest.Name)

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

	// for each container image subtype
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

		// for each container of this type
		for containerIndex, container := range containers {

			// move this?
			img, err := image.NewReference(container.Image)
			if err != nil {
				glog.Error(err)
				a.StringToAdmissionResponse(fmt.Sprintf("Deny %q, invalid image name", container.Image))
				continue
			}

			glog.Infof("Getting policy for container image: %s   namespace: %s", img.String(), namespace)
			containerPolicy, err := c.policyClient.GetPolicyToEnforce(namespace, img.String())
			if err != nil {
				a.ToAdmissionResponse(err)
				continue
			}

			credentialCandidates := c.getPodCredentials(namespace, img, pod)

			digest, deny, err := c.verifiedDigestByPolicy(namespace, img, credentialCandidates, containerPolicy)
			if err != nil {
				// error and return
				a.ToAdmissionResponse(err)
				return a.Flush()
			}
			if deny != nil {
				// deny and continue
				a.ToAdmissionResponse(deny)
				continue
			}
			if digest != nil {
				// convert digest to patch
				glog.Infof("Mutation #: %s %d  Image name: %s", containerType, containerIndex, img.String())
				if strings.Contains(container.Image, img.String()) {
					glog.Infof("Mutated to: %s@sha256:%s", img.NameWithTag(), digest.String())
					patch := types.JSONPatch{
						Op:    "replace",
						Path:  fmt.Sprintf("%s/%s/%d/image", specPath, containerType, containerIndex),
						Value: fmt.Sprintf("%s@sha256:%s", img.NameWithTag(), digest.String()),
					}
					glog.Infof("Patch: %v", patch)
					patches = append(patches, patch)
				}
			}
		}
	}

	if a.HasErrors() {
		glog.Info("Deny")
		return a.Flush()
	}

	// apply patches
	if len(patches) > 0 {
		jsonPatch, err := json.Marshal(patches)
		if err != nil {
			a.StringToAdmissionResponse(fmt.Sprintf("Invalid Patch: %s", err.Error()))
			return a.Flush()
		}
		glog.Infof("Mutation patch: %s", string(jsonPatch))
		a.SetPatch(jsonPatch)
	}
	a.SetAllowed()
	glog.Info("Allow")
	return a.Flush()
}

func (c *Controller) getPodCredentials(namespace string, img *image.Reference, pod corev1.PodSpec) [][]string {
	var creds [][]string
	for _, secret := range pod.ImagePullSecrets {
		username, password, err := c.kubeClientsetWrapper.GetSecretToken(namespace, secret.Name, img.GetHostname())
		if err != nil {
			glog.Error(err)
			continue
		}
		creds = append(creds, []string{username, password})
	}
	return creds
}

func (c *Controller) verifiedDigestByPolicy(namespace string, img *image.Reference, credentials [][]string, policy *securityenforcementv1beta1.Policy) (*bytes.Buffer, error, error) {

	// no policy indicates admission should be allowed, without mutation
	if policy == nil {
		return nil, nil, nil
	}

	var digest *bytes.Buffer
	var deny, err error
	if policy.Simple != nil {
		digest, deny, err = simpleverifier.VerifyByPolicy(img.String(), credentials, policy)
		if err != nil || deny != nil {
			return nil, deny, err
		}
	}

	if policy.Trust.Enabled != nil && *policy.Trust.Enabled == true {
		var notaryDigest *bytes.Buffer
		notaryDigest, deny, err = c.nv.VerifyByPolicy(namespace, img, credentials, policy)
		if err != nil || deny != nil {
			return nil, deny, err
		}
		if notaryDigest != nil {
			if digest != nil && notaryDigest != digest {
				return nil, fmt.Errorf("Notary signs conflicting digest: %v simple: %v", notaryDigest, digest), nil
			}
			digest = notaryDigest
		}
	}

	return digest, nil, nil
}
