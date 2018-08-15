package libvirt

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ghodss/yaml"
	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1alpha1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
)

func masterMachinesRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "manifests/libvirt/99_openshift-cluster-api_master-machines.yaml",
		RebuildHelper: masterMachinesRebuilder,
	}

	parents, err := asset.GetParents(
		ctx,
		getByName,
		"cluster-name",
		"libvirt/uri",
		"machines/master-count",
		"network/cluster-cidr",
	)
	if err != nil {
		return nil, err
	}

	clusterCIDR := string(parents["network/cluster-cidr"].Data)
	clusterName := string(parents["cluster-name"].Data)
	uri := string(parents["libvirt/uri"].Data)

	masterCount, err := strconv.ParseUint(string(parents["machines/master-count"].Data), 10, 32)
	if err != nil {
		return nil, errors.Wrap(err, "parse master count")
	}

	role := "master"
	userDataSecret := fmt.Sprintf("%s-user-data", role)
	poolName := role // FIXME: knob to control this?
	total := int64(masterCount)

	provider, err := provider(uri, clusterName, clusterCIDR, userDataSecret)
	if err != nil {
		return nil, errors.Wrap(err, "create provider")
	}

	var machines []runtime.RawExtension
	for idx := int64(0); idx < total; idx++ {
		machine := clusterapi.Machine{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Machine",
				APIVersion: "cluster.k8s.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s-%d", clusterName, poolName, idx),
				Namespace: "openshift-cluster-api",
				Labels: map[string]string{
					"sigs.k8s.io/cluster-api-cluster":      clusterName,
					"sigs.k8s.io/cluster-api-machine-role": role,
					"sigs.k8s.io/cluster-api-machine-type": role,
				},
			},
			Spec: clusterapi.MachineSpec{
				ProviderConfig: clusterapi.ProviderConfig{
					Value: &runtime.RawExtension{Object: provider},
				},
				// we don't need to set Versions, because we control those via operators.
			},
		}

		machines = append(machines, runtime.RawExtension{Object: &machine})
	}

	list := &metav1.List{
		TypeMeta: metav1.TypeMeta{
			Kind:       "List",
			APIVersion: "v1",
		},
		Items: machines,
	}

	asset.Data, err = yaml.Marshal(list)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func provider(uri, clusterName, clusterCIDR, userDataSecret string) (*libvirtprovider.LibvirtMachineProviderConfig, error) {
	return &libvirtprovider.LibvirtMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "LibvirtMachineProviderConfig",
			APIVersion: "libvirtproviderconfig.k8s.io/v1alpha1",
		},
		DomainMemory: 2048,
		DomainVcpu:   2,
		Ignition: &libvirtprovider.Ignition{
			UserDataSecret: userDataSecret,
		},
		Volume: &libvirtprovider.Volume{
			PoolName:     "default",
			BaseVolumeID: fmt.Sprintf("/var/lib/libvirt/images/%s-base", clusterName),
		},
		NetworkInterfaceName:    clusterName,
		NetworkInterfaceAddress: clusterCIDR,
		Autostart:               false,
		URI:                     uri,
	}, nil
}

func init() {
	installerassets.Rebuilders["manifests/libvirt/99_openshift-cluster-api_master-machines.yaml"] = masterMachinesRebuilder
	installerassets.Defaults["libvirt/machines/master-count"] = installerassets.ConstantDefault([]byte("1"))
}
