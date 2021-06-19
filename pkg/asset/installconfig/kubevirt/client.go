package kubevirt

import (
	"context"
	"fmt"

	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/pkg/errors"
	authv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
	cdiv1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
)

//go:generate mockgen -source=./client.go -destination=./mock/client_generated.go -package=mock

// Client is a wrapper object for actual infra-cluster clients: kubernetes and the kubevirt
type Client interface {
	GetVirtualMachine(ctx context.Context, namespace string, name string) (*kubevirtapiv1.VirtualMachine, error)
	ListVirtualMachine(ctx context.Context, namespace string, opts metav1.ListOptions) (*kubevirtapiv1.VirtualMachineList, error)
	DeleteVirtualMachine(ctx context.Context, namespace string, name string) error
	GetDataVolume(ctx context.Context, namespace string, name string) (*cdiv1.DataVolume, error)
	ListDataVolume(ctx context.Context, namespace string, opts metav1.ListOptions) (*cdiv1.DataVolumeList, error)
	DeleteDataVolume(ctx context.Context, namespace string, name string) error
	GetSecret(ctx context.Context, namespace string, name string) (*corev1.Secret, error)
	ListSecret(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.SecretList, error)
	DeleteSecret(ctx context.Context, namespace string, name string) error
	GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error)
	GetStorageClass(ctx context.Context, name string) (*storagev1.StorageClass, error)
	GetNetworkAttachmentDefinition(ctx context.Context, name string, namespace string) (*unstructured.Unstructured, error)
	CreateSelfSubjectAccessReview(ctx context.Context, reviewObj *authv1.SelfSubjectAccessReview) (*authv1.SelfSubjectAccessReview, error)
	GetHyperConverged(ctx context.Context, name string, namespace string) (*unstructured.Unstructured, error)
}

type client struct {
	kubernetesClient *kubernetes.Clientset
	dynamicClient    dynamic.Interface
}

var (
	vmRes = schema.GroupVersionResource{
		Group:    kubevirtapiv1.GroupVersion.Group,
		Version:  kubevirtapiv1.GroupVersion.Version,
		Resource: "virtualmachines",
	}

	dvRes = schema.GroupVersionResource{
		Group:    cdiv1.SchemeGroupVersion.Group,
		Version:  cdiv1.SchemeGroupVersion.Version,
		Resource: "datavolumes",
	}
)

// LoadKubeConfigContent returns the kubeconfig file content
func LoadKubeConfigContent() ([]byte, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		return nil, err
	}

	// Remove anything that is not related to the current context from the result rawConfig
	currentContextValue := rawConfig.Contexts[rawConfig.CurrentContext]
	if currentContextValue == nil {
		return nil, fmt.Errorf("currentContext is not included in rawConfig.Contexts")
	}

	rawConfig.Contexts = map[string]*clientcmdapi.Context{
		rawConfig.CurrentContext: currentContextValue,
	}

	if v, ok := rawConfig.Clusters[currentContextValue.Cluster]; ok {
		rawConfig.Clusters = map[string]*clientcmdapi.Cluster{
			currentContextValue.Cluster: v,
		}
	}

	if v, ok := rawConfig.AuthInfos[currentContextValue.AuthInfo]; ok {
		rawConfig.AuthInfos = map[string]*clientcmdapi.AuthInfo{
			currentContextValue.AuthInfo: v,
		}
	}

	return clientcmd.Write(rawConfig)
}

// NewClient creates our client wrapper object for the actual kubeVirt and kubernetes clients we use.
func NewClient() (Client, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	restClientConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	kubernetesClient, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(restClientConfig)
	if err != nil {
		return nil, err
	}

	return &client{
		kubernetesClient: kubernetesClient,
		dynamicClient:    dynamicClient,
	}, nil
}

func (c *client) GetVirtualMachine(ctx context.Context, namespace string, name string) (*kubevirtapiv1.VirtualMachine, error) {
	resp, err := c.getResource(ctx, namespace, name, vmRes)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, err
		}
		return nil, errors.Wrap(err, "failed to get VirtualMachine")
	}
	var vm kubevirtapiv1.VirtualMachine
	err = c.fromUnstructedToInterface(*resp, &vm, "VirtualMachine")
	return &vm, err
}

func (c *client) ListVirtualMachine(ctx context.Context, namespace string, opts metav1.ListOptions) (*kubevirtapiv1.VirtualMachineList, error) {
	resp, err := c.listResource(ctx, namespace, vmRes, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list VirtualMachine")
	}
	var vmList kubevirtapiv1.VirtualMachineList
	err = c.fromUnstructedListToInterface(*resp, &vmList, "VirtualMachineList")
	return &vmList, err
}

func (c *client) DeleteVirtualMachine(ctx context.Context, namespace string, name string) error {
	return c.deleteResource(ctx, namespace, name, vmRes)
}

func (c *client) GetDataVolume(ctx context.Context, namespace string, name string) (*cdiv1.DataVolume, error) {
	resp, err := c.getResource(ctx, namespace, name, dvRes)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, err
		}
		return nil, errors.Wrap(err, "failed to get DataVolume")
	}
	var dv cdiv1.DataVolume
	err = c.fromUnstructedToInterface(*resp, &dv, "DataVolume")
	return &dv, err
}

func (c *client) ListDataVolume(ctx context.Context, namespace string, opts metav1.ListOptions) (*cdiv1.DataVolumeList, error) {
	resp, err := c.listResource(ctx, namespace, dvRes, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list DataVolume")
	}
	var dvList cdiv1.DataVolumeList
	err = c.fromUnstructedListToInterface(*resp, &dvList, "DataVolumeList")
	return &dvList, err
}

func (c *client) DeleteDataVolume(ctx context.Context, namespace string, name string) error {
	return c.deleteResource(ctx, namespace, name, dvRes)
}

func (c *client) GetSecret(ctx context.Context, namespace string, name string) (*corev1.Secret, error) {
	return c.kubernetesClient.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *client) ListSecret(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.SecretList, error) {
	return c.kubernetesClient.CoreV1().Secrets(namespace).List(ctx, opts)
}

func (c *client) DeleteSecret(ctx context.Context, namespace string, name string) error {
	return c.kubernetesClient.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (c *client) GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error) {
	return c.kubernetesClient.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
}

func (c *client) GetStorageClass(ctx context.Context, name string) (*storagev1.StorageClass, error) {
	return c.kubernetesClient.StorageV1().StorageClasses().Get(ctx, name, metav1.GetOptions{})
}

func (c *client) GetNetworkAttachmentDefinition(ctx context.Context, name string, namespace string) (*unstructured.Unstructured, error) {
	nadRes := schema.GroupVersionResource{
		Group:    nadv1.SchemeGroupVersion.Group,
		Version:  nadv1.SchemeGroupVersion.Version,
		Resource: "network-attachment-definitions",
	}
	return c.getResource(ctx, namespace, name, nadRes)
}

func (c *client) CreateSelfSubjectAccessReview(ctx context.Context, reviewObj *authv1.SelfSubjectAccessReview) (*authv1.SelfSubjectAccessReview, error) {
	return c.kubernetesClient.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, reviewObj, metav1.CreateOptions{})
}

func (c *client) GetHyperConverged(ctx context.Context, name string, namespace string) (*unstructured.Unstructured, error) {
	resource := schema.GroupVersionResource{
		Group:    "hco.kubevirt.io",
		Version:  "v1beta1",
		Resource: "hyperconvergeds",
	}
	return c.getResource(ctx, namespace, name, resource)
}

func (c *client) createResource(ctx context.Context, obj interface{}, namespace string, resource schema.GroupVersionResource) error {
	resultMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return errors.Wrapf(err, "failed to translate %s to Unstructed (for create operation)", resource.Resource)
	}
	input := unstructured.Unstructured{}
	input.SetUnstructuredContent(resultMap)
	resp, err := c.dynamicClient.Resource(resource).Namespace(namespace).Create(ctx, &input, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to create %s", resource.Resource)
	}
	unstructured := resp.UnstructuredContent()
	return runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, obj)
}

func (c *client) getResource(ctx context.Context, namespace string, name string, resource schema.GroupVersionResource) (*unstructured.Unstructured, error) {
	return c.dynamicClient.Resource(resource).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *client) deleteResource(ctx context.Context, namespace string, name string, resource schema.GroupVersionResource) error {
	return c.dynamicClient.Resource(resource).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (c *client) listResource(ctx context.Context, namespace string, resource schema.GroupVersionResource, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	return c.dynamicClient.Resource(resource).Namespace(namespace).List(ctx, opts)
}

func (c *client) fromUnstructedToInterface(src unstructured.Unstructured, dst interface{}, interfaceType string) error {
	unstructured := src.UnstructuredContent()
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, dst); err != nil {
		return errors.Wrapf(err, "failed to translate unstructed to %s", interfaceType)
	}
	return nil
}

func (c *client) fromUnstructedListToInterface(src unstructured.UnstructuredList, dst interface{}, interfaceType string) error {
	unstructured := src.UnstructuredContent()
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, dst); err != nil {
		return errors.Wrapf(err, "failed to translate unstructed to %s", interfaceType)
	}
	return nil
}
