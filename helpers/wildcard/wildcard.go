// Copyright 2018, 2026 Portieris Authors.
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

// CompareImageRef matches a policy pattern against an image reference string,
// enforcing that the host component of the pattern is matched only against the
// host component of the image reference and never as a substring of the path.
// This prevents an attacker from hosting a malicious image at
// attacker.example.com/x.trusted.example.com/myorg/evil and having it match a
// policy pattern such as *.trusted.example.com/myorg/* by embedding the
// trusted-registry hostname inside the repository path.
//
// If the pattern contains no '/' it is treated as a pure image/tag pattern and
// falls back to CompareAnyTag (e.g. the bare "*" catch-all policy).
//
// For patterns that contain '/', the string is split at the first '/' into a
// host part and a path part, and each is matched independently via
// CompareAnyTag so that wildcard semantics within each segment are preserved.
func CompareImageRef(pattern, imageRef string) bool {
	slashIdx := strings.Index(pattern, "/")
	if slashIdx < 0 {
		// No '/' in pattern — pure tag/image glob or bare "*", use existing logic.
		return CompareAnyTag(pattern, imageRef)
	}

	patternHost := pattern[:slashIdx]
	patternPath := pattern[slashIdx+1:]

	// Split the image reference at its first '/'.
	// A valid image reference always has a host component when a registry is
	// specified; if there is no '/' the image ref itself has no host segment,
	// so it cannot match a host/path policy pattern.
	refSlashIdx := strings.Index(imageRef, "/")
	if refSlashIdx < 0 {
		return false
	}

	refHost := imageRef[:refSlashIdx]
	refPath := imageRef[refSlashIdx+1:]

	// Both sides of the split are matched with CompareAnyTag which preserves
	// wildcard semantics within each segment.
	// The host match uses Compare (not CompareAnyTag) because a ':*' suffix
	// is never meaningful for a hostname.
	return Compare(patternHost, refHost) && CompareAnyTag(patternPath, refPath)
}

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
