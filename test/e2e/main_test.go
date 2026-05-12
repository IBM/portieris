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

package e2e

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	testFramework "github.com/IBM/portieris/test/framework"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	framework *testFramework.Framework
	exitCode  int
	err       error

	noInstall bool

	testTrustImagePolicy, testTrustClusterImagePolicy, testArmada, testVAImagePolicy, testWildcardImagePolicy, testGeneric, testSimpleImagePolicy, testSimpleClusterImagePolicy, testZeroReplica bool
)

const (
	ChartName            = "portieris"
	MutatingWebhookName  = "image-admission-config"
	AdmissionWebhookName = "image-admission-config"
)

// validateRegistryCredentials checks if the IBM Cloud Container Registry credentials are valid
func validateRegistryCredentials(framework *testFramework.Framework) error {
	log.Println("Validating IBM Cloud Container Registry credentials...")

	// Get the secret from default namespace
	secretNames := []string{"all-icr-io", "default-icr-io"}
	var secretData map[string][]byte

	for _, secretName := range secretNames {
		s, err := framework.KubeClient.CoreV1().Secrets("default").Get(context.TODO(), secretName, metav1.GetOptions{})
		if err == nil {
			secretData = s.Data
			log.Printf("Found secret: %s\n", secretName)
			break
		}
	}

	if secretData == nil {
		return fmt.Errorf("no IBM Cloud registry secrets found in default namespace (looked for: %v)", secretNames)
	}

	// Parse the dockerconfigjson
	dockerConfigJSON, ok := secretData[".dockerconfigjson"]
	if !ok {
		return fmt.Errorf("secret does not contain .dockerconfigjson")
	}

	var dockerConfig struct {
		Auths map[string]struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Auth     string `json:"auth"`
		} `json:"auths"`
	}

	if err := json.Unmarshal(dockerConfigJSON, &dockerConfig); err != nil {
		return fmt.Errorf("failed to parse dockerconfigjson: %v", err)
	}

	// Test credentials for each registry
	validRegistries := 0
	for registry, creds := range dockerConfig.Auths {
		// Skip non-ICR registries
		if !strings.Contains(registry, "icr.io") {
			continue
		}

		username := creds.Username
		password := creds.Password

		// If auth field is present, decode it
		if creds.Auth != "" {
			decoded, err := base64.StdEncoding.DecodeString(creds.Auth)
			if err == nil {
				parts := strings.SplitN(string(decoded), ":", 2)
				if len(parts) == 2 {
					username = parts[0]
					password = parts[1]
				}
			}
		}

		// Test the credentials by attempting to get a token
		tokenURL := fmt.Sprintf("https://%s/oauth/token", registry)
		req, err := http.NewRequest("GET", tokenURL, nil)
		if err != nil {
			log.Printf("Warning: failed to create request for %s: %v\n", registry, err)
			continue
		}

		req.SetBasicAuth(username, password)
		req.Header.Set("Service", "registry")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Warning: failed to validate credentials for %s: %v\n", registry, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized {
			// 200 = valid credentials, 401 = endpoint exists but creds invalid
			if resp.StatusCode == http.StatusUnauthorized {
				return fmt.Errorf("invalid credentials for registry %s (HTTP 401). Please update the secret in the default namespace with valid IBM Cloud API key", registry)
			}
			log.Printf("✓ Credentials validated for %s\n", registry)
			validRegistries++
		} else {
			log.Printf("Warning: unexpected response %d from %s\n", resp.StatusCode, registry)
		}
	}

	if validRegistries == 0 {
		return fmt.Errorf("no valid IBM Cloud Container Registry credentials found. Please ensure the secret in default namespace contains valid credentials")
	}

	log.Printf("Successfully validated credentials for %d registries\n", validRegistries)
	return nil
}

func TestMain(m *testing.M) {
	helmChart := flag.String("helmChart", "", "helm chart location")
	flag.BoolVar(&noInstall, "no-install", false, "turns off helm chart installation for quicker feedback loops")
	flag.BoolVar(&testTrustImagePolicy, "trust-image-policy", false, "runs trust tests for image policies")
	flag.BoolVar(&testTrustClusterImagePolicy, "trust-cluster-image-policy", false, "runs trust tests for cluster image policies")
	flag.BoolVar(&testArmada, "armada", false, "runs tests for Armada based installation")
	flag.BoolVar(&testWildcardImagePolicy, "wildcards-image-policy", false, "runs tests for wildcards in image policies")
	flag.BoolVar(&testGeneric, "generic", false, "runs generic enforment tests")
	flag.BoolVar(&testSimpleImagePolicy, "simple-image-policy", false, "runs tests for simple signing policies")
	flag.BoolVar(&testSimpleClusterImagePolicy, "simple-cluster-image-policy", false, "runs tests for simple signing policies")
	flag.BoolVar(&testZeroReplica, "zero-replica", false, "runs tests for zero replica enforcement")

	flag.Parse()

	defer func() {
		if !noInstall {
			if ok := framework.Teardown(); !ok {
				log.Print("framework teardown had some errors\n")
				os.Exit(1)
			}
			os.Exit(exitCode)
		}
	}()

	framework, err = testFramework.New(os.Getenv("KUBECONFIG"), *helmChart, noInstall)
	if err != nil {
		log.Printf("error during framework initialisation: %v\n", err)
		os.Exit(1)
	}

	// Validate registry credentials before running tests
	if err := validateRegistryCredentials(framework); err != nil {
		log.Printf("ERROR: Registry credential validation failed: %v\n", err)
		log.Println("Please update the IBM Cloud Container Registry credentials in the 'default' namespace before running tests.")
		log.Println("This validation prevents test timeouts caused by invalid API keys.")
		os.Exit(1)
	}

	if !noInstall {
		// Check for deployment.
		deploymentName := fmt.Sprintf("%v-%v", framework.HelmRelease, ChartName)
		if err := framework.WaitForDeployment(deploymentName, framework.Namespace, time.Minute); err != nil {
			log.Printf("error waiting for deployment %s in %s to appear: %v\n", deploymentName, framework.Namespace, err)
			os.Exit(1)
		}

		// Check for CRDs.
		if err := framework.WaitForImagePolicyDefinition(time.Minute); err != nil {
			log.Printf("error waiting for ImagePolicyDefinition to appear: %v\n", err)
			os.Exit(1)
		}
		if err := framework.WaitForClusterImagePolicyDefinition(time.Minute); err != nil {
			log.Printf("error waiting for ClusterImagePolicyDefinition to appear: %v\n", err)
			os.Exit(1)
		}

		// Check for mutatingadmissionwebhook.
		if err := framework.WaitForMutatingAdmissionWebhook(MutatingWebhookName, time.Minute); err != nil {
			log.Printf("error waiting for MutatingWebhookConfiguration to appear: %v\n", err)
			os.Exit(1)
		}

		// Check for validatingadmission webhook.
		if err := framework.WaitForValidatingAdmissionWebhook(AdmissionWebhookName, time.Minute); err != nil {
			log.Printf("error waiting for ValidatingWebhookConfiguration to appear: %v\n", err)
			os.Exit(1)
		}
	}

	exitCode = m.Run()
}
