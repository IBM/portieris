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

For both types or resource if there are multiple resources, they are merged together, and can be protected by RBAC policy.

## Installation Default Policies
Default policies are installed when Portieris is installed. You must review and change these according to your requirements.
The installation [default policies](helm/portieris/templates/default/policies.yaml) should be customised. 

## Repository Matching
When an image is evaluated for admission, the set of policies set is wildcard matched on the repository name. If there are multiple matches the most specific match is used. 

## Policy 
A policy consists of an array of objects defining requirements on the image either in `trust:` (Docker Content Trust / Notary V1)  or `simple:` (RedHat's Simple Signing) objects . 

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
The policy requirements are similar to those defined for configuration files consulted when using the RedHat tools [policy requirements](https://github.com/containers/image/blob/master/docs/containers-policy.json.5.md#policy-requirements). However there are some differences, the main difference is that the public key in a "signedBy" requirement is defined in a "KeySecret:" attribute, the value is the name of an in-scope Kubernetes secret containing the public key data. The value of "KeyType" is implied and cannot be provided.

this example requires that images from `icr.io` are signed by the identity with public key in `my-pubkey`:
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
          - type: "signedBy"
            keySecret: my-pubkey
```

this example requires that a given image is singed but allows the location to have changed:
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
          - type: "signedBy"
            keySecret: my-pubkey
            signedIdentity: 
                type: "matchExactRepository"
                dockerRepository: "icr.io/ibm/db2/db2manager"
```




