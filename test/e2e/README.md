# E2E Tests

The E2E tests execute against a real Kubernetes cluster and are typically run by maintainers against an IBM Cloud Kubernetes Service (IKS) cluster before accepting a PR.

The tests can be run with the following steps:

1. [Order a new IKS cluster](https://cloud.ibm.com/docs/containers-cli-plugin?topic=containers-cli-plugin-kubernetes-service-cli#cs_cluster_create) (or use an existing cluster)
1. Get the cluster config: `ibmcloud ks cluster config -c <cluster-name>`
1. Export KUBECONFIG to point at the `kube-config.yaml` for your cluster found in `~/.bluemix/plugins/container-service/clusters`.
1. Create a new IBM Cloud Container Registry (ICCR) namespace owned by the same account as the cluster
1. Export the HUB variable pointing to your new ICCR namespace: `export HUB=uk.icr.io/yournamespace`
1. After completing your code changes, build and push the image to your namespace: `make push`
1. Install portieris into your cluster: `helm install helm/portieris`
1. Run `make e2e.quick`
