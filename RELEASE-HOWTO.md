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

## Determining version number

Versioning should be performed according to the semver spec. Accordingly, the patch version must be used only for non-breaking bug fixes, and new function should be a minor version increase.

Since changes to the Helm chart will change behaviour, any change to the Helm chart must be released in a minor version bump. Patch versions of the container image should be consumable without reviewing the Helm chart for changes.

## Release process

1. Update **VERSION** in the `Makefile`, `Chart.yaml`, and `values.yaml` files to the release version number. 
2. Update the `CHANGELOG.md` file to reference the right version and date. 
3. Update the `go.mod` and `go.sum` files. 
4. Commit the changes.
5. Run `make alltests`.
6. Run `make e2e`, or both `make helm.install.local` and `make e2e.quick`.
7. Publish the image to IBM Cloud Container Registry at `icr.io/portieris`
8. Create a **tag = VERSION** by running `git tag <VERSION>`.
9. Create a release that has the chart as a release artifact. 
