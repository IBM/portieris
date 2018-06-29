#! /bin/bash -e

NAMESPACE=${1:-ibm-system}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE_NAME=$(cat $SCRIPT_DIR/portieris/Chart.yaml | yq -r .name)

echo "Using $SERVICE_NAME as the service name"

CERT_DIR=$SCRIPT_DIR/certs

rm -rf $CERT_DIR
mkdir -p $CERT_DIR

cat > $CERT_DIR/server.conf << EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth, serverAuth
EOF

# Create a certificate authority
openssl genrsa -out $CERT_DIR/caKey.pem 2048
openssl req -x509 -new -nodes -key $CERT_DIR/caKey.pem -days 100000 -out $CERT_DIR/caCert.pem -subj "/CN=${SERVICE_NAME}_ca"

# Create a server certiticate
openssl genrsa -out $CERT_DIR/serverKey.pem 2048
# Note the CN is the DNS name of the service of the webhook.
openssl req -new -key $CERT_DIR/serverKey.pem -out $CERT_DIR/server.csr -subj "/CN=${SERVICE_NAME}.${NAMESPACE}.svc" -config $CERT_DIR/server.conf
openssl x509 -req -in $CERT_DIR/server.csr -CA $CERT_DIR/caCert.pem -CAkey $CERT_DIR/caKey.pem -CAcreateserial -out $CERT_DIR/serverCert.pem -days 100000 -extensions v3_req -extfile $CERT_DIR/server.conf

outfile=$CERT_DIR/certs.go

cat > $outfile << EOF
// Copyright 2018 Portieris Authors.
//
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
EOF

echo "" >> $outfile
echo "// This file was generated using openssl by the gengerts.sh script" >> $outfile
echo "" >> $outfile
echo "package main" >> $outfile

for file in caKey caCert serverKey serverCert; do
	data=$(cat $CERT_DIR/${file}.pem)
	echo "" >> $outfile
	echo "var $file = []byte(\`$data\`)" >> $outfile
done

cp $CERT_DIR/certs.go $GOPATH/src/github.com/IBM/portieris/cmd/trust/certs.go

echo "caBundle: \"$(cat $CERT_DIR/caCert.pem | base64)\""