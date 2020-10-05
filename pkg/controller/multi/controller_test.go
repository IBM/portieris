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
	"testing"

	"github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/IBM/portieris/pkg/verifier/simple"
	notaryverifier "github.com/IBM/portieris/pkg/verifier/trust"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestNewController(t *testing.T) {
	wantKubeWrapper := &mockKubeWrapper{}
	wantPolicyClient := &mockPolicyClient{}
	wantNV := &notaryverifier.Verifier{}
	wantEnforcer := &enforcer{
		kubeClientsetWrapper: wantKubeWrapper,
		nv:                   wantNV,
		sv:                   simple.NewVerifier(),
	}
	wantController := Controller{
		kubeClientsetWrapper: wantKubeWrapper,
		policyClient:         wantPolicyClient,
		Enforcer:             wantEnforcer,
	}

	gotController := NewController(wantKubeWrapper, wantPolicyClient, wantNV)

	assert.Equal(t, wantController, *gotController)
}
