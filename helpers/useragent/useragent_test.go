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
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSet_RoundTrip(t *testing.T) {
	tests := map[string]struct {
		wantStatus int
		wantErr    bool
	}{
		"good path": {
			wantStatus: http.StatusTeapot,
		},
		"error path": {
			wantErr: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := &http.Client{
				Transport: &Set{
					Transport: http.DefaultTransport,
				},
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, r.Header.Get("User-Agent"), "portieris/undefined")
				w.WriteHeader(test.wantStatus)
			}))
			defer ts.Close()

			r, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			require.NoError(t, err)
			if test.wantErr {
				r, err = http.NewRequest(http.MethodGet, "htootyps://notaurl", nil)
				require.NoError(t, err)
			}

			res, err := c.Do(r)

			if (err != nil) != test.wantErr {
				t.Errorf("error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !test.wantErr {
				assert.Equal(t, res.StatusCode, test.wantStatus)
			}
		})
	}
}
