package alicloud

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudPolarDBNodeClasses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPolarDBInstanceClassesRead,

		Schema: map[string]*schema.Schema{
			"pay_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
			},
			"db_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"db_node_class": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values.
			"classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"supported_engines": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"available_resources": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"db_node_class": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"engine": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceAlicloudPolarDBInstanceClassesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := polardb.CreateDescribeDBClusterAvailableResourcesRequest()
	payType := d.Get("pay_type").(string)
	if payType == string(PostPaid) {
		request.PayType = "Postpaid"
	} else if payType == string(PrePaid) {
		request.PayType = "Prepaid"
	}

	dbType, dbTypeGot := d.GetOk("db_type")
	dbVersion, dbVersionGot := d.GetOk("db_version")
	if dbTypeGot && dbVersionGot {
		request.DBType = dbType.(string)
		request.DBVersion = dbVersion.(string)
	}
	if dbNodeClass, ok := d.GetOk("db_node_class"); ok {
		request.DBNodeClass = dbNodeClass.(string)
	}
	if regionId, ok := d.GetOk("region_id"); ok {
		request.RegionId = regionId.(string)
	}
	if zoneId, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = zoneId.(string)
	}

	var response = &polardb.DescribeDBClusterAvailableResourcesResponse{}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
			return polardbClient.DescribeDBClusterAvailableResources(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{Throttling}) {
				time.Sleep(time.Duration(5) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*polardb.DescribeDBClusterAvailableResourcesResponse)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_node_classes", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	ids := []string{}
	var availableClasses []interface{}
	for _, AvailableZone := range response.AvailableZones {
		zondId := AvailableZone.ZoneId
		ids = append(ids, zondId)

		supportedEngines := make([]interface{}, 0)
		for _, supportedEngine := range AvailableZone.SupportedEngines {
			if len(supportedEngine.AvailableResources) == 0 {
				continue
			}
			if dbTypeGot && !strings.Contains(strings.ToLower(supportedEngine.Engine), strings.ToLower(dbType.(string))) {
				continue
			}
			if dbVersionGot && !strings.Contains(strings.ToLower(supportedEngine.Engine), strings.ToLower(dbVersion.(string))) {
				continue
			}
			var dbNodeClasses []map[string]string
			for _, availableResource := range supportedEngine.AvailableResources {
				dbNodeClass := map[string]string{"db_node_class": availableResource.DBNodeClass}
				dbNodeClasses = append(dbNodeClasses, dbNodeClass)
			}
			availableResources := map[string]interface{}{
				"engine":              supportedEngine.Engine,
				"available_resources": dbNodeClasses,
			}
			supportedEngines = append(supportedEngines, availableResources)
			ids = append(ids, supportedEngine.Engine)
		}

		var availableClass map[string]interface{}
		if len(supportedEngines) > 0 {

			availableClass = map[string]interface{}{
				"zone_id":           zondId,
				"supported_engines": supportedEngines,
			}
		}
		if len(availableClass) > 0 {
			availableClasses = append(availableClasses, availableClass)
		}

	}
	d.SetId(dataResourceIdHash(ids))

	err = d.Set("classes", availableClasses)
	if err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok {
		err = writeToFile(output.(string), availableClasses)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}
