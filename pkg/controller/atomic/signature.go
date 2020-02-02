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

// Implementation of verify against signature interface (doesnt seem to have policy flexibility)

package atomic

import (
	"fmt"
	"github.com/containers/image/v5/signature"
)

// VerifyBySignature ...
func VerifyBySignature(publicKey, unverifiedManifest, unverifiedSignature []byte, expectedDockerReference, expectedFingerPrint string) error {

	// without hosts public keys
	mech, _, err := signature.NewEphemeralGPGSigningMechanism(publicKey)
	if err != nil {
		return fmt.Errorf("Error initializing GPG: %v", err)
	}
	defer mech.Close()

	_, err = signature.VerifyDockerManifestSignature(unverifiedSignature, unverifiedManifest, expectedDockerReference, mech, expectedFingerPrint)
	if err != nil {
		return fmt.Errorf("Error verifying signature: %v", err)
	}
	return nil
}
