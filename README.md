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

## Installing Portieris

Portieris is installed using a Helm chart. Before you begin, make sure that you have Kubernetes 1.9 or above, and Helm 2.8 or above installed in your cluster.

To install Portieris:

* Clone the Portieris Git repository to your workstation.
* Change directory into the Portieris Git repository.
<<<<<<< HEAD
* Run `./helm/portieris/gencerts <namespace>`. The `gencerts` script generates new SSL certificates and keys for Portieris. Portieris presents this certificates to the Kubernetes API server when the API server makes admission requests. If you do not generate new certificates, it could be possible for an attacker to spoof Portieris in your cluster.
* Run `helm upgrade --install --name portieris -n portieris --set namespace=<namespace> helm/portieris` (when using `helm3`, the namespace has to exist before running the command).
=======
* Run `./helm/portieris/gencerts`. The `gencerts` script generates new SSL certificates and keys for Portieris. Portieris presents this certificates to the Kubernetes API server when the API server makes admission requests. If you do not generate new certificates, it could be possible for an attacker to spoof Portieris in your cluster.
* [Future feature] If your environment uses dockerconfig at the host to authenticate to your private repository, update the value of secrets.yaml under helm/templates directory
* Run `helm install -n portieris helm/portieris`.
>>>>>>> Issue 51:

## Uninstalling Portieris

You can uninstall Portieris at any time using `helm delete --purge portieris`. Note that all your image security policies are deleted when you uninstall Portieris.

## Image security policies

Image security policies define Portieris' behavior in your cluster. There are two types of policy:

* ImagePolicies can be configured in each Kubernetes namespace, and define Portieris' behavior in that namespace. If an ImagePolicy exists in a namespace, the policies from that namespace are used, even if the ImagePolicy does not have a matching policy for a given image. If a namespace does not have an ImagePolicy, the ClusterImagePolicy is used.
* ClusterImagePolicies are configured at the cluster level, and take effect whenever there is no ImagePolicy in the namespace where the workload is being deployed.

## Configuring image security policies

You can configure custom security policies to control what images can be deployed in your Kubernetes namespaces, and to enforce trust pinning of particular signers.  
Portieris supports sourcing trust data from the following registries without additional configuration in the image policy:
* IBM Cloud Container Registry
* Quay.io
* Docker Hub

To use a different trust server for a repository, you can specify the `trustServer` parameter in your policy:
*Example*
```yaml
apiVersion: securityenforcement.admission.cloud.ibm.com/v1beta1
kind: ImagePolicy
metadata:
  name: allow-custom
spec:
   repositories:
    - name: "icr.io/*"
      policy:
        trust:
          enabled: true
          trustServer: "https://icr.io:4443" # Optional, custom trust server for repository
```  
For more information, see the [IBM Cloud docs](https://cloud.ibm.com/docs/services/Registry?topic=registry-security_enforce#customize_policies).
## Configuring access controls for your security policies

You can configure Kubernetes RBAC rules to define which users and applications have the ability to modify your security policies. For more information, see the [IBM Cloud docs](https://cloud.ibm.com/docs/services/Registry?topic=registry-security_enforce#assign_user_policy).

## Reporting security issues

To report a security issue, DO NOT open an issue. Instead, send your report via email to alchreg@uk.ibm.com privately.
