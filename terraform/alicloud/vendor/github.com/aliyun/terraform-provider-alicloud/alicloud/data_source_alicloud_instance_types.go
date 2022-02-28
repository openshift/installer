package alicloud

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

type instanceTypeWithOriginalPrice struct {
	InstanceType  ecs.InstanceType
	OriginalPrice float64
}

func dataSourceAlicloudInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type_family": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
			},
			"cpu_core_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"memory_size": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"gpu_amount": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"gpu_spec": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  PostPaid,
				// %q must contain a valid InstanceChargeType, expected common.PrePaid, common.PostPaid
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Vpc", "Classic"}, false),
			},
			"spot_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      NoSpot,
				ValidateFunc: validation.StringInSlice([]string{"NoSpot", "SpotAsPriceGo", "SpotWithPriceLimit"}, false),
			},
			"eni_amount": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"kubernetes_node_role": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(KubernetesNodeMaster),
					string(KubernetesNodeWorker),
				}, false),
			},
			"is_outdated": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sorted_by": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CPU",
					"Memory",
					"Price",
				}, false),
			},
			"system_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "cloud_efficiency",
				ValidateFunc: validation.StringInSlice([]string{"cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values.
			"instance_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_core_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"price": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"gpu": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"burstable_instance": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"initial_credit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"baseline_credit": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"eni_amount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"local_storage": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"capacity": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"amount": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"category": {
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

func dataSourceAlicloudInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	zoneId, validZones, _, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return err
	}

	mapInstanceTypes := make(map[string][]string)
	for _, zone := range validZones {
		if zoneId != "" && zoneId != zone.ZoneId {
			continue
		}
		for _, r := range zone.AvailableResources.AvailableResource {
			if r.Type == string(InstanceTypeResource) {
				for _, t := range r.SupportedResources.SupportedResource {
					if t.Status == string(SoldOut) {
						continue
					}

					zones, _ := mapInstanceTypes[t.Value]
					zones = append(zones, zone.ZoneId)
					mapInstanceTypes[t.Value] = zones
				}
			}
		}
	}

	cpu := d.Get("cpu_core_count").(int)
	mem := d.Get("memory_size").(float64)
	family := strings.TrimSpace(d.Get("instance_type_family").(string))
	gpuAmount := d.Get("gpu_amount").(int)
	gpuSpec := d.Get("gpu_spec").(string)

	req := ecs.CreateDescribeInstanceTypesRequest()
	req.InstanceTypeFamily = family

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstanceTypes(req)
	})
	if err != nil {
		return err
	}
	var instanceTypes []instanceTypeWithOriginalPrice
	resp, _ := raw.(*ecs.DescribeInstanceTypesResponse)
	if resp != nil {

		eniAmount := d.Get("eni_amount").(int)
		k8sNode := strings.TrimSpace(d.Get("kubernetes_node_role").(string))
		for _, types := range resp.InstanceTypes.InstanceType {
			if _, ok := mapInstanceTypes[types.InstanceTypeId]; !ok {
				continue
			}

			if cpu > 0 && types.CpuCoreCount != cpu {
				continue
			}

			if mem > 0 && types.MemorySize != mem {
				continue
			}
			if eniAmount > types.EniQuantity {
				continue
			}
			if gpuAmount > 0 && types.GPUAmount != gpuAmount {
				continue
			}
			if gpuSpec != "" && types.GPUSpec != gpuSpec {
				continue
			}
			// Kubernetes node does not support instance types which family is "ecs.t5" and spec less that c2g4
			// Kubernetes master node does not support gpu instance types which family prefixes with "ecs.gn"
			if k8sNode != "" {
				if types.InstanceTypeFamily == "ecs.t5" {
					continue
				}
				if types.CpuCoreCount < 2 || types.MemorySize < 4 {
					continue
				}
				if k8sNode == string(KubernetesNodeMaster) && strings.HasPrefix(types.InstanceTypeFamily, "ecs.gn") {
					continue
				}
			}

			instanceTypes = append(instanceTypes, instanceTypeWithOriginalPrice{
				InstanceType: types,
			})
		}
		sortedBy := d.Get("sorted_by").(string)

		if sortedBy == "Price" && len(instanceTypes) > 0 {
			bssopenapiService := BssopenapiService{client}

			priceList, err := getEcsInstanceTypePrice(bssopenapiService, d.Get("instance_charge_type").(string), instanceTypes)
			if err != nil {
				return WrapError(err)
			}
			for i := 0; i < len(instanceTypes); i++ {
				instanceTypes[i].OriginalPrice = priceList[i]
			}
		}
	}

	return instanceTypesDescriptionAttributes(d, instanceTypes, mapInstanceTypes)
}

func instanceTypesDescriptionAttributes(d *schema.ResourceData, types []instanceTypeWithOriginalPrice, mapTypes map[string][]string) error {
	sortedBy := d.Get("sorted_by").(string)
	if sortedBy != "" {
		sort.SliceStable(types, func(i, j int) bool {
			switch sortedBy {
			case "Price":
				return types[i].OriginalPrice < types[j].OriginalPrice
			case "CPU":
				return types[i].InstanceType.CpuCoreCount < types[j].InstanceType.CpuCoreCount
			case "Memory":
				return types[i].InstanceType.MemorySize < types[j].InstanceType.MemorySize
			}
			return false
		})
	}

	var ids []string
	var s []map[string]interface{}
	for _, t := range types {
		mapping := map[string]interface{}{
			"id":             t.InstanceType.InstanceTypeId,
			"cpu_core_count": t.InstanceType.CpuCoreCount,
			"memory_size":    t.InstanceType.MemorySize,
			"family":         t.InstanceType.InstanceTypeFamily,
			"eni_amount":     t.InstanceType.EniQuantity,
		}
		if sortedBy == "Price" {
			mapping["price"] = fmt.Sprintf("%.4f", t.OriginalPrice)
		}
		zoneIds := mapTypes[t.InstanceType.InstanceTypeId]
		sort.Strings(zoneIds)
		mapping["availability_zones"] = zoneIds
		gpu := map[string]interface{}{
			"amount":   strconv.Itoa(t.InstanceType.GPUAmount),
			"category": t.InstanceType.GPUSpec,
		}
		mapping["gpu"] = gpu
		brust := map[string]interface{}{
			"initial_credit":  strconv.Itoa(t.InstanceType.InitialCredit),
			"baseline_credit": strconv.Itoa(t.InstanceType.BaselineCredit),
		}
		mapping["burstable_instance"] = brust
		local := map[string]interface{}{
			"capacity": strconv.FormatInt(t.InstanceType.LocalStorageCapacity, 10),
			"amount":   strconv.Itoa(t.InstanceType.LocalStorageAmount),
			"category": t.InstanceType.LocalStorageCategory,
		}
		mapping["local_storage"] = local

		ids = append(ids, t.InstanceType.InstanceTypeId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instance_types", s); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func getEcsInstanceTypePrice(bssopenapiService BssopenapiService, instanceChargeType string, instanceTypes []instanceTypeWithOriginalPrice) ([]float64, error) {
	client := bssopenapiService.client
	var modules interface{}
	moduleCode := "InstanceType"
	var payAsYouGo []bssopenapi.GetPayAsYouGoPriceModuleList
	var subsciption []bssopenapi.GetSubscriptionPriceModuleList
	for _, types := range instanceTypes {
		config := fmt.Sprintf("InstanceType:%s,IoOptimized:IoOptimized,ImageOs:linux,Region:%s",
			types.InstanceType.InstanceTypeId, client.RegionId)
		if instanceChargeType == string(PostPaid) {
			payAsYouGo = append(payAsYouGo, bssopenapi.GetPayAsYouGoPriceModuleList{
				ModuleCode: moduleCode,
				Config:     config,
				PriceType:  "Hour",
			})
		} else {
			subsciption = append(subsciption, bssopenapi.GetSubscriptionPriceModuleList{
				ModuleCode: moduleCode,
				Config:     config,
			})

		}
	}

	if len(payAsYouGo) != 0 {
		modules = payAsYouGo
	} else {
		modules = subsciption
	}

	return bssopenapiService.GetInstanceTypePrice("ecs", "", modules)
}
