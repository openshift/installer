// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMIAMServiceID() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMIAMServiceIDCreate,
		Read:     resourceIBMIAMServiceIDRead,
		Update:   resourceIBMIAMServiceIDUpdate,
		Delete:   resourceIBMIAMServiceIDDelete,
		Importer: &schema.ResourceImporter{},

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

func resourceIBMIAMServiceIDCreate(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	createServiceIDOptions := iamidentityv1.CreateServiceIDOptions{
		Name:      &name,
		AccountID: &userDetails.userAccount,
	}

	if d, ok := d.GetOk("description"); ok {
		des := d.(string)
		createServiceIDOptions.Description = &des
	}

	serviceID, resp, err := iamIdentityClient.CreateServiceID(&createServiceIDOptions)
	if err != nil || serviceID == nil {
		log.Printf("Error creating serviceID: %s, %s", err, resp)
		return err
	}
	d.SetId(*serviceID.ID)

	return resourceIBMIAMServiceIDRead(d, meta)
}

func resourceIBMIAMServiceIDRead(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
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
		return err
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

func resourceIBMIAMServiceIDUpdate(d *schema.ResourceData, meta interface{}) error {

	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
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
			return err
		}
	}

	return resourceIBMIAMServiceIDRead(d, meta)

}

func resourceIBMIAMServiceIDDelete(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	serviceIDUUID := d.Id()
	deleteServiceIDOptions := iamidentityv1.DeleteServiceIDOptions{
		ID: &serviceIDUUID,
	}
	resp, err := iamIdentityClient.DeleteServiceID(&deleteServiceIDOptions)
	if err != nil {
		log.Printf("Error deleting serviceID: %s %s", err, resp)
		return err
	}

	d.SetId("")

	return nil
}
