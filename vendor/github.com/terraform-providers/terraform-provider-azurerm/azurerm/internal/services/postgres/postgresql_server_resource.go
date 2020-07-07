package postgres

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPostgreSQLServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPostgreSQLServerCreate,
		Read:   resourceArmPostgreSQLServerRead,
		Update: resourceArmPostgreSQLServerUpdate,
		Delete: resourceArmPostgreSQLServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:          schema.TypeString,
				Optional:      true, // required in 2.0
				Computed:      true, // remove in 2.0
				ConflictsWith: []string{"sku"},
				ValidateFunc: validation.StringInSlice([]string{
					"B_Gen4_1",
					"B_Gen4_2",
					"B_Gen5_1",
					"B_Gen5_2",
					"GP_Gen4_2",
					"GP_Gen4_4",
					"GP_Gen4_8",
					"GP_Gen4_16",
					"GP_Gen4_32",
					"GP_Gen5_2",
					"GP_Gen5_4",
					"GP_Gen5_8",
					"GP_Gen5_16",
					"GP_Gen5_32",
					"GP_Gen5_64",
					"MO_Gen5_2",
					"MO_Gen5_4",
					"MO_Gen5_8",
					"MO_Gen5_16",
					"MO_Gen5_32",
				}, false),
			},

			// remove in 2.0
			"sku": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"sku_name"},
				Deprecated:    "This property has been deprecated in favour of the 'sku_name' property and will be removed in version 2.0 of the provider",
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"B_Gen4_1",
								"B_Gen4_2",
								"B_Gen5_1",
								"B_Gen5_2",
								"GP_Gen4_2",
								"GP_Gen4_4",
								"GP_Gen4_8",
								"GP_Gen4_16",
								"GP_Gen4_32",
								"GP_Gen5_2",
								"GP_Gen5_4",
								"GP_Gen5_8",
								"GP_Gen5_16",
								"GP_Gen5_32",
								"GP_Gen5_64",
								"MO_Gen5_2",
								"MO_Gen5_4",
								"MO_Gen5_8",
								"MO_Gen5_16",
								"MO_Gen5_32",
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"capacity": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validate.IntInSlice([]int{
								1,
								2,
								4,
								8,
								16,
								32,
								64,
							}),
						},

						"tier": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(postgresql.Basic),
								string(postgresql.GeneralPurpose),
								string(postgresql.MemoryOptimized),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"family": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Gen4",
								"Gen5",
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},

			"administrator_login": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"administrator_login_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(postgresql.NineFullStopFive),
					string(postgresql.NineFullStopSix),
					string(postgresql.OneOne),
					string(postgresql.OneZero),
					string(postgresql.OneZeroFullStopZero),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"storage_profile": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_mb": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.IntBetweenAndDivisibleBy(5120, 4194304, 1024),
						},

						"backup_retention_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(7, 35),
						},

<<<<<<< HEAD:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/resource_arm_postgresql_server.go
						"geo_redundant_backup": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Enabled",
								"Disabled",
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"auto_grow": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(postgresql.StorageAutogrowEnabled),
							ValidateFunc: validation.StringInSlice([]string{
								string(postgresql.StorageAutogrowEnabled),
								string(postgresql.StorageAutogrowDisabled),
							}, false),
						},
					},
				},
=======
			"ssl_enforcement_enabled": {
				Type:         schema.TypeBool,
				Optional:     true, // required in 3.0
				ExactlyOneOf: []string{"ssl_enforcement", "ssl_enforcement_enabled"},
>>>>>>> 5aa20dd53... vendor: bump terraform-provider-azure to version v2.17.0:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/postgresql_server_resource.go
			},

			"ssl_enforcement": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(postgresql.SslEnforcementEnumDisabled),
					string(postgresql.SslEnforcementEnumEnabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"threat_detection_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"disabled_alerts": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"Sql_Injection",
									"Sql_Injection_Vulnerability",
									"Access_Anomaly",
									"Data_Exfiltration",
									"Unsafe_Action",
								}, false),
							},
						},

						"email_account_admins": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"email_addresses": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								// todo email validation in code
							},
							Set: schema.HashString,
						},

						"retention_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"storage_account_access_key": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"storage_endpoint": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},

<<<<<<< HEAD:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/resource_arm_postgresql_server.go
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			tier, _ := diff.GetOk("sku.0.tier")
			storageMB, _ := diff.GetOk("storage_profile.0.storage_mb")

			if strings.ToLower(tier.(string)) == "basic" && storageMB.(int) > 1048576 {
				return fmt.Errorf("basic pricing tier only supports upto 1,048,576 MB (1TB) of storage")
			}

			return nil
		},
=======
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("sku_name", func(old, new, meta interface{}) bool {
				oldTier := strings.Split(old.(string), "_")
				newTier := strings.Split(new.(string), "_")
				// If the sku tier was not changed, we don't need fornew
				if oldTier[0] == newTier[0] {
					return false
				}
				// Basic tier could not be changed to other tiers
				if oldTier[0] == "B" || newTier[0] == "B" {
					return true
				}
				return false
			}),
		),
>>>>>>> 5aa20dd53... vendor: bump terraform-provider-azure to version v2.17.0:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/postgresql_server_resource.go
	}
}

func resourceArmPostgreSQLServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	securityClient := meta.(*clients.Client).Postgres.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Server creation.")

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_postgresql_server", *existing.ID)
		}
	}

<<<<<<< HEAD:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/resource_arm_postgresql_server.go
	var sku *postgresql.Sku
	if b, ok := d.GetOk("sku_name"); ok {
		var err error
		sku, err = expandServerSkuName(b.(string))
		if err != nil {
			return fmt.Errorf("error expanding sku_name for PostgreSQL Server %s (Resource Group %q): %v", name, resourceGroup, err)
=======
	mode := postgresql.CreateMode(d.Get("create_mode").(string))
	tlsMin := postgresql.MinimalTLSVersionEnum(d.Get("ssl_minimal_tls_version_enforced").(string))
	source := d.Get("creation_source_server_id").(string)
	version := postgresql.ServerVersion(d.Get("version").(string))

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name` for PostgreSQL Server %s (Resource Group %q): %v", name, resourceGroup, err)
	}

	infraEncrypt := postgresql.InfrastructureEncryptionEnabled
	if v := d.Get("infrastructure_encryption_enabled"); !v.(bool) {
		infraEncrypt = postgresql.InfrastructureEncryptionDisabled
	}

	publicAccess := postgresql.PublicNetworkAccessEnumEnabled
	if v := d.Get("public_network_access_enabled"); !v.(bool) {
		publicAccess = postgresql.PublicNetworkAccessEnumDisabled
	}

	ssl := postgresql.SslEnforcementEnumEnabled
	if v := d.Get("ssl_enforcement_enabled"); !v.(bool) {
		ssl = postgresql.SslEnforcementEnumDisabled
	}

	storage := expandPostgreSQLStorageProfile(d)

	var props postgresql.BasicServerPropertiesForCreate
	switch mode {
	case postgresql.CreateModeDefault:
		admin := d.Get("administrator_login").(string)
		pass := d.Get("administrator_login_password").(string)

		if admin == "" {
			return fmt.Errorf("`administrator_login` must not be empty when `create_mode` is `default`")
		}
		if pass == "" {
			return fmt.Errorf("`administrator_login_password` must not be empty when `create_mode` is `default`")
		}

		if _, ok := d.GetOk("restore_point_in_time"); ok {
			return fmt.Errorf("`restore_point_in_time` cannot be set when `create_mode` is `default`")
		}

		// check admin
		props = &postgresql.ServerPropertiesForDefaultCreate{
			AdministratorLogin:         &admin,
			AdministratorLoginPassword: &pass,
			CreateMode:                 mode,
			InfrastructureEncryption:   infraEncrypt,
			PublicNetworkAccess:        publicAccess,
			MinimalTLSVersion:          tlsMin,
			SslEnforcement:             ssl,
			StorageProfile:             storage,
			Version:                    version,
		}
	case postgresql.CreateModePointInTimeRestore:
		v, ok := d.GetOk("restore_point_in_time")
		if !ok || v.(string) == "" {
			return fmt.Errorf("restore_point_in_time must be set when create_mode is PointInTimeRestore")
		}
		time, _ := time.Parse(time.RFC3339, v.(string)) // should be validated by the schema

		props = &postgresql.ServerPropertiesForRestore{
			CreateMode:     mode,
			SourceServerID: &source,
			RestorePointInTime: &date.Time{
				Time: time,
			},
			InfrastructureEncryption: infraEncrypt,
			PublicNetworkAccess:      publicAccess,
			MinimalTLSVersion:        tlsMin,
			SslEnforcement:           ssl,
			StorageProfile:           storage,
			Version:                  version,
		}
	case postgresql.CreateModeGeoRestore:
		props = &postgresql.ServerPropertiesForGeoRestore{
			CreateMode:               mode,
			SourceServerID:           &source,
			InfrastructureEncryption: infraEncrypt,
			PublicNetworkAccess:      publicAccess,
			MinimalTLSVersion:        tlsMin,
			SslEnforcement:           ssl,
			StorageProfile:           storage,
			Version:                  version,
		}
	case postgresql.CreateModeReplica:
		props = &postgresql.ServerPropertiesForReplica{
			CreateMode:               mode,
			SourceServerID:           &source,
			InfrastructureEncryption: infraEncrypt,
			PublicNetworkAccess:      publicAccess,
			MinimalTLSVersion:        tlsMin,
			SslEnforcement:           ssl,
			Version:                  version,
>>>>>>> 5aa20dd53... vendor: bump terraform-provider-azure to version v2.17.0:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/postgresql_server_resource.go
		}
	} else if _, ok := d.GetOk("sku"); ok {
		sku = expandAzureRmPostgreSQLServerSku(d)
	} else {
		return fmt.Errorf("One of `sku` or `sku_name` must be set for PostgreSQL Server %q (Resource Group %q)", name, resourceGroup)
	}

	properties := postgresql.ServerForCreate{
		Location: &location,
		Properties: &postgresql.ServerPropertiesForDefaultCreate{
			AdministratorLogin:         utils.String(d.Get("administrator_login").(string)),
			AdministratorLoginPassword: utils.String(d.Get("administrator_login_password").(string)),
			Version:                    postgresql.ServerVersion(d.Get("version").(string)),
			SslEnforcement:             postgresql.SslEnforcementEnum(d.Get("ssl_enforcement").(string)),
			StorageProfile:             expandAzureRmPostgreSQLStorageProfile(d),
			CreateMode:                 postgresql.CreateMode("Default"),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Create(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("Error creating PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read PostgreSQL Server %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	if v, ok := d.GetOk("threat_detection_policy"); ok {
		alert := expandSecurityAlertPolicy(v)
		if alert != nil {
			future, err := securityClient.CreateOrUpdate(ctx, resourceGroup, name, *alert)
			if err != nil {
				return fmt.Errorf("error updataing postgres server security alert policy: %v", err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("error waiting for creation/update of postgrest server security alert policy (server %q, resource group %q): %+v", name, resourceGroup, err)
			}
		}
	}

	return resourceArmPostgreSQLServerRead(d, meta)
}

func resourceArmPostgreSQLServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	securityClient := meta.(*clients.Client).Postgres.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// TODO: support for Delta updates

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Server update.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

<<<<<<< HEAD:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/resource_arm_postgresql_server.go
	var sku *postgresql.Sku
	if b, ok := d.GetOk("sku_name"); ok {
		var err error
		sku, err = expandServerSkuName(b.(string))
		if err != nil {
			return fmt.Errorf("error expanding sku_name for PostgreSQL Server %q (Resource Group %q): %v", name, resourceGroup, err)
		}
	} else if _, ok := d.GetOk("sku"); ok {
		sku = expandAzureRmPostgreSQLServerSku(d)
	} else {
		return fmt.Errorf("One of `sku` or `sku_name` must be set for PostgreSQL Server %q (Resource Group %q)", name, resourceGroup)
=======
	ssl := postgresql.SslEnforcementEnumEnabled
	if v := d.Get("ssl_enforcement_enabled"); !v.(bool) {
		ssl = postgresql.SslEnforcementEnumDisabled
>>>>>>> 5aa20dd53... vendor: bump terraform-provider-azure to version v2.17.0:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/postgresql_server_resource.go
	}

	properties := postgresql.ServerUpdateParameters{
		ServerUpdateParametersProperties: &postgresql.ServerUpdateParametersProperties{
			AdministratorLoginPassword: utils.String(d.Get("administrator_login_password").(string)),
<<<<<<< HEAD:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/resource_arm_postgresql_server.go
=======
			PublicNetworkAccess:        publicAccess,
			SslEnforcement:             ssl,
			StorageProfile:             expandPostgreSQLStorageProfile(d),
>>>>>>> 5aa20dd53... vendor: bump terraform-provider-azure to version v2.17.0:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/postgresql_server_resource.go
			Version:                    postgresql.ServerVersion(d.Get("version").(string)),
			SslEnforcement:             postgresql.SslEnforcementEnum(d.Get("ssl_enforcement").(string)),
			StorageProfile:             expandAzureRmPostgreSQLStorageProfile(d),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Update(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("Error updating PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

<<<<<<< HEAD:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/resource_arm_postgresql_server.go
	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read PostgreSQL Server %s (resource group %s) ID", name, resourceGroup)
	}
=======
	if v, ok := d.GetOk("threat_detection_policy"); ok {
		alert := expandSecurityAlertPolicy(v)
		if alert != nil {
			future, err := securityClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, *alert)
			if err != nil {
				return fmt.Errorf("error updataing mssql server security alert policy: %v", err)
			}
>>>>>>> 5aa20dd53... vendor: bump terraform-provider-azure to version v2.17.0:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/postgresql_server_resource.go

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("error waiting for creation/update of postgrest server security alert policy (server %q, resource group %q): %+v", id.Name, id.ResourceGroup, err)
			}
		}
	}

	return resourceArmPostgreSQLServerRead(d, meta)
}

func resourceArmPostgreSQLServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	securityClient := meta.(*clients.Client).Postgres.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["servers"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] PostgreSQL Server %q was not found (resource group %q)", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
<<<<<<< HEAD:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/resource_arm_postgresql_server.go
	}

	d.Set("administrator_login", resp.AdministratorLogin)
	d.Set("version", string(resp.Version))
	d.Set("ssl_enforcement", string(resp.SslEnforcement))
=======
	}

	tier := postgresql.Basic
	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
		tier = sku.Tier
	}

	if props := resp.ServerProperties; props != nil {
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("ssl_enforcement", string(props.SslEnforcement))
		d.Set("ssl_minimal_tls_version_enforced", props.MinimalTLSVersion)
		d.Set("version", string(props.Version))
>>>>>>> 5aa20dd53... vendor: bump terraform-provider-azure to version v2.17.0:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/postgresql_server_resource.go

	if err := d.Set("sku", flattenPostgreSQLServerSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}
	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if err := d.Set("storage_profile", flattenPostgreSQLStorageProfile(resp.StorageProfile)); err != nil {
		return fmt.Errorf("Error setting `storage_profile`: %+v", err)
	}

	// Computed
	d.Set("fqdn", resp.FullyQualifiedDomainName)

<<<<<<< HEAD:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/resource_arm_postgresql_server.go
=======
		// Computed
		d.Set("fqdn", props.FullyQualifiedDomainName)
	}

	// the basic does not support threat detection policies
	if tier == postgresql.GeneralPurpose || tier == postgresql.MemoryOptimized {
		secResp, err := securityClient.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil && !utils.ResponseWasNotFound(secResp.Response) {
			return fmt.Errorf("error making read request to postgres server security alert policy: %+v", err)
		}

		if !utils.ResponseWasNotFound(secResp.Response) {
			block := flattenSecurityAlertPolicy(secResp.SecurityAlertPolicyProperties, d.Get("threat_detection_policy.0.storage_account_access_key").(string))
			if err := d.Set("threat_detection_policy", block); err != nil {
				return fmt.Errorf("setting `threat_detection_policy`: %+v", err)
			}
		}
	}

>>>>>>> 5aa20dd53... vendor: bump terraform-provider-azure to version v2.17.0:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/postgresql_server_resource.go
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmPostgreSQLServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["servers"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error waiting for deletion of PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func expandServerSkuName(skuName string) (*postgresql.Sku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 3 {
		return nil, fmt.Errorf("sku_name (%s) has the wrong number of parts (%d) after splitting on _", skuName, len(parts))
	}

	var tier postgresql.SkuTier
	switch parts[0] {
	case "B":
		tier = postgresql.Basic
	case "GP":
		tier = postgresql.GeneralPurpose
	case "MO":
		tier = postgresql.MemoryOptimized
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, parts[0])
	}

	capacity, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("cannot convert skuname %s capcity %s to int", skuName, parts[2])
	}

	return &postgresql.Sku{
		Name:     utils.String(skuName),
		Tier:     tier,
		Capacity: utils.Int32(int32(capacity)),
		Family:   utils.String(parts[1]),
	}, nil
}

<<<<<<< HEAD:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/resource_arm_postgresql_server.go
func expandAzureRmPostgreSQLServerSku(d *schema.ResourceData) *postgresql.Sku {
	skus := d.Get("sku").([]interface{})
	sku := skus[0].(map[string]interface{})
=======
func expandPostgreSQLStorageProfile(d *schema.ResourceData) *postgresql.StorageProfile {
	storage := postgresql.StorageProfile{}
	if v, ok := d.GetOk("storage_profile"); ok {
		storageprofile := v.([]interface{})[0].(map[string]interface{})
>>>>>>> 5aa20dd53... vendor: bump terraform-provider-azure to version v2.17.0:vendor/github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/postgresql_server_resource.go

	name := sku["name"].(string)
	capacity := sku["capacity"].(int)
	tier := sku["tier"].(string)
	family := sku["family"].(string)

	return &postgresql.Sku{
		Name:     utils.String(name),
		Tier:     postgresql.SkuTier(tier),
		Capacity: utils.Int32(int32(capacity)),
		Family:   utils.String(family),
	}
}

func expandAzureRmPostgreSQLStorageProfile(d *schema.ResourceData) *postgresql.StorageProfile {
	storageprofiles := d.Get("storage_profile").([]interface{})
	storageprofile := storageprofiles[0].(map[string]interface{})

	backupRetentionDays := storageprofile["backup_retention_days"].(int)
	geoRedundantBackup := storageprofile["geo_redundant_backup"].(string)
	storageMB := storageprofile["storage_mb"].(int)
	autoGrow := storageprofile["auto_grow"].(string)

	return &postgresql.StorageProfile{
		BackupRetentionDays: utils.Int32(int32(backupRetentionDays)),
		GeoRedundantBackup:  postgresql.GeoRedundantBackup(geoRedundantBackup),
		StorageMB:           utils.Int32(int32(storageMB)),
		StorageAutogrow:     postgresql.StorageAutogrow(autoGrow),
	}
}

func flattenPostgreSQLServerSku(resp *postgresql.Sku) []interface{} {
	values := map[string]interface{}{}

	if name := resp.Name; name != nil {
		values["name"] = *name
	}

	if capacity := resp.Capacity; capacity != nil {
		values["capacity"] = *capacity
	}

	values["tier"] = string(resp.Tier)

	if family := resp.Family; family != nil {
		values["family"] = *family
	}

	return []interface{}{values}
}

func flattenPostgreSQLStorageProfile(resp *postgresql.StorageProfile) []interface{} {
	values := map[string]interface{}{}

	if storageMB := resp.StorageMB; storageMB != nil {
		values["storage_mb"] = *storageMB
	}

	values["auto_grow"] = string(resp.StorageAutogrow)

	if backupRetentionDays := resp.BackupRetentionDays; backupRetentionDays != nil {
		values["backup_retention_days"] = *backupRetentionDays
	}

	values["geo_redundant_backup"] = string(resp.GeoRedundantBackup)

	return []interface{}{values}
}

func expandSecurityAlertPolicy(i interface{}) *postgresql.ServerSecurityAlertPolicy {
	slice := i.([]interface{})
	if len(slice) == 0 {
		return nil
	}

	block := slice[0].(map[string]interface{})

	state := postgresql.ServerSecurityAlertPolicyStateEnabled
	if !block["enabled"].(bool) {
		state = postgresql.ServerSecurityAlertPolicyStateDisabled
	}

	props := &postgresql.SecurityAlertPolicyProperties{
		State: state,
	}

	if v, ok := block["disabled_alerts"]; ok {
		props.DisabledAlerts = utils.ExpandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := block["email_addresses"]; ok {
		props.EmailAddresses = utils.ExpandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := block["email_account_admins"]; ok {
		props.EmailAccountAdmins = utils.Bool(v.(bool))
	}

	if v, ok := block["retention_days"]; ok {
		props.RetentionDays = utils.Int32(int32(v.(int)))
	}

	if v, ok := block["storage_account_access_key"]; ok && v.(string) != "" {
		props.StorageAccountAccessKey = utils.String(v.(string))
	}

	if v, ok := block["storage_endpoint"]; ok && v.(string) != "" {
		props.StorageEndpoint = utils.String(v.(string))
	}

	return &postgresql.ServerSecurityAlertPolicy{
		SecurityAlertPolicyProperties: props,
	}
}

func flattenSecurityAlertPolicy(props *postgresql.SecurityAlertPolicyProperties, accessKey string) interface{} {
	if props == nil {
		return nil
	}

	// check if its an empty block as in its never been set before
	if props.DisabledAlerts != nil && len(*props.DisabledAlerts) == 1 && (*props.DisabledAlerts)[0] == "" &&
		props.EmailAddresses != nil && len(*props.EmailAddresses) == 1 && (*props.EmailAddresses)[0] == "" &&
		props.StorageAccountAccessKey != nil && *props.StorageAccountAccessKey == "" &&
		props.StorageEndpoint != nil && *props.StorageEndpoint == "" &&
		props.RetentionDays != nil && *props.RetentionDays == 0 &&
		props.EmailAccountAdmins != nil && !*props.EmailAccountAdmins &&
		props.State == postgresql.ServerSecurityAlertPolicyStateDisabled {
		return nil
	}

	block := map[string]interface{}{}

	block["enabled"] = props.State == postgresql.ServerSecurityAlertPolicyStateEnabled

	block["disabled_alerts"] = utils.FlattenStringSlice(props.DisabledAlerts)
	block["email_addresses"] = utils.FlattenStringSlice(props.EmailAddresses)

	if v := props.EmailAccountAdmins; v != nil {
		block["email_account_admins"] = *v
	}
	if v := props.RetentionDays; v != nil {
		block["retention_days"] = *v
	}
	if v := props.StorageEndpoint; v != nil {
		block["storage_endpoint"] = *v
	}

	block["storage_account_access_key"] = accessKey

	return []interface{}{block}
}
