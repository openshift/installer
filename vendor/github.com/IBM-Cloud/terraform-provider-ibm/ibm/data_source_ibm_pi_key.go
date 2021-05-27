// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
)

func dataSourceIBMPIKey() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMPIKeysRead,
		Schema: map[string]*schema.Schema{

			helpers.PIKeyName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "SSHKey Name to be used for pvminstances",
				ValidateFunc: validation.NoZeroValues,
			},
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			//Computed Attributes
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sshkey": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func dataSourceIBMPIKeysRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()

	if err != nil {
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	sshkeyC := instance.NewIBMPIKeyClient(sess, powerinstanceid)
	sshkeydata, err := sshkeyC.Get(d.Get(helpers.PIKeyName).(string), powerinstanceid)

	if err != nil {
		return err
	}

	d.SetId(*sshkeydata.Name)
	d.Set("creation_date", sshkeydata.CreationDate.String())
	d.Set("sshkey", sshkeydata.SSHKey)
	d.Set(helpers.PIKeyName, sshkeydata.Name)
	d.Set(helpers.PICloudInstanceId, powerinstanceid)

	return nil

}
