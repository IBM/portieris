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

package v1beta1_test

import (
	. "github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {

	Describe("when there are not policies", func() {
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
						Spec: PolicySpec{
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
			It("Should return a policy but `trust` and `va` should be nil", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: PolicySpec{
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
				Expect(policy.Va.Enabled).To(BeNil())
			})
		})

		Context("but not policies", func() {
			It("Should return a policy but `trust` and `va` should be nil", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: PolicySpec{
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
				Expect(policy.Va.Enabled).To(BeNil())
			})
		})

		Context("but not `trust` policy", func() {
			It("Should return a policy but `trust` should be nil", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: PolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:notrust",
										Policy: Policy{
											Va: VA{
												Enabled: TruePointer,
											},
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
				Expect(policy.Va.Enabled).ToNot(BeNil())
				Expect(*policy.Va.Enabled).To(BeTrue())
			})
		})

		Context("but not `va` policy", func() {
			It("Should return a policy but `va` should be nil", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: PolicySpec{
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
				Expect(policy.Va.Enabled).To(BeNil())
				Expect(policy.Trust.Enabled).ToNot(BeNil())
				Expect(policy.Trust.TrustServer).ToNot(BeEmpty())
				Expect(*policy.Trust.Enabled).To(BeFalse())
			})
		})

		Context("`trust` and `va` policy are set", func() {
			It("Should return a policy and `trust` and `va` should be `true`", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: PolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:enabled",
										Policy: Policy{
											Trust: Trust{
												Enabled:     TruePointer,
												TrustServer: "https://some-trust-server.com",
											},
											Va: VA{
												Enabled: TruePointer,
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
				Expect(policy.Va.Enabled).ToNot(BeNil())
				Expect(*policy.Va.Enabled).To(BeTrue())
				Expect(policy.Trust.Enabled).ToNot(BeNil())
				Expect(policy.Trust.TrustServer).ToNot(BeEmpty())
				Expect(*policy.Trust.Enabled).To(BeTrue())
			})
		})

		Context("`trust` and `va` policy are set", func() {
			It("Should return a policy and `trust` and `va` should be `false`", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: PolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:disabled",
										Policy: Policy{
											Trust: Trust{
												Enabled: FalsePointer,
											},
											Va: VA{
												Enabled: FalsePointer,
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
				Expect(policy.Va.Enabled).ToNot(BeNil())
				Expect(*policy.Va.Enabled).To(BeFalse())
				Expect(policy.Trust.Enabled).ToNot(BeNil())
				Expect(policy.Trust.TrustServer).To(BeEmpty())
				Expect(*policy.Trust.Enabled).To(BeFalse())
			})
		})
	})

	Describe("when the image has a tag and the repo does not", func() {

		Context("Repository defined, but no policy", func() {
			It("Should find repo by adding `:*` and `trust` and `va` should be nil", func() {
				apl := ImagePolicyList{
					Items: []ImagePolicy{
						{
							Spec: PolicySpec{
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
				Expect(policy.Va.Enabled).To(BeNil())
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
						Spec: PolicySpec{
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
			It("Should return a policy but `trust` and `va` should be nil", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: PolicySpec{
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
				Expect(policy.Va.Enabled).To(BeNil())
			})
		})

		Context("but not policies", func() {
			It("Should return a policy but `trust` and `va` should be nil", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: PolicySpec{
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
				Expect(policy.Va.Enabled).To(BeNil())
			})
		})

		Context("but not `trust` policy", func() {
			It("Should return a policy but `trust` should be nil", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: PolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:notrust",
										Policy: Policy{
											Va: VA{
												Enabled: TruePointer,
											},
										},
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
				Expect(policy.Va.Enabled).ToNot(BeNil())
				Expect(*policy.Va.Enabled).To(BeTrue())
			})
		})

		Context("but not `va` policy", func() {
			It("Should return a policy but `va` should be nil", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: PolicySpec{
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
				policy := apl.FindClusterImagePolicy("test.com/namespace/hello:nova")
				Expect(policy).ToNot(BeNil())
				Expect(policy.Va.Enabled).To(BeNil())
				Expect(policy.Trust.Enabled).ToNot(BeNil())
				Expect(policy.Trust.TrustServer).ToNot(BeEmpty())
				Expect(*policy.Trust.Enabled).To(BeFalse())
			})
		})

		Context("`trust` and `va` policy are set", func() {
			It("Should return a policy and `trust` and `va` should be `true`", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: PolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:enabled",
										Policy: Policy{
											Trust: Trust{
												Enabled:     TruePointer,
												TrustServer: "https://some-trust-server.com",
											},
											Va: VA{
												Enabled: TruePointer,
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
				Expect(policy.Va.Enabled).ToNot(BeNil())
				Expect(*policy.Va.Enabled).To(BeTrue())
				Expect(policy.Trust.Enabled).ToNot(BeNil())
				Expect(policy.Trust.TrustServer).ToNot(BeEmpty())
				Expect(*policy.Trust.Enabled).To(BeTrue())
			})
		})

		Context("`trust` and `va` policy are set", func() {
			It("Should return a policy and `trust` and `va` should be `false`", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: PolicySpec{
								Repositories: []Repository{
									{
										Name: "test.com/namespace/hello:disabled",
										Policy: Policy{
											Trust: Trust{
												Enabled: FalsePointer,
											},
											Va: VA{
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
				Expect(policy.Va.Enabled).ToNot(BeNil())
				Expect(*policy.Va.Enabled).To(BeFalse())
				Expect(policy.Trust.Enabled).ToNot(BeNil())
				Expect(policy.Trust.TrustServer).To(BeEmpty())
				Expect(*policy.Trust.Enabled).To(BeFalse())
			})
		})
	})

	Describe("when the image has a tag and the repo does not", func() {

		Context("Repository defined, but no policy", func() {
			It("Should find repo by adding `:*` and `trust` and `va` should be nil", func() {
				apl := ClusterImagePolicyList{
					Items: []ClusterImagePolicy{
						{
							Spec: PolicySpec{
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
				Expect(policy.Va.Enabled).To(BeNil())
			})
		})
	})
})
