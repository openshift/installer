package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutomationModule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationModuleCreateUpdate,
		Read:   resourceArmAutomationModuleRead,
		Update: resourceArmAutomationModuleCreateUpdate,
		Delete: resourceArmAutomationModuleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"automation_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"module_link": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uri": {
							Type:     schema.TypeString,
							Required: true,
						},

						"hash": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"algorithm": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmAutomationModuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationModuleClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Automation Module creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	accName := d.Get("automation_account_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, accName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Automation Module %q (Account %q / Resource Group %q): %s", name, accName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_automation_module", *existing.ID)
		}
	}

	contentLink := expandModuleLink(d)

	parameters := automation.ModuleCreateOrUpdateParameters{
		ModuleCreateOrUpdateProperties: &automation.ModuleCreateOrUpdateProperties{
			ContentLink: &contentLink,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resGroup, accName, name, parameters); err != nil {
		return err
	}

	// the API returns 'done' but it's not actually finished provisioning yet
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			string(automation.ModuleProvisioningStateActivitiesStored),
			string(automation.ModuleProvisioningStateConnectionTypeImported),
			string(automation.ModuleProvisioningStateContentDownloaded),
			string(automation.ModuleProvisioningStateContentRetrieved),
			string(automation.ModuleProvisioningStateContentStored),
			string(automation.ModuleProvisioningStateContentValidated),
			string(automation.ModuleProvisioningStateCreated),
			string(automation.ModuleProvisioningStateCreating),
			string(automation.ModuleProvisioningStateModuleDataStored),
			string(automation.ModuleProvisioningStateModuleImportRunbookComplete),
			string(automation.ModuleProvisioningStateRunningImportModuleRunbook),
			string(automation.ModuleProvisioningStateStartingImportModuleRunbook),
			string(automation.ModuleProvisioningStateUpdating),
		},
		Target: []string{
			string(automation.ModuleProvisioningStateSucceeded),
		},
		Timeout:    30 * time.Minute,
		MinTimeout: 30 * time.Second,
		Refresh: func() (interface{}, string, error) {

			resp, err2 := client.Get(ctx, resGroup, accName, name)
			if err2 != nil {
				return resp, "Error", fmt.Errorf("Error retrieving Module %q (Automation Account %q / Resource Group %q): %+v", name, accName, resGroup, err2)
			}

			if properties := resp.ModuleProperties; properties != nil {
				return resp, string(properties.ProvisioningState), nil
			}

			return resp, "Unknown", nil
		},
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Module %q (Automation Account %q / Resource Group %q) to finish provisioning: %+v", name, accName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, accName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Module %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAutomationModuleRead(d, meta)
}

func resourceArmAutomationModuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationModuleClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["modules"]

	resp, err := client.Get(ctx, resGroup, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Automation Module %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("automation_account_name", accName)

	return nil
}

func resourceArmAutomationModuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationModuleClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["modules"]

	resp, err := client.Delete(ctx, resGroup, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Automation Module %q: %+v", name, err)
	}

	return nil
}

func expandModuleLink(d *schema.ResourceData) automation.ContentLink {
	inputs := d.Get("module_link").([]interface{})
	input := inputs[0].(map[string]interface{})
	uri := input["uri"].(string)

	hashes := input["hash"].([]interface{})

	if len(hashes) > 0 {
		hash := hashes[0].(map[string]interface{})
		hashValue := hash["value"].(string)
		hashAlgorithm := hash["algorithm"].(string)

		return automation.ContentLink{
			URI: &uri,
			ContentHash: &automation.ContentHash{
				Algorithm: &hashAlgorithm,
				Value:     &hashValue,
			},
		}
	}

	return automation.ContentLink{
		URI: &uri,
	}
}
