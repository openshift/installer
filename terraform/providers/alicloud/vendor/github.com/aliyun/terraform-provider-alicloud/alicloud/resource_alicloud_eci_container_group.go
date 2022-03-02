package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEciContainerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEciContainerGroupCreate,
		Read:   resourceAlicloudEciContainerGroupRead,
		Update: resourceAlicloudEciContainerGroupUpdate,
		Delete: resourceAlicloudEciContainerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"container_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"containers": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"args": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"commands": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cpu": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  0,
						},
						"environment_vars": {
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
								},
							},
						},
						"gpu": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"image": {
							Type:     schema.TypeString,
							Required: true,
						},
						"image_pull_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "IfNotPresent",
						},
						"memory": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  0,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
							ForceNew: true,
						},
						"ready": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"restart_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_mounts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
						"working_dir": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"cpu": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  2,
			},
			"dns_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name_servers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"options": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"searches": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"eci_security_context": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sysctls": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
							ForceNew: true,
						},
					},
				},
				ForceNew: true,
			},
			"host_aliases": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostnames": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							ForceNew: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
				ForceNew: true,
			},
			"init_containers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"args": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"commands": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cpu": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  0,
						},
						"environment_vars": {
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
								},
							},
						},
						"gpu": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"image": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_pull_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "IfNotPresent",
						},
						"memory": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  0,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
							ForceNew: true,
						},
						"ready": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"restart_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_mounts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
						"working_dir": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"memory": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  4,
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"restart_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Always",
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_file_volume_config_file_to_paths": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"path": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"disk_volume_disk_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disk_volume_fs_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"flex_volume_driver": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"flex_volume_fs_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"flex_volume_options": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"nfs_volume_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"nfs_volume_read_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"nfs_volume_server": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"image_registry_credential": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:     schema.TypeString,
							Required: true,
						},
						"server": {
							Type:     schema.TypeString,
							Required: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudEciContainerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eciService := EciService{client}
	var response map[string]interface{}
	action := "CreateContainerGroup"
	request := make(map[string]interface{})
	conn, err := client.NewEciClient()
	if err != nil {
		return WrapError(err)
	}
	request["ContainerGroupName"] = d.Get("container_group_name")
	Containers := make([]map[string]interface{}, len(d.Get("containers").(*schema.Set).List()))
	for i, ContainersValue := range d.Get("containers").(*schema.Set).List() {
		ContainersMap := ContainersValue.(map[string]interface{})
		Containers[i] = make(map[string]interface{})
		Containers[i]["Arg"] = ContainersMap["args"]
		Containers[i]["Command"] = ContainersMap["commands"]
		Containers[i]["Cpu"] = ContainersMap["cpu"]
		EnvironmentVars := make([]map[string]interface{}, len(ContainersMap["environment_vars"].([]interface{})))
		for i, EnvironmentVarsValue := range ContainersMap["environment_vars"].([]interface{}) {
			EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
			EnvironmentVars[i] = make(map[string]interface{})
			EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
			EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
		}
		Containers[i]["EnvironmentVar"] = EnvironmentVars

		Containers[i]["Gpu"] = ContainersMap["gpu"]
		Containers[i]["Image"] = ContainersMap["image"]
		Containers[i]["ImagePullPolicy"] = ContainersMap["image_pull_policy"]
		Containers[i]["Memory"] = ContainersMap["memory"]
		Containers[i]["Name"] = ContainersMap["name"]
		Ports := make([]map[string]interface{}, len(ContainersMap["ports"].([]interface{})))
		for i, PortsValue := range ContainersMap["ports"].([]interface{}) {
			PortsMap := PortsValue.(map[string]interface{})
			Ports[i] = make(map[string]interface{})
			Ports[i]["Port"] = PortsMap["port"]
			Ports[i]["Protocol"] = PortsMap["protocol"]
		}
		Containers[i]["Port"] = Ports

		VolumeMounts := make([]map[string]interface{}, len(ContainersMap["volume_mounts"].([]interface{})))
		for i, VolumeMountsValue := range ContainersMap["volume_mounts"].([]interface{}) {
			VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
			VolumeMounts[i] = make(map[string]interface{})
			VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
			VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
		}
		Containers[i]["VolumeMount"] = VolumeMounts

		Containers[i]["WorkingDir"] = ContainersMap["working_dir"]
	}
	request["Container"] = Containers

	if v, ok := d.GetOk("cpu"); ok {
		request["Cpu"] = v
	}

	if v, ok := d.GetOk("dns_config"); ok {
		if v != nil {
			dnsConfigMap := make(map[string]interface{})
			for _, dnsConfig := range v.(*schema.Set).List() {
				dnsConfigArg := dnsConfig.(map[string]interface{})
				dnsConfigMap["NameServers"] = dnsConfigArg["name_servers"]
				if dnsConfigArg["options"] != nil {
					optionsMaps := make([]map[string]interface{}, 0)
					for _, options := range dnsConfigArg["options"].([]interface{}) {
						optionsMap := make(map[string]interface{})
						optionsArg := options.(map[string]interface{})
						optionsMap["Name"] = optionsArg["name"]
						optionsMap["Value"] = optionsArg["value"]
						optionsMaps = append(optionsMaps, optionsMap)
					}
					dnsConfigMap["Options"] = optionsMaps
				}
				dnsConfigMap["Searches"] = dnsConfigArg["searches"]
			}
			request["DnsConfig"] = dnsConfigMap
		}
	}
	if v, ok := d.GetOk("eci_security_context"); ok {
		if v != nil {
			eciSecurityContextMap := make(map[string]interface{})
			for _, eciSecurityContext := range v.(*schema.Set).List() {
				eciSecurityContextArg := eciSecurityContext.(map[string]interface{})
				if eciSecurityContextArg["sysctls"] != nil {
					sysctlsMaps := make([]map[string]interface{}, 0)
					for _, sysctls := range eciSecurityContextArg["sysctls"].([]interface{}) {
						sysctlsMap := make(map[string]interface{})
						sysctlsArg := sysctls.(map[string]interface{})
						sysctlsMap["Name"] = sysctlsArg["name"]
						sysctlsMap["Value"] = sysctlsArg["value"]
						sysctlsMaps = append(sysctlsMaps, sysctlsMap)
					}
					eciSecurityContextMap["Sysctls"] = sysctlsMaps
				}
			}
			request["EciSecurityContext"] = eciSecurityContextMap
		}
	}
	if v, ok := d.GetOk("host_aliases"); ok {
		HostAliases := make([]map[string]interface{}, len(v.([]interface{})))
		for i, HostAliasesValue := range v.([]interface{}) {
			HostAliasesMap := HostAliasesValue.(map[string]interface{})
			HostAliases[i] = make(map[string]interface{})
			HostAliases[i]["Hostname"] = HostAliasesMap["hostnames"]
			HostAliases[i]["Ip"] = HostAliasesMap["ip"]
		}
		request["HostAliase"] = HostAliases

	}

	if v, ok := d.GetOk("init_containers"); ok {
		InitContainers := make([]map[string]interface{}, len(v.([]interface{})))
		for i, InitContainersValue := range v.([]interface{}) {
			InitContainersMap := InitContainersValue.(map[string]interface{})
			InitContainers[i] = make(map[string]interface{})
			InitContainers[i]["Arg"] = InitContainersMap["args"]
			InitContainers[i]["Command"] = InitContainersMap["commands"]
			InitContainers[i]["Cpu"] = InitContainersMap["cpu"]
			EnvironmentVars := make([]map[string]interface{}, len(InitContainersMap["environment_vars"].([]interface{})))
			for i, EnvironmentVarsValue := range InitContainersMap["environment_vars"].([]interface{}) {
				EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
				EnvironmentVars[i] = make(map[string]interface{})
				EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
				EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
			}
			InitContainers[i]["EnvironmentVar"] = EnvironmentVars

			InitContainers[i]["Gpu"] = InitContainersMap["gpu"]
			InitContainers[i]["Image"] = InitContainersMap["image"]
			InitContainers[i]["ImagePullPolicy"] = InitContainersMap["image_pull_policy"]
			InitContainers[i]["Memory"] = InitContainersMap["memory"]
			InitContainers[i]["Name"] = InitContainersMap["name"]
			Ports := make([]map[string]interface{}, len(InitContainersMap["ports"].([]interface{})))
			for i, PortsValue := range InitContainersMap["ports"].([]interface{}) {
				PortsMap := PortsValue.(map[string]interface{})
				Ports[i] = make(map[string]interface{})
				Ports[i]["Port"] = PortsMap["port"]
				Ports[i]["Protocol"] = PortsMap["protocol"]
			}
			InitContainers[i]["Port"] = Ports

			VolumeMounts := make([]map[string]interface{}, len(InitContainersMap["volume_mounts"].([]interface{})))
			for i, VolumeMountsValue := range InitContainersMap["volume_mounts"].([]interface{}) {
				VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
				VolumeMounts[i] = make(map[string]interface{})
				VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
				VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
			}
			InitContainers[i]["VolumeMount"] = VolumeMounts

			InitContainers[i]["WorkingDir"] = InitContainersMap["working_dir"]
		}
		request["InitContainer"] = InitContainers

	}

	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}

	if v, ok := d.GetOk("memory"); ok {
		request["Memory"] = v
	}

	if v, ok := d.GetOk("ram_role_name"); ok {
		request["RamRoleName"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("restart_policy"); ok {
		request["RestartPolicy"] = v
	}

	request["SecurityGroupId"] = d.Get("security_group_id")
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if v, ok := d.GetOk("volumes"); ok {
		Volumes := make([]map[string]interface{}, len(v.([]interface{})))
		for i, VolumesValue := range v.([]interface{}) {
			VolumesMap := VolumesValue.(map[string]interface{})
			Volumes[i] = make(map[string]interface{})
			ConfigFileVolumeConfigFileToPaths := make([]map[string]interface{}, len(VolumesMap["config_file_volume_config_file_to_paths"].([]interface{})))
			for i, ConfigFileVolumeConfigFileToPathsValue := range VolumesMap["config_file_volume_config_file_to_paths"].([]interface{}) {
				ConfigFileVolumeConfigFileToPathsMap := ConfigFileVolumeConfigFileToPathsValue.(map[string]interface{})
				ConfigFileVolumeConfigFileToPaths[i] = make(map[string]interface{})
				ConfigFileVolumeConfigFileToPaths[i]["Content"] = ConfigFileVolumeConfigFileToPathsMap["content"]
				ConfigFileVolumeConfigFileToPaths[i]["Path"] = ConfigFileVolumeConfigFileToPathsMap["path"]
			}
			Volumes[i]["ConfigFileVolume.ConfigFileToPath"] = ConfigFileVolumeConfigFileToPaths

			Volumes[i]["DiskVolume.DiskId"] = VolumesMap["disk_volume_disk_id"]
			Volumes[i]["DiskVolume.FsType"] = VolumesMap["disk_volume_fs_type"]
			Volumes[i]["FlexVolume.Driver"] = VolumesMap["flex_volume_driver"]
			Volumes[i]["FlexVolume.FsType"] = VolumesMap["flex_volume_fs_type"]
			Volumes[i]["FlexVolume.Options"] = VolumesMap["flex_volume_options"]
			Volumes[i]["NFSVolume.Path"] = VolumesMap["nfs_volume_path"]
			Volumes[i]["NFSVolume.Server"] = VolumesMap["nfs_volume_server"]
			Volumes[i]["Name"] = VolumesMap["name"]
			Volumes[i]["Type"] = VolumesMap["type"]
		}
		request["Volume"] = Volumes

	}

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return WrapError(err)
		}
		request["VSwitchId"] = vswitchId
		if request["ZoneId"] == nil {
			request["ZoneId"] = vsw.ZoneId
		}
	}
	if v, ok := d.GetOk("image_registry_credential"); ok {
		imageRegisryCredentialMaps := make([]map[string]interface{}, 0)
		for _, raw := range v.(*schema.Set).List() {
			obj := raw.(map[string]interface{})
			imageRegisryCredentialMaps = append(imageRegisryCredentialMaps, map[string]interface{}{
				"Password": obj["password"],
				"Server":   obj["server"],
				"UserName": obj["user_name"],
			})
		}
		request["ImageRegistryCredential"] = imageRegisryCredentialMaps
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("CreateContainerGroup")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-08"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eci_container_group", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(response["ContainerGroupId"]))
	stateConf := BuildStateConf([]string{}, []string{"Running", "Succeeded"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, eciService.EciContainerGroupStateRefreshFunc(d.Id(), []string{"Failed", "ScheduleFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEciContainerGroupRead(d, meta)
}
func resourceAlicloudEciContainerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eciService := EciService{client}
	object, err := eciService.DescribeEciContainerGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eci_openapi_container_group eciService.DescribeEciContainerGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("container_group_name", object["ContainerGroupName"])

	containers := make([]map[string]interface{}, 0)
	if containersList, ok := object["Containers"].([]interface{}); ok {
		for _, v := range containersList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"args":              m1["Args"],
					"commands":          m1["Commands"],
					"cpu":               m1["Cpu"],
					"gpu":               m1["Gpu"],
					"image":             m1["Image"],
					"image_pull_policy": m1["ImagePullPolicy"],
					"memory":            m1["Memory"],
					"name":              m1["Name"],
					"ready":             m1["Ready"],
					"restart_count":     m1["RestartCount"],
					"working_dir":       m1["WorkingDir"],
				}
				if m1["EnvironmentVars"] != nil {
					environmentVarsMaps := make([]map[string]interface{}, 0)
					for _, environmentVarsValue := range m1["EnvironmentVars"].([]interface{}) {
						environmentVars := environmentVarsValue.(map[string]interface{})
						environmentVarsMap := map[string]interface{}{
							"key":   environmentVars["Key"],
							"value": environmentVars["Value"],
						}
						environmentVarsMaps = append(environmentVarsMaps, environmentVarsMap)
					}
					temp1["environment_vars"] = environmentVarsMaps
				}
				if m1["Ports"] != nil {
					portsMaps := make([]map[string]interface{}, 0)
					for _, portsValue := range m1["Ports"].([]interface{}) {
						ports := portsValue.(map[string]interface{})
						portsMap := map[string]interface{}{
							"port":     ports["Port"],
							"protocol": ports["Protocol"],
						}
						portsMaps = append(portsMaps, portsMap)
					}
					temp1["ports"] = portsMaps
				}
				if m1["VolumeMounts"] != nil {
					volumeMountsMaps := make([]map[string]interface{}, 0)
					for _, volumeMountsValue := range m1["VolumeMounts"].([]interface{}) {
						volumeMounts := volumeMountsValue.(map[string]interface{})
						volumeMountsMap := map[string]interface{}{
							"mount_path": volumeMounts["MountPath"],
							"name":       volumeMounts["Name"],
							"read_only":  volumeMounts["ReadOnly"],
						}
						volumeMountsMaps = append(volumeMountsMaps, volumeMountsMap)
					}
					temp1["volume_mounts"] = volumeMountsMaps
				}
				containers = append(containers, temp1)

			}
		}
	}
	if err := d.Set("containers", containers); err != nil {
		return WrapError(err)
	}
	d.Set("cpu", object["Cpu"])

	dnsConfigSli := make([]map[string]interface{}, 0)
	if len(object["DnsConfig"].(map[string]interface{})) > 0 {
		dnsConfig := object["DnsConfig"]
		dnsConfigMap := make(map[string]interface{})
		dnsConfigMap["name_servers"] = dnsConfig.(map[string]interface{})["NameServers"]

		optionsSli := make([]map[string]interface{}, 0)
		if len(dnsConfig.(map[string]interface{})["Options"].([]interface{})) > 0 {
			for _, options := range dnsConfig.(map[string]interface{})["Options"].([]interface{}) {
				optionsMap := make(map[string]interface{})
				optionsMap["name"] = options.(map[string]interface{})["Name"]
				optionsMap["value"] = options.(map[string]interface{})["Value"]
				optionsSli = append(optionsSli, optionsMap)
			}
		}
		dnsConfigMap["options"] = optionsSli
		dnsConfigMap["searches"] = dnsConfig.(map[string]interface{})["Searches"]
		dnsConfigSli = append(dnsConfigSli, dnsConfigMap)
	}
	d.Set("dns_config", dnsConfigSli)

	eciSecurityContextSli := make([]map[string]interface{}, 0)
	if len(object["EciSecurityContext"].(map[string]interface{})) > 0 {
		eciSecurityContext := object["EciSecurityContext"]
		eciSecurityContextMap := make(map[string]interface{})

		sysctlsSli := make([]map[string]interface{}, 0)
		if len(eciSecurityContext.(map[string]interface{})["Sysctls"].([]interface{})) > 0 {
			for _, sysctls := range eciSecurityContext.(map[string]interface{})["Sysctls"].([]interface{}) {
				sysctlsMap := make(map[string]interface{})
				sysctlsMap["name"] = sysctls.(map[string]interface{})["Name"]
				sysctlsMap["value"] = sysctls.(map[string]interface{})["Value"]
				sysctlsSli = append(sysctlsSli, sysctlsMap)
			}
		}
		eciSecurityContextMap["sysctls"] = sysctlsSli
		eciSecurityContextSli = append(eciSecurityContextSli, eciSecurityContextMap)
	}
	d.Set("eci_security_context", eciSecurityContextSli)

	hostAliases := make([]map[string]interface{}, 0)
	if hostAliasesList, ok := object["HostAliases"].([]interface{}); ok {
		for _, v := range hostAliasesList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"hostnames": m1["Hostnames"],
					"ip":        m1["Ip"],
				}
				hostAliases = append(hostAliases, temp1)

			}
		}
	}
	if err := d.Set("host_aliases", hostAliases); err != nil {
		return WrapError(err)
	}

	initContainers := make([]map[string]interface{}, 0)
	if initContainersList, ok := object["InitContainers"].([]interface{}); ok {
		for _, v := range initContainersList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"args":              m1["Args"],
					"commands":          m1["Command"],
					"cpu":               m1["Cpu"],
					"gpu":               m1["Gpu"],
					"image":             m1["Image"],
					"image_pull_policy": m1["ImagePullPolicy"],
					"memory":            m1["Memory"],
					"name":              m1["Name"],
					"ready":             m1["Ready"],
					"restart_count":     m1["RestartCount"],
					"working_dir":       m1["WorkingDir"],
				}
				if m1["EnvironmentVars"] != nil {
					environmentVarsMaps := make([]map[string]interface{}, 0)
					for _, environmentVarsValue := range m1["EnvironmentVars"].([]interface{}) {
						environmentVars := environmentVarsValue.(map[string]interface{})
						environmentVarsMap := map[string]interface{}{
							"key":   environmentVars["Key"],
							"value": environmentVars["Value"],
						}
						environmentVarsMaps = append(environmentVarsMaps, environmentVarsMap)
					}
					temp1["environment_vars"] = environmentVarsMaps
				}
				if m1["Ports"] != nil {
					portsMaps := make([]map[string]interface{}, 0)
					for _, portsValue := range m1["Ports"].([]interface{}) {
						ports := portsValue.(map[string]interface{})
						portsMap := map[string]interface{}{
							"port":     ports["Port"],
							"protocol": ports["Protocol"],
						}
						portsMaps = append(portsMaps, portsMap)
					}
					temp1["ports"] = portsMaps
				}
				if m1["VolumeMounts"] != nil {
					volumeMountsMaps := make([]map[string]interface{}, 0)
					for _, volumeMountsValue := range m1["VolumeMounts"].([]interface{}) {
						volumeMounts := volumeMountsValue.(map[string]interface{})
						volumeMountsMap := map[string]interface{}{
							"mount_path": volumeMounts["MountPath"],
							"name":       volumeMounts["Name"],
							"read_only":  volumeMounts["ReadOnly"],
						}
						volumeMountsMaps = append(volumeMountsMaps, volumeMountsMap)
					}
					temp1["volume_mounts"] = volumeMountsMaps
				}
				initContainers = append(initContainers, temp1)

			}
		}
	}
	if err := d.Set("init_containers", initContainers); err != nil {
		return WrapError(err)
	}
	d.Set("instance_type", object["InstanceType"])
	d.Set("memory", object["Memory"])
	d.Set("ram_role_name", object["RamRoleName"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("restart_policy", object["RestartPolicy"])
	d.Set("security_group_id", object["SecurityGroupId"])
	d.Set("status", object["Status"])
	d.Set("tags", tagsToMap(object["Tags"]))
	d.Set("vswitch_id", object["VSwitchId"])

	volumes := make([]map[string]interface{}, 0)
	if volumesList, ok := object["Volumes"].([]interface{}); ok {
		for _, v := range volumesList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"disk_volume_disk_id":  m1["DiskVolumeDiskId"],
					"disk_volume_fs_type":  m1["DiskVolumeFsType"],
					"flex_volume_driver":   m1["FlexVolumeDriver"],
					"flex_volume_fs_type":  m1["FlexVolumeFsType"],
					"flex_volume_options":  m1["FlexVolumeOptions"],
					"nfs_volume_path":      m1["NFSVolumePath"],
					"nfs_volume_read_only": m1["NFSVolumeReadOnly"],
					"nfs_volume_server":    m1["NFSVolumeServer"],
					"name":                 m1["Name"],
					"type":                 m1["Type"],
				}
				if m1["ConfigFileVolumeConfigFileToPaths"] != nil {
					configFileVolumeConfigFileToPathsMaps := make([]map[string]interface{}, 0)
					for _, configFileVolumeConfigFileToPathsValue := range m1["ConfigFileVolumeConfigFileToPaths"].([]interface{}) {
						configFileVolumeConfigFileToPaths := configFileVolumeConfigFileToPathsValue.(map[string]interface{})
						configFileVolumeConfigFileToPathsMap := map[string]interface{}{
							"content": configFileVolumeConfigFileToPaths["Content"],
							"path":    configFileVolumeConfigFileToPaths["Path"],
						}
						configFileVolumeConfigFileToPathsMaps = append(configFileVolumeConfigFileToPathsMaps, configFileVolumeConfigFileToPathsMap)
					}
					temp1["config_file_volume_config_file_to_paths"] = configFileVolumeConfigFileToPathsMaps
				}
				volumes = append(volumes, temp1)

			}
		}
	}
	if err := d.Set("volumes", volumes); err != nil {
		return WrapError(err)
	}
	d.Set("zone_id", object["ZoneId"])
	return nil
}
func resourceAlicloudEciContainerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eciService := EciService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ContainerGroupId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("containers") {
		update = true
		Containers := make([]map[string]interface{}, len(d.Get("containers").(*schema.Set).List()))
		for i, ContainersValue := range d.Get("containers").(*schema.Set).List() {
			ContainersMap := ContainersValue.(map[string]interface{})
			Containers[i] = make(map[string]interface{})
			Containers[i]["Arg"] = ContainersMap["args"]
			Containers[i]["Command"] = ContainersMap["commands"]
			Containers[i]["Cpu"] = ContainersMap["cpu"]
			EnvironmentVars := make([]map[string]interface{}, len(ContainersMap["environment_vars"].([]interface{})))
			for i, EnvironmentVarsValue := range ContainersMap["environment_vars"].([]interface{}) {
				EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
				EnvironmentVars[i] = make(map[string]interface{})
				EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
				EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
			}
			Containers[i]["EnvironmentVar"] = EnvironmentVars

			Containers[i]["Gpu"] = ContainersMap["gpu"]
			Containers[i]["Image"] = ContainersMap["image"]
			Containers[i]["ImagePullPolicy"] = ContainersMap["image_pull_policy"]
			Containers[i]["Memory"] = ContainersMap["memory"]
			Containers[i]["Name"] = ContainersMap["name"]
			Ports := make([]map[string]interface{}, len(ContainersMap["ports"].([]interface{})))
			for i, PortsValue := range ContainersMap["ports"].([]interface{}) {
				PortsMap := PortsValue.(map[string]interface{})
				Ports[i] = make(map[string]interface{})
				Ports[i]["Port"] = PortsMap["port"]
				Ports[i]["Protocol"] = PortsMap["protocol"]
			}
			Containers[i]["Port"] = Ports

			VolumeMounts := make([]map[string]interface{}, len(ContainersMap["volume_mounts"].([]interface{})))
			for i, VolumeMountsValue := range ContainersMap["volume_mounts"].([]interface{}) {
				VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
				VolumeMounts[i] = make(map[string]interface{})
				VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
				VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
			}
			Containers[i]["VolumeMount"] = VolumeMounts

			Containers[i]["WorkingDir"] = ContainersMap["working_dir"]
		}
		request["Container"] = Containers

	}
	if d.HasChange("cpu") {
		update = true
		request["Cpu"] = d.Get("cpu")
	}
	if d.HasChange("dns_config") {
		update = true
		if d.Get("dns_config") != nil {
			dnsConfigMap := make(map[string]interface{})
			for _, dnsConfig := range d.Get("dns_config").(*schema.Set).List() {
				dnsConfigArg := dnsConfig.(map[string]interface{})
				dnsConfigMap["NameServers"] = dnsConfigArg["name_servers"]
				if dnsConfigArg["options"] != nil {
					optionsMaps := make([]map[string]interface{}, 0)
					for _, options := range dnsConfigArg["options"].([]interface{}) {
						optionsMap := make(map[string]interface{})
						optionsArg := options.(map[string]interface{})
						optionsMap["Name"] = optionsArg["name"]
						optionsMap["Value"] = optionsArg["value"]
						optionsMaps = append(optionsMaps, optionsMap)
					}
					dnsConfigMap["Options"] = optionsMaps
				}
				dnsConfigMap["Searches"] = dnsConfigArg["searches"]
			}
			request["DnsConfig"] = dnsConfigMap
		}
	}
	if d.HasChange("init_containers") {
		update = true
		InitContainers := make([]map[string]interface{}, len(d.Get("init_containers").([]interface{})))
		for i, InitContainersValue := range d.Get("init_containers").([]interface{}) {
			InitContainersMap := InitContainersValue.(map[string]interface{})
			InitContainers[i] = make(map[string]interface{})
			InitContainers[i]["Arg"] = InitContainersMap["args"]
			InitContainers[i]["Command"] = InitContainersMap["commands"]
			InitContainers[i]["Cpu"] = InitContainersMap["cpu"]
			EnvironmentVars := make([]map[string]interface{}, len(InitContainersMap["environment_vars"].([]interface{})))
			for i, EnvironmentVarsValue := range InitContainersMap["environment_vars"].([]interface{}) {
				EnvironmentVarsMap := EnvironmentVarsValue.(map[string]interface{})
				EnvironmentVars[i] = make(map[string]interface{})
				EnvironmentVars[i]["Key"] = EnvironmentVarsMap["key"]
				EnvironmentVars[i]["Value"] = EnvironmentVarsMap["value"]
			}
			InitContainers[i]["EnvironmentVar"] = EnvironmentVars

			InitContainers[i]["Gpu"] = InitContainersMap["gpu"]
			InitContainers[i]["Image"] = InitContainersMap["image"]
			InitContainers[i]["ImagePullPolicy"] = InitContainersMap["image_pull_policy"]
			InitContainers[i]["Memory"] = InitContainersMap["memory"]
			InitContainers[i]["Name"] = InitContainersMap["name"]
			Ports := make([]map[string]interface{}, len(InitContainersMap["ports"].([]interface{})))
			for i, PortsValue := range InitContainersMap["ports"].([]interface{}) {
				PortsMap := PortsValue.(map[string]interface{})
				Ports[i] = make(map[string]interface{})
				Ports[i]["Port"] = PortsMap["port"]
				Ports[i]["Protocol"] = PortsMap["protocol"]
			}
			InitContainers[i]["Port"] = Ports

			VolumeMounts := make([]map[string]interface{}, len(InitContainersMap["volume_mounts"].([]interface{})))
			for i, VolumeMountsValue := range InitContainersMap["volume_mounts"].([]interface{}) {
				VolumeMountsMap := VolumeMountsValue.(map[string]interface{})
				VolumeMounts[i] = make(map[string]interface{})
				VolumeMounts[i]["MountPath"] = VolumeMountsMap["mount_path"]
				VolumeMounts[i]["Name"] = VolumeMountsMap["name"]
			}
			InitContainers[i]["VolumeMount"] = VolumeMounts

			InitContainers[i]["WorkingDir"] = InitContainersMap["working_dir"]
		}
		request["InitContainer"] = InitContainers

	}
	if d.HasChange("memory") {
		update = true
		request["Memory"] = d.Get("memory")
	}
	if d.HasChange("restart_policy") {
		update = true
		request["RestartPolicy"] = d.Get("restart_policy")
	}
	if d.HasChange("tags") {
		update = true
		count := 1
		for key, value := range d.Get("tags").(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if d.HasChange("volumes") {
		update = true
		Volumes := make([]map[string]interface{}, len(d.Get("volumes").([]interface{})))
		for i, VolumesValue := range d.Get("volumes").([]interface{}) {
			VolumesMap := VolumesValue.(map[string]interface{})
			Volumes[i] = make(map[string]interface{})
			ConfigFileVolumeConfigFileToPaths := make([]map[string]interface{}, len(VolumesMap["config_file_volume_config_file_to_paths"].([]interface{})))
			for i, ConfigFileVolumeConfigFileToPathsValue := range VolumesMap["config_file_volume_config_file_to_paths"].([]interface{}) {
				ConfigFileVolumeConfigFileToPathsMap := ConfigFileVolumeConfigFileToPathsValue.(map[string]interface{})
				ConfigFileVolumeConfigFileToPaths[i] = make(map[string]interface{})
				ConfigFileVolumeConfigFileToPaths[i]["Content"] = ConfigFileVolumeConfigFileToPathsMap["content"]
				ConfigFileVolumeConfigFileToPaths[i]["Path"] = ConfigFileVolumeConfigFileToPathsMap["path"]
			}
			Volumes[i]["ConfigFileVolume.ConfigFileToPath"] = ConfigFileVolumeConfigFileToPaths

			Volumes[i]["NFSVolume.Path"] = VolumesMap["nfs_volume_path"]
			Volumes[i]["NFSVolume.Server"] = VolumesMap["nfs_volume_server"]
			Volumes[i]["Name"] = VolumesMap["name"]
			Volumes[i]["Type"] = VolumesMap["type"]
		}
		request["Volume"] = Volumes

	}
	if d.HasChange("image_registry_credential") {
		update = true
		if v, ok := d.GetOk("image_registry_credential"); ok {
			imageRegisryCredentialMaps := make([]map[string]interface{}, 0)
			for _, raw := range v.(*schema.Set).List() {
				obj := raw.(map[string]interface{})
				imageRegisryCredentialMaps = append(imageRegisryCredentialMaps, map[string]interface{}{
					"Password": obj["password"],
					"Server":   obj["server"],
					"UserName": obj["user_name"],
				})
			}
			request["ImageRegistryCredential"] = imageRegisryCredentialMaps
		}
	}
	if update {
		action := "UpdateContainerGroup"
		conn, err := client.NewEciClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running", "Succeeded"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, eciService.EciContainerGroupStateRefreshFunc(d.Id(), []string{"Failed", "ScheduleFailed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudEciContainerGroupRead(d, meta)
}
func resourceAlicloudEciContainerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteContainerGroup"
	var response map[string]interface{}
	conn, err := client.NewEciClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ContainerGroupId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
