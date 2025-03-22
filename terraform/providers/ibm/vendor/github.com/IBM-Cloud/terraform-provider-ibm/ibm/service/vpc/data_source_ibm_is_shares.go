// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIbmIsShares() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsSharesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the share.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource group of the share.",
			},
			"shares": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of file shares.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_transit_encryption_modes": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Allowed transit encryption modes",
						},
						"access_control_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The access control mode for the share",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the file share is created.",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this share.",
						},
						"encryption": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of encryption used for this file share.",
						},
						"encryption_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key used to encrypt this file share. The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this share.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this file share.",
						},
						"iops": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum input/output operation performance bandwidth per second for the file share.",
						},
						"latest_sync": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information about the latest synchronization for this file share.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"completed_at": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The completed date and time of last synchronization between the replica share and its source.",
									},
									"data_transferred": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The data transferred (in bytes) in the last synchronization between the replica and its source.",
									},
									"started_at": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The start date and time of last synchronization between the replica share and its source.",
									},
								},
							},
						},
						"latest_job": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The latest job associated with this file share.This property will be absent if no jobs have been created for this file share.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status of the file share job.The enumerated values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the file share job on which the unexpected property value was encountered.* `cancelled`: This job has been cancelled.* `failed`: This job has failed.* `queued`: This job is queued.* `running`: This job is running.* `succeeded`: This job completed successfully.",
									},
									"status_reasons": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The reasons for the file share job status (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"code": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "A snake case string succinctly identifying the status reason.",
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
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the file share job.The enumerated values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the file share job on which the unexpected property value was encountered.* `replication_failover`: This is a share replication failover job.* `replication_init`: This is a share replication is initialization job.* `replication_split`: This is a share replication split job.",
									},
								},
							},
						},
						"lifecycle_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the file share.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this file share. If unspecified, the name will be a hyphenated list of randomly-selected words.",
						},
						"profile": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name of the profile this file share uses.",
						},
						"replica_share": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The replica file share for this source file share.This property will be present when the `replication_role` is `source`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this file share.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
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
										Description: "The URL for this file share.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this file share.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this file share.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"replication_cron_spec": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cron specification for the file share replication schedule.This property will be present when the `replication_role` is `replica`.",
						},
						"replication_role": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The replication role of the file share.* `none`: This share is not participating in replication.* `replica`: This share is a replication target.* `source`: This share is a replication source.",
						},
						"replication_status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The replication status of the file share.* `active`: This share is actively participating in replication, and the replica's data is up-to-date with the replication schedule.* `failover_pending`: This share is performing a replication failover.* `initializing`: This share is initializing replication.* `none`: This share is not participating in replication.* `split_pending`: This share is performing a replication split.",
						},
						"replication_status_reasons": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current replication status (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the status reason.",
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
						"resource_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the resource group for this file share.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the file share rounded up to the next gigabyte.",
						},
						"source_share": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The source file share for this replica file share.This property will be present when the `replication_role` is `replica`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this file share.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
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
										Description: "The URL for this file share.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this file share.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this file share.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"mount_targets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Mount targets for the file share.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
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
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this share target.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this share target.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this share target.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of resource referenced.",
									},
								},
							},
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name of the zone this file share will reside in.",
						},
						isFileShareAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access management tags",
						},
						isFileShareTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of tags",
						},
						"origin_share": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The origin share this accessor share is referring to.This property will be present when the `accessor_binding_role` is `accessor`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this file share.",
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
										Description: "The URL for this file share.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this file share.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this share. The name is unique across all shares in the region.",
									},
									"remote": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"account": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates that the referenced resource is remote to thisaccount, and identifies the owning account.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique identifier for this account.",
															},
															"resource_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The resource type.",
															},
														},
													},
												},
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
						"accessor_binding_role": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The accessor binding role of this file share:- `none`: This file share is not participating in access with another file share- `origin`: This file share is the origin for one or more file shares  (which may be in other accounts)- `accessor`: This file share is providing access to another file share  (which may be in another account).",
						},
						"accessor_bindings": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The accessor bindings for this file share. Each accessor binding identifies a resource (possibly in another account) with access to this file share's data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this share accessor binding.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this share accessor binding.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"snapshot_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of snapshots for this share.",
						},
						"snapshot_size": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total size (in gigabytes) of snapshots used for this file share.",
						},
						"source_snapshot": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The snapshot from which this share was cloned.This property will be present when the share was created from a snapshot.The resources supported by this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in thefuture.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this share snapshot.",
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
										Description: "The URL for this share snapshot.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this share snapshot.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this share snapshot. The name is unique across all snapshots for the file share.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
					},
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages.",
			},
		},
	}
}

func dataSourceIbmIsSharesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	shareName := ""
	if shareNameIntf, ok := d.GetOk("name"); ok {
		shareName = shareNameIntf.(string)
	}
	resGrp := ""
	if resGrpIntf, ok := d.GetOk("resource_group"); ok {
		resGrp = resGrpIntf.(string)
	}
	listSharesOptions := &vpcv1.ListSharesOptions{}

	if shareName != "" {
		listSharesOptions.Name = &shareName
	}
	if resGrp != "" {
		listSharesOptions.ResourceGroupID = &resGrp
	}
	start := ""
	allrecs := []vpcv1.Share{}
	totalCount := 0
	for {
		if start != "" {
			listSharesOptions.Start = &start
		}
		shareCollection, response, err := vpcClient.ListSharesWithContext(context, listSharesOptions)
		if err != nil {
			log.Printf("[DEBUG] ListSharesWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		if totalCount == 0 {
			totalCount = int(*shareCollection.TotalCount)
		}
		start = flex.GetNext(shareCollection.Next)
		allrecs = append(allrecs, shareCollection.Shares...)
		if start == "" {
			break
		}
	}

	d.SetId(dataSourceIbmIsSharesID(d))

	if allrecs != nil {
		shares, err := dataSourceShareCollectionFlattenShares(meta, allrecs)
		if err != nil {
			return err
		}
		errSettingShares := d.Set("shares", shares)
		if errSettingShares != nil {
			return diag.FromErr(fmt.Errorf("Error setting shares %s", errSettingShares))
		}
	}
	if err = d.Set("total_count", totalCount); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting total_count: %s", err))
	}

	return nil
}

// dataSourceIbmIsSharesID returns a reasonable ID for the list.
func dataSourceIbmIsSharesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceShareCollectionFlattenShares(meta interface{}, result []vpcv1.Share) (shares []map[string]interface{}, diag diag.Diagnostics) {
	for _, sharesItem := range result {
		shareMap, err := dataSourceShareCollectionSharesToMap(meta, sharesItem)
		if err != nil {
			return nil, err
		}
		shares = append(shares, shareMap)
	}

	return shares, nil
}

func dataSourceShareCollectionSharesToMap(meta interface{}, sharesItem vpcv1.Share) (sharesMap map[string]interface{}, diag diag.Diagnostics) {
	sharesMap = map[string]interface{}{}

	if sharesItem.CreatedAt != nil {
		sharesMap["created_at"] = sharesItem.CreatedAt.String()
	}
	if sharesItem.CRN != nil {
		sharesMap["crn"] = sharesItem.CRN
	}
	if sharesItem.Encryption != nil {
		sharesMap["encryption"] = sharesItem.Encryption
	}
	if sharesItem.EncryptionKey != nil && sharesItem.EncryptionKey.CRN != nil {
		sharesMap["encryption_key"] = *sharesItem.EncryptionKey.CRN
	}
	if sharesItem.Href != nil {
		sharesMap["href"] = sharesItem.Href
	}
	if sharesItem.ID != nil {
		sharesMap["id"] = sharesItem.ID
	}
	if sharesItem.Iops != nil {
		sharesMap["iops"] = sharesItem.Iops
	}
	latest_syncs := []map[string]interface{}{}
	if sharesItem.LatestSync != nil {
		latest_sync := make(map[string]interface{})
		latest_sync["completed_at"] = flex.DateTimeToString(sharesItem.LatestSync.CompletedAt)
		if sharesItem.LatestSync.DataTransferred != nil {
			latest_sync["data_transferred"] = *sharesItem.LatestSync.DataTransferred
		}
		latest_sync["started_at"] = flex.DateTimeToString(sharesItem.LatestSync.CompletedAt)
		latest_syncs = append(latest_syncs, latest_sync)
	}
	sharesMap["latest_sync"] = latest_syncs
	if sharesItem.LifecycleState != nil {
		sharesMap["lifecycle_state"] = sharesItem.LifecycleState
	}
	if sharesItem.LatestJob != nil {
		sharesMap["latest_job"] = dataSourceShareFlattenLatestJob(*sharesItem.LatestJob)
	}
	if sharesItem.Name != nil {
		sharesMap["name"] = sharesItem.Name
	}
	if sharesItem.Profile != nil {
		sharesMap["profile"] = *sharesItem.Profile.Name
	}
	if sharesItem.ReplicaShare != nil {
		sharesMap["replica_share"] = dataSourceShareFlattenReplicaShare(*sharesItem.ReplicaShare)
	}
	if sharesItem.ReplicationCronSpec != nil {
		sharesMap["replication_cron_spec"] = *sharesItem.ReplicationCronSpec
	}
	if sharesItem.AccessControlMode != nil {
		sharesMap["access_control_mode"] = *&sharesItem.AccessControlMode
	}
	if !core.IsNil(sharesItem.AllowedTransitEncryptionModes) {
		sharesMap["allowed_transit_encryption_modes"] = sharesItem.AllowedTransitEncryptionModes
	}
	if sharesItem.AccessorBindingRole != nil {
		sharesMap["accessor_binding_role"] = sharesItem.AccessorBindingRole
	}
	accessorBindings := []map[string]interface{}{}
	for _, accessorBindingsItem := range sharesItem.AccessorBindings {
		accessorBindingsItemMap := ResourceIBMIsShareShareAccessorBindingReferenceToMap(&accessorBindingsItem)
		accessorBindings = append(accessorBindings, accessorBindingsItemMap)
	}
	sharesMap["accessor_bindings"] = accessorBindings

	if !core.IsNil(sharesItem.OriginShare) {
		originShareMap := ResourceIBMIsShareShareReferenceToMap(sharesItem.OriginShare)

		sharesMap["origin_share"] = []map[string]interface{}{originShareMap}
	}
	sharesMap["replication_role"] = *sharesItem.ReplicationRole
	sharesMap["replication_status"] = *sharesItem.ReplicationStatus

	if sharesItem.ReplicationStatusReasons != nil {
		sharesMap["replication_status_reasons"] = dataSourceShareFlattenReplicationStatusReasons(sharesItem.ReplicationStatusReasons)
	}
	if sharesItem.ResourceGroup != nil {
		sharesMap["resource_group"] = *sharesItem.ResourceGroup.ID
	}
	if sharesItem.ResourceType != nil {
		sharesMap["resource_type"] = sharesItem.ResourceType
	}
	if sharesItem.Size != nil {
		sharesMap["size"] = sharesItem.Size
	}
	if sharesItem.SourceShare != nil {
		sharesMap["source_share"] = dataSourceShareFlattenSourceShare(*sharesItem.SourceShare)
	}
	if sharesItem.MountTargets != nil {
		targetsList := []map[string]interface{}{}
		for _, targetsItem := range sharesItem.MountTargets {
			targetsList = append(targetsList, dataSourceShareCollectionSharesTargetsToMap(targetsItem))
		}
		sharesMap["mount_targets"] = targetsList
	}
	if sharesItem.Zone != nil {
		sharesMap["zone"] = *sharesItem.Zone.Name
	}
	if sharesItem.SnapshotCount != nil {
		sharesMap["snapshot_count"] = flex.IntValue(sharesItem.SnapshotCount)
	}
	if sharesItem.SnapshotSize != nil {
		sharesMap["snapshot_size"] = flex.IntValue(sharesItem.SnapshotSize)
	}
	sourceSnapshot := []map[string]interface{}{}
	if sharesItem.SourceSnapshot != nil {
		modelMap, err := DataSourceIBMIsShareShareSourceSnapshotToMap(sharesItem.SourceSnapshot)
		if err != nil {
			return nil, flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_share", "read", "source_snapshot-to-map").GetDiag()
		}
		sourceSnapshot = append(sourceSnapshot, modelMap)
	}
	sharesMap["source_snapshot"] = sourceSnapshot
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *sharesItem.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error gettings shares (%s) access tags: %s", *sharesItem.ID, err)
	}

	if sharesItem.UserTags != nil {
		sharesMap[isFileShareTags] = sharesItem.UserTags
	}
	sharesMap[isFileShareAccessTags] = accesstags
	return sharesMap, nil
}

func dataSourceShareCollectionSharesTargetsToMap(targetsItem vpcv1.ShareMountTargetReference) (targetsMap map[string]interface{}) {
	targetsMap = map[string]interface{}{}

	if targetsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceShareTargetsDeletedToMap(*targetsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		targetsMap["deleted"] = deletedList
	}
	if targetsItem.Href != nil {
		targetsMap["href"] = targetsItem.Href
	}
	if targetsItem.ID != nil {
		targetsMap["id"] = targetsItem.ID
	}
	if targetsItem.Name != nil {
		targetsMap["name"] = targetsItem.Name
	}
	if targetsItem.ResourceType != nil {
		targetsMap["resource_type"] = targetsItem.ResourceType
	}

	return targetsMap
}
