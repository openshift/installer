package libvirt

import (
	"context"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
)

func workerMachineSetsRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "manifests/libvirt/99_openshift-cluster-api_worker-machinesets.yaml",
		RebuildHelper: workerMachineSetsRebuilder,
	}

	parents, err := asset.GetParents(
		ctx,
		getByName,
		"cluster-name",
		"libvirt/uri",
		"network/cluster-cidr",
	)
	if err != nil {
		return nil, err
	}

	clusterCIDR := string(parents["network/cluster-cidr"].Data)
	clusterName := string(parents["cluster-name"].Data)
	uri := string(parents["libvirt/uri"].Data)

	role := "worker"
	userDataSecret := fmt.Sprintf("%s-user-data", role)
	poolName := role  // FIXME: knob to control this
	total := int64(1) // FIXME: knob to control this

	provider, err := provider(uri, clusterName, clusterCIDR, userDataSecret)
	if err != nil {
		return nil, errors.Wrap(err, "create provider")
	}

	name := fmt.Sprintf("%s-%s-%d", clusterName, poolName, 0)
	machineSet := clusterapi.MachineSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MachineSet",
			APIVersion: "cluster.k8s.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "openshift-cluster-api",
			Labels: map[string]string{
				"sigs.k8s.io/cluster-api-cluster":      clusterName,
				"sigs.k8s.io/cluster-api-machine-role": role,
				"sigs.k8s.io/cluster-api-machine-type": role,
			},
		},
		Spec: clusterapi.MachineSetSpec{
			Replicas: pointer.Int32Ptr(int32(total)),
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"sigs.k8s.io/cluster-api-machineset": name,
					"sigs.k8s.io/cluster-api-cluster":    clusterName,
				},
			},
			Template: clusterapi.MachineTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"sigs.k8s.io/cluster-api-machineset":   name,
						"sigs.k8s.io/cluster-api-cluster":      clusterName,
						"sigs.k8s.io/cluster-api-machine-role": role,
						"sigs.k8s.io/cluster-api-machine-type": role,
					},
				},
				Spec: clusterapi.MachineSpec{
					ProviderConfig: clusterapi.ProviderConfig{
						Value: &runtime.RawExtension{Object: provider},
					},
					// we don't need to set Versions, because we control those via cluster operators.
				},
			},
		},
	}

	asset.Data, err = yaml.Marshal(machineSet)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	installerassets.Rebuilders["manifests/libvirt/99_openshift-cluster-api_worker-machinesets.yaml"] = workerMachineSetsRebuilder
}
