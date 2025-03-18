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

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsPrivatePathServiceGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsPrivatePathServiceGatewaysRead,

		Schema: map[string]*schema.Schema{
			"private_path_service_gateways": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of private path service gateways.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the private path service gateway was created.",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this private path service gateway.",
						},
						"default_access_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy to use for bindings from accounts without an explicit account policy.",
						},
						"endpoint_gateway_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of endpoint gateways using this private path service gateway.",
						},
						"endpoint_gateway_binding_auto_delete": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether endpoint gateway bindings will be automatically deleted after endpoint_gateway_binding_auto_delete_timeout hours have passed.",
						},
						"endpoint_gateway_binding_auto_delete_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "If endpoint_gateway_binding_auto_delete is true, the hours after which endpoint gateway bindings will be automatically deleted.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this private path service gateway.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this private path service gateway.",
						},
						"lifecycle_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the private path service gateway.",
						},
						"load_balancer": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The load balancer for this private path service gateway.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The load balancer's CRN.",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
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
										Description: "The load balancer's canonical URL.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this load balancer.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this load balancer. The name is unique across all load balancers in the VPC.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this private path service gateway. The name is unique across all private path service gateways in the VPC.",
						},
						"published": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the availability of this private path service gateway- `true`: Any account can request access to this private path service gateway.- `false`: Access is restricted to the account that created this private path service gateway.",
						},
						"region": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The region served by this private path service gateway.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this region.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this region.",
									},
								},
							},
						},
						"resource_group": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this private path service gateway.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this resource group.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this resource group.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this resource group.",
									},
								},
							},
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"service_endpoints": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The fully qualified domain names for this private path service gateway.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"vpc": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The VPC this private path service gateway resides in.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this VPC.",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
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
										Description: "The URL for this VPC.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this VPC.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this VPC. The name is unique across all VPCs in the region.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"zonal_affinity": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this private path service gateway has zonal affinity.- `true`:  Traffic to the service from a zone will favor service endpoints in           the same zone.- `false`: Traffic to the service from a zone will be load balanced across all zones           in the region the service resides in.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsPrivatePathServiceGatewaysRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listPrivatePathServiceGatewaysOptions := &vpcv1.ListPrivatePathServiceGatewaysOptions{}

	privatePathServiceGatewayCollection, response, err := vpcClient.ListPrivatePathServiceGatewaysWithContext(context, listPrivatePathServiceGatewaysOptions)
	if err != nil {
		log.Printf("[DEBUG] ListPrivatePathServiceGatewaysWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListPrivatePathServiceGatewaysWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMIsPrivatePathServiceGatewaysID(d))

	privatePathServiceGateways := []map[string]interface{}{}
	if privatePathServiceGatewayCollection.PrivatePathServiceGateways != nil {
		for _, modelItem := range privatePathServiceGatewayCollection.PrivatePathServiceGateways {
			modelMap, err := dataSourceIBMIsPrivatePathServiceGatewaysPrivatePathServiceGatewayToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			privatePathServiceGateways = append(privatePathServiceGateways, modelMap)
		}
	}
	if err = d.Set("private_path_service_gateways", privatePathServiceGateways); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting private_path_service_gateways %s", err))
	}

	return nil
}

// dataSourceIBMIsPrivatePathServiceGatewaysID returns a reasonable ID for the list.
func dataSourceIBMIsPrivatePathServiceGatewaysID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMIsPrivatePathServiceGatewaysPrivatePathServiceGatewayCollectionFirstToMap(model *vpcv1.PrivatePathServiceGatewayCollectionFirst) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewaysPrivatePathServiceGatewayCollectionNextToMap(model *vpcv1.PrivatePathServiceGatewayCollectionNext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewaysPrivatePathServiceGatewayToMap(model *vpcv1.PrivatePathServiceGateway) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CreatedAt != nil {
		modelMap["created_at"] = flex.DateTimeToString(model.CreatedAt)
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.DefaultAccessPolicy != nil {
		modelMap["default_access_policy"] = *model.DefaultAccessPolicy
	}
	if model.EndpointGatewayCount != nil {
		modelMap["endpoint_gateway_count"] = *model.EndpointGatewayCount
	}
	if model.EndpointGatewayBindingAutoDelete != nil {
		modelMap["endpoint_gateway_binding_auto_delete"] = *model.EndpointGatewayBindingAutoDelete
	}
	if model.EndpointGatewayBindingAutoDeleteTimeout != nil {
		modelMap["endpoint_gateway_binding_auto_delete_timeout"] = *model.EndpointGatewayBindingAutoDeleteTimeout
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.LifecycleState != nil {
		modelMap["lifecycle_state"] = *model.LifecycleState
	}
	if model.LoadBalancer != nil {
		loadBalancerMap, err := dataSourceIBMIsPrivatePathServiceGatewaysLoadBalancerReferenceToMap(model.LoadBalancer)
		if err != nil {
			return modelMap, err
		}
		modelMap["load_balancer"] = []map[string]interface{}{loadBalancerMap}
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Published != nil {
		modelMap["published"] = *model.Published
	}
	// if model.Remote != nil && model.Remote.Region != nil {
	// 	regionMap, err := dataSourceIBMIsPrivatePathServiceGatewaysRegionReferenceToMap(model.Remote.Region)
	// 	if err != nil {
	// 		return modelMap, err
	// 	}
	// 	modelMap["region"] = []map[string]interface{}{regionMap}
	// }
	if model.ResourceGroup != nil {
		resourceGroupMap, err := dataSourceIBMIsPrivatePathServiceGatewaysResourceGroupReferenceToMap(model.ResourceGroup)
		if err != nil {
			return modelMap, err
		}
		modelMap["resource_group"] = []map[string]interface{}{resourceGroupMap}
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	if model.ServiceEndpoints != nil {
		modelMap["service_endpoints"] = model.ServiceEndpoints
	}
	if model.VPC != nil {
		vpcMap, err := dataSourceIBMIsPrivatePathServiceGatewaysVPCReferenceToMap(model.VPC)
		if err != nil {
			return modelMap, err
		}
		modelMap["vpc"] = []map[string]interface{}{vpcMap}
	}
	if model.ZonalAffinity != nil {
		modelMap["zonal_affinity"] = *model.ZonalAffinity
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewaysLoadBalancerReferenceToMap(model *vpcv1.LoadBalancerReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsPrivatePathServiceGatewaysLoadBalancerReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewaysLoadBalancerReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewaysRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewaysResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewaysVPCReferenceToMap(model *vpcv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsPrivatePathServiceGatewaysVPCReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewaysVPCReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
