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
	"fmt"
	"os"

	"github.com/containers/image/v5/signature"
)

// InsecureAcceptAnything .
var InsecureAcceptAnything, _ = signature.NewPolicyFromBytes([]byte(`{ "default": [{"type": "insecureAcceptAnything"}] }`))

// Reject .
var Reject, _ = signature.NewPolicyFromBytes([]byte(`{ "default": [{"type": "reject"}] }`))

// NewPolicyFromString .
func NewPolicyFromString(policyString string) (*signature.Policy, error) {
	return signature.NewPolicyFromBytes([]byte(policyString))
}

// SafePolicyFromString .
func SafePolicyFromString(policyString string) *signature.Policy {
	policy, err := signature.NewPolicyFromBytes([]byte(policyString))
	if err != nil || policy == nil {
		fmt.Fprintf(os.Stderr, "Policy creation failed %v", err)
		os.Exit(1)
	}
	return policy
}
