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
	"time"

	"github.com/IBM/portieris/test/framework"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
)

func buildDaemonSet(t *testing.T, fw *framework.Framework, manifestLocation, namespace string, expectCreateFail bool) *appsv1.DaemonSet {
	manifest, err := fw.LoadDaemonSetManifest(manifestLocation)
	if err != nil {
		t.Fatalf("Error loading manifest: %v", err)
	}
	if manifest == nil {
		t.Fatalf("Error loading manifest: manifest is nil")
	}
	if err = fw.CreateDaemonSet(namespace, manifest); err != nil {
		if !expectCreateFail {
			t.Fatalf("Error creating %q daemonsets in %v: %v", manifest.Name, namespace, err)
		} else {
			return nil
		}
	}
	fw.WaitForDaemonSetPods(manifest.Name, namespace, time.Minute)
	daemonset, err := fw.GetDaemonSets(manifest.Name, namespace)
	if err != nil {
		t.Fatalf("Error getting %q daemonsets in %v: %v", manifest.Name, namespace, err)
	}
	return daemonset
}

// TestDaemonSetRunnable tests whether a manifest is deployable to the specified namespace
func TestDaemonSetRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	daemonset := buildDaemonSet(t, fw, manifestLocation, namespace, false)
	defer fw.DeleteDaemonSet(daemonset.Name, daemonset.Namespace)
	if !assert.Equal(t, daemonset.Status.DesiredNumberScheduled, daemonset.Status.NumberReady, "DaemonSet failed: available replicas did not match expected replicas") {
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
	}
}

// TestDaemonSetNotRunnable tests whether a manifest is deployable to the specified namespace
func TestDaemonSetNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	daemonset := buildDaemonSet(t, fw, manifestLocation, namespace, true)
	if daemonset != nil {
		defer fw.DeleteDaemonSet(daemonset.Name, daemonset.Namespace)
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
		t.Errorf("Expected daemonset creation to fail")
	}
}
