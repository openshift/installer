package alicloud

import (
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBClusterCreate,
		Read:   resourceAlicloudPolarDBClusterRead,
		Update: resourceAlicloudPolarDBClusterUpdate,
		Delete: resourceAlicloudPolarDBClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"db_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_node_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"modify_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Upgrade", "Downgrade"}, false),
				Optional:     true,
				Default:      "Upgrade",
			},
			"db_node_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2, 16),
				Default:      2,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pay_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
				Optional:     true,
				Default:      PostPaid,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  RenewNotRenewal,
				ValidateFunc: validation.StringInSlice([]string{
					string(RenewAutoRenewal),
					string(RenewNormal),
					string(RenewNotRenewal)}, false),
				DiffSuppressFunc: polardbPostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12, 24, 36}),
				DiffSuppressFunc: polardbPostPaidAndRenewDiffSuppressFunc,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				DiffSuppressFunc: polardbPostPaidDiffSuppressFunc,
			},
			"db_cluster_ip_array": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_cluster_ip_array_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "default",
						},
						"security_ips": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"modify_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"Cover", "Append", "Delete"}, false),
						},
					},
				},
			},
			"security_ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"collector_status": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Enable", "Disabled"}, false),
				Optional:     true,
				Computed:     true,
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
			},
			"tde_status": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Enabled", "Disabled"}, false),
				Optional:     true,
				Default:      "Disabled",
			},
			"encrypt_new_tables": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"ON", "OFF"}, false),
				Optional:     true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudPolarDBClusterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	request, err := buildPolarDBCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}
	raw, err := client.WithPolarDBClient(func(polarClient *polardb.Client) (interface{}, error) {
		return polarClient.CreateDBCluster(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_cluster", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.CreateDBClusterResponse)
	d.SetId(response.DBClusterId)

	// wait cluster status change from Creating to running
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudPolarDBClusterUpdate(d, meta)
}

func resourceAlicloudPolarDBClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	d.Partial(true)

	if d.HasChange("parameters") {
		if err := polarDBService.ModifyParameters(d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("parameters")
	}

	if err := polarDBService.setClusterTags(d); err != nil {
		return WrapError(err)
	}

	conn, err := client.NewPolarDBClient()
	if err != nil {
		return WrapError(err)
	}

	payType := d.Get("pay_type").(string)
	if !d.IsNewResource() && d.HasChange("pay_type") {
		action := "TransformDBClusterPayType"
		request := map[string]interface{}{
			"RegionId":    client.RegionId,
			"DBClusterId": d.Id(),
			"PayType":     convertPolarDBPayTypeUpdateRequest(payType),
		}
		if payType == string(PrePaid) {
			period := d.Get("period").(int)
			request["UsedTime"] = strconv.Itoa(period)
			request["Period"] = Month
			if period > 9 {
				request["UsedTime"] = strconv.Itoa(period / 12)
				request["Period"] = Year
			}
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-08-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if payType == string(PrePaid) {
			d.SetPartial("period")
		}
		d.SetPartial("pay_type")
	}

	if (d.Get("pay_type").(string) == string(PrePaid)) &&
		(d.HasChange("renewal_status") || d.HasChange("auto_renew_period")) {
		status := d.Get("renewal_status").(string)
		request := polardb.CreateModifyAutoRenewAttributeRequest()
		request.DBClusterIds = d.Id()
		request.RenewalStatus = status

		if status == string(RenewAutoRenewal) {
			period := d.Get("auto_renew_period").(int)
			request.Duration = strconv.Itoa(period)
			request.PeriodUnit = string(Month)
			if period > 9 {
				request.Duration = strconv.Itoa(period / 12)
				request.PeriodUnit = string(Year)
			}
		}
		//wait asynchronously cluster payType
		if err := polarDBService.WaitForPolarDBPayType(d.Id(), "Prepaid", DefaultLongTimeout); err != nil {
			return WrapError(err)
		}
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyAutoRenewAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("renewal_status")
		d.SetPartial("auto_renew_period")
	}

	if d.HasChange("maintain_time") {
		request := polardb.CreateModifyDBClusterMaintainTimeRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = d.Id()
		request.MaintainTime = d.Get("maintain_time").(string)

		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBClusterMaintainTime(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("maintain_time")
	}

	if d.HasChange("db_cluster_ip_array") {

		if err := polarDBService.ModifyDBAccessWhitelistSecurityIps(d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("db_cluster_ip_array")
	}

	if d.HasChange("security_ips") {
		ipList := expandStringList(d.Get("security_ips").(*schema.Set).List())

		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}

		if err := polarDBService.ModifyDBSecurityIps(d.Id(), ipstr); err != nil {
			return WrapError(err)
		}
		d.SetPartial("security_ips")
	}

	if d.HasChange("db_node_count") {
		cluster, err := polarDBService.DescribePolarDBCluster(d.Id())
		if err != nil {
			return WrapError(err)
		}
		currentDbNodeCount := len(cluster.DBNodes.DBNode)
		expectDbNodeCount := d.Get("db_node_count").(int)
		if expectDbNodeCount > currentDbNodeCount {
			//create node
			expandDbNodes := &[]polardb.CreateDBNodesDBNode{
				{
					TargetClass: cluster.DBNodeClass,
				},
			}
			request := polardb.CreateCreateDBNodesRequest()
			request.RegionId = client.RegionId
			request.DBClusterId = d.Id()
			request.DBNode = expandDbNodes
			raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
				return polarDBClient.CreateDBNodes(request)
			})

			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			if err != nil {
				return WrapErrorf(
					err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			response, _ := raw.(*polardb.CreateDBNodesResponse)
			// wait cluster status change from DBNodeCreating to running
			stateConf := BuildStateConf([]string{"DBNodeCreating"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(response.DBClusterId, []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, response.DBClusterId)
			}
		} else {
			//delete node
			deleteDbNodeId := ""
			for _, dbNode := range cluster.DBNodes.DBNode {
				if dbNode.DBNodeRole == "Reader" {
					deleteDbNodeId = dbNode.DBNodeId
				}
			}
			request := polardb.CreateDeleteDBNodesRequest()
			request.RegionId = client.RegionId
			request.DBClusterId = d.Id()
			request.DBNodeId = &[]string{
				deleteDbNodeId,
			}

			raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
				return polarDBClient.DeleteDBNodes(request)
			})

			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			stateConf := BuildStateConf([]string{"DBNodeDeleting"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err = stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

	}

	if d.HasChange("collector_status") {
		request := polardb.CreateModifyDBClusterAuditLogCollectorRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = d.Id()
		request.CollectorStatus = d.Get("collector_status").(string)

		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBClusterAuditLogCollector(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("collector_status")
	}

	if v, ok := d.GetOk("db_type"); ok && v.(string) == "MySQL" {
		if d.HasChange("tde_status") {
			if v, ok := d.GetOk("tde_status"); ok && v.(string) != "Disabled" {
				action := "ModifyDBClusterTDE"
				request := map[string]interface{}{
					"DBClusterId": d.Id(),
					"TDEStatus":   convertPolarDBTdeStatusUpdateRequest(v.(string)),
				}
				if s, ok := d.GetOk("encrypt_new_tables"); ok && s.(string) != "" {
					request["EncryptNewTables"] = s.(string)
				}
				//retry
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-08-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				//wait tde status 'Enabled'

				stateConf := BuildStateConf([]string{}, []string{"Enabled"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, polarDBService.PolarDBClusterTDEStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
				d.SetPartial("tde_status")
				d.SetPartial("encrypt_new_tables")
			}
		}
	}

	if d.HasChange("security_group_ids") {
		securityGroupsList := expandStringList(d.Get("security_group_ids").(*schema.Set).List())
		securityGroupsStr := strings.Join(securityGroupsList[:], COMMA_SEPARATED)

		request := polardb.CreateModifyDBClusterAccessWhitelistRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = d.Id()
		request.WhiteListType = "SecurityGroup"
		request.SecurityGroupIds = securityGroupsStr
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBClusterAccessWhitelist(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("security_group_ids")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudPolarDBClusterRead(d, meta)
	}

	if d.HasChange("db_node_class") {
		request := polardb.CreateModifyDBNodeClassRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = d.Id()
		request.ModifyType = d.Get("modify_type").(string)
		request.DBNodeTargetClass = d.Get("db_node_class").(string)
		//wait asynchronously cluster nodes num the same
		if err := polarDBService.WaitForPolarDBNodeClass(d.Id(), DefaultLongTimeout); err != nil {
			return WrapError(err)
		}

		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBNodeClass(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		// wait cluster status change from Creating to running
		stateConf := BuildStateConf([]string{"ClassChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("db_node_class")
	}

	if d.HasChange("description") {
		request := polardb.CreateModifyDBClusterDescriptionRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = d.Id()
		request.DBClusterDescription = d.Get("description").(string)

		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBClusterDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("description")
	}

	d.Partial(false)
	return resourceAlicloudPolarDBClusterRead(d, meta)
}

func resourceAlicloudPolarDBClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	clusterAttribute, err := polarDBService.DescribePolarDBClusterAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	cluster, err := polarDBService.DescribePolarDBCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	dbClusterIPArrayName := "default"
	if v, ok := d.GetOk("db_cluster_ip_array_name"); ok {
		dbClusterIPArrayName = v.(string)
	}

	if err = polarDBService.DBClusterIPArrays(d, "db_cluster_ip_array", dbClusterIPArrayName); err != nil {
		return WrapError(err)
	}

	ips, err := polarDBService.DescribeDBSecurityIps(d.Id(), dbClusterIPArrayName)
	if err != nil {
		return WrapError(err)
	}

	d.Set("security_ips", ips)

	//describe endpoints
	if len(ips) == 1 && strings.HasPrefix(ips[0], LOCAL_HOST_IP) {
		d.Set("connection_string", "")
	} else {
		endpoints, err := polarDBService.DescribePolarDBInstanceNetInfo(d.Id())
		if err != nil {
			return WrapError(err)
		}
		for _, endpoint := range endpoints {
			if endpoint.EndpointType == "Cluster" {
				d.Set("connection_string", endpoint.AddressItems[0].ConnectionString)
			}
		}
	}

	d.Set("vswitch_id", clusterAttribute.VSwitchId)
	d.Set("pay_type", getChargeType(clusterAttribute.PayType))
	d.Set("id", clusterAttribute.DBClusterId)
	d.Set("description", clusterAttribute.DBClusterDescription)
	d.Set("db_type", clusterAttribute.DBType)
	d.Set("db_version", clusterAttribute.DBVersion)
	d.Set("maintain_time", clusterAttribute.MaintainTime)
	d.Set("zone_ids", clusterAttribute.ZoneIds)
	d.Set("db_node_class", cluster.DBNodeClass)
	d.Set("db_node_count", len(clusterAttribute.DBNodes))
	d.Set("resource_group_id", clusterAttribute.ResourceGroupId)
	tags, err := polarDBService.DescribeTags(d.Id(), "cluster")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", polarDBService.tagsToMap(tags))

	if clusterAttribute.PayType == string(Prepaid) {
		clusterAutoRenew, err := polarDBService.DescribePolarDBAutoRenewAttribute(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}
		renewPeriod := 1
		if clusterAutoRenew != nil {
			renewPeriod = clusterAutoRenew.Duration
		}
		if clusterAutoRenew != nil && clusterAutoRenew.PeriodUnit == string(Year) {
			renewPeriod = renewPeriod * 12
		}
		d.Set("auto_renew_period", renewPeriod)
		//period, err := computePeriodByUnit(clusterAttribute.CreationTime, clusterAttribute.ExpireTime, d.Get("period").(int), "Month")
		//if err != nil {
		//	return WrapError(err)
		//}
		//d.Set("period", period)
		d.Set("renewal_status", clusterAutoRenew.RenewalStatus)
	}

	if err = polarDBService.RefreshParameters(d); err != nil {
		return WrapError(err)
	}

	clusterCollectStatus, err := polarDBService.DescribeDBAuditLogCollectorStatus(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("collector_status", clusterCollectStatus)

	clusterTDEStatus, err := polarDBService.DescribeDBClusterTDE(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("tde_status", clusterTDEStatus["TDEStatus"])
	d.Set("encrypt_new_tables", clusterTDEStatus["EncryptNewTables"])

	securityGroups, err := polarDBService.DescribeDBSecurityGroups(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("security_group_ids", securityGroups)

	return nil
}

func resourceAlicloudPolarDBClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	cluster, err := polarDBService.DescribePolarDBClusterAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	// Pre paid cluster can not be release.
	if PayType(cluster.PayType) == Prepaid {
		return WrapError(Error("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically."))
	}

	request := polardb.CreateDeleteDBClusterRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Id()
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DeleteDBCluster(request)
		})

		if err != nil && !NotFoundError(err) {
			if IsExpectedErrors(err, []string{"OperationDenied.DBClusterStatus", "OperationDenied.PolarDBClusterStatus", "OperationDenied.ReadPolarDBClusterStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"Creating", "Running", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{}))
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildPolarDBCreateRequest(d *schema.ResourceData, meta interface{}) (*polardb.CreateDBClusterRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := polardb.CreateCreateDBClusterRequest()
	request.RegionId = string(client.Region)
	request.DBType = Trim(d.Get("db_type").(string))
	request.DBVersion = Trim(d.Get("db_version").(string))
	request.DBNodeClass = d.Get("db_node_class").(string)
	request.DBClusterDescription = d.Get("description").(string)
	request.ClientToken = buildClientToken(request.GetActionName())

	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request.ResourceGroupId = v.(string)
	}

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request.ZoneId = Trim(zone.(string))
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))

	if vswitchId != "" {
		request.VSwitchId = vswitchId
		request.ClusterNetworkType = strings.ToUpper(string(Vpc))

		// check vswitchId in zone
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return nil, WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if request.ZoneId != vsw.ZoneId {
			return nil, WrapError(Error("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request.ZoneId))
		}

		request.VPCId = vsw.VpcId
	}

	payType := Trim(d.Get("pay_type").(string))
	request.PayType = string(Postpaid)
	if payType == string(PrePaid) {
		request.PayType = string(Prepaid)
	}
	if PayType(request.PayType) == Prepaid {
		period := d.Get("period").(int)
		request.UsedTime = strconv.Itoa(period)
		request.Period = string(Month)
		if period > 9 {
			request.UsedTime = strconv.Itoa(period / 12)
			request.Period = string(Year)
		}
		if d.Get("renewal_status").(string) != string(RenewNotRenewal) {
			request.AutoRenew = requests.Boolean(strconv.FormatBool(true))
		} else {
			request.AutoRenew = requests.Boolean(strconv.FormatBool(false))
		}
	}

	request.TDEStatus = requests.NewBoolean(convertPolarDBTdeStatusCreateRequest(d.Get("tde_status").(string)))

	return request, nil
}

func convertPolarDBTdeStatusCreateRequest(source string) bool {
	switch source {
	case "Enabled":
		return true
	}
	return false
}

func convertPolarDBTdeStatusUpdateRequest(source string) string {
	switch source {
	case "Enabled":
		return "Enable"
	}
	return "Disable"
}

func convertPolarDBPayTypeUpdateRequest(source string) string {
	switch source {
	case "PrePaid":
		return "Prepaid"
	}
	return "Postpaid"
}
