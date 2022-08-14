// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"reflect"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Define all the constants that matches with the given terrafrom attribute
const (
	// Request Param Constants
	isReservedIPLimit  = "limit"
	isReservedIPSort   = "sort"
	isReservedIPs      = "reserved_ips"
	isReservedIPsCount = "total_count"
)

func DataSourceIBMISReservedIPs() *schema.Resource {
	return &schema.Resource{
		Read: dataSdataSourceIBMISReservedIPsRead,
		Schema: map[string]*schema.Schema{
			/*
				Request Parameters
				==================
				These are mandatory req parameters
			*/
			isSubNetID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The subnet identifier.",
			},
			/*
				Response Parameters
				===================
				All of these are computed and an user doesn't need to provide
				these from outside.
			*/

			isReservedIPs: {
				Type:        schema.TypeList,
				Description: "Collection of reserved IPs in this subnet.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isReservedIPAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address",
						},
						isReservedIPAutoDelete: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If reserved ip shall be deleted automatically",
						},
						isReservedIPCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the reserved IP was created.",
						},
						isReservedIPhref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this reserved IP.",
						},
						isReservedIPID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this reserved IP",
						},
						isReservedIPLifecycleState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the reserved IP",
						},
						isReservedIPName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined or system-provided name for this reserved IP.",
						},
						isReservedIPOwner: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The owner of a reserved IP, defining whether it is managed by the user or the provider.",
						},
						isReservedIPType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						isReservedIPTarget: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reserved IP target id",
						},
					},
				},
			},
			isReservedIPsCount: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages",
			},
		},
	}
}

func dataSdataSourceIBMISReservedIPsRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	subnetID := d.Get(isSubNetID).(string)

	// Flatten all the reserved IPs
	start := ""
	allrecs := []vpcv1.ReservedIP{}
	for {
		options := &vpcv1.ListSubnetReservedIpsOptions{SubnetID: &subnetID}

		if start != "" {
			options.Start = &start
		}

		result, response, err := sess.ListSubnetReservedIps(options)
		if err != nil || response == nil || result == nil {
			return fmt.Errorf("[ERROR] Error fetching reserved ips %s\n%s", err, response)
		}
		start = flex.GetNext(result.Next)
		allrecs = append(allrecs, result.ReservedIps...)
		if start == "" {
			break
		}
	}

	// Now store all the reserved IP info with their response tags
	reservedIPs := []map[string]interface{}{}
	for _, data := range allrecs {
		ipsOutput := map[string]interface{}{}
		ipsOutput[isReservedIPAddress] = *data.Address
		ipsOutput[isReservedIPAutoDelete] = *data.AutoDelete
		ipsOutput[isReservedIPCreatedAt] = (*data.CreatedAt).String()
		ipsOutput[isReservedIPhref] = *data.Href
		ipsOutput[isReservedIPID] = *data.ID
		ipsOutput[isReservedIPLifecycleState] = data.LifecycleState
		ipsOutput[isReservedIPName] = *data.Name
		ipsOutput[isReservedIPOwner] = *data.Owner
		ipsOutput[isReservedIPType] = *data.ResourceType
		if data.Target != nil {
			targetIntf := data.Target
			switch reflect.TypeOf(targetIntf).String() {
			case "*vpcv1.ReservedIPTargetEndpointGatewayReference":
				{
					target := targetIntf.(*vpcv1.ReservedIPTargetEndpointGatewayReference)
					ipsOutput[isReservedIPTarget] = target.ID
				}
			case "*vpcv1.ReservedIPTargetNetworkInterfaceReferenceTargetContext":
				{
					target := targetIntf.(*vpcv1.ReservedIPTargetNetworkInterfaceReferenceTargetContext)
					ipsOutput[isReservedIPTarget] = target.ID
				}
			case "*vpcv1.ReservedIPTargetLoadBalancerReference":
				{
					target := targetIntf.(*vpcv1.ReservedIPTargetLoadBalancerReference)
					ipsOutput[isReservedIPTarget] = target.ID
				}
			case "*vpcv1.ReservedIPTargetVPNGatewayReference":
				{
					target := targetIntf.(*vpcv1.ReservedIPTargetVPNGatewayReference)
					ipsOutput[isReservedIPTarget] = target.ID
				}
			case "*vpcv1.ReservedIPTarget":
				{
					target := targetIntf.(*vpcv1.ReservedIPTarget)
					ipsOutput[isReservedIPTarget] = target.ID
				}
			}
		}
		reservedIPs = append(reservedIPs, ipsOutput)
	}

	d.SetId(time.Now().UTC().String()) // This is not any reserved ip or subnet id but state id
	d.Set(isReservedIPs, reservedIPs)
	d.Set(isReservedIPsCount, len(reservedIPs))
	d.Set(isSubNetID, subnetID)
	return nil
}
