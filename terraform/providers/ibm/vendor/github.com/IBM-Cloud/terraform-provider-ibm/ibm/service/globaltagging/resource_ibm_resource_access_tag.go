// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package globaltagging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMResourceAccessTag() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMResourceAccessTagCreate,
		Read:     resourceIBMResourceAccessTagRead,
		Delete:   resourceIBMResourceAccessTagDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_resource_access_tag", "name"),
				Set:          flex.ResourceIBMVPCHash,
				Description:  "Name of the access tag",
			},
			tagType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the tag(access)",
			},
		},
		DeprecationMessage: "ibm_resource_access_tag has been deprecated. Use ibm_iam_access_tag instead.",
	}
}

func ResourceIBMResourceAccessTagValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-]):([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-])$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmResourceAccessTagValidator := validate.ResourceValidator{ResourceName: "ibm_resource_access_tag", Schema: validateSchema}
	return &ibmResourceAccessTagValidator
}

func resourceIBMResourceAccessTagCreate(d *schema.ResourceData, meta interface{}) error {

	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		return fmt.Errorf("Error getting global tagging client settings: %s", err)
	}

	tagName := d.Get("name").(string)
	add := make([]string, 0)
	add = append(add, tagName)
	accessTagType := "access"
	createTagOptions := &globaltaggingv1.CreateTagOptions{
		TagType:  &accessTagType,
		TagNames: add,
	}
	results, _, err := gtClient.CreateTag(createTagOptions)
	if err != nil {
		return err
	}
	if results != nil {
		errMap := make([]globaltaggingv1.CreateTagResultsResultsItem, 0)
		for _, res := range results.Results {
			if res.IsError != nil && *res.IsError {
				errMap = append(errMap, res)
			}
		}
		if len(errMap) > 0 {
			output, err := json.MarshalIndent(errMap, "", "    ")
			log.Printf("err is %s", err)
			return fmt.Errorf("[ERROR] Error while creating access tag(%s) : %s", tagName, string(output))
		}
	}

	d.SetId(tagName)
	d.Set(tagType, accessTagType)

	return nil
}

func resourceIBMResourceAccessTagRead(d *schema.ResourceData, meta interface{}) error {
	tagName := d.Id()
	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		return fmt.Errorf("Error getting global tagging client settings: %s", err)
	}
	accessTagType := "access"
	listTagsOptions := &globaltaggingv1.ListTagsOptions{
		TagType: &accessTagType,
	}
	taggingResult, _, err := gtClient.ListTags(listTagsOptions)
	if err != nil {
		return err
	}

	var taglist []string
	for _, item := range taggingResult.Items {
		taglist = append(taglist, *item.Name)
	}
	existingAccessTags := flex.NewStringSet(flex.ResourceIBMVPCHash, taglist)
	if !existingAccessTags.Contains(tagName) {
		d.SetId("")
		return nil
	}
	d.Set("name", tagName)
	d.Set(tagType, accessTagType)
	return nil
}

func resourceIBMResourceAccessTagDelete(d *schema.ResourceData, meta interface{}) error {

	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting global tagging client settings: %s", err)
	}
	tagName := d.Get("name").(string)
	accessTagType := "access"

	deleteTagOptions := &globaltaggingv1.DeleteTagOptions{
		TagName: &tagName,
		TagType: &accessTagType,
	}

	results, resp, err := gtClient.DeleteTag(deleteTagOptions)

	if err != nil {
		return fmt.Errorf("[ERROR] Error while deleting access tag(%s) : %v\n%v", tagName, err, resp)
	}
	if results != nil {
		errMap := make([]globaltaggingv1.DeleteTagResultsItem, 0)
		for _, res := range results.Results {
			if res.IsError != nil && *res.IsError {
				errMap = append(errMap, res)
			}
		}
		if len(errMap) > 0 {
			output, err := json.MarshalIndent(errMap, "", "    ")
			log.Printf("err is %s", err)
			return fmt.Errorf("[ERROR] Error while deleting access tag(%s) : %s", tagName, string(output))
		}
	}

	d.SetId("")
	return nil
}
