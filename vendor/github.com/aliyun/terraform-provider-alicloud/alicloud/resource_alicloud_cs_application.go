package alicloud

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"sort"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCSApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSApplicationCreate,
		Read:   resourceAlicloudCSApplicationRead,
		Update: resourceAlicloudCSApplicationUpdate,
		Delete: resourceAlicloudCSApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,63}$`), "should be 1-64 characters long, and can contain numbers, English letters and hyphens, but cannot start with hyphens."),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.HasChange("version")
				},
			},
			"template": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					yaml, _ := normalizeYamlString(v)
					return yaml
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if !d.HasChange("version") {
						return true
					}
					equal, _ := compareYamlTemplateAreEquivalent(old, new)
					return equal
				},
				ValidateFunc: validateYamlString,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1.0",
			},
			"environment": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.HasChange("version")
				},
			},
			"latest_image": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.HasChange("version")
				},
			},
			"blue_green": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.IsNewResource()
				},
			},
			"blue_green_confirm": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.Get("blue_green").(bool)
				},
			},
			"services": {
				Type: schema.TypeList,
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"default_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCSApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	clusterName := Trim(d.Get("cluster_name").(string))

	args := &cs.ProjectCreationArgs{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Template:    d.Get("template").(string),
		Version:     d.Get("version").(string),
		LatestImage: d.Get("latest_image").(bool),
	}
	if environment, ok := d.GetOk("environment"); ok {
		env := make(map[string]string)
		for k, v := range environment.(map[string]interface{}) {
			env[k] = v.(string)
		}
		args.Environment = env
	}
	invoker := NewInvoker()
	if err := invoker.Run(func() error {
		cluster, certs, err := csService.GetContainerClusterAndCertsByName(clusterName)
		if err == nil {
			_, err = client.WithCsProjectClient(cluster.ClusterID, cluster.MasterURL, *certs, func(csProjectClient *cs.ProjectClient) (interface{}, error) {
				return nil, csProjectClient.CreateProject(args)
			})
		}
		return err
	}); err != nil {
		return fmt.Errorf("Creating container application got an error: %#v", err)
	}

	d.SetId(fmt.Sprintf("%s%s%s", clusterName, COLON_SEPARATED, args.Name))

	if err := csService.WaitForContainerApplication(clusterName, args.Name, Running, DefaultTimeoutMedium); err != nil {
		return fmt.Errorf("Waitting for container application %#v got an error: %#v", cs.Running, err)
	}

	return resourceAlicloudCSApplicationRead(d, meta)
}

func resourceAlicloudCSApplicationRead(d *schema.ResourceData, meta interface{}) error {
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	application, err := csService.DescribeContainerApplication(parts[0], parts[1])

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("cluster_name", parts[0])
	d.Set("name", application.Name)
	d.Set("description", application.Description)
	d.Set("version", application.Version)
	d.Set("template", application.Template)
	env := make(map[string]string)
	for k, v := range application.Environment {
		if k == "COMPOSE_PROJECT_NAME" || k == "ACS_PROJECT_VERSION" {
			continue
		}
		if k == "ACS_DEFAULT_DOMAIN" {
			d.Set("default_domain", v)
			continue
		}
		env[k] = v
	}
	d.Set("environment", env)
	var services []map[string]interface{}
	serviceIds := make([]string, len(application.Services))
	for _, s := range application.Services {
		serviceIds = append(serviceIds, s.ID)
	}
	if len(serviceIds) > 0 {
		sort.Strings(serviceIds)
	}
	for _, id := range serviceIds {
		for _, s := range application.Services {
			if s.ID != id {
				continue
			}
			mapping := map[string]interface{}{
				"id":      s.ID,
				"name":    s.Name,
				"status":  s.CurrentState,
				"version": s.Version,
			}
			services = append(services, mapping)
		}
	}

	if err := d.Set("services", services); err != nil {
		return err
	}

	return nil
}

func resourceAlicloudCSApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	clusterName := parts[0]
	invoker := NewInvoker()
	args := &cs.ProjectUpdationArgs{
		Name:        parts[1],
		Description: d.Get("description").(string),
		Template:    d.Get("template").(string),
		Version:     d.Get("version").(string),
		LatestImage: d.Get("latest_image").(bool),
	}

	if environment, ok := d.GetOk("environment"); ok {
		env := make(map[string]string)
		for k, v := range environment.(map[string]interface{}) {
			env[k] = v.(string)
		}
		args.Environment = env
	}

	blue_green := d.Get("blue_green").(bool)
	if blue_green {
		args.UpdateMethod = "blue-green"
	}

	d.Partial(true)
	update := false
	if d.HasChange("description") {
		update = true
		d.SetPartial("description")
	}

	if d.HasChange("template") {
		update = true
		d.SetPartial("template")
	}

	if d.HasChange("environment") {
		update = true
		d.SetPartial("environment")
	}

	if d.HasChange("version") {
		update = true
		d.SetPartial("version")
	}

	if d.HasChange("latest_image") {
		update = true
	}

	if d.HasChange("blue_green") {
		update = true
		d.SetPartial("blue_green")
	}

	if !d.HasChange("version") && !blue_green {
		cluster, certs, err := csService.GetContainerClusterAndCertsByName(clusterName)
		if err == nil {
			_, err = client.WithCsProjectClient(cluster.ClusterID, cluster.MasterURL, *certs, func(csProjectClient *cs.ProjectClient) (interface{}, error) {
				return nil, csProjectClient.RollBackBlueGreenProject(parts[1], true)
			})
		}
		return fmt.Errorf("Rollbacking container application blue-green got an error: %#v", err)
	}
	if update {
		for {
			err := invoker.Run(func() error {
				cluster, certs, err := csService.GetContainerClusterAndCertsByName(clusterName)
				if err == nil {
					_, err = client.WithCsProjectClient(cluster.ClusterID, cluster.MasterURL, *certs, func(csProjectClient *cs.ProjectClient) (interface{}, error) {
						return nil, csProjectClient.UpdateProject(args)
					})
				}
				return err
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"Conflicts with unconfirmed updates for operation"}) {
					err := invoker.Run(func() error {
						cluster, certs, err := csService.GetContainerClusterAndCertsByName(clusterName)
						if err == nil {
							_, err = client.WithCsProjectClient(cluster.ClusterID, cluster.MasterURL, *certs, func(csProjectClient *cs.ProjectClient) (interface{}, error) {
								return nil, csProjectClient.RollBackBlueGreenProject(parts[1], true)
							})
						}
						return err
					})
					if err != nil {
						return fmt.Errorf("Rollbacking container application blue-green got an error: %#v", err)
					}
					err = csService.WaitForContainerApplication(parts[0], parts[1], Running, DefaultTimeoutMedium)
					if err != nil {
						return fmt.Errorf("After rolling back blue-green project, waitting for container application %#v got an error: %#v", Running, err)
					}
					continue
				}
				return fmt.Errorf("Updating container application got an error: %#v", err)
			}
			break
		}
	}

	if d.Get("blue_green_confirm").(bool) {
		err := invoker.Run(func() error {
			cluster, certs, err := csService.GetContainerClusterAndCertsByName(clusterName)
			if err == nil {
				_, err = client.WithCsProjectClient(cluster.ClusterID, cluster.MasterURL, *certs, func(csProjectClient *cs.ProjectClient) (interface{}, error) {
					return nil, csProjectClient.ConfirmBlueGreenProject(parts[1], true)
				})
			}
			return err
		})
		if err != nil {
			return fmt.Errorf("Confirmming container application blue-green got an error: %#v", err)
		}
	}

	if err := csService.WaitForContainerApplication(parts[0], parts[1], Running, DefaultTimeoutMedium); err != nil {
		return fmt.Errorf("After updating, waitting for container application %#v got an error: %#v", Running, err)
	}

	d.Partial(false)

	return resourceAlicloudCSApplicationRead(d, meta)
}

func resourceAlicloudCSApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	clusterName := parts[0]

	appName := parts[1]
	invoker := NewInvoker()

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		err := invoker.Run(func() error {
			cluster, certs, err := csService.GetContainerClusterAndCertsByName(clusterName)
			if err == nil {
				_, err = client.WithCsProjectClient(cluster.ClusterID, cluster.MasterURL, *certs, func(csProjectClient *cs.ProjectClient) (interface{}, error) {
					return nil, csProjectClient.DeleteProject(appName, true, false)
				})
			}
			return err
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Not Found"}) {
				return nil
			}
			if !IsExpectedErrors(err, []string{"Unable to reach primary cluster manager", AliyunGoClientFailure}) {
				return resource.NonRetryableError(fmt.Errorf("Deleting container application %s got an error: %#v.", appName, err))
			}
		}

		var project cs.GetProjectResponse
		if err := invoker.Run(func() error {
			cluster, certs, err := csService.GetContainerClusterAndCertsByName(clusterName)
			if err != nil {
				return err
			}
			raw, err := client.WithCsProjectClient(cluster.ClusterID, cluster.MasterURL, *certs, func(csProjectClient *cs.ProjectClient) (interface{}, error) {
				return csProjectClient.GetProject(appName)
			})
			if err != nil {
				return err
			}
			project, _ = raw.(cs.GetProjectResponse)
			return nil
		}); err != nil {
			if IsExpectedErrors(err, []string{"Not Found", "Unable to reach primary cluster manager"}) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Getting container application %s got an error: %#v.", appName, err))
		}
		if project.Name == "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting container application %s timeout: %#v.", appName, err))
	})
}
