package appid

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var supportedAppIDCDTemplates = []string{"USER_VERIFICATION", "RESET_PASSWORD", "WELCOME", "PASSWORD_CHANGED", "MFA_VERIFICATION"}

func DataSourceIBMAppIDCloudDirectoryTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAppIDCloudDirectoryTemplateRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"template_name": {
				Description:  "The type of email template. This can be `USER_VERIFICATION`, `WELCOME`, `PASSWORD_CHANGED`, `RESET_PASSWORD` or `MFA_VERIFICATION`",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(supportedAppIDCDTemplates, false),
			},
			"language": {
				Description: "Preferred language for resource. Format as described at RFC5646. According to the configured languages codes returned from the `GET /management/v4/{tenantId}/config/ui/languages API`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "en",
			},
			"subject": {
				Description: "The subject of the email",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"html_body": {
				Description: "The HTML body of the email",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"base64_encoded_html_body": {
				Description: "The HTML body of the email encoded in Base64",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"plain_text_body": {
				Description: "The text body of the email.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppIDCloudDirectoryTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	templateName := d.Get("template_name").(string)
	language := d.Get("language").(string)

	template, resp, err := appIDClient.GetTemplateWithContext(ctx, &appid.GetTemplateOptions{
		TenantID:     &tenantID,
		TemplateName: &templateName,
		Language:     &language,
	})

	if err != nil {
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

	d.SetId(fmt.Sprintf("%s/%s/%s", tenantID, templateName, language))

	return nil
}
