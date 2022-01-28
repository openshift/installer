package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAlidnsRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlidnsRecordsRead,
		Schema: map[string]*schema.Schema{
			"direction": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"key_word": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"line": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "default",
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rr_key_word": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"search_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENABLE", "DISABLE"}, false),
				Default:      "ENABLE",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "NS", "MX", "TXT", "CNAME", "SRV", "AAAA", "CAA", "REDIRECT_URL", "FORWORD_URL"}, false),
			},
			"type_key_word": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"value_key_word": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rr_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"value_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"locked": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlidnsRecordsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateDescribeDomainRecordsRequest()
	if v, ok := d.GetOk("direction"); ok {
		request.Direction = v.(string)
	}
	request.DomainName = d.Get("domain_name").(string)
	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("key_word"); ok {
		request.KeyWord = v.(string)
	}
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	if v, ok := d.GetOk("line"); ok {
		request.Line = v.(string)
	}
	if v, ok := d.GetOk("order_by"); ok {
		request.OrderBy = v.(string)
	}
	if v, ok := d.GetOk("rr_key_word"); ok {
		request.RRKeyWord = v.(string)
	}
	if v, ok := d.GetOk("search_mode"); ok {
		request.SearchMode = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		request.Status = v.(string)
	}
	if v, ok := d.GetOk("type"); ok {
		request.Type = v.(string)
	}
	if v, ok := d.GetOk("type_key_word"); ok {
		request.TypeKeyWord = v.(string)
	}
	if v, ok := d.GetOk("value_key_word"); ok {
		request.ValueKeyWord = v.(string)
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []alidns.Record
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var rrRegex *regexp.Regexp
	if v, ok := d.GetOk("rr_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		rrRegex = r
	}
	var valueRegex *regexp.Regexp
	if v, ok := d.GetOk("value_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		valueRegex = r
	}
	for {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.DescribeDomainRecords(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alidns_records", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*alidns.DescribeDomainRecordsResponse)

		for _, item := range response.DomainRecords.Record {
			if rrRegex != nil {
				if !rrRegex.MatchString(item.RR) {
					continue
				}
			}
			if valueRegex != nil {
				if !valueRegex.MatchString(item.Value) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.RecordId]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.DomainRecords.Record) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))

	for i, object := range objects {
		mapping := map[string]interface{}{
			"domain_name": object.DomainName,
			"line":        object.Line,
			"locked":      object.Locked,
			"priority":    object.Priority,
			"rr":          object.RR,
			"id":          object.RecordId,
			"record_id":   object.RecordId,
			"status":      object.Status,
			"ttl":         object.TTL,
			"type":        object.Type,
			"value":       object.Value,
			"remark":      object.Remark,
		}
		ids[i] = object.RecordId
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("records", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
