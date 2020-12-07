# RELEASE HOWTO

A runbook style document for project owners to have some consistency in the process used to make a new release and set expectations for contributors.

## Version Number
Each release will see an increment in the version number acording to [SEMVER](https://semver.org/)
Releases will be even minor numbers, odd minor numbers indicate development version (building toward next even number release)

## Release Process
1. Update VERSION in Makefile, Chart.yaml and values.yaml to the intended (even minor version) release version. Check the CHANGELOG.md references the right version and date. 
1. Run make alltests
1. Run make e2e or helm.install & e2e.quick 
1. Publish the image to dockerhub (currently).
1. Commit go.mod and go.sum 
1. Create a tag = VERSION
1. Create a release with the chart as a release artifact. 

## Post Release 
1. As the first commit after release update VERSION in Makefile, version in Chart.yaml, heading in CHANGELOG.md, and tag in values.yaml to the next odd minor version number to indicate that further commits are clearly contributions to the next release. Consequently installers from source helm should check out the release tag commit to find the appropriate image by default.
