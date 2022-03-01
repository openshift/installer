package alicloud

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	roacs "github.com/alibabacloud-go/cs-20151215/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"strconv"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	aliyungoecs "github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/yaml.v2"
)

const (
	KubernetesClusterNetworkTypeFlannel = "flannel"
	KubernetesClusterNetworkTypeTerway  = "terway"

	KubernetesClusterLoggingTypeSLS = "SLS"
)

var (
	KubernetesClusterNodeCIDRMasksByDefault = 24
)

func resourceAlicloudCSKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSKubernetesCreate,
		Read:   resourceAlicloudCSKubernetesRead,
		Update: resourceAlicloudCSKubernetesUpdate,
		Delete: resourceAlicloudCSKubernetesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(1, 63),
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "Terraform-Creation",
				ValidateFunc:  validation.StringLenBetween(0, 37),
				ConflictsWith: []string{"name"},
				Deprecated:    "Field 'name_prefix' has been deprecated from provider version 1.75.0.",
			},
			// master configurations
			"master_vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems: 3,
				MaxItems: 5,
			},
			"master_instance_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 3,
				MaxItems: 5,
			},
			"master_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ValidateFunc: validation.IntBetween(40, 500),
			},
			"master_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
			},
			"master_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
				DiffSuppressFunc: masterDiskPerformanceLevelDiffSuppressFunc,
			},
			"master_disk_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"master_instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:      PostPaid,
			},
			"master_period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Month,
				ValidateFunc:     validation.StringInSlice([]string{"Week", "Month"}, false),
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				// must be a valid period, expected [1-9], 12, 24, 36, 48 or 60,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36, 48, 60})),
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_auto_renew": {
				Type:             schema.TypeBool,
				Default:          false,
				Optional:         true,
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			// worker configurations
			"worker_vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems: 1,
			},
			"worker_instance_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 10,
			},
			"worker_number": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"worker_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ValidateFunc: validation.IntBetween(20, 32768),
			},
			"worker_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
			},
			"worker_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
				DiffSuppressFunc: workerDiskPerformanceLevelDiffSuppressFunc,
			},
			"worker_disk_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"worker_data_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          40,
				ValidateFunc:     validation.IntBetween(20, 32768),
				DiffSuppressFunc: workerDataDiskSizeSuppressFunc,
			},
			"worker_data_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"worker_data_disks": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"encrypted": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"performance_level": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"worker_instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:      PostPaid,
			},
			"worker_period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Month,
				ValidateFunc:     validation.StringInSlice([]string{"Week", "Month"}, false),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"worker_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36, 48, 60})),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"worker_auto_renew": {
				Type:             schema.TypeBool,
				Default:          false,
				Optional:         true,
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"worker_auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"exclude_autoscaler_nodes": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			// global configurations
			// Terway network
			"pod_vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MaxItems: 10,
			},
			// Flannel network
			"pod_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_cidr_mask": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      KubernetesClusterNodeCIDRMasksByDefault,
				ValidateFunc: validation.IntBetween(24, 28),
			},
			"new_nat_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"key_name", "kms_encrypted_password"},
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "kms_encrypted_password"},
			},
			"kms_encrypted_password": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "key_name"},
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"user_ca": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_ssh": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"load_balancer_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"slb.s1.small", "slb.s2.small", "slb.s2.medium", "slb.s3.small", "slb.s3.medium", "slb.s3.large"}, false),
				Default:      "slb.s1.small",
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"install_cloud_monitor": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// cpu policy options of kubelet
			"cpu_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "static"}, false),
			},
			"proxy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"iptables", "ipvs"}, false),
			},
			"addons": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"slb_internet_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"os_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Linux",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Windows", "Linux"}, false),
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "CentOS",
				ForceNew: true,
			},
			"node_port_range": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "30000-32767",
				ForceNew: true,
			},
			"runtime": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "docker",
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "19.03.5",
						},
					},
				},
			},
			"cluster_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "cluster.local",
				ForceNew:    true,
				Description: "cluster local domain",
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"effect": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"NoSchedule", "NoExecute", "PreferNoSchedule"}, false),
						},
					},
				},
			},
			"rds_instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"custom_san": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// computed parameters
			"kube_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_ca_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_authority": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"connections": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_server_internet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_server_intranet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"slb_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'slb_id' has been deprecated from provider version 1.9.2. New field 'slb_internet' replaces it.",
			},
			"slb_internet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slb_intranet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_enterprise_security_group": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"security_group_id"},
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_nodes": {
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
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"worker_nodes": {
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
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			// remove parameters below
			// mix vswitch_ids between master and worker is not a good guidance to create cluster
			"worker_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'worker_instance_type' has been removed from provider version 1.75.0. New field 'worker_instance_types' replaces it.",
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems: 3,
				MaxItems: 5,
				Removed:  "Field 'vswitch_ids' has been removed from provider version 1.75.0. New field 'master_vswitch_ids' and 'worker_vswitch_ids' replace it.",
			},
			// single instance type would cause extra troubles
			"master_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'master_instance_type' has been removed from provider version 1.75.0. New field 'master_instance_types' replaces it.",
			},
			// force update is a high risk operation
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Removed:  "Field 'force_update' has been removed from provider version 1.75.0.",
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// single az would be never supported.
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'vswitch_id' has been removed from provider version 1.75.0. New field 'master_vswitch_ids' and 'worker_vswitch_ids' replaces it.",
			},
			// worker_numbers in array is a hell of management
			"worker_numbers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:    schema.TypeInt,
					Default: 3,
				},
				MinItems: 1,
				MaxItems: 3,
				Removed:  "Field 'worker_numbers' has been removed from provider version 1.75.0. New field 'worker_number' replaces it.",
			},
			"nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Removed:  "Field 'nodes' has been removed from provider version 1.9.4. New field 'master_nodes' replaces it.",
			},
			// too hard to use this config
			"log_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{KubernetesClusterLoggingTypeSLS}, false),
							Required:     true,
						},
						"project": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Removed: "Field 'log_config' has been removed from provider version 1.75.0. New field 'addons' replaces it.",
			},
			"cluster_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{KubernetesClusterNetworkTypeFlannel, KubernetesClusterNetworkTypeTerway}, false),
				Removed:      "Field 'cluster_network_type' has been removed from provider version 1.75.0. New field 'addons' replaces it.",
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_name_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^customized,[a-z0-9]([-a-z0-9\.])*,([5-9]|[1][0-2]),([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`), "Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test."),
			},
			"worker_ram_role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_account_issuer": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"api_audiences": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"retain_resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAlicloudCSKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()

	var requestInfo *cs.Client
	var raw interface{}

	// prepare args and set default value
	args, err := buildKubernetesArgs(d, meta)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes", "PrepareKubernetesClusterArgs", err)
	}

	if err = invoker.Run(func() error {
		raw, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			args.RegionId = common.Region(client.RegionId)
			args.ClusterType = cs.DelicatedKubernetes
			return csClient.CreateDelicatedKubernetesCluster(args)
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes", "CreateKubernetesCluster", raw)
	}

	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Params"] = args
		addDebug("CreateKubernetesCluster", raw, requestInfo, requestMap)
	}

	cluster, ok := raw.(*cs.ClusterCommonResponse)
	if ok != true {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_kubernetes", "ParseKubernetesClusterResponse", raw)
	}

	d.SetId(cluster.ClusterID)

	// reset interval to 10s
	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCSKubernetesRead(d, meta)
}

func resourceAlicloudCSKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	d.Partial(true)
	invoker := NewInvoker()
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		var requestInfo cs.ModifyClusterArgs
		requestInfo.ResourceGroupId = d.Get("resource_group_id").(string)

		response, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.ModifyCluster(d.Id(), &requestInfo)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyCluster", DenverdinoAliyungo)
		}
		addDebug("ModifyCluster", response, requestInfo)
		d.SetPartial("resource_group_id")
	}
	if d.HasChange("worker_number") && !d.IsNewResource() {
		oldV, newV := d.GetChange("worker_number")

		oldValue, ok := oldV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("worker_number old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := newV.(int)
		if ok != true {
			return WrapErrorf(fmt.Errorf("worker_number new value can not be parsed"), "parseError %d", newValue)
		}

		// scale out cluster.
		if newValue > oldValue {
			password := d.Get("password").(string)
			if password == "" {
				if v := d.Get("kms_encrypted_password").(string); v != "" {
					kmsService := KmsService{client}
					decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
					if err != nil {
						return WrapError(err)
					}
					password = decryptResp
				}
			}

			keyPair := d.Get("key_name").(string)

			args := &cs.ScaleOutKubernetesClusterRequest{
				KeyPair:             keyPair,
				LoginPassword:       password,
				ImageId:             d.Get("image_id").(string),
				UserData:            d.Get("user_data").(string),
				Count:               int64(newValue) - int64(oldValue),
				WorkerVSwitchIds:    expandStringList(d.Get("worker_vswitch_ids").([]interface{})),
				WorkerInstanceTypes: expandStringList(d.Get("worker_instance_types").([]interface{})),
			}

			if v := d.Get("user_data").(string); v != "" {
				_, base64DecodeError := base64.StdEncoding.DecodeString(v)
				if base64DecodeError == nil {
					args.UserData = v
				} else {
					args.UserData = base64.StdEncoding.EncodeToString([]byte(v))
				}
			}

			if v, ok := d.GetOk("worker_instance_charge_type"); ok {
				args.WorkerInstanceChargeType = v.(string)
				if args.WorkerInstanceChargeType == string(PrePaid) {
					args.WorkerAutoRenew = d.Get("worker_auto_renew").(bool)
					args.WorkerAutoRenewPeriod = d.Get("worker_auto_renew_period").(int)
					args.WorkerPeriod = d.Get("worker_period").(int)
					args.WorkerPeriodUnit = d.Get("worker_period_unit").(string)
				}
			}

			if v, ok := d.GetOk("worker_disk_category"); ok {
				args.WorkerSystemDiskCategory = aliyungoecs.DiskCategory(v.(string))
			}

			if v, ok := d.GetOk("worker_disk_size"); ok {
				args.WorkerSystemDiskSize = int64(v.(int))
			}

			if v, ok := d.GetOk("worker_disk_snapshot_policy_id"); ok {
				args.WorkerSnapshotPolicyId = v.(string)
			}

			if v, ok := d.GetOk("worker_disk_performance_level"); ok {
				args.WorkerSystemDiskPerformanceLevel = v.(string)
			}

			if dds, ok := d.GetOk("worker_data_disks"); ok {
				disks := dds.([]interface{})
				createDataDisks := make([]cs.DataDisk, 0, len(disks))
				for _, e := range disks {
					pack := e.(map[string]interface{})
					dataDisk := cs.DataDisk{
						Size:                 pack["size"].(string),
						DiskName:             pack["name"].(string),
						Category:             pack["category"].(string),
						Device:               pack["device"].(string),
						AutoSnapshotPolicyId: pack["auto_snapshot_policy_id"].(string),
						KMSKeyId:             pack["kms_key_id"].(string),
						Encrypted:            pack["encrypted"].(string),
						PerformanceLevel:     pack["performance_level"].(string),
					}
					createDataDisks = append(createDataDisks, dataDisk)
				}
				args.WorkerDataDisks = createDataDisks
			}

			if _, ok := d.GetOk("tags"); ok {
				if tags, err := ConvertCsTags(d); err == nil {
					args.Tags = tags
				}
				d.SetPartial("tags")
			}

			if d.HasChange("taints") && !d.IsNewResource() {
				args.Taints = expandKubernetesTaintsConfig(d.Get("taints").([]interface{}))
			}

			if d.HasChange("runtime") && !d.IsNewResource() {
				args.Runtime = expandKubernetesRuntimeConfig(d.Get("runtime").(map[string]interface{}))
			}

			if d.HasChange("rds_instances") && !d.IsNewResource() {
				args.RdsInstances = expandStringList(d.Get("rds_instances").([]interface{}))
			}

			if d.HasChange("cpu_policy") && !d.IsNewResource() {
				args.CpuPolicy = d.Get("cpu_policy").(string)
			}

			if d.HasChange("install_cloud_monitor") && !d.IsNewResource() {
				args.CloudMonitorFlags = d.Get("install_cloud_monitor").(bool)
			}

			if d.HasChange("image_id") && !d.IsNewResource() {
				args.ImageId = d.Get("image_id").(string)
			}

			var resoponse interface{}
			if err := invoker.Run(func() error {
				var err error
				resoponse, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
					resp, err := csClient.ScaleOutKubernetesCluster(d.Id(), args)
					return resp, err
				})
				return err
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ResizeKubernetesCluster", DenverdinoAliyungo)
			}
			if debugOn() {
				resizeRequestMap := make(map[string]interface{})
				resizeRequestMap["ClusterId"] = d.Id()
				resizeRequestMap["Args"] = args
				addDebug("ResizeKubernetesCluster", resoponse, resizeRequestMap)
			}

			stateConf := BuildStateConf([]string{"scaling"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("worker_number")
			d.SetPartial("worker_data_disks")
			d.SetPartial("worker_disk_category")
			d.SetPartial("worker_disk_size")
			d.SetPartial("worker_disk_snapshot_policy_id")
			d.SetPartial("worker_disk_performance_level")
		}

		// remove cluster nodes.
		if newValue < oldValue {
			nodes, err := removeKubernetesNodes(d, meta)
			if err != nil {
				return WrapErrorf(fmt.Errorf("node removed failed"), "node:%++v, err:%++v", nodes, err)
			}
		}
	}

	// modify cluster name
	if !d.IsNewResource() && (d.HasChange("name") || d.HasChange("name_prefix")) {
		var clusterName string
		if v, ok := d.GetOk("name"); ok {
			clusterName = v.(string)
		} else {
			clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
		}
		var requestInfo *cs.Client
		var response interface{}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				return nil, csClient.ModifyClusterName(d.Id(), clusterName)
			})
			response = raw
			return err
		}); err != nil && !IsExpectedErrors(err, []string{"ErrorClusterNameAlreadyExist"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyClusterName", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			requestMap["ClusterName"] = clusterName
			addDebug("ModifyClusterName", response, requestInfo, requestMap)
		}
		d.SetPartial("name")
		d.SetPartial("name_prefix")
	}

	// modify cluster deletion protection
	if !d.IsNewResource() && d.HasChange("deletion_protection") {
		var requestInfo cs.ModifyClusterArgs
		if v, ok := d.GetOk("deletion_protection"); ok {
			requestInfo.DeletionProtection = v.(bool)
		}

		var response interface{}
		if err := invoker.Run(func() error {
			_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				return nil, csClient.ModifyCluster(d.Id(), &requestInfo)
			})
			return err
		}); err != nil && !IsExpectedErrors(err, []string{"ErrorModifyDeletionProtectionFailed"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyCluster", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			requestMap["deletion_protection"] = requestInfo.DeletionProtection
			addDebug("ModifyCluster", response, requestInfo, requestMap)
		}
		d.SetPartial("deletion_protection")
	}

	// modify cluster maintenance window
	if !d.IsNewResource() && d.HasChange("maintenance_window") {
		var mw cs.MaintenanceWindow
		if v := d.Get("maintenance_window").([]interface{}); len(v) > 0 {
			mw = expandMaintenanceWindowConfig(v)
		}
		_ = modifyMaintenanceWindow(d, meta, mw)
	}

	// modify cluster tag
	if d.HasChange("tags") {
		err := updateKubernetesClusterTag(d, meta)
		if err != nil {
			return WrapErrorf(err, ResponseCodeMsg, d.Id(), "ModifyClusterTags", AlibabaCloudSdkGoERROR)
		}
	}

	// migrate cluster to pro form standard
	if d.HasChange("cluster_spec") {
		oldValue, newValue := d.GetChange("cluster_spec")
		o, ok := oldValue.(string)
		if ok != true {
			return WrapErrorf(fmt.Errorf("cluster_spec old value can not be parsed"), "parseError %d", oldValue)
		}
		n, ok := newValue.(string)
		if ok != true {
			return WrapErrorf(fmt.Errorf("cluster_pec new value can not be parsed"), "parseError %d", newValue)
		}

		if o == "ack.standard" && strings.Contains(n, "pro") {
			err := migrateAlicloudManagedKubernetesCluster(d, meta)
			if err != nil {
				return WrapErrorf(err, ResponseCodeMsg, d.Id(), "MigrateCluster", AlibabaCloudSdkGoERROR)
			}
		}
	}

	d.SetPartial("tags")

	UpgradeAlicloudKubernetesCluster(d, meta)
	d.Partial(false)
	return resourceAlicloudCSKubernetesRead(d, meta)
}

func resourceAlicloudCSKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()
	object, err := csService.DescribeCsKubernetes(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("vpc_id", object.VpcId)
	d.Set("security_group_id", object.SecurityGroupId)
	d.Set("version", object.CurrentVersion)
	d.Set("worker_ram_role_name", object.WorkerRamRoleName)
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("cluster_spec", object.ClusterSpec)
	d.Set("deletion_protection", object.DeletionProtection)

	if err := d.Set("tags", flattenTagsConfig(object.Tags)); err != nil {
		return WrapError(err)
	}
	if d.Get("os_type") == "" {
		d.Set("os_type", "Linux")
	}
	if d.Get("platform") == "" {
		d.Set("platform", "CentOS")
	}
	if d.Get("node_port_range") == "" {
		d.Set("node_port_range", "30000-32767")
	}
	if d.Get("cluster_domain") == "" {
		d.Set("cluster_domain", "cluster.local")
	}
	if d.Get("load_balancer_spec") == "" {
		d.Set("load_balancer_spec", "slb.s1.small")
	}
	// d.Set("os_type", object.OSType)
	// d.Set("platform", object.Platform)
	// d.Set("timezone", object.TimeZone)
	// d.Set("cluster_domain", object.ClusterDomin)
	// d.Set("custom_san",object.CustomSAN)
	// d.Set("runtime", object.Runtime)
	// d.Set("taints", object.Taits)
	// d.Set("rds_instances", object.RdsInstances)
	// d.Set("node_port_range", object.NodePortRange)
	d.Set("maintenance_window", flattenMaintenanceWindowConfig(&object.MaintenanceWindow))

	var masterNodes []map[string]interface{}
	var workerNodes []map[string]interface{}
	var defaultNodePoolId string
	var nodePoolDetails interface{}
	// get the default nodepool id
	if err := invoker.Run(func() error {
		var err error
		nodePoolDetails, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			nodePools, err := csClient.DescribeClusterNodePools(d.Id())
			return *nodePools, err
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DescribeClusterNodePools", DenverdinoAliyungo)
	}

	for _, v := range nodePoolDetails.([]cs.NodePoolDetail) {
		if v.BasicNodePool.NodePoolInfo.Name == "default-nodepool" {
			defaultNodePoolId = v.BasicNodePool.NodePoolInfo.NodePoolId
		}
	}

	pageNumber := 1
	for {
		if defaultNodePoolId == "" {
			break
		}
		var result []cs.KubernetesNodeType
		var pagination *cs.PaginationResult
		var requestInfo *cs.Client
		var response interface{}
		if err := invoker.Run(func() error {
			raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
				requestInfo = csClient
				nodes, paginationResult, err := csClient.GetKubernetesClusterNodes(d.Id(), common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge}, defaultNodePoolId)
				return []interface{}{nodes, paginationResult}, err
			})
			response = raw
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", DenverdinoAliyungo)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			requestMap["Pagination"] = common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge}
			addDebug("GetKubernetesClusterNodes", response, requestInfo, requestMap)
		}
		result, _ = response.([]interface{})[0].([]cs.KubernetesNodeType)
		pagination, _ = response.([]interface{})[1].(*cs.PaginationResult)

		if pageNumber == 1 && (len(result) == 0 || result[0].InstanceId == "") {
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				if err := invoker.Run(func() error {
					raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
						requestInfo = csClient
						nodes, _, err := csClient.GetKubernetesClusterNodes(d.Id(), common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge}, defaultNodePoolId)
						return nodes, err
					})
					response = raw
					return err
				}); err != nil {
					return resource.NonRetryableError(err)
				}
				tmp, _ := response.([]cs.KubernetesNodeType)
				if len(tmp) > 0 && tmp[0].InstanceId != "" {
					result = tmp
				}

				for _, stableState := range cs.NodeStableClusterState {
					// If cluster is in NodeStableClusteState, node list will not change
					if object.State == string(stableState) {
						if debugOn() {
							requestMap := make(map[string]interface{})
							requestMap["ClusterId"] = d.Id()
							requestMap["Pagination"] = common.Pagination{PageNumber: pageNumber, PageSize: PageSizeLarge}
							addDebug("GetKubernetesClusterNodes", response, requestInfo, requestMap)
						}
						return nil
					}
				}
				time.Sleep(5 * time.Second)
				return resource.RetryableError(Error("[ERROR] There is no any nodes in kubernetes cluster %s.", d.Id()))
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", DenverdinoAliyungo)
			}

		}

		if d.Get("exclude_autoscaler_nodes") != nil && d.Get("exclude_autoscaler_nodes").(bool) {
			result, err = knockOffAutoScalerNodes(result, meta)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", AlibabaCloudSdkGoERROR)
			}
		}

		for _, node := range result {
			mapping := map[string]interface{}{
				"id":         node.InstanceId,
				"name":       node.InstanceName,
				"private_ip": node.IpAddress[0],
			}
			if node.InstanceRole == "Master" {
				masterNodes = append(masterNodes, mapping)
			} else {
				workerNodes = append(workerNodes, mapping)
			}
		}

		if len(result) < pagination.PageSize {
			break
		}
		pageNumber += 1
	}

	d.Set("master_nodes", masterNodes)
	d.Set("worker_nodes", workerNodes)
	d.Set("worker_number", int64(len(workerNodes)))

	// Get slb information and set connect
	connection := make(map[string]string)
	masterURL := object.MasterURL
	endPoint := make(map[string]string)
	_ = json.Unmarshal([]byte(masterURL), &endPoint)
	connection["api_server_internet"] = endPoint["api_server_endpoint"]
	connection["api_server_intranet"] = endPoint["intranet_api_server_endpoint"]
	if endPoint["api_server_endpoint"] != "" {
		connection["master_public_ip"] = strings.Split(strings.Split(endPoint["api_server_endpoint"], ":")[1], "/")[2]
	}
	if object.Profile != EdgeProfile {
		connection["service_domain"] = fmt.Sprintf("*.%s.%s.alicontainer.com", d.Id(), object.RegionId)
	}

	d.Set("connections", connection)
	d.Set("slb_internet", connection["master_public_ip"])
	if endPoint["intranet_api_server_endpoint"] != "" {
		d.Set("slb_intranet", strings.Split(strings.Split(endPoint["intranet_api_server_endpoint"], ":")[1], "/")[2])
	}

	// set nat gateway
	natRequest := vpc.CreateDescribeNatGatewaysRequest()
	natRequest.VpcId = object.VpcId
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeNatGateways(natRequest)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), natRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(natRequest.GetActionName(), raw, natRequest.RpcRequest, natRequest)
	nat, _ := raw.(*vpc.DescribeNatGatewaysResponse)
	if nat != nil && len(nat.NatGateways.NatGateway) > 0 {
		d.Set("nat_gateway_id", nat.NatGateways.NatGateway[0].NatGatewayId)
	}

	// get cluster conn certs
	// If the cluster is failed, there is no need to get cluster certs
	if object.State == "failed" {
		return nil
	}
	var requestInfo *cs.Client
	var response interface{}
	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.GetClusterCerts(d.Id())
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetClusterCerts", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["Id"] = d.Id()
		addDebug("GetClusterCerts", response, requestInfo, requestMap)
	}
	cert, _ := response.(cs.ClusterCerts)

	// write cluster conn authority to local file
	if ce, ok := d.GetOk("client_cert"); ok && ce.(string) != "" {
		if err := writeToFile(ce.(string), cert.Cert); err != nil {
			return WrapError(err)
		}
	}
	if key, ok := d.GetOk("client_key"); ok && key.(string) != "" {
		if err := writeToFile(key.(string), cert.Key); err != nil {
			return WrapError(err)
		}
	}
	if ca, ok := d.GetOk("cluster_ca_cert"); ok && ca.(string) != "" {
		if err := writeToFile(ca.(string), cert.CA); err != nil {
			return WrapError(err)
		}
	}

	var config *cs.ClusterConfig
	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeClusterUserConfig(d.Id(), false)
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetClusterConfig", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["Id"] = d.Id()
		addDebug("GetClusterConfig", response, requestInfo, requestMap)
	}
	config, _ = response.(*cs.ClusterConfig)

	if file, ok := d.GetOk("kube_config"); ok && file.(string) != "" {
		if err := writeToFile(file.(string), config.Config); err != nil {
			return WrapError(err)
		}
	}

	// write cluster conn authority to tf state
	if err := d.Set("certificate_authority", flattenAlicloudCSCertificate(config)); err != nil {
		return fmt.Errorf("error setting certificate_authority: %s", err)
	}

	return nil
}

func resourceAlicloudCSKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	csService := CsService{meta.(*connectivity.AliyunClient)}
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}

	args := &roacs.DeleteClusterRequest{}
	if v := d.Get("retain_resources"); len(v.([]interface{})) > 0 {
		args.RetainResources = tea.StringSlice(expandStringList(v.([]interface{})))
	}

	_, err = client.DeleteCluster(tea.String(d.Id()), args)
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, ResourceName, "DeleteCluster", AliyunTablestoreGoSdk)
	}

	stateConf := BuildStateConf([]string{"running", "deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildKubernetesArgs(d *schema.ResourceData, meta interface{}) (*cs.DelicatedKubernetesClusterCreationRequest, error) {
	client := meta.(*connectivity.AliyunClient)

	vpcService := VpcService{client}

	var vswitchID string
	if list := expandStringList(d.Get("worker_vswitch_ids").([]interface{})); len(list) > 0 {
		vswitchID = list[0]
	} else {
		vswitchID = ""
	}

	var vpcId string
	if vswitchID != "" {
		vsw, err := vpcService.DescribeVSwitch(vswitchID)
		if err != nil {
			return nil, err
		}
		vpcId = vsw.VpcId
	}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
	}

	addons := make([]cs.Addon, 0)
	if v, ok := d.GetOk("addons"); ok {
		all, ok := v.([]interface{})
		if ok {
			for _, a := range all {
				addon, ok := a.(map[string]interface{})
				if ok {
					addons = append(addons, cs.Addon{
						Name:     addon["name"].(string),
						Config:   addon["config"].(string),
						Disabled: addon["disabled"].(bool),
					})
				}
			}
		}
	}

	var apiAudiences string
	if d.Get("api_audiences") != nil {
		if list := expandStringList(d.Get("api_audiences").([]interface{})); len(list) > 0 {
			apiAudiences = strings.Join(list, ",")
		}
	}

	creationArgs := &cs.DelicatedKubernetesClusterCreationRequest{
		ClusterArgs: cs.ClusterArgs{
			DisableRollback:    true,
			Name:               clusterName,
			DeletionProtection: d.Get("deletion_protection").(bool),
			VpcId:              vpcId,
			// the params below is ok to be empty
			KubernetesVersion:         d.Get("version").(string),
			NodeCidrMask:              strconv.Itoa(d.Get("node_cidr_mask").(int)),
			KeyPair:                   d.Get("key_name").(string),
			ServiceCidr:               d.Get("service_cidr").(string),
			CloudMonitorFlags:         d.Get("install_cloud_monitor").(bool),
			SecurityGroupId:           d.Get("security_group_id").(string),
			IsEnterpriseSecurityGroup: d.Get("is_enterprise_security_group").(bool),
			EndpointPublicAccess:      d.Get("slb_internet_enabled").(bool),
			SnatEntry:                 d.Get("new_nat_gateway").(bool),
			Addons:                    addons,
			ApiAudiences:              apiAudiences,
		},
	}

	if lbSpec, ok := d.GetOk("load_balancer_spec"); ok {
		creationArgs.LoadBalancerSpec = lbSpec.(string)
	}

	if osType, ok := d.GetOk("os_type"); ok {
		creationArgs.OsType = osType.(string)
	}

	if platform, ok := d.GetOk("platform"); ok {
		creationArgs.Platform = platform.(string)
	}

	if timezone, ok := d.GetOk("timezone"); ok {
		creationArgs.Timezone = timezone.(string)
	}

	if clusterDomain, ok := d.GetOk("cluster_domain"); ok {
		creationArgs.ClusterDomain = clusterDomain.(string)
	}

	if customSan, ok := d.GetOk("custom_san"); ok {
		creationArgs.CustomSAN = customSan.(string)
	}

	if imageId, ok := d.GetOk("image_id"); ok {
		creationArgs.ClusterArgs.ImageId = imageId.(string)
	}
	if nodeNameMode, ok := d.GetOk("node_name_mode"); ok {
		creationArgs.ClusterArgs.NodeNameMode = nodeNameMode.(string)
	}
	if saIssuer, ok := d.GetOk("service_account_issuer"); ok {
		creationArgs.ClusterArgs.ServiceAccountIssuer = saIssuer.(string)
	}
	if resourceGroupId, ok := d.GetOk("resource_group_id"); ok {
		creationArgs.ClusterArgs.ResourceGroupId = resourceGroupId.(string)
	}

	if v := d.Get("user_data").(string); v != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(v)
		if base64DecodeError == nil {
			creationArgs.UserData = v
		} else {
			creationArgs.UserData = base64.StdEncoding.EncodeToString([]byte(v))
		}
	}

	if _, ok := d.GetOk("pod_vswitch_ids"); ok {
		creationArgs.PodVswitchIds = expandStringList(d.Get("pod_vswitch_ids").([]interface{}))
	} else {
		creationArgs.ContainerCidr = d.Get("pod_cidr").(string)
	}

	if password := d.Get("password").(string); password == "" {
		if v, ok := d.GetOk("kms_encrypted_password"); ok && v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v.(string), d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return nil, WrapError(err)
			}
			password = decryptResp
		}
		creationArgs.LoginPassword = password
	} else {
		creationArgs.LoginPassword = password
	}

	if tags, err := ConvertCsTags(d); err == nil {
		creationArgs.Tags = tags
	}
	// CA default is empty
	if userCa, ok := d.GetOk("user_ca"); ok {
		userCaContent, err := loadFileContent(userCa.(string))
		if err != nil {
			return nil, fmt.Errorf("reading user_ca file failed %s", err)
		}
		creationArgs.UserCa = string(userCaContent)
	}

	// set proxy mode and default is ipvs
	if proxyMode := d.Get("proxy_mode").(string); proxyMode != "" {
		creationArgs.ProxyMode = cs.ProxyMode(proxyMode)
	} else {
		creationArgs.ProxyMode = cs.ProxyMode(cs.IPVS)
	}

	// dedicated kubernetes must provide master_vswitch_ids
	if _, ok := d.GetOk("master_vswitch_ids"); ok {
		creationArgs.MasterArgs = cs.MasterArgs{
			MasterCount:              len(d.Get("master_vswitch_ids").([]interface{})),
			MasterVSwitchIds:         expandStringList(d.Get("master_vswitch_ids").([]interface{})),
			MasterInstanceTypes:      expandStringList(d.Get("master_instance_types").([]interface{})),
			MasterSystemDiskCategory: aliyungoecs.DiskCategory(d.Get("master_disk_category").(string)),
			MasterSystemDiskSize:     int64(d.Get("master_disk_size").(int)),
			// TODO support other params
		}
	}

	if v, ok := d.GetOk("master_disk_snapshot_policy_id"); ok && v != "" {
		creationArgs.MasterArgs.MasterSnapshotPolicyId = v.(string)
	}

	if v, ok := d.GetOk("master_disk_performance_level"); ok && v != "" {
		creationArgs.MasterArgs.MasterSystemDiskPerformanceLevel = v.(string)
	}

	if v, ok := d.GetOk("master_instance_charge_type"); ok {
		creationArgs.MasterInstanceChargeType = v.(string)
		if creationArgs.MasterInstanceChargeType == string(PrePaid) {
			creationArgs.MasterAutoRenew = d.Get("master_auto_renew").(bool)
			creationArgs.MasterAutoRenewPeriod = d.Get("master_auto_renew_period").(int)
			creationArgs.MasterPeriod = d.Get("master_period").(int)
			creationArgs.MasterPeriodUnit = d.Get("master_period_unit").(string)
		}
	}

	var workerDiskSize int64
	if d.Get("worker_disk_size") != nil {
		workerDiskSize = int64(d.Get("worker_disk_size").(int))
	}

	if v, ok := d.GetOk("worker_vswitch_ids"); ok {
		creationArgs.WorkerArgs.WorkerVSwitchIds = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("worker_instance_types"); ok {
		creationArgs.WorkerArgs.WorkerInstanceTypes = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("worker_number"); ok {
		creationArgs.WorkerArgs.NumOfNodes = int64(v.(int))
	}
	if v, ok := d.GetOk("worker_disk_category"); ok {
		creationArgs.WorkerArgs.WorkerSystemDiskCategory = aliyungoecs.DiskCategory(v.(string))
	}
	if v, ok := d.GetOk("worker_disk_snapshot_policy_id"); ok && v != "" {
		creationArgs.WorkerArgs.WorkerSnapshotPolicyId = v.(string)
	}
	if v, ok := d.GetOk("worker_disk_performance_level"); ok && v != "" {
		creationArgs.WorkerArgs.WorkerSystemDiskPerformanceLevel = v.(string)
	}

	if dds, ok := d.GetOk("worker_data_disks"); ok {
		disks := dds.([]interface{})
		createDataDisks := make([]cs.DataDisk, 0, len(disks))
		for _, e := range disks {
			pack := e.(map[string]interface{})
			dataDisk := cs.DataDisk{
				Size:                 pack["size"].(string),
				DiskName:             pack["name"].(string),
				Category:             pack["category"].(string),
				Device:               pack["device"].(string),
				AutoSnapshotPolicyId: pack["auto_snapshot_policy_id"].(string),
				KMSKeyId:             pack["kms_key_id"].(string),
				Encrypted:            pack["encrypted"].(string),
				PerformanceLevel:     pack["performance_level"].(string),
			}
			createDataDisks = append(createDataDisks, dataDisk)
		}
		creationArgs.WorkerDataDisks = createDataDisks
	}
	if workerDiskSize != 0 {
		creationArgs.WorkerArgs.WorkerSystemDiskSize = workerDiskSize
	}

	if v, ok := d.GetOk("worker_instance_charge_type"); ok {
		creationArgs.WorkerInstanceChargeType = v.(string)
		if creationArgs.WorkerInstanceChargeType == string(PrePaid) {
			creationArgs.WorkerAutoRenew = d.Get("worker_auto_renew").(bool)
			creationArgs.WorkerAutoRenewPeriod = d.Get("worker_auto_renew_period").(int)
			creationArgs.WorkerPeriod = d.Get("worker_period").(int)
			creationArgs.WorkerPeriodUnit = d.Get("worker_period_unit").(string)
		}
	}

	if v, ok := d.GetOk("cluster_spec"); ok {
		creationArgs.ClusterSpec = v.(string)
	}

	if encryptionProviderKey, ok := d.GetOk("encryption_provider_key"); ok {
		creationArgs.EncryptionProviderKey = encryptionProviderKey.(string)
	}

	if rdsInstances, ok := d.GetOk("rds_instances"); ok {
		creationArgs.RdsInstances = expandStringList(rdsInstances.([]interface{}))
	}

	if nodePortRange, ok := d.GetOk("node_port_range"); ok {
		creationArgs.NodePortRange = nodePortRange.(string)
	}

	if runtime, ok := d.GetOk("runtime"); ok {
		if v := runtime.(map[string]interface{}); len(v) > 0 {
			creationArgs.Runtime = expandKubernetesRuntimeConfig(v)
		}
	}

	if taints, ok := d.GetOk("taints"); ok {
		if v := taints.([]interface{}); len(v) > 0 {
			creationArgs.Taints = expandKubernetesTaintsConfig(v)
		}
	}

	// Cluster maintenance window. Effective only in the professional managed cluster
	if v, ok := d.GetOk("maintenance_window"); ok {
		creationArgs.MaintenanceWindow = expandMaintenanceWindowConfig(v.([]interface{}))
	}

	// Configure control plane log. Effective only in the professional managed cluster
	if v, ok := d.GetOk("control_plane_log_components"); ok {
		creationArgs.ControlplaneComponents = expandStringList(v.([]interface{}))
		// ttl default is 30 days
		creationArgs.ControlplaneLogTTL = "30"
	}
	if v, ok := d.GetOk("control_plane_log_ttl"); ok {
		creationArgs.ControlplaneLogTTL = v.(string)
	}
	if v, ok := d.GetOk("control_plane_log_project"); ok {
		creationArgs.ControlplaneLogProject = v.(string)
	}

	return creationArgs, nil
}

func knockOffAutoScalerNodes(nodes []cs.KubernetesNodeType, meta interface{}) ([]cs.KubernetesNodeType, error) {
	log.Printf("[DEBUG] start to knock off auto scaler nodes %++v\n", nodes)
	client := meta.(*connectivity.AliyunClient)
	realNodesMap := make(map[string]ecs.Instance)
	result := make([]cs.KubernetesNodeType, 0)
	instanceIds := make([]interface{}, 0)

	if len(nodes) == 0 {
		return result, nil
	}

	for _, node := range nodes {
		instanceIds = append(instanceIds, node.InstanceId)
	}

	request := ecs.CreateDescribeInstancesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(len(nodes))
	request.InstanceIds = convertListToJsonString(instanceIds)

	// get all the nodes in use
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstances(request)
	})

	if err != nil {
		return result, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*ecs.DescribeInstancesResponse)

	// filter out all autoscaler nodes
	for _, instance := range response.Instances.Instance {
		flags := false
		for _, v := range instance.Tags.Tag {
			if v.TagKey == defaultScalingGroupTag && v.TagValue == "true" {
				flags = true
			}
		}
		if flags != true {
			realNodesMap[instance.InstanceId] = instance
		}
	}

	// get the target node list
	for _, node := range nodes {
		if _, ok := realNodesMap[node.InstanceId]; ok {
			result = append(result, node)
		}
	}

	return result, nil
}

func expandKubernetesTaintsConfig(l []interface{}) []cs.Taint {
	config := []cs.Taint{}

	for _, v := range l {
		if m, ok := v.(map[string]interface{}); ok {
			config = append(config, cs.Taint{
				Key:    m["key"].(string),
				Value:  m["value"].(string),
				Effect: cs.Effect(m["effect"].(string)),
			})
		}
	}

	return config
}

func expandKubernetesRuntimeConfig(l map[string]interface{}) cs.Runtime {
	config := cs.Runtime{}

	if v, ok := l["name"]; ok && v != "" {
		config.Name = v.(string)
	}
	if v, ok := l["version"]; ok && v != "" {
		config.Version = v.(string)
	}

	return config
}

func removeKubernetesNodes(d *schema.ResourceData, meta interface{}) ([]string, error) {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	invoker := NewInvoker()
	// remove nodes count
	o, n := d.GetChange("worker_number")
	count := o.(int) - n.(int)

	var result []cs.KubernetesNodeType
	var response interface{}
	var defaultNodePoolId string

	// get the default nodepool id
	if err := invoker.Run(func() error {
		var err error
		response, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			nodePools, err := csClient.DescribeClusterNodePools(d.Id())
			return *nodePools, err
		})
		return err
	}); err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", DenverdinoAliyungo)
	}

	for _, v := range response.([]cs.NodePoolDetail) {
		if v.BasicNodePool.NodePoolInfo.Name == "default-nodepool" {
			defaultNodePoolId = v.BasicNodePool.NodePoolInfo.NodePoolId
		}
	}

	// list all nodes of default nodepool
	if err := invoker.Run(func() error {
		var err error
		response, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			nodes, _, err := csClient.GetKubernetesClusterNodes(d.Id(), common.Pagination{PageNumber: 1, PageSize: PageSizeLarge}, defaultNodePoolId)
			return nodes, err
		})
		return err
	}); err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", DenverdinoAliyungo)
	}

	ret := response.([]cs.KubernetesNodeType)

	// filter out autoscaler nodes
	var err error
	result, err = knockOffAutoScalerNodes(ret, meta)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetKubernetesClusterNodes", AlibabaCloudSdkGoERROR)
	}

	// filter out Master nodes
	var allNodeName []string
	var allHostName []string
	for _, value := range result {
		if value.InstanceRole == "Worker" {
			allNodeName = append(allNodeName, value.NodeName)
			allHostName = append(allHostName, value.InstanceId)
		}
	}

	// remove nodes
	removeNodesName := allNodeName[:count]
	removeNodesArgs := &cs.DeleteKubernetesClusterNodesRequest{
		Nodes:       removeNodesName,
		ReleaseNode: true,
		DrainNode:   false,
	}
	if err := invoker.Run(func() error {
		var err error
		response, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			resp, err := csClient.DeleteKubernetesClusterNodes(d.Id(), removeNodesArgs)
			return resp, err
		})
		return err
	}); err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, d.Id(), "RemoveClusterNodes", DenverdinoAliyungo)
	}

	stateConf := BuildStateConf([]string{"removing"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return nil, WrapErrorf(err, IdMsg, d.Id())
	}

	d.SetPartial("worker_number")

	return allHostName[len(allHostName)-count:], nil
}

func flattenAlicloudCSCertificate(certificate *cs.ClusterConfig) map[string]string {
	if certificate == nil {
		return map[string]string{}
	}

	kubeConfig := make(map[string]interface{})
	_ = yaml.Unmarshal([]byte(certificate.Config), &kubeConfig)

	m := make(map[string]string)
	m["cluster_cert"] = kubeConfig["clusters"].([]interface{})[0].(map[interface{}]interface{})["cluster"].(map[interface{}]interface{})["certificate-authority-data"].(string)
	m["client_cert"] = kubeConfig["users"].([]interface{})[0].(map[interface{}]interface{})["user"].(map[interface{}]interface{})["client-certificate-data"].(string)
	m["client_key"] = kubeConfig["users"].([]interface{})[0].(map[interface{}]interface{})["user"].(map[interface{}]interface{})["client-key-data"].(string)

	return m
}

// ACK pro maintenance window
func expandMaintenanceWindowConfig(l []interface{}) (config cs.MaintenanceWindow) {
	if len(l) == 0 || l[0] == nil {
		return
	}

	m := l[0].(map[string]interface{})

	if v, ok := m["enable"]; ok {
		config.Enable = v.(bool)
	}
	if v, ok := m["maintenance_time"]; ok && v != "" {
		config.MaintenanceTime = cs.MaintenanceTime(v.(string))
	}
	if v, ok := m["duration"]; ok && v != "" {
		config.Duration = v.(string)
	}
	if v, ok := m["weekly_period"]; ok && v != "" {
		config.WeeklyPeriod = cs.WeeklyPeriod(v.(string))
	}

	return
}

func flattenMaintenanceWindowConfig(config *cs.MaintenanceWindow) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	m = append(m, map[string]interface{}{
		"enable":           config.Enable,
		"maintenance_time": config.MaintenanceTime,
		"duration":         config.Duration,
		"weekly_period":    config.WeeklyPeriod,
	})

	return
}

func modifyMaintenanceWindow(d *schema.ResourceData, meta interface{}, mw cs.MaintenanceWindow) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()

	var response interface{}
	var requestInfo cs.ModifyClusterArgs

	requestInfo.MaintenanceWindow = mw

	if err := invoker.Run(func() error {
		_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.ModifyCluster(d.Id(), &requestInfo)
		})
		return err
	}); err != nil && !IsExpectedErrors(err, []string{"ErrorModifyMaintenanceWindowFailed"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyCluster", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = d.Id()
		requestMap["maintenance_window"] = requestInfo.DeletionProtection
		addDebug("ModifyCluster", response, requestInfo, requestMap)
	}
	d.SetPartial("maintenance_window")

	return nil
}
