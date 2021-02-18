---

copyright:
  years: 2020, 2021
lastupdated: "2021-02-18"

---

# How to create a new release

Find out about the process that project owners must follow to set expectations and ensure consistency when you are creating a new release.

## Version number

Each release version number must increase incrementally and comply with [Semantic Versioning (SEMVER)](https://semver.org/). 

Use the following format `vx.x.x`. 

## Release process

1. Update **VERSION** in the `Makefile`, `Chart.yaml`, and `values.yaml` files to the release version number. 
2. Update the `CHANGELOG.md` file to reference the right version and date. 
3. Update the `go.mod` and `go.sum` files. 
4. Commit the changes.
5. Run `make alltests`.
6. Run `make e2e`, or both `make helm.install.local` and `make e2e.quick`.
7. Publish the image to Docker Hub.
8. Create a **tag = VERSION** by running `git tag <VERSION>`.
9. Create a release that has the chart as a release artifact. 
