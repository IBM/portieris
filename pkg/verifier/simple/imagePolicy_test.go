// Copyright 2020, 2022 Portieris Authors.
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

// Implementation of verify against image policy interface

package simple

import (
	"testing"

	"github.com/IBM/portieris/helpers/credential"
	"github.com/containers/image/v5/signature"
	"github.com/stretchr/testify/assert"
)

var policyRequirementInsecure = signature.NewPRInsecureAcceptAnything()

// Cover error paths - good path is covered in e2e tests
func TestVerifyByPolicy(t *testing.T) {
	tests := []struct {
		name        string
		image       string
		credentials credential.Credentials
		policies    *signature.PolicyRequirement
		wantErr     bool
		errMsg      string
		wantDeny    bool
		denyMsg     string
	}{
		{
			name:        "bad image",
			image:       "blahBLAHblah",
			credentials: credential.Credentials{},
			policies:    &policyRequirementInsecure,
			wantErr:     true,
			errMsg:      "name must be lowercase",
		},
		{
			name:        "no creds",
			image:       "docker.io/library/busybox",
			credentials: credential.Credentials{},
			policies:    &policyRequirementInsecure,
			wantErr:     false,
		},
		{
			name:        "extra creds",
			image:       "docker.io/library/busybox",
			credentials: credential.Credentials{{Username: "user", Password: "password"}},
			policies:    &policyRequirementInsecure,
			wantErr:     false,
		},
		{
			name:        "bad registry",
			image:       "nonsuch.io/library/busybox",
			credentials: credential.Credentials{{Username: "user", Password: "password"}},
			policies:    &policyRequirementInsecure,
			wantErr:     true,
			errMsg:      "pinging docker registry ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := &signature.Policy{
				Default: signature.PolicyRequirements{signature.NewPRReject()},
				Transports: map[string]signature.PolicyTransportScopes{
					"docker": {
						"": {*tt.policies},
					},
				},
			}
			digest, deny, err := verifier{}.VerifyByPolicy(tt.image, tt.credentials, "", policy)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg, "unexpected error")
			} else {
				assert.NoError(t, err)
			}
			if tt.wantDeny {
				assert.Error(t, deny)
				assert.Contains(t, deny.Error(), tt.denyMsg, "unexpected deny")
			} else {
				assert.NoError(t, deny)
			}
			if !tt.wantErr && !tt.wantDeny {
				assert.NotEmpty(t, digest, "should have a digest")
			}
		})
	}
}
