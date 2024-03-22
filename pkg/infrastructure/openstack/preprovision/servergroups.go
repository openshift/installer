package preprovision

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/sirupsen/logrus"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"

	"github.com/openshift/installer/pkg/asset/installconfig"
	tfvars "github.com/openshift/installer/pkg/tfvars/openstack"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/openstack/defaults"
)

// ServerGroups creates server groups referenced by name in the Machine
// manifests if they don't exist already. The newly created server groups have
// the policy defined in the install-config's machine-pools.
func ServerGroups(_ context.Context, installConfig *installconfig.InstallConfig, machineManifests []capo.OpenStackMachine) error {
	logrus.Debugf("Creating the server groups")
	computeClient, err := defaults.NewServiceClient("compute", defaults.DefaultClientOpts(installConfig.Config.Platform.OpenStack.Cloud))
	if err != nil {
		return fmt.Errorf("failed to build an OpenStack client: %w", err)
	}
	computeClient.Microversion = "2.64"

	var masterMP, workerMP *openstack.MachinePool
	{
		if mp := installConfig.Config.ControlPlane; mp != nil {
			masterMP = mp.Platform.OpenStack
		}
		if mp := installConfig.Config.WorkerMachinePool(); mp != nil {
			workerMP = mp.Platform.OpenStack
		}
	}

	// serverGroups is the set of server groups to be created.
	serverGroups := make(map[string]openstack.ServerGroupPolicy, len(machineManifests))
	for _, machine := range machineManifests {
		fmt.Println("DEBUG: machine", machine.GetName(), "roles:", machine.GetAnnotations())
		var policy openstack.ServerGroupPolicy
		if _, ok := machine.Labels["cluster.x-k8s.io/control-plane"]; ok {
			policy = tfvars.GetServerGroupPolicy(masterMP, installConfig.Config.OpenStack.DefaultMachinePlatform)
		} else {
			policy = tfvars.GetServerGroupPolicy(workerMP, installConfig.Config.OpenStack.DefaultMachinePlatform)
		}

		if sgFilter := machine.Spec.ServerGroup; sgFilter != nil {
			if name := sgFilter.Name; name != "" {
				if visitedPolicy, ok := serverGroups[name]; ok {
					if policy != visitedPolicy {
						return fmt.Errorf("server group %q is referenced with different policies in the install-config machine-pools", name)
					}
				} else {
					serverGroups[name] = policy
				}
			}
		}
	}

	// Remove existing server groups from serverGroups.
	servergroups.List(computeClient, nil).EachPage(func(p pagination.Page) (bool, error) {
		sgs, err := servergroups.ExtractServerGroups(p)
		if err != nil {
			return false, err
		}

		for i := range sgs {
			delete(serverGroups, sgs[i].Name)
		}
		return true, nil
	})

	// Create the server groups referenced by name in the Machine manifests, that don't exist already.
	for name, policy := range serverGroups {
		logrus.Debugf("Creating server group %q with policy %q", name, policy)
		if _, err := servergroups.Create(computeClient, servergroups.CreateOpts{
			Name:   name,
			Policy: string(policy),
		}).Extract(); err != nil {
			return fmt.Errorf("failed to create server group %q: %w", name, err)
		}
	}
	return nil
}
