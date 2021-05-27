// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/models"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMIAMServiceID() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMIAMServiceIDCreate,
		Read:     resourceIBMIAMServiceIDRead,
		Update:   resourceIBMIAMServiceIDUpdate,
		Delete:   resourceIBMIAMServiceIDDelete,
		Exists:   resourceIBMIAMServiceIDExists,
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
		},
	}
}

func resourceIBMIAMServiceIDCreate(d *schema.ResourceData, meta interface{}) error {
	iamClient, err := meta.(ClientSession).IAMAPI()
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	userDetails, err := meta.(ClientSession).BluemixUserDetails()

	boundTo := crn.New(userDetails.cloudName, userDetails.cloudType)
	boundTo.ScopeType = crn.ScopeAccount
	boundTo.Scope = userDetails.userAccount

	request := models.ServiceID{
		Name:    name,
		BoundTo: boundTo.String(),
	}

	if des, ok := d.GetOk("description"); ok {
		request.Description = des.(string)
	}

	serviceID, err := iamClient.ServiceIds().Create(request)
	if err != nil {
		return fmt.Errorf("Error creating serviceID: %s", err)
	}

	d.SetId(serviceID.UUID)

	return resourceIBMIAMServiceIDRead(d, meta)
}

func resourceIBMIAMServiceIDRead(d *schema.ResourceData, meta interface{}) error {
	iamClient, err := meta.(ClientSession).IAMAPI()
	if err != nil {
		return err
	}
	serviceIDUUID := d.Id()

	serviceID, err := iamClient.ServiceIds().Get(serviceIDUUID)
	if err != nil {
		return fmt.Errorf("Error retrieving serviceID: %s", err)
	}

	d.Set("name", serviceID.Name)
	d.Set("description", serviceID.Description)
	d.Set("crn", serviceID.CRN)
	d.Set("version", serviceID.Version)
	d.Set("iam_id", serviceID.IAMID)

	return nil
}

func resourceIBMIAMServiceIDUpdate(d *schema.ResourceData, meta interface{}) error {

	iamClient, err := meta.(ClientSession).IAMAPI()
	if err != nil {
		return err
	}
	serviceIDUUID := d.Id()

	hasChange := false
	updateReq := models.ServiceID{}

	if d.HasChange("name") {
		updateReq.Name = d.Get("name").(string)
		hasChange = true
	}

	if d.HasChange("description") {
		updateReq.Description = d.Get("description").(string)
		hasChange = true
	}

	if hasChange {
		_, err = iamClient.ServiceIds().Update(serviceIDUUID, updateReq, "*")
		if err != nil {
			return fmt.Errorf("Error updating serviceID: %s", err)
		}
	}

	return resourceIBMIAMServiceIDRead(d, meta)

}

func resourceIBMIAMServiceIDDelete(d *schema.ResourceData, meta interface{}) error {
	iamClient, err := meta.(ClientSession).IAMAPI()
	if err != nil {
		return err
	}

	serviceIDUUID := d.Id()

	err = iamClient.ServiceIds().Delete(serviceIDUUID)
	if err != nil {
		return fmt.Errorf("Error deleting serviceID: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceIBMIAMServiceIDExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iamClient, err := meta.(ClientSession).IAMAPI()
	if err != nil {
		return false, err
	}
	serviceIDUUID := d.Id()

	serviceID, err := iamClient.ServiceIds().Get(serviceIDUUID)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return serviceID.UUID == serviceIDUUID, nil
}

func CloudName(region models.Region) string {
	regionID := region.ID
	if regionID == "" {
		return ""
	}

	splits := strings.Split(regionID, ":")
	if len(splits) != 3 {
		return ""
	}

	customer := splits[0]
	if customer != "ibm" {
		return customer
	}

	deployment := splits[1]
	switch {
	case deployment == "yp":
		return "bluemix"
	case strings.HasPrefix(deployment, "ys"):
		return "staging"
	default:
		return ""
	}
}

func CloudType(region models.Region) string {
	return region.Type
}

func GenerateBoundToCRN(region models.Region, accountID string) crn.CRN {
	var boundTo crn.CRN
	if region.Type == "dedicated" {
		// cname and ctype are hard coded for dedicated
		boundTo = crn.New("bluemix", "public")
	} else {
		boundTo = crn.New(CloudName(region), CloudType(region))
	}

	boundTo.ScopeType = crn.ScopeAccount
	boundTo.Scope = accountID
	return boundTo
}
