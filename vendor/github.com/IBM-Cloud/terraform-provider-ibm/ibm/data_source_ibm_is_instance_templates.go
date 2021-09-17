// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isInstanceTemplates                     = "templates"
	isInstanceTemplatesFirst                = "first"
	isInstanceTemplatesHref                 = "href"
	isInstanceTemplatesCrn                  = "crn"
	isInstanceTemplatesLimit                = "limit"
	isInstanceTemplatesNext                 = "next"
	isInstanceTemplatesTotalCount           = "total_count"
	isInstanceTemplatesName                 = "name"
	isInstanceTemplatesPortSpeed            = "port_speed"
	isInstanceTemplatesPortType             = "type"
	isInstanceTemplatesPortValue            = "value"
	isInstanceTemplatesDeleteVol            = "delete_volume_on_instance_delete"
	isInstanceTemplatesVol                  = "volume"
	isInstanceTemplatesMemory               = "memory"
	isInstanceTemplatesMemoryValue          = "value"
	isInstanceTemplatesMemoryType           = "type"
	isInstanceTemplatesMemoryValues         = "values"
	isInstanceTemplatesMemoryDefault        = "default"
	isInstanceTemplatesMemoryMin            = "min"
	isInstanceTemplatesMemoryMax            = "max"
	isInstanceTemplatesMemoryStep           = "step"
	isInstanceTemplatesSocketCount          = "socket_count"
	isInstanceTemplatesSocketValue          = "value"
	isInstanceTemplatesSocketType           = "type"
	isInstanceTemplatesSocketValues         = "values"
	isInstanceTemplatesSocketDefault        = "default"
	isInstanceTemplatesSocketMin            = "min"
	isInstanceTemplatesSocketMax            = "max"
	isInstanceTemplatesSocketStep           = "step"
	isInstanceTemplatesVcpuArch             = "vcpu_architecture"
	isInstanceTemplatesVcpuArchType         = "type"
	isInstanceTemplatesVcpuArchValue        = "value"
	isInstanceTemplatesVcpuCount            = "vcpu_count"
	isInstanceTemplatesVcpuCountValue       = "value"
	isInstanceTemplatesVcpuCountType        = "type"
	isInstanceTemplatesVcpuCountValues      = "values"
	isInstanceTemplatesVcpuCountDefault     = "default"
	isInstanceTemplatesVcpuCountMin         = "min"
	isInstanceTemplatesVcpuCountMax         = "max"
	isInstanceTemplatesVcpuCountStep        = "step"
	isInstanceTemplatesStart                = "start"
	isInstanceTemplatesVersion              = "version"
	isInstanceTemplatesGeneration           = "generation"
	isInstanceTemplatesBootVolumeAttachment = "boot_volume_attachment"

	isInstanceTemplateVPC                     = "vpc"
	isInstanceTemplateZone                    = "zone"
	isInstanceTemplateProfile                 = "profile"
	isInstanceTemplateKeys                    = "keys"
	isInstanceTemplateVolumeAttachments       = "volume_attachments"
	isInstanceTemplateNetworkInterfaces       = "network_interfaces"
	isInstanceTemplatePrimaryNetworkInterface = "primary_network_interface"
	isInstanceTemplateNicName                 = "name"
	isInstanceTemplateNicPortSpeed            = "port_speed"
	isInstanceTemplateNicAllowIPSpoofing      = "allow_ip_spoofing"
	isInstanceTemplateNicPrimaryIpv4Address   = "primary_ipv4_address"
	isInstanceTemplateNicSecondaryAddress     = "secondary_addresses"
	isInstanceTemplateNicSecurityGroups       = "security_groups"
	isInstanceTemplateNicSubnet               = "subnet"
	isInstanceTemplateNicFloatingIPs          = "floating_ips"
	isInstanceTemplateUserData                = "user_data"
	isInstanceTemplateGeneration              = "generation"
	isInstanceTemplateImage                   = "image"
	isInstanceTemplateResourceGroup           = "resource_group"
	isInstanceTemplateName                    = "name"
	isInstanceTemplateDeleteVolume            = "delete_volume_on_instance_delete"
	isInstanceTemplateVolAttName              = "name"
	isInstanceTemplateVolAttVolume            = "volume"
)

func dataSourceIBMISInstanceTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstanceTemplatesRead,
		Schema: map[string]*schema.Schema{
			isInstanceTemplates: {
				Type:        schema.TypeList,
				Description: "Collection of instance templates",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplatesName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplatesHref: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplatesCrn: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplateVPC: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplateZone: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplateProfile: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplateKeys: {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						isInstanceTemplateVolumeAttachments: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceTemplatesDeleteVol: {
										Type:     schema.TypeBool,
										Computed: true,
									},
									isInstanceTemplatesName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplatesVol: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplateVolAttVolPrototype: {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isInstanceTemplateVolAttVolIops: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The maximum I/O operations per second (IOPS) for the volume.",
												},
												isInstanceTemplateVolAttVolProfile: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The  globally unique name for the volume profile to use for this volume.",
												},
												isInstanceTemplateVolAttVolCapacity: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The capacity of the volume in gigabytes. The specified minimum and maximum capacity values for creating or updating volumes may expand in the future.",
												},
												isInstanceTemplateVolAttVolEncryptionKey: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.",
												},
											},
										},
									},
								},
							},
						},
						isInstanceTemplatePrimaryNetworkInterface: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceTemplateNicName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplateNicPrimaryIpv4Address: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplateNicSecurityGroups: {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},
									isInstanceTemplateNicSubnet: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						isInstanceTemplateNetworkInterfaces: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceTemplateNicName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplateNicPrimaryIpv4Address: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplateNicSecurityGroups: {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},
									isInstanceTemplateNicSubnet: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						isInstanceTemplateUserData: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplateImage: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceTemplatesBootVolumeAttachment: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceTemplatesDeleteVol: {
										Type:     schema.TypeBool,
										Computed: true,
									},
									isInstanceTemplatesName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplatesVol: {
										Type:     schema.TypeString,
										Computed: true,
									},
									isInstanceTemplateBootSize: {
										Type:     schema.TypeInt,
										Computed: true,
									},
									isInstanceTemplateBootProfile: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						isInstanceTemplateResourceGroup: {
							Type:     schema.TypeString,
							Computed: true,
						},
						"placement_target": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The placement restrictions to use for the virtual server instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this dedicated host.",
									},
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this dedicated host.",
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this dedicated host.",
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

func dataSourceIBMISInstanceTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}
	listInstanceTemplatesOptions := &vpcv1.ListInstanceTemplatesOptions{}
	availableTemplates, _, err := instanceC.ListInstanceTemplates(listInstanceTemplatesOptions)
	if err != nil {
		return err
	}
	templates := make([]map[string]interface{}, 0)
	for _, instTempl := range availableTemplates.Templates {
		template := map[string]interface{}{}
		instance := instTempl.(*vpcv1.InstanceTemplate)
		template["id"] = instance.ID
		template[isInstanceTemplatesHref] = instance.Href
		template[isInstanceTemplatesCrn] = instance.CRN
		template[isInstanceTemplateName] = instance.Name
		template[isInstanceTemplateUserData] = instance.UserData

		if instance.Keys != nil {
			keys := []string{}
			for _, intfc := range instance.Keys {
				instanceKeyIntf := intfc.(*vpcv1.KeyIdentity)
				keys = append(keys, *instanceKeyIntf.ID)
			}
			template[isInstanceTemplateKeys] = keys
		}
		if instance.Profile != nil {
			instanceProfileIntf := instance.Profile
			identity := instanceProfileIntf.(*vpcv1.InstanceProfileIdentity)
			template[isInstanceTemplateProfile] = identity.Name
		}

		if instance.PlacementTarget != nil {
			placementTargetList := []map[string]interface{}{}
			placementTargetMap := dataSourceInstanceTemplateCollectionTemplatesPlacementTargetToMap(*instance.PlacementTarget.(*vpcv1.InstancePlacementTargetPrototype))
			placementTargetList = append(placementTargetList, placementTargetMap)
			template["placement_target"] = placementTargetList
		}

		if instance.PrimaryNetworkInterface != nil {
			interfaceList := make([]map[string]interface{}, 0)
			currentPrimNic := map[string]interface{}{}
			currentPrimNic[isInstanceTemplateNicName] = *instance.PrimaryNetworkInterface.Name
			if instance.PrimaryNetworkInterface.PrimaryIpv4Address != nil {
				currentPrimNic[isInstanceTemplateNicPrimaryIpv4Address] = *instance.PrimaryNetworkInterface.PrimaryIpv4Address
			}
			subInf := instance.PrimaryNetworkInterface.Subnet
			subnetIdentity := subInf.(*vpcv1.SubnetIdentity)
			currentPrimNic[isInstanceTemplateNicSubnet] = *subnetIdentity.ID

			if len(instance.PrimaryNetworkInterface.SecurityGroups) != 0 {
				secgrpList := []string{}
				for i := 0; i < len(instance.PrimaryNetworkInterface.SecurityGroups); i++ {
					secGrpInf := instance.PrimaryNetworkInterface.SecurityGroups[i]
					secGrpIdentity := secGrpInf.(*vpcv1.SecurityGroupIdentity)
					secgrpList = append(secgrpList, string(*secGrpIdentity.ID))
				}
				currentPrimNic[isInstanceTemplateNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
			}
			interfaceList = append(interfaceList, currentPrimNic)
			template[isInstanceTemplatePrimaryNetworkInterface] = interfaceList
		}

		if instance.NetworkInterfaces != nil {
			interfacesList := make([]map[string]interface{}, 0)
			for _, intfc := range instance.NetworkInterfaces {
				currentNic := map[string]interface{}{}
				currentNic[isInstanceTemplateNicName] = *intfc.Name
				if intfc.PrimaryIpv4Address != nil {
					currentNic[isInstanceTemplateNicPrimaryIpv4Address] = *intfc.PrimaryIpv4Address
				}
				//currentNic[isInstanceTemplateNicAllowIpSpoofing] = intfc.AllowIpSpoofing
				subInf := intfc.Subnet
				subnetIdentity := subInf.(*vpcv1.SubnetIdentity)
				currentNic[isInstanceTemplateNicSubnet] = *subnetIdentity.ID
				if len(intfc.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(intfc.SecurityGroups); i++ {
						secGrpInf := intfc.SecurityGroups[i]
						secGrpIdentity := secGrpInf.(*vpcv1.SecurityGroupIdentity)
						secgrpList = append(secgrpList, string(*secGrpIdentity.ID))
					}
					currentNic[isInstanceTemplateNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
				}

				interfacesList = append(interfacesList, currentNic)
			}
			template[isInstanceTemplateNetworkInterfaces] = interfacesList
		}

		if instance.Image != nil {
			imageInf := instance.Image
			imageIdentity := imageInf.(*vpcv1.ImageIdentity)
			template[isInstanceTemplateImage] = imageIdentity.ID
		}

		if instance.VPC != nil {
			vpcInf := instance.VPC
			vpcRef := vpcInf.(*vpcv1.VPCIdentity)
			template[isInstanceTemplateVPC] = vpcRef.ID
		}

		if instance.Zone != nil {
			zoneInf := instance.Zone
			zone := zoneInf.(*vpcv1.ZoneIdentity)
			template[isInstanceTemplateZone] = zone.Name
		}

		interfacesList := make([]map[string]interface{}, 0)
		if instance.VolumeAttachments != nil {
			for _, volume := range instance.VolumeAttachments {
				volumeAttach := map[string]interface{}{}
				volumeAttach[isInstanceTemplateVolAttName] = *volume.Name
				volumeAttach[isInstanceTemplateDeleteVolume] = *volume.DeleteVolumeOnInstanceDelete
				volumeIntf := volume.Volume
				volumeInst := volumeIntf.(*vpcv1.VolumeAttachmentVolumePrototypeInstanceContext)
				newVolumeArr := []map[string]interface{}{}
				newVolume := map[string]interface{}{}

				if volumeInst.ID != nil {
					volumeAttach[isInstanceTemplateVolAttVolume] = *volumeInst.ID
				}

				if volumeInst.Capacity != nil {
					newVolume[isInstanceTemplateVolAttVolCapacity] = *volumeInst.Capacity
				}
				if volumeInst.Profile != nil {
					profile := volumeInst.Profile.(*vpcv1.VolumeProfileIdentity)
					newVolume[isInstanceTemplateVolAttVolProfile] = profile.Name
				}

				if volumeInst.Iops != nil {
					newVolume[isInstanceTemplateVolAttVolIops] = *volumeInst.Iops
				}
				if volumeInst.EncryptionKey != nil {
					encryptionKey := volumeInst.EncryptionKey.(*vpcv1.EncryptionKeyIdentity)
					newVolume[isInstanceTemplateVolAttVolEncryptionKey] = *encryptionKey.CRN
				}
				newVolumeArr = append(newVolumeArr, newVolume)
				volumeAttach[isInstanceTemplateVolAttVolPrototype] = newVolumeArr

				interfacesList = append(interfacesList, volumeAttach)
			}
			template[isInstanceTemplateVolumeAttachments] = interfacesList
		}

		if instance.BootVolumeAttachment != nil {
			bootVolList := make([]map[string]interface{}, 0)
			bootVol := map[string]interface{}{}

			bootVol[isInstanceTemplatesDeleteVol] = *instance.BootVolumeAttachment.DeleteVolumeOnInstanceDelete
			if instance.BootVolumeAttachment.Volume != nil {
				volumeIntf := instance.BootVolumeAttachment.Volume
				bootVol[isInstanceTemplatesName] = volumeIntf.Name
				bootVol[isInstanceTemplatesVol] = volumeIntf.Name
				bootVol[isInstanceTemplateBootSize] = volumeIntf.Capacity
				if instance.BootVolumeAttachment.Volume.Profile != nil {
					volProfIntf := instance.BootVolumeAttachment.Volume.Profile
					volProfInst := volProfIntf.(*vpcv1.VolumeProfileIdentity)
					bootVol[isInstanceTemplateBootProfile] = volProfInst.Name
				}
			}
			bootVolList = append(bootVolList, bootVol)
			template[isInstanceTemplatesBootVolumeAttachment] = bootVolList
		}

		if instance.ResourceGroup != nil {
			rg := instance.ResourceGroup
			template[isInstanceTemplateResourceGroup] = rg.ID
		}

		templates = append(templates, template)
	}
	d.SetId(dataSourceIBMISInstanceTemplatesID(d))
	d.Set(isInstanceTemplates, templates)
	return nil
}

// dataSourceIBMISInstanceTemplatesID returns a reasonable ID for a instance templates list.
func dataSourceIBMISInstanceTemplatesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceInstanceTemplateCollectionTemplatesPlacementTargetToMap(placementTargetItem vpcv1.InstancePlacementTargetPrototype) (placementTargetMap map[string]interface{}) {
	placementTargetMap = map[string]interface{}{}

	if placementTargetItem.ID != nil {
		placementTargetMap["id"] = placementTargetItem.ID
	}
	if placementTargetItem.CRN != nil {
		placementTargetMap["crn"] = placementTargetItem.CRN
	}
	if placementTargetItem.Href != nil {
		placementTargetMap["href"] = placementTargetItem.Href
	}

	return placementTargetMap
}
