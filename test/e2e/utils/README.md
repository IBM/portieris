# utils
--
    import "github.com/IBM/portieris/test/e2e/utils"


## Usage

#### func  CheckIfTesting

```go
func CheckIfTesting(t *testing.T, boolToCheck bool)
```

#### func  CleanUpClusterImagePolicyTest

```go
func CleanUpClusterImagePolicyTest(t *testing.T, fw *framework.Framework, clusterPolicy, namespace string)
```

#### func  CleanUpImagePolicyTest

```go
func CleanUpImagePolicyTest(t *testing.T, fw *framework.Framework, namespace string)
```

#### func  CreateClusterImagePolicyAndNamespace

```go
func CreateClusterImagePolicyAndNamespace(t *testing.T, fw *framework.Framework, manifestPath string) (*v1beta1.ClusterImagePolicy, *corev1.Namespace)
```

#### func  CreateImagePolicyInstalledNamespace

```go
func CreateImagePolicyInstalledNamespace(t *testing.T, fw *framework.Framework, manifestPath string) *corev1.Namespace
```

#### func  CreateSecret

```go
func CreateSecret(t *testing.T, fw *framework.Framework, manifestPath, namespace string)
```

#### func  DeleteThenReturnClusterImagePolicy

```go
func DeleteThenReturnClusterImagePolicy(t *testing.T, fw *framework.Framework, clusterImagePolicy string) *v1beta1.ClusterImagePolicy
```
DeleteThenReturnClusterImagePolicy is used for temporary deletion of a cluster
image policy for a given test The returned ClusterImagePolicy should be used to
recreate after the test is complete using a defer

#### func  DumpEvents

```go
func DumpEvents(t *testing.T, fw *framework.Framework, namespace string)
```

#### func  DumpPolicies

```go
func DumpPolicies(t *testing.T, fw *framework.Framework, namespace string)
```

#### func  KillPod

```go
func KillPod(t *testing.T, fw *framework.Framework, namespace string)
```
KillPod kills first pod return in podlist in the given namespace

#### func  TestCronJobNotRunnable

```go
func TestCronJobNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestCronJobNotRunnable tests whether a manifest is deployable to the specified
namespace

#### func  TestCronJobRunnable

```go
func TestCronJobRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestCronJobRunnable tests whether a manifest is deployable to the specified
namespace

#### func  TestCurrentDeployStatus

```go
func TestCurrentDeployStatus(t *testing.T, fw *framework.Framework, namespace, deploymentName string)
```
TestCurrentDeployStatus checks the deployment currently has the expected number
of replicas

#### func  TestDaemonSetNotRunnable

```go
func TestDaemonSetNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestDaemonSetNotRunnable tests whether a manifest is deployable to the specified
namespace

#### func  TestDaemonSetRunnable

```go
func TestDaemonSetRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestDaemonSetRunnable tests whether a manifest is deployable to the specified
namespace

#### func  TestDeploymentNotRunnable

```go
func TestDeploymentNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestDeploymentNotRunnable tests whether a manifest is deployable to the
specified namespace

#### func  TestDeploymentNotRunnableOnPatch

```go
func TestDeploymentNotRunnableOnPatch(t *testing.T, fw *framework.Framework, deploymentName, patchString, namespace string)
```
TestDeploymentNotRunnableOnPatch tests whether a deplomyent is not runnable
after a patch

#### func  TestDeploymentNotRunnableOnReplace

```go
func TestDeploymentNotRunnableOnReplace(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestDeploymentNotRunnableOnReplace tests whether a deplomyent is not runnable
after a replace

#### func  TestDeploymentRunnable

```go
func TestDeploymentRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestDeploymentRunnable tests whether a manifest is deployable to the specified
namespace

#### func  TestJobNotRunnable

```go
func TestJobNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestJobNotRunnable tests whether a manifest is deployable to the specified
namespace

#### func  TestJobRunnable

```go
func TestJobRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestJobRunnable tests whether a manifest is deployable to the specified
namespace

#### func  TestPodNotRunnable

```go
func TestPodNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestPodNotRunnable tests whether a manifest is deployable to the specified
namespace

#### func  TestPodRunnable

```go
func TestPodRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestPodRunnable tests whether a manifest is deployable to the specified
namespace

#### func  TestReplicaSetNotRunnable

```go
func TestReplicaSetNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestReplicaSetNotRunnable tests whether a manifest is deployable to the
specified namespace

#### func  TestReplicaSetRunnable

```go
func TestReplicaSetRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestReplicaSetRunnable tests whether a manifest is deployable to the specified
namespace

#### func  TestReplicationControllerNotRunnable

```go
func TestReplicationControllerNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestReplicationControllerNotRunnable tests whether a manifest is deployable to
the specified namespace

#### func  TestReplicationControllerRunnable

```go
func TestReplicationControllerRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestReplicationControllerRunnable tests whether a manifest is deployable to the
specified namespace

#### func  TestStartDeployNoDelete

```go
func TestStartDeployNoDelete(t *testing.T, fw *framework.Framework, manifestLocation, namespace string) (deploymentName string)
```
TestStartDeployNoDelete starts a deployment and only deletes on failure

#### func  TestStatefulSetNotRunnable

```go
func TestStatefulSetNotRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestStatefulSetNotRunnable tests whether a manifest is deployable to the
specified namespace

#### func  TestStatefulSetRunnable

```go
func TestStatefulSetRunnable(t *testing.T, fw *framework.Framework, manifestLocation, namespace string)
```
TestStatefulSetRunnable tests whether a manifest is deployable to the specified
namespace

#### func  UpdateImagePolicy

```go
func UpdateImagePolicy(t *testing.T, fw *framework.Framework, manifestPath, namespace, oldPolicy string)
```
