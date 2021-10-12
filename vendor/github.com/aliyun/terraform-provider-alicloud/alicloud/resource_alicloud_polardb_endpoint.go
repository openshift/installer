package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudPolarDBEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBEndpointCreate,
		Read:   resourceAlicloudPolarDBEndpointRead,
		Update: resourceAlicloudPolarDBEndpointUpdate,
		Delete: resourceAlicloudPolarDBEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Custom", "Primary", "Cluster"}, false),
				Default:      "Custom",
			},
			"nodes": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"read_write_mode": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"ReadWrite", "ReadOnly"}, false),
				Optional:     true,
				Computed:     true,
			},
			"auto_add_new_nodes": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Enable", "Disable"}, false),
				Optional:     true,
				Computed:     true,
			},
			"endpoint_config": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"ssl_enabled": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Enable", "Disable", "Update"}, false),
				Optional:     true,
			},
			"ssl_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Public", "Private", "Inner"}, false),
			},
			"ssl_expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_auto_rotate": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enable", "Disable"}, false),
			},
			"ssl_certificate_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudPolarDBEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	clusterId := d.Get("db_cluster_id").(string)
	endpointType := d.Get("endpoint_type").(string)
	request := polardb.CreateCreateDBClusterEndpointRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = clusterId
	request.EndpointType = endpointType
	if nodes, ok := d.GetOk("nodes"); ok {
		nodes := expandStringList(nodes.(*schema.Set).List())
		dbNodes := strings.Join(nodes, ",")
		request.Nodes = dbNodes
	}
	if readWriteMode, ok := d.GetOk("read_write_mode"); ok {
		request.ReadWriteMode = readWriteMode.(string)
	}
	if autoAddNewNodes, ok := d.GetOk("auto_add_new_nodes"); ok {
		request.AutoAddNewNodes = autoAddNewNodes.(string)
	}
	if endpointConfig, ok := d.GetOk("endpoint_config"); ok {
		endpointConfig, err := json.Marshal(endpointConfig)
		if err != nil {
			return WrapError(err)
		}
		request.EndpointConfig = string(endpointConfig)
	}

	enpoints, err := polarDBService.DescribePolarDBInstanceNetInfo(clusterId)
	if err != nil {
		return WrapError(err)
	}
	oldEndpoints := make([]interface{}, 0)
	for _, value := range enpoints {
		oldEndpoints = append(oldEndpoints, value.DBEndpointId)
	}
	oldEndpointIds := schema.NewSet(schema.HashString, oldEndpoints)

	var raw interface{}
	err = resource.Retry(8*time.Minute, func() *resource.RetryError {
		raw, err = client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.CreateDBClusterEndpoint(request)
		})
		if err != nil {
			OperationDeniedDBStatus = append(OperationDeniedDBStatus, "ClusterEndpoint.StatusNotValid")
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_endpoint", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	dbEndpointId, err := polarDBService.WaitForPolarDBEndpoints(d, Active, oldEndpointIds, DefaultTimeoutMedium)
	if err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, dbEndpointId))

	return resourceAlicloudPolarDBEndpointUpdate(d, meta)
}

func resourceAlicloudPolarDBEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	parts, errParse := ParseResourceId(d.Id(), 2)
	if errParse != nil {
		return WrapError(errParse)
	}
	dbClusterId := parts[0]
	dbEndpointId := parts[1]
	object, err := polarDBService.DescribePolarDBClusterEndpoint(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("db_cluster_id", dbClusterId)
	d.Set("endpoint_type", object.EndpointType)
	nodes := strings.Split(object.Nodes, ",")
	d.Set("nodes", nodes)

	var autoAddNewNodes string
	var readWriteMode string
	if object.EndpointType == "Primary" {
		autoAddNewNodes = "Disable"
		readWriteMode = "ReadWrite"
	} else {
		autoAddNewNodes = object.AutoAddNewNodes
		readWriteMode = object.ReadWriteMode
	}
	d.Set("auto_add_new_nodes", autoAddNewNodes)
	d.Set("read_write_mode", readWriteMode)

	if err = polarDBService.RefreshEndpointConfig(d); err != nil {
		return WrapError(err)
	}

	dbClusterSSL, err := polarDBService.DescribePolarDBClusterSSL(d)

	var sslConnectionString string
	var sslExpireTime string
	var sslEnabled string
	if len(dbClusterSSL.Items) < 1 {
		sslConnectionString = ""
		sslExpireTime = ""
		sslEnabled = ""
	} else if len(dbClusterSSL.Items) == 1 && dbClusterSSL.Items[0].DBEndpointId == "" {
		sslConnectionString = dbClusterSSL.Items[0].SSLConnectionString
		sslExpireTime = dbClusterSSL.Items[0].SSLExpireTime
		sslEnabled = convertPolarDBSSLEnableResponse(dbClusterSSL.Items[0].SSLEnabled)
	} else {
		for _, item := range dbClusterSSL.Items {
			if item.DBEndpointId == dbEndpointId {
				sslConnectionString = item.SSLConnectionString
				sslExpireTime = item.SSLExpireTime
				sslEnabled = convertPolarDBSSLEnableResponse(item.SSLEnabled)
			}
		}
	}
	sslAutoRotate := dbClusterSSL.SSLAutoRotate

	if err := d.Set("ssl_connection_string", sslConnectionString); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ssl_expire_time", sslExpireTime); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ssl_auto_rotate", sslAutoRotate); err != nil {
		return WrapError(err)
	}
	if sslEnabled == "Enable" {
		d.Set("ssl_certificate_url", "https://apsaradb-public.oss-ap-southeast-1.aliyuncs.com/ApsaraDB-CA-Chain.zip?file=ApsaraDB-CA-Chain.zip&regionId="+polarDBService.client.RegionId)
	} else {
		d.Set("ssl_certificate_url", "")
	}
	return nil
}

func resourceAlicloudPolarDBEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	dbClusterId := parts[0]
	dbEndpointId := parts[1]
	if d.HasChange("nodes") || d.HasChange("read_write_mode") || d.HasChange("auto_add_new_nodes") || d.HasChange("endpoint_config") {
		modifyEndpointRequest := polardb.CreateModifyDBClusterEndpointRequest()
		modifyEndpointRequest.RegionId = client.RegionId
		modifyEndpointRequest.DBClusterId = dbClusterId
		modifyEndpointRequest.DBEndpointId = dbEndpointId

		configItem := make(map[string]string)
		if d.HasChange("nodes") {
			nodes := expandStringList(d.Get("nodes").(*schema.Set).List())
			dbNodes := strings.Join(nodes, ",")
			modifyEndpointRequest.Nodes = dbNodes
			configItem["Nodes"] = dbNodes
		}
		if d.HasChange("read_write_mode") {
			modifyEndpointRequest.ReadWriteMode = d.Get("read_write_mode").(string)
			configItem["ReadWriteMode"] = d.Get("read_write_mode").(string)
		}
		if d.HasChange("auto_add_new_nodes") {
			modifyEndpointRequest.AutoAddNewNodes = d.Get("auto_add_new_nodes").(string)
			configItem["AutoAddNewNodes"] = d.Get("auto_add_new_nodes").(string)
		}
		if d.HasChange("endpoint_config") {
			endpointConfig, err := json.Marshal(d.Get("endpoint_config"))
			if err != nil {
				return WrapError(err)
			}
			modifyEndpointRequest.EndpointConfig = string(endpointConfig)
			configItem["EndpointConfig"] = string(endpointConfig)
		}

		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
				return polarDBClient.ModifyDBClusterEndpoint(modifyEndpointRequest)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"EndpointStatus.NotSupport", "OperationDenied.DBClusterStatus"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(modifyEndpointRequest.GetActionName(), raw, modifyEndpointRequest.RpcRequest, modifyEndpointRequest)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifyEndpointRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		// wait cluster endpoint config modified
		if err := polarDBService.WaitPolardbEndpointConfigEffect(
			d.Id(), configItem, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("ssl_enabled") || d.HasChange("net_type") || d.HasChange("ssl_auto_rotate") {
		if d.Get("ssl_enabled") == "" && d.Get("net_type") != "" {
			return WrapErrorf(Error("Need to specify ssl_enabled as Enable or Disable, if you want to modify the net_type."), DefaultErrorMsg, d.Id(), "ModifyDBClusterSSL", ProviderERROR)
		}
		modifySSLRequest := polardb.CreateModifyDBClusterSSLRequest()
		modifySSLRequest.SSLEnabled = d.Get("ssl_enabled").(string)
		modifySSLRequest.NetType = d.Get("net_type").(string)
		modifySSLRequest.DBClusterId = dbClusterId
		modifySSLRequest.DBEndpointId = dbEndpointId
		modifySSLRequest.SSLAutoRotate = d.Get("ssl_auto_rotate").(string)
		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
				return polarDBClient.ModifyDBClusterSSL(modifySSLRequest)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"EndpointStatus.NotSupport", "OperationDenied.DBClusterStatus"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(modifySSLRequest.GetActionName(), raw, modifySSLRequest.RpcRequest, modifySSLRequest)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), modifySSLRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		// wait cluster status change from SSL_MODIFYING to Running
		stateConf := BuildStateConf([]string{"SSL_MODIFYING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(dbClusterId, []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, dbClusterId)
		}
	}

	return resourceAlicloudPolarDBEndpointRead(d, meta)
}

func resourceAlicloudPolarDBEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	parts, errParse := ParseResourceId(d.Id(), 2)
	if errParse != nil {
		return WrapError(errParse)
	}
	object, err := polarDBService.DescribePolarDBClusterEndpoint(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if object.EndpointType != "Custom" {
		return WrapErrorf(Error(fmt.Sprintf("%s type endpoint can not be deleted.", object.EndpointType)), DefaultErrorMsg, d.Id(), "DeleteDBClusterEndpoint", ProviderERROR)
	}

	request := polardb.CreateDeleteDBClusterEndpointRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = parts[0]
	request.DBEndpointId = parts[1]

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DeleteDBClusterEndpoint(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBClusterStatus", "EndpointStatus.NotSupport", "ClusterEndpoint.StatusNotValid"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound", "InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	endpointIds := schema.NewSet(schema.HashString, make([]interface{}, 0))
	dbEndpoint, err := polarDBService.WaitForPolarDBEndpoints(d, Deleted, endpointIds, DefaultTimeoutMedium)
	if dbEndpoint != "" || err != nil {
		return WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR)
	}
	return nil
}

func convertPolarDBSSLEnableResponse(source string) string {
	switch source {
	case "Enabled":
		return "Enable"
	}
	return "Disable"
}
