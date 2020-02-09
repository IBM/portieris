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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var policySignedByStuart = SafePolicyFromString(`{ "default": [ { "type": "signedBy",
"keyType": "GPGKeys",
"keyData": "mQENBF15FJUBCAC+RDRL14lFAeVUAQrsg7XU3tLEb6Goy+XADZL1VLOgjDNqbkM8UnHRlGAVcMkui/vaiF/PHQchIc64vbQFjHsswxuNiRpL1n72k3dq9fQkdE5uMFtgm/LYlqFJDOhdFWarUUvBW1rTAwZAxWQSsZGGzTasSzA2JtiAR51qAMF3JZxV6RARvIAf4XqdVTG/LhbA15GTDx4zGI30hb29pVV6d6nV+qEvXP4QTOQ27dBv8ZN1d8rDSQI7fhb7xoXt6xqsSjFl+rgCCyoRbCCWpdQIhcBLqK4O8MEYp2M+D5YpO8WV4OM9EDx9YhFpsNaOirzfd1ZQZ+vUpT7qFq2kqen1ABEBAAG0KFN0dWFydCBIYXl0b24gPHN0dWFydC5oYXl0b25AdWsuaWJtLmNvbT6JAVQEEwEIAD4WIQR3TcmcAGUBN1Ici7Pxx2Awu2yqjQUCXXkUlQIbAwUJA8JnAAULCQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRDxx2Awu2yqjbtkCACaUuWxKuuw1+kDKy06Ir1/+mrPrNHiPndmmOxvVrhTJqmukXHSq3HgXxJWCnU+ubhmnKV7StK8pG8bNSFVtTxVKfcedQGZlvQY6/avGfd7BysKpFQl9QjAwojcirVFmOzA/bfVY1lGUivnxOUwzPsznngl+fsG3s9VEYYnry8DoeewR6Xy4d+EB/phTK61Oh+gB7Gic3wnf5HJMKWUl4GchyzPpxRi2az8tfBS9tkHNaqtzIb5QV9mZ2/LQ/opXvE56yyM9jRaQqKdeO6MtQus2AI8w2NYl3PNFK/NcblcJKEMlNJ6d+j/stk/mCNmRvebutiDYuTKCjqmW1lYleP9uQENBF15FJUBCACckqwmDBKPp93nXSyJzH8Di9cC7cL58Q6pGjcwG4GhanfbDxR0eDem/l2Ccn3lVoBdSM8P5SGRCbQdgUNfreHofjp6idcFg/rkjc2Q5BS+fQ0HDfFuLMnS3eKuwFbRSHtNKDP/fKiIgKzx4ra55S7lgVX8Skh11acFHkuH+9xpeV+bv84F28TCZ+pL+G2XYRqYKNvAnGB5PmCfUwZJlgJEu29F7sYiplYD5nIWBSz0ZwzWM+wSGCdntgxYuw+7c+3vfOwsgAOpgqXXNHwpRSd1xazbTpu8Kz1nWeZ8w8aPmYKuo9+ucMbpzYpqmyiXb1DiHbxOVsE3ZM6kBIyl7H5HABEBAAGJATwEGAEIACYWIQR3TcmcAGUBN1Ici7Pxx2Awu2yqjQUCXXkUlQIbDAUJA8JnAAAKCRDxx2Awu2yqjQ4xCACRYNG/6JpKuOjsU/LSpw8GrBNjFMlzNdiPOdHiW/gglBbMJB3LJJrM4TvMcFsqmuKUh1j7/gO9GUhm3VIRxZXxmble0sEh5n6Tpz0HoZb2ndvi+tqbMm1ufDP9pbIXOZzdksywrAX3283vjDUTlDog7qYBzQEG6TK68RGDKGobDtBIoR9S/enHoAkrWONKJ9uyJw2cIpx72MPXiMqP6vnLExdgp01NoEQx1UPfy/Y9gJ5aGaUUBDG7i6twpeTo9XFyJihrU5tFfrzT6iuGggxFfJoCgxVAKzXJnGTulcClquAOmMCFKqxbkOTIUy0uATSGF4pIvGu0Edi0GzvfCKST"
} ]
}`)

const username = "iamapikey"
const password = "KtUf1Gjc9aZRx1E8KNxbZ9hEWcd3XhQ3xRdzJt5PwlTi"

func TestVerifyByPolicy_GoodPath(t *testing.T) {

	digest, err := VerifyByPolicy("uk.icr.io/stuarts/busybox:latest", policySignedByStuart, username, password)
	assert.Nil(t, err, "error")
	assert.True(t, strings.HasPrefix(digest, "sha256:"), "expected sha256: prefix")
}

func TestVerifyByPolicy_GoodPathOther(t *testing.T) {

	_, err := VerifyByPolicy("uk.icr.io/stuarts/other:latest", policySignedByStuart, username, password)
	assert.NotNil(t, err, "expected error")
}
