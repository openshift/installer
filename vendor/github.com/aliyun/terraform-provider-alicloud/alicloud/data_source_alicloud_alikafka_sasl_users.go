package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAlikafkaSaslUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlikafkaSaslUsersRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlikafkaSaslUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	request := alikafka.CreateDescribeSaslUsersRequest()
	request.InstanceId = d.Get("instance_id").(string)
	request.RegionId = client.RegionId

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.DescribeSaslUsers(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alikafka_sasl_users", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alikafka.DescribeSaslUsersResponse)

	var filteredSaslUsers []alikafka.SaslUserVO
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r, err = regexp.Compile(nameRegex.(string))
			if err != nil {
				return WrapError(err)
			}
		}
		for _, saslUser := range response.SaslUserList.SaslUserVO {
			if r != nil && !r.MatchString(saslUser.Username) {
				continue
			}

			filteredSaslUsers = append(filteredSaslUsers, saslUser)
		}
	} else {
		filteredSaslUsers = response.SaslUserList.SaslUserVO
	}
	return alikafkaSaslUsersDecriptionAttributes(d, filteredSaslUsers, meta)
}

func alikafkaSaslUsersDecriptionAttributes(d *schema.ResourceData, saslUsersInfo []alikafka.SaslUserVO, meta interface{}) error {

	var names []string
	var s []map[string]interface{}

	for _, item := range saslUsersInfo {
		mapping := map[string]interface{}{
			"username": item.Username,
			"password": item.Password,
		}

		names = append(names, item.Username)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(names))

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("users", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
