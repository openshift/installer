// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsSnapshotConsistencyGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsSnapshotConsistencyGroupCreate,
		ReadContext:   resourceIBMIsSnapshotConsistencyGroupRead,
		UpdateContext: resourceIBMIsSnapshotConsistencyGroupUpdate,
		DeleteContext: resourceIBMIsSnapshotConsistencyGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"delete_snapshots_on_delete": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether deleting the snapshot consistency group will also delete the snapshots in the group.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_snapshot_consistency_group", "name"),
				Description:  "The name for this snapshot consistency group. The name is unique across all snapshot consistency groups in the region.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Resource group Id",
			},
			"snapshots": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The member snapshots that are data-consistent with respect to captured time. (may be[deleted](https://cloud.ibm.com/apidocs/vpc#deleted-resources)).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name for this snapshot. The name is unique across all snapshots in the region.",
						},
						"source_volume": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The volume to create this snapshot from.",
						},
						"tags": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         flex.ResourceIBMVPCHash,
							Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_snapshot_consistency_group", "tags")},
							Description: "User Tags for the snapshot",
						},
					},
				},
			},
			"snapshot_reference": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The member snapshots that are data-consistent with respect to captured time. (may be[deleted](https://cloud.ibm.com/apidocs/vpc#deleted-resources)).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN of this snapshot.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this snapshot.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this snapshot.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this snapshot. The name is unique across all snapshots in the region.",
						},
						"remote": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource is remote to this region,and identifies the native region.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this region.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this region.",
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"backup_policy_plan": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "If present, the backup policy plan which created this snapshot consistency group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this backup policy plan.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this backup policy plan.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this backup policy plan. The name is unique across all plans in the backup policy.",
						},
						"remote": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource is remote to this region,and identifies the native region.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this region.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this region.",
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this snapshot consistency group was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN of this snapshot consistency group.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this snapshot consistency group.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of this snapshot consistency group.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"service_tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The [service tags](https://cloud.ibm.com/apidocs/tagging#types-of-tags)[`is.instance:` prefix](https://cloud.ibm.com/docs/vpc?topic=vpc-snapshots-vpc-faqs) associated with this snapshot consistency group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_snapshot_consistency_group", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "Snapshot Consistency Group tags list",
			},
			"access_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_snapshot_consistency_group", "access_tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},
		},
	}
}

func ResourceIBMIsSnapshotConsistencyGroupValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
	)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "access_tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-]):([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-])$`,
			MinValueLength:             1,
			MaxValueLength:             128})
	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_snapshot_consistency_group", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsSnapshotConsistencyGroupCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	createSnapshotConsistencyGroupOptions := &vpcv1.CreateSnapshotConsistencyGroupOptions{}
	snapshotConsistencyGroupPrototype := &vpcv1.SnapshotConsistencyGroupPrototype{}

	if _, ok := d.GetOk("delete_snapshots_on_delete"); ok {
		snapshotConsistencyGroupPrototype.DeleteSnapshotsOnDelete = core.BoolPtr(d.Get("delete_snapshots_on_delete").(bool))
	}

	if _, ok := d.GetOk("name"); ok {
		snapshotConsistencyGroupPrototype.Name = core.StringPtr(d.Get("name").(string))
	}

	if rgrp, ok := d.GetOk("resource_group"); ok {
		rg := rgrp.(string)
		snapshotConsistencyGroupPrototype.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}
	var snapshotConsistencyGroupPrototypeSnapshotsItemArray []vpcv1.SnapshotPrototypeSnapshotConsistencyGroupContext
	snapshotsArray := d.Get("snapshots").([]interface{})
	for _, snapshot := range snapshotsArray {
		snapshotVal := snapshot.(map[string]interface{})
		snapshotConsistencyGroupPrototypeSnapshotsItem := &vpcv1.SnapshotPrototypeSnapshotConsistencyGroupContext{}

		volume := snapshotVal["source_volume"].(string)
		snapshotConsistencyGroupPrototypeSnapshotsItem.SourceVolume = &vpcv1.VolumeIdentity{
			ID: &volume,
		}
		name := snapshotVal["name"].(string)
		if name != "" {
			snapshotConsistencyGroupPrototypeSnapshotsItem.Name = &name
		}

		userTags := snapshotVal["tags"].(*schema.Set)
		if userTags != nil && userTags.Len() != 0 {
			userTagsArray := make([]string, userTags.Len())
			for i, userTag := range userTags.List() {
				userTagStr := userTag.(string)
				userTagsArray[i] = userTagStr
			}
			schematicTags := os.Getenv("IC_ENV_TAGS")
			var envTags []string
			if schematicTags != "" {
				envTags = strings.Split(schematicTags, ",")
				userTagsArray = append(userTagsArray, envTags...)
			}
			snapshotConsistencyGroupPrototypeSnapshotsItem.UserTags = userTagsArray
		}
		snapshotConsistencyGroupPrototypeSnapshotsItemArray = append(snapshotConsistencyGroupPrototypeSnapshotsItemArray, *snapshotConsistencyGroupPrototypeSnapshotsItem)
	}
	snapshotConsistencyGroupPrototype.Snapshots = snapshotConsistencyGroupPrototypeSnapshotsItemArray

	createSnapshotConsistencyGroupOptions.SnapshotConsistencyGroupPrototype = snapshotConsistencyGroupPrototype
	snapshotConsistencyGroup, response, err := vpcClient.CreateSnapshotConsistencyGroupWithContext(context, createSnapshotConsistencyGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateSnapshotConsistencyGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateSnapshotConsistencyGroupWithContext failed %s\n%s", err, response))
	}

	d.SetId(*snapshotConsistencyGroup.ID)

	_, err = isWaitForSnapshotConsistencyGroupAvailable(vpcClient, d.Id(), d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return diag.FromErr(err)
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk("tags"); ok || v != "" {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *snapshotConsistencyGroup.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc snapshot consistency group (%s) tags: %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("access_tags"); ok {
		oldList, newList := d.GetChange("access_tags")
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *snapshotConsistencyGroup.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource snapshot consistency group (%s) access tags: %s", d.Id(), err)
		}
	}

	return resourceIBMIsSnapshotConsistencyGroupRead(context, d, meta)
}

func isWaitForSnapshotConsistencyGroupAvailable(sess *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Snapshot Consistency Group(%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"stable", "failed"},
		Refresh:    isSnapshotConsistencyGroupRefreshFunc(sess, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isSnapshotConsistencyGroupRefreshFunc(vpcClient *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getSnapshotConsistencyGroupOptions := &vpcv1.GetSnapshotConsistencyGroupOptions{}
		getSnapshotConsistencyGroupOptions.SetID(id)

		snapshotConsistencyGroup, response, err := vpcClient.GetSnapshotConsistencyGroup(getSnapshotConsistencyGroupOptions)
		if err != nil {
			log.Printf("[DEBUG] GetSnapshotConsistencyGroupWithContext failed %s\n%s", err, response)
			return nil, "failed", fmt.Errorf("[ERROR] Error GetSnapshotConsistencyGroupWithContext failed : %s\n%s", err, response)
		}

		if *snapshotConsistencyGroup.LifecycleState == "stable" {
			return snapshotConsistencyGroup, *snapshotConsistencyGroup.LifecycleState, nil
		} else if *snapshotConsistencyGroup.LifecycleState == "failed" {
			return snapshotConsistencyGroup, *snapshotConsistencyGroup.LifecycleState,
				fmt.Errorf("Snapshot Consistency Group (%s) went into failed state during the operation \n [WARNING] Running terraform apply again will remove the tainted snapshot consistency group and attempt to create the snapshot again replacing the previous configuration", *snapshotConsistencyGroup.ID)
		}

		return snapshotConsistencyGroup, "pending", nil
	}
}

func isWaitForSnapshotConsistencyGroupUpdate(sess *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Snapshot Consistency Group (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating"},
		Target:     []string{"stable", "failed"},
		Refresh:    isSnapshotUpdateConsistencyGroupRefreshFunc(sess, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func isSnapshotUpdateConsistencyGroupRefreshFunc(vpcClient *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getSnapshotConsistencyGroupOptions := &vpcv1.GetSnapshotConsistencyGroupOptions{}
		getSnapshotConsistencyGroupOptions.SetID(id)

		snapshotConsistencyGroup, response, err := vpcClient.GetSnapshotConsistencyGroup(getSnapshotConsistencyGroupOptions)
		if err != nil {
			log.Printf("[DEBUG] GetSnapshotConsistencyGroupWithContext failed %s\n%s", err, response)
			return nil, "failed", fmt.Errorf("[ERROR] Error GetSnapshotConsistencyGroupWithContext failed : %s\n%s", err, response)
		}

		if *snapshotConsistencyGroup.LifecycleState == "stable" || *snapshotConsistencyGroup.LifecycleState == "failed" {
			return snapshotConsistencyGroup, *snapshotConsistencyGroup.LifecycleState, nil
		} else if *snapshotConsistencyGroup.LifecycleState == isSnapshotFailed {
			return snapshotConsistencyGroup, *snapshotConsistencyGroup.LifecycleState, fmt.Errorf("Snapshot Consistency Group (%s) went into failed state during the operation \n [WARNING] Running terraform apply again will remove the tainted snapshot consistency group and attempt to create the snapshot consistency group again replacing the previous configuration", *snapshotConsistencyGroup.ID)
		}

		return snapshotConsistencyGroup, isSnapshotUpdating, nil
	}
}

func isWaitForSnapshotConsistencyGroupDeleted(sess *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Snapshot Consistency Group (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted", "failed"},
		Refresh:    isSnapshotDeleteConsistencyGroupRefreshFunc(sess, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isSnapshotDeleteConsistencyGroupRefreshFunc(vpcClient *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Refresh function for Snapshot Consistency Group delete.")

		getSnapshotConsistencyGroupOptions := &vpcv1.GetSnapshotConsistencyGroupOptions{}
		getSnapshotConsistencyGroupOptions.SetID(id)

		snapshotConsistencyGroup, response, err := vpcClient.GetSnapshotConsistencyGroup(getSnapshotConsistencyGroupOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return snapshotConsistencyGroup, "deleted", nil
			}
			return nil, "failed", fmt.Errorf("[ERROR] The Snapshot Consistency Group %s failed to delete: %s\n%s", id, err, response)
		}
		return snapshotConsistencyGroup, *snapshotConsistencyGroup.LifecycleState, nil
	}
}

func resourceIBMIsSnapshotConsistencyGroupRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getSnapshotConsistencyGroupOptions := &vpcv1.GetSnapshotConsistencyGroupOptions{}
	getSnapshotConsistencyGroupOptions.SetID(d.Id())

	snapshotConsistencyGroup, response, err := vpcClient.GetSnapshotConsistencyGroupWithContext(context, getSnapshotConsistencyGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetSnapshotConsistencyGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSnapshotConsistencyGroupWithContext failed %s\n%s", err, response))
	}

	if !core.IsNil(snapshotConsistencyGroup.DeleteSnapshotsOnDelete) {
		if err = d.Set("delete_snapshots_on_delete", snapshotConsistencyGroup.DeleteSnapshotsOnDelete); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting delete_snapshots_on_delete: %s", err))
		}
	}
	if !core.IsNil(snapshotConsistencyGroup.Name) {
		if err = d.Set("name", snapshotConsistencyGroup.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}
	}
	if !core.IsNil(snapshotConsistencyGroup.ResourceGroup) && !core.IsNil(snapshotConsistencyGroup.ResourceGroup.ID) {
		d.Set("resource_group", snapshotConsistencyGroup.ResourceGroup.ID)
	}
	if !core.IsNil(snapshotConsistencyGroup.Snapshots) {
		snapshots := []map[string]interface{}{}
		for _, snapshotsItem := range snapshotConsistencyGroup.Snapshots {
			snapshotsItemMap, err := resourceIBMIsSnapshotConsistencyGroupSnapshotConsistencyGroupSnapshotsItemToMap(&snapshotsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			snapshots = append(snapshots, snapshotsItemMap)
		}
		if err = d.Set("snapshot_reference", snapshots); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting snapshots: %s", err))
		}
	}
	if !core.IsNil(snapshotConsistencyGroup.BackupPolicyPlan) {
		backupPolicyPlanMap, err := resourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanReferenceToMap(snapshotConsistencyGroup.BackupPolicyPlan)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("backup_policy_plan", []map[string]interface{}{backupPolicyPlanMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting backup_policy_plan: %s", err))
		}
	} else {
		if err = d.Set("backup_policy_plan", []map[string]interface{}{}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting backup_policy_plan: %s", err))
		}
	}
	if err = d.Set("created_at", flex.DateTimeToString(snapshotConsistencyGroup.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("crn", snapshotConsistencyGroup.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("href", snapshotConsistencyGroup.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", snapshotConsistencyGroup.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", snapshotConsistencyGroup.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if err = d.Set("service_tags", snapshotConsistencyGroup.ServiceTags); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting service_tags: %s", err))
	}
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *snapshotConsistencyGroup.CRN, "", isUserTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc snapshot consistency group (%s) tags: %s", d.Id(), err)
	}
	d.Set("tags", tags)

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *snapshotConsistencyGroup.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource VPC snapshot consistency group (%s) access tags: %s", d.Id(), err)
	}
	d.Set("access_tags", accesstags)
	return nil
}

func resourceIBMIsSnapshotConsistencyGroupUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateSnapshotConsistencyGroupOptions := &vpcv1.UpdateSnapshotConsistencyGroupOptions{}
	updateSnapshotConsistencyGroupOptions.SetID(d.Id())
	hasChange := false

	if d.HasChange("tags") {
		getSnapshotConsistencyGroupOptions := &vpcv1.GetSnapshotConsistencyGroupOptions{}
		getSnapshotConsistencyGroupOptions.SetID(d.Id())

		snapshotConsistencyGroup, response, err := vpcClient.GetSnapshotConsistencyGroupWithContext(context, getSnapshotConsistencyGroupOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			log.Printf("[DEBUG] GetSnapshotConsistencyGroupWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("GetSnapshotConsistencyGroupWithContext failed %s\n%s", err, response))
		}

		oldList, newList := d.GetChange("tags")
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *snapshotConsistencyGroup.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on update of resource vpc snapshot consistency group (%s) tags: %s", d.Id(), err)
		}
	}
	if d.HasChange("access_tags") {
		getSnapshotConsistencyGroupOptions := &vpcv1.GetSnapshotConsistencyGroupOptions{}
		getSnapshotConsistencyGroupOptions.SetID(d.Id())

		snapshotConsistencyGroup, response, err := vpcClient.GetSnapshotConsistencyGroupWithContext(context, getSnapshotConsistencyGroupOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			log.Printf("[DEBUG] GetSnapshotConsistencyGroupWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("GetSnapshotConsistencyGroupWithContext failed %s\n%s", err, response))
		}

		oldList, newList := d.GetChange("access_tags")
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *snapshotConsistencyGroup.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on update of resource VPC snapshot consistency group  (%s) access tags: %s", d.Id(), err)
		}
	}

	patchVals := &vpcv1.SnapshotConsistencyGroupPatch{}
	if d.HasChange("delete_snapshots_on_delete") {
		newDeleteSnapshotsOnDelete := d.Get("delete_snapshots_on_delete").(bool)
		patchVals.DeleteSnapshotsOnDelete = &newDeleteSnapshotsOnDelete
		hasChange = true
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}

	getSnapshotConsistencyGroupOptions := &vpcv1.GetSnapshotConsistencyGroupOptions{}
	getSnapshotConsistencyGroupOptions.SetID(d.Id())

	_, response, err := vpcClient.GetSnapshotConsistencyGroupWithContext(context, getSnapshotConsistencyGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetSnapshotConsistencyGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSnapshotConsistencyGroupWithContext failed %s\n%s", err, response))
	}
	updateSnapshotConsistencyGroupOptions.SetIfMatch(response.Headers.Get("Etag"))

	if hasChange {
		updateSnapshotConsistencyGroupOptions.SnapshotConsistencyGroupPatch, _ = patchVals.AsPatch()
		_, response, err := vpcClient.UpdateSnapshotConsistencyGroupWithContext(context, updateSnapshotConsistencyGroupOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateSnapshotConsistencyGroupWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateSnapshotConsistencyGroupWithContext failed %s\n%s", err, response))
		}
		_, err = isWaitForSnapshotConsistencyGroupUpdate(vpcClient, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMIsSnapshotConsistencyGroupRead(context, d, meta)
}

func resourceIBMIsSnapshotConsistencyGroupDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteSnapshotConsistencyGroupOptions := &vpcv1.DeleteSnapshotConsistencyGroupOptions{}

	deleteSnapshotConsistencyGroupOptions.SetID(d.Id())

	_, response, err := vpcClient.DeleteSnapshotConsistencyGroupWithContext(context, deleteSnapshotConsistencyGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteSnapshotConsistencyGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteSnapshotConsistencyGroupWithContext failed %s\n%s", err, response))
	}

	_, err = isWaitForSnapshotConsistencyGroupDeleted(vpcClient, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceIBMIsSnapshotConsistencyGroupResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func resourceIBMIsSnapshotConsistencyGroupSnapshotConsistencyGroupSnapshotsItemToMap(model *vpcv1.SnapshotReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsSnapshotConsistencyGroupSnapshotReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Remote != nil {
		remoteMap, err := resourceIBMIsSnapshotConsistencyGroupSnapshotRemoteToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = model.ResourceType
	}
	return modelMap, nil
}

func resourceIBMIsSnapshotConsistencyGroupSnapshotReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}

func resourceIBMIsSnapshotConsistencyGroupSnapshotRemoteToMap(model *vpcv1.SnapshotRemote) (map[string]interface{}, error) {
	regionMap := make(map[string]interface{})
	if model.Region != nil {
		regionMap, err := resourceIBMIsSnapshotConsistencyGroupRegionReferenceToMap(model.Region)
		if err != nil {
			return regionMap, err
		}
		return regionMap, nil
	}
	return regionMap, nil
}

func resourceIBMIsSnapshotConsistencyGroupRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	modelMap["name"] = model.Name
	return modelMap, nil
}

func resourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanReferenceToMap(model *vpcv1.BackupPolicyPlanReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	if model.Remote != nil {
		remoteMap, err := resourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanRemoteToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func resourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanRemoteToMap(model *vpcv1.BackupPolicyPlanRemote) (map[string]interface{}, error) {
	regionMap := make(map[string]interface{})
	if model.Region != nil {
		regionMap, err := resourceIBMIsSnapshotConsistencyGroupRegionReferenceToMap(model.Region)
		if err != nil {
			return regionMap, err
		}
		return regionMap, nil
	}
	return regionMap, nil
}
