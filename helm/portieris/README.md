# Container Image Security Enforcement

This chart installs Container Image Security Enforcement for IBM Cloud Container Service in your cluster.

## Prerequisites

* Kubernetes v1.9+
* Tiller v2.8+

## Chart details

This chart:
* Installs Container Image Security Enforcement.
* Configures Kubernetes admission webhooks to direct admission requests to Container Image Security Enforcement.
* Adds a resource definition for security policies.
* Adds a default cluster-wide security policy, and a default security policy in the kube-system and ibm-system Kubernetes namespaces.

## Installing the chart

```
helm repo add ibm https://registry.stage1.ng.bluemix.net/helm/ibm
helm install -n cise ibm/ibmcloud-image-enforcement
```

For full installation instructions, see [Installing security enforcement in your cluster](https://console.bluemix.net/docs/services/Registry/registry_security_enforce.html#sec_enforcer_install).

## Default security policies

This chart installs default security policies in your cluster. You can modify the default policies or replace them with your own. For more information, see [Default policies](https://console.bluemix.net/docs/services/Registry/registry_security_enforce.html#default_policies).

You can apply access control policies to limit who can modify Image Security Enforcement policies in your cluster. See [Controlling who can customize policies](https://console.bluemix.net/docs/services/Registry/registry_security_enforce.html#assign_user_policy).

## Customizing security policies

You can add your own security policies, scoped to a Kubernetes namespace or the entire cluster. Cluster policies are used when no namespace scoped policies exist in the Kubernetes namespace you are deploying to.

For information about configuring security policies, and an explanation of the security policy resources, see [Customizing policies](https://console.bluemix.net/docs/services/Registry/registry_security_enforce.html#customize_policies).

## Removing the chart

1. Container Image Security Enforcement uses Hyperkube to remove some configuration from your cluster when you remove it. Before you can remove Container Image Security Enforcement, you must make sure that Hyperkube is allowed to run. Make sure that the policy for the ibm-system namespace allows the `hyperkube` image.
    ```yaml
    - name: quay.io/coreos/hyperkube
      policies:
    ```
    Alternatively, disable Container Image Security Enforcement entirely.
    ```
    kubectl delete MutatingWebhookConfiguration image-admission-config 
    kubectl delete ValidatingWebhookConfiguration image-admission-config
    ```
2. Remove the resource definitions for your security policies. When you delete the resource definitions, your security policies are also deleted.
    ```
    kubectl delete crd clusterimagepolicies.securityenforcement.admission.cloud.ibm.com imagepolicies.securityenforcement.admission.cloud.ibm.com
    ```
3. Remove the chart.
    ```
    helm delete --purge cise
    ```
