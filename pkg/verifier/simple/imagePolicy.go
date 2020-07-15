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

// Implementation of verify against containers/image policy interface

package simple

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
)

// VerifyByPolicy verifies the image according to the supplied policy and returns the verified digest, verify error or processing error
func VerifyByPolicy(imageToVerify string, credentials [][]string, registriesConfigDir string, simplePolicy *signature.Policy) (*bytes.Buffer, error, error) {

	policyContext, err := signature.NewPolicyContext(simplePolicy)
	if err != nil {
		return nil, nil, err
	}
	imageReference, err := docker.ParseReference(`//` + imageToVerify)
	if err != nil {
		return nil, nil, err
	}
	// if expensive, make instance
	systemContext := &types.SystemContext{
		RootForImplicitAbsolutePaths: "/nowhere",  // read nothing from files
		DockerRegistryUserAgent:      "portieris", // add version?
		RegistriesDirPath:            registriesConfigDir,
	}

	switch len(credentials) {
	case 0:
		return verifyAttempt(imageReference, systemContext, policyContext)

	case 1:
		systemContext.DockerAuthConfig = &types.DockerAuthConfig{
			Username: credentials[0][0],
			Password: credentials[0][1],
		}
		return verifyAttempt(imageReference, systemContext, policyContext)
	}

	for _, credential := range credentials {
		systemContext.DockerAuthConfig = &types.DockerAuthConfig{
			Username: credential[0],
			Password: credential[1],
		}

		digest, deny, err := verifyAttempt(imageReference, systemContext, policyContext)
		if deny != nil {
			switch deny.(type) {
			case *docker.ErrUnauthorizedForCredentials:
				continue
			}
		}
		return digest, deny, err
	}

	return nil, nil, fmt.Errorf("Deny %q, no valid ImagePullSecret, %d tried", imageToVerify, len(credentials))
}

func verifyAttempt(imageReference types.ImageReference, systemContext *types.SystemContext, policyContext *signature.PolicyContext) (*bytes.Buffer, error, error) {
	imageSource, err := imageReference.NewImageSource(context.Background(), systemContext)
	if err != nil {
		return nil, nil, err
	}

	unparsedImage := image.UnparsedInstance(imageSource, nil)
	_, deny := policyContext.IsRunningImageAllowed(context.Background(), unparsedImage)
	if deny != nil {
		return nil, deny, nil
	}

	// get the digest
	m, _, err := unparsedImage.Manifest(context.Background())
	digest, err := manifest.Digest(m)
	if err != nil {
		return nil, nil, err
	}
	return bytes.NewBufferString(strings.TrimPrefix(digest.String(), "sha256:")), nil, nil
}
