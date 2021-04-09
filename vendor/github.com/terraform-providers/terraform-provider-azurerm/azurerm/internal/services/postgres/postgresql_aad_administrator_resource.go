package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2020-01-01/postgresql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourcePostgreSQLAdministrator() *schema.Resource {
	return &schema.Resource{
		Create: resourcePostgreSQLAdministratorCreateUpdate,
		Read:   resourcePostgreSQLAdministratorRead,
		Update: resourcePostgreSQLAdministratorCreateUpdate,
		Delete: resourcePostgreSQLAdministratorDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AzureActiveDirectoryAdministratorID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"login": {
				Type:     schema.TypeString,
				Required: true,
			},

			"object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourcePostgreSQLAdministratorCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServerAdministratorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverName := d.Get("server_name").(string)
	resGroup := d.Get("resource_group_name").(string)
	login := d.Get("login").(string)
	objectId := uuid.FromStringOrNil(d.Get("object_id").(string))
	tenantId := uuid.FromStringOrNil(d.Get("tenant_id").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, serverName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing PostgreSQL AD Administrator (Resource Group %q, Server %q): %+v", resGroup, serverName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_postgresql_active_directory_administrator", *existing.ID)
		}
	}

	parameters := postgresql.ServerAdministratorResource{
		ServerAdministratorProperties: &postgresql.ServerAdministratorProperties{
			AdministratorType: utils.String("ActiveDirectory"),
			Login:             utils.String(login),
			Sid:               &objectId,
			TenantID:          &tenantId,
		},
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, serverName, parameters)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for PostgreSQL AD Administrator (Resource Group %q, Server %q): %+v", resGroup, serverName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for PostgreSQL AD Administrator (Resource Group %q, Server %q): %+v", resGroup, serverName, err)
	}

	resp, err := client.Get(ctx, resGroup, serverName)
	if err != nil {
		return fmt.Errorf("Error issuing get request for PostgreSQL AD Administrator (Resource Group %q, Server %q): %+v", resGroup, serverName, err)
	}

	d.SetId(*resp.ID)

	return nil
}

func resourcePostgreSQLAdministratorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServerAdministratorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AzureActiveDirectoryAdministratorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading PostgreSQL AD administrator %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading PostgreSQL AD administrator: %+v", err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("server_name", id.ServerName)

	if props := resp.ServerAdministratorProperties; props != nil {
		d.Set("login", props.Login)
		d.Set("object_id", props.Sid.String())
		d.Set("tenant_id", props.TenantID.String())
	}

	return nil
}

func resourcePostgreSQLAdministratorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServerAdministratorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AzureActiveDirectoryAdministratorID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		return fmt.Errorf("deleting AD Administrator (PostgreSQL Server %q / Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of AD Administrator (PostgreSQL Server %q / Resource Group %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	return nil
}
