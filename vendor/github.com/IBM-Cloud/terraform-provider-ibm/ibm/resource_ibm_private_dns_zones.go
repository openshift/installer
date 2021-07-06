// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	pdnsInstanceID      = "instance_id"
	pdnsZoneName        = "name"
	pdnsZoneDescription = "description"
	pdnsZoneLabel       = "label"
	pdnsZoneCreatedOn   = "created_on"
	pdnsZoneModifiedOn  = "modified_on"
	pdnsZoneState       = "state"
	pdnsZoneID          = "zone_id"
)

func resourceIBMPrivateDNSZone() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMPrivateDNSZoneCreate,
		Read:     resourceIBMPrivateDNSZoneRead,
		Update:   resourceIBMPrivateDNSZoneUpdate,
		Delete:   resourceIBMPrivateDNSZoneDelete,
		Exists:   resourceIBMPrivateDNSZoneExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID",
			},

			pdnsZoneID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone ID",
			},

			pdnsZoneName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone name",
			},

			pdnsZoneDescription: {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Zone description",
			},

			pdnsZoneState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone state",
			},

			pdnsZoneLabel: {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Label",
			},

			pdnsZoneCreatedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date",
			},

			pdnsZoneModifiedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification date",
			},
		},
	}
}

func resourceIBMPrivateDNSZoneCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	var (
		instanceID      string
		zoneName        string
		zoneDescription string
		zoneLabel       string
	)

	instanceID = d.Get(pdnsInstanceID).(string)
	zoneName = d.Get(pdnsZoneName).(string)
	if v, ok := d.GetOk(pdnsZoneDescription); ok {
		zoneDescription = v.(string)
	}
	if v, ok := d.GetOk(pdnsZoneLabel); ok {
		zoneLabel = v.(string)
	}
	createZoneOptions := sess.NewCreateDnszoneOptions(instanceID)
	createZoneOptions.SetName(zoneName)
	createZoneOptions.SetDescription(zoneDescription)
	createZoneOptions.SetLabel(zoneLabel)
	response, detail, err := sess.CreateDnszone(createZoneOptions)
	if err != nil {
		return fmt.Errorf("Error creating pdns zone:%s\n%s", err, detail)
	}

	d.SetId(fmt.Sprintf("%s/%s", *response.InstanceID, *response.ID))
	d.Set(pdnsZoneID, *response.ID)

	return resourceIBMPrivateDNSZoneRead(d, meta)
}

func resourceIBMPrivateDNSZoneRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	idSet := strings.Split(d.Id(), "/")
	getZoneOptions := sess.NewGetDnszoneOptions(idSet[0], idSet[1])
	response, detail, err := sess.GetDnszone(getZoneOptions)
	if err != nil {
		return fmt.Errorf("Error fetching pdns zone:%s\n%s", err, detail)
	}

	d.Set(pdnsZoneID, response.ID)
	d.Set(pdnsInstanceID, response.InstanceID)
	d.Set(pdnsZoneName, response.Name)
	d.Set(pdnsZoneDescription, response.Description)
	d.Set(pdnsZoneLabel, response.Label)
	d.Set(pdnsZoneCreatedOn, response.CreatedOn)
	d.Set(pdnsZoneModifiedOn, response.ModifiedOn)
	d.Set(pdnsZoneState, response.State)

	return nil
}

func resourceIBMPrivateDNSZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	idSet := strings.Split(d.Id(), "/")

	// Check DNS zone is present?
	getZoneOptions := sess.NewGetDnszoneOptions(idSet[0], idSet[1])
	_, response, err := sess.GetDnszone(getZoneOptions)
	if err != nil {
		return fmt.Errorf("Error fetching pdns zone:%s\n%s", err, response)
	}

	// Update DNS zone if attributes has any change

	if d.HasChange(pdnsZoneLabel) || d.HasChange(pdnsZoneDescription) {
		updateZoneOptions := sess.NewUpdateDnszoneOptions(idSet[0], idSet[1])
		description := d.Get(pdnsZoneDescription).(string)
		label := d.Get(pdnsZoneLabel).(string)
		updateZoneOptions.SetDescription(description)
		updateZoneOptions.SetLabel(label)

		_, detail, err := sess.UpdateDnszone(updateZoneOptions)

		if err != nil {
			return fmt.Errorf("Error updating pdns zone:%s\n%s", err, detail)
		}
	}

	return resourceIBMPrivateDNSZoneRead(d, meta)
}

func resourceIBMPrivateDNSZoneDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	idSet := strings.Split(d.Id(), "/")

	deleteZoneOptions := sess.NewDeleteDnszoneOptions(idSet[0], idSet[1])
	response, err := sess.DeleteDnszone(deleteZoneOptions)
	if err != nil {
		return fmt.Errorf("Error deleting pdns zone:%s\n%s", err, response)
	}

	d.SetId("")
	return nil
}

func resourceIBMPrivateDNSZoneExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return false, err
	}

	idSet := strings.Split(d.Id(), "/")
	getZoneOptions := sess.NewGetDnszoneOptions(idSet[0], idSet[1])
	_, response, err := sess.GetDnszone(getZoneOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
