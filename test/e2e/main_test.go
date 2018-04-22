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

package e2e

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	testFramework "github.com/IBM/portieris/test/framework"
)

var (
	framework *testFramework.Framework
	exitCode  int
	err       error

	noInstall bool

	testTrustImagePolicy, testTrustClusterImagePolicy, testArmada, testVAImagePolicy, testVAClusterImagePolicy, testWildcardImagePolicy, testGeneric bool
)

const (
	ChartName            = "ibmcloud-image-enforcement"
	MutatingWebhookName  = "image-admission-config"
	AdmissionWebhookName = "image-admission-config"
)

func TestMain(m *testing.M) {
	kubeconfig := flag.String("kubeconfig", "", "kube config path")
	helmChart := flag.String("helmChart", "", "helm chart location")
	flag.BoolVar(&noInstall, "no-install", false, "turns off helm chart installation for quicker feedback loops")
	flag.BoolVar(&testTrustImagePolicy, "trust-image-policy", false, "runs trust tests for image policies")
	flag.BoolVar(&testTrustClusterImagePolicy, "trust-cluster-image-policy", false, "runs trust tests for cluster image policies")
	flag.BoolVar(&testArmada, "armada", false, "runs tests for Armada based installation")
	flag.BoolVar(&testVAImagePolicy, "va-image-policy", false, "runs va tests for image policies")
	flag.BoolVar(&testVAClusterImagePolicy, "va-cluster-image-policy", false, "runs va tests for cluster image policies")
	flag.BoolVar(&testWildcardImagePolicy, "wildcards-image-policy", false, "runs tests for wildcards in image policies")
	flag.BoolVar(&testGeneric, "generic", false, "runs generic enforment tests")

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

	framework, err = testFramework.New(*kubeconfig, *helmChart, noInstall)
	if err != nil {
		log.Printf("error during framework initialisation: %v\n", err)
		os.Exit(1)
	}

	if !noInstall {
		// Check for deployment
		if err := framework.WaitForDeployment(fmt.Sprintf("%v-%v", framework.HelmRelease, ChartName), framework.Namespace, time.Minute); err != nil {
			log.Printf("error waiting for deployment to appear: %v\n", err)
			os.Exit(1)
		}

		// Check for CRDs
		if err := framework.WaitForImagePolicyDefinition(time.Minute); err != nil {
			log.Printf("error waiting for ImagePolicyDefinition to appear: %v\n", err)
			os.Exit(1)
		}
		if err := framework.WaitForClusterImagePolicyDefinition(time.Minute); err != nil {
			log.Printf("error waiting for ClusterImagePolicyDefinition to appear: %v\n", err)
			os.Exit(1)
		}

		// Check for mutatingadmissionwebhook
		if err := framework.WaitForMutatingAdmissionWebhook(MutatingWebhookName, time.Minute); err != nil {
			log.Printf("error waiting for MutatingWebhookConfiguration to appear: %v\n", err)
			os.Exit(1)
		}

		// check for validatingadmission webhook
		if err := framework.WaitForValidatingAdmissionWebhook(AdmissionWebhookName, time.Minute); err != nil {
			log.Printf("error waiting for ValidatingWebhookConfiguration to appear: %v\n", err)
			os.Exit(1)
		}
	}

	exitCode = m.Run()
}
