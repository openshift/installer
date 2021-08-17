package openstack

import (
	"github.com/sirupsen/logrus"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	operv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/asset/installconfig/openstack/validation"
	"github.com/openshift/installer/pkg/quota"
	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	openstackprovider "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"
)

// These numbers should reflect what is documented here:
// https://github.com/openshift/installer/tree/master/docs/user/openstack
// https://github.com/openshift/installer/blob/master/docs/user/openstack/kuryr.md
// Number of ports, routers, subnets and routers here don't include the constraints needed
// for each machine, which are calculated later
var minNetworkConstraint = buildNetworkConstraint(4, 0, 0, 0, 2, 56)
var minNetworkConstraintWithKuryr = buildNetworkConstraint(1490, 0, 249, 249, 249, 996)

func buildNetworkConstraint(ports, routers, subnets, networks, securityGroups, securityGroupRules int64) []quota.Constraint {
	return []quota.Constraint{
		{Name: "Port", Count: ports},
		{Name: "Router", Count: routers},
		{Name: "Subnet", Count: subnets},
		{Name: "Network", Count: networks},
		{Name: "SecurityGroup", Count: securityGroups},
		{Name: "SecurityGroupRule", Count: securityGroupRules},
	}
}

func getNetworkConstraints(networkType string) []quota.Constraint {
	if networkType == string(operv1.NetworkTypeKuryr) {
		return minNetworkConstraintWithKuryr
	}
	return minNetworkConstraint
}

// Constraints returns a list of quota constraints based on the InstallConfig.
// These constraints can be used to check if there is enough quota for creating a cluster
// for the install config.
func Constraints(ci *validation.CloudInfo, controlPlanes []machineapi.Machine, computes []machineapi.MachineSet, networkType string) []quota.Constraint {
	var constraints []quota.Constraint

	for i := 0; i < len(controlPlanes); i++ {
		constraints = append(constraints, machineConstraints(ci, &controlPlanes[i], networkType)...)
	}
	constraints = append(constraints, instanceConstraint(int64(len(controlPlanes))))

	for i := 0; i < len(computes); i++ {
		constraints = append(constraints, machineSetConstraints(ci, &computes[i], networkType)...)
	}
	constraints = append(constraints, instanceConstraint(int64(len(computes))))
	constraints = append(constraints, getNetworkConstraints(networkType)...)

	// If the cluster is using pre-provisioned networks, then the quota constraints should be
	// null because the installer doesn't need to create any resources.
	if ci.MachinesSubnet == nil {
		constraints = append(constraints, networkConstraint(1), routerConstraint(1), subnetConstraint(1))
	}

	return aggregate(constraints)
}

func getOpenstackProviderSpec(spec *machineapi.ProviderSpec) *openstackprovider.OpenstackProviderSpec {
	if spec.Value == nil {
		logrus.Warnf("Empty ProviderSpec")
		return nil
	}

	return spec.Value.Object.(*openstackprovider.OpenstackProviderSpec)
}

func machineConstraints(ci *validation.CloudInfo, machine *machineapi.Machine, networkType string) []quota.Constraint {
	osps := getOpenstackProviderSpec(&machine.Spec.ProviderSpec)
	if osps == nil {
		logrus.Warnf("Skipping quota validation for Machine %s: Invalid ProviderSpec", machine.Name)
		return nil
	}

	flavorInfo, ok := ci.Flavors[osps.Flavor]
	if !ok {
		// This will result in a separate validation failure
		logrus.Warnf("Skipping quota validation for Machine %s: Flavor '%s' is not valid",
			machine.Name, osps.Flavor)
		return nil
	}
	flavor := flavorInfo.Flavor
	return []quota.Constraint{machineFlavorCoresToQuota(&flavor), machineFlavorRAMToQuota(&flavor), portConstraint(int64(len(osps.Networks)))}
}

func machineSetConstraints(ci *validation.CloudInfo, ms *machineapi.MachineSet, networkType string) []quota.Constraint {
	osps := getOpenstackProviderSpec(&ms.Spec.Template.Spec.ProviderSpec)
	if osps == nil {
		logrus.Warnf("Skipping quota validation for MachineSet %s: Invalid ProviderSpec", ms.Name)
		return nil
	}

	replicas := ms.Spec.Replicas
	if replicas == nil {
		// We defensively check for nil Replicas here, but this should have
		// already been defaulted if omitted.

		logrus.Warnf("Skipping quota validation for MachineSet %s due to unspecified replica count", ms.Name)
		return nil
	}

	flavorInfo, ok := ci.Flavors[osps.Flavor]
	if !ok {
		// This will result in a separate validation failure
		logrus.Warnf("Skipping quota validation for MachineSet %s: Flavor '%s' is not valid", ms.Name, osps.Flavor)
		return nil
	}
	flavor := flavorInfo.Flavor

	coresConstraint := machineFlavorCoresToQuota(&flavor)
	coresConstraint.Count = coresConstraint.Count * int64(*replicas)
	ramConstraint := machineFlavorRAMToQuota(&flavor)
	ramConstraint.Count = ramConstraint.Count * int64(*replicas)
	portConstraint := portConstraint((int64(len(osps.Networks)) * int64(*replicas)))

	return []quota.Constraint{coresConstraint, ramConstraint, portConstraint}
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

func machineFlavorCoresToQuota(f *flavors.Flavor) quota.Constraint {
	return generateConstraint("Cores", int64(f.VCPUs))
}

func machineFlavorRAMToQuota(f *flavors.Flavor) quota.Constraint {
	return generateConstraint("RAM", int64(f.RAM))
}

func instanceConstraint(count int64) quota.Constraint {
	return generateConstraint("Instances", count)
}

func portConstraint(count int64) quota.Constraint {
	return generateConstraint("Port", count)
}

func routerConstraint(count int64) quota.Constraint {
	return generateConstraint("Router", count)
}

func networkConstraint(count int64) quota.Constraint {
	return generateConstraint("Network", count)
}

func subnetConstraint(count int64) quota.Constraint {
	return generateConstraint("Subnet", count)
}

func generateConstraint(name string, count int64) quota.Constraint {
	return quota.Constraint{
		Name:  name,
		Count: count,
	}
}
