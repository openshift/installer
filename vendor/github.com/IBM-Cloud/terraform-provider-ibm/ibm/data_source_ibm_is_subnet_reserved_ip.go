// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func dataSourceIBMISReservedIP() *schema.Resource {
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
		return fmt.Errorf("Error fetching the reserved IP %s\n%s", err, response)
	}

	d.SetId(*reserveIP.ID)
	d.Set(isReservedIPAutoDelete, *reserveIP.AutoDelete)
	d.Set(isReservedIPCreatedAt, (*reserveIP.CreatedAt).String())
	d.Set(isReservedIPhref, *reserveIP.Href)
	d.Set(isReservedIPName, *reserveIP.Name)
	d.Set(isReservedIPOwner, *reserveIP.Owner)
	d.Set(isReservedIPType, *reserveIP.ResourceType)
	if reserveIP.Target != nil {
		target, ok := reserveIP.Target.(*vpcv1.ReservedIPTarget)
		if ok {
			d.Set(isReservedIPTarget, target.ID)
		}
	}
	return nil // By default there should be no error
}
