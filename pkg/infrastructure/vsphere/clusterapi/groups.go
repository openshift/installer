package clusterapi

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/vim25/types"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

func createVMGroup(ctx context.Context, session *session.Session, cluster, vmGroup string) error {
	clusterObj, err := session.Finder.ClusterComputeResource(ctx, cluster)
	if err != nil {
		logrus.Errorf("unable to find cluster compute resource: %v", err)
		return nil
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
		logrus.Errorf("unable to reconfigure cluster: %v", err)
		return nil
	}

	err = task.Wait(ctx)
	if err != nil {
		logrus.Errorf("unable to wait for cluster reconfiguration task: %v", err)
		return nil
	}
	return nil
}

func createVMHostAffinityRule(ctx context.Context, session *session.Session, cluster, hostGroup, vmGroup, rule string) error {
	clusterObj, err := session.Finder.ClusterComputeResource(ctx, cluster)
	if err != nil {
		logrus.Errorf("unable to find cluster compute resource: %v", err)
		return nil
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
		logrus.Errorf("unable to reconfigure cluster: %v", err)
		return nil
	}

	err = task.Wait(ctx)
	if err != nil {
		logrus.Errorf("unable to wait for cluster reconfiguration task: %v", err)
		return nil
	}
	return nil
}
