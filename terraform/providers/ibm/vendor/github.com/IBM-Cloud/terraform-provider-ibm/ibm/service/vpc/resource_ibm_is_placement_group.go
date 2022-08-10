// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const (
	isPlacementGroupDeleting          = "deleting"
	isPlacementGroupStable            = "stable"
	isPlacementGroupFailed            = "failed"
	isPlacementGroupDeleteDone        = "done"
	isPlacementGroupPending           = "pending"
	isPlacementGroupWaiting           = "waiting"
	isPlacementGroupUpdating          = "updating"
	isPlacementGroupResourcesAttached = "resources_attached"
	isPlacementGroupSuspended         = "suspended"

	isPlacementGroupTags       = "tags"
	isPlacementGroupAccessTags = "access_tags"
)

func ResourceIbmIsPlacementGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsPlacementGroupCreate,
		ReadContext:   resourceIbmIsPlacementGroupRead,
		UpdateContext: resourceIbmIsPlacementGroupUpdate,
		DeleteContext: resourceIbmIsPlacementGroupDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),
		Schema: map[string]*schema.Schema{
			"strategy": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_placement_group", "strategy"),
				Description:  "The strategy for this placement group- `host_spread`: place on different compute hosts- `power_spread`: place on compute hosts that use different power sourcesThe enumerated values for this property may expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the placement group on which the unexpected strategy was encountered.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_placement_group", "name"),
				Description:  "The unique user-defined name for this placement group. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The unique identifier of the resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.",
			},
			isPlacementGroupTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_placement_group", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags",
			},
			isPlacementGroupAccessTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_placement_group", isPlacementGroupAccessTags)},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the placement group was created.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this placement group.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this placement group.",
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the placement group.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func ResourceIbmIsPlacementGroupValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "strategy",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "host_spread, power_spread",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 isPlacementGroupAccessTags,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([ ]*[A-Za-z0-9:_.-]+[ ]*)+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_placement_group", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsPlacementGroupCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createPlacementGroupOptions := &vpcv1.CreatePlacementGroupOptions{}

	createPlacementGroupOptions.SetStrategy(d.Get("strategy").(string))
	createPlacementGroupOptions.SetName(d.Get("name").(string))

	if resourceGroupIntf, ok := d.GetOk("resource_group"); ok && resourceGroupIntf.(string) != "" {
		resourceGroup := resourceGroupIntf.(string)
		resourceGroupIdentity := &vpcv1.ResourceGroupIdentity{
			ID: &resourceGroup,
		}
		createPlacementGroupOptions.SetResourceGroup(resourceGroupIdentity)
	}

	placementGroup, response, err := vpcClient.CreatePlacementGroupWithContext(context, createPlacementGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] CreatePlacementGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(*placementGroup.ID)

	_, err = isWaitForPlacementGroupAvailable(vpcClient, d.Id(), d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for placement group to be available %s", err))
	}
	if _, ok := d.GetOk(isPlacementGroupTags); ok {
		oldList, newList := d.GetChange(isPlacementGroupTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *placementGroup.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error creating placement group (%s) tags: %s", d.Id(), err)
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk(isPlacementGroupAccessTags); ok {
		oldList, newList := d.GetChange(isPlacementGroupAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *placementGroup.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error creating placement group (%s) access tags: %s", d.Id(), err)
			return diag.FromErr(err)
		}
	}
	return resourceIbmIsPlacementGroupRead(context, d, meta)
}

func resourceIbmIsPlacementGroupRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getPlacementGroupOptions := &vpcv1.GetPlacementGroupOptions{}

	getPlacementGroupOptions.SetID(d.Id())

	placementGroup, response, err := vpcClient.GetPlacementGroupWithContext(context, getPlacementGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetPlacementGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	if err = d.Set("strategy", placementGroup.Strategy); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting strategy: %s", err))
	}
	if err = d.Set("name", placementGroup.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if placementGroup.ResourceGroup != nil {
		if err = d.Set("resource_group", *placementGroup.ResourceGroup.ID); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_group: %s", err))
		}
	}
	if err = d.Set("created_at", placementGroup.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("crn", placementGroup.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	if err = d.Set("href", placementGroup.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", placementGroup.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", placementGroup.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *placementGroup.CRN, "", isUserTagType)
	if err != nil {
		log.Printf(
			"Error getting placement group (%s) tags: %s", d.Id(), err)
	}

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *placementGroup.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error getting placement group (%s) access tags: %s", d.Id(), err)
	}

	d.Set(isPlacementGroupTags, tags)
	d.Set(isPlacementGroupAccessTags, accesstags)
	return nil
}

func resourceIbmIsPlacementGroupUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updatePlacementGroupOptions := &vpcv1.UpdatePlacementGroupOptions{}

	updatePlacementGroupOptions.SetID(d.Id())

	hasChange := false

	placementGroupPatchModel := &vpcv1.PlacementGroupPatch{}
	if d.HasChange("name") {
		plName := d.Get("name").(string)
		placementGroupPatchModel.Name = &plName
		hasChange = true
	}
	if hasChange {
		placementGroupPatch, err := placementGroupPatchModel.AsPatch()
		if err != nil {
			log.Printf("[DEBUG] Error calling AsPatch for PlacementGroupPatch %s", err)
			return diag.FromErr(err)
		}
		updatePlacementGroupOptions.SetPlacementGroupPatch(placementGroupPatch)
		_, response, err := vpcClient.UpdatePlacementGroupWithContext(context, updatePlacementGroupOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdatePlacementGroupWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
	}
	if d.HasChange(isPlacementGroupTags) {
		oldList, newList := d.GetChange(isPlacementGroupTags)
		err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, d.Get("crn").(string), "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on update of resource subnet (%s) tags: %s", d.Id(), err)
		}
	}

	if d.HasChange(isPlacementGroupAccessTags) {
		oldList, newList := d.GetChange(isPlacementGroupAccessTags)
		err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, d.Get("crn").(string), "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on update of resource subnet (%s) access tags: %s", d.Id(), err)
		}
	}
	return resourceIbmIsPlacementGroupRead(context, d, meta)
}

func resourceIbmIsPlacementGroupDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deletePlacementGroupOptions := &vpcv1.DeletePlacementGroupOptions{}

	deletePlacementGroupOptions.SetID(d.Id())

	response, err := vpcClient.DeletePlacementGroupWithContext(context, deletePlacementGroupOptions)
	if err != nil {
		if response.StatusCode == 409 {
			_, err = isWaitForPlacementGroupDeleteRetry(vpcClient, d, d.Id())
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error deleting PLacementGroup: %s", err))
			}
		} else {
			return diag.FromErr(fmt.Errorf("[ERROR] Error deleting PLacementGroup: %s\n%s", err, response))
		}
	}
	_, err = isWaitForPlacementGroupDelete(vpcClient, d, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func isWaitForPlacementGroupDelete(vpcClient *vpcv1.VpcV1, d *schema.ResourceData, id string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isPlacementGroupDeleting, isPlacementGroupStable, isPlacementGroupPending, isPlacementGroupWaiting, isPlacementGroupUpdating},
		Target:  []string{isPlacementGroupDeleteDone, ""},
		Refresh: func() (interface{}, string, error) {
			getPlacementGroupOptions := &vpcv1.GetPlacementGroupOptions{
				ID: &id,
			}

			placementGroup, response, err := vpcClient.GetPlacementGroup(getPlacementGroupOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return placementGroup, isPlacementGroupDeleteDone, nil
				} else if response != nil && response.StatusCode == 409 {
					return placementGroup, *placementGroup.LifecycleState, fmt.Errorf("[ERROR] The  PLacementGroup %s failed to delete: %v", id, err)
				}
				return nil, "", fmt.Errorf("[ERROR] Error Getting PLacementGroup: %s\n%s", err, response)
			}
			if *placementGroup.LifecycleState == isPlacementGroupFailed {
				return placementGroup, *placementGroup.LifecycleState, fmt.Errorf("[ERROR] The  PLacementGroup %s failed to delete: %v", id, err)
			}
			return placementGroup, isPlacementGroupDeleting, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isWaitForPlacementGroupDeleteRetry(vpcClient *vpcv1.VpcV1, d *schema.ResourceData, id string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isPlacementGroupResourcesAttached},
		Target:  []string{isPlacementGroupDeleting, isPlacementGroupDeleteDone, ""},
		Refresh: func() (interface{}, string, error) {
			deletePlacementGroupOptions := &vpcv1.DeletePlacementGroupOptions{}

			deletePlacementGroupOptions.SetID(id)

			response, err := vpcClient.DeletePlacementGroup(deletePlacementGroupOptions)
			if err != nil {
				if response != nil && response.StatusCode == 409 {
					return response, isPlacementGroupResourcesAttached, err
				} else if response != nil && response.StatusCode == 404 {
					return response, isPlacementGroupDeleteDone, nil
				}
				return response, "", fmt.Errorf("[ERROR] Error deleting PLacementGroup: %s\n%s", err, response)
			}
			return response, isPlacementGroupDeleting, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isWaitForPlacementGroupAvailable(vpcClient *vpcv1.VpcV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for placement group (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isPlacementGroupPending, isPlacementGroupWaiting, isPlacementGroupUpdating},
		Target:     []string{isPlacementGroupFailed, isPlacementGroupStable, isPlacementGroupSuspended, ""},
		Refresh:    isPlacementGroupRefreshFunc(vpcClient, id, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isPlacementGroupRefreshFunc(vpcClient *vpcv1.VpcV1, id string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getinsOptions := &vpcv1.GetPlacementGroupOptions{
			ID: &id,
		}
		placementGroup, response, err := vpcClient.GetPlacementGroup(getinsOptions)
		if placementGroup == nil || err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error getting placementGroup : %s\n%s", err, response)
		}

		d.Set("lifecycle_state", *placementGroup.LifecycleState)

		if *placementGroup.LifecycleState == isPlacementGroupSuspended || *placementGroup.LifecycleState == isPlacementGroupFailed {

			return placementGroup, *placementGroup.LifecycleState, fmt.Errorf("status of placement group is %s : \n%s", *placementGroup.LifecycleState, response)

		}
		return placementGroup, *placementGroup.LifecycleState, nil
	}
}
