package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBaseDomain(t *testing.T) {
	platform := Platform{}
	zoneID := "/subscriptions/<subid>/resourceGroups/<rg_name>/providers/Microsoft.Network/dnszones/<zone_name>"
	platform.SetBaseDomain(zoneID)
	assert.Equal(t, "<rg_name>", platform.BaseDomainResourceGroupName)
}
