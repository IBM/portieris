// Copyright 2018 IBM
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
	"sync"
)

// FakeRegistry .
type FakeRegistry struct {
	GetContentTrustTokenStub        func(registryToken, imageRepo, hostname string) (string, error)
	getContentTrustTokenMutex       sync.RWMutex
	getContentTrustTokenArgsForCall []struct {
		registryToken string
		imageRepo     string
		hostname      string
	}
	getContentTrustTokenReturns struct {
		token string
		err   error
	}
}

// GetContentTrustToken ...
func (fake *FakeRegistry) GetContentTrustToken(registryToken, imageRepo, hostname string) (string, error) {
	fake.getContentTrustTokenMutex.Lock()
	fake.getContentTrustTokenArgsForCall = append(fake.getContentTrustTokenArgsForCall, struct {
		registryToken string
		imageRepo     string
		hostname      string
	}{registryToken, imageRepo, hostname})
	fake.getContentTrustTokenMutex.Unlock()
	if fake.GetContentTrustTokenStub != nil {
		return fake.GetContentTrustTokenStub(registryToken, imageRepo, hostname)
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
