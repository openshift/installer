// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsInstanceClusterNetworkAttachments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsInstanceClusterNetworkAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The virtual server instance identifier.",
			},
			"cluster_network_attachments": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A page of ordered cluster network attachments (sorted based on the `before` property) for the instance. A cluster network attachment represents a device to which a cluster network interface is attached.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"before": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The instance cluster network attachment that is immediately before. If absent, this is thelast instance cluster network attachment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this instance cluster network attachment.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this instance cluster network attachment.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this instance cluster network attachment. The name is unique across all network attachments for the instance.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"cluster_network_interface": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The cluster network interface for this instance cluster network attachment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
										Description: "The URL for this cluster network interface.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this cluster network interface.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this cluster network interface. The name is unique across all interfaces in the cluster network.",
									},
									"primary_ip": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The primary IP for this cluster network interface.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP address.If the address is pending allocation, the value will be `0.0.0.0`.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 addresses in the future.",
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
													Description: "The URL for this cluster network subnet reserved IP.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this cluster network subnet reserved IP.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this cluster network subnet reserved IP. The name is unique across all reserved IPs in a cluster network subnet.",
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
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
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
													Description: "The URL for this cluster network subnet.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this cluster network subnet.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this cluster network subnet. The name is unique across all cluster network subnets in the cluster network.",
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
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this instance cluster network attachment.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this instance cluster network attachment.",
						},
						"lifecycle_reasons": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current `lifecycle_state` (if any).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A reason code for this lifecycle state:- `internal_error`: internal error (contact IBM support)- `resource_suspended_by_provider`: The resource has been suspended (contact IBM  support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
									},
									"message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the reason for this lifecycle state.",
									},
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about the reason for this lifecycle state.",
									},
								},
							},
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the instance cluster network attachment.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this instance cluster network attachment. The name is unique across all network attachments for the instance.",
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

func dataSourceIBMIsInstanceClusterNetworkAttachmentsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_cluster_network_attachments", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listInstanceClusterNetworkAttachmentsOptions := &vpcv1.ListInstanceClusterNetworkAttachmentsOptions{}

	listInstanceClusterNetworkAttachmentsOptions.SetInstanceID(d.Get("instance_id").(string))

	var pager *vpcv1.InstanceClusterNetworkAttachmentsPager
	pager, err = vpcClient.NewInstanceClusterNetworkAttachmentsPager(listInstanceClusterNetworkAttachmentsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_cluster_network_attachments", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("InstanceClusterNetworkAttachmentsPager.GetAll() failed %s", err), "(Data) ibm_is_instance_cluster_network_attachments", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsInstanceClusterNetworkAttachmentsID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentToMap(&modelItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_cluster_network_attachments", "read", "Instances-to-map").GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("cluster_network_attachments", mapSlice); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting cluster_network_attachments %s", err), "(Data) ibm_is_instance_cluster_network_attachments", "read", "cluster_network_attachments-set").GetDiag()
	}

	return nil
}

// dataSourceIBMIsInstanceClusterNetworkAttachmentsID returns a reasonable ID for the list.
func dataSourceIBMIsInstanceClusterNetworkAttachmentsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentToMap(model *vpcv1.InstanceClusterNetworkAttachment) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Before != nil {
		beforeMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentBeforeToMap(model.Before)
		if err != nil {
			return modelMap, err
		}
		modelMap["before"] = []map[string]interface{}{beforeMap}
	}
	clusterNetworkInterfaceMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkInterfaceReferenceToMap(model.ClusterNetworkInterface)
	if err != nil {
		return modelMap, err
	}
	modelMap["cluster_network_interface"] = []map[string]interface{}{clusterNetworkInterfaceMap}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range model.LifecycleReasons {
		lifecycleReasonsItemMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentLifecycleReasonToMap(&lifecycleReasonsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	modelMap["lifecycle_reasons"] = lifecycleReasons
	modelMap["lifecycle_state"] = *model.LifecycleState
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentBeforeToMap(model *vpcv1.InstanceClusterNetworkAttachmentBefore) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkInterfaceReferenceToMap(model *vpcv1.ClusterNetworkInterfaceReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentsDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	primaryIPMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReservedIPReferenceToMap(model.PrimaryIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
	modelMap["resource_type"] = *model.ResourceType
	subnetMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReferenceToMap(model.Subnet)
	if err != nil {
		return modelMap, err
	}
	modelMap["subnet"] = []map[string]interface{}{subnetMap}
	return modelMap, nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentsDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReservedIPReferenceToMap(model *vpcv1.ClusterNetworkSubnetReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = *model.Address
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentsDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReferenceToMap(model *vpcv1.ClusterNetworkSubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsInstanceClusterNetworkAttachmentsDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentLifecycleReasonToMap(model *vpcv1.InstanceClusterNetworkAttachmentLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
