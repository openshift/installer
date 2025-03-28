// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventstreams

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/eventstreams-go-sdk/pkg/adminrestv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// A quota in an Event Streams service instance.
// The ID is the CRN with the last two components "quota:entity".
// The producer_byte_rate and consumer_byte_rate are the two quota properties, and must be at least -1;
// -1 means no quota applied.
func ResourceIBMEventStreamsQuota() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEventStreamsQuotaCreate,
		ReadContext:   resourceIBMEventStreamsQuotaRead,
		UpdateContext: resourceIBMEventStreamsQuotaUpdate,
		DeleteContext: resourceIBMEventStreamsQuotaDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"resource_instance_id": {
				Type:        schema.TypeString,
				Description: "The ID or the CRN of the Event Streams service instance",
				Required:    true,
				ForceNew:    true,
			},
			"entity": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The entity for which the quota is set; 'default' or IAM ID",
			},
			"producer_byte_rate": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(-1),
				Description:  "The producer quota in bytes per second, -1 means no quota",
			},
			"consumer_byte_rate": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(-1),
				Description:  "The consumer quota in bytes per second, -1 means no quota",
			},
		},
	}
}

func resourceIBMEventStreamsQuotaCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminrestClient, instanceCRN, entity, err := getQuotaClientInstanceEntity(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "Error getting Event Streams instance", "ibm_event_streams_quota", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createQuotaOptions := &adminrestv1.CreateQuotaOptions{}
	createQuotaOptions.SetEntityName(entity)
	pbr := d.Get("producer_byte_rate").(int)
	cbr := d.Get("consumer_byte_rate").(int)
	if pbr == -1 && cbr == -1 {
		return diag.FromErr(fmt.Errorf("Quota for %s cannot be created: producer_byte_rate and consumer_byte_rate are both -1 (no quota)", entity))
	}
	if pbr != -1 {
		createQuotaOptions.SetProducerByteRate(int64(pbr))
	}
	if cbr != -1 {
		createQuotaOptions.SetConsumerByteRate(int64(cbr))
	}

	response, err := adminrestClient.CreateQuotaWithContext(context, createQuotaOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateQuota failed with response: %s", response), "ibm_event_streams_quota", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.SetId(getQuotaID(instanceCRN, entity))

	return resourceIBMEventStreamsQuotaRead(context, d, meta)
}

func resourceIBMEventStreamsQuotaRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminrestClient, instanceCRN, entity, err := getQuotaClientInstanceEntity(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "Error getting Event Streams instance", "ibm_event_streams_quota", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getQuotaOptions := &adminrestv1.GetQuotaOptions{}
	getQuotaOptions.SetEntityName(entity)
	quota, response, err := adminrestClient.GetQuotaWithContext(context, getQuotaOptions)
	if err != nil || quota == nil {
		d.SetId("")
		var tfErr *flex.TerraformProblem
		if response != nil && response.StatusCode == 404 {
			tfErr = flex.TerraformErrorf(err, fmt.Sprintf("Quota for '%s' does not exist", entity), "ibm_event_streams_quota", "read")
		} else {
			tfErr = flex.TerraformErrorf(err, fmt.Sprintf("GetQuota failed with response: %s", response), "ibm_event_streams_quota", "read")
		}
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.Set("resource_instance_id", instanceCRN)
	d.Set("entity", entity)
	d.Set("producer_byte_rate", getQuotaValue(quota.ProducerByteRate))
	d.Set("consumer_byte_rate", getQuotaValue(quota.ConsumerByteRate))
	d.SetId(getQuotaID(instanceCRN, entity))

	return nil
}

func resourceIBMEventStreamsQuotaUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("producer_byte_rate") || d.HasChange("consumer_byte_rate") {
		adminrestClient, _, entity, err := getQuotaClientInstanceEntity(d, meta)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "Error getting Event Streams instance", "ibm_event_streams_quota", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		updateQuotaOptions := &adminrestv1.UpdateQuotaOptions{}
		updateQuotaOptions.SetEntityName(entity)
		updateQuotaOptions.SetProducerByteRate(int64(d.Get("producer_byte_rate").(int)))
		updateQuotaOptions.SetConsumerByteRate(int64(d.Get("consumer_byte_rate").(int)))

		response, err := adminrestClient.UpdateQuotaWithContext(context, updateQuotaOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateQuota failed with response: %s", response), "ibm_event_streams_quota", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	return resourceIBMEventStreamsQuotaRead(context, d, meta)
}

func resourceIBMEventStreamsQuotaDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminrestClient, _, entity, err := getQuotaClientInstanceEntity(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "Error getting Event Streams instance", "ibm_event_streams_quota", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteQuotaOptions := &adminrestv1.DeleteQuotaOptions{}
	deleteQuotaOptions.SetEntityName(entity)

	response, err := adminrestClient.DeleteQuotaWithContext(context, deleteQuotaOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteQuota failed with response: %s", response), "ibm_event_streams_quota", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.SetId("")
	return nil
}
