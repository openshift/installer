package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2020-01-01/postgresql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/validate"
	resourcesClient "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/client"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourcePostgreSQLServerKey() *schema.Resource {
	return &schema.Resource{
		Create: resourcePostgreSQLServerKeyCreateUpdate,
		Read:   resourcePostgreSQLServerKeyRead,
		Update: resourcePostgreSQLServerKeyCreateUpdate,
		Delete: resourcePostgreSQLServerKeyDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ServerKeyID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerID,
			},

			"key_vault_key_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.NestedItemId,
			},
		},
	}
}

func getPostgreSQLServerKeyName(ctx context.Context, keyVaultsClient *client.Client, resourcesClient *resourcesClient.Client, keyVaultKeyURI string) (*string, error) {
	keyVaultKeyID, err := keyVaultParse.ParseNestedItemID(keyVaultKeyURI)
	if err != nil {
		return nil, err
	}
	keyVaultIDRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, keyVaultKeyID.KeyVaultBaseUrl)
	if err != nil {
		return nil, err
	}
	// function azure.GetKeyVaultIDFromBaseUrl returns nil with nil error when it does not find the keyvault by the keyvault URL
	if keyVaultIDRaw == nil {
		return nil, fmt.Errorf("cannot get the keyvault ID from keyvault URL %q", keyVaultKeyID.KeyVaultBaseUrl)
	}
	keyVaultID, err := keyVaultParse.VaultID(*keyVaultIDRaw)
	if err != nil {
		return nil, err
	}
	return utils.String(fmt.Sprintf("%s_%s_%s", keyVaultID.Name, keyVaultKeyID.Name, keyVaultKeyID.Version)), nil
}

func resourcePostgreSQLServerKeyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	keysClient := meta.(*clients.Client).Postgres.ServerKeysClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverID, err := parse.ServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}
	keyVaultKeyURI := d.Get("key_vault_key_id").(string)
	name, err := getPostgreSQLServerKeyName(ctx, keyVaultsClient, resourcesClient, keyVaultKeyURI)
	if err != nil {
		return fmt.Errorf("cannot compose name for PostgreSQL Server Key (Resource Group %q / Server %q): %+v", serverID.ResourceGroup, serverID.Name, err)
	}

	locks.ByName(serverID.Name, postgreSQLServerResourceName)
	defer locks.UnlockByName(serverID.Name, postgreSQLServerResourceName)

	if d.IsNewResource() {
		// This resource is a singleton, but its name can be anything.
		// If you create a new key with different name with the old key, the service will not give you any warning but directly replace the old key with the new key.
		// Therefore sometimes you cannot get the old key using the GET API since you may not know the name of the old key
		resp, err := keysClient.List(ctx, serverID.ResourceGroup, serverID.Name)
		if err != nil {
			return fmt.Errorf("listing existing PostgreSQL Server Keys in Resource Group %q / Server %q: %+v", serverID.ResourceGroup, serverID.Name, err)
		}
		keys := resp.Values()
		if len(keys) > 1 {
			return fmt.Errorf("expecting at most one PostgreSQL Server Key, but got %q", len(keys))
		}
		if len(keys) == 1 && keys[0].ID != nil && *keys[0].ID != "" {
			return tf.ImportAsExistsError("azurerm_postgresql_server_key", *keys[0].ID)
		}
	}

	param := postgresql.ServerKey{
		ServerKeyProperties: &postgresql.ServerKeyProperties{
			ServerKeyType: utils.String("AzureKeyVault"),
			URI:           utils.String(d.Get("key_vault_key_id").(string)),
		},
	}

	future, err := keysClient.CreateOrUpdate(ctx, serverID.Name, *name, param, serverID.ResourceGroup)
	if err != nil {
		return fmt.Errorf("creating/updating PostgreSQL Server Key %q (Resource Group %q / Server %q): %+v", *name, serverID.ResourceGroup, serverID.Name, err)
	}
	if err := future.WaitForCompletionRef(ctx, keysClient.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of PostgreSQL Server Key %q (Resource Group %q / Server %q): %+v", *name, serverID.ResourceGroup, serverID.Name, err)
	}

	resp, err := keysClient.Get(ctx, serverID.ResourceGroup, serverID.Name, *name)
	if err != nil {
		return fmt.Errorf("retrieving PostgreSQL Server Key %q (Resource Group %q / Server %q): %+v", *name, serverID.ResourceGroup, serverID.Name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("nil or empty ID returned for PostgreSQL Server Key %q (Resource Group %q / Server %q): %+v", *name, serverID.ResourceGroup, serverID.Name, err)
	}

	d.SetId(*resp.ID)
	return resourcePostgreSQLServerKeyRead(d, meta)
}

func resourcePostgreSQLServerKeyRead(d *schema.ResourceData, meta interface{}) error {
	serversClient := meta.(*clients.Client).Postgres.ServersClient
	keysClient := meta.(*clients.Client).Postgres.ServerKeysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerKeyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := keysClient.Get(ctx, id.ResourceGroup, id.ServerName, id.KeyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] PostgreSQL Server Key %q was not found (Resource Group %q / Server %q)", id.KeyName, id.ResourceGroup, id.ServerName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving PostgreSQL Server Key %q (Resource Group %q / Server %q): %+v", id.KeyName, id.ResourceGroup, id.ServerName, err)
	}

	respServer, err := serversClient.Get(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		return fmt.Errorf("cannot get MySQL Server ID: %+v", err)
	}

	d.Set("server_id", respServer.ID)
	if props := resp.ServerKeyProperties; props != nil {
		d.Set("key_vault_key_id", props.URI)
	}

	return nil
}

func resourcePostgreSQLServerKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServerKeysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerKeyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ServerName, postgreSQLServerResourceName)
	defer locks.UnlockByName(id.ServerName, postgreSQLServerResourceName)

	future, err := client.Delete(ctx, id.ServerName, id.KeyName, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("deleting PostgreSQL Server Key %q (Resource Group %q / Server %q): %+v", id.KeyName, id.ResourceGroup, id.ServerName, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of PostgreSQL Server Key %q (Resource Group %q / Server %q): %+v", id.KeyName, id.ResourceGroup, id.ServerName, err)
	}

	return nil
}
