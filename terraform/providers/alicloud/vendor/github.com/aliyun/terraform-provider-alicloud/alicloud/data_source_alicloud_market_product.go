package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/market"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudProduct() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudProductRead,

		Schema: map[string]*schema.Schema{
			"product_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"available_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values.
			"product": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"skus": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sku_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sku_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"package_versions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												// Currently, the API products can return package_version, but others can not for ensure.
												"package_version": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"package_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"images": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"image_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"image_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"region_id": {
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
				},
			},
		},
	}
}

func dataSourceAlicloudProductRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := market.CreateDescribeProductRequest()
	request.Code = d.Get("product_code").(string)
	raw, err := client.WithMarketClient(func(marketClient *market.Client) (interface{}, error) {
		return marketClient.DescribeProduct(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_market_product", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*market.DescribeProductResponse)
	return productDescriptionAttributes(d, response)
}

func productDescriptionAttributes(d *schema.ResourceData, product *market.DescribeProductResponse) error {
	var s []map[string]interface{}

	var skus []map[string]interface{}
	for _, sku := range product.ProductSkus.ProductSku {
		skuMapping := map[string]interface{}{
			"sku_code": sku.Code,
			"sku_name": sku.Name,
		}
		var imageIdMap []map[string]interface{}
		var pvMapSli []map[string]interface{}
		for _, module := range sku.Modules.Module {
			if module.Code == "package_version" {
				for _, property := range module.Properties.Property {
					if property.Key == "package_version" {
						for _, value := range property.PropertyValues.PropertyValue {
							pvMapping := map[string]interface{}{
								"package_version": value.Value,
								"package_name":    value.DisplayName,
							}
							pvMapSli = append(pvMapSli, pvMapping)
						}
					}
				}
			} else if module.Code == "img_id" {
				availableRegion := d.Get("available_region").(string)
				for _, property := range module.Properties.Property {
					if property.Key == "img_id" {
						for _, value := range property.PropertyValues.PropertyValue {
							regionId := getRegionByImageIdPrefix(value.Value)
							if availableRegion != "" && availableRegion != regionId {
								continue
							}
							imageIdMapping := map[string]interface{}{
								"image_id":   value.Value,
								"image_name": value.DisplayName,
								"region_id":  regionId,
							}
							imageIdMap = append(imageIdMap, imageIdMapping)
						}
					}
				}
			}
		}
		skuMapping["package_versions"] = pvMapSli
		skuMapping["images"] = imageIdMap
		skus = append(skus, skuMapping)
	}
	mapping := map[string]interface{}{
		"code":        product.Code,
		"name":        product.Name,
		"description": product.Description,
		"skus":        skus,
	}

	s = append(s, mapping)

	d.SetId(dataResourceIdHash([]string{product.Code}))
	if err := d.Set("product", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
