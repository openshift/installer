package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/karbon"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixKarbonCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceNutanixKarbonClusterRead,
		SchemaVersion: 1,
		Schema:        KarbonClusterDataSourceMap(),
	}
}

func dataSourceNutanixKarbonClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection
	conn := meta.(*Client).KarbonAPI
	setTimeout(meta)
	// Make request to the API
	karbonClusterID, iok := d.GetOk("karbon_cluster_id")
	karbonClusterNameInput, nok := d.GetOk("karbon_cluster_name")
	if !iok && !nok {
		return diag.Errorf("please provide one of karbon_cluster_id or karbon_cluster_name attributes")
	}
	var err error
	var resp *karbon.ClusterIntentResponse

	if iok {
		resp, err = conn.Cluster.GetKarbonCluster(karbonClusterID.(string))
	} else {
		resp, err = conn.Cluster.GetKarbonCluster(karbonClusterNameInput.(string))
	}

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	karbonClusterName := *resp.Name
	flattenedEtcdNodepool, err := flattenNodePools(d, conn, "etcd_node_pool", karbonClusterName, resp.ETCDConfig.NodePools)
	if err != nil {
		return diag.FromErr(err)
	}
	flattenedWorkerNodepool, err := flattenNodePools(d, conn, "worker_node_pool", karbonClusterName, resp.WorkerConfig.NodePools)
	if err != nil {
		return diag.FromErr(err)
	}
	flattenedMasterNodepool, err := flattenNodePools(d, conn, "master_node_pool", karbonClusterName, resp.MasterConfig.NodePools)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", utils.StringValue(resp.Name)); err != nil {
		return diag.Errorf("error setting name for Karbon Cluster %s: %s", d.Id(), err)
	}
	if err = d.Set("status", utils.StringValue(resp.Status)); err != nil {
		return diag.Errorf("error setting status for Karbon Cluster %s: %s", d.Id(), err)
	}
	if err = d.Set("version", utils.StringValue(resp.Version)); err != nil {
		return diag.Errorf("error setting version for Karbon Cluster %s: %s", d.Id(), err)
	}
	if err = d.Set("kubeapi_server_ipv4_address", utils.StringValue(resp.KubeAPIServerIPv4Address)); err != nil {
		return diag.Errorf("error setting kubeapi_server_ipv4_address for Karbon Cluster %s: %s", d.Id(), err)
	}
	if err = d.Set("deployment_type", resp.MasterConfig.DeploymentType); err != nil {
		return diag.Errorf("error setting deployment_type for Karbon Cluster %s: %s", d.Id(), err)
	}
	if err = d.Set("worker_node_pool", flattenedWorkerNodepool); err != nil {
		return diag.Errorf("error setting worker_node_pool for Karbon Cluster %s: %s", d.Id(), err)
	}
	if err = d.Set("etcd_node_pool", flattenedEtcdNodepool); err != nil {
		return diag.Errorf("error setting etcd_node_pool for Karbon Cluster %s: %s", d.Id(), err)
	}
	if err = d.Set("master_node_pool", flattenedMasterNodepool); err != nil {
		return diag.Errorf("error setting master_node_pool for Karbon Cluster %s: %s", d.Id(), err)
	}
	if err = d.Set("uuid", utils.StringValue(resp.UUID)); err != nil {
		return diag.Errorf("error setting uuid for Karbon Cluster %s: %s", d.Id(), err)
	}
	d.SetId(*resp.UUID)

	return nil
}

func KarbonClusterDataSourceMap() map[string]*schema.Schema {
	kcsm := KarbonClusterElementDataSourceMap()
	kcsm["karbon_cluster_id"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"karbon_cluster_name"},
	}
	kcsm["karbon_cluster_name"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"karbon_cluster_id"},
	}
	return kcsm
}

func KarbonClusterElementDataSourceMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"uuid": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deployment_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"kubeapi_server_ipv4_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"etcd_node_pool":   nodePoolDatasourceSchema(),
		"master_node_pool": nodePoolDatasourceSchema(),
		"worker_node_pool": nodePoolDatasourceSchema(),
	}
}

func nodePoolDatasourceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"node_os_version": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"num_instances": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"ahv_config": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"cpu": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"disk_mib": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"memory_mib": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"network_uuid": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"prism_element_cluster_uuid": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"nodes": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"hostname": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"ipv4_address": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}
