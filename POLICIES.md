# Portieris Policies

## Custom Resource Definitions

Portieris defines two custom resource types:

* ClusterImagePolicy at cluster scope
  - example:
```
apiVersion: securityenforcement.admission.cloud.ibm.com/v1beta1
kind: ClusterImagePolicy
metadata:
  name: portieris-default-cluster-image-policy
spec:
  repositories:
  - name: '*'
    policy:
```
This resource provides a default when no ImagePolicy is defined. This example is an empty policy which allows all images.

* ImagePolicy at namespace scope
  - example:
```
apiVersion: securityenforcement.admission.cloud.ibm.com/v1beta1
kind: ImagePolicy
metadata:
  name: allow-all
spec:
   repositories:
    - name: "icr.io/*"
      policy:
         simple:
          - type: "insecureAcceptAnything"
```
This resource enables different namespaces to be treated appropriately, the example shows all images from the "icr.io" registry being allowed, simple signature verification will be invoked but the policy is totally permissive. The overall effect is similar to an empty policy.

In both cases multiple policy resources are merged, and can be protected by RBAC policy.

## Repository Matching

When an image is evaluated for admission, the set of policies set is wildcard matched on the repository name. Only if no match is found from ImagePolicies are the ClusterImagePolicies consulted. If there are multiple matches the most specific match is used. 

## Policy 

A policy consists of an array of objects defining requirements on the image either in `trust:` (Docker Content Trust / Notary V1)  or `simple:` (RedHat's Simple Signing) objects . 

### Simple

The policy requirements are similar to those defined for configuration files consulted when using the RedHat tools [policy requirements](https://github.com/containers/image/blob/master/docs/containers-policy.json.5.md#policy-requirements) however there are some differences. The main difference is that the public key in a "signedBy" requirement is defined in a "KeySecret:" attribute, the value is the name of an in-scope Kubernetes secret containing the public key data. The value of "KeyType" is implied and cannot be provided.

examples:
```
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


```
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




