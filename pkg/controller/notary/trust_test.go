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
	"fmt"

	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/IBM/portieris/pkg/notary/fakenotary"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	notaryclient "github.com/theupdateframework/notary/client"
	"github.com/theupdateframework/notary/tuf/data"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("Trust", func() {

	BeforeEach(func() {
		resetAllFakes()
	})

	Describe("getSecretToken", func() {
		var (
			fakeRepo         *fakenotary.FakeRepository
			fakeErrorMessage string
			fakeErr          error
			server           string
			image            string
			notaryToken      string
			targetName       string
			signerPublicKey  string
		)

		BeforeEach(func() {
			fakeRepo = &fakenotary.FakeRepository{}
			fakeErrorMessage = "FAKE_ERROR"
			fakeErr = fmt.Errorf("%s", fakeErrorMessage)
			server = "http://localhost"
			image = "test"
			notaryToken = "token"
			targetName = "target"
			signerPublicKey = `
-----BEGIN PUBLIC KEY-----
MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgHkUcBWbd/CLtcDiAiAjfWScAkcM
InK8c+vl8puAwTNYPWU/sK5xk+W+i/Bw2mb6G+VB9udHAPqH/AJgHnBMBuVdo4g+
/G05zIZY03pN1fKl/DCzuckoWVkRUO7grAV5euHcxKwppDa1AmvxcjSGoH8D+zeh
Aaq77T1fhZkFO5uTAgMBAAE=
-----END PUBLIC KEY-----
`
		})

		It("should return an error if it fails to get the repo", func() {
			trust.GetNotaryRepoReturns(nil, fakeErr)
			ctrl = NewController(kubeWrapper, policyClient, trust, cr)
			_, err := ctrl.getDigest(server, image, notaryToken, targetName, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fakeErrorMessage))
		})

		It("should return an error if it fails to get the target by name", func() {
			fakeRepo.GetAllTargetMetadataByNameReturns(nil, fakeErr)
			trust.GetNotaryRepoReturns(fakeRepo, nil)
			ctrl = NewController(kubeWrapper, policyClient, trust, cr)
			_, err := ctrl.getDigest(server, image, notaryToken, targetName, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fakeErrorMessage))
		})

		It("should return an error if there are not targets", func() {
			fakeRepo.GetAllTargetMetadataByNameReturns([]notaryclient.TargetSignedStruct{}, nil)
			trust.GetNotaryRepoReturns(fakeRepo, nil)
			ctrl = NewController(kubeWrapper, policyClient, trust, cr)
			_, err := ctrl.getDigest(server, image, notaryToken, targetName, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("No signed targets found"))
		})

		Context("when there is not a list of signers", func() {
			It("should return the digest if there is a target but not required signers", func() {
				publicKey := data.NewPublicKey("sha256", []byte("abc"))
				fakeRepo.GetAllTargetMetadataByNameReturns([]notaryclient.TargetSignedStruct{
					{
						Target: notaryclient.Target{
							Hashes: data.Hashes{"sha256": []byte("1234567890")},
						},
						Role: data.DelegationRole{
							BaseRole: data.BaseRole{
								Name: "targets/releases",
								Keys: map[string]data.PublicKey{"abc": publicKey},
							},
						},
					},
				}, nil)
				trust.GetNotaryRepoReturns(fakeRepo, nil)
				ctrl = NewController(kubeWrapper, policyClient, trust, cr)
				digest, err := ctrl.getDigest(server, image, notaryToken, targetName, nil)
				Expect(err).ToNot(HaveOccurred())
				Expect(digest.String()).To(Equal("31323334353637383930"))
			})
		})

		Context("when there is a list of signers", func() {
			It("should fail if it can't parse a public key", func() {
				publicKey := data.NewPublicKey("sha256", []byte("abc"))
				fakeRepo.GetAllTargetMetadataByNameReturns([]notaryclient.TargetSignedStruct{
					{
						Target: notaryclient.Target{
							Hashes: data.Hashes{"sha256": []byte("1234567890")},
						},
						Role: data.DelegationRole{
							BaseRole: data.BaseRole{
								Name: "targets/wibble",
								Keys: map[string]data.PublicKey{"261144b64ca3413e7fb3fd509099f1b92df19d4e4158e709fbaa2f8fc22f7191": publicKey},
							},
						},
					},
					{
						Target: notaryclient.Target{
							Hashes: data.Hashes{"sha256": []byte("1234567890")},
						},
						Role: data.DelegationRole{
							BaseRole: data.BaseRole{
								Name: "targets/releases",
								Keys: map[string]data.PublicKey{"abc": publicKey},
							},
						},
					},
				}, nil)
				trust.GetNotaryRepoReturns(fakeRepo, nil)
				ctrl = NewController(kubeWrapper, policyClient, trust, cr)
				_, err := ctrl.getDigest(server, image, notaryToken, targetName, []Signer{
					{
						signer:    "wibble",
						publicKey: "invalid signer public key",
					},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("no valid public key found"))
			})

			It("should fail if the key ids are not the same", func() {
				publicKey := data.NewPublicKey("sha256", []byte("abc"))
				fakeRepo.GetAllTargetMetadataByNameReturns([]notaryclient.TargetSignedStruct{
					{
						Target: notaryclient.Target{
							Hashes: data.Hashes{"sha256": []byte("1234567890")},
						},
						Role: data.DelegationRole{
							BaseRole: data.BaseRole{
								Name: "targets/wibble",
								Keys: map[string]data.PublicKey{"different key id": publicKey},
							},
						},
					},
					{
						Target: notaryclient.Target{
							Hashes: data.Hashes{"sha256": []byte("1234567890")},
						},
						Role: data.DelegationRole{
							BaseRole: data.BaseRole{
								Name: "targets/releases",
								Keys: map[string]data.PublicKey{"abc": publicKey},
							},
						},
					},
				}, nil)
				trust.GetNotaryRepoReturns(fakeRepo, nil)
				ctrl = NewController(kubeWrapper, policyClient, trust, cr)
				_, err := ctrl.getDigest(server, image, notaryToken, targetName, []Signer{
					{
						signer:    "wibble",
						publicKey: signerPublicKey,
					},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Public keys are different"))
			})

			It("should fail if the signer does not have a public key in the role", func() {
				publicKey := data.NewPublicKey("sha256", []byte("abc"))
				fakeRepo.GetAllTargetMetadataByNameReturns([]notaryclient.TargetSignedStruct{
					{
						Target: notaryclient.Target{
							Hashes: data.Hashes{"sha256": []byte("1234567890")},
						},
						Role: data.DelegationRole{
							BaseRole: data.BaseRole{
								Name: "targets/wibble",
								Keys: map[string]data.PublicKey{"different key id": publicKey},
							},
						},
					},
					{
						Target: notaryclient.Target{
							Hashes: data.Hashes{"sha256": []byte("1234567890")},
						},
						Role: data.DelegationRole{
							BaseRole: data.BaseRole{
								Name: "targets/releases",
								Keys: map[string]data.PublicKey{"abc": publicKey},
							},
						},
					},
				}, nil)
				trust.GetNotaryRepoReturns(fakeRepo, nil)
				ctrl = NewController(kubeWrapper, policyClient, trust, cr)
				_, err := ctrl.getDigest(server, image, notaryToken, targetName, []Signer{
					{
						signer: "wibble",
					},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("PublicKey not found in role wibble"))
			})

			It("should fail if the signer does not have a role", func() {
				publicKey := data.NewPublicKey("sha256", []byte("abc"))
				fakeRepo.GetAllTargetMetadataByNameReturns([]notaryclient.TargetSignedStruct{
					{
						Target: notaryclient.Target{
							Hashes: data.Hashes{"sha256": []byte("1234567890")},
						},
						Role: data.DelegationRole{
							BaseRole: data.BaseRole{
								Name: "targets/wibble",
								Keys: map[string]data.PublicKey{"different key id": publicKey},
							},
						},
					},
					{
						Target: notaryclient.Target{
							Hashes: data.Hashes{"sha256": []byte("1234567890")},
						},
						Role: data.DelegationRole{
							BaseRole: data.BaseRole{
								Name: "targets/releases",
								Keys: map[string]data.PublicKey{"abc": publicKey},
							},
						},
					},
				}, nil)
				trust.GetNotaryRepoReturns(fakeRepo, nil)
				ctrl = NewController(kubeWrapper, policyClient, trust, cr)
				_, err := ctrl.getDigest(server, image, notaryToken, targetName, []Signer{
					{
						// signer: "wibble",
						publicKey: signerPublicKey,
					},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no signature found for role"))
			})

			It("should return a digest", func() {
				publicKey := data.NewPublicKey("sha256", []byte("abc"))
				fakeRepo.GetAllTargetMetadataByNameReturns([]notaryclient.TargetSignedStruct{
					{
						Target: notaryclient.Target{
							Hashes: data.Hashes{"sha256": []byte("1234567890")},
						},
						Role: data.DelegationRole{
							BaseRole: data.BaseRole{
								Name: "targets/wibble",
								Keys: map[string]data.PublicKey{"261144b64ca3413e7fb3fd509099f1b92df19d4e4158e709fbaa2f8fc22f7191": publicKey},
							},
						},
					},
					{
						Target: notaryclient.Target{
							Hashes: data.Hashes{"sha256": []byte("1234567890")},
						},
						Role: data.DelegationRole{
							BaseRole: data.BaseRole{
								Name: "targets/releases",
								Keys: map[string]data.PublicKey{"whatever, don't care": publicKey},
							},
						},
					},
				}, nil)
				trust.GetNotaryRepoReturns(fakeRepo, nil)
				ctrl = NewController(kubeWrapper, policyClient, trust, cr)
				digest, err := ctrl.getDigest(server, image, notaryToken, targetName, []Signer{
					{
						signer:    "wibble",
						publicKey: signerPublicKey,
					},
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(digest.String()).To(Equal("31323334353637383930"))
			})

		})

	})

	Describe("getSignerSecret", func() {

		var (
			namespace  string
			secretName string
			fakeSecret *corev1.Secret
		)

		BeforeEach(func() {
			namespace = metav1.NamespaceDefault
			secretName = "my-secret"
			fakeSecret = &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
					Namespace: namespace,
				},
				Data: map[string][]byte{
					".dockerconfigjson": []byte(`{}`),
				},
			}
		})

		It("should return an error if there is not secret", func() {
			ctrl = NewController(kubeWrapper, policyClient, trust, cr)
			signer, err := ctrl.getSignerSecret(namespace, "no-secret")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`secrets "no-secret" not found`))
			Expect(signer).To(Equal(Signer{}))
		})

		It("should return an error if the name is empty", func() {
			fakeSecret.Data = map[string][]byte{
				"name":      []byte(""),
				"publicKey": []byte("key"),
			}
			kubeClientset = k8sfake.NewSimpleClientset(fakeSecret)
			kubeWrapper = kubernetes.NewKubeClientsetWrapper(kubeClientset)
			ctrl = NewController(kubeWrapper, policyClient, trust, cr)
			signer, err := ctrl.getSignerSecret(namespace, "my-secret")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("name or publicKey field in secret my-secret is empty"))
			Expect(signer).To(Equal(Signer{}))
		})

		It("should return an error if the publicKey is empty", func() {
			fakeSecret.Data = map[string][]byte{
				"name":      []byte("signer"),
				"publicKey": []byte(""),
			}
			kubeClientset = k8sfake.NewSimpleClientset(fakeSecret)
			kubeWrapper = kubernetes.NewKubeClientsetWrapper(kubeClientset)
			ctrl = NewController(kubeWrapper, policyClient, trust, cr)
			signer, err := ctrl.getSignerSecret(namespace, "my-secret")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("name or publicKey field in secret my-secret is empty"))
			Expect(signer).To(Equal(Signer{}))
		})

		It("should return the signer if the name and publicKey are not empty", func() {
			fakeSecret.Data = map[string][]byte{
				"name":      []byte("signer"),
				"publicKey": []byte("key"),
			}
			kubeClientset = k8sfake.NewSimpleClientset(fakeSecret)
			kubeWrapper = kubernetes.NewKubeClientsetWrapper(kubeClientset)
			ctrl = NewController(kubeWrapper, policyClient, trust, cr)
			signer, err := ctrl.getSignerSecret(namespace, "my-secret")
			Expect(err).ToNot(HaveOccurred())
			Expect(signer).To(Equal(Signer{signer: "signer", publicKey: "key"}))
		})
	})

})
