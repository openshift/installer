package openstack

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
)

const (
	CloudName             = "openstack"
	CredentialsSecretName = "openstack-cloud-credentials"
)

func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) (*capiutils.GenerateClusterAssetsOutput, error) {
	manifests := []*asset.RuntimeFile{}
	openstackInstallConfig := installConfig.Config.OpenStack
	openStackCluster := &capo.OpenStackCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
			Labels: map[string]string{
				capi.ClusterNameLabel: clusterID.InfraID,
			},
		},
		Spec: capo.OpenStackClusterSpec{
			CloudName: CloudName,
			// TODO(stephenfin): Create credentials
			IdentityRef: &capo.OpenStackIdentityReference{
				Kind: "Secret",
				Name: CredentialsSecretName,
			},
			// We disable management of most networking resources since either
			// we (the installer) will create them, or the user will have
			// pre-created them as part of a "Bring Your Own Network (BYON)"
			// configuration
			ManagedSecurityGroups:      false,
			DisableAPIServerFloatingIP: true,
			// TODO(stephenfin): update when we support dual-stack (there are
			// potentially *two* IPs here)
			APIServerFixedIP:  openstackInstallConfig.APIVIPs[0],
			DNSNameservers:    openstackInstallConfig.ExternalDNS,
			ExternalNetworkID: openstackInstallConfig.ExternalNetwork,
			Tags: []string{
				fmt.Sprintf("openshiftClusterID=%s", clusterID.InfraID),
			},
		},
	}
	if openstackInstallConfig.ControlPlanePort != nil {
		// TODO(maysa): update when BYO dual-stack is supported in CAPO
		openStackCluster.Spec.Network.ID = openstackInstallConfig.ControlPlanePort.Network.ID
		openStackCluster.Spec.Network.Name = openstackInstallConfig.ControlPlanePort.Network.Name
		openStackCluster.Spec.Subnet.ID = openstackInstallConfig.ControlPlanePort.FixedIPs[0].Subnet.ID
		openStackCluster.Spec.Subnet.Name = openstackInstallConfig.ControlPlanePort.FixedIPs[0].Subnet.Name
	} else {
		openStackCluster.Spec.NodeCIDR = capiutils.CIDRFromInstallConfig(installConfig).String()
	}
	openStackCluster.SetGroupVersionKind(capo.GroupVersion.WithKind("OpenStackCluster"))

	manifests = append(manifests, &asset.RuntimeFile{
		Object: openStackCluster,
		File:   asset.File{Filename: "02_infra-cluster.yaml"},
	})

	// TODO(stephenfin): Create credentials request/cloud secret

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRef: &corev1.ObjectReference{
			APIVersion: "infrastructure.cluster.x-k8s.io/v1alpha7",
			Kind:       "OpenStackCluster",
			Name:       openStackCluster.Name,
			Namespace:  openStackCluster.Namespace,
		},
	}, nil
}
