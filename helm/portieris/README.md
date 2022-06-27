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

### Define certificates

You will need the following:
* A CA certificate
* A server certificate, signed by the CA
* A private key for the server certificate

Note: If you want, you can use the `gencerts` script to generate these automatically for you.

Next, you will need to add the contents of these three files to your values.json. Under `UseGeneratedCerts` you will add them to `tlsCert`, `tlsKey`, and `caCert` respectively. Also be sure to set `enabled` to true. The resulting section of the values file should look something like this:
```
UseGeneratedCerts:
  enabled: true
  tlsCert: |
    -----BEGIN CERTIFICATE-----
    ...certificate data...
    -----END CERTIFICATE-----
  tlsKey: |
    -----BEGIN CERTIFICATE-----
    ...key data...
    -----END CERTIFICATE-----
  caCert: |
    -----BEGIN CERTIFICATE-----
    ...certificate data...
    -----END CERTIFICATE-----
```
Include every line of the certificate/key files, making sure you use proper indentation before each line so as to produce a valid YAML.

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
