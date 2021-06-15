// Copyright 2018, 2021 Portieris Authors.
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

package kubernetes

import (
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
)

var codec = serializer.NewCodecFactory(runtime.NewScheme())

var _ WrapperInterface = &Wrapper{}

// WrapperInterface is the interface for a wrapper around kubeclientset that includes some helper functions for applying behaviour to kube resources
type WrapperInterface interface {
	kubernetes.Interface
	GetPodSpec(*admissionv1.AdmissionRequest) (string, *corev1.PodSpec, error)
	GetSecretToken(namespace, secretName, registry string) (string, string, error)
	GetSecretKey(namespace, secretName string) ([]byte, error)
	GetBasicCredentials(namespace, secretName string) (string, string, error)
}

// Wrapper is a wrapper around kubeclientset that includes some helper functions for applying behaviour to kube resources
type Wrapper struct {
	kubernetes.Interface
}

// NewKubeClientsetWrapper creates a wrapper from the kubeclientset passed in
func NewKubeClientsetWrapper(kubeClientset kubernetes.Interface) *Wrapper {
	return &Wrapper{kubeClientset}
}
