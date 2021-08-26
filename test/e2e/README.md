---

copyright:
  years: 2020, 2021
lastupdated: "2021-08-26"

---

# E2E tests

The end-to-end (E2E) tests run against a Kubernetes cluster. The tests are typically run by maintainers against an IBM Cloud Kubernetes Service cluster before a pull request (PR) is accepted.

Before you begin, contact a maintainer to find out what to set for exporting the `E2E_ACCOUNT_HEADER`. You'll need this information for one of the steps.

Run a test by completing the following steps.

1. Order an IBM Cloud Kubernetes Service cluster, see [Creating clusters](https://cloud.ibm.com/docs/containers?topic=containers-clusters), or use an existing cluster.
2. Get the cluster configuration by running the following [`ibmcloud ks cluster config`](https://cloud.ibm.com/docs/containers?topic=containers-kubernetes-service-cli#cs_cluster_config) command. 
   
   ```
   ibmcloud ks cluster config -c <cluster-name>
   ```
   
3. Export KUBECONFIG to point at the `kube-config.yaml` file for your cluster that is in `~/.bluemix/plugins/container-service/clusters`.

   For example,

   ```
   export KUBECONFIG=~/.bluemix/plugins/container-service/clusters/mycluster-xxxx/kube-config.yaml
   ```

5. Create an IBM Cloud Container Registry namespace that is owned by the same account as the cluster, see [Setting up a namespace](https://cloud.ibm.com/docs/Registry?topic=Registry-registry_setup_cli_namespace#registry_namespace_setup).
6. Export the `HUB` variable that points to your IBM Cloud Container Registry namespace.
   
   For example,
   
   ```
   export HUB=uk.icr.io/yournamespace
   ```
   
8. After completing your code changes, build and push the image to your namespace by running the following command. 
   
   ```
   make push
   ```
   
8. Install Portieris into your cluster.

   ```
   helm install helm/portieris
   ```
   
11. Export `E2E_ACCOUNT_HEADER`. (Contact a maintainer to find out what to set.)
12. Run the following command.
   
    ```
    make e2e.quick
    ```
