// Package baremetal generates Machine objects for bare metal.
package baremetal

import (
	"fmt"
	"path"
	"strings"

	baremetalprovider "github.com/metal3-io/cluster-api-provider-baremetal/pkg/apis/baremetal/v1alpha1"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != baremetal.Name {
		return nil, fmt.Errorf("non bare metal configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != baremetal.Name {
		return nil, fmt.Errorf("non bare metal machine-pool: %q", poolPlatform)
	}
	clustername := config.ObjectMeta.Name
	platform := config.Platform.BareMetal

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	provider := provider(clustername, platform, osImage, userDataSecret)
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		machine := machineapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.openshift.io/v1beta1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-machine-api",
				Name:      fmt.Sprintf("%s-%s-%d", clustername, pool.Name, idx),
				Labels: map[string]string{
					"machine.openshift.io/cluster-api-cluster":      clustername,
					"machine.openshift.io/cluster-api-machine-role": role,
					"machine.openshift.io/cluster-api-machine-type": role,
				},
			},
			Spec: machineapi.MachineSpec{
				ProviderSpec: machineapi.ProviderSpec{
					Value: &runtime.RawExtension{Object: provider},
				},
				// we don't need to set Versions, because we control those via cluster operators.
			},
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

func provider(clusterName string, platform *baremetal.Platform, osImage string, userDataSecret string) *baremetalprovider.BareMetalMachineProviderSpec {
	// The rhcos-downloader container launched by the baremetal-operator downloads the image,
	// compresses it to speed up deployments and makes it available on platform.ClusterProvisioningIP, via http
	// osImage looks like:
	//   https://releases-art-rhcos.svc.ci.openshift.org/art/storage/releases/rhcos-4.2/42.80.20190725.1/rhcos-42.80.20190725.1-openstack.qcow2
	// But the cached URL looks like:
	//   http://172.22.0.3:6180/images/rhcos-42.80.20190725.1-openstack.qcow2/rhcos-42.80.20190725.1-compressed.qcow2
	// See https://github.com/openshift/ironic-rhcos-downloader for more details
	imageFilename := path.Base(osImage)
	compressedImageFilename := strings.Replace(imageFilename, "openstack", "compressed", 1)
	cacheImageURL := fmt.Sprintf("http://%s:6180/images/%s/%s", platform.ClusterProvisioningIP, imageFilename, compressedImageFilename)
	cacheChecksumURL := fmt.Sprintf("%s.md5sum", cacheImageURL)
	return &baremetalprovider.BareMetalMachineProviderSpec{
		Image: baremetalprovider.Image{
			URL:      cacheImageURL,
			Checksum: cacheChecksumURL,
		},
		UserData: &corev1.SecretReference{Name: userDataSecret},
	}
}
