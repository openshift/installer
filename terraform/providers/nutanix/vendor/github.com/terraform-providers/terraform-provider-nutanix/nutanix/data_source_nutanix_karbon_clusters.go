package nutanix

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixKarbonClusters() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceNutanixKarbonClustersRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: KarbonClusterElementDataSourceMap(),
				},
			},
		},
	}
}

func dataSourceNutanixKarbonClustersRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).KarbonAPI
	setTimeout(meta)
	// Make request to the API
	resp, err := conn.Cluster.ListKarbonClusters()
	if err != nil {
		d.SetId("")
		return nil
	}

	clusters := make([]map[string]interface{}, len(*resp))

	for k, v := range *resp {
		cluster := make(map[string]interface{})
		if err != nil {
			return fmt.Errorf("error searching for cluster via legacy API: %s", err)
		}
		karbonClusterName := *v.Name
		flattenedEtcdNodepool, err := flattenNodePools(d, conn, "etcd_node_pool", karbonClusterName, v.ETCDConfig.NodePools)
		if err != nil {
			return err
		}
		flattenedWorkerNodepool, err := flattenNodePools(d, conn, "worker_node_pool", karbonClusterName, v.WorkerConfig.NodePools)
		if err != nil {
			return err
		}
		flattenedMasterNodepool, err := flattenNodePools(d, conn, "master_node_pool", karbonClusterName, v.MasterConfig.NodePools)
		if err != nil {
			return err
		}
		cluster["name"] = utils.StringValue(v.Name)

		cluster["status"] = utils.StringValue(v.Status)

		// Must use legacy API because GA API reports different version
		cluster["version"] = utils.StringValue(v.Version)
		// cluster["version"] = utils.StringValue(respLegacy.K8sConfig.Version)
		cluster["kubeapi_server_ipv4_address"] = utils.StringValue(v.KubeAPIServerIPv4Address)
		cluster["deployment_type"] = v.MasterConfig.DeploymentType
		cluster["worker_node_pool"] = flattenedWorkerNodepool

		cluster["etcd_node_pool"] = flattenedEtcdNodepool
		cluster["master_node_pool"] = flattenedMasterNodepool
		cluster["uuid"] = utils.StringValue(v.UUID)
		clusters[k] = cluster
	}

	if err := d.Set("clusters", clusters); err != nil {
		return fmt.Errorf("failed to set clusters output: %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}
