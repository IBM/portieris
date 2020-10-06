// Copyright 2018, 2020 Portieris Authors.
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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/IBM/portieris/helpers/credential"
	"github.com/IBM/portieris/helpers/image"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/IBM/portieris/pkg/policy"
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
	// policyClient is a securityenforcementclientset with a wrapper for retrieving the relevant policy spec
	policyClient policy.Interface
	// Enforcer is used to check that containers satisfy constraints set by a policy
	Enforcer
}

// NewController creates a new controller object from the various clients passed in
func NewController(kubeWrapper kubernetes.WrapperInterface, policyClient policy.Interface, nv *notaryverifier.Verifier) *Controller {
	enforcer := NewEnforcer(kubeWrapper, nv)
	return &Controller{
		kubeClientsetWrapper: kubeWrapper,
		policyClient:         policyClient,
		Enforcer:             enforcer,
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

	return c.admitPod(admissionRequest.Namespace, podSpecLocation, *ps)
}

func (c *Controller) admitPod(namespace, specPath string, pod corev1.PodSpec) *admissionv1beta1.AdmissionResponse {
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

		newPatches, denials, err := c.getPatchesForContainers(containerType, namespace, specPath, pod, containers)
		a.StringsToAdmissionResponse(denials)
		if err != nil {
			a.ToAdmissionResponse(err)
			a.Flush()
		}
		patches = append(patches, newPatches...)
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

func (c *Controller) getPatchesForContainers(containerType, namespace, specPath string, pod corev1.PodSpec, containers []corev1.Container) ([]types.JSONPatch, []string, error) {
	patches := []types.JSONPatch{}
	denials := []string{}

	// for each container of this type
	for containerIndex, container := range containers {
		img, err := image.NewReference(container.Image)
		if err != nil {
			glog.Error(err)
			denials = append(denials, fmt.Sprintf("Deny %q, invalid image name", container.Image))
			continue
		}

		glog.Infof("Getting policy for container image: %s   namespace: %s", img.String(), namespace)
		containerPolicy, err := c.policyClient.GetPolicyToEnforce(namespace, img.String())
		if err != nil {
			denials = append(denials, err.Error())
			continue
		}

		credentialCandidates := c.getPodCredentials(namespace, img, pod)

		digest, deny, err := c.Enforcer.DigestByPolicy(namespace, img, credentialCandidates, containerPolicy)
		if err != nil {
			return patches, denials, err
		}
		if deny != nil {
			denials = append(denials, deny.Error())
			continue
		}
		if digest != nil {
			// convert digest to patch
			glog.Infof("Mutation #: %s %d  Image name: %s", containerType, containerIndex, img.String())
			if strings.Contains(container.Image, img.String()) {
				// ISSUE: https://github.com/IBM/portieris/issues/90
				glog.Infof("Mutated to: %s@sha256:%s", img.NameWithoutTag(), digest.String())
				patch := types.JSONPatch{
					Op:    "replace",
					Path:  fmt.Sprintf("%s/%s/%d/image", specPath, containerType, containerIndex),
					Value: fmt.Sprintf("%s@sha256:%s", img.NameWithoutTag(), digest.String()),
				}
				glog.Infof("Patch: %v", patch)
				patches = append(patches, patch)
			}
		}
	}

	return patches, denials, nil
}

func (c *Controller) getPodCredentials(namespace string, img *image.Reference, pod corev1.PodSpec) credential.Credentials {
	var creds credential.Credentials
	for _, secret := range pod.ImagePullSecrets {
		username, password, err := c.kubeClientsetWrapper.GetSecretToken(namespace, secret.Name, img.GetHostname())
		if err != nil {
			glog.Error(err)
			continue
		}
		cred := credential.Credential{
			Username: username,
			Password: password,
		}
		creds = append(creds, cred)
		glog.Infof("ImagePullSecret %s/%s found", namespace, secret.Name)
	}
	return creds
}
