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

package trust

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/IBM/portieris/helpers/image"
	"github.com/IBM/portieris/helpers/oauth"
	securityenforcementv1beta1 "github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/IBM/portieris/pkg/notary"
	registryclient "github.com/IBM/portieris/pkg/registry"
	"github.com/golang/glog"
	store "github.com/theupdateframework/notary/storage"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var codec = serializer.NewCodecFactory(runtime.NewScheme())

// Verifier is the notary controller
type Verifier struct {
	// kubeClientsetWrapper is a standard kubernetes clientset with a wrapper for retrieving podSpec from a given object
	kubeClientsetWrapper kubernetes.WrapperInterface
	// Trust Client
	trust notary.Interface
	// Container Registry client
	cr registryclient.Interface
}

// NewVerifier creates a new controller object from the various clients passed in
func NewVerifier(kubeWrapper kubernetes.WrapperInterface, trust notary.Interface, cr registryclient.Interface) *Verifier {
	return &Verifier{
		kubeClientsetWrapper: kubeWrapper,
		trust:                trust,
		cr:                   cr,
	}
}

// VerifyByPolicy ...
func (v *Verifier) VerifyByPolicy(namespace string, img *image.Reference, credentials [][]string, policy *securityenforcementv1beta1.Policy) (*bytes.Buffer, error, error) {
	notaryURL := policy.Trust.TrustServer
	var err error
	if notaryURL == "" {
		notaryURL, err = img.GetContentTrustURL()
		if err != nil {
			return nil, nil, fmt.Errorf("Trust Server/Image Configuration Error: %v", err.Error())
		}
	}

	var signers []Signer
	if policy.Trust.SignerSecrets != nil {
		// Generate a []Singer with the values for each signerSecret
		signers = make([]Signer, len(policy.Trust.SignerSecrets))
		for i, secretName := range policy.Trust.SignerSecrets {
			signers[i], err = v.getSignerSecret(namespace, secretName.Name)
			if err != nil {
				return nil, nil, fmt.Errorf("Deny %q, could not get signerSecret from your cluster, %s", img.String(), err.Error())
			}
		}
	}

	official := !strings.ContainsRune(img.RepoName(), '/')

	hostname := img.GetHostname()
	port := img.GetPort()
	if port != "" {
		port = ":" + port
	}
	repo := hostname + port

	resp, err := oauth.CheckAuthRequired(notaryURL, repo, img.RepoName(), official)
	if err != nil {
		glog.Error(err)
		return nil, nil, fmt.Errorf("Some error occurred while checking if authentication is required to fetch target metadata")
	}

	if resp.StatusCode == http.StatusUnauthorized {
		glog.Infof("Need to get token for %s to fetch target metadata", notaryURL)
	} else if resp.StatusCode == http.StatusOK {
		glog.Infof("No need to fetch token for %s to get the target metadata", notaryURL)
	} else {
		glog.Infof("Status code: %v was returned", resp.StatusCode)
		return nil, nil, fmt.Errorf("Status code: %v was returned while checking if authentication is required to fetch target metadata", resp.StatusCode)
	}

	glog.Infof("Status code: %v returned for repo: %v", resp.StatusCode, img.NameWithoutTag())

	var challengeSlice []oauth.Challenge

	if resp.StatusCode == http.StatusUnauthorized {
		challengeSlice = oauth.ResponseChallenges(resp)
		if err != nil {
			glog.Error(err)
			return nil, nil, fmt.Errorf("Some error occurred when fetching challenge slice %s", err.Error())
		}
	}

	for _, cred := range credentials {
		notaryToken, err := v.cr.GetContentTrustToken(cred[0], cred[1], img.NameWithoutTag(), challengeSlice)
		if err != nil {
			glog.Error(err)
			continue
		}

		// Get image digest
		glog.Infof("getting signed image... %v", img.RepoName())
		// notaryToken will be blank for unauthorized calls
		var image string
		if official {
			image = "docker.io/library/" + img.RepoName()
		} else {
			image = img.NameWithoutTag()
		}
		glog.Infof("Image: %v", image)

		digest, err := v.getDigest(notaryURL, image, notaryToken, img.GetTag(), signers)
		if err != nil {
			if strings.Contains(err.Error(), "401") {
				continue
			}

			if _, ok := err.(store.ErrServerUnavailable); ok {
				glog.Errorf("Trust server unavailable: %v", err)
				return nil, nil, fmt.Errorf("Deny %q, failed to get content trust information: %s", img.String(), err.Error())
			}
			return nil, fmt.Errorf("Deny %q, failed to get content trust information: %s", img.String(), err.Error()), nil
		}
		return digest, nil, nil
	}

	// if no credentials defined and pulling signed images from public docker
	notaryToken, err := v.cr.GetContentTrustToken("", "", img.NameWithoutTag(), challengeSlice)
	if err != nil {
		glog.Error(err)
		return nil, nil, fmt.Errorf("Some error occurred while trying to fetch token for unauthenticated pubilc pull %s", err.Error())
	}

	// Get image digest
	glog.Infof("getting signed image for %v", img.RepoName())
	// notaryToken will be blank for unauthorized calls
	var image string
	if official {
		image = "docker.io/library/" + img.RepoName()
	} else {
		image = img.NameWithoutTag()
	}
	glog.Infof("Image: %v and tag: %v", image, img.GetTag())
	glog.Infof("Notary URL: %v", notaryURL)
	digest, err := v.getDigest(notaryURL, image, notaryToken, img.GetTag(), signers)
	if err != nil {
		glog.Infof(err.Error())
		if strings.Contains(err.Error(), "401") {
			return nil, fmt.Errorf("Deny %q, no valid ImagePullSecret defined for %s", img.String(), img.String()), nil
		}

		if _, ok := err.(store.ErrServerUnavailable); ok {
			glog.Errorf("Trust server unavailable: %v", err)
			return nil, nil, fmt.Errorf("Deny %q, failed to get content trust information: %s", img.String(), err.Error())
		}
		return nil, fmt.Errorf("Deny %q, failed to get content trust information: %s", img.String(), err.Error()), nil
	}
	return digest, nil, nil
}
