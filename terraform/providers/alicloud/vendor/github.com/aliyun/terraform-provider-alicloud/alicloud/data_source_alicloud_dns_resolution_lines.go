package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// https://help.aliyun.com/document_detail/69017.html
func dataSourceAlicloudDnsResolutionLines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDnsResolutionLinesRead,

		Schema: map[string]*schema.Schema{
			"line_codes": {
				Optional: true,
				Type:     schema.TypeList,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"line_display_names": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"line_names": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_client_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// computed values
			"lines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"line_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDnsResolutionLinesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidns.CreateDescribeSupportLinesRequest()

	request.RegionId = client.RegionId
	request.Lang = d.Get("lang").(string)
	request.UserClientIp = d.Get("user_client_ip").(string)
	request.DomainName = d.Get("domain_name").(string)

	raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.DescribeSupportLines(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dns_support_lines", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeSupportLinesResponse)
	allLines := response.RecordLines.RecordLine

	var filteredLines []alidns.RecordLine
	for _, line := range allLines {
		if v, ok := d.GetOk("line_codes"); ok {
			lineCodes, _ := v.([]interface{})
			if !hasEqualString(lineCodes, line.LineCode) {
				continue
			}
		}
		if v, ok := d.GetOk("line_display_names"); ok {
			lineDisplayNames, _ := v.([]interface{})
			if !hasEqualString(lineDisplayNames, line.LineDisplayName) {
				continue
			}
		}
		if v, ok := d.GetOk("line_names"); ok {
			lineNames, _ := v.([]interface{})
			if !hasEqualString(lineNames, line.LineName) {
				continue
			}
		}
		filteredLines = append(filteredLines, line)
	}

	return resolutionLinesDecriptionAttributes(d, filteredLines, meta)
}

func resolutionLinesDecriptionAttributes(d *schema.ResourceData, Lines []alidns.RecordLine, meta interface{}) error {
	var s []map[string]interface{}
	var lineCodes []string
	var lineDisplayNames []string

	for _, line := range Lines {
		mapping := map[string]interface{}{
			"line_code":         line.LineCode,
			"line_display_name": line.LineDisplayName,
			"line_name":         line.LineName,
		}
		lineCodes = append(lineCodes, line.LineCode)
		lineDisplayNames = append(lineDisplayNames, line.LineDisplayName)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(lineCodes))

	if err := d.Set("lines", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("line_codes", lineCodes); err != nil {
		return WrapError(err)
	}

	if err := d.Set("line_display_names", lineDisplayNames); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

// to check if a slice has a specific string
func hasEqualString(s []interface{}, a string) bool {
	var isEqual bool
	for _, code := range s {
		fatherCode, _ := code.(string)
		if fatherCode == a {
			isEqual = true
			return isEqual
		}
	}
	return isEqual
}
