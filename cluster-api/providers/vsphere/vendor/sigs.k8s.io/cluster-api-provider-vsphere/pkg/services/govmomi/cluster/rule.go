/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

import (
	"github.com/pkg/errors"
	"github.com/vmware/govmomi/vim25/types"
	"k8s.io/utils/pointer"
)

type Rule interface {
	Disabled() bool

	IsMandatory() bool
}

type vmHostAffinityRule struct {
	*types.ClusterVmHostRuleInfo
}

func (v vmHostAffinityRule) IsMandatory() bool {
	return pointer.BoolDeref(v.Mandatory, false)
}

func (v vmHostAffinityRule) Disabled() bool {
	if v.Enabled == nil {
		return true
	}
	return negate(*v.Enabled)
}

func negate(input bool) bool {
	return !input
}

func VerifyAffinityRule(ctx computeClusterContext, clusterName, hostGroupName, vmGroupName string) (Rule, error) {
	rules, err := listRules(ctx, clusterName)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to list rules for compute cluster %s", clusterName)
	}

	for _, rule := range rules {
		if vmHostRuleInfo, ok := rule.(*types.ClusterVmHostRuleInfo); ok {
			if vmHostRuleInfo.AffineHostGroupName == hostGroupName &&
				vmHostRuleInfo.VmGroupName == vmGroupName {
				return vmHostAffinityRule{vmHostRuleInfo}, nil
			}
		}
	}
	return nil, errors.New("no matching affinity rule found/exists")
}

func listRules(ctx computeClusterContext, clusterName string) ([]types.BaseClusterRuleInfo, error) {
	ccr, err := ctx.GetSession().Finder.ClusterComputeResource(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	clusterConfigInfoEx, err := ccr.Configuration(ctx)
	if err != nil {
		return nil, err
	}
	return clusterConfigInfoEx.Rule, nil
}
