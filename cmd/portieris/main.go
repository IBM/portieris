// Copyright 2018,2021 Portieris Authors.
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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	kube "github.com/IBM/portieris/helpers/kube"
	"github.com/IBM/portieris/internal/info"
	"github.com/IBM/portieris/pkg/controller/multi"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/IBM/portieris/pkg/metrics"
	notaryclient "github.com/IBM/portieris/pkg/notary"
	registryclient "github.com/IBM/portieris/pkg/registry"
	notaryverifier "github.com/IBM/portieris/pkg/verifier/trust"
	"github.com/IBM/portieris/pkg/webhook"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	mkdir := flag.String("mkdir", "", "create directories needed for Portieris to run")
	kubeconfig := flag.String("kubeconfig", "", "location of kubeconfig file to use for an out-of-cluster kube client configuration")
	imageSigningPublicKeySecretNamespace := flag.String("pubkey-secret-namespace-override", "", "specify a namespace that will override where to look for the image signing public key secret")

	flag.Parse() // glog flags

	if mkdir != nil && *mkdir != "" {
		dirs := strings.Split(*mkdir, ",")
		for i := range dirs {
			dir := dirs[i]
			fmt.Printf("making dir %s\n", dir)
			// make a directory and required parents, permission will be determined by UMASK
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to create: %s, %v\n", dir, err)
				os.Exit(1)
			}
			// explicitly set the required permission
			err = os.Chmod(dir, 0775)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to change mode: %s, %v\n", dir, err)
				os.Exit(1)
			}
		}
		os.Exit(0)
	}

	glog.Info("Starting portieris ", info.Version)

	kubeClientConfig := kube.GetKubeClientConfig(kubeconfig)
	kubeClientset := kube.GetKubeClient(kubeClientConfig)
	kubeWrapper := kubernetes.NewKubeClientsetWrapper(kubeClientset)
	policyClient := kube.GetPolicyClient(kubeClientConfig)

	ca, err := ioutil.ReadFile("/etc/certs/ca.pem")
	if err != nil {
		if os.IsNotExist(err) {
			glog.Info("CA not provided at /etc/certs/ca.pem, will use default system pool")
		} else {
			glog.Fatal("Could not read /etc/certs/ca.pem", err)
		}
	}
	trust, err := notaryclient.NewClient(".trust", ca)
	if err != nil {
		glog.Fatal("Could not get trust client", err)
	}

	serverCert, err := ioutil.ReadFile("/etc/certs/tls.crt")
	if err != nil {
		glog.Fatal("Could not read /etc/certs/tls.crt", err)
	}
	serverKey, err := ioutil.ReadFile("/etc/certs/tls.key")
	if err != nil {
		glog.Fatal("Could not read /etc/certs/tls.key", err)
	}

	cr := registryclient.NewClient()
	nv := notaryverifier.NewVerifier(kubeWrapper, trust, cr)
	pmetrics := metrics.NewMetrics()
	controller := multi.NewController(kubeWrapper, policyClient, nv, pmetrics, *imageSigningPublicKeySecretNamespace)

	// Setup http handler for metrics
	go func() {
		r := mux.NewRouter()
		r.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8080", r)
	}()

	webhook := webhook.NewServer("policy", controller, serverCert, serverKey)
	webhook.Run()
}
