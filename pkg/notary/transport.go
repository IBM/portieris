// Copyright 2018, 2026 Portieris Authors.
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

package notary

import (
	"net/http"
)

// headerTransport is a simple RoundTripper that adds headers to HTTP requests.
// It uses the standard library's req.Clone() for safe request copying.
type headerTransport struct {
	base    http.RoundTripper
	headers http.Header
}

// RoundTrip implements the http.RoundTripper interface.
// It clones the request and adds the configured headers before forwarding to the base transport.
func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	req2 := req.Clone(req.Context())

	// Add custom headers
	for key, values := range t.headers {
		req2.Header[key] = values
	}

	return t.base.RoundTrip(req2)
}
