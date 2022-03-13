package alicloud

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpnConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpnConnectionsRead,

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

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"customer_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
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
			},

			// Computed values
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_gateway_id": {
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
						"remote_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"effect_immediately": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"ipsec_config": {
							Type:     schema.TypeList,
							Optional: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipsec_enc_alg": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipsec_auth_alg": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipsec_pfs": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ipsec_lifetime": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},

						"ike_config": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"psk": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_version": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_mode": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_enc_alg": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_auth_alg": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_pfs": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_lifetime": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"ike_local_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ike_remote_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpnConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := vpc.CreateDescribeVpnConnectionsRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var allVpnConns []vpc.VpnConnection

	if v, ok := d.GetOk("vpn_gateway_id"); ok && v.(string) != "" {
		request.VpnGatewayId = v.(string)
	}

	if v, ok := d.GetOk("customer_gateway_id"); ok && v.(string) != "" {
		request.CustomerGatewayId = v.(string)
	}

	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpnConnections(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpn_connections", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVpnConnectionsResponse)

		if len(response.VpnConnections.VpnConnection) < 1 {
			break
		}

		allVpnConns = append(allVpnConns, response.VpnConnections.VpnConnection...)

		if len(response.VpnConnections.VpnConnection) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredVpnConns []vpc.VpnConnection
	var reg *regexp.Regexp
	var ids []string
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, item := range v.([]interface{}) {
			if item == nil {
				continue
			}
			ids = append(ids, strings.Trim(item.(string), " "))
		}
	}
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		if r, err := regexp.Compile(nameRegex.(string)); err == nil {
			reg = r
		} else {
			return WrapError(err)
		}
	}

	for _, vpnConn := range allVpnConns {
		if reg != nil {
			if !reg.MatchString(vpnConn.Name) {
				continue
			}
		}
		if ids != nil && len(ids) != 0 {
			for _, id := range ids {
				if vpnConn.VpnConnectionId == id {
					filteredVpnConns = append(filteredVpnConns, vpnConn)
				}
			}
		} else {
			filteredVpnConns = append(filteredVpnConns, vpnConn)
		}
	}

	return vpnConnectionsDecriptionAttributes(d, filteredVpnConns, meta)
}

func vpnConnectionsDecriptionAttributes(d *schema.ResourceData, vpnSetTypes []vpc.VpnConnection, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, conn := range vpnSetTypes {
		mapping := map[string]interface{}{
			"customer_gateway_id": conn.CustomerGatewayId,
			"vpn_gateway_id":      conn.VpnGatewayId,
			"id":                  conn.VpnConnectionId,
			"name":                conn.Name,
			"local_subnet":        conn.LocalSubnet,
			"remote_subnet":       conn.RemoteSubnet,
			"create_time":         TimestampToStr(conn.CreateTime),
			"effect_immediately":  conn.EffectImmediately,
			"status":              conn.Status,
			"ike_config":          vpnGatewayService.ParseIkeConfig(conn.IkeConfig),
			"ipsec_config":        vpnGatewayService.ParseIpsecConfig(conn.IpsecConfig),
		}
		ids = append(ids, conn.VpnConnectionId)
		names = append(names, conn.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("connections", s); err != nil {
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
