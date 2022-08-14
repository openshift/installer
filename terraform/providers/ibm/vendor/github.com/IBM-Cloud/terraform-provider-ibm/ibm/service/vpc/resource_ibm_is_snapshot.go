// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isSnapshotName             = "name"
	isSnapshotResourceGroup    = "resource_group"
	isSnapshotSourceVolume     = "source_volume"
	isSnapshotSourceImage      = "source_image"
	isSnapshotUserTags         = "tags"
	isSnapshotCRN              = "crn"
	isSnapshotHref             = "href"
	isSnapshotEncryption       = "encryption"
	isSnapshotEncryptionKey    = "encryption_key"
	isSnapshotOperatingSystem  = "operating_system"
	isSnapshotLCState          = "lifecycle_state"
	isSnapshotMinCapacity      = "minimum_capacity"
	isSnapshotResourceType     = "resource_type"
	isSnapshotSize             = "size"
	isSnapshotBootable         = "bootable"
	isSnapshotDeleting         = "deleting"
	isSnapshotDeleted          = "deleted"
	isSnapshotAvailable        = "stable"
	isSnapshotFailed           = "failed"
	isSnapshotPending          = "pending"
	isSnapshotSuspended        = "suspended"
	isSnapshotUpdating         = "updating"
	isSnapshotWaiting          = "waiting"
	isSnapshotCapturedAt       = "captured_at"
	isSnapshotBackupPolicyPlan = "backup_policy_plan"
)

func ResourceIBMSnapshot() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISSnapshotCreate,
		Read:     resourceIBMISSnapshotRead,
		Update:   resourceIBMISSnapshotUpdate,
		Delete:   resourceIBMISSnapshotDelete,
		Exists:   resourceIBMISSnapshotExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{

			isSnapshotName: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_snapshot", isSnapshotName),
				Description:  "Snapshot name",
			},

			isSnapshotResourceGroup: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Resource group info",
			},

			isSnapshotSourceVolume: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Snapshot source volume",
			},

			isSnapshotSourceImage: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If present, the image id from which the data on this volume was most directly provisioned.",
			},

			isSnapshotOperatingSystem: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique name for the operating system included in this image",
			},

			isSnapshotBootable: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if a boot volume attachment can be created with a volume created from this snapshot",
			},

			isSnapshotLCState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Snapshot lifecycle state",
			},
			isSnapshotCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},
			isSnapshotEncryption: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Encryption type of the snapshot",
			},
			isSnapshotEncryptionKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A reference to the root key used to wrap the data encryption key for the source volume.",
			},

			isSnapshotHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for the snapshot",
			},

			isSnapshotMinCapacity: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Minimum capacity of the snapshot",
			},
			isSnapshotResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type of the snapshot",
			},

			isSnapshotSize: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the snapshot",
			},

			isSnapshotUserTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_snapshot", isSnapshotUserTags)},
				Set:         flex.ResourceIBMVPCHash,
				Description: "User Tags for the snapshot",
			},

			isSnapshotBackupPolicyPlan: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "If present, the backup policy plan which created this snapshot.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and provides some supplementary information.",
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
							Description: "The unique user-defined name for this backup policy plan.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced",
						},
					},
				},
			},
		},
	}
}

func ResourceIBMISSnapshotValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isSnapshotName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isSnapshotUserTags,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})
	ibmISSnapshotResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_snapshot", Schema: validateSchema}
	return &ibmISSnapshotResourceValidator
}

func resourceIBMISSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.CreateSnapshotOptions{}
	snapshotprototypeoptions := &vpcv1.SnapshotPrototypeSnapshotBySourceVolume{}
	if snapshotName, ok := d.GetOk(isSnapshotName); ok {
		name := snapshotName.(string)
		snapshotprototypeoptions.Name = &name
	}
	if sourceVolume, ok := d.GetOk(isSnapshotSourceVolume); ok {
		sv := sourceVolume.(string)
		snapshotprototypeoptions.SourceVolume = &vpcv1.VolumeIdentity{
			ID: &sv,
		}
	}
	if grp, ok := d.GetOk(isVPCResourceGroup); ok {
		rg := grp.(string)
		snapshotprototypeoptions.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}

	var userTags *schema.Set
	if v, ok := d.GetOk(isSnapshotUserTags); ok {
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
			snapshotprototypeoptions.UserTags = userTagsArray
		}
	}
	options.SnapshotPrototype = snapshotprototypeoptions

	log.Printf("[DEBUG] Snapshot create")

	snapshot, response, err := sess.CreateSnapshot(options)
	if err != nil || snapshot == nil {
		return fmt.Errorf("[ERROR] Error creating Snapshot %s\n%s", err, response)
	}

	d.SetId(*snapshot.ID)
	log.Printf("[INFO] Snapshot : %s", *snapshot.ID)

	_, err = isWaitForSnapshotAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return err
	}

	return resourceIBMISSnapshotRead(d, meta)
}

func isWaitForSnapshotAvailable(sess *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Snapshot (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isSnapshotPending},
		Target:     []string{isSnapshotAvailable, isSnapshotFailed},
		Refresh:    isSnapshotRefreshFunc(sess, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isSnapshotRefreshFunc(sess *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getSnapshotOptions := &vpcv1.GetSnapshotOptions{
			ID: &id,
		}
		snapshot, response, err := sess.GetSnapshot(getSnapshotOptions)
		if err != nil {
			return nil, isSnapshotFailed, fmt.Errorf("[ERROR] Error getting Snapshot : %s\n%s", err, response)
		}

		if *snapshot.LifecycleState == isSnapshotAvailable {
			return snapshot, *snapshot.LifecycleState, nil
		} else if *snapshot.LifecycleState == isSnapshotFailed {
			return snapshot, *snapshot.LifecycleState, fmt.Errorf("Snapshot (%s) went into failed state during the operation \n [WARNING] Running terraform apply again will remove the tainted snapshot and attempt to create the snapshot again replacing the previous configuration", *snapshot.ID)
		}

		return snapshot, isSnapshotPending, nil
	}
}

func resourceIBMISSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	err := snapshotGet(d, meta, id)
	if err != nil {
		return err
	}
	return nil
}

func snapshotGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getSnapshotOptions := &vpcv1.GetSnapshotOptions{
		ID: &id,
	}
	snapshot, response, err := sess.GetSnapshot(getSnapshotOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Snapshot : %s\n%s", err, response)
	}

	d.SetId(*snapshot.ID)
	d.Set(isSnapshotName, *snapshot.Name)
	d.Set(isSnapshotHref, *snapshot.Href)
	d.Set(isSnapshotCRN, *snapshot.CRN)
	d.Set(isSnapshotMinCapacity, *snapshot.MinimumCapacity)
	d.Set(isSnapshotSize, *snapshot.Size)
	d.Set(isSnapshotEncryption, *snapshot.Encryption)
	d.Set(isSnapshotLCState, *snapshot.LifecycleState)
	d.Set(isSnapshotResourceType, *snapshot.ResourceType)
	d.Set(isSnapshotBootable, *snapshot.Bootable)
	if snapshot.UserTags != nil {
		if err = d.Set(isSnapshotUserTags, snapshot.UserTags); err != nil {
			return fmt.Errorf("Error setting user tags: %s", err)
		}
	}
	if snapshot.ResourceGroup != nil && snapshot.ResourceGroup.ID != nil {
		d.Set(isSnapshotResourceGroup, *snapshot.ResourceGroup.ID)
	}
	if snapshot.SourceVolume != nil && snapshot.SourceVolume.ID != nil {
		d.Set(isSnapshotSourceVolume, *snapshot.SourceVolume.ID)
	}

	if snapshot.SourceImage != nil && snapshot.SourceImage.ID != nil {
		d.Set(isSnapshotSourceImage, *snapshot.SourceImage.ID)
	}

	if snapshot.OperatingSystem != nil && snapshot.OperatingSystem.Name != nil {
		d.Set(isSnapshotOperatingSystem, *snapshot.OperatingSystem.Name)
	}
	backupPolicyPlanList := []map[string]interface{}{}
	if snapshot.BackupPolicyPlan != nil {
		backupPolicyPlan := map[string]interface{}{}
		if snapshot.BackupPolicyPlan.Deleted != nil {
			snapshotBackupPolicyPlanDeletedMap := map[string]interface{}{}
			snapshotBackupPolicyPlanDeletedMap["more_info"] = snapshot.BackupPolicyPlan.Deleted.MoreInfo
			backupPolicyPlan["deleted"] = []map[string]interface{}{snapshotBackupPolicyPlanDeletedMap}
		}
		backupPolicyPlan["href"] = snapshot.BackupPolicyPlan.Href
		backupPolicyPlan["id"] = snapshot.BackupPolicyPlan.ID
		backupPolicyPlan["name"] = snapshot.BackupPolicyPlan.Name
		backupPolicyPlan["resource_type"] = snapshot.BackupPolicyPlan.ResourceType
		backupPolicyPlanList = append(backupPolicyPlanList, backupPolicyPlan)
	}
	d.Set(isSnapshotBackupPolicyPlan, backupPolicyPlanList)
	return nil
}

func resourceIBMISSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	name := ""
	hasChanged := false

	if d.HasChange(isSnapshotName) {
		name = d.Get(isSnapshotName).(string)
		hasChanged = true
	}
	err := snapshotUpdate(d, meta, id, name, hasChanged)
	if err != nil {
		return err
	}
	return resourceIBMISSnapshotRead(d, meta)
}

func snapshotUpdate(d *schema.ResourceData, meta interface{}, id, name string, hasChanged bool) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getSnapshotOptions := &vpcv1.GetSnapshotOptions{
		ID: &id,
	}
	_, response, err := sess.GetSnapshot(getSnapshotOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting Snapshot : %s\n%s", err, response)
	}
	eTag := response.Headers.Get("ETag")

	updateSnapshotOptions := &vpcv1.UpdateSnapshotOptions{
		ID: &id,
	}
	updateSnapshotOptions.IfMatch = &eTag

	// user tags update
	if d.HasChange(isSnapshotUserTags) {
		var userTags *schema.Set
		if v, ok := d.GetOk(isSnapshotUserTags); ok {

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
				snapshotPatchModel := &vpcv1.SnapshotPatch{}
				snapshotPatchModel.UserTags = userTagsArray
				snapshotPatch, err := snapshotPatchModel.AsPatch()
				if err != nil {
					return fmt.Errorf("Error calling asPatch for SnapshotPatch: %s", err)
				}
				updateSnapshotOptions.SnapshotPatch = snapshotPatch
				_, response, err := sess.UpdateSnapshot(updateSnapshotOptions)
				if err != nil {
					return fmt.Errorf("Error updating Snapshot : %s\n%s", err, response)
				}
				_, err = isWaitForSnapshotUpdate(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange(isSnapshotName) {
		updateSnapshotOptions := &vpcv1.UpdateSnapshotOptions{
			ID: &id,
		}
		snapshotPatchModel := &vpcv1.SnapshotPatch{
			Name: &name,
		}
		snapshotPatch, err := snapshotPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for SnapshotPatch: %s", err)
		}
		updateSnapshotOptions.SnapshotPatch = snapshotPatch
		_, response, err := sess.UpdateSnapshot(updateSnapshotOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating Snapshot : %s\n%s", err, response)
		}
		_, err = isWaitForSnapshotUpdate(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}
	return nil
}

func isWaitForSnapshotUpdate(sess *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Snapshot (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isSnapshotUpdating},
		Target:     []string{isSnapshotAvailable, isSnapshotFailed},
		Refresh:    isSnapshotUpdateRefreshFunc(sess, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func isSnapshotUpdateRefreshFunc(sess *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getSnapshotOptions := &vpcv1.GetSnapshotOptions{
			ID: &id,
		}
		snapshot, response, err := sess.GetSnapshot(getSnapshotOptions)
		if err != nil {
			return nil, isSnapshotFailed, fmt.Errorf("[ERROR] Error getting Snapshot : %s\n%s", err, response)
		}

		if *snapshot.LifecycleState == isSnapshotAvailable || *snapshot.LifecycleState == isSnapshotFailed {
			return snapshot, *snapshot.LifecycleState, nil
		} else if *snapshot.LifecycleState == isSnapshotFailed {
			return snapshot, *snapshot.LifecycleState, fmt.Errorf("Snapshot (%s) went into failed state during the operation \n [WARNING] Running terraform apply again will remove the tainted snapshot and attempt to create the snapshot again replacing the previous configuration", *snapshot.ID)
		}

		return snapshot, isSnapshotUpdating, nil
	}
}

func resourceIBMISSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	err := snapshotDelete(d, meta, id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func snapshotDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getSnapshotOptions := &vpcv1.GetSnapshotOptions{
		ID: &id,
	}
	_, response, err := sess.GetSnapshot(getSnapshotOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Snapshot (%s): %s\n%s", id, err, response)
	}

	deleteSnapshotOptions := &vpcv1.DeleteSnapshotOptions{
		ID: &id,
	}
	response, err = sess.DeleteSnapshot(deleteSnapshotOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting Snapshot : %s\n%s", err, response)
	}
	_, err = isWaitForSnapshotDeleted(sess, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForSnapshotDeleted(sess *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Snapshot (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isSnapshotDeleting},
		Target:     []string{isSnapshotDeleted, isSnapshotFailed},
		Refresh:    isSnapshotDeleteRefreshFunc(sess, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isSnapshotDeleteRefreshFunc(sess *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Refresh function for Snapshot delete.")
		getSnapshotOptions := &vpcv1.GetSnapshotOptions{
			ID: &id,
		}
		snapshot, response, err := sess.GetSnapshot(getSnapshotOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return snapshot, isSnapshotDeleted, nil
			}
			return nil, isSnapshotFailed, fmt.Errorf("[ERROR] The Snapshot %s failed to delete: %s\n%s", id, err, response)
		}
		return snapshot, *snapshot.LifecycleState, nil
	}
}

func resourceIBMISSnapshotExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id := d.Id()
	exists, err := snapshotExists(d, meta, id)
	return exists, err
}

func snapshotExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getSnapshotOptions := &vpcv1.GetSnapshotOptions{
		ID: &id,
	}
	_, response, err := sess.GetSnapshot(getSnapshotOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting Snapshot: %s\n%s", err, response)
	}
	return true, nil
}
