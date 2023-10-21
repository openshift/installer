package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudLogProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudLogProjectsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Normal", "Disable"}, true),
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modify_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudLogProjectsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	logService := LogService{client}

	var objects []*sls.LogProject
	var logProjectNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		logProjectNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response []string
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.ListProject()
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response, _ = raw.([]string)
		return nil
	})
	addDebug("ListProject", response)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_projects", "ListProject", AliyunLogGoSdkERROR)
	}

	for _, projectName := range response {
		if logProjectNameRegex != nil {
			if !logProjectNameRegex.MatchString(projectName) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[projectName]; !ok {
				continue
			}
		}
		project, err := logService.DescribeLogProject(projectName)
		if err != nil {
			if NotFoundError(err) {
				log.Printf("The project '%s' is no exist! \n", projectName)
				continue
			}
			return WrapError(err)
		}
		if statusOk && status.(string) != "" && status.(string) != project.Status {
			continue
		}
		objects = append(objects, project)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":               object.Name,
			"description":      object.Description,
			"last_modify_time": object.LastModifyTime,
			"owner":            object.Owner,
			"project_name":     object.Name,
			"region":           object.Region,
			"status":           object.Status,
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("projects", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
