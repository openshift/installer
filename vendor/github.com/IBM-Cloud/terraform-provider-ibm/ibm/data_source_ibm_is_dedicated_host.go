// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func dataSourceIbmIsDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmIsDedicatedHostRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of this dedicated host",
			},
			"host_group": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the dedicated host group this dedicated host belongs to",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the resource group this dedicated host belongs to",
			},
			"available_memory": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of memory in gibibytes that is currently available for instances.",
			},
			"available_vcpu": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The available VCPU for the dedicated host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"architecture": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VCPU architecture.",
						},
						"count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of VCPUs assigned.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the dedicated host was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this dedicated host.",
			},
			"disks": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of the dedicated host's disks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"available": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The remaining space left for instance placement in GB (gigabytes).",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the disk was created.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this disk.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this disk.",
						},
						"instance_disks": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance disks that are on this dedicated host disk.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
										Description: "The URL for this instance disk.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this instance disk.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this disk.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"interface_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The disk interface used for attaching the diskThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of this dedicated host disk.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined or system-provided name for this disk.",
						},
						"provisionable": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this dedicated host disk is available for instance disk creation.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
						"size": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the disk in GB (gigabytes).",
						},
						"supported_instance_interface_types": &schema.Schema{
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
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this dedicated host.",
			},
			"instance_placement_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, instances can be placed on this dedicated host.",
			},
			"instances": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of instances that are allocated to this dedicated host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this virtual server instance.",
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
							Description: "The URL for this virtual server instance.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this virtual server instance.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this virtual server instance (and default system hostname).",
						},
					},
				},
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the dedicated host resource.",
			},
			"memory": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total amount of memory in gibibytes for this host.",
			},
			"profile": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The profile this dedicated host uses.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this dedicated host.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this dedicated host profile.",
						},
					},
				},
			},
			"provisionable": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this dedicated host is available for instance creation.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
			"socket_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of sockets for this host.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The administrative state of the dedicated host.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the dedicated host on which the unexpected property value was encountered.",
			},
			"supported_instance_profiles": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of instance profiles that can be used by instances placed on this dedicated host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this virtual server instance profile.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this virtual server instance profile.",
						},
					},
				},
			},
			"vcpu": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The total VCPU of the dedicated host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"architecture": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VCPU architecture.",
						},
						"count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of VCPUs assigned.",
						},
					},
				},
			},
			"zone": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique name of the zone this dedicated host resides in.",
			},
		},
	}
}

func dataSourceIbmIsDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	listDedicatedHostsOptions := &vpcv1.ListDedicatedHostsOptions{}
	hostgroupid := d.Get("host_group").(string)
	listDedicatedHostsOptions.DedicatedHostGroupID = &hostgroupid
	if resgrpid, ok := d.GetOk("resource_group"); ok {
		resgrpidstr := resgrpid.(string)
		listDedicatedHostsOptions.ResourceGroupID = &resgrpidstr
	}
	dedicatedHostCollection, response, err := vpcClient.ListDedicatedHostsWithContext(context.TODO(), listDedicatedHostsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListDedicatedHostsWithContext failed %s\n%s", err, response)
		return err
	}
	name := d.Get("name").(string)
	if len(dedicatedHostCollection.DedicatedHosts) != 0 {
		dedicatedHost := vpcv1.DedicatedHost{}
		for _, data := range dedicatedHostCollection.DedicatedHosts {
			if *data.Name == name {
				dedicatedHost = data
				d.SetId(*dedicatedHost.ID)

				if err = d.Set("available_memory", dedicatedHost.AvailableMemory); err != nil {
					return fmt.Errorf("Error setting available_memory: %s", err)
				}

				if dedicatedHost.AvailableVcpu != nil {
					err = d.Set("available_vcpu", dataSourceDedicatedHostFlattenAvailableVcpu(*dedicatedHost.AvailableVcpu))
					if err != nil {
						return fmt.Errorf("Error setting available_vcpu %s", err)
					}
				}
				if err = d.Set("created_at", dedicatedHost.CreatedAt.String()); err != nil {
					return fmt.Errorf("Error setting created_at: %s", err)
				}
				if err = d.Set("crn", dedicatedHost.CRN); err != nil {
					return fmt.Errorf("Error setting crn: %s", err)
				}
				if dedicatedHost.Disks != nil {
					err = d.Set("disks", dataSourceDedicatedHostFlattenDisks(dedicatedHost.Disks))
					if err != nil {
						return fmt.Errorf("Error setting disks %s", err)
					}
				}
				if dedicatedHost.Group != nil {
					err = d.Set("host_group", *dedicatedHost.Group.ID)
					if err != nil {
						return fmt.Errorf("Error setting group %s", err)
					}
				}
				if err = d.Set("href", dedicatedHost.Href); err != nil {
					return fmt.Errorf("Error setting href: %s", err)
				}
				if err = d.Set("instance_placement_enabled", dedicatedHost.InstancePlacementEnabled); err != nil {
					return fmt.Errorf("Error setting instance_placement_enabled: %s", err)
				}

				if dedicatedHost.Instances != nil {
					err = d.Set("instances", dataSourceDedicatedHostFlattenInstances(dedicatedHost.Instances))
					if err != nil {
						return fmt.Errorf("Error setting instances %s", err)
					}
				}
				if err = d.Set("lifecycle_state", dedicatedHost.LifecycleState); err != nil {
					return fmt.Errorf("Error setting lifecycle_state: %s", err)
				}
				if err = d.Set("memory", dedicatedHost.Memory); err != nil {
					return fmt.Errorf("Error setting memory: %s", err)
				}
				if err = d.Set("name", dedicatedHost.Name); err != nil {
					return fmt.Errorf("Error setting name: %s", err)
				}

				if dedicatedHost.Profile != nil {
					err = d.Set("profile", dataSourceDedicatedHostFlattenProfile(*dedicatedHost.Profile))
					if err != nil {
						return fmt.Errorf("Error setting profile %s", err)
					}
				}
				if err = d.Set("provisionable", dedicatedHost.Provisionable); err != nil {
					return fmt.Errorf("Error setting provisionable: %s", err)
				}

				if dedicatedHost.ResourceGroup != nil {
					err = d.Set("resource_group", *dedicatedHost.ResourceGroup.ID)
					if err != nil {
						return fmt.Errorf("Error setting resource_group %s", err)
					}
				}
				if err = d.Set("resource_type", dedicatedHost.ResourceType); err != nil {
					return fmt.Errorf("Error setting resource_type: %s", err)
				}
				if err = d.Set("socket_count", dedicatedHost.SocketCount); err != nil {
					return fmt.Errorf("Error setting socket_count: %s", err)
				}
				if err = d.Set("state", dedicatedHost.State); err != nil {
					return fmt.Errorf("Error setting state: %s", err)
				}

				if dedicatedHost.SupportedInstanceProfiles != nil {
					err = d.Set("supported_instance_profiles", dataSourceDedicatedHostFlattenSupportedInstanceProfiles(dedicatedHost.SupportedInstanceProfiles))
					if err != nil {
						return fmt.Errorf("Error setting supported_instance_profiles %s", err)
					}
				}

				if dedicatedHost.Vcpu != nil {
					err = d.Set("vcpu", dataSourceDedicatedHostFlattenVcpu(*dedicatedHost.Vcpu))
					if err != nil {
						return fmt.Errorf("Error setting vcpu %s", err)
					}
				}

				if dedicatedHost.Zone != nil {
					err = d.Set("zone", *dedicatedHost.Zone.Name)
					if err != nil {
						return fmt.Errorf("Error setting zone %s", err)
					}
				}

				return nil
			}
		}
	}
	return fmt.Errorf("No Dedicated Host found with name %s", name)
}

// dataSourceIbmIsDedicatedHostID returns a reasonable ID for the list.
func dataSourceIbmIsDedicatedHostID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceDedicatedHostFlattenAvailableVcpu(result vpcv1.Vcpu) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostAvailableVcpuToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostAvailableVcpuToMap(availableVcpuItem vpcv1.Vcpu) (availableVcpuMap map[string]interface{}) {
	availableVcpuMap = map[string]interface{}{}

	if availableVcpuItem.Architecture != nil {
		availableVcpuMap["architecture"] = availableVcpuItem.Architecture
	}
	if availableVcpuItem.Count != nil {
		availableVcpuMap["count"] = availableVcpuItem.Count
	}

	return availableVcpuMap
}

func dataSourceDedicatedHostFlattenGroup(result vpcv1.DedicatedHostGroupReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostGroupToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostGroupToMap(groupItem vpcv1.DedicatedHostGroupReference) (groupMap map[string]interface{}) {
	groupMap = map[string]interface{}{}

	if groupItem.CRN != nil {
		groupMap["crn"] = groupItem.CRN
	}
	if groupItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceDedicatedHostGroupDeletedToMap(*groupItem.Deleted)
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

func dataSourceDedicatedHostGroupDeletedToMap(deletedItem vpcv1.DedicatedHostGroupReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceDedicatedHostFlattenInstances(result []vpcv1.InstanceReference) (instances []map[string]interface{}) {
	for _, instancesItem := range result {
		instances = append(instances, dataSourceDedicatedHostInstancesToMap(instancesItem))
	}

	return instances
}

func dataSourceDedicatedHostInstancesToMap(instancesItem vpcv1.InstanceReference) (instancesMap map[string]interface{}) {
	instancesMap = map[string]interface{}{}

	if instancesItem.CRN != nil {
		instancesMap["crn"] = instancesItem.CRN
	}
	if instancesItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceDedicatedHostInstancesDeletedToMap(*instancesItem.Deleted)
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

func dataSourceDedicatedHostInstancesDeletedToMap(deletedItem vpcv1.InstanceReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceDedicatedHostFlattenProfile(result vpcv1.DedicatedHostProfileReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostProfileToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostProfileToMap(profileItem vpcv1.DedicatedHostProfileReference) (profileMap map[string]interface{}) {
	profileMap = map[string]interface{}{}

	if profileItem.Href != nil {
		profileMap["href"] = profileItem.Href
	}
	if profileItem.Name != nil {
		profileMap["name"] = profileItem.Name
	}

	return profileMap
}

func dataSourceDedicatedHostFlattenSupportedInstanceProfiles(result []vpcv1.InstanceProfileReference) (supportedInstanceProfiles []map[string]interface{}) {
	for _, supportedInstanceProfilesItem := range result {
		supportedInstanceProfiles = append(supportedInstanceProfiles, dataSourceDedicatedHostSupportedInstanceProfilesToMap(supportedInstanceProfilesItem))
	}

	return supportedInstanceProfiles
}

func dataSourceDedicatedHostSupportedInstanceProfilesToMap(supportedInstanceProfilesItem vpcv1.InstanceProfileReference) (supportedInstanceProfilesMap map[string]interface{}) {
	supportedInstanceProfilesMap = map[string]interface{}{}

	if supportedInstanceProfilesItem.Href != nil {
		supportedInstanceProfilesMap["href"] = supportedInstanceProfilesItem.Href
	}
	if supportedInstanceProfilesItem.Name != nil {
		supportedInstanceProfilesMap["name"] = supportedInstanceProfilesItem.Name
	}

	return supportedInstanceProfilesMap
}

func dataSourceDedicatedHostFlattenVcpu(result vpcv1.Vcpu) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostVcpuToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostVcpuToMap(vcpuItem vpcv1.Vcpu) (vcpuMap map[string]interface{}) {
	vcpuMap = map[string]interface{}{}

	if vcpuItem.Architecture != nil {
		vcpuMap["architecture"] = vcpuItem.Architecture
	}
	if vcpuItem.Count != nil {
		vcpuMap["count"] = vcpuItem.Count
	}

	return vcpuMap
}

func dataSourceDedicatedHostFlattenDisks(result []vpcv1.DedicatedHostDisk) (disks []map[string]interface{}) {
	for _, disksItem := range result {
		disks = append(disks, dataSourceDedicatedHostDisksToMap(disksItem))
	}

	return disks
}

func dataSourceDedicatedHostDisksToMap(disksItem vpcv1.DedicatedHostDisk) (disksMap map[string]interface{}) {
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

func dataSourceDedicatedHostDisksInstanceDisksToMap(instanceDisksItem vpcv1.InstanceDiskReference) (instanceDisksMap map[string]interface{}) {
	instanceDisksMap = map[string]interface{}{}

	if instanceDisksItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceDedicatedHostInstanceDisksDeletedToMap(*instanceDisksItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		instanceDisksMap["deleted"] = deletedList
	}
	if instanceDisksItem.Href != nil {
		instanceDisksMap["href"] = instanceDisksItem.Href
	}
	if instanceDisksItem.ID != nil {
		instanceDisksMap["id"] = instanceDisksItem.ID
	}
	if instanceDisksItem.Name != nil {
		instanceDisksMap["name"] = instanceDisksItem.Name
	}
	if instanceDisksItem.ResourceType != nil {
		instanceDisksMap["resource_type"] = instanceDisksItem.ResourceType
	}

	return instanceDisksMap
}

func dataSourceDedicatedHostInstanceDisksDeletedToMap(deletedItem vpcv1.InstanceDiskReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
