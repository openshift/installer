package validation

import (
	"fmt"
	"strings"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

type platformValidation func(*openstack.Platform, *types.InstallConfig) error

func validateFailureDomainNamesNotEmpty(platform *openstack.Platform, _ *types.InstallConfig) error {
	for _, failureDomain := range platform.FailureDomains {
		if failureDomain.Name == "" {
			return fmt.Errorf("must specify a failure domain name")
		}
	}

	return nil
}

func validateFailureDomainNamesUnique(platform *openstack.Platform, _ *types.InstallConfig) error {
	var (
		names      = map[string]struct{}{}
		duplicates = []string{}
	)

	for _, failureDomain := range platform.FailureDomains {
		if failureDomain.Name != "" {
			if _, ok := names[failureDomain.Name]; ok {
				duplicates = append(duplicates, failureDomain.Name)
			} else {
				names[failureDomain.Name] = struct{}{}
			}
		}
	}

	if len(duplicates) > 0 {
		return fmt.Errorf("failure domain names must be unique. Found duplicates: %s", strings.Join(duplicates, ", "))
	}

	return nil
}

func validateFailureDomainMachinesSubnetDependency(platform *openstack.Platform, _ *types.InstallConfig) error {
	if platform.MachinesSubnet == "" {
		for _, failureDomain := range platform.FailureDomains {
			if failureDomain.Subnet != "" {
				return fmt.Errorf("must specify a machinesSubnet when failure domain subnets are specified")
			}
		}
	}

	return nil
}
