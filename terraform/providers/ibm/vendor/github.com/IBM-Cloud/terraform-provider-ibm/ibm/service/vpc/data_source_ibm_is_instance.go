// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/ScaleFT/sshkeys"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/ssh"
)

const (
	isInstancePEM                       = "private_key"
	isInstancePassphrase                = "passphrase"
	isInstanceInitPassword              = "password"
	isInstanceInitKeys                  = "keys"
	isInstanceNicPrimaryIP              = "primary_ip"
	isInstanceNicReservedIpAddress      = "address"
	isInstanceNicReservedIpHref         = "href"
	isInstanceNicReservedIpAutoDelete   = "auto_delete"
	isInstanceNicReservedIpName         = "name"
	isInstanceNicReservedIpId           = "reserved_ip"
	isInstanceNicReservedIpResourceType = "resource_type"
)

func DataSourceIBMISInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstanceRead,

		Schema: map[string]*schema.Schema{

			isInstanceAvailablePolicyHostFailure: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The availability policy to use for this virtual server instance. The action to perform if the compute host experiences a failure.",
			},

			isInstanceName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name",
			},

			isInstanceMetadataServiceEnabled: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the metadata service endpoint is available to the virtual server instance",
			},
			isInstanceMetadataService: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The metadata service configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceMetadataServiceEnabled1: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the metadata service endpoint will be available to the virtual server instance",
						},

						isInstanceMetadataServiceProtocol: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The communication protocol to use for the metadata service endpoint. Applies only when the metadata service is enabled.",
						},

						isInstanceMetadataServiceRespHopLimit: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The hop limit (IP time to live) for IP response packets from the metadata service",
						},
					},
				},
			},
			isInstancePEM: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance Private Key file",
			},

			isInstancePassphrase: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Passphrase for Instance Private Key file",
			},

			isInstanceInitPassword: {
				Type:        schema.TypeString,
				Sensitive:   true,
				Computed:    true,
				Description: "password for Windows Instance",
			},

			isInstanceInitKeys: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance keys",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance key id",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance key name",
						},
					},
				},
			},

			isInstanceVPC: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC id",
			},

			isInstanceZone: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone name",
			},

			isInstanceProfile: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Profile info",
			},

			isInstanceTotalVolumeBandwidth: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of bandwidth (in megabits per second) allocated exclusively to instance storage volumes",
			},

			isInstanceBandwidth: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total bandwidth (in megabits per second) shared across the instance's network interfaces and storage volumes",
			},

			isInstanceTotalNetworkBandwidth: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of bandwidth (in megabits per second) allocated exclusively to instance network interfaces.",
			},

			isInstanceLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the virtual server instance.",
			},
			isInstanceLifecycleReasons: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current lifecycle_state (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceLifecycleReasonsCode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the reason for this lifecycle state.",
						},

						isInstanceLifecycleReasonsMessage: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this lifecycle state.",
						},

						isInstanceLifecycleReasonsMoreInfo: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this lifecycle state.",
						},
					},
				},
			},

			isInstanceTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "list of tags for the instance",
			},
			isInstanceAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "list of access tags for the instance",
			},
			isInstanceBootVolume: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance Boot Volume",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume id",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume name",
						},
						"device": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume device",
						},
						"volume_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume's volume id",
						},
						"volume_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume's volume name",
						},
						"volume_crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume's volume CRN",
						},
					},
				},
			},

			isInstanceCatalogOffering: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The catalog offering or offering version to use when provisioning this virtual server instance. If an offering is specified, the latest version of that offering will be used. The specified offering or offering version may be in a different account in the same enterprise, subject to IAM policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceCatalogOfferingOfferingCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifies a catalog offering by a unique CRN property",
						},
						isInstanceCatalogOfferingVersionCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifies a version of a catalog offering by a unique CRN property",
						},
					},
				},
			},

			isInstanceVolumeAttachments: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance Volume Attachments",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Volume Attachment id",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Volume Attachment name",
						},
						"volume_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume's volume id",
						},
						"volume_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume's volume name",
						},
						"volume_crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Boot Volume's volume CRN",
						},
					},
				},
			},

			isInstancePrimaryNetworkInterface: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Primary Network interface info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Primary Network Interface id",
						},
						isInstanceNicName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Primary Network Interface name",
						},
						isInstanceNicPortSpeed: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance Primary Network Interface port speed",
						},
						isInstanceNicPrimaryIpv4Address: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Primary Network Interface IPV4 Address",
						},
						isInstanceNicPrimaryIP: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The primary IP address to bind to the network interface. This can be specified using an existing reserved IP, or a prototype object for a new reserved IP.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceNicReservedIpAddress: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address to reserve, which must not already be reserved on the subnet.",
									},
									isInstanceNicReservedIpHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP",
									},
									isInstanceNicReservedIpName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isInstanceNicReservedIpId: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifies a reserved IP by a unique property.",
									},
									isInstanceNicReservedIpResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type",
									},
								},
							},
						},
						isInstanceNicSecurityGroups: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "Instance Primary Network Interface Security groups",
						},
						isInstanceNicSubnet: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Primary Network Interface subnet",
						},
					},
				},
			},

			isInstanceNetworkInterfaces: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance Network interface info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Network Interface id",
						},
						isInstanceNicName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Network Interface name",
						},
						isInstanceNicPrimaryIpv4Address: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Network Interface IPV4 Address",
						},
						isInstanceNicPrimaryIP: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The primary IP address to bind to the network interface. This can be specified using an existing reserved IP, or a prototype object for a new reserved IP.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceNicReservedIpAddress: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address to reserve, which must not already be reserved on the subnet.",
									},
									isInstanceNicReservedIpHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP",
									},
									isInstanceNicReservedIpName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isInstanceNicReservedIpId: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifies a reserved IP by a unique property.",
									},
									isInstanceNicReservedIpResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type",
									},
								},
							},
						},
						isInstanceNicSecurityGroups: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "Instance Network Interface Security Groups",
						},
						isInstanceNicSubnet: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Network Interface subnet",
						},
					},
				},
			},

			isInstanceImage: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance Image",
			},

			isInstanceVolumes: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of volumes",
			},

			isInstanceResourceGroup: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance resource group",
			},

			isInstanceCPU: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance vCPU",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceCPUArch: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance vCPU Architecture",
						},
						isInstanceCPUCount: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance vCPU count",
						},
						// Added for AMD support, manufacturer details.
						isInstanceCPUManufacturer: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance vCPU Manufacturer",
						},
					},
				},
			},

			isInstanceGpu: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance GPU",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceGpuCount: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance GPU Count",
						},
						isInstanceGpuMemory: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance GPU Memory",
						},
						isInstanceGpuManufacturer: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance GPU Manufacturer",
						},
						isInstanceGpuModel: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance GPU Model",
						},
					},
				},
			},

			isInstanceMemory: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance memory",
			},

			"numa_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of NUMA nodes this virtual server instance is provisioned on. This property may be absent if the instance's `status` is not `running`.",
			},

			isInstanceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "instance status",
			},

			isInstanceStatusReasons: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current status (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceStatusReasonsCode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason",
						},

						isInstanceStatusReasonsMessage: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason",
						},

						isInstanceStatusReasonsMoreInfo: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about this status reason",
						},
					},
				},
			},

			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			flex.ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			IsInstanceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			flex.ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
			isInstanceDisks: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of the instance's disks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the disk was created.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this instance disk.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this instance disk.",
						},
						"interface_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The disk interface used for attaching the disk.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this disk.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the disk in GB (gigabytes).",
						},
					},
				},
			},
			"placement_target": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The placement restrictions for the virtual server instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this dedicated host group.",
						},
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this dedicated host group.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this dedicated host group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this dedicated host group. If unspecified, the name will be a hyphenated list of randomly-selected words.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISInstanceRead(d *schema.ResourceData, meta interface{}) error {

	name := d.Get(isInstanceName).(string)

	err := instanceGetByName(d, meta, name)
	if err != nil {
		return err
	}
	return nil
}

func instanceGetByName(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	listInstancesOptions := &vpcv1.ListInstancesOptions{
		Name: &name,
	}

	instances, response, err := sess.ListInstances(listInstancesOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Fetching Instances %s\n%s", err, response)
	}
	allrecs := instances.Instances

	if len(allrecs) == 0 {
		return fmt.Errorf("[ERROR] No Instance found with name %s", name)
	}
	instance := allrecs[0]
	d.SetId(*instance.ID)
	id := *instance.ID

	// catalog
	if instance.CatalogOffering != nil {
		versionCrn := *instance.CatalogOffering.Version.CRN
		catalogList := make([]map[string]interface{}, 0)
		catalogMap := map[string]interface{}{}
		catalogMap[isInstanceCatalogOfferingVersionCrn] = versionCrn
		catalogList = append(catalogList, catalogMap)
		d.Set(isInstanceCatalogOffering, catalogList)
	}

	d.Set(isInstanceName, *instance.Name)
	if instance.Profile != nil {
		d.Set(isInstanceProfile, *instance.Profile.Name)
	}
	if instance.MetadataService != nil {
		d.Set(isInstanceMetadataServiceEnabled, instance.MetadataService.Enabled)

		metadataService := []map[string]interface{}{}
		metadataServiceMap := map[string]interface{}{}

		metadataServiceMap[isInstanceMetadataServiceEnabled1] = instance.MetadataService.Enabled
		if instance.MetadataService.Protocol != nil {
			metadataServiceMap[isInstanceMetadataServiceProtocol] = instance.MetadataService.Protocol
		}
		if instance.MetadataService.ResponseHopLimit != nil {
			metadataServiceMap[isInstanceMetadataServiceRespHopLimit] = instance.MetadataService.ResponseHopLimit
		}

		metadataService = append(metadataService, metadataServiceMap)
		d.Set(isInstanceMetadataService, metadataService)
	}

	if instance.AvailabilityPolicy != nil && instance.AvailabilityPolicy.HostFailure != nil {
		d.Set(isInstanceAvailablePolicyHostFailure, *instance.AvailabilityPolicy.HostFailure)
	}
	cpuList := make([]map[string]interface{}, 0)
	if instance.Vcpu != nil {
		currentCPU := map[string]interface{}{}
		currentCPU[isInstanceCPUArch] = *instance.Vcpu.Architecture
		currentCPU[isInstanceCPUCount] = *instance.Vcpu.Count
		currentCPU[isInstanceCPUManufacturer] = *instance.Vcpu.Manufacturer // Added for AMD support, manufacturer details.
		cpuList = append(cpuList, currentCPU)
	}
	d.Set(isInstanceCPU, cpuList)

	if instance.PlacementTarget != nil {
		placementTargetMap := resourceIbmIsInstanceInstancePlacementToMap(*instance.PlacementTarget.(*vpcv1.InstancePlacementTarget))
		d.Set("placement_target", []map[string]interface{}{placementTargetMap})
	}

	d.Set(isInstanceMemory, *instance.Memory)
	if instance.NumaCount != nil {
		d.Set("numa_count", *instance.NumaCount)
	}
	gpuList := make([]map[string]interface{}, 0)
	if instance.Gpu != nil {
		currentGpu := map[string]interface{}{}
		currentGpu[isInstanceGpuManufacturer] = instance.Gpu.Manufacturer
		currentGpu[isInstanceGpuModel] = instance.Gpu.Model
		currentGpu[isInstanceGpuCount] = instance.Gpu.Count
		currentGpu[isInstanceGpuMemory] = instance.Gpu.Memory
		gpuList = append(gpuList, currentGpu)
		d.Set(isInstanceGpu, gpuList)
	}

	if instance.Bandwidth != nil {
		d.Set(isInstanceBandwidth, int(*instance.Bandwidth))
	}

	if instance.TotalNetworkBandwidth != nil {
		d.Set(isInstanceTotalNetworkBandwidth, int(*instance.TotalNetworkBandwidth))
	}

	if instance.TotalVolumeBandwidth != nil {
		d.Set(isInstanceTotalVolumeBandwidth, int(*instance.TotalVolumeBandwidth))
	}

	if instance.Disks != nil {
		d.Set(isInstanceDisks, dataSourceInstanceFlattenDisks(instance.Disks))
	}

	if instance.PrimaryNetworkInterface != nil {
		primaryNicList := make([]map[string]interface{}, 0)
		currentPrimNic := map[string]interface{}{}
		currentPrimNic["id"] = *instance.PrimaryNetworkInterface.ID
		currentPrimNic[isInstanceNicName] = *instance.PrimaryNetworkInterface.Name

		// reserved ip changes
		primaryIpList := make([]map[string]interface{}, 0)
		currentPrimIp := map[string]interface{}{}
		if instance.PrimaryNetworkInterface.PrimaryIP.Address != nil {
			currentPrimNic[isInstanceNicPrimaryIpv4Address] = *instance.PrimaryNetworkInterface.PrimaryIP.Address
			currentPrimIp[isInstanceNicReservedIpAddress] = *instance.PrimaryNetworkInterface.PrimaryIP.Address
		}
		if instance.PrimaryNetworkInterface.PrimaryIP.Href != nil {
			currentPrimIp[isInstanceNicReservedIpHref] = *instance.PrimaryNetworkInterface.PrimaryIP.Href
		}
		if instance.PrimaryNetworkInterface.PrimaryIP.Name != nil {
			currentPrimIp[isInstanceNicReservedIpName] = *instance.PrimaryNetworkInterface.PrimaryIP.Name
		}
		if instance.PrimaryNetworkInterface.PrimaryIP.ID != nil {
			currentPrimIp[isInstanceNicReservedIpId] = *instance.PrimaryNetworkInterface.PrimaryIP.ID
		}
		if instance.PrimaryNetworkInterface.PrimaryIP.ResourceType != nil {
			currentPrimIp[isInstanceNicReservedIpResourceType] = *instance.PrimaryNetworkInterface.PrimaryIP.ResourceType
		}
		primaryIpList = append(primaryIpList, currentPrimIp)
		currentPrimNic[isInstanceNicPrimaryIP] = primaryIpList

		getnicoptions := &vpcv1.GetInstanceNetworkInterfaceOptions{
			InstanceID: &id,
			ID:         instance.PrimaryNetworkInterface.ID,
		}
		insnic, response, err := sess.GetInstanceNetworkInterface(getnicoptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error getting network interfaces attached to the instance %s\n%s", err, response)
		}
		if insnic.PortSpeed != nil {
			currentPrimNic[isInstanceNicPortSpeed] = *insnic.PortSpeed
		}
		currentPrimNic[isInstanceNicSubnet] = *insnic.Subnet.ID
		if len(insnic.SecurityGroups) != 0 {
			secgrpList := []string{}
			for i := 0; i < len(insnic.SecurityGroups); i++ {
				secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
			}
			currentPrimNic[isInstanceNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
		}

		primaryNicList = append(primaryNicList, currentPrimNic)
		d.Set(isInstancePrimaryNetworkInterface, primaryNicList)
	}

	if instance.NetworkInterfaces != nil {
		interfacesList := make([]map[string]interface{}, 0)
		for _, intfc := range instance.NetworkInterfaces {
			if *intfc.ID != *instance.PrimaryNetworkInterface.ID {
				currentNic := map[string]interface{}{}
				currentNic["id"] = *intfc.ID
				currentNic[isInstanceNicName] = *intfc.Name

				// reserved ip changes
				primaryIpList := make([]map[string]interface{}, 0)
				currentPrimIp := map[string]interface{}{}
				if intfc.PrimaryIP.Address != nil {
					currentPrimIp[isInstanceNicReservedIpAddress] = *intfc.PrimaryIP.Address
					currentNic[isInstanceNicPrimaryIpv4Address] = *intfc.PrimaryIP.Address
				}
				if intfc.PrimaryIP.Href != nil {
					currentPrimIp[isInstanceNicReservedIpHref] = *intfc.PrimaryIP.Href
				}
				if intfc.PrimaryIP.Name != nil {
					currentPrimIp[isInstanceNicReservedIpName] = *intfc.PrimaryIP.Name
				}
				if intfc.PrimaryIP.ID != nil {
					currentPrimIp[isInstanceNicReservedIpId] = *intfc.PrimaryIP.ID
				}
				if intfc.PrimaryIP.ResourceType != nil {
					currentPrimIp[isInstanceNicReservedIpResourceType] = *intfc.PrimaryIP.ResourceType
				}
				primaryIpList = append(primaryIpList, currentPrimIp)
				currentNic[isInstanceNicPrimaryIP] = primaryIpList

				getnicoptions := &vpcv1.GetInstanceNetworkInterfaceOptions{
					InstanceID: &id,
					ID:         intfc.ID,
				}
				insnic, response, err := sess.GetInstanceNetworkInterface(getnicoptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error getting network interfaces attached to the instance %s\n%s", err, response)
				}
				currentNic[isInstanceNicSubnet] = *insnic.Subnet.ID
				if len(insnic.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(insnic.SecurityGroups); i++ {
						secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
					}
					currentNic[isInstanceNicSecurityGroups] = flex.NewStringSet(schema.HashString, secgrpList)
				}
				interfacesList = append(interfacesList, currentNic)

			}
		}

		d.Set(isInstanceNetworkInterfaces, interfacesList)
	}

	var rsaKey *rsa.PrivateKey
	if instance.Image != nil {
		d.Set(isInstanceImage, *instance.Image.ID)
		image := *instance.Image.Name
		res := strings.Contains(image, "windows")
		if res {
			if privatekey, ok := d.GetOk(isInstancePEM); ok {
				keyFlag := privatekey.(string)
				keybytes := []byte(keyFlag)

				if keyFlag != "" {
					block, err := pem.Decode(keybytes)
					if block == nil {
						return fmt.Errorf("[ERROR] Failed to load the private key from the given key contents. Instead of the key file path, please make sure the private key is pem format (%v)", err)
					}
					isEncrypted := false
					if block.Type == "OPENSSH PRIVATE KEY" {
						var err error
						isEncrypted, err = isOpenSSHPrivKeyEncrypted(block.Bytes)
						if err != nil {
							return fmt.Errorf("[ERROR] Failed to check if the provided open ssh key is encrypted or not %s", err)
						}
					} else {
						isEncrypted = x509.IsEncryptedPEMBlock(block)
					}
					passphrase := ""
					var privateKey interface{}
					if isEncrypted {
						if pass, ok := d.GetOk(isInstancePassphrase); ok {
							passphrase = pass.(string)
						} else {
							return fmt.Errorf("[ERROR] Mandatory field 'passphrase' not provided")
						}
						var err error
						privateKey, err = sshkeys.ParseEncryptedRawPrivateKey(keybytes, []byte(passphrase))
						if err != nil {
							return fmt.Errorf("[ERROR] Fail to decrypting the private key: %s", err)
						}
					} else {
						var err error
						privateKey, err = sshkeys.ParseEncryptedRawPrivateKey(keybytes, nil)
						if err != nil {
							return fmt.Errorf("[ERROR] Fail to decrypting the private key: %s", err)
						}
					}
					var ok bool
					rsaKey, ok = privateKey.(*rsa.PrivateKey)
					if !ok {
						return fmt.Errorf("[ERROR] Failed to convert to RSA private key")
					}
				}
			}
		}
	}

	getInstanceInitializationOptions := &vpcv1.GetInstanceInitializationOptions{
		ID: &id,
	}
	initParms, response, err := sess.GetInstanceInitialization(getInstanceInitializationOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Getting instance Initialization: %s\n%s", err, response)
	}
	if initParms.Keys != nil {
		initKeyList := make([]map[string]interface{}, 0)
		for _, key := range initParms.Keys {
			initKey := map[string]interface{}{}
			id := ""
			if key.ID != nil {
				id = *key.ID
			}
			initKey["id"] = id
			name := ""
			if key.Name != nil {
				name = *key.Name
			}
			initKey["name"] = name
			initKeyList = append(initKeyList, initKey)
			break

		}
		d.Set(isInstanceInitKeys, initKeyList)
	}
	//set the lifecycle status, reasons
	if instance.LifecycleState != nil {
		d.Set(isInstanceLifecycleState, *instance.LifecycleState)
	}
	if instance.LifecycleReasons != nil {
		d.Set(isInstanceLifecycleReasons, dataSourceInstanceFlattenLifecycleReasons(instance.LifecycleReasons))
	}

	if initParms.Password != nil && initParms.Password.EncryptedPassword != nil {
		ciphertext := *initParms.Password.EncryptedPassword
		password := base64.StdEncoding.EncodeToString(ciphertext)
		if rsaKey != nil {
			rng := rand.Reader
			clearPassword, err := rsa.DecryptPKCS1v15(rng, rsaKey, ciphertext)
			if err != nil {
				return fmt.Errorf("[ERROR] Can not decrypt the password with the given key, %s", err)
			}
			password = string(clearPassword)
		}
		d.Set(isInstanceInitPassword, password)
	}

	d.Set(isInstanceStatus, *instance.Status)
	//set the status reasons
	if instance.StatusReasons != nil {
		statusReasonsList := make([]map[string]interface{}, 0)
		for _, sr := range instance.StatusReasons {
			currentSR := map[string]interface{}{}
			if sr.Code != nil && sr.Message != nil {
				currentSR[isInstanceStatusReasonsCode] = *sr.Code
				currentSR[isInstanceStatusReasonsMessage] = *sr.Message
				if sr.MoreInfo != nil {
					currentSR[isInstanceStatusReasonsMoreInfo] = *sr.MoreInfo
				}
				statusReasonsList = append(statusReasonsList, currentSR)
			}
		}
		d.Set(isInstanceStatusReasons, statusReasonsList)
	}
	d.Set(isInstanceVPC, *instance.VPC.ID)
	d.Set(isInstanceZone, *instance.Zone.Name)

	var volumes []string
	volumes = make([]string, 0)
	if instance.VolumeAttachments != nil {
		for _, volume := range instance.VolumeAttachments {
			if volume.Volume != nil && *volume.Volume.ID != *instance.BootVolumeAttachment.Volume.ID {
				volumes = append(volumes, *volume.Volume.ID)
			}
		}
	}
	d.Set(isInstanceVolumes, flex.NewStringSet(schema.HashString, volumes))
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
		d.Set(isInstanceVolumeAttachments, volList)
	}
	if instance.BootVolumeAttachment != nil {
		bootVolList := make([]map[string]interface{}, 0)
		bootVol := map[string]interface{}{}
		bootVol["id"] = *instance.BootVolumeAttachment.ID
		bootVol["name"] = *instance.BootVolumeAttachment.Name
		if instance.BootVolumeAttachment.Device != nil {
			bootVol["device"] = *instance.BootVolumeAttachment.Device.ID
		}
		if instance.BootVolumeAttachment.Volume != nil {
			bootVol["volume_name"] = *instance.BootVolumeAttachment.Volume.Name
			bootVol["volume_id"] = *instance.BootVolumeAttachment.Volume.ID
			bootVol["volume_crn"] = *instance.BootVolumeAttachment.Volume.CRN
		}
		bootVolList = append(bootVolList, bootVol)
		d.Set(isInstanceBootVolume, bootVolList)
	}
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *instance.CRN, "", isInstanceUserTagType)
	if err != nil {
		log.Printf(
			"[ERROR] Error on get of resource vpc Instance (%s) tags: %s", d.Id(), err)
	}
	d.Set(isInstanceTags, tags)
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *instance.CRN, "", isInstanceAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc Instance (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isInstanceAccessTags, accesstags)

	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, controller+"/vpc-ext/compute/vs")
	d.Set(flex.ResourceName, instance.Name)
	d.Set(flex.ResourceCRN, instance.CRN)
	d.Set(IsInstanceCRN, instance.CRN)
	d.Set(flex.ResourceStatus, instance.Status)
	if instance.ResourceGroup != nil {
		d.Set(isInstanceResourceGroup, instance.ResourceGroup.ID)
		d.Set(flex.ResourceGroupName, instance.ResourceGroup.Name)
	}
	return nil

}

const opensshv1Magic = "openssh-key-v1"

type opensshPrivateKey struct {
	CipherName   string
	KdfName      string
	KdfOpts      string
	NumKeys      uint32
	PubKey       string
	PrivKeyBlock string
}

func isOpenSSHPrivKeyEncrypted(data []byte) (bool, error) {
	magic := append([]byte(opensshv1Magic), 0)
	if !bytes.Equal(magic, data[0:len(magic)]) {
		return false, errors.New("[ERROR] Invalid openssh private key format")
	}
	content := data[len(magic):]

	privKey := opensshPrivateKey{}

	if err := ssh.Unmarshal(content, &privKey); err != nil {
		return false, err
	}

	if privKey.KdfName == "none" && privKey.CipherName == "none" {
		return false, nil
	}
	return true, nil
}

func dataSourceInstanceFlattenDisks(result []vpcv1.InstanceDisk) (disks []map[string]interface{}) {
	for _, disksItem := range result {
		disks = append(disks, dataSourceInstanceDisksToMap(disksItem))
	}

	return disks
}

func dataSourceInstanceDisksToMap(disksItem vpcv1.InstanceDisk) (disksMap map[string]interface{}) {
	disksMap = map[string]interface{}{}

	if disksItem.CreatedAt != nil {
		disksMap["created_at"] = disksItem.CreatedAt.String()
	}
	if disksItem.Href != nil {
		disksMap["href"] = disksItem.Href
	}
	if disksItem.ID != nil {
		disksMap["id"] = disksItem.ID
	}
	if disksItem.InterfaceType != nil {
		disksMap["interface_type"] = disksItem.InterfaceType
	}
	if disksItem.Name != nil {
		disksMap["name"] = disksItem.Name
	}
	if disksItem.ResourceType != nil {
		disksMap["resource_type"] = disksItem.ResourceType
	}
	if disksItem.Size != nil {
		disksMap["size"] = disksItem.Size
	}

	return disksMap
}
func dataSourceInstanceFlattenLifecycleReasons(lifecycleReasons []vpcv1.InstanceLifecycleReason) (lifecycleReasonsList []map[string]interface{}) {
	lifecycleReasonsList = make([]map[string]interface{}, 0)
	for _, lr := range lifecycleReasons {
		currentLR := map[string]interface{}{}
		if lr.Code != nil && lr.Message != nil {
			currentLR[isInstanceLifecycleReasonsCode] = *lr.Code
			currentLR[isInstanceLifecycleReasonsMessage] = *lr.Message
			if lr.MoreInfo != nil {
				currentLR[isInstanceLifecycleReasonsMoreInfo] = *lr.MoreInfo
			}
			lifecycleReasonsList = append(lifecycleReasonsList, currentLR)
		}
	}
	return lifecycleReasonsList
}
