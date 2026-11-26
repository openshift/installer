package azure

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/openshift/installer/pkg/asset/installconfig/azure/mock"
)

func TestVMAvailabilityZonesPerInstanceType(t *testing.T) {
	const (
		region       = "centralus"
		zonedType    = "Standard_D8s_v5"
		unzonedType  = "Standard_Unzoned"
		narrowedType = "Standard_D2s_v5"
	)

	mockCtrl := gomock.NewController(t)
	azureClient := mock.NewMockAPI(mockCtrl)
	// Each instance type is looked up exactly once, and each gets its own
	// answer rather than the first one cached.
	azureClient.EXPECT().GetAvailabilityZones(gomock.Any(), region, zonedType).Return([]string{"3", "1", "2"}, nil).Times(1)
	azureClient.EXPECT().GetAvailabilityZones(gomock.Any(), region, narrowedType).Return([]string{"2", "1"}, nil).Times(1)
	azureClient.EXPECT().GetAvailabilityZones(gomock.Any(), region, unzonedType).Return(nil, nil).Times(1)

	meta := &Metadata{Region: region, client: azureClient}

	for range 2 {
		zones, err := meta.VMAvailabilityZones(context.Background(), zonedType)
		assert.NoError(t, err)
		assert.Equal(t, []string{"1", "2", "3"}, zones, "zones should be sorted")

		zones, err = meta.VMAvailabilityZones(context.Background(), narrowedType)
		assert.NoError(t, err)
		assert.Equal(t, []string{"1", "2"}, zones)

		zones, err = meta.VMAvailabilityZones(context.Background(), unzonedType)
		assert.NoError(t, err)
		assert.Empty(t, zones)
	}
}

func TestVMAvailabilityZonesError(t *testing.T) {
	const region = "centralus"

	mockCtrl := gomock.NewController(t)
	azureClient := mock.NewMockAPI(mockCtrl)
	azureClient.EXPECT().GetAvailabilityZones(gomock.Any(), region, "Standard_Broken").
		Return(nil, fmt.Errorf("boom")).AnyTimes()

	meta := &Metadata{Region: region, client: azureClient}

	zones, err := meta.VMAvailabilityZones(context.Background(), "Standard_Broken")
	assert.Nil(t, zones)
	assert.ErrorContains(t, err, "error retrieving Availability Zones: boom")
}
