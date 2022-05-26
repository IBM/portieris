---

copyright:
  years: 2018, 2021
lastupdated: "2021-08-26"

---

# Portieris

Use this chart to install Portieris in your cluster.

## Prerequisites

* Kubernetes v1.16, or later
* Helm v3.0, or later

## Chart details

This chart does the following tasks:

* Installs Portieris.
* Adds a resource definition for security policies.
* Adds some default security policies.
* Configures Kubernetes admission webhooks to direct admission requests to Portieris.

## Installing the chart

### Regenerate certificates

This installation uses the default certificates. To avoid using the default certificates, check out the source project at the release level and run the `./gencerts` script.

**Important** If you don't run the `./gencerts` script, you will deploy with certificates that are publicly accessible on GitHub.

### IBM Cloud Kubernetes Service

If you're deploying onto an IBM Cloud cluster, Portieris automatically creates policies that deploy the various Kubernetes components and includes a policy rule that allows all images without verification. Change the policy that allows everything because it is insecure, but keep the IBM Cloud specific policies.

```
helm install --create-namespace -n portieris --set PolicySet=IKS .
```

### Other Kubernetes clusters

If you're deploying onto a generic cluster, Portieris does not create default policies, that allows all images without verification. This policy prevents Portieris from stopping you deploying to your cluster. Update this policy to something more restrictive.

```
helm install --create-namespace -n portieris . --set PolicySet=None --debug
```

For installation instructions, see [Installing Portieris](https://github.com/IBM/portieris/blob/master/README.md#installing-portieris).

## Default security policies

This chart installs default security policies in your cluster. You can modify the default policies or replace them with your own. The default policies are defined in the [chart policies templates](https://github.com/IBM/portieris/blob/master/helm/portieris/templates/policies.yaml).

Apply access control policies to limit who can modify Portieris policies in your cluster, see [Configuring access controls for your security policies](https://github.com/IBM/portieris/blob/master/README.md#configuring-access-controls-for-your-security-policies).

## Customizing security policies

You can add your own security policies, scoped to a Kubernetes namespace or the entire cluster. Cluster policies are used when no namespace scoped policies exist in the Kubernetes namespace that you are deploying to.

For more information about configuring security policies, and an explanation of the security policy resources, see [Portieris policies](https://github.com/IBM/portieris/blob/master/POLICIES.md).

## Removing the chart

To remove the chart, run the following command.

    ```
    helm delete -n portieris portieris
    ```
