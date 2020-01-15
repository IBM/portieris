#! /bin/bash

RELEASE_NAME=${1:-portieris}
NAMESPACE=${2:-ibm-system}

echo Purging release ${RELEASE_NAME} in ${NAMESPACE}

kubectl delete MutatingWebhookConfiguration image-admission-config --ignore-not-found=true
kubectl delete ValidatingWebhookConfiguration image-admission-config --ignore-not-found=true

kubectl delete crd clusterimagepolicies.securityenforcement.admission.cloud.ibm.com imagepolicies.securityenforcement.admission.cloud.ibm.com --ignore-not-found=true

kubectl delete jobs -n ${NAMESPACE} create-admission-webhooks create-armada-image-policies create-crds validate-crd-creation --ignore-not-found=true

helm delete --purge ${RELEASE_NAME}

kubectl delete jobs -n ${NAMESPACE} delete-admission-webhooks delete-crds --ignore-not-found=true
