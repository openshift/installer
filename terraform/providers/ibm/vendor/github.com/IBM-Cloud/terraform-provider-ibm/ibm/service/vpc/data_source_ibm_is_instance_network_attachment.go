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

func DataSourceIBMIsInstanceNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsInstanceNetworkAttachmentRead,

		Schema: map[string]*schema.Schema{
			"instance": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The virtual server instance identifier.",
			},
			"network_attachment": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance network attachment identifier.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the instance network attachment was created.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this instance network attachment.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the instance network attachment.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this instance network attachment. The name is unique across all network attachments for the instance.",
			},
			"port_speed": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port speed for this instance network attachment in Mbps.",
			},
			"primary_ip": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The primary IP address of the virtual network interface for the instance networkattachment.",
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
				Description: "The subnet of the virtual network interface for the instance network attachment.",
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
				Description: "The instance network attachment type.",
			},
			"virtual_network_interface": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The virtual network interface for this instance network attachment.",
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
		},
	}
}

func dataSourceIBMIsInstanceNetworkAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getInstanceNetworkAttachmentOptions := &vpcv1.GetInstanceNetworkAttachmentOptions{}

	getInstanceNetworkAttachmentOptions.SetInstanceID(d.Get("instance").(string))
	getInstanceNetworkAttachmentOptions.SetID(d.Get("network_attachment").(string))

	instanceByNetworkAttachment, response, err := vpcClient.GetInstanceNetworkAttachmentWithContext(context, getInstanceNetworkAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] GetInstanceNetworkAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetInstanceNetworkAttachmentWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getInstanceNetworkAttachmentOptions.InstanceID, *getInstanceNetworkAttachmentOptions.ID))

	if err = d.Set("created_at", flex.DateTimeToString(instanceByNetworkAttachment.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}

	if err = d.Set("href", instanceByNetworkAttachment.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}

	if err = d.Set("lifecycle_state", instanceByNetworkAttachment.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
	}

	if err = d.Set("name", instanceByNetworkAttachment.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if err = d.Set("port_speed", flex.IntValue(instanceByNetworkAttachment.PortSpeed)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting port_speed: %s", err))
	}

	primaryIP := []map[string]interface{}{}
	if instanceByNetworkAttachment.PrimaryIP != nil {
		modelMap, err := dataSourceIBMIsInstanceNetworkAttachmentReservedIPReferenceToMap(instanceByNetworkAttachment.PrimaryIP)
		if err != nil {
			return diag.FromErr(err)
		}
		primaryIP = append(primaryIP, modelMap)
	}
	if err = d.Set("primary_ip", primaryIP); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting primary_ip %s", err))
	}

	if err = d.Set("resource_type", instanceByNetworkAttachment.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}

	subnet := []map[string]interface{}{}
	if instanceByNetworkAttachment.Subnet != nil {
		modelMap, err := dataSourceIBMIsInstanceNetworkAttachmentSubnetReferenceToMap(instanceByNetworkAttachment.Subnet)
		if err != nil {
			return diag.FromErr(err)
		}
		subnet = append(subnet, modelMap)
	}
	if err = d.Set("subnet", subnet); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting subnet %s", err))
	}

	if err = d.Set("type", instanceByNetworkAttachment.Type); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
	}

	virtualNetworkInterface := []map[string]interface{}{}
	if instanceByNetworkAttachment.VirtualNetworkInterface != nil {
		modelMap, err := dataSourceIBMIsInstanceNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(instanceByNetworkAttachment.VirtualNetworkInterface)
		if err != nil {
			return diag.FromErr(err)
		}
		virtualNetworkInterface = append(virtualNetworkInterface, modelMap)
	}
	if err = d.Set("virtual_network_interface", virtualNetworkInterface); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting virtual_network_interface %s", err))
	}

	return nil
}

func dataSourceIBMIsInstanceNetworkAttachmentReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsInstanceNetworkAttachmentReservedIPReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsInstanceNetworkAttachmentReservedIPReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsInstanceNetworkAttachmentSubnetReferenceToMap(model *vpcv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsInstanceNetworkAttachmentSubnetReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsInstanceNetworkAttachmentSubnetReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsInstanceNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(model *vpcv1.VirtualNetworkInterfaceReferenceAttachmentContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}
