// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cloudfoundry

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/api/mccp/mccpv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMSpace() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMSpaceCreate,
		Read:     resourceIBMSpaceRead,
		Update:   resourceIBMSpaceUpdate,
		Delete:   resourceIBMSpaceDelete,
		Exists:   resourceIBMSpaceExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name for the space",
			},
			"org": {
				Description: "The org this space belongs to",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"auditors": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The IBMID of the users who will have auditor role in this space, ex - user@example.com",
				Set:         flex.ResourceIBMVPCHash,
			},
			"managers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The IBMID of the users who will have manager role in this space, ex - user@example.com",
				Set:         flex.ResourceIBMVPCHash,
			},
			"developers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The IBMID of the users who will have developer role in this space, ex - user@example.com",
				Set:         flex.ResourceIBMVPCHash,
			},
			"space_quota": {
				Description: "The name of the Space Quota Definition",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
		DeprecationMessage: "This service is deprecated.",
	}
}

func resourceIBMSpaceCreate(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(conns.ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	org := d.Get("org").(string)
	name := d.Get("name").(string)

	req := mccpv2.SpaceCreateRequest{
		Name: name,
	}

	orgFields, err := cfClient.Organizations().FindByName(org, conns.BluemixRegion)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving org: %s", err)
	}
	req.OrgGUID = orgFields.GUID

	if spaceQuota, ok := d.GetOk("space_quota"); ok {
		quota, err := cfClient.SpaceQuotas().FindByName(spaceQuota.(string), orgFields.GUID)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving space quota: %s", err)
		}
		req.SpaceQuotaGUID = quota.GUID
	}

	spaceAPI := cfClient.Spaces()
	space, err := spaceAPI.Create(req)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating space: %s", err)
	}

	spaceGUID := space.Metadata.GUID
	d.SetId(spaceGUID)

	if developerSet := d.Get("developers").(*schema.Set); len(developerSet.List()) > 0 {
		developers := flex.ExpandStringList(developerSet.List())
		for _, d := range developers {
			_, err := spaceAPI.AssociateDeveloper(spaceGUID, d)
			if err != nil {
				return fmt.Errorf("[ERROR] Error associating developer %s with space %s : %s", d, spaceGUID, err)
			}
		}
	}

	if auditorSet := d.Get("auditors").(*schema.Set); len(auditorSet.List()) > 0 {
		auditors := flex.ExpandStringList(auditorSet.List())
		for _, d := range auditors {
			_, err := spaceAPI.AssociateAuditor(spaceGUID, d)
			if err != nil {
				return fmt.Errorf("[ERROR] Error associating auditor %s with space %s : %s", d, spaceGUID, err)
			}
		}

	}
	if managerSet := d.Get("managers").(*schema.Set); len(managerSet.List()) > 0 {
		managers := flex.ExpandStringList(managerSet.List())
		for _, d := range managers {
			_, err := spaceAPI.AssociateManager(spaceGUID, d)
			if err != nil {
				return fmt.Errorf("[ERROR] Error associating manager %s with space %s : %s", d, spaceGUID, err)
			}
		}
	}

	return resourceIBMSpaceRead(d, meta)
}

func resourceIBMSpaceRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(conns.ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	spaceGUID := d.Id()

	spaceAPI := cfClient.Spaces()
	orgAPI := cfClient.Organizations()
	spaceDetails, err := spaceAPI.Get(spaceGUID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving space: %s", err)
	}

	auditors, err := spaceAPI.ListAuditors(spaceGUID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving auditors in the space: %s", err)
	}

	managers, err := spaceAPI.ListManagers(spaceGUID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving managers in the space: %s", err)
	}

	developers, err := spaceAPI.ListDevelopers(spaceGUID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving developers in space: %s", err)
	}

	d.Set("auditors", flex.FlattenSpaceRoleUsers(auditors))
	d.Set("managers", flex.FlattenSpaceRoleUsers(managers))
	d.Set("developers", flex.FlattenSpaceRoleUsers(developers))

	if spaceDetails.Entity.SpaceQuotaGUID != "" {
		sqAPI := cfClient.SpaceQuotas()
		quota, err := sqAPI.Get(spaceDetails.Entity.SpaceQuotaGUID)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving quotas details for space: %s", err)
		}
		d.Set("space_quota", quota.Entity.Name)
	}
	d.Set("name", spaceDetails.Entity.Name)
	org, err := orgAPI.Get(spaceDetails.Entity.OrgGUID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving Organization details for space: %s", err)
	}
	d.Set("org", org.Entity.Name)
	return nil
}

func resourceIBMSpaceUpdate(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(conns.ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	id := d.Id()

	req := mccpv2.SpaceUpdateRequest{}
	if d.HasChange("name") {
		req.Name = helpers.String(d.Get("name").(string))
	}

	api := cfClient.Spaces()
	_, err = api.Update(id, req)
	if err != nil {
		return fmt.Errorf("[ERROR] Error updating space: %s", err)
	}

	err = updateAuditors(api, id, d)
	if err != nil {
		return err
	}
	err = updateManagers(api, id, d)
	if err != nil {
		return err
	}
	err = updateDevelopers(api, id, d)
	if err != nil {
		return err
	}
	return resourceIBMSpaceRead(d, meta)
}

func resourceIBMSpaceDelete(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(conns.ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	id := d.Id()

	err = cfClient.Spaces().Delete(id, false)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting space: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMSpaceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cfClient, err := meta.(conns.ClientSession).MccpAPI()
	if err != nil {
		return false, err
	}
	id := d.Id()

	space, err := cfClient.Spaces().Get(id)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("[ERROR] Error getting space: %s", err)
	}

	return space.Metadata.GUID == id, nil
}

func updateDevelopers(api mccpv2.Spaces, spaceGUID string, d *schema.ResourceData) error {
	if !d.HasChange("developers") {
		return nil
	}
	var remove, add []string
	o, n := d.GetChange("developers")
	os := o.(*schema.Set)
	ns := n.(*schema.Set)

	remove = flex.ExpandStringList(os.Difference(ns).List())
	add = flex.ExpandStringList(ns.Difference(os).List())

	if len(add) > 0 {
		for _, d := range add {
			_, err := api.AssociateDeveloper(spaceGUID, d)
			if err != nil {
				return fmt.Errorf("[ERROR] Error associating developer %s with space %s : %s", d, spaceGUID, err)
			}
		}
	}
	if len(remove) > 0 {
		for _, d := range remove {
			err := api.DisassociateDeveloper(spaceGUID, d)
			if err != nil {
				return fmt.Errorf("[ERROR] Error dis-associating developer %s with space %s : %s", d, spaceGUID, err)
			}
		}
	}
	return nil
}

func updateManagers(api mccpv2.Spaces, spaceGUID string, d *schema.ResourceData) error {
	if !d.HasChange("managers") {
		return nil
	}
	var remove, add []string
	o, n := d.GetChange("managers")
	os := o.(*schema.Set)
	ns := n.(*schema.Set)

	remove = flex.ExpandStringList(os.Difference(ns).List())
	add = flex.ExpandStringList(ns.Difference(os).List())

	if len(add) > 0 {
		for _, d := range add {
			_, err := api.AssociateManager(spaceGUID, d)
			if err != nil {
				return fmt.Errorf("[ERROR] Error associating manager %s with space %s : %s", d, spaceGUID, err)
			}
		}
	}
	if len(remove) > 0 {
		for _, d := range remove {
			err := api.DisassociateManager(spaceGUID, d)
			if err != nil {
				return fmt.Errorf("[ERROR] Error dis-associating manager %s with space %s : %s", d, spaceGUID, err)
			}
		}
	}
	return nil
}
func updateAuditors(api mccpv2.Spaces, spaceGUID string, d *schema.ResourceData) error {
	if !d.HasChange("auditors") {
		return nil
	}
	var remove, add []string
	o, n := d.GetChange("auditors")
	os := o.(*schema.Set)
	ns := n.(*schema.Set)

	remove = flex.ExpandStringList(os.Difference(ns).List())
	add = flex.ExpandStringList(ns.Difference(os).List())

	if len(add) > 0 {
		for _, d := range add {
			_, err := api.AssociateAuditor(spaceGUID, d)
			if err != nil {
				return fmt.Errorf("[ERROR] Error associating auditor %s with space %s : %s", d, spaceGUID, err)
			}
		}
	}
	if len(remove) > 0 {
		for _, d := range remove {
			err := api.DisassociateAuditor(spaceGUID, d)
			if err != nil {
				return fmt.Errorf("[ERROR] Error dis-associating auditor %s with space %s : %s", d, spaceGUID, err)
			}
		}
	}
	return nil
}
