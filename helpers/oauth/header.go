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

package oauth

import (
	"net/http"
	"strings"
)

// Challenge carries information from a WWW-Authenticate response header.
// See RFC 2617.
type Challenge struct {
	// Scheme is the auth-scheme according to RFC 2617
	Scheme string

	// Parameters are the auth-params according to RFC 2617
	Parameters map[string]string
}

// Octet types from RFC 2616.
type octetType byte

var octetTypes [256]octetType

const (
	isToken octetType = 1 << iota
	isSpace
)

func init() {
	// OCTET      = <any 8-bit sequence of data>
	// CHAR       = <any US-ASCII character (octets 0 - 127)>
	// CTL        = <any US-ASCII control character (octets 0 - 31) and DEL (127)>
	// CR         = <US-ASCII CR, carriage return (13)>
	// LF         = <US-ASCII LF, linefeed (10)>
	// SP         = <US-ASCII SP, space (32)>
	// HT         = <US-ASCII HT, horizontal-tab (9)>
	// <">        = <US-ASCII double-quote mark (34)>
	// CRLF       = CR LF
	// LWS        = [CRLF] 1*( SP | HT )
	// TEXT       = <any OCTET except CTLs, but including LWS>
	// separators = "(" | ")" | "<" | ">" | "@" | "," | ";" | ":" | "\" | <">
	//              | "/" | "[" | "]" | "?" | "=" | "{" | "}" | SP | HT
	// token      = 1*<any CHAR except CTLs or separators>
	// qdtext     = <any TEXT except <">>

	for c := 0; c < 256; c++ {
		var t octetType
		isCtl := c <= 31 || c == 127
		isChar := 0 <= c && c <= 127
		isSeparator := strings.ContainsRune(" \t\"(),/:;<=>?@[]\\{}", rune(c))
		if strings.ContainsRune(" \t\r\n", rune(c)) {
			t |= isSpace
		}
		if isChar && !isCtl && !isSeparator {
			t |= isToken
		}
		octetTypes[c] = t
	}
}

// ParseAuthHeader takes the header info from a response and gives back realm and service information
func ParseAuthHeader(header http.Header) []Challenge {
	challenges := []Challenge{}
	for _, h := range header[http.CanonicalHeaderKey("WWW-Authenticate")] {
		v, p := parseValueAndParams(h)
		if v != "" {
			challenges = append(challenges, Challenge{Scheme: v, Parameters: p})
		}
	}
	return challenges
}

func parseValueAndParams(header string) (value string, params map[string]string) {
	params = make(map[string]string)
	value, s := expectToken(header)
	if value == "" {
		return
	}
	value = strings.ToLower(value)
	s = "," + skipSpace(s)
	for strings.HasPrefix(s, ",") {
		var pkey string
		pkey, s = expectToken(skipSpace(s[1:]))
		if pkey == "" {
			return
		}
		if !strings.HasPrefix(s, "=") {
			return
		}
		var pvalue string
		pvalue, s = expectTokenOrQuoted(s[1:])
		if pvalue == "" {
			return
		}
		pkey = strings.ToLower(pkey)
		params[pkey] = pvalue
		s = skipSpace(s)
	}
	return
}

func skipSpace(s string) (rest string) {
	i := 0
	for ; i < len(s); i++ {
		if octetTypes[s[i]]&isSpace == 0 {
			break
		}
	}
	return s[i:]
}

func expectToken(s string) (token, rest string) {
	i := 0
	for ; i < len(s); i++ {
		if octetTypes[s[i]]&isToken == 0 {
			break
		}
	}
	return s[:i], s[i:]
}

func expectTokenOrQuoted(s string) (value, rest string) {
	if !strings.HasPrefix(s, "\"") {
		return expectToken(s)
	}
	s = s[1:]
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			return s[:i], s[i+1:]
		case '\\':
			p := make([]byte, len(s)-1)
			j := copy(p, s[:i])
			escape := true
			for i = i + 1; i < len(s); i++ {
				b := s[i]
				switch {
				case escape:
					escape = false
					p[j] = b
					j++
				case b == '\\':
					escape = true
				case b == '"':
					return string(p[:j]), s[i+1:]
				default:
					p[j] = b
					j++
				}
			}
			return "", ""
		}
	}
	return "", ""
}
