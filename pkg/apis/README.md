---

copyright:
  years: 2021
lastupdated: "2021-05-12"

---

# Image Policy CRD (CustomResourceDefinition)

## Regenerate the CRD files with code-generator

**Important:** You only have to regenerate the customer resource definition (CRD) files if you're a developer and you're changing the design of the CRD policies.

Some scripts from the [k8s.io/code-generator](https://github.com/kubernetes/code-generator) repository are used to generate the clientset, informers, listers, and deep-copy functions.

### Prerequisites

Clone the code-generator repository that is compatible with type of the API that you're using.

```bash
make code-generator
```

#### Generate the CRDs

```bash
make regenerate
```
