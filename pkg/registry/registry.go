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

package registry

import (
	"fmt"

	"github.com/docker/distribution/registry/client/auth/challenge"
	"github.com/golang/glog"

	"github.com/IBM/portieris/helpers/oauth"
)

// Client .
type Client struct{}

// Interface .
type Interface interface {
	GetContentTrustToken(username, password, imageRepo string, challengeSlice []challenge.Challenge) (string, error)
}

// NewClient .
func NewClient() Interface {
	return &Client{}
}

// GetContentTrustToken .
func (c Client) GetContentTrustToken(username, password, imageRepo string, challengeSlice []challenge.Challenge) (string, error) {
	token, err := oauth.Request(password, imageRepo, username, challengeSlice)
	if err != nil {
		return "", err
	}
	useToken := ""
	if token.Token != "" {
		glog.Info("Using token")
		useToken = token.Token
	} else if token.AccessToken != "" {
		glog.Info("Using access token")
		useToken = token.AccessToken
	} else {
		return "", fmt.Errorf("neither token or access token could be found")
	}
	return useToken, nil
}
