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
	"github.com/IBM/portieris/pkg/metrics"
	"github.com/IBM/portieris/pkg/verifier/simple"
	notaryverifier "github.com/IBM/portieris/pkg/verifier/trust"
	"github.com/IBM/portieris/pkg/verifier/vulnerability"
	"github.com/IBM/portieris/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	k8sv1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type mockPolicyClient struct {
	mock.Mock
}

func (mpc *mockPolicyClient) GetPolicyToEnforce(namespace, image string) (*v1beta1.Policy, error) {
	args := mpc.Called(namespace, image)
	return args.Get(0).(*v1beta1.Policy), args.Error(1)
}

type mockKubeWrapper struct {
	mock.Mock
	kubernetes.Interface
}

func (mkw *mockKubeWrapper) GetPodSpec(req *k8sv1beta1.AdmissionRequest) (string, *corev1.PodSpec, error) {
	args := mkw.Called(req)
	return args.String(0), args.Get(1).(*corev1.PodSpec), args.Error(2)
}

func (mkw *mockKubeWrapper) GetSecretToken(namespace, secretName, registry string) (string, string, error) {
	args := mkw.Called(namespace, secretName, registry)
	return args.String(0), args.String(1), args.Error(2)
}

func (mkw *mockKubeWrapper) GetSecretKey(namespace, secretName string) ([]byte, error) {
	args := mkw.Called(namespace, secretName)
	return args.Get(0).([]byte), args.Error(1)
}

func (mkw *mockKubeWrapper) GetBasicCredentials(namespace, secretName string) (string, string, error) {
	args := mkw.Called(namespace, secretName)
	return args.String(0), args.String(1), args.Error(2)
}

type mockEnforcer struct {
	mock.Mock
}

func (me *mockEnforcer) DigestByPolicy(namespace string, img *image.Reference, credentials credential.Credentials, policy *v1beta1.Policy) (*bytes.Buffer, error, error) {
	args := me.Called(namespace, img, credentials, policy)
	return args.Get(0).(*bytes.Buffer), args.Error(1), args.Error(2)
}

func (me *mockEnforcer) VulnerabilityPolicy(img *image.Reference, credentials credential.Credentials, policy *v1beta1.Policy) vulnerability.ScanResponse {
	args := me.Called(img, credentials, policy)
	return args.Get(0).(vulnerability.ScanResponse)
}

func TestNewController(t *testing.T) {
	wantKubeWrapper := &mockKubeWrapper{}
	wantPolicyClient := &mockPolicyClient{}
	wantNV := &notaryverifier.Verifier{}
	wantScannerFactory := vulnerability.NewScannerFactory()
	wantEnforcer := &enforcer{
		kubeClientsetWrapper: wantKubeWrapper,
		nv:                   wantNV,
		scannerFactory:       &wantScannerFactory,
		sv:                   simple.NewVerifier(),
	}
	wantMetrics := metrics.NewMetrics()
	defer wantMetrics.UnregisterAll()

	wantController := Controller{
		kubeClientsetWrapper: wantKubeWrapper,
		policyClient:         wantPolicyClient,
		Enforcer:             wantEnforcer,
		PMetrics:             wantMetrics,
	}

	gotController := NewController(wantKubeWrapper, wantPolicyClient, wantNV, wantMetrics)

	assert.Equal(t, wantController, *gotController)
}

func TestController_getPatchesForContainers(t *testing.T) {
	type getPolicyToEnforceMock struct {
		outPolicy *v1beta1.Policy
		outErr    error
	}
	type enforcerVulnerabilityPolicyMock struct {
		outScanResponse vulnerability.ScanResponse
	}
	type enforceDigestByPolicyMock struct {
		outDigest string
		outDeny   error
		outErr    error
	}
	type mocks struct {
		getPolicyToEnforce          *getPolicyToEnforceMock
		inImage                     string
		credentials                 credential.Credentials
		enforcerVulnerabilityPolicy *enforcerVulnerabilityPolicyMock
		enforceDigestByPolicy       *enforceDigestByPolicyMock
	}
	tests := []struct {
		name             string
		containerType    string
		namespace        string
		specPath         string
		imagePullSecrets []corev1.LocalObjectReference
		containers       []corev1.Container
		mocks            []mocks
		wantPatches      []types.JSONPatch
		wantDenials      map[string][]string
		wantErr          error
	}{
		{
			name:        "No containers, return no patches or denials",
			containers:  []corev1.Container{},
			wantPatches: []types.JSONPatch{},
			wantDenials: map[string][]string{},
			wantErr:     nil,
		},
		{
			name: "Invalid image name in container, deny",
			containers: []corev1.Container{
				{Image: "Invalid&Image%Name"},
			},
			wantPatches: []types.JSONPatch{},
			wantDenials: map[string][]string{
				"invalidimagename": {"Deny \"Invalid&Image%Name\", invalid image name"},
			},
			wantErr: nil,
		},
		{
			name:      "Fail to get policy, deny",
			namespace: "some-namespace",
			containers: []corev1.Container{
				{Image: "icr.io/some-namespace/image:tag"},
			},
			mocks: []mocks{
				{
					inImage: "icr.io/some-namespace/image:tag",
					getPolicyToEnforce: &getPolicyToEnforceMock{
						outErr: fmt.Errorf("no sorry"),
					},
				},
			},
			wantPatches: []types.JSONPatch{},
			wantDenials: map[string][]string{
				"icr.io/some-namespace/image:tag": {"no sorry"},
			},
			wantErr: nil,
		},
		{
			name:      "Fail to get policy, deny, multiple containers",
			namespace: "some-namespace",
			containers: []corev1.Container{
				{
					Image: "icr.io/some-namespace/image:tag",
				},
				{
					Image: "icr.io/some-namespace/anotherimage:tag@sha256:def",
				},
			},
			mocks: []mocks{
				{
					inImage: "icr.io/some-namespace/image:tag",
					getPolicyToEnforce: &getPolicyToEnforceMock{
						outErr: fmt.Errorf("no sorry"),
					},
				},
				{
					inImage: "icr.io/some-namespace/anotherimage:tag@sha256:def",
					getPolicyToEnforce: &getPolicyToEnforceMock{
						outErr: fmt.Errorf("no sorry"),
					},
				},
			},
			wantPatches: []types.JSONPatch{},
			wantDenials: map[string][]string{
				"icr.io/some-namespace/image:tag":        {"no sorry"},
				"icr.io/some-namespace/anotherimage:def": {"no sorry"},
			},
			wantErr: nil,
		},
		{
			name:      "Not vulnerable with no signing enforcement",
			namespace: "some-namespace",
			imagePullSecrets: []corev1.LocalObjectReference{
				{
					Name: "wibble",
				},
			},
			containers: []corev1.Container{
				{Image: "icr.io/some-namespace/image:tag"},
			},
			mocks: []mocks{
				{
					inImage: "icr.io/some-namespace/image:tag",
					getPolicyToEnforce: &getPolicyToEnforceMock{
						outPolicy: &v1beta1.Policy{},
					},
					credentials: credential.Credentials{{
						Username: "wibble",
						Password: "dibble",
					}},
					enforcerVulnerabilityPolicy: &enforcerVulnerabilityPolicyMock{
						outScanResponse: vulnerability.ScanResponse{
							CanDeploy: true,
						},
					},
					enforceDigestByPolicy: &enforceDigestByPolicyMock{},
				},
			},
			wantPatches: []types.JSONPatch{},
			wantDenials: map[string][]string{
				"icr.io/some-namespace/image:tag": {},
			},
			wantErr: nil,
		},
		{
			name:      "digest by policy errors",
			namespace: "some-namespace",
			imagePullSecrets: []corev1.LocalObjectReference{
				{
					Name: "wibble",
				},
			},
			containers: []corev1.Container{
				{Image: "icr.io/some-namespace/image:tag"},
			},
			mocks: []mocks{
				{
					inImage: "icr.io/some-namespace/image:tag",
					getPolicyToEnforce: &getPolicyToEnforceMock{
						outPolicy: &v1beta1.Policy{},
					},
					credentials: credential.Credentials{{
						Username: "wibble",
						Password: "dibble",
					}},
					enforcerVulnerabilityPolicy: &enforcerVulnerabilityPolicyMock{
						outScanResponse: vulnerability.ScanResponse{
							CanDeploy: true,
						},
					},
					enforceDigestByPolicy: &enforceDigestByPolicyMock{
						outErr: fmt.Errorf("failed"),
					},
				},
			},
			wantPatches: []types.JSONPatch{},
			wantDenials: map[string][]string{
				"icr.io/some-namespace/image:tag": {},
			},
			wantErr: fmt.Errorf("failed"),
		},
		{
			name:      "digest by policy says denied",
			namespace: "some-namespace",
			imagePullSecrets: []corev1.LocalObjectReference{
				{
					Name: "wibble",
				},
			},
			containers: []corev1.Container{
				{Image: "icr.io/some-namespace/image:tag"},
			},
			mocks: []mocks{
				{
					inImage: "icr.io/some-namespace/image:tag",
					getPolicyToEnforce: &getPolicyToEnforceMock{
						outPolicy: &v1beta1.Policy{},
					},
					credentials: credential.Credentials{{
						Username: "wibble",
						Password: "dibble",
					}},
					enforcerVulnerabilityPolicy: &enforcerVulnerabilityPolicyMock{
						outScanResponse: vulnerability.ScanResponse{
							CanDeploy: true,
						},
					},
					enforceDigestByPolicy: &enforceDigestByPolicyMock{
						outDeny: fmt.Errorf("I don't think so"),
					},
				},
			},
			wantPatches: []types.JSONPatch{},
			wantDenials: map[string][]string{
				"icr.io/some-namespace/image:tag": {"I don't think so"},
			},
			wantErr: nil,
		},
		{
			name:      "digest by policy returns a digest",
			namespace: "some-namespace",
			imagePullSecrets: []corev1.LocalObjectReference{
				{
					Name: "wibble",
				},
			},
			containers: []corev1.Container{
				{Image: "icr.io/some-namespace/image:tag"},
			},
			mocks: []mocks{
				{
					inImage: "icr.io/some-namespace/image:tag",
					getPolicyToEnforce: &getPolicyToEnforceMock{
						outPolicy: &v1beta1.Policy{},
					},
					credentials: credential.Credentials{{
						Username: "wibble",
						Password: "dibble",
					}},
					enforcerVulnerabilityPolicy: &enforcerVulnerabilityPolicyMock{
						outScanResponse: vulnerability.ScanResponse{
							CanDeploy: true,
						},
					},
					enforceDigestByPolicy: &enforceDigestByPolicyMock{
						outDigest: "somedigest",
					},
				},
			},
			wantPatches: []types.JSONPatch{},
			wantDenials: map[string][]string{
				"icr.io/some-namespace/image:somedigest": {},
			},
			wantErr: nil,
		},
		{
			name:      "allowed by signing but vulnerable",
			namespace: "some-namespace",
			imagePullSecrets: []corev1.LocalObjectReference{
				{
					Name: "wibble",
				},
			},
			containers: []corev1.Container{
				{Image: "icr.io/some-namespace/image:tag"},
			},
			mocks: []mocks{
				{
					inImage: "icr.io/some-namespace/image:tag",
					getPolicyToEnforce: &getPolicyToEnforceMock{
						outPolicy: &v1beta1.Policy{},
					},
					credentials: credential.Credentials{{
						Username: "wibble",
						Password: "dibble",
					}},
					enforcerVulnerabilityPolicy: &enforcerVulnerabilityPolicyMock{
						outScanResponse: vulnerability.ScanResponse{
							CanDeploy:  false,
							DenyReason: "I don't want to",
						},
					},
					enforceDigestByPolicy: &enforceDigestByPolicyMock{
						outDigest: "somedigest",
					},
				},
			},
			wantPatches: []types.JSONPatch{},
			wantDenials: map[string][]string{
				"icr.io/some-namespace/image:somedigest": {"I don't want to"},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			podSpec := corev1.PodSpec{
				ImagePullSecrets: tt.imagePullSecrets,
			}

			policyClient := mockPolicyClient{}
			policyClient.Test(t)
			defer policyClient.AssertExpectations(t)

			kubeWrapper := mockKubeWrapper{}
			kubeWrapper.Test(t)
			defer kubeWrapper.AssertExpectations(t)

			enforcer := mockEnforcer{}
			enforcer.Test(t)
			defer enforcer.AssertExpectations(t)

			for idx, m := range tt.mocks {
				imageName := m.inImage
				var img *image.Reference
				if imageName != "" {
					var err error
					img, err = image.NewReference(imageName)
					require.NoError(t, err)
				}
				policy := m.getPolicyToEnforce.outPolicy
				namespace := tt.namespace
				if m.getPolicyToEnforce != nil {
					err := m.getPolicyToEnforce.outErr
					policyClient.
						On("GetPolicyToEnforce", namespace, img.String()).Return(policy, err).Once()
				}

				creds := m.credentials
				for idx, cred := range creds {
					secretName := tt.imagePullSecrets[idx].Name
					kubeWrapper.On("GetSecretToken", namespace, secretName, img.GetHostname()).
						Return(cred.Username, cred.Password, nil).Once()
				}

				if m.enforcerVulnerabilityPolicy != nil {
					response := m.enforcerVulnerabilityPolicy.outScanResponse
					enforcer.On("VulnerabilityPolicy", img, creds, policy).Return(response).Once()
				}

				if m.enforceDigestByPolicy != nil {
					var digest *bytes.Buffer
					if m.enforceDigestByPolicy.outDigest != "" {
						digest = &bytes.Buffer{}
						digest.Write([]byte(m.enforceDigestByPolicy.outDigest))
						wantPatch := types.JSONPatch{
							Op:    "replace",
							Path:  fmt.Sprintf("%s/%s/%d/image", tt.specPath, tt.containerType, idx),
							Value: fmt.Sprintf("%s@sha256:%s", img.NameWithoutTag(), m.enforceDigestByPolicy.outDigest),
						}
						tt.wantPatches = append(tt.wantPatches, wantPatch)
					}
					deny := m.enforceDigestByPolicy.outDeny
					err := m.enforceDigestByPolicy.outErr
					enforcer.On("DigestByPolicy", namespace, img, creds, policy).Return(digest, deny, err).Once()
				}
			}

			c := &Controller{
				policyClient:         &policyClient,
				Enforcer:             &enforcer,
				kubeClientsetWrapper: &kubeWrapper,
				PMetrics:             metrics.NewMetrics(),
			}
			defer c.PMetrics.UnregisterAll()

			gotPatches, gotDenials, gotErr := c.getPatchesForContainers(tt.containerType, tt.namespace, tt.specPath, podSpec, tt.containers)

			assert.Equal(t, tt.wantPatches, gotPatches)
			assert.Equal(t, tt.wantDenials, gotDenials)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
