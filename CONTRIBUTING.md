# Contributing to Portieris

## Reporting security issues

The Portieris maintainers take security extremely seriously. To report a security issue, DO NOT open an issue. Instead, send your report via email to alchreg@uk.ibm.com privately.

## Reporting a general issue

Step one in reporting a problem with Portieris should always be to check the [issues page](https://github.com/ibm/portieris/issues) to see if your issue has been reported already. If an issue exists, don't open another one. Instead, react to the initial post inside that issue with a thumbs-up (+1) emoji. If you have additional debugging information, feel free to add a comment with that information.

If you don't find an existing issue, feel free to create one. Make sure to:

1. Use a title that accurately describes your problem, and allows other users to find it easily.
2. If you're reporting an error, include any relevant output of your Portieris containers, and the message that you receive on the command line.
3. If possible, create a workload that allows other people to quickly reproduce the issue, and include it in your report.
4. Provide as much relevant information about your cluster's configuration. For example, if your cluster is behind a firewall and you're having network connectivity problems, check whether your firewall is blocking the traffic and include any relevant logs.

## Contributing small changes

If you've found a bug or a typo and want to fix it yourself, make sure to report it first, but clearly indicate that you intend to fix it yourself. As above, if the bug is a security issue, report it privately rather than opening an issue. As you write the fix, add a test to reproduce the problem, and make sure that it passes before you submit your pull request.

## Contributing major changes

Before you start work on a new feature or an architectural change, open an issue in which you describe what you are planning to do, and tag it `enhancement`, to open a dialogue with the maintainers. Once you and the rest of the community have agreed on what should be done, start work and submit a PR as normal.

Try to keep your PRs as small as possible. It's a lot easier to review many smaller PRs than one enormous one. If you need to make a large change, consider whether you can break the change down into a number of smaller changes.

## Changes to policy files

If you add new elements to the Portieris policy files (pkg/apis/securityenforcement/v1beta1/*) you will need to run the code generator. Please see: [Generator Readme](pkg/apis/securityenforcement/v1beta1/README.md)

## Coding style

Generally, code should be written using idiomatic Go. Check out [Effective Go](https://golang.org/doc/effective_go.html) for loads of well-written tips on how to write Go... effectively!

## Testing changes

Your code must pass Travis before it can be merged. Before you submit a PR, you should make sure that the tests pass. You can run the tests on your workstation using `make`. To download the prerequisites for running the tests, run `make test-deps` from the project root.

You can run the following commands:

* `make vet` runs `go vet` on the project. Vet looks for issues in the code, such as useless variable assignments. Vet failures should be fixed by hand.
* `make fmt` runs `go fmt` on the project. Fmt ensures that code is laid out correctly according to Go standards. Fmt errors can usually be fixed automatically using the command in the failure message.
* `make lint` runs `go lint` on the project. Lint looks for coding style errors, such as public functions that are not documented correctly. Lint failures should be fixed by hand.
* `make copyright-check` runs a script to ensure that all source code files in the project have the correct copyright statement. You can run `make copyright` to fix issues if they are found.
* `make test` runs `go test` against each package in the project.
