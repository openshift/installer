package aws

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	capi "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck //CORS-3563

	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	quotatypes "github.com/openshift/installer/pkg/asset/quota/types"
	"github.com/openshift/installer/pkg/quota"
	"github.com/openshift/installer/pkg/types"
)

// Constraints returns a list of quota constraints based on the InstallConfig.
// These constraints can be used to check if there is enough quota for creating a cluster
// for the install config.
func Constraints(config *types.InstallConfig, controlPlanes []quotatypes.MachineInfo, computes []quotatypes.MachineInfo, instanceTypes map[string]InstanceTypeInfo) []quota.Constraint {
	var ret []quota.Constraint
	for _, gen := range []constraintGenerator{
		network(config, append(controlPlanes, computes...)),
		controlPlane(config, controlPlanes, instanceTypes),
		compute(config, computes, instanceTypes),
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

func network(config *types.InstallConfig, machines []quotatypes.MachineInfo) func() []quota.Constraint {
	return func() []quota.Constraint {
		zones := sets.NewString()
		for _, m := range machines {
			zones.Insert(m.AvailabilityZone)
		}

		var ret []quota.Constraint
		ret = append(ret, quota.Constraint{
			Name:   "vpc/L-E79EC296", // 2 sg
			Region: config.Platform.AWS.Region,
			Count:  2,
		})
		if len(config.Platform.AWS.VPC.Subnets) == 0 {
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

func controlPlane(config *types.InstallConfig, machines []quotatypes.MachineInfo, instanceTypes map[string]InstanceTypeInfo) func() []quota.Constraint {
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

func compute(config *types.InstallConfig, machines []quotatypes.MachineInfo, instanceTypes map[string]InstanceTypeInfo) func() []quota.Constraint {
	return func() []quota.Constraint {
		var ret []quota.Constraint
		for _, m := range machines {
			q := machineTypeToQuota(m.InstanceType, instanceTypes)
			q.Count *= m.Replicas
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
		logrus.Warnf("%s", warnMessage)
		return quota.Constraint{Name: "ec2/L-7295265B", Count: 0}
	}
	r := regexp.MustCompile(`^([A-Za-z]+)[0-9]`)
	match := r.FindStringSubmatch(strings.ToLower(t))
	if match == nil {
		logrus.Warnf("%s", warnMessage)
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
		logrus.Warnf("%s", warnMessage)
		return quota.Constraint{Name: "ec2/L-7295265B", Count: 0}
	}
}

// MachineInfoFromMAPIMachines converts MAPI Machine objects to MachineInfo.
func MachineInfoFromMAPIMachines(mapiMachines []machinev1beta1.Machine) []quotatypes.MachineInfo {
	infos := make([]quotatypes.MachineInfo, 0, len(mapiMachines))
	for _, m := range mapiMachines {
		providerConfig := m.Spec.ProviderSpec.Value.Object.(*machinev1beta1.AWSMachineProviderConfig)
		infos = append(infos, quotatypes.MachineInfo{
			InstanceType:     providerConfig.InstanceType,
			AvailabilityZone: providerConfig.Placement.AvailabilityZone,
			Replicas:         1,
		})
	}
	return infos
}

// MachineInfoFromMAPIMachineSets converts MAPI MachineSet objects to MachineInfo.
func MachineInfoFromMAPIMachineSets(mapiMachineSets []machinev1beta1.MachineSet) []quotatypes.MachineInfo {
	infos := make([]quotatypes.MachineInfo, 0, len(mapiMachineSets))
	for _, ms := range mapiMachineSets {
		providerConfig := ms.Spec.Template.Spec.ProviderSpec.Value.Object.(*machinev1beta1.AWSMachineProviderConfig)
		infos = append(infos, quotatypes.MachineInfo{
			InstanceType:     providerConfig.InstanceType,
			AvailabilityZone: providerConfig.Placement.AvailabilityZone,
			Replicas:         int64(*ms.Spec.Replicas),
		})
	}
	return infos
}

// MachineInfoFromCAPIMachineSets converts CAPI MachineSet and AWSMachineTemplate objects to MachineInfo.
func MachineInfoFromCAPIMachineSets(capiMachineSets []capi.MachineSet, capiTemplates []capa.AWSMachineTemplate) []quotatypes.MachineInfo {
	templateInstanceTypes := make(map[string]string, len(capiTemplates))
	for _, t := range capiTemplates {
		templateInstanceTypes[t.Name] = t.Spec.Template.Spec.InstanceType
	}
	infos := make([]quotatypes.MachineInfo, 0, len(capiMachineSets))
	for _, ms := range capiMachineSets {
		infos = append(infos, quotatypes.MachineInfo{
			InstanceType: templateInstanceTypes[ms.Spec.Template.Spec.InfrastructureRef.Name],
			Replicas:     int64(*ms.Spec.Replicas),
		})
	}
	return infos
}
