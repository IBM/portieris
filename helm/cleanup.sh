#!/usr/bin/env bash

RELEASE_NAME=${1:-portieris}
NAMESPACE=${2:-portieris}

echo Deleting release "${RELEASE_NAME}" in "${NAMESPACE}"

kubectl delete MutatingWebhookConfiguration image-admission-config --ignore-not-found=true
kubectl delete ValidatingWebhookConfiguration image-admission-config --ignore-not-found=true

kubectl delete crd clusterimagepolicies.securityenforcement.admission.cloud.ibm.com imagepolicies.securityenforcement.admission.cloud.ibm.com --ignore-not-found=true

helm delete "${RELEASE_NAME}" --no-hooks --namespace "${NAMESPACE}"
