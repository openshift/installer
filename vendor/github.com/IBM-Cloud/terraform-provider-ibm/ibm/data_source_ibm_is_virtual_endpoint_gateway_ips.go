// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMISEndpointGatewayIPs() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceIBMISEndpointGatewayIPsRead,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			isVirtualEndpointGatewayID: {
				Type:     schema.TypeString,
				Required: true,
			},
			isVirtualEndpointGatewayIPs: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVirtualEndpointGatewayIPID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway IP id",
						},
						isVirtualEndpointGatewayIPName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway IP name",
						},
						isVirtualEndpointGatewayIPResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway IP resource type",
						},
						isVirtualEndpointGatewayIPCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway IP created date and time",
						},
						isVirtualEndpointGatewayIPAutoDelete: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Endpoint gateway IP auto delete",
						},
						isVirtualEndpointGatewayIPAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint gateway IP address",
						},
						isVirtualEndpointGatewayIPTarget: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Endpoint gateway detail",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVirtualEndpointGatewayIPTargetID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IPs target id",
									},
									isVirtualEndpointGatewayIPTargetName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IPs target name",
									},
									isVirtualEndpointGatewayIPTargetResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Endpoint gateway resource type",
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

func dataSourceIBMISEndpointGatewayIPsRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	gatewayID := d.Get(isVirtualEndpointGatewayID).(string)

	start := ""
	allrecs := []vpcv1.ReservedIP{}
	for {
		options := sess.NewListEndpointGatewayIpsOptions(gatewayID)
		if start != "" {
			options.Start = &start
		}
		result, response, err := sess.ListEndpointGatewayIps(options)
		if err != nil {
			return fmt.Errorf("Error fetching endpoint gateway ips %s\n%s", err, response)
		}
		start = GetNext(result.Next)
		allrecs = append(allrecs, result.Ips...)
		if start == "" {
			break
		}
	}
	endpointGatewayIPs := []map[string]interface{}{}
	for _, ip := range allrecs {
		ipsOutput := map[string]interface{}{}
		ipsOutput[isVirtualEndpointGatewayIPID] = *ip.ID
		ipsOutput[isVirtualEndpointGatewayIPName] = *ip.Name
		ipsOutput[isVirtualEndpointGatewayIPCreatedAt] = (*ip.CreatedAt).String()
		ipsOutput[isVirtualEndpointGatewayIPAddress] = *ip.Address
		ipsOutput[isVirtualEndpointGatewayIPAutoDelete] = *ip.AutoDelete
		ipsOutput[isVirtualEndpointGatewayIPResourceType] = *ip.ResourceType
		ipsOutput[isVirtualEndpointGatewayIPTarget] =
			flattenEndpointGatewayIPTarget(ip.Target.(*vpcv1.ReservedIPTarget))

		endpointGatewayIPs = append(endpointGatewayIPs, ipsOutput)
	}
	d.SetId(dataSourceIBMISEndpointGatewayIPsCheckID(d))
	d.Set(isVirtualEndpointGatewayIPs, endpointGatewayIPs)
	return nil
}

// dataSourceIBMISEndpointGatewayIPsCheckID returns a reasonable ID for dns zones list.
func dataSourceIBMISEndpointGatewayIPsCheckID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
