package sagemaker

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceImageCreate,
		Read:   resourceImageRead,
		Update: resourceImageUpdate,
		Delete: resourceImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"image_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 63),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9](-*[a-zA-Z0-9])*$`), "Valid characters are a-z, A-Z, 0-9, and - (hyphen)."),
				),
			},
			"role_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceImageCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SageMakerConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	name := d.Get("image_name").(string)
	input := &sagemaker.CreateImageInput{
		ImageName: aws.String(name),
		RoleArn:   aws.String(d.Get("role_arn").(string)),
	}

	if v, ok := d.GetOk("display_name"); ok {
		input.DisplayName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if len(tags) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	// for some reason even if the operation is retried the same error response is given even though the role is valid. a short sleep before creation solves it.
	time.Sleep(1 * time.Minute)
	_, err := conn.CreateImage(input)
	if err != nil {
		return fmt.Errorf("creating SageMaker Image %s: %w", name, err)
	}

	d.SetId(name)

	if _, err := WaitImageCreated(conn, d.Id()); err != nil {
		return fmt.Errorf("waiting for SageMaker Image (%s) to be created: %w", d.Id(), err)
	}

	return resourceImageRead(d, meta)
}

func resourceImageRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SageMakerConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	image, err := FindImageByName(conn, d.Id())
	if err != nil {
		if tfawserr.ErrMessageContains(err, sagemaker.ErrCodeResourceNotFound, "does not exist") {
			d.SetId("")
			log.Printf("[WARN] Unable to find SageMaker Image (%s); removing from state", d.Id())
			return nil
		}
		return fmt.Errorf("reading SageMaker Image (%s): %w", d.Id(), err)

	}

	arn := aws.StringValue(image.ImageArn)
	d.Set("image_name", image.ImageName)
	d.Set("arn", arn)
	d.Set("role_arn", image.RoleArn)
	d.Set("display_name", image.DisplayName)
	d.Set("description", image.Description)

	tags, err := ListTags(conn, arn)

	if err != nil {
		return fmt.Errorf("listing tags for SageMaker Image (%s): %w", d.Id(), err)
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

func resourceImageUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SageMakerConn
	needsUpdate := false

	input := &sagemaker.UpdateImageInput{
		ImageName: aws.String(d.Id()),
	}

	var deleteProperties []*string

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			input.Description = aws.String(v.(string))
		} else {
			deleteProperties = append(deleteProperties, aws.String("Description"))
			input.DeleteProperties = deleteProperties
		}
		needsUpdate = true
	}

	if d.HasChange("display_name") {
		if v, ok := d.GetOk("display_name"); ok {
			input.DisplayName = aws.String(v.(string))
		} else {
			deleteProperties = append(deleteProperties, aws.String("DisplayName"))
			input.DeleteProperties = deleteProperties
		}
		needsUpdate = true
	}

	if needsUpdate {
		log.Printf("[DEBUG] sagemaker Image update config: %#v", *input)
		_, err := conn.UpdateImage(input)
		if err != nil {
			return fmt.Errorf("updating SageMaker Image: %w", err)
		}

		if _, err := WaitImageCreated(conn, d.Id()); err != nil {
			return fmt.Errorf("waiting for SageMaker Image (%s) to update: %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("updating SageMaker Image (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceImageRead(d, meta)
}

func resourceImageDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SageMakerConn

	input := &sagemaker.DeleteImageInput{
		ImageName: aws.String(d.Id()),
	}

	if _, err := conn.DeleteImage(input); err != nil {
		if tfawserr.ErrMessageContains(err, sagemaker.ErrCodeResourceNotFound, "No Image with the name") {
			return nil
		}
		return fmt.Errorf("deleting SageMaker Image (%s): %w", d.Id(), err)
	}

	if _, err := WaitImageDeleted(conn, d.Id()); err != nil {
		if tfawserr.ErrMessageContains(err, sagemaker.ErrCodeResourceNotFound, "does not exist") {
			return nil
		}
		return fmt.Errorf("waiting for SageMaker Image (%s) to delete: %w", d.Id(), err)
	}

	return nil
}
