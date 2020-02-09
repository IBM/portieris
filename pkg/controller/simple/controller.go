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

	"github.com/containers/image/v5/signature"

	"github.com/IBM/portieris/helpers/image"
	"github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
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

			imagePolicy, err := policyTransformer(policy.Simple)
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

func policyTransformer(inPolicy *v1beta1.Simple) (*signature.Policy, error) {
	switch inPolicy.Type {

	case "reject":
		return Reject, nil

	case "signedBy":
		policyString :=
			fmt.Sprintf(`{ "default": [ { "type": "signedBy",
		"keyType": "GPGKeys",
		"keyData": "mQENBF15FJUBCAC+RDRL14lFAeVUAQrsg7XU3tLEb6Goy+XADZL1VLOgjDNqbkM8UnHRlGAVcMkui/vaiF/PHQchIc64vbQFjHsswxuNiRpL1n72k3dq9fQkdE5uMFtgm/LYlqFJDOhdFWarUUvBW1rTAwZAxWQSsZGGzTasSzA2JtiAR51qAMF3JZxV6RARvIAf4XqdVTG/LhbA15GTDx4zGI30hb29pVV6d6nV+qEvXP4QTOQ27dBv8ZN1d8rDSQI7fhb7xoXt6xqsSjFl+rgCCyoRbCCWpdQIhcBLqK4O8MEYp2M+D5YpO8WV4OM9EDx9YhFpsNaOirzfd1ZQZ+vUpT7qFq2kqen1ABEBAAG0KFN0dWFydCBIYXl0b24gPHN0dWFydC5oYXl0b25AdWsuaWJtLmNvbT6JAVQEEwEIAD4WIQR3TcmcAGUBN1Ici7Pxx2Awu2yqjQUCXXkUlQIbAwUJA8JnAAULCQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRDxx2Awu2yqjbtkCACaUuWxKuuw1+kDKy06Ir1/+mrPrNHiPndmmOxvVrhTJqmukXHSq3HgXxJWCnU+ubhmnKV7StK8pG8bNSFVtTxVKfcedQGZlvQY6/avGfd7BysKpFQl9QjAwojcirVFmOzA/bfVY1lGUivnxOUwzPsznngl+fsG3s9VEYYnry8DoeewR6Xy4d+EB/phTK61Oh+gB7Gic3wnf5HJMKWUl4GchyzPpxRi2az8tfBS9tkHNaqtzIb5QV9mZ2/LQ/opXvE56yyM9jRaQqKdeO6MtQus2AI8w2NYl3PNFK/NcblcJKEMlNJ6d+j/stk/mCNmRvebutiDYuTKCjqmW1lYleP9uQENBF15FJUBCACckqwmDBKPp93nXSyJzH8Di9cC7cL58Q6pGjcwG4GhanfbDxR0eDem/l2Ccn3lVoBdSM8P5SGRCbQdgUNfreHofjp6idcFg/rkjc2Q5BS+fQ0HDfFuLMnS3eKuwFbRSHtNKDP/fKiIgKzx4ra55S7lgVX8Skh11acFHkuH+9xpeV+bv84F28TCZ+pL+G2XYRqYKNvAnGB5PmCfUwZJlgJEu29F7sYiplYD5nIWBSz0ZwzWM+wSGCdntgxYuw+7c+3vfOwsgAOpgqXXNHwpRSd1xazbTpu8Kz1nWeZ8w8aPmYKuo9+ucMbpzYpqmyiXb1DiHbxOVsE3ZM6kBIyl7H5HABEBAAGJATwEGAEIACYWIQR3TcmcAGUBN1Ici7Pxx2Awu2yqjQUCXXkUlQIbDAUJA8JnAAAKCRDxx2Awu2yqjQ4xCACRYNG/6JpKuOjsU/LSpw8GrBNjFMlzNdiPOdHiW/gglBbMJB3LJJrM4TvMcFsqmuKUh1j7/gO9GUhm3VIRxZXxmble0sEh5n6Tpz0HoZb2ndvi+tqbMm1ufDP9pbIXOZzdksywrAX3283vjDUTlDog7qYBzQEG6TK68RGDKGobDtBIoR9S/enHoAkrWONKJ9uyJw2cIpx72MPXiMqP6vnLExdgp01NoEQx1UPfy/Y9gJ5aGaUUBDG7i6twpeTo9XFyJihrU5tFfrzT6iuGggxFfJoCgxVAKzXJnGTulcClquAOmMCFKqxbkOTIUy0uATSGF4pIvGu0Edi0GzvfCKST"
	} ]
}`)
		return NewPolicyFromString(policyString)
	default:
	case "insecureAcceptAny":
	}
	return InsecureAcceptAnything, nil
}
