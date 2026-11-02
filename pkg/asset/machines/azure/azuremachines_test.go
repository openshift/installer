package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/types"
	aztypes "github.com/openshift/installer/pkg/types/azure"
)

func TestGenerateMachinesWithCapacityReservationGroup(t *testing.T) {
	testCases := []struct {
		name       string
		group      *aztypes.CapacityReservationGroup
		expectedID string
	}{
		{
			name: "subscription specified",
			group: &aztypes.CapacityReservationGroup{
				SubscriptionID: "reservation-subscription",
				ResourceGroup:  "reservation-resource-group",
				Name:           "reservation-group",
			},
			expectedID: "/subscriptions/reservation-subscription/resourceGroups/reservation-resource-group/providers/Microsoft.Compute/capacityReservationGroups/reservation-group",
		},
		{
			name: "subscription omitted falls back to the cluster subscription",
			group: &aztypes.CapacityReservationGroup{
				ResourceGroup: "reservation-resource-group",
				Name:          "reservation-group",
			},
			expectedID: "/subscriptions/cluster-subscription/resourceGroups/reservation-resource-group/providers/Microsoft.Compute/capacityReservationGroups/reservation-group",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalSubscriptionID := tc.group.SubscriptionID
			replicas := int64(1)
			input := &MachineInput{
				Platform: &aztypes.Platform{},
				Pool: &types.MachinePool{
					Name:     types.MachinePoolControlPlaneRoleName,
					Replicas: &replicas,
					Platform: types.MachinePoolPlatform{
						Azure: &aztypes.MachinePool{
							Zones:                    []string{"1"},
							Identity:                 &aztypes.VMIdentity{},
							CapacityReservationGroup: tc.group,
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
						assert.Equal(t, tc.expectedID, *machine.Spec.CapacityReservationGroupID)
					}
				}
			}

			// The group is shared with the install-config asset, so the
			// fallback must not write back to it.
			assert.Equal(t, originalSubscriptionID, tc.group.SubscriptionID)
		})
	}
}

func TestGenerateMachinesWithoutCapacityReservationGroup(t *testing.T) {
	replicas := int64(1)
	input := &MachineInput{
		Platform: &aztypes.Platform{},
		Pool: &types.MachinePool{
			Name:     types.MachinePoolControlPlaneRoleName,
			Replicas: &replicas,
			Platform: types.MachinePoolPlatform{
				Azure: &aztypes.MachinePool{
					Zones:    []string{"1"},
					Identity: &aztypes.VMIdentity{},
				},
			},
		},
	}

	files, err := GenerateMachines("test-cluster", "cluster-resource-group", "cluster-subscription", nil, input)
	assert.NoError(t, err)

	azureMachineCount := 0
	for _, file := range files {
		machine, ok := file.Object.(*capz.AzureMachine)
		if !ok {
			continue
		}
		azureMachineCount++
		assert.Nil(t, machine.Spec.CapacityReservationGroupID)
	}
	assert.Equal(t, 2, azureMachineCount)
}
