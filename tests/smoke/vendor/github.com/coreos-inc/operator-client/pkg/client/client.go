package client

import (
	"github.com/golang/glog"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Interface assertion.
var _ Interface = &Client{}

// Client is a kubernetes client that can talk to the API server.
type Client struct {
	config *rest.Config
	*kubernetes.Clientset
}

// NewClient creates a kubernetes client or bails out on on failures.
func NewClient(kubeconfig string) Interface {
	var config *rest.Config
	var err error

	if kubeconfig != "" {
		glog.V(4).Infof("Loading kube client config from path %q", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		glog.V(4).Infof("Using in-cluster kube client config")
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		glog.Fatalf("Cannot load config for REST client: %v", err)
	}

	return &Client{config, kubernetes.NewForConfigOrDie(config)}
}

// KubernetesInterface returns the Kubernetes interface.
func (c *Client) KubernetesInterface() kubernetes.Interface {
	return c.Clientset
}
