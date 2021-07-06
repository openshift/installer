// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"errors"
	"log"
	"time"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_p_vm_instances"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMPICapture() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMPICaptureCreate,
		Read:   resourceIBMPICaptureRead,
		Update: resourceIBMPICaptureUpdate,
		Delete: resourceIBMPICaptureDelete,
		//Exists:   resourceIBMPICaptureExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: " Cloud Instance ID - This is the service_instance_id.",
			},

			helpers.PIInstanceName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance Name of the Power VM",
			},

			helpers.PIInstanceCaptureName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the capture to create. Note : this must be unique",
			},

			helpers.PIInstanceCaptureDestination: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of destination to store the image capture to",
				ValidateFunc: validateAllowedStringValue([]string{"image-catalog", "cloud-storage", "both"}),
			},

			helpers.PIInstanceCaptureVolumeIds: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "List of volume names that need to be passed in the input",
			},

			helpers.PIInstanceCaptureCloudStorageRegion: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "List of Regions to use",
				ValidateFunc: validateAllowedStringValue([]string{"us-south", "us-east", "us-de"}),
			},

			helpers.PIInstanceCaptureCloudStorageAccessKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of Cloud Storage Access Key",
			},
			helpers.PIInstanceCaptureCloudStorageSecretKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the Cloud Storage Secret Key",
			},
			helpers.PIInstanceCaptureCloudStorageImagePath: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the Image Path",
			},
		},
	}
}

func resourceIBMPICaptureCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	name := d.Get(helpers.PIInstanceName).(string)
	capturename := d.Get(helpers.PIInstanceCaptureName).(string)
	capturedestination := d.Get(helpers.PIInstanceCaptureDestination).(string)
	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)

	cloudstorageImagePath := d.Get(helpers.PIInstanceCaptureCloudStorageImagePath).(string)
	if cloudstorageImagePath == "" {
		log.Printf("CloudImagePath is not provided")

	}

	cloudstorageregion := d.Get(helpers.PIInstanceCaptureCloudStorageRegion).(string)
	if cloudstorageregion == "" {
		log.Printf("CloudStorageRegion is not provided")
	}

	client := st.NewIBMPIInstanceClient(sess, powerinstanceid)

	body := &models.PVMInstanceCapture{
		CaptureDestination:    ptrToString(capturedestination),
		CaptureName:           ptrToString(capturename),
		CaptureVolumeIds:      nil,
		CloudStorageAccessKey: "",
		CloudStorageImagePath: cloudstorageImagePath,
		//CloudStorageRegion:   ptrToString(cloudstorageregion),
		CloudStorageSecretKey: "",
	}

	captureinfo, err := client.CaptureInstanceToImageCatalog(name, powerinstanceid, &p_cloud_p_vm_instances.PcloudPvminstancesCapturePostParams{
		Body: body,
	}, createTimeOut)

	log.Printf("Printing the data from the capture %+v", &captureinfo)

	if err != nil {
		return errors.New("The capture cannot be performed")
	}

	// If this is an image catalog then we need to check what the status is

	imageClient := st.NewIBMPIImageClient(sess, powerinstanceid)
	imagedata, err := imageClient.Get(d.Get(helpers.PIInstanceCaptureName).(string), powerinstanceid)

	if err != nil {
		return err
	}
	log.Printf("Printing the data %s - %s", *imagedata.ImageID, imagedata.State)

	_, err = isWaitForImageCaptureAvailable(client, *imagedata.ImageID, powerinstanceid, d.Timeout(schema.TimeoutCreate))

	//_, err = isWaitForIBMPIVolumeAvailable(client, d.Id(), powerinstanceid, d.Timeout(schema.TimeoutCreate))
	//if err != nil {
	//	return err
	//}
	return nil
	//return resourceIBMPIVolumeAttachRead(d, meta)

}

func resourceIBMPICaptureRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceIBMPICaptureUpdate(d *schema.ResourceData, meta interface{}) error {

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

	_, err = isWaitForIBMPIVolumeAvailable(client, *volrequest.VolumeID, powerinstanceid, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	return resourceIBMPIVolumeRead(d, meta)
}

func resourceIBMPICaptureDelete(d *schema.ResourceData, meta interface{}) error {

	sess, _ := meta.(ClientSession).IBMPISession()
	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	client := st.NewIBMPIVolumeClient(sess, powerinstanceid)

	err := client.Delete(d.Id(), powerinstanceid, deleteTimeOut)
	if err != nil {
		return err
	}

	// wait for power volume states to be back as available. if it's attached it will be in-use
	d.SetId("")
	return nil
}

func isWaitForImageCaptureAvailable(client *st.IBMPIInstanceClient, s string, s2 string, timeout time.Duration) (interface{}, error) {

	return nil, nil
}
