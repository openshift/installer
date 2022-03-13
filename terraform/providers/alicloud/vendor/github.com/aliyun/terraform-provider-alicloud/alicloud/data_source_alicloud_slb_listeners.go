package alicloud

import (
	"strconv"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudSlbListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbListenersRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"frontend_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description_regex": {
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
			"slb_listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frontend_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backend_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_slave_server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"persistence_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"established_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sticky_session": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sticky_session_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cookie_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cookie": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_connect_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_connect_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_http_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gzip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ca_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x_forwarded_for": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x_forwarded_for_slb_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x_forwarded_for_slb_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"x_forwarded_for_slb_proto": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// http https
						"idle_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						// http https
						"request_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						// https
						"enable_http2": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// https
						"tls_cipher_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSlbListenersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerId = d.Get("load_balancer_id").(string)

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_listeners", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeLoadBalancerAttributeResponse)
	var filteredListenersTemp []slb.ListenerPortAndProtocol
	port := -1
	if v, ok := d.GetOk("frontend_port"); ok && v.(int) != 0 {
		port = v.(int)
	}
	protocol := ""
	if v, ok := d.GetOk("protocol"); ok && v.(string) != "" {
		protocol = v.(string)
	}
	var r *regexp.Regexp
	if despRegex, ok := d.GetOk("description_regex"); ok && despRegex.(string) != "" {
		r = regexp.MustCompile(despRegex.(string))
	}
	if port != -1 || protocol != "" || r != nil {
		for _, listener := range response.ListenerPortsAndProtocol.ListenerPortAndProtocol {
			if port != -1 && listener.ListenerPort != port {
				continue
			}
			if protocol != "" && listener.ListenerProtocol != protocol {
				continue
			}
			if r != nil && !r.MatchString(listener.Description) {
				continue
			}

			filteredListenersTemp = append(filteredListenersTemp, listener)
		}
	} else {
		filteredListenersTemp = response.ListenerPortsAndProtocol.ListenerPortAndProtocol
	}

	return slbListenersDescriptionAttributes(d, filteredListenersTemp, meta)
}

func slbListenersDescriptionAttributes(d *schema.ResourceData, listeners []slb.ListenerPortAndProtocol, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var ids []string
	var s []map[string]interface{}

	for _, listener := range listeners {
		mapping := map[string]interface{}{
			"frontend_port": listener.ListenerPort,
			"protocol":      listener.ListenerProtocol,
			"description":   listener.Description,
		}

		loadBalancerId := d.Get("load_balancer_id").(string)
		switch Protocol(listener.ListenerProtocol) {
		case Http:
			request := slb.CreateDescribeLoadBalancerHTTPListenerAttributeRequest()
			request.LoadBalancerId = loadBalancerId
			request.ListenerPort = requests.NewInteger(listener.ListenerPort)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.DescribeLoadBalancerHTTPListenerAttribute(request)
			})
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_listeners", request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ := raw.(*slb.DescribeLoadBalancerHTTPListenerAttributeResponse)
			mapping["backend_port"] = response.BackendServerPort
			mapping["status"] = response.Status
			mapping["bandwidth"] = response.Bandwidth
			mapping["scheduler"] = response.Scheduler
			mapping["server_group_id"] = response.VServerGroupId
			mapping["sticky_session"] = response.StickySession
			mapping["sticky_session_type"] = response.StickySessionType
			mapping["cookie_timeout"] = response.CookieTimeout
			mapping["cookie"] = response.Cookie
			mapping["health_check"] = response.HealthCheck
			mapping["health_check_domain"] = response.HealthCheckDomain
			mapping["health_check_uri"] = response.HealthCheckURI
			mapping["health_check_connect_port"] = response.HealthCheckConnectPort
			mapping["healthy_threshold"] = response.HealthyThreshold
			mapping["unhealthy_threshold"] = response.UnhealthyThreshold
			mapping["health_check_timeout"] = response.HealthCheckTimeout
			mapping["health_check_interval"] = response.HealthCheckInterval
			mapping["health_check_http_code"] = response.HealthCheckHttpCode
			mapping["gzip"] = response.Gzip
			mapping["x_forwarded_for"] = response.XForwardedFor
			mapping["x_forwarded_for_slb_ip"] = response.XForwardedForSLBIP
			mapping["x_forwarded_for_slb_id"] = response.XForwardedForSLBID
			mapping["x_forwarded_for_slb_proto"] = response.XForwardedForProto
			mapping["idle_timeout"] = response.IdleTimeout
			mapping["request_timeout"] = response.RequestTimeout
		case Https:
			request := slb.CreateDescribeLoadBalancerHTTPSListenerAttributeRequest()
			request.LoadBalancerId = loadBalancerId
			request.ListenerPort = requests.NewInteger(listener.ListenerPort)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.DescribeLoadBalancerHTTPSListenerAttribute(request)
			})
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_listeners", request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ := raw.(*slb.DescribeLoadBalancerHTTPSListenerAttributeResponse)
			mapping["backend_port"] = response.BackendServerPort
			mapping["status"] = response.Status
			mapping["security_status"] = response.SecurityStatus
			mapping["bandwidth"] = response.Bandwidth
			mapping["scheduler"] = response.Scheduler
			mapping["server_group_id"] = response.VServerGroupId
			mapping["sticky_session"] = response.StickySession
			mapping["sticky_session_type"] = response.StickySessionType
			mapping["cookie_timeout"] = response.CookieTimeout
			mapping["cookie"] = response.Cookie
			mapping["health_check"] = response.HealthCheck
			mapping["health_check_domain"] = response.HealthCheckDomain
			mapping["health_check_uri"] = response.HealthCheckURI
			mapping["health_check_connect_port"] = response.HealthCheckConnectPort
			mapping["healthy_threshold"] = response.HealthyThreshold
			mapping["unhealthy_threshold"] = response.UnhealthyThreshold
			mapping["health_check_timeout"] = response.HealthCheckTimeout
			mapping["health_check_interval"] = response.HealthCheckInterval
			mapping["health_check_http_code"] = response.HealthCheckHttpCode
			mapping["gzip"] = response.Gzip
			mapping["ssl_certificate_id"] = response.ServerCertificateId
			mapping["server_certificate_id"] = response.ServerCertificateId
			mapping["ca_certificate_id"] = response.CACertificateId
			mapping["x_forwarded_for"] = response.XForwardedFor
			mapping["x_forwarded_for_slb_ip"] = response.XForwardedForSLBIP
			mapping["x_forwarded_for_slb_id"] = response.XForwardedForSLBID
			mapping["x_forwarded_for_slb_proto"] = response.XForwardedForProto
			mapping["idle_timeout"] = response.IdleTimeout
			mapping["request_timeout"] = response.RequestTimeout
			mapping["enable_http2"] = response.EnableHttp2
			mapping["tls_cipher_policy"] = response.TLSCipherPolicy
		case Tcp:
			request := slb.CreateDescribeLoadBalancerTCPListenerAttributeRequest()
			request.LoadBalancerId = loadBalancerId
			request.ListenerPort = requests.NewInteger(listener.ListenerPort)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.DescribeLoadBalancerTCPListenerAttribute(request)
			})
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_listeners", request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ := raw.(*slb.DescribeLoadBalancerTCPListenerAttributeResponse)
			mapping["backend_port"] = response.BackendServerPort
			mapping["status"] = response.Status
			mapping["bandwidth"] = response.Bandwidth
			mapping["scheduler"] = response.Scheduler
			mapping["server_group_id"] = response.VServerGroupId
			mapping["master_slave_server_group_id"] = response.MasterSlaveServerGroupId
			mapping["persistence_timeout"] = response.PersistenceTimeout
			mapping["established_timeout"] = response.EstablishedTimeout
			mapping["health_check"] = response.HealthCheck
			mapping["health_check_type"] = response.HealthCheckType
			mapping["health_check_domain"] = response.HealthCheckDomain
			mapping["health_check_uri"] = response.HealthCheckURI
			mapping["health_check_connect_port"] = response.HealthCheckConnectPort
			mapping["health_check_connect_timeout"] = response.HealthCheckConnectTimeout
			mapping["healthy_threshold"] = response.HealthyThreshold
			mapping["unhealthy_threshold"] = response.UnhealthyThreshold
			mapping["health_check_interval"] = response.HealthCheckInterval
			mapping["health_check_http_code"] = response.HealthCheckHttpCode
		case Udp:
			request := slb.CreateDescribeLoadBalancerUDPListenerAttributeRequest()
			request.LoadBalancerId = loadBalancerId
			request.ListenerPort = requests.NewInteger(listener.ListenerPort)
			raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
				return slbClient.DescribeLoadBalancerUDPListenerAttribute(request)
			})
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_listeners", request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ := raw.(*slb.DescribeLoadBalancerUDPListenerAttributeResponse)
			mapping["backend_port"] = response.BackendServerPort
			mapping["status"] = response.Status
			mapping["bandwidth"] = response.Bandwidth
			mapping["scheduler"] = response.Scheduler
			mapping["server_group_id"] = response.VServerGroupId
			mapping["master_slave_server_group_id"] = response.MasterSlaveServerGroupId
			mapping["persistence_timeout"] = response.PersistenceTimeout
			mapping["health_check"] = response.HealthCheck
			mapping["health_check_connect_port"] = response.HealthCheckConnectPort
			mapping["health_check_connect_timeout"] = response.HealthCheckConnectTimeout
			mapping["healthy_threshold"] = response.HealthyThreshold
			mapping["unhealthy_threshold"] = response.UnhealthyThreshold
			mapping["health_check_interval"] = response.HealthCheckInterval
		}

		ids = append(ids, strconv.Itoa(listener.ListenerPort))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slb_listeners", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
