// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isBareMetalServerProfiles = "profiles"
)

func DataSourceIBMIsBareMetalServerProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsBareMetalServerProfilesRead,

		Schema: map[string]*schema.Schema{

			isBareMetalServerProfiles: {
				Type:        schema.TypeList,
				Description: "List of BMS profile maps",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						isBareMetalServerProfileName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this bare metal server profile",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
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
				},
			},
		},
	}
}

func dataSourceIBMIsBareMetalServerProfilesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	start := ""
	allrecs := []vpcv1.BareMetalServerProfile{}
	for {
		listBMSProfilesOptions := &vpcv1.ListBareMetalServerProfilesOptions{}
		if start != "" {
			listBMSProfilesOptions.Start = &start
		}
		availableProfiles, response, err := sess.ListBareMetalServerProfilesWithContext(context, listBMSProfilesOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error fetching Bare Metal Server Profiles %s\n%s", err, response))
		}
		start = flex.GetNext(availableProfiles.Next)
		allrecs = append(allrecs, availableProfiles.Profiles...)
		if start == "" {
			break
		}
	}

	profilesInfo := make([]map[string]interface{}, 0)
	for _, profile := range allrecs {

		l := map[string]interface{}{
			isBareMetalServerProfileName:   *profile.Name,
			"id":                           *profile.Name,
			isBareMetalServerProfileFamily: *profile.Family,
		}
		l[isBareMetalServerProfileHref] = *profile.Href
		if profile.Bandwidth != nil {
			bwList := make([]map[string]interface{}, 0)
			bw := profile.Bandwidth.(*vpcv1.BareMetalServerProfileBandwidth)
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
			l[isBareMetalServerProfileBandwidth] = bwList
		}
		if profile.CpuArchitecture != nil {
			caList := make([]map[string]interface{}, 0)
			ca := profile.CpuArchitecture
			architecture := map[string]interface{}{
				isBareMetalServerProfileType:  *ca.Type,
				isBareMetalServerProfileValue: *ca.Value,
			}
			caList = append(caList, architecture)
			l[isBareMetalServerProfileCPUArchitecture] = caList
		}
		if profile.CpuCoreCount != nil {
			ccList := make([]map[string]interface{}, 0)
			cc := profile.CpuCoreCount.(*vpcv1.BareMetalServerProfileCpuCoreCount)
			coreCount := map[string]interface{}{
				isBareMetalServerProfileType:  *cc.Type,
				isBareMetalServerProfileValue: *cc.Value,
			}
			ccList = append(ccList, coreCount)
			l[isBareMetalServerProfileCPUCoreCount] = ccList
		}
		consoleTypes := []map[string]interface{}{}
		if profile.ConsoleTypes != nil {
			modelMap, err := dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileConsoleTypesToMap(profile.ConsoleTypes)
			if err != nil {
				return diag.FromErr(err)
			}
			consoleTypes = append(consoleTypes, modelMap)
		}
		l["console_types"] = consoleTypes

		networkInterfaceCount := []map[string]interface{}{}
		if profile.NetworkInterfaceCount != nil {
			modelMap, err := dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkInterfaceCountToMap(profile.NetworkInterfaceCount)
			if err != nil {
				return diag.FromErr(err)
			}
			networkInterfaceCount = append(networkInterfaceCount, modelMap)
		}
		l["network_interface_count"] = networkInterfaceCount

		if profile.CpuSocketCount != nil {
			scList := make([]map[string]interface{}, 0)
			sc := profile.CpuSocketCount.(*vpcv1.BareMetalServerProfileCpuSocketCount)
			socketCount := map[string]interface{}{
				isBareMetalServerProfileType:  *sc.Type,
				isBareMetalServerProfileValue: *sc.Value,
			}
			scList = append(scList, socketCount)
			l[isBareMetalServerProfileCPUSocketCount] = scList
		}

		// vni
		virtualNetworkInterfacesSupported := []map[string]interface{}{}
		if profile.VirtualNetworkInterfacesSupported != nil {
			modelMap, err := dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileVirtualNetworkInterfacesSupportedToMap(profile.VirtualNetworkInterfacesSupported)
			if err != nil {
				return diag.FromErr(err)
			}
			virtualNetworkInterfacesSupported = append(virtualNetworkInterfacesSupported, modelMap)
		}
		l["virtual_network_interfaces_supported"] = virtualNetworkInterfacesSupported
		networkAttachmentCount := []map[string]interface{}{}
		if profile.NetworkAttachmentCount != nil {
			modelMap, err := dataSourceIBMIsBareMetalServerProfileBareMetalServerProfileNetworkAttachmentCountToMap(profile.NetworkAttachmentCount)
			if err != nil {
				return diag.FromErr(err)
			}
			networkAttachmentCount = append(networkAttachmentCount, modelMap)
		}
		l["network_attachment_count"] = networkAttachmentCount

		if profile.Memory != nil {
			memList := make([]map[string]interface{}, 0)
			mem := profile.Memory.(*vpcv1.BareMetalServerProfileMemory)
			m := map[string]interface{}{
				isBareMetalServerProfileType:  *mem.Type,
				isBareMetalServerProfileValue: *mem.Value,
			}
			memList = append(memList, m)
			l[isBareMetalServerProfileMemory] = memList
		}
		l[isBareMetalServerProfileRT] = profile.ResourceType
		if profile.SupportedTrustedPlatformModuleModes != nil {
			list := make([]map[string]interface{}, 0)
			var stpmmlist []string
			for _, item := range profile.SupportedTrustedPlatformModuleModes.Values {
				stpmmlist = append(stpmmlist, item)
			}
			m := map[string]interface{}{
				isBareMetalServerProfileType: *profile.SupportedTrustedPlatformModuleModes.Type,
			}
			m[isBareMetalServerProfileValues] = stpmmlist
			list = append(list, m)
			l[isBareMetalServerProfileSTPMMs] = list
		}
		if profile.OsArchitecture != nil {
			list := make([]map[string]interface{}, 0)
			var valuelist []string
			for _, item := range profile.OsArchitecture.Values {
				valuelist = append(valuelist, item)
			}
			m := map[string]interface{}{
				isBareMetalServerProfileDefault: *profile.OsArchitecture.Default,
				isBareMetalServerProfileType:    *profile.OsArchitecture.Type,
			}
			m[isBareMetalServerProfileValues] = valuelist
			list = append(list, m)
			l[isBareMetalServerProfileOS] = list
		}

		if profile.Disks != nil {
			list := make([]map[string]interface{}, 0)
			for _, disk := range profile.Disks {
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
			l[isBareMetalServerProfileDisks] = list
		}

		profilesInfo = append(profilesInfo, l)
	}
	d.SetId(dataSourceIBMIsBMSProfilesID(d))
	d.Set(isBareMetalServerProfiles, profilesInfo)
	return nil
}

// dataSourceIBMIsBMSProfilesID returns a reasonable ID for a BMS Profile list.
func dataSourceIBMIsBMSProfilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
