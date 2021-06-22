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

package framework

import (
	"fmt"

	"time"

	batchv1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// ListCronJobs retrieves all jobs associated with the installed Helm release
func (f *Framework) ListCronJobs() (*batchv1.CronJobList, error) {
	opts := f.getHelmReleaseSelectorListOptions()
	return f.KubeClient.BatchV1beta1().CronJobs(corev1.NamespaceAll).List(opts)
}

// GetCronJob retrieves the specified deployment
func (f *Framework) GetCronJob(name, namespace string) (*batchv1.CronJob, error) {
	return f.KubeClient.BatchV1beta1().CronJobs(namespace).Get(name, metav1.GetOptions{})
}

// LoadCronJobManifest takes a manifest and decodes it into a CronJob object
func (f *Framework) LoadCronJobManifest(pathToManifest string) (*batchv1.CronJob, error) {
	manifest, err := openFile(pathToManifest)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file %q: %v", pathToManifest, err)
	}
	job := batchv1.CronJob{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&job); err != nil {
		return nil, fmt.Errorf("Unable to decode file %q: %v", pathToManifest, err)
	}
	return &job, nil
}

// CreateCronJob creates a CronJob resource and then waits for it to appear
func (f *Framework) CreateCronJob(namespace string, job *batchv1.CronJob) error {
	if _, err := f.KubeClient.BatchV1beta1().CronJobs(namespace).Create(job); err != nil {
		return err
	}
	return f.WaitForCronJob(job.Name, namespace, 2*time.Minute)
}

// DeleteCronJob deletes the specified deployment
func (f *Framework) DeleteCronJob(name, namespace string) error {
	return f.KubeClient.BatchV1beta1().CronJobs(namespace).Delete(name, &metav1.DeleteOptions{})
}

// WaitForCronJob waits until job deployment has completed
func (f *Framework) WaitForCronJob(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		job, err := f.GetCronJob(name, namespace)
		if err != nil {
			return false, err
		}
		if len(job.Status.Active) == 0 {
			return false, nil
		}
		return true, nil
	})
}
