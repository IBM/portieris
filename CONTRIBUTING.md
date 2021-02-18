---

copyright:
  years: 2018, 2021
lastupdated: "2021-02-18"

---

# Contributing to Portieris

## Reporting security issues

**Important** To report a security issue, don't open an issue. Instead, send your report by email privately to `alchreg@uk.ibm.com`.

## Reporting a general issue

Before you report an issue, check the [issues page](https://github.com/ibm/portieris/issues) to see whether an issue is already open. If an issue exists, don't open another one. Instead, react to the initial post inside that issue with a thumbs-up (+1) emoji. If you have any other debugging information, add a comment with that information.

If you don't find an existing issue, create one. Ensure that you take the following criteria into consideration:

- Use a title that accurately describes your problem, and allows other users to find it easily.
- If you're reporting an error, include any relevant output from your Portieris containers, and the message that you receive on the command line.
- If possible, create a workload that allows other people to quickly reproduce the issue, and include it in your report.
- Provide as much relevant information about your cluster's configuration as you can. For example, if your cluster is behind a firewall and you're having network connectivity problems, check whether your firewall is blocking the traffic and include any relevant logs.

## Contributing small changes

If you've found a bug or a typo and want to fix it yourself, report it first, but clearly indicate that you intend to fix it yourself. 

**Important** If the bug is a security issue, report it privately by email to `alchreg@uk.ibm.com` rather than opening an issue.

When you write the fix, add a test to reproduce the problem, and make sure that it passes the test before you submit your pull request (PR).

## Contributing major changes

Before you start work on a new feature or an architectural change, open an issue and describe what you plan to do, and, to open a dialogue with the maintainers, tag the issue with `enhancement`. After you and the rest of the community have agreed on what to do, start work and submit a PR.

Try to keep your PRs as small as possible. It's a lot easier to review many smaller PRs than one enormous one. If you want to make a large change, consider whether you can break the change down into a number of smaller changes.

## Changes to policy files

If you add new elements to the Portieris policy files, (`pkg/apis/securityenforcement/v1beta1/*`), you have to run the code generator, see [Generator Readme](pkg/apis/securityenforcement/v1beta1/README.md).

## Coding style

Write your code in Idiomatic Go. Check out [Effective Go](https://golang.org/doc/effective_go.html) for some tips about how to write Go effectively.

## Testing changes

Your code must pass a Travis check before it can be merged. Before you submit a PR, ensure that the tests pass. You can run the tests on your workstation by using the `make` command. To download the prerequisites for running the tests, run the `make test-deps` command from the project root.

To test your code, you can run the following commands:

- `make vet` runs the `go vet` command on the project. The `vet` command looks for issues in the code, such as ineffective variable assignments. Any `vet` failures must be fixed by hand.
- `make fmt` runs the `go fmt` command on the project. The `fmt` command ensures that code is laid out in the correct format according to Go standards. Any `fmt` errors can usually be fixed automatically by using the command that's in the failure message.
- `make lint` runs the `go lint` command on the project. The `lint` command looks for coding style errors, such as public functions that are not documented correctly. Any `lint` failures must be fixed by hand.
- `make copyright-check` runs a script to ensure that all source code files in the project have the correct copyright statement. If any issues are found, you can run the `make copyright` command to fix the issues.
- `make test` runs the `go test` command against each package in the project. The `test` command tests the package.
