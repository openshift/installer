/*
Copyright 2021 Red Hat.

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

package util

import (
	"context"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/library-go/pkg/config/clusterstatus"
	"github.com/openshift/library-go/pkg/config/leaderelection"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

// GetLeaderElectionConfig returns leader election configs defaults based on the cluster topology
func GetLeaderElectionConfig(restcfg *rest.Config, leaderElection configv1.LeaderElection) configv1.LeaderElection {

	userExplicitlySetLeaderElectionValues := leaderElection.LeaseDuration.Duration != 0 ||
		leaderElection.RenewDeadline.Duration != 0 ||
		leaderElection.RetryPeriod.Duration != 0

	// Defaults follow conventions
	// https://github.com/openshift/enhancements/blob/master/CONVENTIONS.md#high-availability
	defaultLeaderElection := leaderelection.LeaderElectionDefaulting(
		leaderElection,
		"", "",
	)

	// If user has not supplied any leader election values and leader election is not disabled
	// Fetch cluster infra status to determine if we should be using SNO LE config
	if !userExplicitlySetLeaderElectionValues && !leaderElection.Disable {
		if infra, err := clusterstatus.GetClusterInfraStatus(context.TODO(), restcfg); err == nil && infra != nil {
			if infra.ControlPlaneTopology == configv1.SingleReplicaTopologyMode {
				return leaderelection.LeaderElectionSNOConfig(defaultLeaderElection)
			}
		} else {
			klog.Warningf("unable to get cluster infrastructure status, using HA cluster values for leader election: %v", err)
		}
	}

	return defaultLeaderElection
}
