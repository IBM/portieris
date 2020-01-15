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

package wildcard

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	Describe("Wildcard test", func() {
		Context("when the pattern is empty", func() {
			It("match should be true if str is empty", func() {
				match := Compare("", "")
				Expect(match).To(Equal(true))
			})
			It("match should be false if str is not empty", func() {
				match := Compare("", "abcdefgh")
				Expect(match).To(Equal(false))
			})
		})

		Context("when the pattern has no wildcards", func() {
			It("should be true when strings match", func() {
				match := Compare("abcdefgh", "abcdefgh")
				Expect(match).To(Equal(true))
			})
			It("should be false when strings do not", func() {
				match := Compare("abcdefgh", "stuvwxyz")
				Expect(match).To(Equal(false))
			})
		})

		Context("when the pattern is a single wildcard)", func() {
			It("should match anything", func() {
				match := Compare("*", "stuvwxyz")
				Expect(match).To(Equal(true))
			})
			It("should match anything", func() {
				match := Compare("*", "abcd*fgh")
				Expect(match).To(Equal(true))
			})
		})

		Context("when the pattern has a leading wildcard)", func() {
			It("should match if end of pattern and str match", func() {
				match := Compare("*wxyz", "stuvwxyz")
				Expect(match).To(Equal(true))
			})
			It("should not match if end of pattern and str do not match", func() {
				match := Compare("*efgh", "stuvwxyz")
				Expect(match).To(Equal(false))
			})
		})

		Context("when the pattern has multiple wildcards)", func() {
			It("multiple wildcards to match", func() {
				match := Compare("s*wx*z", "stuvwxyz")
				Expect(match).To(Equal(true))
			})
			It("multiple wildcards to not match", func() {
				match := Compare("a*de*h", "stuvwxyz")
				Expect(match).To(Equal(false))
			})
			It("multiple wildcards to not match", func() {
				match := Compare("a*de*h", "abcwxyz")
				Expect(match).To(Equal(false))
			})
		})

		Context("when the pattern has a traiing wildcard)", func() {
			It("should match if beginning of pattern and str match", func() {
				match := Compare("stuv*", "stuvwxyz")
				Expect(match).To(Equal(true))
			})
			It("should not match if beginning of pattern and str do not match", func() {
				match := Compare("efgh*", "stuvwxyz")
				Expect(match).To(Equal(false))
			})
		})
	})

	Describe("CompareAnyTag - Compare with no tags supplied", func() {
		It("should match if str has a tag specified", func() {
			match := CompareAnyTag("stuvwxyz", "stuvwxyz:latest")
			Expect(match).To(Equal(true))
		})
		It("should match if pattern has a wildcard tag", func() {
			match := CompareAnyTag("stuvwxyz:*", "stuvwxyz:latest")
			Expect(match).To(Equal(true))
		})
		It("should match if pattern has a wildcard at namespace level", func() {
			match := CompareAnyTag("uk.icr.io/*", "uk.icr.io/imagesec-demo/alpinegood:demo")
			Expect(match).To(Equal(true))
		})
		It("should match if pattern has a wildcard at image level", func() {
			match := CompareAnyTag("uk.icr.io/imagesec-demo/*", "uk.icr.io/imagesec-demo/alpinegood:demo")
			Expect(match).To(Equal(true))
		})
		It("should match if pattern doesn't have the tag", func() {
			match := CompareAnyTag("uk.icr.io/imagesec-demo/alpinegood", "uk.icr.io/imagesec-demo/alpinegood:demo")
			Expect(match).To(Equal(true))
		})
		It("should not match if base in incorrect", func() {
			match := CompareAnyTag("abcdefgh", "stuvwxyz:latest")
			Expect(match).To(Equal(false))
		})
	})
})
