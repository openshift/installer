// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamaccessgroup

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIAMAccessGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIAMAccessGroupCreate,
		ReadContext:   resourceIBMIAMAccessGroupRead,
		UpdateContext: resourceIBMIAMAccessGroupUpdate,
		DeleteContext: resourceIBMIAMAccessGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the access group",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the access group",
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of the access group",
			},
		},
	}
}

func resourceIBMIAMAccessGroupCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	creatAccessGroupOptions := iamAccessGroupsClient.NewCreateAccessGroupOptions(userDetails.UserAccount, name)
	if des, ok := d.GetOk("description"); ok {
		description := des.(string)
		creatAccessGroupOptions.Description = &description
	}
	agrp, detailedResponse, err := iamAccessGroupsClient.CreateAccessGroup(creatAccessGroupOptions)
	if err != nil || agrp == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error creating access group: %s. API Response: %s", err, detailedResponse))
	}

	d.SetId(*agrp.ID)

	return resourceIBMIAMAccessGroupRead(context, d, meta)
}

func resourceIBMIAMAccessGroupRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}
	agrpID := d.Id()
	getAccessGroupOptions := iamAccessGroupsClient.NewGetAccessGroupOptions(agrpID)
	getAccessGroupOptions.SetShowCRN(true)
	var agrp *iamaccessgroupsv2.Group
	var detailedResponse *core.DetailedResponse
	err = resource.RetryContext(context, 5*time.Second, func() *resource.RetryError {
		agrp, detailedResponse, err = iamAccessGroupsClient.GetAccessGroup(getAccessGroupOptions)
		if err != nil || agrp == nil {
			if detailedResponse != nil && detailedResponse.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if conns.IsResourceTimeoutError(err) {
		agrp, detailedResponse, err = iamAccessGroupsClient.GetAccessGroup(getAccessGroupOptions)
	}
	if err != nil || agrp == nil {
		if detailedResponse != nil && detailedResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving access group: %s. API Response: %s", err, detailedResponse))
	}
	version := detailedResponse.GetHeaders().Get("etag")
	d.Set("name", agrp.Name)
	d.Set("description", agrp.Description)
	d.Set("version", version)
	d.Set("crn", agrp.CRN)
	return nil
}

func resourceIBMIAMAccessGroupUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}
	agrpID := d.Id()

	hasChange := false
	version := d.Get("version").(string)
	updateAccessGroupOptions := iamAccessGroupsClient.NewUpdateAccessGroupOptions(agrpID, version)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateAccessGroupOptions.Name = &name
		hasChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateAccessGroupOptions.Description = &description
		hasChange = true
	}

	if hasChange {
		agrp, detailedResponse, err := iamAccessGroupsClient.UpdateAccessGroup(updateAccessGroupOptions)
		if err != nil || agrp == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating access group: %s. API Response: %s", err, detailedResponse))
		}
	}

	return resourceIBMIAMAccessGroupRead(context, d, meta)

}

func resourceIBMIAMAccessGroupDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	agID := d.Id()
	force := true
	deleteAccessGroupOptions := iamAccessGroupsClient.NewDeleteAccessGroupOptions(agID)
	deleteAccessGroupOptions.SetForce(force)
	detailedResponse, err := iamAccessGroupsClient.DeleteAccessGroup(deleteAccessGroupOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error deleting access group: %s, API Response: %s", err, detailedResponse))
	}

	d.SetId("")

	return nil
}
