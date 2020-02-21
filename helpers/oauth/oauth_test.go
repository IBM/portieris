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

	"github.com/IBM/portieris/helpers/image"
	challenge "github.com/docker/distribution/registry/client/auth/challenge"
)

func TestHappyPathWithAuth(t *testing.T) {
	notaryURL := "https://notary.docker.io"
	img, _ := image.NewReference("nginx")

	resp, err := CheckAuthRequired(notaryURL, *img)

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
	img, _ := image.NewReference("docker.io/library/nginx")

	resp, err := CheckAuthRequired(notaryURL, *img)

	if err != nil {
		t.Errorf("Some error occurred: %s", err.Error())
	}

	if expected := http.StatusUnauthorized; resp.StatusCode != expected {
		t.Errorf("Unexpected status code: %v, expected: %v", resp.StatusCode, expected)
	}
}

func TestSadPathWithAuth(t *testing.T) {
	notaryURL := "https://invalid.docker.io"
	img, _ := image.NewReference("docker.io/library/nginx")

	_, err := CheckAuthRequired(notaryURL, *img)

	if err != nil {
		expected := "Get https://invalid.docker.io/v2/docker.io/library/docker.io/library/nginx/_trust/tuf/root.json: dial tcp: lookup invalid.docker.io: no such host"
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
	img, _ := image.NewReference("nginx")

	resp, _ := CheckAuthRequired(notaryURL, *img)

	challengeSlice := challenge.ResponseChallenges(resp)

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
	img, _ := image.NewReference("us.icr.io/molepigeon/testimage")

	resp, _ := CheckAuthRequired(notaryURL, *img)

	challengeSlice := challenge.ResponseChallenges(resp)

	_, err := Request(password, repo, username, challengeSlice)

	if err != nil {
		//todo(chaosaffe): Actually check that we get the expected error
		// do nothing as we are expecting error of 401
	}

}

func TestSadPathWithRequestInvalidURL(t *testing.T) {
	notaryURL := "https://usa.icr.io:4443"
	img, _ := image.NewReference("us.icr.io/molepigeon/testimage")

	_, err := CheckAuthRequired(notaryURL, *img)

	if err != nil {
		if expected := "Get https://usa.icr.io:4443/v2/us.icr.io/molepigeon/us.icr.io/molepigeon/testimage/_trust/tuf/root.json: x509: certificate is valid for icr.io, va.icr.io, registry.bluemix.net, va.bluemix.net, cp.icr.io, registry.marketplace.redhat.com, not usa.icr.io"; err.Error() != expected {
			t.Errorf("Unexpected error message: %v, expected: %v", err.Error(), expected)
		}
	}

}

func TestSadPathWithRequestMissingWWWAuthenticate(t *testing.T) {
	password := ""
	username := ""
	repo := "nginx"

	notaryURL := "https://notary.docker.io"
	img, _ := image.NewReference("nginx")

	resp, _ := CheckAuthRequired(notaryURL, *img)

	challengeSlice := challenge.ResponseChallenges(resp)

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

	challengeSlice := []challenge.Challenge{challenge.Challenge{Scheme: "test", Parameters: nil}}

	_, err := Request(password, repo, username, challengeSlice)

	if err != nil {
		if expected := "unable to fetch oauth realm and service header details"; err.Error() != expected {
			t.Errorf("Unexpected error message: %v, expected: %v", err.Error(), expected)
		}
	}
}
