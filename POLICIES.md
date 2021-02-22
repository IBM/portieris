---

copyright:
  years: 2020, 2021
lastupdated: "2021-02-18"

---

# Portieris policies

## Image policy resources

Portieris defines two custom resource types for policy. Image policy resources and cluster image policy resources.

For both types of resource, if multiple resources exist, they are merged together and can be protected by a role-based access control (RBAC) policy.

* Image policy resources, `ImagePolicy`, are configured in a Kubernetes namespace and define Portieris' behavior in that namespace. If image policy resources exist in a namespace, the policies from those image policy resources are used exclusively. If no match exists for the workload image in `ImagePolicy`, cluster image policy resources, `ClusterImagePolicy`, are not examined. Images in deployed workloads are wildcard matched against the set of policies defined, if no policy matches the workload image, deployment is denied.

  The following example allows any image from the `icr.io` registry with no further checks (the policy is empty).
 
  ```yaml
  apiVersion: portieris.cloud.ibm.com/v1
  kind: ImagePolicy
  metadata:
    name: allow-all-icrio
  spec:
     repositories:
      - name: "icr.io/*"
        policy:
  ```

* Cluster image policy resources, `ClusterImagePolicy`, are configured at the cluster level, and take effect whenever no image policy resource, `ImagePolicy`, is defined in the namespace where the workload is deployed. These cluster image policy resources have the same structure as namespace image policy resources and, if no matching policy is found for an image, deployment is denied.

  The following example allows all images from all registries with no checks.

  ```yaml
  apiVersion: portieris.cloud.ibm.com/v1
  kind: ClusterImagePolicy
  metadata:
    name: portieris-default-cluster-image-policy
  spec:
    repositories:
    - name: '*'
      policy:
  ```

## Installation default policies

Default policies are installed when Portieris is installed. You must review and change these according to your requirements.
You must customise the installation [default policies](helm/portieris/templates/default/policies.yaml).

## Repository matching

When an image is evaluated for admission, the set of policies is wildcard matched on the repository name. If multiple matches are found, the most specific match is used.

## Policy

A policy consists of an array of objects that define requirements on the image by using either `trust:` (Docker Content Trust and Notary v1), `simple:` (Red Hat Simple Signing), or `vulnerability:` objects.

**Important** If your policy was developed before Portieris v0.10.0, the policy has the API version: `apiVersion: securityenforcement.admission.cloud.ibm.com/v1beta1`. To ensure that the policy is enforced, you must update the API version to `apiVersion: portieris.cloud.ibm.com/v1`.

### Image mutation option

You can also set a mutate image, `mutateImage: bool`, behavior preference for each policy. The default value is `true`, which is also the original behavior and means that on successful admission the container's image property is mutated to ensure that the immutable digest form of the image is used. If the value is `false`, the original image reference is retained with the consequences described in [README](README.md#image-mutation-option).

**Example**

```yaml
apiVersion: portieris.cloud.ibm.com/v1
kind: ImagePolicy
metadata:
  name: signedby-me
spec:
   repositories:
    - name: "icr.io/*"
      policy:
        mutateImage: false
        simple:
          requirements:
          - type: "signedBy"
            keySecret: my-pubkey
```

### `trust` (Docker Content Trust and Notary)

Portieris supports sourcing trust data from the following registries without additional configuration in the image policy:

* IBM Cloud Container Registry
* Quay.io
* Docker Hub

To use a different trust server for a repository, you can specify the `trustServer` parameter in your policy:

**Example**

```yaml
apiVersion: portieris.cloud.ibm.com/v1
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

### `simple` (Red Hat simple signing)

The policy requirements are similar to those defined for the configuration files that are consulted when you're using the Red Hat tools [policy requirements](https://github.com/containers/image/blob/master/docs/containers-policy.json.5.md#policy-requirements). However, the main difference is that the public key in a `signedBy` requirement is defined in a `keySecret` attribute, the value is the name of an in-scope Kubernetes secret that contains a public key block. The value of `keyType`, `keyPath`, and `keyData`, see [policy requirements](https://github.com/containers/image/blob/master/docs/containers-policy.json.5.md#policy-requirements), can't be provided. If multiple keys are present in the key ring, the requirement is satisfied if the signature is signed by any of them.

To export a single public key identified by `<finger-print>` from Gnu Privacy Guard (GPG) and create a KeySecret from it, you can use the following script.

```bash
gpg --export --armour <finger-print> > my.pubkey
kubectl create secret generic my-pubkey --from-file=key=my.pubkey
```

When you create the secret, ensure that you're creating the key with a value of `key`.

```
kubectl create secret generic my-pubkey --from-file=key=<your_pub_key>
```

The following example requires that images from `icr.io` are signed by the identity with public key in `my-pubkey`.

```yaml
apiVersion: portieris.cloud.ibm.com/v1
kind: ImagePolicy
metadata:
  name: signedby-me
spec:
   repositories:
    - name: "icr.io/*"
      policy:
        simple:
          requirements:
          - type: "signedBy"
            keySecret: my-pubkey
```

The following example requires that a specific image is signed, but it allows the registry location to change. In this pattern, a policy for each image is required to exactly define the new location.

```yaml
apiVersion: portieris.cloud.ibm.com/v1
kind: ImagePolicy
metadata:
  name: db2mirror
spec:
   repositories:
    - name: "uk.icr.io/mymirror/db2/db2manager:6.1.0.0"
      policy:
        simple:
          requirements:
          - type: "signedBy"
            keySecret: db2-pubkey
            signedIdentity:
                type: "matchExactRepository"
                dockerRepository: "icr.io/ibm/db2/db2manager"
```

The following example allows many images with a common origin that have moved with their original signature, for example, by mirroring, in this case from `icr.io/db2` to `registry.myco.com:5000/mymirror/ibmdb2`.

```yaml
apiVersion: portieris.cloud.ibm.com/v1
kind: ImagePolicy
metadata:
  name: remapDB2Example
spec:
   repositories:
    - name: "registry.myco.com:5000/mymirror/ibmdb2"
      policy:
        simple:
          requirements:
          - type: "signedBy"
            keySecret: db2-pubkey
            signedIdentity:
                type: "remapIdentity"
                prefix: "registry.myco.com:5000/mymirror/ibmdb2"
                signedPrefix: "icr.io/db2"
```

You can also specify the location of signature storage for registries that don't support the registry extension. Where `storeSecret` identifies an in-scope Kubernetes secret that contains `username` and `password` data items that are used to authenticate with the server referenced in `storeURL`.

```yaml
apiVersion: portieris.cloud.ibm.com/v1
kind: ImagePolicy
metadata:
  name: db2mirror
spec:
   repositories:
    - name: "uk.icr.io/mymirror/db2/db2manager:6.1.0.0"
      policy:
        simple:
          storeURL: https://server.another.com/signatures
          storeSecret: another-secret
          requirements:
          - type: "signedBy"
            keySecret: db2-pubkey
```

### `vulnerability`

Vulnerability policies enable you to admit or deny pod admission based on the security status of the container images within the pod. Vulnerability-based admission is available for [Vulnerability Advisor for IBM Cloud Container Registry](https://cloud.ibm.com/docs/Registry?topic=va-va_index). Vulnerability Advisor is available for any image in [IBM Cloud Container Registry](https://www.ibm.com/cloud/container-registry).

Example policy:

```yaml
apiVersion: portieris.cloud.ibm.com/v1
kind: ImagePolicy
metadata:
  name: block-vulnerable-images
spec:
   repositories:
    - name: "uk.icr.io/*"
      policy:
        vulnerability:
          ICCRVA:
            enabled: true
            account: "an-IBM-Cloud-account-id"
```

#### Vulnerability Advisor for IBM Cloud Container Registry details

For each `container` in the pod that is being considered for admission, a [vulnerability status](https://cloud.ibm.com/apidocs/container-registry/va#imagestatusquerypath) report is retrieved for the `image` that is specified by the container.

The optional `account` parameter specifies the IBM Cloud account from where exemptions are retreived for images that match the policy repository name.

If the report returns an overall status of `OK`, `WARN`, or `UNSUPPORTED` the pod is allowed. In the event of any other status, or any error condition, the pod is denied.

**Note** Recently pushed images to the registry that have not completed scanning are denied admission.
