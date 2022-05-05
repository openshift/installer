package pushnotification

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/push-notifications-go-sdk/pushservicev1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMPNApplicationChrome() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationChromeRead,

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
		DeprecationMessage: "This service is deprecated. For more information about the deprecation of this service, see here https://www.ibm.com/cloud/blog/announcements/ibm-push-notifications-deprecation",
	}
}

func dataSourceApplicationChromeRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	pushServiceClient, err := meta.(conns.ClientSession).PushServiceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getChromeWebConfOptions := &pushservicev1.GetChromeWebConfOptions{}

	guid := d.Get("guid").(string)
	getChromeWebConfOptions.SetApplicationID(guid)

	chromeWebConf, response, err := pushServiceClient.GetChromeWebConfWithContext(context, getChromeWebConfOptions)
	if err != nil {
		log.Printf("[DEBUG] GetChromeWebConfWithContext failed %s\n%d", err, response.StatusCode)
		return diag.FromErr(err)
	}

	d.SetId(guid)
	if err = d.Set("server_key", chromeWebConf.ApiKey); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting server_key: %s", err))
	}
	if err = d.Set("web_site_url", chromeWebConf.WebSiteURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting web_site_url: %s", err))
	}

	return nil
}
