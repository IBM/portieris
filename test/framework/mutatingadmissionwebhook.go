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

package framework

import (
	"context"
	"log"
	"time"

	v1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForMutatingAdmissionWebhook waits until the specified mutating admission webhook is created or the timeout is reached.
func (f *Framework) WaitForMutatingAdmissionWebhook(name string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.KubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(context.TODO(), name, metav1.GetOptions{}); err != nil {
			return false, err
		}
		log.Printf("Found MutatingWebhookConfiguration %q", name)
		return true, nil
	})
}

// ListMutatingAdmissionWebhooks lists the mutating admission webhooks that are associated with the installed Helm release.
func (f *Framework) ListMutatingAdmissionWebhooks() (*v1.MutatingWebhookConfigurationList, error) {
	opts := f.getHelmReleaseSelectorListOptions()
	return f.KubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().List(context.TODO(), opts)
}
