package instance

import (
	"fmt"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_instances"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPICloudInstanceClient ...
type IBMPICloudInstanceClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

// NewIBMPICloudInstanceClient ...
func NewIBMPICloudInstanceClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPICloudInstanceClient {
	return &IBMPICloudInstanceClient{
		session:         sess,
		powerinstanceid: powerinstanceid,
	}
}

// Get information about a cloud instance
func (f *IBMPICloudInstanceClient) Get(powerinstanceid string) (*models.CloudInstance, error) {
	params := p_cloud_instances.NewPcloudCloudinstancesGetParams().WithCloudInstanceID(powerinstanceid)
	resp, err := f.session.Power.PCloudInstances.PcloudCloudinstancesGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return nil, fmt.Errorf("Failed to Get Cloud Instance %s: %s", powerinstanceid, err)
	}
	return resp.Payload, nil
}

// Update a cloud instance
func (f *IBMPICloudInstanceClient) Update(powerinstanceid string, updateparams *p_cloud_instances.PcloudCloudinstancesPutParams) (*models.CloudInstance, error) {
	params := p_cloud_instances.NewPcloudCloudinstancesPutParamsWithTimeout(f.session.Timeout).WithCloudInstanceID(powerinstanceid).WithBody(updateparams.Body)
	resp, err := f.session.Power.PCloudInstances.PcloudCloudinstancesPut(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return nil, fmt.Errorf("Failed to Update Cloud Instance %s: %s", powerinstanceid, err)

	}
	return resp.Payload, nil
}

// Delete a Cloud instance
func (f *IBMPICloudInstanceClient) Delete(powerinstanceid string) (models.Object, error) {
	params := p_cloud_instances.NewPcloudCloudinstancesDeleteParams().WithCloudInstanceID(powerinstanceid)
	resp, err := f.session.Power.PCloudInstances.PcloudCloudinstancesDelete(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Delete Cloud Instance %s: %s", powerinstanceid, err)
	}
	return resp.Payload, nil
}
