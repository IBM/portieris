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
	"context"
	"fmt"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
)

// VerifyByPolicy returns the verified digest or the verification error
func VerifyByPolicy(imageToVerify string, policy *signature.Policy, username, password string) (string, error) {

	dockerAuth := &types.DockerAuthConfig{
		Username: username,
		Password: password,
	}
	systemContext := &types.SystemContext{
		RootForImplicitAbsolutePaths: "/nowhere", // read nothing from files
		DockerAuthConfig:             dockerAuth,
		DockerRegistryUserAgent:      "portieris",
	}

	imageReference, err := docker.ParseReference(`//` + imageToVerify)
	if err != nil {
		return "", err
	}
	imageSource, err := imageReference.NewImageSource(context.Background(), systemContext)
	if err != nil {
		return "", err
	}
	unparsedImage := image.UnparsedInstance(imageSource, nil)

	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		return "", err
	}

	allowed, err := policyContext.IsRunningImageAllowed(context.Background(), unparsedImage)
	if err != nil {
		return "", err
	}
	// redundant?
	if !allowed {
		return "", fmt.Errorf("not allowed")
	}

	// get the digest
	m, _, err := unparsedImage.Manifest(context.Background())
	digest, err := manifest.Digest(m)
	if err != nil {
		return "", err
	}
	return digest.String(), nil
}
