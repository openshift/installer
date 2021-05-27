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

func dataSourceIbmIsDedicatedHostProfiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmIsDedicatedHostProfilesRead,

		Schema: map[string]*schema.Schema{
			"first": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link to the first page of resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for a page of resources.",
						},
					},
				},
			},
			"limit": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of resources that can be returned by the request.",
			},
			"next": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link to the next page of resources. This property is present for all pagesexcept the last page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for a page of resources.",
						},
					},
				},
			},
			"profiles": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of dedicated host profiles.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"class": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product class this dedicated host profile belongs to.",
						},
						"disks": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of the dedicated host profile's disks.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interface_type": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type for this profile field.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The interface of the disk for a dedicated host with this profileThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
												},
											},
										},
									},
									"quantity": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The number of disks of this type for a dedicated host with this profile.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type for this profile field.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The value for this profile field.",
												},
											},
										},
									},
									"size": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The size of the disk in GB (gigabytes).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type for this profile field.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The size of the disk in GB (gigabytes).",
												},
											},
										},
									},
									"supported_instance_interface_types": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type for this profile field.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The instance disk interfaces supported for a dedicated host with this profile.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
						"family": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product family this dedicated host profile belongs toThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this dedicated host.",
						},
						"memory": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The permitted values for this profile field.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this dedicated host profile.",
						},
						"socket_count": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The permitted values for this profile field.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"supported_instance_profiles": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Array of instance profiles that can be used by instances placed on dedicated hosts with this profile.",
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
						"vcpu_architecture": &schema.Schema{
							Type: schema.TypeList,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The VCPU architecture for a dedicated host with this profile.",
									},
								},
							},
						},
						"vcpu_count": &schema.Schema{
							Type: schema.TypeList,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The permitted values for this profile field.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
					},
				},
			},
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages.",
			},
		},
	}
}

func dataSourceIbmIsDedicatedHostProfilesRead(d *schema.ResourceData, meta interface{}) error {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return err
	}

	listDedicatedHostProfilesOptions := &vpcv1.ListDedicatedHostProfilesOptions{}

	dedicatedHostProfileCollection, response, err := vpcClient.ListDedicatedHostProfilesWithContext(context.TODO(), listDedicatedHostProfilesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListDedicatedHostProfilesWithContext failed %s\n%s", err, response)
		return err
	}

	if dedicatedHostProfileCollection.First != nil {
		err = d.Set("first", dataSourceDedicatedHostProfileCollectionFlattenFirst(*dedicatedHostProfileCollection.First))
		if err != nil {
			return fmt.Errorf("Error setting first %s", err)
		}
	}
	if err = d.Set("limit", dedicatedHostProfileCollection.Limit); err != nil {
		return fmt.Errorf("Error setting limit: %s", err)
	}

	if dedicatedHostProfileCollection.Next != nil {
		err = d.Set("next", dataSourceDedicatedHostProfileCollectionFlattenNext(*dedicatedHostProfileCollection.Next))
		if err != nil {
			return fmt.Errorf("Error setting next %s", err)
		}
	}

	if len(dedicatedHostProfileCollection.Profiles) != 0 {

		d.SetId(dataSourceIbmIsDedicatedHostProfilesID(d))

		if dedicatedHostProfileCollection.Profiles != nil {
			err = d.Set("profiles", dataSourceDedicatedHostProfileCollectionFlattenProfiles(dedicatedHostProfileCollection.Profiles))
			if err != nil {
				return fmt.Errorf("Error setting profiles %s", err)
			}
		}
		if err = d.Set("total_count", dedicatedHostProfileCollection.TotalCount); err != nil {
			return fmt.Errorf("Error setting total_count: %s", err)
		}
	}
	return nil
}

// dataSourceIbmIsDedicatedHostProfilesID returns a reasonable ID for the list.
func dataSourceIbmIsDedicatedHostProfilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceDedicatedHostProfileCollectionFlattenFirst(result vpcv1.DedicatedHostProfileCollectionFirst) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostProfileCollectionFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostProfileCollectionFirstToMap(firstItem vpcv1.DedicatedHostProfileCollectionFirst) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceDedicatedHostProfileCollectionFlattenNext(result vpcv1.DedicatedHostProfileCollectionNext) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostProfileCollectionNextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostProfileCollectionNextToMap(nextItem vpcv1.DedicatedHostProfileCollectionNext) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}

func dataSourceDedicatedHostProfileCollectionFlattenProfiles(result []vpcv1.DedicatedHostProfile) (profiles []map[string]interface{}) {
	for _, profilesItem := range result {
		profiles = append(profiles, dataSourceDedicatedHostProfileCollectionProfilesToMap(profilesItem))
	}

	return profiles
}

func dataSourceDedicatedHostProfileCollectionProfilesToMap(profilesItem vpcv1.DedicatedHostProfile) (profilesMap map[string]interface{}) {
	profilesMap = map[string]interface{}{}

	if profilesItem.Class != nil {
		profilesMap["class"] = profilesItem.Class
	}
	if profilesItem.Disks != nil {
		disksList := []map[string]interface{}{}
		for _, disksItem := range profilesItem.Disks {
			disksList = append(disksList, dataSourceDedicatedHostProfileCollectionProfilesDisksToMap(disksItem))
		}
		profilesMap["disks"] = disksList
	}
	if profilesItem.Family != nil {
		profilesMap["family"] = profilesItem.Family
	}
	if profilesItem.Href != nil {
		profilesMap["href"] = profilesItem.Href
	}
	if profilesItem.Memory != nil {
		memoryList := []map[string]interface{}{}
		memoryMap := dataSourceDedicatedHostProfileCollectionProfilesMemoryToMap(*profilesItem.Memory.(*vpcv1.DedicatedHostProfileMemory))
		memoryList = append(memoryList, memoryMap)
		profilesMap["memory"] = memoryList
	}
	if profilesItem.Name != nil {
		profilesMap["name"] = profilesItem.Name
	}
	if profilesItem.SocketCount != nil {
		socketCountList := []map[string]interface{}{}
		socketCountMap := dataSourceDedicatedHostProfileCollectionProfilesSocketCountToMap(*profilesItem.SocketCount.(*vpcv1.DedicatedHostProfileSocket))
		socketCountList = append(socketCountList, socketCountMap)
		profilesMap["socket_count"] = socketCountList
	}
	if profilesItem.SupportedInstanceProfiles != nil {
		supportedInstanceProfilesList := []map[string]interface{}{}
		for _, supportedInstanceProfilesItem := range profilesItem.SupportedInstanceProfiles {
			supportedInstanceProfilesList = append(supportedInstanceProfilesList, dataSourceDedicatedHostProfileCollectionProfilesSupportedInstanceProfilesToMap(supportedInstanceProfilesItem))
		}
		profilesMap["supported_instance_profiles"] = supportedInstanceProfilesList
	}
	if profilesItem.VcpuArchitecture != nil {
		vcpuArchitectureList := []map[string]interface{}{}
		vcpuArchitectureMap := dataSourceDedicatedHostProfileCollectionProfilesVcpuArchitectureToMap(*profilesItem.VcpuArchitecture)
		vcpuArchitectureList = append(vcpuArchitectureList, vcpuArchitectureMap)
		profilesMap["vcpu_architecture"] = vcpuArchitectureList
	}
	if profilesItem.VcpuCount != nil {
		vcpuCountList := []map[string]interface{}{}
		vcpuCountMap := dataSourceDedicatedHostProfileCollectionProfilesVcpuCountToMap(*profilesItem.VcpuCount.(*vpcv1.DedicatedHostProfileVcpu))
		vcpuCountList = append(vcpuCountList, vcpuCountMap)
		profilesMap["vcpu_count"] = vcpuCountList
	}

	return profilesMap
}

func dataSourceDedicatedHostProfileCollectionProfilesMemoryToMap(memoryItem vpcv1.DedicatedHostProfileMemory) (memoryMap map[string]interface{}) {
	memoryMap = map[string]interface{}{}

	if memoryItem.Type != nil {
		memoryMap["type"] = memoryItem.Type
	}
	if memoryItem.Value != nil {
		memoryMap["value"] = memoryItem.Value
	}
	if memoryItem.Default != nil {
		memoryMap["default"] = memoryItem.Default
	}
	if memoryItem.Max != nil {
		memoryMap["max"] = memoryItem.Max
	}
	if memoryItem.Min != nil {
		memoryMap["min"] = memoryItem.Min
	}
	if memoryItem.Step != nil {
		memoryMap["step"] = memoryItem.Step
	}
	if memoryItem.Values != nil {
		memoryMap["values"] = memoryItem.Values
	}

	return memoryMap
}

func dataSourceDedicatedHostProfileCollectionProfilesSocketCountToMap(socketCountItem vpcv1.DedicatedHostProfileSocket) (socketCountMap map[string]interface{}) {
	socketCountMap = map[string]interface{}{}

	if socketCountItem.Type != nil {
		socketCountMap["type"] = socketCountItem.Type
	}
	if socketCountItem.Value != nil {
		socketCountMap["value"] = socketCountItem.Value
	}
	if socketCountItem.Default != nil {
		socketCountMap["default"] = socketCountItem.Default
	}
	if socketCountItem.Max != nil {
		socketCountMap["max"] = socketCountItem.Max
	}
	if socketCountItem.Min != nil {
		socketCountMap["min"] = socketCountItem.Min
	}
	if socketCountItem.Step != nil {
		socketCountMap["step"] = socketCountItem.Step
	}
	if socketCountItem.Values != nil {
		socketCountMap["values"] = socketCountItem.Values
	}

	return socketCountMap
}

func dataSourceDedicatedHostProfileCollectionProfilesSupportedInstanceProfilesToMap(supportedInstanceProfilesItem vpcv1.InstanceProfileReference) (supportedInstanceProfilesMap map[string]interface{}) {
	supportedInstanceProfilesMap = map[string]interface{}{}

	if supportedInstanceProfilesItem.Href != nil {
		supportedInstanceProfilesMap["href"] = supportedInstanceProfilesItem.Href
	}
	if supportedInstanceProfilesItem.Name != nil {
		supportedInstanceProfilesMap["name"] = supportedInstanceProfilesItem.Name
	}

	return supportedInstanceProfilesMap
}

func dataSourceDedicatedHostProfileCollectionProfilesVcpuArchitectureToMap(vcpuArchitectureItem vpcv1.DedicatedHostProfileVcpuArchitecture) (vcpuArchitectureMap map[string]interface{}) {
	vcpuArchitectureMap = map[string]interface{}{}

	if vcpuArchitectureItem.Type != nil {
		vcpuArchitectureMap["type"] = vcpuArchitectureItem.Type
	}
	if vcpuArchitectureItem.Value != nil {
		vcpuArchitectureMap["value"] = vcpuArchitectureItem.Value
	}

	return vcpuArchitectureMap
}

func dataSourceDedicatedHostProfileCollectionProfilesVcpuCountToMap(vcpuCountItem vpcv1.DedicatedHostProfileVcpu) (vcpuCountMap map[string]interface{}) {
	vcpuCountMap = map[string]interface{}{}

	if vcpuCountItem.Type != nil {
		vcpuCountMap["type"] = vcpuCountItem.Type
	}
	if vcpuCountItem.Value != nil {
		vcpuCountMap["value"] = vcpuCountItem.Value
	}
	if vcpuCountItem.Default != nil {
		vcpuCountMap["default"] = vcpuCountItem.Default
	}
	if vcpuCountItem.Max != nil {
		vcpuCountMap["max"] = vcpuCountItem.Max
	}
	if vcpuCountItem.Min != nil {
		vcpuCountMap["min"] = vcpuCountItem.Min
	}
	if vcpuCountItem.Step != nil {
		vcpuCountMap["step"] = vcpuCountItem.Step
	}
	if vcpuCountItem.Values != nil {
		vcpuCountMap["values"] = vcpuCountItem.Values
	}

	return vcpuCountMap
}

func dataSourceDedicatedHostProfileCollectionProfilesDisksToMap(disksItem vpcv1.DedicatedHostProfileDisk) (disksMap map[string]interface{}) {
	disksMap = map[string]interface{}{}

	if disksItem.InterfaceType != nil {
		interfaceTypeList := []map[string]interface{}{}
		interfaceTypeMap := dataSourceDedicatedHostProfileCollectionDisksInterfaceTypeToMap(*disksItem.InterfaceType)
		interfaceTypeList = append(interfaceTypeList, interfaceTypeMap)
		disksMap["interface_type"] = interfaceTypeList
	}
	if disksItem.Quantity != nil {
		quantityList := []map[string]interface{}{}
		quantityMap := dataSourceDedicatedHostProfileCollectionDisksQuantityToMap(*disksItem.Quantity)
		quantityList = append(quantityList, quantityMap)
		disksMap["quantity"] = quantityList
	}
	if disksItem.Size != nil {
		sizeList := []map[string]interface{}{}
		sizeMap := dataSourceDedicatedHostProfileCollectionDisksSizeToMap(*disksItem.Size)
		sizeList = append(sizeList, sizeMap)
		disksMap["size"] = sizeList
	}
	if disksItem.SupportedInstanceInterfaceTypes != nil {
		supportedInstanceInterfaceTypesList := []map[string]interface{}{}
		supportedInstanceInterfaceTypesMap := dataSourceDedicatedHostProfileCollectionDisksSupportedInstanceInterfaceTypesToMap(*disksItem.SupportedInstanceInterfaceTypes)
		supportedInstanceInterfaceTypesList = append(supportedInstanceInterfaceTypesList, supportedInstanceInterfaceTypesMap)
		disksMap["supported_instance_interface_types"] = supportedInstanceInterfaceTypesList
	}

	return disksMap
}

func dataSourceDedicatedHostProfileCollectionDisksInterfaceTypeToMap(interfaceTypeItem vpcv1.DedicatedHostProfileDiskInterface) (interfaceTypeMap map[string]interface{}) {
	interfaceTypeMap = map[string]interface{}{}

	if interfaceTypeItem.Type != nil {
		interfaceTypeMap["type"] = interfaceTypeItem.Type
	}
	if interfaceTypeItem.Value != nil {
		interfaceTypeMap["value"] = interfaceTypeItem.Value
	}

	return interfaceTypeMap
}

func dataSourceDedicatedHostProfileCollectionDisksQuantityToMap(quantityItem vpcv1.DedicatedHostProfileDiskQuantity) (quantityMap map[string]interface{}) {
	quantityMap = map[string]interface{}{}

	if quantityItem.Type != nil {
		quantityMap["type"] = quantityItem.Type
	}
	if quantityItem.Value != nil {
		quantityMap["value"] = quantityItem.Value
	}

	return quantityMap
}

func dataSourceDedicatedHostProfileCollectionDisksSizeToMap(sizeItem vpcv1.DedicatedHostProfileDiskSize) (sizeMap map[string]interface{}) {
	sizeMap = map[string]interface{}{}

	if sizeItem.Type != nil {
		sizeMap["type"] = sizeItem.Type
	}
	if sizeItem.Value != nil {
		sizeMap["value"] = sizeItem.Value
	}

	return sizeMap
}

func dataSourceDedicatedHostProfileCollectionDisksSupportedInstanceInterfaceTypesToMap(supportedInstanceInterfaceTypesItem vpcv1.DedicatedHostProfileDiskSupportedInterfaces) (supportedInstanceInterfaceTypesMap map[string]interface{}) {
	supportedInstanceInterfaceTypesMap = map[string]interface{}{}

	if supportedInstanceInterfaceTypesItem.Type != nil {
		supportedInstanceInterfaceTypesMap["type"] = supportedInstanceInterfaceTypesItem.Type
	}
	if supportedInstanceInterfaceTypesItem.Value != nil {
		supportedInstanceInterfaceTypesMap["value"] = supportedInstanceInterfaceTypesItem.Value
	}

	return supportedInstanceInterfaceTypesMap
}
