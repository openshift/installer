package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/types"
	aztypes "github.com/openshift/installer/pkg/types/azure"
)

func TestGenerateMachinesWithCapacityReservationGroup(t *testing.T) {
	replicas := int64(1)
	group := &aztypes.CapacityReservationGroup{
		SubscriptionID: "reservation-subscription",
		ResourceGroup:  "reservation-resource-group",
		Name:           "reservation-group",
	}
	input := &MachineInput{
		Platform: &aztypes.Platform{},
		Pool: &types.MachinePool{
			Name:     types.MachinePoolControlPlaneRoleName,
			Replicas: &replicas,
			Platform: types.MachinePoolPlatform{
				Azure: &aztypes.MachinePool{
					Zones:                    []string{"1"},
					Identity:                 &aztypes.VMIdentity{},
					CapacityReservationGroup: group,
				},
			},
		},
	}

	files, err := GenerateMachines("test-cluster", "cluster-resource-group", "cluster-subscription", nil, input)
	assert.NoError(t, err)

	var azureMachines []*capz.AzureMachine
	for _, file := range files {
		if machine, ok := file.Object.(*capz.AzureMachine); ok {
			azureMachines = append(azureMachines, machine)
		}
	}

	if assert.Len(t, azureMachines, 2) {
		for _, machine := range azureMachines {
			if assert.NotNil(t, machine.Spec.CapacityReservationGroupID) {
				assert.Equal(t, group.ToID(), *machine.Spec.CapacityReservationGroupID)
			}
		}
	}
}
