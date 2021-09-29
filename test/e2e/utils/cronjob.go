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

func buildCronJob(t *testing.T, fw *framework.Framework, manifestLocation, namespace string, expectCreateFail bool) *batchv1.CronJob {
	manifest, err := fw.LoadCronJobManifest(manifestLocation)
	if err != nil {
		t.Fatalf("Error loading manifest: %v", err)
	}
	if manifest == nil {
		t.Fatalf("Error loading manifest: manifest is nil")
	}
	if err = fw.CreateCronJob(namespace, manifest); err != nil {
		if !expectCreateFail {
			t.Fatalf("Failed to create cronjob on success path: %s", err)
		} else {
			return nil
		}
	}

	fw.WaitForCronJob(manifest.Name, namespace, time.Second*61)
	cronjob, err := fw.GetCronJob(manifest.Name, namespace)
	if err != nil {
		t.Fatalf("Error getting %q cronjob in %v: %v", manifest.Name, namespace, err)
	}
	return cronjob
}

// TestCronJobRunnable tests whether a manifest is deployable to the specified namespace.
func TestCronJobRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	cronjob := buildCronJob(t, fw, manifestLocation, namespace, false)
	defer fw.DeleteCronJob(cronjob.Name, cronjob.Namespace)
	if !assert.NotZero(t, len(cronjob.Status.Active), "CronJob failed ") {
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
	}

}

// TestCronJobNotRunnable tests whether a manifest is deployable to the specified namespace.
func TestCronJobNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	cronjob := buildCronJob(t, fw, manifestLocation, namespace, true)
	if cronjob != nil {
		if !assert.Zero(t, len(cronjob.Status.Active), "CronJob Running ") {
			DumpEvents(t, fw, namespace)
			DumpPolicies(t, fw, namespace)

		}
	}
}
