package openstack

import (
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	"github.com/sirupsen/logrus"

	machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"
	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset/installconfig/openstack/validation"
	"github.com/openshift/installer/pkg/quota"
)

// These numbers should reflect what is documented here:
// https://github.com/openshift/installer/tree/master/docs/user/openstack
// Number of ports, routers, subnets and routers here don't include the constraints needed
// for each machine, which are calculated later
var minNetworkConstraint = buildNetworkConstraint(4, 0, 0, 0, 2, 56)

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

func getNetworkConstraints() []quota.Constraint {
	return minNetworkConstraint
}

// Constraints returns a list of quota constraints based on the InstallConfig.
// These constraints can be used to check if there is enough quota for creating a cluster
// for the install config.
func Constraints(ci *validation.CloudInfo, controlPlanes []machineapi.Machine, computes []machineapi.MachineSet) []quota.Constraint {
	var constraints []quota.Constraint

	for i := 0; i < len(controlPlanes); i++ {
		constraints = append(constraints, machineConstraints(ci, &controlPlanes[i])...)
	}
	constraints = append(constraints, instanceConstraint(int64(len(controlPlanes))))

	for i := 0; i < len(computes); i++ {
		constraints = append(constraints, machineSetConstraints(ci, &computes[i])...)
	}
	constraints = append(constraints, instanceConstraint(int64(len(computes))))
	constraints = append(constraints, getNetworkConstraints()...)

	// If the cluster is using pre-provisioned networks, then the quota constraints should be
	// null because the installer doesn't need to create any resources.
	if len(ci.ControlPlanePortSubnets) == 0 {
		constraints = append(constraints, networkConstraint(1), routerConstraint(1), subnetConstraint(1))
	}
	// if the cluster does not have worker nodes then reduce the server group value from 2 to 1
	numServerGroups := int64(2)
	if len(computes) == 0 {
		numServerGroups--
	}
	constraints = append(constraints, serverGroupsConstraint(numServerGroups))
	constraints = append(constraints, serverGroupMembersConstraint(int64(len(controlPlanes)+len(computes))))

	return aggregate(constraints)
}

func getOpenstackProviderSpec(spec *machineapi.ProviderSpec) *machinev1alpha1.OpenstackProviderSpec {
	if spec.Value == nil {
		logrus.Warnf("Empty ProviderSpec")
		return nil
	}

	return spec.Value.Object.(*machinev1alpha1.OpenstackProviderSpec)
}

func machineConstraints(ci *validation.CloudInfo, machine *machineapi.Machine) []quota.Constraint {
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

func machineSetConstraints(ci *validation.CloudInfo, ms *machineapi.MachineSet) []quota.Constraint {
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

func serverGroupsConstraint(count int64) quota.Constraint {
	return generateConstraint("ServerGroups", count)
}

func serverGroupMembersConstraint(count int64) quota.Constraint {
	return generateConstraint("ServerGroupMembers", count)
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
