package alicloud

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudPolarDBEndpointAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBEndpointAddressCreate,
		Read:   resourceAlicloudPolarDBEndpointAddressRead,
		Update: resourceAlicloudPolarDBEndpointAddressUpdate,
		Delete: resourceAlicloudPolarDBEndpointAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_endpoint_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"net_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Public"}, false),
				Default:      "Public",
				ForceNew:     true,
			},
			"connection_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9\\-]{4,28}[a-z0-9]$`), "The prefix must be 6 to 30 characters in length, and can contain lowercase letters, digits, and hyphens (-), must start with a letter and end with a digit or letter."),
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudPolarDBEndpointAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	clusterId := d.Get("db_cluster_id").(string)
	dbEndpointId := d.Get("db_endpoint_id").(string)
	netType := d.Get("net_type").(string)

	request := polardb.CreateCreateDBEndpointAddressRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = clusterId
	request.DBEndpointId = dbEndpointId
	request.NetType = netType
	if v, ok := d.GetOk("connection_prefix"); ok && v.(string) != "" {
		request.ConnectionStringPrefix = v.(string)
	}
	var raw interface{}
	var err error
	err = resource.Retry(8*time.Minute, func() *resource.RetryError {
		raw, err = client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.CreateDBEndpointAddress(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_endpoint_address", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, dbEndpointId))

	if err := polarDBService.WaitForPolarDBConnection(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	// wait instance running after allocating
	if err := polarDBService.WaitForPolarDBInstance(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudPolarDBEndpointAddressRead(d, meta)
}

func resourceAlicloudPolarDBEndpointAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	object, err := polarDBService.DescribePolarDBConnection(d.Id())

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("db_cluster_id", parts[0])
	d.Set("db_endpoint_id", parts[1])
	d.Set("port", object.Port)
	d.Set("net_type", object.NetType)
	d.Set("connection_string", object.ConnectionString)
	d.Set("ip_address", object.IPAddress)
	prefix := strings.Split(object.ConnectionString, ".")
	d.Set("connection_prefix", prefix[0])

	return nil
}

func resourceAlicloudPolarDBEndpointAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("connection_prefix") {
		request := polardb.CreateModifyDBEndpointAddressRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = parts[0]
		request.DBEndpointId = parts[1]
		object, err := polarDBService.DescribePolarDBConnection(d.Id())
		if err != nil {
			return WrapError(err)
		}
		request.NetType = object.NetType
		request.ConnectionStringPrefix = d.Get("connection_prefix").(string)
		request.NetType = d.Get("net_type").(string)
		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
				return polarDBClient.ModifyDBEndpointAddress(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"EndpointStatus.NotSupport", "OperationDenied.DBClusterStatus"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		// wait instance connection_prefix modify success
		if err := polarDBService.WaitForPolarDBConnectionPrefix(d.Id(), request.ConnectionStringPrefix, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
	}
	return resourceAlicloudPolarDBEndpointAddressRead(d, meta)
}

func resourceAlicloudPolarDBEndpointAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	parts, errParse := ParseResourceId(d.Id(), 2)
	if errParse != nil {
		return WrapError(errParse)
	}
	request := polardb.CreateDeleteDBEndpointAddressRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = parts[0]
	request.DBEndpointId = parts[1]

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		object, err := polarDBService.DescribePolarDBConnection(d.Id())
		if err != nil {
			return resource.NonRetryableError(WrapError(err))
		}
		request.NetType = object.NetType
		var raw interface{}
		raw, err = client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DeleteDBEndpointAddress(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBClusterStatus", "EndpointStatus.NotSupport"}) {
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

	return polarDBService.WaitForPolarDBConnection(d.Id(), Deleted, DefaultTimeoutMedium)
}
