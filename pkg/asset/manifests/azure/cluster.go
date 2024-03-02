package azure

import (
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils/cidr"
)

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) (*capiutils.GenerateClusterAssetsOutput, error) {
	manifests := []*asset.RuntimeFile{}
	mainCIDR := capiutils.CIDRFromInstallConfig(installConfig)

	session, err := installConfig.Azure.Session()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Azure session")
	}

	subnets, err := cidr.SplitIntoSubnetsIPv4(mainCIDR.String(), 2)
	if err != nil {
		return nil, errors.Wrap(err, "failed to split CIDR into subnets")
	}

	// CAPZ expects the capz-system to be created.
	azureNamespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "capz-system"}}
	manifests = append(manifests, &asset.RuntimeFile{
		Object: azureNamespace,
		File:   asset.File{Filename: "00_azure-namespace.yaml"},
	})

	azureCluster := &capz.AzureCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: capz.AzureClusterSpec{
			ResourceGroup: fmt.Sprintf("%s-rg", clusterID.InfraID),
			AzureClusterClassSpec: capz.AzureClusterClassSpec{
				SubscriptionID:   session.Credentials.SubscriptionID,
				Location:         installConfig.Config.Azure.Region,
				AzureEnvironment: string(installConfig.Azure.CloudName),
				IdentityRef: &corev1.ObjectReference{
					APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
					Kind:       "AzureClusterIdentity",
					Name:       clusterID.InfraID,
				},
			},
			NetworkSpec: capz.NetworkSpec{
				Vnet: capz.VnetSpec{
					ID: installConfig.Config.Azure.VirtualNetwork,
					VnetClassSpec: capz.VnetClassSpec{
						CIDRBlocks: []string{
							mainCIDR.String(),
						},
					},
				},
				Subnets: capz.Subnets{
					{
						SubnetClassSpec: capz.SubnetClassSpec{
							Name: "control-plane-subnet",
							Role: capz.SubnetControlPlane,
							CIDRBlocks: []string{
								subnets[0].String(),
							},
						},
					},
					{
						SubnetClassSpec: capz.SubnetClassSpec{
							Name: "worker-subnet",
							Role: capz.SubnetNode,
							CIDRBlocks: []string{
								subnets[1].String(),
							},
						},
					},
				},
			},
		},
	}
	manifests = append(manifests, &asset.RuntimeFile{
		Object: azureCluster,
		File:   asset.File{Filename: "02_azure-cluster.yaml"},
	})

	azureClientSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID + "-azure-client-secret",
			Namespace: capiutils.Namespace,
		},
		StringData: map[string]string{
			"clientSecret": session.Credentials.ClientSecret,
		},
	}
	manifests = append(manifests, &asset.RuntimeFile{
		Object: azureClientSecret,
		File:   asset.File{Filename: "01_azure-client-secret.yaml"},
	})

	id := &capz.AzureClusterIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: capz.AzureClusterIdentitySpec{
			Type:              capz.ManualServicePrincipal,
			AllowedNamespaces: &capz.AllowedNamespaces{}, // Allow all namespaces.
			ClientID:          session.Credentials.ClientID,
			ClientSecret: corev1.SecretReference{
				Name:      azureClientSecret.Name,
				Namespace: azureClientSecret.Namespace,
			},
			TenantID: session.Credentials.TenantID,
		},
	}
	manifests = append(manifests, &asset.RuntimeFile{
		Object: id,
		File:   asset.File{Filename: "01_aws-cluster-controller-identity-default.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRef: &corev1.ObjectReference{
			APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
			Kind:       "AzureCluster",
			Name:       azureCluster.Name,
			Namespace:  azureCluster.Namespace,
		},
	}, nil
}
