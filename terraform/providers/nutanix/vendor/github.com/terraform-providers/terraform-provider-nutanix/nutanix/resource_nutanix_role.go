package nutanix

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spf13/cast"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixRoleCreate,
		Read:   resourceNutanixRoleRead,
		Update: resourceNutanixRoleUpdate,
		Delete: resourceNutanixRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"spec_hash": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
			"owner_reference": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"project_reference": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "project",
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission_reference_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "permission",
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
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["kind"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["uuid"].(string)))
					return hashcode.String(buf.String())
				},
			},
		},
	}
}

func resourceNutanixRoleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	request := &v3.Role{}
	spec := &v3.RoleSpec{}
	metadata := &v3.Metadata{}
	role := &v3.RoleResources{}

	name, nameOk := d.GetOk("name")
	permissions, permissionsOk := d.GetOk("permission_reference_list")

	if !nameOk && !permissionsOk {
		return fmt.Errorf("please provide the required `name` and `permission_reference_list`  attribute")
	}

	if err := getMetadataAttributesV2(d, metadata, "role"); err != nil {
		return err
	}

	spec.Name = utils.StringPtr(name.(string))
	if desc, descOk := d.GetOk("description"); descOk {
		spec.Description = utils.StringPtr(desc.(string))
	}
	role.PermissionReferenceList = validateArrayRef(permissions.(*schema.Set), nil)

	if name, ok := d.GetOk("name"); ok {
		spec.Name = utils.StringPtr(name.(string))
	}
	spec.Resources = role
	request.Metadata = metadata
	request.Spec = spec

	resp, err := conn.V3.CreateRole(request)
	if err != nil {
		return fmt.Errorf("error creating Nutanix Role %s: %+v", utils.StringValue(spec.Name), err)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the Role to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING", "PENDING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		id := d.Id()
		d.SetId("")
		return fmt.Errorf("error waiting for role  id (%s) to create: %+v", id, err)
	}

	// Setting Description because in Get request is not present.
	d.Set("description", utils.StringValue(resp.Spec.Description))

	d.SetId(utils.StringValue(resp.Metadata.UUID))

	return resourceNutanixRoleRead(d, meta)
}

func resourceNutanixRoleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	id := d.Id()
	resp, err := conn.V3.GetRole(id)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		errDel := resourceNutanixRoleDelete(d, meta)
		if errDel != nil {
			return fmt.Errorf("error deleting role (%s) after read error: %+v", id, errDel)
		}
		d.SetId("")
		return fmt.Errorf("error reading role id (%s): %+v", id, err)
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}
	if err := d.Set("project_reference", flattenReferenceValuesList(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", flattenReferenceValuesList(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	d.Set("api_version", resp.APIVersion)

	if status := resp.Status; status != nil {
		if err := d.Set("name", utils.StringValue(resp.Status.Name)); err != nil {
			return err
		}
		if err := d.Set("description", utils.StringValue(resp.Status.Description)); err != nil {
			return err
		}
		if err := d.Set("state", utils.StringValue(resp.Status.State)); err != nil {
			return err
		}

		if res := status.Resources; res != nil {
			if err := d.Set("permission_reference_list", flattenArrayReferenceValues(status.Resources.PermissionReferenceList)); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceNutanixRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	request := &v3.Role{}
	metadata := &v3.Metadata{}
	res := &v3.RoleResources{}
	spec := &v3.RoleSpec{}

	id := d.Id()
	response, err := conn.V3.GetRole(id)

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}
		return fmt.Errorf("error retrieving for role id (%s) :%+v", id, err)
	}

	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if response.Spec != nil {
		spec = response.Spec

		if response.Spec.Resources != nil {
			res = response.Spec.Resources
		}
	}

	if d.HasChange("categories") {
		metadata.Categories = expandCategories(d.Get("categories"))
	}
	if d.HasChange("owner_reference") {
		metadata.OwnerReference = validateRefList(d.Get("owner_reference").([]interface{}), nil)
	}
	if d.HasChange("project_reference") {
		metadata.ProjectReference = validateRefList(d.Get("project_reference").([]interface{}), utils.StringPtr("project"))
	}
	if d.HasChange("name") {
		spec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		spec.Description = utils.StringPtr(d.Get("description").(string))
	}

	if d.HasChange("permission_reference_list") {
		res.PermissionReferenceList = validateArrayRef(d.Get("permission_reference_list"), nil)
	}

	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	resp, errUpdate := conn.V3.UpdateRole(d.Id(), request)
	if errUpdate != nil {
		return fmt.Errorf("error updating role id %s): %s", d.Id(), errUpdate)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for role (%s) to update: %s", d.Id(), err)
	}
	// Setting Description because in Get request is not present.
	d.Set("description", utils.StringValue(resp.Spec.Description))

	return resourceNutanixRoleRead(d, meta)
}

func resourceNutanixRoleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	resp, err := conn.V3.DeleteRole(d.Id())
	if err != nil {
		return fmt.Errorf("error deleting role id %s): %s", d.Id(), err)
	}

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING", "DELETED_PENDING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, cast.ToString(resp.Status.ExecutionContext.TaskUUID)),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for role (%s) to update: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceNutanixRoleExists(conn *v3.Client, name string) (*string, error) {
	var accessUUID *string

	filter := fmt.Sprintf("name==%s", name)
	accessList, err := conn.V3.ListAllRole(filter)

	if err != nil {
		return nil, err
	}

	for _, access := range accessList.Entities {
		if utils.StringValue(access.Status.Name) == name {
			accessUUID = access.Metadata.UUID
		}
	}
	return accessUUID, nil
}
