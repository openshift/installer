package instance

import (
	"fmt"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_tenants"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPITenantClient ...
type IBMPITenantClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

// NewIBMPITenantClient ...
func NewIBMPITenantClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPITenantClient {
	return &IBMPITenantClient{
		session:         sess,
		powerinstanceid: powerinstanceid,
	}
}

// Get ..
func (f *IBMPITenantClient) Get(tenantid string) (*models.Tenant, error) {
	params := p_cloud_tenants.NewPcloudTenantsGetParams().WithTenantID(f.session.UserAccount).WithTenantID(tenantid)
	resp, err := f.session.Power.PCloudTenants.PcloudTenantsGet(params, ibmpisession.NewAuth(f.session, tenantid))

	if err != nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to perform get operation... %s", err)
	}
	return resp.Payload, nil
}
