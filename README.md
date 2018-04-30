# Portieris

![Portieris logo](./logos/text_and_logo.svg)

[![Travis badge](https://api.travis-ci.org/IBM/portieris.svg?branch=master)](https://travis-ci.org/IBM/portieris)

Portieris is a Kubernetes admission controller for enforcing Content Trust. You can create image security policies for each Kubernetes namespace, or at the cluster level, and enforce different levels of trust for different images.

## How it works

Portieris uses a Kubernetes Mutating Admission Webhook to modify your Kubernetes resources at the point of creation, to ensure that Kubernetes pulls the signed version. When configured to do so, it enforces trust pinning, and blocks the creation of resources that use untrusted images.

If your cloud provider provides a [Notary](https://github.com/theupdateframework/notary) server (sometimes referred to as Content Trust), Portieris accesses trust data in that Notary server that corresponds to the image that you are deploying.

When you create or edit a workload, the Kubernetes API server sends a request to Portieris. The AdmissionRequest contains the content of your workload. For each image in your workload, Portieris finds a matching security policy. If trust enforcement is enabled in your policy, Portieris pulls signature information for your image from the corresponding Notary server and, if a signed version of the image exists, creates a JSON patch to edit the image name in the workload to the signed image by digest. If a signer is defined in the policy, Portieris additionally checks that the image is signed by the specified role, and verifies that the specified key was used to sign the image.

If any image in your workload is not signed where trust is enforced, or is not signed by the correct role when a signer is defined, the entire workload is prevented from deploying.

Portieris receives AdmissionRequests for creation of or edits to all types of workload. To prevent Portieris from impacting auto-recovery, it approves requests where a parent exists.

Portieris' Admission Webhook is configured to fail closed. Three instances of Portieris make sure that it is able to approve its own upgrades and auto-recovery. If all instance of Portieris are unavailable, Kubernetes will not auto-recover it, and you must delete the MutatingAdmissionWebhook to allow Portieris to recover.

## Current limitations

Portieris only supports sourcing trust data from the IBM Cloud Notary servers. You can deploy images from other servers while Portieris is installed, but you must disable trust enforcement for these images. The default ClusterImagePolicy enforces trust disabled for "*". We are looking to support additional Notary servers. If you want a particular Notary server to be supported, [raise an issue](https://github.com/ibm/portieris/issues/new).

## Installing Portieris

Portieris is installed using a Helm chart. Before you begin, make sure that you have Kubernetes 1.9 or above, and Helm 2.8 or above installed in your cluster.

To install Portieris:

* Clone the Portieris Git repository to your workstation.
* Change directory into the Portieris Git repository, then run `helm install -n portieris helm/portieris`.

## Uninstalling Portieris

You can uninstall Portieris at any time using `helm delete --purge portieris`. Note that all your image security policies are deleted when you uninstall Portieris.

## Image security policies

Image security policies define Portieris' behavior in your cluster. There are two types of policy:

* ImagePolicies can be configured in each Kubernetes namespace, and define Portieris' behavior in that namespace. If an ImagePolicy exists in a namespace, the policies from that namespace are used, even if the ImagePolicy does not have a matching policy for a given image. If a namespace does not have an ImagePolicy, the ClusterImagePolicy is used.
* ClusterImagePolicies are configured at the cluster level, and take effect whenever there is no ImagePolicy in the namespace where the workload is being deployed.

## Configuring image security policies

You can configure custom security policies to control what images can be deployed in your Kubernetes namespaces, and to enforce trust pinning of particular signers. For more information, see the [IBM Cloud docs](https://console.bluemix.net/docs/services/Registry/registry_security_enforce.html#customize_policies).

## Configuring access controls for your security policies

You can configure Kubernetes RBAC rules to define which users and applications have the ability to modify your security policies. For more information, see the [IBM Cloud docs](https://console.bluemix.net/docs/services/Registry/registry_security_enforce.html#assign_user_policy).

## Reporting security issues

To report a security issue, DO NOT open an issue. Instead, send your report via email to alchreg@uk.ibm.com privately.
