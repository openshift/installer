// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package catalogmanagement

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func DataSourceIBMCmObject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCmObjectRead,

		Schema: map[string]*schema.Schema{
			"catalog_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Catalog identifier.",
			},
			"object_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Object identifier.",
			},
			"catalog_object_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "unique id.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The programmatic name of this object.",
			},
			"rev": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloudant revision.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn for this specific object.",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The url for this specific object.",
			},
			"parent_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The parent for this specific object.",
			},
			"label_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A map of translated strings, by language code.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display name in the requested language.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of tags associated with this catalog.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time this catalog was created.",
			},
			"updated": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time this catalog was last updated.",
			},
			"short_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Short description in the requested language.",
			},
			"short_description_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A map of translated strings, by language code.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"kind": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kind of object.",
			},
			"publish": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Publish information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permit_ibm_public_publish": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is it permitted to request publishing to IBM or Public.",
						},
						"ibm_approved": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if this offering has been approved for use by all IBMers.",
						},
						"public_approved": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if this offering has been approved for use by all IBM Cloud users.",
						},
						"portal_approval_record": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The portal's approval record ID.",
						},
						"portal_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The portal UI URL.",
						},
					},
				},
			},
			"state": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Offering state.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "one of: new, validated, account-published, ibm-published, public-published.",
						},
						"current_entered": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date and time of current request.",
						},
						"pending": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "one of: new, validated, account-published, ibm-published, public-published.",
						},
						"pending_requested": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date and time of pending request.",
						},
						"previous": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "one of: new, validated, account-published, ibm-published, public-published.",
						},
					},
				},
			},
			"catalog_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the catalog.",
			},
			"data": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stringified map of data values for this object.",
			},
		},
	}
}

func dataSourceIBMCmObjectRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getObjectOptions := &catalogmanagementv1.GetObjectOptions{}

	getObjectOptions.SetCatalogIdentifier(d.Get("catalog_id").(string))
	getObjectOptions.SetObjectIdentifier(d.Get("object_id").(string))

	catalogObject, response, err := catalogManagementClient.GetObjectWithContext(context, getObjectOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetObjectWithContext failed %s\n%s", err, response), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*getObjectOptions.ObjectIdentifier)

	if err = d.Set("catalog_object_id", catalogObject.ID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_object_id: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("name", catalogObject.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("rev", catalogObject.Rev); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rev: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", catalogObject.CRN); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("url", catalogObject.URL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting url: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("parent_id", catalogObject.ParentID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting parent_id: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if catalogObject.LabelI18n != nil {
		if err = d.Set("label_i18n", catalogObject.LabelI18n); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting label_i18n: %s", err), "(Data) ibm_cm_object", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("label", catalogObject.Label); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting label: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("created", flex.DateTimeToString(catalogObject.Created)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("updated", flex.DateTimeToString(catalogObject.Updated)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("short_description", catalogObject.ShortDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting short_description: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if catalogObject.ShortDescriptionI18n != nil {
		if err = d.Set("short_description_i18n", catalogObject.ShortDescriptionI18n); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting short_description_i18n: %s", err), "(Data) ibm_cm_object", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("kind", catalogObject.Kind); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting kind: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	publish := []map[string]interface{}{}
	if catalogObject.Publish != nil {
		modelMap, err := dataSourceIBMCmObjectPublishObjectToMap(catalogObject.Publish)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_object", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		publish = append(publish, modelMap)
	}
	if err = d.Set("publish", publish); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting publish: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	state := []map[string]interface{}{}
	if catalogObject.State != nil {
		modelMap, err := dataSourceIBMCmObjectStateToMap(catalogObject.State)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_cm_object", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		state = append(state, modelMap)
	}
	if err = d.Set("state", state); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("catalog_name", catalogObject.CatalogName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_name: %s", err), "(Data) ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if catalogObject.Data != nil {
		dataString, err := json.Marshal(catalogObject.Data)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting data, error with json marshal: %s", err), "(Data) ibm_cm_object", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("data", string(dataString)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting data: %s", err), "(Data) ibm_cm_object", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return nil
}

func dataSourceIBMCmObjectPublishObjectToMap(model *catalogmanagementv1.PublishObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PermitIBMPublicPublish != nil {
		modelMap["permit_ibm_public_publish"] = *model.PermitIBMPublicPublish
	}
	if model.IBMApproved != nil {
		modelMap["ibm_approved"] = *model.IBMApproved
	}
	if model.PublicApproved != nil {
		modelMap["public_approved"] = *model.PublicApproved
	}
	return modelMap, nil
}

func dataSourceIBMCmObjectStateToMap(model *catalogmanagementv1.State) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Current != nil {
		modelMap["current"] = *model.Current
	}
	if model.CurrentEntered != nil {
		modelMap["current_entered"] = model.CurrentEntered.String()
	}
	if model.Pending != nil {
		modelMap["pending"] = *model.Pending
	}
	if model.PendingRequested != nil {
		modelMap["pending_requested"] = model.PendingRequested.String()
	}
	if model.Previous != nil {
		modelMap["previous"] = *model.Previous
	}
	return modelMap, nil
}
