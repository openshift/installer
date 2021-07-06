// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	cisWAFPackages = "waf_packages"
)

func dataSourceIBMCISWAFPackages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISWAFPackagesRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS Zone CRN",
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "CIS domain id",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisWAFPackages: {
				Type:        schema.TypeList,
				Description: "Collection of GLB Health check/monitor detail",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CIS WAF package id",
						},
						cisWAFPackageID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAF pakcage ID",
						},
						cisWAFPackageName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAF pakcage name",
						},
						cisWAFPackageDetectionMode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAF pakcage detection mode",
						},
						cisWAFPackageDescription: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WAF pakcage description",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMCISWAFPackagesRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisWAFPackageClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)

	opt := cisClient.NewListWafPackagesOptions()
	result, resp, err := cisClient.ListWafPackages(opt)
	if err != nil {
		log.Printf("Error listing waf packages detail: %s", resp)
		return err
	}

	packages := make([]interface{}, 0)
	for _, instance := range result.Result {
		packageDetail := make(map[string]interface{}, 0)
		packageDetail["id"] = convertCisToTfThreeVar(*instance.ID, zoneID, crn)
		packageDetail[cisWAFPackageID] = *instance.ID
		packageDetail[cisWAFPackageName] = *instance.Name
		packageDetail[cisWAFPackageDetectionMode] = *instance.DetectionMode

		if instance.Description != nil {
			packageDetail[cisWAFPackageDescription] = *instance.Description
		}
		packages = append(packages, packageDetail)
	}
	d.SetId(dataSourceIBMCISWAFPackagesCheckID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisWAFPackages, packages)
	return nil
}

func dataSourceIBMCISWAFPackagesCheckID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
