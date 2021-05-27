package instance

import (
	"fmt"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_system_pools"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPISystemPoolClient ...
type IBMPISystemPoolClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

// NewIBMPISystemPoolClient ...
func NewIBMPISystemPoolClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPISystemPoolClient {
	return &IBMPISystemPoolClient{
		sess, powerinstanceid,
	}
}

//Get the System Pools
func (f *IBMPISystemPoolClient) Get(powerinstanceid string) (models.SystemPools, error) {
	params := p_cloud_system_pools.NewPcloudSystempoolsGetParamsWithTimeout(f.session.Timeout).WithCloudInstanceID(powerinstanceid)
	resp, err := f.session.Power.PCloudSystemPools.PcloudSystempoolsGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to perform get operation... %s", err)
	}
	return resp.Payload, nil
}
