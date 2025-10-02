package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"

	mgValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSubnetServiceEndpointStoragePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubnetServiceEndpointStoragePolicyCreateUpdate,
		Read:   resourceSubnetServiceEndpointStoragePolicyRead,
		Update: resourceSubnetServiceEndpointStoragePolicyCreateUpdate,
		Delete: resourceSubnetServiceEndpointStoragePolicyDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SubnetServiceEndpointStoragePolicyID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SubnetServiceEndpointStoragePolicyName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": location.Schema(),

			"definition": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.SubnetServiceEndpointStoragePolicyDefinitionName,
						},

						"service_resources": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.Any(
									azure.ValidateResourceID,
									mgValidate.ManagementGroupID,
								),
							},
						},

						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 140),
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceSubnetServiceEndpointStoragePolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewSubnetServiceEndpointStoragePolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.ServiceEndpointPolicyName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing %s: %+v", resourceId, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_subnet_service_endpoint_storage_policy", resourceId.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	param := network.ServiceEndpointPolicy{
		Location: &location,
		ServiceEndpointPolicyPropertiesFormat: &network.ServiceEndpointPolicyPropertiesFormat{
			ServiceEndpointPolicyDefinitions: expandServiceEndpointPolicyDefinitions(d.Get("definition").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.ServiceEndpointPolicyName, param)
	if err != nil {
		return fmt.Errorf("creating Subnet Service Endpoint Storage Policy %q (Resource Group %q): %+v", resourceId.ServiceEndpointPolicyName, resourceId.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Subnet Service Endpoint Storage Policy %q (Resource Group %q): %+v", resourceId.ServiceEndpointPolicyName, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())

	return resourceSubnetServiceEndpointStoragePolicyRead(d, meta)
}

func resourceSubnetServiceEndpointStoragePolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubnetServiceEndpointStoragePolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceEndpointPolicyName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Subnet Service Endpoint Storage Policy %q was not found in Resource Group %q - removing from state!", id.ServiceEndpointPolicyName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Subnet Service Endpoint Storage Policy %q (Resource Group %q): %+v", id.ServiceEndpointPolicyName, id.ResourceGroup, err)
	}

	d.Set("name", id.ServiceEndpointPolicyName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if prop := resp.ServiceEndpointPolicyPropertiesFormat; prop != nil {
		if err := d.Set("definition", flattenServiceEndpointPolicyDefinitions(prop.ServiceEndpointPolicyDefinitions)); err != nil {
			return fmt.Errorf("setting `definition`: %v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSubnetServiceEndpointStoragePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubnetServiceEndpointStoragePolicyID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.ServiceEndpointPolicyName); err != nil {
		return fmt.Errorf("deleting Subnet Service Endpoint Storage Policy %q (Resource Group %q): %+v", id.ServiceEndpointPolicyName, id.ResourceGroup, err)
	}

	return nil
}

func expandServiceEndpointPolicyDefinitions(input []interface{}) *[]network.ServiceEndpointPolicyDefinition {
	if len(input) == 0 {
		return nil
	}

	output := make([]network.ServiceEndpointPolicyDefinition, 0)
	for _, e := range input {
		e := e.(map[string]interface{})
		output = append(output, network.ServiceEndpointPolicyDefinition{
			Name: utils.String(e["name"].(string)),
			ServiceEndpointPolicyDefinitionPropertiesFormat: &network.ServiceEndpointPolicyDefinitionPropertiesFormat{
				Description:      utils.String(e["description"].(string)),
				Service:          utils.String("Microsoft.Storage"),
				ServiceResources: utils.ExpandStringSlice(e["service_resources"].(*schema.Set).List()),
			},
		})
	}

	return &output
}

func flattenServiceEndpointPolicyDefinitions(input *[]network.ServiceEndpointPolicyDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	for _, e := range *input {
		name := ""
		if e.Name != nil {
			name = *e.Name
		}

		var (
			description     = ""
			serviceResource = []interface{}{}
		)
		if b := e.ServiceEndpointPolicyDefinitionPropertiesFormat; b != nil {
			if b.Description != nil {
				description = *b.Description
			}
			serviceResource = utils.FlattenStringSlice(b.ServiceResources)
		}

		output = append(output, map[string]interface{}{
			"name":              name,
			"description":       description,
			"service_resources": serviceResource,
		})
	}

	return output
}
