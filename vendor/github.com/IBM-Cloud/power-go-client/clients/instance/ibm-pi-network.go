package instance

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_networks"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

// IBMPINetworkClient ...
type IBMPINetworkClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

// NewIBMPINetworkClient ...
func NewIBMPINetworkClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPINetworkClient {
	return &IBMPINetworkClient{
		session:         sess,
		powerinstanceid: powerinstanceid,
	}
}

// Get ...
func (f *IBMPINetworkClient) Get(id, powerinstanceid string, timeout time.Duration) (*models.Network, error) {
	params := p_cloud_networks.NewPcloudNetworksGetParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid).WithNetworkID(id)
	resp, err := f.session.Power.PCloudNetworks.PcloudNetworksGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Get PI Network %s :%s", id, err)
	}
	return resp.Payload, nil
}

// Create ...
func (f *IBMPINetworkClient) Create(name string, networktype string, cidr string, dnsservers []string, gateway string, startip string, endip string, powerinstanceid string, timeout time.Duration) (*models.Network, *models.Network, error) {

	var body = models.NetworkCreate{
		Type: &networktype,
		Name: name,
	}
	if networktype == "vlan" {
		var ipbody = []*models.IPAddressRange{
			{EndingIPAddress: &endip, StartingIPAddress: &startip}}
		if ipbody != nil {
			body.IPAddressRanges = ipbody
		}
		if &gateway != nil {
			body.Gateway = gateway
		}
		if &cidr != nil {
			body.Cidr = cidr
		}
	}
	if dnsservers != nil {
		body.DNSServers = dnsservers
	}
	params := p_cloud_networks.NewPcloudNetworksPostParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid).WithBody(&body)
	_, resp, err := f.session.Power.PCloudNetworks.PcloudNetworksPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp == nil || resp.Payload == nil {
		return nil, nil, fmt.Errorf("Failed to Create PI Network %s :%s", name, err)
	}

	return resp.Payload, nil, nil
}

// GetPublic ...
func (f *IBMPINetworkClient) GetPublic(powerinstanceid string, timeout time.Duration) (*models.Networks, error) {

	filterQuery := "type=\"pub-vlan\""
	params := p_cloud_networks.NewPcloudNetworksGetallParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid).WithFilter(&filterQuery)

	resp, err := f.session.Power.PCloudNetworks.PcloudNetworksGetall(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Get all PI Networks in a power instance %s :%s", powerinstanceid, err)
	}
	return resp.Payload, nil
}

// Delete ...
func (f *IBMPINetworkClient) Delete(id string, powerinstanceid string, timeout time.Duration) error {
	params := p_cloud_networks.NewPcloudNetworksDeleteParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid).WithNetworkID(id)
	_, err := f.session.Power.PCloudNetworks.PcloudNetworksDelete(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil {
		return fmt.Errorf("Failed to Delete PI Network %s :%s", id, err)
	}
	return nil
}

// New Function for Ports

//GetAllPort ...
func (f *IBMPINetworkClient) GetAllPort(id string, powerinstanceid string, timeout time.Duration) (*models.NetworkPorts, error) {

	params := p_cloud_networks.NewPcloudNetworksPortsGetallParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid).WithNetworkID(id)
	resp, err := f.session.Power.PCloudNetworks.PcloudNetworksPortsGetall(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Get all PI Network Ports %s :%s", id, err)
	}
	return resp.Payload, nil

}

// GetPort ...
func (f *IBMPINetworkClient) GetPort(id string, powerinstanceid string, networkPortID string, timeout time.Duration) (*models.NetworkPort, error) {
	params := p_cloud_networks.NewPcloudNetworksPortsGetParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid).WithNetworkID(id).WithPortID(networkPortID)
	resp, err := f.session.Power.PCloudNetworks.PcloudNetworksPortsGet(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to Get PI Network Ports %s :%s", networkPortID, err)
	}
	return resp.Payload, nil

}

//CreatePort ...
func (f *IBMPINetworkClient) CreatePort(id string, powerinstanceid string, networkportdef *p_cloud_networks.PcloudNetworksPortsPostParams, timeout time.Duration) (*models.NetworkPort, error) {
	params := p_cloud_networks.NewPcloudNetworksPortsPostParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid).WithNetworkID(id).WithBody(networkportdef.Body)
	resp, err := f.session.Power.PCloudNetworks.PcloudNetworksPortsPost(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to create the network port for network %s cloudinstance id [%s]", id, powerinstanceid)
	}
	return resp.Payload, nil
}

// DeletePort ...
func (f *IBMPINetworkClient) DeletePort(networkid string, powerinstanceid string, portid string, timeout time.Duration) (*models.Object, error) {
	params := p_cloud_networks.NewPcloudNetworksPortsDeleteParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid).WithNetworkID(networkid).WithPortID(portid)
	resp, err := f.session.Power.PCloudNetworks.PcloudNetworksPortsDelete(params, ibmpisession.NewAuth(f.session, powerinstanceid))

	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to create the network port %s for network %s cloudinstance id [%s]", portid, networkid, powerinstanceid)
	}
	return &resp.Payload, nil
}

//AttachPort to the PVM Instance
func (f *IBMPINetworkClient) AttachPort(powerinstanceid, networkID, portID, description, pvminstanceid string, timeout time.Duration) (*models.NetworkPort, error) {

	var body = models.NetworkPortUpdate{}
	if &description != nil {
		body.Description = &description
	}
	if &pvminstanceid != nil {
		body.PvmInstanceID = &pvminstanceid
	}

	params := p_cloud_networks.NewPcloudNetworksPortsPutParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid).WithNetworkID(networkID).WithPortID(portID).WithBody(&body)
	resp, err := f.session.Power.PCloudNetworks.PcloudNetworksPortsPut(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to attach the port [%s] to network %s the pvminstance [%s]", portID, networkID, pvminstanceid)
	}
	return resp.Payload, nil
}

// DetachPort from the PVM Instance
func (f *IBMPINetworkClient) DetachPort(powerinstanceid, networkID, portID string, timeout time.Duration) (*models.NetworkPort, error) {
	emptyPVM := ""
	body := &models.NetworkPortUpdate{
		PvmInstanceID: &emptyPVM,
	}
	params := p_cloud_networks.NewPcloudNetworksPortsPutParamsWithTimeout(timeout).WithCloudInstanceID(powerinstanceid).WithNetworkID(networkID).WithPortID(portID).WithBody(body)
	resp, err := f.session.Power.PCloudNetworks.PcloudNetworksPortsPut(params, ibmpisession.NewAuth(f.session, powerinstanceid))
	if err != nil || resp == nil || resp.Payload == nil {
		return nil, fmt.Errorf("Failed to detach the port [%s] to network %s ", portID, networkID)
	}

	return resp.Payload, nil
}
