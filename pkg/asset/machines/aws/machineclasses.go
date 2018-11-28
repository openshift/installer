// Package aws generates Machine objects for aws.
package aws

import (
	"fmt"

	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/pkg/errors"
)

// MachineClasses returns a list of MachineClasses for a machinepool.
func MachineClasses(config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string) ([]clusterapi.MachineClass, error) {
	if configPlatform := config.Platform.Name(); configPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS machine-pool: %q", poolPlatform)
	}
	clusterName := config.ObjectMeta.Name
	platform := config.Platform.AWS
	mpool := pool.Platform.AWS

	var machineClasses []clusterapi.MachineClass
	for idx := range mpool.Zones {
		machineClass, err := machineClass(config.ClusterID, clusterName, platform, mpool, idx, role, userDataSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create machineClass")
		}
		machineClasses = append(machineClasses, *machineClass)
	}

	return machineClasses, nil
}
