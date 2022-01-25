package alicloud

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"time"

	aliyungoecs "github.com/denverdino/aliyungo/ecs"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	EdgeKubernetesDefaultTimeoutInMinutes = 60
	EdgeProfile                           = "Edge"
)

func resourceAlicloudCSEdgeKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSEdgeKubernetesCreate,
		Read:   resourceAlicloudCSKubernetesRead,
		Update: resourceAlicloudCSEdgeKubernetesUpdate,
		Delete: resourceAlicloudCSKubernetesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(EdgeKubernetesDefaultTimeoutInMinutes * time.Minute),
			Update: schema.DefaultTimeout(EdgeKubernetesDefaultTimeoutInMinutes * time.Minute),
			Delete: schema.DefaultTimeout(EdgeKubernetesDefaultTimeoutInMinutes * time.Minute),
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
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
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
			//cloud worker number
			"worker_number": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
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
				ValidateFunc: validation.StringInSlice([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD), string(DiskCloudESSD)}, false),
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
			"proxy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ipvs",
				ValidateFunc: validation.StringInSlice([]string{"iptables", "ipvs"}, false),
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
			// global configurations
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
				ConflictsWith: []string{"key_name"},
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password"},
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
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
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
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"rds_instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// computed parameters start
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
			// computed params end

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
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
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

func resourceAlicloudCSEdgeKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
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
			args.Profile = EdgeProfile
			return csClient.CreateManagedKubernetesCluster(&cs.ManagedKubernetesClusterCreationRequest{
				ClusterArgs: args.ClusterArgs,
				WorkerArgs:  args.WorkerArgs,
			})
		})
		response = raw
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cs_edge_kubernetes", "CreateKubernetesCluster", response)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Args"] = args
		addDebug("CreateKubernetesCluster", response, requestInfo, requestMap)
	}
	cluster, _ := response.(*cs.ClusterCommonResponse)
	d.SetId(cluster.ClusterID)

	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Minute, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCSKubernetesRead(d, meta)
}

func resourceAlicloudCSEdgeKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	d.Partial(true)
	invoker := NewInvoker()
	//scale up cloud worker nodes
	var resp interface{}
	if d.HasChanges("worker_number") {
		oldV, newV := d.GetChange("worker_number")
		oldValue, ok := oldV.(int)
		if !ok {
			return WrapErrorf(fmt.Errorf("worker_number old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := newV.(int)
		if !ok {
			return WrapErrorf(fmt.Errorf("worker_number new value can not be parsed"), "parseError %d", oldValue)
		}

		// Edge cluster node support remove nodes.
		if newValue < oldValue {
			return WrapErrorf(fmt.Errorf("worker_number can not be less than before"), "scaleOutCloudWorkersFailed %d:%d", newValue, oldValue)
		}

		// scale out cluster.
		if newValue > oldValue {
			password := d.Get("password").(string)
			keyPair := d.Get("key_name").(string)

			args := &cs.ScaleOutKubernetesClusterRequest{
				KeyPair:             keyPair,
				LoginPassword:       password,
				Count:               int64(newValue) - int64(oldValue),
				WorkerVSwitchIds:    expandStringList(d.Get("worker_vswitch_ids").([]interface{})),
				WorkerInstanceTypes: expandStringList(d.Get("worker_instance_types").([]interface{})),
			}

			if userData, ok := d.GetOk("user_data"); ok {
				_, base64DecodeError := base64.StdEncoding.DecodeString(userData.(string))
				if base64DecodeError == nil {
					args.UserData = userData.(string)
				} else {
					args.UserData = base64.StdEncoding.EncodeToString([]byte(userData.(string)))
				}
			}

			if imageID, ok := d.GetOk("image_id"); ok {
				args.ImageId = imageID.(string)

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

			if d.HasChange("tags") && !d.IsNewResource() {
				if tags, err := ConvertCsTags(d); err == nil {
					args.Tags = tags
				}
				d.SetPartial("tags")
			}

			if err := invoker.Run(func() error {
				var err error
				resp, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
					resp, err := csClient.ScaleOutKubernetesCluster(d.Id(), args)
					return resp, err
				})
				return err
			}); err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ScaleOutCloudWorkers", DenverdinoAliyungo)
			}
			if debugOn() {
				resizeRequestMap := make(map[string]interface{})
				resizeRequestMap["ClusterId"] = d.Id()
				resizeRequestMap["Args"] = args
				addDebug("ResizeKubernetesCluster", resp, resizeRequestMap)
			}
			stateConf := BuildStateConf([]string{"scaling"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))

			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("worker_data_disks")
			d.SetPartial("worker_number")
			d.SetPartial("worker_disk_category")
			d.SetPartial("worker_disk_size")
			d.SetPartial("worker_disk_snapshot_policy_id")
			d.SetPartial("worker_disk_performance_level")
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

	// modify cluster tag
	if d.HasChange("tags") {
		err := updateKubernetesClusterTag(d, meta)
		if err != nil {
			return WrapErrorf(err, ResponseCodeMsg, d.Id(), "ModifyClusterTags", AlibabaCloudSdkGoERROR)
		}
	}
	d.SetPartial("tags")

	// upgrade cluster version
	err := UpgradeAlicloudKubernetesCluster(d, meta)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpgradeClusterVersion", DenverdinoAliyungo)
	}
	d.Partial(false)
	return resourceAlicloudCSKubernetesRead(d, meta)
}
