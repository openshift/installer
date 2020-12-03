package openstack

import (
	"github.com/openshift/installer/pkg/asset/installconfig/openstack/validation"
	"github.com/openshift/installer/pkg/quota"
	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	openstackprovider "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"
)

// Constraints returns a list of quota constraints based on the InstallConfig.
// These constraints can be used to check if there is enough quota for creating a cluster
// for the install config.
func Constraints(ci *validation.CloudInfo, controlPlanes []machineapi.Machine, computes []machineapi.MachineSet) []quota.Constraint {
	ctrplConfigs := make([]*openstackprovider.OpenstackProviderSpec, len(controlPlanes))
	for i, m := range controlPlanes {
		ctrplConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*openstackprovider.OpenstackProviderSpec)
	}
	computeReplicas := make([]int64, len(computes))
	computeConfigs := make([]*openstackprovider.OpenstackProviderSpec, len(computes))
	for i, m := range computes {
		computeReplicas[i] = int64(*m.Spec.Replicas)
		computeConfigs[i] = m.Spec.Template.Spec.ProviderSpec.Value.Object.(*openstackprovider.OpenstackProviderSpec)
	}

	var ret []quota.Constraint
	ret = controlPlane(ci, ctrplConfigs)
	ret = append(ret, compute(ci, computeReplicas, computeConfigs)...)
	return aggregate(ret)
}

func controlPlane(ci *validation.CloudInfo, machines []*openstackprovider.OpenstackProviderSpec) []quota.Constraint {
	var ret []quota.Constraint
	for _, m := range machines {
		flavorInfo := ci.Flavors[m.Flavor]
		ret = append(ret, machineFlavorCoresToQuota(flavorInfo, ci), machineFlavorRAMToQuota(flavorInfo, ci))
	}
	ret = append(ret, instanceConstraint(int64(len(machines))))
	return ret
}

func compute(ci *validation.CloudInfo, replicas []int64, machines []*openstackprovider.OpenstackProviderSpec) []quota.Constraint {
	var ret []quota.Constraint
	for idx, m := range machines {
		flavorInfo := ci.Flavors[m.Flavor]
		coresConstraint := machineFlavorCoresToQuota(flavorInfo, ci)
		coresConstraint.Count = coresConstraint.Count * replicas[idx]
		ramConstraint := machineFlavorRAMToQuota(flavorInfo, ci)
		ramConstraint.Count = ramConstraint.Count * replicas[idx]
		ret = append(ret, instanceConstraint(int64(replicas[idx])), coresConstraint, ramConstraint)
	}
	return ret
}

func aggregate(quotas []quota.Constraint) []quota.Constraint {
	counts := map[string]int64{}
	for _, q := range quotas {
		counts[q.Name] = counts[q.Name] + q.Count
	}
	aggregatedQuotas := make([]quota.Constraint, 0, len(counts))
	for n, c := range counts {
		aggregatedQuotas = append(aggregatedQuotas, quota.Constraint{Name: n, Count: c})
	}
	return aggregatedQuotas
}

//constraintGenerator generates a list of constraints.
type constraintGenerator func() []quota.Constraint

func machineFlavorCoresToQuota(f validation.Flavor, ci *validation.CloudInfo) quota.Constraint {
	return generateConstraint("Cores", int64(f.VCPUs))
}

func machineFlavorRAMToQuota(f validation.Flavor, ci *validation.CloudInfo) quota.Constraint {
	return generateConstraint("RAM", int64(f.RAM))
}

func instanceConstraint(count int64) quota.Constraint {
	return generateConstraint("Instances", count)
}

func generateConstraint(name string, count int64) quota.Constraint {
	return quota.Constraint{
		Name:  name,
		Count: count,
	}
}
