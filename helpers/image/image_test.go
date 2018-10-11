// Copyright 2018 Portieris Authors.
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

package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReference(t *testing.T) {
	type expectations struct {
		ReferenceError  bool
		Hostname        string
		Port            string
		HasIBMRepo      bool
		RegistryURL     string
		ContentTrustURL string
		ContentTrustErr bool
		Tag             string
		Digest          string
		NameWithTag     string
		NameWithoutTag  string
		String          string
	}
	tests := []struct {
		name   string
		in     string
		expect expectations
	}{
		{
			name: "errors on invalid input",
			in:   "?",
			expect: expectations{
				ReferenceError: true,
			},
		},
		{
			name: "parses an image without a tag",
			in:   "test.com/namespace/name",
			expect: expectations{
				Hostname:        "test.com",
				HasIBMRepo:      false,
				Port:            "",
				Tag:             "latest",
				Digest:          "",
				NameWithTag:     "test.com/namespace/name:latest",
				NameWithoutTag:  "test.com/namespace/name",
				String:          "test.com/namespace/name",
				RegistryURL:     "https://test.com",
				ContentTrustErr: true,
			},
		},
		{
			name: "parses an image with a digest",
			in:   "test.com:8080/namespace/name@sha256:1234567890",
			expect: expectations{
				Hostname:        "test.com",
				HasIBMRepo:      false,
				Port:            "8080",
				Tag:             "latest",
				Digest:          "1234567890",
				NameWithTag:     "test.com:8080/namespace/name:latest",
				NameWithoutTag:  "test.com:8080/namespace/name",
				String:          "test.com:8080/namespace/name@sha256:1234567890",
				RegistryURL:     "https://test.com:8080",
				ContentTrustErr: true,
			},
		},
		{
			name: "parses an image with a tag",
			in:   "test.com/namespace/name:v1",
			expect: expectations{
				Hostname:        "test.com",
				HasIBMRepo:      false,
				Port:            "",
				Tag:             "v1",
				Digest:          "",
				NameWithTag:     "test.com/namespace/name:v1",
				NameWithoutTag:  "test.com/namespace/name",
				String:          "test.com/namespace/name:v1",
				RegistryURL:     "https://test.com",
				ContentTrustErr: true,
			},
		},
		{
			name: "parses an image with a tag and a digest",
			in:   "test.com:8080/namespace/name:v1@sha256:1234567890",
			expect: expectations{
				Hostname:        "test.com",
				HasIBMRepo:      false,
				Port:            "8080",
				Tag:             "v1",
				Digest:          "1234567890",
				NameWithTag:     "test.com:8080/namespace/name:v1",
				NameWithoutTag:  "test.com:8080/namespace/name",
				String:          "test.com:8080/namespace/name:v1@sha256:1234567890",
				RegistryURL:     "https://test.com:8080",
				ContentTrustErr: true,
			},
		},
		{
			name: "parses an image from Docker Hub with a tag and a digest",
			in:   "namespace/name:v1@sha256:1234567890",
			expect: expectations{
				Hostname:        "docker.io",
				HasIBMRepo:      false,
				Port:            "",
				Tag:             "v1",
				Digest:          "1234567890",
				NameWithTag:     "namespace/name:v1",
				NameWithoutTag:  "namespace/name",
				String:          "namespace/name:v1@sha256:1234567890",
				RegistryURL:     "https://docker.io",
				ContentTrustErr: false,
				ContentTrustURL: "https://notary.docker.io",
			},
		},
		{
			name: "parses a Docker Hub public image with a tag and a digest",
			in:   "ubuntu:v1@sha256:1234567890",
			expect: expectations{
				Hostname:        "docker.io",
				HasIBMRepo:      false,
				Port:            "",
				Tag:             "v1",
				Digest:          "1234567890",
				NameWithTag:     "ubuntu:v1",
				NameWithoutTag:  "ubuntu",
				String:          "ubuntu:v1@sha256:1234567890",
				RegistryURL:     "https://docker.io",
				ContentTrustErr: false,
				ContentTrustURL: "https://notary.docker.io",
			},
		},
		{
			name: "parses an IBM image",
			in:   "registry.ng.bluemix.net/namespace/name",
			expect: expectations{
				Hostname:        "registry.ng.bluemix.net",
				HasIBMRepo:      true,
				Port:            "",
				Tag:             "latest",
				Digest:          "",
				NameWithTag:     "registry.ng.bluemix.net/namespace/name:latest",
				NameWithoutTag:  "registry.ng.bluemix.net/namespace/name",
				String:          "registry.ng.bluemix.net/namespace/name",
				RegistryURL:     "https://registry.ng.bluemix.net",
				ContentTrustErr: false,
				ContentTrustURL: "https://registry.ng.bluemix.net:4443",
			},
		},
		{
			name: "parses a quay.io image",
			in:   "quay.io/namespace/name",
			expect: expectations{
				Hostname:        "quay.io",
				HasIBMRepo:      false,
				Port:            "",
				Tag:             "latest",
				Digest:          "",
				NameWithTag:     "quay.io/namespace/name:latest",
				NameWithoutTag:  "quay.io/namespace/name",
				String:          "quay.io/namespace/name",
				RegistryURL:     "https://quay.io",
				ContentTrustErr: false,
				ContentTrustURL: "https://quay.io:443",
			},
		},
		{
			name: "parses an ICR image",
			in:   "us.icr.io/namespace/name",
			expect: expectations{
				Hostname:        "us.icr.io",
				HasIBMRepo:      true,
				Port:            "",
				Tag:             "latest",
				Digest:          "",
				NameWithTag:     "us.icr.io/namespace/name:latest",
				NameWithoutTag:  "us.icr.io/namespace/name",
				String:          "us.icr.io/namespace/name",
				RegistryURL:     "https://us.icr.io",
				ContentTrustErr: false,
				ContentTrustURL: "https://us.icr.io:4443",
			},
		},
		{
			name: "parses a staging ICR image",
			in:   "stg.icr.io/namespace/name",
			expect: expectations{
				Hostname:        "stg.icr.io",
				HasIBMRepo:      true,
				Port:            "",
				Tag:             "latest",
				Digest:          "",
				NameWithTag:     "stg.icr.io/namespace/name:latest",
				NameWithoutTag:  "stg.icr.io/namespace/name",
				String:          "stg.icr.io/namespace/name",
				RegistryURL:     "https://stg.icr.io",
				ContentTrustErr: false,
				ContentTrustURL: "https://stg.icr.io:4443",
			},
		},
		{
			name: "parses an ICR image with a port",
			in:   "de.icr.io:8080/namespace/name",
			expect: expectations{
				Hostname:        "de.icr.io",
				HasIBMRepo:      true,
				Port:            "8080",
				Tag:             "latest",
				Digest:          "",
				NameWithTag:     "de.icr.io:8080/namespace/name:latest",
				NameWithoutTag:  "de.icr.io:8080/namespace/name",
				String:          "de.icr.io:8080/namespace/name",
				RegistryURL:     "https://de.icr.io:8080",
				ContentTrustErr: false,
				ContentTrustURL: "https://de.icr.io:4443",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image, err := NewReference(tt.in)
			if tt.expect.ReferenceError {
				assert.Error(t, err)
				return
			}

			if assert.NoError(t, err) {
				assert.Equal(t, tt.expect.Hostname, image.GetHostname(), "Hostname")
				assert.Equal(t, tt.expect.Port, image.GetPort(), "Port")
				assert.Equal(t, tt.expect.HasIBMRepo, image.HasIBMRepo(), "HasIBMRepo")
				assert.Equal(t, tt.expect.RegistryURL, image.GetRegistryURL(), "GetRegistryURL")
				trustURL, trustErr := image.GetContentTrustURL()
				if tt.expect.ContentTrustErr {
					assert.Error(t, trustErr, "GetContentTrust err")
				} else {
					assert.NoError(t, trustErr, "GetContentTrust err")
					assert.Equal(t, tt.expect.ContentTrustURL, trustURL, "GetContentTrust")
				}
				assert.Equal(t, tt.expect.Tag, image.GetTag(), "Tag")
				assert.Equal(t, tt.expect.Digest, image.GetDigest(), "Digest")
				assert.Equal(t, tt.expect.NameWithTag, image.NameWithTag(), "NameWithTag")
				assert.Equal(t, tt.expect.NameWithoutTag, image.NameWithoutTag(), "NameWithoutTag")
				assert.Equal(t, tt.expect.String, image.String(), "String")
			}
		})
	}
}
