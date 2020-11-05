package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
)

func main() {
	MakeRequest(true)
}

func MakeRequest(rsa bool) {
	if rsa {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA}, MaxVersion: tls.VersionTLS12}
	}
	resp, err := http.Get("https://jonasburster.de:9091/")
	if err != nil {
		log.Println(err)

	}
	if resp.TLS != nil {
		certificates := resp.TLS.PeerCertificates
		if len(certificates) > 0 {
			cert := certificates[0]
			c, _ := x509.ParseCertificate(cert.Raw)
			log.Println(c.PublicKeyAlgorithm)
		}
	}
}

// resp, clientErr := client.Do(req)
// if clientErr != nil {
//     panic(clientErr)
// }
// if resp.TLS != nil {
//     certificates := resp.TLS.PeerCertificates
//     if len(certificates) > 0 {
//         // you probably want certificates[0]
//         cert := certificates[0]
//     }
// }

// docker run \
//   -v /home/joe/repos/mono/ecc-haproxy/realCerts:/var/ssl \
//   -p 80:80 \
//   -e DOMAINS=jonasburster.de \
//   --rm \
//   asamoshkin/letsencrypt-certgen issue

// docker run  \
//   -v /home/joe/repos/mono/ecc-haproxy/realCerts/realCertsRsa:/acme.sh  \
//   --net=host \
//   neilpang/acme.sh  --issue -d jonasburster.de  --standalone --keylength ec-256
