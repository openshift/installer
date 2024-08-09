package vsphere

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
)

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) (*capiutils.GenerateClusterAssetsOutput, error) {
	manifests := []*asset.RuntimeFile{}

	assetOutput := &capiutils.GenerateClusterAssetsOutput{}

	for index, vcenter := range installConfig.Config.VSphere.VCenters {
		vsphereCreds := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("vsphere-creds-%d", index),
				Namespace: capiutils.Namespace,
			},
			Data: make(map[string][]byte),
		}
		vsphereCreds.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))

		vsphereCreds.Data["username"] = []byte(vcenter.Username)
		vsphereCreds.Data["password"] = []byte(vcenter.Password)

		manifests = append(manifests, &asset.RuntimeFile{
			Object: vsphereCreds,
			File:   asset.File{Filename: fmt.Sprintf("01_%v.yaml", vsphereCreds.Name)},
		})

		vsphereCluster := &capv.VSphereCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%v-%d", clusterID.InfraID, index),
				Namespace: capiutils.Namespace,
			},
			Spec: capv.VSphereClusterSpec{
				Server: fmt.Sprintf("https://%s", vcenter.Server),
				ControlPlaneEndpoint: capv.APIEndpoint{
					Host: fmt.Sprintf("api.%s.%s", installConfig.Config.ObjectMeta.Name, installConfig.Config.BaseDomain),
					Port: 6443,
				},
				IdentityRef: &capv.VSphereIdentityReference{
					Kind: capv.SecretKind,
					Name: vsphereCreds.Name,
				},
			},
		}
		vsphereCluster.SetGroupVersionKind(capv.GroupVersion.WithKind("VSphereCluster"))
		manifests = append(manifests, &asset.RuntimeFile{
			Object: vsphereCluster,
			File:   asset.File{Filename: fmt.Sprintf("01_vsphere-cluster-%d.yaml", index)},
		})

		infra := &corev1.ObjectReference{
			APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
			Kind:       "VSphereCluster",
			Name:       vsphereCluster.Name,
			Namespace:  capiutils.Namespace,
		}

		assetOutput.InfrastructureRefs = append(assetOutput.InfrastructureRefs, infra)
	}

	assetOutput.Manifests = manifests

	return assetOutput, nil
}
