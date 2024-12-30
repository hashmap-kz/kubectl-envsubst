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
	"path/filepath"
)

// KubeClient is a dynamic type of client for kubernetes
type KubeClient struct {
	client          *dynamic.DynamicClient
	discoveryMapper *restmapper.DeferredDiscoveryRESTMapper
}

func GetClient() (*KubeClient, error) {

	kubeConfigFile, err := getKubeconfigPath()
	if err != nil {
		return nil, err
	}

	client, err := newKubeClient(kubeConfigFile)
	return client, nil
}

// NewKubeClient creates an instance of KubeClient
func newKubeClient(kubeConfigFile string) (*KubeClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile)
	if err != nil {
		return nil, err
	}

	// create the dynamic client
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// create a discovery client to map dynamic API resources
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	discoveryClient := memory.NewMemCacheClient(clientset.Discovery())
	discoveryMapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)

	return &KubeClient{client: client, discoveryMapper: discoveryMapper}, nil
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

		log.Printf("%s %q", obj.GetKind(), obj.GetName())
		return nil
	})

}

// getKubeconfigPath determines the default location of kube config (~/.kube/config)
func getKubeconfigPath() (string, error) {
	if kubeconfigEnvPath := os.Getenv("KUBECONFIG"); kubeconfigEnvPath != "" {
		return kubeconfigEnvPath, nil
	}

	if home := UserHomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config"), nil
	}

	return "", fmt.Errorf("could not determine kubeconfig directory")
}

// UserHomeDir returns the current use home directory
func UserHomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
