#!/bin/bash

# List curves
# openssl ecparam -list_curves

# https://stackoverflow.com/questions/15686821/generate-ec-keypair-from-openssl-command-line

# https://www.scottbrady91.com/OpenSSL/Creating-Elliptical-Curve-Keys-using-OpenSSL

# https://www.guyrutenberg.com/2013/12/28/creating-self-signed-ecdsa-ssl-certificate-using-openssl/

# openssl ecparam -name secp256r1 -genkey -param_enc explicit -out private-key.pem
# openssl req -new -x509 -key private-key.pem -out server.pem -days 730

# examine 
# openssl ecparam -in private-key.pem -text -noout
# openssl x509 -in server.pem -text -noout

# cat private-key.pem server.pem > server-private.pem

#  https://gist.github.com/marta-krzyk-dev/83168c9a8e985e5b3b1b14a98b533b9c
#  the only one which is working for browsers
openssl ecparam -name secp256r1 -genkey -param_enc explicit -out private-key.pem
openssl	ecparam -genkey -name secp256r1 -noout -out private-key.pem

# -subj '/C=DE/ST=IT/L=IT/O=IT/OU=IT/CN=localhost' not working on git bash windows
openssl req -new -x509 -key private-key.pem -out certificate.pem -days 900000 -subj '//C=DE\ST=IT\L=IT\O=IT\OU=IT\CN=localhost'
cat certificate.pem private-key.pem > server.pem