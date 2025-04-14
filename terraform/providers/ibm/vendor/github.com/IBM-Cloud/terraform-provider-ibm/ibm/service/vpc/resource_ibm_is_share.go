// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-beta-go-sdk/vpcbetav1"
)

const (
	isFileShareAccessTags             = "access_tags"
	isFileShareTags                   = "tags"
	IsFileShareReplicationRoleNone    = "none"
	IsFileShareReplicationRoleSource  = "source"
	IsFileShareReplicationRoleReplica = "replica"
)

func ResourceIbmIsShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsShareCreate,
		ReadContext:   resourceIbmIsShareRead,
		UpdateContext: resourceIbmIsShareUpdate,
		DeleteContext: resourceIbmIsShareDelete,
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
				},
			),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceSharesValidate(diff)
				}),
		),

		Schema: map[string]*schema.Schema{
			"encryption_key": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"size"},
				ForceNew:     true,
				Computed:     true,
				Description:  "The CRN of the key to use for encrypting this file share.If no encryption key is provided, the share will not be encrypted.",
			},
			"initial_owner": {
				Type:         schema.TypeList,
				Optional:     true,
				MinItems:     1,
				MaxItems:     1,
				RequiredWith: []string{"size"},
				ForceNew:     true,
				Description:  "The owner assigned to the file share at creation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gid": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "The initial group identifier for the file share.",
						},
						"uid": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "The initial user identifier for the file share.",
						},
					},
				},
			},
			"resource_group": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"size"},
				ForceNew:     true,
				Computed:     true,
				Description:  "The unique identifier of the resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.",
			},
			"access_control_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"size"},
				Computed:     true,
				Description:  "The access control mode for the share:",
			},
			"size": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ExactlyOneOf:  []string{"size", "source_share"},
				ConflictsWith: []string{"replication_cron_spec", "source_share"},
				ValidateFunc:  validate.InvokeValidator("ibm_is_share", "size"),
				Description:   "The size of the file share rounded up to the next gigabyte.",
			},
			"mount_targets": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The share targets for this file share.Share targets mounted from a replica must be mounted read-only.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of this mount target",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Href of this mount target",
						},
						"transit_encryption": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The transit encryption mode.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The user-defined name for this share target. Names must be unique within the share the share target resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
						},
						"virtual_network_interface": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "VNI for mount target.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CRN of virtual network interface",
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "href of virtual network interface",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of this VNI",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Name of this VNI",
									},
									"primary_ip": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "VNI for mount target.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"reserved_ip": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "ID of reserved IP",
												},
												"address": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The IP address to reserve, which must not already be reserved on the subnet.",
												},
												"auto_delete": {
													Type:        schema.TypeBool,
													Optional:    true,
													Computed:    true,
													Description: "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Name for reserved IP",
												},
												"resource_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Resource type of primary ip",
												},
												"href": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "href of primary ip",
												},
											},
										},
									},
									"resource_group": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Resource group id",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource type of primary ip",
									},
									"security_groups": {
										Type:        schema.TypeSet,
										Computed:    true,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
										Description: "The security groups to use for this virtual network interface.",
									},
									"subnet": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The associated subnet. Required if primary_ip is not specified.",
									},
								},
							},
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type of mount target",
						},
						"vpc": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The unique identifier of the VPC in which instances can mount the file share using this share target.This property will be removed in a future release.The `subnet` property should be used instead.",
						},
					},
				},
			},
			"iops": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_share", "iops"),
				Description:  "The maximum input/output operation performance bandwidth per second for the file share.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_share", "name"),
				Description:  "The unique user-defined name for this file share. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"profile": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The globally unique name for this share profile.",
			},
			"replica_share": &schema.Schema{
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"source_share"},
				Description:   "Configuration for a replica file share to create and associate with this file share. Ifunspecified, a replica may be subsequently added by creating a new file share with a`source_share` referencing this file share.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this replica share.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The href for this replica share.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of this replica file share.",
						},
						"iops": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_share", "iops"),
							Description:  "The maximum input/output operation per second (IOPS) for the file share.",
						},
						"name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_share", "name"),
							Description:  "The unique user-defined name for this file share.",
						},
						"profile": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Share profile name.",
						},
						"replication_cron_spec": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The cron specification for the file share replication schedule.Replication of a share can be scheduled to occur at most once per hour.",
						},
						"replication_role": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The replication role of the file share.",
						},
						"replication_status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The replication status of the file share.",
						},
						"replication_status_reasons": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current replication status.",
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
						"mount_targets": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "The share targets for this replica file share.Share targets mounted from a replica must be mounted read-only.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "href of mount target",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of this share target.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The user-defined name for this share target. Names must be unique within the share the share target resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
									},
									"virtual_network_interface": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "VNI for mount target.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "CRN of virtual network interface",
												},
												"href": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "href of virtual network interface",
												},
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ID of this VNI",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Name of this VNI",
												},
												"primary_ip": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "VNI for mount target.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"reserved_ip": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "ID of reserved IP",
															},
															"address": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "The IP address to reserve, which must not already be reserved on the subnet.",
															},
															"auto_delete": {
																Type:        schema.TypeBool,
																Optional:    true,
																Computed:    true,
																Description: "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Name for reserved IP",
															},
															"resource_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Resource type of primary ip",
															},
															"href": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "href of primary ip",
															},
														},
													},
												},
												"resource_group": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Resource group id",
												},
												"resource_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Resource type of primary ip",
												},
												"security_groups": {
													Type:        schema.TypeSet,
													Computed:    true,
													Optional:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Set:         schema.HashString,
													Description: "The security groups to use for this virtual network interface.",
												},
												"subnet": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The associated subnet. Required if primary_ip is not specified.",
												},
											},
										},
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource type of virtual network interface",
									},
									"vpc": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The ID of the VPC in which instances can mount the file share using this share target.This property will be removed in a future release.The `subnet` property should be used instead.",
									},
								},
							},
						},
						isFileShareTags: {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_share", isFileShareTags)},
							Set:         flex.ResourceIBMVPCHash,
							Description: "User Tags for the replica share",
						},
						isFileShareAccessTags: {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_share", "accesstag")},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access management tags for this replica share",
						},
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The name of the zone this replica file share will reside in. Must be a different zone in the same region as the source share.",
						},
					},
				},
			},
			"source_share": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"replica_share", "size"},
				RequiredWith:  []string{"replication_cron_spec"},
				Description:   "The ID of the source file share for this replica file share. The specified file share must not already have a replica, and must not be a replica.",
			},
			"replication_cron_spec": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressCronSpecDiff,
				Computed:         true,
				RequiredWith:     []string{"source_share"},
				ConflictsWith:    []string{"replica_share", "size"},
				Description:      "The cron specification for the file share replication schedule.Replication of a share can be scheduled to occur at most once per hour.",
			},
			"replication_role": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The replication role of the file share.* `none`: This share is not participating in replication.* `replica`: This share is a replication target.* `source`: This share is a replication source.",
			},
			"replication_status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The replication status of the file share.* `initializing`: This share is initializing replication.* `active`: This share is actively participating in replication.* `failover_pending`: This share is performing a replication failover.* `split_pending`: This share is performing a replication split.* `none`: This share is not participating in replication.* `degraded`: This share's replication sync is degraded.* `sync_pending`: This share is performing a replication sync.",
			},
			"replication_status_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current replication status (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "A snake case string succinctly identifying the status reason.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "An explanation of the status reason.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Link to documentation about this status reason.",
						},
					},
				},
			},
			"last_sync_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the file share was last synchronized to its replica.This property will be present when the `replication_role` is `source`.",
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
							Description: "The type of the file share job.The enumerated values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the file share job on which the unexpected property value was encountered.* `replication_failover`: This is a share replication failover job.* `replication_init`: This is a share replication is initialization job.* `replication_split`: This is a share replication split job.* `replication_sync`: This is a share replication synchronization job.",
						},
					},
				},
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the file share.",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The globally unique name of the zone this file share will reside in.",
			},
			isFileShareTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_share", isFileShareTags)},
				Set:         flex.ResourceIBMVPCHash,
				Description: "User Tags for the file share",
			},
			isFileShareAccessTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_share", "accesstag")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
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
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this share.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
		},
	}
}

func ResourceIbmIsShareValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 1)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "iops",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "100",
			MaxValue:                   "48000",
		},
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
			Identifier:                 "size",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "10",
			MaxValue:                   "32000",
		},
		validate.ValidateSchema{
			Identifier:                 isFileShareTags,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "accesstag",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-]):([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-])$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_share", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsShareCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	createShareOptions := &vpcbetav1.CreateShareOptions{}

	sharePrototype := &vpcbetav1.SharePrototype{}
	if accessControlModeIntf, ok := d.GetOk("access_control_mode"); ok {
		accessControlMode := accessControlModeIntf.(string)
		sharePrototype.AccessControlMode = &accessControlMode
	}
	if sizeIntf, ok := d.GetOk("size"); ok {

		size := int64(sizeIntf.(int))
		sharePrototype.Size = &size
		if encryptionKeyIntf, ok := d.GetOk("encryption_key"); ok {
			encryptionKey := encryptionKeyIntf.(string)
			encryptionKeyIdentity := &vpcbetav1.EncryptionKeyIdentity{
				CRN: &encryptionKey,
			}
			sharePrototype.EncryptionKey = encryptionKeyIdentity
		}
		initial_owner := &vpcbetav1.ShareInitialOwner{}
		if initialOwnerIntf, ok := d.GetOk("initial_owner"); ok {
			initialOwnerMap := initialOwnerIntf.([]interface{})[0].(map[string]interface{})
			if initialOwnerGIDIntf, ok := initialOwnerMap["gid"]; ok {
				initialOwnerGID := int64(initialOwnerGIDIntf.(int))
				initial_owner.Gid = &initialOwnerGID
			}
			if initialOwnerUIDIntf, ok := initialOwnerMap["uid"]; ok {
				initialOwnerUID := int64(initialOwnerUIDIntf.(int))
				initial_owner.Uid = &initialOwnerUID
			}
		}
		if resgrp, ok := d.GetOk("resource_group"); ok {
			resgrpstr := resgrp.(string)
			resourceGroup := &vpcbetav1.ResourceGroupIdentity{
				ID: &resgrpstr,
			}
			sharePrototype.ResourceGroup = resourceGroup
		}
		if replicaShareIntf, ok := d.GetOk("replica_share"); ok {
			replicaShareMap := replicaShareIntf.([]interface{})[0].(map[string]interface{})
			replicaShare := &vpcbetav1.SharePrototypeShareContext{}
			iopsIntf, ok := replicaShareMap["iops"]
			iops := iopsIntf.(int)
			if ok && iops != 0 {
				replicaShare.Iops = core.Int64Ptr(int64(iops))
			}
			if replicaShareMap["name"] != nil {
				replicaShare.Name = core.StringPtr(replicaShareMap["name"].(string))
			}
			if replicaShareMap["profile"] != nil {
				replicaShare.Profile = &vpcbetav1.ShareProfileIdentity{
					Name: core.StringPtr(replicaShareMap["profile"].(string)),
				}
			}
			if replicaShareMap["replication_cron_spec"] != nil {
				replicaShare.ReplicationCronSpec = core.StringPtr(replicaShareMap["replication_cron_spec"].(string))
			}
			if replicaShareMap["zone"] != nil {
				replicaShare.Zone = &vpcbetav1.ZoneIdentity{
					Name: core.StringPtr(replicaShareMap["zone"].(string)),
				}
			}

			replicaTargets, ok := replicaShareMap["mount_targets"]
			if ok {
				var targets []vpcbetav1.ShareMountTargetPrototypeIntf
				targetsIntf := replicaTargets.([]interface{})
				for _, targetIntf := range targetsIntf {
					target := targetIntf.(map[string]interface{})
					targetsItem, err := resourceIbmIsShareMapToShareMountTargetPrototype(d, target)
					if err != nil {
						return diag.FromErr(err)
					}
					targets = append(targets, &targetsItem)
				}
				replicaShare.MountTargets = targets
			}

			var userTags *schema.Set
			if v, ok := replicaShareMap[isFileShareTags]; ok {
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
					replicaShare.UserTags = userTagsArray
				}
			}
			sharePrototype.ReplicaShare = replicaShare
		}
	} else {
		sourceShare := d.Get("source_share").(string)
		if sourceShare != "" {
			sharePrototype.SourceShare = &vpcbetav1.ShareIdentity{
				ID: &sourceShare,
			}
		}
		replicationCronSpec := d.Get("replication_cron_spec").(string)
		sharePrototype.ReplicationCronSpec = &replicationCronSpec
	}

	if iopsIntf, ok := d.GetOk("iops"); ok {
		iops := int64(iopsIntf.(int))
		sharePrototype.Iops = &iops
	}
	if nameIntf, ok := d.GetOk("name"); ok {
		name := nameIntf.(string)
		sharePrototype.Name = &name
	}
	if profileIntf, ok := d.GetOk("profile"); ok {
		profileStr := profileIntf.(string)
		profile := &vpcbetav1.ShareProfileIdentity{
			Name: &profileStr,
		}
		sharePrototype.Profile = profile
	}

	if shareTargetPrototypeIntf, ok := d.GetOk("mount_targets"); ok {
		var targets []vpcbetav1.ShareMountTargetPrototypeIntf
		for _, e := range shareTargetPrototypeIntf.([]interface{}) {
			value := e.(map[string]interface{})
			targetsItem, err := resourceIbmIsShareMapToShareMountTargetPrototype(d, value)
			if err != nil {
				return diag.FromErr(err)
			}
			targets = append(targets, &targetsItem)
		}
		sharePrototype.MountTargets = targets
	}
	if zone, ok := d.GetOk("zone"); ok {
		zonestr := zone.(string)
		zone := &vpcbetav1.ZoneIdentity{
			Name: &zonestr,
		}
		sharePrototype.Zone = zone
	}
	var userTags *schema.Set
	if v, ok := d.GetOk(isFileShareTags); ok {
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
			sharePrototype.UserTags = userTagsArray
		}
	}
	createShareOptions.SetSharePrototype(sharePrototype)
	share, response, err := vpcClient.CreateShareWithContext(context, createShareOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateShareWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	_, err = isWaitForShareAvailable(context, vpcClient, *share.ID, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if share.ReplicaShare != nil && share.ReplicaShare.ID != nil {
		_, err = isWaitForShareAvailable(context, vpcClient, *share.ReplicaShare.ID, d, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		replicaShareAccessTagsSchema := "replica_share.0.access_tags"
		if _, ok := d.GetOk(replicaShareAccessTagsSchema); ok {
			oldList, newList := d.GetChange(replicaShareAccessTagsSchema)
			err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *share.ReplicaShare.CRN, "", isAccessTagType)
			if err != nil {
				log.Printf(
					"Error creating replica file share (%s) access tags: %s", d.Id(), err)
			}
		}
	}
	d.SetId(*share.ID)

	if _, ok := d.GetOk(isFileShareAccessTags); ok {
		oldList, newList := d.GetChange(isFileShareAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *share.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error creating file share (%s) access tags: %s", d.Id(), err)
		}
	}
	return resourceIbmIsShareRead(context, d, meta)
}

func resourceIbmIsShareMapToShareMountTargetPrototype(d *schema.ResourceData, shareTargetPrototypeMap map[string]interface{}) (vpcbetav1.ShareMountTargetPrototype, error) {
	shareTargetPrototype := vpcbetav1.ShareMountTargetPrototype{}

	if nameIntf, ok := shareTargetPrototypeMap["name"]; ok && nameIntf != "" {
		shareTargetPrototype.Name = core.StringPtr(nameIntf.(string))
	}

	if vpcIntf, ok := shareTargetPrototypeMap["vpc"]; ok && vpcIntf != "" {
		vpc := vpcIntf.(string)
		shareTargetPrototype.VPC = &vpcbetav1.VPCIdentity{
			ID: &vpc,
		}
	} else if vniIntf, ok := shareTargetPrototypeMap["virtual_network_interface"]; ok {
		vniPrototype := vpcbetav1.ShareMountTargetVirtualNetworkInterfacePrototype{}
		vniMap := vniIntf.([]interface{})[0].(map[string]interface{})
		vniPrototype, err := ShareMountTargetMapToShareMountTargetPrototype(d, vniMap)
		if err != nil {
			return shareTargetPrototype, err
		}
		shareTargetPrototype.VirtualNetworkInterface = &vniPrototype
	}
	if transitEncryptionIntf, ok := shareTargetPrototypeMap["transit_encryption"]; ok && transitEncryptionIntf != "" {
		transitEncryption := transitEncryptionIntf.(string)
		shareTargetPrototype.TransitEncryption = &transitEncryption
	}
	return shareTargetPrototype, nil
}

func resourceIbmIsShareRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	getShareOptions := &vpcbetav1.GetShareOptions{}

	getShareOptions.SetID(d.Id())

	share, response, err := vpcClient.GetShareWithContext(context, getShareOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetShareWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	if share.EncryptionKey != nil {
		if err = d.Set("encryption_key", *share.EncryptionKey.CRN); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting encryption_key: %s", err))
		}
	}
	if share.AccessControlMode != nil {
		if err = d.Set("access_control_mode", *share.AccessControlMode); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting access_control_mode: %s", err))
		}
	}
	if err = d.Set("iops", flex.IntValue(share.Iops)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting iops: %s", err))
	}
	if err = d.Set("name", share.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if share.Profile != nil {
		if err = d.Set("profile", *share.Profile.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting profile: %s", err))
		}
	}
	if share.ResourceGroup != nil {
		if err = d.Set("resource_group", *share.ResourceGroup.ID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group: %s", err))
		}
	}
	if err = d.Set("size", flex.IntValue(share.Size)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting size: %s", err))
	}
	targets := make([]map[string]interface{}, 0)
	if share.MountTargets != nil {
		for _, targetsItem := range share.MountTargets {
			GetShareMountTargetOptions := &vpcbetav1.GetShareMountTargetOptions{}
			GetShareMountTargetOptions.SetShareID(d.Id())
			GetShareMountTargetOptions.SetID(*targetsItem.ID)

			shareTarget, response, err := vpcClient.GetShareMountTargetWithContext(context, GetShareMountTargetOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					d.SetId("")
					return nil
				}
				log.Printf("[DEBUG] GetShareMountTargetWithContext failed %s\n%s", err, response)
				return diag.FromErr(err)
			}

			targetsItemMap, err := ShareMountTargetToMap(context, vpcClient, d, *shareTarget)
			if err != nil {
				return diag.FromErr(err)
			}
			targets = append(targets, targetsItemMap)
		}
	}
	if err = d.Set("mount_targets", targets); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting mount_targets: %s", err))
	}

	replicaShare := []map[string]interface{}{}
	if share.ReplicaShare != nil && share.ReplicaShare.ID != nil {
		if _, ok := d.GetOk("replica_share"); ok {
			getShareOptions := &vpcbetav1.GetShareOptions{}

			getShareOptions.SetID(*share.ReplicaShare.ID)

			share, response, err := vpcClient.GetShareWithContext(context, getShareOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					d.SetId("")
					return nil
				}
				log.Printf("[DEBUG] GetShareWithContext failed %s\n%s", err, response)
				return diag.FromErr(err)
			}
			replicaShareItem, err := ShareReplicaToMap(context, vpcClient, d, meta, *share)
			if err != nil {
				return diag.FromErr(err)
			}
			replicaShare = append(replicaShare, replicaShareItem)
			d.Set("replica_share", replicaShare)
		}
	}

	if share.Zone != nil {
		if err = d.Set("zone", *share.Zone.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting zone: %s", err))
		}
	}
	if err = d.Set("created_at", share.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("crn", share.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("encryption", share.Encryption); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting encryption: %s", err))
	}
	if err = d.Set("href", share.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", share.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", share.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	// if share.LastSyncAt != nil {
	// 	d.Set("last_sync_at", share.LastSyncAt.String())
	// }
	latest_jobs := []map[string]interface{}{}
	if share.LatestJob != nil {
		latest_job := make(map[string]interface{})
		latest_job["status"] = share.LatestJob.Status
		status_reasons := []map[string]interface{}{}
		for _, status_reason_item := range share.LatestJob.StatusReasons {
			status_reason := make(map[string]interface{})
			status_reason["code"] = status_reason_item.Code
			status_reason["message"] = status_reason_item.Message
			if status_reason_item.MoreInfo != nil {
				status_reason["more_info"] = status_reason_item.MoreInfo
			}
			status_reasons = append(status_reasons, status_reason)
		}
		latest_job["status_reasons"] = status_reasons
		latest_job["type"] = share.LatestJob.Type
		latest_jobs = append(latest_jobs, latest_job)
	}
	d.Set("latest_job", latest_jobs)

	if share.ReplicationCronSpec != nil {
		d.Set("replication_cron_spec", share.ReplicationCronSpec)
	}
	d.Set("replication_role", share.ReplicationRole)
	d.Set("replication_status", share.ReplicationStatus)
	status_reasons := []map[string]interface{}{}
	for _, status_reason_item := range share.ReplicationStatusReasons {
		status_reason := make(map[string]interface{})
		status_reason["code"] = status_reason_item.Code
		status_reason["message"] = status_reason_item.Message
		if status_reason_item.MoreInfo != nil {
			status_reason["more_info"] = status_reason_item.MoreInfo
		}
		status_reasons = append(status_reasons, status_reason)
	}
	d.Set("replication_status_reasons", status_reasons)

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *share.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error getting shares (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isFileShareAccessTags, accesstags)
	// d.Set(isFileShareTags, tags)
	if share.UserTags != nil {
		if err = d.Set(isFileShareTags, share.UserTags); err != nil {
			log.Printf(
				"Error setting shares (%s) user tags: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceIbmIsShareUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	getShareOptions := &vpcbetav1.GetShareOptions{}

	getShareOptions.SetID(d.Id())

	share, response, err := vpcClient.GetShareWithContext(context, getShareOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetShareWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	eTag := response.Headers.Get("ETag")

	err = shareUpdate(vpcClient, context, d, meta, "share", d.Id(), eTag)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("replica_share") {
		if share.ReplicaShare != nil && share.ReplicaShare.ID != nil {
			getShareOptions := &vpcbetav1.GetShareOptions{}
			getShareOptions.SetID(*share.ReplicaShare.ID)

			replicaShare, response, err := vpcClient.GetShareWithContext(context, getShareOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					d.SetId("")
					return nil
				}
				log.Printf("[DEBUG] GetShareWithContext failed %s\n%s", err, response)
				return diag.FromErr(err)
			}
			eTag := response.Headers.Get("ETag")
			err = shareUpdate(vpcClient, context, d, meta, "replica_share", *replicaShare.ID, eTag)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange(isFileShareAccessTags) {
		oldList, newList := d.GetChange(isFileShareAccessTags)
		err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, d.Get("crn").(string), "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error updating shares (%s) access tags: %s", d.Id(), err)
		}
	}

	return resourceIbmIsShareRead(context, d, meta)
}

func resourceIbmIsShareDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	getShareOptions := &vpcbetav1.GetShareOptions{}

	getShareOptions.SetID(d.Id())

	share, response, err := vpcClient.GetShareWithContext(context, getShareOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetShareWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	if share.MountTargets != nil {
		if _, ok := d.GetOk("mount_targets"); ok {
			for _, targetsItem := range share.MountTargets {

				deleteShareMountTargetOptions := &vpcbetav1.DeleteShareMountTargetOptions{}

				deleteShareMountTargetOptions.SetShareID(d.Id())
				deleteShareMountTargetOptions.SetID(*targetsItem.ID)

				_, response, err := vpcClient.DeleteShareMountTargetWithContext(context, deleteShareMountTargetOptions)
				if err != nil {
					log.Printf("[DEBUG] DeleteShareMountTargetWithContext failed %s\n%s", err, response)
					return diag.FromErr(err)
				}
				_, err = isWaitForTargetDelete(context, vpcClient, d, d.Id(), *targetsItem.ID)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if share.ReplicationRole != nil && *share.ReplicationRole == IsFileShareReplicationRoleSource && share.ReplicaShare != nil {

		getShareOptions := &vpcbetav1.GetShareOptions{}
		getShareOptions.SetID(*share.ReplicaShare.ID)

		if _, ok := d.GetOk("replica_share"); ok {
			replicaShare, response, err := vpcClient.GetShareWithContext(context, getShareOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					d.SetId("")
					return nil
				}
				log.Printf("[DEBUG] GetShareWithContext failed %s\n%s", err, response)
				return diag.FromErr(err)
			}
			if replicaShare.MountTargets != nil {
				if _, ok := d.GetOk("replica_share.0.mount_targets"); ok {
					for _, targetsItem := range replicaShare.MountTargets {

						deleteShareMountTargetOptions := &vpcbetav1.DeleteShareMountTargetOptions{}

						deleteShareMountTargetOptions.SetShareID(*replicaShare.ID)
						deleteShareMountTargetOptions.SetID(*targetsItem.ID)

						_, response, err := vpcClient.DeleteShareMountTargetWithContext(context, deleteShareMountTargetOptions)
						if err != nil {
							log.Printf("[DEBUG] DeleteShareMountTargetWithContext failed %s\n%s", err, response)
							return diag.FromErr(err)
						}
						_, err = isWaitForTargetDelete(context, vpcClient, d, d.Id(), *targetsItem.ID)
						if err != nil {
							return diag.FromErr(err)
						}
					}
				}
			}
			replicaShare, response, err = vpcClient.GetShareWithContext(context, getShareOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					d.SetId("")
					return nil
				}
				log.Printf("[DEBUG] GetShareWithContext failed %s\n%s", err, response)
				return diag.FromErr(err)
			}
			replicaETag := response.Headers.Get("ETag")
			deleteShareOptions := &vpcbetav1.DeleteShareOptions{}
			deleteShareOptions.IfMatch = &replicaETag
			deleteShareOptions.SetID(*replicaShare.ID)
			_, response, err = vpcClient.DeleteShareWithContext(context, deleteShareOptions)
			if err != nil {
				log.Printf("[DEBUG] DeleteShareWithContext failed %s\n%s", err, response)
				return diag.FromErr(err)
			}

			_, err = isWaitForShareDelete(context, vpcClient, d, *replicaShare.ID)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	share, response, err = vpcClient.GetShareWithContext(context, getShareOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetShareWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	ETag := response.Headers.Get("ETag")
	deleteShareOptions := &vpcbetav1.DeleteShareOptions{}

	deleteShareOptions.SetID(d.Id())
	deleteShareOptions.IfMatch = &ETag
	_, response, err = vpcClient.DeleteShareWithContext(context, deleteShareOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteShareWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	_, err = isWaitForShareDelete(context, vpcClient, d, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func isWaitForShareAvailable(context context.Context, vpcClient *vpcbetav1.VpcbetaV1, shareid string, d *schema.ResourceData, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for share (%s) to be available.", shareid)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating", "pending", "waiting"},
		Target:     []string{"stable", "failed"},
		Refresh:    isShareRefreshFunc(context, vpcClient, shareid, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isShareRefreshFunc(context context.Context, vpcClient *vpcbetav1.VpcbetaV1, shareid string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		shareOptions := &vpcbetav1.GetShareOptions{}

		shareOptions.SetID(shareid)

		share, response, err := vpcClient.GetShareWithContext(context, shareOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting share: %s\n%s", err, response)
		}
		d.Set("lifecycle_state", *share.LifecycleState)
		if *share.LifecycleState == "stable" || *share.LifecycleState == "failed" {

			return share, *share.LifecycleState, nil

		}
		return share, "pending", nil
	}
}

func isWaitForShareDelete(context context.Context, vpcClient *vpcbetav1.VpcbetaV1, d *schema.ResourceData, shareid string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{"deleting", "stable", "waiting"},
		Target:  []string{"done"},
		Refresh: func() (interface{}, string, error) {
			shareOptions := &vpcbetav1.GetShareOptions{}

			shareOptions.SetID(shareid)

			share, response, err := vpcClient.GetShareWithContext(context, shareOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return share, "done", nil
				}
				return nil, "", fmt.Errorf("Error Getting Target: %s\n%s", err, response)
			}
			if *share.LifecycleState == isInstanceFailed {
				return share, *share.LifecycleState, fmt.Errorf("The  target %s failed to delete: %v", shareid, err)
			}
			return share, "deleting", nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func suppressCronSpecDiff(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("latest_job.0.type").(string) == "replication_failover" && d.Get("latest_job.0.status").(string) == "succeeded"
}

func ShareReplicaToMap(context context.Context, vpcClient *vpcbetav1.VpcbetaV1, d *schema.ResourceData, meta interface{}, shareReplica vpcbetav1.Share) (map[string]interface{}, error) {
	shareReplicaMap := map[string]interface{}{}

	shareReplicaMap["crn"] = shareReplica.CRN
	shareReplicaMap["href"] = shareReplica.Href
	shareReplicaMap["name"] = shareReplica.Name
	shareReplicaMap["id"] = shareReplica.ID
	shareReplicaMap["iops"] = shareReplica.Iops
	shareReplicaMap["replication_cron_spec"] = shareReplica.ReplicationCronSpec
	shareReplicaMap["replication_role"] = shareReplica.ReplicationRole
	shareReplicaMap["profile"] = shareReplica.Profile.Name
	shareReplicaMap["replication_status"] = shareReplica.ReplicationStatus
	shareReplicaMap["zone"] = shareReplica.Zone.Name
	status_reasons := []map[string]interface{}{}
	for _, status_reason_item := range shareReplica.ReplicationStatusReasons {
		status_reason := make(map[string]interface{})
		status_reason["code"] = status_reason_item.Code
		status_reason["message"] = status_reason_item.Message
		if status_reason_item.MoreInfo != nil {
			status_reason["more_info"] = status_reason_item.MoreInfo
		}
		status_reasons = append(status_reasons, status_reason)
	}
	shareReplicaMap["replication_status_reasons"] = status_reasons

	targets := []map[string]interface{}{}
	for _, mountTarget := range shareReplica.MountTargets {
		GetShareMountTargetOptions := &vpcbetav1.GetShareMountTargetOptions{}

		GetShareMountTargetOptions.SetShareID(*shareReplica.ID)
		GetShareMountTargetOptions.SetID(*mountTarget.ID)

		shareTarget, response, err := vpcClient.GetShareMountTargetWithContext(context, GetShareMountTargetOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				d.SetId("")
				return nil, err
			}
			log.Printf("[DEBUG] GetShareMountTargetWithContext failed %s\n%s", err, response)
			return nil, err
		}
		targetsItemMap, err := ShareMountTargetToMap(context, vpcClient, d, *shareTarget)
		if err != nil {
			return nil, err
		}
		targets = append(targets, targetsItemMap)
	}

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *shareReplica.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error getting shares (%s) access tags: %s", d.Id(), err)
	}
	shareReplicaMap[isFileShareAccessTags] = accesstags
	// d.Set(isFileShareTags, tags)
	if shareReplica.UserTags != nil {
		shareReplicaMap[isFileShareTags] = shareReplica.UserTags
	}

	shareReplicaMap["mount_targets"] = targets

	return shareReplicaMap, nil
}

func ShareMountTargetToMap(context context.Context, vpcClient *vpcbetav1.VpcbetaV1, d *schema.ResourceData, shareMountTarget vpcbetav1.ShareMountTarget) (map[string]interface{}, error) {
	mountTarget := map[string]interface{}{}

	mountTarget["name"] = *shareMountTarget.Name
	mountTarget["id"] = *shareMountTarget.ID
	mountTarget["href"] = *shareMountTarget.Href
	mountTarget["resource_type"] = *shareMountTarget.ResourceType
	if shareMountTarget.TransitEncryption != nil {
		mountTarget["transit_encryption"] = *shareMountTarget.TransitEncryption
	}
	if shareMountTarget.VirtualNetworkInterface != nil {
		vni, err := ShareMountTargetVirtualNetworkInterfaceToMap(context, vpcClient, d, *shareMountTarget.VirtualNetworkInterface.ID)
		if err != nil {
			return nil, err
		}
		mountTarget["virtual_network_interface"] = vni
	}
	if shareMountTarget.VPC != nil {
		mountTarget["vpc"] = *shareMountTarget.VPC.ID
	}
	return mountTarget, nil
}

func shareUpdate(vpcClient *vpcbetav1.VpcbetaV1, context context.Context, d *schema.ResourceData, meta interface{}, shareType, shareId, eTag string) error {
	updateShareOptions := &vpcbetav1.UpdateShareOptions{}

	updateShareOptions.SetID(shareId)
	updateShareOptions.IfMatch = &eTag

	hasChange := false
	hasSizeChanged := false
	sharePatchModel := &vpcbetav1.SharePatch{}
	shareNameSchema := ""
	shareIopsSchema := ""
	shareProfileSchema := ""
	shareTagsSchema := ""
	shareMountTargetSchema := ""
	accessTagsSchema := ""
	shareCRN := ""
	replicationCronSpec := ""

	if shareType == "share" {
		shareNameSchema = "name"
		shareIopsSchema = "iops"
		shareProfileSchema = "profile"
		shareTagsSchema = "tags"
		shareMountTargetSchema = "mount_targets"
		accessTagsSchema = "access_tags"
		shareCRN = "crn"
		replicationCronSpec = "replication_cron_spec"
	} else {
		shareNameSchema = "replica_share.0.name"
		shareIopsSchema = "replica_share.0.iops"
		shareProfileSchema = "replica_share.0.profile"
		shareTagsSchema = "replica_share.0.tags"
		shareMountTargetSchema = "replica_share.0.mount_targets"
		accessTagsSchema = "replica_share.0.access_tags"
		shareCRN = "replica_share.0.crn"
		replicationCronSpec = "replica_share.0.replication_cron_spec"
	}
	if d.HasChange(shareNameSchema) {
		name := d.Get(shareNameSchema).(string)
		sharePatchModel.Name = &name
		hasChange = true
	}
	if shareType == "share" {
		if d.HasChange("size") {
			size := int64(d.Get("size").(int))
			hasSizeChanged = true
			sharePatchModel.Size = &size
			hasChange = true
		}
	}
	if shareType == "share" {
		if d.HasChange("access_control_mode") {
			accessControlMode := d.Get("access_control_mode").(string)
			if accessControlMode != "" {
				sharePatchModel.AccessControlMode = &accessControlMode
				hasChange = true
			}
		}
	}
	if d.HasChange(replicationCronSpec) {
		replicationCronSpecStr := d.Get(replicationCronSpec).(string)
		if replicationCronSpecStr != "" {
			sharePatchModel.ReplicationCronSpec = &replicationCronSpecStr
			hasChange = true
		}

	}

	if d.HasChange(shareIopsSchema) {
		iops := int64(d.Get(shareIopsSchema).(int))
		sharePatchModel.Iops = &iops
		hasChange = true
	}

	if d.HasChange(shareProfileSchema) {
		_, new := d.GetChange(shareProfileSchema)
		profile := new.(string)
		sharePatchModel.Profile = &vpcbetav1.ShareProfileIdentity{
			Name: &profile,
		}
		hasChange = true
	}

	if d.HasChange(shareTagsSchema) {
		var userTags *schema.Set
		if v, ok := d.GetOk(shareTagsSchema); ok {

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

				sharePatchModel.UserTags = userTagsArray
			}
		}
		hasChange = true
	}
	if hasChange {

		sharePatch, err := sharePatchModel.AsPatch()

		if err != nil {
			log.Printf("[DEBUG] SharePatch AsPatch failed %s", err)
			return err
		}
		updateShareOptions.SetSharePatch(sharePatch)
		if hasSizeChanged {
			_, err = isWaitForShareAvailable(context, vpcClient, d.Id(), d, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return err
			}
		}
		_, response, err := vpcClient.UpdateShareWithContext(context, updateShareOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateShareWithContext failed %s\n%s", err, response)
			return err
		}
		_, err = isWaitForShareAvailable(context, vpcClient, d.Id(), d, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}
	if d.HasChange(shareMountTargetSchema) && !d.IsNewResource() {
		target_prototype := d.Get(shareMountTargetSchema).([]interface{})
		for targetIdx := range target_prototype {
			targetName := fmt.Sprintf("%s.%d.name", shareMountTargetSchema, targetIdx)
			vniName := fmt.Sprintf("%s.%d.virtual_network_interface.0.name", shareMountTargetSchema, targetIdx)
			vniId := fmt.Sprintf("%s.%d.virtual_network_interface.0.id", shareMountTargetSchema, targetIdx)
			targetId := fmt.Sprintf("%s.%d.id", shareMountTargetSchema, targetIdx)
			securityGroups := fmt.Sprintf("%s.%d.virtual_network_interface.0.security_groups", shareMountTargetSchema, targetIdx)
			vniPrimaryIpName := fmt.Sprintf("%s.%d.virtual_network_interface.0.primary_ip.0.name", shareMountTargetSchema, targetIdx)
			vniPrimaryIpAutoDelete := fmt.Sprintf("%s.%d.virtual_network_interface.0.primary_ip.0.auto_delete", shareMountTargetSchema, targetIdx)
			vniSubnet := fmt.Sprintf("%s.%d.virtual_network_interface.0.subnet", shareMountTargetSchema, targetIdx)
			vniResvedIp := fmt.Sprintf("%s.%d.virtual_network_interface.0.primary_ip.0.reserved_ip", shareMountTargetSchema, targetIdx)
			mountTargetId := d.Get(targetId).(string)
			if d.HasChange(targetName) {
				updateShareTargetOptions := &vpcbetav1.UpdateShareMountTargetOptions{}

				updateShareTargetOptions.SetShareID(shareId)
				updateShareTargetOptions.SetID(mountTargetId)

				shareTargetPatchModel := &vpcbetav1.ShareMountTargetPatch{}

				name := d.Get(targetName).(string)
				shareTargetPatchModel.Name = &name

				shareTargetPatch, err := shareTargetPatchModel.AsPatch()
				if err != nil {
					log.Printf("[DEBUG] ShareTargetPatch AsPatch failed %s", err)
					return err
				}
				updateShareTargetOptions.SetShareMountTargetPatch(shareTargetPatch)
				_, response, err := vpcClient.UpdateShareMountTargetWithContext(context, updateShareTargetOptions)
				if err != nil {
					log.Printf("[DEBUG] UpdateShareTargetWithContext failed %s\n%s", err, response)
					return err
				}
				_, err = WaitForTargetAvailable(context, vpcClient, shareId, mountTargetId, d, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			}

			if d.HasChange(vniName) {
				vniNameStr := d.Get(vniName).(string)
				vniPatchModel := &vpcbetav1.VirtualNetworkInterfacePatch{
					Name: &vniNameStr,
				}
				vniPatch, err := vniPatchModel.AsPatch()
				if err != nil {
					log.Printf("[DEBUG] Virtual network interface AsPatch failed %s", err)
					return err
				}
				shareTargetOptions := &vpcbetav1.GetShareMountTargetOptions{}

				shareTargetOptions.SetShareID(shareId)
				shareTargetOptions.SetID(mountTargetId)
				shareTarget, response, err := vpcClient.GetShareMountTargetWithContext(context, shareTargetOptions)
				if err != nil {
					log.Printf("[DEBUG] GetShareMountTargetWithContext failed %s\n%s", err, response)
					return err
				}
				vniId := *shareTarget.VirtualNetworkInterface.ID
				updateVNIOptions := &vpcbetav1.UpdateVirtualNetworkInterfaceOptions{
					ID:                           &vniId,
					VirtualNetworkInterfacePatch: vniPatch,
				}
				_, response, err = vpcClient.UpdateVirtualNetworkInterfaceWithContext(context, updateVNIOptions)
				if err != nil {
					log.Printf("[DEBUG] UpdateShareTargetWithContext failed %s\n%s", err, response)
					return err
				}
				_, err = WaitForVNIAvailable(vpcClient, *shareTarget.VirtualNetworkInterface.ID, d, d.Timeout(schema.TimeoutCreate))
				if err != nil {
					return err
				}
			}

			if d.HasChange(securityGroups) {
				ovs, nvs := d.GetChange(securityGroups)
				ov := ovs.(*schema.Set)
				nv := nvs.(*schema.Set)
				remove := flex.ExpandStringList(ov.Difference(nv).List())
				add := flex.ExpandStringList(nv.Difference(ov).List())
				networkID := d.Get(vniId).(string)
				if len(add) > 0 {

					for i := range add {
						createsgnicoptions := &vpcbetav1.CreateSecurityGroupTargetBindingOptions{
							SecurityGroupID: &add[i],
							ID:              &networkID,
						}
						_, response, err := vpcClient.CreateSecurityGroupTargetBinding(createsgnicoptions)
						if err != nil {
							return fmt.Errorf("[ERROR] Error while creating security group %q for virtual network interface of share mount target %s\n%s: %q", add[i], d.Id(), err, response)
						}
						_, err = WaitForVNIAvailable(vpcClient, networkID, d, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return err
						}

						_, err = WaitForTargetAvailable(context, vpcClient, shareId, mountTargetId, d, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return err
						}
					}

				}
				if len(remove) > 0 {
					for i := range remove {
						deletesgnicoptions := &vpcbetav1.DeleteSecurityGroupTargetBindingOptions{
							SecurityGroupID: &remove[i],
							ID:              &networkID,
						}
						response, err := vpcClient.DeleteSecurityGroupTargetBinding(deletesgnicoptions)
						if err != nil {
							return fmt.Errorf("[ERROR] Error while removing security group %q for virtual network interface of share mount target %s\n%s: %q", remove[i], d.Id(), err, response)
						}
						_, err = WaitForVNIAvailable(vpcClient, networkID, d, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return err
						}

						_, err = WaitForTargetAvailable(context, vpcClient, shareId, mountTargetId, d, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return err
						}
					}
				}
			}

			if !d.IsNewResource() && (d.HasChange(vniPrimaryIpName) || d.HasChange(vniPrimaryIpAutoDelete)) {
				sess, err := meta.(conns.ClientSession).VpcV1API()
				if err != nil {
					return err
				}
				subnetId := d.Get(vniSubnet).(string)
				ripId := d.Get(vniResvedIp).(string)
				updateripoptions := &vpcbetav1.UpdateSubnetReservedIPOptions{
					SubnetID: &subnetId,
					ID:       &ripId,
				}
				reservedIpPath := &vpcbetav1.ReservedIPPatch{}
				if d.HasChange(vniPrimaryIpName) {
					name := d.Get(vniPrimaryIpName).(string)
					reservedIpPath.Name = &name
				}
				if d.HasChange(vniPrimaryIpAutoDelete) {
					auto := d.Get(vniPrimaryIpAutoDelete).(bool)
					reservedIpPath.AutoDelete = &auto
				}
				reservedIpPathAsPatch, err := reservedIpPath.AsPatch()
				if err != nil {
					return fmt.Errorf("[ERROR] Error calling reserved ip as patch \n%s", err)
				}
				updateripoptions.ReservedIPPatch = reservedIpPathAsPatch
				_, response, err := vpcClient.UpdateSubnetReservedIP(updateripoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error updating instance network interface reserved ip(%s): %s\n%s", ripId, err, response)
				}
				_, err = isWaitForReservedIpAvailable(sess, subnetId, ripId, d.Timeout(schema.TimeoutUpdate), d)
				if err != nil {
					return fmt.Errorf("[ERROR] Error waiting for the reserved IP to be available: %s", err)
				}
			}
		}
	}
	if d.HasChange(accessTagsSchema) {
		oldList, newList := d.GetChange(accessTagsSchema)
		err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, d.Get(shareCRN).(string), "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error updating shares (%s) access tags: %s", d.Id(), err)
		}
	}
	return nil
}
