package instance

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/power-go-client/errors"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/datacenters"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPIDatacentersClient
type IBMPIDatacentersClient struct {
	IBMPIClient
}

// NewIBMPIDatacenterClient
func NewIBMPIDatacenterClient(ctx context.Context, sess *ibmpisession.IBMPISession, cloudInstanceID string) *IBMPIDatacentersClient {
	return &IBMPIDatacentersClient{
		*NewIBMPIClient(ctx, sess, cloudInstanceID),
	}
}
func (f *IBMPIDatacentersClient) Get(datacenterRegion string) (*models.Datacenter, error) {
	params := datacenters.NewV1DatacentersGetParams().WithContext(f.ctx).WithTimeout(helpers.PIGetTimeOut).WithDatacenterRegion(datacenterRegion)
	resp, err := f.session.Power.Datacenters.V1DatacentersGet(params)
	if err != nil {
		return nil, ibmpisession.SDKFailWithAPIError(err, fmt.Errorf(errors.GetDatacenterOperationFailed, f.cloudInstanceID, err))
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get Datacenter %s", datacenterRegion)
	}
	return resp.Payload, nil
}

func (f *IBMPIDatacentersClient) GetAll() (*models.Datacenters, error) {
	params := datacenters.NewV1DatacentersGetallParams().WithContext(f.ctx).WithTimeout(helpers.PICreateTimeOut)
	resp, err := f.session.Power.Datacenters.V1DatacentersGetall(params)

	if err != nil {
		return nil, ibmpisession.SDKFailWithAPIError(err, fmt.Errorf("failed to Get all Datacenters: %w", err))
	}
	if resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("failed to Get all Datacenters")
	}
	return resp.Payload, nil
}
