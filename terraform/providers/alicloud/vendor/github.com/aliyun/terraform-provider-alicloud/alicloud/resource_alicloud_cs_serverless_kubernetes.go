package alicloud

import (
	"fmt"
	"regexp"
	"time"

	roacs "github.com/alibabacloud-go/cs-20151215/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCSServerlessKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSServerlessKubernetesCreate,
		Read:   resourceAlicloudCSServerlessKubernetesRead,
		Update: resourceAlicloudCSServerlessKubernetesUpdate,
		Delete: resourceAlicloudCSServerlessKubernetesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(1, 63),
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(0, 37),
				ConflictsWith: []string{"name"},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'vswitch_id' has been deprecated from provider version 1.91.0. New field 'vswitch_ids' replace it.",
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems:      1,
				ConflictsWith: []string{"vswitch_id"},
			},
			"service_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"new_nat_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"private_zone": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"service_discovery_types"},
				Deprecated:    "Field 'private_zone' has been deprecated from provider version 1.123.1. New field 'service_discovery_types' replace it.",
			},
			"service_discovery_types": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"CoreDNS", "PrivateZone"}, false),
				},
				ConflictsWith: []string{"private_zone"},
			},
			"zone_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"endpoint_public_access_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"kube_config": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"client_cert": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"client_key": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"cluster_ca_cert": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"addons": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"load_balancer_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"slb.s1.small", "slb.s2.small", "slb.s2.medium", "slb.s3.small", "slb.s3.medium", "slb.s3.large"}, false),
				Default:      "slb.s1.small",
			},
			"logging_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "SLS",
			},
			"sls_project_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"retain_resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAlicloudCSServerlessKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()

	csService := CsService{client}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
	}

	tags := make([]cs.Tag, 0)
	tagsMap, ok := d.Get("tags").(map[string]interface{})
	if ok {
		for key, value := range tagsMap {
			if value != nil {
				if v, ok := value.(string); ok {
					tags = append(tags, cs.Tag{
						Key:   key,
						Value: v,
					})
				}
			}
		}
	}

	addons := make([]cs.Addon, 0)
	if v, ok := d.GetOk("addons"); ok {
		all, ok := v.([]interface{})
		if ok {
			for _, a := range all {
				addon, ok := a.(map[string]interface{})
				if ok {
					addons = append(addons, cs.Addon{
						Name:     addon["name"].(string),
						Config:   addon["config"].(string),
						Disabled: addon["disabled"].(bool),
					})
				}
			}
		}
	}

	args := &cs.ServerlessCreationArgs{
		Name:                 clusterName,
		ClusterType:          cs.ClusterTypeServerlessKubernetes,
		RegionId:             client.RegionId,
		VpcId:                d.Get("vpc_id").(string),
		EndpointPublicAccess: d.Get("endpoint_public_access_enabled").(bool),
		NatGateway:           d.Get("new_nat_gateway").(bool),
		SecurityGroupId:      d.Get("security_group_id").(string),
		Addons:               addons,
		KubernetesVersion:    d.Get("version").(string),
		DeletionProtection:   d.Get("deletion_protection").(bool),
		ResourceGroupId:      d.Get("resource_group_id").(string),
	}

	if v, ok := d.GetOk("time_zone"); ok {
		args.TimeZone = v.(string)
	}

	if v, ok := d.GetOk("zone_id"); ok {
		args.ZoneID = v.(string)
	}

	if v, ok := d.GetOk("service_cidr"); ok {
		args.ServiceCIDR = v.(string)
	}

	if v, ok := d.GetOk("logging_type"); ok {
		args.LoggingType = v.(string)
	}

	if v, ok := d.GetOk("sls_project_name"); ok {
		args.SLSProjectName = v.(string)
	}

	if v, ok := d.GetOk("service_discovery_types"); ok {
		args.ServiceDiscoveryTypes = expandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOkExists("private_zone"); ok {
		args.ServiceDiscoveryTypes = []string{}
		if v.(bool) == true {
			args.ServiceDiscoveryTypes = []string{"PrivateZone"}
		}
	}

	if v := d.Get("vswitch_id").(string); v != "" {
		args.VSwitchId = v
	}

	if v, ok := d.GetOk("vswitch_ids"); ok {
		args.VswitchIds = expandStringList(v.([]interface{}))
	}

	if lbSpec, ok := d.GetOk("load_balancer_spec"); ok {
		args.LoadBalancerSpec = lbSpec.(string)
	}

	//set tags
	if len(tags) > 0 {
		args.Tags = tags
	}

	var requestInfo *cs.Client
	var response interface{}
	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.CreateServerlessKubernetesCluster(args)
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_serverless_kubernetes", "CreateServerlessKubernetesCluster", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Args"] = args
		addDebug("CreateServerlessKubernetesCluster", response, requestInfo, requestMap)
	}
	cluster, _ := response.(*cs.ClusterCommonResponse)
	d.SetId(cluster.ClusterID)

	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, csService.CsServerlessKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCSServerlessKubernetesRead(d, meta)
}

func resourceAlicloudCSServerlessKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()
	rosClient, err := client.NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	object, err := csService.DescribeCsServerlessKubernetes(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	vswitchIds := []string{}
	resources, _ := rosClient.DescribeClusterResources(tea.String(d.Id()))
	for _, resource := range resources.Body {
		if tea.StringValue(resource.ResourceType) == "VSWITCH" {
			vswitchIds = append(vswitchIds, tea.StringValue(resource.InstanceId))
		}
	}

	d.Set("name", object.Name)
	d.Set("vpc_id", object.VpcId)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("vswitch_ids", vswitchIds)
	d.Set("security_group_id", object.SecurityGroupId)
	d.Set("deletion_protection", object.DeletionProtection)
	d.Set("version", object.CurrentVersion)
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("cluster_spec", object.ClusterSpec)

	if err := d.Set("tags", flattenTagsConfig(object.Tags)); err != nil {
		return WrapError(err)
	}
	if d.Get("load_balancer_spec") == "" {
		d.Set("load_balancer_spec", "slb.s1.small")
	}
	if d.Get("logging_type") == "" {
		d.Set("logging_type", "SLS")
	}

	var requestInfo *cs.Client
	var response interface{}

	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.GetClusterCerts(d.Id())
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetClusterCerts", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = d.Id()
		addDebug("GetClusterCerts", response, requestInfo, requestMap)
	}
	cert, _ := response.(cs.ClusterCerts)
	if ce, ok := d.GetOk("client_cert"); ok && ce.(string) != "" {
		if err := writeToFile(ce.(string), cert.Cert); err != nil {
			return WrapError(err)
		}
	}
	if key, ok := d.GetOk("client_key"); ok && key.(string) != "" {
		if err := writeToFile(key.(string), cert.Key); err != nil {
			return WrapError(err)
		}
	}
	if ca, ok := d.GetOk("cluster_ca_cert"); ok && ca.(string) != "" {
		if err := writeToFile(ca.(string), cert.CA); err != nil {
			return WrapError(err)
		}
	}

	var config *cs.ClusterConfig
	if file, ok := d.GetOk("kube_config"); ok && file.(string) != "" {
		var requestInfo *cs.Client

		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return csClient.DescribeClusterUserConfig(d.Id(), !d.Get("endpoint_public_access_enabled").(bool))
			})
			response = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetClusterConfig", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			addDebug("GetClusterConfig", response, requestInfo, requestMap)
		}
		config, _ = response.(*cs.ClusterConfig)

		if err := writeToFile(file.(string), config.Config); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceAlicloudCSServerlessKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	// modify cluster tag
	if d.HasChange("tags") {
		err := updateKubernetesClusterTag(d, meta)
		if err != nil {
			return WrapErrorf(err, ResponseCodeMsg, d.Id(), "ModifyClusterTags", AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("tags")
	}

	// upgrade cluster version
	err := UpgradeAlicloudKubernetesCluster(d, meta)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpgradeClusterVersion", DenverdinoAliyungo)
	}

	if err := modifyKubernetesCluster(d, meta); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyCluster", DenverdinoAliyungo)
	}

	d.Partial(false)
	return resourceAlicloudCSServerlessKubernetesRead(d, meta)
}

func resourceAlicloudCSServerlessKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	csService := CsService{meta.(*connectivity.AliyunClient)}
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	args := &roacs.DeleteClusterRequest{}
	if v := d.Get("retain_resources"); len(v.([]interface{})) > 0 {
		args.RetainResources = tea.StringSlice(expandStringList(v.([]interface{})))
	}

	_, err = client.DeleteCluster(tea.String(d.Id()), args)
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "DeleteCluster", AliyunTablestoreGoSdk)
	}

	stateConf := BuildStateConf([]string{"running", "deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, csService.CsServerlessKubernetesInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func modifyKubernetesCluster(d *schema.ResourceData, meta interface{}) error {
	var update bool
	action := "ModifyCluster"
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}

	var modifyClusterRequest cs.ModifyClusterArgs
	if d.HasChange("deletion_protection") {
		update = true
		modifyClusterRequest.DeletionProtection = d.Get("deletion_protection").(bool)
	}

	if update {
		conn, err := meta.(*connectivity.AliyunClient).NewTeaRoaCommonClient(connectivity.OpenAckService)
		if err != nil {
			return WrapError(err)
		}
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2015-12-15"), nil, StringPointer("PUT"), StringPointer("AK"), String(fmt.Sprintf("/api/v2/clusters/%s", d.Id())), nil, nil, modifyClusterRequest, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
					return resource.RetryableError(err)
				}
				addDebug(action, response, nil)
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, nil)
			return nil
		})

		stateConf := BuildStateConf([]string{"updating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return err
		}

		if err != nil {
			return err
		}
	}
	d.SetPartial("deletion_protection")

	return nil
}
