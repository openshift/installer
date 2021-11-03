package alicloud

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	cdn2 "github.com/denverdino/aliyungo/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCdnDomainNew() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCdnDomainCreateNew,
		Read:   resourceAlicloudCdnDomainReadNew,
		Update: resourceAlicloudCdnDomainUpdateNew,
		Delete: resourceAlicloudCdnDomainDeleteNew,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(5, 67),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cdn_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(cdn2.CdnTypes, false),
			},
			"sources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(cdn2.SourceTypes, false),
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      80,
							ValidateFunc: validation.IntInSlice([]int{80, 443}),
						},
						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      20,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  10,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if v, ok := d.GetOk("sources"); ok && len(v.(*schema.Set).List()) > 0 {
									sources := make([]map[string]interface{}, len(v.(*schema.Set).List()))
									byteSources, _ := json.Marshal(v.(*schema.Set).List())
									json.Unmarshal(byteSources, &sources)
									for _, source := range sources {
										if source["type"].(string) == "ipaddr" && formatInt(source["weight"]) != 10 {
											return true
										}
									}
								}
								return false
							},
							ValidateFunc: validation.IntBetween(0, 100),
						},
					},
				},
			},
			"certificate_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_certificate": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"server_certificate_status": {
							Type:         schema.TypeString,
							Default:      "on",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
						},
						"private_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"force_set": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"1", "0"}, false),
						},
						"cert_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cert_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"upload", "cas", "free"}, false),
						},
					},
				},
				MaxItems: 1,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(cdn2.Scopes, false),
			},
			"tags": tagsSchema(),
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCdnDomainCreateNew(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnService := &CdnService{client: client}

	request := cdn.CreateAddCdnDomainRequest()
	request.RegionId = client.RegionId
	if v := d.Get("resource_group_id").(string); v != "" {
		request.ResourceGroupId = v
	}
	request.DomainName = d.Get("domain_name").(string)
	request.CdnType = d.Get("cdn_type").(string)
	if v, ok := d.GetOk("scope"); ok {
		request.Scope = v.(string)
	}

	sources, err := cdnService.convertCdnSourcesToString(d.Get("sources").(*schema.Set).List())
	if err != nil {
		return WrapError(err)
	}
	request.Sources = sources

	raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
		return cdnClient.AddCdnDomain(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cdn_domain", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(request.DomainName)

	err = cdnService.WaitForCdnDomain(d.Id(), Online, DefaultLongTimeout)
	if err != nil {
		return WrapError(err)
	}

	return resourceAlicloudCdnDomainUpdateNew(d, meta)
}

func resourceAlicloudCdnDomainUpdateNew(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnService := &CdnService{client}

	d.Partial(true)

	if !d.IsNewResource() {
		request := cdn.CreateModifyCdnDomainRequest()
		request.RegionId = client.RegionId
		if v := d.Get("resource_group_id").(string); v != "" {
			request.ResourceGroupId = v
		}
		request.DomainName = d.Id()

		if d.HasChange("sources") {
			sources, err := cdnService.convertCdnSourcesToString(d.Get("sources").(*schema.Set).List())
			if err != nil {
				return WrapError(err)
			}
			request.Sources = sources

			raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
				return cdnClient.ModifyCdnDomain(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)

			err = cdnService.WaitForCdnDomain(d.Id(), Online, DefaultTimeoutMedium)
			if err != nil {
				return WrapError(err)
			}
			d.SetPartial("sources")
			d.SetPartial("resource_group_id")

		}
	}

	if d.HasChange("certificate_config") {
		if err := certificateConfigUpdateNew(client, d); err != nil {
			return WrapError(err)
		}
	}

	if err := setCdnTags(client, TagResourceCdn, d); err != nil {
		return WrapError(err)
	}

	d.Partial(false)
	return resourceAlicloudCdnDomainReadNew(d, meta)
}

func resourceAlicloudCdnDomainReadNew(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	cdnService := &CdnService{client: client}
	object, err := cdnService.DescribeCdnDomainNew(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if len(object.SourceModels.SourceModel) > 0 {
		sources := make([]map[string]interface{}, 0)
		for _, model := range object.SourceModels.SourceModel {
			priority, _ := strconv.Atoi(model.Priority)
			weight, _ := strconv.Atoi(model.Weight)
			sources = append(sources, map[string]interface{}{
				"content":  model.Content,
				"port":     model.Port,
				"priority": priority,
				"type":     model.Type,
				"weight":   weight,
			})
		}

		err := d.Set("sources", sources)
		if err != nil {
			return WrapError(err)
		}
	}

	d.Set("domain_name", object.DomainName)
	d.Set("cdn_type", object.CdnType)
	d.Set("scope", object.Scope)
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("cname", object.Cname)

	certInfo, err := cdnService.DescribeDomainCertificateInfo(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	oldConfig := d.Get("certificate_config").([]interface{})
	config := make([]map[string]interface{}, 1)
	serverCertificateStatus := object.ServerCertificateStatus
	if serverCertificateStatus == "" {
		serverCertificateStatus = "off"
	}
	config[0] = map[string]interface{}{
		"server_certificate":        certInfo.ServerCertificate,
		"server_certificate_status": serverCertificateStatus,
		"cert_name":                 certInfo.CertName,
		"cert_type":                 certInfo.CertType,
	}
	if oldConfig != nil && len(oldConfig) > 0 {
		val := oldConfig[0].(map[string]interface{})
		config[0]["private_key"] = val["private_key"]
		config[0]["force_set"] = val["force_set"]
		config[0]["region"] = val["region"]
	}

	d.Set("certificate_config", config)

	tags, err := cdnService.DescribeTags(d.Id(), TagResourceCdn)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", cdnTagsToMap(tags))

	return nil
}

func resourceAlicloudCdnDomainDeleteNew(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnService := CdnService{client}
	request := cdn.CreateDeleteCdnDomainRequest()
	request.RegionId = client.RegionId
	request.DomainName = d.Id()
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
			return cdnClient.DeleteCdnDomain(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceBusy"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(cdnService.WaitForCdnDomain(d.Id(), Deleted, DefaultTimeout))
}

func certificateConfigUpdateNew(client *connectivity.AliyunClient, d *schema.ResourceData) error {
	cdnService := &CdnService{client}
	request := cdn.CreateSetDomainServerCertificateRequest()
	request.RegionId = client.RegionId
	request.DomainName = d.Id()
	v, ok := d.GetOk("certificate_config")
	if !ok {
		request.ServerCertificateStatus = "off"
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
				return cdnClient.SetDomainServerCertificate(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceBusy"}) {
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
		d.SetPartial("certificate_config")
		return nil
	}
	config := v.([]interface{})
	val := config[0].(map[string]interface{})
	request.ServerCertificateStatus = val["server_certificate_status"].(string)

	serverCertificate, okServerCertificate := val["server_certificate"]
	if okServerCertificate {
		request.ServerCertificate = serverCertificate.(string)
	}
	if v, ok := val["private_key"]; ok {
		request.PrivateKey = v.(string)
	}
	if v, ok := val["force_set"]; ok && v.(string) != "" {
		request.ForceSet = v.(string)
	}
	if v, ok := val["cert_name"]; ok && v.(string) != "" {
		request.CertName = v.(string)
	}
	if v, ok := val["cert_type"]; ok && v.(string) != "" {
		request.CertType = v.(string)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
			return cdnClient.SetDomainServerCertificate(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceBusy"}) {
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
	d.SetPartial("certificate_config")
	if serverCertificate != "" && request.ServerCertificateStatus != "off" {
		err := cdnService.WaitForServerCertificateNew(d.Id(), request.ServerCertificate, DefaultTimeout)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}
