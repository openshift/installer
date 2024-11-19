// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package transitgateway

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

const (
	isTransitGatewayConnectionTunnelDeleting  = "deleting"
	isTransitGatewayConnectionTunnelDetaching = "detaching"
	isTransitGatewayConnectionTunnelDeleted   = "detached"
	isTransitGatewayConnectionTunnelPending   = "pending"
	isTransitGatewayConnectionTunnelAttached  = "attached"
)

func ResourceIBMTransitGatewayConnectionRgreTunnel() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMTransitGatewayConnectionRgreTunnelCreate,
		Read:     resourceIBMTransitGatewayConnectionRgreTunnelRead,
		Delete:   resourceIBMTransitGatewayConnectionRgreTunnelDelete,
		Exists:   resourceIBMTransitGatewayConnectionRgreTunnelExists,
		Update:   resourceIBMTransitGatewayConnectionRgreTunnelUpdate,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			tgGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Transit Gateway identifier",
			},
			tgConnectionId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Transit Gateway Connection identifier",
			},
			tgconTunnelName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validate.InvokeValidator("ibm_tg_connection_rgre_tunnel", tgconTunnelName),
				Description:  "The user-defined name for this tunnel connection.",
			},
			tgLocalGatewayIp: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The local gateway IP address.",
			},
			tgLocalTunnelIp: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The local tunnel IP address.",
			},
			tgRemoteGatewayIp: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The remote gateway IP address.",
			},
			tgRemoteTunnelIp: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The remote tunnel IP address.",
			},
			tgZone: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Location of GRE tunnel.",
			},
			tgCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this tunnel was created",
			},
			tgUpdatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this tunnel was last updated",
			},
			tgGreTunnelStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "What is the current configuration state of this tunnel. Possible values: [attached,failed,pending,deleting,detaching,detached]",
			},
			tgGreTunnelId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Transit Gateway tunnel identifier",
			},
			tgMtu: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Only visible for cross account connections, this field represents the status of the request to connect the given network between accounts.Possible values: [pending,approved,rejected,expired,detached]",
			},
			tgLocalBgpAsn: {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The local network BGP ASN.",
			},
			tgRemoteBgpAsn: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The remote network BGP ASN.",
			},
			tgBaseNetworkType: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The type of the base network for the RGRE. It should be i.e classic or VPC",
			},
			tgNetworkAccountID: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the account which owns the network that is being connected. Generally only used if the network is in a different account than the gateway. This field is required for type 'unbound_gre_tunnel' when the associated_network_type is 'classic' and the GRE tunnel is in a different account than the gateway.",
			},
			tgNetworkId: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the network being connected via this connection. This field is required for some types, such as 'vpc' or 'directlink' or 'power_virtual_server'. The value of this is the CRN of the VPC or direct link or power_virtual_server gateway to be connected. This field is required to be unspecified for network type 'classic', 'gre_tunnel', and 'unbound_gre_tunnel'.",
			},
		},
	}
}

func ResourceIBMTransitGatewayConnectionRgreTunnelValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 tgName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-zA-Z]|[a-zA-Z][-_a-zA-Z0-9]*[a-zA-Z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	ibmTransitGatewayConnectionTunnelResourceValidator := validate.ResourceValidator{ResourceName: "ibm_tg_connection_rgre_tunnel", Schema: validateSchema}

	return &ibmTransitGatewayConnectionTunnelResourceValidator
}
func resourceIBMTransitGatewayConnectionRgreTunnelCreate(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	createTransitGatewayConnectionRgreTunnelOptions := &transitgatewayapisv1.CreateTransitGatewayGreTunnelOptions{}

	gatewayId := d.Get(tgGatewayId).(string)
	createTransitGatewayConnectionRgreTunnelOptions.SetTransitGatewayID(gatewayId)

	connectionId := d.Get(tgConnectionId).(string)
	createTransitGatewayConnectionRgreTunnelOptions.SetID(connectionId)

	if _, ok := d.GetOk(tgconTunnelName); ok {
		name := d.Get(tgconTunnelName).(string)
		createTransitGatewayConnectionRgreTunnelOptions.SetName(name)
	}

	if _, ok := d.GetOk(tgLocalGatewayIp); ok {
		localGatewayIp := d.Get(tgLocalGatewayIp).(string)
		createTransitGatewayConnectionRgreTunnelOptions.SetLocalGatewayIp(localGatewayIp)
	}

	if _, ok := d.GetOk(tgLocalTunnelIp); ok {
		localTunnelIp := d.Get(tgLocalTunnelIp).(string)
		createTransitGatewayConnectionRgreTunnelOptions.SetLocalTunnelIp(localTunnelIp)
	}
	if _, ok := d.GetOk(tgRemoteBgpAsn); ok {
		remoteBgpAsn := int64(d.Get(tgRemoteBgpAsn).(int))
		createTransitGatewayConnectionRgreTunnelOptions.SetRemoteBgpAsn(remoteBgpAsn)
	}
	if _, ok := d.GetOk(tgRemoteGatewayIp); ok {
		remoteGatewayIp := d.Get(tgRemoteGatewayIp).(string)
		createTransitGatewayConnectionRgreTunnelOptions.SetRemoteGatewayIp(remoteGatewayIp)
	}
	if _, ok := d.GetOk(tgRemoteTunnelIp); ok {
		remoteTunnelIp := d.Get(tgRemoteTunnelIp).(string)
		createTransitGatewayConnectionRgreTunnelOptions.SetRemoteTunnelIp(remoteTunnelIp)
	}
	if _, ok := d.GetOk(tgZone); ok {
		zoneIdentity := &transitgatewayapisv1.ZoneIdentity{}
		zoneName := d.Get(tgZone).(string)
		zoneIdentity.Name = &zoneName
		createTransitGatewayConnectionRgreTunnelOptions.SetZone(zoneIdentity)
	}

	rGRETunnel, response, err := client.CreateTransitGatewayGreTunnel(createTransitGatewayConnectionRgreTunnelOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Create Transit Gateway connection  rGRE tunnel err %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", gatewayId, connectionId, *rGRETunnel.ID))
	d.Set(tgGreTunnelId, *rGRETunnel.ID)

	return resourceIBMTransitGatewayConnectionRgreTunnelRead(d, meta)
}
func isWaitForTransitGatewayConnectionRgreTunnelAvailable(client *transitgatewayapisv1.TransitGatewayApisV1, id string, timeout time.Duration) (interface{}, error) {

	log.Printf("Waiting for transit gateway connection tunnel (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isTransitGatewayConnectionPending},
		Target:     []string{isTransitGatewayConnectionTunnelAttached, ""},
		Refresh:    isTransitGatewayConnectionRgreTunnelRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
func isTransitGatewayConnectionRgreTunnelRefreshFunc(client *transitgatewayapisv1.TransitGatewayApisV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		parts, err := flex.IdParts(id)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Transit Gateway connection: %s", err)
		}

		gatewayId := parts[0]
		connectionID := parts[1]
		rGRETunnelID := parts[2]

		getTransitGatewayConnectionrGRETunnelOptions := &transitgatewayapisv1.GetTransitGatewayConnectionTunnelsOptions{}
		getTransitGatewayConnectionrGRETunnelOptions.SetTransitGatewayID(gatewayId)
		getTransitGatewayConnectionrGRETunnelOptions.SetID(connectionID)
		getTransitGatewayConnectionrGRETunnelOptions.SetGreTunnelID(rGRETunnelID)
		trGREunnel, response, err := client.GetTransitGatewayConnectionTunnels(getTransitGatewayConnectionrGRETunnelOptions)

		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Transit Gateway Connection Tunnel (%s): %s\n%s", rGRETunnelID, err, response)
		}
		if *trGREunnel.Status == "attached" || *trGREunnel.Status == "failed" {
			return trGREunnel, isTransitGatewayConnectionTunnelAttached, nil
		}
		return trGREunnel, isTransitGatewayConnectionTunnelPending, nil
	}
}
func resourceIBMTransitGatewayConnectionRgreTunnelRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	connectionID := parts[1]
	rGRETunnelID := parts[2]

	getTransitGatewayConnectionOptions := &transitgatewayapisv1.GetTransitGatewayConnectionOptions{}
	getTransitGatewayConnectionOptions.SetTransitGatewayID(gatewayId)
	getTransitGatewayConnectionOptions.SetID(connectionID)
	tgConnection, response, err := client.GetTransitGatewayConnection(getTransitGatewayConnectionOptions)

	// if it's  not across account wait for it became available.
	if tgConnection.RequestStatus == nil {
		_, err = isWaitForTransitGatewayConnectionRgreTunnelAvailable(client, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}

	getTransitGatewayConnectionrGRETunnelOptions := &transitgatewayapisv1.GetTransitGatewayConnectionTunnelsOptions{}
	getTransitGatewayConnectionrGRETunnelOptions.SetTransitGatewayID(gatewayId)
	getTransitGatewayConnectionrGRETunnelOptions.SetID(connectionID)
	getTransitGatewayConnectionrGRETunnelOptions.SetGreTunnelID(rGRETunnelID)
	instance, response, err := client.GetTransitGatewayConnectionTunnels(getTransitGatewayConnectionrGRETunnelOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Transit Gateway Connection  Redundant GRE Tunnel (%s): %s\n%s", rGRETunnelID, err, response)
	}

	if instance.Name != nil {
		d.Set(tgconTunnelName, *instance.Name)
	}
	if instance.BaseNetworkType != nil {
		d.Set(tgBaseNetworkType, *instance.BaseNetworkType)
	}
	if instance.UpdatedAt != nil {
		d.Set(tgUpdatedAt, instance.UpdatedAt.String())
	}
	if instance.CreatedAt != nil {
		d.Set(tgCreatedAt, instance.CreatedAt.String())
	}
	if instance.NetworkID != nil {
		d.Set(tgNetworkId, *instance.NetworkID)
	}
	if instance.Status != nil {
		d.Set(tgGreTunnelStatus, *instance.Status)
	}
	if instance.NetworkAccountID != nil {
		d.Set(tgNetworkAccountID, *instance.NetworkAccountID)
	}
	if instance.Mtu != nil {
		d.Set(tgMtu, *instance.Mtu)
	}
	if instance.RemoteBgpAsn != nil {
		d.Set(tgRemoteBgpAsn, *instance.RemoteBgpAsn)
	}
	if instance.LocalBgpAsn != nil {
		d.Set(tgLocalBgpAsn, *instance.LocalBgpAsn)
	}

	d.Set(tgConnectionId, connectionID)
	d.Set(tgGatewayId, gatewayId)
	d.Set(tgGreTunnelId, *instance.ID)

	return nil
}

func resourceIBMTransitGatewayConnectionRgreTunnelUpdate(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	connectionID := parts[1]
	rGRETunnelID := parts[2]

	updateTransitGatewayConnectionOptions := &transitgatewayapisv1.UpdateTransitGatewayConnectionTunnelsOptions{}
	updateTransitGatewayConnectionOptions.SetTransitGatewayID(gatewayId)
	updateTransitGatewayConnectionOptions.SetID(connectionID)
	updateTransitGatewayConnectionOptions.SetGreTunnelID(rGRETunnelID)

	if d.HasChange(tgconTunnelName) {
		if d.Get(tgconTunnelName) != nil {
			name := d.Get(tgconTunnelName).(string)
			updateTransitGatewayConnectionOptions.Name = &name
		}
	}

	_, response, err := client.UpdateTransitGatewayConnectionTunnels(updateTransitGatewayConnectionOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error in Update Transit Gateway Connection Tunnel: %s\n%s", err, response)
	}

	return resourceIBMTransitGatewayConnectionRgreTunnelRead(d, meta)
}

func resourceIBMTransitGatewayConnectionRgreTunnelDelete(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	connectionID := parts[1]
	rGRETunnelID := parts[2]

	deleteTransitGatewayConnectionRgreTunnelOptions := &transitgatewayapisv1.DeleteTransitGatewayConnectionTunnelsOptions{}

	deleteTransitGatewayConnectionRgreTunnelOptions.SetTransitGatewayID(gatewayId)
	deleteTransitGatewayConnectionRgreTunnelOptions.SetID(connectionID)
	deleteTransitGatewayConnectionRgreTunnelOptions.SetGreTunnelID(rGRETunnelID)

	response, err := client.DeleteTransitGatewayConnectionTunnels(deleteTransitGatewayConnectionRgreTunnelOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error deleting Transit Gateway Connection Tunnel (%s): %s\n%s", rGRETunnelID, err, response)
	}
	_, err = isWaitForTransitGatewayConnectionRgreTunnelDeleted(client, d.Id(), d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForTransitGatewayConnectionRgreTunnelDeleted(client *transitgatewayapisv1.TransitGatewayApisV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for transit gateway Connection Tunnel (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isTransitGatewayConnectionTunnelDeleting, isTransitGatewayConnectionDetaching},
		Target:     []string{"", isTransitGatewayConnectionTunnelDeleted},
		Refresh:    isTransitGatewayConnectionRgreTunnelDeleteRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isTransitGatewayConnectionRgreTunnelDeleteRefreshFunc(client *transitgatewayapisv1.TransitGatewayApisV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] tg gateway connection Tunnel delete function here")
		parts, err := flex.IdParts(id)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Transit Gateway connection Tunnel: %s", err)

		}

		gatewayId := parts[0]
		connectionID := parts[1]
		rGRETunnelID := parts[2]

		getTransitGatewayConnectionrGRETunnelOptions := &transitgatewayapisv1.GetTransitGatewayConnectionTunnelsOptions{}
		getTransitGatewayConnectionrGRETunnelOptions.SetTransitGatewayID(gatewayId)
		getTransitGatewayConnectionrGRETunnelOptions.SetID(connectionID)
		getTransitGatewayConnectionrGRETunnelOptions.SetGreTunnelID(rGRETunnelID)
		trGREunnel, response, err := client.GetTransitGatewayConnectionTunnels(getTransitGatewayConnectionrGRETunnelOptions)

		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return trGREunnel, isTransitGatewayConnectionTunnelDeleted, nil
			}
			return nil, "", fmt.Errorf("[ERROR] Error Getting Transit Gateway Connection Tunnel (%s): %s\n%s", rGRETunnelID, err, response)
		}
		return trGREunnel, isTransitGatewayConnectionTunnelDeleting, err
	}
}
func resourceIBMTransitGatewayConnectionRgreTunnelExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client, err := transitgatewayClient(meta)
	if err != nil {
		return false, err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) < 3 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of gatewayID/ConnectionID/TunnelID", d.Id())
	}

	gatewayId := parts[0]
	connectionID := parts[1]
	rGRETunnelID := parts[2]

	getTransitGatewayConnectionrGRETunnelOptions := &transitgatewayapisv1.GetTransitGatewayConnectionTunnelsOptions{}
	getTransitGatewayConnectionrGRETunnelOptions.SetTransitGatewayID(gatewayId)
	getTransitGatewayConnectionrGRETunnelOptions.SetID(connectionID)
	getTransitGatewayConnectionrGRETunnelOptions.SetGreTunnelID(rGRETunnelID)
	_, response, err := client.GetTransitGatewayConnectionTunnels(getTransitGatewayConnectionrGRETunnelOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error Getting Transit Gateway Connection Tunnel: %s\n%s", err, response)
	}
	return true, nil
}
