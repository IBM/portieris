// Copyright 2018, 2022 Portieris Authors.
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

package trust

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/IBM/portieris/helpers/credential"
	"github.com/IBM/portieris/helpers/image"
	policyv1 "github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/v1"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/IBM/portieris/pkg/notary"
	registryclient "github.com/IBM/portieris/pkg/registry"
	"github.com/golang/glog"
	store "github.com/theupdateframework/notary/storage"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var codec = serializer.NewCodecFactory(runtime.NewScheme())

// Interface is for verifying notary signatures
type Interface interface {
	VerifyByPolicy(string, *image.Reference, credential.Credentials, *policyv1.Policy) (*bytes.Buffer, error, error)
}

// Verifier is the notary controller
type Verifier struct {
	// kubeClientsetWrapper is a standard kubernetes clientset with a wrapper for retrieving podSpec from a given object
	kubeClientsetWrapper kubernetes.WrapperInterface
	// Trust Client
	trust notary.Interface
	// Container Registry client
	cr registryclient.Interface
}

// NewVerifier creates a new controller object from the various clients passed in
func NewVerifier(kubeWrapper kubernetes.WrapperInterface, trust notary.Interface, cr registryclient.Interface) *Verifier {
	return &Verifier{
		kubeClientsetWrapper: kubeWrapper,
		trust:                trust,
		cr:                   cr,
	}
}

// VerifyByPolicy ...
func (v *Verifier) VerifyByPolicy(namespace string, img *image.Reference, credentials credential.Credentials, policy *policyv1.Policy) (*bytes.Buffer, error, error) {
	notaryURL := policy.Trust.TrustServer
	var err error
	if notaryURL == "" {
		notaryURL, err = img.GetContentTrustURL()
		if err != nil {
			return nil, nil, fmt.Errorf("Trust Server/Image Configuration Error: %v", err.Error())
		}
	}

	var signers []Signer
	if policy.Trust.SignerSecrets != nil {
		// Generate a []Singer with the values for each signerSecret
		signers = make([]Signer, len(policy.Trust.SignerSecrets))
		for i, secretName := range policy.Trust.SignerSecrets {
			signers[i], err = v.getSignerSecret(namespace, secretName.Name)
			if err != nil {
				return nil, nil, fmt.Errorf("Deny %q, could not get signerSecret from your cluster, %s", img.String(), err.Error())
			}
		}
	}

	authEndpoint, err := v.trust.CheckAuthRequired(notaryURL, img)
	if err != nil {
		return nil, nil, fmt.Errorf("Deny %q, could not resolve the auth-endpoint, %s", img.String(), err.Error())
	}

        if authEndpoint == nil {
		credentials = append(credentials, credential.Credential{})
	}
	for _, credential := range credentials {
		var notaryToken string

		if authEndpoint != nil {
			notaryToken, err = v.cr.GetContentTrustToken(authEndpoint.URL, credential.Username, credential.Password, authEndpoint.Service, authEndpoint.Scope)
			if err != nil {
				glog.Error(err)
				continue
			}
		}

		digest, err := v.getDigest(notaryURL, img.NameWithoutTag(), notaryToken, img.GetTag(), signers)
		if err != nil {
			if strings.Contains(err.Error(), "401") {
				continue
			}

			if _, ok := err.(store.ErrServerUnavailable); ok {
				glog.Errorf("Trust server unavailable: %v", err)
				return nil, nil, fmt.Errorf("Deny %q, failed to get content trust information: %s", img.String(), err.Error())
			}
			return nil, fmt.Errorf("Deny %q, failed to get content trust information: %s", img.String(), err.Error()), nil
		}
		return digest, nil, nil
	}

	return nil, fmt.Errorf("Deny %q, no valid ImagePullSecret defined for %s", img.String(), img.String()), nil
}
