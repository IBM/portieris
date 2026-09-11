// Copyright 2018, 2026 Portieris Authors.
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

package kube

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/IBM/portieris/internal/info"
)

// minimalKubeconfig writes a valid kubeconfig to a temp file and returns its path.
func minimalKubeconfig(t *testing.T) string {
	t.Helper()
	kubeconfigPath := filepath.Join(t.TempDir(), "kubeconfig")
	kubeconfig := []byte(`
apiVersion: v1
kind: Config
clusters:
- name: test-cluster
  cluster:
    server: https://127.0.0.1
contexts:
- name: test-context
  context:
    cluster: test-cluster
    user: test-user
current-context: test-context
users:
- name: test-user
  user:
    token: test-token
`)
	if err := os.WriteFile(kubeconfigPath, kubeconfig, 0600); err != nil {
		t.Fatalf("failed to write temporary kubeconfig: %v", err)
	}
	return kubeconfigPath
}

// TestGetKubeClientConfigDefaults checks that zero values leave the config
// fields unset, letting client-go apply its own defaults.
func TestGetKubeClientConfigDefaults(t *testing.T) {
	path := minimalKubeconfig(t)
	config := GetKubeClientConfig(&path, 0, 0)
	if config == nil {
		t.Fatal("GetKubeClientConfig returned nil")
	}
	if got, want := config.UserAgent, "portieris/"+info.Version; got != want {
		t.Errorf("config.UserAgent = %q, want %q", got, want)
	}
	if got := config.QPS; got != 0 {
		t.Errorf("config.QPS = %v, want 0 (zero signals client-go to use its own default of 5)", got)
	}
	if got := config.Burst; got != 0 {
		t.Errorf("config.Burst = %d, want 0 (zero signals client-go to use its own default of 10)", got)
	}
}

// TestGetKubeClientConfigExplicitRateLimits checks that both values are applied
// when provided explicitly.
func TestGetKubeClientConfigExplicitRateLimits(t *testing.T) {
	path := minimalKubeconfig(t)
	config := GetKubeClientConfig(&path, 50, 100)
	if config == nil {
		t.Fatal("GetKubeClientConfig returned nil")
	}
	if got, want := config.QPS, float32(50); got != want {
		t.Errorf("config.QPS = %v, want %v", got, want)
	}
	if got, want := config.Burst, 100; got != want {
		t.Errorf("config.Burst = %d, want %d", got, want)
	}
}

// TestGetKubeClientConfigQPSOnlyDerivesBurst checks that burst is auto-derived
// as 2×qps when only qps is set.
func TestGetKubeClientConfigQPSOnlyDerivesBurst(t *testing.T) {
	path := minimalKubeconfig(t)
	config := GetKubeClientConfig(&path, 50, 0)
	if config == nil {
		t.Fatal("GetKubeClientConfig returned nil")
	}
	if got, want := config.QPS, float32(50); got != want {
		t.Errorf("config.QPS = %v, want %v", got, want)
	}
	if got, want := config.Burst, 100; got != want {
		t.Errorf("config.Burst = %d, want %d (2x qps)", got, want)
	}
}

// TestGetKubeClientConfigBurstOnlyDerivesQPS checks that qps is auto-derived
// as burst/2 when only burst is set.
func TestGetKubeClientConfigBurstOnlyDerivesQPS(t *testing.T) {
	path := minimalKubeconfig(t)
	config := GetKubeClientConfig(&path, 0, 100)
	if config == nil {
		t.Fatal("GetKubeClientConfig returned nil")
	}
	if got, want := config.QPS, float32(50); got != want {
		t.Errorf("config.QPS = %v, want %v (burst/2)", got, want)
	}
	if got, want := config.Burst, 100; got != want {
		t.Errorf("config.Burst = %d, want %d", got, want)
	}
}

// TestGetKubeClientConfigEnvVarsBoth checks that KUBE_API_QPS and KUBE_API_BURST
// env vars are read when arguments are left at zero (the production path).
func TestGetKubeClientConfigEnvVarsBoth(t *testing.T) {
	t.Setenv("KUBE_API_QPS", "50")
	t.Setenv("KUBE_API_BURST", "100")
	path := minimalKubeconfig(t)
	config := GetKubeClientConfig(&path, 0, 0)
	if config == nil {
		t.Fatal("GetKubeClientConfig returned nil")
	}
	if got, want := config.QPS, float32(50); got != want {
		t.Errorf("config.QPS = %v, want %v", got, want)
	}
	if got, want := config.Burst, 100; got != want {
		t.Errorf("config.Burst = %d, want %d", got, want)
	}
}

// TestGetKubeClientConfigEnvVarQPSOnly checks that burst is derived from KUBE_API_QPS
// when KUBE_API_BURST is absent.
func TestGetKubeClientConfigEnvVarQPSOnly(t *testing.T) {
	t.Setenv("KUBE_API_QPS", "50")
	t.Setenv("KUBE_API_BURST", "")
	path := minimalKubeconfig(t)
	config := GetKubeClientConfig(&path, 0, 0)
	if config == nil {
		t.Fatal("GetKubeClientConfig returned nil")
	}
	if got, want := config.QPS, float32(50); got != want {
		t.Errorf("config.QPS = %v, want %v", got, want)
	}
	if got, want := config.Burst, 100; got != want {
		t.Errorf("config.Burst = %d, want %d (2x qps)", got, want)
	}
}

// TestGetKubeClientConfigEnvVarBurstOnly checks that qps is derived from KUBE_API_BURST
// when KUBE_API_QPS is absent.
func TestGetKubeClientConfigEnvVarBurstOnly(t *testing.T) {
	t.Setenv("KUBE_API_QPS", "")
	t.Setenv("KUBE_API_BURST", "100")
	path := minimalKubeconfig(t)
	config := GetKubeClientConfig(&path, 0, 0)
	if config == nil {
		t.Fatal("GetKubeClientConfig returned nil")
	}
	if got, want := config.QPS, float32(50); got != want {
		t.Errorf("config.QPS = %v, want %v (burst/2)", got, want)
	}
	if got, want := config.Burst, 100; got != want {
		t.Errorf("config.Burst = %d, want %d", got, want)
	}
}
