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

// Implementation of Verify against openpgp directly without use of containers/image
// library, but copy a lot of it for flexibility but more work
// will need to import (or copy) some defs of signature at least
// incomplete

package simple

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/openpgp"
)

// InvalidSignatureError ...
type InvalidSignatureError struct {
	msg string
}

func (err InvalidSignatureError) Error() string {
	return err.msg
}

// KeyRingFromBytes ...
func KeyRingFromBytes(inKeyRing []byte) (openpgp.EntityList, error) {
	outKeyRing := openpgp.EntityList{}

	keyRing, err := openpgp.ReadKeyRing(bytes.NewReader(inKeyRing))
	if err != nil {
		k, e2 := openpgp.ReadArmoredKeyRing(bytes.NewReader(inKeyRing))
		if e2 != nil {
			return nil, err // The original error  -- FIXME: is this better?
		}
		keyRing = k
	}

	for _, entity := range keyRing {
		if entity.PrimaryKey == nil {
			// Coverage: This should never happen, openpgp.ReadEntity fails with a
			// openpgp.errors.StructuralError instead of returning an entity with this
			// field set to nil.
			continue
		}
		keyRing = append(outKeyRing, entity)
	}
	return outKeyRing, nil
}

// VerifySignature ...
func VerifySignature(signatureReader io.Reader, manifest, publicKey []byte) error {

	keyRing, err := KeyRingFromBytes(publicKey)

	md, err := openpgp.ReadMessage(signatureReader, keyRing, nil, nil)
	if err != nil {
		return err
	}
	if !md.IsSigned {
		return errors.New("not signed")
	}
	_, err = ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		// Coverage: md.UnverifiedBody.Read only fails if the body is encrypted
		// (and possibly also signed, but it _must_ be encrypted) and the signing
		// “modification detection code” detects a mismatch. But in that case,
		// we would expect the signature verification to fail as well, and that is checked
		// first.  Besides, we are not supplying any decryption keys, so we really
		// can never reach this “encrypted data MDC mismatch” path.
		return err
	}
	if md.SignatureError != nil {
		return fmt.Errorf("signature error: %v", md.SignatureError)
	}
	if md.SignedBy == nil {
		return InvalidSignatureError{msg: fmt.Sprintf("Invalid GPG signature: %#v", md.Signature)}
	}
	if md.Signature != nil {
		if md.Signature.SigLifetimeSecs != nil {
			expiry := md.Signature.CreationTime.Add(time.Duration(*md.Signature.SigLifetimeSecs) * time.Second)
			if time.Now().After(expiry) {
				return InvalidSignatureError{msg: fmt.Sprintf("Signature expired on %s", expiry)}
			}
		}
	} else if md.SignatureV3 == nil {
		// Coverage: If md.SignedBy != nil, the final md.UnverifiedBody.Read() either sets one of md.Signature or md.SignatureV3,
		// or sets md.SignatureError.
		return InvalidSignatureError{msg: "Unexpected openpgp.MessageDetails: neither Signature nor SignatureV3 is set"}
	}

	// Uppercase the fingerprint to be compatible with gpgme
	//return content, strings.ToUpper(fmt.Sprintf("%x", md.SignedBy.PublicKey.Fingerprint)), nil

	return nil
}
