package web

import (
	"context"
	"fmt"
	"log"
	"time"

	relayMngt "github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	azValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	relayParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/relay/parse"
	relayValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/relay/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceAppServiceHybridConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppServiceHybridConnectionCreateUpdate,
		Read:   resourceAppServiceHybridConnectionRead,
		Update: resourceAppServiceHybridConnectionCreateUpdate,
		Delete: resourceAppServiceHybridConnectionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.HybridConnectionID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"app_service_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AppServiceName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"relay_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: relayValidate.HybridConnectionID,
			},

			"hostname": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: azValidate.PortNumberOrZero,
			},

			"send_key_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "RootManageSharedAccessKey",
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"namespace_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"relay_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_bus_namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_bus_suffix": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"send_key_value": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func resourceAppServiceHybridConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("app_service_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	relayArmURI := d.Get("relay_id").(string)
	relayId, err := relayParse.HybridConnectionID(relayArmURI)
	if err != nil {
		return fmt.Errorf("Error parsing relay ID %q: %s", relayArmURI, err)
	}
	namespaceName := relayId.NamespaceName
	relayName := relayId.Name

	if d.IsNewResource() {
		existing, err := client.GetHybridConnection(ctx, resourceGroup, name, namespaceName, relayName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing App Service Hybrid Connection %q (Resource Group %q, Namespace %q, Relay Name %q): %s", name, resourceGroup, namespaceName, relayName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_hybrid_connection", *existing.ID)
		}
	}

	port := int32(d.Get("port").(int))

	connectionEnvelope := web.HybridConnection{
		HybridConnectionProperties: &web.HybridConnectionProperties{
			RelayArmURI:  &relayArmURI,
			Hostname:     utils.String(d.Get("hostname").(string)),
			Port:         &port,
			SendKeyName:  utils.String(d.Get("send_key_name").(string)),
			SendKeyValue: utils.String(""), // The service creates this no matter what is sent, but the API requires the field to be set
		},
	}

	hybridConnection, err := client.CreateOrUpdateHybridConnection(ctx, resourceGroup, name, namespaceName, relayName, connectionEnvelope)
	if err != nil {
		return fmt.Errorf("failed creating App Service Hybrid Connection %q (resource group %q): %s", name, resourceGroup, err)
	}

	if hybridConnection.ID == nil && *hybridConnection.ID == "" {
		return fmt.Errorf("failed to read ID for Hybrid Connection %q", name)
	}
	d.SetId(*hybridConnection.ID)

	return resourceAppServiceHybridConnectionRead(d, meta)
}

func resourceAppServiceHybridConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HybridConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Hybrid Connection %q for App Service %q (Resource Group %q) was not found - removing from state", id.HybridConnectionNamespaceName, id.SiteName, id.RelayName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving App Service Hybrid Connection %q in Namespace %q, Resource Group %q: %s", id.SiteName, id.HybridConnectionNamespaceName, id.ResourceGroup, err)
	}
	d.Set("app_service_name", id.SiteName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("namespace_name", id.HybridConnectionNamespaceName)
	d.Set("relay_name", id.RelayName)

	if props := resp.HybridConnectionProperties; props != nil {
		d.Set("port", resp.Port)
		d.Set("service_bus_namespace", resp.ServiceBusNamespace)
		d.Set("send_key_name", resp.SendKeyName)
		d.Set("service_bus_suffix", resp.ServiceBusSuffix)
		d.Set("relay_id", resp.RelayArmURI)
		d.Set("hostname", resp.Hostname)
	}

	// key values are not returned in the response, so we get the primary key from the relay namespace ListKeys func
	if resp.ServiceBusNamespace != nil && resp.SendKeyName != nil {
		relayNSClient := meta.(*clients.Client).Relay.NamespacesClient
		relayNamespaceRG, err := findRelayNamespace(relayNSClient, ctx, *resp.ServiceBusNamespace)
		if err != nil {
			return err
		}
		accessKeys, err := relayNSClient.ListKeys(ctx, relayNamespaceRG, *resp.ServiceBusNamespace, *resp.SendKeyName)
		if err != nil {
			return fmt.Errorf("unable to List Access Keys for Namespace %q (Resource Group %q): %+v", *resp.ServiceBusNamespace, id.ResourceGroup, err)
		} else {
			d.Set("send_key_value", accessKeys.PrimaryKey)
		}
	}

	return nil
}

func resourceAppServiceHybridConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HybridConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DeleteHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting App Service Hybrid Connection %q (Resource Group %q, Relay %q): %+v", id.SiteName, id.ResourceGroup, id.RelayName, err)
		}
	}

	return nil
}

func findRelayNamespace(client *relayMngt.NamespacesClient, ctx context.Context, name string) (string, error) {
	relayNSIterator, err := client.ListComplete(ctx)
	if err != nil {
		return "", fmt.Errorf("listing Relay Namespaces: %+v", err)
	}

	var found *relayMngt.Namespace
	for relayNSIterator.NotDone() {
		namespace := relayNSIterator.Value()
		if namespace.Name != nil && *namespace.Name == name {
			found = &namespace
			break
		}
		if err := relayNSIterator.NextWithContext(ctx); err != nil {
			return "", fmt.Errorf("listing Relay Namespaces: %+v", err)
		}
	}

	if found == nil || found.ID == nil {
		return "", fmt.Errorf("could not find Relay Namespace with name: %q", name)
	}

	id, err := relayParse.NamespaceID(*found.ID)
	if err != nil {
		return "", fmt.Errorf("relay Namespace id not valid: %+v", err)
	}
	return id.ResourceGroup, nil
}
