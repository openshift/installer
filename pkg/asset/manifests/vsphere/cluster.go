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

	vcenter := installConfig.Config.VSphere.VCenters[0]

	vsphereCreds := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "vsphere-creds",
			Namespace: capiutils.Namespace,
		},
		Data: make(map[string][]byte),
	}
	vsphereCreds.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))

	vsphereCreds.Data["username"] = []byte(vcenter.Username)
	vsphereCreds.Data["password"] = []byte(vcenter.Password)

	manifests = append(manifests, &asset.RuntimeFile{
		Object: vsphereCreds,
		File:   asset.File{Filename: "01_vsphere-creds.yaml"},
	})

	vsphereCluster := &capv.VSphereCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
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
				Name: "vsphere-creds",
			},
		},
	}
	vsphereCluster.SetGroupVersionKind(capv.GroupVersion.WithKind("VSphereCluster"))
	manifests = append(manifests, &asset.RuntimeFile{
		Object: vsphereCluster,
		File:   asset.File{Filename: "01_vsphere-cluster.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRef: &corev1.ObjectReference{
			APIVersion: capv.GroupVersion.String(),
			Kind:       "VSphereCluster",
			Name:       clusterID.InfraID,
			Namespace:  capiutils.Namespace,
		},
	}, nil
}
