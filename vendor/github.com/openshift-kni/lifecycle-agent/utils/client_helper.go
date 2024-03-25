package utils

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/openshift-kni/lifecycle-agent/api/seedreconfig"
	"github.com/openshift-kni/lifecycle-agent/internal/common"
	ocp_config_v1 "github.com/openshift/api/config/v1"
	operatorv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	"github.com/samber/lo"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func GetSecretData(ctx context.Context, name, namespace, key string, client runtimeclient.Client) (string, error) {
	secret := &corev1.Secret{}
	if err := client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, secret); err != nil {
		// NOTE: The error is intentionally left unwrapped here, so the caller
		// can check client.IgnoreNotFound on it
		return "", err //nolint:wrapcheck
	}

	data, ok := secret.Data[key]
	if !ok {
		return "", fmt.Errorf("did not find key %s in Secret %s/%s", key, name, namespace)
	}

	return string(data), nil
}

func GetConfigMapData(ctx context.Context, name, namespace, key string, client runtimeclient.Client) (string, error) {
	cm := &corev1.ConfigMap{}
	if err := client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, cm); err != nil {
		return "", fmt.Errorf("failed to get get configMap: %w", err)
	}

	data, ok := cm.Data[key]
	if !ok {
		return "", fmt.Errorf("did not find key %s in ConfigMap", key)
	}

	return data, nil
}

func GetClusterName(ctx context.Context, client runtimeclient.Client) (string, error) {
	installConfig, err := getInstallConfig(ctx, client)
	if err != nil {
		return "", fmt.Errorf("failed to get install config: %w", err)
	}
	return installConfig.Metadata.Name, nil
}

func GetClusterBaseDomain(ctx context.Context, client runtimeclient.Client) (string, error) {
	installConfig, err := getInstallConfig(ctx, client)
	if err != nil {
		return "", fmt.Errorf("failed to get install config: %w", err)
	}
	return installConfig.BaseDomain, nil
}

type ClusterInfo struct {
	OCPVersion               string
	BaseDomain               string
	ClusterName              string
	ClusterID                string
	NodeIP                   string
	ReleaseRegistry          string
	Hostname                 string
	MirrorRegistryConfigured bool
}

func GetClusterInfo(ctx context.Context, client runtimeclient.Client) (*ClusterInfo, error) {
	clusterVersion := &ocp_config_v1.ClusterVersion{}
	if err := client.Get(ctx, types.NamespacedName{Name: "version"}, clusterVersion); err != nil {
		return nil, fmt.Errorf("failed to get clusterversion: %w", err)
	}

	clusterName, err := GetClusterName(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get clusterName: %w", err)
	}

	clusterBaseDomain, err := GetClusterBaseDomain(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get clusterBaseDomain: %w", err)
	}

	node, err := GetSNOMasterNode(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get SNOMasterNode: %w", err)
	}
	ip, err := getNodeInternalIP(*node)
	if err != nil {
		return nil, err
	}
	hostname, err := getNodeHostname(*node)
	if err != nil {
		return nil, err
	}

	releaseRegistry, err := GetReleaseRegistry(ctx, client)
	if err != nil {
		return nil, err
	}

	mirrorRegistrySources, err := GetMirrorRegistrySourceRegistries(ctx, client)
	if err != nil {
		return nil, err
	}

	return &ClusterInfo{
		ClusterName:              clusterName,
		BaseDomain:               clusterBaseDomain,
		OCPVersion:               clusterVersion.Status.Desired.Version,
		ClusterID:                string(clusterVersion.Spec.ClusterID),
		NodeIP:                   ip,
		ReleaseRegistry:          releaseRegistry,
		Hostname:                 hostname,
		MirrorRegistryConfigured: len(mirrorRegistrySources) > 0,
	}, nil
}

// TODO: add dual stuck support
func getNodeInternalIP(node corev1.Node) (string, error) {
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			return addr.Address, nil
		}
	}
	return "", fmt.Errorf("failed to find node internal ip address")
}

func getNodeHostname(node corev1.Node) (string, error) {
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeHostName {
			return addr.Address, nil
		}
	}
	return "", fmt.Errorf("failed to find node hostname")
}

type installConfigMetadata struct {
	Name string `json:"name"`
}

type basicInstallConfig struct {
	BaseDomain string                `json:"baseDomain"`
	Metadata   installConfigMetadata `json:"metadata"`
}

func getInstallConfig(ctx context.Context, client runtimeclient.Client) (*basicInstallConfig, error) {
	cm := &corev1.ConfigMap{}
	err := client.Get(ctx, types.NamespacedName{Name: common.InstallConfigCM, Namespace: common.InstallConfigCMNamespace}, cm)
	if err != nil {
		return nil, fmt.Errorf("could not get configMap: %w", err)
	}

	data, ok := cm.Data["install-config"]
	if !ok {
		return nil, fmt.Errorf("did not find key install-config in configmap")
	}

	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(data)), 4096)
	instConf := &basicInstallConfig{}
	if err := decoder.Decode(instConf); err != nil {
		return nil, fmt.Errorf("failed to decode install config, err: %w", err)
	}
	return instConf, nil
}

func GetCSVDeployment(ctx context.Context, client runtimeclient.Client) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}
	if err := client.Get(ctx,
		types.NamespacedName{
			Name:      common.CsvDeploymentName,
			Namespace: common.CsvDeploymentNamespace},
		deployment); err != nil {
		return nil, fmt.Errorf("failed to get cluster version deployment, err: %w", err)
	}

	return deployment, nil
}

func GetInfrastructure(ctx context.Context, client runtimeclient.Client) (*ocp_config_v1.Infrastructure, error) {
	infrastructure := &ocp_config_v1.Infrastructure{}
	if err := client.Get(ctx,
		types.NamespacedName{
			Name: common.OpenshiftInfraCRName},
		infrastructure); err != nil {
		return nil, fmt.Errorf("failed to get infra CR: %w", err)
	}

	return infrastructure, nil
}

func GetReleaseRegistry(ctx context.Context, client runtimeclient.Client) (string, error) {
	deployment, err := GetCSVDeployment(ctx, client)
	if err != nil {
		return "", err
	}

	return strings.Split(deployment.Spec.Template.Spec.Containers[0].Image, "/")[0], nil
}
func ReadSeedReconfigurationFromFile(path string) (*seedreconfig.SeedReconfiguration, error) {
	data := &seedreconfig.SeedReconfiguration{}
	err := ReadYamlOrJSONFile(path, data)
	return data, err
}

func ExtractRegistryFromImage(image string) string {
	return strings.Split(image, "/")[0]
}

func GetMirrorRegistrySourceRegistries(ctx context.Context, client runtimeclient.Client) ([]string, error) {
	var sourceRegistries []string
	allNamespaces := runtimeclient.ListOptions{Namespace: metav1.NamespaceAll}
	currentIcps := &operatorv1alpha1.ImageContentSourcePolicyList{}
	if err := client.List(ctx, currentIcps, &allNamespaces); err != nil {
		return nil, fmt.Errorf("failed to list ImageContentSourcePolicy: %w", err)
	}
	for _, icsp := range currentIcps.Items {
		for _, rdp := range icsp.Spec.RepositoryDigestMirrors {
			sourceRegistries = append(sourceRegistries, ExtractRegistryFromImage(rdp.Source))
		}
	}
	currentIdms := ocp_config_v1.ImageDigestMirrorSetList{}
	if err := client.List(ctx, &currentIdms, &allNamespaces); err != nil {
		return nil, fmt.Errorf("failed to list ImageDigestMirrorSet: %w", err)
	}

	for _, idms := range currentIdms.Items {
		for _, idm := range idms.Spec.ImageDigestMirrors {
			sourceRegistries = append(sourceRegistries, ExtractRegistryFromImage(idm.Source))
		}
	}
	return sourceRegistries, nil
}

func HasProxy(ctx context.Context, client runtimeclient.Client) (bool, error) {
	proxy := &ocp_config_v1.Proxy{}
	if err := client.Get(ctx, types.NamespacedName{Name: common.OpenshiftProxyCRName}, proxy); err != nil {
		return false, fmt.Errorf("failed to get proxy CR: %w", err)
	}

	if proxy.Spec.HTTPProxy == "" && proxy.Spec.HTTPSProxy == "" && proxy.Spec.NoProxy == "" {
		return false, nil
	}

	return true, nil
}

func ShouldOverrideSeedRegistry(ctx context.Context, client runtimeclient.Client, mirrorRegistryConfigured bool, releaseRegistry string) (bool, error) {
	mirroredRegistries, err := GetMirrorRegistrySourceRegistries(ctx, client)
	if err != nil {
		return false, err
	}
	isMirrorRegistryConfigured := len(mirroredRegistries) > 0

	// if snoa doesn't have mirror registry but seed have we should try to override registry
	if !isMirrorRegistryConfigured && mirrorRegistryConfigured {
		return true, err
	}

	return !lo.Contains(mirroredRegistries, releaseRegistry), nil
}
