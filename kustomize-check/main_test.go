package main

import (
	"sync"
	"testing"
)

func TestRunKustomize(t *testing.T) {
	root := "/home/joe/repos/k8s-setup"
	_, err := runKustomize(root)
	if err != nil {
		t.Error(err)
	}
}

func TestWalk(t *testing.T) {
	dirs, err := walkDir("/home/joe/repos/k8s-setup")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dirs)
}

func TestKustomize(t *testing.T) {
	dirs := []string{
		// "./test/prevent-crd-deletion/staging/",
		// "./test/certwatcher/staging/",
	}
	for _, d := range dirs {
		_, err := kustomize(d)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestRunImages(t *testing.T) {
	images := imagesMutex{
		images: []string{"jonas27test/goserver:latest", "jonas27test/goserver@sha256:8ae229414f942ccfb7c531952bd4a929607d7a0db60746aefebc4ab02860620e", "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/webhook:v0.42.0@sha256:90989eeb6e0ba9c481b1faba3b01bcc70725baa58484c8f6ce9d22cc601e63dc", "ghcr.io/goharbor/harbor-operator:dev_master", "alpine", "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/resolvers:v0.42.0@sha256:eaa7d21d45f0bc1c411823d6a943e668c820f9cf52f1549d188edb89e992f6e0", "quay.io/minio/operator:v4.4.28", "gcr.io/k8s-staging-metrics-server/metrics-server:master", "gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/controller:v0.42.0@sha256:1fa50403c071b768984e23e26d0e68d2f7e470284ef2eb73581ec556bacdad95", "quay.io/jetstack/cert-manager-ctl:v1.10.0", "registry.opensource.zalan.do/acid/pgbouncer:master-18", "quay.io/spotahome/redis-operator:v1.1.1", "registry.k8s.io/ingress-nginx/kube-webhook-certgen:v20220916-gd32f8c343@sha256:39c5b2e3310dc4264d638ad28d9d1d96c4cbb2b2dcfb52368fe4e3c63f61e10f", "ghcr.io/isso-comments/isso:latest", "quay.io/jetstack/cert-manager-cainjector:v1.10.0", "ghcr.io/goharbor/postgres-operator:v1.7.0", "registry.k8s.io/ingress-nginx/controller:v1.5.1@sha256:4ba73c697770664c1e00e9f968de14e08f606ff961c76e5d7033a4a9c593c629", "quay.io/argoproj/argocd:v2.5.4", "registry.opensource.zalan.do/acid/spilo-13:2.1-p1", "jonas27/photos2022:0.1.2@sha256:65a554a12a6d1f5a802fa931205a2ea0d71bdff475ae24a555230ced0c6b69e2", "gcr.io/k8s-staging-gateway-api/admission-server:v0.5.1", "k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.1.1", "redis:7.0.5-alpine", "quay.io/jetstack/cert-manager-controller:v1.10.0", "jonas27/imagechecker:master-1671657213@sha256:12bc2c056cc4f59c976c17e4fde549e113272750bd415f4835e35bfed534c866", "nginx/nginx-ingress:2.4.1", "registry.opensource.zalan.do/acid/logical-backup:v1.7.0", "ghcr.io/dexidp/dex:v2.35.3", "quay.io/jetstack/cert-manager-webhook:v1.10.0"},
		mu:     sync.Mutex{},
	}
	err := runImages(&images)
	if err != nil {
		t.Error(err)
	}
}

func TestImage(t *testing.T) {
	tests := []struct {
		image string
		pass  bool
	}{
		{
			image: "jonas27test/goserver:latest",
			pass:  true,
		},
		{
			image: "quay.io/jetstack/cert-manager-cainjector:v1.7.2",
			pass:  true,
		},
		{
			image: "quay.io/jetstack/cert-manager-cainjector:v1.7.2@sha256:f82b3a5a153d9cabfc115e9ebb92b71851095d299ff8ab46f9677cae53557604",
			pass:  false,
		},
		{
			image: "test/testFalse",
			pass:  false,
		},
	}
	for _, tc := range tests {
		tc := tc // important because of loop var overwrite
		t.Run(tc.image, func(t *testing.T) {
			t.Parallel()
			err := checkImage(tc.image, nil)
			if tc.pass && err != nil {
				t.Error(err)
			} else if !tc.pass && err == nil {
				t.Errorf("test should no pass")
			}
		})
	}
}

// Usage: go test -bench ^BenchmarkImage$ -run=^$ -benchtime=10x
func BenchmarkImage(b *testing.B) {
	image := "jonas27test/goserver:latest"
	for i := 0; i < b.N; i++ {
		err := checkImage(image, nil)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestCustomAuth(t *testing.T) {
	img := "custom.docker.com/test/test:vtest@sha256:dtest"
	auth, err := customAuth()
	if err != nil {
		t.Fatal(err)
	}
	desc, err := remoteImage(img, auth)
	if err != nil {
		t.Log(err)
	}
	t.Log(desc)
}
