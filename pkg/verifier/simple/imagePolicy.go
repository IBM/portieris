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

	"github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
)

// VerifyByPolicy1 ...
func VerifyByPolicy1(imageToVerify string, credentials [][]string, portierisPolicy *v1beta1.Policy) (*bytes.Buffer, error, error) {

	simplePolicy, err := transformPolicy(portierisPolicy.Simple)
	if err != nil {
		return nil, nil, err
	}
	policyContext, err := signature.NewPolicyContext(simplePolicy)
	if err != nil {
		return nil, nil, err
	}
	imageReference, err := docker.ParseReference(`//` + imageToVerify)
	if err != nil {
		return nil, nil, err
	}
	systemContext := &types.SystemContext{
		RootForImplicitAbsolutePaths: "/nowhere", // read nothing from files
		DockerRegistryUserAgent:      "portieris",
	}

	for _, cred := range credentials {
		dockerAuthConfig := &types.DockerAuthConfig{
			Username: cred[0],
			Password: cred[1],
		}
		systemContext.DockerAuthConfig = dockerAuthConfig
		imageSource, err := imageReference.NewImageSource(context.Background(), systemContext)
		if err != nil {
			return nil, nil, err
		}
		unparsedImage := image.UnparsedInstance(imageSource, nil)
		allowed, deny := policyContext.IsRunningImageAllowed(context.Background(), unparsedImage)
		if err != nil {
			switch err.(type) {
			case *docker.ErrUnauthorizedForCredentials:
				continue
			default:
				return nil, deny, nil
			}
		}
		// redundant?
		if !allowed {
			return nil, fmt.Errorf("not allowed"), nil
		}
		// get the digest
		m, _, err := unparsedImage.Manifest(context.Background())
		digest, err := manifest.Digest(m)
		if err != nil {
			return nil, nil, err
		}
		return bytes.NewBufferString(digest.String()), nil, nil
	}
	return nil, nil, fmt.Errorf("Deny %q, no valid ImagePullSecret defined for %s", imageToVerify, imageToVerify)
}
