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

package image

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TODO: Work out why these tests aren't running/rip out ginkgo

func TestImage(t *testing.T) {

	var _ = Describe("Image", func() {

		Describe("When the image name is invalid", func() {
			It("should return an error", func() {
				image, err := NewImageReference("?")
				Expect(err).To(HaveOccurred())
				Expect(image).To(BeNil())
			})
		})

		Describe("When the image is valid without a tag", func() {
			It("should use `latest` as tag", func() {
				image, err := NewImageReference("test.com/namespace/name")
				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
				Expect(image.GetTag()).To(Equal("latest"))
				Expect(image.NameWithTag()).To(Equal("test.com/namespace/name:latest"))
				Expect(image.NameWithoutTag()).To(Equal("test.com/namespace/name"))
				Expect(image.String()).To(Equal("test.com/namespace/name"))
			})
		})

		Describe("When the image is valid without a tag but with digest", func() {
			It("should be OK", func() {
				image, err := NewImageReference("test.com:8080/namespace/name@sha256:1234567890")
				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
				Expect(image.GetTag()).To(Equal("latest"))
				Expect(image.NameWithTag()).To(Equal("test.com:8080/namespace/name:latest"))
				Expect(image.NameWithoutTag()).To(Equal("test.com:8080/namespace/name"))
				Expect(image.String()).To(Equal("test.com:8080/namespace/name@sha256:1234567890"))
				Expect(image.GetDigest()).To(Equal("1234567890"))
			})
		})

		Describe("When the image is valid with a tag", func() {
			It("should not be latest", func() {
				image, err := NewImageReference("test.com/namespace/name:v1")
				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
				Expect(image.GetTag()).To(Equal("v1"))
				Expect(image.NameWithTag()).To(Equal("test.com/namespace/name:v1"))
				Expect(image.NameWithoutTag()).To(Equal("test.com/namespace/name"))
				Expect(image.String()).To(Equal("test.com/namespace/name:v1"))
			})
		})

		Describe("When the image is valid and has a digest", func() {
			It("should be OK", func() {
				image, err := NewImageReference("test.com:8080/namespace/name:v1@sha256:1234567890")
				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
				Expect(image.GetTag()).To(Equal("v1"))
				Expect(image.NameWithTag()).To(Equal("test.com:8080/namespace/name:v1"))
				Expect(image.NameWithoutTag()).To(Equal("test.com:8080/namespace/name"))
				Expect(image.String()).To(Equal("test.com:8080/namespace/name:v1@sha256:1234567890"))
				Expect(image.GetDigest()).To(Equal("1234567890"))

			})
		})

		Describe("When the image is from Docker Hub and has a digest", func() {
			It("should be OK", func() {
				image, err := NewImageReference("namespace/name:v1@sha256:1234567890")
				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
				Expect(image.GetTag()).To(Equal("v1"))
				Expect(image.NameWithTag()).To(Equal("namespace/name:v1"))
				Expect(image.NameWithoutTag()).To(Equal("namespace/name"))
				Expect(image.String()).To(Equal("namespace/name:v1@sha256:1234567890"))
				Expect(image.GetDigest()).To(Equal("1234567890"))
				Expect(image.GetHostname()).To(Equal("docker.io"))
			})
		})

		Describe("When the image is a Docker Hub public image and has a digest", func() {
			It("should be OK", func() {
				image, err := NewImageReference("ubuntu:v1@sha256:1234567890")
				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
				Expect(image.GetTag()).To(Equal("v1"))
				Expect(image.NameWithTag()).To(Equal("ubuntu:v1"))
				Expect(image.NameWithoutTag()).To(Equal("ubuntu"))
				Expect(image.String()).To(Equal("ubuntu:v1@sha256:1234567890"))
				Expect(image.GetDigest()).To(Equal("1234567890"))
				Expect(image.GetHostname()).To(Equal("docker.io"))
			})
		})

		Describe("When the image is valid and it's an IBM repository", func() {
			It("should be OK", func() {
				image, err := NewImageReference("registry.ng.bluemix.net/namespace/name")
				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
				Expect(image.GetHostname()).To(Equal("registry.ng.bluemix.net"))
				Expect(image.GetPort()).To(Equal(""))
				Expect(image.HasIBMRepo()).To(BeTrue())
				Expect(image.GetRegistryURL()).To(Equal("https://registry.ng.bluemix.net"))
				Expect(image.GetContentTrustURL()).To(Equal("https://registry.ng.bluemix.net:4443"))
				vaURL, err := image.GetVAURL()
				Expect(err).ToNot(HaveOccurred())
				Expect(vaURL).To(Equal("https://va.ng.bluemix.net"))
				Expect(image.NameWithTag()).To(Equal("registry.ng.bluemix.net/namespace/name:latest"))
				Expect(image.NameWithoutTag()).To(Equal("registry.ng.bluemix.net/namespace/name"))
				Expect(image.GetTag()).To(Equal("latest"))
			})
		})

		Describe("When the image is valid and it's an IBM repository but the hostname has a port", func() {
			It("should be OK", func() {
				image, err := NewImageReference("registry.ng.bluemix.net:8080/namespace/name:v1")
				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
				Expect(image.GetHostname()).To(Equal("registry.ng.bluemix.net"))
				Expect(image.GetPort()).To(Equal("8080"))
				Expect(image.HasIBMRepo()).To(BeTrue())
				Expect(image.GetRegistryURL()).To(Equal("https://registry.ng.bluemix.net:8080"))
				Expect(image.GetContentTrustURL()).To(Equal("https://registry.ng.bluemix.net:4443"))
				vaURL, err := image.GetVAURL()
				Expect(err).ToNot(HaveOccurred())
				Expect(vaURL).To(Equal("https://va.ng.bluemix.net"))
				Expect(image.GetTag()).To(Equal("v1"))
				Expect(image.NameWithTag()).To(Equal("registry.ng.bluemix.net:8080/namespace/name:v1"))
				Expect(image.NameWithoutTag()).To(Equal("registry.ng.bluemix.net:8080/namespace/name"))
			})
		})

		Describe("When the image is valid but it's not an IBM repository", func() {
			It("should be OK", func() {
				image, err := NewImageReference("test.com/namespace/name")
				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
				Expect(image.GetHostname()).To(Equal("test.com"))
				Expect(image.GetPort()).To(Equal(""))
				Expect(image.HasIBMRepo()).To(BeFalse())
				Expect(image.GetRegistryURL()).To(Equal("https://test.com"))
				Expect(image.GetContentTrustURL()).To(Equal("https://test.com:4443"))
				vaURL, err := image.GetVAURL()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`Deny "test.com/namespace/name", Vulnerability Advisor is not supported for images from this registry: test.com`))
				Expect(vaURL).To(Equal(""))
				Expect(image.GetTag()).To(Equal("latest"))
				Expect(image.NameWithTag()).To(Equal("test.com/namespace/name:latest"))
				Expect(image.NameWithoutTag()).To(Equal("test.com/namespace/name"))
			})
		})

		Describe("When the image is valid but it's not an IBM repository and the hostname has a port", func() {
			It("should be OK", func() {
				image, err := NewImageReference("test.com:8080/namespace/name:v1")
				Expect(err).ToNot(HaveOccurred())
				Expect(image).ToNot(BeNil())
				Expect(image.GetHostname()).To(Equal("test.com"))
				Expect(image.GetPort()).To(Equal("8080"))
				Expect(image.HasIBMRepo()).To(BeFalse())
				Expect(image.GetRegistryURL()).To(Equal("https://test.com:8080"))
				Expect(image.GetContentTrustURL()).To(Equal("https://test.com:4443"))
				vaURL, err := image.GetVAURL()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(`Deny "test.com:8080/namespace/name:v1", Vulnerability Advisor is not supported for images from this registry: test.com`))
				Expect(vaURL).To(Equal(""))
				Expect(image.GetTag()).To(Equal("v1"))
				Expect(image.NameWithTag()).To(Equal("test.com:8080/namespace/name:v1"))
				Expect(image.NameWithoutTag()).To(Equal("test.com:8080/namespace/name"))
			})
		})

	})

}
