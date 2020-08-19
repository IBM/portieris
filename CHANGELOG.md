# Change Log

Notable changes recorded here.
This project adheres to [Semantic Versioning](http://semver.org/).

# 0.7.1
## 2020-06-xx

* Fix the port name in service template ([PR#149](https://github.com/IBM/portieris/pull/149))
* Change the default namespace to portieris ([#117](https://github.com/IBM/portieris/issues/117))
* Support Helm3 and Openshift 4 ([#130](https://github.com/IBM/portieris/pull/130))
* anti-affinity and liveness/readiness probes  ([#66](https://github.com/IBM/portieris/issues/66))
* Support sourcing webhook certificates from cert-manager ([#59](https://github.com/IBM/portieris/issues/59))

# 0.7.0
## 2020-06-09

* Support for reading simple signatures from lookaside storage, ([#93](https://github.com/IBM/portieris/issues/93))

# 0.6.0
## 2020-03-26

* Support for the verification of simple signatures using [containers/image](https://github.com/containers/image). ([#70](https://github.com/IBM/portieris/issues/70))
