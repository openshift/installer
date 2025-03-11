// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cloudfoundry

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMSpace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMSpaceRead,

		Schema: map[string]*schema.Schema{
			"space": {
				Description:  "Space name, for example dev",
				Type:         schema.TypeString,
				Optional:     true,
				Deprecated:   "use name instead",
				ExactlyOneOf: []string{"space", "name"},
			},
			"name": {
				Description:  "Space name, for example dev",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"space", "name"},
			},
			"org": {
				Description: "The org this space belongs to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"auditors": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The IBMID of the users who  have auditor role in this space, ex - user@example.com",
			},
			"managers": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The IBMID of the users who  have manager role in this space, ex - user@example.com",
			},
			"developers": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The IBMID of the users who  have developer role in this space, ex - user@example.com",
			},
		},
		DeprecationMessage: "This service is deprecated.",
	}
}

func dataSourceIBMSpaceRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(conns.ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	orgAPI := cfClient.Organizations()
	spaceAPI := cfClient.Spaces()
	var space string
	if v, ok := d.GetOk("name"); ok {
		space = v.(string)
	}
	if v, ok := d.GetOk("space"); ok {
		space = v.(string)
	}

	org := d.Get("org").(string)

	orgFields, err := orgAPI.FindByName(org, conns.BluemixRegion)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving org: %s", err)
	}
	spaceFields, err := spaceAPI.FindByNameInOrg(orgFields.GUID, space, conns.BluemixRegion)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving space: %s", err)
	}

	spaceGUID := spaceFields.GUID
	d.SetId(spaceGUID)

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

	return nil
}
