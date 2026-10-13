package preprovision

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servergroups"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/sirupsen/logrus"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"

	mapov1alpha1 "github.com/openshift/api/machine/v1alpha1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	tfvars "github.com/openshift/installer/pkg/tfvars/openstack"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/openstack/defaults"
)

// ServerGroups creates server groups referenced by name in the Machine
// manifests if they don't exist already. The newly created server groups have
// the policy defined in the install-config's machine-pools.
func ServerGroups(ctx context.Context, installConfig *installconfig.InstallConfig, capiMachines []capo.OpenStackMachine, mapoWorkerProviderSpecs []mapov1alpha1.OpenstackProviderSpec) error {
	logrus.Debugf("Creating the server groups")
	computeClient, err := defaults.NewServiceClient(ctx, "compute", defaults.DefaultClientOpts(installConfig.Config.Platform.OpenStack.Cloud))
	if err != nil {
		return fmt.Errorf("failed to build an OpenStack client: %w", err)
	}
	computeClient.Microversion = "2.64"

	var masterPolicy, workerPolicy openstack.ServerGroupPolicy
	{
		if masterMP := installConfig.Config.ControlPlane; masterMP != nil {
			masterPolicy = tfvars.GetServerGroupPolicy(masterMP.Platform.OpenStack, installConfig.Config.OpenStack.DefaultMachinePlatform)
		} else {
			masterPolicy = tfvars.GetServerGroupPolicy(nil, installConfig.Config.OpenStack.DefaultMachinePlatform)
		}

		if workerMP := installConfig.Config.WorkerMachinePool(); workerMP != nil {
			workerPolicy = tfvars.GetServerGroupPolicy(workerMP.Platform.OpenStack, installConfig.Config.OpenStack.DefaultMachinePlatform)
		} else {
			workerPolicy = tfvars.GetServerGroupPolicy(nil, installConfig.Config.OpenStack.DefaultMachinePlatform)
		}
	}

	// serverGroups is the set of server groups to be created.
	serverGroups := make(map[string]openstack.ServerGroupPolicy, len(capiMachines)+len(mapoWorkerProviderSpecs))
	for _, machine := range capiMachines {
		if _, ok := machine.Labels["cluster.x-k8s.io/control-plane"]; !ok {
			logrus.Debugf("Found unexpected machine %q among the CAPI Machine manifests", machine.GetName())
			continue
		}
		if sgParam := machine.Spec.ServerGroup; sgParam != nil && sgParam.Filter != nil && sgParam.Filter.Name != nil {
			serverGroupName := *sgParam.Filter.Name
			if visitedPolicy, ok := serverGroups[serverGroupName]; ok {
				if masterPolicy != visitedPolicy {
					return fmt.Errorf("server group %q is referenced with different policies in the install-config machine-pools", serverGroupName)
				}
				continue
			}
			serverGroups[serverGroupName] = masterPolicy
		}
	}
	for _, providerSpec := range mapoWorkerProviderSpecs {
		if serverGroupName := providerSpec.ServerGroupName; serverGroupName != "" {
			if visitedPolicy, ok := serverGroups[serverGroupName]; ok {
				if workerPolicy != visitedPolicy {
					return fmt.Errorf("server group %q is referenced with different policies in the install-config machine-pools", serverGroupName)
				}
				continue
			}
			serverGroups[serverGroupName] = workerPolicy
		}
	}

	// Remove existing server groups from the list of resources to be created.
	if err = servergroups.List(computeClient, nil).EachPage(ctx, func(_ context.Context, p pagination.Page) (bool, error) {
		sgs, err := servergroups.ExtractServerGroups(p)
		if err != nil {
			return false, err
		}

		for i := range sgs {
			delete(serverGroups, sgs[i].Name)
		}
		return true, nil
	}); err != nil {
		return fmt.Errorf("failed to list server groups: %w", err)
	}

	// Create the server groups referenced by name in the Machine manifests, that don't exist already.
	for name, policy := range serverGroups {
		logrus.Debugf("Creating server group %q with policy %q", name, policy)
		if _, err := servergroups.Create(ctx, computeClient, servergroups.CreateOpts{
			Name:   name,
			Policy: string(policy),
		}).Extract(); err != nil {
			return fmt.Errorf("failed to create server group %q: %w", name, err)
		}
	}
	return nil
}
