# Running E2E Tests


1. Set `HUB` and `TAG` environment variables to work with a publically accessible image registry. E.g. `export HUB=liamwhite` and `export TAG=e2e`.
2. Export `KUBECONFIG` to point at an empty cluster.
3. Run `make e2e.local` (or `make e2e.local.armada`if testing IKS) to run the full suite of e2e tests. 
