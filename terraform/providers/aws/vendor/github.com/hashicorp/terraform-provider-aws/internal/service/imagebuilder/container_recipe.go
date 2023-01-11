package imagebuilder

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceContainerRecipe() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerRecipeCreate,
		Read:   resourceContainerRecipeRead,
		Update: resourceContainerRecipeUpdate,
		Delete: resourceContainerRecipeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"component": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"component_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: verify.ValidARN,
						},
						"parameter": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringLenBetween(1, 256),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"container_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DOCKER"}, false),
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1024),
			},
			"dockerfile_template_data": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"dockerfile_template_data", "dockerfile_template_uri"},
				ValidateFunc: validation.StringLenBetween(1, 16000),
			},
			"dockerfile_template_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"dockerfile_template_data", "dockerfile_template_uri"},
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^s3://`), "must begin with s3://"),
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"instance_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"block_device_mapping": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_name": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringLenBetween(1, 1024),
									},
									"ebs": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"delete_on_termination": {
													// Use TypeString to allow an "unspecified" value,
													// since TypeBool only has true/false with false default.
													// The conversion from bare true/false values in
													// configurations to TypeString value is currently safe.
													Type:             schema.TypeString,
													Optional:         true,
													ForceNew:         true,
													DiffSuppressFunc: verify.SuppressEquivalentTypeStringBoolean,
													ValidateFunc:     verify.ValidTypeStringNullableBoolean,
												},
												"encrypted": {
													// Use TypeString to allow an "unspecified" value,
													// since TypeBool only has true/false with false default.
													// The conversion from bare true/false values in
													// configurations to TypeString value is currently safe.
													Type:             schema.TypeString,
													Optional:         true,
													ForceNew:         true,
													DiffSuppressFunc: verify.SuppressEquivalentTypeStringBoolean,
													ValidateFunc:     verify.ValidTypeStringNullableBoolean,
												},
												"iops": {
													Type:         schema.TypeInt,
													Optional:     true,
													ForceNew:     true,
													ValidateFunc: validation.IntBetween(100, 64000),
												},
												"kms_key_id": {
													Type:         schema.TypeString,
													Optional:     true,
													ForceNew:     true,
													ValidateFunc: validation.StringLenBetween(1, 1024),
												},
												"snapshot_id": {
													Type:         schema.TypeString,
													Optional:     true,
													ForceNew:     true,
													ValidateFunc: validation.StringLenBetween(1, 1024),
												},
												"throughput": {
													Type:         schema.TypeInt,
													Optional:     true,
													ForceNew:     true,
													ValidateFunc: validation.IntBetween(125, 1000),
												},
												"volume_size": {
													Type:         schema.TypeInt,
													Optional:     true,
													ForceNew:     true,
													ValidateFunc: validation.IntBetween(1, 16000),
												},
												"volume_type": {
													Type:         schema.TypeString,
													Optional:     true,
													ForceNew:     true,
													ValidateFunc: validation.StringInSlice(imagebuilder.EbsVolumeType_Values(), false),
												},
											},
										},
									},
									"no_device": {
										// Use TypeBool to allow an "unspecified" value of false,
										// since the API uses an empty string ("") as true and
										// this is not compatible with TypeString's zero value.
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"virtual_name": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringLenBetween(1, 1024),
									},
								},
							},
						},
						"image": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(1, 1024),
						},
					},
				},
			},
			"kms_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1024),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_image": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1024),
			},
			"platform": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"target_repository": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repository_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1024),
						},
						"service": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ECR"}, false),
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"working_directory": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 1024),
			},
		},
		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceContainerRecipeCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ImageBuilderConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	input := &imagebuilder.CreateContainerRecipeInput{
		ClientToken: aws.String(resource.UniqueId()),
	}

	if v, ok := d.GetOk("component"); ok && len(v.([]interface{})) > 0 {
		input.Components = expandComponentConfigurations(v.([]interface{}))
	}

	if v, ok := d.GetOk("container_type"); ok {
		input.ContainerType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("dockerfile_template_data"); ok {
		input.DockerfileTemplateData = aws.String(v.(string))
	}

	if v, ok := d.GetOk("dockerfile_template_uri"); ok {
		input.DockerfileTemplateUri = aws.String(v.(string))
	}

	if v, ok := d.GetOk("instance_configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.InstanceConfiguration = expandInstanceConfiguration(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		input.KmsKeyId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		input.Name = aws.String(v.(string))
	}

	if v, ok := d.GetOk("parent_image"); ok {
		input.ParentImage = aws.String(v.(string))
	}

	if len(tags) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	if v, ok := d.GetOk("target_repository"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.TargetRepository = expandTargetContainerRepository(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("version"); ok {
		input.SemanticVersion = aws.String(v.(string))
	}

	if v, ok := d.GetOk("working_directory"); ok {
		input.WorkingDirectory = aws.String(v.(string))
	}

	output, err := conn.CreateContainerRecipe(input)

	if err != nil {
		return fmt.Errorf("error creating Image Builder Container Recipe: %w", err)
	}

	if output == nil {
		return fmt.Errorf("error creating Image Builder Container Recipe: empty response")
	}

	d.SetId(aws.StringValue(output.ContainerRecipeArn))

	return resourceContainerRecipeRead(d, meta)
}

func resourceContainerRecipeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ImageBuilderConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &imagebuilder.GetContainerRecipeInput{
		ContainerRecipeArn: aws.String(d.Id()),
	}

	output, err := conn.GetContainerRecipe(input)

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, imagebuilder.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Image Builder Container Recipe (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting Image Builder Container Recipe (%s): %w", d.Id(), err)
	}

	if output == nil || output.ContainerRecipe == nil {
		return fmt.Errorf("error getting Image Builder Container Recipe (%s): empty response", d.Id())
	}

	containerRecipe := output.ContainerRecipe

	d.Set("arn", containerRecipe.Arn)
	d.Set("component", flattenComponentConfigurations(containerRecipe.Components))
	d.Set("container_type", containerRecipe.ContainerType)
	d.Set("date_created", containerRecipe.DateCreated)
	d.Set("description", containerRecipe.Description)
	d.Set("dockerfile_template_data", containerRecipe.DockerfileTemplateData)
	d.Set("encrypted", containerRecipe.Encrypted)

	if containerRecipe.InstanceConfiguration != nil {
		d.Set("instance_configuration", []interface{}{flattenInstanceConfiguration(containerRecipe.InstanceConfiguration)})
	} else {
		d.Set("instance_configuration", nil)
	}

	d.Set("kms_key_id", containerRecipe.KmsKeyId)
	d.Set("name", containerRecipe.Name)
	d.Set("owner", containerRecipe.Owner)
	d.Set("parent_image", containerRecipe.ParentImage)
	d.Set("platform", containerRecipe.Platform)

	tags := KeyValueTags(containerRecipe.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	d.Set("target_repository", []interface{}{flattenTargetContainerRepository(containerRecipe.TargetRepository)})
	d.Set("version", containerRecipe.Version)
	d.Set("working_directory", containerRecipe.WorkingDirectory)

	return nil
}

func resourceContainerRecipeUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ImageBuilderConn

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating tags for Image Builder Container Recipe (%s): %w", d.Id(), err)
		}
	}

	return resourceContainerRecipeRead(d, meta)
}

func resourceContainerRecipeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ImageBuilderConn

	input := &imagebuilder.DeleteContainerRecipeInput{
		ContainerRecipeArn: aws.String(d.Id()),
	}

	_, err := conn.DeleteContainerRecipe(input)

	if tfawserr.ErrCodeEquals(err, imagebuilder.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Image Builder Container Recipe (%s): %w", d.Id(), err)
	}

	return nil
}

func expandInstanceConfiguration(tfMap map[string]interface{}) *imagebuilder.InstanceConfiguration {
	if tfMap == nil {
		return nil
	}

	apiObject := &imagebuilder.InstanceConfiguration{}

	if v, ok := tfMap["block_device_mapping"].(*schema.Set); ok && v.Len() > 0 {
		apiObject.BlockDeviceMappings = expandInstanceBlockDeviceMappings(v.List())
	}

	if v, ok := tfMap["image"].(string); ok && v != "" {
		apiObject.Image = aws.String(v)
	}

	return apiObject
}

func flattenInstanceConfiguration(apiObject *imagebuilder.InstanceConfiguration) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.BlockDeviceMappings; v != nil {
		tfMap["block_device_mapping"] = flattenInstanceBlockDeviceMappings(v)
	}

	if v := apiObject.Image; v != nil {
		tfMap["image"] = aws.StringValue(v)
	}

	return tfMap
}
