package clusterapi

import (
	"context"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/types"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

func getClusterVmGroups(ctx context.Context, vim25Client *vim25.Client, computeCluster string) ([]*types.ClusterVmGroup, error) {
	finder := find.NewFinder(vim25Client, true)

	ccr, err := finder.ClusterComputeResource(ctx, computeCluster)
	if err != nil {
		return nil, err
	}

	clusterConfig, err := ccr.Configuration(ctx)
	if err != nil {
		return nil, err
	}

	var clusterVmGroup []*types.ClusterVmGroup

	for _, g := range clusterConfig.Group {
		if vmg, ok := g.(*types.ClusterVmGroup); ok {
			clusterVmGroup = append(clusterVmGroup, vmg)
		}
	}
	return clusterVmGroup, nil
}

func createVMGroup(ctx context.Context, session *session.Session, cluster, vmGroup string) error {
	clusterObj, err := session.Finder.ClusterComputeResource(ctx, cluster)
	if err != nil {
		return err
	}

	clusterConfigSpec := &types.ClusterConfigSpecEx{
		GroupSpec: []types.ClusterGroupSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperation("add"),
				},
				Info: &types.ClusterVmGroup{
					ClusterGroupInfo: types.ClusterGroupInfo{
						Name: vmGroup,
					},
				},
			},
		},
	}

	task, err := clusterObj.Reconfigure(ctx, clusterConfigSpec, true)
	if err != nil {
		return err
	}

	return task.Wait(ctx)
}

func createVMHostAffinityRule(ctx context.Context, session *session.Session, cluster, hostGroup, vmGroup, rule string) error {
	clusterObj, err := session.Finder.ClusterComputeResource(ctx, cluster)
	if err != nil {
		return err
	}
	clusterConfigSpec := &types.ClusterConfigSpecEx{
		RulesSpec: []types.ClusterRuleSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperation("add"),
				},
				Info: &types.ClusterVmHostRuleInfo{
					ClusterRuleInfo: types.ClusterRuleInfo{
						Name:        rule,
						Mandatory:   types.NewBool(true),
						Enabled:     types.NewBool(true),
						UserCreated: types.NewBool(true),
					},
					VmGroupName:         vmGroup,
					AffineHostGroupName: hostGroup,
				},
			},
		},
	}

	task, err := clusterObj.Reconfigure(ctx, clusterConfigSpec, true)
	if err != nil {
		return err
	}

	return task.Wait(ctx)
}
