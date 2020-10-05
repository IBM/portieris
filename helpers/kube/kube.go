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

package kube

import (
	"fmt"
	"os"

	securityenforcementclientset "github.com/IBM/portieris/pkg/apis/securityenforcement/client/clientset/versioned"
	"github.com/IBM/portieris/pkg/policy"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeClient creates a kube clientset
func GetKubeClient() *kubernetes.Clientset {
	var config *rest.Config
	var err error

	// If KUBECONFIG ENV var is set, use that kubeconfig file location to create the kube client
	kubeconfig, kubeconfigSet := os.LookupEnv("KUBECONFIG")
	if kubeconfigSet {
		glog.Info(fmt.Sprintf("KUBECONFIG env variable is set to %s, using this for kube client config", kubeconfig))
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		glog.Info("KUBECONFIG env variable is NOT set, defaulting to in-cluster kube client config")
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		glog.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatal(err)
	}
	return clientset
}

// GetPolicyClient creates a policy clientset
func GetPolicyClient() (*policy.Client, error) {
	var config *rest.Config
	var err error

	// If KUBECONFIG ENV var is set, use that kubeconfig file location to create the kube client
	kubeconfig, kubeconfigSet := os.LookupEnv("KUBECONFIG")
	if kubeconfigSet {
		glog.Info(fmt.Sprintf("KUBECONFIG env variable is set to %s, using this for kube client config", kubeconfig))
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		glog.Info("KUBECONFIG env variable is NOT set, defaulting to in-cluster kube client config")
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	// Get admission policy clientset
	clientset, err := securityenforcementclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	policyClient := policy.NewClient(clientset)
	return policyClient, nil
}
