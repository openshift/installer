// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
)

func dataSourceIBMPIKey() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPIKeyRead,
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

func dataSourceIBMPIKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	sshkeyC := instance.NewIBMPIKeyClient(ctx, sess, cloudInstanceID)
	sshkeydata, err := sshkeyC.Get(d.Get(helpers.PIKeyName).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*sshkeydata.Name)
	d.Set("creation_date", sshkeydata.CreationDate.String())
	d.Set("sshkey", sshkeydata.SSHKey)
	d.Set(helpers.PIKeyName, sshkeydata.Name)

	return nil
}
