// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func ResourceIBMPISnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPISnapshotCreate,
		ReadContext:   resourceIBMPISnapshotRead,
		UpdateContext: resourceIBMPISnapshotUpdate,
		DeleteContext: resourceIBMPISnapshotDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			helpers.PISnapshotName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique name of the snapshot",
			},
			helpers.PIInstanceName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name / id of the pvm",
			},
			helpers.PIInstanceVolumeIds: {
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              schema.HashString,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "List of PI volumes",
			},
			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: " Cloud Instance ID - This is the service_instance_id.",
			},
			"pi_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the PVM instance snapshot",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Snapshot description",
				Deprecated:  "This field is deprecated, use pi_description instead",
			},

			// Computed Attributes
			helpers.PISnapshot: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the snapshot",
				Deprecated:  "This field is deprecated, use snapshot_id instead",
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the PVM instance snapshot",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_snapshots": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceIBMPISnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	instanceid := d.Get(helpers.PIInstanceName).(string)
	volids := flex.ExpandStringList((d.Get(helpers.PIInstanceVolumeIds).(*schema.Set)).List())
	name := d.Get(helpers.PISnapshotName).(string)

	var description string
	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}
	if v, ok := d.GetOk("pi_description"); ok {
		description = v.(string)
	}

	client := st.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)

	snapshotBody := &models.SnapshotCreate{Name: &name, Description: description}

	if len(volids) > 0 {
		snapshotBody.VolumeIDs = volids
	} else {
		log.Printf("no volumeids provided. Will snapshot the entire instance")
	}

	snapshotResponse, err := client.CreatePvmSnapShot(instanceid, snapshotBody)
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *snapshotResponse.SnapshotID))

	pisnapclient := st.NewIBMPISnapshotClient(ctx, sess, cloudInstanceID)
	_, err = isWaitForPIInstanceSnapshotAvailable(ctx, pisnapclient, *snapshotResponse.SnapshotID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMPISnapshotRead(ctx, d, meta)
}

func resourceIBMPISnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("Calling the Snapshot Read function post create")
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, snapshotID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	snapshot := st.NewIBMPISnapshotClient(ctx, sess, cloudInstanceID)
	snapshotdata, err := snapshot.Get(snapshotID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(helpers.PISnapshotName, snapshotdata.Name)
	d.Set(helpers.PISnapshot, *snapshotdata.SnapshotID)
	d.Set("snapshot_id", *snapshotdata.SnapshotID)
	d.Set("status", snapshotdata.Status)
	d.Set("creation_date", snapshotdata.CreationDate.String())
	d.Set("volume_snapshots", snapshotdata.VolumeSnapshots)
	d.Set("last_update_date", snapshotdata.LastUpdateDate.String())

	return nil
}

func resourceIBMPISnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log.Printf("Calling the IBM Power Snapshot  update call")
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, snapshotID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := st.NewIBMPISnapshotClient(ctx, sess, cloudInstanceID)

	if d.HasChange(helpers.PISnapshotName) || d.HasChange("description") {
		name := d.Get(helpers.PISnapshotName).(string)
		description := d.Get("description").(string)
		snapshotBody := &models.SnapshotUpdate{Name: name, Description: description}

		_, err := client.Update(snapshotID, snapshotBody)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = isWaitForPIInstanceSnapshotAvailable(ctx, client, snapshotID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMPISnapshotRead(ctx, d, meta)
}

func resourceIBMPISnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, snapshotID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := st.NewIBMPISnapshotClient(ctx, sess, cloudInstanceID)
	snapshot, err := client.Get(snapshotID)
	if err != nil {
		// snapshot does not exist
		d.SetId("")
		return nil
	}

	log.Printf("The snapshot  to be deleted is in the following state .. %s", snapshot.Status)

	err = client.Delete(snapshotID)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = isWaitForPIInstanceSnapshotDeleted(ctx, client, snapshotID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
func isWaitForPIInstanceSnapshotAvailable(ctx context.Context, client *st.IBMPISnapshotClient, id string, timeout time.Duration) (interface{}, error) {

	log.Printf("Waiting for PIInstance Snapshot (%s) to be available and active ", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"in_progress", "BUILD"},
		Target:     []string{"available", "ACTIVE"},
		Refresh:    isPIInstanceSnapshotRefreshFunc(client, id),
		Delay:      30 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceSnapshotRefreshFunc(client *st.IBMPISnapshotClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		snapshotInfo, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		//if pvm.Health.Status == helpers.PIInstanceHealthOk {
		if snapshotInfo.Status == "available" && snapshotInfo.PercentComplete == 100 {
			log.Printf("The snapshot is now available")
			return snapshotInfo, "available", nil

		}
		return snapshotInfo, "in_progress", nil
	}
}

// Delete Snapshot

func isWaitForPIInstanceSnapshotDeleted(ctx context.Context, client *st.IBMPISnapshotClient, id string, timeout time.Duration) (interface{}, error) {

	log.Printf("Waiting for (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PIInstanceDeleting},
		Target:     []string{"Not Found"},
		Refresh:    isPIInstanceSnapshotDeleteRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPIInstanceSnapshotDeleteRefreshFunc(client *st.IBMPISnapshotClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		snapshot, err := client.Get(id)
		if err != nil {
			log.Printf("The snapshot is not found.")
			return snapshot, helpers.PIInstanceNotFound, nil
		}
		return snapshot, helpers.PIInstanceNotFound, nil

	}
}
