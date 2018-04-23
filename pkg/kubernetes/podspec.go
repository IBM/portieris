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
	"fmt"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type object interface {
	metav1.Object
	runtime.Object
}

const (
	podSpecPath      = "/spec"
	templateSpecPath = "/spec/template/spec"
	cronJobSpecPath  = "/spec/jobTemplate/spec/template/spec"
)

// ErrObjectHasParents is returned when the resource being created is the child of another resource
var ErrObjectHasParents = fmt.Errorf("This object has parents")

// ErrObjectHasZeroReplicas is returned when the resource being created has zero replicas
var ErrObjectHasZeroReplicas = fmt.Errorf("This object has zero replicas")

// GetPodSpec retrieves the podspec from the admission request passed in
func (w *Wrapper) GetPodSpec(ar *v1beta1.AdmissionRequest) (string, *corev1.PodSpec, error) {
	ps := corev1.PodSpec{}
	var templateString string

	switch ar.Resource {
	case metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}:
		pod := corev1.Pod{}
		if err := w.decodeObject(ar.Object.Raw, &pod); err != nil {
			return "", nil, err
		}
		ps = pod.Spec
		templateString = podSpecPath
	case metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "replicationcontrollers"}:
		rc := corev1.ReplicationController{}
		if err := w.decodeObject(ar.Object.Raw, &rc); err != nil {
			return "", nil, err
		}
		if rc.Spec.Replicas != nil && *rc.Spec.Replicas == int32(0) {
			return "", nil, ErrObjectHasZeroReplicas
		}
		ps = rc.Spec.Template.Spec
		w.mutateWithSA(ar.Namespace, &ps)
		templateString = templateSpecPath
	case metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "deployments"}:
		deploy := extensionsv1beta1.Deployment{}
		if err := w.decodeObject(ar.Object.Raw, &deploy); err != nil {
			return "", nil, err
		}
		if deploy.Spec.Replicas != nil && *deploy.Spec.Replicas == int32(0) {
			return "", nil, ErrObjectHasZeroReplicas
		}
		ps = deploy.Spec.Template.Spec
		w.mutateWithSA(ar.Namespace, &ps)
		templateString = templateSpecPath
	case metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}:
		deploy := appsv1.Deployment{}
		if err := w.decodeObject(ar.Object.Raw, &deploy); err != nil {
			return "", nil, err
		}
		if deploy.Spec.Replicas != nil && *deploy.Spec.Replicas == int32(0) {
			return "", nil, ErrObjectHasZeroReplicas
		}
		ps = deploy.Spec.Template.Spec
		w.mutateWithSA(ar.Namespace, &ps)
		templateString = templateSpecPath
	case metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"}:
		rs := appsv1.ReplicaSet{}
		if err := w.decodeObject(ar.Object.Raw, &rs); err != nil {
			return "", nil, err
		}
		if rs.Spec.Replicas != nil && *rs.Spec.Replicas == int32(0) {
			return "", nil, ErrObjectHasZeroReplicas
		}
		ps = rs.Spec.Template.Spec
		w.mutateWithSA(ar.Namespace, &ps)
		templateString = templateSpecPath
	case metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "replicasets"}:
		rs := extensionsv1beta1.ReplicaSet{}
		if err := w.decodeObject(ar.Object.Raw, &rs); err != nil {
			return "", nil, err
		}
		if rs.Spec.Replicas != nil && *rs.Spec.Replicas == int32(0) {
			return "", nil, ErrObjectHasZeroReplicas
		}
		ps = rs.Spec.Template.Spec
		w.mutateWithSA(ar.Namespace, &ps)
		templateString = templateSpecPath
	case metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}:
		ds := appsv1.DaemonSet{}
		if err := w.decodeObject(ar.Object.Raw, &ds); err != nil {
			return "", nil, err
		}
		ps = ds.Spec.Template.Spec
		w.mutateWithSA(ar.Namespace, &ps)
		templateString = templateSpecPath
	case metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}:
		sts := appsv1.StatefulSet{}
		if err := w.decodeObject(ar.Object.Raw, &sts); err != nil {
			return "", nil, err
		}
		if sts.Spec.Replicas != nil && *sts.Spec.Replicas == int32(0) {
			return "", nil, ErrObjectHasZeroReplicas
		}
		ps = sts.Spec.Template.Spec
		w.mutateWithSA(ar.Namespace, &ps)
		templateString = templateSpecPath
	case metav1.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"}:
		job := batchv1.Job{}
		if err := w.decodeObject(ar.Object.Raw, &job); err != nil {
			return "", nil, err
		}
		ps = job.Spec.Template.Spec
		w.mutateWithSA(ar.Namespace, &ps)
		templateString = templateSpecPath
	case metav1.GroupVersionResource{Group: "batch", Version: "v1beta1", Resource: "cronjobs"}:
		job := batchv1beta1.CronJob{}
		if err := w.decodeObject(ar.Object.Raw, &job); err != nil {
			return "", nil, err
		}
		ps = job.Spec.JobTemplate.Spec.Template.Spec //:sob:
		w.mutateWithSA(ar.Namespace, &ps)
		templateString = cronJobSpecPath
	default:
		glog.Errorf("Resource not supported: %+v", ar.Resource)
		return "", nil, fmt.Errorf(`The resource "%s/%s/%s" is not supported. Make sure that you are using a supported kubectl version, and that you are using a supported Kubernetes workload type`, ar.Resource.Group, ar.Resource.Version, ar.Resource.Resource)
	}

	return templateString, &ps, nil
}

func (w *Wrapper) decodeObject(raw []byte, object object) error {
	deserializer := codec.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, object); err != nil {
		return err
	}
	if len(object.GetOwnerReferences()) != 0 {
		return ErrObjectHasParents
	}
	return nil
}

func (w *Wrapper) mutateWithSA(ns string, ps *corev1.PodSpec) error {
	if ns == "" || ps == nil || len(ps.ImagePullSecrets) != 0 {
		// Do nothing
		return nil
	}

	name := "default"
	if ps.ServiceAccountName != "" {
		name = ps.ServiceAccountName
	}
	sa, err := w.CoreV1().ServiceAccounts(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	ps.ImagePullSecrets = append(ps.ImagePullSecrets, sa.ImagePullSecrets...)
	return nil
}
