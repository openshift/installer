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

func dataSourceIBMCmOffering() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCmOfferingRead,

		Schema: map[string]*schema.Schema{
			"catalog_identifier": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Catalog identifier.",
				ForceNew:    true,
			},
			"offering_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the catalog containing this offering.",
				ForceNew:    true,
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The url for this specific offering.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn for this specific offering.",
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display Name in the requested language.",
				ForceNew:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The programmatic name of this offering.",
			},
			"offering_icon_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for an icon associated with this offering.",
			},
			"offering_docs_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for an additional docs with this offering.",
			},
			"offering_support_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL to be displayed in the Consumption UI for getting support on this offering.",
			},
			"short_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Short description in the requested language.",
			},
			"long_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Long description in the requested language.",
			},
			"permit_request_ibm_public_publish": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is it permitted to request publishing to IBM or Public.",
			},
			"ibm_publish_approved": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if this offering has been approved for use by all IBMers.",
			},
			"public_publish_approved": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if this offering has been approved for use by all IBM Cloud users.",
			},
			"public_original_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The original offering CRN that this publish entry came from.",
			},
			"publish_public_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the public catalog entry of this offering.",
			},
			"portal_approval_record": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The portal's approval record ID.",
			},
			"portal_ui_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The portal UI URL.",
			},
			"catalog_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the catalog containing this offering.",
			},
			"catalog_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the catalog.",
			},
			"disclaimer": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A disclaimer for this offering.",
			},
			"hidden": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determine if this offering should be displayed in the Consumption UI.",
			},
			"repo_info": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Repository info for offerings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Token for private repos.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public or enterprise GitHub.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMCmOfferingRead(d *schema.ResourceData, meta interface{}) error {
	catalogManagementClient, err := meta.(ClientSession).CatalogManagementV1()
	if err != nil {
		return err
	}

	getOfferingOptions := &catalogmanagementv1.GetOfferingOptions{}

	getOfferingOptions.SetCatalogIdentifier(d.Get("catalog_identifier").(string))
	getOfferingOptions.SetOfferingID(d.Get("offering_id").(string))

	offering, response, err := catalogManagementClient.GetOfferingWithContext(context.TODO(), getOfferingOptions)
	if err != nil {
		log.Printf("[DEBUG] GetOfferingWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId(*offering.ID)
	if err = d.Set("url", offering.URL); err != nil {
		return fmt.Errorf("Error setting url: %s", err)
	}
	if err = d.Set("crn", offering.CRN); err != nil {
		return fmt.Errorf("Error setting crn: %s", err)
	}
	if err = d.Set("label", offering.Label); err != nil {
		return fmt.Errorf("Error setting label: %s", err)
	}
	if err = d.Set("name", offering.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err = d.Set("offering_icon_url", offering.OfferingIconURL); err != nil {
		return fmt.Errorf("Error setting offering_icon_url: %s", err)
	}
	if err = d.Set("offering_docs_url", offering.OfferingDocsURL); err != nil {
		return fmt.Errorf("Error setting offering_docs_url: %s", err)
	}
	if err = d.Set("offering_support_url", offering.OfferingSupportURL); err != nil {
		return fmt.Errorf("Error setting offering_support_url: %s", err)
	}
	if err = d.Set("short_description", offering.ShortDescription); err != nil {
		return fmt.Errorf("Error setting short_description: %s", err)
	}
	if err = d.Set("long_description", offering.LongDescription); err != nil {
		return fmt.Errorf("Error setting long_description: %s", err)
	}
	if err = d.Set("permit_request_ibm_public_publish", offering.PermitRequestIBMPublicPublish); err != nil {
		return fmt.Errorf("Error setting permit_request_ibm_public_publish: %s", err)
	}
	if err = d.Set("ibm_publish_approved", offering.IBMPublishApproved); err != nil {
		return fmt.Errorf("Error setting ibm_publish_approved: %s", err)
	}
	if err = d.Set("public_publish_approved", offering.PublicPublishApproved); err != nil {
		return fmt.Errorf("Error setting public_publish_approved: %s", err)
	}
	if err = d.Set("public_original_crn", offering.PublicOriginalCRN); err != nil {
		return fmt.Errorf("Error setting public_original_crn: %s", err)
	}
	if err = d.Set("publish_public_crn", offering.PublishPublicCRN); err != nil {
		return fmt.Errorf("Error setting publish_public_crn: %s", err)
	}
	if err = d.Set("portal_approval_record", offering.PortalApprovalRecord); err != nil {
		return fmt.Errorf("Error setting portal_approval_record: %s", err)
	}
	if err = d.Set("portal_ui_url", offering.PortalUIURL); err != nil {
		return fmt.Errorf("Error setting portal_ui_url: %s", err)
	}
	if err = d.Set("catalog_id", offering.CatalogID); err != nil {
		return fmt.Errorf("Error setting catalog_id: %s", err)
	}
	if err = d.Set("catalog_name", offering.CatalogName); err != nil {
		return fmt.Errorf("Error setting catalog_name: %s", err)
	}
	if err = d.Set("disclaimer", offering.Disclaimer); err != nil {
		return fmt.Errorf("Error setting disclaimer: %s", err)
	}
	if err = d.Set("hidden", offering.Hidden); err != nil {
		return fmt.Errorf("Error setting hidden: %s", err)
	}

	if offering.RepoInfo != nil {
		repoInfoMap := dataSourceOfferingRepoInfoToMap(*offering.RepoInfo)
		if err = d.Set("repo_info", []map[string]interface{}{repoInfoMap}); err != nil {
			return fmt.Errorf("Error setting repo_info %s", err)
		}
	}

	return nil
}

func dataSourceOfferingRepoInfoToMap(repoInfoItem catalogmanagementv1.RepoInfo) (repoInfoMap map[string]interface{}) {
	repoInfoMap = map[string]interface{}{}

	if repoInfoItem.Token != nil {
		repoInfoMap["token"] = repoInfoItem.Token
	}
	if repoInfoItem.Type != nil {
		repoInfoMap["type"] = repoInfoItem.Type
	}

	return repoInfoMap
}
