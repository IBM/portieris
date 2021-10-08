// Copyright 2018, 2021 Portieris Authors.
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

package trust

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"path"

	"github.com/golang/glog"
	"github.com/theupdateframework/notary/tuf/data"
	"github.com/theupdateframework/notary/tuf/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var releasesRole = data.RoleName(path.Join(data.CanonicalTargetsRole.String(), "releases"))

// Signer struct holds the signer and publicKey from a SignerSecret
type Signer struct {
	signer    string
	publicKey string
}

type foundSigner struct {
	found  bool
	signer Signer
}

// getDigest .
func (v *Verifier) getDigest(server, image, notaryToken, targetName string, signers []Signer) (*bytes.Buffer, error) {
	repo, err := v.trust.GetNotaryRepo(server, image, notaryToken)
	if err != nil {
		return nil, err
	}

	roleNames := make([]string, len(signers))
	for i, val := range signers {
		roleNames[i] = path.Join(data.CanonicalTargetsRole.String(), val.signer)
	}
	rolelist := data.NewRoleList(roleNames)

	var foundSignerByRole = map[data.RoleName]*foundSigner{}
	for i, role := range rolelist {
		foundSignerByRole[role] = &foundSigner{
			signer: signers[i],
		}
	}

	targets, err := repo.GetAllTargetMetadataByName(targetName)
	if err != nil {
		glog.Infof("GetAllTargetMetadataByName returned err: %+v", err)
		return nil, err
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("No signed targets found")
	}

	var digest []byte // holds digest of the signed image

	// Get the digest of the latest signed release
	// A digest is "released" if it's signed by the targets or targets/releases roles
	for _, target := range targets {
		if target.Role.Name == data.CanonicalTargetsRole || target.Role.Name == releasesRole {
			digest = target.Target.Hashes["sha256"]
		}
	}

	if len(rolelist) == 0 {
		glog.Infof("roleList length == 0, returning digest %s", hex.EncodeToString(digest))
	} else {
		for _, target := range targets { // iterate over each target
			// See if a signer was specified for this target
			if role, ok := foundSignerByRole[target.Role.Name]; ok {
				if role.signer.publicKey != "" {
					// Assuming public key is in PEM format and not encoded any further
					keyFromConfig, err := utils.ParsePEMPublicKey([]byte(role.signer.publicKey))
					if err != nil {
						return nil, err
					}
					if _, ok := target.Role.BaseRole.Keys[keyFromConfig.ID()]; !ok {
						glog.Infof("Key %s not found in role key list: %+v", keyFromConfig.ID(), target.Role.BaseRole.ListKeyIDs())
						return nil, fmt.Errorf("Public keys are different")
					}
					// We found a matching KeyID, so mark the role found in the map.
					role.found = true
				} else {
					glog.Infof("PublicKey not found in role %s", role.signer.signer)
					return nil, fmt.Errorf("PublicKey not found in role %s", role.signer.signer)
				}

				// verify that the digest is consistent between all of the roles that we care about
				if !bytes.Equal(digest, target.Target.Hashes["sha256"]) {
					return nil, fmt.Errorf("Incompatible digest")
				}
			}
		}

		// Now iterate over the signers to make sure we hit them all going over targets
		for _, role := range foundSignerByRole {
			if !role.found {
				return nil, fmt.Errorf("no signature found for role %s", role.signer.signer)
			}
		}
	}

	return bytes.NewBufferString(hex.EncodeToString(digest)), nil
}

// Retrieve the username and public key for the given namespace/secret
func (v *Verifier) getSignerSecret(namespace, signerSecretName string) (Signer, error) {

	// Retrieve secret
	secret, err := v.kubeClientsetWrapper.CoreV1().Secrets(namespace).Get(context.TODO(), signerSecretName, metav1.GetOptions{})
	if err != nil {
		glog.Error("Error: ", err)
		return Signer{}, err
	}

	// Parse the returned data.
	signer := string(secret.Data["name"])
	publicKey := string(secret.Data["publicKey"])

	if signer == "" || publicKey == "" {
		return Signer{}, fmt.Errorf("name or publicKey field in secret %s is empty", signerSecretName)
	}

	return Signer{signer: signer, publicKey: publicKey}, nil
}
