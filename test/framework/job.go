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

package framework

import (
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// ListJobs retrieves all jobs associated with the installed Helm release
func (f *Framework) ListJobs() (*batchv1.JobList, error) {
	opts := f.getHelmReleaseSelectorListOptions()
	return f.KubeClient.BatchV1().Jobs(corev1.NamespaceAll).List(opts)
}

// GetJob retrieves the specified deployment
func (f *Framework) GetJob(name, namespace string) (*batchv1.Job, error) {
	return f.KubeClient.BatchV1().Jobs(namespace).Get(name, metav1.GetOptions{})
}

// LoadJobManifest takes a manifest and decodes it into a Job object
func (f *Framework) LoadJobManifest(pathToManifest string) (*batchv1.Job, error) {
	manifest, err := openFile(pathToManifest)
	if err != nil {
		return nil, fmt.Errorf("Unable to open file %q: %v", pathToManifest, err)
	}
	job := batchv1.Job{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&job); err != nil {
		return nil, fmt.Errorf("Unable to decode file %q: %v", pathToManifest, err)
	}
	return &job, nil
}

// CreateJob creates a Job resource and then waits for it to appear
func (f *Framework) CreateJob(namespace string, job *batchv1.Job) error {
	if _, err := f.KubeClient.BatchV1().Jobs(namespace).Create(job); err != nil {
		return err
	}
	return f.WaitForJob(job.Name, namespace, time.Minute)
}

// DeleteJob deletes the specified deployment
func (f *Framework) DeleteJob(name, namespace string) error {
	return f.KubeClient.BatchV1().Jobs(namespace).Delete(name, &metav1.DeleteOptions{})
}

// WaitForJob waits until job deployment has completed
func (f *Framework) WaitForJob(name, namespace string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		job, err := f.GetJob(name, namespace)
		if err != nil {
			return false, err
		}
		if job.Status.Active == 0 && job.Status.Succeeded == 0 {
			return false, nil
		}

		return true, nil
	})
}
