// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isBareMetalServers = "servers"
)

func DataSourceIBMIsBareMetalServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBareMetalServersRead,

		Schema: map[string]*schema.Schema{
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the resource group this bare metal server belongs to",
			},
			"vpc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vpc ID this bare metal server is in",
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vpc name this bare metal server is in",
			},
			"vpc_crn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vpc CRN this bare metal server is in",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the bare metal server",
			},
			"network_interfaces_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the subnet of the bare metal server network interfaces",
			},
			"network_interfaces_subnet_crn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The crn of the subnet of the bare metal server network interfaces",
			},
			"network_interfaces_subnet_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the subnet of the bare metal server network interfaces",
			},
			isBareMetalServers: {
				Type:        schema.TypeList,
				Description: "List of Bare Metal Servers",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bare metal server id",
						},
						isBareMetalServerName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bare metal server name",
						},
						isBareMetalServerBandwidth: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total bandwidth (in megabits per second)",
						},
						isBareMetalServerEnableSecureBoot: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether secure boot is enabled. If enabled, the image must support secure boot or the server will fail to boot.",
						},

						isBareMetalServerTrustedPlatformModule: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isBareMetalServerTrustedPlatformModuleMode: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The trusted platform module mode to use. The specified value must be listed in the bare metal server profile's supported_trusted_platform_module_modes",
									},
									isBareMetalServerTrustedPlatformModuleEnabled: {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether the trusted platform module is enabled.",
									},
									isBareMetalServerTrustedPlatformModuleSupportedModes: {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         flex.ResourceIBMVPCHash,
										Computed:    true,
										Description: "The trusted platform module (TPM) mode:: disabled: No TPM functionality, tpm_2: TPM 2.0. The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered. Enum: [ disabled, tpm_2 ]",
									},
								},
							},
						},
						isBareMetalServerBootTarget: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this bare metal server disk",
						},
						isBareMetalServerCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the bare metal server was created",
						},
						isBareMetalServerCPU: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The bare metal server CPU configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isBareMetalServerCPUArchitecture: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CPU architecture",
									},
									isBareMetalServerCPUCoreCount: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The total number of cores",
									},
									isBareMetalServerCpuSocketCount: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The total number of CPU sockets",
									},
									isBareMetalServerCpuThreadPerCore: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The total number of hardware threads per core",
									},
								},
							},
						},
						isBareMetalServerCRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this bare metal server",
						},
						isBareMetalServerDisks: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The disks for this bare metal server, including any disks that are associated with the boot_target.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isBareMetalServerDiskHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this bare metal server disk",
									},
									isBareMetalServerDiskID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this bare metal server disk",
									},
									isBareMetalServerDiskInterfaceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The disk interface used for attaching the disk. Supported values are [ nvme, sata ]",
									},
									isBareMetalServerDiskName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this disk",
									},
									isBareMetalServerDiskResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type",
									},
									isBareMetalServerDiskSize: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The size of the disk in GB (gigabytes)",
									},
								},
							},
						},
						isBareMetalServerHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this bare metal server",
						},
						isBareMetalServerMemory: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of memory, truncated to whole gibibytes",
						},

						isBareMetalServerPrimaryNetworkInterface: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Primary Network interface info",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									isBareMetalServerNicAllowIPSpoofing: {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether IP spoofing is allowed on this interface.",
									},
									isBareMetalServerNicName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isBareMetalServerNicPortSpeed: {
										Type:     schema.TypeInt,
										Computed: true,
									},
									isBareMetalServerNicHref: {
										Type:       schema.TypeString,
										Computed:   true,
										Deprecated: "This URL of the interface",
									},

									isBareMetalServerNicSecurityGroups: {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},
									isBareMetalServerNicSubnet: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isBareMetalServerNicPrimaryIP: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "IPv4, The IP address. ",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isBareMetalServerNicIpAddress: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The globally unique IP address",
												},
												isBareMetalServerNicIpHref: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this reserved IP",
												},
												isBareMetalServerNicIpName: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
												},
												isBareMetalServerNicIpID: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Identifies a reserved IP by a unique property.",
												},
												isBareMetalServerNicResourceType: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type",
												},
											},
										},
									},
								},
							},
						},

						isBareMetalServerNetworkInterfaces: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									isBareMetalServerNicHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this network interface",
									},
									isBareMetalServerNicAllowIPSpoofing: {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether IP spoofing is allowed on this interface.",
									},
									isBareMetalServerNicName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isBareMetalServerNicSecurityGroups: {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},
									isBareMetalServerNicSubnet: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isBareMetalServerNicPrimaryIP: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "IPv4, The IP address. ",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isBareMetalServerNicIpAddress: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The globally unique IP address",
												},
												isBareMetalServerNicIpHref: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this reserved IP",
												},
												isBareMetalServerNicIpName: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
												},
												isBareMetalServerNicIpID: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Identifies a reserved IP by a unique property.",
												},
												isBareMetalServerNicResourceType: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type",
												},
											},
										},
									},
								},
							},
						},

						isBareMetalServerKeys: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "SSH key Ids for the bare metal server",
						},

						isBareMetalServerImage: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "image id",
						},
						isBareMetalServerProfile: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "profile name",
						},

						isBareMetalServerZone: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone name",
						},

						isBareMetalServerVPC: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC the bare metal server is to be a part of",
						},

						isBareMetalServerResourceGroup: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource group name",
						},
						isBareMetalServerResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type name",
						},

						isBareMetalServerStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bare metal server status",
						},

						isBareMetalServerStatusReasons: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isBareMetalServerStatusReasonsCode: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the status reason",
									},

									isBareMetalServerStatusReasonsMessage: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the status reason",
									},

									isBareMetalServerStatusReasonsMoreInfo: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about this status reason",
									},
								},
							},
						},
						isBareMetalServerTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server", "tags")},
							Set:         flex.ResourceIBMVPCHash,
							Description: "Tags for the Bare metal server",
						},
						isBareMetalServerAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access tags",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISBareMetalServersRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	start := ""
	allrecs := []vpcv1.BareMetalServer{}

	listBareMetalServersOptions := &vpcv1.ListBareMetalServersOptions{}
	if resgroupintf, ok := d.GetOk("resource_group"); ok {
		resGroup := resgroupintf.(string)
		listBareMetalServersOptions.ResourceGroupID = &resGroup
	}
	if nameintf, ok := d.GetOk("name"); ok {
		name := nameintf.(string)
		listBareMetalServersOptions.Name = &name
	}
	if vpcIntf, ok := d.GetOk("vpc"); ok {
		vpcid := vpcIntf.(string)
		listBareMetalServersOptions.VPCID = &vpcid
	}
	if vpcNameIntf, ok := d.GetOk("vpc_name"); ok {
		vpcName := vpcNameIntf.(string)
		listBareMetalServersOptions.VPCName = &vpcName
	}
	if vpcCrnIntf, ok := d.GetOk("vpc_crn"); ok {
		vpcCrn := vpcCrnIntf.(string)
		listBareMetalServersOptions.VPCCRN = &vpcCrn
	}
	// if subnetIntf, ok := d.GetOk("network_interfaces_subnet"); ok {
	// 	subnetId := subnetIntf.(string)
	// 	listBareMetalServersOptions.NetworkInterfacesSubnetID = &subnetId
	// }
	// if subnetNameIntf, ok := d.GetOk("network_interfaces_subnet_name"); ok {
	// 	subnetName := subnetNameIntf.(string)
	// 	listBareMetalServersOptions.NetworkInterfacesSubnetName = &subnetName
	// }
	// if subnetCrnIntf, ok := d.GetOk("network_interfaces_subnet_crn"); ok {
	// 	subnetCrn := subnetCrnIntf.(string)
	// 	listBareMetalServersOptions.NetworkInterfacesSubnetCRN = &subnetCrn
	// }
	for {

		if start != "" {
			listBareMetalServersOptions.Start = &start
		}
		availableServers, response, err := sess.ListBareMetalServersWithContext(context, listBareMetalServersOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error fetching Bare Metal Servers %s\n%s", err, response))
		}
		start = flex.GetNext(availableServers.Next)
		allrecs = append(allrecs, availableServers.BareMetalServers...)
		if start == "" {
			break
		}
	}

	serversInfo := make([]map[string]interface{}, 0)
	for _, bms := range allrecs {

		l := map[string]interface{}{
			isBareMetalServerName: *bms.Name,
		}
		l["id"] = *bms.ID
		l[isBareMetalServerBandwidth] = *bms.Bandwidth
		bmsBootTargetIntf := bms.BootTarget.(*vpcv1.BareMetalServerBootTarget)
		bmsBootTarget := bmsBootTargetIntf.ID
		l[isBareMetalServerBootTarget] = bmsBootTarget
		cpuList := make([]map[string]interface{}, 0)
		if bms.Cpu != nil {
			currentCPU := map[string]interface{}{}
			currentCPU[isBareMetalServerCPUArchitecture] = *bms.Cpu.Architecture
			currentCPU[isBareMetalServerCPUCoreCount] = *bms.Cpu.CoreCount
			currentCPU[isBareMetalServerCpuSocketCount] = *bms.Cpu.SocketCount
			currentCPU[isBareMetalServerCpuThreadPerCore] = *bms.Cpu.ThreadsPerCore
			cpuList = append(cpuList, currentCPU)
		}
		l[isBareMetalServerCPU] = cpuList
		l[isBareMetalServerName] = *bms.Name
		l[isBareMetalServerCRN] = *bms.CRN

		// disks

		diskList := make([]map[string]interface{}, 0)
		if bms.Disks != nil {
			for _, disk := range bms.Disks {
				currentDisk := map[string]interface{}{
					isBareMetalServerDiskHref:          disk.Href,
					isBareMetalServerDiskID:            disk.ID,
					isBareMetalServerDiskInterfaceType: disk.InterfaceType,
					isBareMetalServerDiskName:          disk.Name,
					isBareMetalServerDiskResourceType:  disk.ResourceType,
					isBareMetalServerDiskSize:          disk.Size,
				}
				diskList = append(diskList, currentDisk)
			}
		}
		l[isBareMetalServerDisks] = diskList

		l[isBareMetalServerHref] = *bms.Href
		l[isBareMetalServerMemory] = *bms.Memory
		l[isBareMetalServerProfile] = *bms.Profile.Name

		//enable secure boot
		if bms.EnableSecureBoot != nil {
			l[isBareMetalServerEnableSecureBoot] = bms.EnableSecureBoot
		}

		// tpm
		if bms.TrustedPlatformModule != nil {
			trustedPlatformModuleMap, err := resourceIBMIsBareMetalServerBareMetalServerTrustedPlatformModulePrototypeToMap(bms.TrustedPlatformModule)
			if err != nil {
				return diag.FromErr(err)
			}
			l[isBareMetalServerTrustedPlatformModule] = []map[string]interface{}{trustedPlatformModuleMap}
		}

		//pni

		if bms.PrimaryNetworkInterface != nil && bms.PrimaryNetworkInterface.ID != nil {
			primaryNicList := make([]map[string]interface{}, 0)
			currentPrimNic := map[string]interface{}{}
			currentPrimNic["id"] = *bms.PrimaryNetworkInterface.ID
			currentPrimNic[isBareMetalServerNicHref] = *bms.PrimaryNetworkInterface.Href
			currentPrimNic[isBareMetalServerNicName] = *bms.PrimaryNetworkInterface.Name
			currentPrimNic[isBareMetalServerNicHref] = *bms.PrimaryNetworkInterface.Href
			currentPrimNic[isBareMetalServerNicSubnet] = *bms.PrimaryNetworkInterface.Subnet.ID
			primaryIpList := make([]map[string]interface{}, 0)
			currentIP := map[string]interface{}{}
			if bms.PrimaryNetworkInterface.PrimaryIP.Href != nil {
				currentIP[isBareMetalServerNicIpAddress] = *bms.PrimaryNetworkInterface.PrimaryIP.Address
			}
			if bms.PrimaryNetworkInterface.PrimaryIP.Href != nil {
				currentIP[isBareMetalServerNicIpHref] = *bms.PrimaryNetworkInterface.PrimaryIP.Href
			}
			if bms.PrimaryNetworkInterface.PrimaryIP.Name != nil {
				currentIP[isBareMetalServerNicIpName] = *bms.PrimaryNetworkInterface.PrimaryIP.Name
			}
			if bms.PrimaryNetworkInterface.PrimaryIP.ID != nil {
				currentIP[isBareMetalServerNicIpID] = *bms.PrimaryNetworkInterface.PrimaryIP.ID
			}
			if bms.PrimaryNetworkInterface.PrimaryIP.ResourceType != nil {
				currentIP[isBareMetalServerNicResourceType] = *bms.PrimaryNetworkInterface.PrimaryIP.ResourceType
			}
			primaryIpList = append(primaryIpList, currentIP)
			currentPrimNic[isBareMetalServerNicPrimaryIP] = primaryIpList
			getnicoptions := &vpcv1.GetBareMetalServerNetworkInterfaceOptions{
				BareMetalServerID: bms.ID,
				ID:                bms.PrimaryNetworkInterface.ID,
			}
			bmsnic, response, err := sess.GetBareMetalServerNetworkInterface(getnicoptions)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error getting network interfaces attached to the bare metal server %s\n%s", err, response))
			}

			switch reflect.TypeOf(bmsnic).String() {
			case "*vpcv1.BareMetalServerNetworkInterfaceByPci":
				{
					primNic := bmsnic.(*vpcv1.BareMetalServerNetworkInterfaceByPci)
					currentPrimNic[isInstanceNicAllowIPSpoofing] = *primNic.AllowIPSpoofing
					currentPrimNic[isBareMetalServerNicPortSpeed] = *primNic.PortSpeed
					if len(primNic.SecurityGroups) != 0 {
						secgrpList := []string{}
						for i := 0; i < len(primNic.SecurityGroups); i++ {
							secgrpList = append(secgrpList, string(*(primNic.SecurityGroups[i].ID)))
						}
						currentPrimNic[isInstanceNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
					}
				}
			case "*vpcv1.BareMetalServerNetworkInterfaceByVlan":
				{
					primNic := bmsnic.(*vpcv1.BareMetalServerNetworkInterfaceByVlan)
					currentPrimNic[isInstanceNicAllowIPSpoofing] = *primNic.AllowIPSpoofing
					currentPrimNic[isBareMetalServerNicPortSpeed] = *primNic.PortSpeed

					if len(primNic.SecurityGroups) != 0 {
						secgrpList := []string{}
						for i := 0; i < len(primNic.SecurityGroups); i++ {
							secgrpList = append(secgrpList, string(*(primNic.SecurityGroups[i].ID)))
						}
						currentPrimNic[isInstanceNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
					}
				}
			case "*vpcv1.BareMetalServerNetworkInterfaceByHiperSocket":
				{
					primNic := bmsnic.(*vpcv1.BareMetalServerNetworkInterfaceByHiperSocket)
					currentPrimNic[isInstanceNicAllowIPSpoofing] = *primNic.AllowIPSpoofing

					if len(primNic.SecurityGroups) != 0 {
						secgrpList := []string{}
						for i := 0; i < len(primNic.SecurityGroups); i++ {
							secgrpList = append(secgrpList, string(*(primNic.SecurityGroups[i].ID)))
						}
						currentPrimNic[isInstanceNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
					}
				}
			}

			primaryNicList = append(primaryNicList, currentPrimNic)
			l[isBareMetalServerPrimaryNetworkInterface] = primaryNicList
		}

		//ni

		interfacesList := make([]map[string]interface{}, 0)
		for _, intfc := range bms.NetworkInterfaces {
			if intfc.ID != nil && *intfc.ID != *bms.PrimaryNetworkInterface.ID {
				currentNic := map[string]interface{}{}
				currentNic["id"] = *intfc.ID
				currentNic[isBareMetalServerNicHref] = *intfc.Href
				currentNic[isBareMetalServerNicName] = *intfc.Name
				primaryIpList := make([]map[string]interface{}, 0)
				currentIP := map[string]interface{}{}
				if intfc.PrimaryIP.Href != nil {
					currentIP[isBareMetalServerNicIpAddress] = *intfc.PrimaryIP.Address
				}
				if intfc.PrimaryIP.Href != nil {
					currentIP[isBareMetalServerNicIpHref] = *intfc.PrimaryIP.Href
				}
				if intfc.PrimaryIP.Name != nil {
					currentIP[isBareMetalServerNicIpName] = *intfc.PrimaryIP.Name
				}
				if intfc.PrimaryIP.ID != nil {
					currentIP[isBareMetalServerNicIpID] = *intfc.PrimaryIP.ID
				}
				if intfc.PrimaryIP.ResourceType != nil {
					currentIP[isBareMetalServerNicResourceType] = *intfc.PrimaryIP.ResourceType
				}
				primaryIpList = append(primaryIpList, currentIP)
				currentNic[isBareMetalServerNicPrimaryIP] = primaryIpList
				getnicoptions := &vpcv1.GetBareMetalServerNetworkInterfaceOptions{
					BareMetalServerID: bms.ID,
					ID:                intfc.ID,
				}
				bmsnicintf, response, err := sess.GetBareMetalServerNetworkInterface(getnicoptions)
				if err != nil {
					return diag.FromErr(fmt.Errorf("[ERROR] Error getting network interfaces attached to the bare metal server %s\n%s", err, response))
				}

				switch reflect.TypeOf(bmsnicintf).String() {
				case "*vpcv1.BareMetalServerNetworkInterfaceByPci":
					{
						bmsnic := bmsnicintf.(*vpcv1.BareMetalServerNetworkInterfaceByPci)
						currentNic[isBareMetalServerNicAllowIPSpoofing] = *bmsnic.AllowIPSpoofing
						currentNic[isBareMetalServerNicSubnet] = *bmsnic.Subnet.ID
						if len(bmsnic.SecurityGroups) != 0 {
							secgrpList := []string{}
							for i := 0; i < len(bmsnic.SecurityGroups); i++ {
								secgrpList = append(secgrpList, string(*(bmsnic.SecurityGroups[i].ID)))
							}
							currentNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
						}
					}
				case "*vpcv1.BareMetalServerNetworkInterfaceByVlan":
					{
						bmsnic := bmsnicintf.(*vpcv1.BareMetalServerNetworkInterfaceByVlan)
						currentNic[isBareMetalServerNicAllowIPSpoofing] = *bmsnic.AllowIPSpoofing
						currentNic[isBareMetalServerNicSubnet] = *bmsnic.Subnet.ID
						if len(bmsnic.SecurityGroups) != 0 {
							secgrpList := []string{}
							for i := 0; i < len(bmsnic.SecurityGroups); i++ {
								secgrpList = append(secgrpList, string(*(bmsnic.SecurityGroups[i].ID)))
							}
							currentNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
						}
					}
				case "*vpcv1.BareMetalServerNetworkInterfaceByHiperSocket":
					{
						bmsnic := bmsnicintf.(*vpcv1.BareMetalServerNetworkInterfaceByHiperSocket)
						currentNic[isBareMetalServerNicAllowIPSpoofing] = *bmsnic.AllowIPSpoofing
						currentNic[isBareMetalServerNicSubnet] = *bmsnic.Subnet.ID
						if len(bmsnic.SecurityGroups) != 0 {
							secgrpList := []string{}
							for i := 0; i < len(bmsnic.SecurityGroups); i++ {
								secgrpList = append(secgrpList, string(*(bmsnic.SecurityGroups[i].ID)))
							}
							currentNic[isBareMetalServerNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
						}
					}
				}
				interfacesList = append(interfacesList, currentNic)
			}
		}
		l[isBareMetalServerNetworkInterfaces] = interfacesList
		l[isBareMetalServerCreatedAt] = bms.CreatedAt.String()

		//disks
		l[isBareMetalServerResourceType] = *bms.ResourceType
		l[isBareMetalServerStatus] = *bms.Status
		if bms.StatusReasons != nil {
			statusReasonsList := make([]map[string]interface{}, 0)
			for _, sr := range bms.StatusReasons {
				currentSR := map[string]interface{}{}
				if sr.Code != nil && sr.Message != nil {
					currentSR[isBareMetalServerStatusReasonsCode] = *sr.Code
					currentSR[isBareMetalServerStatusReasonsMessage] = *sr.Message
					if sr.MoreInfo != nil {
						currentSR[isBareMetalServerStatusReasonsMoreInfo] = *sr.MoreInfo
					}
					statusReasonsList = append(statusReasonsList, currentSR)
				}
			}
			l[isBareMetalServerStatusReasons] = statusReasonsList
		}
		l[isBareMetalServerVPC] = *bms.VPC.ID
		l[isBareMetalServerZone] = *bms.Zone.Name

		// set keys and image using initialization

		optionsInitialization := &vpcv1.GetBareMetalServerInitializationOptions{
			ID: bms.ID,
		}

		initialization, response, err := sess.GetBareMetalServerInitialization(optionsInitialization)
		if err != nil || initialization == nil {
			log.Printf("[ERROR] Error getting Bare Metal Server (%s) initialization : %s\n%s", *bms.ID, err, response)
		}

		l[isBareMetalServerImage] = *initialization.Image.ID

		keyListList := []string{}
		for i := 0; i < len(initialization.Keys); i++ {
			keyListList = append(keyListList, string(*(initialization.Keys[i].ID)))
		}
		l[isBareMetalServerKeys] = keyListList

		tags, err := flex.GetGlobalTagsUsingCRN(meta, *bms.CRN, "", isBareMetalServerUserTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on get of resource bare metal server (%s) tags: %s", *bms.ID, err)
		}
		l[isBareMetalServerTags] = tags

		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *bms.CRN, "", isBareMetalServerAccessTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on get of resource bare metal server (%s) access tags: %s", *bms.ID, err)
		}
		l[isBareMetalServerAccessTags] = accesstags

		if bms.ResourceGroup != nil {
			l[isBareMetalServerResourceGroup] = *bms.ResourceGroup.ID
		}
		serversInfo = append(serversInfo, l)
	}
	d.SetId(dataSourceIBMISBareMetalServersID(d))
	d.Set(isBareMetalServers, serversInfo)
	return nil
}

// dataSourceIBMISBareMetalServersID returns a reasonable ID for a Bare Metal Servers list.
func dataSourceIBMISBareMetalServersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
