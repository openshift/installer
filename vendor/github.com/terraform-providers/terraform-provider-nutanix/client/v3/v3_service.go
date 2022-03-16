package v3

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Operations ...
type Operations struct {
	client *client.Client
}

// Service ...
type Service interface {
	CreateVM(createRequest *VMIntentInput) (*VMIntentResponse, error)
	DeleteVM(uuid string) (*DeleteResponse, error)
	GetVM(uuid string) (*VMIntentResponse, error)
	ListVM(getEntitiesRequest *DSMetadata) (*VMListIntentResponse, error)
	UpdateVM(uuid string, body *VMIntentInput) (*VMIntentResponse, error)
	CreateSubnet(createRequest *SubnetIntentInput) (*SubnetIntentResponse, error)
	DeleteSubnet(uuid string) (*DeleteResponse, error)
	GetSubnet(uuid string) (*SubnetIntentResponse, error)
	ListSubnet(getEntitiesRequest *DSMetadata) (*SubnetListIntentResponse, error)
	UpdateSubnet(uuid string, body *SubnetIntentInput) (*SubnetIntentResponse, error)
	CreateImage(createRequest *ImageIntentInput) (*ImageIntentResponse, error)
	DeleteImage(uuid string) (*DeleteResponse, error)
	GetImage(uuid string) (*ImageIntentResponse, error)
	ListImage(getEntitiesRequest *DSMetadata) (*ImageListIntentResponse, error)
	UpdateImage(uuid string, body *ImageIntentInput) (*ImageIntentResponse, error)
	UploadImage(uuid, filepath string) error
	CreateOrUpdateCategoryKey(body *CategoryKey) (*CategoryKeyStatus, error)
	ListCategories(getEntitiesRequest *CategoryListMetadata) (*CategoryKeyListResponse, error)
	DeleteCategoryKey(name string) error
	GetCategoryKey(name string) (*CategoryKeyStatus, error)
	ListCategoryValues(name string, getEntitiesRequest *CategoryListMetadata) (*CategoryValueListResponse, error)
	CreateOrUpdateCategoryValue(name string, body *CategoryValue) (*CategoryValueStatus, error)
	GetCategoryValue(name string, value string) (*CategoryValueStatus, error)
	DeleteCategoryValue(name string, value string) error
	GetCategoryQuery(query *CategoryQueryInput) (*CategoryQueryResponse, error)
	UpdateNetworkSecurityRule(uuid string, body *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error)
	ListNetworkSecurityRule(getEntitiesRequest *DSMetadata) (*NetworkSecurityRuleListIntentResponse, error)
	GetNetworkSecurityRule(uuid string) (*NetworkSecurityRuleIntentResponse, error)
	DeleteNetworkSecurityRule(uuid string) (*DeleteResponse, error)
	CreateNetworkSecurityRule(request *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error)
	ListCluster(getEntitiesRequest *DSMetadata) (*ClusterListIntentResponse, error)
	GetCluster(uuid string) (*ClusterIntentResponse, error)
	UpdateVolumeGroup(uuid string, body *VolumeGroupInput) (*VolumeGroupResponse, error)
	ListVolumeGroup(getEntitiesRequest *DSMetadata) (*VolumeGroupListResponse, error)
	GetVolumeGroup(uuid string) (*VolumeGroupResponse, error)
	DeleteVolumeGroup(uuid string) error
	CreateVolumeGroup(request *VolumeGroupInput) (*VolumeGroupResponse, error)
	ListAllVM(filter string) (*VMListIntentResponse, error)
	ListAllSubnet(filter string) (*SubnetListIntentResponse, error)
	ListAllNetworkSecurityRule(filter string) (*NetworkSecurityRuleListIntentResponse, error)
	ListAllImage(filter string) (*ImageListIntentResponse, error)
	ListAllCluster(filter string) (*ClusterListIntentResponse, error)
	ListAllCategoryValues(categoryName, filter string) (*CategoryValueListResponse, error)
	GetTask(taskUUID string) (*TasksResponse, error)
	GetHost(taskUUID string) (*HostResponse, error)
	ListHost(getEntitiesRequest *DSMetadata) (*HostListResponse, error)
	ListAllHost() (*HostListResponse, error)
	CreateProject(request *Project) (*Project, error)
	GetProject(projectUUID string) (*Project, error)
	ListProject(getEntitiesRequest *DSMetadata) (*ProjectListResponse, error)
	ListAllProject(filter string) (*ProjectListResponse, error)
	UpdateProject(uuid string, body *Project) (*Project, error)
	DeleteProject(uuid string) (*DeleteResponse, error)
	CreateAccessControlPolicy(request *AccessControlPolicy) (*AccessControlPolicy, error)
	GetAccessControlPolicy(accessControlPolicyUUID string) (*AccessControlPolicy, error)
	ListAccessControlPolicy(getEntitiesRequest *DSMetadata) (*AccessControlPolicyListResponse, error)
	ListAllAccessControlPolicy(filter string) (*AccessControlPolicyListResponse, error)
	UpdateAccessControlPolicy(uuid string, body *AccessControlPolicy) (*AccessControlPolicy, error)
	DeleteAccessControlPolicy(uuid string) (*DeleteResponse, error)
	CreateRole(request *Role) (*Role, error)
	GetRole(uuid string) (*Role, error)
	ListRole(getEntitiesRequest *DSMetadata) (*RoleListResponse, error)
	ListAllRole(filter string) (*RoleListResponse, error)
	UpdateRole(uuid string, body *Role) (*Role, error)
	DeleteRole(uuid string) (*DeleteResponse, error)
	CreateUser(request *UserIntentInput) (*UserIntentResponse, error)
	GetUser(userUUID string) (*UserIntentResponse, error)
	UpdateUser(uuid string, body *UserIntentInput) (*UserIntentResponse, error)
	DeleteUser(uuid string) (*DeleteResponse, error)
	ListUser(getEntitiesRequest *DSMetadata) (*UserListResponse, error)
	ListAllUser(filter string) (*UserListResponse, error)
	GetUserGroup(userUUID string) (*UserGroupIntentResponse, error)
	ListUserGroup(getEntitiesRequest *DSMetadata) (*UserGroupListResponse, error)
	ListAllUserGroup(filter string) (*UserGroupListResponse, error)
	GetPermission(permissionUUID string) (*PermissionIntentResponse, error)
	ListPermission(getEntitiesRequest *DSMetadata) (*PermissionListResponse, error)
	ListAllPermission(filter string) (*PermissionListResponse, error)
	GetProtectionRule(uuid string) (*ProtectionRuleResponse, error)
	ListProtectionRules(getEntitiesRequest *DSMetadata) (*ProtectionRulesListResponse, error)
	ListAllProtectionRules(filter string) (*ProtectionRulesListResponse, error)
	CreateProtectionRule(request *ProtectionRuleInput) (*ProtectionRuleResponse, error)
	UpdateProtectionRule(uuid string, body *ProtectionRuleInput) (*ProtectionRuleResponse, error)
	DeleteProtectionRule(uuid string) (*DeleteResponse, error)
	GetRecoveryPlan(uuid string) (*RecoveryPlanResponse, error)
	ListRecoveryPlans(getEntitiesRequest *DSMetadata) (*RecoveryPlanListResponse, error)
	ListAllRecoveryPlans(filter string) (*RecoveryPlanListResponse, error)
	CreateRecoveryPlan(request *RecoveryPlanInput) (*RecoveryPlanResponse, error)
	UpdateRecoveryPlan(uuid string, body *RecoveryPlanInput) (*RecoveryPlanResponse, error)
	DeleteRecoveryPlan(uuid string) (*DeleteResponse, error)
}

/*CreateVM Creates a VM
 * This operation submits a request to create a VM based on the input parameters.
 *
 * @param body
 * @return *VMIntentResponse
 */
func (op Operations) CreateVM(createRequest *VMIntentInput) (*VMIntentResponse, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/vms", createRequest)
	vmIntentResponse := new(VMIntentResponse)

	if err != nil {
		return nil, err
	}

	return vmIntentResponse, op.client.Do(ctx, req, vmIntentResponse)
}

/*DeleteVM Deletes a VM
 * This operation submits a request to delete a op.
 *
 * @param uuid The uuid of the entity.
 * @return error
 */
func (op Operations) DeleteVM(uuid string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/vms/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)

	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
}

/*GetVM Gets a VM
 * This operation gets a op.
 *
 * @param uuid The uuid of the entity.
 * @return *VMIntentResponse
 */
func (op Operations) GetVM(uuid string) (*VMIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/vms/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	vmIntentResponse := new(VMIntentResponse)

	if err != nil {
		return nil, err
	}

	return vmIntentResponse, op.client.Do(ctx, req, vmIntentResponse)
}

/*ListVM Get a list of VMs This operation gets a list of VMs, allowing for sorting and pagination. Note: Entities that have not been created
 * successfully are not listed.
 *
 * @param getEntitiesRequest @return *VmListIntentResponse
 */
func (op Operations) ListVM(getEntitiesRequest *DSMetadata) (*VMListIntentResponse, error) {
	ctx := context.TODO()
	path := "/vms/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	vmListIntentResponse := new(VMListIntentResponse)

	if err != nil {
		return nil, err
	}

	return vmListIntentResponse, op.client.Do(ctx, req, vmListIntentResponse)
}

/*UpdateVM Updates a VM
 * This operation submits a request to update a VM based on the input parameters.
 *
 * @param uuid The uuid of the entity.
 * @param body
 * @return *VMIntentResponse
 */
func (op Operations) UpdateVM(uuid string, body *VMIntentInput) (*VMIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/vms/%s", uuid)
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	vmIntentResponse := new(VMIntentResponse)

	if err != nil {
		return nil, err
	}

	return vmIntentResponse, op.client.Do(ctx, req, vmIntentResponse)
}

/*CreateSubnet Creates a subnet
 * This operation submits a request to create a subnet based on the input parameters. A subnet is a block of IP addresses.
 *
 * @param body
 * @return *SubnetIntentResponse
 */
func (op Operations) CreateSubnet(createRequest *SubnetIntentInput) (*SubnetIntentResponse, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/subnets", createRequest)
	subnetIntentResponse := new(SubnetIntentResponse)

	if err != nil {
		return nil, err
	}

	return subnetIntentResponse, op.client.Do(ctx, req, subnetIntentResponse)
}

/*DeleteSubnet Deletes a subnet
 * This operation submits a request to delete a subnet.
 *
 * @param uuid The uuid of the entity.
 * @return error if exist error
 */
func (op Operations) DeleteSubnet(uuid string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/subnets/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)

	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
}

/*GetSubnet Gets a subnet entity
 * This operation gets a subnet.
 *
 * @param uuid The uuid of the entity.
 * @return *SubnetIntentResponse
 */
func (op Operations) GetSubnet(uuid string) (*SubnetIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/subnets/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	subnetIntentResponse := new(SubnetIntentResponse)

	if err != nil {
		return nil, err
	}

	// Recheck subnet already exist error
	// if *subnetIntentResponse.Status.State == "ERROR" {
	// 	pretty, _ := json.MarshalIndent(subnetIntentResponse.Status.MessageList, "", "  ")
	// 	return nil, fmt.Errorf("error: %s", string(pretty))
	// }

	return subnetIntentResponse, op.client.Do(ctx, req, subnetIntentResponse)
}

/*ListSubnet Gets a list of subnets This operation gets a list of subnets, allowing for sorting and pagination. Note: Entities that have not
 * been created successfully are not listed.
 *
 * @param getEntitiesRequest @return *SubnetListIntentResponse
 */
func (op Operations) ListSubnet(getEntitiesRequest *DSMetadata) (*SubnetListIntentResponse, error) {
	ctx := context.TODO()
	path := "/subnets/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	subnetListIntentResponse := new(SubnetListIntentResponse)

	if err != nil {
		return nil, err
	}

	return subnetListIntentResponse, op.client.Do(ctx, req, subnetListIntentResponse)
}

/*UpdateSubnet Updates a subnet
 * This operation submits a request to update a subnet based on the input parameters.
 *
 * @param uuid The uuid of the entity.
 * @param body
 * @return *SubnetIntentResponse
 */
func (op Operations) UpdateSubnet(uuid string, body *SubnetIntentInput) (*SubnetIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/subnets/%s", uuid)
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	subnetIntentResponse := new(SubnetIntentResponse)

	if err != nil {
		return nil, err
	}

	return subnetIntentResponse, op.client.Do(ctx, req, subnetIntentResponse)
}

/*CreateImage Creates a IMAGE Images are raw ISO, QCOW2, or VMDK files that are uploaded by a user can be attached to a op. An ISO image is
 * attached as a virtual CD-ROM drive, and QCOW2 and VMDK files are attached as SCSI disks. An image has to be explicitly added to the
 * self-service catalog before users can create VMs from it.
 *
 * @param body @return *ImageIntentResponse
 */
func (op Operations) CreateImage(body *ImageIntentInput) (*ImageIntentResponse, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/images", body)
	imageIntentResponse := new(ImageIntentResponse)

	if err != nil {
		return nil, err
	}

	return imageIntentResponse, op.client.Do(ctx, req, imageIntentResponse)
}

/*UploadImage Uplloads a Image Binary file Images are raw ISO, QCOW2, or VMDK files that are uploaded by a user can be attached to a op. An
 * ISO image is attached as a virtual CD-ROM drive, and QCOW2 and VMDK files are attached as SCSI disks. An image has to be explicitly added
 * to the self-service catalog before users can create VMs from it.
 *
 * @param uuid @param filepath
 */
func (op Operations) UploadImage(uuid, filepath string) error {
	ctx := context.Background()

	path := fmt.Sprintf("/images/%s/file", uuid)

	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error: cannot open file: %s", err)
	}
	defer file.Close()

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error: Cannot read file %s", err)
	}

	req, err := op.client.NewUploadRequest(ctx, http.MethodPut, path, fileContents)

	if err != nil {
		return fmt.Errorf("error: Creating request %s", err)
	}

	err = op.client.Do(ctx, req, nil)

	return err
}

/*DeleteImage deletes a IMAGE
 * This operation submits a request to delete a IMAGE.
 *
 * @param uuid The uuid of the entity.
 * @return error if error exists
 */
func (op Operations) DeleteImage(uuid string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/images/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)

	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
}

/*GetImage gets a IMAGE
 * This operation gets a IMAGE.
 *
 * @param uuid The uuid of the entity.
 * @return *ImageIntentResponse
 */
func (op Operations) GetImage(uuid string) (*ImageIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/images/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	imageIntentResponse := new(ImageIntentResponse)

	if err != nil {
		return nil, err
	}

	return imageIntentResponse, op.client.Do(ctx, req, imageIntentResponse)
}

/*ListImage gets a list of IMAGEs This operation gets a list of IMAGEs, allowing for sorting and pagination. Note: Entities that have not
 * been created successfully are not listed.
 *
 * @param getEntitiesRequest @return *ImageListIntentResponse
 */
func (op Operations) ListImage(getEntitiesRequest *DSMetadata) (*ImageListIntentResponse, error) {
	ctx := context.TODO()
	path := "/images/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	imageListIntentResponse := new(ImageListIntentResponse)

	if err != nil {
		return nil, err
	}

	return imageListIntentResponse, op.client.Do(ctx, req, imageListIntentResponse)
}

/*UpdateImage updates a IMAGE
 * This operation submits a request to update a IMAGE based on the input parameters.
 *
 * @param uuid The uuid of the entity.
 * @param body
 * @return *ImageIntentResponse
 */
func (op Operations) UpdateImage(uuid string, body *ImageIntentInput) (*ImageIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/images/%s", uuid)
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	imageIntentResponse := new(ImageIntentResponse)

	if err != nil {
		return nil, err
	}

	return imageIntentResponse, op.client.Do(ctx, req, imageIntentResponse)
}

/*GetCluster gets a CLUSTER
 * This operation gets a CLUSTER.
 *
 * @param uuid The uuid of the entity.
 * @return *ImageIntentResponse
 */
func (op Operations) GetCluster(uuid string) (*ClusterIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/clusters/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	clusterIntentResponse := new(ClusterIntentResponse)

	if err != nil {
		return nil, err
	}

	return clusterIntentResponse, op.client.Do(ctx, req, clusterIntentResponse)
}

/*ListCluster gets a list of CLUSTERS This operation gets a list of CLUSTERS, allowing for sorting and pagination. Note: Entities that have
 * not been created successfully are not listed.
 *
 * @param getEntitiesRequest @return *ClusterListIntentResponse
 */
func (op Operations) ListCluster(getEntitiesRequest *DSMetadata) (*ClusterListIntentResponse, error) {
	ctx := context.TODO()
	path := "/clusters/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	clusterList := new(ClusterListIntentResponse)

	if err != nil {
		return nil, err
	}

	return clusterList, op.client.Do(ctx, req, clusterList)
}

/*UpdateImage updates a CLUSTER
 * This operation submits a request to update a CLUSTER based on the input parameters.
 *
 * @param uuid The uuid of the entity.
 * @param body
 * @return *ImageIntentResponse
 */
// func (op Operations) UpdateImage(uuid string, body *ImageIntentInput) (*ImageIntentResponse, error) {
// 	ctx := context.TODO()

// 	path := fmt.Sprintf("/images/%s", uuid)

// 	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	imageIntentResponse := new(ImageIntentResponse)

// 	err = op.client.Do(ctx, req, imageIntentResponse)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return imageIntentResponse, nil
// }

// CreateOrUpdateCategoryKey ...
func (op Operations) CreateOrUpdateCategoryKey(body *CategoryKey) (*CategoryKeyStatus, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s", utils.StringValue(body.Name))
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	categoryKeyResponse := new(CategoryKeyStatus)

	if err != nil {
		return nil, err
	}

	return categoryKeyResponse, op.client.Do(ctx, req, categoryKeyResponse)
}

/*ListCategories gets a list of Categories This operation gets a list of Categories, allowing for sorting and pagination. Note: Entities
 * that have not been created successfully are not listed.
 *
 * @param getEntitiesRequest @return *ImageListIntentResponse
 */
func (op Operations) ListCategories(getEntitiesRequest *CategoryListMetadata) (*CategoryKeyListResponse, error) {
	ctx := context.TODO()
	path := "/categories/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	categoryKeyListResponse := new(CategoryKeyListResponse)

	if err != nil {
		return nil, err
	}

	return categoryKeyListResponse, op.client.Do(ctx, req, categoryKeyListResponse)
}

/*DeleteCategoryKey Deletes a Category
 * This operation submits a request to delete a op.
 *
 * @param name The name of the entity.
 * @return error
 */
func (op Operations) DeleteCategoryKey(name string) error {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s", name)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

/*GetCategoryKey gets a Category
 * This operation gets a Category.
 *
 * @param name The name of the entity.
 * @return *CategoryKeyStatus
 */
func (op Operations) GetCategoryKey(name string) (*CategoryKeyStatus, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s", name)
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	categoryKeyStatusResponse := new(CategoryKeyStatus)

	if err != nil {
		return nil, err
	}

	return categoryKeyStatusResponse, op.client.Do(ctx, req, categoryKeyStatusResponse)
}

/*ListCategoryValues gets a list of Category values for a specific key This operation gets a list of Categories, allowing for sorting and
 * pagination. Note: Entities that have not been created successfully are not listed.
 *
 * @param name @param getEntitiesRequest @return *CategoryValueListResponse
 */
func (op Operations) ListCategoryValues(name string, getEntitiesRequest *CategoryListMetadata) (*CategoryValueListResponse, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/categories/%s/list", name)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	categoryValueListResponse := new(CategoryValueListResponse)

	if err != nil {
		return nil, err
	}

	return categoryValueListResponse, op.client.Do(ctx, req, categoryValueListResponse)
}

// CreateOrUpdateCategoryValue ...
func (op Operations) CreateOrUpdateCategoryValue(name string, body *CategoryValue) (*CategoryValueStatus, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s/%s", name, utils.StringValue(body.Value))
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	categoryValueResponse := new(CategoryValueStatus)

	if err != nil {
		return nil, err
	}

	return categoryValueResponse, op.client.Do(ctx, req, categoryValueResponse)
}

/*GetCategoryValue gets a Category Value
 * This operation gets a Category Value.
 *
 * @param name The name of the entity.
 * @params value the value of entity that belongs to category key
 * @return *CategoryValueStatus
 */
func (op Operations) GetCategoryValue(name string, value string) (*CategoryValueStatus, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s/%s", name, value)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	categoryValueStatusResponse := new(CategoryValueStatus)

	if err != nil {
		return nil, err
	}

	return categoryValueStatusResponse, op.client.Do(ctx, req, categoryValueStatusResponse)
}

/*DeleteCategoryValue Deletes a Category Value
 * This operation submits a request to delete a op.
 *
 * @param name The name of the entity.
 * @params value the value of entity that belongs to category key
 * @return error
 */
func (op Operations) DeleteCategoryValue(name string, value string) error {
	ctx := context.TODO()

	path := fmt.Sprintf("/categories/%s/%s", name, value)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

/*GetCategoryQuery gets list of entities attached to categories or policies in which categories are used as defined by the filter criteria.
 *
 * @param query Categories query input object.
 * @return *CategoryQueryResponse
 */
func (op Operations) GetCategoryQuery(query *CategoryQueryInput) (*CategoryQueryResponse, error) {
	ctx := context.TODO()

	path := "/categories/query"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, query)
	categoryQueryResponse := new(CategoryQueryResponse)

	if err != nil {
		return nil, err
	}

	return categoryQueryResponse, op.client.Do(ctx, req, categoryQueryResponse)
}

/*CreateNetworkSecurityRule Creates a Network security rule
 * This operation submits a request to create a Network security rule based on the input parameters.
 *
 * @param request
 * @return *NetworkSecurityRuleIntentResponse
 */
func (op Operations) CreateNetworkSecurityRule(request *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error) {
	ctx := context.TODO()

	networkSecurityRuleIntentResponse := new(NetworkSecurityRuleIntentResponse)
	req, err := op.client.NewRequest(ctx, http.MethodPost, "/network_security_rules", request)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleIntentResponse, op.client.Do(ctx, req, networkSecurityRuleIntentResponse)
}

/*DeleteNetworkSecurityRule Deletes a Network security rule
 * This operation submits a request to delete a Network security rule.
 *
 * @param uuid The uuid of the entity.
 * @return void
 */
func (op Operations) DeleteNetworkSecurityRule(uuid string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/network_security_rules/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)

	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
}

/*GetNetworkSecurityRule Gets a Network security rule
 * This operation gets a Network security rule.
 *
 * @param uuid The uuid of the entity.
 * @return *NetworkSecurityRuleIntentResponse
 */
func (op Operations) GetNetworkSecurityRule(uuid string) (*NetworkSecurityRuleIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/network_security_rules/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	networkSecurityRuleIntentResponse := new(NetworkSecurityRuleIntentResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleIntentResponse, op.client.Do(ctx, req, networkSecurityRuleIntentResponse)
}

/*ListNetworkSecurityRule Gets all network security rules This operation gets a list of Network security rules, allowing for sorting and
 * pagination. Note: Entities that have not been created successfully are not listed.
 *
 * @param getEntitiesRequest @return *NetworkSecurityRuleListIntentResponse
 */
func (op Operations) ListNetworkSecurityRule(getEntitiesRequest *DSMetadata) (*NetworkSecurityRuleListIntentResponse, error) {
	ctx := context.TODO()
	path := "/network_security_rules/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	networkSecurityRuleListIntentResponse := new(NetworkSecurityRuleListIntentResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleListIntentResponse, op.client.Do(ctx, req, networkSecurityRuleListIntentResponse)
}

/*UpdateNetworkSecurityRule Updates a Network security rule
 * This operation submits a request to update a Network security rule based on the input parameters.
 *
 * @param uuid The uuid of the entity.
 * @param body
 * @return void
 */
func (op Operations) UpdateNetworkSecurityRule(
	uuid string,
	body *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/network_security_rules/%s", uuid)
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	networkSecurityRuleIntentResponse := new(NetworkSecurityRuleIntentResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleIntentResponse, op.client.Do(ctx, req, networkSecurityRuleIntentResponse)
}

/*CreateVolumeGroup Creates a Volume group
 * This operation submits a request to create a Volume group based on the input parameters.
 *
 * @param request
 * @return *VolumeGroupResponse
 */
func (op Operations) CreateVolumeGroup(request *VolumeGroupInput) (*VolumeGroupResponse, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/volume_groups", request)
	networkSecurityRuleResponse := new(VolumeGroupResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleResponse, op.client.Do(ctx, req, networkSecurityRuleResponse)
}

/*DeleteVolumeGroup Deletes a Volume group
 * This operation submits a request to delete a Volume group.
 *
 * @param uuid The uuid of the entity.
 * @return void
 */
func (op Operations) DeleteVolumeGroup(uuid string) error {
	ctx := context.TODO()

	path := fmt.Sprintf("/volume_groups/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

/*GetVolumeGroup Gets a Volume group
 * This operation gets a Volume group.
 *
 * @param uuid The uuid of the entity.
 * @return *VolumeGroupResponse
 */
func (op Operations) GetVolumeGroup(uuid string) (*VolumeGroupResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/volume_groups/%s", uuid)
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	networkSecurityRuleResponse := new(VolumeGroupResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleResponse, op.client.Do(ctx, req, networkSecurityRuleResponse)
}

/*ListVolumeGroup Gets all network security rules This operation gets a list of Volume groups, allowing for sorting and pagination. Note:
 * Entities that have not been created successfully are not listed.
 *
 * @param getEntitiesRequest @return *VolumeGroupListResponse
 */
func (op Operations) ListVolumeGroup(getEntitiesRequest *DSMetadata) (*VolumeGroupListResponse, error) {
	ctx := context.TODO()
	path := "/volume_groups/list"
	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	networkSecurityRuleListResponse := new(VolumeGroupListResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleListResponse, op.client.Do(ctx, req, networkSecurityRuleListResponse)
}

/*UpdateVolumeGroup Updates a Volume group
 * This operation submits a request to update a Volume group based on the input parameters.
 *
 * @param uuid The uuid of the entity.
 * @param body
 * @return void
 */
func (op Operations) UpdateVolumeGroup(uuid string, body *VolumeGroupInput) (*VolumeGroupResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/volume_groups/%s", uuid)
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	networkSecurityRuleResponse := new(VolumeGroupResponse)

	if err != nil {
		return nil, err
	}

	return networkSecurityRuleResponse, op.client.Do(ctx, req, networkSecurityRuleResponse)
}

const itemsPerPage int64 = 100

func hasNext(ri *int64) bool {
	*ri -= itemsPerPage
	return *ri >= (0 - itemsPerPage)
}

// ListAllVM ...
func (op Operations) ListAllVM(filter string) (*VMListIntentResponse, error) {
	entities := make([]*VMIntentResource, 0)

	resp, err := op.ListVM(&DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("vm"),
		Length: utils.Int64Ptr(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListVM(&DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("vm"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
		}

		resp.Entities = entities
	}

	return resp, nil
}

// ListAllSubnet ...
func (op Operations) ListAllSubnet(filter string) (*SubnetListIntentResponse, error) {
	entities := make([]*SubnetIntentResponse, 0)

	resp, err := op.ListSubnet(&DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("subnet"),
		Length: utils.Int64Ptr(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListSubnet(&DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("subnet"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}

		resp.Entities = entities
	}

	return resp, nil
}

// ListAllNetworkSecurityRule ...
func (op Operations) ListAllNetworkSecurityRule(filter string) (*NetworkSecurityRuleListIntentResponse, error) {
	entities := make([]*NetworkSecurityRuleIntentResource, 0)

	resp, err := op.ListNetworkSecurityRule(&DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("network_security_rule"),
		Length: utils.Int64Ptr(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListNetworkSecurityRule(&DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("network_security_rule"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}

		resp.Entities = entities
	}

	return resp, nil
}

// ListAllImage ...
func (op Operations) ListAllImage(filter string) (*ImageListIntentResponse, error) {
	entities := make([]*ImageIntentResponse, 0)

	resp, err := op.ListImage(&DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("image"),
		Length: utils.Int64Ptr(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListImage(&DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("image"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}

		resp.Entities = entities
	}

	return resp, nil
}

// ListAllCluster ...
func (op Operations) ListAllCluster(filter string) (*ClusterListIntentResponse, error) {
	entities := make([]*ClusterIntentResponse, 0)

	resp, err := op.ListCluster(&DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("cluster"),
		Length: utils.Int64Ptr(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListCluster(&DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("cluster"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}

		resp.Entities = entities
	}

	return resp, nil
}

// ListAllCluster ...
func (op Operations) ListAllCategoryValues(categoryKeyName, filter string) (*CategoryValueListResponse, error) {
	entities := make([]*CategoryValueStatus, 0)

	resp, err := op.ListCategoryValues(categoryKeyName, &CategoryListMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("category"),
		Length: utils.Int64Ptr(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListCategoryValues(categoryKeyName, &CategoryListMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("category"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}

		resp.Entities = entities
	}

	return resp, nil
}

// GetTask ...
func (op Operations) GetTask(taskUUID string) (*TasksResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/tasks/%s", taskUUID)
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	tasksTesponse := new(TasksResponse)

	if err != nil {
		return nil, err
	}

	return tasksTesponse, op.client.Do(ctx, req, tasksTesponse)
}

// GetHost ...
func (op Operations) GetHost(hostUUID string) (*HostResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/hosts/%s", hostUUID)
	host := new(HostResponse)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return host, op.client.Do(ctx, req, host)
}

// ListHost ...
func (op Operations) ListHost(getEntitiesRequest *DSMetadata) (*HostListResponse, error) {
	ctx := context.TODO()
	path := "/hosts/list"

	hostList := new(HostListResponse)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return hostList, op.client.Do(ctx, req, hostList)
}

// ListAllHost ...
func (op Operations) ListAllHost() (*HostListResponse, error) {
	entities := make([]*HostResponse, 0)

	resp, err := op.ListHost(&DSMetadata{
		Kind:   utils.StringPtr("host"),
		Length: utils.Int64Ptr(itemsPerPage),
	})
	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListHost(&DSMetadata{
				Kind:   utils.StringPtr("cluster"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}

		resp.Entities = entities
	}

	return resp, nil
}

/*CreateProject creates a project
 * This operation submits a request to create a project based on the input parameters.
 *
 * @param request *Project
 * @return *Project
 */
func (op Operations) CreateProject(request *Project) (*Project, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/projects", request)
	if err != nil {
		return nil, err
	}

	projectResponse := new(Project)

	return projectResponse, op.client.Do(ctx, req, projectResponse)
}

/*GetProject This operation gets a project.
 *
 * @param uuid The prject uuid - string.
 * @return *Project
 */
func (op Operations) GetProject(projectUUID string) (*Project, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/projects/%s", projectUUID)
	project := new(Project)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return project, op.client.Do(ctx, req, project)
}

/*ListProject gets a list of projects.
 *
 * @param metadata allows create filters to get specific data - *DSMetadata.
 * @return *ProjectListResponse
 */
func (op Operations) ListProject(getEntitiesRequest *DSMetadata) (*ProjectListResponse, error) {
	ctx := context.TODO()
	path := "/projects/list"

	projectList := new(ProjectListResponse)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return projectList, op.client.Do(ctx, req, projectList)
}

/*ListAllProject gets a list of projects
 * This operation gets a list of Projects, allowing for sorting and pagination.
 * Note: Entities that have not been created successfully are not listed.
 * @return *ProjectListResponse
 */
func (op Operations) ListAllProject(filter string) (*ProjectListResponse, error) {
	entities := make([]*Project, 0)

	resp, err := op.ListProject(&DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("project"),
		Length: utils.Int64Ptr(itemsPerPage),
	})
	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListProject(&DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("project"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}

		resp.Entities = entities
	}

	return resp, nil
}

/*UpdateProject Updates a project
 * This operation submits a request to update a existing Project based on the input parameters
 * @param uuid The uuid of the entity - string.
 * @param body - *Project
 * @return *Project, error
 */
func (op Operations) UpdateProject(uuid string, body *Project) (*Project, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/projects/%s", uuid)
	projectInput := new(Project)

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}

	return projectInput, op.client.Do(ctx, req, projectInput)
}

/*DeleteProject Deletes a project
 * This operation submits a request to delete a existing Project.
 *
 * @param uuid The uuid of the entity.
 * @return void
 */
func (op Operations) DeleteProject(uuid string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/projects/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)
	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
}

/*CreateAccessControlPolicy creates a access policy
 * This operation submits a request to create a access policy based on the input parameters.
 *
 * @param request *Access Policy
 * @return *Access Policy
 */
func (op Operations) CreateAccessControlPolicy(request *AccessControlPolicy) (*AccessControlPolicy, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/access_control_policies", request)
	if err != nil {
		return nil, err
	}

	AccessControlPolicyResponse := new(AccessControlPolicy)

	return AccessControlPolicyResponse, op.client.Do(ctx, req, AccessControlPolicyResponse)
}

/*GetAccessControlPolicy This operation gets a AccessControlPolicy.
 *
 * @param uuid The access policy uuid - string.
 * @return *AccessControlPolicy
 */
func (op Operations) GetAccessControlPolicy(accessControlPolicyUUID string) (*AccessControlPolicy, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/access_control_policies/%s", accessControlPolicyUUID)
	AccessControlPolicy := new(AccessControlPolicy)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return AccessControlPolicy, op.client.Do(ctx, req, AccessControlPolicy)
}

/*ListAccessControlPolicy gets a list of AccessControlPolicys.
 *
 * @param metadata allows create filters to get specific data - *DSMetadata.
 * @return *AccessControlPolicyListResponse
 */
func (op Operations) ListAccessControlPolicy(getEntitiesRequest *DSMetadata) (*AccessControlPolicyListResponse, error) {
	ctx := context.TODO()
	path := "/access_control_policies/list"

	AccessControlPolicyList := new(AccessControlPolicyListResponse)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return AccessControlPolicyList, op.client.Do(ctx, req, AccessControlPolicyList)
}

/*ListAllAccessControlPolicy gets a list of AccessControlPolicys
 * This operation gets a list of AccessControlPolicys, allowing for sorting and pagination.
 * Note: Entities that have not been created successfully are not listed.
 * @return *AccessControlPolicyListResponse
 */
func (op Operations) ListAllAccessControlPolicy(filter string) (*AccessControlPolicyListResponse, error) {
	entities := make([]*AccessControlPolicy, 0)

	resp, err := op.ListAccessControlPolicy(&DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("access_control_policy"),
		Length: utils.Int64Ptr(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListAccessControlPolicy(&DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("access_control_policy"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}

		resp.Entities = entities
	}

	return resp, nil
}

/*UpdateAccessControlPolicy Updates a AccessControlPolicy
 * This operation submits a request to update a existing AccessControlPolicy based on the input parameters
 * @param uuid The uuid of the entity - string.
 * @param body - *AccessControlPolicy
 * @return *AccessControlPolicy, error
 */
func (op Operations) UpdateAccessControlPolicy(uuid string, body *AccessControlPolicy) (*AccessControlPolicy, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/access_control_policies/%s", uuid)
	AccessControlPolicyInput := new(AccessControlPolicy)

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}

	return AccessControlPolicyInput, op.client.Do(ctx, req, AccessControlPolicyInput)
}

/*DeleteAccessControlPolicy Deletes a AccessControlPolicy
 * This operation submits a request to delete a existing AccessControlPolicy.
 *
 * @param uuid The uuid of the entity.
 * @return void
 */
func (op Operations) DeleteAccessControlPolicy(uuid string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/access_control_policies/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)

	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
}

/*CreateRole creates a role
 * This operation submits a request to create a role based on the input parameters.
 *
 * @param request *Role
 * @return *Role
 */
func (op Operations) CreateRole(request *Role) (*Role, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/roles", request)
	if err != nil {
		return nil, err
	}

	RoleResponse := new(Role)

	return RoleResponse, op.client.Do(ctx, req, RoleResponse)
}

/*GetRole This operation gets a role.
 *
 * @param uuid The role uuid - string.
 * @return *Role
 */
func (op Operations) GetRole(roleUUID string) (*Role, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/roles/%s", roleUUID)
	Role := new(Role)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return Role, op.client.Do(ctx, req, Role)
}

/*ListRole gets a list of roles.
 *
 * @param metadata allows create filters to get specific data - *DSMetadata.
 * @return *RoleListResponse
 */
func (op Operations) ListRole(getEntitiesRequest *DSMetadata) (*RoleListResponse, error) {
	ctx := context.TODO()
	path := "/roles/list"

	RoleList := new(RoleListResponse)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return RoleList, op.client.Do(ctx, req, RoleList)
}

/*ListAllRole gets a list of Roles
 * This operation gets a list of Roles, allowing for sorting and pagination.
 * Note: Entities that have not been created successfully are not listed.
 * @return *RoleListResponse
 */
func (op Operations) ListAllRole(filter string) (*RoleListResponse, error) {
	entities := make([]*Role, 0)

	resp, err := op.ListRole(&DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("role"),
		Length: utils.Int64Ptr(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListRole(&DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("role"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}

		resp.Entities = entities
	}

	return resp, nil
}

/*UpdateRole Updates a role
 * This operation submits a request to update a existing role based on the input parameters
 * @param uuid The uuid of the entity - string.
 * @param body - *Role
 * @return *Role, error
 */
func (op Operations) UpdateRole(uuid string, body *Role) (*Role, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/roles/%s", uuid)
	RoleInput := new(Role)

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}

	return RoleInput, op.client.Do(ctx, req, RoleInput)
}

/*DeleteRole Deletes a role
 * This operation submits a request to delete a existing role.
 *
 * @param uuid The uuid of the entity.
 * @return void
 */
func (op Operations) DeleteRole(uuid string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/roles/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)

	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
}

/*CreateUser creates a User
 * This operation submits a request to create a userbased on the input parameters.
 *
 * @param request *VMIntentInput
 * @return *UserIntentResponse
 */
func (op Operations) CreateUser(request *UserIntentInput) (*UserIntentResponse, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/users", request)
	if err != nil {
		return nil, err
	}

	UserIntentResponse := new(UserIntentResponse)

	return UserIntentResponse, op.client.Do(ctx, req, UserIntentResponse)
}

/*GetUser This operation gets a User.
 *
 * @param uuid The user uuid - string.
 * @return *User
 */
func (op Operations) GetUser(userUUID string) (*UserIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/users/%s", userUUID)
	User := new(UserIntentResponse)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return User, op.client.Do(ctx, req, User)
}

/*UpdateUser Updates a User
 * This operation submits a request to update a existing User based on the input parameters
 * @param uuid The uuid of the entity - string.
 * @param body - *User
 * @return *User, error
 */
func (op Operations) UpdateUser(uuid string, body *UserIntentInput) (*UserIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/users/%s", uuid)
	UserInput := new(UserIntentResponse)

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}

	return UserInput, op.client.Do(ctx, req, UserInput)
}

/*DeleteUser Deletes a User
 * This operation submits a request to delete a existing User.
 *
 * @param uuid The uuid of the entity.
 * @return void
 */
func (op Operations) DeleteUser(uuid string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/users/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)

	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
}

/*ListUser gets a list of Users.
 *
 * @param metadata allows create filters to get specific data - *DSMetadata.
 * @return *UserListResponse
 */
func (op Operations) ListUser(getEntitiesRequest *DSMetadata) (*UserListResponse, error) {
	ctx := context.TODO()
	path := "/users/list"

	UserList := new(UserListResponse)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return UserList, op.client.Do(ctx, req, UserList)
}

// ListAllUser ...
func (op Operations) ListAllUser(filter string) (*UserListResponse, error) {
	entities := make([]*UserIntentResponse, 0)

	resp, err := op.ListUser(&DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("user"),
		Length: utils.Int64Ptr(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListUser(&DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("user"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
		}

		resp.Entities = entities
	}

	return resp, nil
}

/*GetUserGroup This operation gets a User.
 *
 * @param uuid The user uuid - string.
 * @return *User
 */
func (op Operations) GetUserGroup(userGroupUUID string) (*UserGroupIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/user_groups/%s", userGroupUUID)
	User := new(UserGroupIntentResponse)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return User, op.client.Do(ctx, req, User)
}

/*ListUserGroup gets a list of UserGroups.
 *
 * @param metadata allows create filters to get specific data - *DSMetadata.
 * @return *UserGroupListResponse
 */
func (op Operations) ListUserGroup(getEntitiesRequest *DSMetadata) (*UserGroupListResponse, error) {
	ctx := context.TODO()
	path := "/user_groups/list"

	UserGroupList := new(UserGroupListResponse)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return UserGroupList, op.client.Do(ctx, req, UserGroupList)
}

// ListAllUserGroup ...
func (op Operations) ListAllUserGroup(filter string) (*UserGroupListResponse, error) {
	entities := make([]*UserGroupIntentResponse, 0)

	resp, err := op.ListUserGroup(&DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("user_group"),
		Length: utils.Int64Ptr(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListUserGroup(&DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("user"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
		}

		resp.Entities = entities
	}

	return resp, nil
}

/*GePermission This operation gets a Permission.
 *
 * @param uuid The permission uuid - string.
 * @return *PermissionIntentResponse
 */
func (op Operations) GetPermission(permissionUUID string) (*PermissionIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/permissions/%s", permissionUUID)
	permission := new(PermissionIntentResponse)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return permission, op.client.Do(ctx, req, permission)
}

/*ListPermission gets a list of Permissions.
 *
 * @param metadata allows create filters to get specific data - *DSMetadata.
 * @return *PermissionListResponse
 */
func (op Operations) ListPermission(getEntitiesRequest *DSMetadata) (*PermissionListResponse, error) {
	ctx := context.TODO()
	path := "/permissions/list"

	PermissionList := new(PermissionListResponse)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return PermissionList, op.client.Do(ctx, req, PermissionList)
}

// ListAllPermission ...
func (op Operations) ListAllPermission(filter string) (*PermissionListResponse, error) {
	entities := make([]*PermissionIntentResponse, 0)

	resp, err := op.ListPermission(&DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("permission"),
		Length: utils.Int64Ptr(itemsPerPage),
	})

	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListPermission(&DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("permission"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
		}

		resp.Entities = entities
	}

	return resp, nil
}

//GetProtectionRule ...
func (op Operations) GetProtectionRule(uuid string) (*ProtectionRuleResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/protection_rules/%s", uuid)
	protectionRule := new(ProtectionRuleResponse)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return protectionRule, op.client.Do(ctx, req, protectionRule)
}

//ListProtectionRules ...
func (op Operations) ListProtectionRules(getEntitiesRequest *DSMetadata) (*ProtectionRulesListResponse, error) {
	ctx := context.TODO()
	path := "/protection_rules/list"

	list := new(ProtectionRulesListResponse)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return list, op.client.Do(ctx, req, list)
}

// ListAllProtectionRules ...
func (op Operations) ListAllProtectionRules(filter string) (*ProtectionRulesListResponse, error) {
	entities := make([]*ProtectionRuleResponse, 0)

	resp, err := op.ListProtectionRules(&DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("protection_rule"),
		Length: utils.Int64Ptr(itemsPerPage),
	})
	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListProtectionRules(&DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("protection_rule"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}

		resp.Entities = entities
	}

	return resp, nil
}

//CreateProtectionRule ...
func (op Operations) CreateProtectionRule(createRequest *ProtectionRuleInput) (*ProtectionRuleResponse, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/protection_rules", createRequest)
	protectionRuleResponse := new(ProtectionRuleResponse)

	if err != nil {
		return nil, err
	}

	return protectionRuleResponse, op.client.Do(ctx, req, protectionRuleResponse)
}

//UpdateProtectionRule ...
func (op Operations) UpdateProtectionRule(uuid string, body *ProtectionRuleInput) (*ProtectionRuleResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/protection_rules/%s", uuid)
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	protectionRuleResponse := new(ProtectionRuleResponse)

	if err != nil {
		return nil, err
	}

	return protectionRuleResponse, op.client.Do(ctx, req, protectionRuleResponse)
}

//DeleteProtectionRule ...
func (op Operations) DeleteProtectionRule(uuid string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/protection_rules/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)

	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
}

//GetRecoveryPlan ...
func (op Operations) GetRecoveryPlan(uuid string) (*RecoveryPlanResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/recovery_plans/%s", uuid)
	RecoveryPlan := new(RecoveryPlanResponse)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return RecoveryPlan, op.client.Do(ctx, req, RecoveryPlan)
}

//ListRecoveryPlans ...
func (op Operations) ListRecoveryPlans(getEntitiesRequest *DSMetadata) (*RecoveryPlanListResponse, error) {
	ctx := context.TODO()
	path := "/recovery_plans/list"

	list := new(RecoveryPlanListResponse)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return list, op.client.Do(ctx, req, list)
}

// ListAllRecoveryPlans ...
func (op Operations) ListAllRecoveryPlans(filter string) (*RecoveryPlanListResponse, error) {
	entities := make([]*RecoveryPlanResponse, 0)

	resp, err := op.ListRecoveryPlans(&DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("recovery_plan"),
		Length: utils.Int64Ptr(itemsPerPage),
	})
	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListRecoveryPlans(&DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("recovery_plan"),
				Length: utils.Int64Ptr(itemsPerPage),
				Offset: utils.Int64Ptr(offset),
			})

			if err != nil {
				return nil, err
			}

			entities = append(entities, resp.Entities...)

			offset += itemsPerPage
			log.Printf("[Debug] total=%d, remaining=%d, offset=%d len(entities)=%d\n", totalEntities, remaining, offset, len(entities))
		}

		resp.Entities = entities
	}

	return resp, nil
}

//CreateRecoveryPlan ...
func (op Operations) CreateRecoveryPlan(createRequest *RecoveryPlanInput) (*RecoveryPlanResponse, error) {
	ctx := context.TODO()

	req, err := op.client.NewRequest(ctx, http.MethodPost, "/recovery_plans", createRequest)
	RecoveryPlanResponse := new(RecoveryPlanResponse)

	if err != nil {
		return nil, err
	}

	return RecoveryPlanResponse, op.client.Do(ctx, req, RecoveryPlanResponse)
}

//UpdateRecoveryPlan ...
func (op Operations) UpdateRecoveryPlan(uuid string, body *RecoveryPlanInput) (*RecoveryPlanResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/recovery_plans/%s", uuid)
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, body)
	RecoveryPlanResponse := new(RecoveryPlanResponse)

	if err != nil {
		return nil, err
	}

	return RecoveryPlanResponse, op.client.Do(ctx, req, RecoveryPlanResponse)
}

//DeleteRecoveryPlan ...
func (op Operations) DeleteRecoveryPlan(uuid string) (*DeleteResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/recovery_plans/%s", uuid)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)

	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
}
