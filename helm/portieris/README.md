---

copyright:
  years: 2018, 2021
lastupdated: "2021-02-19"

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
helm install --create-namespace -n portieris .
```

### Other Kubernetes clusters

If you're deploying onto a generic cluster, Portieris automatically creates a policy that allows all images without verification. This policy prevents Portieris from stopping you deploying to your cluster. Update this policy to something more restrictive.

```
helm install --create-namespace -n portieris . --set IBMContainerService=false --debug
```

For installation instructions, see [Installing security enforcement in your cluster](https://cloud.ibm.com/docs/services/Registry?topic=registry-security_enforce#sec_enforcer_install).

## Default security policies

This chart installs default security policies in your cluster. Modify the default policies or replace them with your own. For more information, see [Default policies](https://cloud.ibm.com/docs/services/Registry?topic=registry-security_enforce#default_policies).

Apply access control policies to limit who can modify Portieris policies in your cluster. See [Controlling who can customize policies](https://cloud.ibm.com/docs/services/Registry?topic=registry-security_enforce#assign_user_policy).

## Customizing security policies

You can add your own security policies, scoped to a Kubernetes namespace or the entire cluster. Cluster policies are used when no namespace scoped policies exist in the Kubernetes namespace that you are deploying to.

For more information about configuring security policies, and an explanation of the security policy resources, see [Customizing policies](https://cloud.ibm.com/docs/services/Registry?topic=registry-security_enforce#customize_policies).

## Removing the chart

To remove the chart, run the following command.

    ```
    helm delete -n portieris portieris
    ```
