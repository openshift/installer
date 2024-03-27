package preprovision

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/gophercloud/gophercloud/pagination"
	mapov1alpha1 "github.com/openshift/api/machine/v1alpha1"
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
func ServerGroups(_ context.Context, installConfig *installconfig.InstallConfig, capiMachines []capo.OpenStackMachine, mapoWorkerProviderSpecs []mapov1alpha1.OpenstackProviderSpec) error {
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
	serverGroups := make(map[string]openstack.ServerGroupPolicy, len(capiMachines)+len(mapoWorkerProviderSpecs))
	for _, machine := range capiMachines {
		var policy openstack.ServerGroupPolicy
		if _, ok := machine.Labels["cluster.x-k8s.io/control-plane"]; ok {
			fmt.Println("MYDEBUG: machine", machine.GetName(), "identifies as MASTER")
			policy = tfvars.GetServerGroupPolicy(masterMP, installConfig.Config.OpenStack.DefaultMachinePlatform)
		} else {
			fmt.Println("MYDEBUG: machine", machine.GetName(), "identifies as WORKER")
			// This IF branch is never executed as long as
			// Installer only generates CAPI manifests for masters
			// and bootstrap
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
	for _, providerSpec := range mapoWorkerProviderSpecs {
		var policy openstack.ServerGroupPolicy
		if role := providerSpec.Labels["machine.openshift.io/cluster-api-machine-role"]; role == "worker" {
			fmt.Println("MYDEBUG: machine", providerSpec.GetName(), "identifies as WORKER")
			policy = tfvars.GetServerGroupPolicy(masterMP, installConfig.Config.OpenStack.DefaultMachinePlatform)
		} else {
			fmt.Println("MYDEBUG: providerSpec", providerSpec.GetName(), "has role:", role)
		}

		if name := providerSpec.ServerGroupName; name != "" {
			if visitedPolicy, ok := serverGroups[name]; ok {
				if policy != visitedPolicy {
					return fmt.Errorf("server group %q is referenced with different policies in the install-config machine-pools", name)
				}
			} else {
				serverGroups[name] = policy
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
