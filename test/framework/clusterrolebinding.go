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

import v1beta1 "k8s.io/api/rbac/v1beta1"

// ListClusterRoleBindings lists all the cluster role bindings that are associated with the installed Helm release.
func (f *Framework) ListClusterRoleBindings() (*v1beta1.ClusterRoleBindingList, error) {
	opts := f.getHelmReleaseSelectorListOptions()
	return f.KubeClient.RbacV1beta1().ClusterRoleBindings().List(opts)
}
