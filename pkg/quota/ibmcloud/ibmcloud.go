package ibmcloud

import (
	"sort"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/quota"
	"github.com/openshift/installer/pkg/types"
)

// Constraints returns a list of quota constraints based on the InstallConfig.
// These constraints can be used to check if there is enough quota for creating
// a cluster for the install config.
func Constraints(config *types.InstallConfig, controlPlanes []machineapi.Machine, computes []machineapi.MachineSet) []quota.Constraint {
	computeReplicas := make([]int64, len(computes))
	for i, w := range computes {
		computeReplicas[i] = int64(*w.Spec.Replicas)
	}

	var ret []quota.Constraint
	for _, gen := range []constraintGenerator{
		network(config, len(controlPlanes)),
		instances(config, len(controlPlanes), computeReplicas),
	} {
		ret = append(ret, gen()...)
	}
	return aggregate(ret)
}

// constraintGenerator generates a list of constraints.
type constraintGenerator func() []quota.Constraint

func network(config *types.InstallConfig, controlPlaneCount int) func() []quota.Constraint {
	return func() []quota.Constraint {
		region := config.Platform.IBMCloud.Region

		// Floating IPs: bootstrap + control plane nodes.
		// The public API LB also uses floating IPs when publish is External.
		fipCount := int64(1 + controlPlaneCount)
		if config.Publish == types.ExternalPublishingStrategy || config.Publish == "" {
			fipCount++
		}

		ret := []quota.Constraint{
			{Name: "is/floating-ip", Region: region, Count: fipCount},
			{Name: "is/security-group", Region: region, Count: 6},
			{Name: "is/load-balancer", Region: region, Count: 2},
		}

		if config.Platform.IBMCloud.VPCName == "" {
			ret = append(ret, quota.Constraint{
				Name: "is/vpc", Region: region, Count: 1,
			})
		}

		return ret
	}
}

func instances(config *types.InstallConfig, controlPlaneCount int, computeReplicas []int64) func() []quota.Constraint {
	return func() []quota.Constraint {
		region := config.Platform.IBMCloud.Region

		// control plane + bootstrap
		count := int64(controlPlaneCount + 1)
		for _, r := range computeReplicas {
			count += r
		}

		return []quota.Constraint{
			{Name: "is/instance", Region: region, Count: count},
		}
	}
}

func aggregate(constraints []quota.Constraint) []quota.Constraint {
	if len(constraints) == 0 {
		return nil
	}

	sort.SliceStable(constraints, func(i, j int) bool {
		return constraints[i].Name < constraints[j].Name
	})

	i := 0
	for j := 1; j < len(constraints); j++ {
		if constraints[i].Name == constraints[j].Name && constraints[i].Region == constraints[j].Region {
			constraints[i].Count += constraints[j].Count
		} else {
			i++
			if i != j {
				constraints[i] = constraints[j]
			}
		}
	}
	return constraints[:i+1]
}
