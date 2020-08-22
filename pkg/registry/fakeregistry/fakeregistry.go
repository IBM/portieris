// Copyright 2018, 2020 Portieris Authors.
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

package fakeregistry

import (
	"fmt"
	"sync"

	"github.com/IBM/portieris/pkg/registry"
)

var _ registry.Interface = &FakeRegistry{}

// FakeRegistry .
type FakeRegistry struct {
	GetContentTrustTokenStub        func(username, password, imageRepo, hostname string) (string, error)
	getContentTrustTokenMutex       sync.RWMutex
	getContentTrustTokenArgsForCall []struct {
		username  string
		password  string
		imageRepo string
		hostname  string
	}
	getContentTrustTokenReturns struct {
		token string
		err   error
	}
}

// GetContentTrustToken ...
func (fake *FakeRegistry) GetContentTrustToken(username, password, imageRepo, hostname string) (string, error) {
	fake.getContentTrustTokenMutex.Lock()
	fake.getContentTrustTokenArgsForCall = append(fake.getContentTrustTokenArgsForCall, struct {
		username  string
		password  string
		imageRepo string
		hostname  string
	}{username, password, imageRepo, hostname})
	fake.getContentTrustTokenMutex.Unlock()
	if fake.GetContentTrustTokenStub != nil {
		return fake.GetContentTrustTokenStub(username, password, imageRepo, hostname)
	}
	return fake.getContentTrustTokenReturns.token, fake.getContentTrustTokenReturns.err
}

// NoAnonymousContentTrustTokenStub ...
func (fake *FakeRegistry) NoAnonymousContentTrustTokenStub(username, password, imageRepo, hostname string) (string, error) {
	if username == "" {
		return "", fmt.Errorf("not allowed")
	}
	return fake.getContentTrustTokenReturns.token, fake.getContentTrustTokenReturns.err
}

// GetContentTrustTokenReturns ...
func (fake *FakeRegistry) GetContentTrustTokenReturns(token string, err error) {
	fake.getContentTrustTokenMutex.Lock()
	defer fake.getContentTrustTokenMutex.Unlock()
	fake.getContentTrustTokenReturns = struct {
		token string
		err   error
	}{token, err}
}

// GetRegistryToken ...
func (fake *FakeRegistry) GetRegistryToken(username, password, imageRepo, hostname string) (string, error) {
	return "", fmt.Errorf("not implemented")
}
