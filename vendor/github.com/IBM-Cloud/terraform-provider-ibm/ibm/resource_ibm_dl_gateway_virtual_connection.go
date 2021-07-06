// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/networking-go-sdk/directlinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	dlGatewaysVirtualConnections = "gateway_vcs"
	dlVCNetworkAccount           = "network_account"
	dlVCNetworkId                = "network_id"
	dlVCName                     = "name"
	dlVCType                     = "type"
	dlVCCreatedAt                = "created_at"
	dlVCStatus                   = "status"
	dlGatewayId                  = "gateway"
	ID                           = "id"
	dlVirtualConnectionId        = "virtual_connection_id"
)

func resourceIBMDLGatewayVC() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMdlGatewayVCCreate,
		Read:     resourceIBMdlGatewayVCRead,
		Delete:   resourceIBMdlGatewayVCDelete,
		Exists:   resourceIBMdlGatewayVCExists,
		Update:   resourceIBMdlGatewayVCUpdate,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			dlGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Direct Link gateway identifier",
			},
			dlVCType: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: InvokeValidator("ibm_dl_virtual_connection", dlVCType),
				Description:  "The type of virtual connection.Allowable values (classic,vpc)",
			},
			dlVCName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: InvokeValidator("ibm_dl_virtual_connection", dlVCName),
				Description:  "The user-defined name for this virtual connection. Virtualconnection names are unique within a gateway. This is the name of thevirtual connection itself, the network being connected may have its ownname attribute",
			},
			dlVCNetworkId: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Unique identifier of the target network. For type=vpc virtual connections this is the CRN of the target VPC. This field does not apply to type=classic connections.",
			},
			dlVCCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time resource was created",
			},
			dlVCStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the virtual connection.Possible values: [pending,attached,approval_pending,rejected,expired,deleting,detached_by_network_pending,detached_by_network]",
			},

			dlVCNetworkAccount: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "For virtual connections across two different IBM Cloud Accounts network_account indicates the account that owns the target network.",
			},
			dlVirtualConnectionId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Direct Gateway virtual connection identifier",
			},

			RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the Direct link gateway",
			},
		},
	}
}
func resourceIBMdlGatewayVCValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	vcType := "classic, vpc"
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 dlVCType,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              vcType})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 dlVCName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^([a-zA-Z]|[a-zA-Z][-_a-zA-Z0-9]*[a-zA-Z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	ibmDLGatewayVCResourceValidator := ResourceValidator{ResourceName: "ibm_dl_virtual_connection", Schema: validateSchema}

	return &ibmDLGatewayVCResourceValidator
}
func resourceIBMdlGatewayVCCreate(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}

	createGatewayVCOptions := &directlinkv1.CreateGatewayVirtualConnectionOptions{}

	gatewayId := d.Get(dlGatewayId).(string)
	createGatewayVCOptions.SetGatewayID(gatewayId)
	vcName := d.Get(dlVCName).(string)
	createGatewayVCOptions.SetName(vcName)
	vcType := d.Get(dlVCType).(string)
	createGatewayVCOptions.SetType(vcType)

	if _, ok := d.GetOk(dlVCNetworkId); ok {
		vcNetworkId := d.Get(dlVCNetworkId).(string)
		createGatewayVCOptions.SetNetworkID(vcNetworkId)
	}

	gatewayVC, response, err := directLink.CreateGatewayVirtualConnection(createGatewayVCOptions)
	if err != nil {
		log.Printf("[DEBUG] Create Direct Link Gateway (Dedicated) Virtual connection err %s\n%s", err, response)
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", gatewayId, *gatewayVC.ID))
	d.Set(dlVirtualConnectionId, *gatewayVC.ID)
	return resourceIBMdlGatewayVCRead(d, meta)
}

func resourceIBMdlGatewayVCRead(d *schema.ResourceData, meta interface{}) error {

	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	ID := parts[1]

	getGatewayVirtualConnectionOptions := &directlinkv1.GetGatewayVirtualConnectionOptions{}
	getGatewayVirtualConnectionOptions.SetGatewayID(gatewayId)
	getGatewayVirtualConnectionOptions.SetID(ID)
	instance, response, err := directLink.GetGatewayVirtualConnection(getGatewayVirtualConnectionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Directlink Gateway Connection (%s): %s\n%s", ID, err, response)
	}

	if instance.Name != nil {
		d.Set(dlVCName, *instance.Name)
	}
	if instance.Type != nil {
		d.Set(dlVCType, *instance.Type)
	}
	if instance.NetworkAccount != nil {
		d.Set(dlVCNetworkAccount, *instance.NetworkAccount)
	}
	if instance.NetworkID != nil {
		d.Set(dlVCNetworkId, *instance.NetworkID)
	}
	if instance.CreatedAt != nil {
		d.Set(dlVCCreatedAt, instance.CreatedAt.String())
	}
	if instance.Status != nil {
		d.Set(dlVCStatus, *instance.Status)
	}
	d.Set(dlVirtualConnectionId, *instance.ID)
	d.Set(dlGatewayId, gatewayId)
	getGatewayOptions := &directlinkv1.GetGatewayOptions{
		ID: &gatewayId,
	}
	dlgw, response, err := directLink.GetGateway(getGatewayOptions)
	if err != nil {
		return fmt.Errorf("Error Getting Direct Link Gateway (Dedicated Template): %s\n%s", err, response)
	}
	d.Set(RelatedCRN, *dlgw.Crn)
	return nil
}

func resourceIBMdlGatewayVCUpdate(d *schema.ResourceData, meta interface{}) error {

	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	ID := parts[1]

	getVCOptions := &directlinkv1.GetGatewayVirtualConnectionOptions{
		ID: &ID,
	}
	getVCOptions.SetGatewayID(gatewayId)
	_, detail, err := directLink.GetGatewayVirtualConnection(getVCOptions)

	if err != nil {
		log.Printf("Error fetching Direct Link Gateway (Dedicated Template) Virtual Connection:%s", detail)
		return err
	}

	updateGatewayVCOptions := &directlinkv1.UpdateGatewayVirtualConnectionOptions{}
	updateGatewayVCOptions.ID = &ID
	updateGatewayVCOptions.SetGatewayID(gatewayId)
	if d.HasChange(dlName) {
		if d.Get(dlName) != nil {
			name := d.Get(dlName).(string)
			updateGatewayVCOptions.Name = &name
		}
	}

	_, response, err := directLink.UpdateGatewayVirtualConnection(updateGatewayVCOptions)
	if err != nil {
		log.Printf("[DEBUG] Update Direct Link Gateway (Dedicated) Virtual Connection err %s\n%s", err, response)
		return err
	}

	return resourceIBMdlGatewayVCRead(d, meta)
}

func resourceIBMdlGatewayVCDelete(d *schema.ResourceData, meta interface{}) error {

	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	ID := parts[1]
	delVCOptions := &directlinkv1.DeleteGatewayVirtualConnectionOptions{
		ID: &ID,
	}
	delVCOptions.SetGatewayID(gatewayId)
	response, err := directLink.DeleteGatewayVirtualConnection(delVCOptions)

	if err != nil && response.StatusCode != 404 {
		log.Printf("Error deleting Direct Link Gateway (Dedicated Template) Virtual Connection: %s", response)
		return err
	}

	d.SetId("")
	return nil
}

func resourceIBMdlGatewayVCExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return false, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	gatewayId := parts[0]
	ID := parts[1]

	getVCOptions := &directlinkv1.GetGatewayVirtualConnectionOptions{
		ID: &ID,
	}
	getVCOptions.SetGatewayID(gatewayId)
	_, response, err := directLink.GetGatewayVirtualConnection(getVCOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return false, nil
		}
		return false, fmt.Errorf("Error Getting Direct Link Gateway (Dedicated Template) Virtual Connection: %s\n%s", err, response)
	}

	if response.StatusCode == 404 {
		d.SetId("")
		return false, nil
	}
	return true, nil
}
