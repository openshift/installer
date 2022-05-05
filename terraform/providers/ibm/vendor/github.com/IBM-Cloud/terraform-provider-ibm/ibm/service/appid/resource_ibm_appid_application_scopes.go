package appid

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMAppIDApplicationScopes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMAppIDApplicationScopesCreate,
		ReadContext:   resourceIBMAppIDApplicationScopesRead,
		DeleteContext: resourceIBMAppIDApplicationScopesDelete,
		UpdateContext: resourceIBMAppIDApplicationScopesUpdate,
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
			"client_id": {
				Description: "The `client_id` is a public identifier for applications",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"scopes": {
				Description: "A `scope` is a runtime action in your application that you register with IBM Cloud App ID to create an access permission",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

func resourceIBMAppIDApplicationScopesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	scopes := flex.ExpandStringList(d.Get("scopes").([]interface{}))

	scopeOpts := &appid.PutApplicationsScopesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
		Scopes:   scopes,
	}

	_, resp, err := appIDClient.PutApplicationsScopesWithContext(ctx, scopeOpts)

	if err != nil {
		return diag.Errorf("Error setting application scopes: %s\n%s", err, resp)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, clientID))

	return resourceIBMAppIDApplicationScopesRead(ctx, d, meta)
}

func resourceIBMAppIDApplicationScopesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	idParts := strings.Split(id, "/")

	if len(idParts) < 2 {
		return diag.Errorf("Incorrect ID %s: ID should be a combination of tenantID/clientID", d.Id())
	}

	tenantID := idParts[0]
	clientID := idParts[1]

	scopes, resp, err := appIDClient.GetApplicationScopesWithContext(ctx, &appid.GetApplicationScopesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID application '%s' is not found, removing scopes from state", clientID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error getting AppID application scopes: %s\n%s", err, resp)
	}

	if err := d.Set("scopes", scopes.Scopes); err != nil {
		return diag.Errorf("Error setting application scopes: %s", err)
	}

	d.Set("tenant_id", tenantID)
	d.Set("client_id", clientID)

	return nil
}

func resourceIBMAppIDApplicationScopesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	scopes := flex.ExpandStringList(d.Get("scopes").([]interface{}))

	scopeOpts := &appid.PutApplicationsScopesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
		Scopes:   scopes,
	}

	_, resp, err := appIDClient.PutApplicationsScopesWithContext(ctx, scopeOpts)

	if err != nil {
		return diag.Errorf("Error updating application scopes: %s\n%s", err, resp)
	}

	return resourceIBMAppIDApplicationScopesRead(ctx, d, meta)
}

func resourceIBMAppIDApplicationScopesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	scopeOpts := &appid.PutApplicationsScopesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
		Scopes:   []string{},
	}

	_, resp, err := appIDClient.PutApplicationsScopesWithContext(ctx, scopeOpts)

	if err != nil {
		return diag.Errorf("Error clearing application scopes: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
