package instance

import (
	"fmt"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_tenants_ssh_keys"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPIKeyClient ...
type IBMPIKeyClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

// NewIBMPIKeyClient ...
func NewIBMPIKeyClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPIKeyClient {
	return &IBMPIKeyClient{sess, powerinstanceid}
}

/*
This was a change requested by the IBM cloud Team to move the powerinstanceid out from the provider and pass it in the call
The Power-IAAS API requires the crn to be passed in the header.
*/

// Get Key...
func (f *IBMPIKeyClient) Get(id, powerinstanceid string) (*models.SSHKey, error) {

	var tenantid = f.session.UserAccount
	params := p_cloud_tenants_ssh_keys.NewPcloudTenantsSshkeysGetParams().WithTenantID(tenantid).WithSshkeyName(id)
	resp, err := f.session.Power.PCloudTenantsSSHKeys.PcloudTenantsSshkeysGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Get PI Key %s :%s", id, err)
	}
	return resp.Payload, nil
}

// Create PI Key ...
func (f *IBMPIKeyClient) Create(name string, sshkey, powerinstanceid string) (*models.SSHKey, *models.SSHKey, error) {
	var body = models.SSHKey{
		Name:   &name,
		SSHKey: &sshkey,
	}
	params := p_cloud_tenants_ssh_keys.NewPcloudTenantsSshkeysPostParamsWithTimeout(f.session.Timeout).WithTenantID(f.session.UserAccount).WithBody(&body)
	_, postok, err := f.session.Power.PCloudTenantsSSHKeys.PcloudTenantsSshkeysPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || postok == nil {
		return nil, nil, fmt.Errorf("Failed to Create PI Key %s :%s", name, err)
	}
	return nil, postok.Payload, nil

}

// Delete ...
func (f *IBMPIKeyClient) Delete(id string, powerinstanceid string) error {
	var tenantid = f.session.UserAccount
	params := p_cloud_tenants_ssh_keys.NewPcloudTenantsSshkeysDeleteParamsWithTimeout(f.session.Timeout).WithTenantID(tenantid).WithSshkeyName(id)
	_, err := f.session.Power.PCloudTenantsSSHKeys.PcloudTenantsSshkeysDelete(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return fmt.Errorf("Failed to Delete PI Key %s :%s", id, err)
	}
	return nil
}
