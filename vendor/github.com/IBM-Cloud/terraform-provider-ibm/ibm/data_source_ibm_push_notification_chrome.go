package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/push-notifications-go-sdk/pushservicev1"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMPNApplicationChrome() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationChromeRead,

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique guid of the application using the push service.",
			},
			"server_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A server key that gives the push service an authorized access to Google services that is used for Chrome Web Push.",
			},
			"web_site_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the WebSite / WebApp that should be permitted to subscribe to WebPush.",
			},
		},
	}
}

func dataSourceApplicationChromeRead(d *schema.ResourceData, meta interface{}) error {
	pushServiceClient, err := meta.(ClientSession).PushServiceV1()
	if err != nil {
		return err
	}

	getChromeWebConfOptions := &pushservicev1.GetChromeWebConfOptions{}

	guid := d.Get("guid").(string)
	getChromeWebConfOptions.SetApplicationID(guid)

	chromeWebConf, response, err := pushServiceClient.GetChromeWebConfWithContext(context.TODO(), getChromeWebConfOptions)
	if err != nil {
		log.Printf("[DEBUG] GetChromeWebConfWithContext failed %s\n%d", err, response.StatusCode)
		return err
	}

	d.SetId(guid)
	if err = d.Set("server_key", chromeWebConf.ApiKey); err != nil {
		return fmt.Errorf("Error setting server_key: %s", err)
	}
	if err = d.Set("web_site_url", chromeWebConf.WebSiteURL); err != nil {
		return fmt.Errorf("Error setting web_site_url: %s", err)
	}

	return nil
}
