// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMDNSSecondary() *schema.Resource {
	return &schema.Resource{
		Exists:   resourceIBMDNSSecondaryExists,
		Create:   resourceIBMDNSSecondaryCreate,
		Read:     resourceIBMDNSSecondaryRead,
		Update:   resourceIBMDNSSecondaryUpdate,
		Delete:   resourceIBMDNSSecondaryDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"master_ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Master IP Address",
			},

			"transfer_frequency": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Transfer frequency value",
			},

			"zone_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone name",
			},

			"status_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status ID",
			},

			"status_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status text",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of tags",
			},
		},
	}
}

func resourceIBMDNSSecondaryCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetDnsSecondaryService(sess)

	// prepare creation parameters
	opts := datatypes.Dns_Secondary{
		MasterIpAddress:   sl.String(d.Get("master_ip_address").(string)),
		TransferFrequency: sl.Int(d.Get("transfer_frequency").(int)),
		ZoneName:          sl.String(d.Get("zone_name").(string)),
	}

	// create Dns_Secondary object
	response, err := service.CreateObject(&opts)
	if err != nil {
		return fmt.Errorf("Error creating Dns Secondary Zone: %s", err)
	}

	// populate id
	id := *response.Id
	d.SetId(strconv.Itoa(id))
	log.Printf("[INFO] Created Dns Secondary Zone: %d", id)

	// read remote state
	return resourceIBMDNSSecondaryRead(d, meta)
}

func resourceIBMDNSSecondaryRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetDnsSecondaryService(sess)

	dnsId, _ := strconv.Atoi(d.Id())

	// retrieve remote object state
	dns_domain_secondary, err := service.Id(dnsId).GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving Dns Secondary Zone %d: %s", dnsId, err)
	}

	// populate fields
	d.Set("master_ip_address", *dns_domain_secondary.MasterIpAddress)
	d.Set("transfer_frequency", *dns_domain_secondary.TransferFrequency)
	d.Set("zone_name", *dns_domain_secondary.ZoneName)
	d.Set("status_id", *dns_domain_secondary.StatusId)
	d.Set("status_text", *dns_domain_secondary.StatusText)

	return nil
}

func resourceIBMDNSSecondaryUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	domainId, _ := strconv.Atoi(d.Id())
	hasChange := false

	opts := datatypes.Dns_Secondary{}
	if d.HasChange("master_ip_address") {
		opts.MasterIpAddress = sl.String(d.Get("master_ip_address").(string))
		hasChange = true
	}

	if d.HasChange("transfer_frequency") {
		opts.TransferFrequency = sl.Int(d.Get("transfer_frequency").(int))
		hasChange = true
	}

	if hasChange {
		service := services.GetDnsSecondaryService(sess)
		_, err := service.Id(domainId).EditObject(&opts)

		if err != nil {
			return fmt.Errorf("Error editing DNS secondary zone (%d): %s", domainId, err)
		}
	}

	return nil
}

func resourceIBMDNSSecondaryDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetDnsSecondaryService(sess)

	dnsId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting Dns Secondary Zone: %s", err)
	}

	log.Printf("[INFO] Deleting Dns Secondary Zone: %d", dnsId)
	result, err := service.Id(dnsId).DeleteObject()
	if err != nil {
		return fmt.Errorf("Error deleting Dns Secondary Zone: %s", err)
	}

	if !result {
		return errors.New("Error deleting Dns Secondary Zone")
	}

	d.SetId("")
	return nil
}

func resourceIBMDNSSecondaryExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetDnsSecondaryService(sess)

	dnsId, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(dnsId).GetObject()
	return err == nil && result.Id != nil && *result.Id == dnsId, nil
}
