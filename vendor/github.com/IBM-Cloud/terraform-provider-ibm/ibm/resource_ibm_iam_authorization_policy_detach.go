// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMIAMAuthorizationPolicyDetach() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMIAMAuthorizationPolicyDetachCreate,
		Read:   resourceIBMIAMAuthorizationPolicyDetachRead,
		Delete: resourceIBMIAMAuthorizationPolicyDetachDelete,
		Exists: resourceIBMIAMAuthorizationPolicyDetachExists,

		Schema: map[string]*schema.Schema{
			"authorization_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Authorization policy ID",
			},
		},
	}
}

func resourceIBMIAMAuthorizationPolicyDetachCreate(d *schema.ResourceData, meta interface{}) error {
	iampapClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}
	policyID := d.Get("authorization_policy_id").(string)

	deletePolicyOptions := iampapClient.NewDeletePolicyOptions(
		policyID,
	)
	_, err = iampapClient.DeletePolicy(deletePolicyOptions)
	if err != nil {
		return fmt.Errorf("Error detaching authorization policy: %s", err)
	}

	d.SetId(time.Now().UTC().String())

	return resourceIBMIAMAuthorizationPolicyDetachRead(d, meta)
}

func resourceIBMIAMAuthorizationPolicyDetachRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIBMIAMAuthorizationPolicyDetachDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}

func resourceIBMIAMAuthorizationPolicyDetachExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	if d.Id() == "" {
		return false, nil
	}
	return true, nil
}
