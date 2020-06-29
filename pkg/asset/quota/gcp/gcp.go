package gcp

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	gcpprovider "github.com/openshift/cluster-api-provider-gcp/pkg/apis/gcpprovider/v1beta1"
	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"

	"github.com/openshift/installer/pkg/quota"
	"github.com/openshift/installer/pkg/types"
)

// Constraints returns a list of quota constraints based on the InstallConfig.
// These constraints can be used to check if there is enough quota for creating a cluster
// for the isntall config.
func Constraints(config *types.InstallConfig, controlPlanes []machineapi.Machine, computes []machineapi.MachineSet) []quota.Constraint {
	ctrplConfigs := make([]*gcpprovider.GCPMachineProviderSpec, len(controlPlanes))
	for i, m := range controlPlanes {
		ctrplConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*gcpprovider.GCPMachineProviderSpec)
	}
	computeReplicas := make([]int64, len(computes))
	computeConfigs := make([]*gcpprovider.GCPMachineProviderSpec, len(computes))
	for i, w := range computes {
		computeReplicas[i] = int64(*w.Spec.Replicas)
		computeConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*gcpprovider.GCPMachineProviderSpec)
	}

	var ret []quota.Constraint
	for _, gen := range []constraintGenerator{
		network(config),
		apiExternal(config),
		apiInternal(config),
		controlPlane(config, ctrplConfigs),
		compute(config, computeReplicas, computeConfigs),
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

func network(config *types.InstallConfig) func() []quota.Constraint {
	return func() []quota.Constraint {
		net := []quota.Constraint{{
			Name:   "compute.googleapis.com/networks",
			Region: "global",
			Count:  1,
		}, {
			Name:   "compute.googleapis.com/subnetworks",
			Region: "global",
			Count:  2,
		}, {
			Name:   "compute.googleapis.com/routers",
			Region: "global",
			Count:  1,
		}}

		firewalls := []quota.Constraint{{
			Name:   "compute.googleapis.com/firewalls",
			Region: "global",
			Count:  6,
		}}

		if len(config.Platform.GCP.Network) > 0 {
			return firewalls
		}
		return append(net, firewalls...)
	}
}

func apiExternal(config *types.InstallConfig) func() []quota.Constraint {
	return func() []quota.Constraint {
		if config.Publish == types.InternalPublishingStrategy {
			return nil
		}
		return []quota.Constraint{{
			Name:   "compute.googleapis.com/health_checks",
			Region: "global",
			Count:  1,
		}, {
			Name:   "compute.googleapis.com/forwarding_rules",
			Region: "global",
			Count:  1,
		}, {
			Name:   "compute.googleapis.com/target_pools",
			Region: "global",
			Count:  1,
		}, {
			Name:   "compute.googleapis.com/regional_static_addresses",
			Region: config.Platform.GCP.Region,
			Count:  1,
		}}
	}
}

func apiInternal(config *types.InstallConfig) func() []quota.Constraint {
	return func() []quota.Constraint {
		return []quota.Constraint{{
			Name:   "compute.googleapis.com/health_checks",
			Region: "global",
			Count:  1,
		}, {
			Name:   "compute.googleapis.com/forwarding_rules",
			Region: "global",
			Count:  1,
		}, {
			Name:   "compute.googleapis.com/backend_services",
			Region: "global",
			Count:  1,
		}, {
			Name:   "compute.googleapis.com/regional_static_addresses",
			Region: config.Platform.GCP.Region,
			Count:  1,
		}}
	}
}

func controlPlane(config *types.InstallConfig, machines []*gcpprovider.GCPMachineProviderSpec) func() []quota.Constraint {
	return func() []quota.Constraint {
		var ret []quota.Constraint
		for _, m := range machines {
			q := machineTypeToQuota(m.MachineType)
			q.Region = config.Platform.GCP.Region
			ret = append(ret, q)
		}

		ret = append(ret, quota.Constraint{
			Name:   "iam.googleapis.com/quota/service-account-count",
			Region: "global",
			Count:  1,
		})
		return ret
	}
}

func compute(config *types.InstallConfig, replicas []int64, machines []*gcpprovider.GCPMachineProviderSpec) func() []quota.Constraint {
	return func() []quota.Constraint {
		var ret []quota.Constraint
		for idx, m := range machines {
			q := machineTypeToQuota(m.MachineType)
			q.Count = q.Count * replicas[idx]
			q.Region = config.Platform.GCP.Region
			ret = append(ret, q)
		}

		ret = append(ret, quota.Constraint{
			Name:   "iam.googleapis.com/quota/service-account-count",
			Region: "global",
			Count:  1,
		})
		return ret
	}
}

func others() []quota.Constraint {
	return []quota.Constraint{{
		Name:   "compute.googleapis.com/images",
		Region: "global",
		Count:  1,
	}, {
		Name:   "iam.googleapis.com/quota/service-account-count",
		Region: "global",
		Count:  3,
	}}
}

func machineTypeToQuota(t string) quota.Constraint {
	var class string
	var count int64
	split := strings.Split(t, "-")
	switch len(split) {
	case 3, 4:
		class = split[0]
		if c, err := strconv.ParseInt(split[2], 10, 0); err == nil {
			count = c
		}
	case 2:
		class = split[0]
	}
	switch class {
	case "c2", "m1", "m2", "n2", "n2d":
		return quota.Constraint{Name: fmt.Sprintf("compute.googleapis.com/%s_cpus", class), Count: count}
	case "e2":
		if count == 0 {
			count = 2
		}
		return quota.Constraint{Name: "compute.googleapis.com/cpus", Count: count}
	case "f1", "g1":
		return quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 1}
	default:
		return quota.Constraint{Name: "compute.googleapis.com/cpus", Count: count}
	}
}
