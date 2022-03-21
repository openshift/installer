package nutanix

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	// UserKind Represents kind of resource
	UserKind = "user"
)

var (
	userDelay      = 10 * time.Second
	userMinTimeout = 3 * time.Second
)

func resourceNutanixUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixUserCreate,
		Read:   resourceNutanixUserRead,
		Update: resourceNutanixUserUpdate,
		Delete: resourceNutanixUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
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
			"directory_service_user": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_principal_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							//ValidateFunc: validation.StringInSlice([]string{"role"}, false),
						},
						"directory_service_reference": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "directory_service",
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
						"default_user_principal_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"identity_provider_user": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							//ValidateFunc: validation.StringInSlice([]string{"role"}, false),
						},
						"identity_provider_reference": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "identity_provider",
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
					},
				},
			},
			"user_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_reference_list": {
				Type:     schema.TypeSet,
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
				Type:     schema.TypeSet,
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixUserCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating User: %s", d.Get("name").(string))
	client := meta.(*Client)
	conn := client.API
	timeout := client.WaitTimeout

	if client.WaitTimeout == 0 {
		timeout = 10
	}

	request := &v3.UserIntentInput{}

	metadata := &v3.Metadata{}

	if err := getMetadataAttributes(d, metadata, "user"); err != nil {
		return err
	}

	spec := &v3.UserSpec{
		Resources: &v3.UserResources{
			DirectoryServiceUser: expandDirectoryServiceUser(d),
			IdentityProviderUser: expandIdentityProviderUser(d),
		},
	}

	request.Metadata = metadata
	request.Spec = spec

	// Make request to the API
	resp, err := conn.V3.CreateUser(request)
	if err != nil {
		return fmt.Errorf("error creating Nutanix User: %+v", err)
	}

	UUID := *resp.Metadata.UUID
	// set terraform state
	d.SetId(UUID)

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the Image to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    time.Duration(timeout) * time.Minute,
		Delay:      userDelay,
		MinTimeout: userMinTimeout,
	}

	if _, errw := stateConf.WaitForState(); errw != nil {
		// delErr := resourceNutanixUserDelete(d, meta)
		// if delErr != nil {
		// 	return fmt.Errorf("error waiting for image (%s) to delete in creation: %s", d.Id(), delErr)
		// }
		d.SetId("")
		return fmt.Errorf("error waiting for user (%s) to create: %s", UUID, errw)
	}
	return resourceNutanixUserRead(d, meta)
}

func resourceNutanixUserRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading User: %s", d.Id())

	// Get client connection
	conn := meta.(*Client).API
	uuid := d.Id()

	// Make request to the API
	resp, err := conn.V3.GetUser(uuid)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading user UUID (%s) with error %s", uuid, err)
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err = d.Set("metadata", m); err != nil {
		return fmt.Errorf("error setting metadata for user UUID(%s), %s", d.Id(), err)
	}
	if err = d.Set("categories", c); err != nil {
		return fmt.Errorf("error setting categories for user UUID(%s), %s", d.Id(), err)
	}

	if err = d.Set("owner_reference", flattenReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return fmt.Errorf("error setting owner_reference for user UUID(%s), %s", d.Id(), err)
	}
	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Status.Name))

	if err = d.Set("state", resp.Status.State); err != nil {
		return fmt.Errorf("error setting state for user UUID(%s), %s", d.Id(), err)
	}

	if err = d.Set("directory_service_user", flattenDirectoryServiceUser(resp.Status.Resources.DirectoryServiceUser)); err != nil {
		return fmt.Errorf("error setting directory_service_user for user UUID(%s), %s", d.Id(), err)
	}

	//TODO: change to status when API is fixed
	if err = d.Set("identity_provider_user", flattenIdentityProviderUser(resp.Spec.Resources.IdentityProviderUser)); err != nil {
		return fmt.Errorf("error setting identity_provider_user for user UUID(%s), %s", d.Id(), err)
	}

	if err = d.Set("user_type", resp.Status.Resources.UserType); err != nil {
		return fmt.Errorf("error setting user_type for user UUID(%s), %s", d.Id(), err)
	}

	if err = d.Set("display_name", resp.Status.Resources.DisplayName); err != nil {
		return fmt.Errorf("error setting display_name for user UUID(%s), %s", d.Id(), err)
	}

	if err := d.Set("project_reference_list", flattenArrayReferenceValues(resp.Status.Resources.ProjectsReferenceList)); err != nil {
		return fmt.Errorf("error setting project_reference_list for user UUID(%s), %s", d.Id(), err)
	}

	if err := d.Set("access_control_policy_reference_list", flattenArrayReferenceValues(resp.Status.Resources.AccessControlPolicyReferenceList)); err != nil {
		return fmt.Errorf("error setting access_control_policy_reference_list for user UUID(%s), %s", d.Id(), err)
	}

	//TODO:
	//identity_provider_user

	return nil
}

func resourceNutanixUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	conn := client.API
	timeout := client.WaitTimeout

	if client.WaitTimeout == 0 {
		timeout = 10
	}

	// get state
	request := &v3.UserIntentInput{}
	metadata := &v3.Metadata{}
	spec := &v3.UserSpec{}
	res := &v3.UserResources{}

	response, err := conn.V3.GetUser(d.Id())

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}
		return err
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
		or := d.Get("owner_reference").(map[string]interface{})
		metadata.OwnerReference = validateRef(or)
	}

	if d.HasChange("project_reference") {
		pr := d.Get("project_reference").(map[string]interface{})
		metadata.ProjectReference = validateRef(pr)
	}

	if d.HasChange("directory_service_user") {
		res.DirectoryServiceUser = expandDirectoryServiceUser(d)
	}

	if d.HasChange("identity_provider_user") {
		res.IdentityProviderUser = expandIdentityProviderUser(d)
	}

	request.Metadata = metadata
	request.Spec = spec

	resp, errUpdate := conn.V3.UpdateUser(d.Id(), request)

	if errUpdate != nil {
		return fmt.Errorf("error updating user(%s) %s", d.Id(), errUpdate)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the Image to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    time.Duration(timeout) * time.Minute,
		Delay:      userDelay,
		MinTimeout: userMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		// delErr := resourceNutanixUserDelete(d, meta)
		// if delErr != nil {
		// 	return fmt.Errorf("error waiting for image (%s) to delete in update: %s", d.Id(), delErr)
		// }
		uuid := d.Id()
		d.SetId("")
		return fmt.Errorf("error waiting for user (%s) to update: %s", uuid, err)
	}

	return resourceNutanixUserRead(d, meta)
}

func resourceNutanixUserDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting User: %s", d.Get("display_name").(string))

	client := meta.(*Client)
	conn := client.API
	timeout := client.WaitTimeout

	if client.WaitTimeout == 0 {
		timeout = 10
	}

	UUID := d.Id()

	resp, err := conn.V3.DeleteUser(UUID)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}
		return err
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the Image to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    time.Duration(timeout) * time.Minute,
		Delay:      userDelay,
		MinTimeout: userMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		d.SetId("")
		return fmt.Errorf("error waiting for user (%s) to delete: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] USER DELETED")

	d.SetId("")
	return nil
}

func expandDirectoryServiceUser(d *schema.ResourceData) *v3.DirectoryServiceUser {
	directoryServiceUserState, ok := d.GetOk("directory_service_user")
	if !ok {
		return nil
	}

	directoryServiceUserMap := directoryServiceUserState.([]interface{})[0].(map[string]interface{})
	directoryServiceUser := &v3.DirectoryServiceUser{}

	if upn, ok := directoryServiceUserMap["user_principal_name"]; ok {
		directoryServiceUser.UserPrincipalName = utils.StringPtr(upn.(string))
	}

	if dpr, ok := directoryServiceUserMap["directory_service_reference"]; ok {
		directoryServiceUser.DirectoryServiceReference = expandReference(dpr.([]interface{})[0].(map[string]interface{}))
	}

	if !reflect.DeepEqual(*directoryServiceUser, v3.DirectoryServiceUser{}) {
		return directoryServiceUser
	}
	return nil
}

func expandIdentityProviderUser(d *schema.ResourceData) *v3.IdentityProvider {
	identityProviderState, ok := d.GetOk("identity_provider_user")
	if !ok {
		return nil
	}

	identityProviderMap := identityProviderState.([]interface{})[0].(map[string]interface{})
	identityProvider := &v3.IdentityProvider{}

	if username, ok := identityProviderMap["username"]; ok {
		identityProvider.Username = utils.StringPtr(username.(string))
	}

	if ipr, ok := identityProviderMap["identity_provider_reference"]; ok {
		identityProvider.IdentityProviderReference = expandReference(ipr.([]interface{})[0].(map[string]interface{}))
	}

	if !reflect.DeepEqual(*identityProvider, v3.IdentityProvider{}) {
		return identityProvider
	}
	return nil
}

func flattenDirectoryServiceUser(dsu *v3.DirectoryServiceUser) []interface{} {
	log.Print("in flattenDirectoryServiceUser")
	if dsu != nil {
		directoryServiceUserMap := map[string]interface{}{}

		if dsu.DefaultUserPrincipalName != nil {
			directoryServiceUserMap["default_user_principal_name"] = dsu.DefaultUserPrincipalName
		}

		if dsu.UserPrincipalName != nil {
			directoryServiceUserMap["user_principal_name"] = dsu.UserPrincipalName
		}

		if dsu.DirectoryServiceReference != nil {
			directoryServiceUserMap["directory_service_reference"] = []interface{}{flattenReferenceValues(dsu.DirectoryServiceReference)}
		}
		return []interface{}{directoryServiceUserMap}
	}
	return nil
}

func flattenIdentityProviderUser(ipu *v3.IdentityProvider) []interface{} {
	if ipu != nil {
		identityProviderUserMap := map[string]interface{}{}

		if ipu.Username != nil {
			identityProviderUserMap["username"] = ipu.Username
		}

		if ipu.IdentityProviderReference != nil {
			identityProviderUserMap["identity_provider_reference"] = []interface{}{flattenReferenceValues(ipu.IdentityProviderReference)}
		}

		return []interface{}{identityProviderUserMap}
	}
	return nil
}
