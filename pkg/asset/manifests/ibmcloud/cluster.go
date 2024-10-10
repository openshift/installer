package ibmcloud

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capibmcloud "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
)

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID, imageName string) (*capiutils.GenerateClusterAssetsOutput, error) {
	manifests := []*asset.RuntimeFile{}
	platform := installConfig.Config.Platform.IBMCloud

	// Collect and build information for Cluster manifest
	resourceGroup := clusterID.InfraID
	if platform.ResourceGroupName != "" {
		resourceGroup = platform.ResourceGroupName
	}

	// Create the IBMVPCCluster manifest
	ibmcloudCluster := &capibmcloud.IBMVPCCluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: capibmcloud.GroupVersion.String(),
			Kind:       "IBMVPCCluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: capibmcloud.IBMVPCClusterSpec{
			ControlPlaneEndpoint: capi.APIEndpoint{
				Host: fmt.Sprintf("api.%s.%s", installConfig.Config.ObjectMeta.Name, installConfig.Config.BaseDomain),
				Port: 6443,
			},
			Region:        platform.Region,
			ResourceGroup: resourceGroup,
		},
	}

	ibmcloudCluster.SetGroupVersionKind(capibmcloud.GroupVersion.WithKind("IBMVPCCluster"))

	manifests = append(manifests, &asset.RuntimeFile{
		Object: ibmcloudCluster,
		File:   asset.File{Filename: "01_ibmcloud-cluster.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRefs: []*corev1.ObjectReference{
			{
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
				Kind:       "IBMVPCCluster",
				Name:       ibmcloudCluster.Name,
				Namespace:  ibmcloudCluster.Namespace,
			},
		},
	}, nil
}
