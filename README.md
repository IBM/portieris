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

Portieris is installed using a Helm chart. Before you begin, make sure that you have Kubernetes 1.9 or above, and Helm 2.8 or above (not Helm 3.x) installed in your cluster.

To install Portieris:

* Clone the Portieris Git repository to your workstation.
* Change directory into the Portieris Git repository.
* Run `./helm/portieris/gencerts <namespace>`. The `gencerts` script generates new SSL certificates and keys for Portieris. Portieris presents this certificates to the Kubernetes API server when the API server makes admission requests. If you do not generate new certificates, it could be possible for an attacker to spoof Portieris in your cluster.
* Run `helm upgrade --install portieris --set namespace=<namespace> helm/portieris`.

## Uninstalling Portieris

You can uninstall Portieris at any time by running `helm delete --purge portieris`. Note that all your image security policies are deleted when you uninstall Portieris.

## Image security policies

Image security policies define Portieris' behavior in your cluster. You must configure your own policies in order for Portieris to enforce your desired security posture. [Policies](POLICIES.md) are described separately.

## Reporting security issues

To report a security issue, DO NOT open an issue. Instead, send your report via email to alchreg@uk.ibm.com privately.
