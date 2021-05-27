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

func resourceIBMDNSDomain() *schema.Resource {
	return &schema.Resource{
		Exists:   resourceIBMDNSDomainExists,
		Create:   resourceIBMDNSDomainCreate,
		Read:     resourceIBMDNSDomainRead,
		Update:   resourceIBMDNSDomainUpdate,
		Delete:   resourceIBMDNSDomainDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DNS name",
			},

			"serial": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS serial info",
			},

			"update_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS update date",
			},

			"target": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DNS target info",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "tags associated with reosurce.",
			},
		},
	}
}

func resourceIBMDNSDomainCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetDnsDomainService(sess.SetRetries(0))

	// prepare creation parameters
	opts := datatypes.Dns_Domain{
		Name: sl.String(d.Get("name").(string)),
	}

	opts.ResourceRecords = []datatypes.Dns_Domain_ResourceRecord{}

	if targetString, ok := d.GetOk("target"); ok {
		opts.ResourceRecords = []datatypes.Dns_Domain_ResourceRecord{
			{
				Data: sl.String(targetString.(string)),
				Host: sl.String("@"),
				Ttl:  sl.Int(86400),
				Type: sl.String("a"),
			},
		}
	}

	// create Dns_Domain object
	response, err := service.CreateObject(&opts)
	if err != nil {
		return fmt.Errorf("Error creating Dns Domain: %s", err)
	}

	// populate id
	id := *response.Id
	d.SetId(strconv.Itoa(id))
	log.Printf("[INFO] Created Dns Domain: %d", id)

	// read remote state
	return resourceIBMDNSDomainRead(d, meta)
}

func resourceIBMDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetDnsDomainService(sess)

	dnsId, _ := strconv.Atoi(d.Id())

	// retrieve remote object state
	dns_domain, err := service.Id(dnsId).Mask(
		"id,name,updateDate,resourceRecords",
	).GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving Dns Domain %d: %s", dnsId, err)
	}

	// populate fields
	d.Set("name", dns_domain.Name)
	d.Set("serial", sl.Get(dns_domain.Serial, nil))
	d.Set("update_date", sl.Get(dns_domain.UpdateDate.String(), nil))

	// find a record with host @; that will have the current target.
	for _, record := range dns_domain.ResourceRecords {
		if *record.Type == "a" && *record.Host == "@" {
			d.Set("target", *record.Data)
			break
		}
	}

	return nil
}

func resourceIBMDNSDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	// If the target has been updated, find the corresponding dns record and update its data

	sess := meta.(ClientSession).SoftLayerSession()
	domainService := services.GetDnsDomainService(sess)
	service := services.GetDnsDomainResourceRecordService(sess.SetRetries(0))

	domainId, _ := strconv.Atoi(d.Id())

	if !d.HasChange("target") { // target is the only editable field
		return nil
	}

	newTarget := d.Get("target").(string)

	// retrieve domain state
	domain, err := domainService.Id(domainId).Mask(
		"id,name,updateDate,resourceRecords",
	).GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving DNS resource %d: %s", domainId, err)
	}

	// find a record with host @; that will have the current target.
	var record datatypes.Dns_Domain_ResourceRecord
	for _, record = range domain.ResourceRecords {
		if *record.Type == "a" && *record.Host == "@" {
			break
		}
	}

	if record.Id == nil {
		return fmt.Errorf("Could not find DNS target record for domain %s (%d)",
			sl.Get(domain.Name), sl.Get(domain.Id))
	}

	record.Data = sl.String(newTarget)

	_, err = service.Id(*record.Id).EditObject(&record)

	if err != nil {
		return fmt.Errorf("Error editing DNS target record for domain %s (%d): %s",
			sl.Get(domain.Name), sl.Get(domain.Id), err)
	}

	return nil
}

func resourceIBMDNSDomainDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetDnsDomainService(sess)

	dnsId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting Dns Domain: %s", err)
	}

	log.Printf("[INFO] Deleting Dns Domain: %d", dnsId)
	result, err := service.Id(dnsId).DeleteObject()
	if err != nil {
		return fmt.Errorf("Error deleting Dns Domain: %s", err)
	}

	if !result {
		return errors.New("Error deleting Dns Domain")
	}

	d.SetId("")
	return nil
}

func resourceIBMDNSDomainExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetDnsDomainService(sess)

	dnsId, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(dnsId).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error retrieving domain info: %s", err)
	}
	return result.Id != nil && *result.Id == dnsId, nil
}
