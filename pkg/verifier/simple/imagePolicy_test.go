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

var policy2 = &v1beta1.Policy{
	Simple: &v1beta1.Simple{
		KeyData: "mQENBF15FJUBCAC+RDRL14lFAeVUAQrsg7XU3tLEb6Goy+XADZL1VLOgjDNqbkM8UnHRlGAVcMkui/vaiF/PHQchIc64vbQFjHsswxuNiRpL1n72k3dq9fQkdE5uMFtgm/LYlqFJDOhdFWarUUvBW1rTAwZAxWQSsZGGzTasSzA2JtiAR51qAMF3JZxV6RARvIAf4XqdVTG/LhbA15GTDx4zGI30hb29pVV6d6nV+qEvXP4QTOQ27dBv8ZN1d8rDSQI7fhb7xoXt6xqsSjFl+rgCCyoRbCCWpdQIhcBLqK4O8MEYp2M+D5YpO8WV4OM9EDx9YhFpsNaOirzfd1ZQZ+vUpT7qFq2kqen1ABEBAAG0KFN0dWFydCBIYXl0b24gPHN0dWFydC5oYXl0b25AdWsuaWJtLmNvbT6JAVQEEwEIAD4WIQR3TcmcAGUBN1Ici7Pxx2Awu2yqjQUCXXkUlQIbAwUJA8JnAAULCQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRDxx2Awu2yqjbtkCACaUuWxKuuw1+kDKy06Ir1/+mrPrNHiPndmmOxvVrhTJqmukXHSq3HgXxJWCnU+ubhmnKV7StK8pG8bNSFVtTxVKfcedQGZlvQY6/avGfd7BysKpFQl9QjAwojcirVFmOzA/bfVY1lGUivnxOUwzPsznngl+fsG3s9VEYYnry8DoeewR6Xy4d+EB/phTK61Oh+gB7Gic3wnf5HJMKWUl4GchyzPpxRi2az8tfBS9tkHNaqtzIb5QV9mZ2/LQ/opXvE56yyM9jRaQqKdeO6MtQus2AI8w2NYl3PNFK/NcblcJKEMlNJ6d+j/stk/mCNmRvebutiDYuTKCjqmW1lYleP9uQENBF15FJUBCACckqwmDBKPp93nXSyJzH8Di9cC7cL58Q6pGjcwG4GhanfbDxR0eDem/l2Ccn3lVoBdSM8P5SGRCbQdgUNfreHofjp6idcFg/rkjc2Q5BS+fQ0HDfFuLMnS3eKuwFbRSHtNKDP/fKiIgKzx4ra55S7lgVX8Skh11acFHkuH+9xpeV+bv84F28TCZ+pL+G2XYRqYKNvAnGB5PmCfUwZJlgJEu29F7sYiplYD5nIWBSz0ZwzWM+wSGCdntgxYuw+7c+3vfOwsgAOpgqXXNHwpRSd1xazbTpu8Kz1nWeZ8w8aPmYKuo9+ucMbpzYpqmyiXb1DiHbxOVsE3ZM6kBIyl7H5HABEBAAGJATwEGAEIACYWIQR3TcmcAGUBN1Ici7Pxx2Awu2yqjQUCXXkUlQIbDAUJA8JnAAAKCRDxx2Awu2yqjQ4xCACRYNG/6JpKuOjsU/LSpw8GrBNjFMlzNdiPOdHiW/gglBbMJB3LJJrM4TvMcFsqmuKUh1j7/gO9GUhm3VIRxZXxmble0sEh5n6Tpz0HoZb2ndvi+tqbMm1ufDP9pbIXOZzdksywrAX3283vjDUTlDog7qYBzQEG6TK68RGDKGobDtBIoR9S/enHoAkrWONKJ9uyJw2cIpx72MPXiMqP6vnLExdgp01NoEQx1UPfy/Y9gJ5aGaUUBDG7i6twpeTo9XFyJihrU5tFfrzT6iuGggxFfJoCgxVAKzXJnGTulcClquAOmMCFKqxbkOTIUy0uATSGF4pIvGu0Edi0GzvfCKST",
		KeyType: "GPGKeys",
		Type:    "signedBy",
		SignedIdentity: v1beta1.IdentityRequirement{
			Type: "",
		},
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
