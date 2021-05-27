// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
)

func resourceIBMPIKey() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMPIKeyCreate,
		Read:     resourceIBMPIKeyRead,
		Update:   resourceIBMPIKeyUpdate,
		Delete:   resourceIBMPIKeyDelete,
		Exists:   resourceIBMPIKeyExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			helpers.PIKeyName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key name in the PI instance",
			},

			helpers.PIKey: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI instance key info",
			},
			helpers.PIKeyDate: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date info",
			},

			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI cloud instance ID",
			},

			"key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Key ID in the PI instance",
			},
		},
	}
}

func resourceIBMPIKeyCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	name := d.Get(helpers.PIKeyName).(string)
	sshkey := d.Get(helpers.PIKey).(string)
	client := st.NewIBMPIKeyClient(sess, powerinstanceid)

	sshResponse, _, err := client.Create(name, sshkey, powerinstanceid)
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return fmt.Errorf("Failed to create the key %v", err)

	}

	log.Printf("Printing the sshkey %+v", &sshResponse)

	d.SetId(fmt.Sprintf("%s/%s", powerinstanceid, name))
	return resourceIBMPIKeyRead(d, meta)
}

func resourceIBMPIKeyRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return fmt.Errorf("Failed to obtain the key %v", err)
	}

	powerinstanceid := parts[0]
	sshkeyC := st.NewIBMPIKeyClient(sess, powerinstanceid)
	sshkeydata, err := sshkeyC.Get(parts[1], powerinstanceid)

	if err != nil {
		return err
	}

	d.Set(helpers.PIKeyName, sshkeydata.Name)
	d.Set(helpers.PIKey, sshkeydata.SSHKey)
	d.Set(helpers.PIKeyDate, sshkeydata.CreationDate.String())
	d.Set("key_id", sshkeydata.Name)
	d.Set(helpers.PICloudInstanceId, powerinstanceid)

	return nil

}

func resourceIBMPIKeyUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIBMPIKeyDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	powerinstanceid := parts[0]
	sshkeyC := st.NewIBMPIKeyClient(sess, powerinstanceid)
	err = sshkeyC.Delete(parts[1], powerinstanceid)

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceIBMPIKeyExists(d *schema.ResourceData, meta interface{}) (bool, error) {

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
	client := st.NewIBMPIKeyClient(sess, powerinstanceid)

	key, err := client.Get(parts[1], powerinstanceid)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}
	return *key.Name == name, nil
}
