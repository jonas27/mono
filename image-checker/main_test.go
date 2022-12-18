package main

import (
	"log"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
)

func TestGetImage(t *testing.T) {
	ref := "jonas27test/goserver:v1.0.3@sha256:8ae229414f942ccfb7c531952bd4a929607d7a0db60746aefebc4ab02860620e"
	image, digest, err := splitDigest(ref)
	if err != nil {
		log.Fatalln(err)
	}

	correct := getImage(image, digest)

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
