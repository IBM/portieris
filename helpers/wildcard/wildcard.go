// Copyright 2018 Portieris Authors.
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

package wildcard

import (
	"strings"
)

// Wildcard character
const wildcard = "*"

// Compare will match a string pattern which may contain wildcard
// characters against a string str. The result is a boolean.
func Compare(pattern, str string) bool {
	// Empty pattern can only match empty str
	if pattern == "" {
		return str == pattern
	}

	// If the pattern is a single wildcard character, it matches everything
	if pattern == wildcard {
		return true
	}

	// Break the pattern into pieces by wildcard characters
	pieces := strings.Split(pattern, wildcard)

	// If there is a single piece, there were no wildcard characters,
	// so test for equal strings.
	if len(pieces) == 1 {
		return str == pattern
	}

	// Does pattern start or end in a wildcard?
	startingWildcard := strings.HasPrefix(pattern, wildcard)
	endingWildcard := strings.HasSuffix(pattern, wildcard)

	// Don't test trailing piece as it requires different logic
	numPieces := len(pieces) - 1

	// Go over the pieces and ensure they match.
	for i := 0; i < numPieces; i++ {
		index := strings.Index(str, pieces[i])

		switch i {
		// First piece has different logic
		case 0:
			if !startingWildcard && index != 0 {
				return false
			}
		default:
			// Ensure that the pieces match
			if index < 0 {
				return false
			}
		}

		// Trim piece from str
		str = str[index+len(pieces[i]):]
	}

	// Reached the last piece.
	return endingWildcard || strings.HasSuffix(str, pieces[numPieces])
}

// CompareAnyTag will match a string pattern which may contain wildcard
// characters against a string str. If the compare fails, a successive compare
// is made with a ':*' added to pattern (wildcard the addition of a tag to the pattern).
// The result is a boolean based on the last test.
func CompareAnyTag(pattern, str string) bool {
	if !Compare(pattern, str) {
		return Compare(pattern+":*", str)
	}
	return true
}
