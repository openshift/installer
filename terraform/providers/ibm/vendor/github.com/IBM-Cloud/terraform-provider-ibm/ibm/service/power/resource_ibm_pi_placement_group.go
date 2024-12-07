// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMPIPlacementGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIPlacementGroupCreate,
		ReadContext:   resourceIBMPIPlacementGroupRead,
		UpdateContext: resourceIBMPIPlacementGroupUpdate,
		DeleteContext: resourceIBMPIPlacementGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_PlacementGroupName: {
				Description:  "The name of the placement group.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_PlacementGroupPolicy: {
				Description:  "The value of the group's affinity policy. Valid values are 'affinity' and 'anti-affinity'.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{Affinity, AntiAffinity}),
			},

			// Attributes
			Attr_Members: {
				Computed:    true,
				Description: "The list of server instances IDs that are members of the placement group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeSet,
			},
			Attr_PlacementGroupID: {
				Computed:    true,
				Description: "The placement group ID.",
				Type:        schema.TypeString,
			},
		},
	}
}

func resourceIBMPIPlacementGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	name := d.Get(Arg_PlacementGroupName).(string)
	policy := d.Get(Arg_PlacementGroupPolicy).(string)
	client := instance.NewIBMPIPlacementGroupClient(ctx, sess, cloudInstanceID)
	body := &models.PlacementGroupCreate{
		Name:   &name,
		Policy: &policy,
	}

	response, err := client.Create(body)
	if err != nil || response == nil {
		return diag.FromErr(fmt.Errorf("error creating the shared processor pool: %s", err))
	}

	log.Printf("Printing the placement group %+v", &response)

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *response.ID))
	return resourceIBMPIPlacementGroupRead(ctx, d, meta)
}

func resourceIBMPIPlacementGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := parts[0]
	client := instance.NewIBMPIPlacementGroupClient(ctx, sess, cloudInstanceID)

	response, err := client.Get(parts[1])
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return diag.FromErr(err)
	}

	d.Set(Arg_PlacementGroupName, response.Name)
	d.Set(Arg_PlacementGroupPolicy, response.Policy)
	d.Set(Attr_Members, response.Members)
	d.Set(Attr_PlacementGroupID, response.ID)

	return nil
}

func resourceIBMPIPlacementGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceIBMPIPlacementGroupRead(ctx, d, meta)
}

func resourceIBMPIPlacementGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := parts[0]
	client := instance.NewIBMPIPlacementGroupClient(ctx, sess, cloudInstanceID)
	err = client.Delete(parts[1])

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
