// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
)

const (
	PIKeys    = "keys"
	PIKeyName = "name"
	PIKey     = "ssh_key"
	PIKeyDate = "creation_date"
)

func dataSourceIBMPIKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIKeysRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "PI cloud instance ID",
				ValidateFunc: validation.NoZeroValues,
			},
			// Computed Attributes
			PIKeys: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PIKeyName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User defined name for the SSH key",
						},
						PIKey: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SSH RSA key",
						},
						PIKeyDate: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date of SSH key creation",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	client := st.NewIBMPIKeyClient(ctx, sess, cloudInstanceID)
	sshKeys, err := client.GetAll()
	if err != nil {
		log.Printf("[ERROR] get all keys failed %v", err)
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(sshKeys.SSHKeys))
	for _, sshKey := range sshKeys.SSHKeys {
		key := map[string]interface{}{
			PIKeyName: sshKey.Name,
			PIKey:     sshKey.SSHKey,
			PIKeyDate: sshKey.CreationDate.String(),
		}
		result = append(result, key)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(PIKeys, result)

	return nil
}
