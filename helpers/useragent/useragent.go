// Copyright 2021  Portieris Authors.
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

package useragent

import (
	"net/http"

	"github.com/IBM/portieris/internal/info"
)

// Set is a http.RoundTripper which adds the User-Agent header to all requests.
type Set struct {
	Transport http.RoundTripper
}

// RoundTrip sets the User-Agent on the request and then calls the underlying
// Transport.
func (a *Set) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", "portieris/"+info.Version)
	return a.Transport.RoundTrip(r)
}
