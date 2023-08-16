package softwaredefinedstorage

import (
	"time"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
)

// Cluster configs
type ClusterConfig struct {
	RestConfig *rest.Config
	ClientSet  *kubernetes.Clientset
}

var restConfig *rest.Config
var clientSet *kubernetes.Clientset

// sds timeout
var sdsTimeout time.Duration

// Common Interface for different Software Defined Solutions
type Sds interface {
	PreWorkerReplace(worker v2.Worker) error
	PostWorkerReplace(worker v2.Worker) error
}

// Set global variables frequently used in Pre/Post Worker replace actions
func SetGlobals(config *ClusterConfig, timeout time.Duration) {
	restConfig = config.RestConfig
	clientSet = config.ClientSet
	sdsTimeout = timeout
}
