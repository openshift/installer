// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func dataSourceIBMCmVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCmVersionRead,

		Schema: map[string]*schema.Schema{
			"version_loc_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Catalog identifier.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version's CRN.",
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of content type.",
			},
			"sha": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "hash of the content.",
			},
			"catalog_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Catalog ID.",
			},
			"repo_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Content's repo URL.",
			},
			"source_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Content's source URL (e.g git repo).",
			},
			"tgz_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "File used to on-board this version.",
			},
		},
	}
}

func dataSourceIBMCmVersionRead(d *schema.ResourceData, meta interface{}) error {
	catalogManagementClient, err := meta.(ClientSession).CatalogManagementV1()
	if err != nil {
		return err
	}

	getVersionOptions := &catalogmanagementv1.GetVersionOptions{}

	getVersionOptions.SetVersionLocID(d.Get("version_loc_id").(string))

	offering, response, err := catalogManagementClient.GetVersionWithContext(context.TODO(), getVersionOptions)
	version := offering.Kinds[0].Versions[0]

	if err != nil {
		log.Printf("[DEBUG] GetVersionWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId(*version.VersionLocator)
	if err = d.Set("crn", version.CRN); err != nil {
		return fmt.Errorf("Error setting crn: %s", err)
	}
	if err = d.Set("version", version.Version); err != nil {
		return fmt.Errorf("Error setting version: %s", err)
	}
	if err = d.Set("sha", version.Sha); err != nil {
		return fmt.Errorf("Error setting sha: %s", err)
	}
	if err = d.Set("catalog_id", version.CatalogID); err != nil {
		return fmt.Errorf("Error setting catalog_id: %s", err)
	}
	if err = d.Set("repo_url", version.RepoURL); err != nil {
		return fmt.Errorf("Error setting repo_url: %s", err)
	}
	if err = d.Set("source_url", version.SourceURL); err != nil {
		return fmt.Errorf("Error setting source_url: %s", err)
	}
	if err = d.Set("tgz_url", version.TgzURL); err != nil {
		return fmt.Errorf("Error setting tgz_url: %s", err)
	}

	return nil
}
