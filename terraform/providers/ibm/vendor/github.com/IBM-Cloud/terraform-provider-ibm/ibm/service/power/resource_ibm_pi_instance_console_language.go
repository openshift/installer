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

func ResourceIBMPIInstanceConsoleLanguage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIInstanceConsoleLanguageCreate,
		ReadContext:   resourceIBMPIInstanceConsoleLanguageRead,
		UpdateContext: resourceIBMPIInstanceConsoleLanguageUpdate,
		DeleteContext: resourceIBMPIInstanceConsoleLanguageDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_InstanceName: {
				Description:  "The unique identifier or name of the instance.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_LanguageCode: {
				Description:  "Language code.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceIBMPIInstanceConsoleLanguageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	instanceName := d.Get(Arg_InstanceName).(string)
	code := d.Get(Arg_LanguageCode).(string)

	client := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)

	consoleLanguage := &models.ConsoleLanguage{
		Code: &code,
	}

	_, err = client.UpdateConsoleLanguage(instanceName, consoleLanguage)
	if err != nil {
		log.Printf("[DEBUG] err %s", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, instanceName))

	return resourceIBMPIInstanceConsoleLanguageRead(ctx, d, meta)
}

func resourceIBMPIInstanceConsoleLanguageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// There is no get concept for instance console language
	return nil
}

func resourceIBMPIInstanceConsoleLanguageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange(Arg_LanguageCode) {
		cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
		instanceName := d.Get(Arg_InstanceName).(string)
		code := d.Get(Arg_LanguageCode).(string)

		client := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)

		consoleLanguage := &models.ConsoleLanguage{
			Code: &code,
		}
		_, err = client.UpdateConsoleLanguage(instanceName, consoleLanguage)
		if err != nil {
			log.Printf("[DEBUG] err %s", err)
			return diag.FromErr(err)
		}
	}
	return resourceIBMPIInstanceConsoleLanguageRead(ctx, d, meta)
}

func resourceIBMPIInstanceConsoleLanguageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// There is no delete or unset concept for instance console language
	d.SetId("")
	return nil
}
