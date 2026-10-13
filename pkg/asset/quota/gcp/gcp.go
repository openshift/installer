package gcp

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/quota"
	"github.com/openshift/installer/pkg/types"
)

// Constraints returns a list of quota constraints based on the InstallConfig.
// These constraints can be used to check if there is enough quota for creating a cluster
// for the isntall config.
func Constraints(client *Client, config *types.InstallConfig, controlPlanes []machineapi.Machine, computes []machineapi.MachineSet) []quota.Constraint {
	ctrplConfigs := make([]*machineapi.GCPMachineProviderSpec, len(controlPlanes))
	for i, m := range controlPlanes {
		ctrplConfigs[i] = m.Spec.ProviderSpec.Value.Object.(*machineapi.GCPMachineProviderSpec)
	}
	computeReplicas := make([]int64, len(computes))
	computeConfigs := make([]*machineapi.GCPMachineProviderSpec, len(computes))
	for i, w := range computes {
		computeReplicas[i] = int64(*w.Spec.Replicas)
		computeConfigs[i] = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*machineapi.GCPMachineProviderSpec)
	}

	var ret []quota.Constraint
	for _, gen := range []constraintGenerator{
		network(config),
		apiExternal(config),
		apiInternal(config),
		controlPlane(client, config, ctrplConfigs),
		compute(client, config, computeReplicas, computeConfigs),
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

func controlPlane(client MachineTypeGetter, config *types.InstallConfig, machines []*machineapi.GCPMachineProviderSpec) func() []quota.Constraint {
	return func() []quota.Constraint {
		var ret []quota.Constraint
		for _, m := range machines {
			q := machineTypeToQuota(client, m.Zone, m.MachineType)
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

func compute(client MachineTypeGetter, config *types.InstallConfig, replicas []int64, machines []*machineapi.GCPMachineProviderSpec) func() []quota.Constraint {
	return func() []quota.Constraint {
		var ret []quota.Constraint
		for idx, m := range machines {
			q := machineTypeToQuota(client, m.Zone, m.MachineType)
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

func machineTypeToQuota(client MachineTypeGetter, zone string, machineType string) quota.Constraint {
	var name string
	class := strings.SplitN(machineType, "-", 2)[0]
	switch class {
	case "c2", "m1", "m2", "n2", "n2d":
		name = fmt.Sprintf("compute.googleapis.com/%s_cpus", class)
	default:
		name = "compute.googleapis.com/cpus"
	}

	info, err := client.GetMachineType(zone, machineType)
	if err != nil {
		return quota.Constraint{Name: name, Count: guessMachineCPUCount(machineType)}
	}
	return quota.Constraint{Name: name, Count: info.GuestCpus}
}

// the guess is based on https://cloud.google.com/compute/docs/machine-types
func guessMachineCPUCount(machineType string) int64 {
	split := strings.Split(machineType, "-")
	switch len(split) {
	case 4:
		if c, err := strconv.ParseInt(split[2], 10, 0); err == nil {
			return c
		}
	case 3:
		switch split[0] {
		case "c2", "m1", "m2", "n1", "n2", "n2d", "e2", "g2":
			if c, err := strconv.ParseInt(split[2], 10, 0); err == nil {
				return c
			}
		}
	case 2:
		switch split[0] {
		case "e2":
			return 2
		case "f1", "g1":
			return 1
		}
	}
	return 0
}
