// Copyright 2018, 2026 Portieris Authors.
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

	Describe("CompareImageRef - host/path boundary enforcement (unanchored wildcard bypass)", func() {
		Context("legitimate images that must still be admitted", func() {
			It("should match a subdomain host against *.registry.example.com/myorg/*", func() {
				Expect(CompareImageRef("*.registry.example.com/myorg/*", "us.registry.example.com/myorg/foo:latest")).To(BeTrue())
			})
			It("should match a subdomain host against *.registry.example.com/myrepo/*", func() {
				Expect(CompareImageRef("*.registry.example.com/myrepo/*", "eu.registry.example.com/myrepo/node:v1")).To(BeTrue())
			})
			It("should match registry*.example.com/myorg/* against a legitimate host", func() {
				Expect(CompareImageRef("registry*.example.com/myorg/*", "registry01.example.com/myorg/worker:v1")).To(BeTrue())
			})
			It("should match an exact host pattern registry.example.com/myorg/*", func() {
				Expect(CompareImageRef("registry.example.com/myorg/*", "registry.example.com/myorg/foo:latest")).To(BeTrue())
			})
			It("should match a bare * pattern against any image", func() {
				Expect(CompareImageRef("*", "attacker.example.com/anything:latest")).To(BeTrue())
			})
			It("should match trusted.example.com/* against any image on that host", func() {
				Expect(CompareImageRef("trusted.example.com/*", "trusted.example.com/myorg/myimage:demo")).To(BeTrue())
			})
		})

		Context("attacker images that must be denied", func() {
			It("should not match *.registry.example.com/myorg/* when host is attacker-controlled (tagged)", func() {
				Expect(CompareImageRef("*.registry.example.com/myorg/*", "attacker.com/x.registry.example.com/myorg/malware:latest")).To(BeFalse())
			})
			It("should not match *.registry.example.com/myorg/* when host is attacker-controlled (untagged)", func() {
				Expect(CompareImageRef("*.registry.example.com/myorg/*", "attacker.com/x.registry.example.com/myorg/malware")).To(BeFalse())
			})
			It("should not match registry*.example.com/myorg/* when host is attacker-controlled", func() {
				Expect(CompareImageRef("registry*.example.com/myorg/*", "registry.attacker.com/x.example.com/myorg/evil")).To(BeFalse())
			})
			It("should not match *.registry.example.com/myrepo/* when trusted literal is in path only", func() {
				Expect(CompareImageRef("*.registry.example.com/myrepo/*", "evil.io/fake.registry.example.com/myrepo/pwn:latest")).To(BeFalse())
			})
			It("should not match registry.example.com/myorg/* when host is different", func() {
				Expect(CompareImageRef("registry.example.com/myorg/*", "evilreg.io/myorg/foo:latest")).To(BeFalse())
			})
			It("should not match trusted.example.com/* when a different host embeds trusted.example.com in its path", func() {
				Expect(CompareImageRef("trusted.example.com/*", "evil.io/trusted.example.com/myorg/myimage:demo")).To(BeFalse())
			})
		})
	})
})
