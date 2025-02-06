package aws

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/quota"
	"github.com/openshift/installer/pkg/types"
)

// Constraints returns a list of quota constraints based on the InstallConfig.
// These constraints can be used to check if there is enough quota for creating a cluster
// for the isntall config.
func Constraints(config *types.InstallConfig, controlPlanes []machineapi.Machine, computes []machineapi.MachineSet, instanceTypes map[string]InstanceTypeInfo) []quota.Constraint {
	ctrplConfigs := make([]*machineapi.AWSMachineProviderConfig, len(controlPlanes))
	for i, m := range controlPlanes {
		ctrplConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*machineapi.AWSMachineProviderConfig)
	}
	computeReplicas := make([]int64, len(computes))
	computeConfigs := make([]*machineapi.AWSMachineProviderConfig, len(computes))
	for i, w := range computes {
		computeReplicas[i] = int64(*w.Spec.Replicas)
		computeConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*machineapi.AWSMachineProviderConfig)
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

func network(config *types.InstallConfig, machines []*machineapi.AWSMachineProviderConfig) func() []quota.Constraint {
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
		if len(config.Platform.AWS.DeprecatedSubnets) == 0 {
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

func controlPlane(config *types.InstallConfig, machines []*machineapi.AWSMachineProviderConfig, instanceTypes map[string]InstanceTypeInfo) func() []quota.Constraint {
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

func compute(config *types.InstallConfig, replicas []int64, machines []*machineapi.AWSMachineProviderConfig, instanceTypes map[string]InstanceTypeInfo) func() []quota.Constraint {
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
	warnMessage := fmt.Sprintf("The instance class is unknown for the instance type %q. The vCPU quota check will be skipped.", t)
	if !ok {
		logrus.Warnf(warnMessage)
		return quota.Constraint{Name: "ec2/L-7295265B", Count: 0}
	}
	r := regexp.MustCompile(`^([A-Za-z]+)[0-9]`)
	match := r.FindStringSubmatch(strings.ToLower(t))
	if match == nil {
		logrus.Warnf(warnMessage)
		return quota.Constraint{Name: "ec2/L-7295265B", Count: 0}
	}
	switch match[1] {
	case "a", "c", "d", "h", "i", "is", "im", "m", "r", "t", "z":
		return quota.Constraint{Name: "ec2/L-1216C47A", Count: info.vCPU}
	case "g", "vt":
		return quota.Constraint{Name: "ec2/L-DB2E81BA", Count: info.vCPU}
	case "x":
		return quota.Constraint{Name: "ec2/L-7295265B", Count: info.vCPU}
	default:
		logrus.Warnf(warnMessage)
		return quota.Constraint{Name: "ec2/L-7295265B", Count: 0}
	}
}
