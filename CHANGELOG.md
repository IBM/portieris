---

copyright:
  years: 2020-2024
lastupdated: "2024-05-14"

---

# Change Log

Notable changes recorded here.
This project adheres to [Semantic Versioning](http://semver.org/).

## v-next


## v0.13.16

Released: 2024-06-06

* Remediates CVE-2023-45288, CVE-2023-45289, CVE-2023-45290, CVE-2024-24783 and CVE-2024-24784 in go-toolset
* Remediates CVE-2024-2961 in glibc
* Remediates CVE-2024-3727 in github.com/containers/image

## v0.13.15

Released: 2024-05-14

* Remediates CVE-2024-2961 in glibc

## v0.13.14

* github.com/docker/docker update for CVE-2024-29018
* golang.org/x/net update for CVE-2023-45288

## v0.13.13

* start with go=toolset:1.20.12 also for the installer (consistency)
Note: the build pulls dynamic updates to the builder image currently gets go-toolset:1.20.12-3 which resolves CVE-2024-1394

## v0.13.12

Released 2024-03-06

* Update go-toolset:1.20.12-2
* golang/github.com/opencontainers/runc update for CVE-2024-21626
* Add arm64 image. This makes developing and testing on a M1/2 mac easier

## v0.13.11

Released 2024-01-09

* Update go-toolset:1.20.10-3
* Rebuild/Package updates to remediate CVE-2023-3446 CVE-2023-3817 CVE-2023-5678
* golang.org/x/crypto update for CVE-2023-48795

## v0.13.10

Released 2023-11-07

* Set nonroot user on image iconfig to supress container policy checkers.
* Have nancy run from Dockerfile again.
* Remediate CVE.

## v0.13.9

Released 2023-11-01

* Remediate CVE-2023-44487 CVE-2023-29406 CVE-2023-39325 with go-toolset:1.19.13-2.1698062273
* Resolve a compatibility with GKE versioning in templates/pdb.yaml
* Allow namespace selector for skipping admission webhook

## v0.13.8

Released 2023-10-10

* Remediates CVE-2023-4527 CVE-2023-4806 CVE-2023-4813 CVE-2023-4911 in glibc
## v0.13.7

Released 2023-09-11

* Remediates CVE-2023-3978

## v0.13.6

Released 2023-08-21

* consume ubi8/go-toolset:1.19.10-10
* Remediates CVE-2022-41724 CVE-2022-41725 CVE-2023-24540 CVE-2023-29402 CVE-2023-29403 CVE-2023-29404 CVE-2023-29405
* refactor tests since IBM has removed notary service
* do not test vulnerability policy since IBM has deprecated the API

## v0.13.5

Released 2023-04-11

* Remove vulnerable dependency dgrijalva/jwt-go

## v0.13.4

Released 2023-03-29

* Update to go-toolset:1.18.9-13
* Resolves CVE-2022-4304 CVE-2022-4450 CVE-2023-0215 CVE-2023-0286 with openssl
* Resolves CVE-2023-27561 with runc v1.1.15

## v0.13.3

Released 2023-02-02

* Contributed helm value options: skipCreate certificate issuer (aid seamless upgrade) and optional annotations.
* Update to go-toolset:1.18.9-8
* Fixes problem with portieris version in logs showing the golang version

## v0.13.2

Released 2023-01-25

* Update to go-toolset:1.18.4-20 and ensures go rpm is tracked in final image
* Update go dependencies

## v0.13.1

Released 2022-08-23

* Upgrade runc to v1.1.2 for vulnerability fix
* Build with go-toolset:1.17.12 for vulnerability fix

## v0.13.0

Released 2022-06-28

* resolved ([#51](https://github.com/IBM/portieris/issues/51)) following the oAuth spec
* Build with go-tooolset:1.17.10 resolving CVE-2022-29526 CVE-2022-23772 CVE-2022-24921
* code-generator:1.24 + regenerate code
* Helm chart improvements: Fixes ([#142](https://github.com/IBM/portieris/issues/142))
* options to define podDisruptionBudget and options to use generated certificates directly from values.yml ([PR#379](https://github.com/IBM/portieris/pull/379))
* resolve ([#388](https://github.com/IBM/portieris/issues/388)), remove cluster-admins group from SCC

## v0.12.6

Released 2022-08-23

* Build with go-toolset:1.17.12 for vulnerability fix

## v0.12.5

Released 2022-06-30

* Build with go-toolset:1.17.10 resolving CVE-2022-29526 CVE-2022-23772 CVE-2022-24921

## v0.12.4

Released 2022-04-04

* Rebuild
* Resolves CVE-2022-0778

## v0.12.3

Released 2022-04-04

* Rebuild
* Resolves CVE-2021-3999
* Resolves CVE-2022-23218
* Resolves CVE-2022-23219

## v0.12.2

Released 2022-01-06

* Resolves CVE-2021-3712
* Build with go-toolset:1.16.12

## v0.12.1

Released 2021-11-30

* Resolves CVE-2021-23840
* Resolves CVE-2021-23841
* Resolves CVE-2021-27645
* Resolves CVE-2021-33574
* Resolves CVE-2021-35942
* Supports cert-manager >= 1.6

## v0.12.0

Released 2021-10-11

* Added support for batch/v1/cronjobs and dropped batch/v1alpha1/cronjobs inline with 1.21 apis ([#350](https://github.com/IBM/portieris/issues/350))
* Many more documentation improvements
* Set sane priorityClass ([#352](https://github.com/IBM/portieris/issues/352))
* Build using ubi go toolset (golang 1.15.14), and run in ubi-minimal ([#351](https://github.com/IBM/portieris/issues/351))
* Support ObjectSelectorAdmissionSkip ([#349](https://github.com/IBM/portieris/issues/349))
* Require TLS1.2 on webhook

## v0.11.0

Released 2021-06-16

* Further documentation improvements including godoc
* use current resource versions ([#215](https://github.com/IBM/portieris/issues/215))
* update dependencies and golang version to 1.16.5
* a template fix to correctly identify openshift ([PR#326](https://github.com/IBM/portieris/pull/326))

## v0.10.3

Released 2021-06-22

* update dependencies and golang to 1.16.5

## v0.10.2

Released 2021-03-25

* Documentation improvements
* Add keySecretNamespace policy option ([PR#258](https://github.com/IBM/portieris/pull/258))
* update dependencies and golang to 1.15.10

## v0.10.1

Released 2021-02-10

* Add mutateImage policy option ([244](https://github.com/IBM/portieris/issues/244))
* When skipping checks because a parent resource exists, ensure its a known type ([246](https://github.com/IBM/portieris/issues/246))
* Add version to user-agent consistently ([241](https://github.com/IBM/portieris/issues/241))

## v0.10.0

Released 2021-01-11

* Support verifying images that don't require pull secrets ([#123](https://github.com/IBM/portieris/issues/123))
* Redefine policy Custom Resource Definitions (CRDs) by using v1 with validation, breaking change ([#121](https://github.com/IBM/portieris/issues/121))

## v0.9.5

Released 2020-12-15

* Support remapIdentity simple signature identity type ([#92](https://github.com/IBM/portieris/issues/92))
* Switch to pull image from `icr.io/portieris` ([#205](https://github.com/IBM/portieris/issues/205))
* Get default ClusterImagePolicy setting from values ([PR#233](https://github.com/IBM/portieris/pull/233))

## v0.9.4

Released 2020-12-06

* Update to Go 1.14.12
* Support OpenShift projects that create deployments with blank image names ([#227](https://github.com/IBM/portieris/issues/227))
* Add `webHooks.failurePolicy` `value/option`

## v0.9.2

Released 2020-11-30

* Additional logging container image name ([#216](https://github.com/IBM/portieris/issues/216))
* Accept a `--kubeconfig` command line parameter ([PR#218](https://github.com/IBM/portieris/pull/218))

## v0.9.1

Released 2020-11-23

* Add metrics counting allow and deny events. ([#106](https://github.com/IBM/portieris/issues/162))
* Fix a problem with multiple pull secrets and simple signing ([#209](https://github.com/IBM/portieris/issues/209))

## v0.9.0

Released 2020-11-05

* Introduce a policy type to enforce an image vulnerability check ([#71](https://github.com/IBM/portieris/issues/71))
* Normalise the use of Helm, allow `--create-namespace`, remove webhook on uninstall ([PR#189](https://github.com/IBM/portieris/pull/189))
* Add a default policy for Istio image when it's running on IBM Cloud Kubernetes Service ([PR#198](https://github.com/IBM/portieris/pull/198))
* Fix certificate incompatibility in Kubernetes 1.19 ([#196](https://github.com/IBM/portieris/issues/196))

## v0.8.2

Released 2020-10-12

* Provide an option to run out of cluster ([#180](https://github.com/IBM/portieris/issues/180))

## v0.8.1

Released 2020-09-18

* PR checker fixed to fail when tests fail ([#167](https://github.com/IBM/portieris/issues/167))
* Drop support for Helm 2. You must use Helm 3 to install Portieris ([#141](https://github.com/IBM/portieris/issues/141)) ([#41](https://github.com/IBM/portieris/issues/41)) ([#89](https://github.com/IBM/portieris/issues/89))
* Ability to use a namespace selector for admission webhook ([#112](https://github.com/IBM/portieris/issues/112))
* Correctly decode pull secrets where credentials are in the `auth` field ([#174](https://github.com/IBM/portieris/issues/174))
* Ensure the pre-installation steps create the namespace before the service account ([#181](https://github.com/IBM/portieris/issues/181))

## 0.8.0

Released 2020-09-02

* Fix the port name in service template ([PR#149](https://github.com/IBM/portieris/pull/149))
* Change the default namespace to `portieris` ([#117](https://github.com/IBM/portieris/issues/117))
* Support Helm 3 and Openshift 4 ([PR#130](https://github.com/IBM/portieris/pull/130))
* Anti-affinity and liveness/readiness probes  ([#66](https://github.com/IBM/portieris/issues/66))
* Support sourcing webhook certificates from cert-manager ([#59](https://github.com/IBM/portieris/issues/59))
* Allow anonymous Notary access ([PR#159](https://github.com/IBM/portieris/pull/159))

## 0.7.0

Released 2020-06-09

* Support for reading simple signatures from lookaside storage, ([#93](https://github.com/IBM/portieris/issues/93))

## 0.6.0

Released 2020-03-26

* Support for the verification of simple signatures by using [containers/image](https://github.com/containers/image). ([#70](https://github.com/IBM/portieris/issues/70))
