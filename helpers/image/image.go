// Copyright 2018 IBM
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
	"fmt"
	"net/url"
	"strings"

	"github.com/docker/distribution/reference"
)

// ImageReference .
type ImageReference struct {
	original string
	name     string
	tag      string
	digest   string
	hostname string
	port     string
}

// NewImageReference parses the image name and returns an error if the name is invalid.
func NewImageReference(name string) (*ImageReference, error) {
	var digest string
	original := name
	// Remove the digest so `ParseNamed` doesn't fail, it can't handle short digests.
	if strings.Contains(name, "@sha256:") {
		fields := strings.Split(name, "@sha256:")
		name = fields[0]
		digest = fields[1]
	}
	// Get image name
	ref, err := reference.ParseNamed(name)
	if err != nil {
		return nil, err
	}

	// Get the hostname
	hostname, _ := reference.SplitHostname(ref)
	if hostname == "" {
		// If no domain found, treat it as docker.io
		hostname = "docker.io"
	}
	if !strings.Contains(hostname, ".") {
		// Fix SplitHostname wrongly splitting repositories like molepigeon/wibble
		hostname = "docker.io"
	}
	// Make sure it can be used to build a valid URL
	u, err := url.Parse("http://" + hostname)
	if err != nil {
		return nil, err
	}

	// if the image does not have a tag, use `latest` so we can parse it again.
	image := strings.Replace(name, hostname, "", 1)
	if !strings.Contains(image, ":") {
		name += ":latest"
	}

	// Parse the name again including the tag so we can have a reference.taggedReference object
	// we ommit the error here since we already parsed the original string above.
	ref, _ = reference.ParseNamed(name)

	return &ImageReference{
		original: original,
		name:     ref.Name(),
		tag:      ref.(reference.Tagged).Tag(),
		digest:   digest,
		hostname: u.Hostname(),
		port:     u.Port(),
	}, nil
}

// GetHostname returns the repository hostname of an image
func (i ImageReference) GetHostname() string {
	return i.hostname
}

// GetPort returns the port of the hostname.
func (i ImageReference) GetPort() string {
	return i.port
}

// HasIBMRepo returns true if the image has an IBM repository, otherwise false.
func (i ImageReference) HasIBMRepo() bool {
	prefix := "registry"
	suffix := ".bluemix.net"
	if !strings.HasPrefix(i.hostname, prefix) || !strings.HasSuffix(i.hostname, suffix) {
		return false
	}
	return true
}

// GetRegistryURL returns the Registry URL.
func (i ImageReference) GetRegistryURL() string {
	port := i.port
	if port != "" {
		port = ":" + port
	}
	return "https://" + i.hostname + port
}

// GetContentTrustURL returns the Content Trust URL.
func (i ImageReference) GetContentTrustURL() string {
	// TODO: Add support for notaries from other repos other than IBM
	return "https://" + i.hostname + ":4443"
}

// GetVAURL returns the Vulnerability Advisor URL.
func (i ImageReference) GetVAURL() (string, error) {
	if !i.HasIBMRepo() {
		return "", fmt.Errorf("Deny %q, Vulnerability Advisor is not supported for images from this registry: %s", i.String(), i.GetHostname())
	}
	vaURL := strings.Replace(i.hostname, "registry", "va", 1)
	return "https://" + vaURL, nil
}

// GetTag returns the tag.
func (i ImageReference) GetTag() string {
	return i.tag
}

// GetDigest returns the digest.
func (i ImageReference) GetDigest() string {
	return i.digest
}

// NameWithTag returns the image name with the tag.
func (i ImageReference) NameWithTag() string {
	return i.name + ":" + i.tag
}

// NameWithoutTag returns the image name without the tag.
func (i ImageReference) NameWithoutTag() string {
	return i.name
}

// String returns the original image name.
func (i ImageReference) String() string {
	return i.original
}
