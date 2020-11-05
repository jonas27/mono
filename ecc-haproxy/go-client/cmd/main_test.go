package main

import (
	"log"
	"testing"
)

func BenchmarkECC(b *testing.B) {

	log.Println("test Ecc")
	for i := 0; i < b.N; i++ {
		MakeRequest(false)
	}
}
func BenchmarkRSA(b *testing.B) {
	// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA}, MaxVersion: tls.VersionTLS12}
	log.Println("test RSA")
	for n := 0; n < b.N; n++ {
		MakeRequest(true)
	}
}
