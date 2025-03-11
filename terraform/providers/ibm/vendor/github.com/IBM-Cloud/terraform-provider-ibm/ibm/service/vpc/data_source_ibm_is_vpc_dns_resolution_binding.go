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

func DataSourceIBMIsVPCDnsResolutionBinding() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPCDnsResolutionBindingRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPC identifier.",
			},
			"identifier": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The DNS resolution binding identifier.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the DNS resolution binding was created.",
			},
			"health_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current `health_state` (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the reason for this health state.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this health state.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this health state.",
						},
					},
				},
			},
			"health_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The health of this resource.- `ok`: No abnormal behavior detected- `degraded`: Experiencing compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.",
			},
			"endpoint_gateways": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The endpoint gateways in the bound to VPC that are allowed to participate in this DNS resolution binding.The endpoint gateways may be remote and therefore may not be directly retrievable.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this endpoint gateway.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this endpoint gateway.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this endpoint gateway.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this endpoint gateway. The name is unique across all endpoint gateways in the VPC.",
						},
						"remote": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the referenced resource is remote to thisaccount, and identifies the owning account.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this account.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"region": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the referenced resource is remote to thisregion, and identifies the native region.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this region.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The globally unique name for this region.",
												},
											},
										},
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this DNS resolution binding.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the DNS resolution binding.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this DNS resolution binding. The name is unique across all DNS resolution bindings for the VPC.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"vpc": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPC bound to for DNS resolution.The VPC may be remote and therefore may not be directly retrievable.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this VPC.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this VPC.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPC.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this VPC. The name is unique across all VPCs in the region.",
						},
						"remote": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the referenced resource is remote to thisaccount, and identifies the owning account.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this account.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"region": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the referenced resource is remote to thisregion, and identifies the native region.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this region.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The globally unique name for this region.",
												},
											},
										},
									},
								},
							},
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

func dataSourceIBMIsVPCDnsResolutionBindingRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getVPCDnsResolutionBindingOptions := &vpcv1.GetVPCDnsResolutionBindingOptions{}

	getVPCDnsResolutionBindingOptions.SetVPCID(d.Get("vpc_id").(string))
	getVPCDnsResolutionBindingOptions.SetID(d.Get("identifier").(string))

	vpcdnsResolutionBinding, response, err := sess.GetVPCDnsResolutionBindingWithContext(context, getVPCDnsResolutionBindingOptions)
	if err != nil {
		log.Printf("[DEBUG] GetVPCDnsResolutionBindingWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVPCDnsResolutionBindingWithContext failed %s\n%s", err, response))
	}

	d.SetId(*vpcdnsResolutionBinding.ID)

	if err = d.Set("created_at", flex.DateTimeToString(vpcdnsResolutionBinding.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}

	endpointGateways := []map[string]interface{}{}
	if vpcdnsResolutionBinding.EndpointGateways != nil {
		for _, modelItem := range vpcdnsResolutionBinding.EndpointGateways {
			modelMap, err := dataSourceIBMIsVPCDnsResolutionBindingEndpointGatewayReferenceRemoteToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			endpointGateways = append(endpointGateways, modelMap)
		}
	}
	if err = d.Set("endpoint_gateways", endpointGateways); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting endpoint_gateways %s", err))
	}

	if err = d.Set("href", vpcdnsResolutionBinding.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}

	if err = d.Set("lifecycle_state", vpcdnsResolutionBinding.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
	}

	if err = d.Set("name", vpcdnsResolutionBinding.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if err = d.Set("resource_type", vpcdnsResolutionBinding.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}
	healthReasons := []map[string]interface{}{}
	if vpcdnsResolutionBinding.HealthReasons != nil {
		for _, modelItem := range vpcdnsResolutionBinding.HealthReasons {
			modelMap, err := dataSourceIBMIsVPCDnsResolutionBindingVpcdnsResolutionBindingHealthReasonToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			healthReasons = append(healthReasons, modelMap)
		}
	}
	if err = d.Set("health_reasons", healthReasons); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting health_reasons %s", err))
	}

	if err = d.Set("health_state", vpcdnsResolutionBinding.HealthState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting health_state: %s", err))
	}

	vpc := []map[string]interface{}{}
	if vpcdnsResolutionBinding.VPC != nil {
		modelMap, err := dataSourceIBMIsVPCDnsResolutionBindingVPCReferenceRemoteToMap(vpcdnsResolutionBinding.VPC)
		if err != nil {
			return diag.FromErr(err)
		}
		vpc = append(vpc, modelMap)
	}
	if err = d.Set("vpc", vpc); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting vpc %s", err))
	}

	return nil
}

func dataSourceIBMIsVPCDnsResolutionBindingEndpointGatewayReferenceRemoteToMap(model *vpcv1.EndpointGatewayReferenceRemote) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	if model.Remote != nil {
		remoteMap, err := dataSourceIBMIsVPCDnsResolutionBindingEndpointGatewayRemoteToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsVPCDnsResolutionBindingEndpointGatewayRemoteToMap(model *vpcv1.EndpointGatewayRemote) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := dataSourceIBMIsVPCDnsResolutionBindingAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Region != nil {
		regionMap, err := dataSourceIBMIsVPCDnsResolutionBindingRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func dataSourceIBMIsVPCDnsResolutionBindingAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsVPCDnsResolutionBindingRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIBMIsVPCDnsResolutionBindingVPCReferenceRemoteToMap(model *vpcv1.VPCReferenceRemote) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	if model.Remote != nil {
		remoteMap, err := dataSourceIBMIsVPCDnsResolutionBindingVPCRemoteToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsVPCDnsResolutionBindingVPCRemoteToMap(model *vpcv1.VPCRemote) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := dataSourceIBMIsVPCDnsResolutionBindingAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Region != nil {
		regionMap, err := dataSourceIBMIsVPCDnsResolutionBindingRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func dataSourceIBMIsVPCDnsResolutionBindingVpcdnsResolutionBindingHealthReasonToMap(model *vpcv1.VpcdnsResolutionBindingHealthReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = model.Code
	modelMap["message"] = model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}
