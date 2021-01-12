// Copyright 2020, 2021 Portieris Authors.
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

	"github.com/IBM/portieris/helpers/credential"
	"github.com/containers/common/pkg/retry"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
	"github.com/golang/glog"
)

var retryOpts = &retry.RetryOptions{
	MaxRetry: 2,
	// No Delay because the containers/image library internally retries 429
	// errors, all other retryable errors are suitable for immediate retry.
}

// VerifyByPolicy verifies the image according to the supplied policy and returns the verified digest, verify error or processing error
func (v verifier) VerifyByPolicy(imageToVerify string, credentials credential.Credentials, registriesConfigDir string, simplePolicy *signature.Policy) (*bytes.Buffer, error, error) {

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

	imageSource, err := newImageSourceWithRetry(context.Background(), systemContext, imageReference)
	if err == nil {
		defer imageSource.Close()
		glog.Infof("SimpleSigning verification: anonymous access allowed for image %s, continuing with anonymous verify", imageToVerify)
		return verifyWithRetry(context.Background(), imageSource, policyContext)
	}
	glog.Errorf("SimpleSigning verification: anonymous access denied for image %s, continuing with ImagePullSecrets... Error %v", imageToVerify, err)

	numCreds := len(credentials)
	for i, credential := range credentials {
		dockerAuthConfig := &types.DockerAuthConfig{
			Username: credential.Username,
			Password: credential.Password,
		}
		systemContext.DockerAuthConfig = dockerAuthConfig
		imageSource, err := newImageSourceWithRetry(context.Background(), systemContext, imageReference)
		if err != nil {
			if i+1 == numCreds {
				glog.Errorf("SimpleSigning verification: ImagePullSecret with username %s for image %s failed, no more secrets in scope (secret %d/%d). Failing. Error %v", credential.Username, imageToVerify, i+1, numCreds, err)
				return nil, nil, err
			}
			glog.Warningf("SimpleSigning verification: ImagePullSecret with username %s for image %s failed, trying the next secret in scope (secret %d/%d). Error %v", credential.Username, imageToVerify, i+1, numCreds, err)
			continue
		}
		defer imageSource.Close()
		glog.Infof("SimpleSigning verification: ImagePullSecret with username %s for image %s was valid (secret %d/%d), continuing to next stage", credential.Username, imageToVerify, i+1, numCreds)
		return verifyWithRetry(context.Background(), imageSource, policyContext)
	}

	return nil, nil, fmt.Errorf("Deny %q, no valid ImagePullSecret, %d tried", imageToVerify, len(credentials))
}

func verifyWithRetry(ctx context.Context, imageSource types.ImageSource, policyContext *signature.PolicyContext) (*bytes.Buffer, error, error) {
	var (
		manifest *bytes.Buffer
		deny     error
		attempt  int
	)
	return manifest, deny, retry.RetryIfNecessary(ctx, func() error {
		attempt++
		glog.InfoDepth(6, fmt.Sprintf("Attempt %d to verify image", attempt))
		var err error
		manifest, deny, err = verifyAttempt(ctx, imageSource, policyContext)
		if err != nil {
			glog.ErrorDepth(5, fmt.Sprintf("Attempt %d to verify image failed: %s", attempt, err))
		}
		return err
	}, retryOpts)
}

func verifyAttempt(ctx context.Context, imageSource types.ImageSource, policyContext *signature.PolicyContext) (*bytes.Buffer, error, error) {
	unparsedImage := image.UnparsedInstance(imageSource, nil)
	_, err := policyContext.IsRunningImageAllowed(ctx, unparsedImage)
	switch err.(type) {
	case nil:
	case signature.PolicyRequirementError:
		return nil, err, nil
	default:
		return nil, nil, err
	}

	// get the digest
	m, _, err := unparsedImage.Manifest(ctx)
	if err != nil {
		return nil, nil, err
	}
	digest, err := manifest.Digest(m)
	if err != nil {
		return nil, nil, err
	}
	return bytes.NewBufferString(strings.TrimPrefix(digest.String(), "sha256:")), nil, nil
}

func newImageSourceWithRetry(ctx context.Context, systemContext *types.SystemContext, imageReference types.ImageReference) (types.ImageSource, error) {
	var (
		imageSource types.ImageSource
		attempt     int
	)
	return imageSource, retry.RetryIfNecessary(ctx, func() error {
		attempt++
		glog.InfoDepth(6, fmt.Sprintf("Attempt %d to get image source", attempt))
		var err error
		imageSource, err = imageReference.NewImageSource(context.Background(), systemContext)
		if err != nil {
			glog.ErrorDepth(5, fmt.Sprintf("Attempt %d to get image source for simple verification failed: %s", attempt, err))
		}
		return err
	}, retryOpts)
}
