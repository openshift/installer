package aws

import (
	"context"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
)

func workerMachineSetsRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "manifests/aws/99_openshift-cluster-api_worker-machinesets.yaml",
		RebuildHelper: workerMachineSetsRebuilder,
	}

	parents, err := asset.GetParents(
		ctx,
		getByName,
		"aws/ami",
		"aws/instance-type",
		"aws/region",
		"aws/user-tags",
		"aws/zones",
		"cluster-id",
		"cluster-name",
	)
	if err != nil {
		return nil, err
	}

	ami := string(parents["aws/ami"].Data)
	clusterID := string(parents["cluster-id"].Data)
	clusterName := string(parents["cluster-name"].Data)
	instanceType := string(parents["aws/instance-type"].Data)
	region := string(parents["aws/region"].Data)
	var userTags map[string]string
	err = yaml.Unmarshal(parents["aws/user-tags"].Data, &userTags)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal user tags")
	}

	var zones []string
	err = yaml.Unmarshal(parents["aws/zones"].Data, &zones)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal zones")
	}
	numZones := int64(len(zones))

	role := "worker"
	userDataSecret := fmt.Sprintf("%s-user-data", role)
	poolName := role  // FIXME: knob to control this
	total := int64(3) // FIXME: knob to control this

	var machineSets []runtime.RawExtension
	for idx, zone := range zones {
		name := fmt.Sprintf("%s-%s-%s", clusterName, poolName, zone)

		replicas := int32(total / numZones)
		if int64(idx) < total%numZones {
			replicas++
		}

		provider, err := provider(clusterID, clusterName, region, zone, instanceType, ami, userTags, role, userDataSecret)
		if err != nil {
			return nil, errors.Wrap(err, "create provider")
		}

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
				Replicas: &replicas,
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

		machineSets = append(machineSets, runtime.RawExtension{Object: &machineSet})
	}

	list := &metav1.List{
		TypeMeta: metav1.TypeMeta{
			Kind:       "List",
			APIVersion: "v1",
		},
		Items: machineSets,
	}

	asset.Data, err = yaml.Marshal(list)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	installerassets.Rebuilders["manifests/aws/99_openshift-cluster-api_worker-machinesets.yaml"] = workerMachineSetsRebuilder
}
