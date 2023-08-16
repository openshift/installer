// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
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
	tgGatewayConnections                = "gateway_connections"
	tgNetworkId                         = "network_id"
	tgNetworkType                       = "network_type"
	tgNetworkAccountID                  = "network_account_id"
	tgConectionCreatedAt                = "created_at"
	tgConnectionStatus                  = "status"
	tgGatewayId                         = "gateway"
	isTransitGatewayConnectionDeleting  = "deleting"
	isTransitGatewayConnectionDetaching = "detaching"
	isTransitGatewayConnectionDeleted   = "detached"
	isTransitGatewayConnectionPending   = "pending"
	isTransitGatewayConnectionAttached  = "attached"
	tgRequestStatus                     = "request_status"
	tgConnectionId                      = "connection_id"
	tgBaseConnectionId                  = "base_connection_id"
	tgLocalBgpAsn                       = "local_bgp_asn"
	tgLocalGatewayIp                    = "local_gateway_ip"
	tgLocalTunnelIp                     = "local_tunnel_ip"
	tgRemoteBgpAsn                      = "remote_bgp_asn"
	tgRemoteGatewayIp                   = "remote_gateway_ip"
	tgRemoteTunnelIp                    = "remote_tunnel_ip"
	tgZone                              = "zone"
	tgMtu                               = "mtu"
)

func ResourceIBMTransitGatewayConnection() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMTransitGatewayConnectionCreate,
		Read:     resourceIBMTransitGatewayConnectionRead,
		Delete:   resourceIBMTransitGatewayConnectionDelete,
		Exists:   resourceIBMTransitGatewayConnectionExists,
		Update:   resourceIBMTransitGatewayConnectionUpdate,
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
				Computed:    true,
				Description: "The Transit Gateway Connection identifier",
			},
			tgNetworkType: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tg_connection", tgNetworkType),
				Description:  "Defines what type of network is connected via this connection. Allowable values (classic,directlink,vpc,gre_tunnel,unbound_gre_tunnel,power_virtual_server)",
			},
			tgName: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validate.InvokeValidator("ibm_tg_connection", tgName),
				Description:  "The user-defined name for this transit gateway. If unspecified, the name will be the network name (the name of the VPC in the case of network type 'vpc', and the word Classic, in the case of network type 'classic').",
			},
			tgNetworkId: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the network being connected via this connection. This field is required for some types, such as 'vpc' or 'directlink' or 'power_virtual_server'. The value of this is the CRN of the VPC or direct link or power_virtual_server gateway to be connected. This field is required to be unspecified for network type 'classic', 'gre_tunnel', and 'unbound_gre_tunnel'.",
			},
			tgNetworkAccountID: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the account which owns the network that is being connected. Generally only used if the network is in a different account than the gateway. This field is required for type 'unbound_gre_tunnel' when the associated_network_type is 'classic' and the GRE tunnel is in a different account than the gateway.",
			},
			tgBaseConnectionId: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of a network_type 'classic' connection a tunnel is configured over. This field only applies to network type 'gre_tunnel' connections.",
			},
			tgBaseNetworkType: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The type of network the unbound gre tunnel is targeting. This field is required for network type 'unbound_gre_tunnel'.",
			},
			tgLocalGatewayIp: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The local gateway IP address. This field only applies to network type 'gre_tunnel' and 'unbound_gre_tunnel' connections.",
			},
			tgLocalTunnelIp: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The local tunnel IP address. This field only applies to network type 'gre_tunnel' and 'unbound_gre_tunnel' connections.",
			},
			tgRemoteBgpAsn: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The remote network BGP ASN. This field only applies to network type 'gre_tunnel' and 'unbound_gre_tunnel' connections.",
			},
			tgRemoteGatewayIp: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The remote gateway IP address. This field only applies to network type 'gre_tunnel' and 'unbound_gre_tunnel' connections.",
			},
			tgRemoteTunnelIp: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The remote tunnel IP address. This field only applies to network type 'gre_tunnel' and 'unbound_gre_tunnel' connections.",
			},
			tgZone: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Location of GRE tunnel. This field only applies to network type 'gre_tunnel' and 'unbound_gre_tunnel' connections.",
			},
			tgCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this connection was created",
			},
			tgUpdatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this connection was last updated",
			},
			tgStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "What is the current configuration state of this connection. Possible values: [attached,failed,pending,deleting,detaching,detached]",
			},
			tgRequestStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Only visible for cross account connections, this field represents the status of the request to connect the given network between accounts.Possible values: [pending,approved,rejected,expired,detached]",
			},
			flex.RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the transit gateway",
			},
		},
	}
}
func ResourceIBMTransitGatewayConnectionValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	networkType := "classic, directlink, vpc, gre_tunnel, unbound_gre_tunnel, power_virtual_server"
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 tgNetworkType,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              networkType})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 tgName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-zA-Z]|[a-zA-Z][-_a-zA-Z0-9]*[a-zA-Z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	ibmTransitGatewayConnectionResourceValidator := validate.ResourceValidator{ResourceName: "ibm_tg_connection", Schema: validateSchema}

	return &ibmTransitGatewayConnectionResourceValidator
}
func resourceIBMTransitGatewayConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	createTransitGatewayConnectionOptions := &transitgatewayapisv1.CreateTransitGatewayConnectionOptions{}

	gatewayId := d.Get(tgGatewayId).(string)
	createTransitGatewayConnectionOptions.SetTransitGatewayID(gatewayId)

	if _, ok := d.GetOk(tgName); ok {
		name := d.Get(tgName).(string)
		createTransitGatewayConnectionOptions.SetName(name)
	}

	networkType := d.Get(tgNetworkType).(string)
	createTransitGatewayConnectionOptions.SetNetworkType(networkType)
	if _, ok := d.GetOk(tgNetworkId); ok {
		networkID := d.Get(tgNetworkId).(string)
		createTransitGatewayConnectionOptions.SetNetworkID(networkID)
	}
	if _, ok := d.GetOk(tgNetworkAccountID); ok {
		networkAccId := d.Get(tgNetworkAccountID).(string)
		createTransitGatewayConnectionOptions.SetNetworkAccountID(networkAccId)
	}
	if _, ok := d.GetOk(tgBaseConnectionId); ok {
		baseConnectionId := d.Get(tgBaseConnectionId).(string)
		createTransitGatewayConnectionOptions.SetBaseConnectionID(baseConnectionId)
	}
	if _, ok := d.GetOk(tgBaseNetworkType); ok {
		baseNetworkType := d.Get(tgBaseNetworkType).(string)
		createTransitGatewayConnectionOptions.SetBaseNetworkType(baseNetworkType)
	}
	if _, ok := d.GetOk(tgLocalGatewayIp); ok {
		localGatewayIp := d.Get(tgLocalGatewayIp).(string)
		createTransitGatewayConnectionOptions.SetLocalGatewayIp(localGatewayIp)
	}
	if _, ok := d.GetOk(tgLocalTunnelIp); ok {
		localTunnelIp := d.Get(tgLocalTunnelIp).(string)
		createTransitGatewayConnectionOptions.SetLocalTunnelIp(localTunnelIp)
	}
	if _, ok := d.GetOk(tgRemoteBgpAsn); ok {
		remoteBgpAsn := int64(d.Get(tgRemoteBgpAsn).(int))
		createTransitGatewayConnectionOptions.SetRemoteBgpAsn(remoteBgpAsn)
	}
	if _, ok := d.GetOk(tgRemoteGatewayIp); ok {
		remoteGatewayIp := d.Get(tgRemoteGatewayIp).(string)
		createTransitGatewayConnectionOptions.SetRemoteGatewayIp(remoteGatewayIp)
	}
	if _, ok := d.GetOk(tgRemoteTunnelIp); ok {
		remoteTunnelIp := d.Get(tgRemoteTunnelIp).(string)
		createTransitGatewayConnectionOptions.SetRemoteTunnelIp(remoteTunnelIp)
	}
	if _, ok := d.GetOk(tgZone); ok {
		zoneIdentity := &transitgatewayapisv1.ZoneIdentity{}
		zoneName := d.Get(tgZone).(string)
		zoneIdentity.Name = &zoneName
		createTransitGatewayConnectionOptions.SetZone(zoneIdentity)
	}

	tgConnections, response, err := client.CreateTransitGatewayConnection(createTransitGatewayConnectionOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Create Transit Gateway connection err %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s", gatewayId, *tgConnections.ID))
	d.Set(tgConnectionId, *tgConnections.ID)

	if tgConnections.NetworkAccountID != nil {
		d.Set(tgNetworkAccountID, *tgConnections.NetworkAccountID)
		return resourceIBMTransitGatewayConnectionRead(d, meta)
	}
	_, err = isWaitForTransitGatewayConnectionAvailable(client, d.Id(), d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return err
	}
	return resourceIBMTransitGatewayConnectionRead(d, meta)
}
func isWaitForTransitGatewayConnectionAvailable(client *transitgatewayapisv1.TransitGatewayApisV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for transit gateway connection (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isTransitGatewayConnectionPending},
		Target:     []string{isTransitGatewayConnectionAttached, ""},
		Refresh:    isTransitGatewayConnectionRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
func isTransitGatewayConnectionRefreshFunc(client *transitgatewayapisv1.TransitGatewayApisV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		parts, err := flex.IdParts(id)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Transit Gateway connection: %s", err)
			//	return err
		}

		gatewayId := parts[0]
		ID := parts[1]
		getTransitGatewayConnectionOptions := &transitgatewayapisv1.GetTransitGatewayConnectionOptions{}
		getTransitGatewayConnectionOptions.SetTransitGatewayID(gatewayId)
		getTransitGatewayConnectionOptions.SetID(ID)
		tgConnection, response, err := client.GetTransitGatewayConnection(getTransitGatewayConnectionOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Transit Gateway Connection (%s): %s\n%s", ID, err, response)
		}
		if *tgConnection.Status == "attached" || *tgConnection.Status == "failed" {
			return tgConnection, isTransitGatewayConnectionAttached, nil
		}

		return tgConnection, isTransitGatewayConnectionPending, nil
	}
}
func resourceIBMTransitGatewayConnectionRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	ID := parts[1]

	getTransitGatewayConnectionOptions := &transitgatewayapisv1.GetTransitGatewayConnectionOptions{}
	getTransitGatewayConnectionOptions.SetTransitGatewayID(gatewayId)
	getTransitGatewayConnectionOptions.SetID(ID)
	instance, response, err := client.GetTransitGatewayConnection(getTransitGatewayConnectionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Transit Gateway Connection (%s): %s\n%s", ID, err, response)
	}

	if instance.Name != nil {
		d.Set(tgName, *instance.Name)
	}
	if instance.NetworkType != nil {
		d.Set(tgNetworkType, *instance.NetworkType)
	}
	if instance.UpdatedAt != nil {
		d.Set(tgUpdatedAt, instance.UpdatedAt.String())
	}
	if instance.NetworkID != nil {
		d.Set(tgNetworkId, *instance.NetworkID)
	}
	if instance.CreatedAt != nil {
		d.Set(tgCreatedAt, instance.CreatedAt.String())
	}
	if instance.Status != nil {
		d.Set(tgStatus, *instance.Status)
	}
	if instance.NetworkAccountID != nil {
		d.Set(tgNetworkAccountID, *instance.NetworkAccountID)
	}
	if instance.RequestStatus != nil {
		d.Set(tgRequestStatus, *instance.RequestStatus)
	}
	d.Set(tgConnectionId, *instance.ID)
	d.Set(tgGatewayId, gatewayId)
	getTransitGatewayOptions := &transitgatewayapisv1.GetTransitGatewayOptions{
		ID: &gatewayId,
	}
	tgw, response, err := client.GetTransitGateway(getTransitGatewayOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Getting Transit Gateway : %s\n%s", err, response)
	}
	d.Set(flex.RelatedCRN, *tgw.Crn)

	return nil
}

func resourceIBMTransitGatewayConnectionUpdate(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	ID := parts[1]

	getTransitGatewayConnectionOptions := &transitgatewayapisv1.GetTransitGatewayConnectionOptions{
		ID: &ID,
	}
	getTransitGatewayConnectionOptions.SetTransitGatewayID(gatewayId)

	_, response, err := client.GetTransitGatewayConnection(getTransitGatewayConnectionOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Getting Transit Gateway Connection: %s\n%s", err, response)
	}

	updateTransitGatewayConnectionOptions := &transitgatewayapisv1.UpdateTransitGatewayConnectionOptions{}
	updateTransitGatewayConnectionOptions.ID = &ID
	updateTransitGatewayConnectionOptions.SetTransitGatewayID(gatewayId)
	if d.HasChange(tgName) {
		if d.Get(tgName) != nil {
			name := d.Get(tgName).(string)
			updateTransitGatewayConnectionOptions.Name = &name
		}
	}

	_, response, err = client.UpdateTransitGatewayConnection(updateTransitGatewayConnectionOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error in Update Transit Gateway Connection : %s\n%s", err, response)
	}

	return resourceIBMTransitGatewayConnectionRead(d, meta)
}

func resourceIBMTransitGatewayConnectionDelete(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	ID := parts[1]
	deleteTransitGatewayConnectionOptions := &transitgatewayapisv1.DeleteTransitGatewayConnectionOptions{
		ID: &ID,
	}
	deleteTransitGatewayConnectionOptions.SetTransitGatewayID(gatewayId)
	response, err := client.DeleteTransitGatewayConnection(deleteTransitGatewayConnectionOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error deleting Transit Gateway Connection(%s): %s\n%s", ID, err, response)
	}
	_, err = isWaitForTransitGatewayConnectionDeleted(client, d.Id(), d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func isWaitForTransitGatewayConnectionDeleted(client *transitgatewayapisv1.TransitGatewayApisV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for transit gateway Connection(%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isTransitGatewayConnectionDeleting, isTransitGatewayConnectionDetaching},
		Target:     []string{"", isTransitGatewayConnectionDeleted},
		Refresh:    isTransitGatewayConnectionDeleteRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isTransitGatewayConnectionDeleteRefreshFunc(client *transitgatewayapisv1.TransitGatewayApisV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] tg gateway connection delete function here")
		parts, err := flex.IdParts(id)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Transit Gateway connection: %s", err)

		}

		gatewayId := parts[0]
		ID := parts[1]
		getTransitGatewayConnectionOptions := &transitgatewayapisv1.GetTransitGatewayConnectionOptions{}
		getTransitGatewayConnectionOptions.SetTransitGatewayID(gatewayId)
		getTransitGatewayConnectionOptions.SetID(ID)
		tgConnection, response, err := client.GetTransitGatewayConnection(getTransitGatewayConnectionOptions)

		if err != nil {

			if response != nil && response.StatusCode == 404 {
				return tgConnection, isTransitGatewayConnectionDeleted, nil
			}

			return nil, "", fmt.Errorf("[ERROR] Error Getting Transit Gateway Connection (%s): %s\n%s", ID, err, response)
		}
		return tgConnection, isTransitGatewayConnectionDeleting, err
	}
}
func resourceIBMTransitGatewayConnectionExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client, err := transitgatewayClient(meta)
	if err != nil {
		return false, err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) < 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of gatewayID/ConnectionID", d.Id())
	}
	gatewayId := parts[0]
	ID := parts[1]

	getTransitGatewayConnectionOptions := &transitgatewayapisv1.GetTransitGatewayConnectionOptions{
		ID: &ID,
	}
	getTransitGatewayConnectionOptions.SetTransitGatewayID(gatewayId)
	_, response, err := client.GetTransitGatewayConnection(getTransitGatewayConnectionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error Getting Transit Gateway Connection: %s\n%s", err, response)
	}

	return true, nil
}
