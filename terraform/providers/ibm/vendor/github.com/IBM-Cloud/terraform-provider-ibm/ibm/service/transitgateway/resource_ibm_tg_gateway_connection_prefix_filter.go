// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package transitgateway

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMTransitGatewayConnectionPrefixFilter() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMTransitGatewayConnectionPrefixFilterCreate,
		Read:     resourceIBMTransitGatewayConnectionPrefixFilterRead,
		Delete:   resourceIBMTransitGatewayConnectionPrefixFilterDelete,
		Update:   resourceIBMTransitGatewayConnectionPrefixFilterUpdate,
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
				Description: "The Transit Gateway Connection identifier",
			},
			tgPrefixFilterId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Transit Gateway Connection Prefix Filter identifier",
			},
			tgAction: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tg_connection_prefix_filter", tgAction),
				Description:  "Whether to permit or deny the prefix filter",
			},
			tgBefore: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifier of prefix filter that handles ordering",
			},
			tgCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this prefix filter was created",
			},
			tgGe: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "IP Prefix GE",
			},
			tgLe: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "IP Prefix LE",
			},
			tgPrefix: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP Prefix",
			},
			tgUpdatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this prefix filter was last updated",
			},
		},
	}
}

func ResourceIBMTransitGatewayConnectionPrefixFilterValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	actionValues := "permit, deny"
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 tgAction,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              actionValues})

	ibmTransitGatewayConnectionPrefixFilterResourceValidator := validate.ResourceValidator{ResourceName: "ibm_tg_connection_prefix_filter", Schema: validateSchema}

	return &ibmTransitGatewayConnectionPrefixFilterResourceValidator
}

func resourceIBMTransitGatewayConnectionPrefixFilterCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	gatewayId := d.Get(tgGatewayId).(string)
	connectionId := d.Get(tgConnectionId).(string)

	createPrefixFilterOptions := &transitgatewayapisv1.CreateTransitGatewayConnectionPrefixFilterOptions{}
	createPrefixFilterOptions.SetTransitGatewayID(gatewayId)
	createPrefixFilterOptions.SetID(connectionId)

	action := d.Get(tgAction).(string)
	createPrefixFilterOptions.SetAction(action)
	prefix := d.Get(tgPrefix).(string)
	createPrefixFilterOptions.SetPrefix(prefix)

	if _, ok := d.GetOk(tgBefore); ok {
		before := d.Get(tgBefore).(string)
		createPrefixFilterOptions.SetBefore(before)
	}
	if _, ok := d.GetOk(tgGe); ok {
		ge := int64(d.Get(tgGe).(int))
		createPrefixFilterOptions.SetGe(ge)
	}
	if _, ok := d.GetOk(tgLe); ok {
		le := int64(d.Get(tgLe).(int))
		createPrefixFilterOptions.SetLe(le)
	}

	prefixFilter, response, err := client.CreateTransitGatewayConnectionPrefixFilter(createPrefixFilterOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Create Transit Gateway connection prefix filter err %s\n%s", err, response)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", gatewayId, connectionId, *prefixFilter.ID))

	return resourceIBMTransitGatewayConnectionPrefixFilterRead(d, meta)
}

func resourceIBMTransitGatewayConnectionPrefixFilterRead(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	connectionId := parts[1]
	filterId := parts[2]

	getTransitGatewayConnectionPrefixFilterOptionsModel := &transitgatewayapisv1.GetTransitGatewayConnectionPrefixFilterOptions{}
	getTransitGatewayConnectionPrefixFilterOptionsModel.SetTransitGatewayID(gatewayId)
	getTransitGatewayConnectionPrefixFilterOptionsModel.SetID(connectionId)
	getTransitGatewayConnectionPrefixFilterOptionsModel.SetFilterID(filterId)
	prefixFilter, response, err := client.GetTransitGatewayConnectionPrefixFilter(getTransitGatewayConnectionPrefixFilterOptionsModel)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error while retrieving transit gateway connection prefix filter (%s): %s\n%s", filterId, err, response)
	}

	d.Set(tgPrefixFilterId, *prefixFilter.ID)
	d.Set(tgCreatedAt, prefixFilter.CreatedAt.String())
	d.Set(tgPrefix, prefixFilter.Prefix)

	if prefixFilter.UpdatedAt != nil {
		d.Set(tgUpdatedAt, prefixFilter.UpdatedAt.String())
	}
	if prefixFilter.Before != nil {
		d.Set(tgBefore, prefixFilter.Before)
	}
	if prefixFilter.Ge != nil {
		d.Set(tgGe, prefixFilter.Ge)
	}
	if prefixFilter.Le != nil {
		d.Set(tgLe, prefixFilter.Le)
	}

	return nil
}

func resourceIBMTransitGatewayConnectionPrefixFilterUpdate(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	gatewayId := parts[0]
	connectionId := parts[1]
	filterId := parts[2]

	updatePrefixFilterOptions := &transitgatewayapisv1.UpdateTransitGatewayConnectionPrefixFilterOptions{}
	updatePrefixFilterOptions.SetTransitGatewayID(gatewayId)
	updatePrefixFilterOptions.SetID(connectionId)
	updatePrefixFilterOptions.SetFilterID(filterId)

	if d.HasChange(tgAction) {
		if d.Get(tgAction) != nil {
			action := d.Get(tgAction).(string)
			updatePrefixFilterOptions.SetAction(action)
		}
	}
	if d.HasChange(tgBefore) {
		if d.Get(tgBefore) != nil {
			before := d.Get(tgBefore).(string)
			updatePrefixFilterOptions.SetBefore(before)
		}
	}
	if d.HasChange(tgGe) {
		if d.Get(tgGe) != nil {
			ge := int64(d.Get(tgGe).(int))
			updatePrefixFilterOptions.SetGe(ge)
		}
	}
	if d.HasChange(tgLe) {
		if d.Get(tgLe) != nil {
			le := int64(d.Get(tgLe).(int))
			updatePrefixFilterOptions.SetLe(le)
		}
	}
	if d.HasChange(tgPrefix) {
		if d.Get(tgPrefix) != nil {
			prefix := d.Get(tgPrefix).(string)
			updatePrefixFilterOptions.SetPrefix(prefix)
		}
	}

	_, response, err := client.UpdateTransitGatewayConnectionPrefixFilter(updatePrefixFilterOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error in Update Transit Gateway Connection Prefix Filter (%s): %s\n%s", filterId, err, response)
	}

	return resourceIBMTransitGatewayConnectionPrefixFilterRead(d, meta)
}

func resourceIBMTransitGatewayConnectionPrefixFilterDelete(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	gatewayId := parts[0]
	connectionId := parts[1]
	filterId := parts[2]

	deletePrefixFilterOptions := &transitgatewayapisv1.DeleteTransitGatewayConnectionPrefixFilterOptions{}
	deletePrefixFilterOptions.SetTransitGatewayID(gatewayId)
	deletePrefixFilterOptions.SetID(connectionId)
	deletePrefixFilterOptions.SetFilterID(filterId)

	response, err := client.DeleteTransitGatewayConnectionPrefixFilter(deletePrefixFilterOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error deleting Transit Gateway Connection Prefix Filter(%s): %s\n%s", filterId, err, response)
	}

	d.SetId("")
	return nil
}
