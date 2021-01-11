// Copyright 2020, 2021 Portieris Authors.
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

// Implementation of verify against containers/image policy interface

package simple

import (
	"fmt"

	policyv1 "github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/v1"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/containers/image/v5/signature"
)

// TransformPolicies from Portieris to container/image lib policies
func (v verifier) TransformPolicies(kWrapper kubernetes.WrapperInterface, namespace string, inPolicies []policyv1.SimpleRequirement) (*signature.Policy, error) {
	var policyRequirements []signature.PolicyRequirement

	for _, inPolicy := range inPolicies {
		var policyRequirement signature.PolicyRequirement

		switch inPolicy.Type {
		case "insecureAcceptAnything":
			policyRequirement = signature.NewPRInsecureAcceptAnything()

		case "reject":
			policyRequirement = signature.NewPRReject()

		case "signedBy":
			if inPolicy.KeySecret == "" {
				return nil, fmt.Errorf("KeySecret missing in signedBy requirement")
			}

			secretBytes, err := kWrapper.GetSecretKey(namespace, inPolicy.KeySecret)
			if err != nil {
				return nil, err
			}

			keyData, err := decodeArmoredKey(secretBytes)
			if err != nil {
				return nil, err
			}

			signedIdentity, err := policySignedIdentity(&inPolicy)
			if err != nil {
				return nil, err
			}

			policyRequirement, err = signature.NewPRSignedByKeyData(signature.SBKeyTypeGPGKeys, keyData, signedIdentity)
			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("simple policy invalid Type: %s", inPolicy.Type)
		}
		policyRequirements = append(policyRequirements, policyRequirement)
	}

	return &signature.Policy{
		Default: signature.PolicyRequirements{signature.NewPRReject()},
		Transports: map[string]signature.PolicyTransportScopes{
			"docker": {
				"": policyRequirements,
			},
		},
	}, nil
}

func policySignedIdentity(inPolicy *policyv1.SimpleRequirement) (signature.PolicyReferenceMatch, error) {
	switch inPolicy.SignedIdentity.Type {
	case "":
		return signature.NewPRMMatchRepoDigestOrExact(), nil
	case "matchExact":
		return signature.NewPRMMatchExact(), nil
	case "matchRepository":
		return signature.NewPRMMatchRepository(), nil
	case "matchExactReference":
		return signature.NewPRMExactReference(inPolicy.SignedIdentity.DockerReference)
	case "matchExactRepository":
		return signature.NewPRMExactRepository(inPolicy.SignedIdentity.DockerRepository)
	case "remapIdentity":
		return signature.NewPRMRemapIdentity(inPolicy.SignedIdentity.Prefix, inPolicy.SignedIdentity.SignedPrefix)
	default:
		return nil, fmt.Errorf("invalid SignedIdentity Type: %s", inPolicy.SignedIdentity.Type)
	}
}
