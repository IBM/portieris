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

// This file contains HTTP transport functionality derived from
// github.com/distribution/distribution/registry/client/transport
// which was removed in distribution v3.
//
// Original code: Copyright Docker, Inc.
// Original license: Apache License 2.0
// Source: https://github.com/distribution/distribution/blob/v2.8.3/registry/client/transport/transport.go
//
// Modifications made by Portieris Authors:
// - Adapted to work as a standalone implementation without distribution package dependency
// - Simplified to focus on header modification use case
// - Removed deprecated CancelRequest support and request tracking (unused in this codebase)

package notary

import (
	"net/http"
)

// headerTransport is a custom RoundTripper that adds headers to requests.
type headerTransport struct {
	base    http.RoundTripper
	headers http.Header
}

// RoundTrip implements the http.RoundTripper interface.
// It clones the request, adds headers, and executes the request.
func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	req2 := t.cloneRequest(req)

	// Add custom headers
	for key, values := range t.headers {
		req2.Header[key] = append(req2.Header[key], values...)
	}

	// Execute the request
	return t.base.RoundTrip(req2)
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func (t *headerTransport) cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header to avoid race conditions
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}

// Made with Bob
