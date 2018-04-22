#! /bin/bash

RELEASE_NAME=$1

kubectl delete MutatingWebhookConfiguration image-admission-config --ignore-not-found=true
kubectl delete ValidatingWebhookConfiguration image-admission-config --ignore-not-found=true

kubectl delete crd clusterimagepolicies.securityenforcement.admission.cloud.ibm.com imagepolicies.securityenforcement.admission.cloud.ibm.com --ignore-not-found=true

kubectl delete jobs -n ibm-system create-admission-webhooks create-armada-image-policies create-crds validate-crd-creation --ignore-not-found=true

helm delete --purge $RELEASE_NAME

kubectl delete jobs -n ibm-system delete-admission-webhooks delete-crds --ignore-not-found=true