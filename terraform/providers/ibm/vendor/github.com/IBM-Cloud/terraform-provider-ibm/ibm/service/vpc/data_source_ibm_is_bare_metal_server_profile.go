// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isBareMetalServerProfileName            = "name"
	isBareMetalServerProfileBandwidth       = "bandwidth"
	isBareMetalServerProfileType            = "type"
	isBareMetalServerProfileValue           = "value"
	isBareMetalServerProfileCPUArchitecture = "cpu_architecture"
	isBareMetalServerProfileCPUCoreCount    = "cpu_core_count"
	isBareMetalServerProfileCPUSocketCount  = "cpu_socket_count"
	isBareMetalServerProfileDisks           = "disks"
	isBareMetalServerProfileDiskQuantity    = "quantity"
	isBareMetalServerProfileDiskSize        = "size"
	isBareMetalServerProfileDiskSITs        = "supported_interface_types"
	isBareMetalServerProfileFamily          = "family"
	isBareMetalServerProfileHref            = "href"
	isBareMetalServerProfileMemory          = "memory"
	isBareMetalServerProfileOS              = "os_architecture"
	isBareMetalServerProfileValues          = "values"
	isBareMetalServerProfileDefault         = "default"
	isBareMetalServerProfileRT              = "resource_type"
	isBareMetalServerProfileSIFs            = "supported_image_flags"
	isBareMetalServerProfileSTPMMs          = "supported_trusted_platform_module_modes"
)

func DataSourceIBMIsBareMetalServerProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBMSProfileRead,

		Schema: map[string]*schema.Schema{
			isBareMetalServerProfileName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name for this bare metal server profile",
			},

			// vni

			"virtual_network_interfaces_supported": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Indicates whether this profile supports virtual network interfaces.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The value for this profile field.",
						},
					},
				},
			},
			"network_attachment_count": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
					},
				},
			},

			isBareMetalServerProfileFamily: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product family this bare metal server profile belongs to",
			},
			isBareMetalServerProfileHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this bare metal server profile",
			},
			isBareMetalServerProfileBandwidth: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The total bandwidth (in megabits per second) shared across the network interfaces of a bare metal server with this profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerProfileType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field",
						},

						isBareMetalServerProfileValue: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field",
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
			"network_interface_count": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
					},
				},
			},
			"console_types": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The console type configuration for a bare metal server with this profile.",
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
							Description: "The console types for a bare metal server with this profile.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			isBareMetalServerProfileRT: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type for this bare metal server profile",
			},

			isBareMetalServerProfileCPUArchitecture: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The CPU architecture for a bare metal server with this profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerProfileType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field",
						},

						isBareMetalServerProfileValue: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value for this profile field",
						},
					},
				},
			},

			isBareMetalServerProfileCPUSocketCount: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The number of CPU sockets for a bare metal server with this profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerProfileType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field",
						},

						isBareMetalServerProfileValue: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field",
						},
					},
				},
			},

			isBareMetalServerProfileCPUCoreCount: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The CPU core count for a bare metal server with this profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerProfileType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field",
						},

						isBareMetalServerProfileValue: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field",
						},
					},
				},
			},
			isBareMetalServerProfileMemory: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The memory (in gibibytes) for a bare metal server with this profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerProfileType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field",
						},

						isBareMetalServerProfileValue: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field",
						},
					},
				},
			},

			isBareMetalServerProfileSTPMMs: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An array of supported trusted platform module (TPM) modes for this bare metal server profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerProfileType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field",
						},

						isBareMetalServerProfileValues: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "The supported trusted platform module (TPM) modes",
						},
					},
				},
			},
			isBareMetalServerProfileOS: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The supported OS architecture(s) for a bare metal server with this profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerProfileDefault: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default for this profile field",
						},
						isBareMetalServerProfileType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field",
						},

						isBareMetalServerProfileValues: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "The supported OS architecture(s) for a bare metal server with this profile",
						},
					},
				},
			},
			isBareMetalServerProfileDisks: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of the bare metal server profile's disks",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerProfileDiskQuantity: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The number of disks of this configuration for a bare metal server with this profile",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isBareMetalServerProfileType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field",
									},
									isBareMetalServerProfileValue: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field",
									},
								},
							},
						},

						isBareMetalServerProfileDiskSize: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The size of the disk in GB (gigabytes)",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isBareMetalServerProfileType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field",
									},
									isBareMetalServerProfileValue: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field",
									},
								},
							},
						},
						isBareMetalServerProfileDiskSITs: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The disk interface used for attaching the disk.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isBareMetalServerProfileDefault: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
									},
									isBareMetalServerProfileType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field",
									},
									isBareMetalServerProfileValues: {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "The supported disk interfaces used for attaching the disk",
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         flex.ResourceIBMVPCHash,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISBMSProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	options := &vpcv1.GetBareMetalServerProfileOptions{
		Name: &name,
	}
	bmsProfile, response, err := sess.GetBareMetalServerProfileWithContext(context, options)
	if err != nil || bmsProfile == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error Getting Bare Metal Server Profile (%s): %s\n%s", name, err, response))
	}
	d.SetId(*bmsProfile.Name)
	d.Set(isBareMetalServerProfileName, *bmsProfile.Name)
	d.Set(isBareMetalServerProfileFamily, *bmsProfile.Family)
	d.Set(isBareMetalServerProfileHref, *bmsProfile.Href)
	if bmsProfile.Bandwidth != nil {
		bwList := make([]map[string]interface{}, 0)
		bw := bmsProfile.Bandwidth.(*vpcv1.BareMetalServerProfileBandwidth)
		bandwidth := map[string]interface{}{}
		if bw.Type != nil {
			bandwidth[isBareMetalServerProfileType] = *bw.Type
		}
		if bw.Value != nil {
			bandwidth[isBareMetalServerProfileValue] = *bw.Value
		}
		if bw.Values != nil && len(bw.Values) > 0 {
			bandwidth[isBareMetalServerProfileValues] = bw.Values
		}
		if bw.Default != nil {
			bandwidth["default"] = flex.IntValue(bw.Default)
		}
		if bw.Max != nil {
			bandwidth["max"] = flex.IntValue(bw.Max)
		}
		if bw.Min != nil {
			bandwidth["min"] = flex.IntValue(bw.Min)
		}
		if bw.Step != nil {
			bandwidth["step"] = flex.IntValue(bw.Step)
		}
		bwList = append(bwList, bandwidth)
		d.Set(isBareMetalServerProfileBandwidth, bwList)
	}
	consoleTypes := []map[string]interface{}{}
	if bmsProfile.ConsoleTypes != nil {
		modelMap, err := dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileConsoleTypesToMap(bmsProfile.ConsoleTypes)
		if err != nil {
			return diag.FromErr(err)
		}
		consoleTypes = append(consoleTypes, modelMap)
	}
	if err = d.Set("console_types", consoleTypes); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting console_types %s", err))
	}

	networkInterfaceCount := []map[string]interface{}{}
	if bmsProfile.NetworkInterfaceCount != nil {
		modelMap, err := dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkInterfaceCountToMap(bmsProfile.NetworkInterfaceCount)
		if err != nil {
			return diag.FromErr(err)
		}
		networkInterfaceCount = append(networkInterfaceCount, modelMap)
	}
	if err = d.Set("network_interface_count", networkInterfaceCount); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting network_interface_count %s", err))
	}

	if bmsProfile.CpuArchitecture != nil {
		caList := make([]map[string]interface{}, 0)
		ca := bmsProfile.CpuArchitecture
		architecture := map[string]interface{}{
			isBareMetalServerProfileType:  *ca.Type,
			isBareMetalServerProfileValue: *ca.Value,
		}
		caList = append(caList, architecture)
		d.Set(isBareMetalServerProfileCPUArchitecture, caList)
	}
	if bmsProfile.CpuCoreCount != nil {
		ccList := make([]map[string]interface{}, 0)
		cc := bmsProfile.CpuCoreCount.(*vpcv1.BareMetalServerProfileCpuCoreCount)
		coreCount := map[string]interface{}{
			isBareMetalServerProfileType:  *cc.Type,
			isBareMetalServerProfileValue: *cc.Value,
		}
		ccList = append(ccList, coreCount)
		d.Set(isBareMetalServerProfileCPUCoreCount, ccList)
	}
	if bmsProfile.CpuSocketCount != nil {
		scList := make([]map[string]interface{}, 0)
		sc := bmsProfile.CpuSocketCount.(*vpcv1.BareMetalServerProfileCpuSocketCount)
		socketCount := map[string]interface{}{
			isBareMetalServerProfileType:  *sc.Type,
			isBareMetalServerProfileValue: *sc.Value,
		}
		scList = append(scList, socketCount)
		d.Set(isBareMetalServerProfileCPUSocketCount, scList)
	}

	if bmsProfile.Memory != nil {
		memList := make([]map[string]interface{}, 0)
		mem := bmsProfile.Memory.(*vpcv1.BareMetalServerProfileMemory)
		m := map[string]interface{}{
			isBareMetalServerProfileType:  *mem.Type,
			isBareMetalServerProfileValue: *mem.Value,
		}
		memList = append(memList, m)
		d.Set(isBareMetalServerProfileMemory, memList)
	}
	d.Set(isBareMetalServerProfileRT, bmsProfile.ResourceType)
	if bmsProfile.SupportedTrustedPlatformModuleModes != nil {
		list := make([]map[string]interface{}, 0)
		var stpmmlist []string
		for _, item := range bmsProfile.SupportedTrustedPlatformModuleModes.Values {
			stpmmlist = append(stpmmlist, item)
		}
		m := map[string]interface{}{
			isBareMetalServerProfileType: *bmsProfile.SupportedTrustedPlatformModuleModes.Type,
		}
		m[isBareMetalServerProfileValues] = stpmmlist
		list = append(list, m)
		d.Set(isBareMetalServerProfileSTPMMs, list)
	}
	if bmsProfile.OsArchitecture != nil {
		list := make([]map[string]interface{}, 0)
		var valuelist []string
		for _, item := range bmsProfile.OsArchitecture.Values {
			valuelist = append(valuelist, item)
		}
		m := map[string]interface{}{
			isBareMetalServerProfileDefault: *bmsProfile.OsArchitecture.Default,
			isBareMetalServerProfileType:    *bmsProfile.OsArchitecture.Type,
		}
		m[isBareMetalServerProfileValues] = valuelist
		list = append(list, m)
		d.Set(isBareMetalServerProfileOS, list)
	}

	if bmsProfile.Disks != nil {
		list := make([]map[string]interface{}, 0)
		for _, disk := range bmsProfile.Disks {
			qlist := make([]map[string]interface{}, 0)
			slist := make([]map[string]interface{}, 0)
			sitlist := make([]map[string]interface{}, 0)
			quantity := disk.Quantity.(*vpcv1.BareMetalServerProfileDiskQuantity)
			q := make(map[string]interface{})
			q[isBareMetalServerProfileType] = *quantity.Type
			q[isBareMetalServerProfileValue] = *quantity.Value
			qlist = append(qlist, q)
			size := disk.Size.(*vpcv1.BareMetalServerProfileDiskSize)
			s := map[string]interface{}{
				isBareMetalServerProfileType:  *size.Type,
				isBareMetalServerProfileValue: *size.Value,
			}
			slist = append(slist, s)
			sit := map[string]interface{}{
				isBareMetalServerProfileDefault: *disk.SupportedInterfaceTypes.Default,
				isBareMetalServerProfileType:    *disk.SupportedInterfaceTypes.Type,
			}
			var valuelist []string
			for _, item := range disk.SupportedInterfaceTypes.Values {
				valuelist = append(valuelist, item)
			}
			sit[isBareMetalServerProfileValues] = valuelist
			sitlist = append(sitlist, sit)
			sz := map[string]interface{}{
				isBareMetalServerProfileDiskQuantity: qlist,
				isBareMetalServerProfileDiskSize:     slist,
				isBareMetalServerProfileDiskSITs:     sitlist,
			}
			list = append(list, sz)
		}
		d.Set(isBareMetalServerProfileDisks, list)

		// vni
		virtualNetworkInterfacesSupported := []map[string]interface{}{}
		if bmsProfile.VirtualNetworkInterfacesSupported != nil {
			modelMap, err := dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileVirtualNetworkInterfacesSupportedToMap(bmsProfile.VirtualNetworkInterfacesSupported)
			if err != nil {
				return diag.FromErr(err)
			}
			virtualNetworkInterfacesSupported = append(virtualNetworkInterfacesSupported, modelMap)
		}
		if err = d.Set("virtual_network_interfaces_supported", virtualNetworkInterfacesSupported); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting virtual_network_interfaces_supported %s", err))
		}
		networkAttachmentCount := []map[string]interface{}{}
		if bmsProfile.NetworkAttachmentCount != nil {
			modelMap, err := dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkAttachmentCountToMap(bmsProfile.NetworkAttachmentCount)
			if err != nil {
				return diag.FromErr(err)
			}
			networkAttachmentCount = append(networkAttachmentCount, modelMap)
		}
		if err = d.Set("network_attachment_count", networkAttachmentCount); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting network_attachment_count %s", err))
		}
	}

	return nil
}

func dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileConsoleTypesToMap(model *vpcv1.BareMetalServerProfileConsoleTypes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["values"] = model.Values
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkInterfaceCountToMap(model vpcv1.BareMetalServerProfileNetworkInterfaceCountIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.BareMetalServerProfileNetworkInterfaceCountRange); ok {
		return dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkInterfaceCountRangeToMap(model.(*vpcv1.BareMetalServerProfileNetworkInterfaceCountRange))
	} else if _, ok := model.(*vpcv1.BareMetalServerProfileNetworkInterfaceCountDependent); ok {
		return dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkInterfaceCountDependentToMap(model.(*vpcv1.BareMetalServerProfileNetworkInterfaceCountDependent))
	} else if _, ok := model.(*vpcv1.BareMetalServerProfileNetworkInterfaceCount); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.BareMetalServerProfileNetworkInterfaceCount)
		if model.Max != nil {
			modelMap["max"] = flex.IntValue(model.Max)
		}
		if model.Min != nil {
			modelMap["min"] = flex.IntValue(model.Min)
		}
		if model.Type != nil {
			modelMap["type"] = model.Type
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.BareMetalServerProfileNetworkInterfaceCountIntf subtype encountered")
	}
}

func dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkInterfaceCountRangeToMap(model *vpcv1.BareMetalServerProfileNetworkInterfaceCountRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Max != nil {
		modelMap["max"] = flex.IntValue(model.Max)
	}
	if model.Min != nil {
		modelMap["min"] = flex.IntValue(model.Min)
	}
	modelMap["type"] = model.Type
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkInterfaceCountDependentToMap(model *vpcv1.BareMetalServerProfileNetworkInterfaceCountDependent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileVirtualNetworkInterfacesSupportedToMap(model *vpcv1.BareMetalServerProfileVirtualNetworkInterfacesSupported) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["value"] = model.Value
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkAttachmentCountToMap(model vpcv1.BareMetalServerProfileNetworkAttachmentCountIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.BareMetalServerProfileNetworkAttachmentCountRange); ok {
		return dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkAttachmentCountRangeToMap(model.(*vpcv1.BareMetalServerProfileNetworkAttachmentCountRange))
	} else if _, ok := model.(*vpcv1.BareMetalServerProfileNetworkAttachmentCountDependent); ok {
		return dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkAttachmentCountDependentToMap(model.(*vpcv1.BareMetalServerProfileNetworkAttachmentCountDependent))
	} else if _, ok := model.(*vpcv1.BareMetalServerProfileNetworkAttachmentCount); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.BareMetalServerProfileNetworkAttachmentCount)
		if model.Max != nil {
			modelMap["max"] = flex.IntValue(model.Max)
		}
		if model.Min != nil {
			modelMap["min"] = flex.IntValue(model.Min)
		}
		if model.Type != nil {
			modelMap["type"] = model.Type
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.BareMetalServerProfileNetworkAttachmentCountIntf subtype encountered")
	}
}

func dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkAttachmentCountRangeToMap(model *vpcv1.BareMetalServerProfileNetworkAttachmentCountRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Max != nil {
		modelMap["max"] = flex.IntValue(model.Max)
	}
	if model.Min != nil {
		modelMap["min"] = flex.IntValue(model.Min)
	}
	modelMap["type"] = model.Type
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkAttachmentCountDependentToMap(model *vpcv1.BareMetalServerProfileNetworkAttachmentCountDependent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	return modelMap, nil
}
