package ibm

import (
	"fmt"

	"github.com/IBM/push-notifications-go-sdk/pushservicev1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMPNApplicationChrome() *schema.Resource {
	return &schema.Resource{
		Read:     resourceApplicationChromeRead,
		Create:   resourceApplicationChromeCreate,
		Update:   resourceApplicationChromeUpdate,
		Delete:   resourceApplicationChromeDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique guid of the push notification instance.",
			},
			"server_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A server key that gives the push service an authorized access to Google services that is used for Chrome Web Push.",
			},
			"web_site_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the WebSite / WebApp that should be permitted to subscribe to WebPush.",
			},
		},
	}
}

func resourceApplicationChromeCreate(d *schema.ResourceData, meta interface{}) error {
	pnClient, err := meta.(ClientSession).PushServiceV1()
	if err != nil {
		return err
	}

	serverKey := d.Get("server_key").(string)
	websiteURL := d.Get("web_site_url").(string)
	guid := d.Get("guid").(string)

	_, response, err := pnClient.SaveChromeWebConf(&pushservicev1.SaveChromeWebConfOptions{
		ApplicationID: &guid,
		ApiKey:        &serverKey,
		WebSiteURL:    &websiteURL,
	})

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error configuring chrome web platform: %s with responce code  %d", err, response.StatusCode)
	}
	d.SetId(guid)

	return resourceApplicationChromeRead(d, meta)
}

func resourceApplicationChromeUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChanges("server_key", "web_site_url") {
		return resourceApplicationChromeCreate(d, meta)
	}
	return nil
}

func resourceApplicationChromeRead(d *schema.ResourceData, meta interface{}) error {
	pnClient, err := meta.(ClientSession).PushServiceV1()
	if err != nil {
		return err
	}

	guid := d.Id()

	chromeWebConf, response, err := pnClient.GetChromeWebConf(&pushservicev1.GetChromeWebConfOptions{
		ApplicationID: &guid,
	})

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error fetching chrome web platform configuration: %s with responce code  %d", err, response.StatusCode)
	}

	d.SetId(guid)

	if response.StatusCode == 200 {
		d.Set("server_key", *chromeWebConf.ApiKey)
		d.Set("web_site_url", *chromeWebConf.WebSiteURL)
	}
	return nil
}

func resourceApplicationChromeDelete(d *schema.ResourceData, meta interface{}) error {
	pnClient, err := meta.(ClientSession).PushServiceV1()
	if err != nil {
		return err
	}
	guid := d.Get("guid").(string)

	response, err := pnClient.DeleteChromeWebConf(&pushservicev1.DeleteChromeWebConfOptions{
		ApplicationID: &guid,
	})

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error deleting chrome web platform configuration: %s with responce code  %d", err, response.StatusCode)
	}

	d.SetId("")

	return nil

}
