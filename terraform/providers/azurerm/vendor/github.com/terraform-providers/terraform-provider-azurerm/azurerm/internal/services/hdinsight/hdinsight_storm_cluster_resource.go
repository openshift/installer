package hdinsight

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hdinsight/parse"

	"github.com/Azure/azure-sdk-for-go/services/hdinsight/mgmt/2018-06-01/hdinsight"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightStormClusterHeadNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         4,
	MaxInstanceCount:         utils.Int(4),
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: utils.Int32(int32(2)),
}

var hdInsightStormClusterWorkerNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	CanSpecifyDisks:         false,
}

var hdInsightStormClusterZookeeperNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         3,
	MaxInstanceCount:         utils.Int(3),
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: utils.Int32(int32(3)),
}

func resourceHDInsightStormCluster() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: `HDInsight 3.6 will be retired on 2020-12-31 - Storm is not supported in HDInsight 4.0 and so this resource will be removed in the next major version of the AzureRM Terraform Provider.
		
More information on the HDInsight 3.6 deprecation can be found at:

https://docs.microsoft.com/en-us/azure/hdinsight/hdinsight-component-versioning#available-versions`,
		Create: resourceHDInsightStormClusterCreate,
		Read:   resourceHDInsightStormClusterRead,
		Update: hdinsightClusterUpdate("Storm", resourceHDInsightStormClusterRead),
		Delete: hdinsightClusterDelete("Storm"),
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
			"name": SchemaHDInsightName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"cluster_version": SchemaHDInsightClusterVersion(),

			"tier": SchemaHDInsightTier(),

			"tls_min_version": SchemaHDInsightTls(),

			"component_version": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storm": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"gateway": SchemaHDInsightsGateway(),

			"metastores": SchemaHDInsightsExternalMetastores(),

			"storage_account": SchemaHDInsightsStorageAccounts(),

			"roles": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"head_node": SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightStormClusterHeadNodeDefinition, true),

						"worker_node": SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightStormClusterWorkerNodeDefinition, true),

						"zookeeper_node": SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightStormClusterZookeeperNodeDefinition, true),
					},
				},
			},

			"tags": tags.Schema(),

			"https_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ssh_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"monitor": SchemaHDInsightsMonitor(),
		},
	}
}

func resourceHDInsightStormClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HDInsight.ClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	extensionsClient := meta.(*clients.Client).HDInsight.ExtensionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := parse.NewClusterID(subscriptionId, resourceGroup, name)
	location := azure.NormalizeLocation(d.Get("location").(string))
	clusterVersion := d.Get("cluster_version").(string)
	t := d.Get("tags").(map[string]interface{})
	tier := hdinsight.Tier(d.Get("tier").(string))
	tls := d.Get("tls_min_version").(string)

	componentVersionsRaw := d.Get("component_version").([]interface{})
	componentVersions := expandHDInsightStormComponentVersion(componentVersionsRaw)

	gatewayRaw := d.Get("gateway").([]interface{})
	configurations := ExpandHDInsightsConfigurations(gatewayRaw)

	metastoresRaw := d.Get("metastores").([]interface{})
	metastores := expandHDInsightsMetastore(metastoresRaw)
	for k, v := range metastores {
		configurations[k] = v
	}

	storageAccountsRaw := d.Get("storage_account").([]interface{})
	storageAccounts, identity, err := ExpandHDInsightsStorageAccounts(storageAccountsRaw, nil)
	if err != nil {
		return fmt.Errorf("failure expanding `storage_account`: %s", err)
	}

	stormRoles := hdInsightRoleDefinition{
		HeadNodeDef:      hdInsightStormClusterHeadNodeDefinition,
		WorkerNodeDef:    hdInsightStormClusterWorkerNodeDefinition,
		ZookeeperNodeDef: hdInsightStormClusterZookeeperNodeDefinition,
	}
	rolesRaw := d.Get("roles").([]interface{})
	roles, err := expandHDInsightRoles(rolesRaw, stormRoles)
	if err != nil {
		return fmt.Errorf("failure expanding `roles`: %+v", err)
	}

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("failure checking for presence of existing HDInsight Storm Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_hdinsight_storm_cluster", *existing.ID)
	}

	params := hdinsight.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Properties: &hdinsight.ClusterCreateProperties{
			Tier:                   tier,
			OsType:                 hdinsight.Linux,
			ClusterVersion:         utils.String(clusterVersion),
			MinSupportedTLSVersion: utils.String(tls),
			ClusterDefinition: &hdinsight.ClusterDefinition{
				Kind:             utils.String("Storm"),
				ComponentVersion: componentVersions,
				Configurations:   configurations,
			},
			StorageProfile: &hdinsight.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &hdinsight.ComputeProfile{
				Roles: roles,
			},
		},
		Tags:     tags.Expand(t),
		Identity: identity,
	}
	future, err := client.Create(ctx, resourceGroup, name, params)
	if err != nil {
		return fmt.Errorf("failure creating HDInsight Storm Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed waiting for creation of HDInsight Storm Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("failure retrieving HDInsight Storm Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("failure reading ID for HDInsight Storm Cluster %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(id.ID())

	// We can only enable monitoring after creation
	if v, ok := d.GetOk("monitor"); ok {
		monitorRaw := v.([]interface{})
		if err := enableHDInsightMonitoring(ctx, extensionsClient, resourceGroup, name, monitorRaw); err != nil {
			return err
		}
	}

	return resourceHDInsightStormClusterRead(d, meta)
}

func resourceHDInsightStormClusterRead(d *schema.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).HDInsight.ClustersClient
	configurationsClient := meta.(*clients.Client).HDInsight.ConfigurationsClient
	extensionsClient := meta.(*clients.Client).HDInsight.ExtensionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Name

	resp, err := clustersClient.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] HDInsight Storm Cluster %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("failure retrieving HDInsight Storm Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// Each call to configurationsClient methods is HTTP request. Getting all settings in one operation
	configurations, err := configurationsClient.List(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("failure retrieving Configuration for HDInsight Storm Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	gateway, exists := configurations.Configurations["gateway"]
	if !exists {
		return fmt.Errorf("failure retrieving gateway for HDInsight Storm Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	// storage_account isn't returned so I guess we just leave it ¯\_(ツ)_/¯
	if props := resp.Properties; props != nil {
		d.Set("cluster_version", props.ClusterVersion)
		d.Set("tier", string(props.Tier))
		d.Set("tls_min_version", props.MinSupportedTLSVersion)

		if def := props.ClusterDefinition; def != nil {
			if err := d.Set("component_version", flattenHDInsightStormComponentVersion(def.ComponentVersion)); err != nil {
				return fmt.Errorf("failure flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", FlattenHDInsightsConfigurations(gateway)); err != nil {
				return fmt.Errorf("failure flattening `gateway`: %+v", err)
			}

			flattenHDInsightsMetastores(d, configurations.Configurations)
		}

		stormRoles := hdInsightRoleDefinition{
			HeadNodeDef:      hdInsightStormClusterHeadNodeDefinition,
			WorkerNodeDef:    hdInsightStormClusterWorkerNodeDefinition,
			ZookeeperNodeDef: hdInsightStormClusterZookeeperNodeDefinition,
		}
		flattenedRoles := flattenHDInsightRoles(d, props.ComputeProfile, stormRoles)
		if err := d.Set("roles", flattenedRoles); err != nil {
			return fmt.Errorf("failure flattening `roles`: %+v", err)
		}

		httpEndpoint := FindHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
		d.Set("https_endpoint", httpEndpoint)
		sshEndpoint := FindHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
		d.Set("ssh_endpoint", sshEndpoint)

		monitor, err := extensionsClient.GetMonitoringStatus(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("failed reading monitor configuration for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		d.Set("monitor", flattenHDInsightMonitoring(monitor))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func expandHDInsightStormComponentVersion(input []interface{}) map[string]*string {
	vs := input[0].(map[string]interface{})
	return map[string]*string{
		"Storm": utils.String(vs["storm"].(string)),
	}
}

func flattenHDInsightStormComponentVersion(input map[string]*string) []interface{} {
	stormVersion := ""
	if v, ok := input["Storm"]; ok {
		if v != nil {
			stormVersion = *v
		}
	}
	return []interface{}{
		map[string]interface{}{
			"storm": stormVersion,
		},
	}
}
