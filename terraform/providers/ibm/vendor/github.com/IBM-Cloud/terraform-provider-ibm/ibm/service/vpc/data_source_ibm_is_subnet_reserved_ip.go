// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Define all the constants that matches with the given terrafrom attribute
const (
	// Request Param Constants
	isSubNetID     = "subnet"
	isReservedIPID = "reserved_ip"

	// Response Param Constants
	isReservedIPAddress    = "address"
	isReservedIPAutoDelete = "auto_delete"
	isReservedIPCreatedAt  = "created_at"
	isReservedIPhref       = "href"
	isReservedIPName       = "name"
	isReservedIPOwner      = "owner"
	isReservedIPType       = "resource_type"
)

func DataSourceIBMISReservedIP() *schema.Resource {
	return &schema.Resource{
		Read: dataSdataSourceIBMISReservedIPRead,
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
			isReservedIPID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The reserved IP identifier.",
			},

			/*
				Response Parameters
				===================
				All of these are computed and an user doesn't need to provide
				these from outside.
			*/

			isReservedIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address",
			},
			isReservedIPLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the reserved IP",
			},
			isReservedIPAutoDelete: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, this reserved IP will be automatically deleted",
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
				Description: "Reserved IP target id.",
			},
			isReservedIPTargetCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn for target.",
			},
			"target_reference": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The target this reserved IP is bound to.If absent, this reserved IP is provider-owned or unbound.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this endpoint gateway.",
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

// dataSdataSourceIBMISReservedIPRead is used when the reserved IPs are read from the vpc
func dataSdataSourceIBMISReservedIPRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	subnetID := d.Get(isSubNetID).(string)
	reservedIPID := d.Get(isReservedIPID).(string)

	options := sess.NewGetSubnetReservedIPOptions(subnetID, reservedIPID)
	reserveIP, response, err := sess.GetSubnetReservedIP(options)

	if err != nil || response == nil || reserveIP == nil {
		return fmt.Errorf("[ERROR] Error fetching the reserved IP %s\n%s", err, response)
	}

	d.SetId(*reserveIP.ID)
	d.Set(isReservedIPAutoDelete, *reserveIP.AutoDelete)
	d.Set(isReservedIPAddress, *reserveIP.Address)
	d.Set(isReservedIPCreatedAt, (*reserveIP.CreatedAt).String())
	d.Set(isReservedIPhref, *reserveIP.Href)
	d.Set(isReservedIPName, *reserveIP.Name)
	d.Set(isReservedIPOwner, *reserveIP.Owner)
	if reserveIP.LifecycleState != nil {
		d.Set(isReservedIPLifecycleState, *reserveIP.LifecycleState)
	}
	d.Set(isReservedIPType, *reserveIP.ResourceType)
	target := []map[string]interface{}{}
	if reserveIP.Target != nil {
		modelMap, err := dataSourceIBMIsReservedIPReservedIPTargetToMap(reserveIP.Target)
		if err != nil {
			return err
		}
		target = append(target, modelMap)
	}
	d.Set("target_reference", target)
	if len(target) > 0 {
		d.Set(isReservedIPTarget, target[0]["id"])
		d.Set(isReservedIPTargetCrn, target[0]["crn"])
	}
	return nil // By default there should be no error
}
