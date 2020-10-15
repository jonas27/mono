#!/bin/bash
REPO="ecc-haproxy/"
REPO_DIR=${PWD%ecc-haproxy*} && REPO_DIR="$REPO_DIR$REPO"
CERTS_DIR="certs" && CERTS_DIR=$REPO_DIR$CERTS_DIR
CONFIG_DIR="haconfig" && CONFIG_DIR=$REPO_DIR$CONFIG_DIR
cd $REPO_DIR
docker build -t haproxy .
docker stop haproxy || true && docker rm haproxy || true
echo "RUN"
echo $CONFIG_DIR
echo $CERTS_DIR
docker run --name haproxy -p 9090:9090 -p 9091:9091 -v $CERTS_DIR:/etc/ssl -v $CONFIG_DIR:/usr/local/etc/haproxy:ro haproxy
