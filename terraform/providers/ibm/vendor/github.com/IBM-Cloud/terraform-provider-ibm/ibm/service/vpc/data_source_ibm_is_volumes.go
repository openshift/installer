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

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const (
	isVolumeArchitecture                                   = "architecture"
	isVolumeDHOnly                                         = "dedicated_host_only"
	isVolumeDisplayName                                    = "display_name"
	isVolumeOSFamily                                       = "family"
	isVolumeOSVendor                                       = "vendor"
	isVolumeOSVersion                                      = "version"
	isVolumeAttachmentState                                = "attachment_state"
	isVolumes                                              = "volumes"
	isVolumesActive                                        = "active"
	isVolumesBandwidth                                     = "bandwidth"
	isVolumesBusy                                          = "busy"
	isVolumesCapacity                                      = "capacity"
	isVolumesCreatedAt                                     = "created_at"
	isVolumesCRN                                           = "crn"
	isVolumesEncryption                                    = "encryption"
	isVolumesEncryptionKey                                 = "encryption_key"
	isVolumesEncryptionKeyCRN                              = "crn"
	isVolumesHref                                          = "href"
	isVolumesId                                            = "id"
	isVolumesIops                                          = "iops"
	isVolumesName                                          = "name"
	isVolumesOperatingSystem                               = "operating_system"
	isVolumesOperatingSystemFamily                         = "operating_system_family"
	isVolumesOperatingSystemArch                           = "operating_system_architecture"
	isVolumesOperatingSystemHref                           = "href"
	isVolumesOperatingSystemName                           = "name"
	isVolumesProfile                                       = "profile"
	isVolumesProfileHref                                   = "href"
	isVolumesProfileName                                   = "name"
	isVolumesResourceGroup                                 = "resource_group"
	isVolumesResourceGroupHref                             = "href"
	isVolumesResourceGroupId                               = "id"
	isVolumesResourceGroupName                             = "name"
	isVolumesSourceImage                                   = "source_image"
	isVolumesSourceImageCRN                                = "crn"
	isVolumesSourceImageDeleted                            = "deleted"
	isVolumesSourceImageDeletedMoreInfo                    = "more_info"
	isVolumesSourceImageHref                               = "href"
	isVolumesSourceImageId                                 = "id"
	isVolumesSourceImageName                               = "name"
	isVolumesSourceSnapshot                                = "source_snapshot"
	isVolumesSourceSnapshotCRN                             = "crn"
	isVolumesSourceSnapshotDeleted                         = "deleted"
	isVolumesSourceSnapshotDeletedMoreInfo                 = "more_info"
	isVolumesSourceSnapshotHref                            = "href"
	isVolumesSourceSnapshotId                              = "id"
	isVolumesSourceSnapshotName                            = "name"
	isVolumesSourceSnapshotResourceType                    = "resource_type"
	isVolumesStatus                                        = "status"
	isVolumesStatusReasons                                 = "status_reasons"
	isVolumesStatusReasonsCode                             = "code"
	isVolumesStatusReasonsMessage                          = "message"
	isVolumesStatusReasonsMoreInfo                         = "more_info"
	isVolumesVolumeAttachments                             = "volume_attachments"
	isVolumesVolumeAttachmentsDeleteVolumeOnInstanceDelete = "delete_volume_on_instance_delete"
	isVolumesVolumeAttachmentsDeleted                      = "deleted"
	isVolumesVolumeAttachmentsDeletedMoreInfo              = "more_info"
	isVolumesVolumeAttachmentsDevice                       = "device"
	isVolumesVolumeAttachmentsDeviceId                     = "id"
	isVolumesVolumeAttachmentsHref                         = "href"
	isVolumesVolumeAttachmentsId                           = "id"
	isVolumesVolumeAttachmentsInstance                     = "instance"
	isVolumesVolumeAttachmentsInstanceCRN                  = "crn"
	isVolumesVolumeAttachmentsInstanceDeleted              = "deleted"
	isVolumesVolumeAttachmentsInstanceDeletedMoreInfo      = "more_info"
	isVolumesVolumeAttachmentsInstanceHref                 = "href"
	isVolumesVolumeAttachmentsInstanceId                   = "id"
	isVolumesVolumeAttachmentsInstanceName                 = "name"
	isVolumesVolumeAttachmentsName                         = "name"
	isVolumesVolumeAttachmentsType                         = "type"
	isVolumesZone                                          = "zone"
	isVolumesZoneHref                                      = "href"
	isVolumesZoneName                                      = "name"
)

func DataSourceIBMIsVolumes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVolumesRead,

		Schema: map[string]*schema.Schema{
			isVolumeAttachmentState: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Attachment state of the Volume.",
			},
			isVolumesEncryption: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Encryption type of Volume.",
			},
			isVolumesOperatingSystemFamily: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Operating system family of the Volume.",
			},
			isVolumesOperatingSystemArch: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Operating system architecture of the Volume.",
			},
			"volume_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Volume name identifier.",
			},
			"zone_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone name identifier.",
			},
			isVolumes: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of volumes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVolumesActive: &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether a running virtual server instance has an attachment to this volume.",
						},
						isVolumeAttachmentState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The attachment state of the volume.",
						},
						isVolumesBandwidth: &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum bandwidth (in megabits per second) for the volume.",
						},
						isVolumesBusy: &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this volume is performing an operation that must be serialized. If an operation specifies that it requires serialization, the operation will fail unless this property is `false`.",
						},
						isVolumesCapacity: &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The capacity to use for the volume (in gigabytes). The specified minimum and maximum capacity values for creating or updating volumes may expand in the future.",
						},
						isVolumesCreatedAt: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the volume was created.",
						},
						isVolumesCRN: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this volume.",
						},
						isVolumesEncryption: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of encryption used on the volume.",
						},
						isVolumesEncryptionKey: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The root key used to wrap the data encryption key for the volume.This property will be present for volumes with an `encryption` type of`user_managed`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVolumesEncryptionKeyCRN: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.",
									},
								},
							},
						},
						isVolumesHref: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this volume.",
						},
						isVolumesId: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this volume.",
						},
						isVolumesIops: &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum I/O operations per second (IOPS) to use for the volume. Applicable only to volumes using a profile `family` of `custom`.",
						},
						isVolumesName: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this volume.",
						},
						isVolumesOperatingSystem: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The operating system associated with this volume. If absent, this volume was notcreated from an image, or the image did not include an operating system.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVolumeArchitecture: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The operating system architecture.",
									},
									isVolumeDHOnly: &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Images with this operating system can only be used on dedicated hosts or dedicated host groups.",
									},
									isVolumeDisplayName: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unique, display-friendly name for the operating system.",
									},
									isVolumeOSFamily: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The software family for this operating system.",
									},

									isVolumesOperatingSystemHref: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this operating system.",
									},
									isVolumesOperatingSystemName: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this operating system.",
									},
									isVolumeOSVendor: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The vendor of the operating system.",
									},
									isVolumeOSVersion: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The major release version of this operating system.",
									},
								},
							},
						},
						isVolumesProfile: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The profile this volume uses.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVolumesProfileHref: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this volume profile.",
									},
									isVolumesProfileName: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this volume profile.",
									},
								},
							},
						},
						isVolumesResourceGroup: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this volume.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVolumesResourceGroupHref: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this resource group.",
									},
									isVolumesResourceGroupId: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this resource group.",
									},
									isVolumesResourceGroupName: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this resource group.",
									},
								},
							},
						},
						isVolumesSourceImage: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The image from which this volume was created (this may be[deleted](https://cloud.ibm.com/apidocs/vpc#deleted-resources)).If absent, this volume was not created from an image.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVolumesSourceImageCRN: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this image.",
									},
									isVolumesSourceImageDeleted: &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isVolumesSourceImageDeletedMoreInfo: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									isVolumesSourceImageHref: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this image.",
									},
									isVolumesSourceImageId: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this image.",
									},
									isVolumesSourceImageName: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined or system-provided name for this image.",
									},
								},
							},
						},
						isVolumesSourceSnapshot: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The snapshot from which this volume was cloned.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVolumesSourceSnapshotCRN: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this snapshot.",
									},
									isVolumesSourceSnapshotDeleted: &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isVolumesSourceSnapshotDeletedMoreInfo: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									isVolumesSourceSnapshotHref: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this snapshot.",
									},
									isVolumesSourceSnapshotId: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this snapshot.",
									},
									isVolumesSourceSnapshotName: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this snapshot.",
									},
									isVolumesSourceSnapshotResourceType: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						isVolumesStatus: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the volume.The enumerated values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the volume on which the unexpected property value was encountered.",
						},
						isVolumesStatusReasons: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current status (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVolumesStatusReasonsCode: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the status reason.",
									},
									isVolumesStatusReasonsMessage: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the status reason.",
									},
									isVolumesStatusReasonsMoreInfo: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about this status reason.",
									},
								},
							},
						},
						isVolumesVolumeAttachments: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The volume attachments for this volume.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVolumesVolumeAttachmentsDeleteVolumeOnInstanceDelete: &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If set to true, when deleting the instance the volume will also be deleted.",
									},
									isVolumesVolumeAttachmentsDeleted: &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isVolumesVolumeAttachmentsDeletedMoreInfo: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									isVolumesVolumeAttachmentsDevice: &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Information about how the volume is exposed to the instance operating system.This property may be absent if the volume attachment's `status` is not `attached`.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isVolumesVolumeAttachmentsDeviceId: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "A unique identifier for the device which is exposed to the instance operating system.",
												},
											},
										},
									},
									isVolumesVolumeAttachmentsHref: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this volume attachment.",
									},
									isVolumesVolumeAttachmentsId: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this volume attachment.",
									},
									isVolumesVolumeAttachmentsInstance: &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The attached instance.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isVolumesVolumeAttachmentsInstanceCRN: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this virtual server instance.",
												},
												isVolumesVolumeAttachmentsInstanceDeleted: &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															isVolumesVolumeAttachmentsInstanceDeletedMoreInfo: &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Link to documentation about deleted resources.",
															},
														},
													},
												},
												isVolumesVolumeAttachmentsInstanceHref: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this virtual server instance.",
												},
												isVolumesVolumeAttachmentsInstanceId: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this virtual server instance.",
												},
												isVolumesVolumeAttachmentsInstanceName: &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The user-defined name for this virtual server instance (and default system hostname).",
												},
											},
										},
									},
									isVolumesVolumeAttachmentsName: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this volume attachment.",
									},
									isVolumesVolumeAttachmentsType: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of volume attachment.",
									},
								},
							},
						},
						isVolumesZone: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The zone this volume resides in.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVolumesZoneHref: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this zone.",
									},
									isVolumesZoneName: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this zone.",
									},
								},
							},
						},
						isVolumeTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "User Tags for the Volume",
						},
						isVolumeAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "Access management tags for the volume instance",
						},
						isVolumeHealthReasons: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVolumeHealthReasonsCode: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the reason for this health state.",
									},

									isVolumeHealthReasonsMessage: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the reason for this health state.",
									},

									isVolumeHealthReasonsMoreInfo: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about the reason for this health state.",
									},
								},
							},
						},
						isVolumeCatalogOffering: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The catalog offering this volume was created from. If a virtual server instance is provisioned with a boot_volume_attachment specifying this volume, the virtual server instance will use this volume's catalog offering, including its pricing plan.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVolumeCatalogOfferingPlanCrn: {
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
									isVolumeCatalogOfferingVersionCrn: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this version of a catalog offering",
									},
								},
							},
						},
						isVolumeHealthState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health of this resource.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsVolumesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	// filters - volume-name and zone-name
	volumeName := d.Get("volume_name").(string)
	zoneName := d.Get("zone_name").(string)
	attachmentState := d.Get(isVolumeAttachmentState).(string)
	encryption := d.Get(isVolumesEncryption).(string)
	operatingSystemFamily := d.Get(isVolumesOperatingSystemFamily).(string)
	operatingSystemArch := d.Get(isVolumesOperatingSystemArch).(string)

	start := ""
	allrecs := []vpcv1.Volume{}
	listVolumesOptions := &vpcv1.ListVolumesOptions{}
	if start != "" {
		listVolumesOptions.Start = &start
	}
	if volumeName != "" {
		listVolumesOptions.Name = &volumeName
	}
	if zoneName != "" {
		listVolumesOptions.ZoneName = &zoneName
	}
	if attachmentState != "" {
		listVolumesOptions.AttachmentState = &attachmentState
	}
	if encryption != "" {
		listVolumesOptions.Encryption = &encryption
	}
	if operatingSystemFamily != "" {
		listVolumesOptions.OperatingSystemFamily = &operatingSystemFamily
	}
	if operatingSystemArch != "" {
		listVolumesOptions.OperatingSystemArchitecture = &operatingSystemArch
	}

	// list
	for {
		if start != "" {
			listVolumesOptions.Start = &start
		}
		volumeCollection, response, err := vpcClient.ListVolumesWithContext(context, listVolumesOptions)
		if err != nil {
			log.Printf("[DEBUG] ListVolumesWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListVolumesWithContext failed %s\n%s", err, response))
		}

		start = flex.GetNext(volumeCollection.Next)
		allrecs = append(allrecs, volumeCollection.Volumes...)

		if start == "" {
			break
		}

	}

	d.SetId(dataSourceIBMIsVolumesID(d))

	err = d.Set(isVolumes, dataSourceVolumeCollectionFlattenVolumes(allrecs, meta))
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error setting volumes %s", err))
	}

	return nil
}

// dataSourceIBMIsVolumesID returns a reasonable ID for the list.
func dataSourceIBMIsVolumesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceVolumeCollectionFlattenVolumes(result []vpcv1.Volume, meta interface{}) (volumes []map[string]interface{}) {
	for _, volumesItem := range result {
		volumes = append(volumes, dataSourceVolumeCollectionVolumesToMap(volumesItem, meta))
	}

	return volumes
}

func dataSourceVolumeCollectionVolumesToMap(volumesItem vpcv1.Volume, meta interface{}) (volumesMap map[string]interface{}) {
	volumesMap = map[string]interface{}{}

	if volumesItem.Active != nil {
		volumesMap[isVolumesActive] = volumesItem.Active
	}
	if volumesItem.AttachmentState != nil {
		volumesMap[isVolumeAttachmentState] = volumesItem.AttachmentState
	}
	if volumesItem.Bandwidth != nil {
		volumesMap[isVolumesBandwidth] = volumesItem.Bandwidth
	}
	if volumesItem.Busy != nil {
		volumesMap[isVolumesBusy] = volumesItem.Busy
	}
	if volumesItem.Capacity != nil {
		volumesMap[isVolumesCapacity] = volumesItem.Capacity
	}
	if volumesItem.CreatedAt != nil {
		volumesMap[isVolumesCreatedAt] = volumesItem.CreatedAt.String()
	}
	if volumesItem.CRN != nil {
		volumesMap[isVolumesCRN] = volumesItem.CRN
	}
	if volumesItem.Encryption != nil {
		volumesMap[isVolumesEncryption] = volumesItem.Encryption
	}
	if volumesItem.EncryptionKey != nil {
		encryptionKeyList := []map[string]interface{}{}
		encryptionKeyMap := dataSourceVolumeCollectionVolumesEncryptionKeyToMap(*volumesItem.EncryptionKey)
		encryptionKeyList = append(encryptionKeyList, encryptionKeyMap)
		volumesMap[isVolumesEncryptionKey] = encryptionKeyList
	}
	if volumesItem.Href != nil {
		volumesMap[isVolumesHref] = volumesItem.Href
	}
	if volumesItem.ID != nil {
		volumesMap[isVolumesId] = volumesItem.ID
	}
	if volumesItem.Iops != nil {
		volumesMap[isVolumesIops] = volumesItem.Iops
	}
	if volumesItem.Name != nil {
		volumesMap[isVolumesName] = volumesItem.Name
	}
	if volumesItem.OperatingSystem != nil {
		operatingSystemList := []map[string]interface{}{}
		operatingSystemMap := dataSourceVolumeCollectionVolumesOperatingSystemToMap(*volumesItem.OperatingSystem)
		operatingSystemList = append(operatingSystemList, operatingSystemMap)
		volumesMap[isVolumesOperatingSystem] = operatingSystemList
	}
	if volumesItem.Profile != nil {
		profileList := []map[string]interface{}{}
		profileMap := dataSourceVolumeCollectionVolumesProfileToMap(*volumesItem.Profile)
		profileList = append(profileList, profileMap)
		volumesMap[isVolumesProfile] = profileList
	}
	if volumesItem.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourceVolumeCollectionVolumesResourceGroupToMap(*volumesItem.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		volumesMap[isVolumesResourceGroup] = resourceGroupList
	}
	if volumesItem.SourceImage != nil {
		sourceImageList := []map[string]interface{}{}
		sourceImageMap := dataSourceVolumeCollectionVolumesSourceImageToMap(*volumesItem.SourceImage)
		sourceImageList = append(sourceImageList, sourceImageMap)
		volumesMap[isVolumesSourceImage] = sourceImageList
	}
	if volumesItem.SourceSnapshot != nil {
		sourceSnapshotList := []map[string]interface{}{}
		sourceSnapshotMap := dataSourceVolumeCollectionVolumesSourceSnapshotToMap(*volumesItem.SourceSnapshot)
		sourceSnapshotList = append(sourceSnapshotList, sourceSnapshotMap)
		volumesMap[isVolumesSourceSnapshot] = sourceSnapshotList
	}
	if volumesItem.Status != nil {
		volumesMap[isVolumesStatus] = volumesItem.Status
	}
	if volumesItem.HealthState != nil {
		volumesMap[isVolumeHealthState] = volumesItem.HealthState
	}
	if volumesItem.StatusReasons != nil {
		statusReasonsList := []map[string]interface{}{}
		for _, statusReasonsItem := range volumesItem.StatusReasons {
			statusReasonsList = append(statusReasonsList, dataSourceVolumeCollectionVolumesStatusReasonsToMap(statusReasonsItem))
		}
		volumesMap[isVolumesStatusReasons] = statusReasonsList
	}
	if volumesItem.HealthReasons != nil {
		healthReasonsList := []map[string]interface{}{}
		for _, healthReasonsItem := range volumesItem.HealthReasons {
			healthReasonsList = append(healthReasonsList, dataSourceVolumeCollectionVolumesHealthReasonsToMap(healthReasonsItem))
		}
		volumesMap[isVolumeHealthReasons] = healthReasonsList
	}
	if volumesItem.CatalogOffering != nil {
		versionCrn := ""
		if volumesItem.CatalogOffering.Version != nil && volumesItem.CatalogOffering.Version.CRN != nil {
			versionCrn = *volumesItem.CatalogOffering.Version.CRN
		}
		catalogList := make([]map[string]interface{}, 0)
		catalogMap := map[string]interface{}{}
		if versionCrn != "" {
			catalogMap[isVolumeCatalogOfferingVersionCrn] = versionCrn
		}
		if volumesItem.CatalogOffering.Plan != nil {
			planCrn := ""
			if volumesItem.CatalogOffering.Plan.CRN != nil {
				planCrn = *volumesItem.CatalogOffering.Plan.CRN
			}
			if planCrn != "" {
				catalogMap[isVolumeCatalogOfferingPlanCrn] = *volumesItem.CatalogOffering.Plan.CRN
			}
			if volumesItem.CatalogOffering.Plan.Deleted != nil {
				deletedMap := resourceIbmIsVolumeCatalogOfferingVersionPlanReferenceDeletedToMap(*volumesItem.CatalogOffering.Plan.Deleted)
				catalogMap["deleted"] = []map[string]interface{}{deletedMap}
			}
		}
		catalogList = append(catalogList, catalogMap)
		volumesMap[isVolumeCatalogOffering] = catalogList
	}
	if volumesItem.VolumeAttachments != nil {
		volumeAttachmentsList := []map[string]interface{}{}
		for _, volumeAttachmentsItem := range volumesItem.VolumeAttachments {
			volumeAttachmentsList = append(volumeAttachmentsList, dataSourceVolumeCollectionVolumesVolumeAttachmentsToMap(volumeAttachmentsItem))
		}
		volumesMap[isVolumesVolumeAttachments] = volumeAttachmentsList
	}
	if volumesItem.Zone != nil {
		zoneList := []map[string]interface{}{}
		zoneMap := dataSourceVolumeCollectionVolumesZoneToMap(*volumesItem.Zone)
		zoneList = append(zoneList, zoneMap)
		volumesMap[isVolumesZone] = zoneList
	}
	if volumesItem.UserTags != nil {
		volumesMap[isVolumeTags] = volumesItem.UserTags
	}
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *volumesItem.CRN, "", isVolumeAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc volume (%s) access tags: %s", *volumesItem.ID, err)
	}
	volumesMap[isVolumeAccessTags] = accesstags
	return volumesMap
}

func dataSourceVolumeCollectionVolumesEncryptionKeyToMap(encryptionKeyItem vpcv1.EncryptionKeyReference) (encryptionKeyMap map[string]interface{}) {
	encryptionKeyMap = map[string]interface{}{}

	if encryptionKeyItem.CRN != nil {
		encryptionKeyMap[isVolumesEncryptionKeyCRN] = encryptionKeyItem.CRN
	}

	return encryptionKeyMap
}

func dataSourceVolumeCollectionVolumesOperatingSystemToMap(operatingSystemItem vpcv1.OperatingSystem) (operatingSystemMap map[string]interface{}) {
	operatingSystemMap = map[string]interface{}{}

	if operatingSystemItem.Architecture != nil {
		operatingSystemMap[isVolumeArchitecture] = operatingSystemItem.Architecture
	}
	if operatingSystemItem.DedicatedHostOnly != nil {
		operatingSystemMap[isVolumeDHOnly] = operatingSystemItem.DedicatedHostOnly
	}
	if operatingSystemItem.DisplayName != nil {
		operatingSystemMap[isVolumeDisplayName] = operatingSystemItem.DisplayName
	}
	if operatingSystemItem.Family != nil {
		operatingSystemMap[isVolumeOSFamily] = operatingSystemItem.Family
	}
	if operatingSystemItem.Href != nil {
		operatingSystemMap[isVolumesOperatingSystemHref] = operatingSystemItem.Href
	}
	if operatingSystemItem.Name != nil {
		operatingSystemMap[isVolumesOperatingSystemName] = operatingSystemItem.Name
	}
	if operatingSystemItem.Vendor != nil {
		operatingSystemMap[isVolumeOSVendor] = operatingSystemItem.Vendor
	}
	if operatingSystemItem.Version != nil {
		operatingSystemMap[isVolumeOSVersion] = operatingSystemItem.Version
	}

	return operatingSystemMap
}

func dataSourceVolumeCollectionVolumesProfileToMap(profileItem vpcv1.VolumeProfileReference) (profileMap map[string]interface{}) {
	profileMap = map[string]interface{}{}

	if profileItem.Href != nil {
		profileMap[isVolumesProfileHref] = profileItem.Href
	}
	if profileItem.Name != nil {
		profileMap[isVolumesProfileName] = profileItem.Name
	}

	return profileMap
}

func dataSourceVolumeCollectionVolumesResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
	resourceGroupMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resourceGroupMap[isVolumesResourceGroupHref] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resourceGroupMap[isVolumesResourceGroupId] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resourceGroupMap[isVolumesResourceGroupName] = resourceGroupItem.Name
	}

	return resourceGroupMap
}

func dataSourceVolumeCollectionVolumesSourceImageToMap(sourceImageItem vpcv1.ImageReference) (sourceImageMap map[string]interface{}) {
	sourceImageMap = map[string]interface{}{}

	if sourceImageItem.CRN != nil {
		sourceImageMap[isVolumesSourceImageCRN] = sourceImageItem.CRN
	}
	if sourceImageItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVolumeCollectionSourceImageDeletedToMap(*sourceImageItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		sourceImageMap[isVolumesSourceImageDeleted] = deletedList
	}
	if sourceImageItem.Href != nil {
		sourceImageMap[isVolumesSourceImageHref] = sourceImageItem.Href
	}
	if sourceImageItem.ID != nil {
		sourceImageMap[isVolumesSourceImageId] = sourceImageItem.ID
	}
	if sourceImageItem.Name != nil {
		sourceImageMap[isVolumesSourceImageName] = sourceImageItem.Name
	}

	return sourceImageMap
}

func dataSourceVolumeCollectionSourceImageDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap[isVolumesSourceImageDeletedMoreInfo] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceVolumeCollectionVolumesSourceSnapshotToMap(sourceSnapshotItem vpcv1.SnapshotReference) (sourceSnapshotMap map[string]interface{}) {
	sourceSnapshotMap = map[string]interface{}{}

	if sourceSnapshotItem.CRN != nil {
		sourceSnapshotMap[isVolumesSourceSnapshotCRN] = sourceSnapshotItem.CRN
	}
	if sourceSnapshotItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVolumeCollectionSourceSnapshotDeletedToMap(*sourceSnapshotItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		sourceSnapshotMap[isVolumesSourceSnapshotDeleted] = deletedList
	}
	if sourceSnapshotItem.Href != nil {
		sourceSnapshotMap[isVolumesSourceSnapshotHref] = sourceSnapshotItem.Href
	}
	if sourceSnapshotItem.ID != nil {
		sourceSnapshotMap[isVolumesSourceSnapshotId] = sourceSnapshotItem.ID
	}
	if sourceSnapshotItem.Name != nil {
		sourceSnapshotMap[isVolumesSourceSnapshotName] = sourceSnapshotItem.Name
	}
	if sourceSnapshotItem.ResourceType != nil {
		sourceSnapshotMap[isVolumesSourceSnapshotResourceType] = sourceSnapshotItem.ResourceType
	}

	return sourceSnapshotMap
}

func dataSourceVolumeCollectionSourceSnapshotDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap[isVolumesSourceSnapshotDeletedMoreInfo] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceVolumeCollectionVolumesStatusReasonsToMap(statusReasonsItem vpcv1.VolumeStatusReason) (statusReasonsMap map[string]interface{}) {
	statusReasonsMap = map[string]interface{}{}

	if statusReasonsItem.Code != nil {
		statusReasonsMap[isVolumesStatusReasonsCode] = statusReasonsItem.Code
	}
	if statusReasonsItem.Message != nil {
		statusReasonsMap[isVolumesStatusReasonsMessage] = statusReasonsItem.Message
	}
	if statusReasonsItem.MoreInfo != nil {
		statusReasonsMap[isVolumesStatusReasonsMoreInfo] = statusReasonsItem.MoreInfo
	}

	return statusReasonsMap
}

func dataSourceVolumeCollectionVolumesHealthReasonsToMap(statusReasonsItem vpcv1.VolumeHealthReason) (healthReasonsMap map[string]interface{}) {
	healthReasonsMap = map[string]interface{}{}

	if statusReasonsItem.Code != nil {
		healthReasonsMap[isVolumeHealthReasonsCode] = statusReasonsItem.Code
	}
	if statusReasonsItem.Message != nil {
		healthReasonsMap[isVolumeHealthReasonsMessage] = statusReasonsItem.Message
	}
	if statusReasonsItem.MoreInfo != nil {
		healthReasonsMap[isVolumeHealthReasonsMoreInfo] = statusReasonsItem.MoreInfo
	}

	return healthReasonsMap
}

func dataSourceVolumeCollectionVolumesVolumeAttachmentsToMap(volumeAttachmentsItem vpcv1.VolumeAttachmentReferenceVolumeContext) (volumeAttachmentsMap map[string]interface{}) {
	volumeAttachmentsMap = map[string]interface{}{}

	if volumeAttachmentsItem.DeleteVolumeOnInstanceDelete != nil {
		volumeAttachmentsMap[isVolumesVolumeAttachmentsDeleteVolumeOnInstanceDelete] = volumeAttachmentsItem.DeleteVolumeOnInstanceDelete
	}
	if volumeAttachmentsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVolumeCollectionVolumeAttachmentsDeletedToMap(*volumeAttachmentsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		volumeAttachmentsMap[isVolumesVolumeAttachmentsDeleted] = deletedList
	}
	if volumeAttachmentsItem.Device != nil {
		deviceList := []map[string]interface{}{}
		deviceMap := dataSourceVolumeCollectionVolumeAttachmentsDeviceToMap(*volumeAttachmentsItem.Device)
		deviceList = append(deviceList, deviceMap)
		volumeAttachmentsMap[isVolumesVolumeAttachmentsDevice] = deviceList
	}
	if volumeAttachmentsItem.Href != nil {
		volumeAttachmentsMap[isVolumesVolumeAttachmentsHref] = volumeAttachmentsItem.Href
	}
	if volumeAttachmentsItem.ID != nil {
		volumeAttachmentsMap[isVolumesVolumeAttachmentsId] = volumeAttachmentsItem.ID
	}
	if volumeAttachmentsItem.Instance != nil {
		instanceList := []map[string]interface{}{}
		instanceMap := dataSourceVolumeCollectionVolumeAttachmentsInstanceToMap(*volumeAttachmentsItem.Instance)
		instanceList = append(instanceList, instanceMap)
		volumeAttachmentsMap[isVolumesVolumeAttachmentsInstance] = instanceList
	}
	if volumeAttachmentsItem.Name != nil {
		volumeAttachmentsMap[isVolumesVolumeAttachmentsName] = volumeAttachmentsItem.Name
	}
	if volumeAttachmentsItem.Type != nil {
		volumeAttachmentsMap[isVolumesVolumeAttachmentsType] = volumeAttachmentsItem.Type
	}

	return volumeAttachmentsMap
}

func dataSourceVolumeCollectionVolumeAttachmentsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap[isVolumesVolumeAttachmentsDeletedMoreInfo] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceVolumeCollectionVolumeAttachmentsDeviceToMap(deviceItem vpcv1.VolumeAttachmentDevice) (deviceMap map[string]interface{}) {
	deviceMap = map[string]interface{}{}

	if deviceItem.ID != nil {
		deviceMap[isVolumesVolumeAttachmentsDeviceId] = deviceItem.ID
	}

	return deviceMap
}

func dataSourceVolumeCollectionVolumeAttachmentsInstanceToMap(instanceItem vpcv1.InstanceReference) (instanceMap map[string]interface{}) {
	instanceMap = map[string]interface{}{}

	if instanceItem.CRN != nil {
		instanceMap[isVolumesVolumeAttachmentsInstanceCRN] = instanceItem.CRN
	}
	if instanceItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVolumeCollectionInstanceDeletedToMap(*instanceItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		instanceMap[isVolumesVolumeAttachmentsInstanceDeleted] = deletedList
	}
	if instanceItem.Href != nil {
		instanceMap[isVolumesVolumeAttachmentsInstanceHref] = instanceItem.Href
	}
	if instanceItem.ID != nil {
		instanceMap[isVolumesVolumeAttachmentsInstanceId] = instanceItem.ID
	}
	if instanceItem.Name != nil {
		instanceMap[isVolumesVolumeAttachmentsInstanceName] = instanceItem.Name
	}

	return instanceMap
}

func dataSourceVolumeCollectionInstanceDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap[isVolumesVolumeAttachmentsInstanceDeletedMoreInfo] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceVolumeCollectionVolumesZoneToMap(zoneItem vpcv1.ZoneReference) (zoneMap map[string]interface{}) {
	zoneMap = map[string]interface{}{}

	if zoneItem.Href != nil {
		zoneMap[isVolumesZoneHref] = zoneItem.Href
	}
	if zoneItem.Name != nil {
		zoneMap[isVolumesZoneName] = zoneItem.Name
	}

	return zoneMap
}
