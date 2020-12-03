package aws

import (
	"sort"

	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	"k8s.io/apimachinery/pkg/util/sets"
	awsprovider "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1beta1"

	"github.com/openshift/installer/pkg/quota"
	"github.com/openshift/installer/pkg/types"
)

// Constraints returns a list of quota constraints based on the InstallConfig.
// These constraints can be used to check if there is enough quota for creating a cluster
// for the isntall config.
func Constraints(config *types.InstallConfig, controlPlanes []machineapi.Machine, computes []machineapi.MachineSet, instanceTypes map[string]InstanceTypeInfo) []quota.Constraint {
	ctrplConfigs := make([]*awsprovider.AWSMachineProviderConfig, len(controlPlanes))
	for i, m := range controlPlanes {
		ctrplConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*awsprovider.AWSMachineProviderConfig)
	}
	computeReplicas := make([]int64, len(computes))
	computeConfigs := make([]*awsprovider.AWSMachineProviderConfig, len(computes))
	for i, w := range computes {
		computeReplicas[i] = int64(*w.Spec.Replicas)
		computeConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*awsprovider.AWSMachineProviderConfig)
	}

	var ret []quota.Constraint
	for _, gen := range []constraintGenerator{
		network(config, append(ctrplConfigs, computeConfigs...)),
		controlPlane(config, ctrplConfigs, instanceTypes),
		compute(config, computeReplicas, computeConfigs, instanceTypes),
		others,
	} {
		ret = append(ret, gen()...)
	}
	return aggregate(ret)
}

func aggregate(quotas []quota.Constraint) []quota.Constraint {
	sort.SliceStable(quotas, func(i, j int) bool {
		return quotas[i].Name < quotas[j].Name
	})

	i := 0
	for j := 1; j < len(quotas); j++ {
		if quotas[i].Name == quotas[j].Name && quotas[i].Region == quotas[j].Region {
			quotas[i].Count += quotas[j].Count
		} else {
			i++
			if i != j {
				quotas[i] = quotas[j]
			}
		}
	}
	return quotas[:i+1]
}

// constraintGenerator generates a list of constraints.
type constraintGenerator func() []quota.Constraint

func network(config *types.InstallConfig, machines []*awsprovider.AWSMachineProviderConfig) func() []quota.Constraint {
	return func() []quota.Constraint {
		zones := sets.NewString()
		for _, m := range machines {
			zones.Insert(m.Placement.AvailabilityZone)
		}

		var ret []quota.Constraint
		ret = append(ret, quota.Constraint{
			Name:   "vpc/L-E79EC296", // 2 sg
			Region: config.Platform.AWS.Region,
			Count:  2,
		})
		if len(config.Platform.AWS.Subnets) == 0 {
			ret = append(ret, []quota.Constraint{{
				Name:   "vpc/L-F678F1CE", // 1 vpc
				Region: config.Platform.AWS.Region,
				Count:  1,
			}, {
				Name:   "vpc/L-A4707A72", // 1 ig
				Region: config.Platform.AWS.Region,
				Count:  1,
			}, {
				Name:   "vpc/L-FE5A380F", // 1 nat gw per AZ
				Region: config.Platform.AWS.Region,
				Count:  1,
			}}...)

			ret = append(ret, quota.Constraint{
				Name:   "ec2/L-0263D0A3", // 1 eip per AZ
				Region: config.Platform.AWS.Region,
				Count:  int64(zones.Len()),
			})
		}

		return ret
	}
}

func controlPlane(config *types.InstallConfig, machines []*awsprovider.AWSMachineProviderConfig, instanceTypes map[string]InstanceTypeInfo) func() []quota.Constraint {
	return func() []quota.Constraint {
		var ret []quota.Constraint
		for _, m := range machines {
			q := machineTypeToQuota(m.InstanceType, instanceTypes)
			q.Region = config.Platform.AWS.Region
			ret = append(ret, q)
		}
		return ret
	}
}

func compute(config *types.InstallConfig, replicas []int64, machines []*awsprovider.AWSMachineProviderConfig, instanceTypes map[string]InstanceTypeInfo) func() []quota.Constraint {
	return func() []quota.Constraint {
		var ret []quota.Constraint
		for idx, m := range machines {
			q := machineTypeToQuota(m.InstanceType, instanceTypes)
			q.Count = q.Count * replicas[idx]
			q.Region = config.Platform.AWS.Region
			ret = append(ret, q)
		}
		return ret
	}
}

func others() []quota.Constraint {
	return []quota.Constraint{}
}

func machineTypeToQuota(t string, instanceTypes map[string]InstanceTypeInfo) quota.Constraint {
	info, ok := instanceTypes[t]
	if !ok {
		return quota.Constraint{Name: "ec2/L-7295265B", Count: 0}
	}
	class := string(t[0])
	switch class {
	case "a", "c", "d", "h", "i", "m", "r", "t", "z":
		return quota.Constraint{Name: "ec2/L-1216C47A", Count: info.vCPU}
	default:
		return quota.Constraint{Name: "ec2/L-7295265B", Count: info.vCPU}
	}
}
