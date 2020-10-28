#!/usr/bin/env bash
# This tests both ciphers with sigalgs https://github.com/openssl/openssl/issues/10131

# check default is ECDSA
echo | openssl s_client -showcerts -servername localhost -connect localhost:9091 2>/dev/null | openssl x509 -inform pem -noout -text | grep "Public Key Algorithm"
# Use RSA
echo | openssl s_client -sigalgs RSA-PSS+SHA256 -showcerts -servername localhost -connect localhost:9091 2>/dev/null | openssl x509 -inform pem -noout -text | grep "Public Key Algorithm"