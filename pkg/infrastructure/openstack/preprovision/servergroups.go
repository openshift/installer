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

	var serverGroups map[string]openstack.ServerGroupPolicy
	{
		var openstackWorkerMachinePool *openstack.MachinePool
		{
			if installConfig.Config.WorkerMachinePool() != nil && installConfig.Config.WorkerMachinePool().Platform.OpenStack != nil {
				openstackWorkerMachinePool = installConfig.Config.WorkerMachinePool().Platform.OpenStack
			} else {
				openstackWorkerMachinePool = installConfig.Config.Platform.OpenStack.DefaultMachinePlatform
			}
		}

		policy := openstack.SGPolicySoftAffinity
		if p := openstackWorkerMachinePool.ServerGroupPolicy; p.IsSet() {
			policy = p
		}

		serverGroups = make(map[string]openstack.ServerGroupPolicy, len(openstackWorkerMachinePool.Zones)+1)
		for _, zone := range openstackWorkerMachinePool.Zones {
			if zone == "" {
				serverGroups["worker"] = policy
			} else {
				serverGroups["worker-"+zone] = policy
			}
		}
	}

	{
		var openstackMasterMachinePool *openstack.MachinePool
		if installConfig.Config.ControlPlane != nil && installConfig.Config.ControlPlane.Platform.OpenStack != nil {
			openstackMasterMachinePool = installConfig.Config.ControlPlane.Platform.OpenStack
		} else {
			openstackMasterMachinePool = installConfig.Config.Platform.OpenStack.DefaultMachinePlatform
		}

		policy := openstack.SGPolicySoftAffinity
		if p := openstackMasterMachinePool.ServerGroupPolicy; p.IsSet() {
			policy = p
		}

		serverGroups["master"] = policy
	}

	for role, policy := range serverGroups {
		name := infraID + "-" + role
		_, err := servergroups.Create(computeClient, servergroups.CreateOpts{
			Name:   name,
			Policy: string(policy),
		}).Extract()
		if err != nil {
			return fmt.Errorf("failed to create server group %q: %w", name, err)
		}
	}
	return nil
}
