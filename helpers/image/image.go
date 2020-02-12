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
	"fmt"
	"net/url"
	"strings"

	"github.com/IBM/portieris/helpers/trustmap"
	"github.com/docker/distribution/reference"
)

// Reference .
type Reference struct {
	original string
	name     string
	tag      string
	digest   string
	hostname string
	port     string
	repo     string
}

// NewReference parses the image name and returns an error if the name is invalid.
func NewReference(name string) (*Reference, error) {
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
	hostname, repoName := reference.SplitHostname(ref)
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

	return &Reference{
		original: original,
		name:     ref.Name(),
		tag:      ref.(reference.Tagged).Tag(),
		digest:   digest,
		hostname: u.Hostname(),
		port:     u.Port(),
		repo:     repoName,
	}, nil
}

// GetHostname returns the repository hostname of an image
func (r Reference) GetHostname() string {
	return r.hostname
}

// GetPort returns the port of the hostname.
func (r Reference) GetPort() string {
	return r.port
}

// HasIBMRepo returns true if the image has an IBM repository, otherwise false.
func (r Reference) HasIBMRepo() bool {
	if strings.HasPrefix(r.hostname, "registry") && strings.HasSuffix(r.hostname, ".bluemix.net") {
		return true
	}
	if strings.HasSuffix(r.hostname, "icr.io") {
		return true
	}
	return false
}

// GetRegistryURL returns the Registry URL.
func (r Reference) GetRegistryURL() string {
	port := r.port
	if port != "" {
		port = ":" + port
	}
	return "https://" + r.hostname + port
}

// GetContentTrustURL returns the Content Trust URL.
func (r Reference) GetContentTrustURL() (string, error) {
	var output string
	var err error
	for registry, trustServerFn := range trustmap.TrustServerMap {
		if strings.HasSuffix(r.hostname, registry) {
			output = trustServerFn(registry, r.hostname)
		}
	}
	if output == "" {
		err = fmt.Errorf("no trust server could be found")
	}
	return output, err
}

// GetTag returns the tag.
func (r Reference) GetTag() string {
	return r.tag
}

// GetDigest returns the digest.
func (r Reference) GetDigest() string {
	return r.digest
}

// NameWithTag returns the image name with the tag.
func (r Reference) NameWithTag() string {
	return r.name + ":" + r.tag
}

// NameWithoutTag returns the image name without the tag.
func (r Reference) NameWithoutTag() string {
	return r.name
}

// RepoName returns the repo name
func (r Reference) RepoName() string {
	return r.repo
}

// String returns the original image name.
func (r Reference) String() string {
	return r.original
}
