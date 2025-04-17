// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMCISDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISDomainRead,

		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_domain",
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
				Computed:    true,
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
				Optional: true,
			},
			cisDomainCnameSuffix: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}
func DataSourceIBMCISDomainValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISDomainValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_domain",
		Schema:       validateSchema}
	return &iBMCISDomainValidator
}
func dataSourceIBMCISDomainRead(d *schema.ResourceData, meta interface{}) error {
	var zoneFound bool
	cisClient, err := meta.(conns.ClientSession).CisZonesV1ClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	cisClient.Crn = core.StringPtr(crn)
	zoneName := d.Get(cisDomain).(string)

	opt := cisClient.NewListZonesOptions()
	opt.SetPage(1)       // list all zones in one page
	opt.SetPerPage(1000) // maximum allowed limit is 1000 per page
	zones, resp, err := cisClient.ListZones(opt)
	if err != nil {
		log.Printf("dataSourcCISdomainRead - ListZones Failed %s\n", resp)
		return err
	}

	for _, zone := range zones.Result {
		if *zone.Name == zoneName {
			d.SetId(flex.ConvertCisToTfTwoVar(*zone.ID, crn))
			d.Set(cisID, crn)
			d.Set(cisDomain, *zone.Name)
			d.Set(cisDomainStatus, *zone.Status)
			d.Set(cisDomainPaused, *zone.Paused)
			d.Set(cisDomainNameServers, zone.NameServers)
			d.Set(cisDomainOriginalNameServers, zone.OriginalNameServers)
			d.Set(cisDomainID, *zone.ID)
			d.Set(cisDomainType, *zone.Type)

			if cisDomainType == "partial" {
				d.Set(cisDomainVerificationKey, *zone.VerificationKey)
				d.Set(cisDomainCnameSuffix, *zone.CnameSuffix)
			}
			zoneFound = true
		}
	}

	if !zoneFound {
		return fmt.Errorf("[ERROR] Given zone does not exist. Please specify correct domain")
	}

	return nil
}
