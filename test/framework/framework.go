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

package framework

import (
	"fmt"
	"net/http"

	"log"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	securityenforcementclientset "github.com/IBM/portieris/pkg/apis/securityenforcement/client/clientset/versioned"
	customResourceDefinitionClientSet "github.com/kubernetes/apiextensions-apiserver/pkg/client/clientset/internalclientset/typed/apiextensions/internalversion"

	// Needed for testing using oidc (Armada)
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

const (
	ns = "ibm-system"

	helmReleaseName = "cise"
	helmChartName   = "ibmcloud-image-enforcement"

	imagePolicyCRDName        = "imagepolicies.securityenforcement.admission.cloud.ibm.com"
	clusterImagePolicyCRDName = "clusterimagepolicies.securityenforcement.admission.cloud.ibm.com"
)

// Framework is an e2e test framework esponsible for installing and deleting of the helm chart
// It also providers helper functions for talking to Kube clusters
type Framework struct {
	KubeClient                     kubernetes.Interface
	ImagePolicyClient              securityenforcementclientset.Interface
	ClusterImagePolicyClient       securityenforcementclientset.Interface
	CustomResourceDefinitionClient customResourceDefinitionClientSet.CustomResourceDefinitionInterface
	HTTPClient                     *http.Client
	Namespace                      string
	HelmRelease                    string
	HelmChart                      string
}

// New installs the specific helm chart into the Kube cluster of the kubeconfig
func New(kubeconfig, helmChart string, noInstall bool) (*Framework, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("unable to build Kube config: %v", err)
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create Kube client: %v", err)
	}
	httpClient := kubeClient.CoreV1().RESTClient().(*rest.RESTClient).Client
	if err != nil {
		return nil, fmt.Errorf("unable to create Kube client: %v", err)
	}
	imagePolicyClient, err := securityenforcementclientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create ImagePolicy client: %v", err)
	}
	clusterImagePolicyClient, err := securityenforcementclientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create ClusterImagePolicy client: %v", err)
	}
	apiExtenstionsClient, err := customResourceDefinitionClientSet.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create CustomResourceDefinition client: %v", err)
	}
	customResourceDefintionClient := apiExtenstionsClient.CustomResourceDefinitions()

	fw := &Framework{
		KubeClient:                     kubeClient,
		HTTPClient:                     httpClient,
		ImagePolicyClient:              imagePolicyClient,
		ClusterImagePolicyClient:       clusterImagePolicyClient,
		CustomResourceDefinitionClient: customResourceDefintionClient,
		Namespace:                      ns,
		HelmRelease:                    helmReleaseName,
		HelmChart:                      helmChart,
	}
	if !noInstall {
		// TODO: Delete all of the relevant resources rather than just check they aren't there
		if !fw.verifyCleanedUp() {
			return nil, fmt.Errorf("FAILED: Detected remains of a previous test run")
		}
		if err := fw.installChart(); err != nil {
			return nil, fmt.Errorf("unable to install helm chart: %v", err)
		}
	}
	return fw, nil
}

// Teardown deletes the chart and then verifies everything has been cleaned up
func (f *Framework) Teardown() bool {
	if err := f.deleteChart(); err != nil {
		log.Printf("error deleting helm chart: %v", err)
	}
	return f.verifyCleanedUp()
}

func (f *Framework) verifyCleanedUp() bool {
	ok := true

	// Verify Webhooks have been cleaned up
	if valWebhooks, err := f.ListValidatingAdmissionWebhooks(); err != nil {
		ok = false
		log.Printf("Error listing ValidatingAdmissionWebhook: %v", err)
	} else if len(valWebhooks.Items) != 0 {
		ok = false
		log.Printf("FAILED: ValidatingAdmissionWebhooks were still present")
		for _, webhook := range valWebhooks.Items {
			fmt.Printf("\t\t\t\t- %v\n", webhook.Name)
		}
	}
	if mutWebhooks, err := f.ListMutatingAdmissionWebhooks(); err != nil {
		ok = false
		log.Printf("Error listing MutatingAdmissionWebhook: %v", err)
	} else if len(mutWebhooks.Items) != 0 {
		ok = false
		log.Printf("FAILED: MutatingAdmissionWebhook were still present")
		for _, webhook := range mutWebhooks.Items {
			fmt.Printf("\t\t\t\t- %v\n", webhook.Name)
		}
	}

	// Verify CRDs have been cleaned up
	if imagePolicyDefinition, err := f.GetImagePolicyDefinition(); err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok {
			if statusErr.Status().Code != http.StatusNotFound {
				ok = false
				log.Printf("Error getting ImagePolicyDefinition: %v", err)
			}
		}
	} else if imagePolicyDefinition != nil {
		ok = false
		log.Printf("FAILED: ImagePolicyDefinition was still present")
		fmt.Printf("\t\t\t\t- %v\n", imagePolicyDefinition.Name)
	}
	if clusterImagePolicyDefinition, err := f.GetClusterImagePolicyDefinition(); err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok {
			if statusErr.Status().Code != http.StatusNotFound {
				ok = false
				log.Printf("Error getting ClusterImagePolicyDefinition: %v", err)
			}
		}
	} else if clusterImagePolicyDefinition != nil {
		ok = false
		log.Printf("FAILED: ClusterImagePolicyDefinition was still present")
		fmt.Printf("\t\t\t\t- %v\n", clusterImagePolicyDefinition.Name)
	}

	// Verify Deployments have been cleaned up
	if deployments, err := f.ListDeployments(); err != nil {
		ok = false
		log.Printf("Error listing Deployments: %v", err)
	} else if len(deployments.Items) != 0 {
		ok = false
		log.Printf("FAILED: Deployments were still present")
		for _, deployment := range deployments.Items {
			fmt.Printf("\t\t\t\t- %v\n", deployment.Name)
		}
	}

	// Verify Services have been cleaned up
	if services, err := f.ListServices(); err != nil {
		ok = false
		log.Printf("Error listing Services: %v", err)
	} else if len(services.Items) != 0 {
		ok = false
		log.Printf("FAILED: Services were still present")
		for _, service := range services.Items {
			fmt.Printf("\t\t\t\t- %v\n", service.Name)
		}
	}

	// Verify Jobs have been cleaned up
	if jobs, err := f.ListJobs(); err != nil {
		ok = false
		log.Printf("Error listing Jobs: %v", err)
	} else if len(jobs.Items) != 0 {
		ok = false
		log.Printf("FAILED: Jobs were still present")
		for _, job := range jobs.Items {
			fmt.Printf("\t\t\t\t- %v\n", job.Name)
		}
	}

	// Verify ConfigMaps have been cleaned up
	if cms, err := f.ListConfigMaps(); err != nil {
		ok = false
		log.Printf("Error listing ConfigMaps: %v", err)
	} else if len(cms.Items) != 0 {
		ok = false
		log.Printf("FAILED: ConfigMaps were still present")
		for _, cm := range cms.Items {
			fmt.Printf("\t\t\t\t- %v\n", cm.Name)
		}
	}

	// Verify ServiceAccounts have been cleaned up
	if sas, err := f.ListServiceAccounts(); err != nil {
		ok = false
		log.Printf("Error listing ServiceAccounts : %v", err)
	} else if len(sas.Items) != 0 {
		ok = false
		log.Printf("FAILED: ServiceAccounts were still present")
		for _, sa := range sas.Items {
			fmt.Printf("\t\t\t\t- %v\n", sa.Name)
		}
	}

	// Verify ClusterRoles have been cleaned up
	if crs, err := f.ListClusterRoles(); err != nil {
		ok = false
		log.Printf("Error listing ClusterRoles: %v", err)
	} else if len(crs.Items) != 0 {
		ok = false
		log.Printf("FAILED: ClusterRoles were still present")
		for _, cr := range crs.Items {
			fmt.Printf("\t\t\t\t- %v\n", cr.Name)
		}
	}

	// Verify ClusterRoleBindings have been cleaned up
	if crbs, err := f.ListClusterRoleBindings(); err != nil {
		ok = false
		log.Printf("Error listing ClusterRoleBindings: %v", err)
	} else if len(crbs.Items) != 0 {
		ok = false
		log.Printf("FAILED: ClusterRoleBindings were still present")
		for _, crb := range crbs.Items {
			fmt.Printf("\t\t\t\t- %v\n", crb.Name)
		}
	}

	return ok
}
