package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudSslVpnServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSslVpnServersRead,

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

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_ip_pool": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cipher": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proto": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"compress": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"internet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSslVpnServersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateDescribeSslVpnServersRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var allSslVpnServers []vpc.SslVpnServer

	if v, ok := d.GetOk("vpn_gateway_id"); ok && v.(string) != "" {
		request.VpnGatewayId = v.(string)
	}

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeSslVpnServers(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ssl_vpn_servers", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeSslVpnServersResponse)
		if len(response.SslVpnServers.SslVpnServer) < 1 {
			break
		}
		allSslVpnServers = append(allSslVpnServers, response.SslVpnServers.SslVpnServer...)
		if len(response.SslVpnServers.SslVpnServer) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredSslVpnServers []vpc.SslVpnServer
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

	for _, sslVpnServer := range allSslVpnServers {
		if reg != nil {
			if !reg.MatchString(sslVpnServer.Name) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[sslVpnServer.SslVpnServerId]; !ok {
				continue
			}
		}

		filteredSslVpnServers = append(filteredSslVpnServers, sslVpnServer)
	}

	return sslVpnServersDecriptionAttributes(d, filteredSslVpnServers, meta)
}

func sslVpnServersDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.SslVpnServer, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, vpn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"id":              vpn.SslVpnServerId,
			"vpn_gateway_id":  vpn.VpnGatewayId,
			"local_subnet":    vpn.LocalSubnet,
			"name":            vpn.Name,
			"client_ip_pool":  vpn.ClientIpPool,
			"create_time":     TimestampToStr(vpn.CreateTime),
			"cipher":          vpn.Cipher,
			"proto":           vpn.Proto,
			"port":            vpn.Port,
			"compress":        vpn.Compress,
			"connections":     vpn.Connections,
			"max_connections": vpn.MaxConnections,
			"internet_ip":     vpn.InternetIp,
		}
		ids = append(ids, vpn.SslVpnServerId)
		names = append(names, vpn.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("servers", s); err != nil {
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
