// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisWebhookID     = "webhook_id"
	cisWebhookName   = "name"
	cisWebhookURL    = "url"
	cisWebhookType   = "type"
	cisWebhookSecret = "secret"
)

func ResourceIBMCISWebhooks() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISWebhookCreate,
		Read:     ResourceIBMCISWebhookRead,
		Update:   ResourceIBMCISWebhookUpdate,
		Delete:   ResourceIBMCISWebhookDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
			},
			cisWebhookID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Webhook ID",
			},
			cisWebhookName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Webhook Name",
			},
			cisWebhookURL: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Webhook URL",
			},
			cisWebhookType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Webhook Type",
			},
			cisWebhookSecret: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "API key needed to use the webhook",
			},
		},
	}
}
func ResourceIBMCISWebhookCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisWebhookSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the cisWebhookSession %s", err)
	}

	crn := d.Get(cisID).(string)
	sess.Crn = core.StringPtr(crn)

	opt := sess.NewCreateAlertWebhookOptions()

	if name, ok := d.GetOk(cisWebhookName); ok {
		opt.SetName(name.(string))
	}
	if url, ok := d.GetOk(cisWebhookURL); ok {
		opt.SetURL(url.(string))
	}
	if secret, ok := d.GetOk(cisWebhookSecret); ok {
		opt.SetSecret((secret.(string)))
	}
	result, resp, err := sess.CreateAlertWebhook(opt)
	if err != nil || result == nil {
		return fmt.Errorf("[ERROR] Error creating Webhooks  %s %s", err, resp)
	}
	d.SetId(flex.ConvertCisToTfTwoVar(*result.Result.ID, crn))
	return ResourceIBMCISWebhookRead(d, meta)

}
func ResourceIBMCISWebhookRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisWebhookSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the cisWebhookSession %s", err)
	}
	webhooksID, crn, err := flex.ConvertTftoCisTwoVar(d.Id())
	if err != nil {
		return err
	}
	sess.Crn = core.StringPtr(crn)
	opt := sess.NewGetWebhookOptions(webhooksID)

	result, response, err := sess.GetWebhook(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting webhook detail %s, %s", err, response)
	}
	d.Set(cisID, crn)
	d.Set(cisWebhookID, result.Result.ID)
	d.Set(cisWebhookName, result.Result.Name)
	d.Set(cisWebhookURL, result.Result.URL)
	d.Set(cisWebhookType, result.Result.Type)
	return nil
}
func ResourceIBMCISWebhookUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisWebhookSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while updating the webhook %s", err)
	}
	webhooksID, crn, err := flex.ConvertTftoCisTwoVar(d.Id())
	if err != nil {
		return err
	}
	sess.Crn = core.StringPtr(crn)
	opt := sess.NewUpdateAlertWebhookOptions(webhooksID)
	if d.HasChange(cisWebhookName) ||
		d.HasChange(cisWebhookURL) ||
		d.HasChange(cisWebhookSecret) {

		if name, ok := d.GetOk(cisWebhookName); ok {
			opt.SetName(name.(string))
		}
		if url, ok := d.GetOk(cisWebhookURL); ok {
			opt.SetURL(url.(string))
		}
		if secret, ok := d.GetOk(cisWebhookSecret); ok {
			opt.SetSecret((secret.(string)))
		}

		result, _, err := sess.UpdateAlertWebhook(opt)
		if err != nil || result == nil {
			return fmt.Errorf("[ERROR] Error updating the Webhook %s", err)
		}
	}
	return ResourceIBMCISWebhookRead(d, meta)
}
func ResourceIBMCISWebhookDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisWebhookSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while Deleting the webhook %s", err)
	}
	webhooksID, crn, err := flex.ConvertTftoCisTwoVar(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the webhook ID %s", err)
	}
	sess.Crn = core.StringPtr(crn)

	opt := sess.NewDeleteWebhookOptions(webhooksID)

	_, response, err := sess.DeleteWebhook(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error deleting the Webhook %s:%s", err, response)
	}
	return nil

}
