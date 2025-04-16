// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSnapshotClone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISSnapshotCloneRead,

		Schema: map[string]*schema.Schema{
			isSnapshot: {
				Type:     schema.TypeString,
				Required: true,
			},

			isSnapshotCloneAvailable: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this snapshot clone is available for use.",
			},

			isSnapshotCloneCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this snapshot clone was created.",
			},

			isSnapshotCloneZone: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The zone this snapshot clone resides in.",
			},
		},
	}
}

func dataSourceIBMISSnapshotCloneRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get(isSnapshot).(string)
	zone := d.Get(isSnapshotCloneZone).(string)
	err := getSnapshotClone(context, d, meta, id, zone)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getSnapshotClone(context context.Context, d *schema.ResourceData, meta interface{}, id, zone string) error {
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}

	getSnapshotCloneOptions := &vpcv1.GetSnapshotCloneOptions{
		ID:       &id,
		ZoneName: &zone,
	}

	clone, response, err := sess.GetSnapshotCloneWithContext(context, getSnapshotCloneOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error fetching snapshot(%s) clone(%s) %s\n%s", id, zone, err, response)
	}

	if clone != nil && clone.Zone != nil {
		d.SetId(*clone.Zone.Name)
		d.Set(isSnapshotCloneZone, *clone.Zone.Name)
		d.Set(isSnapshotCloneAvailable, *clone.Available)
		if clone.CreatedAt != nil {
			d.Set(isSnapshotCloneCreatedAt, flex.DateTimeToString(clone.CreatedAt))
		}
	} else {
		return fmt.Errorf("[ERROR] No snapshot(%s) clone(%s) found", id, zone)
	}
	return nil
}
