package keyvault

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceKeyVaultSecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeyVaultSecretRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.NestedItemName,
			},

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.VaultID,
			},

			"value": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceKeyVaultSecretRead(d *schema.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := parse.VaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("Error looking up Secret %q vault url from id %q: %+v", name, *keyVaultId, err)
	}

	// we always want to get the latest version
	resp, err := client.GetSecret(ctx, *keyVaultBaseUri, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("KeyVault Secret %q (KeyVault URI %q) does not exist", name, *keyVaultBaseUri)
		}
		return fmt.Errorf("Error making Read request on Azure KeyVault Secret %s: %+v", name, err)
	}

	// the version may have changed, so parse the updated id
	respID, err := parse.ParseNestedItemID(*resp.ID)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	d.Set("name", respID.Name)
	d.Set("key_vault_id", keyVaultId.ID())
	d.Set("value", resp.Value)
	d.Set("version", respID.Version)
	d.Set("content_type", resp.ContentType)

	return tags.FlattenAndSet(d, resp.Tags)
}
