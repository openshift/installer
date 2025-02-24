// Copyright IBM Corp. 2017, 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package globaltagging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIamAccessTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIamAccessTagCreate,
		ReadContext:   resourceIBMIamAccessTagRead,
		DeleteContext: resourceIBMIamAccessTagDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_access_tag", "name"),
				Set:          flex.ResourceIBMVPCHash,
				Description:  "Name of the access tag",
			},
			tagType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the tag(access)",
			},
		},
	}
}

func ResourceIBMIamAccessTagValidator() *validate.ResourceValidator {

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

	ibmIamAccessTagValidator := validate.ResourceValidator{ResourceName: "ibm_iam_access_tag", Schema: validateSchema}
	return &ibmIamAccessTagValidator
}

func resourceIBMIamAccessTagCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_iam_access_tag", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	tagName := d.Get("name").(string)
	add := make([]string, 0)
	add = append(add, tagName)
	accessTagType := "access"
	createTagOptions := &globaltaggingv1.CreateTagOptions{
		TagType:  &accessTagType,
		TagNames: add,
	}
	results, _, err := gtClient.CreateTagWithContext(context, createTagOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_iam_access_tag", "create")
		return tfErr.GetDiag()
	}
	if results != nil {
		errMap := make([]globaltaggingv1.CreateTagResultsResultsItem, 0)
		for _, res := range results.Results {
			if res.IsError != nil && *res.IsError {
				errMap = append(errMap, res)
			}
		}
		if len(errMap) > 0 {
			output, _ := json.MarshalIndent(errMap, "", "    ")
			return diag.FromErr(fmt.Errorf("Error while creating access tag(%s) : %s", tagName, string(output)))
		}
	}

	d.SetId(tagName)
	d.Set(tagType, accessTagType)

	return nil
}

func resourceIBMIamAccessTagRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tagName := d.Id()
	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_iam_access_tag", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	accessTagType := "access"
	listTagsOptions := &globaltaggingv1.ListTagsOptions{
		TagType: &accessTagType,
	}
	taggingResult, _, err := gtClient.ListTagsWithContext(context, listTagsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_iam_access_tag", "read")
		return tfErr.GetDiag()
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

func resourceIBMIamAccessTagDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_iam_access_tag", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	tagName := d.Get("name").(string)
	accessTagType := "access"

	deleteTagOptions := &globaltaggingv1.DeleteTagOptions{
		TagName: &tagName,
		TagType: &accessTagType,
	}

	results, resp, err := gtClient.DeleteTagWithContext(context, deleteTagOptions)

	if err != nil {
		return diag.FromErr(fmt.Errorf("Error while deleting access tag calling api (%s) : %v\n%v", tagName, err, resp))
	}
	if results != nil {
		errMap := make([]globaltaggingv1.DeleteTagResultsItem, 0)
		for _, res := range results.Results {
			if res.IsError != nil && *res.IsError {
				errMap = append(errMap, res)
			}
		}
		if len(errMap) > 0 {
			output, _ := json.MarshalIndent(errMap, "", "    ")
			return diag.FromErr(fmt.Errorf("Error while deleting access tag in results (%s) : %s", tagName, string(output)))
		}
	}

	d.SetId("")
	return nil
}
