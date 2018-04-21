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

package fakenotary

import (
	"sync"

	notaryclient "github.com/theupdateframework/notary/client"
)

// FakeNotary .
type FakeNotary struct {
	GetNotaryRepoStub        func(server, image, notaryToken string) (notaryclient.Repository, error)
	getNotaryRepoMutex       sync.RWMutex
	GetNotaryRepoArgsForCall []struct {
		server      string
		image       string
		notaryToken string
	}
	getNotaryRepoReturns []struct {
		notaryRepo notaryclient.Repository
		err        error
	}
}

// GetNotaryRepo ...
func (fake *FakeNotary) GetNotaryRepo(server, image, notaryToken string) (notaryclient.Repository, error) {
	fake.getNotaryRepoMutex.Lock()
	fake.GetNotaryRepoArgsForCall = append(fake.GetNotaryRepoArgsForCall, struct {
		server      string
		image       string
		notaryToken string
	}{server, image, notaryToken})
	fake.getNotaryRepoMutex.Unlock()
	if fake.GetNotaryRepoStub != nil {
		return fake.GetNotaryRepoStub(server, image, notaryToken)
	}

	if len(fake.getNotaryRepoReturns) < 1 {
		panic("GetNotaryRepo called before it is stubbed")
	}
	returns := fake.getNotaryRepoReturns[0]
	fake.getNotaryRepoReturns = fake.getNotaryRepoReturns[1:]

	return returns.notaryRepo, returns.err
}

// GetNotaryRepoReturns ...
func (fake *FakeNotary) GetNotaryRepoReturns(notaryRepo notaryclient.Repository, err error) {
	fake.getNotaryRepoMutex.Lock()
	defer fake.getNotaryRepoMutex.Unlock()
	fake.getNotaryRepoReturns = append(fake.getNotaryRepoReturns, struct {
		notaryRepo notaryclient.Repository
		err        error
	}{notaryRepo, err})
}
