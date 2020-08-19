# RELEASE HOWTO

A runbook style document for project owners to have some consistency in the process used to make a new release and set expectations for contributors.

## Version Number
Each release will see an increment in the version number acording to [SEMVER](https://semver.org/)

## Release Process
1. Update VERSION in Makefile, Chart.yaml and values.yaml to the intended release version. Check the CHANGELOG.md references the right version. 
1. Run make alltest
1. Run make e2e
1. Publish the image to dockerhub (currently).
1. Create a tag = VERSION
1. Create a release with the chart as a release artifact. 

## Post Release 
1. As the first commit after release update VERSION in Makefile, Chart.yaml and values.yaml with the suffix '+' such that so that further commits are clearly contributions to the next release.  
