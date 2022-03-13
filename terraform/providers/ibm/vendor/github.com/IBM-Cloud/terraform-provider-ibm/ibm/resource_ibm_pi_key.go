// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

func resourceIBMPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIKeyCreate,
		ReadContext:   resourceIBMPIKeyRead,
		UpdateContext: resourceIBMPIKeyUpdate,
		DeleteContext: resourceIBMPIKeyDelete,
		Importer:      &schema.ResourceImporter{},

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

func resourceIBMPIKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	name := d.Get(helpers.PIKeyName).(string)
	sshkey := d.Get(helpers.PIKey).(string)

	client := st.NewIBMPIKeyClient(ctx, sess, cloudInstanceID)
	body := &models.SSHKey{
		Name:   &name,
		SSHKey: &sshkey,
	}
	sshResponse, err := client.Create(body)
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return diag.FromErr(err)
	}

	log.Printf("Printing the sshkey %+v", *sshResponse)

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, name))
	return resourceIBMPIKeyRead(ctx, d, meta)
}

func resourceIBMPIKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, key, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	sshkeyC := st.NewIBMPIKeyClient(ctx, sess, cloudInstanceID)
	sshkeydata, err := sshkeyC.Get(key)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(helpers.PIKeyName, sshkeydata.Name)
	d.Set(helpers.PIKey, sshkeydata.SSHKey)
	d.Set(helpers.PIKeyDate, sshkeydata.CreationDate.String())
	d.Set("key_id", sshkeydata.Name)

	return nil

}
func resourceIBMPIKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceIBMPIKeyRead(ctx, d, meta)
}
func resourceIBMPIKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, key, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	sshkeyC := st.NewIBMPIKeyClient(ctx, sess, cloudInstanceID)
	err = sshkeyC.Delete(key)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
