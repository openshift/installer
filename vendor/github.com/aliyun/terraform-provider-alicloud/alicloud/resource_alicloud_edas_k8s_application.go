package alicloud

import (
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEdasK8sApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEdasK8sApplicationCreate,
		Read:   resourceAlicloudEdasK8sApplicationRead,
		Update: resourceAlicloudEdasK8sApplicationUpdate,
		Delete: resourceAlicloudEdasK8sApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"application_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"replicas": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"image_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"package_url"},
			},
			"package_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"FatJar", "War", "Image"}, false),
				Default:      "Image",
			},
			"application_descriotion": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"limit_mem": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"requests_mem": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"command": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"command_args": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"internet_slb_id": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"internet_slb_protocol": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false),
			},
			"internet_slb_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"internet_target_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"envs": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"pre_stop": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					e := EdasService{}
					return e.PreStopEqual(old, new)
				},
			},
			"post_start": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					e := EdasService{}
					return e.PostStartEqual(old, new)
				},
			},
			"liveness": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					e := EdasService{}
					return e.LivenessEqual(old, new)
				},
			},
			"readiness": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					e := EdasService{}
					return e.ReadinessEqual(old, new)
				},
			},
			"nas_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mount_descs": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"local_volume": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logical_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"package_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"image_url"},
			},
			"package_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  strconv.FormatInt(time.Now().Unix(), 10),
			},
			"jdk": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"web_container": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"edas_container_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"requests_m_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit_m_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudEdasK8sApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}
	request := edas.CreateInsertK8sApplicationRequest()

	packageType := d.Get("package_type").(string)

	request.AppName = d.Get("application_name").(string)
	request.RegionId = client.RegionId
	request.PackageType = packageType
	request.ClusterId = d.Get("cluster_id").(string)

	if strings.ToLower(packageType) == "image" {
		if v, ok := d.GetOk("image_url"); !ok {
			return WrapError(Error("image_url is needed for creating image k8s application"))
		} else {
			request.ImageUrl = v.(string)
		}
	} else {
		if v, ok := d.GetOk("package_url"); !ok {
			return WrapError(Error("package_url is needed for creating fatjar k8s application"))
		} else {
			request.PackageUrl = v.(string)
		}
		if v, ok := d.GetOk("package_version"); ok {
			request.PackageVersion = v.(string)
		}
		if v, ok := d.GetOk("jdk"); !ok {
			return WrapError(Error("jdk is needed for creating non-image k8s application"))
		} else {
			request.JDK = v.(string)
		}
		if strings.ToLower(packageType) == "war" {
			var webContainer string
			var edasContainer string
			if v, ok := d.GetOk("web_container"); ok {
				webContainer = v.(string)
			}
			if v, ok := d.GetOk("edas_container_version"); ok {
				edasContainer = v.(string)
			}
			if len(webContainer) == 0 && len(edasContainer) == 0 {
				return WrapError(Error("web_container or edas_container_version is needed for creating war k8s application"))
			}
			request.WebContainer = webContainer
			request.EdasContainerVersion = edasContainer
		}
	}

	request.Replicas = requests.NewInteger(d.Get("replicas").(int))

	if v, ok := d.GetOk("application_descriotion"); ok {
		request.ApplicationDescription = v.(string)
	}

	if v, ok := d.GetOk("limit_mem"); ok {
		request.LimitMem = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("requests_mem"); ok {
		request.RequestsMem = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("command"); ok {
		request.Command = v.(string)
	}

	if v, ok := d.GetOk("command_args"); ok {
		commands, err := edasService.GetK8sCommandArgs(v.([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.CommandArgs = commands
	}

	if v, ok := d.GetOk("internet_slb_id"); ok {
		request.InternetSlbId = v.(string)
	}

	if v, ok := d.GetOk("internet_slb_protocol"); ok {
		request.InternetSlbProtocol = v.(string)
	}

	if v, ok := d.GetOk("internet_slb_port"); ok {
		request.InternetSlbPort = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("internet_target_port"); ok {
		request.InternetTargetPort = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("envs"); ok {
		envs, err := edasService.GetK8sEnvs(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.Envs = envs
	}

	if v, ok := d.GetOk("pre_stop"); ok {
		request.PreStop = v.(string)
	}

	if v, ok := d.GetOk("post_start"); ok {
		request.PostStart = v.(string)
	}

	if v, ok := d.GetOk("liveness"); ok {
		request.Liveness = v.(string)
	}

	if v, ok := d.GetOk("readiness"); ok {
		request.Readiness = v.(string)
	}

	if v, ok := d.GetOk("nas_id"); ok {
		request.NasId = v.(string)
	}

	if v, ok := d.GetOk("mount_descs"); ok {
		request.MountDescs = v.(string)
	}

	if v, ok := d.GetOk("local_volume"); ok {
		request.LocalVolume = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = v.(string)
	}

	if v, ok := d.GetOk("logical_region_id"); ok {
		request.LogicalRegionId = v.(string)
	}

	if v, ok := d.GetOk("requests_m_cpu"); ok {
		request.RequestsmCpu = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("limit_m_cpu"); ok {
		request.LimitmCpu = requests.NewInteger(v.(int))
	}

	var appId string
	var changeOrderId string

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.InsertK8sApplication(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_edas_k8s_application", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.InsertK8sApplicationResponse)
	appId = response.ApplicationInfo.AppId
	changeOrderId = response.ApplicationInfo.ChangeOrderId
	d.SetId(appId)
	if response.Code != 200 {
		return WrapError(Error("Create k8s application failed for " + response.Message))
	}

	if len(changeOrderId) > 0 {
		stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudEdasK8sApplicationRead(d, meta)
}

func resourceAlicloudEdasK8sApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	response, err := edasService.DescribeEdasK8sApplication(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_edas_k8s_application ecsService.DescribeEdasK8sApplication Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("application_name", response.App.ApplicationName)
	d.Set("application_descriotion", response.App.Description)
	d.Set("cluster_id", response.App.ClusterId)
	d.Set("replicas", response.App.Instances)
	d.Set("package_type", response.App.ApplicationType)
	d.Set("image_url", response.ImageInfo.ImageUrl)
	envs := make(map[string]string)
	for _, e := range response.App.EnvList.Env {
		envs[e.Name] = e.Value
	}
	d.Set("envs", envs)
	d.Set("command", response.App.Cmd)
	d.Set("command_args", response.App.CmdArgs.CmdArg)

	allDeploy := response.DeployGroups.DeployGroup
	for _, v := range allDeploy {
		if len(v.PackageUrl) > 0 {
			d.Set("package_url", v.PackageUrl)
		}
		if len(v.PackageVersion) > 0 {
			d.Set("package_version", v.PackageVersion)
		}
		if len(v.MemoryLimit) > 0 {
			d.Set("limit_mem", v.MemoryLimit)
		}
		if len(v.MemoryRequest) > 0 {
			d.Set("requests_mem", v.MemoryRequest)
		}

		if len(v.CpuRequest) > 0 {
			cpu, err := strconv.Atoi(v.CpuRequest)

			if err != nil {
				return WrapError(err)
			}
			d.Set("requests_m_cpu", cpu*1000)
		}

		for _, c := range v.Components.ComponentsItem {
			if strings.Contains(c.ComponentKey, "JDK") {
				d.Set("jdk", c.ComponentKey)
			}
		}
	}

	if len(response.App.EdasContainerVersion) > 0 {
		d.Set("edas_container_version", response.App.EdasContainerVersion)
	}
	if len(response.Conf.PreStop) > 0 {
		d.Set("pre_stop", response.Conf.PreStop)
	}
	if len(response.Conf.PostStart) > 0 {
		d.Set("post_start", response.Conf.PostStart)
	}
	if len(response.Conf.Liveness) > 0 {
		d.Set("liveness", response.Conf.Liveness)
	}
	if len(response.Conf.Readiness) > 0 {
		d.Set("readiness", response.Conf.Readiness)
	}
	d.Set("namespace", response.NameSpace)
	return nil
}

func resourceAlicloudEdasK8sApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}
	var partialKeys []string
	request := edas.CreateDeployK8sApplicationRequest()
	request.AppId = d.Id()
	request.RegionId = client.RegionId

	d.Partial(true)
	packageType, err := edasService.QueryK8sAppPackageType(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if strings.ToLower(packageType) == "image" {
		if d.HasChange("image_url") {
			partialKeys = append(partialKeys, "image_url")
		}
		request.Image = d.Get("image_url").(string)
		if len(request.Image) == 0 {
			return WrapError(Error("image_url is needed for creating image k8s application"))
		}
	} else {
		if d.HasChange("package_url") {
			partialKeys = append(partialKeys, "package_url")
		}
		request.PackageUrl = d.Get("package_url").(string)
		if len(request.PackageUrl) == 0 {
			return WrapError(Error("package_url is needed for creating fatjar k8s application"))
		}
		request.PackageVersion = d.Get("package_version").(string)

		if d.HasChange("jdk") {
			partialKeys = append(partialKeys, "jdk")
		}
		request.JDK = d.Get("jdk").(string)
		if len(request.JDK) == 0 {
			return WrapError(Error("jdk is needed for creating non-image k8s application"))
		}
		if strings.ToLower(packageType) == "war" {
			var webContainer string
			var edasContainer string
			if d.HasChange("web_container") {
				partialKeys = append(partialKeys, "web_container")
			}
			webContainer = d.Get("web_container").(string)

			if d.HasChange("edas_container_version") {
				partialKeys = append(partialKeys, "edas_container_version")
			}
			edasContainer = d.Get("edas_container_version").(string)
			if len(webContainer) == 0 && len(edasContainer) == 0 {
				return WrapError(Error("web_container or edas_container_version is needed for updating war k8s application"))
			}
			request.WebContainer = webContainer
			request.EdasContainerVersion = edasContainer
		}
	}

	if d.HasChange("replicas") {
		partialKeys = append(partialKeys, "replicas")
	}
	request.Replicas = requests.NewInteger(d.Get("replicas").(int))

	if d.HasChange("limit_mem") {
		partialKeys = append(partialKeys, "limit_mem")
		request.MemoryLimit = requests.NewInteger(d.Get("limit_mem").(int))
	}

	if d.HasChange("requests_mem") {
		partialKeys = append(partialKeys, "requests_mem")
		request.MemoryRequest = requests.NewInteger(d.Get("requests_mem").(int))
	}

	if d.HasChange("command") {
		partialKeys = append(partialKeys, "command")
		request.Command = d.Get("command").(string)
	}

	if d.HasChange("command_args") {
		partialKeys = append(partialKeys, "command_args")
		commands, err := edasService.GetK8sCommandArgsForDeploy(d.Get("command_args").([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.Args = commands
	}

	if d.HasChange("envs") {
		partialKeys = append(partialKeys, "envs")
		envs, err := edasService.GetK8sEnvs(d.Get("envs").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.Envs = envs
	}

	if d.HasChange("pre_stop") {
		if !edasService.PreStopEqual(d.GetChange("pre_stop")) {
			partialKeys = append(partialKeys, "pre_stop")
			request.PreStop = d.Get("pre_stop").(string)
		}
	}

	if d.HasChange("post_start") {
		if !edasService.PostStartEqual(d.GetChange("post_start")) {
			partialKeys = append(partialKeys, "post_start")
			request.PostStart = d.Get("post_start").(string)
		}
	}

	if d.HasChange("liveness") {
		if !edasService.LivenessEqual(d.GetChange("liveness")) {
			partialKeys = append(partialKeys, "liveness")
			request.Liveness = d.Get("liveness").(string)
		}
	}

	if d.HasChange("readiness") {
		if !edasService.ReadinessEqual(d.GetChange("readiness")) {
			partialKeys = append(partialKeys, "readiness")
			request.Readiness = d.Get("readiness").(string)
		}
	}

	if d.HasChange("nas_id") {
		partialKeys = append(partialKeys, "nas_id")
		request.NasId = d.Get("nas_id").(string)
	}

	if d.HasChange("mount_descs") {
		partialKeys = append(partialKeys, "mount_descs")
		request.MountDescs = d.Get("mount_descs").(string)
	}

	if d.HasChange("local_volume") {
		partialKeys = append(partialKeys, "local_volume")
		request.LocalVolume = d.Get("local_volume").(string)
	}

	if d.HasChange("requests_m_cpu") {
		partialKeys = append(partialKeys, "requests_m_cpu")
		request.McpuRequest = requests.NewInteger(d.Get("requests_m_cpu").(int))
	}

	if d.HasChange("limit_m_cpu") {
		partialKeys = append(partialKeys, "limit_m_cpu")
		request.McpuLimit = requests.NewInteger(d.Get("limit_m_cpu").(int))
	}

	if len(partialKeys) > 0 {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeployK8sApplication(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)

		response, _ := raw.(*edas.DeployK8sApplicationResponse)
		changeOrderId := response.ChangeOrderId
		if response.Code != 200 {
			return WrapError(Error("deploy k8s application failed for " + response.Message))
		}

		if len(changeOrderId) > 0 {
			stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		for _, key := range partialKeys {
			d.SetPartial(key)
		}
	}
	d.Partial(false)
	return resourceAlicloudEdasK8sApplicationRead(d, meta)
}

func resourceAlicloudEdasK8sApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	edasService := EdasService{client}

	request := edas.CreateDeleteK8sApplicationRequest()
	request.RegionId = client.RegionId
	request.AppId = d.Id()

	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeleteK8sApplication(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		response := raw.(*edas.DeleteK8sApplicationResponse)
		if response.Code != 200 {
			return resource.NonRetryableError(Error("Delete k8s application failed for " + response.Message))
		}
		changeOrderId := response.ChangeOrderId

		if len(changeOrderId) > 0 {
			stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return resource.NonRetryableError(WrapErrorf(err, IdMsg, d.Id()))
			}
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
