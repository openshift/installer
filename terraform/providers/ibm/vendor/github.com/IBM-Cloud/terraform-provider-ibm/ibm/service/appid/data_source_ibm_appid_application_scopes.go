package appid

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDApplicationScopes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAppIDApplicationScopesRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"client_id": {
				Description: "The `client_id` is a public identifier for applications",
				Type:        schema.TypeString,
				Required:    true,
			},
			"scopes": {
				Description: "A `scope` is a runtime action in your application that you register with IBM Cloud App ID to create an access permission",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceIBMAppIDApplicationScopesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	scopes, resp, err := appIDClient.GetApplicationScopesWithContext(ctx, &appid.GetApplicationScopesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID application scopes: %s\n%s", err, resp)
	}

	if err := d.Set("scopes", scopes.Scopes); err != nil {
		return diag.Errorf("Error setting AppID application scopes: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, clientID))
	return nil
}
