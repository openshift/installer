package preprovision

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/openstack/defaults"
)

// ServerGroups uses installconfig information to pre-create server groups
// referred to in Machines.
func ServerGroups(_ context.Context, installConfig *installconfig.InstallConfig, infraID string) error {
	computeClient, err := defaults.NewServiceClient("compute", defaults.DefaultClientOpts(installConfig.Config.Platform.OpenStack.Cloud))
	if err != nil {
		return fmt.Errorf("failed to build an OpenStack client: %w", err)
	}
	computeClient.Microversion = "2.64"

	defaultPolicy := openstack.SGPolicySoftAntiAffinity
	if defaultMP := installConfig.Config.OpenStack.DefaultMachinePlatform; defaultMP != nil {
		if defaultMP.ServerGroupPolicy.IsSet() {
			defaultPolicy = defaultMP.ServerGroupPolicy
		}
	}

	for role, policy := range map[string]openstack.ServerGroupPolicy{
		"master": func() openstack.ServerGroupPolicy {
			if installConfig.Config.ControlPlane != nil {
				if installConfig.Config.ControlPlane.Platform.OpenStack != nil {
					if p := installConfig.Config.ControlPlane.Platform.OpenStack.ServerGroupPolicy; p.IsSet() {
						return p
					}
				}
			}
			return defaultPolicy
		}(),
		"worker": func() openstack.ServerGroupPolicy {
			if installConfig.Config.WorkerMachinePool() != nil {
				if installConfig.Config.WorkerMachinePool().Platform.OpenStack != nil {
					if p := installConfig.Config.WorkerMachinePool().Platform.OpenStack.ServerGroupPolicy; p.IsSet() {
						return p
					}
				}
			}
			return defaultPolicy
		}(),
	} {
		_, err := servergroups.Create(computeClient, servergroups.CreateOpts{
			Name:   infraID + "-" + role,
			Policy: string(policy),
		}).Extract()
		if err != nil {
			return fmt.Errorf("failed to create the %s server group: %w", role, err)
		}
	}
	return nil
}
