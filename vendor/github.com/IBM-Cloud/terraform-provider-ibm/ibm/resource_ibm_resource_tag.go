// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	resourceID   = "resource_id"
	tags         = "tags"
	resourceType = "resource_type"
	tagType      = "tag_type"
	acccountID   = "acccount_id"
	service      = "service"
	crnRegex     = "^crn:.+:.+:.+:.+:.+:$"
)

func resourceIBMResourceTag() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMResourceTagCreate,
		Read:     resourceIBMResourceTagRead,
		Update:   resourceIBMResourceTagUpdate,
		Delete:   resourceIBMResourceTagDelete,
		Importer: &schema.ResourceImporter{},

		CustomizeDiff: customdiff.Sequence(
			func(diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			resourceID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_resource_tag", resourceID),
				Description:  "CRN of the resource on which the tags should be attached",
			},
			tags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_resource_tag", tags)},
				Set:         resourceIBMVPCHash,
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
				ValidateFunc: validateAllowedStringValue([]string{"service", "access", "user"}),
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

func resourceIBMResourceTagValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 resourceID,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^crn:v1(:[a-zA-Z0-9 \-\._~\*\+,;=!$&'\(\)\/\?#\[\]@]*){8}$|^[0-9]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 tags,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmResourceTagValidator := ResourceValidator{ResourceName: "ibm_resource_tag", Schema: validateSchema}
	return &ibmResourceTagValidator
}

func resourceIBMResourceTagCreate(d *schema.ResourceData, meta interface{}) error {
	var rType, tType string
	resources := []globaltaggingv1.Resource{}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	accountID := userDetails.userAccount

	gtClient, err := meta.(ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		return fmt.Errorf("Error getting global tagging client settings: %s", err)
	}

	resourceID := d.Get(resourceID).(string)
	if v, ok := d.GetOk(resourceType); ok && v != nil {
		rType = v.(string)
	}

	r := globaltaggingv1.Resource{ResourceID: ptrToString(resourceID), ResourceType: ptrToString(rType)}
	resources = append(resources, r)

	var add []string
	if v, ok := d.GetOk(tags); ok {
		tags := v.(*schema.Set)
		for _, t := range tags.List() {
			add = append(add, fmt.Sprint(t))
		}
	}

	schematicTags := os.Getenv("IC_ENV_TAGS")
	var envTags []string
	if schematicTags != "" {
		envTags = strings.Split(schematicTags, ",")
		add = append(add, envTags...)
	}

	AttachTagOptions := &globaltaggingv1.AttachTagOptions{}
	AttachTagOptions.Resources = resources
	AttachTagOptions.TagNames = add
	if v, ok := d.GetOk(tagType); ok && v != nil {
		tType = v.(string)
		AttachTagOptions.TagType = ptrToString(tType)

		if tType == service {
			AttachTagOptions.AccountID = ptrToString(accountID)
		}
	}

	if len(add) > 0 {
		_, resp, err := gtClient.AttachTag(AttachTagOptions)
		if err != nil {
			return fmt.Errorf("Error attaching resource tags >>>>  %v : %s", resp, err)
		}
	}

	crn, err := regexp.Compile(crnRegex)
	if err != nil {
		return err
	}

	if crn.MatchString(resourceID) {
		d.SetId(resourceID)
	} else {
		d.SetId(fmt.Sprintf("%s/%s", resourceID, resourceType))
	}

	return resourceIBMResourceTagRead(d, meta)
}

func resourceIBMResourceTagRead(d *schema.ResourceData, meta interface{}) error {
	var rID, rType, tType string

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	acctID := userDetails.userAccount

	crn, err := regexp.Compile(crnRegex)
	if err != nil {
		return err
	}

	if crn.MatchString(d.Id()) {
		rID = d.Id()
	} else {
		parts, err := vmIdParts(d.Id())
		if err != nil {
			return err
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

	tagList, err := GetGlobalTagsUsingCRN(meta, rID, resourceType, tType)
	if err != nil {
		if apierr, ok := err.(bmxerror.RequestFailure); ok && apierr.StatusCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting resource tags for: %s with error : %s\n", rID, err)
	}

	d.Set(resourceID, rID)
	d.Set(resourceType, rType)
	d.Set(tags, tagList)

	return nil
}

func resourceIBMResourceTagUpdate(d *schema.ResourceData, meta interface{}) error {
	var rID, rType, tType string

	crn, err := regexp.Compile(crnRegex)
	if err != nil {
		return err
	}

	if crn.MatchString(d.Id()) {
		rID = d.Id()
	} else {
		parts, err := vmIdParts(d.Id())
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
		err := UpdateGlobalTagsUsingCRN(oldList, newList, meta, rID, rType, tType)
		if err != nil {
			return fmt.Errorf(
				"Error on create of resource tags: %s", err)
		}
	}

	return resourceIBMResourceTagRead(d, meta)
}

func resourceIBMResourceTagDelete(d *schema.ResourceData, meta interface{}) error {
	var rID, rType string

	crn, err := regexp.Compile(crnRegex)
	if err != nil {
		return err
	}

	if crn.MatchString(d.Id()) {
		rID = d.Id()
	} else {
		parts, err := vmIdParts(d.Id())
		if err != nil {
			return err
		}
		rID = parts[0]
		rType = parts[1]
	}

	gtClient, err := meta.(ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		return fmt.Errorf("Error getting global tagging client settings: %s", err)
	}

	var remove []string
	removeTags := d.Get(tags).(*schema.Set)
	remove = make([]string, len(removeTags.List()))
	for i, v := range removeTags.List() {
		remove[i] = fmt.Sprint(v)
	}

	if len(remove) > 0 {
		resources := []globaltaggingv1.Resource{}
		r := globaltaggingv1.Resource{ResourceID: ptrToString(rID), ResourceType: ptrToString(rType)}
		resources = append(resources, r)

		detachTagOptions := &globaltaggingv1.DetachTagOptions{
			Resources: resources,
			TagNames:  remove,
		}

		_, resp, err := gtClient.DetachTag(detachTagOptions)
		if err != nil {
			return fmt.Errorf("Error detaching resource tags %v: %s\n%s", remove, err, resp)
		}
		for _, v := range remove {
			delTagOptions := &globaltaggingv1.DeleteTagOptions{
				TagName: ptrToString(v),
			}
			_, resp, err := gtClient.DeleteTag(delTagOptions)
			if err != nil {
				return fmt.Errorf("Error deleting resource tag %v: %s\n%s", v, err, resp)
			}
		}
	}
	return nil
}
