// Copyright 2018, 2023 Portieris Authors.
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

package e2e

import (
	"testing"

	"github.com/IBM/portieris/test/e2e/utils"
)

// Temporary check until other registries are supported.
func TestNotary_ClusterImagePolicyRepositories_ThirdPartyTrust(t *testing.T) {
	utils.CheckIfTesting(t, testTrustClusterImagePolicy)
	t.Run("Third party trust is rejected", func(t *testing.T) {
		clusterImagePolicy, namespace := utils.CreateClusterImagePolicyAndNamespace(t, framework, "./testdata/clusterimagepolicy/allow-signed-dockerhub.yaml")
		utils.TestDeploymentNotRunnable(t, framework, "./testdata/deployment/dockerhub-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpClusterImagePolicyTest(t, framework, clusterImagePolicy.Name, namespace.Name)
	})
}
