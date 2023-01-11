package sagemaker

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceAppImageConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppImageConfigCreate,
		Read:   resourceAppImageConfigRead,
		Update: resourceAppImageConfigUpdate,
		Delete: resourceAppImageConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_image_config_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 63),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9](-*[a-zA-Z0-9])*$`), "Valid characters are a-z, A-Z, 0-9, and - (hyphen)."),
				),
			},
			"kernel_gateway_image_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_system_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_gid": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      100,
										ValidateFunc: validation.IntInSlice([]int{0, 100}),
									},
									"default_uid": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      1000,
										ValidateFunc: validation.IntInSlice([]int{0, 1000}),
									},
									"mount_path": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "/home/sagemaker-user",
										ValidateFunc: validation.All(
											validation.StringLenBetween(1, 1024),
											validation.StringMatch(regexp.MustCompile(`^\/.*`), "Must start with `/`."),
										),
									},
								},
							},
						},
						"kernel_spec": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"display_name": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringLenBetween(1, 1024),
									},
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringLenBetween(1, 1024),
									},
								},
							},
						},
					},
				},
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},
		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceAppImageConfigCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SageMakerConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	name := d.Get("app_image_config_name").(string)
	input := &sagemaker.CreateAppImageConfigInput{
		AppImageConfigName: aws.String(name),
	}

	if len(tags) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	if v, ok := d.GetOk("kernel_gateway_image_config"); ok && len(v.([]interface{})) > 0 {
		input.KernelGatewayImageConfig = expandAppImageConfigKernelGatewayImageConfig(v.([]interface{}))
	}

	_, err := conn.CreateAppImageConfig(input)
	if err != nil {
		return fmt.Errorf("creating SageMaker App Image Config %s: %w", name, err)
	}

	d.SetId(name)

	return resourceAppImageConfigRead(d, meta)
}

func resourceAppImageConfigRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SageMakerConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	image, err := FindAppImageConfigByName(conn, d.Id())
	if err != nil {
		if tfawserr.ErrMessageContains(err, sagemaker.ErrCodeResourceNotFound, "does not exist") {
			d.SetId("")
			log.Printf("[WARN] Unable to find SageMaker App Image Config (%s); removing from state", d.Id())
			return nil
		}
		return fmt.Errorf("reading SageMaker App Image Config (%s): %w", d.Id(), err)

	}

	arn := aws.StringValue(image.AppImageConfigArn)
	d.Set("app_image_config_name", image.AppImageConfigName)
	d.Set("arn", arn)

	if err := d.Set("kernel_gateway_image_config", flattenAppImageConfigKernelGatewayImageConfig(image.KernelGatewayImageConfig)); err != nil {
		return fmt.Errorf("setting kernel_gateway_image_config: %w", err)
	}

	tags, err := ListTags(conn, arn)

	if err != nil {
		return fmt.Errorf("listing tags for SageMaker App Image Config (%s): %w", d.Id(), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("setting tags_all: %w", err)
	}

	return nil
}

func resourceAppImageConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SageMakerConn

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("updating SageMaker App Image Config (%s) tags: %w", d.Id(), err)
		}
	}

	if d.HasChange("kernel_gateway_image_config") {

		input := &sagemaker.UpdateAppImageConfigInput{
			AppImageConfigName: aws.String(d.Id()),
		}

		if v, ok := d.GetOk("kernel_gateway_image_config"); ok && len(v.([]interface{})) > 0 {
			input.KernelGatewayImageConfig = expandAppImageConfigKernelGatewayImageConfig(v.([]interface{}))
		}

		log.Printf("[DEBUG] SageMaker App Image Config update config: %#v", *input)
		_, err := conn.UpdateAppImageConfig(input)
		if err != nil {
			return fmt.Errorf("updating SageMaker App Image Config: %w", err)
		}

	}

	return resourceAppImageConfigRead(d, meta)
}

func resourceAppImageConfigDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SageMakerConn

	input := &sagemaker.DeleteAppImageConfigInput{
		AppImageConfigName: aws.String(d.Id()),
	}

	if _, err := conn.DeleteAppImageConfig(input); err != nil {
		if tfawserr.ErrMessageContains(err, sagemaker.ErrCodeResourceNotFound, "does not exist") {
			return nil
		}
		return fmt.Errorf("deleting SageMaker App Image Config (%s): %w", d.Id(), err)
	}

	return nil
}

func expandAppImageConfigKernelGatewayImageConfig(l []interface{}) *sagemaker.KernelGatewayImageConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	config := &sagemaker.KernelGatewayImageConfig{}

	if v, ok := m["kernel_spec"].([]interface{}); ok && len(v) > 0 {
		config.KernelSpecs = expandAppImageConfigKernelGatewayImageConfigKernelSpecs(v)
	}

	if v, ok := m["file_system_config"].([]interface{}); ok && len(v) > 0 {
		config.FileSystemConfig = expandAppImageConfigKernelGatewayImageConfigFileSystemConfig(v)
	}

	return config
}

func expandAppImageConfigKernelGatewayImageConfigFileSystemConfig(l []interface{}) *sagemaker.FileSystemConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	config := &sagemaker.FileSystemConfig{
		DefaultGid: aws.Int64(int64(m["default_gid"].(int))),
		DefaultUid: aws.Int64(int64(m["default_uid"].(int))),
		MountPath:  aws.String(m["mount_path"].(string)),
	}

	return config
}

func expandAppImageConfigKernelGatewayImageConfigKernelSpecs(tfList []interface{}) []*sagemaker.KernelSpec {
	if len(tfList) == 0 {
		return nil
	}

	var kernelSpecs []*sagemaker.KernelSpec

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		kernelSpec := &sagemaker.KernelSpec{
			Name: aws.String(tfMap["name"].(string)),
		}

		if v, ok := tfMap["display_name"].(string); ok && v != "" {
			kernelSpec.DisplayName = aws.String(v)
		}

		kernelSpecs = append(kernelSpecs, kernelSpec)
	}

	return kernelSpecs
}

func flattenAppImageConfigKernelGatewayImageConfig(config *sagemaker.KernelGatewayImageConfig) []map[string]interface{} {
	if config == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{}

	if config.KernelSpecs != nil {
		m["kernel_spec"] = flattenAppImageConfigKernelGatewayImageConfigKernelSpecs(config.KernelSpecs)
	}

	if config.FileSystemConfig != nil {
		m["file_system_config"] = flattenAppImageConfigKernelGatewayImageConfigFileSystemConfig(config.FileSystemConfig)
	}

	return []map[string]interface{}{m}
}

func flattenAppImageConfigKernelGatewayImageConfigFileSystemConfig(config *sagemaker.FileSystemConfig) []map[string]interface{} {
	if config == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"mount_path":  aws.StringValue(config.MountPath),
		"default_gid": aws.Int64Value(config.DefaultGid),
		"default_uid": aws.Int64Value(config.DefaultUid),
	}

	return []map[string]interface{}{m}
}

func flattenAppImageConfigKernelGatewayImageConfigKernelSpecs(kernelSpecs []*sagemaker.KernelSpec) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(kernelSpecs))

	for _, raw := range kernelSpecs {
		kernelSpec := make(map[string]interface{})

		kernelSpec["name"] = aws.StringValue(raw.Name)

		if raw.DisplayName != nil {
			kernelSpec["display_name"] = aws.StringValue(raw.DisplayName)
		}

		res = append(res, kernelSpec)
	}

	return res
}
