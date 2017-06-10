#!/bin/bash
#
# This script generates a number of certificates for testing olb with client
# certificate authentication with signed client certificates.
#
# First, a self-signed CA certificate is created which is used to sign both the
# server and client certificates. Then a server and a client certificate are
# created. The demo/cert/{ca,client,server} directories contain the generated
# certificates and their private keys.
#
# Second, a directory structure for a olb path cert store is created under
# demo/cert/olb/{client,server}. The server directory contains the TLS server
# certificate and private key. The client directory contains the client
# certificate **and** the CA certificate (no private keys). Including the CA
# certificate is necessary since otherwise olb (or the go crypto library)
# cannot verify the client certificate and will respond with the following
# error. Try this by removing the demo/cert/olb/client/ca-cert.pem file and
# restart olb.
#
# http: TLS handshake error from 127.0.0.1:53272: tls: failed to verify client's certificate: x509: certificate signed by unknown authority
#

set -o errexit
set -o nounset

basedir=$(cd $(dirname $0)/.. ; pwd)
certdir=$basedir/demo/cert
openssl=$(which openssl)
[[ -x /usr/local/opt/openssl/bin/openssl ]] && openssl=/usr/local/opt/openssl/bin/openssl

# shorten certdir
certdir=${certdir/$(pwd)\//}

[[ -z "$certdir" ]] && (echo "certdir empty" ; exit 1)
[[ -d "$certdir" ]] && rm -rf "$certdir"
mkdir -p $certdir/{ca,client,server} $certdir/olb/{client,server}

echo "generate CA cert"
$openssl req \
	-x509 -nodes -days 365 -sha256 -newkey rsa:2048 \
	-subj '/C=NL/ST=Noord-Holland/L=Amsterdam/CN=ca' \
	-keyout "$certdir/ca/ca-key.pem" -out "$certdir/ca/ca-cert.pem"

echo "generate client cert"
$openssl req \
	-nodes -days 365 -sha256 -newkey rsa:2048 \
	-subj '/C=NL/ST=Noord-Holland/L=Amsterdam/CN=client' \
	-keyout $certdir/client/client-key.pem -out $certdir/client/client.csr

$openssl x509 \
	-req -set_serial 02 -CA $certdir/ca/ca-cert.pem -CAkey $certdir/ca/ca-key.pem \
	-in $certdir/client/client.csr -out $certdir/client/client-cert.pem

echo "generate server cert"
$openssl req \
	-nodes -days 365 -sha256 -newkey rsa:2048 \
	-subj '/C=NL/ST=Noord-Holland/L=Amsterdam/CN=www.server.com' \
	-keyout $certdir/server/server-key.pem -out $certdir/server/server.csr

$openssl x509 \
	-req -set_serial 03 -CA $certdir/ca/ca-cert.pem -CAkey $certdir/ca/ca-key.pem \
	-in $certdir/server/server.csr -out $certdir/server/server-cert.pem

cp $certdir/ca/ca-cert.pem $certdir/olb/client
cp $certdir/client/client-cert.pem $certdir/olb/client
cp $certdir/server/server-{cert,key}.pem $certdir/olb/server

cat<<EOF

# start olb with path cert store
$basedir/olb \\
 -proxy.cs 'cs=db;type=path;refresh=3s;cert=$certdir/olb/server;clientca=$certdir/olb/client' \\
 -proxy.addr ':9999;cs=db'

# connect with openssl and client cert
$openssl s_client \\
 -tls1_2 -CAfile $certdir/ca/ca-cert.pem \\
 -servername www.server.com -connect localhost:9999 \\
 -cert $certdir/client/client-cert.pem -key $certdir/client/client-key.pem

EOF
