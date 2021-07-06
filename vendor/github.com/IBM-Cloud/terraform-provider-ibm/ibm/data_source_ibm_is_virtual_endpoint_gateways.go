// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ibmISVirtualEndpointGateways = "ibm_is_virtual_endpoint_gateways"
	isVirtualEndpointGateways    = "virtual_endpoint_gateways"
)

func dataSourceIBMISEndpointGateways() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceIBMISEndpointGatewaysRead,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			isVirtualEndpointGateways: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway id",
						},
						isVirtualEndpointGatewayName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway name",
						},
						isVirtualEndpointGatewayResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway resource type",
						},
						isVirtualEndpointGatewayResourceGroupID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource group id",
						},
						isVirtualEndpointGatewayCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway created date and time",
						},
						isVirtualEndpointGatewayHealthState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway health state",
						},
						isVirtualEndpointGatewayLifecycleState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway lifecycle state",
						},
						isVirtualEndpointGatewayIPs: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of reserved IPs bound to an endpoint gateway",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVirtualEndpointGatewayIPsID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this reserved IP",
									},
									isVirtualEndpointGatewayIPsName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined or system-provided name for this reserved IP",
									},
									isVirtualEndpointGatewayIPsResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type(subnet_reserved_ip)",
									},
								},
							},
						},
						isVirtualEndpointGatewayTarget: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Endpoint gateway target",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVirtualEndpointGatewayTargetName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The target name",
									},
									isVirtualEndpointGatewayTargetResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The target resource type",
									},
								},
							},
						},
						isVirtualEndpointGatewayVpcID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC id",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISEndpointGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	start := ""
	allrecs := []vpcv1.EndpointGateway{}
	for {
		options := sess.NewListEndpointGatewaysOptions()
		if start != "" {
			options.Start = &start
		}
		result, response, err := sess.ListEndpointGateways(options)
		if err != nil {
			return fmt.Errorf("Error fetching endpoint gateways %s\n%s", err, response)
		}
		start = GetNext(result.Next)
		allrecs = append(allrecs, result.EndpointGateways...)
		if start == "" {
			break
		}
	}
	endpointGateways := []map[string]interface{}{}
	for _, endpointGateway := range allrecs {
		endpointGatewayOutput := map[string]interface{}{}
		endpointGatewayOutput["id"] = *endpointGateway.ID
		endpointGatewayOutput[isVirtualEndpointGatewayName] = *endpointGateway.Name
		endpointGatewayOutput[isVirtualEndpointGatewayCreatedAt] = (*endpointGateway.CreatedAt).String()
		endpointGatewayOutput[isVirtualEndpointGatewayResourceType] = (*endpointGateway.ResourceType)
		endpointGatewayOutput[isVirtualEndpointGatewayHealthState] = *endpointGateway.HealthState
		endpointGatewayOutput[isVirtualEndpointGatewayLifecycleState] = *endpointGateway.LifecycleState
		endpointGatewayOutput[isVirtualEndpointGatewayResourceGroupID] = *endpointGateway.ResourceGroup.ID
		endpointGatewayOutput[isVirtualEndpointGatewayVpcID] = *endpointGateway.VPC.ID
		endpointGatewayOutput[isVirtualEndpointGatewayTarget] =
			flattenEndpointGatewayTarget(endpointGateway.Target.(*vpcv1.EndpointGatewayTarget))
		endpointGatewayOutput[isVirtualEndpointGatewayIPs] =
			flattenDataSourceIPs(endpointGateway.Ips)
		endpointGateways = append(endpointGateways, endpointGatewayOutput)
	}
	d.SetId(dataSourceIBMISEndpointGatewaysCheckID(d))
	d.Set(isVirtualEndpointGateways, endpointGateways)
	return nil
}

// dataSourceIBMISEndpointGatewaysCheckID returns a reasonable ID for dns zones list.
func dataSourceIBMISEndpointGatewaysCheckID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func flattenDataSourceIPs(ipsList []vpcv1.ReservedIPReference) interface{} {
	ipsListOutput := make([]interface{}, 0)
	for _, item := range ipsList {
		ips := make(map[string]interface{}, 0)
		ips[isVirtualEndpointGatewayIPsID] = *item.ID
		ips[isVirtualEndpointGatewayIPsName] = *item.Name
		ips[isVirtualEndpointGatewayIPsResourceType] = *item.ResourceType

		ipsListOutput = append(ipsListOutput, ips)
	}
	return ipsListOutput
}
