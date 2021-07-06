// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/api/iamuum/iamuumv2"
	"github.com/IBM-Cloud/bluemix-go/models"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMIAMAccessGroup() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMIAMAccessGroupCreate,
		Read:     resourceIBMIAMAccessGroupRead,
		Update:   resourceIBMIAMAccessGroupUpdate,
		Delete:   resourceIBMIAMAccessGroupDelete,
		Exists:   resourceIBMIAMAccessGroupExists,
		Importer: &schema.ResourceImporter{},

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
		},
	}
}

func resourceIBMIAMAccessGroupCreate(d *schema.ResourceData, meta interface{}) error {
	iamuumClient, err := meta.(ClientSession).IAMUUMAPIV2()
	if err != nil {
		return err
	}

	request := models.AccessGroupV2{
		AccessGroup: models.AccessGroup{
			Name: d.Get("name").(string),
		},
	}

	if des, ok := d.GetOk("description"); ok {
		request.Description = des.(string)
	}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	agrp, err := iamuumClient.AccessGroup().Create(request, userDetails.userAccount)
	if err != nil {
		return fmt.Errorf("Error creating access group: %s", err)
	}

	d.SetId(agrp.ID)

	return resourceIBMIAMAccessGroupRead(d, meta)
}

func resourceIBMIAMAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	iamuumClient, err := meta.(ClientSession).IAMUUMAPIV2()
	if err != nil {
		return err
	}
	agrpID := d.Id()

	agrp, version, err := iamuumClient.AccessGroup().Get(agrpID)
	if err != nil {
		return fmt.Errorf("Error retrieving access group: %s", err)
	}

	d.Set("name", agrp.Name)
	d.Set("description", agrp.Description)
	d.Set("version", version)

	return nil
}

func resourceIBMIAMAccessGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	iamuumClient, err := meta.(ClientSession).IAMUUMAPIV2()
	if err != nil {
		return err
	}
	agrpID := d.Id()

	hasChange := false
	updateReq := iamuumv2.AccessGroupUpdateRequest{}

	if d.HasChange("name") {
		updateReq.Name = d.Get("name").(string)
		hasChange = true
	}

	if d.HasChange("description") {
		updateReq.Description = d.Get("description").(string)
		hasChange = true
	}

	if hasChange {
		_, err = iamuumClient.AccessGroup().Update(agrpID, updateReq, d.Get("version").(string))
		if err != nil {
			return fmt.Errorf("Error updating access group: %s", err)
		}
	}

	return resourceIBMIAMAccessGroupRead(d, meta)

}

func resourceIBMIAMAccessGroupDelete(d *schema.ResourceData, meta interface{}) error {
	iamuumClient, err := meta.(ClientSession).IAMUUMAPIV2()
	if err != nil {
		return err
	}

	agID := d.Id()

	err = iamuumClient.AccessGroup().Delete(agID, true)
	if err != nil {
		return fmt.Errorf("Error deleting access group: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceIBMIAMAccessGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iamuumClient, err := meta.(ClientSession).IAMUUMAPIV2()
	if err != nil {
		return false, err
	}
	agID := d.Id()

	agrp, _, err := iamuumClient.AccessGroup().Get(agID)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return agrp.ID == agID, nil
}
