package vsphere

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/types/vsphere"
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

	for _, failureDomain := range installConfig.Config.VSphere.FailureDomains {
		if failureDomain.ZoneType == vsphere.HostGroupFailureDomain {
			dz := &capv.VSphereDeploymentZone{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name: failureDomain.Name,
				},
				Spec: capv.VSphereDeploymentZoneSpec{
					Server:        fmt.Sprintf("https://%s", failureDomain.Server),
					FailureDomain: failureDomain.Name,
					ControlPlane:  ptr.To(true),
					PlacementConstraint: capv.PlacementConstraint{
						ResourcePool: failureDomain.Topology.ResourcePool,
						Folder:       failureDomain.Topology.Folder,
					},
				},
			}

			dz.SetGroupVersionKind(capv.GroupVersion.WithKind("VSphereDeploymentZone"))

			fd := &capv.VSphereFailureDomain{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name: failureDomain.Name,
				},
				Spec: capv.VSphereFailureDomainSpec{
					Region: capv.FailureDomain{
						Name:        failureDomain.Region,
						Type:        capv.FailureDomainType(failureDomain.RegionType),
						TagCategory: "openshift-region",
					},
					Zone: capv.FailureDomain{
						Name:        failureDomain.Zone,
						Type:        capv.FailureDomainType(failureDomain.ZoneType),
						TagCategory: "openshift-zone",
					},
					Topology: capv.Topology{
						Datacenter:     failureDomain.Topology.Datacenter,
						ComputeCluster: &failureDomain.Topology.ComputeCluster,
						Hosts: &capv.FailureDomainHosts{
							VMGroupName:   fmt.Sprintf("%s-%s", clusterID.InfraID, failureDomain.Name),
							HostGroupName: failureDomain.Topology.HostGroup,
						},
						Networks:  failureDomain.Topology.Networks,
						Datastore: failureDomain.Topology.Datastore,
					},
				},
			}
			fd.SetGroupVersionKind(capv.GroupVersion.WithKind("VSphereFailureDomain"))

			manifests = append(manifests, &asset.RuntimeFile{
				Object: fd,
				File:   asset.File{Filename: fmt.Sprintf("01_vsphere-failuredomain-%s.yaml", failureDomain.Name)},
			})

			manifests = append(manifests, &asset.RuntimeFile{
				Object: dz,
				File:   asset.File{Filename: fmt.Sprintf("01_vsphere-deploymentzone-%s.yaml", failureDomain.Name)},
			})
		}
	}

	assetOutput.Manifests = manifests

	return assetOutput, nil
}
