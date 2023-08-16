// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisDomain                    = "domain"
	cisDomainPaused              = "paused"
	cisDomainStatus              = "status"
	cisDomainNameServers         = "name_servers"
	cisDomainOriginalNameServers = "original_name_servers"
	cisDomainType                = "type"
	cisDomainVerificationKey     = "verification_key"
	cisDomainCnameSuffix         = "cname_suffix"
	ibmCISDomain                 = "ibm_cis_domain"
)

func ResourceIBMCISDomain() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeValidator("ibm_cis_domain",
					"cis_id"),
			},
			cisDomain: {
				Type:        schema.TypeString,
				Description: "CISzone - Domain",
				Required:    true,
			},
			cisDomainType: {
				Type:        schema.TypeString,
				Description: "CISzone - Domain Type",
				Default:     "full",
				Optional:    true,
				ValidateFunc: validate.InvokeValidator(ibmCISDomain,
					cisDomainType),
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
			cisDomainVerificationKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			cisDomainCnameSuffix: {
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
	cisClient, err := meta.(conns.ClientSession).CisZonesV1ClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	cisClient.Crn = core.StringPtr(crn)
	zoneName := d.Get(cisDomain).(string)
	zoneType := d.Get(cisDomainType).(string)

	opt := cisClient.NewCreateZoneOptions()
	opt.SetName(zoneName)
	opt.SetType(zoneType)

	result, resp, err := cisClient.CreateZone(opt)
	if err != nil {
		log.Printf("CreateZones Failed %s", resp)
		return err
	}
	d.SetId(flex.ConvertCisToTfTwoVar(*result.Result.ID, crn))
	return resourceCISdomainRead(d, meta)
}

func resourceCISdomainRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisZonesV1ClientSession()
	if err != nil {
		return err
	}

	zoneID, crn, _ := flex.ConvertTftoCisTwoVar(d.Id())
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
	d.Set(cisDomainType, result.Result.Type)

	if cisDomainType == "partial" {
		d.Set(cisDomainVerificationKey, result.Result.VerificationKey)
		d.Set(cisDomainCnameSuffix, result.Result.CnameSuffix)
	}

	return nil
}
func resourceCISdomainExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(conns.ClientSession).CisZonesV1ClientSession()
	if err != nil {
		return false, err
	}

	zoneID, crn, _ := flex.ConvertTftoCisTwoVar(d.Id())
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
	cisClient, err := meta.(conns.ClientSession).CisZonesV1ClientSession()
	if err != nil {
		return err
	}

	zoneID, crn, _ := flex.ConvertTftoCisTwoVar(d.Id())
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

func ResourceIBMCISDomainValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 cisDomainType,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "full, partial"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	ibmCISDomainResourceValidator := validate.ResourceValidator{
		ResourceName: ibmCISDomain,
		Schema:       validateSchema}
	return &ibmCISDomainResourceValidator
}
