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
	"testing"
)

func TestHappyPathWithAuth(t *testing.T) {
	notaryURL := "https://notary.docker.io"
	hostName := "docker.io"
	repoName := "nginx"
	official := true

	resp, err := CheckAuthRequired(notaryURL, hostName, repoName, official)

	if err != nil {
		t.Errorf("Some error occurred: %s", err.Error())
	}

	if expected := http.StatusUnauthorized; resp.StatusCode != expected {
		t.Errorf("Unexpected status code: %v, expected: %v", resp.StatusCode, expected)
	}
}

// TODO: find an endpoint which allows notary tuf info without auth
func TestHappyPathWithAuthUnofficial(t *testing.T) {
	notaryURL := "https://notary.docker.io"
	hostName := "docker.io"
	repoName := "library/nginx"
	official := false

	resp, err := CheckAuthRequired(notaryURL, hostName, repoName, official)

	if err != nil {
		t.Errorf("Some error occurred: %s", err.Error())
	}

	if expected := http.StatusUnauthorized; resp.StatusCode != expected {
		t.Errorf("Unexpected status code: %v, expected: %v", resp.StatusCode, expected)
	}
}

func TestSadPathWithAuth(t *testing.T) {
	notaryURL := "https://invalid.docker.io"
	hostName := "docker.io"
	repoName := "library/nginx"
	official := false

	_, err := CheckAuthRequired(notaryURL, hostName, repoName, official)

	if err != nil {
		expected := "https://invalid.docker.io/v2/docker.io/library/nginx/_trust/tuf/root.json: dial tcp: lookup invalid.docker.io"
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Unexpected error message: %v, expected: %v", err.Error(), expected)
		}
	}
}

func TestHappyPathWithRequest(t *testing.T) {
	password := ""
	username := ""
	repo := "nginx"

	notaryURL := "https://notary.docker.io"
	hostName := "docker.io"
	repoName := "nginx"
	official := true

	resp, _ := CheckAuthRequired(notaryURL, hostName, repoName, official)

	challengeSlice := ResponseChallenges(resp)

	token, err := Request(password, repo, username, challengeSlice)

	if err != nil {
		t.Errorf("Some error occurred: %s", err.Error())
	}

	if token.AccessToken == "" && token.Token == "" {
		t.Errorf("Token not found. Expected access token or token from the response")
	}
}

func TestSadPathWithRequest(t *testing.T) {
	password := ""
	username := ""
	repo := "nginx"

	notaryURL := "https://us.icr.io:4443"
	hostName := "us.icr.io"
	repoName := "molepigeon/testimage"
	official := false

	resp, _ := CheckAuthRequired(notaryURL, hostName, repoName, official)

	challengeSlice := ResponseChallenges(resp)

	_, err := Request(password, repo, username, challengeSlice)

	if err != nil {
		// do nothing as we are expecting error of 401
	}

}

func TestSadPathWithRequestInvalidURL(t *testing.T) {
	notaryURL := "https://usa.icr.io:4443"
	hostName := "us.icr.io"
	repoName := "molepigeon/testimage"
	official := false

	_, err := CheckAuthRequired(notaryURL, hostName, repoName, official)

	if err != nil {
		if expected := "Get https://usa.icr.io:4443/v2/us.icr.io/molepigeon/testimage/_trust/tuf/root.json: x509: certificate is valid for icr.io, va.icr.io, registry.bluemix.net, va.bluemix.net, cp.icr.io, not usa.icr.io"; err.Error() != expected {
			t.Errorf("Unexpected error message: %v, expected: %v", err.Error(), expected)
		}
	}

}

func TestSadPathWithRequestMissingWWWAuthenticate(t *testing.T) {
	password := ""
	username := ""
	repo := "nginx"

	notaryURL := "https://notary.docker.io"
	hostName := "docker.io"
	repoName := "nginx"
	official := true

	resp, _ := CheckAuthRequired(notaryURL, hostName, repoName, official)

	challengeSlice := ResponseChallenges(resp)

	challengeSlice = nil

	_, err := Request(password, repo, username, challengeSlice)

	if err != nil {
		if expected := "unable to fetch www-authenticate header details"; err.Error() != expected {
			t.Errorf("Unexpected error message: %v, expected: %v", err.Error(), expected)
		}
	}
}

func TestSadPathWithRequestMissingRealmAndServiceMissing(t *testing.T) {
	password := ""
	username := ""
	repo := "nginx"

	notaryURL := "https://notary.docker.io"
	hostName := "docker.io"
	repoName := "nginx"
	official := true

	resp, _ := CheckAuthRequired(notaryURL, hostName, repoName, official)

	challengeSlice := ResponseChallenges(resp)

	challengeSlice = nil

	challengeSlice = append(challengeSlice, Challenge{Scheme: "test", Parameters: nil})

	_, err := Request(password, repo, username, challengeSlice)

	if err != nil {
		if expected := "unable to fetch oauth realm and service header details"; err.Error() != expected {
			t.Errorf("Unexpected error message: %v, expected: %v", err.Error(), expected)
		}
	}
}
