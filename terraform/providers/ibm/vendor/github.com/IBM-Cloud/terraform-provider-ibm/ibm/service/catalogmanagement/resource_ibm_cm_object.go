// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package catalogmanagement

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func ResourceIBMCmObject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCmObjectCreate,
		ReadContext:   resourceIBMCmObjectRead,
		UpdateContext: resourceIBMCmObjectUpdate,
		DeleteContext: resourceIBMCmObjectDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"catalog_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Catalog identifier.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The programmatic name of this object.",
			},
			"kind": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Kind of object. Options are \"vpe\", \"preset_configuration\", or \"proxy_source\".",
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
				Optional:    true,
				Description: "The parent for this specific object.",
			},
			"label_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A map of translated strings, by language code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Display name in the requested language.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of tags associated with this catalog.",
				Elem:        &schema.Schema{Type: schema.TypeString},
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
				Optional:    true,
				Description: "Short description in the requested language.",
			},
			"short_description_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A map of translated strings, by language code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
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
				Optional:    true,
				Computed:    true,
				Description: "Stringified map of data values for this object.",
			},
			"rev": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloudant revision.",
			},
			"object_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the object.",
			},
		},
	}
}

func resourceIBMCmObjectCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_object", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createObjectOptions := &catalogmanagementv1.CreateObjectOptions{}

	createObjectOptions.SetCatalogIdentifier(d.Get("catalog_id").(string))
	if _, ok := d.GetOk("name"); ok {
		createObjectOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("parent_id"); ok {
		createObjectOptions.SetParentID(d.Get("parent_id").(string))
	}
	if _, ok := d.GetOk("label"); ok {
		createObjectOptions.SetLabel(d.Get("label").(string))
	}
	if _, ok := d.GetOk("tags"); ok {
		createObjectOptions.SetTags(SIToSS(d.Get("tags").([]interface{})))
	}
	if _, ok := d.GetOk("short_description"); ok {
		createObjectOptions.SetShortDescription(d.Get("short_description").(string))
	}
	if _, ok := d.GetOk("kind"); ok {
		createObjectOptions.SetKind(d.Get("kind").(string))
	}

	catalogObject, response, err := catalogManagementClient.CreateObjectWithContext(context, createObjectOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateObjectWithContext failed %s\n%s", err, response), "ibm_cm_object", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*catalogObject.ID)

	// update data here if provided in resource
	if _, ok := d.GetOk("data"); ok {
		replaceObjectOptions := &catalogmanagementv1.ReplaceObjectOptions{}
		replaceObjectOptions.SetCatalogIdentifier(*catalogObject.CatalogID)
		replaceObjectOptions.SetObjectIdentifier(*catalogObject.ID)
		replaceObjectOptions.SetID(*catalogObject.ID)
		replaceObjectOptions.SetCRN(*catalogObject.CRN)
		replaceObjectOptions.SetURL(*catalogObject.URL)
		replaceObjectOptions.SetCatalogName(*catalogObject.CatalogName)
		replaceObjectOptions.SetCatalogID(*catalogObject.CatalogID)
		replaceObjectOptions.SetKind(*catalogObject.Kind)
		replaceObjectOptions.SetRev(*catalogObject.Rev)
		replaceObjectOptions.SetCreated(*&catalogObject.Created)
		replaceObjectOptions.SetUpdated(*&catalogObject.Updated)
		replaceObjectOptions.SetPublish(*&catalogObject.Publish)
		replaceObjectOptions.SetState(*&catalogObject.State)
		if catalogObject.Label != nil {
			replaceObjectOptions.SetLabel(*catalogObject.Label)
		}
		if catalogObject.Name != nil {
			replaceObjectOptions.SetName(*catalogObject.Name)
		}
		if catalogObject.ParentID != nil {
			replaceObjectOptions.SetParentID(*catalogObject.ParentID)
		}
		if catalogObject.ShortDescription != nil {
			replaceObjectOptions.SetShortDescription(*catalogObject.ShortDescription)
		}
		if catalogObject.Tags != nil {
			replaceObjectOptions.SetTags(catalogObject.Tags)
		}

		dataMap := make(map[string]interface{})
		dataString, err := strconv.Unquote(d.Get("data").(string))
		if err != nil {
			dataString = d.Get("data").(string)
		}
		err = json.Unmarshal([]byte(dataString), &dataMap)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error unmarshalling json %s", err), "ibm_cm_object", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		replaceObjectOptions.SetData(dataMap)

		catalogObject, response, err = catalogManagementClient.ReplaceObjectWithContext(context, replaceObjectOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceObjectWithContext failed %s\n%s", err, response), "ibm_cm_object", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMCmObjectRead(context, d, meta)
}

func resourceIBMCmObjectRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getObjectOptions := &catalogmanagementv1.GetObjectOptions{}

	getObjectOptions.SetCatalogIdentifier(d.Get("catalog_id").(string))
	getObjectOptions.SetObjectIdentifier(d.Id())

	catalogObject, response, err := catalogManagementClient.GetObjectWithContext(context, getObjectOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetObjectWithContext failed %s\n%s", err, response), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("catalog_id", getObjectOptions.CatalogIdentifier); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_id: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("name", catalogObject.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("crn", catalogObject.CRN); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("url", catalogObject.URL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting url: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("parent_id", catalogObject.ParentID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting parent_id: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("label", catalogObject.Label); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting label: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if catalogObject.Tags != nil {
		modifiedTags := []string{}
		for _, tag := range catalogObject.Tags {
			if !strings.HasPrefix(tag, "svc:") && !strings.HasPrefix(tag, "fqdn:") && tag != *catalogObject.ParentID {
				modifiedTags = append(modifiedTags, tag)
			}
		}
		if err = d.Set("tags", modifiedTags); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting tags: %s", err), "ibm_cm_object", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("created", flex.DateTimeToString(catalogObject.Created)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("updated", flex.DateTimeToString(catalogObject.Updated)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("short_description", catalogObject.ShortDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting short_description: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("kind", catalogObject.Kind); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting kind: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	publishMap := map[string]interface{}{}
	if catalogObject.Publish != nil {
		publishMap, err = resourceIBMCmObjectPublishObjectToMap(catalogObject.Publish)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_object", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("publish", []map[string]interface{}{publishMap}); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting publish: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if catalogObject.State != nil {
		stateMap, err := resourceIBMCmObjectStateToMap(catalogObject.State)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_object", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("state", []map[string]interface{}{stateMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state: %s", err), "ibm_cm_object", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("catalog_id", catalogObject.CatalogID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_id: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("catalog_name", catalogObject.CatalogName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting catalog_name: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if catalogObject.Data != nil {
		dataString, err := json.Marshal(catalogObject.Data)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting data, error with json marshal: %s", err), "ibm_cm_object", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("data", string(dataString)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting data: %s", err), "ibm_cm_object", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("rev", catalogObject.Rev); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rev: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("object_id", catalogObject.ID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting object_id: %s", err), "ibm_cm_object", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIBMCmObjectUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_object", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getObjectOptions := &catalogmanagementv1.GetObjectOptions{}

	getObjectOptions.SetCatalogIdentifier(d.Get("catalog_id").(string))
	getObjectOptions.SetObjectIdentifier(d.Id())

	catalogObject, response, err := catalogManagementClient.GetObjectWithContext(context, getObjectOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetObjectWithContext failed %s\n%s", err, response), "ibm_cm_object", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	replaceObjectOptions := &catalogmanagementv1.ReplaceObjectOptions{}

	replaceObjectOptions.SetCatalogIdentifier(*catalogObject.CatalogID)
	replaceObjectOptions.SetObjectIdentifier(*catalogObject.ID)
	replaceObjectOptions.SetID(*catalogObject.ID)
	replaceObjectOptions.SetRev(*catalogObject.Rev)
	if catalogObject.State != nil {
		replaceObjectOptions.SetState(catalogObject.State)
	}
	if catalogObject.Publish != nil {
		replaceObjectOptions.SetPublish(catalogObject.Publish)
	}
	if _, ok := d.GetOk("name"); ok {
		replaceObjectOptions.SetName(d.Get("name").(string))
	} else if catalogObject.Name != nil {
		replaceObjectOptions.SetName(*catalogObject.Name)
	}
	if _, ok := d.GetOk("crn"); ok {
		replaceObjectOptions.SetCRN(d.Get("crn").(string))
	} else if catalogObject.CRN != nil {
		replaceObjectOptions.SetCRN(*catalogObject.CRN)
	}
	if _, ok := d.GetOk("url"); ok {
		replaceObjectOptions.SetURL(d.Get("url").(string))
	} else if catalogObject.URL != nil {
		replaceObjectOptions.SetURL(*catalogObject.URL)
	}
	if _, ok := d.GetOk("parent_id"); ok {
		replaceObjectOptions.SetParentID(d.Get("parent_id").(string))
	} else if catalogObject.ParentID != nil {
		replaceObjectOptions.SetParentID(*catalogObject.ParentID)
	}
	if _, ok := d.GetOk("label"); ok {
		replaceObjectOptions.SetLabel(d.Get("label").(string))
	} else if catalogObject.Label != nil {
		replaceObjectOptions.SetLabel(*catalogObject.Label)
	}
	if _, ok := d.GetOk("tags"); ok {
		replaceObjectOptions.SetTags(SIToSS(d.Get("tags").([]interface{})))
	} else if catalogObject.Tags != nil {
		replaceObjectOptions.SetTags(catalogObject.Tags)
	}
	if _, ok := d.GetOk("created"); ok {
		fmtDateTimeCreated, err := core.ParseDateTime(d.Get("created").(string))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_object", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		replaceObjectOptions.SetCreated(&fmtDateTimeCreated)
	} else if catalogObject.Created != nil {
		replaceObjectOptions.SetCreated(catalogObject.Created)
	}
	if _, ok := d.GetOk("updated"); ok {
		fmtDateTimeUpdated, err := core.ParseDateTime(d.Get("updated").(string))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_object", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		replaceObjectOptions.SetUpdated(&fmtDateTimeUpdated)
	} else if catalogObject.Updated != nil {
		replaceObjectOptions.SetUpdated(catalogObject.Updated)
	}
	if _, ok := d.GetOk("short_description"); ok {
		replaceObjectOptions.SetShortDescription(d.Get("short_description").(string))
	} else if catalogObject.ShortDescription != nil {
		replaceObjectOptions.SetShortDescription(*catalogObject.ShortDescription)
	}
	if _, ok := d.GetOk("kind"); ok {
		replaceObjectOptions.SetKind(d.Get("kind").(string))
	} else if catalogObject.Kind != nil {
		replaceObjectOptions.SetKind(*catalogObject.Kind)
	}
	if _, ok := d.GetOk("catalog_id"); ok {
		replaceObjectOptions.SetCatalogID(d.Get("catalog_id").(string))
	} else if catalogObject.CatalogID != nil {
		replaceObjectOptions.SetCatalogID(*catalogObject.CatalogID)
	}
	if _, ok := d.GetOk("catalog_name"); ok {
		replaceObjectOptions.SetCatalogName(d.Get("catalog_name").(string))
	} else if catalogObject.CatalogName != nil {
		replaceObjectOptions.SetCatalogName(*catalogObject.CatalogName)
	}
	if _, ok := d.GetOk("data"); ok {
		dataMap := make(map[string]interface{})
		dataString, err := strconv.Unquote(d.Get("data").(string))
		if err != nil {
			dataString = d.Get("data").(string)
		}
		err = json.Unmarshal([]byte(dataString), &dataMap)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error unmarshalling json %s", err), "ibm_cm_object", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		replaceObjectOptions.SetData(dataMap)
	} else if catalogObject.Data != nil {
		replaceObjectOptions.SetData(catalogObject.Data)
	}

	_, response, err = catalogManagementClient.ReplaceObjectWithContext(context, replaceObjectOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceObjectWithContext failed %s\n%s", err, response), "ibm_cm_object", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return resourceIBMCmObjectRead(context, d, meta)
}

func resourceIBMCmObjectDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_cm_object", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteObjectOptions := &catalogmanagementv1.DeleteObjectOptions{}

	deleteObjectOptions.SetCatalogIdentifier(d.Get("catalog_id").(string))
	deleteObjectOptions.SetObjectIdentifier(d.Id())

	response, err := catalogManagementClient.DeleteObjectWithContext(context, deleteObjectOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteObjectWithContext failed %s\n%s", err, response), "ibm_cm_object", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIBMCmObjectMapToPublishObject(modelMap map[string]interface{}) (*catalogmanagementv1.PublishObject, error) {
	model := &catalogmanagementv1.PublishObject{}
	if modelMap["permit_ibm_public_publish"] != nil {
		model.PermitIBMPublicPublish = core.BoolPtr(modelMap["permit_ibm_public_publish"].(bool))
	}
	if modelMap["ibm_approved"] != nil {
		model.IBMApproved = core.BoolPtr(modelMap["ibm_approved"].(bool))
	}
	if modelMap["public_approved"] != nil {
		model.PublicApproved = core.BoolPtr(modelMap["public_approved"].(bool))
	}
	return model, nil
}

func resourceIBMCmObjectMapToState(modelMap map[string]interface{}) (*catalogmanagementv1.State, error) {
	model := &catalogmanagementv1.State{}
	if modelMap["current"] != nil && modelMap["current"].(string) != "" {
		model.Current = core.StringPtr(modelMap["current"].(string))
	}
	if modelMap["current_entered"] != nil {

	}
	if modelMap["pending"] != nil && modelMap["pending"].(string) != "" {
		model.Pending = core.StringPtr(modelMap["pending"].(string))
	}
	if modelMap["pending_requested"] != nil {

	}
	if modelMap["previous"] != nil && modelMap["previous"].(string) != "" {
		model.Previous = core.StringPtr(modelMap["previous"].(string))
	}
	return model, nil
}

func resourceIBMCmObjectPublishObjectToMap(model *catalogmanagementv1.PublishObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PermitIBMPublicPublish != nil {
		modelMap["permit_ibm_public_publish"] = model.PermitIBMPublicPublish
	}
	if model.IBMApproved != nil {
		modelMap["ibm_approved"] = model.IBMApproved
	}
	if model.PublicApproved != nil {
		modelMap["public_approved"] = model.PublicApproved
	}
	return modelMap, nil
}

func resourceIBMCmObjectStateToMap(model *catalogmanagementv1.State) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Current != nil {
		modelMap["current"] = model.Current
	}
	if model.CurrentEntered != nil {
		modelMap["current_entered"] = model.CurrentEntered.String()
	}
	if model.Pending != nil {
		modelMap["pending"] = model.Pending
	}
	if model.PendingRequested != nil {
		modelMap["pending_requested"] = model.PendingRequested.String()
	}
	if model.Previous != nil {
		modelMap["previous"] = model.Previous
	}
	return modelMap, nil
}
