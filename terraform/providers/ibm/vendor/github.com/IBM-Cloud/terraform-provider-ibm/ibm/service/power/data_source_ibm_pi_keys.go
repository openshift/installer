// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIKeysRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_Keys: {
				Computed:    true,
				Description: "List of all the SSH keys.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_CreationDate: {
							Computed:    true,
							Description: "Date of SSH key creation.",
							Type:        schema.TypeString,
						},
						Attr_Name: {
							Computed:    true,
							Description: "User defined name for the SSH key.",
							Type:        schema.TypeString,
						},
						Attr_SSHKey: {
							Computed:    true,
							Description: "SSH RSA key.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPIKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	client := instance.NewIBMPIKeyClient(ctx, sess, cloudInstanceID)
	sshKeys, err := client.GetAll()
	if err != nil {
		log.Printf("[ERROR] get all keys failed %v", err)
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(sshKeys.SSHKeys))
	for _, sshKey := range sshKeys.SSHKeys {
		key := map[string]interface{}{
			Attr_CreationDate: sshKey.CreationDate.String(),
			Attr_Name:         sshKey.Name,
			Attr_SSHKey:       sshKey.SSHKey,
		}
		result = append(result, key)
	}
	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(Attr_Keys, result)

	return nil
}
