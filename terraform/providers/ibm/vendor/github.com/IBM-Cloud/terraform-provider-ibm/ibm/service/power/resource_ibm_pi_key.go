// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI cloud instance ID",
			},
			Arg_KeyName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User defined name for the SSH key",
			},
			Arg_Key: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SSH RSA key",
			},

			// Attributes
			Attr_KeyCreationDate: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date of SSH Key creation",
			},
			Attr_KeyID: {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "User defined name for the SSH key (deprecated - replaced by name)",
			},
			Attr_KeyName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User defined name for the SSH key",
			},
			Attr_Key: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SSH RSA key",
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
	sshkey := d.Get(Arg_Key).(string)

	// create key
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
	sshkeyC := st.NewIBMPIKeyClient(ctx, sess, cloudInstanceID)
	sshkeydata, err := sshkeyC.Get(key)
	if err != nil {
		return diag.FromErr(err)
	}

	// set attributes
	d.Set(Attr_KeyName, sshkeydata.Name)
	d.Set(Attr_KeyID, sshkeydata.Name)
	d.Set(Attr_Key, sshkeydata.SSHKey)
	d.Set(Attr_KeyCreationDate, sshkeydata.CreationDate.String())

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
	sshkeyC := st.NewIBMPIKeyClient(ctx, sess, cloudInstanceID)
	err = sshkeyC.Delete(key)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
