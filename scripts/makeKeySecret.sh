#!/bin/sh

# Write a Kubernetes secret containing public key(s) for Portieris to verify simple signatures

NAME=$1
if [ -z "${NAME}" ];then
echo usage: $0 '<secret-name> [<key-id>] [<key-id>]'
echo 'where <key-id> can be obtained from `gpg -k --key-id-format=long`'
fi
shift

cat <<EOF > ${NAME}.yaml
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: ${NAME}
stringData:
  key: |
EOF
gpg --export --armour $* | sed 's/^/    /' >> ${NAME}.yaml
