package nutanix

import (
	"fmt"

	capnv1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
	credentialTypes "github.com/nutanix-cloud-native/prism-go-client/environment/credentials"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
)

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) (*capiutils.GenerateClusterAssetsOutput, error) {
	manifests := []*asset.RuntimeFile{}
	ic := installConfig.Config

	// generate the NutanixCluster manifest.
	ntxCluster := &capnv1.NutanixCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: capnv1.NutanixClusterSpec{
			ControlPlaneEndpoint: capv1.APIEndpoint{
				Host: fmt.Sprintf("api.%s.%s", installConfig.Config.ObjectMeta.Name, installConfig.Config.BaseDomain),
				Port: 6443,
			},
			PrismCentral: &credentialTypes.NutanixPrismEndpoint{
				Address: ic.Platform.Nutanix.PrismCentral.Endpoint.Address,
				Port:    ic.Platform.Nutanix.PrismCentral.Endpoint.Port,
			},
			ControlPlaneFailureDomains: []corev1.LocalObjectReference{},
		},
	}
	ntxCluster.SetGroupVersionKind(capnv1.GroupVersion.WithKind("NutanixCluster"))

	// generate the nutanix-credentials secret manifest.
	// #nosec G101
	credentialsDataFmt := `[{
		"type": "basic_auth",
		"data": {
		  "prismCentral":{
			"username": "%s",
			"password": "%s"
		  }
		}
	  }]`
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
	credSecret.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))

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

	var failureDomains []capnv1.NutanixFailureDomain
	var controlPlaneFailureDomians []corev1.LocalObjectReference
	for _, fd := range ic.Platform.Nutanix.FailureDomains {
		subnets := make([]capnv1.NutanixResourceIdentifier, 0, len(fd.SubnetUUIDs))
		for _, subnetUUID := range fd.SubnetUUIDs {
			subnets = append(subnets, capnv1.NutanixResourceIdentifier{Type: capnv1.NutanixIdentifierUUID, UUID: ptr.To(subnetUUID)})
		}
		_ = append(failureDomains, capnv1.NutanixFailureDomain{
			ObjectMeta: metav1.ObjectMeta{
				Name: fd.Name,
			},
			Spec: capnv1.NutanixFailureDomainSpec{
				Subnets: subnets,
				PrismElementCluster: capnv1.NutanixResourceIdentifier{
					Type: capnv1.NutanixIdentifierUUID,
					UUID: ptr.To(fd.PrismElement.UUID),
				},
			},
		})
		controlPlaneFailureDomians = append(controlPlaneFailureDomians, corev1.LocalObjectReference{
			Name: fd.Name,
		})
	}
	ntxCluster.Spec.ControlPlaneFailureDomains = controlPlaneFailureDomians

	manifests = append(manifests, &asset.RuntimeFile{
		Object: ntxCluster,
		File:   asset.File{Filename: "01_nutanix-cluster.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRefs: []*corev1.ObjectReference{
			{
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
				Kind:       "NutanixCluster",
				Name:       clusterID.InfraID,
				Namespace:  capiutils.Namespace,
			},
		},
	}, nil
}
