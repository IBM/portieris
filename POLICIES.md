# Portieris Policies

## Custom Resource Definitions

Portieris defines 2 resource types:

* ClusterImagePolicy
    - cluster scope
* ImagePolicy
    - namespace scope

ClusteImagePolicy example:
```
```

ImagePolicy example:
```
```

In both cases multiple instances are merged, and can be protected by RBAC policy.

## Repository Matching

When an image is evaluated for admisson, a policy set is matched using a wildcard enabled matched on the repository name, if no match is found in an ImagePolicy the ClusterImagePolicy is consulted. If there are multiple matches the
most specific match is used. 

## Policy 

A policy consists of multiple sections defining requirements on the image either in `trust:` (Docker Content Trust / Notary V1)  or `simple:` (RedHat's Simple Signing) sections. 

### Simple

The policy requirements are similar to those which can be defined in host configuration files consulted when using the RedHat tools []() however with some differences. The main difference is that the public key in a "signedBy" requirement is defined in a "KeySecret:" attribute, the value is the name of an in-scope Kubernetes secret containing the public key data. 

examples:



