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

func ResourceIBMEnSafariDestination() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEnSafariDestinationCreate,
		ReadContext:   resourceIBMEnSafariDestinationRead,
		UpdateContext: resourceIBMEnSafariDestinationUpdate,
		DeleteContext: resourceIBMEnSafariDestinationDelete,
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
			"certificate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Certificate File.",
			},
			"icon_16x16": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Certificate File.",
			},
			"icon_16x16_2x": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Certificate File.",
			},
			"icon_32x32": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Certificate File.",
			},
			"icon_32x32_2x": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Certificate File.",
			},
			"icon_128x128": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Certificate File.",
			},
			"icon_128x128_2x": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Certificate File.",
			},
			"icon_16x16_content_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Certificate File.",
			},
			"icon_16x16_2x_content_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Certificate File.",
			},
			"icon_32x32_content_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Certificate File.",
			},
			"icon_32x32_2x_content_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Certificate File.",
			},
			"icon_128x128_content_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Certificate File.",
			},
			"icon_128x128_2x_content_type": {
				Type:        schema.TypeString,
				Optional:    true,
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
									"password": {
										Type:        schema.TypeString,
										Sensitive:   true,
										Required:    true,
										Description: "The Password for APNS Certificate in case of P12 certificate",
									},
									"url_format_string": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Key ID In case of P8 Certificate",
									},
									"website_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Team ID In case of P8 Certificate",
									},
									"website_push_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The Bundle ID In case of P8 Certificate",
									},
									"website_url": {
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

func resourceIBMEnSafariDestinationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.CreateDestinationOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetName(d.Get("name").(string))

	options.SetType(d.Get("type").(string))

	// options.SetCertificateContentType(d.Get("certificate_content_type").(string))
	if _, ok := d.GetOk("icon_16x16_content_type"); ok {
		options.SetIcon16x16ContentType(d.Get("icon_16x16_content_type").(string))
	}
	if _, ok := d.GetOk("icon_16x16_2x_content_type"); ok {
		options.SetIcon16x162xContentType(d.Get("icon_16x16_2x_content_type").(string))
	}
	if _, ok := d.GetOk("icon_32x32_content_type"); ok {
		options.SetIcon32x32ContentType(d.Get("icon_32x32_content_type").(string))
	}
	if _, ok := d.GetOk("icon_32x32_2x_content_type"); ok {
		options.SetIcon32x322xContentType(d.Get("icon_32x32_2x_content_type").(string))
	}
	if _, ok := d.GetOk("icon_128x128_content_type"); ok {
		options.SetIcon128x128ContentType(d.Get("icon_128x128_content_type").(string))
	}
	if _, ok := d.GetOk("icon_128x128_2x_content_type"); ok {
		options.SetIcon128x1282xContentType(d.Get("icon_128x128_2x_content_type").(string))
	}

	// certificatetype := d.Get("certificate_content_type").(string)

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

	if i16, ok := d.GetOk("icon_16x16"); ok {
		path := i16.(string)
		file, err := os.Open(path)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error opening icon file (%s): %s", path, err))
		}

		icon16_16 := file
		options.SetIcon16x16(icon16_16)
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Failed closing icon file (%s): %s", path, err)
			}
		}()
	}

	if i162x, ok := d.GetOk("icon_16x16_2x"); ok {
		path := i162x.(string)
		file, err := os.Open(path)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error opening icon file (%s): %s", path, err))
		}

		icon16_16_2x := file
		options.SetIcon16x162x(icon16_16_2x)
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Failed closing icon file (%s): %s", path, err)
			}
		}()
	}

	if i32, ok := d.GetOk("icon_32x32"); ok {
		path := i32.(string)
		file, err := os.Open(path)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error opening icon file (%s): %s", path, err))
		}

		icon32_32 := file
		options.SetIcon32x32(icon32_32)
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Failed closing icon file (%s): %s", path, err)
			}
		}()
	}

	if i322x, ok := d.GetOk("icon_32x32_2x"); ok {
		path := i322x.(string)
		file, err := os.Open(path)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error opening icon file (%s): %s", path, err))
		}

		icon32_32_2x := file
		options.SetIcon32x322x(icon32_32_2x)
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Failed closing icon file (%s): %s", path, err)
			}
		}()
	}

	if i128, ok := d.GetOk("icon_128x128"); ok {
		path := i128.(string)
		file, err := os.Open(path)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error opening icon file (%s): %s", path, err))
		}

		icon128_128 := file
		options.SetIcon128x128(icon128_128)
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Failed closing icon file (%s): %s", path, err)
			}
		}()
	}

	if i1282x, ok := d.GetOk("icon_128x128_2x"); ok {
		path := i1282x.(string)
		file, err := os.Open(path)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error opening icon file (%s): %s", path, err))
		}

		icon128_128_2x := file
		options.SetIcon128x1282x(icon128_128_2x)
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Failed closing icon file (%s): %s", path, err)
			}
		}()
	}

	if _, ok := d.GetOk("description"); ok {
		options.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("config"); ok {
		config := SafaridestinationConfigMapToDestinationConfig(d.Get("config.0.params.0").(map[string]interface{}))
		options.SetConfig(&config)
	}

	result, response, err := enClient.CreateDestinationWithContext(context, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("CreateDestinationWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *result.ID))

	return resourceIBMEnSafariDestinationRead(context, d, meta)
}

func resourceIBMEnSafariDestinationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		err = d.Set("config", enSafariDestinationFlattenConfig(*result.Config))
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

func resourceIBMEnSafariDestinationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if ok := d.HasChanges("name", "description", "certificate", "icon_16x16", "icon_16x16_2x", "icon_32x32", "icon_32x32_2x", "icon_128x128", "icon_128x128_2x", "config"); ok {
		options.SetName(d.Get("name").(string))

		if _, ok := d.GetOk("description"); ok {
			options.SetDescription(d.Get("description").(string))
		}

		// certificatetype := d.Get("certificate_content_type").(string)

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

		if i16, ok := d.GetOk("icon_16x16"); ok {
			path := i16.(string)
			file, err := os.Open(path)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error opening icon file (%s): %s", path, err))
			}

			icon16_16 := file
			options.SetIcon16x16(icon16_16)
			defer func() {
				err := file.Close()
				if err != nil {
					log.Printf("[WARN] Failed closing icon file (%s): %s", path, err)
				}
			}()
		}

		if i162x, ok := d.GetOk("icon_16x16_2x"); ok {
			path := i162x.(string)
			file, err := os.Open(path)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error opening icon file (%s): %s", path, err))
			}

			icon16_16_2x := file
			options.SetIcon16x162x(icon16_16_2x)
			defer func() {
				err := file.Close()
				if err != nil {
					log.Printf("[WARN] Failed closing icon file (%s): %s", path, err)
				}
			}()
		}

		if i32, ok := d.GetOk("icon_32x32"); ok {
			path := i32.(string)
			file, err := os.Open(path)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error opening icon file (%s): %s", path, err))
			}

			icon32_32 := file
			options.SetIcon32x32(icon32_32)
			defer func() {
				err := file.Close()
				if err != nil {
					log.Printf("[WARN] Failed closing icon file (%s): %s", path, err)
				}
			}()
		}

		if i322x, ok := d.GetOk("icon_32x32_2x"); ok {
			path := i322x.(string)
			file, err := os.Open(path)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error opening icon file (%s): %s", path, err))
			}

			icon32_32_2x := file
			options.SetIcon32x322x(icon32_32_2x)
			defer func() {
				err := file.Close()
				if err != nil {
					log.Printf("[WARN] Failed closing icon file (%s): %s", path, err)
				}
			}()
		}

		if i128, ok := d.GetOk("icon_128x128"); ok {
			path := i128.(string)
			file, err := os.Open(path)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error opening icon file (%s): %s", path, err))
			}

			icon128_128 := file
			options.SetIcon128x128(icon128_128)
			defer func() {
				err := file.Close()
				if err != nil {
					log.Printf("[WARN] Failed closing icon file (%s): %s", path, err)
				}
			}()
		}

		if i1282x, ok := d.GetOk("icon_128x128_2x"); ok {
			path := i1282x.(string)
			file, err := os.Open(path)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error opening icon file (%s): %s", path, err))
			}

			icon128_128_2x := file
			options.SetIcon128x1282x(icon128_128_2x)
			defer func() {
				err := file.Close()
				if err != nil {
					log.Printf("[WARN] Failed closing icon file (%s): %s", path, err)
				}
			}()
		}
		if _, ok := d.GetOk("config"); ok {
			config := SafaridestinationConfigMapToDestinationConfig(d.Get("config.0.params.0").(map[string]interface{}))
			options.SetConfig(&config)
		}
		_, response, err := enClient.UpdateDestinationWithContext(context, options)
		if err != nil {
			return diag.FromErr(fmt.Errorf("UpdateDestinationWithContext failed %s\n%s", err, response))
		}

		return resourceIBMEnSafariDestinationRead(context, d, meta)
	}

	return nil
}

func resourceIBMEnSafariDestinationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func SafaridestinationConfigMapToDestinationConfig(configParams map[string]interface{}) en.DestinationConfig {
	params := new(en.DestinationConfigParams)
	if configParams["cert_type"] != nil {
		params.CertType = core.StringPtr(configParams["cert_type"].(string))
	}

	if configParams["password"] != nil {
		params.Password = core.StringPtr(configParams["password"].(string))
	}

	if configParams["url_format_string"] != nil {
		params.URLFormatString = core.StringPtr(configParams["url_format_string"].(string))
	}

	if configParams["website_name"] != nil {
		params.WebsiteName = core.StringPtr(configParams["website_name"].(string))
	}

	if configParams["website_push_id"] != nil {
		params.WebsitePushID = core.StringPtr(configParams["website_push_id"].(string))
	}

	if configParams["website_url"] != nil {
		params.WebsiteURL = core.StringPtr(configParams["website_url"].(string))
	}

	destinationConfig := new(en.DestinationConfig)
	destinationConfig.Params = params
	return *destinationConfig
}
