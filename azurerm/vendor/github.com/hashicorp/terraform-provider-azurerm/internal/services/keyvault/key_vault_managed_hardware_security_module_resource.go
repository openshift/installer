package keyvault

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2021-10-01/keyvault" // nolint: staticcheck
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKeyVaultManagedHardwareSecurityModule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmKeyVaultManagedHardwareSecurityModuleCreate,
		Read:   resourceArmKeyVaultManagedHardwareSecurityModuleRead,
		Delete: resourceArmKeyVaultManagedHardwareSecurityModuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedHSMID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ManagedHardwareSecurityModuleName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(keyvault.ManagedHsmSkuNameStandardB1),
				}, false),
			},

			"admin_object_ids": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"purge_protection_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"soft_delete_retention_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      90,
				ValidateFunc: validation.IntBetween(7, 90),
			},

			"hsm_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				//Computed: true,
				Default:  true,
				ForceNew: true,
			},

			"network_acls": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"default_action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(keyvault.NetworkRuleActionAllow),
								string(keyvault.NetworkRuleActionDeny),
							}, false),
						},
						"bypass": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(keyvault.NetworkRuleBypassOptionsNone),
								string(keyvault.NetworkRuleBypassOptionsAzureServices),
							}, false),
						},
					},
				},
			},

			// https://github.com/Azure/azure-rest-api-specs/issues/13365
			"tags": tags.ForceNewSchema(),
		},
	}
}

func resourceArmKeyVaultManagedHardwareSecurityModuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagedHsmClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Println("[INFO] Preparing arguments for Key Vault Managed Hardware Security Module")

	id := parse.NewManagedHSMID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_key_vault_managed_hardware_security_module", id.ID())
	}

	tenantId := uuid.FromStringOrNil(d.Get("tenant_id").(string))
	hsm := keyvault.ManagedHsm{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: &keyvault.ManagedHsmProperties{
			TenantID:                  &tenantId,
			InitialAdminObjectIds:     utils.ExpandStringSlice(d.Get("admin_object_ids").(*pluginsdk.Set).List()),
			CreateMode:                keyvault.CreateModeDefault,
			EnableSoftDelete:          utils.Bool(true),
			SoftDeleteRetentionInDays: utils.Int32(int32(d.Get("soft_delete_retention_days").(int))),
			EnablePurgeProtection:     utils.Bool(d.Get("purge_protection_enabled").(bool)),
			PublicNetworkAccess:       keyvault.PublicNetworkAccessEnabled, // default enabled
			NetworkAcls:               expandMHSMNetworkAcls(d.Get("network_acls").([]interface{})),
		},
		Sku: &keyvault.ManagedHsmSku{
			Family: utils.String("B"),
			Name:   keyvault.ManagedHsmSkuName(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if !d.Get("public_network_access_enabled").(bool) {
		hsm.Properties.PublicNetworkAccess = keyvault.PublicNetworkAccessDisabled
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, hsm)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceArmKeyVaultManagedHardwareSecurityModuleRead(d, meta)
}

func resourceArmKeyVaultManagedHardwareSecurityModuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagedHsmClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedHSMID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[ERROR] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	skuName := ""
	if sku := resp.Sku; sku != nil {
		skuName = string(sku.Name)
	}
	d.Set("sku_name", skuName)

	if props := resp.Properties; props != nil {
		tenantId := ""
		if tid := props.TenantID; tid != nil {
			tenantId = tid.String()
		}
		d.Set("tenant_id", tenantId)
		d.Set("admin_object_ids", utils.FlattenStringSlice(props.InitialAdminObjectIds))
		d.Set("hsm_uri", props.HsmURI)
		d.Set("soft_delete_retention_days", props.SoftDeleteRetentionInDays)
		d.Set("purge_protection_enabled", props.EnablePurgeProtection)

		var publicAccess = true
		if props.PublicNetworkAccess == keyvault.PublicNetworkAccessDisabled {
			publicAccess = false
		}
		d.Set("public_network_access_enabled", publicAccess)

		d.Set("network_acls", flattenMHSMNetworkAcls(props.NetworkAcls))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmKeyVaultManagedHardwareSecurityModuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagedHsmClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedHSMID(d.Id())
	if err != nil {
		return err
	}

	// We need to grab the keyvault hsm to see if purge protection is enabled prior to deletion
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil || resp.Location == nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	// there is an API bug being tracked here: https://github.com/Azure/azure-rest-api-specs/issues/13365
	// taking the statusCode404 as the expected resource deletion result, instead of the error code which triggers retry
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf(
				"waiting for deletion of API Management Service %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	shouldPurge := meta.(*clients.Client).Features.KeyVault.PurgeSoftDeletedHSMsOnDestroy
	if shouldPurge && resp.Properties != nil && utils.NormaliseNilableBool(resp.Properties.EnablePurgeProtection) {
		log.Printf("[DEBUG] cannot purge %s because purge protection is enabled", id)
		return nil
	}

	purgeFuture, err := client.PurgeDeleted(ctx, id.Name, *resp.Location)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = purgeFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf(
				"waiting for purge of %s: %+v", id, err)
		}
	}

	return nil
}

func expandMHSMNetworkAcls(input []interface{}) *keyvault.MHSMNetworkRuleSet {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	res := &keyvault.MHSMNetworkRuleSet{
		Bypass:        keyvault.NetworkRuleBypassOptions(v["bypass"].(string)),
		DefaultAction: keyvault.NetworkRuleAction(v["default_action"].(string)),
	}

	return res
}

func flattenMHSMNetworkAcls(acl *keyvault.MHSMNetworkRuleSet) []interface{} {
	res := map[string]interface{}{
		"bypass":         string(keyvault.NetworkRuleBypassOptionsAzureServices),
		"default_action": string(keyvault.NetworkRuleActionAllow),
	}

	if acl != nil {
		res["bypass"] = string(acl.Bypass)
		res["default_action"] = string(acl.DefaultAction)
	}
	return []interface{}{res}
}
