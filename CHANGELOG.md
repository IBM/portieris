---

copyright:
  years: 2020, 2021
lastupdated: "2021-03-25"

---

# Change Log

Notable changes recorded here.
This project adheres to [Semantic Versioning](http://semver.org/).

# v-next
##

# v0.11.1
## 2021-06-14
* Further documentation improvements including godoc
* Update webhook resources to v1 
* update dependencies and golang version to 1.16.5
* a template fix to correctly identify openshift ([PR#326](https://github.com/IBM/portieris/pull/326))

# v0.10.2
## 2021-03-25
* Documentation improvements 
* Add keySecretNamespace policy option ([PR#258](https://github.com/IBM/portieris/pull/258))
* update dependencies and golang to 1.15.10

# v0.10.1
##
* Add mutateImage policy option ([244](https://github.com/IBM/portieris/issues/244))
* When skipping checks because a parent resource exists, ensure its a known type ([246](https://github.com/IBM/portieris/issues/246))
* Add version to user-agent consistently ([241](https://github.com/IBM/portieris/issues/241))

# v0.10.0
## 2021-01-11
* Support verifying images that don't require pull secrets ([#123](https://github.com/IBM/portieris/issues/123))
* Redefine policy Custom Resource Definitions (CRDs) by using v1 with validation, breaking change ([#121](https://github.com/IBM/portieris/issues/121))
 

# v0.9.5
## 2020-12-15
* Support remapIdentity simple signature identity type ([#92](https://github.com/IBM/portieris/issues/92)) 
* Switch to pull image from `icr.io/portieris` ([#205](https://github.com/IBM/portieris/issues/205))
* Get default ClusterImagePolicy setting from values ([PR#233](https://github.com/IBM/portieris/pull/233))

# v0.9.4
## 2020-12-06
* Update to Go 1.14.12
* Support OpenShift projects that create deployments with blank image names ([#227](https://github.com/IBM/portieris/issues/227))
* Add `webHooks.failurePolicy` `value/option`

# v0.9.2
## 2020-11-30
* Additional logging container image name ([#216](https://github.com/IBM/portieris/issues/216))
* Accept a `--kubeconfig` command line parameter ([PR#218](https://github.com/IBM/portieris/pull/218))

# v0.9.1
## 2020-11-23
* Add metrics counting allow and deny events. ([#106](https://github.com/IBM/portieris/issues/162))
* Fix a problem with multiple pull secrets and simple signing ([#209](https://github.com/IBM/portieris/issues/209))

# v0.9.0
## 2020-11-05

* Introduce a policy type to enforce an image vulnerability check ([#71](https://github.com/IBM/portieris/issues/71))
* Normalise the use of Helm, allow `--create-namespace`, remove webhook on uninstall ([PR#189](https://github.com/IBM/portieris/pull/189))
* Add a default policy for Istio image when it's running on IBM Cloud Kubernetes Service ([PR#198](https://github.com/IBM/portieris/pull/198))
* Fix certificate incompatibility in Kubernetes 1.19 ([#196](https://github.com/IBM/portieris/issues/196))

# v0.8.2
## 2020-10-12
* Provide an option to run out of cluster ([#180](https://github.com/IBM/portieris/issues/180))

# v0.8.1
## 2020-09-18
* PR checker fixed to fail when tests fail ([#167](https://github.com/IBM/portieris/issues/167))
* Drop support for Helm 2. You must use Helm 3 to install Portieris ([#141](https://github.com/IBM/portieris/issues/141)) ([#41](https://github.com/IBM/portieris/issues/41)) ([#89](https://github.com/IBM/portieris/issues/89))
* Ability to use a namespace selector for admission webhook ([#112](https://github.com/IBM/portieris/issues/112))
* Correctly decode pull secrets where credentials are in the `auth` field ([#174](https://github.com/IBM/portieris/issues/174))
* Ensure the pre-installation steps create the namespace before the service account ([#181](https://github.com/IBM/portieris/issues/181))

# 0.8.0
## 2020-09-02
* Fix the port name in service template ([PR#149](https://github.com/IBM/portieris/pull/149))
* Change the default namespace to `portieris` ([#117](https://github.com/IBM/portieris/issues/117))
* Support Helm 3 and Openshift 4 ([PR#130](https://github.com/IBM/portieris/pull/130))
* Anti-affinity and liveness/readiness probes  ([#66](https://github.com/IBM/portieris/issues/66))
* Support sourcing webhook certificates from cert-manager ([#59](https://github.com/IBM/portieris/issues/59))
* Allow anonymous Notary access ([PR#159](https://github.com/IBM/portieris/pull/159))

# 0.7.0
## 2020-06-09
* Support for reading simple signatures from lookaside storage, ([#93](https://github.com/IBM/portieris/issues/93))

# 0.6.0
## 2020-03-26
* Support for the verification of simple signatures by using [containers/image](https://github.com/containers/image). ([#70](https://github.com/IBM/portieris/issues/70))
