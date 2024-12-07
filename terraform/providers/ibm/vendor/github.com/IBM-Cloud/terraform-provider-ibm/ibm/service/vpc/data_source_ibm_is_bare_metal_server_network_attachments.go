// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsBareMetalServerNetworkAttachments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsBareMetalServerNetworkAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"bare_metal_server": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bare metal server identifier.",
			},
			"network_attachments": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of bare metal server network attachments.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the bare metal server network attachment was created.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this bare metal server network attachment.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this bare metal server network attachment.",
						},
						"interface_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network attachment's interface type:- `hipersocket`: a virtual network device that provides high-speed TCP/IP connectivity  within a `s390x` based system- `pci`: a physical PCI device which can only be created or deleted when the bare metal  server is stopped  - Has an `allowed_vlans` property which controls the VLANs that will be permitted    to use the PCI attachment  - Cannot directly use an IEEE 802.1q VLAN tag.- `vlan`: a virtual device, used through a `pci` device that has the `vlan` in its  array of `allowed_vlans`.  - Must use an IEEE 802.1q tag.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the bare metal server network attachment.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this bare metal server network attachment. The name is unique across all network attachments for the bare metal server.",
						},
						"port_speed": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port speed for this bare metal server network attachment in Mbps.",
						},
						"primary_ip": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The primary IP address of the virtual network interface for the bare metal servernetwork attachment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this reserved IP.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this reserved IP. The name is unique across all reserved IPs in a subnet.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"subnet": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The subnet of the virtual network interface for the bare metal server networkattachment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this subnet.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this subnet.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this subnet.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this subnet. The name is unique across all subnets in the VPC.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The bare metal server network attachment type.",
						},
						"virtual_network_interface": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The virtual network interface for this bare metal server network attachment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this virtual network interface.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this virtual network interface.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this virtual network interface.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this virtual network interface. The name is unique across all virtual network interfaces in the VPC.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"allowed_vlans": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"allow_to_float": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if the bare metal server network attachment can automatically float to any other server within the same `resource_group`. The bare metal server network attachment will float automatically if the network detects a GARP or RARP on another bare metal server in the resource group. Applies only to bare metal server network attachments with `vlan` interface type.",
						},
						"vlan": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the 802.1Q VLAN ID tag that must be used for all traffic on this attachment.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	listBareMetalServerNetworkAttachmentsOptions := &vpcv1.ListBareMetalServerNetworkAttachmentsOptions{}

	listBareMetalServerNetworkAttachmentsOptions.SetBareMetalServerID(d.Get("bare_metal_server").(string))

	var pager *vpcv1.BareMetalServerNetworkAttachmentsPager
	pager, err = vpcClient.NewBareMetalServerNetworkAttachmentsPager(listBareMetalServerNetworkAttachmentsOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] BareMetalServerNetworkAttachmentsPager.GetAll() failed %s", err)
		return diag.FromErr(fmt.Errorf("BareMetalServerNetworkAttachmentsPager.GetAll() failed %s", err))
	}

	d.SetId(dataSourceIBMIsBareMetalServerNetworkAttachmentsID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentsBareMetalServerNetworkAttachmentToMap(modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("network_attachments", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting network_attachments %s", err))
	}

	return nil
}

// dataSourceIBMIsBareMetalServerNetworkAttachmentsID returns a reasonable ID for the list.
func dataSourceIBMIsBareMetalServerNetworkAttachmentsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentsBareMetalServerNetworkAttachmentCollectionFirstToMap(model *vpcv1.BareMetalServerNetworkAttachmentCollectionFirst) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentsBareMetalServerNetworkAttachmentToMap(model vpcv1.BareMetalServerNetworkAttachmentIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.BareMetalServerNetworkAttachmentByPci); ok {
		return dataSourceIBMIsBareMetalServerNetworkAttachmentsBareMetalServerNetworkAttachmentByPciToMap(model.(*vpcv1.BareMetalServerNetworkAttachmentByPci))
	} else if _, ok := model.(*vpcv1.BareMetalServerNetworkAttachmentByVlan); ok {
		return dataSourceIBMIsBareMetalServerNetworkAttachmentsBareMetalServerNetworkAttachmentByVlanToMap(model.(*vpcv1.BareMetalServerNetworkAttachmentByVlan))
	} else if _, ok := model.(*vpcv1.BareMetalServerNetworkAttachment); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.BareMetalServerNetworkAttachment)
		modelMap["created_at"] = model.CreatedAt.String()
		modelMap["href"] = model.Href
		modelMap["id"] = model.ID
		modelMap["interface_type"] = model.InterfaceType
		modelMap["lifecycle_state"] = model.LifecycleState
		modelMap["name"] = model.Name
		modelMap["port_speed"] = flex.IntValue(model.PortSpeed)
		primaryIPMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentsReservedIPReferenceToMap(model.PrimaryIP)
		if err != nil {
			return modelMap, err
		}
		modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
		modelMap["resource_type"] = model.ResourceType
		subnetMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentsSubnetReferenceToMap(model.Subnet)
		if err != nil {
			return modelMap, err
		}
		modelMap["subnet"] = []map[string]interface{}{subnetMap}
		modelMap["type"] = model.Type
		virtualNetworkInterfaceMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentsVirtualNetworkInterfaceReferenceAttachmentContextToMap(model.VirtualNetworkInterface)
		if err != nil {
			return modelMap, err
		}
		modelMap["virtual_network_interface"] = []map[string]interface{}{virtualNetworkInterfaceMap}
		if model.AllowedVlans != nil {
			modelMap["allowed_vlans"] = model.AllowedVlans
		}
		if model.AllowToFloat != nil {
			modelMap["allow_to_float"] = model.AllowToFloat
		}
		if model.Vlan != nil {
			modelMap["vlan"] = flex.IntValue(model.Vlan)
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.BareMetalServerNetworkAttachmentIntf subtype encountered")
	}
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentsReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentsReservedIPReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentsReservedIPReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentsSubnetReferenceToMap(model *vpcv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentsSubnetReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentsSubnetReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentsVirtualNetworkInterfaceReferenceAttachmentContextToMap(model *vpcv1.VirtualNetworkInterfaceReferenceAttachmentContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentsBareMetalServerNetworkAttachmentByVlanToMap(model *vpcv1.BareMetalServerNetworkAttachmentByVlan) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["lifecycle_state"] = model.LifecycleState
	modelMap["name"] = model.Name
	modelMap["port_speed"] = flex.IntValue(model.PortSpeed)
	primaryIPMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentsReservedIPReferenceToMap(model.PrimaryIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
	modelMap["resource_type"] = model.ResourceType
	subnetMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentsSubnetReferenceToMap(model.Subnet)
	if err != nil {
		return modelMap, err
	}
	modelMap["subnet"] = []map[string]interface{}{subnetMap}
	modelMap["type"] = model.Type
	virtualNetworkInterfaceMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentsVirtualNetworkInterfaceReferenceAttachmentContextToMap(model.VirtualNetworkInterface)
	if err != nil {
		return modelMap, err
	}
	modelMap["virtual_network_interface"] = []map[string]interface{}{virtualNetworkInterfaceMap}
	modelMap["allow_to_float"] = model.AllowToFloat
	modelMap["interface_type"] = model.InterfaceType
	modelMap["vlan"] = flex.IntValue(model.Vlan)
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentsBareMetalServerNetworkAttachmentByPciToMap(model *vpcv1.BareMetalServerNetworkAttachmentByPci) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["lifecycle_state"] = model.LifecycleState
	modelMap["name"] = model.Name
	modelMap["port_speed"] = flex.IntValue(model.PortSpeed)
	primaryIPMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentsReservedIPReferenceToMap(model.PrimaryIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
	modelMap["resource_type"] = model.ResourceType
	subnetMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentsSubnetReferenceToMap(model.Subnet)
	if err != nil {
		return modelMap, err
	}
	modelMap["subnet"] = []map[string]interface{}{subnetMap}
	modelMap["type"] = model.Type
	virtualNetworkInterfaceMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentsVirtualNetworkInterfaceReferenceAttachmentContextToMap(model.VirtualNetworkInterface)
	if err != nil {
		return modelMap, err
	}
	modelMap["virtual_network_interface"] = []map[string]interface{}{virtualNetworkInterfaceMap}
	modelMap["allowed_vlans"] = model.AllowedVlans
	modelMap["interface_type"] = model.InterfaceType
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentsBareMetalServerNetworkAttachmentCollectionNextToMap(model *vpcv1.BareMetalServerNetworkAttachmentCollectionNext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}
