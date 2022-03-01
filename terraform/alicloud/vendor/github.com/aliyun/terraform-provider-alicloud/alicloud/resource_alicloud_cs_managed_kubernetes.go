package alicloud

import (
	"fmt"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCSManagedKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSManagedKubernetesCreate,
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
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems: 10,
			},
			"worker_number": {
				Type:     schema.TypeInt,
				Optional: true,
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
			"worker_instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:      PostPaid,
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
			"pod_vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MaxItems: 10,
			},
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
			"enable_ssh": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
				Default:      "none",
				ValidateFunc: validation.StringInSlice([]string{"none", "static"}, false),
			},
			"proxy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ipvs",
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
			"slb_internet_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"load_balancer_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"slb.s1.small", "slb.s2.small", "slb.s2.medium", "slb.s3.small", "slb.s3.medium", "slb.s3.large"}, false),
				Default:      "slb.s1.small",
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
			"cluster_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "cluster.local",
				ForceNew:    true,
				Description: "cluster local domain ",
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
			"encryption_provider_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "disk encryption key, only in ack-pro",
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
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// remove parameters below
			// mix vswitch_ids between master and worker is not a good guidance to create cluster
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems: 3,
				MaxItems: 5,
				Removed:  "Field 'vswitch_ids' has been deprecated from provider version 1.75.0. New field 'master_vswitch_ids' and 'worker_vswitch_ids' replace it.",
			},
			// force update is a high risk operation
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Removed:  "Field 'force_update' has been removed from provider version 1.75.0.",
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
			"cluster_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{KubernetesClusterNetworkTypeFlannel, KubernetesClusterNetworkTypeTerway}, false),
				Removed:      "Field 'cluster_network_type' has been removed from provider version 1.75.0. New field 'addons' replaces it.",
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
			"worker_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Removed:  "Field 'worker_instance_type' has been removed from provider version 1.75.0. New field 'worker_instance_types' replaces it.",
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
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ack.standard", "ack.pro.small"}, false),
			},
			"maintenance_window": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"maintenance_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"duration": {
							Type:     schema.TypeString,
							Required: true,
						},
						"weekly_period": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"control_plane_log_ttl": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"control_plane_log_project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"control_plane_log_components": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func resourceAlicloudCSManagedKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	invoker := NewInvoker()
	csService := CsService{client}
	args, err := buildKubernetesArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *cs.Client
	var response interface{}
	if err := invoker.Run(func() error {
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			args.RegionId = common.Region(client.RegionId)
			args.ClusterType = cs.ManagedKubernetes
			return csClient.CreateManagedKubernetesCluster(&cs.ManagedKubernetesClusterCreationRequest{
				ClusterArgs: args.ClusterArgs,
				WorkerArgs:  args.WorkerArgs,
			})
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_managed_kubernetes", "CreateKubernetesCluster", response)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Args"] = args
		addDebug("CreateKubernetesCluster", response, requestInfo, requestMap)
	}
	cluster, _ := response.(*cs.ClusterCommonResponse)
	d.SetId(cluster.ClusterID)

	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCSKubernetesRead(d, meta)
}

func UpgradeAlicloudKubernetesCluster(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("version") {
		nextVersion := d.Get("version").(string)
		args := &cs.UpgradeClusterArgs{
			Version: nextVersion,
		}

		csService := CsService{client}
		err := csService.UpgradeCluster(d.Id(), args)

		if err != nil {
			return WrapError(err)
		}

		d.SetPartial("version")
	}
	return nil
}

func migrateAlicloudManagedKubernetesCluster(d *schema.ResourceData, meta interface{}) error {
	action := "MigrateCluster"
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}

	migrateClusterRequest := map[string]string{
		"type": "ManagedKubernetes",
		"spec": d.Get("cluster_spec").(string),
	}
	conn, err := meta.(*connectivity.AliyunClient).NewTeaRoaCommonClient(connectivity.OpenAckService)
	if err != nil {
		return WrapError(err)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2015-12-15"), nil, StringPointer("POST"), StringPointer("AK"), String(fmt.Sprintf("/clusters/%s/migrate", d.Id())), nil, nil, migrateClusterRequest, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, response, nil)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	stateConf := BuildStateConf([]string{"migrating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	d.SetPartial("cluster_spec")

	return nil
}

func updateKubernetesClusterTag(d *schema.ResourceData, meta interface{}) error {
	action := "ModifyClusterTags"
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}

	var modifyClusterTagsRequest []cs.Tag
	if tags, err := ConvertCsTags(d); err == nil {
		modifyClusterTagsRequest = tags
	}
	d.SetPartial("tags")
	conn, err := meta.(*connectivity.AliyunClient).NewTeaRoaCommonClient(connectivity.OpenAckService)
	if err != nil {
		return WrapError(err)
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequestWithAction(StringPointer(action), StringPointer("2015-12-15"), nil, StringPointer("POST"), StringPointer("AK"), String(fmt.Sprintf("/clusters/%s/tags", d.Id())), nil, nil, modifyClusterTagsRequest, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, response, nil)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	stateConf := BuildStateConf([]string{"updating"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}
