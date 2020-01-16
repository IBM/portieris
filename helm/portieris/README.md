# Portieris

This chart installs Portieris in your cluster.

## Prerequisites

* Kubernetes v1.9+
* Tiller v2.8+

## Chart details

This chart:
* Installs Portieris.
* Adds a resource definition for security policies.
* Adds a default cluster-wide security policy, and a default security policy in the kube-system and ibm-system Kubernetes namespaces.
* Configures Kubernetes admission webhooks to direct admission requests to Portieris.

## Installing the chart

### !!! Regenerate Certs !!!
The install will use the default certs if you do not run the gencerts script. **This means you will deploying with certs that are publically accessible on GitHub.**

Note that this script requires https://github.com/mikefarah/yq to be installed. There are competing versions of yq with different syntax, so check that you have the right one. If `yq r` is a valid command, you have the right one. If `yq -r` is valid, you have a competing version that is incompatible with this script.

```
./gencerts
```

### IBM Cloud Container Service

If you're deploying onto an IBM Cloud cluster Portieris automatically creates policies to allow the various Kubernetes components to be deployed as well as a policy rule to allow all images without verification. The allow everything should be changed because it is insecure but the IBM Cloud specific policies should be kept.

```
helm install -n portieris .
```

### Other Kubernetes Clusters

If you're deploying onto a generic cluster Portieris automatically creates a policy to allow all images without verification. This is to prevent Portieris from preventing you deploying to your cluster. You should update this policy to something more restrictive.

```
helm install -n portieris . --set IBMContainerService=false --debug
```

For full installation instructions, see [Installing security enforcement in your cluster](https://cloud.ibm.com/docs/services/Registry?topic=registry-security_enforce#sec_enforcer_install).

## Default security policies

This chart installs default security policies in your cluster. You should modify the default policies or replace them with your own. For more information, see [Default policies](https://cloud.ibm.com/docs/services/Registry?topic=registry-security_enforce#default_policies).

You should apply access control policies to limit who can modify Portieris policies in your cluster. See [Controlling who can customize policies](https://cloud.ibm.com/docs/services/Registry?topic=registry-security_enforce#assign_user_policy).

## Customizing security policies

You can add your own security policies, scoped to a Kubernetes namespace or the entire cluster. Cluster policies are used when no namespace scoped policies exist in the Kubernetes namespace you are deploying to.

For information about configuring security policies, and an explanation of the security policy resources, see [Customizing policies](https://cloud.ibm.com/docs/services/Registry?topic=registry-security_enforce#customize_policies).

## Removing the chart

1. Portieris uses Hyperkube to remove some configuration from your cluster when you remove it. Before you can remove Portieris, you must make sure that Hyperkube is allowed to run. Make sure that the policy for the ibm-system namespace allows the `hyperkube` image.
    ```yaml
    - name: quay.io/coreos/hyperkube
      policies:
    ```
    Alternatively, disable Portieris manually.
    ```
    kubectl delete MutatingWebhookConfiguration image-admission-config 
    ```
2. Remove the chart.
    ```
    helm delete --purge portieris
    ```
