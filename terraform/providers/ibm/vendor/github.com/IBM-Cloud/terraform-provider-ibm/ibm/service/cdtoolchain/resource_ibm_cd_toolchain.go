// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

package cdtoolchain

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/continuous-delivery-go-sdk/v2/cdtoolchainv2"
	"github.com/IBM/go-sdk-core/v5/core"
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
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain", "description"),
				Description:  "Describes the toolchain.",
			},
			"resource_group_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain", "resource_group_id"),
				Description:  "Resource group where the toolchain is located.",
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
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^(.*?)$`,
			MinValueLength:             0,
			MaxValueLength:             500,
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createToolchainOptions := &cdtoolchainv2.CreateToolchainOptions{}

	createToolchainOptions.SetName(d.Get("name").(string))
	createToolchainOptions.SetResourceGroupID(d.Get("resource_group_id").(string))
	if _, ok := d.GetOk("description"); ok {
		createToolchainOptions.SetDescription(d.Get("description").(string))
	}

	toolchainPost, _, err := cdToolchainClient.CreateToolchainWithContext(context, createToolchainOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateToolchainWithContext failed: %s", err.Error()), "ibm_cd_toolchain", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getToolchainByIDOptions := &cdtoolchainv2.GetToolchainByIDOptions{}

	getToolchainByIDOptions.SetToolchainID(d.Id())

	var toolchain *cdtoolchainv2.Toolchain
	var response *core.DetailedResponse
	err = resource.RetryContext(context, 10*time.Second, func() *resource.RetryError {
		toolchain, response, err = cdToolchainClient.GetToolchainByIDWithContext(context, getToolchainByIDOptions)
		if err != nil || toolchain == nil {
			if response != nil && response.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if conns.IsResourceTimeoutError(err) {
		toolchain, response, err = cdToolchainClient.GetToolchainByIDWithContext(context, getToolchainByIDOptions)
	}
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetToolchainByIDWithContext failed: %s", err.Error()), "ibm_cd_toolchain", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	tags, err := flex.GetTagsUsingCRN(meta, *toolchain.CRN)
	if err != nil {
		log.Printf(
			"Error on get of toolchain (%s) tags: %s", d.Id(), err)
	}
	d.Set("tags", tags)

	if err = d.Set("name", toolchain.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "read", "set-name").GetDiag()
	}
	if !core.IsNil(toolchain.Description) {
		if err = d.Set("description", toolchain.Description); err != nil {
			err = fmt.Errorf("Error setting description: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "read", "set-description").GetDiag()
		}
	}
	if err = d.Set("resource_group_id", toolchain.ResourceGroupID); err != nil {
		err = fmt.Errorf("Error setting resource_group_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "read", "set-resource_group_id").GetDiag()
	}
	if err = d.Set("account_id", toolchain.AccountID); err != nil {
		err = fmt.Errorf("Error setting account_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "read", "set-account_id").GetDiag()
	}
	if err = d.Set("location", toolchain.Location); err != nil {
		err = fmt.Errorf("Error setting location: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "read", "set-location").GetDiag()
	}
	if err = d.Set("crn", toolchain.CRN); err != nil {
		err = fmt.Errorf("Error setting crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "read", "set-crn").GetDiag()
	}
	if err = d.Set("href", toolchain.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "read", "set-href").GetDiag()
	}
	if err = d.Set("ui_href", toolchain.UIHref); err != nil {
		err = fmt.Errorf("Error setting ui_href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "read", "set-ui_href").GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(toolchain.CreatedAt)); err != nil {
		err = fmt.Errorf("Error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "read", "set-created_at").GetDiag()
	}
	if err = d.Set("updated_at", flex.DateTimeToString(toolchain.UpdatedAt)); err != nil {
		err = fmt.Errorf("Error setting updated_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "read", "set-updated_at").GetDiag()
	}
	if err = d.Set("created_by", toolchain.CreatedBy); err != nil {
		err = fmt.Errorf("Error setting created_by: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "read", "set-created_by").GetDiag()
	}

	return nil
}

func resourceIBMCdToolchainUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateToolchainOptions := &cdtoolchainv2.UpdateToolchainOptions{}

	updateToolchainOptions.SetToolchainID(d.Id())

	hasChange := false

	patchVals := &cdtoolchainv2.ToolchainPrototypePatch{}
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
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateToolchainOptions.ToolchainPrototypePatch = ResourceIBMCdToolchainToolchainPrototypePatchAsPatch(patchVals, d)

		_, _, err = cdToolchainClient.UpdateToolchainWithContext(context, updateToolchainOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateToolchainWithContext failed: %s", err.Error()), "ibm_cd_toolchain", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMCdToolchainRead(context, d, meta)
}

func resourceIBMCdToolchainDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteToolchainOptions := &cdtoolchainv2.DeleteToolchainOptions{}

	deleteToolchainOptions.SetToolchainID(d.Id())

	_, err = cdToolchainClient.DeleteToolchainWithContext(context, deleteToolchainOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteToolchainWithContext failed: %s", err.Error()), "ibm_cd_toolchain", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMCdToolchainToolchainPrototypePatchAsPatch(patchVals *cdtoolchainv2.ToolchainPrototypePatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}
	path = "description"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["description"] = nil
	}

	return patch
}
