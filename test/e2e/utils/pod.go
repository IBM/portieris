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

package utils

import (
	"testing"

	"github.com/IBM/portieris/test/framework"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func buildPod(t *testing.T, fw *framework.Framework, manifestLocation, namespace string, expectCreateFail bool) *corev1.Pod {
	manifest, err := fw.LoadPodManifest(manifestLocation)
	if err != nil {
		t.Fatalf("Error loading manifest: %v", err)
	}
	if manifest == nil {
		t.Fatalf("Error loading manifest: manifest is nil")
	}
	if err = fw.CreatePod(namespace, manifest); err != nil {
		if !expectCreateFail {
			t.Fatalf("Failed to create job on success path: %s", err)
		} else {
			return nil
		}
	}
	pod, err := fw.GetPod(manifest.Name, namespace)
	if err != nil {
		t.Fatalf("Error getting %q pod in %v: %v", manifest.Name, namespace, err)
	}
	return pod
}

// TestPodRunnable tests whether a manifest is deployable to the specified namespace
func TestPodRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	pod := buildPod(t, fw, manifestLocation, namespace, false)
	defer fw.DeletePod(pod.Name, pod.Namespace)

	if !assert.Equal(t, corev1.PodRunning, pod.Status.Phase, "Pod failed: Current phase was not running ") {
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
	}
}

// TestPodNotRunnable tests whether a manifest is deployable to the specified namespace
func TestPodNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	pod := buildPod(t, fw, manifestLocation, namespace, true)
	if pod != nil {
		if !assert.Equal(t, corev1.PodFailed, pod.Status.Phase, "Pod Running ") {
			DumpEvents(t, fw, namespace)
			DumpPolicies(t, fw, namespace)

		}
	}
}

// KillPod kills first pod return in podlist in the given namespace
func KillPod(t *testing.T, fw *framework.Framework, namespace string) {
	if err := fw.DeleteRandomPod(namespace); err != nil {
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
		t.Errorf("Failed to delete random pod from namespace")
	}
}
