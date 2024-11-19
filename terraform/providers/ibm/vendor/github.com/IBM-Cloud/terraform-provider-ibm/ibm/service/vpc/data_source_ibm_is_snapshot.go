// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"

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

			isSnapshotCopies: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The copies of this snapshot in other regions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for the copied snapshot.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and provides some supplementary information.",
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
							Description: "The URL for the copied snapshot.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for the copied snapshot.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for the copied snapshot. The name is unique across all snapshots in the copied snapshot's native region.",
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

			isSnapshotName: {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{isSnapshotName, "identifier"},
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_snapshot", isSnapshotName),
				Description:  "Snapshot name",
			},

			"service_tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The [service tags](https://cloud.ibm.com/apidocs/tagging#types-of-tags) prefixed with `is.snapshot:` associated with this snapshot.",
				Elem:        &schema.Schema{Type: schema.TypeString},
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
			isSnapshotSourceSnapshot: {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "If present, the source snapshot this snapshot was created from.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The CRN of the source snapshot.",
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
							Description: "The URL for the source snapshot.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for the source snapshot.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for the source snapshot. The name is unique across all snapshots in the source snapshot's native region.",
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
			isSnapshotConsistencyGroup: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The snapshot consistency group which created this snapshot.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN of this snapshot consistency group.",
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
							Description: "The URL for the snapshot consistency group.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for the snapshot consistency group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for the snapshot consistency group. The name is unique across all snapshot consistency groups in the region.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			isSnapshotSourceImage: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If present, the image id from which the data on this volume was most directly provisioned.",
			},

			isSnapshotAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access tags",
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
			isSnapshotClones: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Zones for creating the snapshot clone",
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
			isSnapshotCatalogOffering: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The catalog offering inherited from the snapshot's source. If a virtual server instance is provisioned with a source_snapshot specifying this snapshot, the virtual server instance will use this snapshot's catalog offering, including its pricing plan.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isSnapshotCatalogOfferingPlanCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this catalog offering version's billing plan",
						},
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and provides some supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						isSnapshotCatalogOfferingVersionCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this version of a catalog offering",
						},
					},
				},
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
				if snapshot.EncryptionKey != nil && snapshot.EncryptionKey.CRN != nil {
					d.Set(isSnapshotEncryptionKey, *snapshot.EncryptionKey.CRN)
				}
				d.Set(isSnapshotLCState, *snapshot.LifecycleState)
				d.Set(isSnapshotResourceType, *snapshot.ResourceType)
				d.Set(isSnapshotBootable, *snapshot.Bootable)

				// source snapshot
				sourceSnapshotList := []map[string]interface{}{}
				if snapshot.SourceSnapshot != nil {
					sourceSnapshot := map[string]interface{}{}
					sourceSnapshot["href"] = snapshot.SourceSnapshot.Href
					if snapshot.SourceSnapshot.Deleted != nil {
						snapshotSourceSnapshotDeletedMap := map[string]interface{}{}
						snapshotSourceSnapshotDeletedMap["more_info"] = snapshot.SourceSnapshot.Deleted.MoreInfo
						sourceSnapshot["deleted"] = []map[string]interface{}{snapshotSourceSnapshotDeletedMap}
					}
					sourceSnapshot["id"] = snapshot.SourceSnapshot.ID
					sourceSnapshot["name"] = snapshot.SourceSnapshot.Name

					sourceSnapshot["resource_type"] = snapshot.SourceSnapshot.ResourceType
					sourceSnapshotList = append(sourceSnapshotList, sourceSnapshot)
				}
				d.Set(isSnapshotSourceSnapshot, sourceSnapshotList)

				// snapshot consistency group
				snapshotConsistencyGroupList := []map[string]interface{}{}
				if snapshot.SnapshotConsistencyGroup != nil {
					snapshotConsistencyGroup := map[string]interface{}{}
					snapshotConsistencyGroup["href"] = snapshot.SnapshotConsistencyGroup.Href
					snapshotConsistencyGroup["crn"] = snapshot.SnapshotConsistencyGroup.CRN
					if snapshot.SnapshotConsistencyGroup.Deleted != nil {
						snapshotConsistencyGroupDeletedMap := map[string]interface{}{}
						snapshotConsistencyGroupDeletedMap["more_info"] = snapshot.SnapshotConsistencyGroup.Deleted.MoreInfo
						snapshotConsistencyGroup["deleted"] = []map[string]interface{}{snapshotConsistencyGroupDeletedMap}
					}
					snapshotConsistencyGroup["id"] = snapshot.SnapshotConsistencyGroup.ID
					snapshotConsistencyGroup["name"] = snapshot.SnapshotConsistencyGroup.Name
					snapshotConsistencyGroup["resource_type"] = snapshot.SnapshotConsistencyGroup.ResourceType
					snapshotConsistencyGroupList = append(snapshotConsistencyGroupList, snapshotConsistencyGroup)
				}
				d.Set(isSnapshotConsistencyGroup, snapshotConsistencyGroupList)

				// snapshot copies
				snapshotCopies := []map[string]interface{}{}
				if snapshot.Copies != nil {
					for _, copiesItem := range snapshot.Copies {
						copiesMap, err := dataSourceIBMIsSnapshotsSnapshotCopiesItemToMap(&copiesItem)
						if err != nil {
							return fmt.Errorf("[ERROR] Error fetching snapshot copies: %s", err)
						}
						snapshotCopies = append(snapshotCopies, copiesMap)
					}
				}
				d.Set(isSnapshotCopies, snapshotCopies)

				if snapshot.UserTags != nil {
					if err = d.Set(isSnapshotUserTags, snapshot.UserTags); err != nil {
						return fmt.Errorf("[ERROR] Error setting user tags: %s", err)
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
				// catalog
				if snapshot.CatalogOffering != nil {
					versionCrn := ""
					if snapshot.CatalogOffering.Version != nil && snapshot.CatalogOffering.Version.CRN != nil {
						versionCrn = *snapshot.CatalogOffering.Version.CRN
					}
					catalogList := make([]map[string]interface{}, 0)
					catalogMap := map[string]interface{}{}
					if versionCrn != "" {
						catalogMap[isSnapshotCatalogOfferingVersionCrn] = versionCrn
					}
					if snapshot.CatalogOffering.Plan != nil {
						planCrn := ""
						if snapshot.CatalogOffering.Plan != nil && snapshot.CatalogOffering.Plan.CRN != nil {
							planCrn = *snapshot.CatalogOffering.Plan.CRN
						}
						if planCrn != "" {
							catalogMap[isSnapshotCatalogOfferingPlanCrn] = planCrn
						}
						if snapshot.CatalogOffering.Plan.Deleted != nil {
							deletedMap := resourceIbmIsSnapshotCatalogOfferingVersionPlanReferenceDeletedToMap(*snapshot.CatalogOffering.Plan.Deleted)
							catalogMap["deleted"] = []map[string]interface{}{deletedMap}
						}
					}
					catalogList = append(catalogList, catalogMap)
					d.Set(isSnapshotCatalogOffering, catalogList)
				}

				var clones []string
				clones = make([]string, 0)
				if snapshot.Clones != nil {
					for _, clone := range snapshot.Clones {
						if clone.Zone != nil && clone.Zone.Name != nil {
							clones = append(clones, *clone.Zone.Name)
						}
					}
				}
				d.Set(isSnapshotClones, flex.NewStringSet(schema.HashString, clones))
				if err = d.Set("service_tags", snapshot.ServiceTags); err != nil {
					return fmt.Errorf("[ERROR] Error setting service_tags: %s", err)
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
				accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *snapshot.CRN, "", isAccessTagType)
				if err != nil {
					log.Printf(
						"[ERROR] Error on get of resource snapshot (%s) access tags: %s", d.Id(), err)
				}
				d.Set(isSnapshotAccessTags, accesstags)
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

		if snapshot.EncryptionKey != nil && snapshot.EncryptionKey.CRN != nil {
			d.Set(isSnapshotEncryptionKey, *snapshot.EncryptionKey.CRN)
		}
		if err = d.Set("service_tags", snapshot.ServiceTags); err != nil {
			return fmt.Errorf("Error setting service_tags: %s", err)
		}
		// source snapshot
		sourceSnapshotList := []map[string]interface{}{}
		if snapshot.SourceSnapshot != nil {
			sourceSnapshot := map[string]interface{}{}
			sourceSnapshot["href"] = snapshot.SourceSnapshot.Href
			sourceSnapshot["crn"] = snapshot.SourceSnapshot.CRN
			if snapshot.SourceSnapshot.Deleted != nil {
				snapshotSourceSnapshotDeletedMap := map[string]interface{}{}
				snapshotSourceSnapshotDeletedMap["more_info"] = snapshot.SourceSnapshot.Deleted.MoreInfo
				sourceSnapshot["deleted"] = []map[string]interface{}{snapshotSourceSnapshotDeletedMap}
			}
			sourceSnapshot["id"] = snapshot.SourceSnapshot.ID
			sourceSnapshot["name"] = snapshot.SourceSnapshot.Name
			sourceSnapshot["resource_type"] = snapshot.SourceSnapshot.ResourceType
			sourceSnapshotList = append(sourceSnapshotList, sourceSnapshot)
		}
		// snapshot copies
		snapshotCopies := []map[string]interface{}{}
		if snapshot.Copies != nil {
			for _, copiesItem := range snapshot.Copies {
				copiesMap, err := dataSourceIBMIsSnapshotsSnapshotCopiesItemToMap(&copiesItem)
				if err != nil {
					return fmt.Errorf("[ERROR] Error fetching snapshot copies: %s", err)
				}
				snapshotCopies = append(snapshotCopies, copiesMap)
			}
		}
		d.Set(isSnapshotCopies, snapshotCopies)

		if snapshot.CapturedAt != nil {
			d.Set(isSnapshotCapturedAt, (*snapshot.CapturedAt).String())
		}
		if snapshot.UserTags != nil {
			if err = d.Set(isSnapshotUserTags, snapshot.UserTags); err != nil {
				return fmt.Errorf("[ERROR] Error setting user tags: %s", err)
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
		var clones []string
		clones = make([]string, 0)
		if snapshot.Clones != nil {
			for _, clone := range snapshot.Clones {
				if clone.Zone != nil && clone.Zone.Name != nil {
					clones = append(clones, *clone.Zone.Name)
				}
			}
		}
		d.Set(isSnapshotClones, flex.NewStringSet(schema.HashString, clones))

		if snapshot.CatalogOffering != nil {
			versionCrn := ""
			if snapshot.CatalogOffering.Version != nil && snapshot.CatalogOffering.Version.CRN != nil {
				versionCrn = *snapshot.CatalogOffering.Version.CRN
			}
			catalogList := make([]map[string]interface{}, 0)
			catalogMap := map[string]interface{}{}
			if versionCrn != "" {
				catalogMap[isSnapshotCatalogOfferingVersionCrn] = versionCrn
			}
			if snapshot.CatalogOffering.Plan != nil {
				planCrn := ""
				if snapshot.CatalogOffering.Plan != nil && snapshot.CatalogOffering.Plan.CRN != nil {
					planCrn = *snapshot.CatalogOffering.Plan.CRN
				}
				if planCrn != "" {
					catalogMap[isSnapshotCatalogOfferingPlanCrn] = planCrn
				}
				if snapshot.CatalogOffering.Plan.Deleted != nil {
					deletedMap := resourceIbmIsSnapshotCatalogOfferingVersionPlanReferenceDeletedToMap(*snapshot.CatalogOffering.Plan.Deleted)
					catalogMap["deleted"] = []map[string]interface{}{deletedMap}
				}
			}
			catalogList = append(catalogList, catalogMap)
			d.Set(isSnapshotCatalogOffering, catalogList)
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
		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *snapshot.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on get of resource snapshot (%s) access tags: %s", d.Id(), err)
		}
		d.Set(isSnapshotAccessTags, accesstags)
		return nil
	}
}

func resourceIbmIsSnapshotCatalogOfferingVersionPlanReferenceDeletedToMap(catalogOfferingVersionPlanReferenceDeleted vpcv1.Deleted) map[string]interface{} {
	catalogOfferingVersionPlanReferenceDeletedMap := map[string]interface{}{}

	catalogOfferingVersionPlanReferenceDeletedMap["more_info"] = catalogOfferingVersionPlanReferenceDeleted.MoreInfo

	return catalogOfferingVersionPlanReferenceDeletedMap
}
