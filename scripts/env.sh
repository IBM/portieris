# probably set these to the same user APIKEY
export PORTIERIS_PULL_APIKEY=
export PORTIERIS_TESTIMAGE_APIKEY=
# charts to test
export VERSION=
# image tag to test e.g. prep-v0.13.19
export TAG=

# name of the secret used to pull portieris made from $REG and $PORTIERIS_PULL_APIKEY
export PULLSECRET=portieris-test
export REG=icr.io
export HUB=${REG}/registry-deploy

# points to kube tests cluster (docker)
export KUBECONFIG=~/.kube/config
