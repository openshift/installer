// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsBareMetalServerNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsBareMetalServerNetworkAttachmentRead,

		Schema: map[string]*schema.Schema{
			"bare_metal_server": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bare metal server identifier.",
			},
			"network_attachment": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bare metal server network attachment identifier.",
			},
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
			"bare_metal_server_network_attachment_id": &schema.Schema{
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
	}
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getBareMetalServerNetworkAttachmentOptions := &vpcv1.GetBareMetalServerNetworkAttachmentOptions{}

	getBareMetalServerNetworkAttachmentOptions.SetBareMetalServerID(d.Get("bare_metal_server").(string))
	getBareMetalServerNetworkAttachmentOptions.SetID(d.Get("network_attachment").(string))

	bareMetalServerNetworkAttachmentIntf, response, err := vpcClient.GetBareMetalServerNetworkAttachmentWithContext(context, getBareMetalServerNetworkAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] GetBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response))
	}
	bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachment)

	d.SetId(fmt.Sprintf("%s/%s", *getBareMetalServerNetworkAttachmentOptions.BareMetalServerID, *getBareMetalServerNetworkAttachmentOptions.ID))

	if err = d.Set("created_at", flex.DateTimeToString(bareMetalServerNetworkAttachment.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("href", bareMetalServerNetworkAttachment.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("bare_metal_server_network_attachment_id", bareMetalServerNetworkAttachment.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting bare_metal_server_network_attachment_id: %s", err))
	}

	if err = d.Set("interface_type", bareMetalServerNetworkAttachment.InterfaceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting interface_type: %s", err))
	}

	if err = d.Set("lifecycle_state", bareMetalServerNetworkAttachment.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}

	if err = d.Set("name", bareMetalServerNetworkAttachment.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("port_speed", flex.IntValue(bareMetalServerNetworkAttachment.PortSpeed)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting port_speed: %s", err))
	}

	primaryIP := []map[string]interface{}{}
	if bareMetalServerNetworkAttachment.PrimaryIP != nil {
		modelMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentReservedIPReferenceToMap(bareMetalServerNetworkAttachment.PrimaryIP)
		if err != nil {
			return diag.FromErr(err)
		}
		primaryIP = append(primaryIP, modelMap)
	}
	if err = d.Set("primary_ip", primaryIP); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting primary_ip %s", err))
	}

	if err = d.Set("resource_type", bareMetalServerNetworkAttachment.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	subnet := []map[string]interface{}{}
	if bareMetalServerNetworkAttachment.Subnet != nil {
		modelMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentSubnetReferenceToMap(bareMetalServerNetworkAttachment.Subnet)
		if err != nil {
			return diag.FromErr(err)
		}
		subnet = append(subnet, modelMap)
	}
	if err = d.Set("subnet", subnet); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting subnet %s", err))
	}

	if err = d.Set("type", bareMetalServerNetworkAttachment.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}

	virtualNetworkInterface := []map[string]interface{}{}
	if bareMetalServerNetworkAttachment.VirtualNetworkInterface != nil {
		modelMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(bareMetalServerNetworkAttachment.VirtualNetworkInterface)
		if err != nil {
			return diag.FromErr(err)
		}
		virtualNetworkInterface = append(virtualNetworkInterface, modelMap)
	}
	if err = d.Set("virtual_network_interface", virtualNetworkInterface); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting virtual_network_interface %s", err))
	}

	if err = d.Set("allow_to_float", bareMetalServerNetworkAttachment.AllowToFloat); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting allow_to_float: %s", err))
	}

	if err = d.Set("vlan", flex.IntValue(bareMetalServerNetworkAttachment.Vlan)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vlan: %s", err))
	}

	return nil
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentReservedIPReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsBareMetalServerNetworkAttachmentReservedIPReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentSubnetReferenceToMap(model *vpcv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentSubnetReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsBareMetalServerNetworkAttachmentSubnetReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(model *vpcv1.VirtualNetworkInterfaceReferenceAttachmentContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}
