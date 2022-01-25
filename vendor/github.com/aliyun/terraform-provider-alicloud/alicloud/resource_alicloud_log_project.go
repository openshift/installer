package alicloud

import (
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudLogProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogProjectCreate,
		Read:   resourceAlicloudLogProjectRead,
		Update: resourceAlicloudLogProjectUpdate,
		Delete: resourceAlicloudLogProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudLogProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	var requestInfo *sls.Client
	request := map[string]string{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
	}
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.CreateProject(request["name"], request["description"])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateProject", raw, requestInfo, request)
		response, _ := raw.(*sls.LogProject)
		d.SetId(response.Name)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_project", "CreateProject", AliyunLogGoSdkERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, logService.LogProjectStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudLogProjectUpdate(d, meta)
}

func resourceAlicloudLogProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	object, err := logService.DescribeLogProject(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.Name)
	d.Set("description", object.Description)
	projectTags, err := logService.DescribeLogProjectTags(object.Name)
	if projectTags != nil {
		tags := map[string]interface{}{}
		for _, tag := range projectTags {
			tags[tag.TagKey] = tag.TagValue
		}
		if err := d.Set("tags", tags); err != nil {
			return WrapError(err)
		}
	}

	return nil
}

func buildTags(projectName string, tags map[string]interface{}) *sls.ResourceTags {
	slsTags := []sls.ResourceTag{}

	for key, value := range tags {
		tag := sls.ResourceTag{Key: key, Value: value.(string)}
		slsTags = append(slsTags, tag)
	}
	projectTags := sls.NewProjectTags(projectName, slsTags)
	return projectTags
}

func deleteProjectTags(client *connectivity.AliyunClient, slsTags []string, projectName string) error {
	var requestInfo *sls.Client
	projectUnTags := sls.NewProjectUnTags(projectName, slsTags)
	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		requestInfo = slsClient
		return nil, slsClient.UnTagResources(projectName, projectUnTags)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, projectName, "DeletaTags", AliyunLogGoSdkERROR)
	}
	addDebug("DeletaTags", raw, requestInfo, map[string]string{
		"name": projectName,
	})
	return nil
}

func resourceAlicloudLogProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	var requestInfo *sls.Client
	logService := LogService{client}
	projectName := d.Get("name").(string)
	request := map[string]string{
		"name":        projectName,
		"description": d.Get("description").(string),
	}
	if d.HasChange("description") {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return slsClient.UpdateProject(request["name"], request["description"])
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateProject", AliyunLogGoSdkERROR)
		}
		addDebug("UpdateProject", raw, requestInfo, request)
	}

	if d.HasChange("tags") {
		projectTags, err := logService.DescribeLogProjectTags(projectName)
		if err != nil {
			return err
		}
		slsTags := []string{}
		for _, value := range projectTags {
			slsTags = append(slsTags, value.TagKey)
		}
		tags := d.Get("tags").(map[string]interface{})
		if tags == nil || len(tags) == 0 {
			if err := deleteProjectTags(client, slsTags, projectName); err != nil {
				return WrapError(err)
			}
		} else {
			if err := deleteProjectTags(client, slsTags, projectName); err != nil {
				return WrapError(err)
			}
			projectNewTags := buildTags(projectName, tags)
			raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				requestInfo = slsClient
				return nil, slsClient.TagResources(projectName, projectNewTags)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateTags", AliyunLogGoSdkERROR)
			}
			addDebug("UpdateTags", raw, requestInfo, request)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)

	return resourceAlicloudLogProjectRead(d, meta)
}

func resourceAlicloudLogProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}
	var requestInfo *sls.Client
	request := map[string]string{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
	}
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteProject(request["name"])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout, "RequestTimeout"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteProject", raw, requestInfo, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ProjectNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteProject", AliyunLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogProject(d.Id(), Deleted, DefaultTimeout))
}
