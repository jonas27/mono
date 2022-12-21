// https://www.cncf.io/blog/2019/10/15/extend-kubernetes-via-a-shared-informer/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	port            = ":8080"
	defaultRegistry = "registry.hub.docker.com"
)

type store struct {
	unavailable      map[string]string
	withDigestAndTag map[string]string
	withDigestOrTag  map[string]string
	wrongDigest      map[string]string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	version := flag.Bool("v", false, "Prints the version of the app.")
	flag.Parse()
	if *version {
		fmt.Printf("Version %s\n", os.Getenv("VERSION"))
		os.Exit(0)
	}

	kubeconfig := "/home/joe/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	stor := store{
		unavailable:      make(map[string]string),
		withDigestAndTag: make(map[string]string),
		withDigestOrTag:  make(map[string]string),
		wrongDigest:      make(map[string]string),
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(pod interface{}) { stor.add(*pod.(*corev1.Pod)) },
		// Omit update, as pod images are immutable.
		DeleteFunc: func(pod interface{}) { stor.remove(*pod.(*corev1.Pod)) },
	})

	go informer.Run(stopper)
	serveMetrics()
}

func (s *store) add(p corev1.Pod) {
	for _, v := range p.Spec.Containers {
		img := v.Image
		_, err := getImage(img)
		key := fmt.Sprintf("%s:%s:%s", p.Namespace, p.Name, img)
		if err != nil {
			s.unavailable[key] = img
			continue
		}
		if !strings.Contains(img, "@") || strings.Count(img, ":") == 1 {
			s.withDigestOrTag[key] = img
			continue
		}
		desc, err := getImage(strings.Split(img, "@")[0])
		if err != nil {
			s.unavailable[key] = img
			continue
		}
		digest := strings.Split(img, "@")[1]
		remoteDigest := desc.Digest.String()
		if remoteDigest == digest {
			s.withDigestAndTag[key] = img
		} else {
			s.wrongDigest[key] = img
		}
	}
	s.updateMetrics()
}
func (s *store) remove(p corev1.Pod) {
	for _, v := range p.Spec.Containers {
		img := v.Image
		key := fmt.Sprintf("%s:%s:%s", p.Namespace, p.Name, img)
		delete(s.unavailable, key)
	}
	s.updateMetrics()
}

func getImage(img string) (*remote.Descriptor, error) {
	ref, err := name.ParseReference(img)
	if err != nil {
		return nil, err
	}
	desc, err := remote.Get(ref)
	if err != nil {
		return nil, err
	}
	return desc, nil
}

func (s *store) updateMetrics() {
	unavailable.Set(float64(len(s.unavailable)))
	withDigesAndtTag.Set(float64(len(s.withDigestAndTag)))
	withDigestOrTag.Set(float64(len(s.withDigestOrTag)))
	wrongDigest.Set(float64(len(s.wrongDigest)))
}
