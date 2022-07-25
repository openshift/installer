package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudSslVpnClientCerts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSslVpnClientCertsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},

			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			},

			"ssl_vpn_server_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"certs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ssl_vpn_server_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeInt,
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

func dataSourceAlicloudSslVpnClientCertsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateDescribeSslVpnClientCertsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var allSslVpnClientCerts []vpc.SslVpnClientCertKey

	if v, ok := d.GetOk("ssl_vpn_server_id"); ok && v.(string) != "" {
		request.SslVpnServerId = v.(string)
	}

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeSslVpnClientCerts(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ssl_vpn_client_certs", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeSslVpnClientCertsResponse)
		if len(response.SslVpnClientCertKeys.SslVpnClientCertKey) < 1 {
			break
		}
		allSslVpnClientCerts = append(allSslVpnClientCerts, response.SslVpnClientCertKeys.SslVpnClientCertKey...)
		if len(response.SslVpnClientCertKeys.SslVpnClientCertKey) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredSslVpnClientCerts []vpc.SslVpnClientCertKey
	var reg *regexp.Regexp

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		if r, err := regexp.Compile(nameRegex.(string)); err == nil {
			reg = r
		} else {
			return WrapError(err)
		}
	}

	for _, sslVpnClientCertKey := range allSslVpnClientCerts {
		if reg != nil {
			if !reg.MatchString(sslVpnClientCertKey.Name) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[sslVpnClientCertKey.SslVpnClientCertId]; !ok {
				continue
			}
		}

		filteredSslVpnClientCerts = append(filteredSslVpnClientCerts, sslVpnClientCertKey)
	}

	return sslVpnClientCertsDecriptionAttributes(d, filteredSslVpnClientCerts, meta)
}

func sslVpnClientCertsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.SslVpnClientCertKey, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"ssl_vpn_server_id": vpn.SslVpnServerId,
			"id":                vpn.SslVpnClientCertId,
			"name":              vpn.Name,
			"end_time":          vpn.EndTime,
			"create_time":       TimestampToStr(vpn.CreateTime),
			"status":            vpn.Status,
		}
		ids = append(ids, vpn.SslVpnClientCertId)
		names = append(names, vpn.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("certs", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
