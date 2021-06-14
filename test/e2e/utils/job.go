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

package utils

import (
	"testing"
	"time"

	"github.com/IBM/portieris/test/framework"
	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
)

func buildJob(t *testing.T, fw *framework.Framework, manifestLocation, namespace string, expectCreateFail bool) *batchv1.Job {
	manifest, err := fw.LoadJobManifest(manifestLocation)
	if err != nil {
		t.Fatalf("Error loading manifest: %v", err)
	}
	if manifest == nil {
		t.Fatalf("Error loading manifest: manifest is nil")
	}
	if err = fw.CreateJob(namespace, manifest); err != nil {
		if !expectCreateFail {
			t.Fatalf("Failed to create job on success path: %s", err)
		} else {
			return nil
		}
	}

	fw.WaitForJob(manifest.Name, namespace, time.Second*30)
	job, err := fw.GetJob(manifest.Name, namespace)
	if err != nil {
		t.Fatalf("Error getting %q job in %v: %v", manifest.Name, namespace, err)
	}
	return job
}

// TestJobRunnable tests whether a manifest is deployable to the specified namespace.
func TestJobRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	job := buildJob(t, fw, manifestLocation, namespace, false)
	defer fw.DeleteJob(job.Name, job.Namespace)
	if !assert.Zero(t, job.Status.Failed, "Job failed ") {
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
	}

}

// TestJobNotRunnable tests whether a manifest is deployable to the specified namespace.
func TestJobNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	job := buildJob(t, fw, manifestLocation, namespace, true)
	if job != nil {
		if !assert.NotZero(t, job.Status.Failed, "Job Running ") {
			DumpEvents(t, fw, namespace)
			DumpPolicies(t, fw, namespace)

		}
	}
}
