// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package transitgateway

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	tgXacGatewayId     = "gateway"
	tgXacConnectionId  = "connection_id"
	tgConnectionAction = "action"
)

func ResourceIBMTransitGatewayConnectionAction() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMTransitGatewayConnectionActionCreate,
		Read:     resourceIBMTransitGatewayConnectionActionRead,
		Delete:   resourceIBMTransitGatewayConnectionActionDelete,
		Importer: &schema.ResourceImporter{},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					// Allow reusing the resource block to approve a new connection
					if diff.HasChange(tgConnectionId) {
						return nil
					}
					if diff.HasChange(tgConnectionAction) {
						o, n := diff.GetChange(tgConnectionAction)
						oldAction := o.(string)
						newAction := n.(string)
						if oldAction == "" {
							// oldAction is empty when performing an action on a connection for the first time
							// We are only concerned with checking if the action changes for existing action/connection resource
							return nil
						}
						return fmt.Errorf("The action for the transit gateway connection has already been performed and cannot be changed from %s to %s", oldAction, newAction)
					}
					return nil
				}),
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
			tgConnectionAction: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tg_connection_action", tgConnectionAction),
				Description:  "The Transit Gateway Connection cross account action",
			},
		},
	}
}

func ResourceIBMTransitGatewayConnectionActionValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	actions := "approve, reject"
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 tgConnectionAction,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              actions})

	ibmTransitGatewayConnectionActionResourceValidator := validate.ResourceValidator{ResourceName: "ibm_tg_connection_action", Schema: validateSchema}

	return &ibmTransitGatewayConnectionActionResourceValidator
}

func resourceIBMTransitGatewayConnectionActionCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	createTransitGatewayConnectionActionsOptions := &transitgatewayapisv1.CreateTransitGatewayConnectionActionsOptions{}
	gatewayId := d.Get(tgXacGatewayId).(string)
	createTransitGatewayConnectionActionsOptions.SetTransitGatewayID(gatewayId)
	connectionId := d.Get(tgXacConnectionId).(string)
	createTransitGatewayConnectionActionsOptions.SetID(connectionId)
	action := d.Get(tgConnectionAction).(string)
	createTransitGatewayConnectionActionsOptions.SetAction(action)

	response, err := client.CreateTransitGatewayConnectionActions(createTransitGatewayConnectionActionsOptions)
	if err != nil {
		return fmt.Errorf("Error performing an action on the Transit Gateway Connection: %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s", gatewayId, connectionId))
	d.Set(tgConnectionId, connectionId)
	return resourceIBMTransitGatewayConnectionActionRead(d, meta)
}

func resourceIBMTransitGatewayConnectionActionRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	gatewayId := parts[0]
	ID := parts[1]

	d.Set(tgConnectionId, ID)
	d.Set(tgGatewayId, gatewayId)

	return nil
}

func resourceIBMTransitGatewayConnectionActionDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
