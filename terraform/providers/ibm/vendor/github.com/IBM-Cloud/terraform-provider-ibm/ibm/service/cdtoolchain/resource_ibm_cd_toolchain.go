// Copyright IBM Corp. 2022, 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cdtoolchain

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/continuous-delivery-go-sdk/cdtoolchainv2"
)

func ResourceIBMCdToolchain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCdToolchainCreate,
		ReadContext:   resourceIBMCdToolchainRead,
		UpdateContext: resourceIBMCdToolchainUpdate,
		DeleteContext: resourceIBMCdToolchainDelete,
		Importer:      &schema.ResourceImporter{},
		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain", "name"),
				Description:  "Toolchain name.",
			},
			"resource_group_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain", "resource_group_id"),
				Description:  "Resource group where toolchain will be created.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain", "description"),
				Description:  "Describes the toolchain.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account ID where toolchain can be found.",
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Toolchain region.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Toolchain CRN.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI that can be used to retrieve toolchain.",
			},
			"ui_href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of a user-facing user interface for this toolchain.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Toolchain creation timestamp.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latest toolchain update timestamp.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity that created the toolchain.",
			},
			"tags": &schema.Schema{
				Type:         schema.TypeSet,
				Optional:     true,
				Computed:     true,
				Elem:         &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain", "tags")},
				Set:          flex.ResourceIBMVPCHash,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain", "tags"),
				Description:  "Toolchain tags.",
			},
		},
	}
}

func ResourceIBMCdToolchainValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([^\x00-\x7F]|[a-zA-Z0-9-._ ])+$`,
			MinValueLength:             0,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "resource_group_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-f]{32}$`,
			MinValueLength:             32,
			MaxValueLength:             32,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^(.*?)$`,
			MinValueLength:             0,
			MaxValueLength:             500,
		},
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "tags",
			Optional:                   true,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_toolchain", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCdToolchainCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createToolchainOptions := &cdtoolchainv2.CreateToolchainOptions{}

	createToolchainOptions.SetName(d.Get("name").(string))
	createToolchainOptions.SetResourceGroupID(d.Get("resource_group_id").(string))
	if _, ok := d.GetOk("description"); ok {
		createToolchainOptions.SetDescription(d.Get("description").(string))
	}

	toolchainPost, response, err := cdToolchainClient.CreateToolchainWithContext(context, createToolchainOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateToolchainWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateToolchainWithContext failed %s\n%s", err, response))
	}

	d.SetId(*toolchainPost.ID)

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk("tags"); ok || v != "" {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *toolchainPost.CRN)
		if err != nil {
			log.Printf(
				"Error on create of toolchain (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceIBMCdToolchainRead(context, d, meta)
}

func resourceIBMCdToolchainRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getToolchainByIDOptions := &cdtoolchainv2.GetToolchainByIDOptions{}

	getToolchainByIDOptions.SetToolchainID(d.Id())

	toolchain, response, err := cdToolchainClient.GetToolchainByIDWithContext(context, getToolchainByIDOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetToolchainByIDWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetToolchainByIDWithContext failed %s\n%s", err, response))
	}

	tags, err := flex.GetTagsUsingCRN(meta, *toolchain.CRN)
	if err != nil {
		log.Printf(
			"Error on get of toolchain (%s) tags: %s", d.Id(), err)
	}
	d.Set("tags", tags)

	if err = d.Set("name", toolchain.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("resource_group_id", toolchain.ResourceGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
	}
	if err = d.Set("description", toolchain.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	if err = d.Set("account_id", toolchain.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}
	if err = d.Set("location", toolchain.Location); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting location: %s", err))
	}
	if err = d.Set("crn", toolchain.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("href", toolchain.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("ui_href", toolchain.UIHref); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ui_href: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(toolchain.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(toolchain.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("created_by", toolchain.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}

	return nil
}

func resourceIBMCdToolchainUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateToolchainOptions := &cdtoolchainv2.UpdateToolchainOptions{}

	updateToolchainOptions.SetToolchainID(d.Id())

	hasChange := false

	patchVals := &cdtoolchainv2.ToolchainPrototypePatch{}
	if d.HasChange("resource_group_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "resource_group_id"))
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	if d.HasChange("description") {
		newDescription := d.Get("description").(string)
		patchVals.Description = &newDescription
		hasChange = true
	}

	if d.HasChange("tags") {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, d.Get("crn").(string))
		if err != nil {
			log.Printf(
				"Error on update of toolchain (%s) tags: %s", d.Id(), err)
		}
	}

	if hasChange {
		updateToolchainOptions.ToolchainPrototypePatch, _ = patchVals.AsPatch()
		_, response, err := cdToolchainClient.UpdateToolchainWithContext(context, updateToolchainOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateToolchainWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateToolchainWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMCdToolchainRead(context, d, meta)
}

func resourceIBMCdToolchainDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteToolchainOptions := &cdtoolchainv2.DeleteToolchainOptions{}

	deleteToolchainOptions.SetToolchainID(d.Id())

	response, err := cdToolchainClient.DeleteToolchainWithContext(context, deleteToolchainOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteToolchainWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteToolchainWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
