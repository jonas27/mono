package main

import (
	"log"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
)

func TestGetImage(t *testing.T) {
	ref := "jonas27test/goserver:v1.0.3@sha256:8ae229414f942ccfb7c531952bd4a929607d7a0db60746aefebc4ab02860620e"
	image, _, err := splitDigest(ref)
	if err != nil {
		log.Fatalln(err)
	}

	correct := getImage(image)

	log.Fatal(correct)
}

func TestInternal(t *testing.T) {
	ref := "jonas27test/goserver:v1.0.3@sha256:8ae229414f942ccfb7c531952bd4a929607d7a0db60746aefebc4ab02860620e"
	digest, err := name.NewDigest(ref)
	if err != nil {
		log.Println(err)
	}
	log.Println(digest)
	tag, err := name.NewTag(ref)
	if err != nil {
		log.Println(err)
	}
	log.Fatal(tag)

}

func TestK8sRegistry(t *testing.T) {
	ref := "registry.k8s.io/ingress-nginx/kube-webhook-certgen:v20220916-gd32f8c343@sha256:39c5b2e3310dc4264d638ad28d9d1d96c4cbb2b2dcfb52368fe4e3c63f61e10f"
	digest := getImage(ref)
	if digest != 

}
