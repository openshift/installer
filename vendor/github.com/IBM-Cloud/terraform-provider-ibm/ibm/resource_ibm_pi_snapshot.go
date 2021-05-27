// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_p_vm_instances"
	"github.com/IBM-Cloud/power-go-client/power/models"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMPISnapshot() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMPISnapshotCreate,
		Read:     resourceIBMPISnapshotRead,
		Update:   resourceIBMPISnapshotUpdate,
		Delete:   resourceIBMPISnapshotDelete,
		Exists:   resourceIBMPISnapshotExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			//Snapshots are created at the pvm instance level

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
				DiffSuppressFunc: applyOnce,
				Description:      "List of PI volumes",
			},

			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: " Cloud Instance ID - This is the service_instance_id.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Snapshot description",
			},
			// Computed Attributes

			helpers.PISnapshot: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the snapshot",
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
	}
}

func resourceIBMPISnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	instanceid := d.Get(helpers.PIInstanceName).(string)
	volids := expandStringList((d.Get(helpers.PIInstanceVolumeIds).(*schema.Set)).List())
	name := d.Get(helpers.PISnapshotName).(string)
	description := d.Get("description").(string)
	if d.Get(description) == "" {
		description = "Testing from Terraform"
	}

	client := st.NewIBMPIInstanceClient(sess, powerinstanceid)

	snapshotBody := &models.SnapshotCreate{Name: &name, Description: description}

	if len(volids) > 0 {
		snapshotBody.VolumeIds = volids
	} else {
		log.Printf("no volumeids provided. Will snapshot the entire instance")
	}

	snapshotResponse, err := client.CreatePvmSnapShot(&p_cloud_p_vm_instances.PcloudPvminstancesSnapshotsPostParams{
		Body: snapshotBody,
	}, instanceid, powerinstanceid, createTimeOut)

	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", powerinstanceid, *snapshotResponse.SnapshotID))
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return fmt.Errorf("failed to get the snapshotid %v", err)
	}

	pisnapclient := st.NewIBMPISnapshotClient(sess, powerinstanceid)
	_, err = isWaitForPIInstanceSnapshotAvailable(pisnapclient, *snapshotResponse.SnapshotID, d.Timeout(schema.TimeoutCreate), powerinstanceid)
	if err != nil {
		return err
	}

	return resourceIBMPISnapshotRead(d, meta)
}

func resourceIBMPISnapshotRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Calling the Snapshot Read function post create")
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	powerinstanceid := parts[0]
	snapshot := st.NewIBMPISnapshotClient(sess, powerinstanceid)
	snapshotdata, err := snapshot.Get(parts[1], powerinstanceid, getTimeOut)

	if err != nil {
		return err
	}

	d.Set(helpers.PISnapshotName, snapshotdata.Name)
	d.Set(helpers.PISnapshot, *snapshotdata.SnapshotID)
	d.Set("status", snapshotdata.Status)
	d.Set("creation_date", snapshotdata.CreationDate.String())
	d.Set("volume_snapshots", snapshotdata.VolumeSnapshots)
	d.Set("last_update_date", snapshotdata.LastUpdateDate.String())

	return nil
}

func resourceIBMPISnapshotUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("Calling the IBM Power Snapshot  update call")
	sess, _ := meta.(ClientSession).IBMPISession()
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	powerinstanceid := parts[0]
	client := st.NewIBMPISnapshotClient(sess, powerinstanceid)

	if d.HasChange(helpers.PISnapshotName) || d.HasChange("description") {
		name := d.Get(helpers.PISnapshotName).(string)
		description := d.Get("description").(string)
		snapshotBody := &models.SnapshotUpdate{Name: name, Description: description}

		_, err := client.Update(parts[1], powerinstanceid, snapshotBody, 60)

		if err != nil {
			return fmt.Errorf("failed to update the snapshot request %v", err)

		}

		_, err = isWaitForPIInstanceSnapshotAvailable(client, parts[1], d.Timeout(schema.TimeoutCreate), powerinstanceid)
		if err != nil {
			return err
		}
	}

	return resourceIBMPISnapshotRead(d, meta)
}

func resourceIBMPISnapshotDelete(d *schema.ResourceData, meta interface{}) error {

	sess, _ := meta.(ClientSession).IBMPISession()
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	powerinstanceid := parts[0]

	client := st.NewIBMPISnapshotClient(sess, powerinstanceid)

	snapshot, err := client.Get(parts[1], powerinstanceid, getTimeOut)
	if err != nil {
		return err
	}

	log.Printf("The snapshot  to be deleted is in the following state .. %s", snapshot.Status)

	snapshotdel_err := client.Delete(parts[1], powerinstanceid, deleteTimeOut)
	if snapshotdel_err != nil {
		return snapshotdel_err
	}

	_, err = isWaitForPIInstanceSnapshotDeleted(client, parts[1], d.Timeout(schema.TimeoutDelete), powerinstanceid)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
func resourceIBMPISnapshotExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return false, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	powerinstanceid := parts[0]
	client := st.NewIBMPISnapshotClient(sess, powerinstanceid)

	snapshotdelete, err := client.Get(parts[1], powerinstanceid, getTimeOut)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	log.Printf("Calling the existing function.. %s", *(snapshotdelete.SnapshotID))

	volumeid := *snapshotdelete.SnapshotID
	return volumeid == parts[1], nil
}

func isWaitForPIInstanceSnapshotAvailable(client *st.IBMPISnapshotClient, id string, timeout time.Duration, powerinstanceid string) (interface{}, error) {

	log.Printf("Waiting for PIInstance Snapshot (%s) to be available and active ", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"in_progress", "BUILD"},
		Target:     []string{"available", "ACTIVE"},
		Refresh:    isPIInstanceSnapshotRefreshFunc(client, id, powerinstanceid),
		Delay:      30 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    60 * time.Minute,
	}

	return stateConf.WaitForState()
}

func isPIInstanceSnapshotRefreshFunc(client *st.IBMPISnapshotClient, id, powerinstanceid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		snapshotInfo, err := client.Get(id, powerinstanceid, getTimeOut)
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

func isWaitForPIInstanceSnapshotDeleted(client *st.IBMPISnapshotClient, id string, timeout time.Duration, powerinstanceid string) (interface{}, error) {

	log.Printf("Waiting for  (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PIInstanceDeleting},
		Target:     []string{"Not Found"},
		Refresh:    isPIInstanceSnapshotDeleteRefreshFunc(client, id, powerinstanceid),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Timeout:    10 * time.Minute,
	}

	return stateConf.WaitForState()
}

func isPIInstanceSnapshotDeleteRefreshFunc(client *st.IBMPISnapshotClient, id, powerinstanceid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		snapshot, err := client.Get(id, powerinstanceid, getTimeOut)
		if err != nil {
			log.Printf("The snapshot is not found.")
			return snapshot, helpers.PIInstanceNotFound, nil

		}
		return snapshot, helpers.PIInstanceNotFound, nil

	}
}
