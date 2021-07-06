// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"errors"
	"fmt"
	"log"
	"time"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (

	/* Fix for PowerVC taking time to attach volume depending on load*/

	attachVolumeTimeOut = 240 * time.Second
)

func resourceIBMPIVolumeAttach() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMPIVolumeAttachCreate,
		Read:   resourceIBMPIVolumeAttachRead,
		Update: resourceIBMPIVolumeAttachUpdate,
		Delete: resourceIBMPIVolumeAttachDelete,
		//Exists:   resourceIBMPowerVolumeExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"volumeattachid": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Optional:    true,
				Description: "Volume attachment ID",
			},

			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: " Cloud Instance ID - This is the service_instance_id.",
			},

			helpers.PIVolumeAttachName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the volume to attach. Note these  volumes should have been created",
			},

			helpers.PIInstanceName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI Instance name",
			},

			helpers.PIVolumeAttachStatus: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			helpers.PIVolumeShareable: {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceIBMPIVolumeAttachCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	name := d.Get(helpers.PIVolumeAttachName).(string)
	servername := d.Get(helpers.PIInstanceName).(string)
	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)

	client := st.NewIBMPIVolumeClient(sess, powerinstanceid)

	volinfo, err := client.Get(name, powerinstanceid, getTimeOut)

	if err != nil {
		return fmt.Errorf("The volume [ %s] cannot be attached since it's not available", name)
	}
	//log.Print("The volume info is %s", volinfo)

	if volinfo.State == "available" || *volinfo.Shareable == true {
		log.Printf(" In the current state the volume can be attached to the instance ")
	}

	if volinfo.State == "in-use" && *volinfo.Shareable == true {

		log.Printf("Volume State /Status is  permitted and hence attaching the volume to the instance")
	}

	if volinfo.State == helpers.PIVolumeAllowableAttachStatus && *volinfo.Shareable == false {

		return errors.New("The volume cannot be attached in the current state. The volume must be in the *available* state. No other states are permissible")
	}

	resp, err := client.Attach(servername, name, powerinstanceid, attachVolumeTimeOut)

	if err != nil {
		return err
	}
	log.Printf("Printing the resp %+v", resp)

	d.SetId(*volinfo.VolumeID)
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return err
	}

	_, err = isWaitForIBMPIVolumeAttachAvailable(client, d.Id(), powerinstanceid, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	//return nil
	return resourceIBMPIVolumeAttachRead(d, meta)
}

func resourceIBMPIVolumeAttachRead(d *schema.ResourceData, meta interface{}) error {
	sess, _ := meta.(ClientSession).IBMPISession()
	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	servername := d.Get(helpers.PIInstanceName).(string)

	client := st.NewIBMPIVolumeClient(sess, powerinstanceid)

	vol, err := client.CheckVolumeAttach(powerinstanceid, servername, d.Id(), getTimeOut)
	if err != nil {
		return err
	}

	//d.SetId(vol.ID.String())
	d.Set(helpers.PIVolumeAttachName, vol.Name)
	d.Set(helpers.PIVolumeSize, vol.Size)
	d.Set(helpers.PIVolumeShareable, vol.Shareable)
	return nil
}

func resourceIBMPIVolumeAttachUpdate(d *schema.ResourceData, meta interface{}) error {

	sess, _ := meta.(ClientSession).IBMPISession()
	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	client := st.NewIBMPIVolumeClient(sess, powerinstanceid)

	name := ""
	if d.HasChange(helpers.PIVolumeAttachName) {
		name = d.Get(helpers.PIVolumeAttachName).(string)
	}

	size := float64(d.Get(helpers.PIVolumeSize).(float64))
	shareable := bool(d.Get(helpers.PIVolumeShareable).(bool))

	volrequest, err := client.Update(d.Id(), name, size, shareable, powerinstanceid, postTimeOut)
	if err != nil {
		return err
	}

	_, err = isWaitForIBMPIVolumeAttachAvailable(client, *volrequest.VolumeID, powerinstanceid, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return resourceIBMPIVolumeRead(d, meta)
}

func resourceIBMPIVolumeAttachDelete(d *schema.ResourceData, meta interface{}) error {

	sess, _ := meta.(ClientSession).IBMPISession()
	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	name := d.Get(helpers.PIVolumeAttachName).(string)
	servername := d.Get(helpers.PIInstanceName).(string)
	client := st.NewIBMPIVolumeClient(sess, powerinstanceid)

	log.Printf("the id of the volume to detach is%s ", d.Id())
	_, err := client.Detach(servername, name, powerinstanceid, deleteTimeOut)
	if err != nil {
		return err
	}

	// wait for power volume states to be back as available. if it's attached it will be in-use
	d.SetId("")
	return nil
}

func isWaitForIBMPIVolumeAttachAvailable(client *st.IBMPIVolumeClient, id, powerinstanceid string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume (%s) to be available for attachment", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PIVolumeProvisioning},
		Target:     []string{helpers.PIVolumeAllowableAttachStatus},
		Refresh:    isIBMPIVolumeAttachRefreshFunc(client, id, powerinstanceid),
		Delay:      10 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    10 * time.Minute,
	}

	return stateConf.WaitForState()
}

func isIBMPIVolumeAttachRefreshFunc(client *st.IBMPIVolumeClient, id, powerinstanceid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vol, err := client.Get(id, powerinstanceid, getTimeOut)
		if err != nil {
			return nil, "", err
		}

		if vol.State == "in-use" {
			return vol, helpers.PIVolumeAllowableAttachStatus, nil
		}

		return vol, helpers.PIVolumeProvisioning, nil
	}
}
