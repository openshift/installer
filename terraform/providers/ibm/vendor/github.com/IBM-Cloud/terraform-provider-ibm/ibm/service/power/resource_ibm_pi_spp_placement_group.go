// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	models "github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func ResourceIBMPISPPPlacementGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPISPPPlacementGroupCreate,
		ReadContext:   resourceIBMPISPPPlacementGroupRead,
		DeleteContext: resourceIBMPISPPPlacementGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			Arg_SPPPlacementGroupName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the SPP placement group",
			},

			Arg_SPPPlacementGroupPolicy: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"affinity", "anti-affinity"}),
				Description:  "Policy of the SPP placement group",
			},

			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "PI cloud instance ID",
			},

			Attr_SPPPlacementGroupMembers: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Member SPP IDs that are the SPP placement group members",
			},

			Attr_SPPPlacementGroupID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SPP placement group ID",
			},
		},
	}
}

func resourceIBMPISPPPlacementGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	name := d.Get(Arg_SPPPlacementGroupName).(string)
	policy := d.Get(Arg_SPPPlacementGroupPolicy).(string)
	client := st.NewIBMPISPPPlacementGroupClient(ctx, sess, cloudInstanceID)
	body := &models.SPPPlacementGroupCreate{
		Name:   &name,
		Policy: &policy,
	}

	response, err := client.Create(body)
	if err != nil || response == nil {
		return diag.Errorf("error creating the spp placement group: %v", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *response.ID))
	return resourceIBMPISPPPlacementGroupRead(ctx, d, meta)
}

func resourceIBMPISPPPlacementGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := parts[0]
	client := st.NewIBMPISPPPlacementGroupClient(ctx, sess, cloudInstanceID)

	response, err := client.Get(parts[1])
	if err != nil || response == nil {
		return diag.Errorf("error reading the spp placement group: %v", err)
	}

	d.Set(Arg_CloudInstanceID, cloudInstanceID)
	d.Set(Arg_SPPPlacementGroupName, response.Name)
	d.Set(Attr_SPPPlacementGroupID, response.ID)
	d.Set(Arg_SPPPlacementGroupPolicy, response.Policy)
	d.Set(Attr_SPPPlacementGroupMembers, response.MemberSharedProcessorPools)

	return nil

}

func resourceIBMPISPPPlacementGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := parts[0]
	client := st.NewIBMPISPPPlacementGroupClient(ctx, sess, cloudInstanceID)
	err = client.Delete(parts[1])

	if err != nil {
		return diag.Errorf("error deleting the spp placement group: %v", err)
	}
	d.SetId("")
	return nil
}
