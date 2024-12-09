// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIbmIsShareDeleteAccessorBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsShareDeleteAccessorBindingCreate,
		ReadContext:   resourceIbmIsShareDeleteAccessorBindingRead,
		UpdateContext: resourceIbmIsShareDeleteAccessorBindingUpdate,
		DeleteContext: resourceIbmIsShareDeleteAccessorBindingDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"share": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file share identifier.",
			},
			"accessor_binding": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "The accessor binding id",
			},
		},
	}
}

func resourceIbmIsShareDeleteAccessorBindingCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	share_id := d.Get("share").(string)

	bindingId := d.Get("accessor_binding").(string)

	deleteAccessBindingOptions := &vpcv1.DeleteShareAccessorBindingOptions{
		ShareID: &share_id,
		ID:      &bindingId,
	}
	response, err := vpcClient.DeleteShareAccessorBindingWithContext(context, deleteAccessBindingOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteShareAccessorBindingWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] DeleteShareAccessorBindingWithContext failed %s\n%s", err, response))
	}

	_, err = isWaitForShareAccessorBindingDeleted(context, vpcClient, share_id, bindingId, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(share_id)
	return nil
}

func isWaitForShareAccessorBindingDeleted(context context.Context, vpcClient *vpcv1.VpcV1, shareid, bindingId string, d *schema.ResourceData, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for share accessor binding (%s) to be deleted.", shareid)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "deleting"},
		Target:     []string{"done"},
		Refresh:    isShareAccessorBindingRefreshFunc(context, vpcClient, shareid, bindingId, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isShareAccessorBindingRefreshFunc(context context.Context, vpcClient *vpcv1.VpcV1, shareid string, bindingId string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		bindingOptions := &vpcv1.GetShareAccessorBindingOptions{
			ShareID: &shareid,
			ID:      &bindingId,
		}

		shareBinding, response, err := vpcClient.GetShareAccessorBindingWithContext(context, bindingOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return shareBinding, "done", nil
			}
			return nil, "", fmt.Errorf("Error Getting Target: %s\n%s", err, response)
		}
		return shareBinding, *shareBinding.LifecycleState, nil
	}
}

func resourceIbmIsShareDeleteAccessorBindingRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceIbmIsShareDeleteAccessorBindingUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceIbmIsShareDeleteAccessorBindingDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}
