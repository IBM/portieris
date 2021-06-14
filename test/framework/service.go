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

import corev1 "k8s.io/api/core/v1"

// ListServices lists all services that are associated with the installed Helm release.
func (f *Framework) ListServices() (*corev1.ServiceList, error) {
	opts := f.getHelmReleaseSelectorListOptions()
	return f.KubeClient.CoreV1().Services(corev1.NamespaceAll).List(opts)
}
