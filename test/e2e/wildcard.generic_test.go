// Copyright 2018 Portieris Authors.
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

func TestWildcards_ImagePolicyRepositories_Wildcards(t *testing.T) {
	utils.CheckIfTesting(t, testWildcardImagePolicy)
	if defaultClusterPolicy := utils.DeleteThenReturnClusterImagePolicy(t, framework, "default"); defaultClusterPolicy != nil {
		defer framework.CreateClusterImagePolicy(defaultClusterPolicy)
	}

	t.Run("Correctly apply policy when single trailing wildcard is used", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-unsigned-trailing-wildcard.yaml")
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Correctly apply policy when single embedded wildcard is used", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-unsigned-embedded-wildcard.yaml")
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})
	t.Run("Correctly apply policy when both embedded and trailing wildcard are used", func(t *testing.T) {
		t.Parallel()
		namespace := utils.CreateImagePolicyInstalledNamespace(t, framework, "./testdata/imagepolicy/allow-unsigned-embedded-trailing-wildcard.yaml")
		utils.TestDeploymentRunnable(t, framework, "./testdata/deployment/global-nginx-unsigned.yaml", namespace.Name)
		utils.CleanUpImagePolicyTest(t, framework, namespace.Name)
	})

}
