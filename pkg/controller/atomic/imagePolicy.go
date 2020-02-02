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

// Implementation of verify against image policy interface

package atomic

import (
	"context"
	"fmt"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
)

// VerifyByPolicy ...
func VerifyByPolicy(imageToVerify string, policyJSON string) error {
	imageReference, err := docker.ParseReference(`//` + imageToVerify)
	if err != nil {
		return err
	}
	systemContext := &types.SystemContext{}

	imageSource, err := imageReference.NewImageSource(context.Background(), systemContext)
	if err != nil {
		return err
	}
	unparsedImage := image.UnparsedInstance(imageSource, nil)

	policy, err := signature.NewPolicyFromBytes([]byte(policyJSON))
	if err != nil {
		return err
	}

	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		return err
	}

	allowed, err := policyContext.IsRunningImageAllowed(context.Background(), unparsedImage)
	if err != nil {
		return err
	}
	if !allowed {
		return fmt.Errorf("not allowed")
	}
	return nil
}
