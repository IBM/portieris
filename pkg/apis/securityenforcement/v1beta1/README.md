## Image Policy CRD (CustomResourceDefinition)

### Regenerate the CRD files with code-generator
You **ONLY** need to do this if you are a developer **AND** you are changing the design of the CRD Policies, otherwise you **DO NOT** need to do this.

We use some scripts from the [k8s.io/code-generator](https://github.com/kubernetes/code-generator) repository to generate the clientset, informers, listers and deep-copy functions.

#### Prerequisites

Clone the code-generator repo compatible with apis used
```bash
git clone https://github.com/kubernetes/code-generator.git ./code-generator
(cd code-generator; git checkout v0.17.3)
```

##### Generate the CRDs
```bash
./code-generator/generate-groups.sh all github.com/IBM/portieris/pkg/apis/securityenforcement/client github.com/IBM/portieris/pkg/apis securityenforcement:v1beta1
```