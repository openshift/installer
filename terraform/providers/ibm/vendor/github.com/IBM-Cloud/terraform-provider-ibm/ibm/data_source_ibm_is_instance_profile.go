// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isInstanceProfileName         = "name"
	isInstanceProfileFamily       = "family"
	isInstanceProfileArchitecture = "architecture"
)

func dataSourceIBMISInstanceProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstanceProfileRead,

		Schema: map[string]*schema.Schema{

			isInstanceProfileName: {
				Type:     schema.TypeString,
				Required: true,
			},

			isInstanceProfileFamily: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product family this virtual server instance profile belongs to.",
			},

			isInstanceProfileArchitecture: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default OS architecture for an instance with this profile.",
			},

			"architecture_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type for the OS architecture.",
			},

			"architecture_values": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The supported OS architecture(s) for an instance with this profile.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"bandwidth": &schema.Schema{
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
			"gpu_count": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "GPU count of this profile",
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
			"gpu_manufacturer": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "GPU manufacturer of this profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"values": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The possible GPU manufacturer(s) for an instance with this profile",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"gpu_memory": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "GPU memory of this profile",
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
			"gpu_model": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "GPU model of this profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"values": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The possible GPU model(s) for an instance with this profile",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"total_volume_bandwidth": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The amount of bandwidth (in megabits per second) allocated exclusively to instance storage volumes. An increase in this value will result in a corresponding decrease to total_network_bandwidth.",
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
			"disks": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of the instance profile's disks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quantity": &schema.Schema{
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
						"size": &schema.Schema{
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
						"supported_interface_types": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The disk interface used for attaching the disk.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The supported disk interfaces used for attaching the disk.",
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
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this virtual server instance profile.",
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
			"port_speed": &schema.Schema{
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
					},
				},
			},
			"vcpu_architecture": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default VCPU architecture for an instance with this profile.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VCPU architecture for an instance with this profile.",
						},
					},
				},
			},
			"vcpu_count": &schema.Schema{
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
		},
	}
}

func dataSourceIBMISInstanceProfileRead(d *schema.ResourceData, meta interface{}) error {

	name := d.Get(isInstanceProfileName).(string)
	err := instanceProfileGet(d, meta, name)
	if err != nil {
		return err
	}
	return nil
}

func instanceProfileGet(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getInstanceProfileOptions := &vpcv1.GetInstanceProfileOptions{
		Name: &name,
	}
	profile, _, err := sess.GetInstanceProfile(getInstanceProfileOptions)
	if err != nil {
		return err
	}
	// For lack of anything better, compose our id from profile name.
	d.SetId(*profile.Name)
	d.Set(isInstanceProfileName, *profile.Name)
	d.Set(isInstanceProfileFamily, *profile.Family)
	if profile.OsArchitecture != nil {
		if profile.OsArchitecture.Default != nil {
			d.Set(isInstanceProfileArchitecture, *profile.OsArchitecture.Default)
		}
		if profile.OsArchitecture.Type != nil {
			d.Set("architecture_type", *profile.OsArchitecture.Type)
		}
		if profile.OsArchitecture.Values != nil {
			d.Set("architecture_values", *&profile.OsArchitecture.Values)
		}

	}
	if profile.Bandwidth != nil {
		err = d.Set("bandwidth", dataSourceInstanceProfileFlattenBandwidth(*profile.Bandwidth.(*vpcv1.InstanceProfileBandwidth)))
		if err != nil {
			return err
		}
	}
	if profile.GpuCount != nil {
		err = d.Set("gpu_count", dataSourceInstanceProfileFlattenGPUCount(*profile.GpuCount.(*vpcv1.InstanceProfileGpu)))
		if err != nil {
			return err
		}
	}
	if profile.GpuMemory != nil {
		err = d.Set("gpu_memory", dataSourceInstanceProfileFlattenGPUMemory(*profile.GpuMemory.(*vpcv1.InstanceProfileGpuMemory)))
		if err != nil {
			return err
		}
	}
	if profile.GpuManufacturer != nil {
		err = d.Set("gpu_manufacturer", dataSourceInstanceProfileFlattenGPUManufacturer(*profile.GpuManufacturer))
		if err != nil {
			return err
		}
	}
	if profile.GpuModel != nil {
		err = d.Set("gpu_model", dataSourceInstanceProfileFlattenGPUModel(*profile.GpuModel))
		if err != nil {
			return err
		}
	}
	if profile.TotalVolumeBandwidth != nil {
		err = d.Set("total_volume_bandwidth", dataSourceInstanceProfileFlattenTotalVolumeBandwidth(*profile.TotalVolumeBandwidth.(*vpcv1.InstanceProfileVolumeBandwidth)))
		if err != nil {
			return err
		}
	}
	if profile.Disks != nil {
		err = d.Set("disks", dataSourceInstanceProfileFlattenDisks(profile.Disks))
		if err != nil {
			return err
		}
	}
	if err = d.Set("href", profile.Href); err != nil {
		return err
	}

	if profile.Memory != nil {
		err = d.Set("memory", dataSourceInstanceProfileFlattenMemory(*profile.Memory.(*vpcv1.InstanceProfileMemory)))
		if err != nil {
			return err
		}
	}
	if profile.PortSpeed != nil {
		err = d.Set("port_speed", dataSourceInstanceProfileFlattenPortSpeed(*profile.PortSpeed.(*vpcv1.InstanceProfilePortSpeed)))
		if err != nil {
			return err
		}
	}

	if profile.VcpuArchitecture != nil {
		err = d.Set("vcpu_architecture", dataSourceInstanceProfileFlattenVcpuArchitecture(*profile.VcpuArchitecture))
		if err != nil {
			return err
		}
	}

	if profile.VcpuCount != nil {
		err = d.Set("vcpu_count", dataSourceInstanceProfileFlattenVcpuCount(*profile.VcpuCount.(*vpcv1.InstanceProfileVcpu)))
		if err != nil {
			return err
		}
	}
	return nil
}

func dataSourceInstanceProfileFlattenBandwidth(result vpcv1.InstanceProfileBandwidth) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileBandwidthToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileBandwidthToMap(bandwidthItem vpcv1.InstanceProfileBandwidth) (bandwidthMap map[string]interface{}) {
	bandwidthMap = map[string]interface{}{}

	if bandwidthItem.Type != nil {
		bandwidthMap["type"] = bandwidthItem.Type
	}
	if bandwidthItem.Value != nil {
		bandwidthMap["value"] = bandwidthItem.Value
	}
	if bandwidthItem.Default != nil {
		bandwidthMap["default"] = bandwidthItem.Default
	}
	if bandwidthItem.Max != nil {
		bandwidthMap["max"] = bandwidthItem.Max
	}
	if bandwidthItem.Min != nil {
		bandwidthMap["min"] = bandwidthItem.Min
	}
	if bandwidthItem.Step != nil {
		bandwidthMap["step"] = bandwidthItem.Step
	}
	if bandwidthItem.Values != nil {
		bandwidthMap["values"] = bandwidthItem.Values
	}

	return bandwidthMap
}

func dataSourceInstanceProfileFlattenGPUCount(result vpcv1.InstanceProfileGpu) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileGPUCountToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileGPUCountToMap(bandwidthItem vpcv1.InstanceProfileGpu) (gpuCountMap map[string]interface{}) {
	gpuCountMap = map[string]interface{}{}

	if bandwidthItem.Type != nil {
		gpuCountMap["type"] = bandwidthItem.Type
	}
	if bandwidthItem.Value != nil {
		gpuCountMap["value"] = bandwidthItem.Value
	}
	if bandwidthItem.Default != nil {
		gpuCountMap["default"] = bandwidthItem.Default
	}
	if bandwidthItem.Max != nil {
		gpuCountMap["max"] = bandwidthItem.Max
	}
	if bandwidthItem.Min != nil {
		gpuCountMap["min"] = bandwidthItem.Min
	}
	if bandwidthItem.Step != nil {
		gpuCountMap["step"] = bandwidthItem.Step
	}
	if bandwidthItem.Values != nil {
		gpuCountMap["values"] = bandwidthItem.Values
	}

	return gpuCountMap
}

func dataSourceInstanceProfileFlattenGPUManufacturer(result vpcv1.InstanceProfileGpuManufacturer) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileGPUManufacturerToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileGPUManufacturerToMap(bandwidthItem vpcv1.InstanceProfileGpuManufacturer) (gpuManufactrerMap map[string]interface{}) {
	gpuManufactrerMap = map[string]interface{}{}

	if bandwidthItem.Type != nil {
		gpuManufactrerMap["type"] = bandwidthItem.Type
	}
	if bandwidthItem.Values != nil {
		gpuManufactrerMap["values"] = bandwidthItem.Values
	}

	return gpuManufactrerMap
}

func dataSourceInstanceProfileFlattenGPUModel(result vpcv1.InstanceProfileGpuModel) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileGPUModelToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileGPUModelToMap(bandwidthItem vpcv1.InstanceProfileGpuModel) (gpuModel map[string]interface{}) {
	gpuModelMap := map[string]interface{}{}

	if bandwidthItem.Type != nil {
		gpuModelMap["type"] = bandwidthItem.Type
	}
	if bandwidthItem.Values != nil {
		gpuModelMap["values"] = bandwidthItem.Values
	}

	return gpuModelMap
}

func dataSourceInstanceProfileFlattenGPUMemory(result vpcv1.InstanceProfileGpuMemory) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileGPUMemoryToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileGPUMemoryToMap(bandwidthItem vpcv1.InstanceProfileGpuMemory) (gpuMemoryMap map[string]interface{}) {
	gpuMemoryMap = map[string]interface{}{}

	if bandwidthItem.Type != nil {
		gpuMemoryMap["type"] = bandwidthItem.Type
	}
	if bandwidthItem.Value != nil {
		gpuMemoryMap["value"] = bandwidthItem.Value
	}
	if bandwidthItem.Default != nil {
		gpuMemoryMap["default"] = bandwidthItem.Default
	}
	if bandwidthItem.Max != nil {
		gpuMemoryMap["max"] = bandwidthItem.Max
	}
	if bandwidthItem.Min != nil {
		gpuMemoryMap["min"] = bandwidthItem.Min
	}
	if bandwidthItem.Step != nil {
		gpuMemoryMap["step"] = bandwidthItem.Step
	}
	if bandwidthItem.Values != nil {
		gpuMemoryMap["values"] = bandwidthItem.Values
	}

	return gpuMemoryMap
}

func dataSourceInstanceProfileFlattenMemory(result vpcv1.InstanceProfileMemory) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileMemoryToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileMemoryToMap(memoryItem vpcv1.InstanceProfileMemory) (memoryMap map[string]interface{}) {
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

func dataSourceInstanceProfileFlattenPortSpeed(result vpcv1.InstanceProfilePortSpeed) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfilePortSpeedToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfilePortSpeedToMap(portSpeedItem vpcv1.InstanceProfilePortSpeed) (portSpeedMap map[string]interface{}) {
	portSpeedMap = map[string]interface{}{}

	if portSpeedItem.Type != nil {
		portSpeedMap["type"] = portSpeedItem.Type
	}
	if portSpeedItem.Value != nil {
		portSpeedMap["value"] = portSpeedItem.Value
	}

	return portSpeedMap
}

func dataSourceInstanceProfileFlattenVcpuArchitecture(result vpcv1.InstanceProfileVcpuArchitecture) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileVcpuArchitectureToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileVcpuArchitectureToMap(vcpuArchitectureItem vpcv1.InstanceProfileVcpuArchitecture) (vcpuArchitectureMap map[string]interface{}) {
	vcpuArchitectureMap = map[string]interface{}{}

	if vcpuArchitectureItem.Default != nil {
		vcpuArchitectureMap["default"] = vcpuArchitectureItem.Default
	}
	if vcpuArchitectureItem.Type != nil {
		vcpuArchitectureMap["type"] = vcpuArchitectureItem.Type
	}
	if vcpuArchitectureItem.Value != nil {
		vcpuArchitectureMap["value"] = vcpuArchitectureItem.Value
	}

	return vcpuArchitectureMap
}

func dataSourceInstanceProfileFlattenVcpuCount(result vpcv1.InstanceProfileVcpu) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileVcpuCountToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileVcpuCountToMap(vcpuCountItem vpcv1.InstanceProfileVcpu) (vcpuCountMap map[string]interface{}) {
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

func dataSourceInstanceProfileFlattenDisks(result []vpcv1.InstanceProfileDisk) (disks []map[string]interface{}) {
	for _, disksItem := range result {
		disks = append(disks, dataSourceInstanceProfileDisksToMap(disksItem))
	}

	return disks
}

func dataSourceInstanceProfileDisksToMap(disksItem vpcv1.InstanceProfileDisk) (disksMap map[string]interface{}) {
	disksMap = map[string]interface{}{}

	if disksItem.Quantity != nil {
		quantityList := []map[string]interface{}{}
		quantityMap := dataSourceInstanceProfileDisksQuantityToMap(*disksItem.Quantity.(*vpcv1.InstanceProfileDiskQuantity))
		quantityList = append(quantityList, quantityMap)
		disksMap["quantity"] = quantityList
	}
	if disksItem.Size != nil {
		sizeList := []map[string]interface{}{}
		sizeMap := dataSourceInstanceProfileDisksSizeToMap(*disksItem.Size.(*vpcv1.InstanceProfileDiskSize))
		sizeList = append(sizeList, sizeMap)
		disksMap["size"] = sizeList
	}
	if disksItem.SupportedInterfaceTypes != nil {
		supportedInterfaceTypesList := []map[string]interface{}{}
		supportedInterfaceTypesMap := dataSourceInstanceProfileDisksSupportedInterfaceTypesToMap(*disksItem.SupportedInterfaceTypes)
		supportedInterfaceTypesList = append(supportedInterfaceTypesList, supportedInterfaceTypesMap)
		disksMap["supported_interface_types"] = supportedInterfaceTypesList
	}

	return disksMap
}

func dataSourceInstanceProfileDisksQuantityToMap(quantityItem vpcv1.InstanceProfileDiskQuantity) (quantityMap map[string]interface{}) {
	quantityMap = map[string]interface{}{}

	if quantityItem.Type != nil {
		quantityMap["type"] = quantityItem.Type
	}
	if quantityItem.Value != nil {
		quantityMap["value"] = quantityItem.Value
	}
	if quantityItem.Default != nil {
		quantityMap["default"] = quantityItem.Default
	}
	if quantityItem.Max != nil {
		quantityMap["max"] = quantityItem.Max
	}
	if quantityItem.Min != nil {
		quantityMap["min"] = quantityItem.Min
	}
	if quantityItem.Step != nil {
		quantityMap["step"] = quantityItem.Step
	}
	if quantityItem.Values != nil {
		quantityMap["values"] = quantityItem.Values
	}

	return quantityMap
}

func dataSourceInstanceProfileDisksSizeToMap(sizeItem vpcv1.InstanceProfileDiskSize) (sizeMap map[string]interface{}) {
	sizeMap = map[string]interface{}{}

	if sizeItem.Type != nil {
		sizeMap["type"] = sizeItem.Type
	}
	if sizeItem.Value != nil {
		sizeMap["value"] = sizeItem.Value
	}
	if sizeItem.Default != nil {
		sizeMap["default"] = sizeItem.Default
	}
	if sizeItem.Max != nil {
		sizeMap["max"] = sizeItem.Max
	}
	if sizeItem.Min != nil {
		sizeMap["min"] = sizeItem.Min
	}
	if sizeItem.Step != nil {
		sizeMap["step"] = sizeItem.Step
	}
	if sizeItem.Values != nil {
		sizeMap["values"] = sizeItem.Values
	}

	return sizeMap
}

func dataSourceInstanceProfileDisksSupportedInterfaceTypesToMap(supportedInterfaceTypesItem vpcv1.InstanceProfileDiskSupportedInterfaces) (supportedInterfaceTypesMap map[string]interface{}) {
	supportedInterfaceTypesMap = map[string]interface{}{}

	if supportedInterfaceTypesItem.Default != nil {
		supportedInterfaceTypesMap["default"] = supportedInterfaceTypesItem.Default
	}
	if supportedInterfaceTypesItem.Type != nil {
		supportedInterfaceTypesMap["type"] = supportedInterfaceTypesItem.Type
	}
	if supportedInterfaceTypesItem.Values != nil {
		supportedInterfaceTypesMap["values"] = supportedInterfaceTypesItem.Values
	}

	return supportedInterfaceTypesMap
}

func dataSourceInstanceProfileFlattenTotalVolumeBandwidth(result vpcv1.InstanceProfileVolumeBandwidth) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileTotalVolumeBandwidthToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileTotalVolumeBandwidthToMap(bandwidthItem vpcv1.InstanceProfileVolumeBandwidth) (bandwidthMap map[string]interface{}) {
	bandwidthMap = map[string]interface{}{}

	if bandwidthItem.Type != nil {
		bandwidthMap["type"] = bandwidthItem.Type
	}
	if bandwidthItem.Value != nil {
		bandwidthMap["value"] = bandwidthItem.Value
	}
	if bandwidthItem.Default != nil {
		bandwidthMap["default"] = bandwidthItem.Default
	}
	if bandwidthItem.Max != nil {
		bandwidthMap["max"] = bandwidthItem.Max
	}
	if bandwidthItem.Min != nil {
		bandwidthMap["min"] = bandwidthItem.Min
	}
	if bandwidthItem.Step != nil {
		bandwidthMap["step"] = bandwidthItem.Step
	}
	if bandwidthItem.Values != nil {
		bandwidthMap["values"] = bandwidthItem.Values
	}

	return bandwidthMap
}
