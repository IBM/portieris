// Copyright 2018 IBM
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

package utils

import (
	"io/ioutil"
	"testing"

	"github.com/IBM/portieris/test/framework"
)

func DumpEvents(t *testing.T, fw *framework.Framework, namespace string) {
	reader := fw.DumpEvents(namespace)
	bytes, _ := ioutil.ReadAll(reader)
	t.Errorf("%s\n", bytes)
}

func DumpPolicies(t *testing.T, fw *framework.Framework, namespace string) {
	reader := fw.DumpPolicies(namespace)
	bytes, _ := ioutil.ReadAll(reader)
	t.Errorf("%s\n", bytes)
}
