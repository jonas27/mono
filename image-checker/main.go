package main

import (
	"errors"
	"log"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

var (
	defaultRegistry = "registry.hub.docker.com"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

func splitDigest(id string) (string, string, error) {
	if strings.Count(id, ":") == 2 {
		return strings.Split(id, "@")[0], strings.Split(id, "@")[1], nil
	}
	return "", "", errors.New("image ID is missing either tag or digest")
}

func getImage(ref string, originalDigest string) bool {
	nref, err := name.ParseReference(ref)
	if err != nil {
		log.Fatalln(err)
	}
	desc, err := remote.Get(nref)
	if err != nil {
		log.Fatal(err)
	}
	if desc.Digest.Algorithm+":"+desc.Digest.Hex == originalDigest {
		return true
	}
	return false
}
