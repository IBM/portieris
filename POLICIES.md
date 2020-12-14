# Portieris Policies

## Image Policy Resources

Portieris defines two custom resource types for policy:

* ImagePolicies can be configured in a Kubernetes namespace, and define Portieris' behavior in that namespace. If ImagePolicies exists in a namespace, the policies from those ImagePolicy resources are used exclusively, if there is no match for the workload image in ImagePolicies ClusterImagePolicies are not examined. Images in deployed workloads are wildcard matched against the set of policies defined, if there is no policy matching the workload image then deployment is denied.
  - this example allows any image from the "icr.io" registry with no further checks (the policy is empty):
```yaml
apiVersion: securityenforcement.admission.cloud.ibm.com/v1beta1
kind: ImagePolicy
metadata:
  name: allow-all-icrio
spec:
   repositories:
    - name: "icr.io/*"
      policy:
```

* ClusterImagePolicies are configured at the cluster level, and take effect whenever there is no ImagePolicy resource defined in the namespace where the workload is being deployed. These resources have the same structure as namespace ImagePolicies and if no matching policy is found for an image deployment is denied.
  - this example allows all images from all registries with no checks:

```yaml
apiVersion: securityenforcement.admission.cloud.ibm.com/v1beta1
kind: ClusterImagePolicy
metadata:
  name: portieris-default-cluster-image-policy
spec:
  repositories:
  - name: '*'
    policy:
```

For both types of resource if there are multiple resources, they are merged together, and can be protected by RBAC policy.

## Installation Default Policies

Default policies are installed when Portieris is installed. You must review and change these according to your requirements.
The installation [default policies](helm/portieris/templates/default/policies.yaml) should be customised.

## Repository Matching

When an image is evaluated for admission, the set of policies set is wildcard matched on the repository name. If there are multiple matches the most specific match is used.

## Policy
A policy consists of an array of objects defining requirements on the image using either `trust:` (Docker Content Trust / Notary V1), `simple:` (RedHat's Simple Signing) or `vulnerability:` objects.

### trust (Docker Content Trust/Notary)
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

### simple (RedHat simple signing)
The policy requirements are similar to those defined for configuration files consulted when using the RedHat tools [policy requirements](https://github.com/containers/image/blob/master/docs/containers-policy.json.5.md#policy-requirements). However there are some differences, the main difference is that the public key in a`signedBy` requirement is defined in a `keySecret` attribute, the value is the name of an in-scope Kubernetes secret containing a public key block. The value of `keyType`, `keyPath` and `keyData` (seen in [policy requirements](https://github.com/containers/image/blob/master/docs/containers-policy.json.5.md#policy-requirements)) cannot be provided. If multiple keys are present in the keyring then the requirement is satisfied if the signature is signed by any of them.

To export a single public key identified by `<finger-print>` from gpg and create a KeySecret from it you could use:
```bash
gpg --export --armour <finger-print> > my.pubkey
kubectl create secret generic my-pubkey --from-file=key=my.pubkey
```

In creating the secret, ensure you are creating the key with a value of `key`, as shown below:

`kubectl create secret generic my-pubkey --from-file=key=<your_pub_key>`

This example requires that images from `icr.io` are signed by the identity with public key in `my-pubkey`:
```yaml
apiVersion: securityenforcement.admission.cloud.ibm.com/v1beta1
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

This example requires that a given image is signed but it allows the registry location to have changed, in this pattern a policy per image is required to exactly define the new location (but see the next example):
```yaml
apiVersion: securityenforcement.admission.cloud.ibm.com/v1beta1
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

This example allows many images with a common origin which have been moved with their original signature, for example by mirroring, in this case from `icr.io/db2` to `registry.myco.com:5000/mymirror/ibmdb2`:
```yaml
apiVersion: securityenforcement.admission.cloud.ibm.com/v1beta1
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


It is also possible to specify the location of signature storage for registries which do not support the registry extension:
```yaml
apiVersion: securityenforcement.admission.cloud.ibm.com/v1beta1
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
where `storeSecret` identifies an in scope Kubernetes secret which contains `username` and `password` data items which are used to authenticate with the server referenced in `storeURL`.

### vulnerability

Vulnerability policies enable you to admit or deny pod admission based on the security status of the container images within the pod. Vulnerability-based admission is available for:
* [Vulnerability Advisor for IBM Cloud Container Registry](https://cloud.ibm.com/docs/Registry?topic=va-va_index)
    * This is available for any image in the [IBM Cloud Container Registry](https://www.ibm.com/uk-en/cloud/container-registry)

Example policy:
```yaml
apiVersion: securityenforcement.admission.cloud.ibm.com/v1beta1
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
For each `container` in the pod being considered for admission, a [vulnerability status](https://cloud.ibm.com/apidocs/container-registry/va#imagestatusquerypath) report is retrieved for the `image` specified by the container.

The optional `account` parameter specifies the IBM Cloud account where exemptions should be fetched from for image matching the policy repository name.

If the report returns an overall status of `OK`, `WARN` or `UNSUPPORTED` the pod is allowed. In the event of any other status, or any error condition the pod is denied.

Please note that images that were recently pushed to the registry and have not yet completed scanning will be denied admission.
