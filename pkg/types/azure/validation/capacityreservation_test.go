package validation

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/azure"
)

func TestValidateCapacityReservationGroup(t *testing.T) {
	valid := func() *azure.CapacityReservationGroup {
		return &azure.CapacityReservationGroup{
			SubscriptionID: subscriptionID,
			ResourceGroup:  resourceGroup,
			Name:           "reservation-group",
		}
	}

	cases := []struct {
		name      string
		group     *azure.CapacityReservationGroup
		cloudName azure.CloudEnvironment
		expected  string
	}{
		{name: "valid", group: valid(), cloudName: azure.PublicCloud},
		{
			name: "valid one-character name",
			group: func() *azure.CapacityReservationGroup {
				group := valid()
				group.Name = "a"
				return group
			}(),
			cloudName: azure.PublicCloud,
		},
		{
			name: "valid 80-character name",
			group: func() *azure.CapacityReservationGroup {
				group := valid()
				group.Name = strings.Repeat("a", 80)
				return group
			}(),
			cloudName: azure.PublicCloud,
		},
		{
			name:      "unsupported on Azure Stack",
			group:     valid(),
			cloudName: azure.StackCloud,
			expected:  `capacityReservationGroup: Forbidden: capacity reservation groups are not supported on Azure Stack`,
		},
		{
			name: "invalid subscription",
			group: func() *azure.CapacityReservationGroup {
				group := valid()
				group.SubscriptionID = "invalid"
				return group
			}(),
			cloudName: azure.PublicCloud,
			expected:  `capacityReservationGroup.subscriptionId: Invalid value: "invalid": invalid subscription ID format`,
		},
		{
			name: "invalid resource group",
			group: func() *azure.CapacityReservationGroup {
				group := valid()
				group.ResourceGroup = ""
				return group
			}(),
			cloudName: azure.PublicCloud,
			expected:  `capacityReservationGroup.resourceGroup: Invalid value: "": invalid resource group format`,
		},
		{
			name: "invalid name",
			group: func() *azure.CapacityReservationGroup {
				group := valid()
				group.Name = ""
				return group
			}(),
			cloudName: azure.PublicCloud,
			expected:  `capacityReservationGroup.name: Invalid value: "": invalid capacity reservation group name format`,
		},
		{
			name: "invalid 81-character name",
			group: func() *azure.CapacityReservationGroup {
				group := valid()
				group.Name = strings.Repeat("a", 81)
				return group
			}(),
			cloudName: azure.PublicCloud,
			expected:  fmt.Sprintf(`capacityReservationGroup.name: Invalid value: %q: invalid capacity reservation group name format`, strings.Repeat("a", 81)),
		},
		{
			name: "invalid name character",
			group: func() *azure.CapacityReservationGroup {
				group := valid()
				group.Name = "reservation/group"
				return group
			}(),
			cloudName: azure.PublicCloud,
			expected:  `capacityReservationGroup.name: Invalid value: "reservation/group": invalid capacity reservation group name format`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateCapacityReservationGroup(tc.group, tc.cloudName, field.NewPath("capacityReservationGroup")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, fmt.Sprintf("^%s$", tc.expected), err)
			}
		})
	}
}
