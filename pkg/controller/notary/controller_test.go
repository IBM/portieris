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

package notary

import (
	"errors"
	"testing"

	securityenforcementv1beta1 "github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/IBM/portieris/pkg/kubernetes/fakewrapper"
	"github.com/IBM/portieris/pkg/policy/fakepolicy"
	"github.com/IBM/portieris/pkg/registry/fakeregistry"
	"github.com/stretchr/testify/assert"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

var fakeAdmissionRequest = &admissionv1beta1.AdmissionRequest{Operation: "blah", Name: "blah", Namespace: "blah"}

func TestController_Admit(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		podSpec     *corev1.PodSpec
		wantAllowed bool
	}{
		{
			name:        "calls through to mutatePodSpec",
			wantAllowed: false,
			podSpec:     &corev1.PodSpec{},
		},
		{
			name:        "Returns allowed if object has parents",
			err:         kubernetes.ErrObjectHasParents,
			wantAllowed: true,
		},
		{
			name:        "Returns allowed if object zero replicas",
			err:         kubernetes.ErrObjectHasParents,
			wantAllowed: true,
		},
		{
			name:        "Returns allowed if object zero replicas",
			err:         errors.New("generic error"),
			wantAllowed: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &fakewrapper.Stub{Err: tt.err, PodSpec: tt.podSpec}
			ctrl := &Controller{kubeClientsetWrapper: stub}
			got := ctrl.Admit(fakeAdmissionRequest)
			assert.Equal(t, tt.wantAllowed, got.Allowed)
		})
	}
}

func generatePodSpec(image string, imagePullSecrets bool) *corev1.PodSpec {
	ips := []corev1.LocalObjectReference{}
	if imagePullSecrets {
		ips = append(ips, corev1.LocalObjectReference{Name: "secret"})
	}
	return &corev1.PodSpec{
		ImagePullSecrets: ips,
		Containers: []corev1.Container{
			corev1.Container{
				Image: image,
			},
		},
	}
}

func generatePolicy(trustEnabled bool) *securityenforcementv1beta1.Policy {
	return &securityenforcementv1beta1.Policy{
		Trust: securityenforcementv1beta1.Trust{
			Enabled: &trustEnabled,
		},
	}
}

func TestController_mutatePodSpec(t *testing.T) {
	tests := []struct {
		name            string
		podSpec         *corev1.PodSpec
		kubeWrapperStub *fakewrapper.Stub
		policyStub      *fakepolicy.Stub
		registryStub    *fakeregistry.Stub
		wantAllowed     bool
		wantMessage     string
	}{
		{
			name:        "Deny if invalid image name",
			podSpec:     generatePodSpec("not-valid!!!!1111!!!", true),
			wantAllowed: false,
			wantMessage: "Deny \"not-valid!!!!1111!!!\", invalid image name",
		},
		{
			name:        "Deny if kube wrapper returns an error",
			podSpec:     generatePodSpec("registry.bluemix.net/liamwhite/test:test", true),
			policyStub:  &fakepolicy.Stub{Err: errors.New("bang")},
			wantAllowed: false,
			wantMessage: "bang",
		},
		{
			name:        "Allow if there is not a policy to enforce",
			podSpec:     generatePodSpec("registry.bluemix.net/liamwhite/test:test", true),
			policyStub:  &fakepolicy.Stub{},
			wantAllowed: true,
		},
		{
			name:        "Allow if policy has trust disabled",
			podSpec:     generatePodSpec("registry.bluemix.net/liamwhite/test:test", true),
			policyStub:  &fakepolicy.Stub{Policy: generatePolicy(false)},
			wantAllowed: true,
		},
		{
			name:        "Deny if there are not imagePullSecrets for the PodSpec when trust is enabled",
			podSpec:     generatePodSpec("registry.bluemix.net/liamwhite/test:test", false),
			policyStub:  &fakepolicy.Stub{Policy: generatePolicy(true)},
			wantAllowed: false,
			wantMessage: "no ImagePullSecret defined",
		},
		{
			name:            "Deny if unable to retrieve token from secret",
			podSpec:         generatePodSpec("registry.bluemix.net/liamwhite/test:test", true),
			policyStub:      &fakepolicy.Stub{Policy: generatePolicy(true)},
			kubeWrapperStub: &fakewrapper.Stub{Err: errors.New("bang")},
			wantAllowed:     false,
			wantMessage:     "no valid ImagePullSecret defined",
		},
		{
			name:            "Deny if unable to retrieve notary token",
			podSpec:         generatePodSpec("registry.bluemix.net/liamwhite/test:test", true),
			policyStub:      &fakepolicy.Stub{Policy: generatePolicy(true)},
			kubeWrapperStub: &fakewrapper.Stub{Token: "a-token"},
			registryStub:    &fakeregistry.Stub{Err: errors.New("bang")},
			wantAllowed:     false,
			wantMessage:     "no valid ImagePullSecret defined",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := NewController(tt.kubeWrapperStub, tt.policyStub, nil, tt.registryStub)
			got := ctrl.mutatePodSpec("", "", *tt.podSpec)
			assert.Equal(t, tt.wantAllowed, got.Allowed)
			if tt.wantMessage != "" {
				assert.Contains(t, got.Result.Message, tt.wantMessage)
			}
		})
	}
}
