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

// Implementation of verify against image policy interface

package simple

import (
	"testing"

	"github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/stretchr/testify/assert"
)

var policyBad = &v1beta1.Policy{
	Simple: &v1beta1.Simple{
		Type: "invalid",
	},
}

var policyInsecure = &v1beta1.Policy{
	Simple: &v1beta1.Simple{
		Type: "insecureAcceptAnything",
	},
}

// Cover error paths - good path is covered in e2e tests
func TestVerifyByPolicy(t *testing.T) {
	tests := []struct {
		name         string
		nockRegistry bool
		image        string
		credentials  [][]string
		policies     *v1beta1.Policy
		wantErr      bool
		errMsg       string
		wantDeny     bool
		denyMsg      string
	}{
		{
			name:        "bad policy",
			image:       "docker.io/library/busybox",
			credentials: [][]string{},
			policies:    policyBad,
			wantErr:     true,
			errMsg:      "policy unexpected type",
		},
		{
			name:        "bad image",
			image:       "blahBLAHblah",
			credentials: [][]string{},
			policies:    policyInsecure,
			wantErr:     true,
			errMsg:      "name must be lowercase",
		},
		{
			name:        "no creds", // fails, future enhancement to cover no-auth registries
			image:       "docker.io/library/busybox",
			credentials: [][]string{},
			policies:    policyInsecure,
			wantErr:     true,
			errMsg:      "no valid ImagePullSecret",
		},
		{
			name:         "bad registry",
			nockRegistry: true,
			image:        "nonsuch.io/library/busybox",
			credentials:  [][]string{{"user", "password"}},
			policies:     policyInsecure,
			wantErr:      true,
			errMsg:       "pinging docker registry ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			digest, deny, err := VerifyByPolicy(tt.image, tt.credentials, tt.policies)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg, "unexpected error")
			} else {
				assert.NoError(t, err)
				//assert.Equal(t, tt.wantList, got)
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
