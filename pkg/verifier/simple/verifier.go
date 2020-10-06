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

package simple

import (
	"bytes"

	"github.com/IBM/portieris/helpers/credential"
	"github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/containers/image/v5/signature"
)

// Verifier is for verifying simple signing
type Verifier interface {
	TransformPolicies(kWrapper kubernetes.WrapperInterface, namespace string, inPolicies []v1beta1.SimpleRequirement) (*signature.Policy, error)
	CreateRegistryDir(storeURL, storeUser, storePassword string) (string, error)
	VerifyByPolicy(imageToVerify string, credentials credential.Credentials, registriesConfigDir string, simplePolicy *signature.Policy) (*bytes.Buffer, error, error)
	RemoveRegistryDir(dirName string) error
}

type verifier struct{}

// NewVerifier creates a new Verfier
func NewVerifier() Verifier {
	return &verifier{}
}
