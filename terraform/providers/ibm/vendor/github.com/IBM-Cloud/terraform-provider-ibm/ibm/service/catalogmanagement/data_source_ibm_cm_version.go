// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package catalogmanagement

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func DataSourceIBMCmVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCmVersionRead,

		Schema: map[string]*schema.Schema{
			"version_loc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Catalog identifier.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version's CRN.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of content type.",
			},
			"sha": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "hash of the content.",
			},
			"catalog_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Catalog ID.",
			},
			"repo_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Content's repo URL.",
			},
			"source_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Content's source URL (e.g git repo).",
			},
			"tgz_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "File used to on-board this version.",
			},
		},
	}
}

func dataSourceIBMCmVersionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getVersionOptions := &catalogmanagementv1.GetVersionOptions{}

	getVersionOptions.SetVersionLocID(d.Get("version_loc_id").(string))

	offering, response, err := catalogManagementClient.GetVersionWithContext(context, getVersionOptions)
	version := offering.Kinds[0].Versions[0]

	if err != nil {
		log.Printf("[DEBUG] GetVersionWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(*version.VersionLocator)
	if err = d.Set("crn", version.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	if err = d.Set("version", version.Version); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting version: %s", err))
	}
	if err = d.Set("sha", version.Sha); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting sha: %s", err))
	}
	if err = d.Set("catalog_id", version.CatalogID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting catalog_id: %s", err))
	}
	if err = d.Set("repo_url", version.RepoURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting repo_url: %s", err))
	}
	if err = d.Set("source_url", version.SourceURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting source_url: %s", err))
	}
	if err = d.Set("tgz_url", version.TgzURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting tgz_url: %s", err))
	}

	return nil
}
