package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"strings"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudFCFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFCFunctionCreate,
		Read:   resourceAlicloudFCFunctionRead,
		Update: resourceAlicloudFCFunctionUpdate,
		Delete: resourceAlicloudFCFunctionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validation.StringLenBetween(1, 128),
			},
			"name_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 122),
			},

			"oss_bucket": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"filename"},
			},

			"oss_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"filename"},
			},

			"filename": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"oss_bucket", "oss_key"},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"code_checksum": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"environment_variables": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"handler": {
				Type:     schema.TypeString,
				Required: true,
			},
			"memory_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      128,
				ValidateFunc: validation.IntBetween(128, 32768),
			},
			"runtime": {
				Type:     schema.TypeString,
				Required: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"function_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"initializer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"initialization_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"instance_concurrency": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"ca_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "e1",
			},
			"custom_container_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": {
							Type:     schema.TypeString,
							Required: true,
						},
						"command": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"args": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudFCFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	serviceName := d.Get("service").(string)
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else if v, ok := d.GetOk("name_prefix"); ok {
		name = resource.PrefixedUniqueId(v.(string))
	} else {
		name = resource.UniqueId()
	}

	request := &fc.CreateFunctionInput{
		ServiceName: StringPointer(serviceName),
	}
	object := fc.FunctionCreateObject{
		FunctionName:          StringPointer(name),
		Description:           StringPointer(d.Get("description").(string)),
		Runtime:               StringPointer(d.Get("runtime").(string)),
		Handler:               StringPointer(d.Get("handler").(string)),
		Timeout:               Int32Pointer(int32(d.Get("timeout").(int))),
		MemorySize:            Int32Pointer(int32(d.Get("memory_size").(int))),
		Initializer:           StringPointer(d.Get("initializer").(string)),
		InitializationTimeout: Int32Pointer(int32(d.Get("initialization_timeout").(int))),
		InstanceType:          StringPointer(d.Get("instance_type").(string)),
	}
	// Set function environment variables.
	if variables := d.Get("environment_variables").(map[string]interface{}); len(variables) > 0 {
		byteVar, err := json.Marshal(variables)
		if err != nil {
			return WrapError(err)
		}
		err = json.Unmarshal(byteVar, &object.EnvironmentVariables)
		if err != nil {
			return WrapError(err)
		}
	}
	if strings.EqualFold(*object.Runtime, "custom-container") {
		// Set custom container config.
		cfg, err := parseCustomContainerConfig(d)
		if err != nil {
			return WrapError(err)
		}
		object.CustomContainerConfig = cfg
	} else {
		// Set function code.
		code, err := getFunctionCode(d)
		if err != nil {
			return WrapError(err)
		}
		object.Code = code
	}
	// Set CA port if the runtime is custom runtime or custom container.
	if strings.EqualFold(*object.Runtime, "custom") || strings.EqualFold(*object.Runtime, "custom-container") {
		object.CAPort = Int32Pointer(int32(d.Get("ca_port").(int)))
	}
	// Disable instance concurrency for python runtime.
	if !strings.EqualFold(*object.Runtime, "python2.7") && !strings.EqualFold(*object.Runtime, "python3") {
		object.InstanceConcurrency = Int32Pointer(int32(d.Get("instance_concurrency").(int)))
	}
	request.FunctionCreateObject = object

	var function *fc.CreateFunctionOutput
	var requestInfo *fc.Client
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.CreateFunction(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"AccessDenied"}) {
				return resource.RetryableError(WrapError(err))
			}
			return resource.NonRetryableError(WrapError(err))
		}
		addDebug("CreateFunction", raw, requestInfo, request)
		function, _ = raw.(*fc.CreateFunctionOutput)
		return nil

	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fc_function", "CreateFunction", FcGoSdk)
	}

	if function == nil {
		return WrapError(Error("Creating function compute function got a empty response"))
	}

	d.SetId(fmt.Sprintf("%s%s%s", serviceName, COLON_SEPARATED, *function.FunctionName))

	return resourceAlicloudFCFunctionRead(d, meta)
}

func resourceAlicloudFCFunctionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	object, err := fcService.DescribeFcFunction(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("service", parts[0])
	d.Set("code_checksum", object.CodeChecksum)
	d.Set("name", object.FunctionName)
	d.Set("function_id", object.FunctionID)
	d.Set("description", object.Description)
	d.Set("handler", object.Handler)
	d.Set("memory_size", object.MemorySize)
	d.Set("runtime", object.Runtime)
	d.Set("timeout", object.Timeout)
	d.Set("last_modified", object.LastModifiedTime)
	d.Set("environment_variables", object.EnvironmentVariables)
	d.Set("initializer", object.Initializer)
	d.Set("initialization_timeout", object.InitializationTimeout)
	d.Set("instance_concurrency", object.InstanceConcurrency)
	d.Set("instance_type", object.InstanceType)
	d.Set("ca_port", object.CAPort)
	var customContainerConfig []map[string]interface{}
	if object.CustomContainerConfig != nil {
		customContainerConfig = append(customContainerConfig, map[string]interface{}{
			"image":   object.CustomContainerConfig.Image,
			"command": object.CustomContainerConfig.Command,
			"args":    object.CustomContainerConfig.Args,
		})
	}
	if err := d.Set("custom_container_config", customContainerConfig); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudFCFunctionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := &fc.UpdateFunctionInput{}

	update := false
	if d.HasChange("filename") || d.HasChange("oss_bucket") || d.HasChange("oss_key") || d.HasChange("code_checksum") {
		update = true
	}
	if d.HasChange("description") {
		update = true
		request.Description = StringPointer(d.Get("description").(string))
	}
	if d.HasChange("handler") {
		update = true
		request.Handler = StringPointer(d.Get("handler").(string))
	}
	if d.HasChange("memory_size") {
		update = true
		request.MemorySize = Int32Pointer(int32(d.Get("memory_size").(int)))
	}
	if d.HasChange("timeout") {
		update = true
		request.Timeout = Int32Pointer(int32(d.Get("timeout").(int)))
	}
	runtime := StringPointer(d.Get("runtime").(string))
	if d.HasChange("runtime") {
		update = true
		request.Runtime = runtime
	}
	if d.HasChange("environment_variables") {
		update = true
		byteVar, err := json.Marshal(d.Get("environment_variables").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		err = json.Unmarshal(byteVar, &request.EnvironmentVariables)
		if err != nil {
			return WrapError(err)
		}
	}
	if d.HasChange("initializer") {
		update = true
		request.Initializer = StringPointer(d.Get("initializer").(string))
	}
	if d.HasChange("initialization_timeout") {
		update = true
		request.InitializationTimeout = Int32Pointer(int32(d.Get("initialization_timeout").(int)))
	}
	if d.HasChange("instance_concurrency") {
		update = true
		request.InstanceConcurrency = Int32Pointer(int32(d.Get("instance_concurrency").(int)))
	}
	if d.HasChange("instance_type") {
		update = true
		request.InstanceType = StringPointer(d.Get("instance_type").(string))
	}
	if d.HasChange("ca_port") {
		update = true
		request.CAPort = Int32Pointer(int32(d.Get("ca_port").(int)))
	}
	if d.HasChange("custom_container_config") {
		update = true
		config, err := parseCustomContainerConfig(d)
		if err != nil {
			return WrapError(err)
		}
		request.CustomContainerConfig = config
	}

	if update {
		split := strings.Split(d.Id(), COLON_SEPARATED)
		request.ServiceName = StringPointer(split[0])
		request.FunctionName = StringPointer(split[1])
		if !strings.EqualFold(*runtime, "custom-container") {
			code, err := getFunctionCode(d)
			if err != nil {
				return WrapError(err)
			}
			request.Code = code
		}
		var requestInfo *fc.Client
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.UpdateFunction(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateFunction", FcGoSdk)
		}
		addDebug("UpdateFunction", raw, requestInfo, request)
	}

	return resourceAlicloudFCFunctionRead(d, meta)
}

func resourceAlicloudFCFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := &fc.DeleteFunctionInput{
		ServiceName:  StringPointer(parts[0]),
		FunctionName: StringPointer(parts[1]),
	}
	var requestInfo *fc.Client
	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.DeleteFunction(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ServiceNotFound", "FunctionNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteFunction", FcGoSdk)
	}
	addDebug("DeleteFunction", raw, requestInfo, request)
	return WrapError(fcService.WaitForFcFunction(d.Id(), Deleted, DefaultTimeout))
}

func getFunctionCode(d *schema.ResourceData) (*fc.Code, error) {
	code := fc.NewCode()
	if filename, ok := d.GetOk("filename"); ok && filename.(string) != "" {
		file, err := loadFileContent(filename.(string))
		if err != nil {
			return code, WrapError(err)
		}
		code.WithZipFile(file)
	} else {
		bucket, bucketOk := d.GetOk("oss_bucket")
		key, keyOk := d.GetOk("oss_key")
		if !bucketOk || !keyOk {
			return code, nil
		}
		code.WithOSSBucketName(bucket.(string)).WithOSSObjectName(key.(string))
	}
	return code, nil
}

// The first return value is nil when the "custom_container_config" is not been set.
func parseCustomContainerConfig(d *schema.ResourceData) (config *fc.CustomContainerConfig, err error) {
	c := fc.NewCustomContainerConfig()
	if v, ok := d.GetOk("custom_container_config"); ok {
		config, ok := v.([]interface{})[0].(map[string]interface{})
		if ok {
			return c.WithImage(config["image"].(string)).WithCommand(config["command"].(string)).WithArgs(config["args"].(string)), nil
		} else {
			return nil, Error("Failed to parse custom_container_config")
		}
	}
	// The "custom_container_config" has not been set.
	return nil, nil
}
