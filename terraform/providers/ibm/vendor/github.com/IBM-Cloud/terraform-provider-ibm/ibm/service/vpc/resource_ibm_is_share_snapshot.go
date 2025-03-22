// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.91.0-d9755c53-20240605-153412
 */

package vpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsShareSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsShareSnapshotCreate,
		ReadContext:   resourceIBMIsShareSnapshotRead,
		UpdateContext: resourceIBMIsShareSnapshotUpdate,
		DeleteContext: resourceIBMIsShareSnapshotDelete,
		Importer:      &schema.ResourceImporter{},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				},
			),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceValidateAccessTags(diff, v)
				}),
		),

		Schema: map[string]*schema.Schema{
			"share": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_share_snapshot", "share"),
				Description:  "The file share identifier.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_share_snapshot", "name"),
				Description:  "The name for this share snapshot. The name is unique across all snapshots for the file share.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_share_snapshot", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "The [user tags](https://cloud.ibm.com/apidocs/tagging#types-of-tags) associated with this share snapshot.",
			},
			"access_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_share_snapshot", "access_tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},
			"backup_policy_plan": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "If present, the backup policy plan which created this share snapshot.",
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
										Computed:    true,
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
							Description: "If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the referenced resource is remote to thisregion, and identifies the native region.",
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
			"captured_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time the data capture for this share snapshot was completed.If absent, this snapshot's data has not yet been captured.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the share snapshot was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this share snapshot.",
			},
			"fingerprint": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fingerprint for this snapshot.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this share snapshot.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of this share snapshot.",
			},
			"minimum_size": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum size of a share created from this snapshot. When a snapshot is created, this will be set to the size of the `source_share`.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group for this file share.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this resource group.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this resource group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this resource group.",
						},
					},
				},
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the share snapshot:- `available`: The share snapshot is available for use.- `failed`: The share snapshot is irrecoverably unusable.- `pending`: The share snapshot is being provisioned and is not yet usable.- `unusable`: The share snapshot is not currently usable (see `status_reasons`)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			},
			"status_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current status (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A reason code for the status:- `encryption_key_deleted`: File share snapshot is unusable  because its `encryption_key` was deleted- `internal_error`: Internal error (contact IBM support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about this status reason.",
						},
					},
				},
			},
			"zone": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The zone this share snapshot resides in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this zone.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this zone.",
						},
					},
				},
			},
			"share_snapshot": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this share snapshot.",
			},
		},
	}
}

func ResourceIBMIsShareSnapshotValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "share",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		validate.ValidateSchema{
			Identifier:                 "access_tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-]):([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-])$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_share_snapshot", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsShareSnapshotCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createShareSnapshotOptions := &vpcv1.CreateShareSnapshotOptions{}

	createShareSnapshotOptions.SetShareID(d.Get("share").(string))
	if _, ok := d.GetOk("name"); ok {
		createShareSnapshotOptions.SetName(d.Get("name").(string))
	}
	var userTags *schema.Set
	if v, ok := d.GetOk("tags"); ok {
		userTags = v.(*schema.Set)
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
			createShareSnapshotOptions.SetUserTags(userTagsArray)
		}
	}
	shareSnapshot, _, err := vpcClient.CreateShareSnapshotWithContext(context, createShareSnapshotOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateShareSnapshotWithContext failed: %s", err.Error()), "ibm_is_share_snapshot", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createShareSnapshotOptions.ShareID, *shareSnapshot.ID))
	_, err = isWaitForShareSnapshotAvailable(context, vpcClient, *createShareSnapshotOptions.ShareID, *shareSnapshot.ID, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	if _, ok := d.GetOk("access_tags"); ok {
		oldList, newList := d.GetChange(isSubnetAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *shareSnapshot.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on create of resource share snapshot (%s) access tags: %s", d.Id(), err)
		}
	}
	return resourceIBMIsShareSnapshotRead(context, d, meta)
}

func resourceIBMIsShareSnapshotRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getShareSnapshotOptions := &vpcv1.GetShareSnapshotOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "sep-id-parts").GetDiag()
	}

	getShareSnapshotOptions.SetShareID(parts[0])
	getShareSnapshotOptions.SetID(parts[1])

	shareSnapshot, response, err := vpcClient.GetShareSnapshotWithContext(context, getShareSnapshotOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetShareSnapshotWithContext failed: %s", err.Error()), "ibm_is_share_snapshot", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(shareSnapshot.Name) {
		if err = d.Set("name", shareSnapshot.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-name").GetDiag()
		}
	}
	if !core.IsNil(shareSnapshot.UserTags) {
		if err = d.Set("tags", shareSnapshot.UserTags); err != nil {
			err = fmt.Errorf("Error setting tags: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-tags").GetDiag()
		}
	}
	backupPolicyPlanMap := make(map[string]interface{}, 1)
	if !core.IsNil(shareSnapshot.BackupPolicyPlan) {
		backupPolicyPlanMap, err = ResourceIBMIsShareSnapshotBackupPolicyPlanReferenceToMap(shareSnapshot.BackupPolicyPlan)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "backup_policy_plan-to-map").GetDiag()
		}
	}
	if err = d.Set("backup_policy_plan", []map[string]interface{}{backupPolicyPlanMap}); err != nil {
		err = fmt.Errorf("Error setting backup_policy_plan: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-backup_policy_plan").GetDiag()
	}

	if !core.IsNil(shareSnapshot.CapturedAt) {
		if err = d.Set("captured_at", flex.DateTimeToString(shareSnapshot.CapturedAt)); err != nil {
			err = fmt.Errorf("Error setting captured_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-captured_at").GetDiag()
		}
	}
	if err = d.Set("created_at", flex.DateTimeToString(shareSnapshot.CreatedAt)); err != nil {
		err = fmt.Errorf("Error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-created_at").GetDiag()
	}
	if err = d.Set("crn", shareSnapshot.CRN); err != nil {
		err = fmt.Errorf("Error setting crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-crn").GetDiag()
	}
	if err = d.Set("fingerprint", shareSnapshot.Fingerprint); err != nil {
		err = fmt.Errorf("Error setting fingerprint: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-fingerprint").GetDiag()
	}
	if err = d.Set("href", shareSnapshot.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-href").GetDiag()
	}
	if err = d.Set("lifecycle_state", shareSnapshot.LifecycleState); err != nil {
		err = fmt.Errorf("Error setting lifecycle_state: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-lifecycle_state").GetDiag()
	}
	if err = d.Set("minimum_size", flex.IntValue(shareSnapshot.MinimumSize)); err != nil {
		err = fmt.Errorf("Error setting minimum_size: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-minimum_size").GetDiag()
	}
	if !core.IsNil(shareSnapshot.ResourceGroup) {
		resourceGroupMap, err := ResourceIBMIsShareSnapshotResourceGroupReferenceToMap(shareSnapshot.ResourceGroup)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "resource_group-to-map").GetDiag()
		}
		if err = d.Set("resource_group", []map[string]interface{}{resourceGroupMap}); err != nil {
			err = fmt.Errorf("Error setting resource_group: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-resource_group").GetDiag()
		}
	}
	if err = d.Set("resource_type", shareSnapshot.ResourceType); err != nil {
		err = fmt.Errorf("Error setting resource_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-resource_type").GetDiag()
	}
	if err = d.Set("status", shareSnapshot.Status); err != nil {
		err = fmt.Errorf("Error setting status: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-status").GetDiag()
	}
	statusReasons := []map[string]interface{}{}
	for _, statusReasonsItem := range shareSnapshot.StatusReasons {
		statusReasonsItemMap, err := ResourceIBMIsShareSnapshotShareSnapshotStatusReasonToMap(&statusReasonsItem)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "status_reasons-to-map").GetDiag()
		}
		statusReasons = append(statusReasons, statusReasonsItemMap)
	}
	if err = d.Set("status_reasons", statusReasons); err != nil {
		err = fmt.Errorf("Error setting status_reasons: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-status_reasons").GetDiag()
	}
	zoneMap, err := ResourceIBMIsShareSnapshotZoneReferenceToMap(shareSnapshot.Zone)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "zone-to-map").GetDiag()
	}
	if err = d.Set("zone", []map[string]interface{}{zoneMap}); err != nil {
		err = fmt.Errorf("Error setting zone: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-zone").GetDiag()
	}
	if err = d.Set("share_snapshot", shareSnapshot.ID); err != nil {
		err = fmt.Errorf("Error setting share_snapshot: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "read", "set-is_share_snapshot_id").GetDiag()
	}
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *shareSnapshot.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"[ERROR] Error on get of resource share snapshot (%s) access tags: %s", d.Id(), err)
	}
	d.Set("access_tags", accesstags)
	return nil
}

func resourceIBMIsShareSnapshotUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateShareSnapshotOptions := &vpcv1.UpdateShareSnapshotOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "update", "sep-id-parts").GetDiag()
	}

	updateShareSnapshotOptions.SetShareID(parts[0])
	updateShareSnapshotOptions.SetID(parts[1])

	hasChange := false

	getShareSnapshotOptions := &vpcv1.GetShareSnapshotOptions{}

	getShareSnapshotOptions.SetShareID(parts[0])
	getShareSnapshotOptions.SetID(parts[1])

	_, response, err := vpcClient.GetShareSnapshotWithContext(context, getShareSnapshotOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetShareSnapshotWithContext failed: %s", err.Error()), "ibm_is_share_snapshot", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	patchVals := &vpcv1.ShareSnapshotPatch{}

	if d.HasChange("tags") {
		var userTags *schema.Set
		if v, ok := d.GetOk("tags"); ok {

			userTags = v.(*schema.Set)
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
				patchVals.UserTags = userTagsArray
				hasChange = true
			}
		}
	}
	updateShareSnapshotOptions.SetIfMatch(response.Headers.Get("Etag"))

	if hasChange {
		updateShareSnapshotOptions.ShareSnapshotPatch, _ = patchVals.AsPatch()
		if _, exists := d.GetOk("tags"); d.HasChange("tags") && !exists {
			updateShareSnapshotOptions.ShareSnapshotPatch["tags"] = nil
		}

		_, _, err = vpcClient.UpdateShareSnapshotWithContext(context, updateShareSnapshotOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateShareSnapshotWithContext failed: %s", err.Error()), "ibm_is_share_snapshot", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		_, err = isWaitForShareSnapshotAvailable(context, vpcClient, parts[0], parts[1], d, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("access_tags") {
		oldList, newList := d.GetChange("access_tags")
		err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, d.Get(isSnapshotCRN).(string), "", isAccessTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on update of resource share snapshot (%s) access tags: %s", d.Id(), err)
		}
	}
	return resourceIBMIsShareSnapshotRead(context, d, meta)
}

func resourceIBMIsShareSnapshotDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteShareSnapshotOptions := &vpcv1.DeleteShareSnapshotOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_share_snapshot", "delete", "sep-id-parts").GetDiag()
	}

	deleteShareSnapshotOptions.SetShareID(parts[0])
	deleteShareSnapshotOptions.SetID(parts[1])

	_, _, err = vpcClient.DeleteShareSnapshotWithContext(context, deleteShareSnapshotOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteShareSnapshotWithContext failed: %s", err.Error()), "ibm_is_share_snapshot", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMIsShareSnapshotBackupPolicyPlanReferenceToMap(model *vpcv1.BackupPolicyPlanReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsShareSnapshotDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Remote != nil {
		remoteMap, err := ResourceIBMIsShareSnapshotBackupPolicyPlanRemoteToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsShareSnapshotDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func ResourceIBMIsShareSnapshotBackupPolicyPlanRemoteToMap(model *vpcv1.BackupPolicyPlanRemote) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Region != nil {
		regionMap, err := ResourceIBMIsShareSnapshotRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func ResourceIBMIsShareSnapshotRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIBMIsShareSnapshotResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIBMIsShareSnapshotShareSnapshotStatusReasonToMap(model *vpcv1.ShareSnapshotStatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func ResourceIBMIsShareSnapshotZoneReferenceToMap(model *vpcv1.ZoneReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func isWaitForShareSnapshotAvailable(context context.Context, vpcClient *vpcv1.VpcV1, shareid string, shareSnapshotId string, d *schema.ResourceData, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for share snapshot (%s) to be available.", shareid)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating", "pending", "waiting"},
		Target:     []string{"available", "failed"},
		Refresh:    isShareSnapshotRefreshFunc(context, vpcClient, shareid, shareSnapshotId, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isShareSnapshotRefreshFunc(context context.Context, vpcClient *vpcv1.VpcV1, shareid, shareSnapshotId string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		shareSnapOptions := &vpcv1.GetShareSnapshotOptions{}

		shareSnapOptions.SetShareID(shareid)
		shareSnapOptions.SetID(shareSnapshotId)

		shareSnapshot, response, err := vpcClient.GetShareSnapshotWithContext(context, shareSnapOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting share snapshot: %s\n%s", err, response)
		}
		d.Set("status", *shareSnapshot.Status)
		if *shareSnapshot.Status == "available" || *shareSnapshot.Status == "failed" {

			if *shareSnapshot.Status == "available" {
				if _, ok := d.GetOk("tags"); ok && len(shareSnapshot.UserTags) == 0 {
					return shareSnapshot, "pending", nil
				}
			}
			return shareSnapshot, *shareSnapshot.Status, nil

		}
		return shareSnapshot, "pending", nil
	}
}
