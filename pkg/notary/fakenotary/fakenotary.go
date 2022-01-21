// Copyright 2018, 2022 Portieris Authors.
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

	"github.com/IBM/portieris/helpers/image"
	"github.com/IBM/portieris/pkg/notary"
	notaryclient "github.com/theupdateframework/notary/client"
)

// FakeNotary .
type FakeNotary struct {
	GetNotaryRepoStub        func(server, image, notaryToken string) (notaryclient.Repository, error)
	getNotaryRepoMutex       sync.RWMutex
	GetNotaryRepoArgsForCall []struct {
		Server      string
		Image       string
		NotaryToken string
	}
	getNotaryRepoReturns []struct {
		notaryRepo notaryclient.Repository
		err        error
	}
	CheckAuthRequiredStub        func(notaryURL string, img *image.Reference) (*notary.AuthEndpoint, error)
	checkAuthRequiredMutex       sync.RWMutex
	CheckAuthRequiredArgsForCall []struct {
		NotaryURL string
		Img       *image.Reference
	}
	checkAuthRequiredReturns []struct {
		endpoint *notary.AuthEndpoint
		err      error
	}
}

// GetNotaryRepo ...
func (fake *FakeNotary) GetNotaryRepo(server, image, notaryToken string) (notaryclient.Repository, error) {
	fake.getNotaryRepoMutex.Lock()
	fake.GetNotaryRepoArgsForCall = append(fake.GetNotaryRepoArgsForCall, struct {
		Server      string
		Image       string
		NotaryToken string
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

// CheckAuthRequired ...
func (fake *FakeNotary) CheckAuthRequired(notaryURL string, img *image.Reference) (*notary.AuthEndpoint, error) {
	fake.checkAuthRequiredMutex.Lock()
	fake.CheckAuthRequiredArgsForCall = append(fake.CheckAuthRequiredArgsForCall, struct {
		NotaryURL string
		Img       *image.Reference
	}{notaryURL, img})
	fake.checkAuthRequiredMutex.Unlock()
	if fake.CheckAuthRequiredStub != nil {
		return fake.CheckAuthRequiredStub(notaryURL, img)
	}

	if len(fake.checkAuthRequiredReturns) < 1 {
		panic("CheckAuthRequired called before it is stubbed")
	}
	returns := fake.checkAuthRequiredReturns[0]
	fake.checkAuthRequiredReturns = fake.checkAuthRequiredReturns[1:]

	return returns.endpoint, returns.err
}

// CheckAuthRequiredReturns ...
func (fake *FakeNotary) CheckAuthRequiredReturns(endpoint *notary.AuthEndpoint, err error) {
	fake.checkAuthRequiredMutex.Lock()
	defer fake.checkAuthRequiredMutex.Unlock()
	fake.checkAuthRequiredReturns = append(fake.checkAuthRequiredReturns, struct {
		endpoint *notary.AuthEndpoint
		err      error
	}{endpoint, err})
}

// DefaultAuthEndpointStub does always return an auth endpoint
func (fake *FakeNotary) DefaultAuthEndpointStub(notaryURL string, img *image.Reference) (*notary.AuthEndpoint, error) {
	return &notary.AuthEndpoint{URL: notaryURL + "/oauth/token", Service: "notary", Scope: "pull"}, nil
}

// NoAuthRequiredStub does not return an auth endpoint to simulate cases where no authentication is required
func (fake *FakeNotary) NoAuthRequiredStub(notaryURL string, img *image.Reference) (*notary.AuthEndpoint, error) {
	return nil, nil
}
