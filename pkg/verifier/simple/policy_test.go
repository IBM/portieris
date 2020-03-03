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
	"testing"

	"github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/stretchr/testify/assert"
)

func TestTransformPolicyInvalidType(t *testing.T) {
	simplePolicy := &v1beta1.Simple{
		Type: "invalid",
	}
	policy, err := transformPolicy(simplePolicy)
	assert.Nil(t, policy, "unexpected")
	assert.Error(t, err, "expected")
}

func TestTransformPolicyRejectType(t *testing.T) {
	simplePolicy := &v1beta1.Simple{
		Type: "reject",
	}
	policy, err := transformPolicy(simplePolicy)
	assert.NotNil(t, policy, "unexpected")
	assert.Nil(t, err, "expected")
}

func TestTransformPolicyAcceptType(t *testing.T) {
	simplePolicy := &v1beta1.Simple{
		Type: "insecureAcceptAnything",
	}
	policy, err := transformPolicy(simplePolicy)
	assert.NotNil(t, policy, "unexpected")
	assert.Nil(t, err, "expected")
}

func TestTransformPolicySignedByType(t *testing.T) {
	simplePolicy := &v1beta1.Simple{
		Type:           "signedBy",
		KeyType:        "GPGKey",
		KeyData:        "somedata",
		SignedIdentity: v1beta1.IdentityRequirement{},
	}
	policy, err := transformPolicy(simplePolicy)
	assert.NotNil(t, policy, "unexpected")
	assert.Nil(t, err, "expected")
}

func TestTransformPolicySignedByTypeInvalidKey(t *testing.T) {
	simplePolicy := &v1beta1.Simple{
		Type:           "signedBy",
		KeyType:        "GPGKey",
		KeyData:        "******",
		SignedIdentity: v1beta1.IdentityRequirement{},
	}
	policy, err := transformPolicy(simplePolicy)
	assert.Nil(t, policy, "unexpected")
	assert.Error(t, err, "expected")
}

func TestTransformPolicySignedByTypeEmptyKey(t *testing.T) {
	simplePolicy := &v1beta1.Simple{
		Type:           "signedBy",
		KeyType:        "GPGKey",
		KeyData:        "",
		SignedIdentity: v1beta1.IdentityRequirement{},
	}
	policy, err := transformPolicy(simplePolicy)
	assert.Nil(t, policy, "unexpected")
	assert.Error(t, err, "expected")
}

func TestTransformPolicySignedByTypeInvalidSignedId(t *testing.T) {
	simplePolicy := &v1beta1.Simple{
		Type:    "signedBy",
		KeyType: "GPGKey",
		KeyData: "somedata",
		SignedIdentity: v1beta1.IdentityRequirement{
			Type: "invalid",
		},
	}
	policy, err := transformPolicy(simplePolicy)
	assert.Nil(t, policy, "unexpected")
	assert.Error(t, err, "expected")
}

func TestTransformPolicySignedByTypeMatchExact(t *testing.T) {
	simplePolicy := &v1beta1.Simple{
		Type:    "signedBy",
		KeyType: "GPGKey",
		KeyData: "somedata",
		SignedIdentity: v1beta1.IdentityRequirement{
			Type: "matchExact",
		},
	}
	policy, err := transformPolicy(simplePolicy)
	assert.NotNil(t, policy, "expected")
	assert.Nil(t, err, "no error")
}

func TestTransformPolicySignedByTypeMatchRepository(t *testing.T) {
	simplePolicy := &v1beta1.Simple{
		Type:    "signedBy",
		KeyType: "GPGKey",
		KeyData: "somedata",
		SignedIdentity: v1beta1.IdentityRequirement{
			Type: "matchRepository",
		},
	}
	policy, err := transformPolicy(simplePolicy)
	assert.NotNil(t, policy, "expected")
	assert.Nil(t, err, "no error")
}

func TestTransformPolicySignedByTypeMatchExactRepository(t *testing.T) {
	simplePolicy := &v1beta1.Simple{
		Type:    "signedBy",
		KeyType: "GPGKey",
		KeyData: "somedata",
		SignedIdentity: v1beta1.IdentityRequirement{
			Type:             "matchExactRepository",
			DockerRepository: "repository",
		},
	}
	policy, err := transformPolicy(simplePolicy)
	assert.NotNil(t, policy, "expected")
	assert.Nil(t, err, "no error")
}

func TestTransformPolicySignedByTypeMatchExactReference(t *testing.T) {
	simplePolicy := &v1beta1.Simple{
		Type:    "signedBy",
		KeyType: "GPGKey",
		KeyData: "somedata",
		SignedIdentity: v1beta1.IdentityRequirement{
			Type:            "matchExactReference",
			DockerReference: "reg.io/image:tag",
		},
	}
	policy, err := transformPolicy(simplePolicy)
	assert.NotNil(t, policy, "expected")
	assert.Nil(t, err, "no error")
}
