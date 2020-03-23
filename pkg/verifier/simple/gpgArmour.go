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

package simple

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/openpgp/armor"
)

func decodeArmoredKey(keyBytes []byte) ([]byte, error) {
	if len(keyBytes) == 0 {
		return nil, fmt.Errorf("Key: empty")
	}
	block, err := armor.Decode(bytes.NewReader(keyBytes))
	if err != nil {
		if err.Error() == "EOF" {
			return nil, fmt.Errorf("Unable to decode key: %v", err)
		}
		return nil, err
	}
	switch block.Type {
	case "PGP PUBLIC KEY BLOCK":
		break
	default:
		return nil, fmt.Errorf("Expected \"PGP PUBLIC KEY BLOCK\" not found: %s", block.Type)
	}
	return ioutil.ReadAll(block.Body)
}
