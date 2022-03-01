package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudFCAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFCAliasCreate,
		Read:   resourceAlicloudFCAliasRead,
		Update: resourceAlicloudFCAliasUpdate,
		Delete: resourceAlicloudFCAliasDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alias_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"routing_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_version_weights": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeFloat},
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudFCAliasCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	serviceName := d.Get("service_name").(string)
	aliasName := d.Get("alias_name").(string)
	serviceVersion := d.Get("service_version").(string)

	request := fc.NewCreateAliasInput(serviceName).
		WithAliasName(aliasName).
		WithVersionID(serviceVersion)

	if description, ok := d.GetOk("description"); ok {
		request = request.WithDescription(description.(string))
	}

	if routingConfig, ok := d.GetOk("routing_config"); ok {
		v := expandFCAliasRoutingConfig(routingConfig.([]interface{}))
		request = request.WithAdditionalVersionWeight(v)
	}

	var response *fc.CreateAliasOutput
	var requestInfo *fc.Client
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.CreateAlias(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"AccessDenied"}) {
				return resource.RetryableError(WrapError(err))
			}
			return resource.NonRetryableError(WrapError(err))
		}
		addDebug("CreateAlias", raw, requestInfo, request)
		response, _ = raw.(*fc.CreateAliasOutput)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fc_alias", "CreateAlias", FcGoSdk)
	}

	if response == nil {
		return WrapError(Error("Creating function compute alias got an empty response"))
	}

	d.SetId(fmt.Sprintf("%s:%s", serviceName, *response.AliasName))

	return resourceAlicloudFCAliasRead(d, meta)
}

func resourceAlicloudFCAliasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	id := d.Id()
	response, err := fcService.DescribeFcAlias(id)
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fc_alias fcService.DescribeFcAlias Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	serviceName := parts[0]
	d.Set("service_name", serviceName)
	d.Set("alias_name", *response.AliasName)
	d.Set("description", *response.Description)
	d.Set("service_version", *response.VersionID)

	if err := d.Set("routing_config", flattenFCAliasRoutingConfig(response.AdditionalVersionWeight)); err != nil {
		return fmt.Errorf("error setting FC alias routing_config: %s", err)
	}

	return nil
}

func resourceAlicloudFCAliasUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return err
	}

	serviceName := parts[0]
	aliasName := parts[1]
	update := false
	request := fc.NewUpdateAliasInput(serviceName, aliasName)

	if d.HasChange("service_version") {
		update = true
		request = request.WithVersionID(d.Get("service_version").(string))
	}

	if d.HasChange("description") {
		update = true
		request = request.WithDescription(d.Get("description").(string))
	}

	if d.HasChange("routing_config") {
		update = true
		if routingConfig, ok := d.GetOk("routing_config"); ok {
			v := expandFCAliasRoutingConfig(routingConfig.([]interface{}))
			request = request.WithAdditionalVersionWeight(v)
		}
	}

	if update {
		var requestInfo *fc.Client
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.UpdateAlias(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateAlias", FcGoSdk)
		}
		addDebug("UpdateAlias", raw, requestInfo, request)
	}

	return resourceAlicloudFCAliasRead(d, meta)
}

func resourceAlicloudFCAliasDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return err
	}

	serviceName := parts[0]
	aliasName := parts[1]
	request := fc.NewDeleteAliasInput(serviceName, aliasName)
	var requestInfo *fc.Client
	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.DeleteAlias(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"AliasNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteAlias", FcGoSdk)
	}
	addDebug("DeleteCustomDomain", raw, requestInfo, request)
	return nil
}

func expandFCAliasRoutingConfig(l []interface{}) map[string]float64 {
	if len(l) > 0 && l[0] != nil {
		m := l[0].(map[string]interface{})
		if v, ok := m["additional_version_weights"]; ok {
			additionalVersionWeigth := expandFloat64Map(v.(map[string]interface{}))
			return additionalVersionWeigth
		}
	}

	return nil
}

func expandFloat64Map(m map[string]interface{}) map[string]float64 {
	float64Map := make(map[string]float64, len(m))
	for k, v := range m {
		float64Map[k] = v.(float64)
	}
	return float64Map
}

func flattenFCAliasRoutingConfig(additionalVersionWeights map[string]float64) []interface{} {
	if additionalVersionWeights == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"additional_version_weights": additionalVersionWeights,
	}

	return []interface{}{m}
}
