#! /bin/bash -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOYMENT=$1
SERVICE_NAME=$(cat $SCRIPT_DIR/ibmcloud-image-enforcement/Chart.yaml | yq -r .name)-$DEPLOYMENT

echo "Using $SERVICE_NAME as the service name"

CERT_DIR=$SCRIPT_DIR/certs

rm -fr $CERT_DIR
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
openssl req -new -key $CERT_DIR/serverKey.pem -out $CERT_DIR/server.csr -subj "/CN=${SERVICE_NAME}.ibm-system.svc" -config $CERT_DIR/server.conf
openssl x509 -req -in $CERT_DIR/server.csr -CA $CERT_DIR/caCert.pem -CAkey $CERT_DIR/caKey.pem -CAcreateserial -out $CERT_DIR/serverCert.pem -days 100000 -extensions v3_req -extfile $CERT_DIR/server.conf

outfile=$CERT_DIR/certs.go

cat > $outfile << EOF
/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Registry, 5737-D42
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/
EOF

echo "" >> $outfile
echo "// This file was generated using openssl by the genCerts.sh script" >> $outfile
echo "" >> $outfile
echo "package main" >> $outfile

for file in caKey caCert serverKey serverCert; do
	data=$(cat $CERT_DIR/${file}.pem)
	echo "" >> $outfile
	echo "var $file = []byte(\`$data\`)" >> $outfile
done

cp $CERT_DIR/certs.go $GOPATH/src/github.com/IBM/portieris/cmd/$DEPLOYMENT/certs.go

echo "Printing out caBundle for $DEPLOYMENT..."
echo -e "caBundle: \"$(cat $CERT_DIR/caCert.pem | base64)\"\n"