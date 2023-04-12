package main

import (
	"testing"
)

func TestCloudSetup(t *testing.T) {
	run("/home/joe/repos/cloud-setup/infrastructure/azure-k8s")
	run("/home/joe/repos/cloud-setup/infrastructure/ck8s")
	run("/home/joe/repos/cloud-setup/infrastructure/k8s")
}

func TestWalk(t *testing.T) {
	dirs, err := walkDir("./test/")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dirs)
}

func TestKustomize(t *testing.T) {
	dirs := []string{
		// "./test/prevent-crd-deletion/staging/",
		"./test/certwatcher/staging/",
	}
	for _, d := range dirs {
		_, err := kustomize(d)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestImage(t *testing.T) {
	dirs := []string{
		// "./test/prevent-crd-deletion/staging/",
		"./test/certwatcher/staging/",
		// "./test/argocd-setup/staging/",
	}
	auth, err := customAuth()
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range dirs {
		yaml, err := kustomize(d)
		if err != nil {
			t.Error(err)
		}
		images := parseImages(yaml)
		err = checkImages(images, auth)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestQuay(t *testing.T) {
	img := []string{"quay.io/jetstack/cert-manager-cainjector:v1.7.2@sha256:f82b3a5a153d9cabfc115e9ebb92b71851095d299ff8ab46f9677cae53557604"}
	err := checkImages(img, nil)
	if err != nil {
		t.Error(err)
	}

}

func TestCustomAuth(t *testing.T) {
	img := "inx.dockreg.net/portal/certwatcher:v1.3.1@sha256:d16ce69a033cbcc90b2252f1fc656d0c3596d3ce8857963e5c63b3e6f9752d6b"
	auth, err := customAuth()
	if err != nil {
		t.Fatal(err)
	}
	desc, err := remoteImage(img, auth)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(desc)
}
