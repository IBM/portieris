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

package kube

import (
	securityenforcementclientset "github.com/IBM/portieris/pkg/apis/securityenforcement/client/clientset/versioned"
	"github.com/IBM/portieris/pkg/policy"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetKubeClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		glog.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatal(err)
	}
	return clientset
}

func GetPolicyClient() (*policy.Client, error) {
	// Get configuration
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// Get admission policy clientset
	clientset, err := securityenforcementclientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	policyClient := policy.NewClient(clientset)
	return policyClient, nil
}
