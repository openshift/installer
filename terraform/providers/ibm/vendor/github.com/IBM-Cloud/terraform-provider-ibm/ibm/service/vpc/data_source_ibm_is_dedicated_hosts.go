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
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIbmIsDedicatedHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsDedicatedHostsRead,

		Schema: map[string]*schema.Schema{
			"host_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the dedicated host group this dedicated host belongs to",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the resource group this dedicated host belongs to",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The zone name this dedicated host is in",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the dedicated host",
			},
			"dedicated_hosts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of dedicated hosts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"available_memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of memory in gibibytes that is currently available for instances.",
						},
						"available_vcpu": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The available VCPU for the dedicated host.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"architecture": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The VCPU architecture.",
									},
									"manufacturer": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The VCPU manufacturer.",
									},
									"count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of VCPUs assigned.",
									},
								},
							},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the dedicated host was created.",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this dedicated host.",
						},
						"disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of the dedicated host's disks.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"available": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The remaining space left for instance placement in GB (gigabytes).",
									},
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time that the disk was created.",
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this disk.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this disk.",
									},
									"instance_disks": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Instance disks that are on this dedicated host disk.",
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
													Description: "The URL for this instance disk.",
												},
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this instance disk.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The user-defined name for this disk.",
												},
												"resource_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"interface_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The disk interface used for attaching the diskThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
									},
									"lifecycle_state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The lifecycle state of this dedicated host disk.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined or system-provided name for this disk.",
									},
									"provisionable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether this dedicated host disk is available for instance disk creation.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of resource referenced.",
									},
									"size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The size of the disk in GB (gigabytes).",
									},
									"supported_instance_interface_types": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The instance disk interfaces supported for this dedicated host disk.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"host_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the dedicated host group this dedicated host is in.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this dedicated host.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this dedicated host.",
						},
						"instance_placement_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to true, instances can be placed on this dedicated host.",
						},
						"instances": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Array of instances that are allocated to this dedicated host.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this virtual server instance.",
									},
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
										Description: "The URL for this virtual server instance.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this virtual server instance.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this virtual server instance (and default system hostname).",
									},
								},
							},
						},
						"lifecycle_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the dedicated host resource.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total amount of memory in gibibytes for this host.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this dedicated host. If unspecified, the name will be a hyphenated list of randomly-selected words.",
						},
						"numa": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The dedicated host NUMA configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The total number of NUMA nodes for this dedicated host",
									},
									"nodes": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The NUMA nodes for this dedicated host.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"available_vcpu": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The available VCPU for this NUMA node.",
												},
												"vcpu": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The total VCPU capacity for this NUMA node.",
												},
											},
										},
									},
								},
							},
						},
						"profile": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The profile this dedicated host uses.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this dedicated host.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this dedicated host profile.",
									},
								},
							},
						},
						"provisionable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this dedicated host is available for instance creation.",
						},
						"resource_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the resource group for this dedicated host.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
						"socket_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of sockets for this host.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The administrative state of the dedicated host.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the dedicated host on which the unexpected property value was encountered.",
						},
						"supported_instance_profiles": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Array of instance profiles that can be used by instances placed on this dedicated host.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this virtual server instance profile.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this virtual server instance profile.",
									},
								},
							},
						},
						"vcpu": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The total VCPU of the dedicated host.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"architecture": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The VCPU architecture.",
									},
									"manufacturer": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The VCPU manufacturer.",
									},
									"count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of VCPUs assigned.",
									},
								},
							},
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name of the zone this dedicated host resides in.",
						},
						isDedicatedHostAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access tags",
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

func dataSourceIbmIsDedicatedHostsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listDedicatedHostsOptions := &vpcv1.ListDedicatedHostsOptions{}
	if hostgroupintf, ok := d.GetOk("host_group"); ok {
		hostgroupid := hostgroupintf.(string)
		listDedicatedHostsOptions.DedicatedHostGroupID = &hostgroupid
	}
	if resgroupintf, ok := d.GetOk("resource_group"); ok {
		resGroup := resgroupintf.(string)
		listDedicatedHostsOptions.ResourceGroupID = &resGroup
	}
	if zoneintf, ok := d.GetOk("zone"); ok {
		zoneName := zoneintf.(string)
		listDedicatedHostsOptions.ZoneName = &zoneName
	}
	if nameintf, ok := d.GetOk("name"); ok {
		name := nameintf.(string)
		listDedicatedHostsOptions.Name = &name
	}
	start := ""
	allrecs := []vpcv1.DedicatedHost{}
	for {
		if start != "" {
			listDedicatedHostsOptions.Start = &start
		}
		dedicatedHostCollection, response, err := vpcClient.ListDedicatedHostsWithContext(context, listDedicatedHostsOptions)
		if err != nil {
			log.Printf("[DEBUG] ListDedicatedHostsWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		start = flex.GetNext(dedicatedHostCollection.Next)
		allrecs = append(allrecs, dedicatedHostCollection.DedicatedHosts...)
		if start == "" {
			break
		}
	}

	if len(allrecs) > 0 {

		d.SetId(dataSourceIbmIsDedicatedHostsID(d))

		err = d.Set("dedicated_hosts", dataSourceDedicatedHostCollectionFlattenDedicatedHosts(allrecs, meta))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting dedicated_hosts %s", err))
		}

		if err = d.Set("total_count", len(allrecs)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting total_count: %s", err))
		}
	}
	return nil
}

// dataSourceIbmIsDedicatedHostsID returns a reasonable ID for the list.
func dataSourceIbmIsDedicatedHostsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceDedicatedHostCollectionFlattenDedicatedHosts(result []vpcv1.DedicatedHost, meta interface{}) (dedicatedHosts []map[string]interface{}) {
	for _, dedicatedHostsItem := range result {
		dedicatedHosts = append(dedicatedHosts, dataSourceDedicatedHostCollectionDedicatedHostsToMap(dedicatedHostsItem, meta))
	}

	return dedicatedHosts
}

func dataSourceDedicatedHostCollectionDedicatedHostsToMap(dedicatedHostsItem vpcv1.DedicatedHost, meta interface{}) (dedicatedHostsMap map[string]interface{}) {
	dedicatedHostsMap = map[string]interface{}{}

	if dedicatedHostsItem.AvailableMemory != nil {
		dedicatedHostsMap["available_memory"] = dedicatedHostsItem.AvailableMemory
	}
	if dedicatedHostsItem.AvailableVcpu != nil {
		availableVcpuList := []map[string]interface{}{}
		availableVcpuMap := dataSourceDedicatedHostCollectionDedicatedHostsAvailableVcpuToMap(*dedicatedHostsItem.AvailableVcpu)
		availableVcpuList = append(availableVcpuList, availableVcpuMap)
		dedicatedHostsMap["available_vcpu"] = availableVcpuList
	}
	if dedicatedHostsItem.CreatedAt != nil {
		dedicatedHostsMap["created_at"] = dedicatedHostsItem.CreatedAt.String()
	}
	if dedicatedHostsItem.CRN != nil {
		dedicatedHostsMap["crn"] = dedicatedHostsItem.CRN
	}
	if dedicatedHostsItem.Disks != nil {
		disksList := []map[string]interface{}{}
		for _, disksItem := range dedicatedHostsItem.Disks {
			disksList = append(disksList, dataSourceDedicatedHostCollectionDedicatedHostsDisksToMap(disksItem))
		}
		dedicatedHostsMap["disks"] = disksList
	}
	if dedicatedHostsItem.Group != nil {
		dedicatedHostsMap["host_group"] = *dedicatedHostsItem.Group.ID
	}
	if dedicatedHostsItem.Href != nil {
		dedicatedHostsMap["href"] = dedicatedHostsItem.Href
	}
	if dedicatedHostsItem.ID != nil {
		dedicatedHostsMap["id"] = dedicatedHostsItem.ID
	}
	if dedicatedHostsItem.InstancePlacementEnabled != nil {
		dedicatedHostsMap["instance_placement_enabled"] = dedicatedHostsItem.InstancePlacementEnabled
	}
	if dedicatedHostsItem.Instances != nil {
		instancesList := []map[string]interface{}{}
		for _, instancesItem := range dedicatedHostsItem.Instances {
			instancesList = append(instancesList, dataSourceDedicatedHostCollectionDedicatedHostsInstancesToMap(instancesItem))
		}
		dedicatedHostsMap["instances"] = instancesList
	}
	if dedicatedHostsItem.LifecycleState != nil {
		dedicatedHostsMap["lifecycle_state"] = dedicatedHostsItem.LifecycleState
	}
	if dedicatedHostsItem.Memory != nil {
		dedicatedHostsMap["memory"] = dedicatedHostsItem.Memory
	}
	if dedicatedHostsItem.Name != nil {
		dedicatedHostsMap["name"] = dedicatedHostsItem.Name
	}
	if dedicatedHostsItem.Numa != nil {
		dedicatedHostsMap["numa"] = dataSourceDedicatedHostFlattenNumaNodes(*dedicatedHostsItem.Numa)
	}
	if dedicatedHostsItem.Profile != nil {
		profileList := []map[string]interface{}{}
		profileMap := dataSourceDedicatedHostCollectionDedicatedHostsProfileToMap(*dedicatedHostsItem.Profile)
		profileList = append(profileList, profileMap)
		dedicatedHostsMap["profile"] = profileList
	}
	if dedicatedHostsItem.Provisionable != nil {
		dedicatedHostsMap["provisionable"] = dedicatedHostsItem.Provisionable
	}
	if dedicatedHostsItem.ResourceGroup != nil {
		dedicatedHostsMap["resource_group"] = *dedicatedHostsItem.ResourceGroup.ID
	}
	if dedicatedHostsItem.ResourceType != nil {
		dedicatedHostsMap["resource_type"] = dedicatedHostsItem.ResourceType
	}
	if dedicatedHostsItem.SocketCount != nil {
		dedicatedHostsMap["socket_count"] = dedicatedHostsItem.SocketCount
	}
	if dedicatedHostsItem.State != nil {
		dedicatedHostsMap["state"] = dedicatedHostsItem.State
	}
	if dedicatedHostsItem.SupportedInstanceProfiles != nil {
		supportedInstanceProfilesList := []map[string]interface{}{}
		for _, supportedInstanceProfilesItem := range dedicatedHostsItem.SupportedInstanceProfiles {
			supportedInstanceProfilesList = append(supportedInstanceProfilesList, dataSourceDedicatedHostCollectionDedicatedHostsSupportedInstanceProfilesToMap(supportedInstanceProfilesItem))
		}
		dedicatedHostsMap["supported_instance_profiles"] = supportedInstanceProfilesList
	}
	if dedicatedHostsItem.Vcpu != nil {
		vcpuList := []map[string]interface{}{}
		vcpuMap := dataSourceDedicatedHostCollectionDedicatedHostsVcpuToMap(*dedicatedHostsItem.Vcpu)
		vcpuList = append(vcpuList, vcpuMap)
		dedicatedHostsMap["vcpu"] = vcpuList
	}
	if dedicatedHostsItem.Zone != nil {
		dedicatedHostsMap["zone"] = *dedicatedHostsItem.Zone.Name
	}
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *dedicatedHostsItem.CRN, "", isDedicatedHostAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource dedicated host (%s) access tags: %s", *dedicatedHostsItem.ID, err)
	}
	dedicatedHostsMap[isDedicatedHostAccessTags] = accesstags

	return dedicatedHostsMap
}

func dataSourceDedicatedHostCollectionDedicatedHostsAvailableVcpuToMap(availableVcpuItem vpcv1.Vcpu) (availableVcpuMap map[string]interface{}) {
	availableVcpuMap = map[string]interface{}{}

	if availableVcpuItem.Architecture != nil {
		availableVcpuMap["architecture"] = availableVcpuItem.Architecture
	}
	// Added AMD Support for the manufacturer.
	if availableVcpuItem.Manufacturer != nil {
		availableVcpuMap["manufacturer"] = availableVcpuItem.Manufacturer
	}
	if availableVcpuItem.Count != nil {
		availableVcpuMap["count"] = availableVcpuItem.Count
	}

	return availableVcpuMap
}

func dataSourceDedicatedHostCollectionDedicatedHostsGroupToMap(groupItem vpcv1.DedicatedHostGroupReference) (groupMap map[string]interface{}) {
	groupMap = map[string]interface{}{}

	if groupItem.CRN != nil {
		groupMap["crn"] = groupItem.CRN
	}
	if groupItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceDedicatedHostCollectionGroupDeletedToMap(*groupItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		groupMap["deleted"] = deletedList
	}
	if groupItem.Href != nil {
		groupMap["href"] = groupItem.Href
	}
	if groupItem.ID != nil {
		groupMap["id"] = groupItem.ID
	}
	if groupItem.Name != nil {
		groupMap["name"] = groupItem.Name
	}
	if groupItem.ResourceType != nil {
		groupMap["resource_type"] = groupItem.ResourceType
	}

	return groupMap
}

func dataSourceDedicatedHostCollectionGroupDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceDedicatedHostCollectionInstancesDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceDedicatedHostCollectionDedicatedHostsInstancesToMap(instancesItem vpcv1.InstanceReference) (instancesMap map[string]interface{}) {
	instancesMap = map[string]interface{}{}

	if instancesItem.CRN != nil {
		instancesMap["crn"] = instancesItem.CRN
	}
	if instancesItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceDedicatedHostCollectionInstancesDeletedToMap(*instancesItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		instancesMap["deleted"] = deletedList
	}
	if instancesItem.Href != nil {
		instancesMap["href"] = instancesItem.Href
	}
	if instancesItem.ID != nil {
		instancesMap["id"] = instancesItem.ID
	}
	if instancesItem.Name != nil {
		instancesMap["name"] = instancesItem.Name
	}

	return instancesMap
}

func dataSourceDedicatedHostCollectionDedicatedHostsProfileToMap(profileItem vpcv1.DedicatedHostProfileReference) (profileMap map[string]interface{}) {
	profileMap = map[string]interface{}{}

	if profileItem.Href != nil {
		profileMap["href"] = profileItem.Href
	}
	if profileItem.Name != nil {
		profileMap["name"] = profileItem.Name
	}

	return profileMap
}

func dataSourceDedicatedHostCollectionDedicatedHostsSupportedInstanceProfilesToMap(supportedInstanceProfilesItem vpcv1.InstanceProfileReference) (supportedInstanceProfilesMap map[string]interface{}) {
	supportedInstanceProfilesMap = map[string]interface{}{}

	if supportedInstanceProfilesItem.Href != nil {
		supportedInstanceProfilesMap["href"] = supportedInstanceProfilesItem.Href
	}
	if supportedInstanceProfilesItem.Name != nil {
		supportedInstanceProfilesMap["name"] = supportedInstanceProfilesItem.Name
	}

	return supportedInstanceProfilesMap
}

func dataSourceDedicatedHostCollectionDedicatedHostsVcpuToMap(vcpuItem vpcv1.Vcpu) (vcpuMap map[string]interface{}) {
	vcpuMap = map[string]interface{}{}

	if vcpuItem.Architecture != nil {
		vcpuMap["architecture"] = vcpuItem.Architecture
	}
	// Added AMD Support for the manufacturer.
	if vcpuItem.Manufacturer != nil {
		vcpuMap["manufacturer"] = vcpuItem.Manufacturer
	}
	if vcpuItem.Count != nil {
		vcpuMap["count"] = vcpuItem.Count
	}

	return vcpuMap
}

func dataSourceDedicatedHostCollectionFlattenFirst(result vpcv1.DedicatedHostCollectionFirst) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostCollectionFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostCollectionFirstToMap(firstItem vpcv1.DedicatedHostCollectionFirst) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceDedicatedHostCollectionFlattenNext(result vpcv1.DedicatedHostCollectionNext) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostCollectionNextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostCollectionNextToMap(nextItem vpcv1.DedicatedHostCollectionNext) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}

func dataSourceDedicatedHostCollectionDedicatedHostsDisksToMap(disksItem vpcv1.DedicatedHostDisk) (disksMap map[string]interface{}) {
	disksMap = map[string]interface{}{}

	if disksItem.Available != nil {
		disksMap["available"] = disksItem.Available
	}
	if disksItem.CreatedAt != nil {
		disksMap["created_at"] = disksItem.CreatedAt.String()
	}
	if disksItem.Href != nil {
		disksMap["href"] = disksItem.Href
	}
	if disksItem.ID != nil {
		disksMap["id"] = disksItem.ID
	}
	if disksItem.InstanceDisks != nil {
		instanceDisksList := []map[string]interface{}{}
		for _, instanceDisksItem := range disksItem.InstanceDisks {
			instanceDisksList = append(instanceDisksList, dataSourceDedicatedHostDisksInstanceDisksToMap(instanceDisksItem))
		}
		disksMap["instance_disks"] = instanceDisksList
	}
	if disksItem.InterfaceType != nil {
		disksMap["interface_type"] = disksItem.InterfaceType
	}
	if disksItem.LifecycleState != nil {
		disksMap["lifecycle_state"] = disksItem.LifecycleState
	}
	if disksItem.Name != nil {
		disksMap["name"] = disksItem.Name
	}
	if disksItem.Provisionable != nil {
		disksMap["provisionable"] = disksItem.Provisionable
	}
	if disksItem.ResourceType != nil {
		disksMap["resource_type"] = disksItem.ResourceType
	}
	if disksItem.Size != nil {
		disksMap["size"] = disksItem.Size
	}
	if disksItem.SupportedInstanceInterfaceTypes != nil {
		disksMap["supported_instance_interface_types"] = disksItem.SupportedInstanceInterfaceTypes
	}

	return disksMap
}
