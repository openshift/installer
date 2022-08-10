// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISSnapshotRead,

		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{isSnapshotName, "identifier"},
				Description:  "Snapshot identifier",
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_snapshot", "identifier"),
			},

			isSnapshotName: {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{isSnapshotName, "identifier"},
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_snapshot", isSnapshotName),
				Description:  "Snapshot name",
			},

			isSnapshotResourceGroup: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group info",
			},

			isSnapshotSourceVolume: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Snapshot source volume id",
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

			isSnapshotCapturedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this snapshot was created",
			},

			isSnapshotUserTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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

func DataSourceIBMISSnapshotValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "identifier",
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isSnapshotName,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString})

	ibmISSnapshotDataSourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_snapshot", Schema: validateSchema}
	return &ibmISSnapshotDataSourceValidator
}

func dataSourceIBMISSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get(isSnapshotName).(string)
	id := d.Get("identifier").(string)
	err := snapshotGetByNameOrID(d, meta, name, id)
	if err != nil {
		return err
	}
	return nil
}

func snapshotGetByNameOrID(d *schema.ResourceData, meta interface{}, name, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if name != "" {
		start := ""
		allrecs := []vpcv1.Snapshot{}
		for {
			listSnapshotOptions := &vpcv1.ListSnapshotsOptions{
				Name: &name,
			}
			if start != "" {
				listSnapshotOptions.Start = &start
			}
			snapshots, response, err := sess.ListSnapshots(listSnapshotOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Fetching snapshots %s\n%s", err, response)
			}
			start = flex.GetNext(snapshots.Next)
			allrecs = append(allrecs, snapshots.Snapshots...)
			if start == "" {
				break
			}
		}

		for _, snapshot := range allrecs {
			if *snapshot.Name == name || *snapshot.ID == id {
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
		}
		return fmt.Errorf("[ERROR] No snapshot found with name %s", name)
	} else {
		getSnapshotOptions := &vpcv1.GetSnapshotOptions{
			ID: &id,
		}
		snapshot, response, err := sess.GetSnapshot(getSnapshotOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error fetching snapshot %s\n%s", err, response)
		}
		if (response != nil && response.StatusCode == 404) || snapshot == nil {
			return fmt.Errorf("[ERROR] No snapshot found with id %s", id)
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
		if snapshot.CapturedAt != nil {
			d.Set(isSnapshotCapturedAt, (*snapshot.CapturedAt).String())
		}
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
}
