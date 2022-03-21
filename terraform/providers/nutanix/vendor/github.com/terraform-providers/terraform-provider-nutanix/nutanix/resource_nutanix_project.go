package nutanix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spf13/cast"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixProjectCreate,
		Read:   resourceNutanixProjectRead,
		Update: resourceNutanixProjectUpdate,
		Delete: resourceNutanixProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"resource_domain": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"units": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"limit": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"resource_type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"account_reference_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Default:  "account",
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"environment_reference_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "environment",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"default_subnet_reference": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "subnet",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"user_reference_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "user",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"external_user_group_reference_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "user_group",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"subnet_reference_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "subnet",
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"external_network_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"owner_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"categories": categoriesSchema(),
			"api_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceNutanixProjectCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	req := &v3.Project{
		Spec:       expandProjectSpec(d),
		Metadata:   expandMetadata(d, "project"),
		APIVersion: d.Get("api_version").(string),
	}

	resp, err := conn.V3.CreateProject(req)
	if err != nil {
		return err
	}

	uuid := *resp.Metadata.UUID
	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the Project to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    vmTimeout,
		Delay:      vmDelay,
		MinTimeout: vmMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForState(); errWaitTask != nil {
		return fmt.Errorf("error waiting for project(%s) to create: %s", uuid, errWaitTask)
	}

	d.SetId(uuid)
	return resourceNutanixProjectRead(d, meta)
}

func resourceNutanixProjectRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	project, err := conn.V3.GetProject(d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return err
	}

	m, c := setRSEntityMetadata(project.Metadata)

	if err := d.Set("name", project.Status.Name); err != nil {
		return fmt.Errorf("error setting `name` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("description", project.Status.Descripion); err != nil {
		return fmt.Errorf("error setting `description` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("state", project.Status.State); err != nil {
		return fmt.Errorf("error setting `state` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("is_default", project.Status.Resources.IsDefault); err != nil {
		return fmt.Errorf("error setting `is_default` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("resource_domain", flattenResourceDomain(project.Spec.Resources.ResourceDomain)); err != nil {
		return fmt.Errorf("error setting `resource_domain` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("account_reference_list", flattenReferenceList(project.Spec.Resources.AccountReferenceList)); err != nil {
		return fmt.Errorf("error setting `account_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("environment_reference_list", flattenReferenceList(project.Spec.Resources.EnvironmentReferenceList)); err != nil {
		return fmt.Errorf("error setting `environment_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("default_subnet_reference", []interface{}{flattenReference(project.Spec.Resources.DefaultSubnetReference)}); err != nil {
		return fmt.Errorf("error setting `default_subnet_reference` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("user_reference_list", flattenReferenceList(project.Spec.Resources.UserReferenceList)); err != nil {
		return fmt.Errorf("error setting `user_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("external_user_group_reference_list",
		flattenReferenceList(project.Spec.Resources.ExternalUserGroupReferenceList)); err != nil {
		return fmt.Errorf("error setting `external_user_group_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("subnet_reference_list", flattenReferenceList(project.Spec.Resources.SubnetReferenceList)); err != nil {
		return fmt.Errorf("error setting `subnet_reference_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("external_network_list", flattenReferenceList(project.Spec.Resources.ExternalNetworkList)); err != nil {
		return fmt.Errorf("error setting `external_network_list` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("metadata", m); err != nil {
		return fmt.Errorf("error setting `metadata` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("project_reference", flattenReferenceValues(project.Metadata.ProjectReference)); err != nil {
		return fmt.Errorf("error setting `project_reference` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("owner_reference", flattenReferenceValues(project.Metadata.OwnerReference)); err != nil {
		return fmt.Errorf("error setting `owner_reference` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("categories", c); err != nil {
		return fmt.Errorf("error setting `categories` for Project(%s): %s", d.Id(), err)
	}
	if err := d.Set("api_version", project.APIVersion); err != nil {
		return fmt.Errorf("error setting `api_version` for Project(%s): %s", d.Id(), err)
	}

	return nil
}

func resourceNutanixProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	project, err := conn.V3.GetProject(d.Id())
	if err != nil {
		return err
	}
	project.Status = nil

	if d.HasChange("name") {
		project.Spec.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		project.Spec.Descripion = d.Get("description").(string)
	}
	if d.HasChange("resource_domain") {
		project.Spec.Resources.ResourceDomain = expandResourceDomain(d)
	}
	if d.HasChange("account_reference_list") {
		project.Spec.Resources.AccountReferenceList = expandReferenceList(d, "account_reference_list")
	}
	if d.HasChange("environment_reference_list") {
		project.Spec.Resources.EnvironmentReferenceList = expandReferenceList(d, "environment_reference_list")
	}
	if d.HasChange("default_subnet_reference") {
		project.Spec.Resources.DefaultSubnetReference = expandReferenceList(d, "default_subnet_reference")[0]
	}
	if d.HasChange("user_reference_list") {
		project.Spec.Resources.UserReferenceList = expandReferenceSet(d, "user_reference_list")
	}
	if d.HasChange("external_user_group_reference_list") {
		project.Spec.Resources.ExternalUserGroupReferenceList = expandReferenceSet(d, "external_user_group_reference_list")
	}
	if d.HasChange("subnet_reference_list") {
		project.Spec.Resources.SubnetReferenceList = expandReferenceList(d, "subnet_reference_list")
	}
	if d.HasChange("external_network_list") {
		project.Spec.Resources.ExternalNetworkList = expandReferenceList(d, "external_network_list")
	}
	if d.HasChange("metadata") || d.HasChange("project_reference") ||
		d.HasChange("owner_reference") || d.HasChange("categories") {
		if err = getMetadataAttributes(d, project.Metadata, "project"); err != nil {
			return fmt.Errorf("error expanding metadata: %+v", err)
		}
	}
	if d.HasChange("api_version") {
		project.APIVersion = d.Get("api_version").(string)
	}

	resp, err := conn.V3.UpdateProject(d.Id(), project)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}
		return err
	}

	uuid := *resp.Metadata.UUID
	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the Project to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    vmTimeout,
		Delay:      vmDelay,
		MinTimeout: vmMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForState(); errWaitTask != nil {
		return fmt.Errorf("error waiting for project(%s) to update: %s", uuid, errWaitTask)
	}

	return resourceNutanixProjectRead(d, meta)
}

func resourceNutanixProjectDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	resp, err := conn.V3.DeleteProject(d.Id())
	if err != nil {
		return fmt.Errorf("error deleting project id %s): %s", d.Id(), err)
	}

	// Wait for the Project to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING", "DELETED_PENDING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, cast.ToString(resp.Status.ExecutionContext.TaskUUID)),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for project (%s) to update: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func expandProjectSpec(d *schema.ResourceData) *v3.ProjectSpec {
	return &v3.ProjectSpec{
		Name:       d.Get("name").(string),
		Descripion: d.Get("description").(string),
		Resources: &v3.ProjectResources{
			ResourceDomain:                 expandResourceDomain(d),
			AccountReferenceList:           expandReferenceList(d, "account_reference_list"),
			EnvironmentReferenceList:       expandReferenceList(d, "environment_reference_list"),
			DefaultSubnetReference:         expandReferenceList(d, "default_subnet_reference")[0],
			UserReferenceList:              expandReferenceSet(d, "user_reference_list"),
			ExternalUserGroupReferenceList: expandReferenceSet(d, "external_user_group_reference_list"),
			SubnetReferenceList:            expandReferenceList(d, "subnet_reference_list"),
			ExternalNetworkList:            expandReferenceList(d, "external_network_list"),
		},
	}
}

func expandResourceDomain(d *schema.ResourceData) *v3.ResourceDomain {
	resourceDomain, ok := d.GetOk("resource_domain")
	if !ok {
		return nil
	}
	resources := cast.ToStringMap(resourceDomain.([]interface{})[0])["resources"].([]interface{})

	rs := make([]*v3.Resources, len(resources))
	for i, resource := range resources {
		r := cast.ToStringMap(resource)
		rs[i] = &v3.Resources{
			Limit:        utils.Int64Ptr(cast.ToInt64(r["limit"])),
			ResourceType: cast.ToString(r["resource_type"]),
		}
	}
	return &v3.ResourceDomain{Resources: rs}
}

func flattenResourceDomain(resourceDomain *v3.ResourceDomain) (res []map[string]interface{}) {
	if resourceDomain != nil {
		if len(resourceDomain.Resources) > 0 {
			resources := make([]map[string]interface{}, len(resourceDomain.Resources))

			for i, r := range resourceDomain.Resources {
				resources[i] = map[string]interface{}{
					"units":         r.Units,
					"value":         cast.ToInt64(r.Value),
					"limit":         cast.ToInt64(r.Limit),
					"resource_type": r.ResourceType,
				}
			}
			res = append(res, map[string]interface{}{"resources": resources})
		}
	}
	return
}

func expandReferenceByMap(reference map[string]interface{}) *v3.ReferenceValues {
	return &v3.ReferenceValues{
		Kind: cast.ToString(reference["kind"]),
		Name: cast.ToString(reference["name"]),
		UUID: cast.ToString(reference["uuid"]),
	}
}

// func expandReference(d *schema.ResourceData, key string) *v3.ReferenceValues {
// 	return expandReferenceByMap(cast.ToStringMap(d.Get(key)))
// }

func expandReferenceList(d *schema.ResourceData, key string) []*v3.ReferenceValues {
	references := d.Get(key).([]interface{})
	list := make([]*v3.ReferenceValues, len(references))

	for i, r := range references {
		list[i] = expandReferenceByMap(cast.ToStringMap(r))
	}
	return list
}

func expandReferenceSet(d *schema.ResourceData, key string) []*v3.ReferenceValues {
	references := d.Get(key).(*schema.Set).List()
	list := make([]*v3.ReferenceValues, len(references))

	for i, r := range references {
		list[i] = expandReferenceByMap(cast.ToStringMap(r))
	}
	return list
}

func expandMetadata(d *schema.ResourceData, kind string) *v3.Metadata {
	metadata := new(v3.Metadata)

	if err := getMetadataAttributes(d, metadata, kind); err != nil {
		log.Printf("Error expanding metadata: %+v", err)
	}
	return metadata
}
