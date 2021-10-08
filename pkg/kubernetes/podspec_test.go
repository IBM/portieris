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
	"testing"

	"github.com/stretchr/testify/assert"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestWrapper_GetPodSpec(t *testing.T) {
	nginxSpec := &corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:  "nginx",
				Image: "docker.io/nginx",
			},
		},
	}

	type ar struct {
		Object    []byte
		Namespace string
		Resource  metav1.GroupVersionResource
	}
	tests := []struct {
		name    string
		ar      ar
		want    string
		want1   *corev1.PodSpec
		wantErr bool
	}{
		{
			name: "Properly handles a pod",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"},
			},
			want:  "/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a malformed pod",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"},
			},
			want:    "/spec/template/spec",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles a replicationcontroller",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":1, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "replicationcontrollers"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Returns the spec for a replicationcontroller with null replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "replicationcontrollers"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a replicationcontroller with zero replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":0, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "replicationcontrollers"},
			},
			wantErr: true,
		},
		{
			name: "Errors for a malformed replicationcontroller",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "replicationcontrollers"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles a deployment",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":1, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Properly handles a deployment with null replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a deployment with zero replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":0, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"},
			},
			wantErr: true,
		},
		{
			name: "Errors for a malformed deployment",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles a beta deployment",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":1, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "deployments"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Properly handles a beta2 deployment",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":1, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "deployments"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a malformed beta2 deployment",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "deployments"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles a beta2 deployment with null replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "deployments"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a beta2 deployment with zero replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":0, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "deployments"},
			},
			wantErr: true,
		},
		{
			name: "Properly handles a beta deployment with null replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "deployments"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a beta deployment with zero replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":0, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "deployments"},
			},
			wantErr: true,
		},
		{
			name: "Errors for a malformed beta deployment",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "deployments"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles an apps beta1 deployment",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":1, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta1", Resource: "deployments"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a malformed an apps beta1 deployment",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta1", Resource: "deployments"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles an apps beta1 deployment with null replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta1", Resource: "deployments"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors foran apps beta1 deployment with zero replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":0, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta1", Resource: "deployments"},
			},
			wantErr: true,
		},
		{
			name: "Properly handles a replicaset",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":1, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Properly handles a replicaset with null replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a replicaset with zero replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":0, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"},
			},
			wantErr: true,
		},
		{
			name: "Errors for a malformed replicaset",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles a beta replicaset",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":1, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "replicasets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Properly handles a beta replicaset with null replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "replicasets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a beta replicaset with zero replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":0, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "replicasets"},
			},
			wantErr: true,
		},
		{
			name: "Errors for a malformed beta replicaset",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "replicasets"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles an apps/v1beta2 replicaset",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":1, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "replicasets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Properly handles an apps/v1beta2 replicaset with null replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "replicasets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for an apps/v1beta2 replicaset with zero replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":0, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "replicasets"},
			},
			wantErr: true,
		},
		{
			name: "Errors for an apps/v1beta2 malformed replicaset",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "replicasets"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles a daemonset",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Properly handles beta daemonset",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "daemonsets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a malformed daemonset",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Errors for a malformed beta daemonset",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "daemonsets"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles an apps/v1beta2 daemonset",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "daemonsets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for an apps/v1beta2 malformed daemonset",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "daemonsets"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles a statefulset",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":1, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Properly handles a statefulset with null replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a statefulset with zero replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":0, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"},
			},
			wantErr: true,
		},
		{
			name: "Errors for a malformed statefulset",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles an apps/v1beta1 statefulset",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":1, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta1", Resource: "statefulsets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Properly handles an apps/v1beta1 statefulset with null replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta1", Resource: "statefulsets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for an apps/v1beta1 statefulset with zero replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":0, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta1", Resource: "statefulsets"},
			},
			wantErr: true,
		},
		{
			name: "Errors for an apps/v1beta1 malformed statefulset",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta1", Resource: "statefulsets"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles an apps/v1beta2 statefulset",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":1, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "statefulsets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Properly handles an apps/v1beta2 statefulset with null replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "statefulsets"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for an apps/v1beta2 statefulset with zero replicas",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"replicas":0, "template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "statefulsets"},
			},
			wantErr: true,
		},
		{
			name: "Errors for an apps/v1beta2 malformed statefulset",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "apps", Version: "v1beta2", Resource: "statefulsets"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles a job",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"},
			},
			want:  "/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a malformed job",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles a v1 cronjob",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"jobTemplate":{"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "batch", Version: "v1", Resource: "cronjobs"},
			},
			want:  "/spec/jobTemplate/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a malformed v1 cronjob",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "batch", Version: "v1", Resource: "cronjobs"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Properly handles a v1beta1 cronjob",
			ar: ar{
				Object:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"jobTemplate":{"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}}}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "batch", Version: "v1beta1", Resource: "cronjobs"},
			},
			want:  "/spec/jobTemplate/spec/template/spec",
			want1: nginxSpec,
		},
		{
			name: "Errors for a malformed v1beta1 cronjob",
			ar: ar{
				Object:    []byte(`lololol`),
				Namespace: "default",
				Resource:  metav1.GroupVersionResource{Group: "batch", Version: "v1beta1", Resource: "cronjobs"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "Errors for an unsupported type",
			ar: ar{
				Object:    []byte(`{}`),
				Namespace: "",
				Resource:  metav1.GroupVersionResource{Group: "wibble", Version: "bibble", Resource: "tibble"},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ar := &admissionv1.AdmissionRequest{
				Resource:  tt.ar.Resource,
				Namespace: tt.ar.Namespace,
				Object: runtime.RawExtension{
					Raw: tt.ar.Object,
				},
			}

			w := &Wrapper{}

			got, got1, err := w.GetPodSpec(ar)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
					assert.Equal(t, tt.want1, got1)
				}
			}
		})
	}
}

func TestWrapper_mutateWithSA(t *testing.T) {
	serviceaccounts := []runtime.Object{
		&corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "default",
				Namespace: "default",
			},
			ImagePullSecrets: []corev1.LocalObjectReference{
				{
					Name: "wibble",
				},
			},
		}, &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "myamazingserviceaccount",
				Namespace: "default",
			},
			ImagePullSecrets: []corev1.LocalObjectReference{
				{
					Name: "dibble",
				},
			},
		},
	}
	tests := []struct {
		name    string
		ns      string
		ps      *corev1.PodSpec
		want    *corev1.PodSpec
		wantErr bool
	}{
		{
			name: "does not mutate if namespace is empty",
			ns:   "",
			ps:   &corev1.PodSpec{},
			want: &corev1.PodSpec{},
		},
		{
			name: "mutates a spec with secrets from the default serviceaccount",
			ns:   "default",
			ps:   &corev1.PodSpec{},
			want: &corev1.PodSpec{
				ImagePullSecrets: []corev1.LocalObjectReference{
					{
						Name: "wibble",
					},
				},
			},
		},
		{
			name: "uses non-default serviceaccounts if specified",
			ns:   "default",
			ps: &corev1.PodSpec{
				ServiceAccountName: "myamazingserviceaccount",
			},
			want: &corev1.PodSpec{
				ImagePullSecrets: []corev1.LocalObjectReference{
					{
						Name: "dibble",
					},
				},
				ServiceAccountName: "myamazingserviceaccount",
			},
		},
		{
			name: "errors without mutation if serviceaccount not found",
			ns:   "default",
			ps: &corev1.PodSpec{
				ServiceAccountName: "gibble",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kubeClientset := k8sfake.NewSimpleClientset(serviceaccounts...)
			w := NewKubeClientsetWrapper(kubeClientset)
			err := w.mutateWithSA(tt.ns, tt.ps)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, tt.ps)
				}
			}
		})
	}
}

func TestWrapper_decodeObject(t *testing.T) {
	pointerToBool := func(in bool) *bool {
		return &in
	}
	tests := []struct {
		name         string
		raw          []byte
		object       object
		want         object
		wantErr      bool
		wantErrEqual string
	}{
		{
			name:   "decodes a pod spec",
			raw:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}`),
			object: &corev1.Pod{},
			want: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "docker.io/nginx",
						},
					},
				},
			},
		},
		{
			name:   "decodes a deployment spec",
			raw:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
			object: &appsv1.Deployment{},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx",
					Namespace: "default",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "nginx",
									Image: "docker.io/nginx",
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "decodes a beta deployment spec",
			raw:    []byte(`{"metadata":{"name":"nginx","namespace":"default"},"spec":{"template":{"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}}}`),
			object: &extensionsv1beta1.Deployment{},
			want: &extensionsv1beta1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx",
					Namespace: "default",
				},
				Spec: extensionsv1beta1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "nginx",
									Image: "docker.io/nginx",
								},
							},
						},
					},
				},
			},
		},
		{
			name:         "returns object has parents if the object has parents",
			raw:          []byte(`{"metadata":{"name":"nginx","namespace":"default","ownerReferences":[{"apiVersion":"extensions/v1beta1","kind":"ReplicaSet","name":"deployment-55d687c698","uid":"e0577bcf-30dd-11e8-83d1-baaf52c27f02","controller":true,"blockOwnerDeletion":true}]},"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}`),
			object:       &corev1.Pod{},
			wantErr:      true,
			wantErrEqual: ErrObjectHasParents.Error(),
		},
		{
			name:   "decodes a pod spec when it's owner of a kind that we do not support",
			raw:    []byte(`{"metadata":{"name":"nginx","namespace":"default","ownerReferences":[{"apiVersion":"customcontroller.v1","kind":"CustomController","name":"customercontroller-55d687c698","uid":"e0577bcf-30dd-11e8-83d1-baaf52c27f02","controller":true,"blockOwnerDeletion":true}]},"spec":{"containers":[{"name":"nginx","image":"docker.io/nginx"}]}}`),
			object: &corev1.Pod{},
			want: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx",
					Namespace: "default",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "customcontroller.v1",
							Kind:               "CustomController",
							Name:               "customercontroller-55d687c698",
							UID:                "e0577bcf-30dd-11e8-83d1-baaf52c27f02",
							Controller:         pointerToBool(true),
							BlockOwnerDeletion: pointerToBool(true),
						},
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "docker.io/nginx",
						},
					},
				},
			},
		},
		{
			name:    "returns error if the object is weird",
			raw:     []byte(`lolololololol`),
			object:  &corev1.Pod{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wrapper{}
			err := w.decodeObject(tt.raw, tt.object)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrEqual != "" {
					assert.EqualError(t, err, tt.wantErrEqual)
				}
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, tt.object)
				}
			}
		})
	}
}
