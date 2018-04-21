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

package kubernetes

import (
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
)

var codec = serializer.NewCodecFactory(runtime.NewScheme())

var _ WrapperInterface = &Wrapper{}

type WrapperInterface interface {
	kubernetes.Interface
	GetPodSpec(*v1beta1.AdmissionRequest) (string, *corev1.PodSpec, error)
	GetSecretToken(namespace, secretName, registry string) (string, error)
}

type Wrapper struct {
	kubernetes.Interface
}

func NewKubeClientsetWrapper(kubeClientset kubernetes.Interface) *Wrapper {
	return &Wrapper{kubeClientset}
}
