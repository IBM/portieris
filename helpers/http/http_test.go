// Copyright 2018, 2022 Portieris Authors.
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

package http

import (
	"bufio"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAuthHeader(t *testing.T) {
	type expectations struct {
		challenges   []Challenge
		errorMessage string
	}
	tests := []struct {
		name   string
		in     []string
		expect expectations
	}{
		{
			name:   "No headers",
			in:     []string{},
			expect: expectations{},
		},
		{
			name: "Basic auth challenge",
			in:   []string{"WWW-Authenticate: Basic realm=\"protected area\""},
			expect: expectations{
				challenges: []Challenge{
					Challenge{
						Scheme: "basic",
						Parameters: map[string]string{
							"realm": "protected area",
						},
					},
				},
			},
		},
		{
			name: "Bearer challenge with quoted strings",
			in: []string{
				"Www-Authenticate: Bearer realm=\"https://private.notary.local/auth\",service=\"notary-server\",scope=\"repository:private.registry.local/image:pull\"",
			},
			expect: expectations{
				challenges: []Challenge{
					Challenge{
						Scheme: "bearer",
						Parameters: map[string]string{
							"realm":   "https://private.notary.local/auth",
							"service": "notary-server",
							"scope":   "repository:private.registry.local/image:pull",
						},
					},
				},
			},
		},
		{
			name: "Bearer challenge with some unquoted strings",
			in: []string{
				"WWW-Authenticate: Bearer realm=\"https://private.notary.local/auth\",service=notary,scope=\"repository:private.registry.local/image:pull\"",
			},
			expect: expectations{
				challenges: []Challenge{
					Challenge{
						Scheme: "bearer",
						Parameters: map[string]string{
							"realm":   "https://private.notary.local/auth",
							"service": "notary",
							"scope":   "repository:private.registry.local/image:pull",
						},
					},
				},
			},
		},
		{
			name: "Bearer challenge with parameters in different order",
			in: []string{
				"WWW-Authenticate: Bearer service=\"notary-server\",scope=\"repository:private.registry.local/image:pull\",realm=\"https://private.notary.local/auth\"",
			},
			expect: expectations{
				challenges: []Challenge{
					Challenge{
						Scheme: "bearer",
						Parameters: map[string]string{
							"realm":   "https://private.notary.local/auth",
							"service": "notary-server",
							"scope":   "repository:private.registry.local/image:pull",
						},
					},
				},
			},
		},
		{
			name: "More complex digest challenge (taken from RFC2617)",
			in: []string{
				"WWW-Authenticate: Digest realm=\"testrealm@host.com\", qop=\"auth,auth-int\", nonce=\"dcd98b7102dd2f0e8b11d0f600bfb0c093\", opaque=\"5ccc069c403ebaf9f0171e9517f40e41\"",
			},
			expect: expectations{
				challenges: []Challenge{
					Challenge{
						Scheme: "digest",
						Parameters: map[string]string{
							"realm":  "testrealm@host.com",
							"qop":    "auth,auth-int",
							"nonce":  "dcd98b7102dd2f0e8b11d0f600bfb0c093",
							"opaque": "5ccc069c403ebaf9f0171e9517f40e41",
						},
					},
				},
			},
		},
		{
			name: "Error - missing auth-scheme",
			in: []string{
				"WWW-Authenticate: ",
			},
			expect: expectations{
				errorMessage: "Unable to parse WWW-Authenticate header '': no auth-scheme found",
			},
		},
		{
			name: "Error - missing parameter value #1",
			in: []string{
				"WWW-Authenticate: Digest param",
			},
			expect: expectations{
				errorMessage: "Unable to parse WWW-Authenticate header 'Digest param': parameter value missing",
			},
		},
		{
			name: "Error - missing parameter value #2",
			in: []string{
				"WWW-Authenticate: Digest param=",
			},
			expect: expectations{
				errorMessage: "Unable to parse WWW-Authenticate header 'Digest param=': parameter value missing",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uut, err := toHTTPResponse(tt.in)

			if assert.NoError(t, err) {
				challenges, err := ParseAuthHeader(uut.Header)

				if tt.expect.errorMessage != "" {
					assert.EqualError(t, err, tt.expect.errorMessage)
					assert.Equal(t, 0, len(challenges))
				} else {
					assert.Equal(t, len(tt.expect.challenges), len(uut.Header))
					assert.Equal(t, len(tt.expect.challenges), len(challenges))

					for i, expectedChallenge := range tt.expect.challenges {
						assert.Equal(t, expectedChallenge.Scheme, challenges[i].Scheme)
						assert.Equal(t, expectedChallenge.Parameters, challenges[i].Parameters)
					}
				}
			}
		})
	}
}

func toHTTPResponse(httpHeaders []string) (*http.Response, error) {
	response := "HTTP/1.0 401 Unauthorized\n"

	for _, httpHeader := range httpHeaders {
		response += httpHeader + "\n"
	}

	response += "\n"

	return http.ReadResponse(bufio.NewReader(strings.NewReader(response)), nil)
}
