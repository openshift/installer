package alicloud

import (
	"encoding/json"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCRNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCRNamespacesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"namespaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_create": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"default_visibility": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudCRNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}
	invoker := NewInvoker()

	var (
		request  *cr.GetNamespaceListRequest
		response *cr.GetNamespaceListResponse
	)
	if err := invoker.Run(func() error {
		raw, err := client.WithCrClient(func(crClient *cr.Client) (interface{}, error) {
			request = cr.CreateGetNamespaceListRequest()
			return crClient.GetNamespaceList(request)
		})
		response, _ = raw.(*cr.GetNamespaceListResponse)
		return err
	}); err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cr_namespaces", "GetNamespaceList", AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), response)
	var crResp crDescribeNamespaceListResponse
	err := json.Unmarshal(response.GetHttpContentBytes(), &crResp)
	if err != nil {
		return WrapError(err)
	}

	var names []string
	var s []map[string]interface{}

	for _, ns := range crResp.Data.Namespace {
		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r, err := regexp.Compile(nameRegex.(string))
			if err != nil {
				return WrapError(err)
			}
			if !r.MatchString(ns.Namespace) {
				continue
			}
		}

		mapping := map[string]interface{}{
			"name": ns.Namespace,
		}

		raw, err := crService.DescribeCrNamespace(ns.Namespace)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		var resp crDescribeNamespaceResponse
		err = json.Unmarshal(raw.GetHttpContentBytes(), &resp)
		if err != nil {
			return WrapError(err)
		}

		mapping["auto_create"] = resp.Data.Namespace.AutoCreate
		mapping["default_visibility"] = resp.Data.Namespace.DefaultVisibility

		names = append(names, ns.Namespace)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("namespaces", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
