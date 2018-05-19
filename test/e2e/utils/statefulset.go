// Copyright 2018 Portieris Authors.
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

func buildStatefulSet(t *testing.T, fw *framework.Framework, manifestLocation, namespace string, expectCreateFail bool) *appsv1.StatefulSet {
	manifest, err := fw.LoadStatefulSetManifest(manifestLocation)
	if err != nil {
		t.Fatalf("Error loading manifest: %v", err)
	}
	if manifest == nil {
		t.Fatalf("Error loading manifest: manifest is nil")
	}
	if err = fw.CreateStatefulSet(namespace, manifest); err != nil {
		if !expectCreateFail {

			t.Fatalf("Error creating %q statefulset in %v: %v", manifest.Name, namespace, err)
		} else {
			return nil
		}
	}
	fw.WaitForStatefulSetPods(manifest.Name, namespace, time.Minute)
	statefulset, err := fw.GetStatefulSet(manifest.Name, namespace)
	if err != nil {
		t.Fatalf("Error getting %q statefulset in %v: %v", manifest.Name, namespace, err)
	}
	return statefulset
}

// TestStatefulSetRunnable tests whether a manifest is deployable to the specified namespace
func TestStatefulSetRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	statefulset := buildStatefulSet(t, fw, manifestLocation, namespace, false)
	defer fw.DeleteStatefulSet(statefulset.Name, statefulset.Namespace)
	if !assert.Equal(t, *statefulset.Spec.Replicas, statefulset.Status.ReadyReplicas, "StatefulSet failed: available replicas did not match expected replicas") {
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
	}
}

// TestStatefulSetNotRunnable tests whether a manifest is deployable to the specified namespace
func TestStatefulSetNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	statefulset := buildStatefulSet(t, fw, manifestLocation, namespace, true)
	if statefulset != nil {
		defer fw.DeleteStatefulSet(statefulset.Name, statefulset.Namespace)

		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
		t.Errorf("Expected statefulset creation to fail")
	}

}
