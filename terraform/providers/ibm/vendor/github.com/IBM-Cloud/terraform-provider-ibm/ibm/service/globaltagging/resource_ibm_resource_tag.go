// Copyright IBM Corp. 2017, 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package globaltagging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

const (
	resourceID   = "resource_id"
	tags         = "tags"
	resourceType = "resource_type"
	tagType      = "tag_type"
	accountID    = "account_id"
	service      = "service"
	replace      = "replace"
)

func ResourceIBMResourceTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMResourceTagCreate,
		ReadContext:   resourceIBMResourceTagRead,
		UpdateContext: resourceIBMResourceTagUpdate,
		DeleteContext: resourceIBMResourceTagDelete,
		Importer:      &schema.ResourceImporter{},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			resourceID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_resource_tag", resourceID),
				Description:  "CRN of the resource on which the tags should be attached",
			},
			tags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_resource_tag", tags)},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags associated with resource instance",
			},
			resourceType: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource type on which the tags should be attached",
			},
			tagType: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_resource_tag", tagType),
				Description:  "Type of the tag. Only allowed values are: user, or service or access (default value : user)",
			},
			accountID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the account that owns the resources to be tagged (required if tag-type is set to service)",
			},
			replace: {
				Type:             schema.TypeBool,
				DiffSuppressFunc: flex.ApplyOnce,
				Optional:         true,
				Default:          false,
				Description:      "If true, it indicates that the attaching operation is a replacement operation",
			},
		},
	}
}

func ResourceIBMResourceTagValidator() *validate.ResourceValidator {
	tagTypeAllowedValues := "service,access,user"
	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 resourceID,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^crn:v1(:[a-zA-Z0-9 \-\._~\*\+,;=!$&'\(\)\/\?#\[\]@]*){8}$|^[0-9]+$`,
			MinValueLength:             1,
			MaxValueLength:             1024})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 tags,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tag_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              tagTypeAllowedValues})

	ibmResourceTagValidator := validate.ResourceValidator{ResourceName: "ibm_resource_tag", Schema: validateSchema}
	return &ibmResourceTagValidator
}

func resourceIBMResourceTagCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var rType, tType string
	resources := []globaltaggingv1.Resource{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_resource_tag", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	accountID := userDetails.UserAccount

	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_resource_tag", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	resourceID := d.Get(resourceID).(string)
	if v, ok := d.GetOk(resourceType); ok && v != nil {
		rType = v.(string)
	}

	r := globaltaggingv1.Resource{ResourceID: flex.PtrToString(resourceID), ResourceType: flex.PtrToString(rType)}
	resources = append(resources, r)

	var add []string
	var news *schema.Set
	if v, ok := d.GetOk(tags); ok {
		tags := v.(*schema.Set)
		news = v.(*schema.Set)
		for _, t := range tags.List() {
			add = append(add, fmt.Sprint(t))
		}
	}

	AttachTagOptions := &globaltaggingv1.AttachTagOptions{}
	AttachTagOptions.Resources = resources
	AttachTagOptions.TagNames = add
	if v, ok := d.GetOk(tagType); ok && v != nil {
		tType = v.(string)
		AttachTagOptions.TagType = flex.PtrToString(tType)

		if tType == service {
			AttachTagOptions.AccountID = flex.PtrToString(accountID)
		}
	}

	if v, ok := d.GetOk(replace); ok && v != nil {
		replace := v.(bool)
		AttachTagOptions.Replace = &replace

	}

	// Fetch tags from schematics only if they are user tags
	if strings.TrimSpace(tagType) == "" || tagType == "user" {
		schematicTags := os.Getenv("IC_ENV_TAGS")
		var envTags []string
		if schematicTags != "" {
			envTags = strings.Split(schematicTags, ",")
			add = append(add, envTags...)
		}
	}

	if len(add) > 0 {
		results, fullResponse, err := gtClient.AttachTagWithContext(context, AttachTagOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_resource_tag", "create")
			return tfErr.GetDiag()
		}

		// Check if there are errors on the attach internal response
		if results != nil {
			errMap := make([]globaltaggingv1.TagResultsItem, 0)
			for _, res := range results.Results {
				if res.IsError != nil && *res.IsError {
					errMap = append(errMap, res)
				}
			}
			if len(errMap) > 0 {
				output, _ := json.MarshalIndent(errMap, "", "    ")
				return diag.FromErr(fmt.Errorf("Error while creating tag: %s - Full response: %s", string(output), fullResponse))
			}
		}
		response, errored := flex.WaitForTagsAvailable(meta, resourceID, resourceType, tagType, news, d.Timeout(schema.TimeoutCreate))
		if errored != nil {
			log.Printf(`[ERROR] Error waiting for resource tags %s : %v
%v`, resourceID, errored, response)
		}
	}

	if strings.HasPrefix(resourceID, "crn:") {
		d.SetId(resourceID)
	} else {
		d.SetId(fmt.Sprintf("%s/%s", resourceID, rType))
	}

	return resourceIBMResourceTagRead(context, d, meta)
}

func resourceIBMResourceTagRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var rID, rType, tType string

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_resource_tag", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	acctID := userDetails.UserAccount

	if strings.HasPrefix(d.Id(), "crn:") {
		rID = d.Id()
	} else {
		parts, err := flex.VmIdParts(d.Id())
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_resource_tag", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if len(parts) < 2 {
			return diag.FromErr(fmt.Errorf("Incorrect ID %s: Id should be a combination of resourceID/resourceType", d.Id()))
		}
		rID = parts[0]
		rType = parts[1]
	}

	if v, ok := d.GetOk(tagType); ok && v != nil {
		tType = v.(string)

		if tType == service {
			d.Set(accountID, acctID)
		}
	}

	tagList, err := flex.GetGlobalTagsUsingSearchAPI(meta, rID, rType, tType)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error getting resource tags for: %s with error : %s", rID, err))
	}

	d.Set(resourceID, rID)
	d.Set(resourceType, rType)
	d.Set(tags, tagList)

	return nil
}

func resourceIBMResourceTagUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var rID, rType, tType string

	if strings.HasPrefix(d.Id(), "crn:") {
		rID = d.Id()
	} else {
		parts, err := flex.VmIdParts(d.Id())
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_resource_tag", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		rID = parts[0]
		rType = parts[1]
	}

	if v, ok := d.GetOk(tagType); ok && v != nil {
		tType = v.(string)
	}

	if _, ok := d.GetOk(tags); ok {
		oldList, newList := d.GetChange(tags)
		err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, rID, rType, tType)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error on create of resource tags: %s", err))
		}
	}

	return resourceIBMResourceTagRead(context, d, meta)
}

func resourceIBMResourceTagDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var rID, rType string
	if strings.HasPrefix(d.Id(), "crn:") {
		rID = d.Id()
	} else {
		parts, err := flex.VmIdParts(d.Id())
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_resource_tag", "delete")
			log.Printf("[ERROR] Error in deleting.\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		rID = parts[0]
		rType = parts[1]
	}
	removeTags := d.Get(tags).(*schema.Set)
	var tType string
	if v, ok := d.GetOk(tagType); ok && v != nil {
		tType = v.(string)
	} else {
		tType = "user"
	}
	if len(removeTags.List()) > 0 {
		err := flex.UpdateGlobalTagsUsingCRN(removeTags, nil, meta, rID, rType, tType)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error on deleting tags: %s", err))
		}
	}
	return nil
}
