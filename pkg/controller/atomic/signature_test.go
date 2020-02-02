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

package atomic

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"golang.org/x/crypto/openpgp/armor"

	"github.com/stretchr/testify/assert"
)

var publicKey1 = `-----BEGIN PGP PUBLIC KEY BLOCK-----

mQENBF15FJUBCAC+RDRL14lFAeVUAQrsg7XU3tLEb6Goy+XADZL1VLOgjDNqbkM8
UnHRlGAVcMkui/vaiF/PHQchIc64vbQFjHsswxuNiRpL1n72k3dq9fQkdE5uMFtg
m/LYlqFJDOhdFWarUUvBW1rTAwZAxWQSsZGGzTasSzA2JtiAR51qAMF3JZxV6RAR
vIAf4XqdVTG/LhbA15GTDx4zGI30hb29pVV6d6nV+qEvXP4QTOQ27dBv8ZN1d8rD
SQI7fhb7xoXt6xqsSjFl+rgCCyoRbCCWpdQIhcBLqK4O8MEYp2M+D5YpO8WV4OM9
EDx9YhFpsNaOirzfd1ZQZ+vUpT7qFq2kqen1ABEBAAG0KFN0dWFydCBIYXl0b24g
PHN0dWFydC5oYXl0b25AdWsuaWJtLmNvbT6JAVQEEwEIAD4WIQR3TcmcAGUBN1Ic
i7Pxx2Awu2yqjQUCXXkUlQIbAwUJA8JnAAULCQgHAgYVCgkICwIEFgIDAQIeAQIX
gAAKCRDxx2Awu2yqjbtkCACaUuWxKuuw1+kDKy06Ir1/+mrPrNHiPndmmOxvVrhT
JqmukXHSq3HgXxJWCnU+ubhmnKV7StK8pG8bNSFVtTxVKfcedQGZlvQY6/avGfd7
BysKpFQl9QjAwojcirVFmOzA/bfVY1lGUivnxOUwzPsznngl+fsG3s9VEYYnry8D
oeewR6Xy4d+EB/phTK61Oh+gB7Gic3wnf5HJMKWUl4GchyzPpxRi2az8tfBS9tkH
NaqtzIb5QV9mZ2/LQ/opXvE56yyM9jRaQqKdeO6MtQus2AI8w2NYl3PNFK/Ncblc
JKEMlNJ6d+j/stk/mCNmRvebutiDYuTKCjqmW1lYleP9uQENBF15FJUBCACckqwm
DBKPp93nXSyJzH8Di9cC7cL58Q6pGjcwG4GhanfbDxR0eDem/l2Ccn3lVoBdSM8P
5SGRCbQdgUNfreHofjp6idcFg/rkjc2Q5BS+fQ0HDfFuLMnS3eKuwFbRSHtNKDP/
fKiIgKzx4ra55S7lgVX8Skh11acFHkuH+9xpeV+bv84F28TCZ+pL+G2XYRqYKNvA
nGB5PmCfUwZJlgJEu29F7sYiplYD5nIWBSz0ZwzWM+wSGCdntgxYuw+7c+3vfOws
gAOpgqXXNHwpRSd1xazbTpu8Kz1nWeZ8w8aPmYKuo9+ucMbpzYpqmyiXb1DiHbxO
VsE3ZM6kBIyl7H5HABEBAAGJATwEGAEIACYWIQR3TcmcAGUBN1Ici7Pxx2Awu2yq
jQUCXXkUlQIbDAUJA8JnAAAKCRDxx2Awu2yqjQ4xCACRYNG/6JpKuOjsU/LSpw8G
rBNjFMlzNdiPOdHiW/gglBbMJB3LJJrM4TvMcFsqmuKUh1j7/gO9GUhm3VIRxZXx
mble0sEh5n6Tpz0HoZb2ndvi+tqbMm1ufDP9pbIXOZzdksywrAX3283vjDUTlDog
7qYBzQEG6TK68RGDKGobDtBIoR9S/enHoAkrWONKJ9uyJw2cIpx72MPXiMqP6vnL
Exdgp01NoEQx1UPfy/Y9gJ5aGaUUBDG7i6twpeTo9XFyJihrU5tFfrzT6iuGggxF
fJoCgxVAKzXJnGTulcClquAOmMCFKqxbkOTIUy0uATSGF4pIvGu0Edi0GzvfCKST
=Q8Ys
-----END PGP PUBLIC KEY BLOCK-----`

var manifest1 = `ewogICAic2NoZW1hVmVyc2lvbiI6IDIsCiAgICJtZWRpYVR5cGUiOiAiYXBwbGljYXRpb24vdm5kLmRvY2tlci5kaXN0cmlidXRpb24ubWFuaWZlc3QudjIranNvbiIsCiAgICJjb25maWciOiB7CiAgICAgICJtZWRpYVR5cGUiOiAiYXBwbGljYXRpb24vdm5kLmRvY2tlci5jb250YWluZXIuaW1hZ2UudjEranNvbiIsCiAgICAgICJzaXplIjogMTQ5NiwKICAgICAgImRpZ2VzdCI6ICJzaGEyNTY6NTk3ODhlZGYxZjNlNzhjZDBlYmU2Y2UxNDQ2ZTlkMTA3ODgyMjVkYjNkZWRjZmQxYTU5Zjc2NGJhZDJiMjY5MCIKICAgfSwKICAgImxheWVycyI6IFsKICAgICAgewogICAgICAgICAibWVkaWFUeXBlIjogImFwcGxpY2F0aW9uL3ZuZC5kb2NrZXIuaW1hZ2Uucm9vdGZzLmRpZmYudGFyLmd6aXAiLAogICAgICAgICAic2l6ZSI6IDcyNzk3OCwKICAgICAgICAgImRpZ2VzdCI6ICJzaGEyNTY6OTBlMDE5NTVlZGNkODVkYWM3OTg1YjcyYTgzNzQ1NDVlYWM2MTdjY2RkZGNjOTkyYjczMmU0M2NkNDI1MzRhZiIKICAgICAgfQogICBdCn0`
var signature1 = `owGbwMvMwMH48XiCwe6cVb2Mpw9EJTHEGQTdqFZKLsosyUxOzFGyqlbKTEnNK8ksqQSxU/KTs1OLdItS01KLUvOSU5WslDKTi/Qy8/WLS0oTi0r0k0qLK5PyK5RqdZQycxPTU5E05SbmZaalFpfopmSmAymg1uKMRCNTMytLQ9M0Y0uDRAtLQ6NUQ7MUk6TUJAszC0tzI4NEQ3NjE4s04zSzFMNE8yQzU0szS/OUNAtTg8QkMyPT1EQjyxRTkGUllQUgxySW5OdmJisk5+eVJGbmpRYpFGem5yWWlBalghTlF5Rk5udBfJVclApUXITQY6pnoGegBDQpMxfousTcAiUrQ1MLAyMjUyMTg9raTkZjFgZGDgZZMUWWct+TcxhSGc2DZLo3w0KPlQkUcgxcnAIwkTu57P8sz9Xc4/q0aXaqHt8d36iz5+5vKsk8v5WJ+aersNlrLp0ms/JDjbava/4wFMw4479zkYTXGuXNJvdfdmyK+G/m7X3y8/U0vZv3c96fdWYwfMJdebLzUsuxO1nld+p3rH29zvr2867e+3IqhVILGiOyO0V3xAa9Ovf+4hcNkzv+7wLkk7hzU/TYpsqe6j/bzvD6WqhCWLq0rb/LgYMO3244mny46fDsU3KLT43FrLCpMdG8+yRO5Wpkr5+q9EaM48/RdzyJsxe3fdVku2f1998GZkZen7WTJd9m1TPtPWy0jPf+0jUWixiLQwpKwmdYFZjXviu4EMkusHweg/6z6xN0dm1imfFjnSrn5Z9XBFcAAA`
var dockerReference1 = "icr.io/stuart/busybox"
var fingerprint1 = "774DC99C00650137521C8BB3F1C76030BB6CAA8D"

var keyFile = "/Users/hayton/Projects/rh-signing/public.gpg"

func TestVerify_GoodPath(t *testing.T) {

	block, err := armor.Decode(bytes.NewReader([]byte(publicKey1)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "public key decode err: %v\n", err)
	}
	if block.Type != "PGP PUBLIC KEY BLOCK" {
		fmt.Fprintf(os.Stderr, "not public key block, %s\n", block.Type)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(block.Body)
	keyBytes := buf.Bytes()

	// keyBytes, err := ioutil.ReadFile(keyFile)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "key read err: %v\n", err)
	// }

	manifestBytes, err := base64.RawStdEncoding.DecodeString(manifest1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "manifest decode err: %v\n", err)
	}
	sigBytes, err := base64.RawStdEncoding.DecodeString(signature1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "signature decode err: %v\n", err)
	}
	err = VerifyBySignature(keyBytes,
		manifestBytes,
		sigBytes,
		dockerReference1,
		fingerprint1)

	assert.Nil(t, err, "error")
}
