package alicloud

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCSServerlessKubernetesClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCSServerlessKubernetesClustersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_public_access_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"deletion_protection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"nat_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connections": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_server_internet": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"api_server_intranet": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"master_public_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tags": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
						}},
				},
			},
		},
	}
}

func dataSourceAlicloudCSServerlessKubernetesClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var allClusterTypes []*cs.ServerlessClusterResponse

	var requestInfo *cs.Client
	invoker := NewInvoker()
	var response interface{}
	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeServerlessKubernetesClusters()
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cs_serverless_kubernetes_clusters", "DescribeServerlessKubernetesClusters", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		addDebug("DescribeServerlessKubernetesClusters", response, requestInfo, requestMap)
	}
	allClusterTypes, _ = response.([]*cs.ServerlessClusterResponse)

	var filteredClusterTypes []*cs.ServerlessClusterResponse
	for _, v := range allClusterTypes {
		if v.ClusterType != cs.ClusterTypeServerlessKubernetes {
			continue
		}
		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r, err := regexp.Compile(nameRegex.(string))
			if err != nil {
				return WrapError(err)
			}
			if !r.MatchString(v.Name) {
				continue
			}
		}
		if ids, ok := d.GetOk("ids"); ok {
			var found bool
			for _, i := range expandStringList(ids.([]interface{})) {
				if v.ClusterId == i {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		filteredClusterTypes = append(filteredClusterTypes, v)
	}

	var filteredKubernetesCluster []*cs.ServerlessClusterResponse

	for _, v := range filteredClusterTypes {
		var serverlessCluster *cs.ServerlessClusterResponse

		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return csClient.DescribeServerlessKubernetesCluster(v.ClusterId)
			})
			response = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cs_serverless_kubernetes_clusters", "DescribeServerlessKubernetesCluster", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["Id"] = v.ClusterId
			addDebug("DescribeServerlessKubernetesCluster", response, requestInfo, requestMap)
		}
		serverlessCluster = response.(*cs.ServerlessClusterResponse)

		filteredKubernetesCluster = append(filteredKubernetesCluster, serverlessCluster)
	}
	return csServerlessKubernetesClusterDescriptionAttributes(d, filteredClusterTypes, meta)
}

func csServerlessKubernetesClusterDescriptionAttributes(d *schema.ResourceData, clusters []*cs.ServerlessClusterResponse, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}

	var ids, names []string
	var s []map[string]interface{}
	for _, ct := range clusters {
		mapping := map[string]interface{}{
			"id":   ct.ClusterId,
			"name": ct.Name,
		}

		if detailedEnabled, ok := d.GetOk("enable_details"); ok && !detailedEnabled.(bool) {
			ids = append(ids, ct.ClusterId)
			names = append(names, ct.Name)
			s = append(s, mapping)
			continue
		}

		mapping["vpc_id"] = ct.VpcId
		mapping["vswitch_id"] = ct.VSwitchId
		mapping["security_group_id"] = ct.SecurityGroupId
		mapping["deletion_protection"] = ct.DeletionProtection
		mapping["tags"] = csService.tagsToMap(ct.Tags)
		//set default value
		mapping["endpoint_public_access_enabled"] = false

		invoker := NewInvoker()
		client := meta.(*connectivity.AliyunClient)

		var response interface{}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				endpoints, err := csClient.GetClusterEndpoints(ct.ClusterId)
				return endpoints, err
			})
			response = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cs_serverless_kubernetes_clusters", "GetClusterEndpoints", DenverdinoAliyungo)
		}
		connection := make(map[string]string)
		if endpoints, ok := response.(cs.ClusterEndpoints); ok && endpoints.ApiServerEndpoint != "" {
			//set public access
			mapping["endpoint_public_access_enabled"] = true

			connection["api_server_internet"] = endpoints.ApiServerEndpoint
			connection["master_public_ip"] = strings.TrimSuffix(strings.TrimPrefix(endpoints.ApiServerEndpoint, "https://"), ":6443")
		}
		if endpoints, ok := response.(cs.ClusterEndpoints); ok && endpoints.IntranetApiServerEndpoint != "" {
			connection["api_server_intranet"] = endpoints.IntranetApiServerEndpoint
		}

		mapping["connections"] = connection

		request := vpc.CreateDescribeNatGatewaysRequest()
		request.VpcId = ct.VpcId

		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNatGateways(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, ct.VpcId, "DescribeNatGateways", AlibabaCloudSdkGoERROR)
		}
		addDebug("DescribeNatGateways", raw, request.RpcRequest, request)
		nat, _ := raw.(*vpc.DescribeNatGatewaysResponse)
		if nat != nil && len(nat.NatGateways.NatGateway) > 0 {
			mapping["nat_gateway_id"] = nat.NatGateways.NatGateway[0].NatGatewayId
		}

		ids = append(ids, ct.ClusterId)
		names = append(names, ct.Name)
		s = append(s, mapping)
	}

	_ = d.Set("ids", ids)
	_ = d.Set("names", names)
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		_ = writeToFile(output.(string), s)
	}

	return nil
}
