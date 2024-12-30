package util

// https://simplerize.com/kubernetes/apply-yaml-manifests-to-kubernetes-using-go-client

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// KubeClient is a dynamic type of client for kubernetes
type KubeClient struct {
	client          *dynamic.DynamicClient
	discoveryMapper *restmapper.DeferredDiscoveryRESTMapper
}

func GetClient() *KubeClient {

	// create the dynamic client
	kubeConfigFile := defaultKubeConfigFile()
	log.Printf("creating dynamic client with config: %s\n", kubeConfigFile)
	client := newKubeClient(kubeConfigFile)

	return client
}

// NewKubeClient creates an instance of KubeClient
func newKubeClient(kubeConfigFile string) *KubeClient {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	// create the dynamic client
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create dynamic client: %w", err)
	}

	// create a discovery client to map dynamic API resources
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create discovery client: %w", err)
	}

	discoveryClient := memory.NewMemCacheClient(clientset.Discovery())
	discoveryMapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	return &KubeClient{client: client, discoveryMapper: discoveryMapper}
}

// Apply applies the given YAML manifests to kubernetes
func (k *KubeClient) Apply(obj *unstructured.Unstructured) error {

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {

		// get GroupVersionResource to invoke the dynamic client
		gvk := obj.GroupVersionKind()
		restMapping, err := k.discoveryMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}
		gvr := restMapping.Resource

		// apply the YAML doc
		namespace := obj.GetNamespace()
		if len(namespace) == 0 {
			namespace = "default"
		}

		applyOpts := metav1.ApplyOptions{FieldManager: "kube-apply"}
		_, err = k.client.Resource(gvr).Namespace(namespace).Apply(context.TODO(), obj.GetName(), obj, applyOpts)

		if err != nil {
			return fmt.Errorf("apply error: %w", err)
		}

		log.Printf("applied YAML for %s %q", obj.GetKind(), obj.GetName())
		return nil
	})

}

// defaultKubeConfigFile determines the default location of kube config (~/.kube/config)
func defaultKubeConfigFile() string {
	e := os.Getenv("KUBECONFIG")
	if e != "" {
		return e
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatalf("unable to find the kube config: %s", err)
	}
	fromFile := filepath.Join(usr.HomeDir, ".kube", "config")
	return fromFile
}
