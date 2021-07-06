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

func dataSourceIBMCmCatalog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCmCatalogRead,

		Schema: map[string]*schema.Schema{
			"catalog_identifier": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID for catalog",
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display Name in the requested language.",
			},
			"short_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description in the requested language.",
			},
			"catalog_icon_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for an icon associated with this catalog.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of tags associated with this catalog.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The url for this specific catalog.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN associated with the catalog.",
			},
			"offerings_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL path to offerings.",
			},
		},
	}
}

func dataSourceIBMCmCatalogRead(d *schema.ResourceData, meta interface{}) error {
	catalogManagementClient, err := meta.(ClientSession).CatalogManagementV1()
	if err != nil {
		return err
	}

	getCatalogOptions := &catalogmanagementv1.GetCatalogOptions{}

	getCatalogOptions.SetCatalogIdentifier(d.Get("catalog_identifier").(string))

	catalog, response, err := catalogManagementClient.GetCatalogWithContext(context.TODO(), getCatalogOptions)
	if err != nil {
		log.Printf("[DEBUG] GetCatalogWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId(*catalog.ID)
	if err = d.Set("label", catalog.Label); err != nil {
		return fmt.Errorf("Error setting label: %s", err)
	}
	if err = d.Set("short_description", catalog.ShortDescription); err != nil {
		return fmt.Errorf("Error setting short_description: %s", err)
	}
	if err = d.Set("catalog_icon_url", catalog.CatalogIconURL); err != nil {
		return fmt.Errorf("Error setting catalog_icon_url: %s", err)
	}
	if err = d.Set("tags", catalog.Tags); err != nil {
		return fmt.Errorf("Error setting tags: %s", err)
	}
	if err = d.Set("url", catalog.URL); err != nil {
		return fmt.Errorf("Error setting url: %s", err)
	}
	if err = d.Set("crn", catalog.CRN); err != nil {
		return fmt.Errorf("Error setting crn: %s", err)
	}
	if err = d.Set("offerings_url", catalog.OfferingsURL); err != nil {
		return fmt.Errorf("Error setting offerings_url: %s", err)
	}
	return nil
}
