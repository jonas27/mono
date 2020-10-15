#!/bin/bash
REPO="ecc-haproxy/"
REPO_DIR=${PWD%ecc-haproxy*} && REPO_DIR="$REPO_DIR$REPO"
CERTS_DIR="certs" && CERTS_DIR=$REPO_DIR$CERTS_DIR
CONFIG_DIR="haconfig" && CONFIG_DIR=$REPO_DIR$CONFIG_DIR
cd $REPO_DIR
docker build -t haproxy .
docker stop haproxy || true && docker rm haproxy || true

echo $CERTS_DIR
echo $CONFIG_DIR

# For Windows  
# Convert forward to backward slash -e 's/\//\\/g' 
CERTS_DIR=$(echo "$CERTS_DIR" | sed -e 's/^\///' -e 's/^./\0:/')
CONFIG_DIR=$(echo "$CONFIG_DIR" | sed -e 's/^\///' -e 's/^./\0:/')

echo $CERTS_DIR
echo $CONFIG_DIR

echo "RUN"
docker run --name haproxy -p 9090:9090 -p 9091:9091 -v $CERTS_DIR:/etc/ssl/certs -v $CONFIG_DIR:/usr/local/etc/haproxy:ro haproxy

# check cert
# echo | openssl s_client -showcerts -servername localhost -connect localhost:9091 2>/dev/null | openssl x509 -inform pem -noout -text
# echo | openssl s_client -cipher 'AES256-SHA256' -showcerts -servername localhost -connect localhost:9091 2>/dev/null | openssl x509 -inform pem -noout -text

# Analyse server.pem
# openssl rsa -in certs/server.pem -text -noout