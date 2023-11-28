// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isVPNGatewayConnectionAdminStateup              = "admin_state_up"
	isVPNGatewayConnectionAdminAuthenticationmode   = "authentication_mode"
	isVPNGatewayConnectionName                      = "name"
	isVPNGatewayConnectionVPNGateway                = "vpn_gateway"
	isVPNGatewayConnection                          = "gateway_connection"
	isVPNGatewayConnectionPeerAddress               = "peer_address"
	isVPNGatewayConnectionPreSharedKey              = "preshared_key"
	isVPNGatewayConnectionLocalCIDRS                = "local_cidrs"
	isVPNGatewayConnectionPeerCIDRS                 = "peer_cidrs"
	isVPNGatewayConnectionIKEPolicy                 = "ike_policy"
	isVPNGatewayConnectionIPSECPolicy               = "ipsec_policy"
	isVPNGatewayConnectionDeadPeerDetectionAction   = "action"
	isVPNGatewayConnectionDeadPeerDetectionInterval = "interval"
	isVPNGatewayConnectionDeadPeerDetectionTimeout  = "timeout"
	isVPNGatewayConnectionStatus                    = "status"
	isVPNGatewayConnectionDeleting                  = "deleting"
	isVPNGatewayConnectionDeleted                   = "done"
	isVPNGatewayConnectionProvisioning              = "provisioning"
	isVPNGatewayConnectionProvisioningDone          = "done"
	isVPNGatewayConnectionMode                      = "mode"
	isVPNGatewayConnectionTunnels                   = "tunnels"
	isVPNGatewayConnectionResourcetype              = "resource_type"
	isVPNGatewayConnectionCreatedat                 = "created_at"
	isVPNGatewayConnectionStatusreasons             = "status_reasons"
)

func ResourceIBMISVPNGatewayConnection() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISVPNGatewayConnectionCreate,
		Read:     resourceIBMISVPNGatewayConnectionRead,
		Update:   resourceIBMISVPNGatewayConnectionUpdate,
		Delete:   resourceIBMISVPNGatewayConnectionDelete,
		Exists:   resourceIBMISVPNGatewayConnectionExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isVPNGatewayConnectionName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_vpn_gateway_connection", isVPNGatewayConnectionName),
				Description:  "VPN Gateway connection name",
			},

			isVPNGatewayConnectionVPNGateway: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPN Gateway info",
			},

			isVPNGatewayConnectionPeerAddress: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPN gateway connection peer address",
			},

			isVPNGatewayConnectionPreSharedKey: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vpn gateway",
			},

			isVPNGatewayConnectionAdminStateup: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "VPN gateway connection admin state",
			},

			isVPNGatewayConnectionLocalCIDRS: {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "VPN gateway connection local CIDRs",
			},

			isVPNGatewayConnectionPeerCIDRS: {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "VPN gateway connection peer CIDRs",
			},

			isVPNGatewayConnectionDeadPeerDetectionAction: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "restart",
				ValidateFunc: validate.InvokeValidator("ibm_is_vpn_gateway_connection", isVPNGatewayConnectionDeadPeerDetectionAction),
				Description:  "Action detection for dead peer detection action",
			},
			isVPNGatewayConnectionDeadPeerDetectionInterval: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validate.InvokeValidator("ibm_is_vpn_gateway_connection", isVPNGatewayConnectionDeadPeerDetectionInterval),
				Description:  "Interval for dead peer detection interval",
			},
			isVPNGatewayConnectionDeadPeerDetectionTimeout: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validate.InvokeValidator("ibm_is_vpn_gateway_connection", isVPNGatewayConnectionDeadPeerDetectionTimeout),
				Description:  "Timeout for dead peer detection",
			},

			isVPNGatewayConnectionIPSECPolicy: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP security policy for vpn gateway connection",
			},

			isVPNGatewayConnectionIKEPolicy: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPN gateway connection IKE Policy",
			},

			isVPNGatewayConnection: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this VPN gateway connection",
			},

			isVPNGatewayConnectionStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPN gateway connection status",
			},
			isVPNGatewayConnectionStatusreasons: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current status (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason.",
						},
						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about this status reason.",
						},
					},
				},
			},
			flex.RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the VPN Gateway resource",
			},

			isVPNGatewayConnectionAdminAuthenticationmode: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The authentication mode",
			},

			isVPNGatewayConnectionResourcetype: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type",
			},

			isVPNGatewayConnectionCreatedat: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this VPN gateway connection was created",
			},

			isVPNGatewayConnectionMode: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The mode of the VPN gateway",
			},

			isVPNGatewayConnectionTunnels: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPN tunnel configuration for this VPN gateway connection (in static route mode)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address of the VPN gateway member in which the tunnel resides",
						},

						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the VPN Tunnel",
						},
					},
				},
			},
		},
	}
}

func ResourceIBMISVPNGatewayConnectionValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	action := "restart, clear, hold, none"

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVPNGatewayConnectionName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVPNGatewayConnectionDeadPeerDetectionAction,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              action})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVPNGatewayConnectionDeadPeerDetectionInterval,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "1",
			MaxValue:                   "86399"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isVPNGatewayConnectionDeadPeerDetectionTimeout,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "2",
			MaxValue:                   "86399"})

	ibmISVPNGatewayConnectionResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_vpn_gateway_connection", Schema: validateSchema}
	return &ibmISVPNGatewayConnectionResourceValidator
}

func resourceIBMISVPNGatewayConnectionCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] VPNGatewayConnection create")
	name := d.Get(isVPNGatewayConnectionName).(string)
	gatewayID := d.Get(isVPNGatewayConnectionVPNGateway).(string)
	peerAddress := d.Get(isVPNGatewayConnectionPeerAddress).(string)
	prephasedKey := d.Get(isVPNGatewayConnectionPreSharedKey).(string)

	stateUp := false
	if _, ok := d.GetOk(isVPNGatewayConnectionAdminStateup); ok {
		stateUp = d.Get(isVPNGatewayConnectionAdminStateup).(bool)
	}
	var interval, timeout int64
	if intvl, ok := d.GetOk(isVPNGatewayConnectionDeadPeerDetectionInterval); ok {
		interval = int64(intvl.(int))
	} else {
		interval = 30
	}

	if tout, ok := d.GetOk(isVPNGatewayConnectionDeadPeerDetectionTimeout); ok {
		timeout = int64(tout.(int))
	} else {
		timeout = 120
	}
	var action string
	if act, ok := d.GetOk(isVPNGatewayConnectionDeadPeerDetectionAction); ok {
		action = act.(string)
	} else {
		action = "none"
	}

	err := vpngwconCreate(d, meta, name, gatewayID, peerAddress, prephasedKey, action, interval, timeout, stateUp)
	if err != nil {
		return err
	}
	return resourceIBMISVPNGatewayConnectionRead(d, meta)
}

func vpngwconCreate(d *schema.ResourceData, meta interface{}, name, gatewayID, peerAddress, prephasedKey, action string, interval, timeout int64, stateUp bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	vpnGatewayConnectionPrototypeModel := &vpcv1.VPNGatewayConnectionPrototype{
		PeerAddress:  &peerAddress,
		Psk:          &prephasedKey,
		AdminStateUp: &stateUp,
		DeadPeerDetection: &vpcv1.VPNGatewayConnectionDpdPrototype{
			Action:   &action,
			Interval: &interval,
			Timeout:  &timeout,
		},
		Name: &name,
	}
	options := &vpcv1.CreateVPNGatewayConnectionOptions{
		VPNGatewayID:                  &gatewayID,
		VPNGatewayConnectionPrototype: vpnGatewayConnectionPrototypeModel,
	}

	if _, ok := d.GetOk(isVPNGatewayConnectionLocalCIDRS); ok {
		localCidrs := flex.ExpandStringList((d.Get(isVPNGatewayConnectionLocalCIDRS).(*schema.Set)).List())
		vpnGatewayConnectionPrototypeModel.LocalCIDRs = localCidrs
	}
	if _, ok := d.GetOk(isVPNGatewayConnectionPeerCIDRS); ok {
		peerCidrs := flex.ExpandStringList((d.Get(isVPNGatewayConnectionPeerCIDRS).(*schema.Set)).List())
		vpnGatewayConnectionPrototypeModel.PeerCIDRs = peerCidrs
	}

	var ikePolicyIdentity, ipsecPolicyIdentity string

	if ikePolicy, ok := d.GetOk(isVPNGatewayConnectionIKEPolicy); ok {
		ikePolicyIdentity = ikePolicy.(string)
		vpnGatewayConnectionPrototypeModel.IkePolicy = &vpcv1.VPNGatewayConnectionIkePolicyPrototype{
			ID: &ikePolicyIdentity,
		}
	} else {
		vpnGatewayConnectionPrototypeModel.IkePolicy = nil
	}
	if ipsecPolicy, ok := d.GetOk(isVPNGatewayConnectionIPSECPolicy); ok {
		ipsecPolicyIdentity = ipsecPolicy.(string)
		vpnGatewayConnectionPrototypeModel.IpsecPolicy = &vpcv1.VPNGatewayConnectionIPsecPolicyPrototype{
			ID: &ipsecPolicyIdentity,
		}
	} else {
		vpnGatewayConnectionPrototypeModel.IpsecPolicy = nil
	}

	vpnGatewayConnectionIntf, response, err := sess.CreateVPNGatewayConnection(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] Create VPN Gateway Connection err %s\n%s", err, response)
	}
	vpnGatewayConnection := vpnGatewayConnectionIntf.(*vpcv1.VPNGatewayConnection)
	d.SetId(fmt.Sprintf("%s/%s", gatewayID, *vpnGatewayConnection.ID))
	log.Printf("[INFO] VPNGatewayConnection : %s/%s", gatewayID, *vpnGatewayConnection.ID)
	return nil
}

func resourceIBMISVPNGatewayConnectionRead(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gID := parts[0]
	gConnID := parts[1]

	err = vpngwconGet(d, meta, gID, gConnID)
	if err != nil {
		return err
	}
	return nil
}

func vpngwconGet(d *schema.ResourceData, meta interface{}, gID, gConnID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.GetVPNGatewayConnectionOptions{
		VPNGatewayID: &gID,
		ID:           &gConnID,
	}
	vpnGatewayConnectionIntf, response, err := sess.GetVPNGatewayConnection(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Vpn Gateway Connection (%s): %s\n%s", gConnID, err, response)
	}
	d.Set(isVPNGatewayConnection, gConnID)
	vpnGatewayConnection := vpnGatewayConnectionIntf.(*vpcv1.VPNGatewayConnection)
	d.Set(isVPNGatewayConnectionName, *vpnGatewayConnection.Name)
	d.Set(isVPNGatewayConnectionVPNGateway, gID)
	d.Set(isVPNGatewayConnectionAdminStateup, *vpnGatewayConnection.AdminStateUp)
	d.Set(isVPNGatewayConnectionPeerAddress, *vpnGatewayConnection.PeerAddress)
	d.Set(isVPNGatewayConnectionPreSharedKey, *vpnGatewayConnection.Psk)

	if vpnGatewayConnection.LocalCIDRs != nil {
		d.Set(isVPNGatewayConnectionLocalCIDRS, flex.FlattenStringList(vpnGatewayConnection.LocalCIDRs))
	}
	if vpnGatewayConnection.PeerCIDRs != nil {
		d.Set(isVPNGatewayConnectionPeerCIDRS, flex.FlattenStringList(vpnGatewayConnection.PeerCIDRs))
	}
	if vpnGatewayConnection.IkePolicy != nil {
		d.Set(isVPNGatewayConnectionIKEPolicy, *vpnGatewayConnection.IkePolicy.ID)
	}
	if vpnGatewayConnection.IpsecPolicy != nil {
		d.Set(isVPNGatewayConnectionIPSECPolicy, *vpnGatewayConnection.IpsecPolicy.ID)
	}
	if vpnGatewayConnection.AuthenticationMode != nil {
		d.Set(isVPNGatewayConnectionAdminAuthenticationmode, *vpnGatewayConnection.AuthenticationMode)
	}
	if vpnGatewayConnection.Status != nil {
		d.Set(isVPNGatewayConnectionStatus, *vpnGatewayConnection.Status)
	}
	if err := d.Set(isVPNGatewayConnectionStatusreasons, resourceVPNGatewayConnectionFlattenLifecycleReasons(vpnGatewayConnection.StatusReasons)); err != nil {
		return fmt.Errorf("[ERROR] Error setting status_reasons: %s", err)
	}
	if vpnGatewayConnection.ResourceType != nil {
		d.Set(isVPNGatewayConnectionResourcetype, *vpnGatewayConnection.ResourceType)
	}
	if vpnGatewayConnection.CreatedAt != nil {
		d.Set(isVPNGatewayConnectionCreatedat, vpnGatewayConnection.CreatedAt.String())
	}

	if vpnGatewayConnection.Mode != nil {
		d.Set(isVPNGatewayConnectionMode, *vpnGatewayConnection.Mode)
	}
	vpcTunnelsList := make([]map[string]interface{}, 0)
	if vpnGatewayConnection.Tunnels != nil {
		for _, vpcTunnel := range vpnGatewayConnection.Tunnels {
			currentTunnel := map[string]interface{}{}
			if vpcTunnel.PublicIP != nil {
				publicIP := *vpcTunnel.PublicIP
				currentTunnel["address"] = *publicIP.Address
			}
			if vpcTunnel.Status != nil {
				currentTunnel["status"] = *vpcTunnel.Status
			}
			vpcTunnelsList = append(vpcTunnelsList, currentTunnel)
		}
	}
	d.Set(isVPNGatewayConnectionTunnels, vpcTunnelsList)

	d.Set(isVPNGatewayConnectionDeadPeerDetectionAction, *vpnGatewayConnection.DeadPeerDetection.Action)
	d.Set(isVPNGatewayConnectionDeadPeerDetectionInterval, *vpnGatewayConnection.DeadPeerDetection.Interval)
	d.Set(isVPNGatewayConnectionDeadPeerDetectionTimeout, *vpnGatewayConnection.DeadPeerDetection.Timeout)
	getVPNGatewayOptions := &vpcv1.GetVPNGatewayOptions{
		ID: &gID,
	}
	vpngatewayIntf, response, err := sess.GetVPNGateway(getVPNGatewayOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Getting VPN Gateway : %s\n%s", err, response)
	}
	vpngateway := vpngatewayIntf.(*vpcv1.VPNGateway)
	d.Set(flex.RelatedCRN, *vpngateway.CRN)
	return nil
}

func resourceIBMISVPNGatewayConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	hasChanged := false

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gID := parts[0]
	gConnID := parts[1]
	err = vpngwconUpdate(d, meta, gID, gConnID, hasChanged)
	if err != nil {
		return err
	}
	return resourceIBMISVPNGatewayConnectionRead(d, meta)
}

func vpngwconUpdate(d *schema.ResourceData, meta interface{}, gID, gConnID string, hasChanged bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	updateVpnGatewayConnectionOptions := &vpcv1.UpdateVPNGatewayConnectionOptions{
		VPNGatewayID: &gID,
		ID:           &gConnID,
	}
	vpnGatewayConnectionPatchModel := &vpcv1.VPNGatewayConnectionPatch{}
	if d.HasChange(isVPNGatewayConnectionName) {
		name := d.Get(isVPNGatewayConnectionName).(string)
		vpnGatewayConnectionPatchModel.Name = &name
		hasChanged = true
	}

	if d.HasChange(isVPNGatewayConnectionPeerAddress) {
		peerAddress := d.Get(isVPNGatewayConnectionPeerAddress).(string)
		vpnGatewayConnectionPatchModel.PeerAddress = &peerAddress
		hasChanged = true
	}

	if d.HasChange(isVPNGatewayConnectionPreSharedKey) {
		psk := d.Get(isVPNGatewayConnectionPreSharedKey).(string)
		vpnGatewayConnectionPatchModel.Psk = &psk
		hasChanged = true
	}

	if d.HasChange(isVPNGatewayConnectionDeadPeerDetectionAction) || d.HasChange(isVPNGatewayConnectionDeadPeerDetectionInterval) || d.HasChange(isVPNGatewayConnectionDeadPeerDetectionTimeout) {
		action := d.Get(isVPNGatewayConnectionDeadPeerDetectionAction).(string)
		interval := int64(d.Get(isVPNGatewayConnectionDeadPeerDetectionInterval).(int))
		timeout := int64(d.Get(isVPNGatewayConnectionDeadPeerDetectionTimeout).(int))

		// Construct an instance of the VPNGatewayConnectionDpdPatch model
		vpnGatewayConnectionDpdPatchModel := new(vpcv1.VPNGatewayConnectionDpdPatch)
		vpnGatewayConnectionDpdPatchModel.Action = &action
		vpnGatewayConnectionDpdPatchModel.Interval = &interval
		vpnGatewayConnectionDpdPatchModel.Timeout = &timeout
		vpnGatewayConnectionPatchModel.DeadPeerDetection = vpnGatewayConnectionDpdPatchModel
		hasChanged = true
	}

	if d.HasChange(isVPNGatewayConnectionIKEPolicy) {
		ikePolicyIdentity := d.Get(isVPNGatewayConnectionIKEPolicy).(string)
		if ikePolicyIdentity == "" {
			var nullPatch *vpcv1.VPNGatewayConnectionIkePolicyPatch
			vpnGatewayConnectionPatchModel.IkePolicy = nullPatch
		} else {
			vpnGatewayConnectionPatchModel.IkePolicy = &vpcv1.VPNGatewayConnectionIkePolicyPatch{
				ID: &ikePolicyIdentity,
			}
		}
		hasChanged = true
	} else {
		vpnGatewayConnectionPatchModel.IkePolicy = nil
	}

	if d.HasChange(isVPNGatewayConnectionIPSECPolicy) {
		ipsecPolicyIdentity := d.Get(isVPNGatewayConnectionIPSECPolicy).(string)
		if ipsecPolicyIdentity == "" {
			var nullPatch *vpcv1.VPNGatewayConnectionIPsecPolicyPatch
			vpnGatewayConnectionPatchModel.IpsecPolicy = nullPatch
		} else {
			vpnGatewayConnectionPatchModel.IpsecPolicy = &vpcv1.VPNGatewayConnectionIPsecPolicyPatch{
				ID: &ipsecPolicyIdentity,
			}
		}
		hasChanged = true
	} else {
		vpnGatewayConnectionPatchModel.IpsecPolicy = nil
	}

	if d.HasChange(isVPNGatewayConnectionAdminStateup) {
		adminStateUp := d.Get(isVPNGatewayConnectionAdminStateup).(bool)
		vpnGatewayConnectionPatchModel.AdminStateUp = &adminStateUp
		hasChanged = true
	}

	if hasChanged {
		vpnGatewayConnectionPatch, err := vpnGatewayConnectionPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for VPNGatewayConnectionPatch: %s", err)
		}
		updateVpnGatewayConnectionOptions.VPNGatewayConnectionPatch = vpnGatewayConnectionPatch
		_, response, err := sess.UpdateVPNGatewayConnection(updateVpnGatewayConnectionOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating Vpn Gateway Connection: %s\n%s", err, response)
		}
	}
	return nil
}

func resourceIBMISVPNGatewayConnectionDelete(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gID := parts[0]
	gConnID := parts[1]

	err = vpngwconDelete(d, meta, gID, gConnID)
	if err != nil {
		return err
	}
	return nil
}

func vpngwconDelete(d *schema.ResourceData, meta interface{}, gID, gConnID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getVpnGatewayConnectionOptions := &vpcv1.GetVPNGatewayConnectionOptions{
		VPNGatewayID: &gID,
		ID:           &gConnID,
	}
	_, response, err := sess.GetVPNGatewayConnection(getVpnGatewayConnectionOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Vpn Gateway Connection(%s): %s\n%s", gConnID, err, response)
	}

	deleteVpnGatewayConnectionOptions := &vpcv1.DeleteVPNGatewayConnectionOptions{
		VPNGatewayID: &gID,
		ID:           &gConnID,
	}
	response, err = sess.DeleteVPNGatewayConnection(deleteVpnGatewayConnectionOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Vpn Gateway Connection : %s\n%s", err, response)
	}

	_, err = isWaitForVPNGatewayConnectionDeleted(sess, gID, gConnID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for Vpn Gateway Connection (%s) is deleted: %s", gConnID, err)
	}

	d.SetId("")
	return nil
}

func isWaitForVPNGatewayConnectionDeleted(vpnGatewayConnection *vpcv1.VpcV1, gID, gConnID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for VPNGatewayConnection (%s) to be deleted.", gConnID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isVPNGatewayConnectionDeleting},
		Target:     []string{"", isVPNGatewayConnectionDeleted},
		Refresh:    isVPNGatewayConnectionDeleteRefreshFunc(vpnGatewayConnection, gID, gConnID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isVPNGatewayConnectionDeleteRefreshFunc(vpnGatewayConnection *vpcv1.VpcV1, gID, gConnID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getVpnGatewayConnectionOptions := &vpcv1.GetVPNGatewayConnectionOptions{
			VPNGatewayID: &gID,
			ID:           &gConnID,
		}
		vpngwcon, response, err := vpnGatewayConnection.GetVPNGatewayConnection(getVpnGatewayConnectionOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return "", isVPNGatewayConnectionDeleted, nil
			}
			return "", "", fmt.Errorf("[ERROR] The Vpn Gateway Connection %s failed to delete: %s\n%s", gConnID, err, response)
		}
		return vpngwcon, isVPNGatewayConnectionDeleting, nil
	}
}

func resourceIBMISVPNGatewayConnectionExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) != 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: ID should be a combination of gID/gConnID", d.Id())
	}

	gID := parts[0]
	gConnID := parts[1]
	exists, err := vpngwconExists(d, meta, gID, gConnID)
	return exists, err
}

func vpngwconExists(d *schema.ResourceData, meta interface{}, gID, gConnID string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	getVpnGatewayConnectionOptions := &vpcv1.GetVPNGatewayConnectionOptions{
		VPNGatewayID: &gID,
		ID:           &gConnID,
	}
	_, response, err := sess.GetVPNGatewayConnection(getVpnGatewayConnectionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting Vpn Gateway Connection: %s\n%s", err, response)
	}
	return true, nil
}

func resourceVPNGatewayConnectionFlattenLifecycleReasons(statusReasons []vpcv1.VPNGatewayConnectionStatusReason) (statusReasonsList []map[string]interface{}) {
	statusReasonsList = make([]map[string]interface{}, 0)
	for _, lr := range statusReasons {
		currentLR := map[string]interface{}{}
		if lr.Code != nil && lr.Message != nil {
			currentLR[isInstanceLifecycleReasonsCode] = *lr.Code
			currentLR[isInstanceLifecycleReasonsMessage] = *lr.Message
			if lr.MoreInfo != nil {
				currentLR[isInstanceLifecycleReasonsMoreInfo] = *lr.MoreInfo
			}
			statusReasonsList = append(statusReasonsList, currentLR)
		}
	}
	return statusReasonsList
}
