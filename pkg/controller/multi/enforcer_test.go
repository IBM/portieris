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
	"testing"

	"github.com/IBM/portieris/helpers/credential"
	"github.com/IBM/portieris/helpers/image"
	"github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/IBM/portieris/pkg/verifier/vulnerability"
	"github.com/containers/image/v5/signature"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockScannerFactory struct {
	mock.Mock
}

func (msf *mockScannerFactory) GetScanners(img image.Reference, credentials credential.Credentials, policy v1beta1.Policy) (scanners []vulnerability.Scanner) {
	args := msf.Called(img, credentials, policy)
	return args.Get(0).([]vulnerability.Scanner)
}

type mockScanner struct {
	mock.Mock
}

func (ms *mockScanner) CanImageDeployBasedOnVulnerabilities(img image.Reference) (vulnerability.ScanResponse, error) {
	args := ms.Called(img)
	return args.Get(0).(vulnerability.ScanResponse), args.Error(1)
}

type mockNotaryVerifier struct {
	mock.Mock
}

func (mnv *mockNotaryVerifier) VerifyByPolicy(namespace string, img *image.Reference, credentials credential.Credentials, policy *v1beta1.Policy) (*bytes.Buffer, error, error) {
	args := mnv.Called(namespace, img, credentials, policy)
	return args.Get(0).(*bytes.Buffer), args.Error(1), args.Error(2)
}

type mockSimpleVerifier struct {
	mock.Mock
}

func (msv *mockSimpleVerifier) TransformPolicies(kWrapper kubernetes.WrapperInterface, namespace string, inPolicies []v1beta1.SimpleRequirement) (*signature.Policy, error) {
	args := msv.Called(kWrapper, namespace, inPolicies)
	return args.Get(0).(*signature.Policy), args.Error(1)
}

func (msv *mockSimpleVerifier) CreateRegistryDir(storeURL, storeUser, storePassword string) (string, error) {
	args := msv.Called(storeURL, storeUser, storePassword)
	return args.String(0), args.Error(1)
}

func (msv *mockSimpleVerifier) RemoveRegistryDir(dirName string) error {
	args := msv.Called(dirName)
	return args.Error(0)
}

func (msv *mockSimpleVerifier) VerifyByPolicy(imageToVerify string, credentials credential.Credentials, registriesConfigDir string, simplePolicy *signature.Policy) (*bytes.Buffer, error, error) {
	args := msv.Called(imageToVerify, credentials, registriesConfigDir, simplePolicy)
	return args.Get(0).(*bytes.Buffer), args.Error(1), args.Error(2)
}

func Test_enforcer_DigestByPolicy(t *testing.T) {
	type transformPoliciesMock struct {
		policy *signature.Policy
		err    error
	}
	type getBasicCredentialsMock struct {
		storeUser     string
		storePassword string
		err           error
	}
	type createRegistryDirMock struct {
		storeConfigDir string
		err            error
	}
	type simpleVerifyByPolicyMock struct {
		digest string
		deny   error
		err    error
	}
	type removeRegistryDirMock struct {
		err error
	}
	tests := []struct {
		name                 string
		namespace            string
		imageName            string
		credentials          credential.Credentials
		policy               *v1beta1.Policy
		transformPolicies    *transformPoliciesMock
		getBasicCredentials  *getBasicCredentialsMock
		createRegistryDir    *createRegistryDirMock
		simpleVerifyByPolicy *simpleVerifyByPolicyMock
		removeRegistryDir    *removeRegistryDirMock
		wantDigest           string
		wantDeny             error
		wantErr              error
	}{
		{
			name:       "No policy, Allow and no mutation",
			namespace:  "wibble",
			imageName:  "icr.io/wibble/some:tag",
			policy:     nil,
			wantDigest: "",
			wantDeny:   nil,
			wantErr:    nil,
		},
		{
			name:      "If TransformPolicies errors, return error",
			namespace: "wibble",
			imageName: "icr.io/wibble/some:tag",
			policy: &v1beta1.Policy{
				Simple: v1beta1.Simple{
					Requirements: []v1beta1.SimpleRequirement{
						{
							Type:      "test",
							KeySecret: "noOneCares",
						},
					},
					StoreURL:    "some.url.com",
					StoreSecret: "someSecret1234",
				},
			},
			transformPolicies: &transformPoliciesMock{
				err: fmt.Errorf("broken"),
			},
			wantDigest: "",
			wantDeny:   nil,
			wantErr:    fmt.Errorf("broken"),
		},
		{
			name:      "If GetBasicCredentials errors, return error",
			namespace: "wibble",
			imageName: "icr.io/wibble/some:tag",
			policy: &v1beta1.Policy{
				Simple: v1beta1.Simple{
					Requirements: []v1beta1.SimpleRequirement{
						{
							Type:      "test",
							KeySecret: "noOneCares",
						},
					},
					StoreURL:    "some.url.com",
					StoreSecret: "someSecret1234",
				},
			},
			transformPolicies: &transformPoliciesMock{},
			getBasicCredentials: &getBasicCredentialsMock{
				err: fmt.Errorf("also broken"),
			},
			wantDigest: "",
			wantDeny:   nil,
			wantErr:    fmt.Errorf("also broken"),
		},
		{
			name:      "If CreateRegistryDir errors, return error",
			namespace: "wibble",
			imageName: "icr.io/wibble/some:tag",
			policy: &v1beta1.Policy{
				Simple: v1beta1.Simple{
					Requirements: []v1beta1.SimpleRequirement{
						{
							Type:      "test",
							KeySecret: "noOneCares",
						},
					},
					StoreURL:    "some.url.com",
					StoreSecret: "someSecret1234",
				},
			},
			transformPolicies: &transformPoliciesMock{},
			getBasicCredentials: &getBasicCredentialsMock{
				storeUser:     "sillyUser",
				storePassword: "password1234",
			},
			createRegistryDir: &createRegistryDirMock{
				err: fmt.Errorf("busted"),
			},
			wantDigest: "",
			wantDeny:   nil,
			wantErr:    fmt.Errorf("busted"),
		},
		{
			name:      "If simple signing VerifyByPolicy errors, return error",
			namespace: "wibble",
			imageName: "icr.io/wibble/some:tag",
			policy: &v1beta1.Policy{
				Simple: v1beta1.Simple{
					Requirements: []v1beta1.SimpleRequirement{
						{
							Type:      "test",
							KeySecret: "noOneCares",
						},
					},
					StoreURL:    "some.url.com",
					StoreSecret: "someSecret1234",
				},
			},
			transformPolicies: &transformPoliciesMock{},
			getBasicCredentials: &getBasicCredentialsMock{
				storeUser:     "sillyUser",
				storePassword: "password1234",
			},
			createRegistryDir: &createRegistryDirMock{
				storeConfigDir: "vault",
			},
			simpleVerifyByPolicy: &simpleVerifyByPolicyMock{
				err: fmt.Errorf("still bust"),
			},
			wantDigest: "",
			wantDeny:   nil,
			wantErr:    fmt.Errorf("simple: still bust"),
		},
		{
			name:      "If simple signing VerifyByPolicy says deny, deny",
			namespace: "wibble",
			imageName: "icr.io/wibble/some:tag",
			policy: &v1beta1.Policy{
				Simple: v1beta1.Simple{
					Requirements: []v1beta1.SimpleRequirement{
						{
							Type:      "test",
							KeySecret: "noOneCares",
						},
					},
					StoreURL:    "some.url.com",
					StoreSecret: "someSecret1234",
				},
			},
			transformPolicies: &transformPoliciesMock{},
			getBasicCredentials: &getBasicCredentialsMock{
				storeUser:     "sillyUser",
				storePassword: "password1234",
			},
			createRegistryDir: &createRegistryDirMock{
				storeConfigDir: "vault",
			},
			simpleVerifyByPolicy: &simpleVerifyByPolicyMock{
				deny: fmt.Errorf("not allowed"),
			},
			removeRegistryDir: &removeRegistryDirMock{},
			wantDigest:        "",
			wantDeny:          fmt.Errorf("simple: policy denied the request: not allowed"),
			wantErr:           nil,
		},
		{
			name:      "Allow access if simple signing is allowed, and no trust policy",
			namespace: "wibble",
			imageName: "icr.io/wibble/some:tag",
			policy: &v1beta1.Policy{
				Simple: v1beta1.Simple{
					Requirements: []v1beta1.SimpleRequirement{
						{
							Type:      "test",
							KeySecret: "noOneCares",
						},
					},
					StoreURL:    "some.url.com",
					StoreSecret: "someSecret1234",
				},
			},
			transformPolicies: &transformPoliciesMock{},
			getBasicCredentials: &getBasicCredentialsMock{
				storeUser:     "sillyUser",
				storePassword: "password1234",
			},
			createRegistryDir: &createRegistryDirMock{
				storeConfigDir: "vault",
			},
			simpleVerifyByPolicy: &simpleVerifyByPolicyMock{
				digest: "sha256@sdfghjkj",
			},
			removeRegistryDir: &removeRegistryDirMock{},
			wantDigest:        "sha256@sdfghjkj",
			wantDeny:          nil,
			wantErr:           nil,
		},
		{
			name:      "If RemoveRegistryDir errors, log it and carry on as normal",
			namespace: "wibble",
			imageName: "icr.io/wibble/some:tag",
			policy: &v1beta1.Policy{
				Simple: v1beta1.Simple{
					Requirements: []v1beta1.SimpleRequirement{
						{
							Type:      "test",
							KeySecret: "noOneCares",
						},
					},
					StoreURL:    "some.url.com",
					StoreSecret: "someSecret1234",
				},
			},
			transformPolicies: &transformPoliciesMock{},
			getBasicCredentials: &getBasicCredentialsMock{
				storeUser:     "sillyUser",
				storePassword: "password1234",
			},
			createRegistryDir: &createRegistryDirMock{
				storeConfigDir: "vault",
			},
			simpleVerifyByPolicy: &simpleVerifyByPolicyMock{
				digest: "sha256@sdfghjkj",
			},
			removeRegistryDir: &removeRegistryDirMock{
				err: fmt.Errorf("whoope, no one cares"),
			},
			wantDigest: "sha256@sdfghjkj",
			wantDeny:   nil,
			wantErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := image.NewReference(tt.imageName)
			require.NoError(t, err)

			kubeWrapper := mockKubeWrapper{}
			kubeWrapper.Test(t)
			defer kubeWrapper.AssertExpectations(t)
			if tt.getBasicCredentials != nil {
				require.NotNil(t, tt.policy)
				kubeWrapper.
					On("GetBasicCredentials", tt.namespace, tt.policy.Simple.StoreSecret).
					Return(tt.getBasicCredentials.storeUser, tt.getBasicCredentials.storePassword, tt.getBasicCredentials.err).
					Once()
			}

			notaryVerfier := mockNotaryVerifier{}
			notaryVerfier.Test(t)
			defer notaryVerfier.AssertExpectations(t)

			simpleVerifier := mockSimpleVerifier{}
			simpleVerifier.Test(t)
			defer simpleVerifier.AssertExpectations(t)
			if tt.transformPolicies != nil {
				require.NotNil(t, tt.policy)
				simpleVerifier.
					On("TransformPolicies", &kubeWrapper, tt.namespace, tt.policy.Simple.Requirements).
					Return(tt.transformPolicies.policy, tt.transformPolicies.err).
					Once()
			}
			if tt.createRegistryDir != nil {
				require.NotNil(t, tt.policy)
				require.NotNil(t, tt.getBasicCredentials)
				user := tt.getBasicCredentials.storeUser
				pass := tt.getBasicCredentials.storePassword
				simpleVerifier.
					On("CreateRegistryDir", tt.policy.Simple.StoreURL, user, pass).
					Return(tt.createRegistryDir.storeConfigDir, tt.createRegistryDir.err).
					Once()
			}
			if tt.simpleVerifyByPolicy != nil {
				require.NotNil(t, tt.transformPolicies)
				inPolicy := tt.transformPolicies.policy
				require.NotNil(t, tt.createRegistryDir)
				inConfigDir := tt.createRegistryDir.storeConfigDir
				digest := bytes.NewBuffer([]byte(tt.simpleVerifyByPolicy.digest))
				simpleVerifier.
					On("VerifyByPolicy", tt.imageName, tt.credentials, inConfigDir, inPolicy).
					Return(digest, tt.simpleVerifyByPolicy.deny, tt.simpleVerifyByPolicy.err).
					Once()
			}
			if tt.removeRegistryDir != nil {
				require.NotNil(t, tt.createRegistryDir)
				simpleVerifier.
					On("RemoveRegistryDir", tt.createRegistryDir.storeConfigDir).
					Return(tt.removeRegistryDir.err).
					Once()
			}

			e := enforcer{
				kubeClientsetWrapper: &kubeWrapper,
				nv:                   &notaryVerfier,
				sv:                   &simpleVerifier,
			}

			gotDigest, gotDeny, gotErr := e.DigestByPolicy(tt.namespace, img, tt.credentials, tt.policy)

			if tt.wantDigest != "" {
				wantDigest := bytes.NewBuffer([]byte(tt.simpleVerifyByPolicy.digest))
				assert.Equal(t, wantDigest, gotDigest)
			}
			assert.Equal(t, tt.wantDeny, gotDeny)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func Test_enforcer_VulnerabilityPolicy(t *testing.T) {
	type canImageDeployBasedOnVulnerabilitiesMock struct {
		response vulnerability.ScanResponse
		err      error
	}
	tests := []struct {
		name         string
		imageName    string
		credentials  credential.Credentials
		policy       *v1beta1.Policy
		scanners     []canImageDeployBasedOnVulnerabilitiesMock
		wantResponse vulnerability.ScanResponse
	}{
		{
			name:         "If policy is nil, allow deploy",
			imageName:    "icr.io/nspc/some:thing",
			policy:       nil,
			wantResponse: vulnerability.ScanResponse{CanDeploy: true},
		},
		{
			name:         "If there are no scanners for the policy, allow deploy",
			imageName:    "icr.io/nspc/some:thing",
			policy:       &v1beta1.Policy{},
			scanners:     []canImageDeployBasedOnVulnerabilitiesMock{},
			wantResponse: vulnerability.ScanResponse{CanDeploy: true},
		},
		{
			name:      "If CanImageDeployBasedOnVulnerabilities errors, deny access",
			imageName: "icr.io/nspc/some:thing",
			policy:    &v1beta1.Policy{},
			scanners: []canImageDeployBasedOnVulnerabilitiesMock{
				{
					err: fmt.Errorf("something's broken something's broken"),
				},
			},
			wantResponse: vulnerability.ScanResponse{
				CanDeploy:  false,
				DenyReason: "something's broken something's broken",
			},
		},
		{
			name:      "If CanImageDeployBasedOnVulnerabilities says denied, deny access",
			imageName: "icr.io/nspc/some:thing",
			policy:    &v1beta1.Policy{},
			scanners: []canImageDeployBasedOnVulnerabilitiesMock{
				{
					response: vulnerability.ScanResponse{
						CanDeploy:  false,
						DenyReason: "because",
					},
				},
			},
			wantResponse: vulnerability.ScanResponse{
				CanDeploy:  false,
				DenyReason: "because",
			},
		},
		{
			name:      "First scanner says yes, but second says no",
			imageName: "icr.io/nspc/some:thing",
			policy:    &v1beta1.Policy{},
			scanners: []canImageDeployBasedOnVulnerabilitiesMock{
				{
					response: vulnerability.ScanResponse{
						CanDeploy: true,
					},
				},
				{
					response: vulnerability.ScanResponse{
						CanDeploy:  false,
						DenyReason: "because",
					},
				},
			},
			wantResponse: vulnerability.ScanResponse{
				CanDeploy:  false,
				DenyReason: "because",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := image.NewReference(tt.imageName)
			require.NoError(t, err)

			scanners := []vulnerability.Scanner{}
			for _, scannerResponse := range tt.scanners {
				scanner := mockScanner{}
				scanner.Test(t)
				defer scanner.AssertExpectations(t)
				scanner.
					On("CanImageDeployBasedOnVulnerabilities", *img).
					Return(scannerResponse.response, scannerResponse.err).
					Once()

				scanners = append(scanners, &scanner)
			}

			scannerFactory := mockScannerFactory{}
			scannerFactory.Test(t)
			defer scannerFactory.AssertExpectations(t)

			if tt.policy != nil {
				scannerFactory.
					On("GetScanners", *img, tt.credentials, *tt.policy).
					Return(scanners).
					Once()
			}

			e := enforcer{
				scannerFactory: &scannerFactory,
			}

			gotResponse := e.VulnerabilityPolicy(img, tt.credentials, tt.policy)

			assert.Equal(t, tt.wantResponse, gotResponse)
		})
	}
}
