// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
)

func resourceIBMPIImage() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMPIImageCreate,
		Read:     resourceIBMPIImageRead,
		Update:   resourceIBMPIImageUpdate,
		Delete:   resourceIBMPIImageDelete,
		Exists:   resourceIBMPIImageExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			helpers.PIImageName: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Image name",
				DiffSuppressFunc: applyOnce,
			},

			helpers.PIInstanceImageName: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Instance image name",
				DiffSuppressFunc: applyOnce,
			},

			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI cloud instance ID",
			},

			// Computed Attribute

			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image ID",
			},
		},
	}
}

func resourceIBMPIImageCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		log.Printf("Failed to get the session")
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	name := d.Get(helpers.PIImageName).(string)
	imageid := d.Get(helpers.PIInstanceImageName).(string)

	client := st.NewIBMPIImageClient(sess, powerinstanceid)

	imageResponse, err := client.Create(name, imageid, powerinstanceid)
	if err != nil {
		return err
	}

	IBMPIImageID := imageResponse.ImageID
	d.SetId(fmt.Sprintf("%s/%s", powerinstanceid, *IBMPIImageID))

	_, err = isWaitForIBMPIImageAvailable(client, *IBMPIImageID, d.Timeout(schema.TimeoutCreate), powerinstanceid)
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return err
	}

	return resourceIBMPIImageRead(d, meta)
}

func resourceIBMPIImageRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	powerinstanceid := parts[0]
	imageC := st.NewIBMPIImageClient(sess, powerinstanceid)
	imagedata, err := imageC.Get(parts[1], powerinstanceid)

	if err != nil {
		return err
	}

	imageid := *imagedata.ImageID
	d.Set("image_id", imageid)
	d.Set(helpers.PICloudInstanceId, powerinstanceid)

	return nil

}

func resourceIBMPIImageUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIBMPIImageDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	powerinstanceid := parts[0]
	imageC := st.NewIBMPIImageClient(sess, powerinstanceid)
	err = imageC.Delete(parts[1], powerinstanceid)

	if err != nil {
		return err
	}
	d.SetId("")
	return nil

}

func resourceIBMPIImageExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return false, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	name := parts[1]
	powerinstanceid := parts[0]
	client := st.NewIBMPIImageClient(sess, powerinstanceid)

	image, err := client.Get(parts[1], powerinstanceid)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}
	return *image.ImageID == name, nil
}

func isWaitForIBMPIImageAvailable(client *st.IBMPIImageClient, id string, timeout time.Duration, powerinstanceid string) (interface{}, error) {
	log.Printf("Waiting for Power Image (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PIImageQueStatus},
		Target:     []string{helpers.PIImageActiveStatus},
		Refresh:    isIBMPIImageRefreshFunc(client, id, powerinstanceid),
		Timeout:    timeout,
		Delay:      20 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isIBMPIImageRefreshFunc(client *st.IBMPIImageClient, id, powerinstanceid string) resource.StateRefreshFunc {

	log.Printf("Calling the isIBMPIImageRefreshFunc Refresh Function....")
	return func() (interface{}, string, error) {
		image, err := client.Get(id, powerinstanceid)
		if err != nil {
			return nil, "", err
		}

		if image.State == "active" {

			return image, helpers.PIImageActiveStatus, nil
		}

		return image, helpers.PIImageQueStatus, nil
	}
}
