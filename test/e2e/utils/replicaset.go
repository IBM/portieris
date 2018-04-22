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

func buildReplicaSet(t *testing.T, fw *framework.Framework, manifestLocation, namespace string, expectCreateFail bool) *appsv1.ReplicaSet {
	manifest, err := fw.LoadReplicaSetManifest(manifestLocation)
	if err != nil {
		t.Fatalf("Error loading manifest: %v", err)
	}
	if manifest == nil {
		t.Fatalf("Error loading manifest: manifest is nil")
	}
	if err = fw.CreateReplicaSet(namespace, manifest); err != nil {
		if !expectCreateFail {
			t.Fatalf("Error creating %q replicaset in %v: %v", manifest.Name, namespace, err)
		} else {
			return nil
		}
	}
	fw.WaitForReplicaSetPods(manifest.Name, namespace, time.Minute)
	replicaset, err := fw.GetReplicaSet(manifest.Name, namespace)
	if err != nil {
		t.Fatalf("Error getting %q replicaset in %v: %v", manifest.Name, namespace, err)
	}
	return replicaset
}

// TestReplicaSetRunnable tests whether a manifest is deployable to the specified namespace
func TestReplicaSetRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	replicaset := buildReplicaSet(t, fw, manifestLocation, namespace, false)
	defer fw.DeleteReplicaSet(replicaset.Name, replicaset.Namespace)
	if !assert.Equal(t, *replicaset.Spec.Replicas, replicaset.Status.AvailableReplicas, "ReplicaSet failed: available replicas did not match expected replicas") {
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
	}
}

// TestReplicaSetNotRunnable tests whether a manifest is deployable to the specified namespace
func TestReplicaSetNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	replicaset := buildReplicaSet(t, fw, manifestLocation, namespace, true)
	if replicaset != nil {
		defer fw.DeleteReplicaSet(replicaset.Name, replicaset.Namespace)
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
		t.Errorf("Expected replicaset creation to fail")

	}
}
