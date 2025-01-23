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

func DataSourceIBMIsInstanceNetworkAttachments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsInstanceNetworkAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"instance": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The virtual server instance identifier.",
			},
			"network_attachments": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of instance network attachments.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this instance network attachment.",
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
				},
			},
		},
	}
}

func dataSourceIBMIsInstanceNetworkAttachmentsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	listInstanceNetworkAttachmentsOptions := &vpcv1.ListInstanceNetworkAttachmentsOptions{}

	listInstanceNetworkAttachmentsOptions.SetInstanceID(d.Get("instance").(string))

	instanceNetworkAttachmentCollection, response, err := vpcClient.ListInstanceNetworkAttachmentsWithContext(context, listInstanceNetworkAttachmentsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListInstanceNetworkAttachmentsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListInstanceNetworkAttachmentsWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMIsInstanceNetworkAttachmentsID(d))

	networkAttachments := []map[string]interface{}{}
	if instanceNetworkAttachmentCollection.NetworkAttachments != nil {
		for _, modelItem := range instanceNetworkAttachmentCollection.NetworkAttachments {
			modelMap, err := dataSourceIBMIsInstanceNetworkAttachmentsInstanceNetworkAttachmentToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			networkAttachments = append(networkAttachments, modelMap)
		}
	}
	if err = d.Set("network_attachments", networkAttachments); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting network_attachments %s", err))
	}

	return nil
}

// dataSourceIBMIsInstanceNetworkAttachmentsID returns a reasonable ID for the list.
func dataSourceIBMIsInstanceNetworkAttachmentsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMIsInstanceNetworkAttachmentsInstanceNetworkAttachmentToMap(model *vpcv1.InstanceNetworkAttachment) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["lifecycle_state"] = model.LifecycleState
	modelMap["name"] = model.Name
	modelMap["port_speed"] = flex.IntValue(model.PortSpeed)
	primaryIPMap, err := dataSourceIBMIsInstanceNetworkAttachmentsReservedIPReferenceToMap(model.PrimaryIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
	modelMap["resource_type"] = model.ResourceType
	subnetMap, err := dataSourceIBMIsInstanceNetworkAttachmentsSubnetReferenceToMap(model.Subnet)
	if err != nil {
		return modelMap, err
	}
	modelMap["subnet"] = []map[string]interface{}{subnetMap}
	modelMap["type"] = model.Type
	virtualNetworkInterfaceMap, err := dataSourceIBMIsInstanceNetworkAttachmentsVirtualNetworkInterfaceReferenceAttachmentContextToMap(model.VirtualNetworkInterface)
	if err != nil {
		return modelMap, err
	}
	modelMap["virtual_network_interface"] = []map[string]interface{}{virtualNetworkInterfaceMap}
	return modelMap, nil
}

func dataSourceIBMIsInstanceNetworkAttachmentsReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsInstanceNetworkAttachmentsReservedIPReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsInstanceNetworkAttachmentsReservedIPReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsInstanceNetworkAttachmentsSubnetReferenceToMap(model *vpcv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsInstanceNetworkAttachmentsSubnetReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsInstanceNetworkAttachmentsSubnetReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsInstanceNetworkAttachmentsVirtualNetworkInterfaceReferenceAttachmentContextToMap(model *vpcv1.VirtualNetworkInterfaceReferenceAttachmentContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}
