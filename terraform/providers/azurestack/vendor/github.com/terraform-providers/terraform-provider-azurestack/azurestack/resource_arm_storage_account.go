package azurestack

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/storage/mgmt/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const blobStorageAccountDefaultAccessTier = "Hot"

func resourceArmStorageAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageAccountCreate,
		Read:   resourceArmStorageAccountRead,
		Update: resourceArmStorageAccountUpdate,
		Delete: resourceArmStorageAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageAccountName,
			},

			"resource_group_name": resourceGroupNameDiffSuppressSchema(),

			"location": locationSchema(),

			"account_kind": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.Storage),
					string(storage.BlobStorage),
				}, true),
				Default: string(storage.Storage),
			},

			"account_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Deprecated:       "This field has been split into `account_tier` and `account_replication_type`",
				ValidateFunc:     validateArmStorageAccountType,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			// Constants not in the 2017-03-09 profile
			"account_tier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
					"Premium",
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			// Constants not in 2017-03-09 profile
			"account_replication_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"LRS",
					"ZRS",
					"GRS",
					"RAGRS",
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			// Only valid for BlobStorage accounts, defaults to "Hot" in create function
			// Currently not supported in Azure Stack to create anything but "Storage" v1
			// This type does not allow you to select Hot or Cold storage options
			// "access_tier": {
			//	 Type:     schema.TypeString,
			//	 Optional: true,
			//	 Computed: true,
			//	 ValidateFunc: validation.StringInSlice([]string{
			//		 string(storage.Cool),
			//		 string(storage.Hot),
			//	 }, true),
			// },

			// Constants not in 2017-03-09 profile
			"account_encryption_source": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string("Microsoft.Storage"),
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.MicrosoftKeyvault),
					string(storage.MicrosoftStorage),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"custom_domain": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"use_subdomain": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"enable_blob_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"primary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_blob_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_blob_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_queue_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_queue_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_table_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_table_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// NOTE: The API does not appear to expose a secondary file endpoint
			"primary_file_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_blob_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_blob_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}

}

func resourceArmStorageAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).storageServiceClient

	resourceGroupName := d.Get("resource_group_name").(string)
	storageAccountName := d.Get("name").(string)
	accountKind := d.Get("account_kind").(string)

	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	accountTier := d.Get("account_tier").(string)
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)

	// Not supported by the profile in the same struct as the original, both of the
	// following commented lines will be read and set later on the correct
	// structs
	// storageAccountEncryptionSource := d.Get("account_encryption_source").(string)
	// enableBlobEncryption := d.Get("enable_blob_encryption").(bool)

	parameters := storage.AccountCreateParameters{
		Location: &location,
		Sku: &storage.Sku{
			Name: storage.SkuName(storageType),
		},
		Tags: *expandTags(tags),
		Kind: storage.Kind(accountKind),

		// If any paramers are specified withouth the right values this will fail
		AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{},
	}

	if _, ok := d.GetOk("custom_domain"); ok {
		parameters.CustomDomain = expandStorageAccountCustomDomain(d)
	}

	// BlobStorage does not support ZRS
	if accountKind == string(storage.BlobStorage) {
		if string(parameters.Sku.Name) == string(storage.StandardZRS) {
			return fmt.Errorf("A `account_replication_type` of `ZRS` isn't supported for Blob Storage accounts.")
		}
		accessTier, ok := d.GetOk("access_tier")
		if !ok {
			// default to "Hot"
			accessTier = blobStorageAccountDefaultAccessTier
		}

		parameters.AccountPropertiesCreateParameters.AccessTier = storage.AccessTier(accessTier.(string))

		enableBlobEncryption := d.Get("enable_blob_encryption").(bool)

		if enableBlobEncryption {
			// if the encryption is enabled, then set the arguments
			storageAccountEncryptionSource := d.Get("account_encryption_source").(string)
			parameters.AccountPropertiesCreateParameters.Encryption =
				&storage.Encryption{
					Services: &storage.EncryptionServices{
						Blob: &storage.EncryptionService{
							Enabled: utils.Bool(enableBlobEncryption),
						}},
					KeySource: storage.KeySource(storageAccountEncryptionSource),
				}
		}

	}

	// Create
	ctx := meta.(*ArmClient).StopContext
	future, err := client.Create(ctx, resourceGroupName, storageAccountName, parameters)
	if err != nil {
		return fmt.Errorf(
			"Error creating Azure Storage Account %q: %+v",
			storageAccountName, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf(
			"Error while waiting for Azure Storage Account %q: %+v",
			storageAccountName, err)
	}

	account, err := future.Result(client)
	if err != nil {
		return fmt.Errorf(
			"Error while fetching Azure Storage Account %q: %+v",
			storageAccountName, err)
	}

	log.Printf("[INFO] storage account %q ID: %q", storageAccountName, *account.ID)
	d.SetId(*account.ID)

	return resourceArmStorageAccountRead(d, meta)
}

// resourceArmStorageAccountUpdate is unusual in the ARM API where most resources have a combined
// and idempotent operation for CreateOrUpdate. In particular updating all of the parameters
// available requires a call to Update per parameter...
func resourceArmStorageAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).storageServiceClient
	// Create
	ctx := meta.(*ArmClient).StopContext
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	storageAccountName := id.Path["storageaccounts"]
	// https://github.com/terraform-providers/terraform-provider-azurestack/issues/98
	// it appears the casing of the Resource ID's changed in Azure Stack version 1905
	// as such we need to confirm both casings
	if storageAccountName == "" {
		storageAccountName = id.Path["storageAccounts"]
	}
	resourceGroupName := id.ResourceGroup

	accountTier := d.Get("account_tier").(string)
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)
	accountKind := d.Get("account_kind").(string)

	if accountKind == string(storage.BlobStorage) {
		if storageType == string(storage.StandardZRS) {
			return fmt.Errorf("A `account_replication_type` of `ZRS` isn't supported for Blob Storage accounts.")
		}
	}

	d.Partial(true)

	if d.HasChange("account_replication_type") {
		sku := storage.Sku{
			Name: storage.SkuName(storageType),
		}

		opts := storage.AccountUpdateParameters{
			Sku: &sku,
		}

		_, err := client.Update(ctx, resourceGroupName, storageAccountName, opts)
		if err != nil {
			return fmt.Errorf("Error updating Azure Storage Account type %q: %+v", storageAccountName, err)
		}

		d.SetPartial("account_replication_type")
	}

	// if d.HasChange("access_tier") {
	// 	accessTier := d.Get("access_tier").(string)

	// 	opts := storage.AccountUpdateParameters{
	// 		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
	// 			AccessTier: storage.AccessTier(accessTier),
	// 		},
	// 	}

	// 	_, err := client.Update(ctx, resourceGroupName, storageAccountName, opts)
	// 	if err != nil {
	// 		return fmt.Errorf("Error updating Azure Storage Account access_tier %q: %+v", storageAccountName, err)
	// 	}

	// 	d.SetPartial("access_tier")
	// }

	if d.HasChange("tags") {
		tags := d.Get("tags").(map[string]interface{})

		opts := storage.AccountUpdateParameters{
			Tags: *expandTags(tags),
		}

		_, err := client.Update(ctx, resourceGroupName, storageAccountName, opts)
		if err != nil {
			return fmt.Errorf("Error updating Azure Storage Account tags %q: %+v", storageAccountName, err)
		}

		d.SetPartial("tags")
	}

	if d.HasChange("enable_blob_encryption") {
		encryptionSource := d.Get("account_encryption_source").(string)

		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				Encryption: &storage.Encryption{
					Services:  &storage.EncryptionServices{},
					KeySource: storage.KeySource(encryptionSource),
				},
			},
		}

		if d.HasChange("enable_blob_encryption") {
			enableEncryption := d.Get("enable_blob_encryption").(bool)
			opts.Encryption.Services.Blob = &storage.EncryptionService{
				Enabled: utils.Bool(enableEncryption),
			}

			d.SetPartial("enable_blob_encryption")
		}

		_, err := client.Update(ctx, resourceGroupName, storageAccountName, opts)
		if err != nil {
			return fmt.Errorf("Error updating Azure Storage Account Encryption %q: %+v", storageAccountName, err)
		}
	}

	if d.HasChange("custom_domain") {
		customDomain := expandStorageAccountCustomDomain(d)
		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				CustomDomain: customDomain,
			},
		}

		_, err := client.Update(ctx, resourceGroupName, storageAccountName, opts)
		if err != nil {
			return fmt.Errorf("Error updating Azure Storage Account Custom Domain %q: %+v", storageAccountName, err)
		}
	}

	d.Partial(false)
	return nil
}

func resourceArmStorageAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).storageServiceClient
	endpointSuffix := meta.(*ArmClient).environment.StorageEndpointSuffix
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["storageaccounts"]
	// https://github.com/terraform-providers/terraform-provider-azurestack/issues/98
	// it appears the casing of the Resource ID's changed in Azure Stack version 1905
	// as such we need to confirm both casings
	if name == "" {
		name = id.Path["storageAccounts"]
	}
	resGroup := id.ResourceGroup

	resp, err := client.GetProperties(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading the state of AzurStack Storage Account %q: %+v", name, err)
	}
	// (resGroup, name)
	keys, err := client.ListKeys(ctx, resGroup, name)
	if err != nil {
		return err
	}

	accessKeys := *keys.Keys
	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureStackNormalizeLocation(*resp.Location))
	d.Set("account_kind", resp.Kind)

	if sku := resp.Sku; sku != nil {
		d.Set("account_type", sku.Name)
		d.Set("account_tier", sku.Tier)
		d.Set("account_replication_type", strings.Split(fmt.Sprintf("%v", sku.Name), "_")[1])
	}

	if props := resp.AccountProperties; props != nil {
		// Currently not supported on Azure Stack
		// d.Set("access_tier", props.AccessTier)

		if customDomain := props.CustomDomain; customDomain != nil {
			if err := d.Set("custom_domain", flattenStorageAccountCustomDomain(customDomain)); err != nil {
				return fmt.Errorf("Error flattening `custom_domain`: %+v", err)
			}
		}

		if encryption := props.Encryption; encryption != nil {
			if services := encryption.Services; services != nil {
				if blob := services.Blob; blob != nil {
					d.Set("enable_blob_encryption", blob.Enabled)
				}
			}
			d.Set("account_encryption_source", string(encryption.KeySource))
		}

		// Computed
		d.Set("primary_location", props.PrimaryLocation)
		d.Set("secondary_location", props.SecondaryLocation)

		if len(accessKeys) > 0 {
			pcs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *accessKeys[0].Value, endpointSuffix)
			d.Set("primary_connection_string", pcs)
		}

		if len(accessKeys) > 1 {
			scs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *accessKeys[1].Value, endpointSuffix)
			d.Set("secondary_connection_string", scs)
		}

		if endpoints := props.PrimaryEndpoints; endpoints != nil {
			d.Set("primary_blob_endpoint", endpoints.Blob)
			d.Set("primary_queue_endpoint", endpoints.Queue)
			d.Set("primary_table_endpoint", endpoints.Table)
			d.Set("primary_file_endpoint", endpoints.File)

			pscs := fmt.Sprintf("DefaultEndpointsProtocol=https;BlobEndpoint=%s;AccountName=%s;AccountKey=%s",
				*endpoints.Blob, *resp.Name, *accessKeys[0].Value)
			d.Set("primary_blob_connection_string", pscs)
		}

		if endpoints := props.SecondaryEndpoints; endpoints != nil {
			if blob := endpoints.Blob; blob != nil {
				d.Set("secondary_blob_endpoint", blob)
				sscs := fmt.Sprintf("DefaultEndpointsProtocol=https;BlobEndpoint=%s;AccountName=%s;AccountKey=%s",
					*blob, *resp.Name, *accessKeys[1].Value)
				d.Set("secondary_blob_connection_string", sscs)
			} else {
				d.Set("secondary_blob_endpoint", "")
				d.Set("secondary_blob_connection_string", "")
			}

			if endpoints.Queue != nil {
				d.Set("secondary_queue_endpoint", endpoints.Queue)
			} else {
				d.Set("secondary_queue_endpoint", "")
			}

			if endpoints.Table != nil {
				d.Set("secondary_table_endpoint", endpoints.Table)
			} else {
				d.Set("secondary_table_endpoint", "")
			}
		}
	}

	d.Set("primary_access_key", accessKeys[0].Value)
	d.Set("secondary_access_key", accessKeys[1].Value)

	flattenAndSetTags(d, &resp.Tags)

	return nil
}

func resourceArmStorageAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).storageServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	name := id.Path["storageaccounts"]
	// https://github.com/terraform-providers/terraform-provider-azurestack/issues/98
	// it appears the casing of the Resource ID's changed in Azure Stack version 1905
	// as such we need to confirm both casings
	if name == "" {
		name = id.Path["storageAccounts"]
	}
	resGroup := id.ResourceGroup

	_, err = client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error issuing AzureStack delete request for storage account %q: %+v", name, err)
	}

	return nil
}

func expandStorageAccountCustomDomain(d *schema.ResourceData) *storage.CustomDomain {
	domains := d.Get("custom_domain").([]interface{})
	if len(domains) == 0 {
		return &storage.CustomDomain{
			Name: utils.String(""),
		}
	}

	domain := domains[0].(map[string]interface{})
	name := domain["name"].(string)
	useSubDomain := domain["use_subdomain"].(bool)
	return &storage.CustomDomain{
		Name:             utils.String(name),
		UseSubDomainName: utils.Bool(useSubDomain),
	}
}

func flattenStorageAccountCustomDomain(input *storage.CustomDomain) []interface{} {
	domain := make(map[string]interface{})

	domain["name"] = *input.Name
	// use_subdomain isn't returned

	return []interface{}{domain}
}

func validateArmStorageAccountName(v interface{}, k string) (ws []string, es []error) {
	input := v.(string)

	if !regexp.MustCompile(`\A([a-z0-9]{3,24})\z`).MatchString(input) {
		es = append(es, fmt.Errorf("name can only consist of lowercase letters and numbers, and must be between 3 and 24 characters long"))
	}

	return
}

func validateArmStorageAccountType(v interface{}, k string) (ws []string, es []error) {
	validAccountTypes := []string{"standard_lrs", "standard_zrs",
		"standard_grs", "standard_ragrs", "premium_lrs"}

	input := strings.ToLower(v.(string))

	for _, valid := range validAccountTypes {
		if valid == input {
			return
		}
	}

	es = append(es, fmt.Errorf("Invalid storage account type %q", input))
	return
}
