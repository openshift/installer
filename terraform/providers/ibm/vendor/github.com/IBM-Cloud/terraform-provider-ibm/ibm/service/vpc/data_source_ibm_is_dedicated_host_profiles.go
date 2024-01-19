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

func DataSourceIbmIsDedicatedHostProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsDedicatedHostProfilesRead,

		Schema: map[string]*schema.Schema{
			"profiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of dedicated host profiles.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product class this dedicated host profile belongs to.",
						},
						"disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of the dedicated host profile's disks.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interface_type": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type for this profile field.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The interface of the disk for a dedicated host with this profileThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
												},
											},
										},
									},
									"quantity": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The number of disks of this type for a dedicated host with this profile.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type for this profile field.",
												},
												"value": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The value for this profile field.",
												},
											},
										},
									},
									"size": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The size of the disk in GB (gigabytes).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type for this profile field.",
												},
												"value": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The size of the disk in GB (gigabytes).",
												},
											},
										},
									},
									"supported_instance_interface_types": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type for this profile field.",
												},
												"value": {
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
						"family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product family this dedicated host profile belongs toThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this dedicated host.",
						},
						"memory": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": {
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
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this dedicated host profile.",
						},
						"socket_count": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": {
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
						"supported_instance_profiles": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Array of instance profiles that can be used by instances placed on dedicated hosts with this profile.",
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
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the dedicated host profile.",
						},
						"vcpu_architecture": {
							Type: schema.TypeList,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The VCPU architecture for a dedicated host with this profile.",
									},
								},
							},
						},
						"vcpu_count": {
							Type: schema.TypeList,

							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": {
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
						"vcpu_manufacturer": &schema.Schema{
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
										Description: "TThe VCPU manufacturer for a dedicated host with this profile.",
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

func dataSourceIbmIsDedicatedHostProfilesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listDedicatedHostProfilesOptions := &vpcv1.ListDedicatedHostProfilesOptions{}

	start := ""
	allrecs := []vpcv1.DedicatedHostProfile{}
	for {
		if start != "" {
			listDedicatedHostProfilesOptions.Start = &start
		}
		dedicatedHostProfileCollection, response, err := vpcClient.ListDedicatedHostProfilesWithContext(context, listDedicatedHostProfilesOptions)
		if err != nil {
			log.Printf("[DEBUG] ListDedicatedHostProfilesWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		start = flex.GetNext(dedicatedHostProfileCollection.Next)
		allrecs = append(allrecs, dedicatedHostProfileCollection.Profiles...)
		if start == "" {
			break
		}
	}

	if len(allrecs) > 0 {

		d.SetId(dataSourceIbmIsDedicatedHostProfilesID(d))

		err = d.Set("profiles", dataSourceDedicatedHostProfileCollectionFlattenProfiles(allrecs))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting profiles %s", err))
		}

		if err = d.Set("total_count", len(allrecs)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting total_count: %s", err))
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
	if profilesItem.Status != nil {
		profilesMap["status"] = *profilesItem.Status
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
	// AMD Support, changes for manufacturer details.
	if profilesItem.VcpuManufacturer != nil {
		vcpuManufacturerList := []map[string]interface{}{}
		vcpuManufacturerMap := dataSourceDedicatedHostProfileCollectionProfilesVcpuManufacturerToMap(*profilesItem.VcpuManufacturer)
		vcpuManufacturerList = append(vcpuManufacturerList, vcpuManufacturerMap)
		profilesMap["vcpu_manufacturer"] = vcpuManufacturerList
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

// AMD Changes, manufacturer details added.
func dataSourceDedicatedHostProfileCollectionProfilesVcpuManufacturerToMap(vcpuManufacturerItem vpcv1.DedicatedHostProfileVcpuManufacturer) (vcpuManufacturerMap map[string]interface{}) {
	vcpuManufacturerMap = map[string]interface{}{}

	if vcpuManufacturerItem.Type != nil {
		vcpuManufacturerMap["type"] = vcpuManufacturerItem.Type
	}
	if vcpuManufacturerItem.Value != nil {
		vcpuManufacturerMap["value"] = vcpuManufacturerItem.Value
	}

	return vcpuManufacturerMap
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
