// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package globaltagging

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

const (
	resourceID   = "resource_id"
	tags         = "tags"
	resourceType = "resource_type"
	tagType      = "tag_type"
	acccountID   = "acccount_id"
	service      = "service"
)

func ResourceIBMResourceTag() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMResourceTagCreate,
		Read:     resourceIBMResourceTagRead,
		Update:   resourceIBMResourceTagUpdate,
		Delete:   resourceIBMResourceTagDelete,
		Importer: &schema.ResourceImporter{},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),

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
			acccountID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the account that owns the resources to be tagged (required if tag-type is set to service)",
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

func resourceIBMResourceTagCreate(d *schema.ResourceData, meta interface{}) error {
	var rType, tType string
	resources := []globaltaggingv1.Resource{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	accountID := userDetails.UserAccount

	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting global tagging client settings: %s", err)
	}

	resourceID := d.Get(resourceID).(string)
	if v, ok := d.GetOk(resourceType); ok && v != nil {
		rType = v.(string)
	}

	r := globaltaggingv1.Resource{ResourceID: flex.PtrToString(resourceID), ResourceType: flex.PtrToString(rType)}
	resources = append(resources, r)

	var add []string
	if v, ok := d.GetOk(tags); ok {
		tags := v.(*schema.Set)
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
		_, resp, err := gtClient.AttachTag(AttachTagOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error attaching resource tags : %v\n%s", resp, err)
		}
	}

	if strings.HasPrefix(resourceID, "crn:") {
		d.SetId(resourceID)
	} else {
		d.SetId(fmt.Sprintf("%s/%s", resourceID, rType))
	}

	return resourceIBMResourceTagRead(d, meta)
}

func resourceIBMResourceTagRead(d *schema.ResourceData, meta interface{}) error {
	var rID, rType, tType string

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	acctID := userDetails.UserAccount

	if strings.HasPrefix(d.Id(), "crn:") {
		rID = d.Id()
	} else {
		parts, err := flex.VmIdParts(d.Id())
		if err != nil {
			return err
		}
		if len(parts) < 2 {
			return fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of resourceID/resourceType", d.Id())
		}
		rID = parts[0]
		rType = parts[1]
	}

	if v, ok := d.GetOk(tagType); ok && v != nil {
		tType = v.(string)

		if tType == service {
			d.Set(acccountID, acctID)
		}
	}

	tagList, err := flex.GetGlobalTagsUsingSearchAPI(meta, rID, rType, tType)
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting resource tags for: %s with error : %s", rID, err)
	}

	d.Set(resourceID, rID)
	d.Set(resourceType, rType)
	d.Set(tags, tagList)

	return nil
}

func resourceIBMResourceTagUpdate(d *schema.ResourceData, meta interface{}) error {
	var rID, rType, tType string

	if strings.HasPrefix(d.Id(), "crn:") {
		rID = d.Id()
	} else {
		parts, err := flex.VmIdParts(d.Id())
		if err != nil {
			return err
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
			return fmt.Errorf("[ERROR] Error on create of resource tags: %s", err)
		}
	}

	return resourceIBMResourceTagRead(d, meta)
}

func resourceIBMResourceTagDelete(d *schema.ResourceData, meta interface{}) error {
	var rID, rType string

	if strings.HasPrefix(d.Id(), "crn:") {
		rID = d.Id()
	} else {
		parts, err := flex.VmIdParts(d.Id())
		if err != nil {
			return err
		}
		rID = parts[0]
		rType = parts[1]
	}

	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting global tagging client settings: %s", err)
	}

	var remove []string
	removeTags := d.Get(tags).(*schema.Set)
	remove = make([]string, len(removeTags.List()))
	for i, v := range removeTags.List() {
		remove[i] = fmt.Sprint(v)
	}

	if len(remove) > 0 {
		resources := []globaltaggingv1.Resource{}
		r := globaltaggingv1.Resource{ResourceID: flex.PtrToString(rID), ResourceType: flex.PtrToString(rType)}
		resources = append(resources, r)

		detachTagOptions := &globaltaggingv1.DetachTagOptions{
			Resources: resources,
			TagNames:  remove,
		}

		_, resp, err := gtClient.DetachTag(detachTagOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error detaching resource tags %v: %s\n%s", remove, err, resp)
		}
		for _, v := range remove {
			delTagOptions := &globaltaggingv1.DeleteTagOptions{
				TagName: flex.PtrToString(v),
			}
			_, resp, err := gtClient.DeleteTag(delTagOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error deleting resource tag %v: %s\n%s", v, err, resp)
			}
		}
	}
	return nil
}
