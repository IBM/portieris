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
	"log"
	"time"

	v1beta1 "k8s.io/api/admissionregistration/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForValidatingAdmissionWebhook waits until the specified ValidationAdmissionWebhook is created or the timeout is reached
func (f *Framework) WaitForValidatingAdmissionWebhook(name string, timeout time.Duration) error {
	return wait.Poll(time.Second*5, timeout, func() (bool, error) {
		if _, err := f.KubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Get(name, metav1.GetOptions{}); err != nil {
			return false, err
		}
		log.Printf("Found ValidatingWebhookConfiguration %q", name)
		return true, nil
	})
}

// ListValidatingAdmissionWebhooks retrieves all ValidatingAdmissionWebhooks associated with the installed Helm release
func (f *Framework) ListValidatingAdmissionWebhooks() (*v1beta1.ValidatingWebhookConfigurationList, error) {
	opts := f.getHelmReleaseSelectorListOptions()
	return f.KubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().List(opts)
}
