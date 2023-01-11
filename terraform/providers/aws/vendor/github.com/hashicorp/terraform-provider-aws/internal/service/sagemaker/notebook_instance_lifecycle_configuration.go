package sagemaker

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceNotebookInstanceLifeCycleConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotebookInstanceLifeCycleConfigurationCreate,
		Read:   resourceNotebookInstanceLifeCycleConfigurationRead,
		Update: resourceNotebookInstanceLifeCycleConfigurationUpdate,
		Delete: resourceNotebookInstanceLifeCycleConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validName,
			},

			"on_create": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 16384),
			},

			"on_start": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 16384),
			},
		},
	}
}

func resourceNotebookInstanceLifeCycleConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SageMakerConn

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		name = resource.UniqueId()
	}

	createOpts := &sagemaker.CreateNotebookInstanceLifecycleConfigInput{
		NotebookInstanceLifecycleConfigName: aws.String(name),
	}

	// on_create is technically a list of NotebookInstanceLifecycleHook elements, but the list has to be length 1
	// (same for on_start)
	if v, ok := d.GetOk("on_create"); ok {
		hook := &sagemaker.NotebookInstanceLifecycleHook{Content: aws.String(v.(string))}
		createOpts.OnCreate = []*sagemaker.NotebookInstanceLifecycleHook{hook}
	}

	if v, ok := d.GetOk("on_start"); ok {
		hook := &sagemaker.NotebookInstanceLifecycleHook{Content: aws.String(v.(string))}
		createOpts.OnStart = []*sagemaker.NotebookInstanceLifecycleHook{hook}
	}

	log.Printf("[DEBUG] SageMaker notebook instance lifecycle configuration create config: %#v", *createOpts)
	_, err := conn.CreateNotebookInstanceLifecycleConfig(createOpts)
	if err != nil {
		return fmt.Errorf("creating SageMaker notebook instance lifecycle configuration: %s", err)
	}
	d.SetId(name)

	return resourceNotebookInstanceLifeCycleConfigurationRead(d, meta)
}

func resourceNotebookInstanceLifeCycleConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SageMakerConn

	request := &sagemaker.DescribeNotebookInstanceLifecycleConfigInput{
		NotebookInstanceLifecycleConfigName: aws.String(d.Id()),
	}

	lifecycleConfig, err := conn.DescribeNotebookInstanceLifecycleConfig(request)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, "ValidationException") {
			log.Printf("[INFO] unable to find the SageMaker notebook instance lifecycle configuration (%s); therefore it is removed from the state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading SageMaker notebook instance lifecycle configuration %s: %s", d.Id(), err)
	}

	if err := d.Set("name", lifecycleConfig.NotebookInstanceLifecycleConfigName); err != nil {
		return fmt.Errorf("setting name for SageMaker notebook instance lifecycle configuration (%s): %s", d.Id(), err)
	}

	if len(lifecycleConfig.OnCreate) > 0 && lifecycleConfig.OnCreate[0] != nil {
		if err := d.Set("on_create", lifecycleConfig.OnCreate[0].Content); err != nil {
			return fmt.Errorf("setting on_create for SageMaker notebook instance lifecycle configuration (%s): %s", d.Id(), err)
		}
	}

	if len(lifecycleConfig.OnStart) > 0 && lifecycleConfig.OnStart[0] != nil {
		if err := d.Set("on_start", lifecycleConfig.OnStart[0].Content); err != nil {
			return fmt.Errorf("setting on_start for SageMaker notebook instance lifecycle configuration (%s): %s", d.Id(), err)
		}
	}

	if err := d.Set("arn", lifecycleConfig.NotebookInstanceLifecycleConfigArn); err != nil {
		return fmt.Errorf("setting arn for SageMaker notebook instance lifecycle configuration (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceNotebookInstanceLifeCycleConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SageMakerConn

	updateOpts := &sagemaker.UpdateNotebookInstanceLifecycleConfigInput{
		NotebookInstanceLifecycleConfigName: aws.String(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("on_create"); ok {
		onCreateHook := &sagemaker.NotebookInstanceLifecycleHook{Content: aws.String(v.(string))}
		updateOpts.OnCreate = []*sagemaker.NotebookInstanceLifecycleHook{onCreateHook}
	}

	if v, ok := d.GetOk("on_start"); ok {
		onStartHook := &sagemaker.NotebookInstanceLifecycleHook{Content: aws.String(v.(string))}
		updateOpts.OnStart = []*sagemaker.NotebookInstanceLifecycleHook{onStartHook}
	}

	_, err := conn.UpdateNotebookInstanceLifecycleConfig(updateOpts)
	if err != nil {
		return fmt.Errorf("updating SageMaker Notebook Instance Lifecycle Configuration: %s", err)
	}
	return resourceNotebookInstanceLifeCycleConfigurationRead(d, meta)
}

func resourceNotebookInstanceLifeCycleConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SageMakerConn

	deleteOpts := &sagemaker.DeleteNotebookInstanceLifecycleConfigInput{
		NotebookInstanceLifecycleConfigName: aws.String(d.Id()),
	}
	log.Printf("[INFO] Deleting SageMaker Notebook Instance Lifecycle Configuration: %s", d.Id())

	_, err := conn.DeleteNotebookInstanceLifecycleConfig(deleteOpts)
	if err != nil {

		if tfawserr.ErrCodeEquals(err, "ValidationException") {
			return nil
		}

		return fmt.Errorf("deleting SageMaker Notebook Instance Lifecycle Configuration: %s", err)
	}
	return nil
}
