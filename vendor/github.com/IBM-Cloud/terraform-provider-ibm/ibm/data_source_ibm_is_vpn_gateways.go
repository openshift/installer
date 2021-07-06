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
	isvpnGateways            = "vpn_gateways"
	isVPNGatewayResourceType = "resource_type"
	isVPNGatewayCrn          = "crn"
)

func dataSourceIBMISVPNGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMVPNGatewaysRead,

		Schema: map[string]*schema.Schema{

			isvpnGateways: {
				Type:        schema.TypeList,
				Description: "Collection of VPN Gateways",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVPNGatewayName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPN Gateway instance name",
						},
						isVPNGatewayCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this VPN gateway was created",
						},
						isVPNGatewayCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPN gateway's CRN",
						},
						isVPNGatewayMembers: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of VPN gateway members",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The public IP address assigned to the VPN gateway member",
									},

									"private_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The private IP address assigned to the VPN gateway member",
									},

									"role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The high availability role assigned to the VPN gateway member",
									},

									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status of the VPN gateway member",
									},
								},
							},
						},

						isVPNGatewayResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},

						isVPNGatewayStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the VPN gateway",
						},

						isVPNGatewaySubnet: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPNGateway subnet info",
						},
						isVPNGatewayResourceGroup: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "resource group identifiers ",
						},
						isVPNGatewayMode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: " VPN gateway mode(policy/route) ",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMVPNGatewaysRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	listvpnGWOptions := sess.NewListVPNGatewaysOptions()

	start := ""
	allrecs := []vpcv1.VPNGatewayIntf{}
	for {
		if start != "" {
			listvpnGWOptions.Start = &start
		}
		availableVPNGateways, detail, err := sess.ListVPNGateways(listvpnGWOptions)
		if err != nil {
			return fmt.Errorf("Error reading list of VPN Gateways:%s\n%s", err, detail)
		}
		start = GetNext(availableVPNGateways.Next)
		allrecs = append(allrecs, availableVPNGateways.VPNGateways...)
		if start == "" {
			break
		}
	}

	vpngateways := make([]map[string]interface{}, 0)
	for _, instance := range allrecs {
		gateway := map[string]interface{}{}
		data := instance.(*vpcv1.VPNGateway)
		gateway[isVPNGatewayName] = *data.Name
		gateway[isVPNGatewayCreatedAt] = data.CreatedAt.String()
		gateway[isVPNGatewayResourceType] = *data.ResourceType
		gateway[isVPNGatewayStatus] = *data.Status
		gateway[isVPNGatewayMode] = *data.Mode
		gateway[isVPNGatewayResourceGroup] = *data.ResourceGroup.ID
		gateway[isVPNGatewaySubnet] = *data.Subnet.ID
		gateway[isVPNGatewayCrn] = *data.CRN

		if data.Members != nil {
			vpcMembersIpsList := make([]map[string]interface{}, 0)
			for _, memberIP := range data.Members {
				currentMemberIP := map[string]interface{}{}
				if memberIP.PublicIP != nil {
					currentMemberIP["address"] = *memberIP.PublicIP.Address
					currentMemberIP["role"] = *memberIP.Role
					currentMemberIP["status"] = *memberIP.Status
					vpcMembersIpsList = append(vpcMembersIpsList, currentMemberIP)
				}
				if memberIP.PrivateIP != nil {
					currentMemberIP["private_address"] = *memberIP.PrivateIP.Address
				}
			}
			gateway[isVPNGatewayMembers] = vpcMembersIpsList
		}

		vpngateways = append(vpngateways, gateway)
	}

	d.SetId(dataSourceIBMVPNGatewaysID(d))
	d.Set(isvpnGateways, vpngateways)
	return nil
}

// dataSourceIBMVPNGatewaysID returns a reasonable ID  list.
func dataSourceIBMVPNGatewaysID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
