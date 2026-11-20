package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapacityReservationGroupToID(t *testing.T) {
	group := &CapacityReservationGroup{
		SubscriptionID: "00000000-0000-0000-0000-000000000000",
		ResourceGroup:  "reservation-rg",
		Name:           "reservation-group",
	}

	assert.Equal(t,
		"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/reservation-rg/providers/Microsoft.Compute/capacityReservationGroups/reservation-group",
		group.ToID())
}
