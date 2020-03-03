// Copyright 2018,2020 Portieris Authors.
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
	"io/ioutil"
	"os"

	notaryclient "github.com/IBM/portieris/pkg/notary"

	kube "github.com/IBM/portieris/helpers/kube"
	"github.com/IBM/portieris/pkg/controller/multi"
	"github.com/IBM/portieris/pkg/kubernetes"
	registryclient "github.com/IBM/portieris/pkg/registry"
	"github.com/IBM/portieris/pkg/webhook"
	"github.com/golang/glog"
)

func main() {
	flag.Parse() // glog flags
	kubeClientset := kube.GetKubeClient()
	kubeWrapper := kubernetes.NewKubeClientsetWrapper(kubeClientset)
	policyClient, err := kube.GetPolicyClient()
	if err != nil {
		glog.Fatal("Could not get policy client", err)
	}

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

	serverCert, err := ioutil.ReadFile("/etc/certs/serverCert.pem")
	if err != nil {
		glog.Fatal("Could not read /etc/certs/serverCert.pem", err)
	}
	serverKey, err := ioutil.ReadFile("/etc/certs/serverKey.pem")
	if err != nil {
		glog.Fatal("Could not read /etc/certs/serverKey.pem", err)
	}

	cr := registryclient.NewClient()
	controller := multi.NewController(kubeWrapper, policyClient, trust, cr)
	webhook := webhook.NewServer("policy", controller, serverCert, serverKey)
	webhook.Run()
}
