# Change Log

Notable changes recorded here.
This project adheres to [Semantic Versioning](http://semver.org/).

# v0.8.2
## 2020-10-12

* Provide option to run out of cluster ([#180](https://github.com/IBM/portieris/issues/180))

# v0.8.1
## 2020-09-18

* PR checker fixed to fail when tests fail ([#167](https://github.com/IBM/portieris/issues/167))
* Drop support for Helm 2. You must now use Helm 3 to install Portieris ([#141](https://github.com/IBM/portieris/issues/141)) ([#41](https://github.com/IBM/portieris/issues/41)) ([#89](https://github.com/IBM/portieris/issues/89))
* Ability to use a namespace selector for admission webhook ([#112](https://github.com/IBM/portieris/issues/112))
* Correctly decode pull secrets where credentials are in the auth field ([#174](https://github.com/IBM/portieris/issues/174))
* Ensure the pre-install steps create the namespace before the serviceaccount ([#181](https://github.com/IBM/portieris/issues/181))

# 0.8.0
## 2020-09-02

* Fix the port name in service template ([PR#149](https://github.com/IBM/portieris/pull/149))
* Change the default namespace to portieris ([#117](https://github.com/IBM/portieris/issues/117))
* Support Helm3 and Openshift 4 ([PR#130](https://github.com/IBM/portieris/pull/130))
* anti-affinity and liveness/readiness probes  ([#66](https://github.com/IBM/portieris/issues/66))
* Support sourcing webhook certificates from cert-manager ([#59](https://github.com/IBM/portieris/issues/59))
* Allow anonymous notary access ([PR#159](https://github.com/IBM/portieris/pull/159))

# 0.7.0
## 2020-06-09

* Support for reading simple signatures from lookaside storage, ([#93](https://github.com/IBM/portieris/issues/93))

# 0.6.0
## 2020-03-26

* Support for the verification of simple signatures using [containers/image](https://github.com/containers/image). ([#70](https://github.com/IBM/portieris/issues/70))
