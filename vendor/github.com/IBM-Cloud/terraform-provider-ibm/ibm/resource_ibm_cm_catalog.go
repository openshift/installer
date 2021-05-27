// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func resourceIBMCmCatalog() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCmCatalogCreate,
		Read:     resourceIBMCmCatalogRead,
		Delete:   resourceIBMCmCatalogDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Display Name in the requested language.",
			},
			"short_description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description in the requested language.",
			},
			"catalog_icon_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "URL for an icon associated with this catalog.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
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

func resourceIBMCmCatalogCreate(d *schema.ResourceData, meta interface{}) error {
	catalogManagementClient, err := meta.(ClientSession).CatalogManagementV1()
	if err != nil {
		return err
	}

	createCatalogOptions := &catalogmanagementv1.CreateCatalogOptions{}

	if _, ok := d.GetOk("label"); ok {
		createCatalogOptions.SetLabel(d.Get("label").(string))
	}
	if _, ok := d.GetOk("short_description"); ok {
		createCatalogOptions.SetShortDescription(d.Get("short_description").(string))
	}
	if _, ok := d.GetOk("catalog_icon_url"); ok {
		createCatalogOptions.SetCatalogIconURL(d.Get("catalog_icon_url").(string))
	}
	if _, ok := d.GetOk("tags"); ok {
		createCatalogOptions.SetTags(d.Get("tags").([]string))
	}

	catalog, response, err := catalogManagementClient.CreateCatalog(createCatalogOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateCatalog failed %s\n%s", err, response)
		return err
	}

	d.SetId(*catalog.ID)

	return resourceIBMCmCatalogRead(d, meta)
}

func resourceIBMCmCatalogRead(d *schema.ResourceData, meta interface{}) error {
	catalogManagementClient, err := meta.(ClientSession).CatalogManagementV1()
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] client is a nil pointer: %v\n", catalogManagementClient == nil)

	getCatalogOptions := &catalogmanagementv1.GetCatalogOptions{}

	getCatalogOptions.SetCatalogIdentifier(d.Id())

	catalog, response, err := catalogManagementClient.GetCatalog(getCatalogOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetCatalog failed %s\n%s", err, response)
		return err
	}
	if err = d.Set("label", catalog.Label); err != nil {
		return fmt.Errorf("Error setting label: %s", err)
	}
	if err = d.Set("short_description", catalog.ShortDescription); err != nil {
		return fmt.Errorf("Error setting short_description: %s", err)
	}
	if err = d.Set("catalog_icon_url", catalog.CatalogIconURL); err != nil {
		return fmt.Errorf("Error setting catalog_icon_url: %s", err)
	}
	if catalog.Tags != nil {
		if err = d.Set("tags", catalog.Tags); err != nil {
			return fmt.Errorf("Error setting tags: %s", err)
		}
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

func resourceIBMCmCatalogDelete(d *schema.ResourceData, meta interface{}) error {
	catalogManagementClient, err := meta.(ClientSession).CatalogManagementV1()
	if err != nil {
		return err
	}

	deleteCatalogOptions := &catalogmanagementv1.DeleteCatalogOptions{}

	deleteCatalogOptions.SetCatalogIdentifier(d.Id())

	response, err := catalogManagementClient.DeleteCatalog(deleteCatalogOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteCatalog failed %s\n%s", err, response)
		return err
	}

	d.SetId("")

	return nil
}
