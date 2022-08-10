// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
)

func ResourceIBMEnAPNSDestination() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEnAPNSDestinationCreate,
		ReadContext:   resourceIBMEnAPNSDestinationRead,
		UpdateContext: resourceIBMEnAPNSDestinationUpdate,
		DeleteContext: resourceIBMEnAPNSDestinationDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Destintion name.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of Destination type push_ios.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Destination description.",
			},
			"certificate_content_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Certificate Content Type to be set p8/p12.",
			},
			"certificate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Certificate File.",
			},
			"config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Payload describing a destination configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"params": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The Certificate Type for IOS, the values are p8/p12.",
									},
									"is_sandbox": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "The flag to determine sandbox or production environment.",
									},
									"password": {
										Type:        schema.TypeString,
										Sensitive:   true,
										Optional:    true,
										Description: "The Password for APNS Certificate in case of P12 certificate",
									},
									"key_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Key ID In case of P8 Certificate",
									},
									"team_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Team ID In case of P8 Certificate",
									},
									"bundle_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Bundle ID In case of P8 Certificate",
									},
								},
							},
						},
					},
				},
			},
			"destination_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination ID",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time.",
			},
			"subscription_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of subscriptions.",
			},
			"subscription_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of subscriptions.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIBMEnAPNSDestinationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.CreateDestinationOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetName(d.Get("name").(string))

	options.SetType(d.Get("type").(string))

	options.SetCertificateContentType(d.Get("certificate_content_type").(string))

	certificatetype := d.Get("certificate_content_type").(string)

	if c, ok := d.GetOk("certificate"); ok {
		path := c.(string)
		file, err := os.Open(path)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error opening Certificate file (%s): %s", path, err))
		}

		certificate := file
		options.SetCertificate(certificate)
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Failed closing Certificate file (%s): %s", path, err)
			}
		}()
	}

	if _, ok := d.GetOk("description"); ok {
		options.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("config"); ok {
		config := APNSdestinationConfigMapToDestinationConfig(d.Get("config.0.params.0").(map[string]interface{}), certificatetype)
		options.SetConfig(&config)
	}

	result, response, err := enClient.CreateDestinationWithContext(context, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("CreateDestinationWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *result.ID))

	return resourceIBMEnAPNSDestinationRead(context, d, meta)
}

func resourceIBMEnAPNSDestinationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.GetDestinationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	result, response, err := enClient.GetDestinationWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("GetDestinationWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("instance_guid", options.InstanceID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting instance_guid: %s", err))
	}

	if err = d.Set("destination_id", options.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting destination_id: %s", err))
	}

	if err = d.Set("name", result.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if err = d.Set("type", result.Type); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
	}

	if err = d.Set("description", result.Description); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
	}

	if result.Config != nil {
		err = d.Set("config", enAPNSDestinationFlattenConfig(*result.Config))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting config %s", err))
		}
	}

	if err = d.Set("updated_at", flex.DateTimeToString(result.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	if err = d.Set("subscription_count", flex.IntValue(result.SubscriptionCount)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting subscription_count: %s", err))
	}

	if err = d.Set("subscription_names", result.SubscriptionNames); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting subscription_names: %s", err))
	}

	return nil
}

func resourceIBMEnAPNSDestinationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.UpdateDestinationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	if ok := d.HasChanges("name", "description", "certificate", "config"); ok {
		options.SetName(d.Get("name").(string))

		if _, ok := d.GetOk("description"); ok {
			options.SetDescription(d.Get("description").(string))
		}

		certificatetype := d.Get("certificate_content_type").(string)

		if c, ok := d.GetOk("certificate"); ok {
			path := c.(string)
			file, err := os.Open(path)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error opening Certificate file (%s): %s", path, err))
			}

			certificate := file
			options.SetCertificate(certificate)
			defer func() {
				err := file.Close()
				if err != nil {
					log.Printf("[WARN] Failed closing Certificate file (%s): %s", path, err)
				}
			}()
		}

		if _, ok := d.GetOk("config"); ok {
			config := APNSdestinationConfigMapToDestinationConfig(d.Get("config.0.params.0").(map[string]interface{}), certificatetype)
			options.SetConfig(&config)
		}
		_, response, err := enClient.UpdateDestinationWithContext(context, options)
		if err != nil {
			return diag.FromErr(fmt.Errorf("UpdateDestinationWithContext failed %s\n%s", err, response))
		}

		return resourceIBMEnAPNSDestinationRead(context, d, meta)
	}

	return nil
}

func resourceIBMEnAPNSDestinationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.DeleteDestinationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	response, err := enClient.DeleteDestinationWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("DeleteDestinationWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func APNSdestinationConfigMapToDestinationConfig(configParams map[string]interface{}, certificatetype string) en.DestinationConfig {
	params := new(en.DestinationConfigParams)
	if certificatetype == "p8" {
		if configParams["cert_type"] != nil {
			params.CertType = core.StringPtr(configParams["cert_type"].(string))
		}

		if configParams["is_sandbox"] != nil {
			params.IsSandbox = core.BoolPtr(configParams["is_sandbox"].(bool))
		}

		if configParams["key_id"] != nil {
			params.KeyID = core.StringPtr(configParams["key_id"].(string))
		}

		if configParams["team_id"] != nil {
			params.TeamID = core.StringPtr(configParams["team_id"].(string))
		}

		if configParams["bundle_id"] != nil {
			params.BundleID = core.StringPtr(configParams["bundle_id"].(string))
		}

	}

	if certificatetype == "p12" {
		if configParams["cert_type"] != nil {
			params.CertType = core.StringPtr(configParams["cert_type"].(string))
		}

		if configParams["is_sandbox"] != nil {
			params.IsSandbox = core.BoolPtr(configParams["is_sandbox"].(bool))
		}

		if configParams["password"] != nil {
			params.Password = core.StringPtr(configParams["password"].(string))
		}

	}

	destinationConfig := new(en.DestinationConfig)
	destinationConfig.Params = params
	return *destinationConfig
}
