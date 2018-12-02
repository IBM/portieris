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

package notary

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Notary", func() {
	var (
		trust Interface
	)

	BeforeEach(func() {
		trust, _ = NewClient(trustDir, nil)
	})

	Describe("Getting the notary repo", func() {
		It("should return an error", func() {
			_, err := trust.GetNotaryRepo("server", "image", "notaryToken")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("HTTPStore requires an absolute baseURL"))
		})
	})

})
