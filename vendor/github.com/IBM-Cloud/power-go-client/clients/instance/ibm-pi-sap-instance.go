package instance

import (
	"fmt"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_s_a_p"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPISAPInstanceClient ...
type IBMPISAPInstanceClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

// NewIBMPISAPInstanceClient ...
func NewIBMPISAPInstanceClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPISAPInstanceClient {
	return &IBMPISAPInstanceClient{
		session:         sess,
		powerinstanceid: powerinstanceid,
	}
}

//Create SAP System
func (f *IBMPISAPInstanceClient) Create(sapdef *p_cloud_s_a_p.PcloudSapPostParams, id, powerinstanceid string) (*models.PVMInstanceList, error) {

	params := p_cloud_s_a_p.NewPcloudSapPostParamsWithTimeout(f.session.Timeout).WithCloudInstanceID(powerinstanceid).WithBody(sapdef.Body)
	sapok, sapcreated, sapaccepted, err := f.session.Power.PCloudSAP.PcloudSapPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return nil, fmt.Errorf("Failed to Create Sap Instance %s", err)
	}

	if sapok != nil && len(sapok.Payload) > 0 {
		return &sapok.Payload, nil
	}
	if sapcreated != nil && len(sapcreated.Payload) > 0 {
		return &sapcreated.Payload, nil
	}
	if sapaccepted != nil && len(sapaccepted.Payload) > 0 {
		return &sapaccepted.Payload, nil
	}

	//return &postok.Payload, nil
	return nil, fmt.Errorf("No response Returned ")
}
