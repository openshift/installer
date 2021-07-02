// Package baremetal generates Machine objects for bare metal.
package baremetal

import (
	"fmt"
	"net"
	"net/url"
	"path"
	"strings"

	baremetalprovider "github.com/metal3-io/cluster-api-provider-baremetal/pkg/apis/baremetal/v1alpha1"
	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"github.com/pkg/errors"
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
	platform := config.Platform.BareMetal

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	provider, err := provider(platform, osImage, userDataSecret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create provider")
	}
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		machine := machineapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.openshift.io/v1beta1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-machine-api",
				Name:      fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx),
				Labels: map[string]string{
					"machine.openshift.io/cluster-api-cluster":      clusterID,
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

func provider(platform *baremetal.Platform, osImage string, userDataSecret string) (*baremetalprovider.BareMetalMachineProviderSpec, error) {
	// The machine-os-downloader container launched by the baremetal-operator downloads the image,
	// compresses it to speed up deployments and makes it available on platform.ClusterProvisioningIP, via http
	// osImage looks like:
	//   https://releases-art-rhcos.svc.ci.openshift.org/art/storage/releases/rhcos-4.2/42.80.20190725.1/rhcos-42.80.20190725.1-openstack.qcow2?sha256sum=123
	// And the cached URL looks like:
	//   http://172.22.0.3:6180/images/rhcos-42.80.20190725.1-openstack.qcow2/cached-rhcos-42.80.20190725.1-openstack.qcow2
	// See https://github.com/openshift/ironic-rhcos-downloader for more details
	// The image is now formatted with a query string containing the sha256sum, we strip that here
	// and it will be consumed for validation in ironic-machine-os-downloader
	imageURL, err := url.Parse(osImage)
	if err != nil {
		return nil, errors.Wrap(err, "invalid osImage URL format")
	}
	imageURL.RawQuery = ""
	imageURL.Fragment = ""
	// We strip any .gz/.xz suffix because ironic-machine-os-downloader unzips the image
	// ref https://github.com/openshift/ironic-rhcos-downloader/pull/12
	imageFilename := path.Base(strings.TrimSuffix(imageURL.String(), ".gz"))
	imageFilename = strings.TrimSuffix(imageFilename, ".xz")
	cachedImageFilename := "cached-" + imageFilename

	cacheImageIP := platform.ClusterProvisioningIP
	if platform.ProvisioningNetwork == baremetal.DisabledProvisioningNetwork && platform.ClusterProvisioningIP == "" {
		cacheImageIP = platform.APIVIP
	}
	cacheImageURL := fmt.Sprintf("http://%s/images/%s/%s", net.JoinHostPort(cacheImageIP, "6181"), imageFilename, cachedImageFilename)
	cacheChecksumURL := fmt.Sprintf("%s.md5sum", cacheImageURL)
	config := &baremetalprovider.BareMetalMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "baremetal.cluster.k8s.io/v1alpha1",
			Kind:       "BareMetalMachineProviderSpec",
		},
		Image: baremetalprovider.Image{
			URL:      cacheImageURL,
			Checksum: cacheChecksumURL,
		},
		UserData: &corev1.SecretReference{Name: userDataSecret},
	}
	return config, nil
}
