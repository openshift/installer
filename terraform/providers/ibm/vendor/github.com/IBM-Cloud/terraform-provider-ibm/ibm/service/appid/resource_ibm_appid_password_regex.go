package appid

import (
	"context"
	b64 "encoding/base64"
	"log"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMAppIDPasswordRegex() *schema.Resource {
	return &schema.Resource{
		Description:   "The regular expression used by App ID for password strength validation",
		CreateContext: resourceIBMAppIDPasswordRegexCreate,
		ReadContext:   resourceIBMAppIDPasswordRegexRead,
		DeleteContext: resourceIBMAppIDPasswordRegexDelete,
		UpdateContext: resourceIBMAppIDPasswordRegexUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"base64_encoded_regex": {
				Description: "The regex expression rule for acceptable password encoded in base64",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"error_message": {
				Description: "Custom error message",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"regex": {
				Description: "The escaped regex expression rule for acceptable password",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceIBMAppIDPasswordRegexRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Id()

	pw, resp, err := appIDClient.GetCloudDirectoryPasswordRegexWithContext(ctx, &appid.GetCloudDirectoryPasswordRegexOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing Password Regex from state", tenantID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error loading AppID Cloud Directory password regex: %s\n%s", err, resp)
	}

	if pw.Base64EncodedRegex != nil {
		d.Set("base64_encoded_regex", *pw.Base64EncodedRegex)
	}

	if pw.Regex != nil {
		d.Set("regex", *pw.Regex)
	}

	if pw.ErrorMessage != nil {
		d.Set("error_message", *pw.ErrorMessage)
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDPasswordRegexCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	regex := d.Get("regex").(string)

	input := &appid.SetCloudDirectoryPasswordRegexOptions{
		TenantID:           &tenantID,
		Base64EncodedRegex: helpers.String(b64.StdEncoding.EncodeToString([]byte(regex))),
	}

	if msg, ok := d.GetOk("error_message"); ok {
		input.ErrorMessage = helpers.String(msg.(string))
	}

	_, resp, err := appIDClient.SetCloudDirectoryPasswordRegexWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error setting AppID Cloud Directory password regex: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDPasswordRegexRead(ctx, d, meta)
}

func resourceIBMAppIDPasswordRegexUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceIBMAppIDPasswordRegexCreate(ctx, d, meta)
}

func resourceIBMAppIDPasswordRegexDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.SetCloudDirectoryPasswordRegexOptions{
		TenantID:           &tenantID,
		Base64EncodedRegex: helpers.String(""),
	}

	_, resp, err := appIDClient.SetCloudDirectoryPasswordRegexWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error resetting AppID Cloud Directory password regex: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
