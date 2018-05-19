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
	corev1 "k8s.io/api/core/v1"
)

func buildReplicationController(t *testing.T, fw *framework.Framework, manifestLocation, namespace string, expectCreateFail bool) *corev1.ReplicationController {
	manifest, err := fw.LoadReplicationControllerManifest(manifestLocation)
	if err != nil {
		t.Fatalf("Error loading manifest: %v", err)
	}
	if manifest == nil {
		t.Fatalf("Error loading manifest: manifest is nil")
	}
	if err = fw.CreateReplicationController(namespace, manifest); err != nil {
		if !expectCreateFail {

			t.Fatalf("Error creating %q replicationcontroller in %v: %v", manifest.Name, namespace, err)
		} else {
			return nil
		}
	}
	fw.WaitForReplicationControllerPods(manifest.Name, namespace, time.Minute)
	replicationcontroller, err := fw.GetReplicationController(manifest.Name, namespace)
	if err != nil {
		t.Fatalf("Error getting %q replicationcontroller in %v: %v", manifest.Name, namespace, err)
	}
	return replicationcontroller
}

// TestReplicationControllerRunnable tests whether a manifest is deployable to the specified namespace
func TestReplicationControllerRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	replicationcontroller := buildReplicationController(t, fw, manifestLocation, namespace, false)
	defer fw.DeleteReplicationController(replicationcontroller.Name, replicationcontroller.Namespace)
	if !assert.Equal(t, *replicationcontroller.Spec.Replicas, replicationcontroller.Status.AvailableReplicas, "ReplicationController failed: available replicas did not match expected replicas") {
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
	}
}

// TestReplicationControllerNotRunnable tests whether a manifest is deployable to the specified namespace
func TestReplicationControllerNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) {
	replicationcontroller := buildReplicationController(t, fw, manifestLocation, namespace, true)
	if replicationcontroller != nil {
		defer fw.DeleteReplicationController(replicationcontroller.Name, replicationcontroller.Namespace)
		DumpEvents(t, fw, namespace)
		DumpPolicies(t, fw, namespace)
		t.Errorf("Expected replicacontroller creation to fail")
	}

}
