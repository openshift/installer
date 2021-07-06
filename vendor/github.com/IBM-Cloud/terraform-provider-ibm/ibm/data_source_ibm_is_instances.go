// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isInstances = "instances"
)

func dataSourceIBMISInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstancesRead,

		Schema: map[string]*schema.Schema{
			"vpc_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc"},
				Description:   "Name of the vpc to filter the instances attached to it",
			},

			"vpc": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc_name"},
				Description:   "VPC ID to filter the instances attached to it",
			},

			isInstances: {
				Type:        schema.TypeList,
				Description: "List of instances",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance id",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance memory",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance status",
						},
						"resource_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance resource group",
						},
						"vpc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "vpc attached to the instance",
						},
						"boot_volume": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance Boot Volume",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Boot volume id",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Boot volume name",
									},
									"device": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Boot volume device",
									},
									"volume_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Boot volume's volume id",
									},
									"volume_crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Boot volume's volume CRN",
									},
								},
							},
						},

						"volume_attachments": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance Volume Attachments",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance volume Attachment id",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance volume Attachment name",
									},
									"volume_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance volume Attachment's volume id",
									},
									"volume_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance volume Attachment's volume name",
									},
									"volume_crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance volume Attachment's volume CRN",
									},
								},
							},
						},

						"primary_network_interface": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance Primary Network Interface",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Primary Network interface id",
									},
									isInstanceNicName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Primary Network interface name",
									},
									isInstanceNicPrimaryIpv4Address: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Primary Network interface IPV4 Address",
									},
									isInstanceNicSecurityGroups: {
										Type:        schema.TypeSet,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
										Description: "Instance Primary Network interface security groups",
									},
									isInstanceNicSubnet: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Primary Network interface subnet",
									},
								},
							},
						},
						"network_interfaces": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance Network Interfaces",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Network interface id",
									},
									isInstanceNicName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Network interface name",
									},
									isInstanceNicPrimaryIpv4Address: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Network interface IPV4 Address",
									},
									isInstanceNicSecurityGroups: {
										Type:        schema.TypeSet,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
										Description: "Instance Network interface security groups",
									},
									isInstanceNicSubnet: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Network interface subnet",
									},
								},
							},
						},
						"profile": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Profile",
						},
						"vcpu": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance vcpu",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"architecture": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance vcpu architecture",
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance vcpu count",
									},
								},
							},
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance zone",
						},
						"image": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Instance Image",
						},
						isInstanceDisks: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of the instance's disks.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"created_at": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time that the disk was created.",
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
									"interface_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The disk interface used for attaching the disk.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
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
									"size": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The size of the disk in GB (gigabytes).",
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

func dataSourceIBMISInstancesRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	if userDetails.generation == 1 {
		err := classicInstancesList(d, meta)
		if err != nil {
			return err
		}
	} else {
		err := instancesList(d, meta)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicInstancesList(d *schema.ResourceData, meta interface{}) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcclassicv1.Instance{}
	for {
		listInstancesOptions := &vpcclassicv1.ListInstancesOptions{}
		if start != "" {
			listInstancesOptions.Start = &start
		}
		instances, response, err := sess.ListInstances(listInstancesOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching Instances %s\n%s", err, response)
		}
		start = GetNext(instances.Next)
		allrecs = append(allrecs, instances.Instances...)
		if start == "" {
			break
		}
	}
	instancesInfo := make([]map[string]interface{}, 0)
	for _, instance := range allrecs {
		id := *instance.ID
		l := map[string]interface{}{}
		l["id"] = id
		l["name"] = *instance.Name
		l["memory"] = *instance.Memory
		l["status"] = *instance.Status
		l["resource_group"] = *instance.ResourceGroup.ID
		l["vpc"] = *instance.VPC.ID

		if instance.BootVolumeAttachment != nil {
			bootVolList := make([]map[string]interface{}, 0)
			bootVol := map[string]interface{}{}
			bootVol["id"] = *instance.BootVolumeAttachment.ID
			bootVol["name"] = *instance.BootVolumeAttachment.Name
			if instance.BootVolumeAttachment.Device != nil {
				bootVol["device"] = *instance.BootVolumeAttachment.Device.ID
			}
			if instance.BootVolumeAttachment.Volume != nil {
				bootVol["volume_id"] = *instance.BootVolumeAttachment.Volume.ID
				bootVol["volume_crn"] = *instance.BootVolumeAttachment.Volume.CRN
			}
			bootVolList = append(bootVolList, bootVol)
			l["boot_volume"] = bootVolList
		}

		if instance.VolumeAttachments != nil {
			volList := make([]map[string]interface{}, 0)
			for _, volume := range instance.VolumeAttachments {
				vol := map[string]interface{}{}
				if volume.Volume != nil {
					vol["id"] = *volume.ID
					vol["volume_id"] = *volume.Volume.ID
					vol["name"] = *volume.Name
					vol["volume_name"] = *volume.Volume.Name
					vol["volume_crn"] = *volume.Volume.CRN
					volList = append(volList, vol)
				}
			}
			l["volume_attachments"] = volList
		}

		if instance.PrimaryNetworkInterface != nil {
			primaryNicList := make([]map[string]interface{}, 0)
			currentPrimNic := map[string]interface{}{}
			currentPrimNic["id"] = *instance.PrimaryNetworkInterface.ID
			currentPrimNic[isInstanceNicName] = *instance.PrimaryNetworkInterface.Name
			currentPrimNic[isInstanceNicPrimaryIpv4Address] = *instance.PrimaryNetworkInterface.PrimaryIpv4Address
			getnicoptions := &vpcclassicv1.GetInstanceNetworkInterfaceOptions{
				InstanceID: &id,
				ID:         instance.PrimaryNetworkInterface.ID,
			}
			insnic, response, err := sess.GetInstanceNetworkInterface(getnicoptions)
			if err != nil {
				return fmt.Errorf("Error getting network interfaces attached to the instance %s\n%s", err, response)
			}
			currentPrimNic[isInstanceNicSubnet] = *insnic.Subnet.ID
			if len(insnic.SecurityGroups) != 0 {
				secgrpList := []string{}
				for i := 0; i < len(insnic.SecurityGroups); i++ {
					secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
				}
				currentPrimNic[isInstanceNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
			}

			primaryNicList = append(primaryNicList, currentPrimNic)
			l["primary_network_interface"] = primaryNicList
		}

		if instance.NetworkInterfaces != nil {
			interfacesList := make([]map[string]interface{}, 0)
			for _, intfc := range instance.NetworkInterfaces {
				if *intfc.ID != *instance.PrimaryNetworkInterface.ID {
					currentNic := map[string]interface{}{}
					currentNic["id"] = *intfc.ID
					currentNic[isInstanceNicName] = *intfc.Name
					currentNic[isInstanceNicPrimaryIpv4Address] = *intfc.PrimaryIpv4Address
					getnicoptions := &vpcclassicv1.GetInstanceNetworkInterfaceOptions{
						InstanceID: &id,
						ID:         intfc.ID,
					}
					insnic, response, err := sess.GetInstanceNetworkInterface(getnicoptions)
					if err != nil {
						return fmt.Errorf("Error getting network interfaces attached to the instance %s\n%s", err, response)
					}
					currentNic[isInstanceNicSubnet] = *insnic.Subnet.ID
					if len(insnic.SecurityGroups) != 0 {
						secgrpList := []string{}
						for i := 0; i < len(insnic.SecurityGroups); i++ {
							secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
						}
						currentNic[isInstanceNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
					}
					interfacesList = append(interfacesList, currentNic)
				}
			}
			l["network_interfaces"] = interfacesList
		}

		l["profile"] = *instance.Profile.Name

		cpuList := make([]map[string]interface{}, 0)
		if instance.Vcpu != nil {
			currentCPU := map[string]interface{}{}
			currentCPU["architecture"] = *instance.Vcpu.Architecture
			currentCPU["count"] = *instance.Vcpu.Count
			cpuList = append(cpuList, currentCPU)
		}
		l["vcpu"] = cpuList

		l["zone"] = *instance.Zone.Name
		if instance.Image != nil {
			l["image"] = *instance.Image.ID
		}
		instancesInfo = append(instancesInfo, l)
	}
	d.SetId(dataSourceIBMISInstancesID(d))
	d.Set(isInstances, instancesInfo)
	return nil
}

func instancesList(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	var vpcName, vpcID string

	if vpc, ok := d.GetOk("vpc_name"); ok {
		vpcName = vpc.(string)
	}

	if vpc, ok := d.GetOk("vpc"); ok {
		vpcID = vpc.(string)
	}

	start := ""
	allrecs := []vpcv1.Instance{}
	for {
		listInstancesOptions := &vpcv1.ListInstancesOptions{}
		if start != "" {
			listInstancesOptions.Start = &start
		}

		if vpcName != "" {
			listInstancesOptions.VPCName = &vpcName
		}
		if vpcID != "" {
			listInstancesOptions.VPCID = &vpcID
		}

		instances, response, err := sess.ListInstances(listInstancesOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching Instances %s\n%s", err, response)
		}
		start = GetNext(instances.Next)
		allrecs = append(allrecs, instances.Instances...)
		if start == "" {
			break
		}
	}
	instancesInfo := make([]map[string]interface{}, 0)
	for _, instance := range allrecs {
		id := *instance.ID
		l := map[string]interface{}{}
		l["id"] = id
		l["name"] = *instance.Name
		l["memory"] = *instance.Memory
		l["status"] = *instance.Status
		l["resource_group"] = *instance.ResourceGroup.ID
		l["vpc"] = *instance.VPC.ID

		if instance.BootVolumeAttachment != nil {
			bootVolList := make([]map[string]interface{}, 0)
			bootVol := map[string]interface{}{}
			bootVol["id"] = *instance.BootVolumeAttachment.ID
			bootVol["name"] = *instance.BootVolumeAttachment.Name
			if instance.BootVolumeAttachment.Device != nil {
				bootVol["device"] = *instance.BootVolumeAttachment.Device.ID
			}
			if instance.BootVolumeAttachment.Volume != nil {
				bootVol["volume_id"] = *instance.BootVolumeAttachment.Volume.ID
				bootVol["volume_crn"] = *instance.BootVolumeAttachment.Volume.CRN
			}
			bootVolList = append(bootVolList, bootVol)
			l["boot_volume"] = bootVolList
		}

		if instance.VolumeAttachments != nil {
			volList := make([]map[string]interface{}, 0)
			for _, volume := range instance.VolumeAttachments {
				vol := map[string]interface{}{}
				if volume.Volume != nil {
					vol["id"] = *volume.ID
					vol["volume_id"] = *volume.Volume.ID
					vol["name"] = *volume.Name
					vol["volume_name"] = *volume.Volume.Name
					vol["volume_crn"] = *volume.Volume.CRN
					volList = append(volList, vol)
				}
			}
			l["volume_attachments"] = volList
		}

		if instance.PrimaryNetworkInterface != nil {
			primaryNicList := make([]map[string]interface{}, 0)
			currentPrimNic := map[string]interface{}{}
			currentPrimNic["id"] = *instance.PrimaryNetworkInterface.ID
			currentPrimNic[isInstanceNicName] = *instance.PrimaryNetworkInterface.Name
			currentPrimNic[isInstanceNicPrimaryIpv4Address] = *instance.PrimaryNetworkInterface.PrimaryIpv4Address
			getnicoptions := &vpcv1.GetInstanceNetworkInterfaceOptions{
				InstanceID: &id,
				ID:         instance.PrimaryNetworkInterface.ID,
			}
			insnic, response, err := sess.GetInstanceNetworkInterface(getnicoptions)
			if err != nil {
				return fmt.Errorf("Error getting network interfaces attached to the instance %s\n%s", err, response)
			}
			currentPrimNic[isInstanceNicSubnet] = *insnic.Subnet.ID
			if len(insnic.SecurityGroups) != 0 {
				secgrpList := []string{}
				for i := 0; i < len(insnic.SecurityGroups); i++ {
					secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
				}
				currentPrimNic[isInstanceNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
			}

			primaryNicList = append(primaryNicList, currentPrimNic)
			l["primary_network_interface"] = primaryNicList
		}

		if instance.NetworkInterfaces != nil {
			interfacesList := make([]map[string]interface{}, 0)
			for _, intfc := range instance.NetworkInterfaces {
				if *intfc.ID != *instance.PrimaryNetworkInterface.ID {
					currentNic := map[string]interface{}{}
					currentNic["id"] = *intfc.ID
					currentNic[isInstanceNicName] = *intfc.Name
					currentNic[isInstanceNicPrimaryIpv4Address] = *intfc.PrimaryIpv4Address
					getnicoptions := &vpcv1.GetInstanceNetworkInterfaceOptions{
						InstanceID: &id,
						ID:         intfc.ID,
					}
					insnic, response, err := sess.GetInstanceNetworkInterface(getnicoptions)
					if err != nil {
						return fmt.Errorf("Error getting network interfaces attached to the instance %s\n%s", err, response)
					}
					currentNic[isInstanceNicSubnet] = *insnic.Subnet.ID
					if len(insnic.SecurityGroups) != 0 {
						secgrpList := []string{}
						for i := 0; i < len(insnic.SecurityGroups); i++ {
							secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
						}
						currentNic[isInstanceNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
					}
					interfacesList = append(interfacesList, currentNic)
				}
			}
			l["network_interfaces"] = interfacesList
		}

		l["profile"] = *instance.Profile.Name

		cpuList := make([]map[string]interface{}, 0)
		if instance.Vcpu != nil {
			currentCPU := map[string]interface{}{}
			currentCPU["architecture"] = *instance.Vcpu.Architecture
			currentCPU["count"] = *instance.Vcpu.Count
			cpuList = append(cpuList, currentCPU)
		}
		l["vcpu"] = cpuList

		l["zone"] = *instance.Zone.Name
		if instance.Image != nil {
			l["image"] = *instance.Image.ID
		}

		if instance.Disks != nil {
			l[isInstanceDisks] = dataSourceInstanceFlattenDisks(instance.Disks)
		}

		instancesInfo = append(instancesInfo, l)
	}
	d.SetId(dataSourceIBMISInstancesID(d))
	d.Set(isInstances, instancesInfo)
	return nil
}

// dataSourceIBMISInstancesID returns a reasonable ID for a Instance list.
func dataSourceIBMISInstancesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
