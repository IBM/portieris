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

package multi

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	securityenforcementfake "github.com/IBM/portieris/pkg/apis/securityenforcement/client/clientset/versioned/fake"
	securityenforcementv1beta1 "github.com/IBM/portieris/pkg/apis/securityenforcement/v1beta1"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/IBM/portieris/pkg/notary/fakenotary"
	"github.com/IBM/portieris/pkg/policy"
	"github.com/IBM/portieris/pkg/registry/fakeregistry"
	"github.com/IBM/portieris/pkg/webhook"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	k8sfake "k8s.io/client-go/kubernetes/fake"

	"testing"
)

func TestAdmissioncontroller(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Admissioncontroller Suite")
}

var _ = BeforeSuite(func() {
	tempTrustDir = "./test_trust"
})

var _ = AfterSuite(func() {
	os.RemoveAll(tempTrustDir)
})

var (
	tempTrustDir        string
	ctrl                *Controller
	kubeClientset       *k8sfake.Clientset
	kubeWrapper         kubernetes.WrapperInterface
	secClientset        *securityenforcementfake.Clientset
	policyClient        *policy.Client
	kubeObjects         []runtime.Object
	imageObjects        []runtime.Object
	clusterimageObjects []runtime.Object
	trust               *fakenotary.FakeNotary
	cr                  *fakeregistry.FakeRegistry
	wh                  *webhook.Server
)

// resetAllFakes should be call before any test
func resetAllFakes() {
	kubeObjects = []runtime.Object{}
	kubeClientset = k8sfake.NewSimpleClientset(kubeObjects...)
	kubeWrapper = kubernetes.NewKubeClientsetWrapper(kubeClientset)
	imageObjects = []runtime.Object{}
	secClientset = securityenforcementfake.NewSimpleClientset(imageObjects...)
	policyClient = policy.NewClient(secClientset)
	trust = &fakenotary.FakeNotary{}
	cr = &fakeregistry.FakeRegistry{}
	ctrl = NewController(kubeWrapper, policyClient, trust, cr)
	wh = webhook.NewServer("notary", ctrl, []byte{}, []byte{})
}

// newImagePolicyFromYAMLOrJSON .
func newImagePolicyFromYAMLOrJSON(payload *bytes.Buffer, namespace string) *securityenforcementv1beta1.ImagePolicy {
	decoder := yaml.NewYAMLOrJSONDecoder(payload, 50)
	policy := &securityenforcementv1beta1.ImagePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
		},
	}
	decoder.Decode(policy)
	return policy
}

// createNewImagePolicy creates a new ImagePolicy in the given namespace
func createNewImagePolicy(repo, namespace string) *securityenforcementv1beta1.ImagePolicy {
	return &securityenforcementv1beta1.ImagePolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: securityenforcementv1beta1.SchemeGroupVersion.String(),
			Kind:       "ImagePolicy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-test", namespace),
			Namespace: namespace,
		},
		Spec: securityenforcementv1beta1.PolicySpec{
			Repositories: []securityenforcementv1beta1.Repository{
				{
					Name: repo,
					Policy: securityenforcementv1beta1.Policy{
						Trust: securityenforcementv1beta1.Trust{
							Enabled:       securityenforcementv1beta1.FalsePointer,
							SignerSecrets: []securityenforcementv1beta1.Signer{},
						},
					},
				},
			},
		},
	}
}

// newClusterImagePolicyFromYAMLOrJSON .
func newClusterImagePolicyFromYAMLOrJSON(payload *bytes.Buffer, namespace string) *securityenforcementv1beta1.ClusterImagePolicy {
	decoder := yaml.NewYAMLOrJSONDecoder(payload, 50)
	policy := &securityenforcementv1beta1.ClusterImagePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
		},
	}
	err := decoder.Decode(policy)
	Expect(err).ToNot(HaveOccurred())
	return policy
}

// createClusterImagePolicy creates a new ClusterImagePolicy in the given namespace
func createClusterImagePolicy(repo, namespace string) *securityenforcementv1beta1.ClusterImagePolicy {
	return &securityenforcementv1beta1.ClusterImagePolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: securityenforcementv1beta1.SchemeGroupVersion.String(),
			Kind:       "ClusterImagePolicy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-test", namespace),
			Namespace: namespace,
		},
		Spec: securityenforcementv1beta1.PolicySpec{
			Repositories: []securityenforcementv1beta1.Repository{
				{
					Name: repo,
					Policy: securityenforcementv1beta1.Policy{
						Trust: securityenforcementv1beta1.Trust{
							Enabled:       securityenforcementv1beta1.FalsePointer,
							SignerSecrets: []securityenforcementv1beta1.Signer{},
						},
					},
				},
			},
		},
	}
}

func newFakeSecret(secretName, namespace, registry string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			".dockerconfigjson": []byte(
				fmt.Sprintf(`{
			  "auths": {
			    "%s": {
			      "username": "token",
			      "password": "registry-token",
			      "email": "email@email.com",
			      "auth": "auth-token"
			    }
			  }
			}`, registry)),
		},
	}
}

// newFakeRequest creates a new http request
func newFakeRequest(image string) *http.Request {
	// TODO: Delete what we don't need for unit tests
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(fmt.Sprintf(`
		{
		  "kind": "AdmissionReview",
		  "apiVersion": "admission.k8s.io/v1beta1",
		  "request": {
		    "uid": "ed782967-1c99-11e8-936d-08002789d446",
		    "kind": {
		      "group": "",
		      "version": "v1",
		      "kind": "Pod"
		    },
		    "resource": {
		      "group": "",
		      "version": "v1",
		      "resource": "pods"
		    },
		    "namespace": "default",
		    "operation": "CREATE",
		    "userInfo": {
		      "username": "minikube-user",
		      "groups": [
		        "system:masters",
		        "system:authenticated"
		      ]
		    },
		    "object": {
		      "metadata": {
		        "name": "nginx",
		        "namespace": "default",
		        "creationTimestamp": null
		      },
		      "spec": {
		        "volumes": [
		          {
		            "name": "default-token-xff4f",
		            "secret": {
		              "secretName": "default-token-xff4f"
		            }
		          }
		        ],
		        "containers": [
		          {
		            "name": "nginx",
		            "image": "%s",
		            "ports": [
		              {
		                "hostPort": 8080,
		                "containerPort": 8080,
		                "protocol": "TCP"
		              }
		            ],
		            "resources": {},
		            "volumeMounts": [
		              {
		                "name": "default-token-xff4f",
		                "readOnly": true,
		                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
		              }
		            ],
		            "terminationMessagePath": "/dev/termination-log",
		            "terminationMessagePolicy": "File",
		            "imagePullPolicy": "Always"
		          }
		        ],
		        "restartPolicy": "Always",
		        "terminationGracePeriodSeconds": 30,
		        "dnsPolicy": "ClusterFirst",
		        "serviceAccountName": "default",
		        "serviceAccount": "default",
		        "hostNetwork": true,
		        "securityContext": {},
		        "imagePullSecrets": [
		          {
								"name": "regsecret"
		          }
		        ],
		        "schedulerName": "default-scheduler",
		        "tolerations": [
		          {
		            "key": "node.kubernetes.io/not-ready",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          },
		          {
		            "key": "node.kubernetes.io/unreachable",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          }
		        ]
		      },
		      "status": {}
		    },
		    "oldObject": null
		  }
		}`, image)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func newFakeRequestMulitpleSecretsBadSecond(image string) *http.Request {
	// TODO: Delete what we don't need for unit tests
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(fmt.Sprintf(`
		{
		  "kind": "AdmissionReview",
		  "apiVersion": "admission.k8s.io/v1beta1",
		  "request": {
		    "uid": "ed782967-1c99-11e8-936d-08002789d446",
		    "kind": {
		      "group": "",
		      "version": "v1",
		      "kind": "Pod"
		    },
		    "resource": {
		      "group": "",
		      "version": "v1",
		      "resource": "pods"
		    },
		    "namespace": "default",
		    "operation": "CREATE",
		    "userInfo": {
		      "username": "minikube-user",
		      "groups": [
		        "system:masters",
		        "system:authenticated"
		      ]
		    },
		    "object": {
		      "metadata": {
		        "name": "nginx",
		        "namespace": "default",
		        "creationTimestamp": null
		      },
		      "spec": {
		        "volumes": [
		          {
		            "name": "default-token-xff4f",
		            "secret": {
		              "secretName": "default-token-xff4f"
		            }
		          }
		        ],
		        "containers": [
		          {
		            "name": "nginx",
		            "image": "%s",
		            "ports": [
		              {
		                "hostPort": 8080,
		                "containerPort": 8080,
		                "protocol": "TCP"
		              }
		            ],
		            "resources": {},
		            "volumeMounts": [
		              {
		                "name": "default-token-xff4f",
		                "readOnly": true,
		                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
		              }
		            ],
		            "terminationMessagePath": "/dev/termination-log",
		            "terminationMessagePolicy": "File",
		            "imagePullPolicy": "Always"
		          }
		        ],
		        "restartPolicy": "Always",
		        "terminationGracePeriodSeconds": 30,
		        "dnsPolicy": "ClusterFirst",
		        "serviceAccountName": "default",
		        "serviceAccount": "default",
		        "hostNetwork": true,
		        "securityContext": {},
		        "imagePullSecrets": [
							{
								"name": "regsecret"
							},
							{
								"name": "badregsecret"
							}
		        ],
		        "schedulerName": "default-scheduler",
		        "tolerations": [
		          {
		            "key": "node.kubernetes.io/not-ready",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          },
		          {
		            "key": "node.kubernetes.io/unreachable",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          }
		        ]
		      },
		      "status": {}
		    },
		    "oldObject": null
		  }
		}`, image)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func newFakeRequestMulitpleSecrets(image string) *http.Request {
	// TODO: Delete what we don't need for unit tests
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(fmt.Sprintf(`
		{
		  "kind": "AdmissionReview",
		  "apiVersion": "admission.k8s.io/v1beta1",
		  "request": {
		    "uid": "ed782967-1c99-11e8-936d-08002789d446",
		    "kind": {
		      "group": "",
		      "version": "v1",
		      "kind": "Pod"
		    },
		    "resource": {
		      "group": "",
		      "version": "v1",
		      "resource": "pods"
		    },
		    "namespace": "default",
		    "operation": "CREATE",
		    "userInfo": {
		      "username": "minikube-user",
		      "groups": [
		        "system:masters",
		        "system:authenticated"
		      ]
		    },
		    "object": {
		      "metadata": {
		        "name": "nginx",
		        "namespace": "default",
		        "creationTimestamp": null
		      },
		      "spec": {
		        "volumes": [
		          {
		            "name": "default-token-xff4f",
		            "secret": {
		              "secretName": "default-token-xff4f"
		            }
		          }
		        ],
		        "containers": [
		          {
		            "name": "nginx",
		            "image": "%s",
		            "ports": [
		              {
		                "hostPort": 8080,
		                "containerPort": 8080,
		                "protocol": "TCP"
		              }
		            ],
		            "resources": {},
		            "volumeMounts": [
		              {
		                "name": "default-token-xff4f",
		                "readOnly": true,
		                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
		              }
		            ],
		            "terminationMessagePath": "/dev/termination-log",
		            "terminationMessagePolicy": "File",
		            "imagePullPolicy": "Always"
		          }
		        ],
		        "restartPolicy": "Always",
		        "terminationGracePeriodSeconds": 30,
		        "dnsPolicy": "ClusterFirst",
		        "serviceAccountName": "default",
		        "serviceAccount": "default",
		        "hostNetwork": true,
		        "securityContext": {},
		        "imagePullSecrets": [
							{
								"name": "badregsecret"
							},
							{
								"name": "regsecret"
							}
		        ],
		        "schedulerName": "default-scheduler",
		        "tolerations": [
		          {
		            "key": "node.kubernetes.io/not-ready",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          },
		          {
		            "key": "node.kubernetes.io/unreachable",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          }
		        ]
		      },
		      "status": {}
		    },
		    "oldObject": null
		  }
		}`, image)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func newFakeRequestMultipleValidSecrets(image string) *http.Request {
	// TODO: Delete what we don't need for unit tests
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(fmt.Sprintf(`
		{
		  "kind": "AdmissionReview",
		  "apiVersion": "admission.k8s.io/v1beta1",
		  "request": {
		    "uid": "ed782967-1c99-11e8-936d-08002789d446",
		    "kind": {
		      "group": "",
		      "version": "v1",
		      "kind": "Pod"
		    },
		    "resource": {
		      "group": "",
		      "version": "v1",
		      "resource": "pods"
		    },
		    "namespace": "default",
		    "operation": "CREATE",
		    "userInfo": {
		      "username": "minikube-user",
		      "groups": [
		        "system:masters",
		        "system:authenticated"
		      ]
		    },
		    "object": {
		      "metadata": {
		        "name": "nginx",
		        "namespace": "default",
		        "creationTimestamp": null
		      },
		      "spec": {
		        "volumes": [
		          {
		            "name": "default-token-xff4f",
		            "secret": {
		              "secretName": "default-token-xff4f"
		            }
		          }
		        ],
		        "containers": [
		          {
		            "name": "nginx",
		            "image": "%s",
		            "ports": [
		              {
		                "hostPort": 8080,
		                "containerPort": 8080,
		                "protocol": "TCP"
		              }
		            ],
		            "resources": {},
		            "volumeMounts": [
		              {
		                "name": "default-token-xff4f",
		                "readOnly": true,
		                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
		              }
		            ],
		            "terminationMessagePath": "/dev/termination-log",
		            "terminationMessagePolicy": "File",
		            "imagePullPolicy": "Always"
		          }
		        ],
		        "restartPolicy": "Always",
		        "terminationGracePeriodSeconds": 30,
		        "dnsPolicy": "ClusterFirst",
		        "serviceAccountName": "default",
		        "serviceAccount": "default",
		        "hostNetwork": true,
		        "securityContext": {},
		        "imagePullSecrets": [
							{
								"name": "regsecret1"
							},
							{
								"name": "regsecret"
							}
		        ],
		        "schedulerName": "default-scheduler",
		        "tolerations": [
		          {
		            "key": "node.kubernetes.io/not-ready",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          },
		          {
		            "key": "node.kubernetes.io/unreachable",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          }
		        ]
		      },
		      "status": {}
		    },
		    "oldObject": null
		  }
		}`, image)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// newFakeRequest creates a new http request
func newFakeRequestMultiContainer(image, image1 string) *http.Request {
	// TODO: Delete what we don't need for unit tests
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(fmt.Sprintf(`
		{
		  "kind": "AdmissionReview",
		  "apiVersion": "admission.k8s.io/v1beta1",
		  "request": {
		    "uid": "ed782967-1c99-11e8-936d-08002789d446",
		    "kind": {
		      "group": "",
		      "version": "v1",
		      "kind": "Pod"
		    },
		    "resource": {
		      "group": "",
		      "version": "v1",
		      "resource": "pods"
		    },
		    "namespace": "default",
		    "operation": "CREATE",
		    "userInfo": {
		      "username": "minikube-user",
		      "groups": [
		        "system:masters",
		        "system:authenticated"
		      ]
		    },
		    "object": {
		      "metadata": {
		        "name": "nginx",
		        "namespace": "default",
		        "creationTimestamp": null
		      },
		      "spec": {
		        "volumes": [
		          {
		            "name": "default-token-xff4f",
		            "secret": {
		              "secretName": "default-token-xff4f"
		            }
		          }
		        ],
		        "containers": [
		          {
		            "name": "nginx",
		            "image": "%s",
		            "ports": [
		              {
		                "hostPort": 8080,
		                "containerPort": 8080,
		                "protocol": "TCP"
		              }
		            ],
		            "resources": {},
		            "volumeMounts": [
		              {
		                "name": "default-token-xff4f",
		                "readOnly": true,
		                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
		              }
		            ],
		            "terminationMessagePath": "/dev/termination-log",
		            "terminationMessagePolicy": "File",
		            "imagePullPolicy": "Always"
							},
							{
		            "name": "statsd",
		            "image": "%s",
		            "ports": [
		              {
		                "hostPort": 8080,
		                "containerPort": 8080,
		                "protocol": "TCP"
		              }
		            ],
		            "resources": {},
		            "volumeMounts": [
		              {
		                "name": "default-token-xff4f",
		                "readOnly": true,
		                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
		              }
		            ],
		            "terminationMessagePath": "/dev/termination-log",
		            "terminationMessagePolicy": "File",
		            "imagePullPolicy": "Always"
		          }
		        ],
		        "restartPolicy": "Always",
		        "terminationGracePeriodSeconds": 30,
		        "dnsPolicy": "ClusterFirst",
		        "serviceAccountName": "default",
		        "serviceAccount": "default",
		        "hostNetwork": true,
		        "securityContext": {},
		        "imagePullSecrets": [
		          {
		            "name": "regsecret"
		          }
		        ],
		        "schedulerName": "default-scheduler",
		        "tolerations": [
		          {
		            "key": "node.kubernetes.io/not-ready",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          },
		          {
		            "key": "node.kubernetes.io/unreachable",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          }
		        ]
		      },
		      "status": {}
		    },
		    "oldObject": null
		  }
		}`, image, image1)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// newFakeRequest creates a new http request
func newFakeRequestMultiContainerMultiSecret(image, image1 string) *http.Request {
	// TODO: Delete what we don't need for unit tests
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(fmt.Sprintf(`
		{
		  "kind": "AdmissionReview",
		  "apiVersion": "admission.k8s.io/v1beta1",
		  "request": {
		    "uid": "ed782967-1c99-11e8-936d-08002789d446",
		    "kind": {
		      "group": "",
		      "version": "v1",
		      "kind": "Pod"
		    },
		    "resource": {
		      "group": "",
		      "version": "v1",
		      "resource": "pods"
		    },
		    "namespace": "default",
		    "operation": "CREATE",
		    "userInfo": {
		      "username": "minikube-user",
		      "groups": [
		        "system:masters",
		        "system:authenticated"
		      ]
		    },
		    "object": {
		      "metadata": {
		        "name": "nginx",
		        "namespace": "default",
		        "creationTimestamp": null
		      },
		      "spec": {
		        "volumes": [
		          {
		            "name": "default-token-xff4f",
		            "secret": {
		              "secretName": "default-token-xff4f"
		            }
		          }
		        ],
		        "containers": [
		          {
		            "name": "nginx",
		            "image": "%s",
		            "ports": [
		              {
		                "hostPort": 8080,
		                "containerPort": 8080,
		                "protocol": "TCP"
		              }
		            ],
		            "resources": {},
		            "volumeMounts": [
		              {
		                "name": "default-token-xff4f",
		                "readOnly": true,
		                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
		              }
		            ],
		            "terminationMessagePath": "/dev/termination-log",
		            "terminationMessagePolicy": "File",
		            "imagePullPolicy": "Always"
							},
							{
		            "name": "statsd",
		            "image": "%s",
		            "ports": [
		              {
		                "hostPort": 8080,
		                "containerPort": 8080,
		                "protocol": "TCP"
		              }
		            ],
		            "resources": {},
		            "volumeMounts": [
		              {
		                "name": "default-token-xff4f",
		                "readOnly": true,
		                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
		              }
		            ],
		            "terminationMessagePath": "/dev/termination-log",
		            "terminationMessagePolicy": "File",
		            "imagePullPolicy": "Always"
		          }
		        ],
		        "restartPolicy": "Always",
		        "terminationGracePeriodSeconds": 30,
		        "dnsPolicy": "ClusterFirst",
		        "serviceAccountName": "default",
		        "serviceAccount": "default",
		        "hostNetwork": true,
		        "securityContext": {},
		        "imagePullSecrets": [
							{
								"name": "regsecret"
							},
							{
								"name": "regsecret3"
							}
		        ],
		        "schedulerName": "default-scheduler",
		        "tolerations": [
		          {
		            "key": "node.kubernetes.io/not-ready",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          },
		          {
		            "key": "node.kubernetes.io/unreachable",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          }
		        ]
		      },
		      "status": {}
		    },
		    "oldObject": null
		  }
		}`, image, image1)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// newFakeRequest creates a new http request
func newFakeRequestInitContainer(initImage, image string) *http.Request {
	// TODO: Delete what we don't need for unit tests
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(fmt.Sprintf(`
		{
		  "kind": "AdmissionReview",
		  "apiVersion": "admission.k8s.io/v1beta1",
		  "request": {
		    "uid": "ed782967-1c99-11e8-936d-08002789d446",
		    "kind": {
		      "group": "",
		      "version": "v1",
		      "kind": "Pod"
		    },
		    "resource": {
		      "group": "",
		      "version": "v1",
		      "resource": "pods"
		    },
		    "namespace": "default",
		    "operation": "CREATE",
		    "userInfo": {
		      "username": "minikube-user",
		      "groups": [
		        "system:masters",
		        "system:authenticated"
		      ]
		    },
		    "object": {
		      "metadata": {
		        "name": "nginx",
		        "namespace": "default",
		        "creationTimestamp": null
		      },
		      "spec": {
		        "volumes": [
		          {
		            "name": "default-token-xff4f",
		            "secret": {
		              "secretName": "default-token-xff4f"
		            }
		          }
		        ],
                "initContainers" : [
				{
		            "name": "nginx",
		            "image": "%s",
		            "ports": [
		              {
		                "hostPort": 8080,
		                "containerPort": 8080,
		                "protocol": "TCP"
		              }
		            ],
		            "resources": {},
		            "volumeMounts": [
		              {
		                "name": "default-token-xff4f",
		                "readOnly": true,
		                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
		              }
		            ],
		            "terminationMessagePath": "/dev/termination-log",
		            "terminationMessagePolicy": "File",
		            "imagePullPolicy": "Always"
					}
                ],
		        "containers": [
                 {
		            "name": "statsd",
		            "image": "%s",
		            "ports": [
		              {
		                "hostPort": 8080,
		                "containerPort": 8080,
		                "protocol": "TCP"
		              }
		            ],
		            "resources": {},
		            "volumeMounts": [
		              {
		                "name": "default-token-xff4f",
		                "readOnly": true,
		                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
		              }
		            ],
		            "terminationMessagePath": "/dev/termination-log",
		            "terminationMessagePolicy": "File",
		            "imagePullPolicy": "Always"
		          }
		        ],
		        "restartPolicy": "Always",
		        "terminationGracePeriodSeconds": 30,
		        "dnsPolicy": "ClusterFirst",
		        "serviceAccountName": "default",
		        "serviceAccount": "default",
		        "hostNetwork": true,
		        "securityContext": {},
		        "imagePullSecrets": [
		          {
		            "name": "regsecret"
		          }
		        ],
		        "schedulerName": "default-scheduler",
		        "tolerations": [
		          {
		            "key": "node.kubernetes.io/not-ready",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          },
		          {
		            "key": "node.kubernetes.io/unreachable",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          }
		        ]
		      },
		      "status": {}
		    },
		    "oldObject": null
		  }
		}`, initImage, image)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// newFakeRequest creates a new http request
func newFakeRequestDeployment(image string) *http.Request {
	// TODO: Delete what we don't need for unit tests
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(fmt.Sprintf(`
		{
		  "kind": "AdmissionReview",
		  "apiVersion": "admission.k8s.io/v1beta1",
		  "request": {
		    "uid": "ed782967-1c99-11e8-936d-08002789d446",
		    "kind": {
		      "group": "apps",
		      "version": "v1",
		      "kind": "Deployment"
		    },
		    "resource": {
		      "group": "apps",
		      "version": "v1",
		      "resource": "deployments"
		    },
		    "namespace": "default",
		    "operation": "CREATE",
		    "userInfo": {
		      "username": "minikube-user",
		      "groups": [
		        "system:masters",
		        "system:authenticated"
		      ]
		    },
		    "object": {
		      "metadata": {
		        "name": "nginx",
		        "namespace": "default",
		        "creationTimestamp": null
		      },
		      "spec": {
						"replicas":1,
						"template":{
							"spec":{
								"volumes": [
									{
										"name": "default-token-xff4f",
										"secret": {
											"secretName": "default-token-xff4f"
		            }
		          }
							],
							"containers": [
										{
											"name": "statsd",
											"image": "%s",
											"ports": [
												{
													"hostPort": 8080,
													"containerPort": 8080,
													"protocol": "TCP"
												}
												],
												"resources": {},
												"volumeMounts": [
													{
														"name": "default-token-xff4f",
														"readOnly": true,
														"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
													}
													],
													"terminationMessagePath": "/dev/termination-log",
													"terminationMessagePolicy": "File",
													"imagePullPolicy": "Always"
												}
												],
												"restartPolicy": "Always",
												"terminationGracePeriodSeconds": 30,
												"dnsPolicy": "ClusterFirst",
												"serviceAccountName": "default",
												"serviceAccount": "default",
												"hostNetwork": true,
												"securityContext": {},
												"imagePullSecrets": [
													{
														"name": "regsecret"
													}
													],
													"schedulerName": "default-scheduler",
													"tolerations": [
														{
															"key": "node.kubernetes.io/not-ready",
															"operator": "Exists",
															"effect": "NoExecute",
															"tolerationSeconds": 300
														},
														{
															"key": "node.kubernetes.io/unreachable",
															"operator": "Exists",
															"effect": "NoExecute",
															"tolerationSeconds": 300
														}
														]
													},
													"status": {}
												},
												"oldObject": null
											}
										}
									}
								}`, image)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func newFakeRequestDeploymentWithZeroReplicas(image string) *http.Request {
	// TODO: Delete what we don't need for unit tests
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(fmt.Sprintf(`
		{
		  "kind": "AdmissionReview",
		  "apiVersion": "admission.k8s.io/v1beta1",
		  "request": {
		    "uid": "ed782967-1c99-11e8-936d-08002789d446",
		    "kind": {
		      "group": "apps",
		      "version": "v1",
		      "kind": "Deployment"
		    },
		    "resource": {
		      "group": "apps",
		      "version": "v1",
		      "resource": "deployments"
		    },
		    "namespace": "default",
		    "operation": "CREATE",
		    "userInfo": {
		      "username": "minikube-user",
		      "groups": [
		        "system:masters",
		        "system:authenticated"
		      ]
		    },
		    "object": {
		      "metadata": {
		        "name": "nginx",
		        "namespace": "default",
		        "creationTimestamp": null
		      },
		      "spec": {
						"replicas":0,
						"template":{
							"spec":{
								"volumes": [
									{
										"name": "default-token-xff4f",
										"secret": {
											"secretName": "default-token-xff4f"
		            }
		          }
							],
							"containers": [
										{
											"name": "statsd",
											"image": "%s",
											"ports": [
												{
													"hostPort": 8080,
													"containerPort": 8080,
													"protocol": "TCP"
												}
												],
												"resources": {},
												"volumeMounts": [
													{
														"name": "default-token-xff4f",
														"readOnly": true,
														"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
													}
													],
													"terminationMessagePath": "/dev/termination-log",
													"terminationMessagePolicy": "File",
													"imagePullPolicy": "Always"
												}
												],
												"restartPolicy": "Always",
												"terminationGracePeriodSeconds": 30,
												"dnsPolicy": "ClusterFirst",
												"serviceAccountName": "default",
												"serviceAccount": "default",
												"hostNetwork": true,
												"securityContext": {},
												"imagePullSecrets": [
													{
														"name": "regsecret"
													}
													],
													"schedulerName": "default-scheduler",
													"tolerations": [
														{
															"key": "node.kubernetes.io/not-ready",
															"operator": "Exists",
															"effect": "NoExecute",
															"tolerationSeconds": 300
														},
														{
															"key": "node.kubernetes.io/unreachable",
															"operator": "Exists",
															"effect": "NoExecute",
															"tolerationSeconds": 300
														}
														]
													},
													"status": {}
												},
												"oldObject": null
											}
										}
									}
								}`, image)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// newFakeRequest creates a new http request
func newFakeRequestWithParents(image string) *http.Request {
	// TODO: Delete what we don't need for unit tests
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(fmt.Sprintf(`
		{
		  "kind": "AdmissionReview",
		  "apiVersion": "admission.k8s.io/v1beta1",
		  "request": {
		    "uid": "ed782967-1c99-11e8-936d-08002789d446",
		    "kind": {
		      "group": "",
		      "version": "v1",
		      "kind": "Pod"
		    },
		    "resource": {
		      "group": "",
		      "version": "v1",
		      "resource": "pods"
		    },
		    "namespace": "default",
		    "operation": "CREATE",
		    "userInfo": {
		      "username": "minikube-user",
		      "groups": [
		        "system:masters",
		        "system:authenticated"
		      ]
		    },
		    "object": {
		      "metadata": {
		        "name": "nginx",
		        "namespace": "default",
						"creationTimestamp": null,
						"ownerReferences":[{"apiVersion":"extensions/v1beta1","kind":"ReplicaSet","name":"deployment-55d687c698","uid":"e0577bcf-30dd-11e8-83d1-baaf52c27f02","controller":true,"blockOwnerDeletion":true}]
		      },
		      "spec": {
		        "volumes": [
		          {
		            "name": "default-token-xff4f",
		            "secret": {
		              "secretName": "default-token-xff4f"
		            }
		          }
		        ],
		        "containers": [
		          {
		            "name": "nginx",
		            "image": "%s",
		            "ports": [
		              {
		                "hostPort": 8080,
		                "containerPort": 8080,
		                "protocol": "TCP"
		              }
		            ],
		            "resources": {},
		            "volumeMounts": [
		              {
		                "name": "default-token-xff4f",
		                "readOnly": true,
		                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
		              }
		            ],
		            "terminationMessagePath": "/dev/termination-log",
		            "terminationMessagePolicy": "File",
		            "imagePullPolicy": "Always"
		          }
		        ],
		        "restartPolicy": "Always",
		        "terminationGracePeriodSeconds": 30,
		        "dnsPolicy": "ClusterFirst",
		        "serviceAccountName": "default",
		        "serviceAccount": "default",
		        "hostNetwork": true,
		        "securityContext": {},
		        "imagePullSecrets": [
		          {
		            "name": "regsecret"
		          }
		        ],
		        "schedulerName": "default-scheduler",
		        "tolerations": [
		          {
		            "key": "node.kubernetes.io/not-ready",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          },
		          {
		            "key": "node.kubernetes.io/unreachable",
		            "operator": "Exists",
		            "effect": "NoExecute",
		            "tolerationSeconds": 300
		          }
		        ]
		      },
		      "status": {}
		    },
		    "oldObject": null
		  }
		}`, image)))
	req.Header.Set("Content-Type", "application/json")
	return req
}
