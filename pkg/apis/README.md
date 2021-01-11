# Image Policy CRD (CustomResourceDefinition)

## Regenerate the CRD files with code-generator

You **ONLY** need to do this if you are a developer **AND** you are changing the design of the CRD Policies, otherwise you **DO NOT** need to do this.

We use some scripts from the [k8s.io/code-generator](https://github.com/kubernetes/code-generator) repository to generate the clientset, informers, listers and deep-copy functions.

### Prerequisites

Clone the code-generator repo compatible with apis used

```bash
make code-generator
```

#### Generate the CRDs

```bash
make regenerate
```
