package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/market"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudMarketOrder() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMarketOrderCreate,
		Read:   resourceAlicloudMarketOrderRead,
		Delete: resourceAlicloudMarketOrderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"product_code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(8, 12),
			},
			"pay_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PostPaid",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PrePaid", "PostPaid"}, false),
			},
			"duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 256),
			},
			"pricing_cycle": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Day", "Month", "Year"}, false),
			},
			"package_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quantity": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1,
			},
			"coupon_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"components": {
				Type:     schema.TypeMap,
				Optional: true,
				Default:  map[string]interface{}{},
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudMarketOrderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request, err := buildAliyunMarketOrderArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	request.RegionId = client.RegionId
	raw, err := client.WithMarketClient(func(marketClient *market.Client) (interface{}, error) {
		return marketClient.CreateOrder(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_market_order", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*market.CreateOrderResponse)

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(response.OrderId)
	return resourceAlicloudMarketOrderRead(d, meta)
}

func resourceAlicloudMarketOrderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	marketService := MarketService{client}

	object, err := marketService.DescribeMarketOrder(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	err = d.Set("product_code", object.ProductCode)
	if err != nil {
		return WrapError(err)
	}
	orderPayType := "PostPaid"
	if strings.Split(object.ProductSkuCode, "-")[1] == "prepay" {
		orderPayType = "PrePaid"
	}
	err = d.Set("pay_type", orderPayType)
	if err != nil {
		return WrapError(err)
	}
	err = d.Set("quantity", object.Quantity)
	if err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudMarketOrderDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func buildAliyunMarketOrderArgs(d *schema.ResourceData, meta interface{}) (*market.CreateOrderRequest, error) {
	packageVersion := d.Get("package_version").(string)
	couponId := ""
	if v, ok := d.GetOk("coupon_id"); ok {
		couponId = v.(string)
	}
	duration := d.Get("duration").(int)
	pricingCycle := d.Get("pricing_cycle").(string)
	productCode := d.Get("product_code").(string)
	quantity := d.Get("quantity").(int)
	orderPayType := "postpay"
	if d.Get("pay_type").(string) == "PrePaid" {
		orderPayType = "prepay"
	}
	skuCode := fmt.Sprintf("%s-%s", productCode, orderPayType)
	components := d.Get("components").(map[string]interface{})
	components["package_version"] = packageVersion
	componentsJson, _ := json.Marshal(components)
	commodity := fmt.Sprintf(`{
    "components": %s,
    "couponId": "%s", 
    "duration": %d, 
    "pricingCycle": "%s",  
    "productCode": "%s",
    "quantity": %d,	
    "skuCode": "%s"  
	}`, componentsJson, couponId, duration, pricingCycle, productCode, quantity, skuCode)
	request := market.CreateCreateOrderRequest()
	request.ClientToken = buildClientToken(request.GetActionName())
	request.OrderType = "INSTANCE_BUY"
	request.PaymentType = "AUTO"
	request.Commodity = commodity
	return request, nil
}
