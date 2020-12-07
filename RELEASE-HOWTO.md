# RELEASE HOWTO

A runbook style document for project owners to have some consistency in the process used to make a new release and set expectations for contributors.

## Version Number
Each release will see an increment in the version number acording to [SEMVER](https://semver.org/)
Use vx.x.x going forward. 

## Release Process
1. Update VERSION in Makefile, Chart.yaml and values.yaml to the intended release version. Update the CHANGELOG.md to reference the right version and date. Update go.mod and go.sum, commit all this.
1. Run `make alltests`
1. Run `make e2e` or `make helm.install.local` & `make e2e.quick`  
1. Publish the image to dockerhub.
1. Create a tag = VERSION `git tag <VERSION>`.
1. Create a release with the chart as a release artifact. 

