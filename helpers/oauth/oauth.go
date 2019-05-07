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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/glog"
)

var client = &http.Client{
	Timeout: 10 * time.Minute,
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 10,
		TLSHandshakeTimeout: 5 * time.Second,
	},
}

// Request is a helper for getting an OAuth token from the Registry OAuth Service.
// Takes the following as input:
//   token               - Auth token being used for the request
//   repo                - Repo you are requesting access too e.g. bainsy88/busybox
//   username            - Username for the OAuth request, identifies the type of token being passed in. Valid usernames are token (for registry token), iambearer, iamapikey, bearer (UAA bearer (legacy)), iamrefresh
//   writeAccessRequired - Whether or not you require write (push and delete) access as well as read (pull)
//   service             - The service you are retrieving the OAuth token for. Current services are either "notary" or "registry"
//   hostname            - Hostname of the registry you wish to call e.g. https://icr.io
// Returns:
//   *auth.TokenResponse - Details of the type is here https://github.ibm.com/alchemy-registry/registry-types/tree/master/auth#type-tokenresponse
//                         Token is the element you will need to forward to the registry/notary as part of a Bearer Authorization Header
//   error
func Request(token string, repo string, username string, writeAccessRequired bool, challengeSlice []Challenge) (*TokenResponse, error) {
	var actions string
	//If you want to verify if a the credential supplied has read and write access to the repo we ask oauth for pull,push and *
	if writeAccessRequired {
		actions = "pull,push,*"
	} else {
		actions = "pull"
	}

	oauthEndpoint := ""
	service := ""

	if challengeSlice == nil {
		errMessage := "unable to fetch www-authenticate header details"
		glog.Errorf(errMessage)
		return nil, fmt.Errorf(errMessage)
	}

	for _, challenge := range challengeSlice {
		oauthEndpoint = challenge.Parameters["realm"]
		service = challenge.Parameters["service"]
	}

	if oauthEndpoint == "" || service == "" {
		errMessage := "unable to fetch oauth realm and service header details"
		glog.Errorf(errMessage)
		return nil, fmt.Errorf(errMessage)
	}

	glog.Infof("Calling oauth endpoint: %s for registry service: %s", oauthEndpoint, service)
	var resp *http.Response
	var err error
	if oauthEndpoint == "https://auth.docker.io/token" {
		resp, err = client.Get(oauthEndpoint)
		if err != nil {
			glog.Errorf("Error sending GET request to registry-oauth: %v", err)
			return nil, err
		}
	} else {
		resp, err = client.PostForm(oauthEndpoint, url.Values{
			"service":    {service},
			"grant_type": {"password"},
			"client_id":  {"testclient"},
			"username":   {username},
			"password":   {token},
			"scope":      {"repository:" + repo + ":" + actions},
		})
		if err != nil {
			glog.Errorf("Error sending POST request to registry-oauth: %v", err)
			return nil, err
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Unexpected, read body for more information and close. It is the upstream callers
		// responsibility to close the response body if an error is not returned.
		glog.Errorf("Received non-success status code %v", resp.StatusCode)
		var body []byte
		if resp.Body != nil {
			body, _ = ioutil.ReadAll(resp.Body)
		}
		return nil, fmt.Errorf("Request to OAuth failed with status code: %v and body: %s", resp.StatusCode, body)
	}

	tokenResponse := TokenResponse{}
	bytes, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bytes, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshall OAuth response: %s", err)
	}

	return &tokenResponse, nil
}

// CheckAuthRequired - checks if the given image needs to be authenticated to fetch metadata or not and returns the response
func CheckAuthRequired(notaryURL, hostName, repoName string) (*http.Response, error) {
	glog.Infof("Notary URL: %s Hostname %s RepoName %s", notaryURL, hostName, repoName)
	// Github issue 51 Fix
	req, err := http.NewRequest("GET", notaryURL+"/v2/"+hostName+"/"+repoName+"/_trust/tuf/root.json", nil)

	resp, err := client.Do(req)

	if err != nil {
		glog.Errorf("Failed to query v2 tuf endpoint for notaryURL: %s", notaryURL)
		return nil, err
	}

	return resp, nil
}
