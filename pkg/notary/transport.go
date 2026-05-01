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
// - Maintains all original functionality: request tracking, cancellation support, and proper cleanup

package notary

import (
	"io"
	"net/http"
	"sync"
)

// headerTransport is a custom RoundTripper that adds headers to requests.
// It maintains request tracking for proper cancellation support and cleanup of resources.
type headerTransport struct {
	base    http.RoundTripper
	headers http.Header
	mu      sync.Mutex                      // guards modReq
	modReq  map[*http.Request]*http.Request // original -> modified
}

// RoundTrip implements the http.RoundTripper interface.
// It clones the request, adds headers, tracks the request mapping for cancellation,
// and wraps the response body to clean up the mapping when done.
func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	req2 := t.cloneRequest(req)

	// Add custom headers
	for key, values := range t.headers {
		req2.Header[key] = append(req2.Header[key], values...)
	}

	// Track the modified request for cancellation support
	t.setModReq(req, req2)

	// Execute the request
	res, err := t.base.RoundTrip(req2)
	if err != nil {
		t.setModReq(req, nil)
		return nil, err
	}

	// Wrap the response body to clean up the request mapping when done
	res.Body = &onEOFReader{
		rc: res.Body,
		fn: func() { t.setModReq(req, nil) },
	}

	return res, nil
}

// CancelRequest cancels an in-flight request by closing its connection.
// This is deprecated in favor of context cancellation, but maintained for compatibility.
func (t *headerTransport) CancelRequest(req *http.Request) {
	type canceler interface {
		CancelRequest(*http.Request)
	}
	if cr, ok := t.base.(canceler); ok {
		t.mu.Lock()
		modReq := t.modReq[req]
		delete(t.modReq, req)
		t.mu.Unlock()
		if modReq != nil {
			cr.CancelRequest(modReq)
		}
	}
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

// setModReq tracks or removes the mapping between original and modified requests
func (t *headerTransport) setModReq(orig, mod *http.Request) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.modReq == nil {
		t.modReq = make(map[*http.Request]*http.Request)
	}
	if mod == nil {
		delete(t.modReq, orig)
	} else {
		t.modReq[orig] = mod
	}
}

// onEOFReader wraps a ReadCloser to call a function when EOF is reached or Close is called.
// This ensures proper cleanup of request tracking.
type onEOFReader struct {
	rc io.ReadCloser
	fn func()
}

func (r *onEOFReader) Read(p []byte) (n int, err error) {
	n, err = r.rc.Read(p)
	if err == io.EOF {
		r.runFunc()
	}
	return
}

func (r *onEOFReader) Close() error {
	err := r.rc.Close()
	r.runFunc()
	return err
}

func (r *onEOFReader) runFunc() {
	if fn := r.fn; fn != nil {
		fn()
		r.fn = nil
	}
}
