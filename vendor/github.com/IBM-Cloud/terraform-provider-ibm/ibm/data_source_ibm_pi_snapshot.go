// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceIBMPISnapshot() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMPISnapshotRead,
		Schema: map[string]*schema.Schema{

			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			helpers.PIInstanceName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			//Computed Attributes

			"pvm_snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"percent_complete": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_updated_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_snapshots": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPISnapshotRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	powerinstancename := d.Get(helpers.PIInstanceName).(string)
	snapshot := instance.NewIBMPIInstanceClient(sess, powerinstanceid)
	snapshotData, err := snapshot.GetSnapShotVM(powerinstanceid, powerinstancename, getTimeOut)

	if err != nil {
		return err
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set("pvm_snapshots", flattenPVMSnapshotInstances(snapshotData.Snapshots))

	return nil

}

func flattenPVMSnapshotInstances(list []*models.Snapshot) []map[string]interface{} {
	log.Printf("Calling the flattensnapshotinstances call with list %d", len(list))
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"id":                *i.SnapshotID,
			"name":              *i.Name,
			"description":       i.Description,
			"creation_date":     i.CreationDate.String(),
			"last_updated_date": i.LastUpdateDate.String(),
			"action":            i.Action,
			"percent_complete":  i.PercentComplete,
			"status":            i.Status,
			"volume_snapshots":  i.VolumeSnapshots,
		}

		result = append(result, l)
	}

	return result
}
