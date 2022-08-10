// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cdtoolchain

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/continuous-delivery-go-sdk/cdtoolchainv2"
)

func ResourceIBMCdToolchain() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIBMCdToolchainCreate,
		ReadContext:   ResourceIBMCdToolchainRead,
		UpdateContext: ResourceIBMCdToolchainUpdate,
		DeleteContext: ResourceIBMCdToolchainDelete,
		Importer:      &schema.ResourceImporter{},

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
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tags associated with the toolchain.",
				Elem:        &schema.Schema{Type: schema.TypeString},
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
			Regexp:                     `^([^\\x00-\\x7F]|[a-zA-Z0-9-._ ])+$`,
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
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_toolchain", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIBMCdToolchainCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	postToolchainResponse, response, err := cdToolchainClient.CreateToolchainWithContext(context, createToolchainOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateToolchainWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateToolchainWithContext failed %s\n%s", err, response))
	}

	d.SetId(*postToolchainResponse.ID)

	return ResourceIBMCdToolchainRead(context, d, meta)
}

func ResourceIBMCdToolchainRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getToolchainByIDOptions := &cdtoolchainv2.GetToolchainByIDOptions{}

	getToolchainByIDOptions.SetToolchainID(d.Id())

	getToolchainByIDResponse, response, err := cdToolchainClient.GetToolchainByIDWithContext(context, getToolchainByIDOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetToolchainByIDWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetToolchainByIDWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("name", getToolchainByIDResponse.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("resource_group_id", getToolchainByIDResponse.ResourceGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
	}
	if err = d.Set("description", getToolchainByIDResponse.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	if err = d.Set("account_id", getToolchainByIDResponse.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}
	if err = d.Set("location", getToolchainByIDResponse.Location); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting location: %s", err))
	}
	if err = d.Set("crn", getToolchainByIDResponse.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("href", getToolchainByIDResponse.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(getToolchainByIDResponse.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(getToolchainByIDResponse.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("created_by", getToolchainByIDResponse.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("tags", getToolchainByIDResponse.Tags); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting tags: %s", err))
	}

	return nil
}

func ResourceIBMCdToolchainUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateToolchainOptions := &cdtoolchainv2.UpdateToolchainOptions{}

	updateToolchainOptions.SetToolchainID(d.Id())

	hasChange := false

	if d.HasChange("resource_group_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "resource_group_id"))
	}
	if d.HasChange("name") {
		updateToolchainOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("description") {
		updateToolchainOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}

	if hasChange {
		response, err := cdToolchainClient.UpdateToolchainWithContext(context, updateToolchainOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateToolchainWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateToolchainWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIBMCdToolchainRead(context, d, meta)
}

func ResourceIBMCdToolchainDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
