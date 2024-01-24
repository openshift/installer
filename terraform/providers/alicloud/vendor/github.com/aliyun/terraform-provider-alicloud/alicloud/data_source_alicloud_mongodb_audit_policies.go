package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudMongodbAuditPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMongodbAuditPoliciesRead,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"audit_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMongodbAuditPoliciesRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	MongoDBService := MongoDBService{client}
	dbInstanceId := d.Get("db_instance_id")
	object, err := MongoDBService.DescribeMongodbAuditPolicy(dbInstanceId.(string))
	if err != nil {
		if NotFoundError(err) {
			d.SetId("MongodbAuditPolicy")
			return nil
		}
		return WrapError(err)
	}

	s := make([]map[string]interface{}, 0)
	mapping := map[string]interface{}{
		"id":             fmt.Sprint(object["DBInstanceId"]),
		"db_instance_id": fmt.Sprint(object["DBInstanceId"]),
		"audit_status":   convertMongodbAuditPolicyResponse(object["LogAuditStatus"].(string)),
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
