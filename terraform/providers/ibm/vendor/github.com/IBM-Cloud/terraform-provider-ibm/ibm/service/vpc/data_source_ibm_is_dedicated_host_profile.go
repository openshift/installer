// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIbmIsDedicatedHostProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsDedicatedHostProfileRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The globally unique name for this virtual server instance profile.",
			},
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
							Description: "The VCPU architecture for a dedicated host with this profile.",
						},
					},
				},
			},
			"vcpu_count": {
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
			"vcpu_manufacturer": &schema.Schema{
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
							Description: "The VCPU manufacturer for a dedicated host with this profile.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmIsDedicatedHostProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	getDedicatedHostProfileOptions := &vpcv1.GetDedicatedHostProfileOptions{
		Name: &name,
	}
	dedicatedHostProfile, response, err := vpcClient.GetDedicatedHostProfileWithContext(context, getDedicatedHostProfileOptions)
	if err != nil {
		log.Printf("[DEBUG] ListDedicatedHostProfilesWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	if dedicatedHostProfile == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] No Dedicated Host Profile found with name %s", name))
	}
	d.SetId(dataSourceIbmIsDedicatedHostProfileID(d))

	if err = d.Set("class", dedicatedHostProfile.Class); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting class: %s", err))
	}

	if dedicatedHostProfile.Disks != nil {
		err = d.Set("disks", dataSourceDedicatedHostProfileFlattenDisks(dedicatedHostProfile.Disks))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting disks %s", err))
		}
	}
	if dedicatedHostProfile.Status != nil {
		d.Set("status", dedicatedHostProfile.Status)
	}
	if err = d.Set("family", dedicatedHostProfile.Family); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting family: %s", err))
	}
	if err = d.Set("href", dedicatedHostProfile.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}

	if dedicatedHostProfile.Memory != nil {
		err = d.Set("memory", dataSourceDedicatedHostProfileFlattenMemory(*dedicatedHostProfile.Memory.(*vpcv1.DedicatedHostProfileMemory)))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting memory %s", err))
		}
	}

	if dedicatedHostProfile.SocketCount != nil {
		err = d.Set("socket_count", dataSourceDedicatedHostProfileFlattenSocketCount(*dedicatedHostProfile.SocketCount.(*vpcv1.DedicatedHostProfileSocket)))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting socket_count %s", err))
		}
	}

	if dedicatedHostProfile.SupportedInstanceProfiles != nil {
		err = d.Set("supported_instance_profiles", dataSourceDedicatedHostProfileFlattenSupportedInstanceProfiles(dedicatedHostProfile.SupportedInstanceProfiles))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting supported_instance_profiles %s", err))
		}
	}

	if dedicatedHostProfile.VcpuArchitecture != nil {
		err = d.Set("vcpu_architecture", dataSourceDedicatedHostProfileFlattenVcpuArchitecture(*dedicatedHostProfile.VcpuArchitecture))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting vcpu_architecture %s", err))
		}
	}

	if dedicatedHostProfile.VcpuCount != nil {
		err = d.Set("vcpu_count", dataSourceDedicatedHostProfileFlattenVcpuCount(*dedicatedHostProfile.VcpuCount.(*vpcv1.DedicatedHostProfileVcpu)))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting vcpu_count %s", err))
		}
	}

	// Changes for the AMD Support, manufacturer information.
	if dedicatedHostProfile.VcpuManufacturer != nil {
		err = d.Set("vcpu_manufacturer", dataSourceDedicatedHostProfileFlattenVcpuManufacturer(*dedicatedHostProfile.VcpuManufacturer))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting vcpu_architecture %s", err))
		}
	}

	return nil

}

// dataSourceIbmIsDedicatedHostProfileID returns a reasonable ID for the list.
func dataSourceIbmIsDedicatedHostProfileID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceDedicatedHostProfileFlattenMemory(result vpcv1.DedicatedHostProfileMemory) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostProfileMemoryToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostProfileMemoryToMap(memoryItem vpcv1.DedicatedHostProfileMemory) (memoryMap map[string]interface{}) {
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

func dataSourceDedicatedHostProfileFlattenSocketCount(result vpcv1.DedicatedHostProfileSocket) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostProfileSocketCountToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostProfileSocketCountToMap(socketCountItem vpcv1.DedicatedHostProfileSocket) (socketCountMap map[string]interface{}) {
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

func dataSourceDedicatedHostProfileFlattenSupportedInstanceProfiles(result []vpcv1.InstanceProfileReference) (supportedInstanceProfiles []map[string]interface{}) {
	for _, supportedInstanceProfilesItem := range result {
		supportedInstanceProfiles = append(supportedInstanceProfiles, dataSourceDedicatedHostProfileSupportedInstanceProfilesToMap(supportedInstanceProfilesItem))
	}

	return supportedInstanceProfiles
}

func dataSourceDedicatedHostProfileSupportedInstanceProfilesToMap(supportedInstanceProfilesItem vpcv1.InstanceProfileReference) (supportedInstanceProfilesMap map[string]interface{}) {
	supportedInstanceProfilesMap = map[string]interface{}{}

	if supportedInstanceProfilesItem.Href != nil {
		supportedInstanceProfilesMap["href"] = supportedInstanceProfilesItem.Href
	}
	if supportedInstanceProfilesItem.Name != nil {
		supportedInstanceProfilesMap["name"] = supportedInstanceProfilesItem.Name
	}

	return supportedInstanceProfilesMap
}

func dataSourceDedicatedHostProfileFlattenVcpuArchitecture(result vpcv1.DedicatedHostProfileVcpuArchitecture) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostProfileVcpuArchitectureToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostProfileVcpuArchitectureToMap(vcpuArchitectureItem vpcv1.DedicatedHostProfileVcpuArchitecture) (vcpuArchitectureMap map[string]interface{}) {
	vcpuArchitectureMap = map[string]interface{}{}

	if vcpuArchitectureItem.Type != nil {
		vcpuArchitectureMap["type"] = vcpuArchitectureItem.Type
	}
	if vcpuArchitectureItem.Value != nil {
		vcpuArchitectureMap["value"] = vcpuArchitectureItem.Value
	}

	return vcpuArchitectureMap
}

// Changes for AMD Support, manufacturer details.
func dataSourceDedicatedHostProfileFlattenVcpuManufacturer(result vpcv1.DedicatedHostProfileVcpuManufacturer) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostProfileVcpuManufacturerToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

// AMD Support for manufacturer
func dataSourceDedicatedHostProfileVcpuManufacturerToMap(vcpuManufacturerItem vpcv1.DedicatedHostProfileVcpuManufacturer) (vcpuManufacturerMap map[string]interface{}) {
	vcpuManufacturerMap = map[string]interface{}{}

	if vcpuManufacturerItem.Type != nil {
		vcpuManufacturerMap["type"] = vcpuManufacturerItem.Type
	}
	if vcpuManufacturerItem.Value != nil {
		vcpuManufacturerMap["value"] = vcpuManufacturerItem.Value
	}

	return vcpuManufacturerMap
}

func dataSourceDedicatedHostProfileFlattenVcpuCount(result vpcv1.DedicatedHostProfileVcpu) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostProfileVcpuCountToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostProfileVcpuCountToMap(vcpuCountItem vpcv1.DedicatedHostProfileVcpu) (vcpuCountMap map[string]interface{}) {
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

func dataSourceDedicatedHostProfileFlattenDisks(result []vpcv1.DedicatedHostProfileDisk) (disks []map[string]interface{}) {
	for _, disksItem := range result {
		disks = append(disks, dataSourceDedicatedHostProfileDisksToMap(disksItem))
	}

	return disks
}

func dataSourceDedicatedHostProfileDisksToMap(disksItem vpcv1.DedicatedHostProfileDisk) (disksMap map[string]interface{}) {
	disksMap = map[string]interface{}{}

	if disksItem.InterfaceType != nil {
		interfaceTypeList := []map[string]interface{}{}
		interfaceTypeMap := dataSourceDedicatedHostProfileDisksInterfaceTypeToMap(*disksItem.InterfaceType)
		interfaceTypeList = append(interfaceTypeList, interfaceTypeMap)
		disksMap["interface_type"] = interfaceTypeList
	}
	if disksItem.Quantity != nil {
		quantityList := []map[string]interface{}{}
		quantityMap := dataSourceDedicatedHostProfileDisksQuantityToMap(*disksItem.Quantity)
		quantityList = append(quantityList, quantityMap)
		disksMap["quantity"] = quantityList
	}
	if disksItem.Size != nil {
		sizeList := []map[string]interface{}{}
		sizeMap := dataSourceDedicatedHostProfileDisksSizeToMap(*disksItem.Size)
		sizeList = append(sizeList, sizeMap)
		disksMap["size"] = sizeList
	}
	if disksItem.SupportedInstanceInterfaceTypes != nil {
		supportedInstanceInterfaceTypesList := []map[string]interface{}{}
		supportedInstanceInterfaceTypesMap := dataSourceDedicatedHostProfileDisksSupportedInstanceInterfaceTypesToMap(*disksItem.SupportedInstanceInterfaceTypes)
		supportedInstanceInterfaceTypesList = append(supportedInstanceInterfaceTypesList, supportedInstanceInterfaceTypesMap)
		disksMap["supported_instance_interface_types"] = supportedInstanceInterfaceTypesList
	}

	return disksMap
}

func dataSourceDedicatedHostProfileDisksInterfaceTypeToMap(interfaceTypeItem vpcv1.DedicatedHostProfileDiskInterface) (interfaceTypeMap map[string]interface{}) {
	interfaceTypeMap = map[string]interface{}{}

	if interfaceTypeItem.Type != nil {
		interfaceTypeMap["type"] = interfaceTypeItem.Type
	}
	if interfaceTypeItem.Value != nil {
		interfaceTypeMap["value"] = interfaceTypeItem.Value
	}

	return interfaceTypeMap
}

func dataSourceDedicatedHostProfileDisksQuantityToMap(quantityItem vpcv1.DedicatedHostProfileDiskQuantity) (quantityMap map[string]interface{}) {
	quantityMap = map[string]interface{}{}

	if quantityItem.Type != nil {
		quantityMap["type"] = quantityItem.Type
	}
	if quantityItem.Value != nil {
		quantityMap["value"] = quantityItem.Value
	}

	return quantityMap
}

func dataSourceDedicatedHostProfileDisksSizeToMap(sizeItem vpcv1.DedicatedHostProfileDiskSize) (sizeMap map[string]interface{}) {
	sizeMap = map[string]interface{}{}

	if sizeItem.Type != nil {
		sizeMap["type"] = sizeItem.Type
	}
	if sizeItem.Value != nil {
		sizeMap["value"] = sizeItem.Value
	}

	return sizeMap
}

func dataSourceDedicatedHostProfileDisksSupportedInstanceInterfaceTypesToMap(supportedInstanceInterfaceTypesItem vpcv1.DedicatedHostProfileDiskSupportedInterfaces) (supportedInstanceInterfaceTypesMap map[string]interface{}) {
	supportedInstanceInterfaceTypesMap = map[string]interface{}{}

	if supportedInstanceInterfaceTypesItem.Type != nil {
		supportedInstanceInterfaceTypesMap["type"] = supportedInstanceInterfaceTypesItem.Type
	}
	if supportedInstanceInterfaceTypesItem.Value != nil {
		supportedInstanceInterfaceTypesMap["value"] = supportedInstanceInterfaceTypesItem.Value
	}

	return supportedInstanceInterfaceTypesMap
}
