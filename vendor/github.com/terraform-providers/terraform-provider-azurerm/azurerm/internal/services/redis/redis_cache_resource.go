package redis

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2020-06-01/redis"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	networkParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redis/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redis/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceRedisCache() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisCacheCreate,
		Read:   resourceRedisCacheRead,
		Update: resourceRedisCacheUpdate,
		Delete: resourceRedisCacheDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.CacheID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: azure.NormalizeLocation,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"zones": azure.SchemaMultipleZones(),

			"capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"family": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validate.CacheFamily,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(redis.Basic),
					string(redis.Standard),
					string(redis.Premium),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"minimum_tls_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  redis.OneFullStopZero,
				ValidateFunc: validation.StringInSlice([]string{
					string(redis.OneFullStopZero),
					string(redis.OneFullStopOne),
					string(redis.OneFullStopTwo),
				}, false),
			},

			"shard_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"enable_non_ssl_port": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.SubnetID,
			},

			"private_static_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"redis_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"maxclients": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"maxmemory_delta": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"maxmemory_reserved": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"maxmemory_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "volatile-lru",
							ValidateFunc: validate.MaxMemoryPolicy,
						},

						"maxfragmentationmemory_reserved": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"rdb_backup_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"rdb_backup_frequency": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validate.CacheBackupFrequency,
						},

						"rdb_backup_max_snapshot_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"rdb_storage_connection_string": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},

						"notify_keyspace_events": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"aof_backup_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"aof_storage_connection_string_0": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},

						"aof_storage_connection_string_1": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"enable_authentication": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"patch_schedule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day_of_week": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
						},
						"start_hour_utc": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 23),
						},
					},
				},
			},

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ssl_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"primary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"public_network_access_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceRedisCacheCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Client
	patchClient := meta.(*clients.Client).Redis.PatchSchedulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Azure ARM Redis Cache creation.")

	location := azure.NormalizeLocation(d.Get("location").(string))
	enableNonSSLPort := d.Get("enable_non_ssl_port").(bool)

	capacity := int32(d.Get("capacity").(int))
	family := redis.SkuFamily(d.Get("family").(string))
	sku := redis.SkuName(d.Get("sku_name").(string))

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	id := parse.NewCacheID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.RediName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Redis Instance %s (resource group %q): %+v", id.RediName, id.ResourceGroup, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_redis_cache", id.ID())
	}

	patchSchedule := expandRedisPatchSchedule(d)
	redisConfiguration, err := expandRedisConfiguration(d)
	if err != nil {
		return fmt.Errorf("parsing Redis Configuration: %+v", err)
	}

	publicNetworkAccess := redis.Enabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = redis.Disabled
	}

	parameters := redis.CreateParameters{
		Location: utils.String(location),
		CreateProperties: &redis.CreateProperties{
			EnableNonSslPort: utils.Bool(enableNonSSLPort),
			Sku: &redis.Sku{
				Capacity: utils.Int32(capacity),
				Family:   family,
				Name:     sku,
			},
			MinimumTLSVersion:   redis.TLSVersion(d.Get("minimum_tls_version").(string)),
			RedisConfiguration:  redisConfiguration,
			PublicNetworkAccess: publicNetworkAccess,
		},
		Tags: expandedTags,
	}

	if v, ok := d.GetOk("shard_count"); ok {
		shardCount := int32(v.(int))
		parameters.ShardCount = &shardCount
	}

	if v, ok := d.GetOk("private_static_ip_address"); ok {
		parameters.StaticIP = utils.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		parsed, parseErr := networkParse.SubnetID(v.(string))
		if parseErr != nil {
			return err
		}

		locks.ByName(parsed.VirtualNetworkName, network.VirtualNetworkResourceName)
		defer locks.UnlockByName(parsed.VirtualNetworkName, network.VirtualNetworkResourceName)

		locks.ByName(parsed.Name, network.SubnetResourceName)
		defer locks.UnlockByName(parsed.Name, network.SubnetResourceName)

		parameters.SubnetID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("zones"); ok {
		parameters.Zones = azure.ExpandZones(v.([]interface{}))
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.RediName, parameters)
	if err != nil {
		return fmt.Errorf("creating Redis Cache %q (Resource Group %q): %v", id.RediName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the create of Redis Cache %q (Resource Group %q): %v", id.RediName, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for Redis Cache %q (Resource Group %q) to become available", id.RediName, id.ResourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Scaling", "Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    redisStateRefreshFunc(ctx, client, id.ResourceGroup, id.RediName),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutCreate),
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Redis Cache %q (Resource Group %q) to become available: %s", id.RediName, id.ResourceGroup, err)
	}

	d.SetId(id.ID())

	if patchSchedule != nil {
		if _, err = patchClient.CreateOrUpdate(ctx, id.ResourceGroup, id.RediName, *patchSchedule); err != nil {
			return fmt.Errorf("setting Redis Patch Schedule: %+v", err)
		}
	}

	return resourceRedisCacheRead(d, meta)
}

func resourceRedisCacheUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Client
	patchClient := meta.(*clients.Client).Redis.PatchSchedulesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Azure ARM Redis Cache update.")

	id, err := parse.CacheID(d.Id())
	if err != nil {
		return err
	}

	enableNonSSLPort := d.Get("enable_non_ssl_port").(bool)

	capacity := int32(d.Get("capacity").(int))
	family := redis.SkuFamily(d.Get("family").(string))
	sku := redis.SkuName(d.Get("sku_name").(string))

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	parameters := redis.UpdateParameters{
		UpdateProperties: &redis.UpdateProperties{
			MinimumTLSVersion: redis.TLSVersion(d.Get("minimum_tls_version").(string)),
			EnableNonSslPort:  utils.Bool(enableNonSSLPort),
			Sku: &redis.Sku{
				Capacity: utils.Int32(capacity),
				Family:   family,
				Name:     sku,
			},
		},
		Tags: expandedTags,
	}

	if v, ok := d.GetOk("shard_count"); ok {
		if d.HasChange("shard_count") {
			shardCount := int32(v.(int))
			parameters.ShardCount = &shardCount
		}
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := redis.Enabled
		if !d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = redis.Disabled
		}
		parameters.PublicNetworkAccess = publicNetworkAccess
	}

	if d.HasChange("redis_configuration") {
		redisConfiguration, err := expandRedisConfiguration(d)
		if err != nil {
			return fmt.Errorf("parsing Redis Configuration: %+v", err)
		}
		parameters.RedisConfiguration = redisConfiguration
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.RediName, parameters); err != nil {
		return fmt.Errorf("updating Redis Cache %q (Resource Group %q): %+v", id.RediName, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for Redis Cache %q (Resource Group %q) to become available", id.RediName, id.ResourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Scaling", "Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    redisStateRefreshFunc(ctx, client, id.ResourceGroup, id.RediName),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutUpdate),
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Redis Cache %q (Resource Group %q) to become available: %+v", id.RediName, id.ResourceGroup, err)
	}

	patchSchedule := expandRedisPatchSchedule(d)

	if patchSchedule == nil || len(*patchSchedule.ScheduleEntries.ScheduleEntries) == 0 {
		_, err = patchClient.Delete(ctx, id.ResourceGroup, id.RediName)
		if err != nil {
			return fmt.Errorf("deleting Redis Patch Schedule: %+v", err)
		}
	} else {
		_, err = patchClient.CreateOrUpdate(ctx, id.ResourceGroup, id.RediName, *patchSchedule)
		if err != nil {
			return fmt.Errorf("setting Redis Patch Schedule: %+v", err)
		}
	}

	return resourceRedisCacheRead(d, meta)
}

func resourceRedisCacheRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Client
	patchSchedulesClient := meta.(*clients.Client).Redis.PatchSchedulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CacheID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.RediName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Redis Cache %q (Resource Group %q): %+v", id.RediName, id.ResourceGroup, err)
	}

	keysResp, err := client.ListKeys(ctx, id.ResourceGroup, id.RediName)
	if err != nil {
		return fmt.Errorf("listing keys for Redis Cache %q (Resource Group %q): %+v", id.RediName, id.ResourceGroup, err)
	}

	schedule, err := patchSchedulesClient.Get(ctx, id.ResourceGroup, id.RediName)
	if err == nil {
		patchSchedule := flattenRedisPatchSchedules(schedule)
		if err = d.Set("patch_schedule", patchSchedule); err != nil {
			return fmt.Errorf("setting `patch_schedule`: %+v", err)
		}
	}

	d.Set("name", id.RediName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if zones := resp.Zones; zones != nil {
		d.Set("zones", zones)
	}

	if sku := resp.Sku; sku != nil {
		d.Set("capacity", sku.Capacity)
		d.Set("family", sku.Family)
		d.Set("sku_name", sku.Name)
	}

	if props := resp.Properties; props != nil {
		d.Set("ssl_port", props.SslPort)
		d.Set("hostname", props.HostName)
		d.Set("minimum_tls_version", string(props.MinimumTLSVersion))
		d.Set("port", props.Port)
		d.Set("enable_non_ssl_port", props.EnableNonSslPort)
		if props.ShardCount != nil {
			d.Set("shard_count", props.ShardCount)
		}
		d.Set("private_static_ip_address", props.StaticIP)

		subnetId := ""
		if props.SubnetID != nil {
			parsed, err := networkParse.SubnetID(*props.SubnetID)
			if err != nil {
				return err
			}

			subnetId = parsed.ID()
		}
		d.Set("subnet_id", subnetId)

		d.Set("public_network_access_enabled", props.PublicNetworkAccess == redis.Enabled)
	}

	redisConfiguration, err := flattenRedisConfiguration(resp.RedisConfiguration)
	if err != nil {
		return fmt.Errorf("flattening `redis_configuration`: %+v", err)
	}
	if err := d.Set("redis_configuration", redisConfiguration); err != nil {
		return fmt.Errorf("setting `redis_configuration`: %+v", err)
	}

	d.Set("primary_access_key", keysResp.PrimaryKey)
	d.Set("secondary_access_key", keysResp.SecondaryKey)

	if props := resp.Properties; props != nil {
		enableSslPort := !*props.EnableNonSslPort
		d.Set("primary_connection_string", getRedisConnectionString(*props.HostName, *props.SslPort, *keysResp.PrimaryKey, enableSslPort))
		d.Set("secondary_connection_string", getRedisConnectionString(*props.HostName, *props.SslPort, *keysResp.SecondaryKey, enableSslPort))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceRedisCacheDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Client
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CacheID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.RediName)
	if err != nil {
		return fmt.Errorf("retrieving Redis Cache %q (Resource Group %q): %+v", id.RediName, id.ResourceGroup, err)
	}
	if read.Properties == nil {
		return fmt.Errorf("retrieving Redis Cache %q (Resource Group %q): `properties` was nil", id.RediName, id.ResourceGroup)
	}
	props := *read.Properties
	if subnetID := props.SubnetID; subnetID != nil {
		parsed, parseErr := networkParse.SubnetID(*subnetID)
		if parseErr != nil {
			return err
		}

		locks.ByName(parsed.VirtualNetworkName, network.VirtualNetworkResourceName)
		defer locks.UnlockByName(parsed.VirtualNetworkName, network.VirtualNetworkResourceName)

		locks.ByName(parsed.Name, network.SubnetResourceName)
		defer locks.UnlockByName(parsed.Name, network.SubnetResourceName)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.RediName)
	if err != nil {
		return fmt.Errorf("deleting Redis Cache %q (Resource Group %q): %+v", id.RediName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of Redis Cache %q (Resource Group %q): %+v", id.RediName, id.ResourceGroup, err)
		}
	}

	return nil
}

func redisStateRefreshFunc(ctx context.Context, client *redis.Client, resourceGroupName string, sgName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, sgName)
		if err != nil {
			return nil, "", fmt.Errorf("polling for status of Redis Cache %q (Resource Group %q): %+v", sgName, resourceGroupName, err)
		}

		return res, string(res.ProvisioningState), nil
	}
}

func expandRedisConfiguration(d *schema.ResourceData) (map[string]*string, error) {
	output := make(map[string]*string)

	input := d.Get("redis_configuration").([]interface{})
	if len(input) == 0 || input[0] == nil {
		return output, nil
	}
	raw := input[0].(map[string]interface{})

	if v := raw["maxclients"].(int); v > 0 {
		output["maxclients"] = utils.String(strconv.Itoa(v))
	}

	if v := raw["maxmemory_delta"].(int); v > 0 {
		output["maxmemory-delta"] = utils.String(strconv.Itoa(v))
	}

	if v := raw["maxmemory_reserved"].(int); v > 0 {
		output["maxmemory-reserved"] = utils.String(strconv.Itoa(v))
	}

	if v := raw["maxmemory_policy"].(string); v != "" {
		output["maxmemory-policy"] = utils.String(v)
	}

	if v := raw["maxfragmentationmemory_reserved"].(int); v > 0 {
		output["maxfragmentationmemory-reserved"] = utils.String(strconv.Itoa(v))
	}

	// RDB Backup
	if v := raw["rdb_backup_enabled"].(bool); v {
		if connStr := raw["rdb_storage_connection_string"].(string); connStr == "" {
			return nil, fmt.Errorf("The rdb_storage_connection_string property must be set when rdb_backup_enabled is true")
		}
		output["rdb-backup-enabled"] = utils.String(strconv.FormatBool(v))
	}

	if v := raw["rdb_backup_frequency"].(int); v > 0 {
		output["rdb-backup-frequency"] = utils.String(strconv.Itoa(v))
	}

	if v := raw["rdb_backup_max_snapshot_count"].(int); v > 0 {
		output["rdb-backup-max-snapshot-count"] = utils.String(strconv.Itoa(v))
	}

	if v := raw["rdb_storage_connection_string"].(string); v != "" {
		output["rdb-storage-connection-string"] = utils.String(v)
	}

	if v := raw["notify_keyspace_events"].(string); v != "" {
		output["notify-keyspace-events"] = utils.String(v)
	}

	// AOF Backup
	if v := raw["aof_backup_enabled"].(bool); v {
		output["aof-backup-enabled"] = utils.String(strconv.FormatBool(v))
	}

	if v := raw["aof_storage_connection_string_0"].(string); v != "" {
		output["aof-storage-connection-string-0"] = utils.String(v)
	}

	if v := raw["aof_storage_connection_string_1"].(string); v != "" {
		output["aof-storage-connection-string-1"] = utils.String(v)
	}

	authEnabled := raw["enable_authentication"].(bool)
	// Redis authentication can only be disabled if it is launched inside a VNET.
	if _, isPrivate := d.GetOk("subnet_id"); !isPrivate {
		if !authEnabled {
			return nil, fmt.Errorf("Cannot set `enable_authentication` to `false` when `subnet_id` is not set")
		}
	} else {
		value := isAuthNotRequiredAsString(authEnabled)
		output["authnotrequired"] = utils.String(value)
	}
	return output, nil
}

func expandRedisPatchSchedule(d *schema.ResourceData) *redis.PatchSchedule {
	v, ok := d.GetOk("patch_schedule")
	if !ok {
		return nil
	}

	scheduleValues := v.([]interface{})
	entries := make([]redis.ScheduleEntry, 0)
	for _, scheduleValue := range scheduleValues {
		vals := scheduleValue.(map[string]interface{})
		dayOfWeek := vals["day_of_week"].(string)
		startHourUtc := vals["start_hour_utc"].(int)

		entry := redis.ScheduleEntry{
			DayOfWeek:    redis.DayOfWeek(dayOfWeek),
			StartHourUtc: utils.Int32(int32(startHourUtc)),
		}
		entries = append(entries, entry)
	}

	schedule := redis.PatchSchedule{
		ScheduleEntries: &redis.ScheduleEntries{
			ScheduleEntries: &entries,
		},
	}
	return &schedule
}

func flattenRedisConfiguration(input map[string]*string) ([]interface{}, error) {
	outputs := make(map[string]interface{}, len(input))

	if v := input["maxclients"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `maxclients` %q: %+v", *v, err)
		}
		outputs["maxclients"] = i
	}
	if v := input["maxmemory-delta"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `maxmemory-delta` %q: %+v", *v, err)
		}
		outputs["maxmemory_delta"] = i
	}
	if v := input["maxmemory-reserved"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `maxmemory-reserved` %q: %+v", *v, err)
		}
		outputs["maxmemory_reserved"] = i
	}
	if v := input["maxmemory-policy"]; v != nil {
		outputs["maxmemory_policy"] = *v
	}

	if v := input["maxfragmentationmemory-reserved"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `maxfragmentationmemory-reserved` %q: %+v", *v, err)
		}
		outputs["maxfragmentationmemory_reserved"] = i
	}

	// delta, reserved, enabled, frequency,, count,
	if v := input["rdb-backup-enabled"]; v != nil {
		b, err := strconv.ParseBool(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `rdb-backup-enabled` %q: %+v", *v, err)
		}
		outputs["rdb_backup_enabled"] = b
	}
	if v := input["rdb-backup-frequency"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `rdb-backup-frequency` %q: %+v", *v, err)
		}
		outputs["rdb_backup_frequency"] = i
	}
	if v := input["rdb-backup-max-snapshot-count"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `rdb-backup-max-snapshot-count` %q: %+v", *v, err)
		}
		outputs["rdb_backup_max_snapshot_count"] = i
	}
	if v := input["rdb-storage-connection-string"]; v != nil {
		outputs["rdb_storage_connection_string"] = *v
	}
	if v := input["notify-keyspace-events"]; v != nil {
		outputs["notify_keyspace_events"] = *v
	}

	if v := input["aof-backup-enabled"]; v != nil {
		b, err := strconv.ParseBool(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `aof-backup-enabled` %q: %+v", *v, err)
		}
		outputs["aof_backup_enabled"] = b
	}
	if v := input["aof-storage-connection-string-0"]; v != nil {
		outputs["aof_storage_connection_string_0"] = *v
	}
	if v := input["aof-storage-connection-string-1"]; v != nil {
		outputs["aof_storage_connection_string_1"] = *v
	}

	// `authnotrequired` is not set for instances launched outside a VNET
	outputs["enable_authentication"] = true
	if v := input["authnotrequired"]; v != nil {
		outputs["enable_authentication"] = isAuthRequiredAsBool(*v)
	}

	return []interface{}{outputs}, nil
}

func isAuthRequiredAsBool(not_required string) bool {
	value := strings.ToLower(not_required)
	output := map[string]bool{
		"yes": false,
		"no":  true,
	}
	return output[value]
}

func isAuthNotRequiredAsString(auth_required bool) string {
	output := map[bool]string{
		true:  "no",
		false: "yes",
	}
	return output[auth_required]
}

func flattenRedisPatchSchedules(schedule redis.PatchSchedule) []interface{} {
	outputs := make([]interface{}, 0)

	for _, entry := range *schedule.ScheduleEntries.ScheduleEntries {
		output := make(map[string]interface{})

		output["day_of_week"] = string(entry.DayOfWeek)
		output["start_hour_utc"] = int(*entry.StartHourUtc)

		outputs = append(outputs, output)
	}

	return outputs
}

func getRedisConnectionString(redisHostName string, sslPort int32, accessKey string, enableSslPort bool) string {
	return fmt.Sprintf("%s:%d,password=%s,ssl=%t,abortConnect=False", redisHostName, sslPort, accessKey, enableSslPort)
}
