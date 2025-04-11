package openstack

import (
	"fmt"
	"os"

	"github.com/gophercloud/utils/v2/openstack/clientconfig"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/optional"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
)

const (
	cloudName = "openstack"
)

// GenerateClusterAssets generates the cluster manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) (*capiutils.GenerateClusterAssetsOutput, error) {
	manifests := []*asset.RuntimeFile{}
	openstackInstallConfig := installConfig.Config.OpenStack

	var (
		externalNetwork        *capo.NetworkParam
		disableExternalNetwork optional.Bool
	)
	if e := openstackInstallConfig.ExternalNetwork; e != "" {
		externalNetwork = &capo.NetworkParam{Filter: &capo.NetworkFilter{Name: e}}
	} else {
		disableExternalNetwork = ptr.To(true)
	}

	openStackCluster := &capo.OpenStackCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
			Labels: map[string]string{
				capi.ClusterNameLabel: clusterID.InfraID,
			},
		},
		Spec: capo.OpenStackClusterSpec{
			IdentityRef: capo.OpenStackIdentityReference{
				Name:      clusterID.InfraID + "-cloud-config",
				CloudName: cloudName,
			},
			// We disable management of most networking resources since either
			// we (the installer) will create them, or the user will have
			// pre-created them as part of a "Bring Your Own Network (BYON)"
			// configuration
			ManagedSecurityGroups:      nil,
			DisableAPIServerFloatingIP: ptr.To(true),
			// TODO(stephenfin): update when we support dual-stack (there are
			// potentially *two* IPs here)
			APIServerFixedIP:       &openstackInstallConfig.APIVIPs[0],
			ExternalNetwork:        externalNetwork,
			DisableExternalNetwork: disableExternalNetwork,
			Tags: []string{
				fmt.Sprintf("openshiftClusterID=%s", clusterID.InfraID),
			},
		},
	}
	if cpPort := openstackInstallConfig.ControlPlanePort; cpPort != nil {
		if networkID := cpPort.Network.ID; networkID != "" {
			openStackCluster.Spec.Network = &capo.NetworkParam{ID: &networkID}
		} else if networkName := cpPort.Network.Name; networkName != "" {
			openStackCluster.Spec.Network = &capo.NetworkParam{Filter: &capo.NetworkFilter{Name: networkName}}
		}
		openStackCluster.Spec.Subnets = make([]capo.SubnetParam, len(cpPort.FixedIPs))
		for i := range cpPort.FixedIPs {
			if subnetID := cpPort.FixedIPs[i].Subnet.ID; subnetID != "" {
				openStackCluster.Spec.Subnets[i] = capo.SubnetParam{ID: &subnetID}
			} else {
				openStackCluster.Spec.Subnets[i] = capo.SubnetParam{Filter: &capo.SubnetFilter{Name: cpPort.FixedIPs[i].Subnet.Name}}
			}
		}
	} else {
		openStackCluster.Spec.ManagedSubnets = []capo.SubnetSpec{
			{
				CIDR:           capiutils.CIDRFromInstallConfig(installConfig).String(),
				DNSNameservers: openstackInstallConfig.ExternalDNS,
			},
		}
	}
	openStackCluster.SetGroupVersionKind(capo.SchemeGroupVersion.WithKind("OpenStackCluster"))

	manifests = append(manifests, &asset.RuntimeFile{
		Object: openStackCluster,
		File:   asset.File{Filename: "02_infra-cluster.yaml"},
	})

	cloudConfig, err := generateCloudConfig(installConfig)
	if err != nil {
		return nil, err
	}

	openStackIdentity := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID + "-cloud-config",
			Namespace: capiutils.Namespace,
		},
		Data: cloudConfig,
	}
	openStackIdentity.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))

	manifests = append(manifests, &asset.RuntimeFile{
		Object: openStackIdentity,
		File:   asset.File{Filename: "02_openstack-cloud-config.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRefs: []*corev1.ObjectReference{
			{
				APIVersion: capo.SchemeGroupVersion.String(),
				Kind:       "OpenStackCluster",
				Name:       openStackCluster.Name,
				Namespace:  openStackCluster.Namespace,
			},
		},
	}, nil
}

func generateCloudConfig(installConfig *installconfig.InstallConfig) (map[string][]byte, error) {
	opts := new(clientconfig.ClientOpts)
	opts.Cloud = installConfig.Config.Platform.OpenStack.Cloud

	cloud, err := clientconfig.GetCloudFromYAML(opts)
	if err != nil {
		return nil, err
	}

	// We need to replace the local cacert path with the one used by CAPO
	caCert := []byte{}
	if cloud.CACertFile != "" {
		caCert, err = os.ReadFile(cloud.CACertFile)
		if err != nil {
			return nil, err
		}

		// TODO: Verify this path. This is taken from CAPO directly
		// https://github.com/kubernetes-sigs/cluster-api-provider-openstack/blob/main/templates/env.rc
		cloud.CACertFile = "/etc/certs/cacert"
	}

	clouds := make(map[string]map[string]*clientconfig.Cloud)
	clouds["clouds"] = map[string]*clientconfig.Cloud{
		cloudName: cloud,
	}

	cloudsYAML, err := yaml.Marshal(clouds)
	if err != nil {
		return nil, err
	}

	creds := map[string][]byte{
		"clouds.yaml": cloudsYAML,
	}
	if len(caCert) != 0 {
		creds["cacert"] = caCert
	}

	return creds, nil
}
