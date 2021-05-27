// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMSecurityGroupCreate,
		Read:     resourceIBMSecurityGroupRead,
		Update:   resourceIBMSecurityGroupUpdate,
		Delete:   resourceIBMSecurityGroupDelete,
		Exists:   resourceIBMSecurityGroupExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "Security group name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Security group description",
			},
		},
	}
}

func resourceIBMSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSecurityGroupService(sess.SetRetries(0))

	name := d.Get("name").(string)
	var description string
	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	} else {
		description = ""
	}

	groups, err := services.GetAccountService(sess).
		Filter(filter.Path("securityGroups.name").Eq(name).Build()).
		GetSecurityGroups()

	if nil == err && len(groups) > 0 {
		group := groups[0]
		id := *group.Id
		d.SetId(fmt.Sprintf("%d", id))
		editSG := false

		oldDesc := ""
		if group.Description != nil {
			oldDesc = *group.Description
		}
		if oldDesc != description {
			group.Description = sl.String(description)
			editSG = true
		}

		if editSG {
			_, err = service.EditObject(&group)
			return err
		}

		return nil
	}
	sg := &datatypes.Network_SecurityGroup{
		Name:        sl.String(name),
		Description: sl.String(description),
	}
	res, err := service.CreateObject(sg)
	if err != nil {
		return fmt.Errorf("Error creating Security Group: %s", err)
	}

	d.SetId(strconv.Itoa(*res.Id))
	log.Printf("[INFO] Security Group: %d", *res.Id)

	return resourceIBMSecurityGroupRead(d, meta)
}

func resourceIBMSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSecurityGroupService(sess)

	groupID, _ := strconv.Atoi(d.Id())
	group, err := service.Id(groupID).GetObject()
	if err != nil {
		// If the group is somehow already destroyed, mark as
		// succesfully gone
		if err, ok := err.(sl.Error); ok && err.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Security Group: %s", err)
	}

	d.Set("name", group.Name)
	d.Set("description", group.Description)
	return nil
}

func resourceIBMSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSecurityGroupService(sess)

	groupID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	group, err := services.GetNetworkSecurityGroupService(sess).Id(groupID).GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving Security Group: %s", err)
	}

	if d.HasChange("description") {
		group.Description = sl.String(d.Get("description").(string))
	}

	if d.HasChange("name") {
		group.Name = sl.String(d.Get("name").(string))
	}
	_, err = service.Id(groupID).EditObject(&group)
	if err != nil {
		return fmt.Errorf("Error editing Security Group: %s", err)
	}
	return resourceIBMSecurityGroupRead(d, meta)
}

func resourceIBMSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSecurityGroupService(sess)

	groupID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	log.Printf("[INFO] Deleting Security Group: %d", groupID)
	_, err = service.Id(groupID).DeleteObject()
	if err != nil {
		return fmt.Errorf("Error deleting Security Group: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMSecurityGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetNetworkSecurityGroupService(sess)

	groupID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(groupID).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}
	return result.Id != nil && *result.Id == groupID, nil
}
