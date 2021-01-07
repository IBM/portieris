// Copyright 2018, 2020 Portieris Authors.
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

package v1_test

import (
	. "github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {

	Describe("when there are no policies", func() {
		It("policy should be nil", func() {
			apl := ImagePolicyList{}
			policy := apl.FindImagePolicy("test.com/hello")
			Expect(policy).To(BeNil())
		})
	})

	Describe("when the image does not have a matching repo", func() {
		It("policy should be nil", func() {
			apl := ImagePolicyList{
				Items: []ImagePolicy{
					{
						Spec: ImagePolicySpec{
							Repositories: []Repository{
								{
									Name: "test.com/*",
								},
							},
						},
					},
				},
			}
			policy := apl.FindImagePolicy("unknown.com/hello")
			Expect(policy).To(BeNil())
		})
	})

	Describe("when the image has a matching repo", func() {
		Context("but not policies", func() {
			It("Should return a policy but `trust.enable` should be nil and `simple.requirements` should be empty", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/*",
									},
								},
							},
						},
					},
				}
				policy := apl.FindImagePolicy("test.com/hello")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).To(BeNil())
				Expect(policy.Trust.TrustServer).To(BeEmpty())
				Expect(policy.Simple.Requirements).To(BeEmpty())
			})
		})

		Context("but not policies", func() {
			It("Should return a policy but `trust.enabled`/`vulnerability.ibmva.enabled` should be nil and `simple.requirements` should be empty", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name:   "test.com/namespace/hello:nopolicies",
										Policy: Policy{},
									},
								},
							},
						},
					},
				}
				policy := apl.FindImagePolicy("test.com/namespace/hello:nopolicies")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).To(BeNil())
				Expect(policy.Trust.TrustServer).To(BeEmpty())
				Expect(policy.Simple.Requirements).To(BeEmpty())
				Expect(policy.Vulnerability.ICCRVA.Enabled).To(BeNil())
			})
		})

		Context("but not `trust` policy", func() {
			It("Should return a policy but `trust.enabled` should be nil", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:notrust",
										Policy: Policy{
											Simple: Simple{Requirements: []SimpleRequirement{{}}},
										},
									},
								},
							},
						},
					},
				}
				policy := apl.FindImagePolicy("test.com/namespace/hello:notrust")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).To(BeNil())
				Expect(policy.Trust.TrustServer).To(BeEmpty())
				Expect(policy.Simple.Requirements).NotTo(BeEmpty())
			})
		})

		Context("but not `simple` policy", func() {
			It("Should return a policy but `simple.requirements` should be empty", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:nova",
										Policy: Policy{
											Trust: Trust{
												Enabled:     FalsePointer,
												TrustServer: "https://some-trust-server.com",
											},
										},
									},
								},
							},
						},
					},
				}
				policy := apl.FindImagePolicy("test.com/namespace/hello:nova")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).ToNot(BeNil())
				Expect(policy.Trust.TrustServer).ToNot(BeEmpty())
				Expect(*policy.Trust.Enabled).To(BeFalse())
				Expect(policy.Simple.Requirements).To(BeEmpty())
			})
		})

		Context("`trust` and `simple` policy are set", func() {
			It("Should return a policy and `trust.enabled` should be `true` and `simple.requirements` should be set", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:enabled",
										Policy: Policy{
											Trust: Trust{
												Enabled:     TruePointer,
												TrustServer: "https://some-trust-server.com",
											},
											Simple: Simple{
												Requirements: []SimpleRequirement{{}},
											},
										},
									},
								},
							},
						},
					},
				}
				policy := apl.FindImagePolicy("test.com/namespace/hello:enabled")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).ToNot(BeNil())
				Expect(policy.Trust.TrustServer).ToNot(BeEmpty())
				Expect(*policy.Trust.Enabled).To(BeTrue())
				Expect(policy.Simple.Requirements).NotTo(BeEmpty())
			})
		})

		Context("`trust` policy are enabled", func() {
			It("Should return a policy and `trust.enabled` should be true", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:disabled",
										Policy: Policy{
											Trust: Trust{
												Enabled: TruePointer,
											},
										},
									},
								},
							},
						},
					},
				}
				policy := apl.FindImagePolicy("test.com/namespace/hello:disabled")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).ToNot(BeNil())
				Expect(policy.Trust.TrustServer).To(BeEmpty())
				Expect(*policy.Trust.Enabled).To(BeTrue())
				Expect(policy.Simple.Requirements).To(BeEmpty())
			})
		})

		Context("`vulnerability.ICCRVA` policy is enabled", func() {
			It("Should return a policy and `vulnerability.ICCRVA.enabled` should be true", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:disabled",
										Policy: Policy{
											Vulnerability: Vulnerability{ICCRVA{
												Enabled: TruePointer,
												Account: "123",
											}},
										},
									},
								},
							},
						},
					},
				}
				policy := apl.FindImagePolicy("test.com/namespace/hello:disabled")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).To(BeNil())
				Expect(policy.Simple.Requirements).To(BeEmpty())
				Expect(policy.Vulnerability.ICCRVA.Enabled).ToNot(BeNil())
				Expect(policy.Vulnerability.ICCRVA.Account).To(Equal("123"))
			})
		})
	})

	Describe("when the image has a tag and the repo does not", func() {

		Context("Repository defined, but no policy", func() {
			It("Should find repo by adding `:*` and `trust.enabled` should be nil and `simple.requirements` should be empty", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/image",
									},
								},
							},
						},
					},
				}
				policy := apl.FindImagePolicy("test.com/namespace/image:tag")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).To(BeNil())
				Expect(policy.Trust.TrustServer).To(BeEmpty())
				Expect(policy.Simple.Requirements).To(BeEmpty())
			})
		})
	})

	Describe("when there are not cluster policies", func() {
		It("Should not fail but the policy should be nil", func() {
			apl := ClusterImagePolicyList{}
			policy := apl.FindClusterImagePolicy("test.com/hello")
			Expect(policy).To(BeNil())
		})
	})

	Describe("when the image does not have a matching repo", func() {
		It("Should not fail but the policy should be nil", func() {
			apl := ClusterImagePolicyList{
				Items: []ClusterImagePolicy{
					{
						Spec: ImagePolicySpec{
							Repositories: []Repository{
								{
									Name: "test.com/*",
								},
							},
						},
					},
				},
			}
			policy := apl.FindClusterImagePolicy("unknown.com/hello")
			Expect(policy).To(BeNil())
		})
	})

	Describe("when the image has a matching repo", func() {

		Context("but not policies", func() {
			It("Should return a policy but `trust.eabled` `simple.requirements` `vcheck.type` should be nil/empty", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/*",
									},
								},
							},
						},
					},
				}
				policy := apl.FindClusterImagePolicy("test.com/hello")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).To(BeNil())
				Expect(policy.Trust.TrustServer).To(BeEmpty())
				Expect(policy.Simple.Requirements).To(BeEmpty())
			})
		})

		Context("but not policies", func() {
			It("Should return a policy but `trust.enabled` should be nil `simple.requirements` should be empty", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name:   "test.com/namespace/hello:nopolicies",
										Policy: Policy{},
									},
								},
							},
						},
					},
				}
				policy := apl.FindClusterImagePolicy("test.com/namespace/hello:nopolicies")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).To(BeNil())
				Expect(policy.Trust.TrustServer).To(BeEmpty())
				Expect(policy.Simple.Requirements).To(BeEmpty())
			})
		})

		Context("but not `trust` or `simple` policy", func() {
			It("Should return a policy but `trust.enabled` should be nil and `simple.requirements` should be empty", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name:   "test.com/namespace/hello:notrust",
										Policy: Policy{},
									},
								},
							},
						},
					},
				}
				policy := apl.FindClusterImagePolicy("test.com/namespace/hello:notrust")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).To(BeNil())
				Expect(policy.Trust.TrustServer).To(BeEmpty())
				Expect(policy.Simple.Requirements).To(BeEmpty())
			})
		})

		Context("both trust and simple policy", func() {
			It("Should return a policy", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:nova",
										Policy: Policy{
											Trust: Trust{
												Enabled:     TruePointer,
												TrustServer: "https://some-trust-server.com",
											},
											Simple: Simple{
												Requirements: []SimpleRequirement{{}},
											},
										},
									},
								},
							},
						},
					},
				}
				policy := apl.FindClusterImagePolicy("test.com/namespace/hello:nova")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).ToNot(BeNil())
				Expect(policy.Trust.TrustServer).ToNot(BeEmpty())
				Expect(*policy.Trust.Enabled).To(BeTrue())
				Expect(policy.Simple.Requirements).NotTo(BeEmpty())
			})
		})

		Context("`trust` policy are set", func() {
			It("Should return a policy and `trust.enabled` should be `true`", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:enabled",
										Policy: Policy{
											Trust: Trust{
												Enabled:     TruePointer,
												TrustServer: "https://some-trust-server.com",
											},
										},
									},
								},
							},
						},
					},
				}
				policy := apl.FindClusterImagePolicy("test.com/namespace/hello:enabled")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).ToNot(BeNil())
				Expect(policy.Trust.TrustServer).ToNot(BeEmpty())
				Expect(*policy.Trust.Enabled).To(BeTrue())
			})
		})

		Context("`trust.enabled` is false", func() {
			It("Should return a policy and `trust.enabled` should be `false`", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:disabled",
										Policy: Policy{
											Trust: Trust{
												Enabled: FalsePointer,
											},
										},
									},
								},
							},
						},
					},
				}
				policy := apl.FindClusterImagePolicy("test.com/namespace/hello:disabled")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).ToNot(BeNil())
				Expect(policy.Trust.TrustServer).To(BeEmpty())
				Expect(*policy.Trust.Enabled).To(BeFalse())
			})
		})
	})

	Describe("when the image has a tag and the repo does not", func() {
		Context("Repository defined, but no policy", func() {
			It("Should find repo by adding `:*` and `trust.enabled` should be nil and `vcheck.type` should be empty", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: ImagePolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/image",
									},
								},
							},
						},
					},
				}
				policy := apl.FindClusterImagePolicy("test.com/namespace/image:tag")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Trust.Enabled).To(BeNil())
				Expect(policy.Trust.TrustServer).To(BeEmpty())
			})
		})
	})
})
