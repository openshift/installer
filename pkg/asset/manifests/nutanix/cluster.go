package nutanix

import (
	"fmt"

	capnv1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	credentialTypes "github.com/nutanix-cloud-native/prism-go-client/environment/credentials"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	//capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

var credentialsDataFmt = `[{
  "type": "basic_auth",
  "data": {
	"prismCentral":{
	  "username": "%s",
	  "password": "%s"
	}
  }
}]`

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) (*capiutils.GenerateClusterAssetsOutput, error) {
	manifests := []*asset.RuntimeFile{}

	ic := installConfig.Config
	ntxCluster := &capnv1.NutanixCluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
			Kind:       "NutanixCluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: capnv1.NutanixClusterSpec{
			PrismCentral: &credentialTypes.NutanixPrismEndpoint{
				Address: ic.Platform.Nutanix.PrismCentral.Endpoint.Address,
				Port:    ic.Platform.Nutanix.PrismCentral.Endpoint.Port,
			},
			FailureDomains: []capnv1.NutanixFailureDomain{},
		},
	}

	stringData := make(map[string]string, 1)
	stringData["credentials"] = fmt.Sprintf(credentialsDataFmt, ic.Platform.Nutanix.PrismCentral.Username, ic.Platform.Nutanix.PrismCentral.Password)
	credSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nutanix-credentials",
			Namespace: capiutils.Namespace,
		},
		StringData: stringData,
	}

	manifests = append(manifests, &asset.RuntimeFile{
		Object: credSecret,
		File:   asset.File{Filename: "01_nutanix-creds.yaml"},
	})

	ntxCluster.Spec.PrismCentral.CredentialRef = &credentialTypes.NutanixCredentialReference{
		Kind:      credentialTypes.SecretKind,
		Name:      "nutanix-credentials",
		Namespace: capiutils.Namespace,
	}

	if ic.AdditionalTrustBundle != "" {
		ntxCluster.Spec.PrismCentral.AdditionalTrustBundle = &credentialTypes.NutanixTrustBundleReference{
			Kind: credentialTypes.NutanixTrustBundleKindString,
			Data: ic.AdditionalTrustBundle,
		}
	}

	for _, fd := range ic.Platform.Nutanix.FailureDomains {
		subnets := make([]capnv1.NutanixResourceIdentifier, 0, len(fd.SubnetUUIDs))
		for _, subnetUUID := range fd.SubnetUUIDs {
			subnets = append(subnets, capnv1.NutanixResourceIdentifier{Type: capnv1.NutanixIdentifierUUID, UUID: &subnetUUID})
		}

		ntxCluster.Spec.FailureDomains = append(ntxCluster.Spec.FailureDomains, capnv1.NutanixFailureDomain{
			Name: fd.Name,
			Cluster: capnv1.NutanixResourceIdentifier{
				Type: capnv1.NutanixIdentifierUUID,
				UUID: &fd.PrismElement.UUID,
			},
			Subnets: subnets,
		})
	}

	manifests = append(manifests, &asset.RuntimeFile{
		Object: ntxCluster,
		File:   asset.File{Filename: "01_nutanix-cluster.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRef: &corev1.ObjectReference{
			APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
			Kind:       "NutanixCluster",
			Name:       clusterID.InfraID,
			Namespace:  capiutils.Namespace,
		},
	}, nil
}
