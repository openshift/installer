package alicloud

import (
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudKvstoreConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKvstoreConnectionsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 1,
				MinItems: 1,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"upgradeable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudKvstoreConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := r_kvstore.CreateDescribeDBInstanceNetInfoRequest()
	var objects []r_kvstore.InstanceNetInfo

	var response *r_kvstore.DescribeDBInstanceNetInfoResponse
	for _, id := range d.Get("ids").([]interface{}) {
		request.InstanceId = id.(string)
	}
	raw, err := client.WithRKvstoreClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeDBInstanceNetInfo(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kvstore_connections", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ = raw.(*r_kvstore.DescribeDBInstanceNetInfoResponse)

	for _, item := range response.NetInfoItems.InstanceNetInfo {
		if item.DBInstanceNetType != "0" {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"connection_string":    object.ConnectionString,
			"db_instance_net_type": object.DBInstanceNetType,
			"expired_time":         object.ExpiredTime,
			"ip_address":           object.IPAddress,
			"id":                   request.InstanceId,
			"instance_id":          request.InstanceId,
			"port":                 object.Port,
			"upgradeable":          object.Upgradeable,
			"vpc_id":               object.VPCId,
			"vpc_instance_id":      object.VPCInstanceId,
			"vswitch_id":           object.VSwitchId,
		}
		ids = append(ids, request.InstanceId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("connections", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
