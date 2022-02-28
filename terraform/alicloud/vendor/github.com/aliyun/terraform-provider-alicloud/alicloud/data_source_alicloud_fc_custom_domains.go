package alicloud

import (
	"regexp"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudFcCustomDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFcCustomDomainsRead,

		Schema: map[string]*schema.Schema{
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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"function_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"qualifier": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"methods": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"cert_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"certificate": {
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

func dataSourceAlicloudFcCustomDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var ids []string
	var names []string
	var customDomainMappings []map[string]interface{}
	nextToken := ""
	for {
		request := fc.NewListCustomDomainsInput()
		if nextToken != "" {
			request.NextToken = &nextToken
		}
		var requestInfo *fc.Client
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.ListCustomDomains(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_fc_custom_domains", "ListCustomDomains", FcGoSdk)
		}
		addDebug("ListCustomDomains", raw, requestInfo, request)
		response, _ := raw.(*fc.ListCustomDomainsOutput)

		if response.CustomDomains == nil || len(response.CustomDomains) < 1 {
			break
		}

		for _, domain := range response.CustomDomains {
			mapping := map[string]interface{}{
				"id":                 *domain.DomainName,
				"domain_name":        *domain.DomainName,
				"account_id":         *domain.AccountID,
				"protocol":           *domain.Protocol,
				"api_version":        *domain.APIVersion,
				"created_time":       domain.CreatedTime,
				"last_modified_time": domain.LastModifiedTime,
			}

			var routeConfigMappings []map[string]interface{}
			for _, v := range domain.RouteConfig.Routes {
				routeConfigMappings = append(routeConfigMappings, map[string]interface{}{
					"path":          *v.Path,
					"service_name":  *v.ServiceName,
					"function_name": *v.FunctionName,
					"qualifier":     *v.Qualifier,
					"methods":       v.Methods,
				})
			}
			mapping["route_config"] = routeConfigMappings

			var certConfigMappings []map[string]interface{}
			if domain.CertConfig != nil && domain.CertConfig.CertName != nil {
				certConfigMappings = append(certConfigMappings, map[string]interface{}{
					"cert_name":   *domain.CertConfig.CertName,
					"certificate": *domain.CertConfig.Certificate,
				})
			}
			mapping["cert_config"] = certConfigMappings

			// Filter by name.
			nameRegex, ok := d.GetOk("name_regex")
			if ok && nameRegex.(string) != "" {
				var r *regexp.Regexp
				if nameRegex != "" {
					r, err = regexp.Compile(nameRegex.(string))
					if err != nil {
						return WrapError(err)
					}
				}
				if r != nil && !r.MatchString(mapping["domain_name"].(string)) {
					continue
				}
			}
			// Filter by id.
			if len(idsMap) > 0 {
				if _, ok := idsMap[*domain.DomainName]; !ok {
					continue
				}
			}
			customDomainMappings = append(customDomainMappings, mapping)
			ids = append(ids, *domain.DomainName)
			names = append(names, *domain.DomainName)
		}

		nextToken = ""
		if response.NextToken != nil {
			nextToken = *response.NextToken
		}
		if nextToken == "" {
			break
		}
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("domains", customDomainMappings); err != nil {
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
		writeToFile(output.(string), customDomainMappings)
	}
	return nil
}
