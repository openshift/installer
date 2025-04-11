// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMEnSMTPSetting() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"smtp_config_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"settings": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnets": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The SMTP allowed Ips.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
		DeprecationMessage: "The resource has been deprecated since the support for legacy allowlisting has been deprecated. The support has been enabled via Context-based-restrictions. For detailed information, please refer here: https://cloud.ibm.com/docs/event-notifications?topic=event-notifications-en-smtp-configurations#en-smtp-configurations-cbr",
	}
}
