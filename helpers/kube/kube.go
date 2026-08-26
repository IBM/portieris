// Copyright 2018, 2021 Portieris Authors.
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
	"strconv"

	"github.com/IBM/portieris/internal/info"
	portierisclientset "github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/client/clientset/versioned"
	"github.com/IBM/portieris/pkg/policy"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeClientConfig creates a kube client config.
// If KUBE_API_QPS or KUBE_API_BURST are set, they override the client-go defaults (QPS=5, Burst=10).
func GetKubeClientConfig(kubeconfigFileLoc *string, qps float32, burst int) *rest.Config {
	var config *rest.Config
	var err error

	if kubeconfigFileLoc != nil && *kubeconfigFileLoc != "" {
		// If --kubeconfig command-line flag is set, use that kubeconfig file location to create the kube client
		glog.Info(fmt.Sprintf("--kubeconfig command line flag set to %s", *kubeconfigFileLoc))
		// need to confirm that the specified file actually exists before using it
		if _, err = os.Stat(*kubeconfigFileLoc); err == nil {
			glog.Info(fmt.Sprintf("Using %s for kube client config", *kubeconfigFileLoc))
			config, err = clientcmd.BuildConfigFromFlags("", *kubeconfigFileLoc)
		} else {
			glog.Fatal(fmt.Sprintf("%s is not a valid file location", *kubeconfigFileLoc))
		}
	} else if kubeconfig, kubeconfigSet := os.LookupEnv("KUBECONFIG"); kubeconfigSet {
		// If KUBECONFIG ENV var is set, use that kubeconfig file location to create the kube client
		glog.Info(fmt.Sprintf("KUBECONFIG env variable is set to %s", kubeconfig))
		// need to confirm that the specified file actually exists before using it
		if _, err = os.Stat(kubeconfig); err == nil {
			glog.Info(fmt.Sprintf("Using %s for kube client config", kubeconfig))
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		} else {
			glog.Fatal(fmt.Sprintf("%s is not a valid file location", kubeconfig))
		}
	} else {
		// If neither the --kubeconfig flag or the KUBECONFIG env var are set, default to using an in-cluster kube client configuration
		glog.Info("No --kubeconfig flag found and KUBECONFIG env variable is NOT set, defaulting to in-cluster kube client config")
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		glog.Fatal(err)
	}

	config.UserAgent = "portieris/" + info.Version

	if qps == 0 {
		if v, parseErr := strconv.ParseFloat(os.Getenv("KUBE_API_QPS"), 32); parseErr == nil && v > 0 {
			qps = float32(v)
		}
	}
	if burst == 0 {
		if v, atoiErr := strconv.Atoi(os.Getenv("KUBE_API_BURST")); atoiErr == nil && v > 0 {
			burst = v
		}
	}

	if qps > 0 && burst == 0 {
		burst = int(qps * 2)
		glog.Warningf("KUBE_API_BURST not set; derived burst=%d from qps=%.1f (2x ratio)", burst, qps)
	}
	if burst > 0 && qps == 0 {
		qps = float32(burst) / 2
		glog.Warningf("KUBE_API_QPS not set; derived qps=%.1f from burst=%d (2x ratio)", qps, burst)
	}
	config.QPS = qps
	config.Burst = burst

	logQPS, logBurst := config.QPS, config.Burst
	if logQPS == 0 {
		logQPS = 5
	}
	if logBurst == 0 {
		logBurst = 10
	}
	glog.Infof("kube client rate limits: qps=%.1f burst=%d", logQPS, logBurst)

	return config
}

// GetKubeClient creates a kube clientset
func GetKubeClient(config *rest.Config) *kubernetes.Clientset {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatal(err)
	}
	return clientset
}

// GetPolicyClient creates a policy clientset
func GetPolicyClient(config *rest.Config) *policy.Client {
	clientset, err := portierisclientset.NewForConfig(config)
	if err != nil {
		glog.Fatal(err)
	}
	policyClient := policy.NewClient(clientset)
	return policyClient
}
