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
	"fmt"
	"testing"

	"github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
)

var keyValidGPG = "mQENBF15FJUBCAC+RDRL14lFAeVUAQrsg7XU3tLEb6Goy+XADZL1VLOgjDNqbkM8\n" +
	"UnHRlGAVcMkui/vaiF/PHQchIc64vbQFjHsswxuNiRpL1n72k3dq9fQkdE5uMFtg\n" +
	"m/LYlqFJDOhdFWarUUvBW1rTAwZAxWQSsZGGzTasSzA2JtiAR51qAMF3JZxV6RAR\n" +
	"vIAf4XqdVTG/LhbA15GTDx4zGI30hb29pVV6d6nV+qEvXP4QTOQ27dBv8ZN1d8rD\n" +
	"SQI7fhb7xoXt6xqsSjFl+rgCCyoRbCCWpdQIhcBLqK4O8MEYp2M+D5YpO8WV4OM9\n" +
	"EDx9YhFpsNaOirzfd1ZQZ+vUpT7qFq2kqen1ABEBAAG0KFN0dWFydCBIYXl0b24g\n" +
	"PHN0dWFydC5oYXl0b25AdWsuaWJtLmNvbT6JAVQEEwEIAD4WIQR3TcmcAGUBN1Ic\n" +
	"i7Pxx2Awu2yqjQUCXXkUlQIbAwUJA8JnAAULCQgHAgYVCgkICwIEFgIDAQIeAQIX\n" +
	"gAAKCRDxx2Awu2yqjbtkCACaUuWxKuuw1+kDKy06Ir1/+mrPrNHiPndmmOxvVrhT\n" +
	"JqmukXHSq3HgXxJWCnU+ubhmnKV7StK8pG8bNSFVtTxVKfcedQGZlvQY6/avGfd7\n" +
	"BysKpFQl9QjAwojcirVFmOzA/bfVY1lGUivnxOUwzPsznngl+fsG3s9VEYYnry8D\n" +
	"oeewR6Xy4d+EB/phTK61Oh+gB7Gic3wnf5HJMKWUl4GchyzPpxRi2az8tfBS9tkH\n" +
	"NaqtzIb5QV9mZ2/LQ/opXvE56yyM9jRaQqKdeO6MtQus2AI8w2NYl3PNFK/Ncblc\n" +
	"JKEMlNJ6d+j/stk/mCNmRvebutiDYuTKCjqmW1lYleP9uQENBF15FJUBCACckqwm\n" +
	"DBKPp93nXSyJzH8Di9cC7cL58Q6pGjcwG4GhanfbDxR0eDem/l2Ccn3lVoBdSM8P\n" +
	"5SGRCbQdgUNfreHofjp6idcFg/rkjc2Q5BS+fQ0HDfFuLMnS3eKuwFbRSHtNKDP/\n" +
	"fKiIgKzx4ra55S7lgVX8Skh11acFHkuH+9xpeV+bv84F28TCZ+pL+G2XYRqYKNvA\n" +
	"nGB5PmCfUwZJlgJEu29F7sYiplYD5nIWBSz0ZwzWM+wSGCdntgxYuw+7c+3vfOws\n" +
	"gAOpgqXXNHwpRSd1xazbTpu8Kz1nWeZ8w8aPmYKuo9+ucMbpzYpqmyiXb1DiHbxO\n" +
	"VsE3ZM6kBIyl7H5HABEBAAGJATwEGAEIACYWIQR3TcmcAGUBN1Ici7Pxx2Awu2yq\n" +
	"jQUCXXkUlQIbDAUJA8JnAAAKCRDxx2Awu2yqjQ4xCACRYNG/6JpKuOjsU/LSpw8G\n" +
	"rBNjFMlzNdiPOdHiW/gglBbMJB3LJJrM4TvMcFsqmuKUh1j7/gO9GUhm3VIRxZXx\n" +
	"mble0sEh5n6Tpz0HoZb2ndvi+tqbMm1ufDP9pbIXOZzdksywrAX3283vjDUTlDog\n" +
	"7qYBzQEG6TK68RGDKGobDtBIoR9S/enHoAkrWONKJ9uyJw2cIpx72MPXiMqP6vnL\n" +
	"Exdgp01NoEQx1UPfy/Y9gJ5aGaUUBDG7i6twpeTo9XFyJihrU5tFfrzT6iuGggxF\n" +
	"fJoCgxVAKzXJnGTulcClquAOmMCFKqxbkOTIUy0uATSGF4pIvGu0Edi0GzvfCKST"

//	"=Q8Ys" (armour checksum)
var keyValidB64 = "ZG8geW91IGtub3cgbWUK"
var keyInvalid = "******"

// Wrapper is a wrapper around kubeclientset that includes some helper functions for applying behaviour to kube resources
type TestWrapper struct {
	kubernetes.WrapperInterface
}

func (w *TestWrapper) GetSecretKey(namespace, secretName string) ([]byte, error) {
	switch secretName {
	case "validKeySecret":
		return []byte("-----BEGIN PGP PUBLIC KEY BLOCK-----\n" +
			"\n" +
			"mQENBF15FJUBCAC+RDRL14lFAeVUAQrsg7XU3tLEb6Goy+XADZL1VLOgjDNqbkM8\n" +
			"UnHRlGAVcMkui/vaiF/PHQchIc64vbQFjHsswxuNiRpL1n72k3dq9fQkdE5uMFtg\n" +
			"m/LYlqFJDOhdFWarUUvBW1rTAwZAxWQSsZGGzTasSzA2JtiAR51qAMF3JZxV6RAR\n" +
			"vIAf4XqdVTG/LhbA15GTDx4zGI30hb29pVV6d6nV+qEvXP4QTOQ27dBv8ZN1d8rD\n" +
			"SQI7fhb7xoXt6xqsSjFl+rgCCyoRbCCWpdQIhcBLqK4O8MEYp2M+D5YpO8WV4OM9\n" +
			"EDx9YhFpsNaOirzfd1ZQZ+vUpT7qFq2kqen1ABEBAAG0KFN0dWFydCBIYXl0b24g\n" +
			"PHN0dWFydC5oYXl0b25AdWsuaWJtLmNvbT6JAVQEEwEIAD4WIQR3TcmcAGUBN1Ic\n" +
			"i7Pxx2Awu2yqjQUCXXkUlQIbAwUJA8JnAAULCQgHAgYVCgkICwIEFgIDAQIeAQIX\n" +
			"gAAKCRDxx2Awu2yqjbtkCACaUuWxKuuw1+kDKy06Ir1/+mrPrNHiPndmmOxvVrhT\n" +
			"JqmukXHSq3HgXxJWCnU+ubhmnKV7StK8pG8bNSFVtTxVKfcedQGZlvQY6/avGfd7\n" +
			"BysKpFQl9QjAwojcirVFmOzA/bfVY1lGUivnxOUwzPsznngl+fsG3s9VEYYnry8D\n" +
			"oeewR6Xy4d+EB/phTK61Oh+gB7Gic3wnf5HJMKWUl4GchyzPpxRi2az8tfBS9tkH\n" +
			"NaqtzIb5QV9mZ2/LQ/opXvE56yyM9jRaQqKdeO6MtQus2AI8w2NYl3PNFK/Ncblc\n" +
			"JKEMlNJ6d+j/stk/mCNmRvebutiDYuTKCjqmW1lYleP9uQENBF15FJUBCACckqwm\n" +
			"DBKPp93nXSyJzH8Di9cC7cL58Q6pGjcwG4GhanfbDxR0eDem/l2Ccn3lVoBdSM8P\n" +
			"5SGRCbQdgUNfreHofjp6idcFg/rkjc2Q5BS+fQ0HDfFuLMnS3eKuwFbRSHtNKDP/\n" +
			"fKiIgKzx4ra55S7lgVX8Skh11acFHkuH+9xpeV+bv84F28TCZ+pL+G2XYRqYKNvA\n" +
			"nGB5PmCfUwZJlgJEu29F7sYiplYD5nIWBSz0ZwzWM+wSGCdntgxYuw+7c+3vfOws\n" +
			"gAOpgqXXNHwpRSd1xazbTpu8Kz1nWeZ8w8aPmYKuo9+ucMbpzYpqmyiXb1DiHbxO\n" +
			"VsE3ZM6kBIyl7H5HABEBAAGJATwEGAEIACYWIQR3TcmcAGUBN1Ici7Pxx2Awu2yq\n" +
			"jQUCXXkUlQIbDAUJA8JnAAAKCRDxx2Awu2yqjQ4xCACRYNG/6JpKuOjsU/LSpw8G\n" +
			"rBNjFMlzNdiPOdHiW/gglBbMJB3LJJrM4TvMcFsqmuKUh1j7/gO9GUhm3VIRxZXx\n" +
			"mble0sEh5n6Tpz0HoZb2ndvi+tqbMm1ufDP9pbIXOZzdksywrAX3283vjDUTlDog\n" +
			"7qYBzQEG6TK68RGDKGobDtBIoR9S/enHoAkrWONKJ9uyJw2cIpx72MPXiMqP6vnL\n" +
			"Exdgp01NoEQx1UPfy/Y9gJ5aGaUUBDG7i6twpeTo9XFyJihrU5tFfrzT6iuGggxF\n" +
			"fJoCgxVAKzXJnGTulcClquAOmMCFKqxbkOTIUy0uATSGF4pIvGu0Edi0GzvfCKST\n" +
			"=Q8Ys\n" +
			"-----END PGP PUBLIC KEY BLOCK-----\n"), nil

	case "badKeySecret":
		return []byte("*****"), nil

	case "emptyKeySecret":
		return []byte{}, nil

	case "missingSecret":
		return nil, fmt.Errorf("secret not found")
	}
	glog.Fatalf("Unexpected test secretName: %s", secretName)
	return nil, nil
}

func TestTransformPolicy(t *testing.T) {
	tests := []struct {
		name           string
		image          string
		simplePolicies []v1beta1.SimpleRequirement
		wantErr        bool
		errMsg         string
	}{
		{
			name: "invalid simple type",
			simplePolicies: []v1beta1.SimpleRequirement{{
				Type: "invalid",
			}},
			wantErr: true,
			errMsg:  "invalid Type:",
		},
		{
			name: "reject type",
			simplePolicies: []v1beta1.SimpleRequirement{{
				Type: "reject",
			}},
			wantErr: false,
		},
		{
			name: "accept type",
			simplePolicies: []v1beta1.SimpleRequirement{{
				Type: "insecureAcceptAnything",
			}},
			wantErr: false,
		},
		{
			name: "signedByInvalidKey",
			simplePolicies: []v1beta1.SimpleRequirement{{
				Type:      "signedBy",
				KeySecret: "badKeySecret",
			}},
			wantErr: true,
			errMsg:  "EOF", // TODO
		},
		{
			name: "signedByValidKey",
			simplePolicies: []v1beta1.SimpleRequirement{{
				Type:      "signedBy",
				KeySecret: "validKeySecret",
			}},
			wantErr: false,
		},
		{
			name: "signedByMissingKey",
			simplePolicies: []v1beta1.SimpleRequirement{{
				Type:      "signedBy",
				KeySecret: "missingSecret",
			}},
			wantErr: true,
			errMsg:  "secret not found",
		},
		{
			name: "signedByNoKey",
			simplePolicies: []v1beta1.SimpleRequirement{{
				Type: "signedBy",
			}},
			wantErr: true,
			errMsg:  "KeySecret missing in signedBy requirement",
		},
		{
			name: "signedByEmptyKey",
			simplePolicies: []v1beta1.SimpleRequirement{{
				Type:      "signedBy",
				KeySecret: "emptyKeySecret",
			}},
			wantErr: true,
			errMsg:  "Key: empty",
		},
		{
			name: "invalid SignedIdentity Type",
			simplePolicies: []v1beta1.SimpleRequirement{{
				Type:      "signedBy",
				KeySecret: "validKeySecret",
				SignedIdentity: v1beta1.IdentityRequirement{
					Type: "invalid",
				},
			}},
			wantErr: true,
			errMsg:  "invalid SignedIdentity Type:",
		},
		{
			name: "SignedIdentity matchExact",
			simplePolicies: []v1beta1.SimpleRequirement{{
				Type:      "signedBy",
				KeySecret: "validKeySecret",
				SignedIdentity: v1beta1.IdentityRequirement{
					Type: "matchExact",
				},
			}},
			wantErr: false,
		},
		{
			name: "SignedIdentity matchRepository",
			simplePolicies: []v1beta1.SimpleRequirement{{
				Type:      "signedBy",
				KeySecret: "validKeySecret",
				SignedIdentity: v1beta1.IdentityRequirement{
					Type: "matchRepository",
				},
			}},
			wantErr: false,
		},
		{
			name: "SignedIdentity matchExactRepository",
			simplePolicies: []v1beta1.SimpleRequirement{{
				Type:      "signedBy",
				KeySecret: "validKeySecret",
				SignedIdentity: v1beta1.IdentityRequirement{
					Type:             "matchExactRepository",
					DockerRepository: "repository",
				},
			}},
			wantErr: false,
		},
		{
			name: "SignedIdentity matchExactReference",
			simplePolicies: []v1beta1.SimpleRequirement{{
				Type:      "signedBy",
				KeySecret: "validKeySecret",
				SignedIdentity: v1beta1.IdentityRequirement{
					Type:            "matchExactReference",
					DockerReference: "reg.io/image:tag",
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testKubeWrapper = &TestWrapper{}
			policy, err := verifier{}.TransformPolicies(testKubeWrapper, "namespace", tt.simplePolicies)
			if tt.wantErr {
				assert.Error(t, err, "error expected")
				if err != nil {
					assert.Contains(t, err.Error(), tt.errMsg, "unexpected error message")
				}
				assert.Nil(t, policy, "policy result unexpected")
			} else {
				assert.NoError(t, err, "error unexpected")
				assert.NotNil(t, policy, "policy result expected")
			}
		})
	}
}
