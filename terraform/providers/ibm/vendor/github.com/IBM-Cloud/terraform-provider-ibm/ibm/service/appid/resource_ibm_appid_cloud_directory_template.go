package appid

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMAppIDCloudDirectoryTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMAppIDCloudDirectoryTemplateCreate,
		ReadContext:   resourceIBMAppIDCloudDirectoryTemplateRead,
		DeleteContext: resourceIBMAppIDCloudDirectoryTemplateDelete,
		UpdateContext: resourceIBMAppIDCloudDirectoryTemplateUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"template_name": {
				Description:  "The type of email template. This can be `USER_VERIFICATION`, `WELCOME`, `PASSWORD_CHANGED`, `RESET_PASSWORD` or `MFA_VERIFICATION`",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(supportedAppIDCDTemplates, false),
				ForceNew:     true,
			},
			"language": {
				Description: "Preferred language for resource. Format as described at RFC5646. According to the configured languages codes returned from the `GET /management/v4/{tenantId}/config/ui/languages API`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "en",
				ForceNew:    true,
			},
			"subject": {
				Description: "The subject of the email",
				Type:        schema.TypeString,
				Required:    true,
			},
			"html_body": {
				Description: "The HTML body of the email",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"base64_encoded_html_body": {
				Description: "The HTML body of the email encoded in Base64",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"plain_text_body": {
				Description: "The text body of the email.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceIBMAppIDCloudDirectoryTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	idParts := strings.Split(id, "/")

	if len(idParts) < 3 {
		return diag.Errorf("Incorrect ID %s: ID should be a combination of tenantID/templateName/language", id)
	}

	tenantID := idParts[0]
	templateName := idParts[1]
	language := idParts[2]

	template, resp, err := appIDClient.GetTemplateWithContext(ctx, &appid.GetTemplateOptions{
		TenantID:     &tenantID,
		TemplateName: &templateName,
		Language:     &language,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing Cloud Directory template from state", tenantID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error loading AppID Cloud Directory template: %s\n%s", err, resp)
	}

	if template.Subject != nil {
		d.Set("subject", *template.Subject)
	}

	if template.HTMLBody != nil {
		d.Set("html_body", *template.HTMLBody)
	}

	if template.Base64EncodedHTMLBody != nil {
		d.Set("base64_encoded_html_body", *template.Base64EncodedHTMLBody)
	}

	if template.PlainTextBody != nil {
		d.Set("plain_text_body", *template.PlainTextBody)
	}

	d.Set("tenant_id", tenantID)
	d.Set("template_name", templateName)
	d.Set("language", language)

	return nil
}

func resourceIBMAppIDCloudDirectoryTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	templateName := d.Get("template_name").(string)
	language := d.Get("language").(string)

	input := &appid.UpdateTemplateOptions{
		TenantID:     &tenantID,
		TemplateName: &templateName,
		Language:     &language,
		Subject:      helpers.String(d.Get("subject").(string)),
	}

	if htmlBody, ok := d.GetOk("html_body"); ok {
		// do not set HTMLBody, otherwise might run into issues with Cloudflare filtering
		input.Base64EncodedHTMLBody = helpers.String(b64.StdEncoding.EncodeToString([]byte(htmlBody.(string))))
	}

	if textBody, ok := d.GetOk("plain_text_body"); ok {
		input.PlainTextBody = helpers.String(textBody.(string))
	}

	_, resp, err := appIDClient.UpdateTemplateWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error updating AppID Cloud Directory email template: %s\n%s", err, resp)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", tenantID, templateName, language))

	return resourceIBMAppIDCloudDirectoryTemplateRead(ctx, d, meta)
}

func resourceIBMAppIDCloudDirectoryTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	templateName := d.Get("template_name").(string)
	language := d.Get("language").(string)

	resp, err := appIDClient.DeleteTemplateWithContext(ctx, &appid.DeleteTemplateOptions{
		TenantID:     &tenantID,
		TemplateName: &templateName,
		Language:     &language,
	})

	if err != nil {
		return diag.Errorf("Error deleting AppID Cloud Directory email template: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}

func resourceIBMAppIDCloudDirectoryTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// this is just a configuration, can reuse create method
	return resourceIBMAppIDCloudDirectoryTemplateCreate(ctx, d, m)
}
