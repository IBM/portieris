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

// Implementation of verify against containers/image policy interface

package simple

import (
	"encoding/base64"
	"fmt"

	"github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/containers/image/v5/signature"
)

// TransformPolicy .
func TransformPolicy(inPolicy *v1beta1.Simple) (*signature.Policy, error) {
	var policyRequirement signature.PolicyRequirement

	switch inPolicy.Type {
	case "insecureAcceptAnything":
		policyRequirement = signature.NewPRInsecureAcceptAnything()

	case "reject":
		policyRequirement = signature.NewPRReject()

	case "signedBy":
		keyData, err := policyKeyData(inPolicy.KeyData)
		if err != nil {
			return nil, err
		}
		signedIdentity, err := policySignedIdentity(inPolicy)
		if err != nil {
			return nil, err
		}
		policyRequirement, err = signature.NewPRSignedByKeyData(signature.SBKeyTypeGPGKeys, keyData, signedIdentity)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("simple policy unexpected type %s", inPolicy.Type)
	}

	policy := &signature.Policy{
		Default: signature.PolicyRequirements{signature.NewPRReject()},
		Transports: map[string]signature.PolicyTransportScopes{
			"docker": {
				"": {
					policyRequirement,
				},
			},
		},
	}
	return policy, nil
}

func policySignedIdentity(inPolicy *v1beta1.Simple) (signature.PolicyReferenceMatch, error) {
	switch inPolicy.SignedIdentity.Type {
	case "":
		return signature.NewPRMMatchRepoDigestOrExact(), nil
	case "matchExact":
		return signature.NewPRMMatchExact(), nil
	case "matchRepository":
		return signature.NewPRMMatchRepository(), nil
	case "matchExactReference":
		prm, err := signature.NewPRMExactReference(inPolicy.SignedIdentity.DockerReference)
		if err != nil {
			return nil, err
		}
		return prm, nil
	case "matchExactRepository":
		prm, err := signature.NewPRMExactRepository(inPolicy.SignedIdentity.DockerRepository)
		if err != nil {
			return nil, err
		}
		return prm, nil
	default:
		return nil, fmt.Errorf("unexpected SignedIdentityType: %s", inPolicy.SignedIdentity.Type)
	}
}

func policyKeyData(base64data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(base64data)
}
