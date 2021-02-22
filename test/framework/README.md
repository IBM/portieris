---

copyright:
  years: 2018, 2021
lastupdated: "2021-02-22"

---

# Framework

```
-- import "github.com/IBM/portieris/test/framework"
```

## Usage

### type Framework

```go
type Framework struct {
	KubeClient                     kubernetes.Interface
	ImagePolicyClient              securityenforcementclientset.Interface
	ClusterImagePolicyClient       securityenforcementclientset.Interface
	CustomResourceDefinitionClient customResourceDefinitionClientSet.CustomResourceDefinitionInterface
	HTTPClient                     *http.Client
	Namespace                      string
	HelmRelease                    string
	HelmChart                      string
}
```

`Framework` is an end-to-end test framework that is responsible for installing and deleting the Helm chart. It also provides helper functions for talking to Kube clusters.

### func  New

```go
func New(kubeconfig, helmChart string, noInstall bool) (*Framework, error)
```

`New` installs the specific Helm chart into the Kube cluster of the `kubeconfig`.

### func (*Framework) CreateClusterImagePolicy

```go
func (f *Framework) CreateClusterImagePolicy(clusterImagePolicy *v1beta1.ClusterImagePolicy) error
```

`CreateClusterImagePolicy` creates a cluster image policy.

### func (*Framework) CreateCronJob

```go
func (f *Framework) CreateCronJob(namespace string, job *batchv1.CronJob) error
```

`CreateCronJob` creates a scheduled job, CronJob, and waits for it to show.

### func (*Framework) CreateDaemonSet

```go
func (f *Framework) CreateDaemonSet(namespace string, daemonset *v1.DaemonSet) error
```

`CreateDaemonSet` creates a DaemonSet and waits for it to show.

### func (*Framework) CreateDeployment

```go
func (f *Framework) CreateDeployment(namespace string, deployment *v1.Deployment) error
```

`CreateDeployment` creates a deployment and waits for it to show.

### func (*Framework) CreateImagePolicy

```go
func (f *Framework) CreateImagePolicy(namespace string, imagePolicy *v1beta1.ImagePolicy) error
```

`CreateImagePolicy` creates an image policy.

### func (*Framework) CreateJob

```go
func (f *Framework) CreateJob(namespace string, job *batchv1.Job) error
```

`CreateJob` creates a job and waits for it to show.

### func (*Framework) CreateNamespace

```go
func (f *Framework) CreateNamespace(name string) (*corev1.Namespace, error)
```

`CreateNamespace` creates a namespace.

### func (*Framework) CreateNamespaceWithIPS

```go
func (f *Framework) CreateNamespaceWithIPS(name string) (*corev1.Namespace, error)
```

`CreateNamespaceWithIPS` creates a namespace, service account, and IP addresses to pull from the Global region of IBM Cloud Container Registry. The `bluemix-default-secret-international` ImagePullSecret from the default namespace is used.

### func (*Framework) CreatePod

```go
func (f *Framework) CreatePod(namespace string, pod *corev1.Pod) error
```

`CreatePod` creates a pod and waits for it to show.

### func (*Framework) CreateReplicaSet

```go
func (f *Framework) CreateReplicaSet(namespace string, replicaset *v1.ReplicaSet) error
```

`CreateReplicaSet` creates a ReplicaSet and waits for it to show.

### func (*Framework) CreateReplicationController

```go
func (f *Framework) CreateReplicationController(namespace string, replicationcontroller *corev1.ReplicationController) error
```

`CreateReplicationController` creates a replication controller and waits for it to show.

### func (*Framework) CreateSecret

```go
func (f *Framework) CreateSecret(namespace string, secret *corev1.Secret) error
```

`CreateSecret` creates a secret and waits for it to show.

### func (*Framework) CreateStatefulSet

```go
func (f *Framework) CreateStatefulSet(namespace string, statefulset *v1.StatefulSet) error
```

`CreateStatefulSet` creates a StatefulSet and waits for it to show.

### func (*Framework) DeleteClusterImagePolicy

```go
func (f *Framework) DeleteClusterImagePolicy(name string) error
```

`DeleteClusterImagePolicy` deletes the specified cluster image policy.

### func (*Framework) DeleteCronJob

```go
func (f *Framework) DeleteCronJob(name, namespace string) error
```

`DeleteCronJob` deletes the specified CronJob.

### func (*Framework) DeleteDaemonSet

```go
func (f *Framework) DeleteDaemonSet(name, namespace string) error
```

`DeleteDaemonSet` deletes the specified DaemonSet.

### func (*Framework) DeleteDeployment

```go
func (f *Framework) DeleteDeployment(name, namespace string) error
```

`DeleteDeployment` deletes the specified deployment.

### func (*Framework) DeleteImagePolicy

```go
func (f *Framework) DeleteImagePolicy(name, namespace string) error
```

`DeleteImagePolicy` deletes the image policy.

### func (*Framework) DeleteJob

```go
func (f *Framework) DeleteJob(name, namespace string) error
```

`DeleteJob` deletes the specified job.

### func (*Framework) DeleteNamespace

```go
func (f *Framework) DeleteNamespace(name string) error
```

`DeleteNamespace` deletes the specified namespace.

### func (*Framework) DeletePod

```go
func (f *Framework) DeletePod(name, namespace string) error
```

`DeletePod` deletes the specified pod.

### func (*Framework) DeleteRandomPod

```go
func (f *Framework) DeleteRandomPod(namespace string) error
```

`DeleteRandomPod` deletes the first pod that is returned in the list of pods for a specified namespace.

### func (*Framework) DeleteReplicaSet

```go
func (f *Framework) DeleteReplicaSet(name, namespace string) error
```

`DeleteReplicaSet` deletes the specified ReplicaSet.

### func (*Framework) DeleteReplicationController

```go
func (f *Framework) DeleteReplicationController(name, namespace string) error
```

`DeleteReplicationController` deletes the specified replication controller.

### func (*Framework) DeleteStatefulSet

```go
func (f *Framework) DeleteStatefulSet(name, namespace string) error
```

`DeleteStatefulSet` deletes the specified StatefulSet.

### func (*Framework) DumpEvents

```go
func (f *Framework) DumpEvents(namespace string) io.Reader
```

`DumpEvents` returns a reader that has events that are written to a specified namespace.

### func (*Framework) DumpPolicies

```go
func (f *Framework) DumpPolicies(namespace string) io.Reader
```

`DumpPolicies` returns a reader that has all cluster and image policies present in it.

### func (*Framework) GenerateTestAnnotation

```go
func (f *Framework) GenerateTestAnnotation() string
```

`GenerateTestAnnotation` returns a unique test annotation that is used to patch resources.

### func (*Framework) GetClusterImagePolicy

```go
func (f *Framework) GetClusterImagePolicy(name string) (*v1beta1.ClusterImagePolicy, error)
```

`GetClusterImagePolicy` retrieves the cluster image policy.

### func (*Framework) GetClusterImagePolicyDefinition

```go
func (f *Framework) GetClusterImagePolicyDefinition() (*apiextensions.CustomResourceDefinition, error)
```

`GetClusterImagePolicyDefinition` retrieves the cluster image policy custom resource definition.

### func (*Framework) GetCronJob

```go
func (f *Framework) GetCronJob(name, namespace string) (*batchv1.CronJob, error)
```

`GetCronJob` retrieves the specified CronJob.

### func (*Framework) GetDaemonSets

```go
func (f *Framework) GetDaemonSets(name, namespace string) (*v1.DaemonSet, error)
```

`GetDaemonSets` retrieves the specified DaemonSets.

### func (*Framework) GetDeployment

```go
func (f *Framework) GetDeployment(name, namespace string) (*v1.Deployment, error)
```

`GetDeployment` retrieves the specified deployment.

### func (*Framework) GetImagePolicy

```go
func (f *Framework) GetImagePolicy(name, namespace string) (*v1beta1.ImagePolicy, error)
```

`GetImagePolicy` retrieves the image policy.

### func (*Framework) GetImagePolicyDefinition

```go
func (f *Framework) GetImagePolicyDefinition() (*apiextensions.CustomResourceDefinition, error)
```

`GetImagePolicyDefinition` retrieves the image policy custom resource definition.

### func (*Framework) GetJob

```go
func (f *Framework) GetJob(name, namespace string) (*batchv1.Job, error)
```

`GetJob` retrieves the specified job.

### func (*Framework) GetNamespace

```go
func (f *Framework) GetNamespace(name string) (*corev1.Namespace, error)
```

`GetNamespace` retrieves the specified namespace.

### func (*Framework) GetPod

```go
func (f *Framework) GetPod(name, namespace string) (*corev1.Pod, error)
```

`GetPod` retrieves the specified pod.

### func (*Framework) GetReplicaSet

```go
func (f *Framework) GetReplicaSet(name, namespace string) (*v1.ReplicaSet, error)
```

`GetReplicaSet` retrieves the specified ReplicaSet.

### func (*Framework) GetReplicationController

```go
func (f *Framework) GetReplicationController(name, namespace string) (*corev1.ReplicationController, error)
```

`GetReplicationController` retrieves the specified replication controller.

### func (*Framework) GetSecret

```go
func (f *Framework) GetSecret(name, namespace string) (*corev1.Secret, error)
```

`GetSecret` retrieves the specified secret.

### func (*Framework) GetStatefulSet

```go
func (f *Framework) GetStatefulSet(name, namespace string) (*v1.StatefulSet, error)
```

`GetStatefulSet` retrieves the specified StatefulSet.

### func (*Framework) ListClusterImagePolicies

```go
func (f *Framework) ListClusterImagePolicies() (*v1beta1.ClusterImagePolicyList, error)
```

`ListClusterImagePolicies` lists the cluster image policies.

### func (*Framework) ListClusterRoleBindings

```go
func (f *Framework) ListClusterRoleBindings() (*v1beta1.ClusterRoleBindingList, error)
```

`ListClusterRoleBindings` lists all the cluster role bindings that are associated with the installed Helm release.

### func (*Framework) ListClusterRoles

```go
func (f *Framework) ListClusterRoles() (*v1beta1.ClusterRoleList, error)
```

`ListClusterRoles` lists all the cluster roles that are associated with the installed Helm release.

### func (*Framework) ListConfigMaps

```go
func (f *Framework) ListConfigMaps() (*corev1.ConfigMapList, error)
```

`ListConfigMaps` lists all the configuration maps that are associated with the installed Helm release.

### func (*Framework) ListCronJobs

```go
func (f *Framework) ListCronJobs() (*batchv1.CronJobList, error)
```

`ListCronJobs` lists all the CronJobs that are associated with the installed Helm release.

### func (*Framework) ListDaemonSet

```go
func (f *Framework) ListDaemonSet() (*v1.DaemonSetList, error)
```
`ListDaemonSet` lists all DaemonSets that are associated with the installed Helm release.

### func (*Framework) ListDeployments

```go
func (f *Framework) ListDeployments() (*v1.DeploymentList, error)
```

`ListDeployments` lists all deployments that are associated with the installed Helm release.

### func (*Framework) ListImagePolicies

```go
func (f *Framework) ListImagePolicies(namespace string) (*v1beta1.ImagePolicyList, error)
```

`ListImagePolicies` lists all image polices in a specified namespace.

### func (*Framework) ListJobs

```go
func (f *Framework) ListJobs() (*batchv1.JobList, error)
```

`ListJobs` lists all jobs that are associated with the installed Helm release.

### func (*Framework) ListMutatingAdmissionWebhooks

```go
func (f *Framework) ListMutatingAdmissionWebhooks() (*v1beta1.MutatingWebhookConfigurationList, error)
```

`ListMutatingAdmissionWebhooks` lists all mutating admission webhooks that are associated with the installed Helm release.

### func (*Framework) ListReplicaSet

```go
func (f *Framework) ListReplicaSet() (*v1.ReplicaSetList, error)
```

`ListReplicaSet` lists all ReplicaSets that are associated with the installed Helm release.

### func (*Framework) ListReplicationController

```go
func (f *Framework) ListReplicationController() (*corev1.ReplicationControllerList, error)
```

`ListReplicationController` lists all replication controllers that are associated with the installed Helm release.

### func (*Framework) ListServiceAccounts

```go
func (f *Framework) ListServiceAccounts() (*corev1.ServiceAccountList, error)
```

`ListServiceAccounts` lists all service accounts that are associated with the installed Helm release.

### func (*Framework) ListServices

```go
func (f *Framework) ListServices() (*corev1.ServiceList, error)
```

`ListServices` lists all services that are associated with the installed Helm release.

### func (*Framework) ListStatefulSet

```go
func (f *Framework) ListStatefulSet() (*v1.StatefulSetList, error)
```

`ListStatefulSet` lists all StatefulSets that are associated with the installed Helm release.

### func (*Framework) ListValidatingAdmissionWebhooks

```go
func (f *Framework) ListValidatingAdmissionWebhooks() (*v1beta1.ValidatingWebhookConfigurationList, error)
```

`ListValidatingAdmissionWebhooks` lists all validating admission webhooks that are associated with the installed Helm release.

### func (*Framework) LoadClusterImagePolicyManifest

```go
func (f *Framework) LoadClusterImagePolicyManifest(pathToManifest string) (*v1beta1.ClusterImagePolicy, error)
```

`LoadClusterImagePolicyManifest` takes a manifest and decodes it into a cluster image policy object.

### func (*Framework) LoadCronJobManifest

```go
func (f *Framework) LoadCronJobManifest(pathToManifest string) (*batchv1.CronJob, error)
```
`LoadCronJobManifest` takes a manifest and decodes it into a CronJob object.

### func (*Framework) LoadDaemonSetManifest

```go
func (f *Framework) LoadDaemonSetManifest(pathToManifest string) (*v1.DaemonSet, error)
```

`LoadDaemonSetManifest` takes a manifest and decodes it into a DaemonSet object.

### func (*Framework) LoadDeploymentManifest

```go
func (f *Framework) LoadDeploymentManifest(pathToManifest string) (*v1.Deployment, error)
```

`LoadDeploymentManifest` takes a manifest and decodes it into a deployment object.

### func (*Framework) LoadImagePolicyManifest

```go
func (f *Framework) LoadImagePolicyManifest(pathToManifest string) (*v1beta1.ImagePolicy, error)
```

`LoadImagePolicyManifest` takes a manifest and decodes it into an image policy object.

### func (*Framework) LoadJobManifest

```go
func (f *Framework) LoadJobManifest(pathToManifest string) (*batchv1.Job, error)
```

`LoadJobManifest` takes a manifest and decodes it into a job object.

### func (*Framework) LoadPodManifest

```go
func (f *Framework) LoadPodManifest(pathToManifest string) (*corev1.Pod, error)
```

`LoadPodManifest` takes a manifest and decodes it into a pod object.

### func (*Framework) LoadReplicaSetManifest

```go
func (f *Framework) LoadReplicaSetManifest(pathToManifest string) (*v1.ReplicaSet, error)
```

`LoadReplicaSetManifest` takes a manifest and decodes it into a ReplicaSet object.

### func (*Framework) LoadReplicationControllerManifest

```go
func (f *Framework) LoadReplicationControllerManifest(pathToManifest string) (*corev1.ReplicationController, error)
```

`LoadReplicationControllerManifest` takes a manifest and decodes it into a replication controller object.

### func (*Framework) LoadSecretManifest

```go
func (f *Framework) LoadSecretManifest(pathToManifest string) (*corev1.Secret, error)
```

`LoadSecretManifest` takes a manifest and decodes it into a secret object.

### func (*Framework) LoadStatefulSetManifest

```go
func (f *Framework) LoadStatefulSetManifest(pathToManifest string) (*v1.StatefulSet, error)
```

`LoadStatefulSetManifest` takes a manifest and decodes it into a StatefulSet object.

### func (*Framework) PatchDeployment

```go
func (f *Framework) PatchDeployment(name, namespace, patch string) (*v1.Deployment, error)
```

`PatchDeployment` patches the specified deployment.

### func (*Framework) ReplaceDeployment

```go
func (f *Framework) ReplaceDeployment(namespace string, deployment *v1.Deployment) (*v1.Deployment, error)
```

`ReplaceDeployment` replaces the specified deployment.

### func (*Framework) Teardown

```go
func (f *Framework) Teardown() bool
```

`Teardown` deletes the chart and verifies that everything is cleaned up.

### func (*Framework) UpdateImagePolicy

```go
func (f *Framework) UpdateImagePolicy(namespace string, imagePolicy *v1beta1.ImagePolicy) error
```

`UpdateImagePolicy` updates the image policy.

### func (*Framework) WaitForClusterImagePolicy

```go
func (f *Framework) WaitForClusterImagePolicy(name string, timeout time.Duration) error
```

`WaitForClusterImagePolicy` waits until the cluster image policy is created or the timeout is reached.

### func (*Framework) WaitForClusterImagePolicyDefinition

```go
func (f *Framework) WaitForClusterImagePolicyDefinition(timeout time.Duration) error
```

`WaitForClusterImagePolicyDefinition` waits until the cluster image policy customer resource definition is created or the timeout is reached.

### func (*Framework) WaitForCronJob

```go
func (f *Framework) WaitForCronJob(name, namespace string, timeout time.Duration) error
```

`WaitForCronJob` waits until the CronJob deployment is complete.

### func (*Framework) WaitForDaemonSet

```go
func (f *Framework) WaitForDaemonSet(name, namespace string, timeout time.Duration) error
```

`WaitForDaemonSet` waits until the specified DaemonSet is created or the timeout is reached.

### func (*Framework) WaitForDaemonSetPods

```go
func (f *Framework) WaitForDaemonSetPods(name, namespace string, timeout time.Duration) error
```

`WaitForDaemonSetPods` waits until the specified DeamonSet pods are created or the timeout is reached.

### func (*Framework) WaitForDeployment

```go
func (f *Framework) WaitForDeployment(name, namespace string, timeout time.Duration) error
```

`WaitForDeployment` waits until the specified deployment is created or the timeout is reached.

### func (*Framework) WaitForDeploymentPods

```go
func (f *Framework) WaitForDeploymentPods(name, namespace string, timeout time.Duration) error
```

`WaitForDeploymentPods` waits until the specified deployment's pods are created or the timeout is reached.

### func (*Framework) WaitForImagePolicy

```go
func (f *Framework) WaitForImagePolicy(name, namespace string, timeout time.Duration) error
```

`WaitForImagePolicy` waits until the image policy is created or the timeout is reached.

### func (*Framework) WaitForImagePolicyDefinition

```go
func (f *Framework) WaitForImagePolicyDefinition(timeout time.Duration) error
```

`WaitForImagePolicyDefinition` waits until the image policy customer resource definition is created or the timeout is reached.

### func (*Framework) WaitForJob

```go
func (f *Framework) WaitForJob(name, namespace string, timeout time.Duration) error
```

`WaitForJob` waits until the job deployment completes.

### func (*Framework) WaitForMutatingAdmissionWebhook

```go
func (f *Framework) WaitForMutatingAdmissionWebhook(name string, timeout time.Duration) error
```

`WaitForMutatingAdmissionWebhook` waits until the specified mutating admission webhook is created or the timeout is reached.

### func (*Framework) WaitForNamespace

```go
func (f *Framework) WaitForNamespace(name string, timeout time.Duration) error
```

`WaitForNamespace` waits until the specified namespace is created or the timeout is reached.

### func (*Framework) WaitForPod

```go
func (f *Framework) WaitForPod(name, namespace string, timeout time.Duration) error
```

`WaitForPod` waits until the pod deployment completes.

### func (*Framework) WaitForPodDelete

```go
func (f *Framework) WaitForPodDelete(name, namespace string, timeout time.Duration) error
```

`WaitForPodDelete` waits until the pod is deleted.

### func (*Framework) WaitForReplicaSet

```go
func (f *Framework) WaitForReplicaSet(name, namespace string, timeout time.Duration) error
```

`WaitForReplicaSet` waits until the specified ReplicaSet is created or the timeout is reached.

### func (*Framework) WaitForReplicaSetPods

```go
func (f *Framework) WaitForReplicaSetPods(name, namespace string, timeout time.Duration) error
```

`WaitForReplicaSetPods` waits until the specified ReplicaSet pods are created or the timeout is reached.

### func (*Framework) WaitForReplicationController

```go
func (f *Framework) WaitForReplicationController(name, namespace string, timeout time.Duration) error
```

`WaitForReplicationController` waits until the specified replication controller is created or the timeout is reached.

### func (*Framework) WaitForReplicationControllerPods

```go
func (f *Framework) WaitForReplicationControllerPods(name, namespace string, timeout time.Duration) error
```

`WaitForReplicationControllerPods` waits until the specified replication controller pods are created or the timeout is reached.

### func (*Framework) WaitForSecret

```go
func (f *Framework) WaitForSecret(name, namespace string, timeout time.Duration) error
```

`WaitForSecret` waits until the specified secret is created or the timeout is reached.

### func (*Framework) WaitForStatefulSet

```go
func (f *Framework) WaitForStatefulSet(name, namespace string, timeout time.Duration) error
```

`WaitForStatefulSet` waits until the specified StatefulSet is created or the timeout is reached.

### func (*Framework) WaitForStatefulSetPods

```go
func (f *Framework) WaitForStatefulSetPods(name, namespace string, timeout time.Duration) error
```

`WaitForStatefulSetPods` waits until the specified StatefulSet pods are created or the timeout is reached.

### func (*Framework) WaitForValidatingAdmissionWebhook

```go
func (f *Framework) WaitForValidatingAdmissionWebhook(name string, timeout time.Duration) error
```

`WaitForValidatingAdmissionWebhook` waits until the specified validating admission webhook is created or the timeout is reached.
