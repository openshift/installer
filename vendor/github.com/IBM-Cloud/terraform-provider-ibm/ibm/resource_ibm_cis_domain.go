// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	cisDomain                    = "domain"
	cisDomainPaused              = "paused"
	cisDomainStatus              = "status"
	cisDomainNameServers         = "name_servers"
	cisDomainOriginalNameServers = "original_name_servers"
)

func resourceIBMCISDomain() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS object id",
				Required:    true,
			},
			cisDomain: {
				Type:        schema.TypeString,
				Description: "CISzone - Domain",
				Required:    true,
			},
			cisDomainPaused: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			cisDomainStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			cisDomainNameServers: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			cisDomainOriginalNameServers: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			cisDomainID: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Create:   resourceCISdomainCreate,
		Read:     resourceCISdomainRead,
		Exists:   resourceCISdomainExists,
		Update:   resourceCISdomainUpdate,
		Delete:   resourceCISdomainDelete,
		Importer: &schema.ResourceImporter{},
	}
}

func resourceCISdomainCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisZonesV1ClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	cisClient.Crn = core.StringPtr(crn)
	zoneName := d.Get(cisDomain).(string)

	opt := cisClient.NewCreateZoneOptions()
	opt.SetName(zoneName)
	result, resp, err := cisClient.CreateZone(opt)
	if err != nil {
		log.Printf("CreateZones Failed %s", resp)
		return err
	}
	d.SetId(convertCisToTfTwoVar(*result.Result.ID, crn))
	return resourceCISdomainRead(d, meta)
}

func resourceCISdomainRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisZonesV1ClientSession()
	if err != nil {
		return err
	}

	zoneID, crn, err := convertTftoCisTwoVar(d.Id())
	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	opt := cisClient.NewGetZoneOptions(zoneID)
	result, resp, err := cisClient.GetZone(opt)
	if err != nil {
		log.Printf("[WARN] Error getting zone %v\n", resp)
		return err
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, result.Result.ID)
	d.Set(cisDomain, result.Result.Name)
	d.Set(cisDomainStatus, result.Result.Status)
	d.Set(cisDomainPaused, result.Result.Paused)
	d.Set(cisDomainNameServers, result.Result.NameServers)
	d.Set(cisDomainOriginalNameServers, result.Result.OriginalNameServers)

	return nil
}
func resourceCISdomainExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(ClientSession).CisZonesV1ClientSession()
	if err != nil {
		return false, err
	}

	zoneID, crn, err := convertTftoCisTwoVar(d.Id())
	log.Println("resource exist :", d.Id())
	if err != nil {
		return false, err
	}
	log.Println("resource exist :", d.Id())
	cisClient.Crn = core.StringPtr(crn)
	opt := cisClient.NewGetZoneOptions(zoneID)
	_, resp, err := cisClient.GetZone(opt)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] zone is not found")
			return false, nil
		}
		log.Printf("[WARN] Error getting zone %v\n", resp)
		return false, err
	}
	return true, nil
}

func resourceCISdomainUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceCISdomainRead(d, meta)
}

func resourceCISdomainDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisZonesV1ClientSession()
	if err != nil {
		return err
	}

	zoneID, crn, err := convertTftoCisTwoVar(d.Id())
	log.Println("resource delete :", d.Id())

	if err != nil {
		return err
	}
	cisClient.Crn = core.StringPtr(crn)
	opt := cisClient.NewGetZoneOptions(zoneID)
	_, resp, err := cisClient.GetZone(opt)
	if err != nil {
		log.Printf("[WARN] Error getting zone %v\n", resp)
		return err
	}
	delOpt := cisClient.NewDeleteZoneOptions(zoneID)
	_, resp, err = cisClient.DeleteZone(delOpt)
	if err != nil {
		log.Printf("[ERR] Error deleting zone %v\n", resp)
		return err
	}
	return nil
}
