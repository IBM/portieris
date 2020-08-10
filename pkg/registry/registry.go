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
	"github.com/IBM/portieris/helpers/oauth"
)

// Client .
type Client struct{}

// Interface .
type Interface interface {
	GetContentTrustToken(username, password, imageRepo, hostname string) (string, error)
	GetRegistryToken(username, password, imageRepo, hostname string) (string, error)
}

// NewClient .
func NewClient() Interface {
	return &Client{}
}

// GetContentTrustToken .
func (c Client) GetContentTrustToken(username, password, imageRepo, hostname string) (string, error) {
	if username == "" && password == "" {
		return "", nil
	}
	token, err := oauth.Request(password, imageRepo, username, false, "notary", hostname)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

// GetRegistryToken .
func (c Client) GetRegistryToken(username, password, imageRepo, hostname string) (string, error) {
	token, err := oauth.Request(password, imageRepo, username, false, "registry", hostname)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}
