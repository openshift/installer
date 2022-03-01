package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudFCCustomDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFCCustomDomainCreate,
		Read:   resourceAlicloudFCCustomDomainRead,
		Update: resourceAlicloudFCCustomDomainUpdate,
		Delete: resourceAlicloudFCCustomDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"route_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"function_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"qualifier": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"methods": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"cert_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"private_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"certificate": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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
		},
	}
}

func resourceAlicloudFCCustomDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	name := d.Get("domain_name").(string)

	request := &fc.CreateCustomDomainInput{
		DomainName: StringPointer(name),
		Protocol:   StringPointer(d.Get("protocol").(string)),
	}

	request.WithRouteConfig(parseRouteConfig(d))

	request.WithCertConfig(parseCertConfig(d))

	var response *fc.CreateCustomDomainOutput
	var requestInfo *fc.Client
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.CreateCustomDomain(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"AccessDenied"}) {
				return resource.RetryableError(WrapError(err))
			}
			return resource.NonRetryableError(WrapError(err))
		}
		addDebug("CreateCusomDomain", raw, requestInfo, request)
		response, _ = raw.(*fc.CreateCustomDomainOutput)
		return nil

	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fc_custom_domain", "CreateCustomDomain", FcGoSdk)
	}

	if response == nil {
		return WrapError(Error("Creating function compute custom domain got a empty response"))
	}

	d.SetId(*response.DomainName)

	return resourceAlicloudFCCustomDomainRead(d, meta)
}

func resourceAlicloudFCCustomDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	object, err := fcService.DescribeFcCustomDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fc_custom_domain fcService.DescribeFcCustomDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain_name", object.DomainName)
	d.Set("account_id", object.AccountID)
	d.Set("protocol", object.Protocol)
	d.Set("api_version", object.APIVersion)
	d.Set("created_time", object.CreatedTime)
	d.Set("last_modified_time", object.LastModifiedTime)

	var routeConfig []map[string]interface{}
	if object.RouteConfig != nil {
		for _, v := range object.RouteConfig.Routes {
			routeConfig = append(routeConfig, map[string]interface{}{
				"path":          *v.Path,
				"service_name":  *v.ServiceName,
				"function_name": *v.FunctionName,
				"qualifier":     *v.Qualifier,
				"methods":       v.Methods,
			})
		}
	}
	if err := d.Set("route_config", routeConfig); err != nil {
		return WrapError(err)
	}

	var certConfig []map[string]interface{}
	if object.CertConfig != nil {
		if object.CertConfig.CertName != nil && object.CertConfig.Certificate != nil {
			oldConfig := d.Get("cert_config").([]interface{})
			certConfig = append(certConfig, map[string]interface{}{
				"cert_name":   *object.CertConfig.CertName,
				"certificate": *object.CertConfig.Certificate,
				// The FC service will not return private key crendential for security reason.
				// Read it from the terraform file.
				"private_key": oldConfig[0].(map[string]interface{})["private_key"],
			})
		} else if object.CertConfig.CertName == nil && object.CertConfig.Certificate == nil {
			// Skip the null cert config.
		} else {
			b, _ := json.Marshal(object)
			return WrapError(Error(fmt.Sprintf("Illegal cert config: %s", string(b))))
		}
	}
	if err := d.Set("cert_config", certConfig); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudFCCustomDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	domainName := d.Id()
	request := fc.NewUpdateCustomDomainInput(domainName)

	update := false
	if d.HasChange("protocol") {
		update = true
		request.Protocol = StringPointer(d.Get("protocol").(string))
	}

	if d.HasChange("route_config") {
		update = true
		request.WithRouteConfig(parseRouteConfig(d))
	}

	if d.HasChange("cert_config") {
		update = true
		request.WithCertConfig(parseCertConfig(d))
	}

	if update {
		var requestInfo *fc.Client
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.UpdateCustomDomain(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateCustomDomain", FcGoSdk)
		}
		addDebug("UpdateCustomDomain", raw, requestInfo, request)
	}

	return resourceAlicloudFCCustomDomainRead(d, meta)
}

func resourceAlicloudFCCustomDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}
	domainName := d.Id()
	request := fc.NewDeleteCustomDomainInput(domainName)
	var requestInfo *fc.Client
	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.DeleteCustomDomain(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DomainNameNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, domainName, "DeleteCustomDomain", FcGoSdk)
	}
	addDebug("DeleteCustomDomain", raw, requestInfo, request)
	return WrapError(fcService.WaitForFcCustomDomain(domainName, Deleted, DefaultTimeout))
}

func parseRouteConfig(d *schema.ResourceData) (config *fc.RouteConfig) {
	if v, ok := d.GetOk("route_config"); ok {
		routeList := v.([]interface{})
		var pathConfigList []fc.PathConfig
		for _, route := range routeList {
			m := route.(map[string]interface{})
			pathConfig := fc.NewPathConfig().
				WithPath(m["path"].(string)).
				WithServiceName(m["service_name"].(string)).
				WithFunctionName(m["function_name"].(string)).
				WithQualifier(m["qualifier"].(string)).
				WithMethods(expandStringList(m["methods"].([]interface{})))
			pathConfigList = append(pathConfigList, *pathConfig)
		}
		config = fc.NewRouteConfig().WithRoutes(pathConfigList)
	}
	return config
}

func parseCertConfig(d *schema.ResourceData) (config *fc.CertConfig) {
	if v, ok := d.GetOk("cert_config"); ok {
		m := v.([]interface{})[0].(map[string]interface{})
		config = &fc.CertConfig{
			CertName:    StringPointer(m["cert_name"].(string)),
			PrivateKey:  StringPointer(m["private_key"].(string)),
			Certificate: StringPointer(m["certificate"].(string)),
		}
	}
	return config
}
