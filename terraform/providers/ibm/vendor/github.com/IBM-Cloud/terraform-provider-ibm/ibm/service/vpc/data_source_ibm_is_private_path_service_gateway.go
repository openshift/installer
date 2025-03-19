// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsPrivatePathServiceGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsPrivatePathServiceGatewayRead,

		Schema: map[string]*schema.Schema{
			"private_path_service_gateway": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"private_path_service_gateway_name", "private_path_service_gateway"},
				Description:  "The private path service gateway identifier.",
			},
			"private_path_service_gateway_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"private_path_service_gateway_name", "private_path_service_gateway"},
				Description:  "The private path service gateway name.",
			},
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
	}
}

func dataSourceIBMIsPrivatePathServiceGatewayRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	var privatePathServiceGateway *vpcv1.PrivatePathServiceGateway
	if ppsgId, ok := d.GetOk("private_path_service_gateway"); ok {

		getPrivatePathServiceGatewayOptions := &vpcv1.GetPrivatePathServiceGatewayOptions{}

		getPrivatePathServiceGatewayOptions.SetID(ppsgId.(string))

		ppsg, response, err := vpcClient.GetPrivatePathServiceGatewayWithContext(context, getPrivatePathServiceGatewayOptions)
		if err != nil {
			log.Printf("[DEBUG] GetPrivatePathServiceGatewayWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("GetPrivatePathServiceGatewayWithContext failed %s\n%s", err, response))
		}
		privatePathServiceGateway = ppsg
	} else {
		ppsgName := d.Get("private_path_service_gateway_name").(string)

		listPrivatePathServiceGatewaysOptions := &vpcv1.ListPrivatePathServiceGatewaysOptions{}

		privatePathServiceGatewayCollection, response, err := vpcClient.ListPrivatePathServiceGatewaysWithContext(context, listPrivatePathServiceGatewaysOptions)
		if err != nil {
			log.Printf("[DEBUG] ListPrivatePathServiceGatewaysWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListPrivatePathServiceGatewaysWithContext failed %s\n%s", err, response))
		}
		if privatePathServiceGatewayCollection.PrivatePathServiceGateways != nil {
			for _, ppsgItem := range privatePathServiceGatewayCollection.PrivatePathServiceGateways {
				if *ppsgItem.Name == ppsgName {
					privatePathServiceGateway = &ppsgItem
				}
			}
		}
	}
	d.SetId(*privatePathServiceGateway.ID)

	if err = d.Set("created_at", flex.DateTimeToString(privatePathServiceGateway.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("crn", privatePathServiceGateway.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if err = d.Set("default_access_policy", privatePathServiceGateway.DefaultAccessPolicy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting default_access_policy: %s", err))
	}

	if err = d.Set("endpoint_gateway_count", flex.IntValue(privatePathServiceGateway.EndpointGatewayCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting endpoint_gateway_count: %s", err))
	}
	if err = d.Set("endpoint_gateway_binding_auto_delete", privatePathServiceGateway.EndpointGatewayBindingAutoDelete); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting endpoint_gateway_binding_auto_delete: %s", err))
	}
	if err = d.Set("endpoint_gateway_binding_auto_delete_timeout", privatePathServiceGateway.EndpointGatewayBindingAutoDeleteTimeout); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting endpoint_gateway_binding_auto_delete_timeout: %s", err))
	}
	if err = d.Set("href", privatePathServiceGateway.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("lifecycle_state", privatePathServiceGateway.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}

	loadBalancer := []map[string]interface{}{}
	if privatePathServiceGateway.LoadBalancer != nil {
		modelMap, err := dataSourceIBMIsPrivatePathServiceGatewayLoadBalancerReferenceToMap(privatePathServiceGateway.LoadBalancer)
		if err != nil {
			return diag.FromErr(err)
		}
		loadBalancer = append(loadBalancer, modelMap)
	}
	if err = d.Set("load_balancer", loadBalancer); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting load_balancer %s", err))
	}

	if err = d.Set("name", privatePathServiceGateway.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("published", privatePathServiceGateway.Published); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting published: %s", err))
	}

	region := []map[string]interface{}{}
	// if privatePathServiceGateway.Remote != nil && privatePathServiceGateway.Remote.Region != nil {
	// 	modelMap, err := dataSourceIBMIsPrivatePathServiceGatewayRegionReferenceToMap(privatePathServiceGateway.Remote.Region)
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	}
	// 	region = append(region, modelMap)
	// }
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region %s", err))
	}

	resourceGroup := []map[string]interface{}{}
	if privatePathServiceGateway.ResourceGroup != nil {
		modelMap, err := dataSourceIBMIsPrivatePathServiceGatewayResourceGroupReferenceToMap(privatePathServiceGateway.ResourceGroup)
		if err != nil {
			return diag.FromErr(err)
		}
		resourceGroup = append(resourceGroup, modelMap)
	}
	if err = d.Set("resource_group", resourceGroup); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group %s", err))
	}

	if err = d.Set("resource_type", privatePathServiceGateway.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if err = d.Set("service_endpoints", privatePathServiceGateway.ServiceEndpoints); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting service_endpoints: %s", err))
	}
	vpc := []map[string]interface{}{}
	if privatePathServiceGateway.VPC != nil {
		modelMap, err := dataSourceIBMIsPrivatePathServiceGatewayVPCReferenceToMap(privatePathServiceGateway.VPC)
		if err != nil {
			return diag.FromErr(err)
		}
		vpc = append(vpc, modelMap)
	}
	if err = d.Set("vpc", vpc); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vpc %s", err))
	}

	if err = d.Set("zonal_affinity", privatePathServiceGateway.ZonalAffinity); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting zonal_affinity: %s", err))
	}

	return nil
}

func dataSourceIBMIsPrivatePathServiceGatewayLoadBalancerReferenceToMap(model *vpcv1.LoadBalancerReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsPrivatePathServiceGatewayLoadBalancerReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsPrivatePathServiceGatewayLoadBalancerReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewayRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewayResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
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

func dataSourceIBMIsPrivatePathServiceGatewayVPCReferenceToMap(model *vpcv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsPrivatePathServiceGatewayVPCReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsPrivatePathServiceGatewayVPCReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
