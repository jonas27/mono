package main

import (
	"log"
	"testing"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
)

func TestGetImage(t *testing.T) {
	image := "jonas27test/goserver:v1.0.3@sha256:8ae229414f942ccfb7c531952bd4a929607d7a0db60746aefebc4ab02860620e"
	desc, err := getImage(image)
	if err != nil {
		log.Fatalln(err)
	}
	digest := "sha256:8ae229414f942ccfb7c531952bd4a929607d7a0db60746aefebc4ab02860620e"
	if digest != desc.Digest.String() {
		log.Fatal(desc.Digest.String())
	}
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

// Image can be later loaded with docker load -i <dir>
func TestSaveOCI(t *testing.T) {
	img := "ubuntu/zookeeper"
	desc, err := getImage(img)
	if err != nil {
		log.Fatalln(err)
	}
	image, err := desc.Image()
	if err != nil {
		log.Fatalln(err)
	}
	dst := "./test/test.tar"
	if err := crane.SaveOCI(image, dst); err != nil {
		log.Fatalf("pulling %s: %s", img, err.Error())
	}
}
