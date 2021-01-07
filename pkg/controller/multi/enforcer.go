// Copyright 2020 Portieris Authors.
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

package multi

import (
	"bytes"
	"fmt"

	"github.com/IBM/portieris/helpers/credential"
	"github.com/IBM/portieris/helpers/image"
	portieriscloudibmcomv1 "github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/v1"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/IBM/portieris/pkg/verifier/simple"
	notaryverifier "github.com/IBM/portieris/pkg/verifier/trust"
	"github.com/IBM/portieris/pkg/verifier/vulnerability"
	"github.com/golang/glog"
)

// Enforcer is an interface that enforces pod admission based on a configured policy
type Enforcer interface {
	DigestByPolicy(string, *image.Reference, credential.Credentials, *portieriscloudibmcomv1.Policy) (*bytes.Buffer, error, error)
	VulnerabilityPolicy(*image.Reference, credential.Credentials, *portieriscloudibmcomv1.Policy) vulnerability.ScanResponse
}

type enforcer struct {
	// kubeClientsetWrapper is a standard kubernetes clientset with a wrapper for retrieving podSpec from a given object
	kubeClientsetWrapper kubernetes.WrapperInterface
	// nv notary signing verifier
	nv notaryverifier.Interface
	// simple signing verifier
	sv simple.Verifier
	// scannerFactory creates new vulnerabilities scanners according to the policy
	scannerFactory vulnerability.ScannerFactory
}

// NewEnforcer returns an enforce that wraps the kubenetes interface and a notary verifier
func NewEnforcer(kubeClientsetWrapper kubernetes.WrapperInterface, nv *notaryverifier.Verifier) Enforcer {
	scannerFactory := vulnerability.NewScannerFactory()
	return &enforcer{
		kubeClientsetWrapper: kubeClientsetWrapper,
		nv:                   nv,
		sv:                   simple.NewVerifier(),
		scannerFactory:       &scannerFactory,
	}
}

func (e enforcer) DigestByPolicy(namespace string, img *image.Reference, credentials credential.Credentials, policy *securityenforcementv1beta1.Policy) (*bytes.Buffer, error, error) {
	// no policy indicates admission should be allowed, without mutation
	if policy == nil {
		return nil, nil, nil
	}

	var digest *bytes.Buffer
	var deny, err error
	if len(policy.Simple.Requirements) > 0 {
		glog.Infof("policy.Simple %v", policy.Simple)
		simplePolicy, err := e.sv.TransformPolicies(e.kubeClientsetWrapper, namespace, policy.Simple.Requirements)
		if err != nil {
			return nil, nil, err
		}
		storeUser, storePassword, err := e.kubeClientsetWrapper.GetBasicCredentials(namespace, policy.Simple.StoreSecret)
		if err != nil {
			return nil, nil, err
		}
		storeConfigDir, err := e.sv.CreateRegistryDir(policy.Simple.StoreURL, storeUser, storePassword)
		if err != nil {
			return nil, nil, err
		}
		digest, deny, err = e.sv.VerifyByPolicy(img.String(), credentials, storeConfigDir, simplePolicy)
		if err != nil {
			return nil, nil, fmt.Errorf("simple: %v", err)
		}
		err = e.sv.RemoveRegistryDir(storeConfigDir)
		if err != nil {
			glog.Warningf("failed to remove %s, %v", storeConfigDir, err)
		}
		if deny != nil {
			return nil, fmt.Errorf("simple: policy denied the request: %v", deny), nil
		}
	}

	if policy.Trust.Enabled != nil && *policy.Trust.Enabled {
		glog.Infof("policy.Trust %v", policy.Trust)
		var notaryDigest *bytes.Buffer
		notaryDigest, deny, err = e.nv.VerifyByPolicy(namespace, img, credentials, policy)
		if err != nil {
			return nil, nil, fmt.Errorf("trust: %v", err)
		}
		if deny != nil {
			return nil, fmt.Errorf("trust: policy denied the request: %v", deny), nil
		}
		glog.Infof("DCT digest: %v", notaryDigest)
		if notaryDigest != nil {
			if digest != nil && notaryDigest != digest {
				return nil, fmt.Errorf("Notary signs conflicting digest: %v simple: %v", notaryDigest, digest), nil
			}
			digest = notaryDigest
		}
	}

	return digest, nil, nil
}

func (e *enforcer) VulnerabilityPolicy(img *image.Reference, credentials credential.Credentials, policy *securityenforcementv1beta1.Policy) vulnerability.ScanResponse {
	if policy == nil {
		glog.Warningf("vulnerability: No policy for image %q so allow", img.String())
		return vulnerability.ScanResponse{CanDeploy: true}
	}

	scanners := e.scannerFactory.GetScanners(*img, credentials, *policy)
	// Loop round all scanners and check if the image can be deployed
	// If any scanner returns either an error, or a CanDeploy=false, the pod will not be admitted
	for _, scanner := range scanners {
		response, err := scanner.CanImageDeployBasedOnVulnerabilities(*img)
		if err != nil {
			return vulnerability.ScanResponse{CanDeploy: false, DenyReason: err.Error()}
		}
		if !response.CanDeploy {
			return response
		}
	}

	return vulnerability.ScanResponse{CanDeploy: true}
}
