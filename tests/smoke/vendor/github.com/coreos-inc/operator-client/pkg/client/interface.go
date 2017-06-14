package client

import (
	"github.com/coreos-inc/operator-client/pkg/types"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
	v1beta1extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/pkg/runtime"
)

// Interface is the top level generic interface for the
// operator client.
type Interface interface {
	KubernetesInterface() kubernetes.Interface

	TPRClient
	TPRKindClient
	AppVersionClient
	MigrationStatusClient
	DaemonSetClient
	PodClient
	DeploymentClient
	ServiceClient
	NodeClient
}

// TPRClient holds methods for manipulating Third Party Resources.
type TPRClient interface {
	GetThirdPartyResource(apiGroup, version, namespace, resourceKind, resourceName string) (*runtime.Unstructured, error)
	GetThirdPartyResourceRaw(apiGroup, version, namespace, resourceKind, resourceName string) ([]byte, error)
	CreateThirdPartyResource(item *runtime.Unstructured) error
	CreateThirdPartyResourceRaw(apiGroup, version, namespace, kind string, data []byte) error
	CreateThirdPartyResourceRawIfNotFound(apiGroup, version, namespace, kind, name string, data []byte) (bool, error)
	UpdateThirdPartyResource(item *runtime.Unstructured) error
	UpdateThirdPartyResourceRaw(apiGroup, version, namespace, resourceKind, resourceName string, data []byte) error
	CreateOrUpdateThirdpartyResourceRaw(apiGroup, version, namespace, resourceKind, resourceName string, data []byte) error
	DeleteThirdPartyResource(apiGroup, version, namespace, resourceKind, resourceName string) error
	AtomicModifyThirdPartyResource(apiGroup, version, namespace, resourceKind, resourceName string, f TPRModifier, data interface{}) error
	SetFailureStatus(name string, failureStatus *types.FailureStatus) error
	SetTaskStatuses(name string, ts []types.TaskStatus) error
	UpdateTaskStatus(name string, ts types.TaskStatus) error
}

// TPRKindClient contains methods for manipulating and creating TPR Kinds.
type TPRKindClient interface {
	GetThirdPartyResourceKind(name string) (*v1beta1extensions.ThirdPartyResource, error)
	CreateThirdPartyResourceKind(tpr *v1beta1extensions.ThirdPartyResource) error
	DeleteThirdPartyResourceKind(name string) error
	EnsureThirdPartyResourceKind(tpr *v1beta1extensions.ThirdPartyResource) error
}

// AppVersionClient contains methods for the AppVersion resource.
type AppVersionClient interface {
	AtomicUpdateAppVersion(name string, fn types.AppVersionModifier) (*types.AppVersion, error)
	UpdateAppVersion(*types.AppVersion) (*types.AppVersion, error)
	GetAppVersion(name string) (*types.AppVersion, error)
}

// MigrationStatusClient contains methods for the MigrationStatus resource.
type MigrationStatusClient interface {
	GetMigrationStatus(name string) (*types.MigrationStatus, error)
	CreateMigrationStatus(*types.MigrationStatus) (*types.MigrationStatus, error)
	UpdateMigrationStatus(*types.MigrationStatus) (*types.MigrationStatus, error)
}

// DaemonSetClient contains methods for the DaemonSet resource.
type DaemonSetClient interface {
	CreateDaemonSet(*v1beta1extensions.DaemonSet) (*v1beta1extensions.DaemonSet, error)
	GetDaemonSet(namespace, name string) (*v1beta1extensions.DaemonSet, error)
	UpdateDaemonSetObject(*v1beta1extensions.DaemonSet) (*v1beta1extensions.DaemonSet, error)
	DaemonSetRollingUpdate(*v1beta1extensions.DaemonSet, UpdateOpts) (bool, error)
	AvailablePodsForDaemonSet(*v1beta1extensions.DaemonSet) (int, error)
	UnavailablePodsForDaemonSet(*v1beta1extensions.DaemonSet) (int, error)
	NumberOfDesiredPodsForDaemonSet(*v1beta1extensions.DaemonSet) (int, error)
}

// PodClient contains methods for the Pod resource.
type PodClient interface {
	DeletePod(namespace, name string) error
	ListPodsWithSelector(namespace string, selector *unversioned.LabelSelector) (*v1.PodList, error)
}

// DeploymentClient contains methods for the Deployment resource.
type DeploymentClient interface {
	GetDeployment(namespace, name string) (*v1beta1extensions.Deployment, error)
	CreateDeployment(*v1beta1extensions.Deployment) (*v1beta1extensions.Deployment, error)
	UpdateDeployment(*v1beta1extensions.Deployment) (*v1beta1extensions.Deployment, error)
	DeploymentRollingUpdate(*v1beta1extensions.Deployment, UpdateOpts) (bool, error)
}

// ServiceClient contains methods for the Service resource.
type ServiceClient interface {
	GetService(namespace, name string) (*v1.Service, error)
	UpdateService(*v1.Service) (*v1.Service, error)
}

// NodeClient contains methods for the Node resource.
type NodeClient interface {
	GetNode(name string) (*v1.Node, error)
	UpdateNode(*v1.Node) (*v1.Node, error)
	DrainNode(*v1.Node) error
	UnCordonNode(*v1.Node) (*v1.Node, error)
	CordonNode(*v1.Node) (*v1.Node, error)
}
