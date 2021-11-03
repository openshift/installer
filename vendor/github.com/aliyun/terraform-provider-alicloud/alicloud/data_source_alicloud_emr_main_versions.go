package alicloud

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/emr"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudEmrMainVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEmrMainVersionsRead,

		Schema: map[string]*schema.Schema{
			"emr_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"main_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"emr_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAlicloudEmrMainVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := emr.CreateListEmrMainVersionRequest()
	if emrVersion, ok := d.GetOk("emr_version"); ok {
		request.EmrVersion = strings.TrimSpace(emrVersion.(string))
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)

	raw, err := client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
		return emrClient.ListEmrMainVersion(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_emr_main_versions", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	var (
		mainVersions []emr.EmrMainVersion
		clusterTypes = make(map[string][]string)
	)
	response, _ := raw.(*emr.ListEmrMainVersionResponse)
	if response != nil {
		// get clusterInfo of specific emr version
		var (
			versionRequest  = emr.CreateDescribeEmrMainVersionRequest()
			versionResponse *emr.DescribeEmrMainVersionResponse
			versionRaw      interface{}
		)
		clusterTypeFilter := func(filter []interface{}, source []emr.ClusterTypeInfo) (result []string) {
			if len(source) == 0 {
				return
			}
			if len(filter) == 0 {
				for _, c := range source {
					result = append(result, c.ClusterType)
				}
				return
			}
			sourceMapping := make(map[string]struct{})
			for _, s := range source {
				sourceMapping[s.ClusterType] = struct{}{}
			}
			for _, f := range filter {
				if _, ok := sourceMapping[f.(string)]; !ok {
					return nil
				}
				result = append(result, f.(string))
			}
			return
		}
		for _, v := range response.EmrMainVersionList.EmrMainVersion {
			if v.EmrVersion == "" {
				continue
			}
			versionRequest.EmrVersion = v.EmrVersion
			versionRaw, err = client.WithEmrClient(func(emrClient *emr.Client) (interface{}, error) {
				return emrClient.DescribeEmrMainVersion(versionRequest)
			})
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_emr_main_versions", request.GetActionName(), AlibabaCloudSdkGoERROR)
			}

			versionResponse, _ = versionRaw.(*emr.DescribeEmrMainVersionResponse)
			if versionResponse == nil {
				continue
			}
			var (
				clusterTypeInfo = versionResponse.EmrMainVersion.ClusterTypeInfoList.ClusterTypeInfo
				types           []string
			)

			// filter by specific clusterType
			if types = clusterTypeFilter(d.Get("cluster_type").([]interface{}), clusterTypeInfo); len(types) == 0 {
				continue
			}
			clusterTypes[v.EmrVersion] = types
		}

		mainVersions = response.EmrMainVersionList.EmrMainVersion
	}

	return emrClusterMainVersionAttributes(d, clusterTypes, mainVersions)
}

func emrClusterMainVersionAttributes(d *schema.ResourceData, clusterTypes map[string][]string, mainVersions []emr.EmrMainVersion) error {
	var (
		ids []string
		s   []map[string]interface{}
	)

	for _, version := range mainVersions {
		// if display is false, ignore it
		if !version.Display {
			continue
		}

		ct := clusterTypes[version.EmrVersion]
		if len(ct) == 0 {
			continue
		}

		mapping := map[string]interface{}{
			"image_id":      version.ImageId,
			"emr_version":   version.EmrVersion,
			"cluster_types": ct,
		}

		s = append(s, mapping)
		ids = append(ids, version.EmrVersion)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("main_versions", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
