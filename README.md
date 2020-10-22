![Portieris logo](./logos/text_and_logo.svg)

[![Travis badge](https://api.travis-ci.org/IBM/portieris.svg?branch=master)](https://travis-ci.org/IBM/portieris)

Portieris is a Kubernetes admission controller for the enforcment of image security policies. You can create image security policies for each Kubernetes namespace, or at the cluster level, and enforce different rules for different images.

## How it works

Portieris uses a Kubernetes Mutating Admission Webhook to modify your Kubernetes resources, at the point of creation, to ensure that Kubernetes runs only policy compliant images. When configured to do so, it can enforce Docker Content Trust with optional trust pinning, or can verify signatures created using RedHat's simple signing model and will prevent the creation of resources using untrusted or unverified images.

If your cloud provider provides a [Notary](https://github.com/theupdateframework/notary) server (sometimes referred to as Content Trust), Portieris accesses trust data in that Notary server that corresponds to the image that you are deploying. In order to verify RedHat simple signatures they must be accessible via registry extension APIs or a configured signature store.

When you create or edit a workload, the Kubernetes API server sends a request to Portieris. The AdmissionRequest contains the content of your workload. For each image in your workload, Portieris finds a matching security policy.

If trust enforcement is enabled in the policy, Portieris pulls signature information for your image from the corresponding Notary server and, if a signed version of the image exists, creates a JSON patch to edit the image name in the workload to the signed image by digest. If a signer is defined in the policy, Portieris additionally checks that the image is signed by the specified role, and verifies that the specified key was used to sign the image.

If simple signing is specified by the policy, Portieris will verify the signature using using the public key and identity rules supplied in the policy and if verified similarly mutates the image name to a digest reference to ensure that concurrent tag changes cannot influence the image being pulled.

While it is possible to require both Notary trust and simple signing, the two methods must agree on the signed digest for the image. If the two methods return different signed digests, the image is denied. It is not possible to allow alternative signing methods.

If any image in your workload does not satisfy the policy the entire workload is prevented from deploying.

Portieris receives AdmissionRequests for creation of or edits to all types of workload. To prevent Portieris from impacting auto-recovery, it approves requests where a parent exists.

Portieris' Admission Webhook is configured to fail closed. Three instances of Portieris make sure that it is able to approve its own upgrades and auto-recovery. If all instance of Portieris are unavailable, Kubernetes will not auto-recover it, and you must delete the MutatingAdmissionWebhook to allow Portieris to recover.

## Installing Portieris

Portieris is installed using a Helm chart. Before you begin, make sure that you have Kubernetes 1.16 or above on your cluster and Helm 3.0 or above installed on your workstation.

To install Portieris in the default namespace (portieris):

* Clone the Portieris Git repository to your workstation. `git clone git@github.com:IBM/portieris.git`
* Change directory into the Portieris Git repository.
* Checkout the tag commit that you want to install , example: `git checkout 0.7.0`
* Run `./helm/portieris/gencerts`. The `gencerts` script generates new SSL certificates and keys for Portieris. Portieris presents this certificates to the Kubernetes API server when the API server makes admission requests. If you do not generate new certificates, it could be possible for an attacker to spoof Portieris in your cluster.
* Run `helm install portieris --create-namespace --namespace portieris helm/portieris`. 

You can use a different namespace if you choose including an existing one if you omit the `--create-namespace` option, but note that the namespace forms part of the webhook certificate common name so you need to generate the certificate for the target namespace.

* Run `./helm/portieris/gencerts <namespace>`.
* Run `helm install portieris --create-namespace --namespace <namespace> helm/portieris`.

To manage certificates through an installed [cert-manager](https://cert-manager.io/):

* Run `helm install portieris --set UseCertManager=true helm/portieris`.

By default, Portieris' admission webhook runs in all namespaces including its own install namespace, so that Portieris is able to review all the pods in the cluster. However, this can prevent the cluster from self healing in the event where Portieris becomes unavailable. Portieris also supports skipping namespaces with a certain label set. You can enable this by adding `--set AllowAdmissionSkip=true` to your install command, but make sure to control who can add labels to namespaces and who can access namespaces with this label so that a malicious party cannot use this label to bypass Portieris.

## Uninstalling Portieris

You can uninstall Portieris, at any time, by running `helm delete portieris --namespace <namespace>`.

**Uninstall Notes**: 
1. All your image security policies are deleted when you uninstall Portieris.
1. If you no longer need the namespace it will have to be manually deleted. i.e. `kubectl delete namespace/<namespace>`
1. If you have issues uninstalling portieris, via helm, try running the cleanup script: `helm/cleanup.sh portieris <namespace>`

## Image security policies

Image security policies define Portieris' behavior in your cluster. You must configure your own policies in order for Portieris to enforce your desired security posture. [Policies](POLICIES.md) are described separately.

## Configuring access controls for your security policies

You can configure Kubernetes RBAC rules to define which users and applications have the ability to modify your security policies. For more information, see the [IBM Cloud docs](https://cloud.ibm.com/docs/services/Registry?topic=registry-security_enforce#assign_user_policy).

If Portieris is installed with `AllowAdmissionSkip=true`, you can prevent Portieris' admission webhook from being called in specific namespaces by labelling the namespace with `securityenforcement.admission.cloud.ibm.com/namespace: skip`. Doing so would allow pods in that namespace to recover when the admission webhook is down, but note that no policies are applied in that namespace. For example, the Portieris install namespace is configured with this label to allow Portieris itself to recover when it is down. Make sure to control who can add labels to namespaces and who can access namespaces with this label so that a malicious party cannot use this label to bypass Portieris.

## Reporting security issues

To report a security issue, DO NOT open an issue. Instead, send your report via email to alchreg@uk.ibm.com privately.
