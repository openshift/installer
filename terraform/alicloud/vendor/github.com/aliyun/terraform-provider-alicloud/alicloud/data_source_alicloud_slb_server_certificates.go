package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudSlbServerCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbServerCertificatesRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"tags": tagsSchema(),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"certificates": {
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
						"fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subject_alternative_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"alicloud_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alicloud_certificate_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"is_alicloud_certificate": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func severCertificateTagsMappings(d *schema.ResourceData, id string, meta interface{}) map[string]string {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	tags, err := slbService.DescribeTags(id, nil, TagResourceCertificate)

	if err != nil {
		return nil
	}

	return slbTagsToMap(tags)
}

func dataSourceAlicloudSlbServerCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := slb.CreateDescribeServerCertificatesRequest()
	request.RegionId = client.RegionId
	tags := d.Get("tags").(map[string]interface{})
	if tags != nil && len(tags) > 0 {
		Tags := make([]slb.DescribeServerCertificatesTag, 0, len(tags))
		for k, v := range tags {
			certificatesTag := slb.DescribeServerCertificatesTag{
				Key:   k,
				Value: v.(string),
			}
			Tags = append(Tags, certificatesTag)
		}
		request.Tag = &Tags
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeServerCertificates(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_slb_server_certificates", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeServerCertificatesResponse)
	var filteredTemp []slb.ServerCertificate
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r, err = regexp.Compile(nameRegex.(string))
			if err != nil {
				return WrapError(err)
			}
		}
		for _, certificate := range response.ServerCertificates.ServerCertificate {
			if r != nil && !r.MatchString(certificate.ServerCertificateName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[certificate.ServerCertificateId]; !ok {
					continue
				}
			}

			filteredTemp = append(filteredTemp, certificate)
		}
	} else {
		filteredTemp = response.ServerCertificates.ServerCertificate
	}

	return slbServerCertificatesDescriptionAttributes(d, filteredTemp, meta)
}

func slbServerCertificatesDescriptionAttributes(d *schema.ResourceData, certificates []slb.ServerCertificate, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, certificate := range certificates {

		var subjectAlternativeNames []string
		if certificate.SubjectAlternativeNames.SubjectAlternativeName != nil {
			subjectAlternativeNames = certificate.SubjectAlternativeNames.SubjectAlternativeName
		}
		mapping := map[string]interface{}{
			"id":                        certificate.ServerCertificateId,
			"name":                      certificate.ServerCertificateName,
			"fingerprint":               certificate.Fingerprint,
			"common_name":               certificate.CommonName,
			"subject_alternative_names": subjectAlternativeNames,
			"expired_time":              certificate.ExpireTime,
			"expired_timestamp":         certificate.ExpireTimeStamp,
			"created_time":              certificate.CreateTime,
			"created_timestamp":         certificate.CreateTimeStamp,
			"alicloud_certificate_id":   certificate.AliCloudCertificateId,
			"alicloud_certificate_name": certificate.AliCloudCertificateName,
			"is_alicloud_certificate":   certificate.IsAliCloudCertificate == 1,
			"resource_group_id":         certificate.ResourceGroupId,
			"tags":                      severCertificateTagsMappings(d, certificate.ServerCertificateId, meta),
		}
		ids = append(ids, certificate.ServerCertificateId)
		names = append(names, certificate.ServerCertificateName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("certificates", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
