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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	challenge "github.com/docker/distribution/registry/client/auth/challenge"
	"github.com/golang/glog"
)

// GetHTTPClient gets an http client to use for getting an oauth token
// Takes the following as input:
//   customFile - Path to custom ca certificate
// Returns:
//   *http.Client
func GetHTTPClient(customFile string) *http.Client {
	rootCA, err := x509.SystemCertPool()
	customCA, err := ioutil.ReadFile(customFile)
	if err != nil {
		if os.IsNotExist(err) {
			glog.Infof("CA not provided at %s, will use default system pool", customFile)
		} else {
			glog.Fatalf("Could not read %s: %s", customFile, err)
		}
	} else {
		rootCA.AppendCertsFromPEM(customCA)
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			DisableKeepAlives:   false,
			MaxIdleConnsPerHost: 10,
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig: &tls.Config{
				// Avoid fallback by default to SSL protocols < TLS1.2
				MinVersion:               tls.VersionTLS12,
				PreferServerCipherSuites: true,
				RootCAs:                  rootCA,
			},
		},
	}
	return client
}

// Request is a helper for getting an OAuth token from the Registry OAuth Service.
// Takes the following as input:
//   token               - Auth token being used for the request
//   repo                - Repo you are requesting access too e.g. bainsy88/busybox
//   username            - Username for the OAuth request, identifies the type of token being passed in. Valid usernames are token (for registry token), iambearer, iamapikey, bearer (UAA bearer (legacy)), iamrefresh
//   service             - The service you are retrieving the OAuth token for. Current services are either "notary" or "registry"
//   hostname            - Hostname of the registry you wish to call e.g. https://icr.io
// Returns:
//   *auth.TokenResponse - Details of the type is here https://github.ibm.com/alchemy-registry/registry-types/tree/master/auth#type-tokenresponse
//                         Token is the element you will need to forward to the registry/notary as part of a Bearer Authorization Header
//   error
func Request(token string, repo string, username string, challengeSlice []challenge.Challenge) (*TokenResponse, error) {

	client := GetHTTPClient("/etc/certs/ca.pem")

	// Github issue 51 Fix
	req, err := http.NewRequest("GET", hostname+"/v2/", nil)

	resp, err := client.Do(req)

	if err != nil {
		glog.Errorf("Failed to query v2 endpoint for hostname: %s", hostname)
		return nil, err
	}

	challengeSlice := ParseAuthHeader(resp.Header)

	oauthEndpoint := ""
	service := ""
	scope := ""

	if challengeSlice == nil {
		errMessage := "unable to fetch www-authenticate header details"
		glog.Errorf(errMessage)
		return nil, fmt.Errorf(errMessage)
	}

	for _, challenge := range challengeSlice {
		oauthEndpoint = challenge.Parameters["realm"]
		service = challenge.Parameters["service"]
		scope = challenge.Parameters["scope"]
	}

	if oauthEndpoint == "" || service == "" {
		errMessage := "unable to fetch oauth realm and service header details"
		glog.Errorf(errMessage)
		return nil, fmt.Errorf(errMessage)
	}

	glog.Infof("Calling oauth endpoint: %s for registry service: %s and scope %s", oauthEndpoint, service, scope)
	var resp *http.Response
	var err error
	resp, err = client.PostForm(oauthEndpoint, url.Values{
		"service":    {service},
		"grant_type": {"password"},
		"client_id":  {"portieris-client"},
		"username":   {username},
		"password":   {token},
		"scope":      {scope},
	})
	if err != nil {
		glog.Errorf("Error sending POST request to registry-oauth: %v", err)
		return nil, err
	}

	// TODO: confirm if status code of 405 needs to be handled in the below block
	if resp.StatusCode == 404 || resp.StatusCode == 405 {
		glog.Info("Calling: " + oauthEndpoint + "?service=" + service + "&scope=" + scope)
		getURL, err := url.Parse(oauthEndpoint)
		if err != nil {
			return nil, err
		}
		q := getURL.Query()
		q.Set("service", service)
		q.Set("scope", scope)
		getURL.RawQuery = q.Encode()
		resp, err = client.Get(getURL.String())
		if err != nil {
			glog.Errorf("Error sending GET request to registry-oauth: %v", err)
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
func CheckAuthRequired(notaryURL, hostName, repoName string, official bool) (*http.Response, error) {
	glog.Infof("Notary URL: %s Hostname %s RepoName %s", notaryURL, hostName, repoName)
	// Github issue 51 Fix
	var req *http.Request
	var err error
	if hostName == "docker.io" && official {
		req, err = http.NewRequest("GET", notaryURL+"/v2/"+hostName+"/library/"+repoName+"/_trust/tuf/root.json", nil)
	} else {
		req, err = http.NewRequest("GET", notaryURL+"/v2/"+hostName+"/"+repoName+"/_trust/tuf/root.json", nil)
	}

	resp, err := client.Do(req)

	if err != nil {
		glog.Errorf("Failed to query v2 tuf endpoint for notaryURL: %s", notaryURL)
		return nil, err
	}

	return resp, nil
}
