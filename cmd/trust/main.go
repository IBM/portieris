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
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	notaryclient "github.com/IBM/portieris/pkg/notary"

	kube "github.com/IBM/portieris/helpers/kube"
	"github.com/IBM/portieris/pkg/controller/multi"
	"github.com/IBM/portieris/pkg/kubernetes"
	registryclient "github.com/IBM/portieris/pkg/registry"
	"github.com/IBM/portieris/pkg/webhook"
	"github.com/golang/glog"
)

func main() {
	mkdir := flag.String("mkdir", "", "create directories needed for Portieris to run")

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
