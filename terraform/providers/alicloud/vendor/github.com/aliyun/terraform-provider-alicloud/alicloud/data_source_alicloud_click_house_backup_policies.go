package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudClickHouseBackupPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudClickHouseBackupPoliciesRead,
		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_retention_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"preferred_backup_period": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"preferred_backup_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudClickHouseBackupPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbClusterId := d.Get("db_cluster_id")
	clickhouseService := ClickhouseService{client}
	object, err := clickhouseService.DescribeClickHouseBackupPolicy(dbClusterId.(string))
	if err != nil {
		if NotFoundError(err) {
			d.SetId("ClickHouseBackupPolicy")
			return nil
		}
		return WrapError(err)
	}
	s := make([]map[string]interface{}, 0)
	mapping := map[string]interface{}{
		"backup_retention_period": formatInt(object["BackupRetentionPeriod"]),
		"id":                      fmt.Sprint(dbClusterId),
		"db_cluster_id":           fmt.Sprint(dbClusterId),
		"preferred_backup_period": strings.Split(object["PreferredBackupPeriod"].(string), ","),
		"preferred_backup_time":   object["PreferredBackupTime"],
		"status":                  fmt.Sprint(object["Switch"]),
	}
	s = append(s, mapping)

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	if err := d.Set("policies", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
