// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIAMServiceID() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIAMServiceIDCreate,
		ReadContext:   resourceIBMIAMServiceIDRead,
		UpdateContext: resourceIBMIAMServiceIDUpdate,
		DeleteContext: resourceIBMIAMServiceIDDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the serviceID",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the serviceID",
			},

			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "version of the serviceID",
			},

			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "crn of the serviceID",
			},

			"iam_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IAM ID of the serviceID",
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceIBMIAMServiceIDCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	createServiceIDOptions := iamidentityv1.CreateServiceIDOptions{
		Name:      &name,
		AccountID: &userDetails.UserAccount,
	}

	if d, ok := d.GetOk("description"); ok {
		des := d.(string)
		createServiceIDOptions.Description = &des
	}

	serviceID, resp, err := iamIdentityClient.CreateServiceID(&createServiceIDOptions)
	if err != nil || serviceID == nil {
		log.Printf("Error creating serviceID: %s, %s", err, resp)
		return diag.FromErr(fmt.Errorf("[ERROR] Error creating serviceID: %s %s", err, resp))
	}
	d.SetId(*serviceID.ID)

	return resourceIBMIAMServiceIDRead(context, d, meta)
}

func resourceIBMIAMServiceIDRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	serviceIDUUID := d.Id()
	getServiceIDOptions := iamidentityv1.GetServiceIDOptions{
		ID: &serviceIDUUID,
	}
	serviceID, resp, err := iamIdentityClient.GetServiceID(&getServiceIDOptions)
	if err != nil || serviceID == nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("Error retrieving serviceID: %s %s", err, resp)
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving serviceID: %s %s", err, resp))
	}
	if serviceID.Name != nil {
		d.Set("name", *serviceID.Name)
	}
	if serviceID.Description != nil {
		d.Set("description", *serviceID.Description)
	}
	if serviceID.CRN != nil {
		d.Set("crn", *serviceID.CRN)
	}
	if serviceID.EntityTag != nil {
		d.Set("version", serviceID.EntityTag)
	}
	if serviceID.IamID != nil {
		d.Set("iam_id", serviceID.IamID)
	}
	if serviceID.Locked != nil {
		d.Set("locked", serviceID.Locked)
	}
	return nil
}

func resourceIBMIAMServiceIDUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	serviceIDUUID := d.Id()

	hasChange := false
	ifMatch := "*"
	updateServiceIDOptions := iamidentityv1.UpdateServiceIDOptions{
		ID:      &serviceIDUUID,
		IfMatch: &ifMatch,
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateServiceIDOptions.Name = &name
		hasChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateServiceIDOptions.Description = &description
		hasChange = true
	}

	if hasChange {
		_, resp, err := iamIdentityClient.UpdateServiceID(&updateServiceIDOptions)
		if err != nil {
			log.Printf("Error updating serviceID: %s, %s", err, resp)
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating serviceID: %s %s", err, resp))
		}
	}

	return resourceIBMIAMServiceIDRead(context, d, meta)

}

func resourceIBMIAMServiceIDDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	serviceIDUUID := d.Id()
	deleteServiceIDOptions := iamidentityv1.DeleteServiceIDOptions{
		ID: &serviceIDUUID,
	}
	resp, err := iamIdentityClient.DeleteServiceID(&deleteServiceIDOptions)
	if err != nil {
		log.Printf("Error deleting serviceID: %s %s", err, resp)
		return diag.FromErr(fmt.Errorf("[ERROR] Error deleting serviceID: %s %s", err, resp))
	}

	d.SetId("")

	return nil
}
