package nutanix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixUserGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixUserGroupRead,
		Schema: map[string]*schema.Schema{
			"user_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user_group_name", "user_group_distinguished_name"},
			},
			"user_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user_group_id", "user_group_distinguished_name"},
			},
			"user_group_distinguished_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user_group_id", "user_group_name"},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
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
			"categories": categoriesSchema(),
			"owner_reference": {
				Type:     schema.TypeMap,
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
				Type:     schema.TypeMap,
				Optional: true,
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
			"user_group_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"directory_service_user_group": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"distinguished_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_service_reference": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uuid": {
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
						"default_user_principal_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"project_reference_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
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
			"access_control_policy_reference_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Group: %s", d.Id())

	// Get client connection
	conn := meta.(*Client).API

	uuid, iok := d.GetOk("user_group_id")
	name, nok := d.GetOk("user_group_name")
	dname, dnok := d.GetOk("user_group_distinguished_name")

	if !iok && !nok && !dnok {
		return fmt.Errorf("please provide one of user_group_id, user_group_distinguished_name or user_group_name attributes")
	}

	var reqErr error
	var resp *v3.UserGroupIntentResponse

	if iok {
		resp, reqErr = findUserGroupByUUID(conn, uuid.(string))
	}

	if dnok {
		resp, reqErr = findUserGroupByDistinguishedName(conn, dname.(string))
	}

	if nok {
		resp, reqErr = findUserGroupByName(conn, name.(string))
	}

	if reqErr != nil {
		if strings.Contains(fmt.Sprint(reqErr), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}
		return fmt.Errorf("error reading group with error %s", reqErr)
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return fmt.Errorf("error setting metadata for group UUID(%s), %s", d.Id(), err)
	}
	if err := d.Set("categories", c); err != nil {
		return fmt.Errorf("error setting categories for group UUID(%s), %s", d.Id(), err)
	}

	if err := d.Set("owner_reference", flattenReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return fmt.Errorf("error setting owner_reference for group UUID(%s), %s", d.Id(), err)
	}
	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Status.Resources.DisplayName))

	if err := d.Set("state", resp.Status.State); err != nil {
		return fmt.Errorf("error setting state for group UUID(%s), %s", d.Id(), err)
	}

	if err := d.Set("directory_service_user_group", flattenDirectoryServiceUserGroup(resp.Status.Resources.DirectoryServiceUserGroup)); err != nil {
		return fmt.Errorf("error setting state for group UUID(%s), %s", d.Id(), err)
	}

	if err := d.Set("user_group_type", resp.Status.Resources.UserGroupType); err != nil {
		return fmt.Errorf("error setting state for group UUID(%s), %s", d.Id(), err)
	}

	if err := d.Set("display_name", resp.Status.Resources.DisplayName); err != nil {
		return fmt.Errorf("error setting state for group UUID(%s), %s", d.Id(), err)
	}

	if err := d.Set("project_reference_list", flattenArrayReferenceValues(resp.Status.Resources.ProjectsReferenceList)); err != nil {
		return fmt.Errorf("error setting state for group UUID(%s), %s", d.Id(), err)
	}

	refe := flattenArrayReferenceValues(resp.Status.Resources.AccessControlPolicyReferenceList)

	if err := d.Set("access_control_policy_reference_list", refe); err != nil {
		return fmt.Errorf("error setting state for group UUID(%s), %s", d.Id(), err)
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}

func flattenDirectoryServiceUserGroup(dsu *v3.DirectoryServiceUserGroup) []interface{} {
	if dsu != nil {
		directoryServiceUserMap := map[string]interface{}{}

		if dsu.DistinguishedName != nil {
			directoryServiceUserMap["distinguished_name"] = dsu.DistinguishedName
		}

		if dsu.DirectoryServiceReference != nil {
			directoryServiceUserMap["directory_service_reference"] = []interface{}{flattenReferenceValues(dsu.DirectoryServiceReference)}
		}
		return []interface{}{directoryServiceUserMap}
	}
	return nil
}

func findUserGroupByName(conn *v3.Client, name string) (*v3.UserGroupIntentResponse, error) {
	return findUserGroupByAttribute(conn, matchUserGroupByName, name)
}

func findUserGroupByDistinguishedName(conn *v3.Client, name string) (*v3.UserGroupIntentResponse, error) {
	return findUserGroupByAttribute(conn, matchUserGroupByDistinguishedName, name)
}

func findUserGroupByAttribute(conn *v3.Client, matches func(*v3.UserGroupIntentResponse, string) bool, targetAttributeValue string) (*v3.UserGroupIntentResponse, error) {
	//filter := fmt.Sprintf("name==%s", name)
	resp, err := conn.V3.ListAllUserGroup("")
	if err != nil {
		return nil, err
	}

	entities := resp.Entities

	found := make([]*v3.UserGroupIntentResponse, 0)
	for _, v := range entities {
		// if *v.Status.Resources.DisplayName == targetAttributeValue {
		if matches(v, targetAttributeValue) {
			found = append(found, v)
		}
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("your query returned more than one result. Please use uuid argument instead")
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("user with the given name, not found")
	}

	return findUserGroupByUUID(conn, *found[0].Metadata.UUID)
}

func findUserGroupByUUID(conn *v3.Client, uuid string) (*v3.UserGroupIntentResponse, error) {
	log.Printf("finding group via uuid: %s", uuid)
	return conn.V3.GetUserGroup(uuid)
}

func matchUserGroupByDistinguishedName(userGroup *v3.UserGroupIntentResponse, name string) bool {
	if userGroup != nil &&
		userGroup.Status != nil &&
		userGroup.Status.Resources != nil &&
		userGroup.Status.Resources.DirectoryServiceUserGroup != nil &&
		userGroup.Status.Resources.DirectoryServiceUserGroup.DistinguishedName != nil &&
		*userGroup.Status.Resources.DirectoryServiceUserGroup.DistinguishedName == name {
		return true
	}
	return false
}

func matchUserGroupByName(userGroup *v3.UserGroupIntentResponse, name string) bool {
	if userGroup != nil &&
		userGroup.Status != nil &&
		userGroup.Status.Resources != nil &&
		userGroup.Status.Resources.DisplayName != nil &&
		*userGroup.Status.Resources.DisplayName == name {
		return true
	}
	return false
}
