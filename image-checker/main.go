// https://www.cncf.io/blog/2019/10/15/extend-kubernetes-via-a-shared-informer/

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// InvalidMessage will be return to the user.
	InvalidMessage  = "image tag AND digest are either not set or do not match"
	port            = ":8080"
	defaultRegistry = "registry.hub.docker.com"
)

type store struct {
	pods map[string]v1.Pod
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set the kubernetes config file path as environment variable
	// kubeconfig := os.Getenv("KUBECONFIG")
	kubeconfig := "/home/joe/.kube/config"

	// Create the client configuration
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		// logger.Panic(err.Error())
		os.Exit(1)
	}

	// Create the client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		// logger.Panic(err.Error())
		os.Exit(1)
	}

	// Create the shared informer factory and use the client to connect to
	// Kubernetes
	factory := informers.NewSharedInformerFactory(clientset, 0)

	// Get the informer for the right resource, in this case a Pod
	informer := factory.Core().V1().Pods().Informer()

	// Create a channel to stops the shared informer gracefully
	stopper := make(chan struct{})
	defer close(stopper)

	// Kubernetes serves an utility to handle API crashes
	defer runtime.HandleCrash()

	stor := store{
		pods: make(map[string]v1.Pod),
	}

	// This is the part where your custom code gets triggered based on the
	// event that the shared informer catches
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		// When a new pod gets created
		AddFunc: func(pod interface{}) { stor.add(*pod.(*corev1.Pod)) },
		// When a pod gets deleted
		DeleteFunc: func(pod interface{}) { stor.remove(*pod.(*corev1.Pod)) },
	})

	// You need to start the informer, in my case, it runs in the background
	informer.Run(stopper)
}

func (s *store) add(pod corev1.Pod) {
	s.pods[key(pod)] = pod
	for _, v := range pod.Spec.Containers {
		// log.Println(v.Image)
		validate(v.Image)
	}
}
func (s *store) remove(pod corev1.Pod) {
	delete(s.pods, key(pod))
}

func key(pod corev1.Pod) string {
	// log.Println(pod.Name)
	return fmt.Sprintf("%s:%s", pod.ObjectMeta.Name, pod.ObjectMeta.Namespace)
}

func validate(img string) bool {
	log.Println(img)
	refWithTag, digest, err := splitDigest(img)
	if err != nil {
		log.Fatal("no digest")
	}
	desc := getImage(refWithTag)
	remoteDigest := desc.Digest.Algorithm + ":" + desc.Digest.Hex
	if remoteDigest == digest {
		return true
	}
	return false
}

func splitDigest(id string) (string, string, error) {
	if strings.Count(id, ":") == 2 {
		return strings.Split(id, "@")[0], strings.Split(id, "@")[1], nil
	}
	return "", "", errors.New("image ID is missing either tag or digest")
}

func getImage(ref string) *remote.Descriptor {
	log.Println(ref)
	nref, err := name.ParseReference(ref)
	if err != nil {
		log.Fatalln(err)
	}
	desc, err := remote.Get(nref)
	if err != nil {
		log.Fatal(err)
	}
	return desc
}
