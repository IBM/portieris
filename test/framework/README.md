# framework
--
    import "github.com/IBM/portieris/test/framework"


## Usage

#### type Framework

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

Framework is an e2e test framework esponsible for installing and deleting of the
helm chart It also providers helper functions for talking to Kube clusters

#### func  New

```go
func New(kubeconfig, helmChart string, noInstall bool) (*Framework, error)
```
New installs the specific helm chart into the Kube cluster of the kubeconfig

#### func (*Framework) CreateClusterImagePolicy

```go
func (f *Framework) CreateClusterImagePolicy(clusterImagePolicy *v1beta1.ClusterImagePolicy) error
```
CreateClusterImagePolicy creates the ClusterImagePolicy

#### func (*Framework) CreateCronJob

```go
func (f *Framework) CreateCronJob(namespace string, job *batchv1.CronJob) error
```
CreateCronJob creates a CronJob resource and then waits for it to appear

#### func (*Framework) CreateDaemonSet

```go
func (f *Framework) CreateDaemonSet(namespace string, daemonset *v1.DaemonSet) error
```
CreateDaemonSet creates a daemonset resource and then waits for it to appear

#### func (*Framework) CreateDeployment

```go
func (f *Framework) CreateDeployment(namespace string, deployment *v1.Deployment) error
```
CreateDeployment creates a deployment resource and then waits for it to appear

#### func (*Framework) CreateImagePolicy

```go
func (f *Framework) CreateImagePolicy(namespace string, imagePolicy *v1beta1.ImagePolicy) error
```
CreateImagePolicy creates the ImagePolicy

#### func (*Framework) CreateJob

```go
func (f *Framework) CreateJob(namespace string, job *batchv1.Job) error
```
CreateJob creates a Job resource and then waits for it to appear

#### func (*Framework) CreateNamespace

```go
func (f *Framework) CreateNamespace(name string) (*corev1.Namespace, error)
```
CreateNamespace creates a namespace

#### func (*Framework) CreateNamespaceWithIPS

```go
func (f *Framework) CreateNamespaceWithIPS(name string) (*corev1.Namespace, error)
```
CreateNamespaceWithIPS creates a namespace, service account and IPS to pull from
the IBM Cloud Container Registry Global region It uses the
bluemix-default-secret-international imagePullSecret from the default namespace

#### func (*Framework) CreatePod

```go
func (f *Framework) CreatePod(namespace string, pod *corev1.Pod) error
```
CreatePod creates a Replicaset resource and then waits for it to appear

#### func (*Framework) CreateReplicaSet

```go
func (f *Framework) CreateReplicaSet(namespace string, replicaset *v1.ReplicaSet) error
```
CreateReplicaSet creates a Replicaset resource and then waits for it to appear

#### func (*Framework) CreateReplicationController

```go
func (f *Framework) CreateReplicationController(namespace string, replicationcontroller *corev1.ReplicationController) error
```
CreateReplicationController creates a Replicaset resource and then waits for it
to appear

#### func (*Framework) CreateSecret

```go
func (f *Framework) CreateSecret(namespace string, secret *corev1.Secret) error
```
CreateSecret creates a secret resource and then waits for it to appear

#### func (*Framework) CreateStatefulSet

```go
func (f *Framework) CreateStatefulSet(namespace string, statefulset *v1.StatefulSet) error
```
CreateStatefulSet creates a StatefulSet resource and then waits for it to appear

#### func (*Framework) DeleteClusterImagePolicy

```go
func (f *Framework) DeleteClusterImagePolicy(name string) error
```
DeleteClusterImagePolicy deletes the specified ClusterImagePolicy

#### func (*Framework) DeleteCronJob

```go
func (f *Framework) DeleteCronJob(name, namespace string) error
```
DeleteCronJob deletes the specified deployment

#### func (*Framework) DeleteDaemonSet

```go
func (f *Framework) DeleteDaemonSet(name, namespace string) error
```
DeleteDaemonSet deletes the specified deployment

#### func (*Framework) DeleteDeployment

```go
func (f *Framework) DeleteDeployment(name, namespace string) error
```
DeleteDeployment deletes the specified deployment

#### func (*Framework) DeleteImagePolicy

```go
func (f *Framework) DeleteImagePolicy(name, namespace string) error
```
DeleteImagePolicy deletes the ImagePolicy

#### func (*Framework) DeleteJob

```go
func (f *Framework) DeleteJob(name, namespace string) error
```
DeleteJob deletes the specified deployment

#### func (*Framework) DeleteNamespace

```go
func (f *Framework) DeleteNamespace(name string) error
```
DeleteNamespace deletes the specified namespace

#### func (*Framework) DeletePod

```go
func (f *Framework) DeletePod(name, namespace string) error
```
DeletePod deletes the specified deployment

#### func (*Framework) DeleteRandomPod

```go
func (f *Framework) DeleteRandomPod(namespace string) error
```
DeleteRandomPod deletes first pod returned in pod list for a given namespace

#### func (*Framework) DeleteReplicaSet

```go
func (f *Framework) DeleteReplicaSet(name, namespace string) error
```
DeleteReplicaSet deletes the specified deployment

#### func (*Framework) DeleteReplicationController

```go
func (f *Framework) DeleteReplicationController(name, namespace string) error
```
DeleteReplicationController deletes the specified deployment

#### func (*Framework) DeleteStatefulSet

```go
func (f *Framework) DeleteStatefulSet(name, namespace string) error
```
DeleteStatefulSet deletes the specified deployment

#### func (*Framework) DumpEvents

```go
func (f *Framework) DumpEvents(namespace string) io.Reader
```
DumpEvents returns a reader that will have events for a given namespace written
to

#### func (*Framework) DumpPolicies

```go
func (f *Framework) DumpPolicies(namespace string) io.Reader
```
DumpPolicies returns a reader that will have all cluster and image policies
present in it

#### func (*Framework) GenerateTestAnnotation

```go
func (f *Framework) GenerateTestAnnotation() string
```
GenerateTestAnnotation returns a unique test annotation for patching resources

#### func (*Framework) GetClusterImagePolicy

```go
func (f *Framework) GetClusterImagePolicy(name string) (*v1beta1.ClusterImagePolicy, error)
```
GetClusterImagePolicy retrieves the ClusterImagePolicy

#### func (*Framework) GetClusterImagePolicyDefinition

```go
func (f *Framework) GetClusterImagePolicyDefinition() (*apiextensions.CustomResourceDefinition, error)
```
GetClusterImagePolicyDefinition retrieves the ClusterImagePolicy CRD

#### func (*Framework) GetCronJob

```go
func (f *Framework) GetCronJob(name, namespace string) (*batchv1.CronJob, error)
```
GetCronJob retrieves the specified deployment

#### func (*Framework) GetDaemonSets

```go
func (f *Framework) GetDaemonSets(name, namespace string) (*v1.DaemonSet, error)
```
GetDaemonSets retrieves the specified deployment

#### func (*Framework) GetDeployment

```go
func (f *Framework) GetDeployment(name, namespace string) (*v1.Deployment, error)
```
GetDeployment retrieves the specified deployment

#### func (*Framework) GetImagePolicy

```go
func (f *Framework) GetImagePolicy(name, namespace string) (*v1beta1.ImagePolicy, error)
```
GetImagePolicy retrieves the ImagePolicy

#### func (*Framework) GetImagePolicyDefinition

```go
func (f *Framework) GetImagePolicyDefinition() (*apiextensions.CustomResourceDefinition, error)
```
GetImagePolicyDefinition retrieves the ImagePolicy CRD

#### func (*Framework) GetJob

```go
func (f *Framework) GetJob(name, namespace string) (*batchv1.Job, error)
```
GetJob retrieves the specified deployment

#### func (*Framework) GetNamespace

```go
func (f *Framework) GetNamespace(name string) (*corev1.Namespace, error)
```
GetNamespace retrieves the specified namespace

#### func (*Framework) GetPod

```go
func (f *Framework) GetPod(name, namespace string) (*corev1.Pod, error)
```
GetPod retrieves the specified deployment

#### func (*Framework) GetReplicaSet

```go
func (f *Framework) GetReplicaSet(name, namespace string) (*v1.ReplicaSet, error)
```
GetReplicaSet retrieves the specified deployment

#### func (*Framework) GetReplicationController

```go
func (f *Framework) GetReplicationController(name, namespace string) (*corev1.ReplicationController, error)
```
GetReplicationController retrieves the specified deployment

#### func (*Framework) GetSecret

```go
func (f *Framework) GetSecret(name, namespace string) (*corev1.Secret, error)
```
GetSecret retrieves the specified secret

#### func (*Framework) GetStatefulSet

```go
func (f *Framework) GetStatefulSet(name, namespace string) (*v1.StatefulSet, error)
```
GetStatefulSet retrieves the specified deployment

#### func (*Framework) ListClusterImagePolicies

```go
func (f *Framework) ListClusterImagePolicies() (*v1beta1.ClusterImagePolicyList, error)
```
ListClusterImagePolicies creates the ClusterImagePolicy

#### func (*Framework) ListClusterRoleBindings

```go
func (f *Framework) ListClusterRoleBindings() (*v1beta1.ClusterRoleBindingList, error)
```
ListClusterRoleBindings retrieves all cluster role bindings associated with the
installed Helm release

#### func (*Framework) ListClusterRoles

```go
func (f *Framework) ListClusterRoles() (*v1beta1.ClusterRoleList, error)
```
ListClusterRoles retrieves all cluster roles associated with the installed Helm
release

#### func (*Framework) ListConfigMaps

```go
func (f *Framework) ListConfigMaps() (*corev1.ConfigMapList, error)
```
ListConfigMaps retrieves all config maps associated with the installed Helm
release

#### func (*Framework) ListCronJobs

```go
func (f *Framework) ListCronJobs() (*batchv1.CronJobList, error)
```
ListCronJobs retrieves all jobs associated with the installed Helm release

#### func (*Framework) ListDaemonSet

```go
func (f *Framework) ListDaemonSet() (*v1.DaemonSetList, error)
```
ListDaemonSet retrieves all daemonset associated with the installed Helm release

#### func (*Framework) ListDeployments

```go
func (f *Framework) ListDeployments() (*v1.DeploymentList, error)
```
ListDeployments retrieves all deployments associated with the installed Helm
release

#### func (*Framework) ListImagePolicies

```go
func (f *Framework) ListImagePolicies(namespace string) (*v1beta1.ImagePolicyList, error)
```
ListImagePolicies lists all ImagePolicies in a given namespace

#### func (*Framework) ListJobs

```go
func (f *Framework) ListJobs() (*batchv1.JobList, error)
```
ListJobs retrieves all jobs associated with the installed Helm release

#### func (*Framework) ListMutatingAdmissionWebhooks

```go
func (f *Framework) ListMutatingAdmissionWebhooks() (*v1beta1.MutatingWebhookConfigurationList, error)
```
ListMutatingAdmissionWebhooks retrieves all Mutating Admission Webhooks
associated with the installed Helm release

#### func (*Framework) ListReplicaSet

```go
func (f *Framework) ListReplicaSet() (*v1.ReplicaSetList, error)
```
ListReplicaSet retrieves all Replicaset associated with the installed Helm
release

#### func (*Framework) ListReplicationController

```go
func (f *Framework) ListReplicationController() (*corev1.ReplicationControllerList, error)
```
ListReplicationController retrieves all Replicaset associated with the installed
Helm release

#### func (*Framework) ListServiceAccounts

```go
func (f *Framework) ListServiceAccounts() (*corev1.ServiceAccountList, error)
```
ListServiceAccounts retrieves all service accounts associated with the installed
Helm release

#### func (*Framework) ListServices

```go
func (f *Framework) ListServices() (*corev1.ServiceList, error)
```
ListServices retrieves all services associated with the installed Helm release

#### func (*Framework) ListStatefulSet

```go
func (f *Framework) ListStatefulSet() (*v1.StatefulSetList, error)
```
ListStatefulSet retrieves all StatefulSet associated with the installed Helm
release

#### func (*Framework) ListValidatingAdmissionWebhooks

```go
func (f *Framework) ListValidatingAdmissionWebhooks() (*v1beta1.ValidatingWebhookConfigurationList, error)
```
ListValidatingAdmissionWebhooks retrieves all ValidatingAdmissionWebhooks
associated with the installed Helm release

#### func (*Framework) LoadClusterImagePolicyManifest

```go
func (f *Framework) LoadClusterImagePolicyManifest(pathToManifest string) (*v1beta1.ClusterImagePolicy, error)
```
LoadClusterImagePolicyManifest takes a manifest and decodes it into a
ImagePolicy object

#### func (*Framework) LoadCronJobManifest

```go
func (f *Framework) LoadCronJobManifest(pathToManifest string) (*batchv1.CronJob, error)
```
LoadCronJobManifest takes a manifest and decodes it into a CronJob object

#### func (*Framework) LoadDaemonSetManifest

```go
func (f *Framework) LoadDaemonSetManifest(pathToManifest string) (*v1.DaemonSet, error)
```
LoadDaemonSetManifest takes a manifest and decodes it into a daemonset object

#### func (*Framework) LoadDeploymentManifest

```go
func (f *Framework) LoadDeploymentManifest(pathToManifest string) (*v1.Deployment, error)
```
LoadDeploymentManifest takes a manifest and decodes it into a deployment object

#### func (*Framework) LoadImagePolicyManifest

```go
func (f *Framework) LoadImagePolicyManifest(pathToManifest string) (*v1beta1.ImagePolicy, error)
```
LoadImagePolicyManifest takes a manifest and decodes it into a ImagePolicy
object

#### func (*Framework) LoadJobManifest

```go
func (f *Framework) LoadJobManifest(pathToManifest string) (*batchv1.Job, error)
```
LoadJobManifest takes a manifest and decodes it into a Job object

#### func (*Framework) LoadPodManifest

```go
func (f *Framework) LoadPodManifest(pathToManifest string) (*corev1.Pod, error)
```
LoadPodManifest takes a manifest and decodes it into a Replicaset object

#### func (*Framework) LoadReplicaSetManifest

```go
func (f *Framework) LoadReplicaSetManifest(pathToManifest string) (*v1.ReplicaSet, error)
```
LoadReplicaSetManifest takes a manifest and decodes it into a Replicaset object

#### func (*Framework) LoadReplicationControllerManifest

```go
func (f *Framework) LoadReplicationControllerManifest(pathToManifest string) (*corev1.ReplicationController, error)
```
LoadReplicationControllerManifest takes a manifest and decodes it into a
Replicaset object

#### func (*Framework) LoadSecretManifest

```go
func (f *Framework) LoadSecretManifest(pathToManifest string) (*corev1.Secret, error)
```
LoadSecretManifest takes a manifest and decodes it into a deployment object

#### func (*Framework) LoadStatefulSetManifest

```go
func (f *Framework) LoadStatefulSetManifest(pathToManifest string) (*v1.StatefulSet, error)
```
LoadStatefulSetManifest takes a manifest and decodes it into a StatefulSet
object

#### func (*Framework) PatchDeployment

```go
func (f *Framework) PatchDeployment(name, namespace, patch string) (*v1.Deployment, error)
```
PatchDeployment patches the specified deployment

#### func (*Framework) ReplaceDeployment

```go
func (f *Framework) ReplaceDeployment(namespace string, deployment *v1.Deployment) (*v1.Deployment, error)
```
ReplaceDeployment patches the specified deployment

#### func (*Framework) Teardown

```go
func (f *Framework) Teardown() bool
```
Teardown deletes the chart and then verifies everything has been cleaned up

#### func (*Framework) UpdateImagePolicy

```go
func (f *Framework) UpdateImagePolicy(namespace string, imagePolicy *v1beta1.ImagePolicy) error
```
UpdateImagePolicy creates the ImagePolicy

#### func (*Framework) WaitForClusterImagePolicy

```go
func (f *Framework) WaitForClusterImagePolicy(name string, timeout time.Duration) error
```
WaitForClusterImagePolicy waits until the ClusterImagePolicy is created or the
timeout is reached

#### func (*Framework) WaitForClusterImagePolicyDefinition

```go
func (f *Framework) WaitForClusterImagePolicyDefinition(timeout time.Duration) error
```
WaitForClusterImagePolicyDefinition waits until the ClusterImagePolicy CRD is
created or the timeout is reached

#### func (*Framework) WaitForCronJob

```go
func (f *Framework) WaitForCronJob(name, namespace string, timeout time.Duration) error
```
WaitForCronJob waits until job deployment has completed

#### func (*Framework) WaitForDaemonSet

```go
func (f *Framework) WaitForDaemonSet(name, namespace string, timeout time.Duration) error
```
WaitForDaemonSet waits until the specified daemonset is created or the timeout
is reached

#### func (*Framework) WaitForDaemonSetPods

```go
func (f *Framework) WaitForDaemonSetPods(name, namespace string, timeout time.Duration) error
```
WaitForDaemonSetPods waits until the specified deployment's pods are created or
the timeout is reached

#### func (*Framework) WaitForDeployment

```go
func (f *Framework) WaitForDeployment(name, namespace string, timeout time.Duration) error
```
WaitForDeployment waits until the specified deployment is created or the timeout
is reached

#### func (*Framework) WaitForDeploymentPods

```go
func (f *Framework) WaitForDeploymentPods(name, namespace string, timeout time.Duration) error
```
WaitForDeploymentPods waits until the specified deployment's pods are created or
the timeout is reached

#### func (*Framework) WaitForImagePolicy

```go
func (f *Framework) WaitForImagePolicy(name, namespace string, timeout time.Duration) error
```
WaitForImagePolicy waits until the ImagePolicy is created or the timeout is
reached

#### func (*Framework) WaitForImagePolicyDefinition

```go
func (f *Framework) WaitForImagePolicyDefinition(timeout time.Duration) error
```
WaitForImagePolicyDefinition waits until the ImagePolicy CRD is created or the
timeout is reached

#### func (*Framework) WaitForJob

```go
func (f *Framework) WaitForJob(name, namespace string, timeout time.Duration) error
```
WaitForJob waits until job deployment has completed

#### func (*Framework) WaitForMutatingAdmissionWebhook

```go
func (f *Framework) WaitForMutatingAdmissionWebhook(name string, timeout time.Duration) error
```
WaitForMutatingAdmissionWebhook waits until the specified
MutatingAdmissionWebhook is created or the timeout is reached

#### func (*Framework) WaitForNamespace

```go
func (f *Framework) WaitForNamespace(name string, timeout time.Duration) error
```
WaitForNamespace waits until the specified namespace is created or the timeout
is reached

#### func (*Framework) WaitForPod

```go
func (f *Framework) WaitForPod(name, namespace string, timeout time.Duration) error
```
WaitForPod waits until pod deployment has completed

#### func (*Framework) WaitForPodDelete

```go
func (f *Framework) WaitForPodDelete(name, namespace string, timeout time.Duration) error
```
WaitForPodDelete waits until pod has been deleted

#### func (*Framework) WaitForReplicaSet

```go
func (f *Framework) WaitForReplicaSet(name, namespace string, timeout time.Duration) error
```
WaitForReplicaSet waits until the specified Replicaset is created or the timeout
is reached

#### func (*Framework) WaitForReplicaSetPods

```go
func (f *Framework) WaitForReplicaSetPods(name, namespace string, timeout time.Duration) error
```
WaitForReplicaSetPods waits until the specified deployment's pods are created or
the timeout is reached

#### func (*Framework) WaitForReplicationController

```go
func (f *Framework) WaitForReplicationController(name, namespace string, timeout time.Duration) error
```
WaitForReplicationController waits until the specified Replicaset is created or
the timeout is reached

#### func (*Framework) WaitForReplicationControllerPods

```go
func (f *Framework) WaitForReplicationControllerPods(name, namespace string, timeout time.Duration) error
```
WaitForReplicationControllerPods waits until the specified deployment's pods are
created or the timeout is reached

#### func (*Framework) WaitForSecret

```go
func (f *Framework) WaitForSecret(name, namespace string, timeout time.Duration) error
```
WaitForSecret waits until the specified deployment is created or the timeout is
reached

#### func (*Framework) WaitForStatefulSet

```go
func (f *Framework) WaitForStatefulSet(name, namespace string, timeout time.Duration) error
```
WaitForStatefulSet waits until the specified StatefulSet is created or the
timeout is reached

#### func (*Framework) WaitForStatefulSetPods

```go
func (f *Framework) WaitForStatefulSetPods(name, namespace string, timeout time.Duration) error
```
WaitForStatefulSetPods waits until the specified deployment's pods are created
or the timeout is reached

#### func (*Framework) WaitForValidatingAdmissionWebhook

```go
func (f *Framework) WaitForValidatingAdmissionWebhook(name string, timeout time.Duration) error
```
WaitForValidatingAdmissionWebhook waits until the specified
ValidationAdmissionWebhook is created or the timeout is reached
