// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMPIKey() *schema.Resource {
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
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_KeyName: {
				Description:  "User defined name for the SSH key.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_SSHKey: {
				Description:  "SSH RSA key.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_CreationDate: {
				Computed:    true,
				Description: "Date of SSH Key creation.",
				Type:        schema.TypeString,
			},
			Attr_Name: {
				Computed:    true,
				Description: "User defined name for the SSH key.",
				Type:        schema.TypeString,
			},
			Attr_Key: {
				Computed:    true,
				Description: "SSH RSA key.",
				Type:        schema.TypeString,
			},
		},
	}
}

func resourceIBMPIKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// session
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	// arguments
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	name := d.Get(Arg_KeyName).(string)
	sshkey := d.Get(Arg_SSHKey).(string)

	// create key
	client := instance.NewIBMPIKeyClient(ctx, sess, cloudInstanceID)
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
	// session
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	// arguments
	cloudInstanceID, key, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// get key
	sshkeyC := instance.NewIBMPIKeyClient(ctx, sess, cloudInstanceID)
	sshkeydata, err := sshkeyC.Get(key)
	if err != nil {
		return diag.FromErr(err)
	}

	// set attributes
	d.Set(Attr_CreationDate, sshkeydata.CreationDate.String())
	d.Set(Attr_Key, sshkeydata.SSHKey)
	d.Set(Attr_Name, sshkeydata.Name)

	return nil
}

func resourceIBMPIKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceIBMPIKeyRead(ctx, d, meta)
}

func resourceIBMPIKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// session
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	// arguments
	cloudInstanceID, key, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// delete key
	sshkeyC := instance.NewIBMPIKeyClient(ctx, sess, cloudInstanceID)
	err = sshkeyC.Delete(key)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
