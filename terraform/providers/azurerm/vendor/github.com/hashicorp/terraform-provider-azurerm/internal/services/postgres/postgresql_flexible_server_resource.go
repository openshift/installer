package postgres

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2021-06-01/serverrestart"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2022-12-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2018-09-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	ServerMaintenanceWindowEnabled  = "Enabled"
	ServerMaintenanceWindowDisabled = "Disabled"
)

var postgresqlFlexibleServerResourceName = "azurerm_postgresql_flexible_server"

func resourcePostgresqlFlexibleServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePostgresqlFlexibleServerCreate,
		Read:   resourcePostgresqlFlexibleServerRead,
		Update: resourcePostgresqlFlexibleServerUpdate,
		Delete: resourcePostgresqlFlexibleServerDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(1 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(1 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(1 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := servers.ParseFlexibleServerID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FlexibleServerName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"administrator_login": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace, validate.AdminUsernames),
			},

			"administrator_password": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"authentication": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"active_directory_auth_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"password_auth_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
						"tenant_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
							RequiredWith: []string{
								"authentication.0.active_directory_auth_enabled",
							},
						},
					},
				},
			},

			"sku_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.FlexibleServerSkuName,
			},

			"storage_mb": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{32768, 65536, 131072, 262144, 524288, 1048576, 2097152, 4194304, 8388608, 16777216}),
			},

			"version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(servers.PossibleValuesForServerVersion(), false),
			},

			"zone": commonschema.ZoneSingleOptional(),

			"create_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(servers.CreateModeDefault),
					string(servers.CreateModePointInTimeRestore),
					string(servers.CreateModeReplica),
				}, false),
			},

			"delegated_subnet_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.SubnetID,
			},

			"private_dns_zone_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				// This is `computed`, because there is a breaking change to require this field when setting vnet.
				// For existing fs who don't want to be recreated, they could contact service team to manually migrate to the private dns zone
				// We need to ignore the diff when remote is set private dns zone
				ForceNew:     true,
				ValidateFunc: privatezones.ValidatePrivateDnsZoneID,
			},

			"point_in_time_restore_time_in_utc": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"source_server_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: servers.ValidateFlexibleServerID,
			},

			"maintenance_window": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"day_of_week": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 6),
						},

						"start_hour": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 23),
						},

						"start_minute": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 59),
						},
					},
				},
			},

			"backup_retention_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(7, 35),
			},

			"geo_redundant_backup_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"high_availability": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"mode": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(servers.HighAvailabilityModeZoneRedundant),
								string(servers.HighAvailabilityModeSameZone),
							}, false),
						},

						"standby_availability_zone": commonschema.ZoneSingleOptional(),
					},
				},
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"replication_role": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(servers.ReplicationRoleNone),
				}, false),
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"customer_managed_key": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_key_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
							RequiredWith: []string{
								"identity",
								"customer_managed_key.0.primary_user_assigned_identity_id",
							},
						},
						"primary_user_assigned_identity_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourcePostgresqlFlexibleServerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := servers.NewFlexibleServerID(subscriptionId, resourceGroup, name)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_postgresql_flexible_server", id.ID())
	}

	createMode := d.Get("create_mode").(string)

	if _, ok := d.GetOk("replication_role"); ok {
		return fmt.Errorf("`replication_role` cannot be set while creating")
	}

	if servers.CreateMode(createMode) == servers.CreateModePointInTimeRestore {
		if _, ok := d.GetOk("source_server_id"); !ok {
			return fmt.Errorf("`source_server_id` is required when `create_mode` is `PointInTimeRestore`")
		}
		if _, ok := d.GetOk("point_in_time_restore_time_in_utc"); !ok {
			return fmt.Errorf("`point_in_time_restore_time_in_utc` is required when `create_mode` is `PointInTimeRestore`")
		}
	}

	if servers.CreateMode(createMode) == servers.CreateModeReplica {
		if _, ok := d.GetOk("source_server_id"); !ok {
			return fmt.Errorf("`source_server_id` is required when `create_mode` is `Replica`")
		}
	}

	if createMode == "" || servers.CreateMode(createMode) == servers.CreateModeDefault {
		if _, ok := d.GetOk("administrator_login"); !ok {
			return fmt.Errorf("`administrator_login` is required when `create_mode` is `Default`")
		}
		if _, ok := d.GetOk("administrator_password"); !ok {
			return fmt.Errorf("`administrator_password` is required when `create_mode` is `Default`")
		}
		if _, ok := d.GetOk("sku_name"); !ok {
			return fmt.Errorf("`sku_name` is required when `create_mode` is `Default`")
		}
		if _, ok := d.GetOk("version"); !ok {
			return fmt.Errorf("`version` is required when `create_mode` is `Default`")
		}
		if _, ok := d.GetOk("storage_mb"); !ok {
			return fmt.Errorf("`storage_mb` is required when `create_mode` is `Default`")
		}
	}

	sku, err := expandFlexibleServerSku(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name` for %s: %v", id, err)
	}

	parameters := servers.Server{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &servers.ServerProperties{
			Network:          expandArmServerNetwork(d),
			Storage:          expandArmServerStorage(d),
			HighAvailability: expandFlexibleServerHighAvailability(d.Get("high_availability").([]interface{}), true),
			Backup:           expandArmServerBackup(d),
			DataEncryption:   expandFlexibleServerDataEncryption(d.Get("customer_managed_key").([]interface{})),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("administrator_login"); ok && v.(string) != "" {
		parameters.Properties.AdministratorLogin = utils.String(v.(string))
	}

	if v, ok := d.GetOk("administrator_password"); ok && v.(string) != "" {
		parameters.Properties.AdministratorLoginPassword = utils.String(v.(string))
	}

	if createMode != "" {
		createModeAttr := servers.CreateMode(createMode)
		parameters.Properties.CreateMode = &createModeAttr
	}

	if v, ok := d.GetOk("version"); ok && v.(string) != "" {
		version := servers.ServerVersion(v.(string))
		parameters.Properties.Version = &version
	}

	if v, ok := d.GetOk("zone"); ok && v.(string) != "" {
		parameters.Properties.AvailabilityZone = utils.String(v.(string))
	}

	if v, ok := d.GetOk("source_server_id"); ok && v.(string) != "" {
		parameters.Properties.SourceServerResourceId = utils.String(v.(string))
	}

	pointInTimeUTC := d.Get("point_in_time_restore_time_in_utc").(string)
	if pointInTimeUTC != "" {
		v, err := time.Parse(time.RFC3339, pointInTimeUTC)
		if err != nil {
			return fmt.Errorf("unable to parse `point_in_time_restore_time_in_utc` value")
		}
		parameters.Properties.PointInTimeUTC = utils.String(v.Format(time.RFC3339))
	}

	// if create with `password_auth_enabled` set to `false`, the service will not accept `administrator_login`.
	// so we create it with  `password_auth_enabled` set to `true`, then set it to `false` in an additional update.
	if authRaw, ok := d.GetOk("authentication"); ok {
		authConfig := expandFlexibleServerAuthConfig(authRaw.([]interface{}))
		passwordAuthEnabled := servers.PasswordAuthEnumEnabled
		authConfig.PasswordAuth = &passwordAuthEnabled
		parameters.Properties.AuthConfig = authConfig
	}

	identity, err := expandFlexibleServerIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`")
	}
	parameters.Identity = identity

	if err = client.CreateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	requireAdditionalUpdate := false
	updateProperties := servers.ServerPropertiesForUpdate{}
	if authRaw, ok := d.GetOk("authentication"); ok {
		authConfig := expandFlexibleServerAuthConfig(authRaw.([]interface{}))
		if authConfig != nil && authConfig.PasswordAuth != nil && *authConfig.PasswordAuth == servers.PasswordAuthEnumDisabled {
			requireAdditionalUpdate = true
			updateProperties.AuthConfig = authConfig
		}
	}

	// `maintenance_window` could only be updated with, could not be created with
	if v, ok := d.GetOk("maintenance_window"); ok {
		requireAdditionalUpdate = true
		updateProperties.MaintenanceWindow = expandArmServerMaintenanceWindow(v.([]interface{}))
	}

	if requireAdditionalUpdate {
		update := servers.ServerForUpdate{
			Properties: &updateProperties,
		}
		if err = client.UpdateThenPoll(ctx, id, update); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourcePostgresqlFlexibleServerRead(d, meta)
}

func resourcePostgresqlFlexibleServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseFlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Postgresql Flexibleserver %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.FlexibleServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))

		if props := model.Properties; props != nil {
			d.Set("administrator_login", props.AdministratorLogin) // if pwdEnabled is set to false, then the service does not return the value of AdministratorLogin
			d.Set("zone", props.AvailabilityZone)
			d.Set("version", props.Version)
			d.Set("fqdn", props.FullyQualifiedDomainName)

			// Currently, `replicationRole` is set to `Primary` when `createMode` is `Replica` and `replicationRole` is updated to `None`. Service team confirmed it should be set to `None` for this scenario. See more details from https://github.com/Azure/azure-rest-api-specs/issues/22499
			d.Set("replication_role", d.Get("replication_role").(string))

			if network := props.Network; network != nil {
				publicNetworkAccess := false
				if network.PublicNetworkAccess != nil {
					publicNetworkAccess = *network.PublicNetworkAccess == servers.ServerPublicNetworkAccessStateEnabled
				}
				d.Set("public_network_access_enabled", publicNetworkAccess)
				d.Set("delegated_subnet_id", network.DelegatedSubnetResourceId)
				d.Set("private_dns_zone_id", network.PrivateDnsZoneArmResourceId)
			}

			if err := d.Set("maintenance_window", flattenArmServerMaintenanceWindow(props.MaintenanceWindow)); err != nil {
				return fmt.Errorf("setting `maintenance_window`: %+v", err)
			}

			if storage := props.Storage; storage != nil && storage.StorageSizeGB != nil {
				d.Set("storage_mb", (*storage.StorageSizeGB * 1024))
			}

			if backup := props.Backup; backup != nil {
				d.Set("backup_retention_days", backup.BackupRetentionDays)

				geoRedundantBackup := false
				if backup.GeoRedundantBackup != nil {
					geoRedundantBackup = *backup.GeoRedundantBackup == servers.GeoRedundantBackupEnumEnabled
				}
				d.Set("geo_redundant_backup_enabled", geoRedundantBackup)
			}

			if err := d.Set("high_availability", flattenFlexibleServerHighAvailability(props.HighAvailability)); err != nil {
				return fmt.Errorf("setting `high_availability`: %+v", err)
			}

			if props.AuthConfig != nil {
				d.Set("authentication", flattenFlexibleServerAuthConfig(props.AuthConfig))
			}

			cmk, err := flattenFlexibleServerDataEncryption(props.DataEncryption)
			if err != nil {
				return fmt.Errorf("flattening `customer_managed_key`: %+v", err)
			}
			if err := d.Set("customer_managed_key", cmk); err != nil {
				return fmt.Errorf("setting `customer_managed_key`: %+v", err)
			}

			id, err := flattenFlexibleServerIdentity(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			if err := d.Set("identity", id); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}
		}

		sku, err := flattenFlexibleServerSku(model.Sku)
		if err != nil {
			return fmt.Errorf("flattening `sku_name` for %s: %v", id, err)
		}

		d.Set("sku_name", sku)

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourcePostgresqlFlexibleServerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseFlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	parameters := servers.ServerForUpdate{
		Properties: &servers.ServerPropertiesForUpdate{},
	}

	var requireFailover bool
	// failover is only supported when `zone` and `high_availability.0.standby_availability_zone` are exchanged with each other
	if d.HasChanges("zone", "high_availability") {
		resp, err := client.Get(ctx, *id)
		if err != nil || resp.Model == nil {
			return err
		}
		props := resp.Model.Properties

		if d.HasChange("zone") {
			if !d.HasChange("high_availability.0.standby_availability_zone") {
				return fmt.Errorf("`zone` can only be changed when exchanged with the zone specified in `high_availability.0.standby_availability_zone`")
			} else {
				// zone can only be changed when it is swapped for an existing high_availability.0.standby_availability_zone - a failover is triggered to make it the new primary availability zone
				// compare current values of zone and high_availability.0.standby_availability_zone with new values and only allow update/failover if the values of zone and an existing high_availability.0.standby_availability_zone have been swapped
				var newZone, newHAStandbyZone string
				newZone = d.Get("zone").(string)
				newHAStandbyZone = d.Get("high_availability.0.standby_availability_zone").(string)
				if props != nil && props.AvailabilityZone != nil && props.HighAvailability != nil && props.HighAvailability.StandbyAvailabilityZone != nil {
					if newZone == *props.HighAvailability.StandbyAvailabilityZone && newHAStandbyZone == *props.AvailabilityZone {
						requireFailover = true
					} else {
						return fmt.Errorf("`zone` can only be changed when exchanged with the zone specified in `high_availability.0.standby_availability_zone`")
					}
				}
			}

			// changes can occur in high_availability.0.standby_availability_zone when zone has not changed in the case where a high_availability block has been newly added or a high_availability block is removed, meaning HA is now disabled
		} else if d.HasChange("high_availability.0.standby_availability_zone") {
			if props != nil && props.HighAvailability != nil && props.HighAvailability.Mode != nil {
				// if HA Mode is currently "ZoneRedundant" and is still set to "ZoneRedundant", high_availability.0.standby_availability_zone cannot be changed
				if *props.HighAvailability.Mode == servers.HighAvailabilityModeZoneRedundant && !d.HasChange("high_availability.0.mode") {
					return fmt.Errorf("an existing `high_availability.0.standby_availability_zone` can only be changed when exchanged with the zone specified in `zone`")
				}
				// if high_availability.0.mode changes from "ZoneRedundant", an existing high_availability block has been removed as this is a required field
				// if high_availability.0.mode is not currently "ZoneRedundant", this must be a newly added block
			}
		}
	}

	if d.HasChange("replication_role") {
		createMode := d.Get("create_mode").(string)
		replicationRole := d.Get("replication_role").(string)
		if createMode == string(servers.CreateModeReplica) && replicationRole == string(servers.ReplicationRoleNone) {
			replicationRole := servers.ReplicationRoleNone
			parameters := servers.ServerForUpdate{
				Properties: &servers.ServerPropertiesForUpdate{
					ReplicationRole: &replicationRole,
				},
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating `replication_role` for %s: %+v", *id, err)
			}
		} else {
			return fmt.Errorf("`replication_role` only can be updated to `None` for replica server")
		}
	}

	if d.HasChange("administrator_password") {
		parameters.Properties.AdministratorLoginPassword = utils.String(d.Get("administrator_password").(string))
	}

	if d.HasChange("authentication") {
		parameters.Properties.AuthConfig = expandFlexibleServerAuthConfig(d.Get("authentication").([]interface{}))
	}

	if d.HasChange("storage_mb") {
		parameters.Properties.Storage = expandArmServerStorage(d)
	}

	if d.HasChange("backup_retention_days") {
		parameters.Properties.Backup = expandArmServerBackup(d)
	}

	if d.HasChange("maintenance_window") {
		parameters.Properties.MaintenanceWindow = expandArmServerMaintenanceWindow(d.Get("maintenance_window").([]interface{}))
	}

	if d.HasChange("sku_name") {
		sku, err := expandFlexibleServerSku(d.Get("sku_name").(string))
		if err != nil {
			return fmt.Errorf("expanding `sku_name` for %s: %v", id, err)
		}
		parameters.Sku = sku
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("high_availability") {
		parameters.Properties.HighAvailability = expandFlexibleServerHighAvailability(d.Get("high_availability").([]interface{}), false)
	}

	if d.HasChange("customer_managed_key") {
		parameters.Properties.DataEncryption = expandFlexibleServerDataEncryption(d.Get("customer_managed_key").([]interface{}))
	}

	if d.HasChange("identity") {
		identity, err := expandFlexibleServerIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity` for Mysql Flexible Server %s (Resource Group %q): %v", id.FlexibleServerName, id.ResourceGroupName, err)
		}
		parameters.Identity = identity
	}

	if err = client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if requireFailover {
		restartClient := meta.(*clients.Client).Postgres.ServerRestartClient

		restartServerId := serverrestart.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName)
		failoverMode := serverrestart.FailoverModePlannedFailover
		restartParameters := serverrestart.RestartParameter{
			RestartWithFailover: utils.Bool(true),
			FailoverMode:        &failoverMode,
		}

		if err = restartClient.ServersRestartThenPoll(ctx, restartServerId, restartParameters); err != nil {
			return fmt.Errorf("failing over %s: %+v", *id, err)
		}
	}

	return resourcePostgresqlFlexibleServerRead(d, meta)
}

func resourcePostgresqlFlexibleServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseFlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandArmServerNetwork(d *pluginsdk.ResourceData) *servers.Network {
	network := servers.Network{}

	if v, ok := d.GetOk("delegated_subnet_id"); ok {
		network.DelegatedSubnetResourceId = utils.String(v.(string))
	}

	if v, ok := d.GetOk("private_dns_zone_id"); ok {
		network.PrivateDnsZoneArmResourceId = utils.String(v.(string))
	}

	return &network
}

func expandArmServerMaintenanceWindow(input []interface{}) *servers.MaintenanceWindow {
	if len(input) == 0 {
		return &servers.MaintenanceWindow{
			CustomWindow: utils.String(ServerMaintenanceWindowDisabled),
		}
	}
	v := input[0].(map[string]interface{})

	maintenanceWindow := servers.MaintenanceWindow{
		CustomWindow: utils.String(ServerMaintenanceWindowEnabled),
		StartHour:    utils.Int64(int64(v["start_hour"].(int))),
		StartMinute:  utils.Int64(int64(v["start_minute"].(int))),
		DayOfWeek:    utils.Int64(int64(v["day_of_week"].(int))),
	}

	return &maintenanceWindow
}

func expandArmServerStorage(d *pluginsdk.ResourceData) *servers.Storage {
	storage := servers.Storage{}

	if v, ok := d.GetOk("storage_mb"); ok {
		storage.StorageSizeGB = utils.Int64(int64(v.(int) / 1024))
	}

	return &storage
}

func expandArmServerBackup(d *pluginsdk.ResourceData) *servers.Backup {
	backup := servers.Backup{}

	if v, ok := d.GetOk("backup_retention_days"); ok {
		backup.BackupRetentionDays = utils.Int64(int64(v.(int)))
	}

	geoRedundantEnabled := servers.GeoRedundantBackupEnumDisabled
	if geoRedundantBackupEnabled := d.Get("geo_redundant_backup_enabled").(bool); geoRedundantBackupEnabled {
		geoRedundantEnabled = servers.GeoRedundantBackupEnumEnabled
	}

	backup.GeoRedundantBackup = &geoRedundantEnabled

	return &backup
}

func expandFlexibleServerSku(name string) (*servers.Sku, error) {
	if name == "" {
		return nil, nil
	}
	parts := strings.SplitAfterN(name, "_", 2)

	var tier servers.SkuTier
	switch strings.TrimSuffix(parts[0], "_") {
	case "B":
		tier = servers.SkuTierBurstable
	case "GP":
		tier = servers.SkuTierGeneralPurpose
	case "MO":
		tier = servers.SkuTierMemoryOptimized
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", name, parts[0])
	}

	return &servers.Sku{
		Name: parts[1],
		Tier: tier,
	}, nil
}

func flattenFlexibleServerSku(sku *servers.Sku) (string, error) {
	if sku == nil || sku.Tier == "" {
		return "", nil
	}

	var tier string
	switch sku.Tier {
	case servers.SkuTierBurstable:
		tier = "B"
	case servers.SkuTierGeneralPurpose:
		tier = "GP"
	case servers.SkuTierMemoryOptimized:
		tier = "MO"
	default:
		return "", fmt.Errorf("sku_name has unknown sku tier %s", sku.Tier)
	}

	return strings.Join([]string{tier, sku.Name}, "_"), nil
}

func flattenArmServerMaintenanceWindow(input *servers.MaintenanceWindow) []interface{} {
	if input == nil || input.CustomWindow == nil || *input.CustomWindow == ServerMaintenanceWindowDisabled {
		return make([]interface{}, 0)
	}

	var dayOfWeek int64
	if input.DayOfWeek != nil {
		dayOfWeek = *input.DayOfWeek
	}
	var startHour int64
	if input.StartHour != nil {
		startHour = *input.StartHour
	}
	var startMinute int64
	if input.StartMinute != nil {
		startMinute = *input.StartMinute
	}
	return []interface{}{
		map[string]interface{}{
			"day_of_week":  dayOfWeek,
			"start_hour":   startHour,
			"start_minute": startMinute,
		},
	}
}

func expandFlexibleServerHighAvailability(inputs []interface{}, isCreate bool) *servers.HighAvailability {
	if len(inputs) == 0 || inputs[0] == nil {
		highAvailability := servers.HighAvailabilityModeDisabled
		return &servers.HighAvailability{
			Mode: &highAvailability,
		}
	}

	input := inputs[0].(map[string]interface{})

	mode := servers.HighAvailabilityMode(input["mode"].(string))
	result := servers.HighAvailability{
		Mode: &mode,
	}

	// service team confirmed it doesn't support to update `high_availability.0.standby_availability_zone` after the PostgreSQL Flexible Server resource is created
	if isCreate {
		if v, ok := input["standby_availability_zone"]; ok && v.(string) != "" {
			result.StandbyAvailabilityZone = utils.String(v.(string))
		}
	}

	return &result
}

func flattenFlexibleServerHighAvailability(ha *servers.HighAvailability) []interface{} {
	if ha == nil || ha.Mode == nil || *ha.Mode == servers.HighAvailabilityModeDisabled {
		return []interface{}{}
	}

	var zone string
	if ha.StandbyAvailabilityZone != nil {
		zone = *ha.StandbyAvailabilityZone
	}

	return []interface{}{
		map[string]interface{}{
			"mode":                      string(*ha.Mode),
			"standby_availability_zone": zone,
		},
	}
}

func expandFlexibleServerAuthConfig(authRaw []interface{}) *servers.AuthConfig {
	if len(authRaw) == 0 || authRaw[0] == nil {
		return nil
	}

	authConfigs := authRaw[0].(map[string]interface{})
	out := servers.AuthConfig{}

	activeDirectoryAuthEnabled := servers.ActiveDirectoryAuthEnumDisabled
	if authConfigs["active_directory_auth_enabled"].(bool) {
		activeDirectoryAuthEnabled = servers.ActiveDirectoryAuthEnumEnabled
	}
	out.ActiveDirectoryAuth = &activeDirectoryAuthEnabled

	passwordAuthEnabled := servers.PasswordAuthEnumDisabled
	if authConfigs["password_auth_enabled"].(bool) {
		passwordAuthEnabled = servers.PasswordAuthEnumEnabled
	}
	out.PasswordAuth = &passwordAuthEnabled

	if tenantId, ok := authConfigs["tenant_id"].(string); ok {
		out.TenantId = &tenantId
	}

	return &out
}

func flattenFlexibleServerAuthConfig(ac *servers.AuthConfig) interface{} {
	out := make(map[string]interface{}, 0)

	if ac == nil {
		return out
	}

	aadEnabled := false
	if ac.ActiveDirectoryAuth != nil {
		aadEnabled = *ac.ActiveDirectoryAuth == servers.ActiveDirectoryAuthEnumEnabled
	}
	out["active_directory_auth_enabled"] = aadEnabled

	// It is by design if PasswordAuthEnabled is not returned or undefined, we consider it as true.
	pwdEnabled := true
	if ac.PasswordAuth != nil {
		pwdEnabled = *ac.PasswordAuth == servers.PasswordAuthEnumEnabled
	}
	out["password_auth_enabled"] = pwdEnabled

	if ac.TenantId != nil {
		out["tenant_id"] = *ac.TenantId
	}

	result := make([]interface{}, 0)
	result = append(result, out)
	return result
}

func expandFlexibleServerDataEncryption(input []interface{}) *servers.DataEncryption {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	det := servers.ArmServerKeyTypeAzureKeyVault
	dataEncryption := servers.DataEncryption{
		Type:                          &det,
		PrimaryKeyURI:                 utils.String(v["key_vault_key_id"].(string)),
		PrimaryUserAssignedIdentityId: utils.String(v["primary_user_assigned_identity_id"].(string)),
	}

	return &dataEncryption
}

func flattenFlexibleServerDataEncryption(de *servers.DataEncryption) ([]interface{}, error) {
	if de == nil || *de.Type != servers.ArmServerKeyTypeAzureKeyVault {
		return []interface{}{}, nil
	}

	item := map[string]interface{}{}
	if de.PrimaryKeyURI != nil {
		item["key_vault_key_id"] = *de.PrimaryKeyURI
	}
	if identity := de.PrimaryUserAssignedIdentityId; identity != nil {
		parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(*identity)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", *identity, err)
		}
		item["primary_user_assigned_identity_id"] = parsed.ID()
	}

	return []interface{}{item}, nil
}

func expandFlexibleServerIdentity(input []interface{}) (*servers.UserAssignedIdentity, error) {
	expanded, err := identity.ExpandUserAssignedMap(input)
	if err != nil || expanded.Type != identity.TypeUserAssigned {
		return nil, err
	}

	idUserAssigned := servers.IdentityTypeUserAssigned
	out := servers.UserAssignedIdentity{
		Type: idUserAssigned,
	}
	if expanded.Type == identity.TypeUserAssigned {
		ids := make(map[string]servers.UserIdentity)
		for k := range expanded.IdentityIds {
			ids[k] = servers.UserIdentity{}
		}
		out.UserAssignedIdentities = &ids
	}

	return &out, nil
}

func flattenFlexibleServerIdentity(input *servers.UserAssignedIdentity) (*[]interface{}, error) {
	var transform *identity.UserAssignedMap

	if input != nil {
		transform = &identity.UserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		for k := range *input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{}
		}
	}

	return identity.FlattenUserAssignedMap(transform)
}
