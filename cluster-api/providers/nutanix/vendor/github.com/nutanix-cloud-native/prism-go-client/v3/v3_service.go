package v3

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nutanix-cloud-native/prism-go-client"
	"github.com/nutanix-cloud-native/prism-go-client/internal"
	"github.com/nutanix-cloud-native/prism-go-client/utils"
	"github.com/nutanix-cloud-native/prism-go-client/v3/models"
)

// Operations ...
type Operations struct {
	client *internal.Client
}

// Service ...
type Service interface {
	CreateVM(ctx context.Context, createRequest *VMIntentInput) (*VMIntentResponse, error)
	DeleteVM(ctx context.Context, uuid string) (*DeleteResponse, error)
	GetVM(ctx context.Context, uuid string) (*VMIntentResponse, error)
	ListVM(ctx context.Context, getEntitiesRequest *DSMetadata) (*VMListIntentResponse, error)
	UpdateVM(ctx context.Context, uuid string, body *VMIntentInput) (*VMIntentResponse, error)
	CreateSubnet(ctx context.Context, createRequest *SubnetIntentInput) (*SubnetIntentResponse, error)
	DeleteSubnet(ctx context.Context, uuid string) (*DeleteResponse, error)
	GetSubnet(ctx context.Context, uuid string) (*SubnetIntentResponse, error)
	ListSubnet(ctx context.Context, getEntitiesRequest *DSMetadata) (*SubnetListIntentResponse, error)
	UpdateSubnet(ctx context.Context, uuid string, body *SubnetIntentInput) (*SubnetIntentResponse, error)
	CreateImage(ctx context.Context, createRequest *ImageIntentInput) (*ImageIntentResponse, error)
	DeleteImage(ctx context.Context, uuid string) (*DeleteResponse, error)
	GetImage(ctx context.Context, uuid string) (*ImageIntentResponse, error)
	ListImage(ctx context.Context, getEntitiesRequest *DSMetadata) (*ImageListIntentResponse, error)
	UpdateImage(ctx context.Context, uuid string, body *ImageIntentInput) (*ImageIntentResponse, error)
	UploadImage(ctx context.Context, uuid, filepath string) error
	CreateOrUpdateCategoryKey(ctx context.Context, body *CategoryKey) (*CategoryKeyStatus, error)
	ListCategories(ctx context.Context, getEntitiesRequest *CategoryListMetadata) (*CategoryKeyListResponse, error)
	DeleteCategoryKey(ctx context.Context, name string) error
	GetCategoryKey(ctx context.Context, name string) (*CategoryKeyStatus, error)
	ListCategoryValues(ctx context.Context, name string, getEntitiesRequest *CategoryListMetadata) (*CategoryValueListResponse, error)
	CreateOrUpdateCategoryValue(ctx context.Context, name string, body *CategoryValue) (*CategoryValueStatus, error)
	GetCategoryValue(ctx context.Context, name string, value string) (*CategoryValueStatus, error)
	DeleteCategoryValue(ctx context.Context, name string, value string) error
	GetCategoryQuery(ctx context.Context, query *CategoryQueryInput) (*CategoryQueryResponse, error)
	UpdateNetworkSecurityRule(ctx context.Context, uuid string, body *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error)
	ListNetworkSecurityRule(ctx context.Context, getEntitiesRequest *DSMetadata) (*NetworkSecurityRuleListIntentResponse, error)
	GetNetworkSecurityRule(ctx context.Context, uuid string) (*NetworkSecurityRuleIntentResponse, error)
	DeleteNetworkSecurityRule(ctx context.Context, uuid string) (*DeleteResponse, error)
	CreateNetworkSecurityRule(ctx context.Context, request *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error)
	ListCluster(ctx context.Context, getEntitiesRequest *DSMetadata) (*ClusterListIntentResponse, error)
	GetCluster(ctx context.Context, uuid string) (*ClusterIntentResponse, error)
	UpdateVolumeGroup(ctx context.Context, uuid string, body *VolumeGroupInput) (*VolumeGroupResponse, error)
	ListVolumeGroup(ctx context.Context, getEntitiesRequest *DSMetadata) (*VolumeGroupListResponse, error)
	GetVolumeGroup(ctx context.Context, uuid string) (*VolumeGroupResponse, error)
	DeleteVolumeGroup(ctx context.Context, uuid string) error
	CreateVolumeGroup(ctx context.Context, request *VolumeGroupInput) (*VolumeGroupResponse, error)
	ListAllVM(ctx context.Context, filter string) (*VMListIntentResponse, error)
	ListAllSubnet(ctx context.Context, filter string, clientSideFilters []*prismgoclient.AdditionalFilter) (*SubnetListIntentResponse, error)
	ListAllNetworkSecurityRule(ctx context.Context, filter string) (*NetworkSecurityRuleListIntentResponse, error)
	ListAllImage(ctx context.Context, filter string) (*ImageListIntentResponse, error)
	ListAllCluster(ctx context.Context, filter string) (*ClusterListIntentResponse, error)
	ListAllCategoryValues(ctx context.Context, categoryName, filter string) (*CategoryValueListResponse, error)
	GetTask(ctx context.Context, taskUUID string) (*TasksResponse, error)
	GetHost(ctx context.Context, taskUUID string) (*HostResponse, error)
	ListHost(ctx context.Context, getEntitiesRequest *DSMetadata) (*HostListResponse, error)
	ListAllHost(ctx context.Context) (*HostListResponse, error)
	CreateProject(ctx context.Context, request *Project) (*Project, error)
	GetProject(ctx context.Context, projectUUID string) (*Project, error)
	ListProject(ctx context.Context, getEntitiesRequest *DSMetadata) (*ProjectListResponse, error)
	ListAllProject(ctx context.Context, filter string) (*ProjectListResponse, error)
	UpdateProject(ctx context.Context, uuid string, body *Project) (*Project, error)
	DeleteProject(ctx context.Context, uuid string) (*DeleteResponse, error)
	CreateAccessControlPolicy(ctx context.Context, request *AccessControlPolicy) (*AccessControlPolicy, error)
	GetAccessControlPolicy(ctx context.Context, accessControlPolicyUUID string) (*AccessControlPolicy, error)
	ListAccessControlPolicy(ctx context.Context, getEntitiesRequest *DSMetadata) (*AccessControlPolicyListResponse, error)
	ListAllAccessControlPolicy(ctx context.Context, filter string) (*AccessControlPolicyListResponse, error)
	UpdateAccessControlPolicy(ctx context.Context, uuid string, body *AccessControlPolicy) (*AccessControlPolicy, error)
	DeleteAccessControlPolicy(ctx context.Context, uuid string) (*DeleteResponse, error)
	CreateRole(ctx context.Context, request *Role) (*Role, error)
	GetRole(ctx context.Context, uuid string) (*Role, error)
	ListRole(ctx context.Context, getEntitiesRequest *DSMetadata) (*RoleListResponse, error)
	ListAllRole(ctx context.Context, filter string) (*RoleListResponse, error)
	UpdateRole(ctx context.Context, uuid string, body *Role) (*Role, error)
	DeleteRole(ctx context.Context, uuid string) (*DeleteResponse, error)
	CreateUser(ctx context.Context, request *UserIntentInput) (*UserIntentResponse, error)
	GetUser(ctx context.Context, userUUID string) (*UserIntentResponse, error)
	UpdateUser(ctx context.Context, uuid string, body *UserIntentInput) (*UserIntentResponse, error)
	DeleteUser(ctx context.Context, uuid string) (*DeleteResponse, error)
	ListUser(ctx context.Context, getEntitiesRequest *DSMetadata) (*UserListResponse, error)
	ListAllUser(ctx context.Context, filter string) (*UserListResponse, error)
	GetCurrentLoggedInUser(ctx context.Context) (*UserIntentResponse, error)
	GetUserGroup(ctx context.Context, userUUID string) (*UserGroupIntentResponse, error)
	ListUserGroup(ctx context.Context, getEntitiesRequest *DSMetadata) (*UserGroupListResponse, error)
	ListAllUserGroup(ctx context.Context, filter string) (*UserGroupListResponse, error)
	GetPermission(ctx context.Context, permissionUUID string) (*PermissionIntentResponse, error)
	ListPermission(ctx context.Context, getEntitiesRequest *DSMetadata) (*PermissionListResponse, error)
	ListAllPermission(ctx context.Context, filter string) (*PermissionListResponse, error)
	GetProtectionRule(ctx context.Context, uuid string) (*ProtectionRuleResponse, error)
	ListProtectionRules(ctx context.Context, getEntitiesRequest *DSMetadata) (*ProtectionRulesListResponse, error)
	ListAllProtectionRules(ctx context.Context, filter string) (*ProtectionRulesListResponse, error)
	CreateProtectionRule(ctx context.Context, request *ProtectionRuleInput) (*ProtectionRuleResponse, error)
	UpdateProtectionRule(ctx context.Context, uuid string, body *ProtectionRuleInput) (*ProtectionRuleResponse, error)
	DeleteProtectionRule(ctx context.Context, uuid string) (*DeleteResponse, error)
	ProcessProtectionRule(ctx context.Context, uuid string) error
	GetRecoveryPlan(ctx context.Context, uuid string) (*RecoveryPlanResponse, error)
	ListRecoveryPlans(ctx context.Context, getEntitiesRequest *DSMetadata) (*RecoveryPlanListResponse, error)
	ListAllRecoveryPlans(ctx context.Context, filter string) (*RecoveryPlanListResponse, error)
	CreateRecoveryPlan(ctx context.Context, request *RecoveryPlanInput) (*RecoveryPlanResponse, error)
	UpdateRecoveryPlan(ctx context.Context, uuid string, body *RecoveryPlanInput) (*RecoveryPlanResponse, error)
	DeleteRecoveryPlan(ctx context.Context, uuid string) (*DeleteResponse, error)
	GetServiceGroup(ctx context.Context, uuid string) (*ServiceGroupResponse, error)
	ListAllServiceGroups(ctx context.Context, filter string) (*ServiceGroupListResponse, error)
	CreateServiceGroup(ctx context.Context, request *ServiceGroupInput) (*Reference, error)
	UpdateServiceGroup(ctx context.Context, uuid string, body *ServiceGroupInput) error
	DeleteServiceGroup(ctx context.Context, uuid string) error
	GetAddressGroup(ctx context.Context, uuid string) (*AddressGroupResponse, error)
	ListAddressGroups(ctx context.Context, getEntitiesRequest *DSMetadata) (*AddressGroupListResponse, error)
	ListAllAddressGroups(ctx context.Context, filter string) (*AddressGroupListResponse, error)
	DeleteAddressGroup(ctx context.Context, uuid string) error
	CreateAddressGroup(ctx context.Context, request *AddressGroupInput) (*Reference, error)
	UpdateAddressGroup(ctx context.Context, uuid string, body *AddressGroupInput) error
	GetRecoveryPlanJob(ctx context.Context, uuid string) (*RecoveryPlanJobIntentResponse, error)
	GetRecoveryPlanJobStatus(ctx context.Context, uuid string, status string) (*RecoveryPlanJobExecutionStatus, error)
	ListRecoveryPlanJobs(ctx context.Context, getEntitiesRequest *DSMetadata) (*RecoveryPlanJobListResponse, error)
	DeleteRecoveryPlanJob(ctx context.Context, uuid string) error
	CreateRecoveryPlanJob(ctx context.Context, request *RecoveryPlanJobIntentInput) (*RecoveryPlanJobResponse, error)
	PerformRecoveryPlanJobAction(ctx context.Context, uuid string, action string, request *RecoveryPlanJobActionRequest) (*RecoveryPlanJobResponse, error)
	GroupsGetEntities(ctx context.Context, request *GroupsGetEntitiesRequest) (*GroupsGetEntitiesResponse, error)
	GetAvailabilityZone(ctx context.Context, uuid string) (*AvailabilityZoneIntentResponse, error)
	GetPrismCentral(ctx context.Context) (*models.PrismCentral, error)
}

/*CreateVM Creates a VM
 * This operation submits a request to create a VM based on the input parameters.
 *
 * @param body
 * @return *VMIntentResponse
 */
func (op Operations) CreateVM(ctx context.Context, createRequest *VMIntentInput) (*VMIntentResponse, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/vms", createRequest)
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
func (op Operations) DeleteVM(ctx context.Context, uuid string) (*DeleteResponse, error) {
	path := fmt.Sprintf("/vms/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
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
func (op Operations) GetVM(ctx context.Context, uuid string) (*VMIntentResponse, error) {
	path := fmt.Sprintf("/vms/%s", uuid)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) ListVM(ctx context.Context, getEntitiesRequest *DSMetadata) (*VMListIntentResponse, error) {
	path := "/vms/list"

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
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
func (op Operations) UpdateVM(ctx context.Context, uuid string, body *VMIntentInput) (*VMIntentResponse, error) {
	path := fmt.Sprintf("/vms/%s", uuid)
	req, err := op.client.NewRequest(http.MethodPut, path, body)
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
func (op Operations) CreateSubnet(ctx context.Context, createRequest *SubnetIntentInput) (*SubnetIntentResponse, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/subnets", createRequest)
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
func (op Operations) DeleteSubnet(ctx context.Context, uuid string) (*DeleteResponse, error) {
	path := fmt.Sprintf("/subnets/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
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
func (op Operations) GetSubnet(ctx context.Context, uuid string) (*SubnetIntentResponse, error) {
	path := fmt.Sprintf("/subnets/%s", uuid)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) ListSubnet(ctx context.Context, getEntitiesRequest *DSMetadata) (*SubnetListIntentResponse, error) {
	path := "/subnets/list"

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
	subnetListIntentResponse := new(SubnetListIntentResponse)

	if err != nil {
		return nil, err
	}
	baseSearchPaths := []string{"metadata", "status", "status.resources"}

	return subnetListIntentResponse, op.client.DoWithFilters(ctx, req, subnetListIntentResponse, getEntitiesRequest.ClientSideFilters, baseSearchPaths)
}

/*UpdateSubnet Updates a subnet
 * This operation submits a request to update a subnet based on the input parameters.
 *
 * @param uuid The uuid of the entity.
 * @param body
 * @return *SubnetIntentResponse
 */
func (op Operations) UpdateSubnet(ctx context.Context, uuid string, body *SubnetIntentInput) (*SubnetIntentResponse, error) {
	path := fmt.Sprintf("/subnets/%s", uuid)
	req, err := op.client.NewRequest(http.MethodPut, path, body)
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
func (op Operations) CreateImage(ctx context.Context, body *ImageIntentInput) (*ImageIntentResponse, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/images", body)
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
func (op Operations) UploadImage(ctx context.Context, uuid, filepath string) error {
	path := fmt.Sprintf("/images/%s/file", uuid)

	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error: cannot open file: %s", err)
	}
	defer file.Close()

	req, err := op.client.NewUploadRequest(http.MethodPut, path, file)
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
func (op Operations) DeleteImage(ctx context.Context, uuid string) (*DeleteResponse, error) {
	path := fmt.Sprintf("/images/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
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
func (op Operations) GetImage(ctx context.Context, uuid string) (*ImageIntentResponse, error) {
	path := fmt.Sprintf("/images/%s", uuid)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) ListImage(ctx context.Context, getEntitiesRequest *DSMetadata) (*ImageListIntentResponse, error) {
	path := "/images/list"

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
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
func (op Operations) UpdateImage(ctx context.Context, uuid string, body *ImageIntentInput) (*ImageIntentResponse, error) {
	path := fmt.Sprintf("/images/%s", uuid)
	req, err := op.client.NewRequest(http.MethodPut, path, body)
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
func (op Operations) GetCluster(ctx context.Context, uuid string) (*ClusterIntentResponse, error) {
	path := fmt.Sprintf("/clusters/%s", uuid)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) ListCluster(ctx context.Context, getEntitiesRequest *DSMetadata) (*ClusterListIntentResponse, error) {
	path := "/clusters/list"

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
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

// 	req, err := op.internal.NewRequest(ctx, http.MethodPut, path, body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	imageIntentResponse := new(ImageIntentResponse)

// 	err = op.internal.Do(ctx, req, imageIntentResponse)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return imageIntentResponse, nil
// }

// CreateOrUpdateCategoryKey ...
func (op Operations) CreateOrUpdateCategoryKey(ctx context.Context, body *CategoryKey) (*CategoryKeyStatus, error) {
	path := fmt.Sprintf("/categories/%s", utils.StringValue(body.Name))
	req, err := op.client.NewRequest(http.MethodPut, path, body)
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
func (op Operations) ListCategories(ctx context.Context, getEntitiesRequest *CategoryListMetadata) (*CategoryKeyListResponse, error) {
	path := "/categories/list"

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
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
func (op Operations) DeleteCategoryKey(ctx context.Context, name string) error {
	path := fmt.Sprintf("/categories/%s", name)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
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
func (op Operations) GetCategoryKey(ctx context.Context, name string) (*CategoryKeyStatus, error) {
	path := fmt.Sprintf("/categories/%s", name)
	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) ListCategoryValues(ctx context.Context, name string, getEntitiesRequest *CategoryListMetadata) (*CategoryValueListResponse, error) {
	path := fmt.Sprintf("/categories/%s/list", name)

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
	categoryValueListResponse := new(CategoryValueListResponse)

	if err != nil {
		return nil, err
	}

	return categoryValueListResponse, op.client.Do(ctx, req, categoryValueListResponse)
}

// CreateOrUpdateCategoryValue ...
func (op Operations) CreateOrUpdateCategoryValue(ctx context.Context, name string, body *CategoryValue) (*CategoryValueStatus, error) {
	path := fmt.Sprintf("/categories/%s/%s", name, utils.StringValue(body.Value))
	req, err := op.client.NewRequest(http.MethodPut, path, body)
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
func (op Operations) GetCategoryValue(ctx context.Context, name string, value string) (*CategoryValueStatus, error) {
	path := fmt.Sprintf("/categories/%s/%s", name, value)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) DeleteCategoryValue(ctx context.Context, name string, value string) error {
	path := fmt.Sprintf("/categories/%s/%s", name, value)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
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
func (op Operations) GetCategoryQuery(ctx context.Context, query *CategoryQueryInput) (*CategoryQueryResponse, error) {
	path := "/category/query"

	req, err := op.client.NewRequest(http.MethodPost, path, query)
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
func (op Operations) CreateNetworkSecurityRule(ctx context.Context, request *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error) {
	networkSecurityRuleIntentResponse := new(NetworkSecurityRuleIntentResponse)
	req, err := op.client.NewRequest(http.MethodPost, "/network_security_rules", request)
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
func (op Operations) DeleteNetworkSecurityRule(ctx context.Context, uuid string) (*DeleteResponse, error) {
	path := fmt.Sprintf("/network_security_rules/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
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
func (op Operations) GetNetworkSecurityRule(ctx context.Context, uuid string) (*NetworkSecurityRuleIntentResponse, error) {
	path := fmt.Sprintf("/network_security_rules/%s", uuid)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) ListNetworkSecurityRule(ctx context.Context, getEntitiesRequest *DSMetadata) (*NetworkSecurityRuleListIntentResponse, error) {
	path := "/network_security_rules/list"

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
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
func (op Operations) UpdateNetworkSecurityRule(ctx context.Context, uuid string, body *NetworkSecurityRuleIntentInput) (*NetworkSecurityRuleIntentResponse, error) {
	path := fmt.Sprintf("/network_security_rules/%s", uuid)
	req, err := op.client.NewRequest(http.MethodPut, path, body)
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
func (op Operations) CreateVolumeGroup(ctx context.Context, request *VolumeGroupInput) (*VolumeGroupResponse, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/volume_groups", request)
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
func (op Operations) DeleteVolumeGroup(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("/volume_groups/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
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
func (op Operations) GetVolumeGroup(ctx context.Context, uuid string) (*VolumeGroupResponse, error) {
	path := fmt.Sprintf("/volume_groups/%s", uuid)
	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) ListVolumeGroup(ctx context.Context, getEntitiesRequest *DSMetadata) (*VolumeGroupListResponse, error) {
	path := "/volume_groups/list"
	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
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
func (op Operations) UpdateVolumeGroup(ctx context.Context, uuid string, body *VolumeGroupInput) (*VolumeGroupResponse, error) {
	path := fmt.Sprintf("/volume_groups/%s", uuid)
	req, err := op.client.NewRequest(http.MethodPut, path, body)
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
func (op Operations) ListAllVM(ctx context.Context, filter string) (*VMListIntentResponse, error) {
	entities := make([]*VMIntentResource, 0)

	resp, err := op.ListVM(ctx, &DSMetadata{
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
			resp, err = op.ListVM(ctx, &DSMetadata{
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
func (op Operations) ListAllSubnet(ctx context.Context, filter string, clientSideFilters []*prismgoclient.AdditionalFilter) (*SubnetListIntentResponse, error) {
	entities := make([]*SubnetIntentResponse, 0)

	resp, err := op.ListSubnet(ctx, &DSMetadata{
		Filter:            &filter,
		Kind:              utils.StringPtr("subnet"),
		Length:            utils.Int64Ptr(itemsPerPage),
		ClientSideFilters: clientSideFilters,
	})
	if err != nil {
		return nil, err
	}

	totalEntities := utils.Int64Value(resp.Metadata.TotalMatches)
	remaining := totalEntities
	offset := utils.Int64Value(resp.Metadata.Offset)

	if totalEntities > itemsPerPage {
		for hasNext(&remaining) {
			resp, err = op.ListSubnet(ctx, &DSMetadata{
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
func (op Operations) ListAllNetworkSecurityRule(ctx context.Context, filter string) (*NetworkSecurityRuleListIntentResponse, error) {
	entities := make([]*NetworkSecurityRuleIntentResource, 0)

	resp, err := op.ListNetworkSecurityRule(ctx, &DSMetadata{
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
			resp, err = op.ListNetworkSecurityRule(ctx, &DSMetadata{
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
func (op Operations) ListAllImage(ctx context.Context, filter string) (*ImageListIntentResponse, error) {
	entities := make([]*ImageIntentResponse, 0)

	resp, err := op.ListImage(ctx, &DSMetadata{
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
			resp, err = op.ListImage(ctx, &DSMetadata{
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
func (op Operations) ListAllCluster(ctx context.Context, filter string) (*ClusterListIntentResponse, error) {
	entities := make([]*ClusterIntentResponse, 0)

	resp, err := op.ListCluster(ctx, &DSMetadata{
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
			resp, err = op.ListCluster(ctx, &DSMetadata{
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
func (op Operations) ListAllCategoryValues(ctx context.Context, categoryName, filter string) (*CategoryValueListResponse, error) {
	entities := make([]*CategoryValueStatus, 0)

	resp, err := op.ListCategoryValues(ctx, categoryName, &CategoryListMetadata{
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
			resp, err = op.ListCategoryValues(ctx, categoryName, &CategoryListMetadata{
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
func (op Operations) GetTask(ctx context.Context, taskUUID string) (*TasksResponse, error) {
	path := fmt.Sprintf("/tasks/%s", taskUUID)
	req, err := op.client.NewRequest(http.MethodGet, path, nil)
	tasksTesponse := new(TasksResponse)

	if err != nil {
		return nil, err
	}

	return tasksTesponse, op.client.Do(ctx, req, tasksTesponse)
}

// GetHost ...
func (op Operations) GetHost(ctx context.Context, hostUUID string) (*HostResponse, error) {
	path := fmt.Sprintf("/hosts/%s", hostUUID)
	host := new(HostResponse)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return host, op.client.Do(ctx, req, host)
}

// ListHost ...
func (op Operations) ListHost(ctx context.Context, getEntitiesRequest *DSMetadata) (*HostListResponse, error) {
	path := "/hosts/list"

	hostList := new(HostListResponse)

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return hostList, op.client.Do(ctx, req, hostList)
}

// ListAllHost ...
func (op Operations) ListAllHost(ctx context.Context) (*HostListResponse, error) {
	entities := make([]*HostResponse, 0)

	resp, err := op.ListHost(ctx, &DSMetadata{
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
			resp, err = op.ListHost(ctx, &DSMetadata{
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
func (op Operations) CreateProject(ctx context.Context, request *Project) (*Project, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/projects", request)
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
func (op Operations) GetProject(ctx context.Context, projectUUID string) (*Project, error) {
	path := fmt.Sprintf("/projects/%s", projectUUID)
	project := new(Project)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) ListProject(ctx context.Context, getEntitiesRequest *DSMetadata) (*ProjectListResponse, error) {
	path := "/projects/list"

	projectList := new(ProjectListResponse)

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
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
func (op Operations) ListAllProject(ctx context.Context, filter string) (*ProjectListResponse, error) {
	entities := make([]*Project, 0)

	resp, err := op.ListProject(ctx, &DSMetadata{
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
			resp, err = op.ListProject(ctx, &DSMetadata{
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
func (op Operations) UpdateProject(ctx context.Context, uuid string, body *Project) (*Project, error) {
	path := fmt.Sprintf("/projects/%s", uuid)
	projectInput := new(Project)

	req, err := op.client.NewRequest(http.MethodPut, path, body)
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
func (op Operations) DeleteProject(ctx context.Context, uuid string) (*DeleteResponse, error) {
	path := fmt.Sprintf("/projects/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
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
func (op Operations) CreateAccessControlPolicy(ctx context.Context, request *AccessControlPolicy) (*AccessControlPolicy, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/access_control_policies", request)
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
func (op Operations) GetAccessControlPolicy(ctx context.Context, accessControlPolicyUUID string) (*AccessControlPolicy, error) {
	path := fmt.Sprintf("/access_control_policies/%s", accessControlPolicyUUID)
	AccessControlPolicy := new(AccessControlPolicy)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) ListAccessControlPolicy(ctx context.Context, getEntitiesRequest *DSMetadata) (*AccessControlPolicyListResponse, error) {
	path := "/access_control_policies/list"

	AccessControlPolicyList := new(AccessControlPolicyListResponse)

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
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
func (op Operations) ListAllAccessControlPolicy(ctx context.Context, filter string) (*AccessControlPolicyListResponse, error) {
	entities := make([]*AccessControlPolicy, 0)

	resp, err := op.ListAccessControlPolicy(ctx, &DSMetadata{
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
			resp, err = op.ListAccessControlPolicy(ctx, &DSMetadata{
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
func (op Operations) UpdateAccessControlPolicy(ctx context.Context, uuid string, body *AccessControlPolicy) (*AccessControlPolicy, error) {
	path := fmt.Sprintf("/access_control_policies/%s", uuid)
	AccessControlPolicyInput := new(AccessControlPolicy)

	req, err := op.client.NewRequest(http.MethodPut, path, body)
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
func (op Operations) DeleteAccessControlPolicy(ctx context.Context, uuid string) (*DeleteResponse, error) {
	path := fmt.Sprintf("/access_control_policies/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
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
func (op Operations) CreateRole(ctx context.Context, request *Role) (*Role, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/roles", request)
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
func (op Operations) GetRole(ctx context.Context, roleUUID string) (*Role, error) {
	path := fmt.Sprintf("/roles/%s", roleUUID)
	Role := new(Role)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) ListRole(ctx context.Context, getEntitiesRequest *DSMetadata) (*RoleListResponse, error) {
	path := "/roles/list"

	RoleList := new(RoleListResponse)

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
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
func (op Operations) ListAllRole(ctx context.Context, filter string) (*RoleListResponse, error) {
	entities := make([]*Role, 0)

	resp, err := op.ListRole(ctx, &DSMetadata{
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
			resp, err = op.ListRole(ctx, &DSMetadata{
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
func (op Operations) UpdateRole(ctx context.Context, uuid string, body *Role) (*Role, error) {
	path := fmt.Sprintf("/roles/%s", uuid)
	RoleInput := new(Role)

	req, err := op.client.NewRequest(http.MethodPut, path, body)
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
func (op Operations) DeleteRole(ctx context.Context, uuid string) (*DeleteResponse, error) {
	path := fmt.Sprintf("/roles/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
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
func (op Operations) CreateUser(ctx context.Context, request *UserIntentInput) (*UserIntentResponse, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/users", request)
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
func (op Operations) GetUser(ctx context.Context, userUUID string) (*UserIntentResponse, error) {
	path := fmt.Sprintf("/users/%s", userUUID)
	User := new(UserIntentResponse)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) UpdateUser(ctx context.Context, uuid string, body *UserIntentInput) (*UserIntentResponse, error) {
	path := fmt.Sprintf("/users/%s", uuid)
	UserInput := new(UserIntentResponse)

	req, err := op.client.NewRequest(http.MethodPut, path, body)
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
func (op Operations) DeleteUser(ctx context.Context, uuid string) (*DeleteResponse, error) {
	path := fmt.Sprintf("/users/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
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
func (op Operations) ListUser(ctx context.Context, getEntitiesRequest *DSMetadata) (*UserListResponse, error) {
	path := "/users/list"

	UserList := new(UserListResponse)

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return UserList, op.client.Do(ctx, req, UserList)
}

// ListAllUser ...
func (op Operations) ListAllUser(ctx context.Context, filter string) (*UserListResponse, error) {
	entities := make([]*UserIntentResponse, 0)

	resp, err := op.ListUser(ctx, &DSMetadata{
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
			resp, err = op.ListUser(ctx, &DSMetadata{
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

/*GetCurrentLoggedInUser This operation gets the user info for the currently logged in User.
 *
 * @param context
 * @return *User
 */
func (op Operations) GetCurrentLoggedInUser(ctx context.Context) (*UserIntentResponse, error) {
	path := "/users/me"
	req, err := op.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	user := new(UserIntentResponse)
	if err := op.client.Do(ctx, req, user); err != nil {
		return nil, err
	}
	return user, nil
}

/*GetUserGroup This operation gets a User.
 *
 * @param uuid The user uuid - string.
 * @return *User
 */
func (op Operations) GetUserGroup(ctx context.Context, userGroupUUID string) (*UserGroupIntentResponse, error) {
	path := fmt.Sprintf("/user_groups/%s", userGroupUUID)
	User := new(UserGroupIntentResponse)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) ListUserGroup(ctx context.Context, getEntitiesRequest *DSMetadata) (*UserGroupListResponse, error) {
	path := "/user_groups/list"

	UserGroupList := new(UserGroupListResponse)

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return UserGroupList, op.client.Do(ctx, req, UserGroupList)
}

// ListAllUserGroup ...
func (op Operations) ListAllUserGroup(ctx context.Context, filter string) (*UserGroupListResponse, error) {
	entities := make([]*UserGroupIntentResponse, 0)

	resp, err := op.ListUserGroup(ctx, &DSMetadata{
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
			resp, err = op.ListUserGroup(ctx, &DSMetadata{
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
func (op Operations) GetPermission(ctx context.Context, permissionUUID string) (*PermissionIntentResponse, error) {
	path := fmt.Sprintf("/permissions/%s", permissionUUID)
	permission := new(PermissionIntentResponse)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
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
func (op Operations) ListPermission(ctx context.Context, getEntitiesRequest *DSMetadata) (*PermissionListResponse, error) {
	path := "/permissions/list"

	PermissionList := new(PermissionListResponse)

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return PermissionList, op.client.Do(ctx, req, PermissionList)
}

// ListAllPermission ...
func (op Operations) ListAllPermission(ctx context.Context, filter string) (*PermissionListResponse, error) {
	entities := make([]*PermissionIntentResponse, 0)

	resp, err := op.ListPermission(ctx, &DSMetadata{
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
			resp, err = op.ListPermission(ctx, &DSMetadata{
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

// GetProtectionRule ...
func (op Operations) GetProtectionRule(ctx context.Context, uuid string) (*ProtectionRuleResponse, error) {
	path := fmt.Sprintf("/protection_rules/%s", uuid)
	protectionRule := new(ProtectionRuleResponse)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return protectionRule, op.client.Do(ctx, req, protectionRule)
}

// ListProtectionRules ...
func (op Operations) ListProtectionRules(ctx context.Context, getEntitiesRequest *DSMetadata) (*ProtectionRulesListResponse, error) {
	path := "/protection_rules/list"

	list := new(ProtectionRulesListResponse)

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return list, op.client.Do(ctx, req, list)
}

// ListAllProtectionRules ...
func (op Operations) ListAllProtectionRules(ctx context.Context, filter string) (*ProtectionRulesListResponse, error) {
	entities := make([]*ProtectionRuleResponse, 0)

	resp, err := op.ListProtectionRules(ctx, &DSMetadata{
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
			resp, err = op.ListProtectionRules(ctx, &DSMetadata{
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

// CreateProtectionRule ...
func (op Operations) CreateProtectionRule(ctx context.Context, createRequest *ProtectionRuleInput) (*ProtectionRuleResponse, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/protection_rules", createRequest)
	protectionRuleResponse := new(ProtectionRuleResponse)

	if err != nil {
		return nil, err
	}

	return protectionRuleResponse, op.client.Do(ctx, req, protectionRuleResponse)
}

// UpdateProtectionRule ...
func (op Operations) UpdateProtectionRule(ctx context.Context, uuid string, body *ProtectionRuleInput) (*ProtectionRuleResponse, error) {
	path := fmt.Sprintf("/protection_rules/%s", uuid)
	req, err := op.client.NewRequest(http.MethodPut, path, body)
	protectionRuleResponse := new(ProtectionRuleResponse)

	if err != nil {
		return nil, err
	}

	return protectionRuleResponse, op.client.Do(ctx, req, protectionRuleResponse)
}

// DeleteProtectionRule ...
func (op Operations) DeleteProtectionRule(ctx context.Context, uuid string) (*DeleteResponse, error) {
	path := fmt.Sprintf("/protection_rules/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)

	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
}

/*ProcessProtectionRule triggers the evaluation of a processing rule
 * immediately.
 *
 * @param uuid is the uuid of the protection rule to process.
 */
func (op Operations) ProcessProtectionRule(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("/protection_rules/%s/process", uuid)

	req, err := op.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

// GetRecoveryPlan ...
func (op Operations) GetRecoveryPlan(ctx context.Context, uuid string) (*RecoveryPlanResponse, error) {
	path := fmt.Sprintf("/recovery_plans/%s", uuid)
	RecoveryPlan := new(RecoveryPlanResponse)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return RecoveryPlan, op.client.Do(ctx, req, RecoveryPlan)
}

// ListRecoveryPlans ...
func (op Operations) ListRecoveryPlans(ctx context.Context, getEntitiesRequest *DSMetadata) (*RecoveryPlanListResponse, error) {
	path := "/recovery_plans/list"

	list := new(RecoveryPlanListResponse)

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return list, op.client.Do(ctx, req, list)
}

// ListAllRecoveryPlans ...
func (op Operations) ListAllRecoveryPlans(ctx context.Context, filter string) (*RecoveryPlanListResponse, error) {
	entities := make([]*RecoveryPlanResponse, 0)

	resp, err := op.ListRecoveryPlans(ctx, &DSMetadata{
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
			resp, err = op.ListRecoveryPlans(ctx, &DSMetadata{
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

// CreateRecoveryPlan ...
func (op Operations) CreateRecoveryPlan(ctx context.Context, createRequest *RecoveryPlanInput) (*RecoveryPlanResponse, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/recovery_plans", createRequest)
	RecoveryPlanResponse := new(RecoveryPlanResponse)

	if err != nil {
		return nil, err
	}

	return RecoveryPlanResponse, op.client.Do(ctx, req, RecoveryPlanResponse)
}

// UpdateRecoveryPlan ...
func (op Operations) UpdateRecoveryPlan(ctx context.Context, uuid string, body *RecoveryPlanInput) (*RecoveryPlanResponse, error) {
	path := fmt.Sprintf("/recovery_plans/%s", uuid)
	req, err := op.client.NewRequest(http.MethodPut, path, body)
	RecoveryPlanResponse := new(RecoveryPlanResponse)

	if err != nil {
		return nil, err
	}

	return RecoveryPlanResponse, op.client.Do(ctx, req, RecoveryPlanResponse)
}

// DeleteRecoveryPlan ...
func (op Operations) DeleteRecoveryPlan(ctx context.Context, uuid string) (*DeleteResponse, error) {
	path := fmt.Sprintf("/recovery_plans/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
	deleteResponse := new(DeleteResponse)

	if err != nil {
		return nil, err
	}

	return deleteResponse, op.client.Do(ctx, req, deleteResponse)
}

func (op Operations) GetServiceGroup(ctx context.Context, uuid string) (*ServiceGroupResponse, error) {
	path := fmt.Sprintf("/service_groups/%s", uuid)
	ServiceGroup := new(ServiceGroupResponse)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return ServiceGroup, op.client.Do(ctx, req, ServiceGroup)
}

func (op Operations) CreateServiceGroup(ctx context.Context, request *ServiceGroupInput) (*Reference, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/service_groups", request)
	ServiceGroup := new(Reference)

	if err != nil {
		return nil, err
	}

	return ServiceGroup, op.client.Do(ctx, req, ServiceGroup)
}

func (op Operations) DeleteServiceGroup(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("/service_groups/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

func (op Operations) ListAllServiceGroups(ctx context.Context, filter string) (*ServiceGroupListResponse, error) {
	entities := make([]*ServiceGroupListEntry, 0)

	resp, err := op.listServiceGroups(ctx, &DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("service_group"),
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
			resp, err = op.listServiceGroups(ctx, &DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("service_group"),
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

func (op Operations) listServiceGroups(ctx context.Context, getEntitiesRequest *DSMetadata) (*ServiceGroupListResponse, error) {
	path := "/service_groups/list"

	list := new(ServiceGroupListResponse)

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return list, op.client.Do(ctx, req, list)
}

func (op Operations) UpdateServiceGroup(ctx context.Context, uuid string, body *ServiceGroupInput) error {
	path := fmt.Sprintf("/service_groups/%s", uuid)
	req, err := op.client.NewRequest(http.MethodPut, path, body)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

func (op Operations) GetAddressGroup(ctx context.Context, uuid string) (*AddressGroupResponse, error) {
	path := fmt.Sprintf("/address_groups/%s", uuid)
	AddressGroup := new(AddressGroupResponse)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return AddressGroup, op.client.Do(ctx, req, AddressGroup)
}

func (op Operations) ListAllAddressGroups(ctx context.Context, filter string) (*AddressGroupListResponse, error) {
	entities := make([]*AddressGroupListEntry, 0)

	resp, err := op.ListAddressGroups(ctx, &DSMetadata{
		Filter: &filter,
		Kind:   utils.StringPtr("address_group"),
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
			resp, err = op.ListAddressGroups(ctx, &DSMetadata{
				Filter: &filter,
				Kind:   utils.StringPtr("address_group"),
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

func (op Operations) ListAddressGroups(ctx context.Context, getEntitiesRequest *DSMetadata) (*AddressGroupListResponse, error) {
	path := "/address_groups/list"

	list := new(AddressGroupListResponse)

	req, err := op.client.NewRequest(http.MethodPost, path, getEntitiesRequest)
	if err != nil {
		return nil, err
	}

	return list, op.client.Do(ctx, req, list)
}

func (op Operations) DeleteAddressGroup(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("/address_groups/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

func (op Operations) CreateAddressGroup(ctx context.Context, request *AddressGroupInput) (*Reference, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/address_groups", request)
	AddressGroup := new(Reference)

	if err != nil {
		return nil, err
	}

	return AddressGroup, op.client.Do(ctx, req, AddressGroup)
}

func (op Operations) UpdateAddressGroup(ctx context.Context, uuid string, body *AddressGroupInput) error {
	path := fmt.Sprintf("/address_groups/%s", uuid)
	req, err := op.client.NewRequest(http.MethodPut, path, body)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

/*Creates a recovery plan job.
 * This operation creates a new recovery plan job based on the inputs in the 'request'.
 *
 * @param request Pointer to a specification of type RecoveryPlanJobIntentInput.
 */
func (op Operations) CreateRecoveryPlanJob(ctx context.Context, request *RecoveryPlanJobIntentInput) (*RecoveryPlanJobResponse, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/recovery_plan_jobs", request)
	response := new(RecoveryPlanJobResponse)

	if err != nil {
		return nil, err
	}

	return response, op.client.Do(ctx, req, response)
}

/*Deletes a recovery plan job.
 * This operation deletes the new recovery plan job identified by 'uuid'.
 *
 * @param uuid UUID of the recovery plan job to be deleted.
 */
func (op Operations) DeleteRecoveryPlanJob(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("/recovery_plan_jobs/%s", uuid)

	req, err := op.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return op.client.Do(ctx, req, nil)
}

/*Perform a recovery plan job action.
 * This operation initiates the 'action' on the recovery plan job identified
 * by 'uuid' and governed by the specification in 'request'.
 *
 * @param uuid UUID of the recovery plan job.
 * @param action one of {'cleanup', 'rerun'}.
 * @param request pointer to the specification of type RecoveryPlanJobActionRequest.
 */
func (op Operations) PerformRecoveryPlanJobAction(ctx context.Context, uuid string, action string,
	request *RecoveryPlanJobActionRequest,
) (*RecoveryPlanJobResponse, error) {
	path := fmt.Sprintf("/recovery_plan_jobs/%s/%s", uuid, action)
	response := new(RecoveryPlanJobResponse)

	req, err := op.client.NewRequest(http.MethodPost, path, request)
	if err != nil {
		return nil, err
	}

	return response, op.client.Do(ctx, req, response)
}

/*Get a recovery plan job.
 * This operation gets the recovery plan job identified by 'uuid'.
 *
 * @param uuid UUID of the recovery plan job.
 */
func (op Operations) GetRecoveryPlanJob(ctx context.Context, uuid string) (*RecoveryPlanJobIntentResponse, error) {
	path := fmt.Sprintf("/recovery_plan_jobs/%s", uuid)
	response := new(RecoveryPlanJobIntentResponse)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return response, op.client.Do(ctx, req, response)
}

/*Get the status of a recovery plan job.
 * This operation gets the execution status of a recovery plan job identified by 'uuid'.
 *
 * @param uuid UUID of the recovery plan job.
 * @param status is one of {'execution_status', 'cleanup_status'}
 */
func (op Operations) GetRecoveryPlanJobStatus(ctx context.Context, uuid string, status string) (*RecoveryPlanJobExecutionStatus, error) {
	path := fmt.Sprintf("/recovery_plan_jobs/%s/%s", uuid, status)
	executionStatus := new(RecoveryPlanJobExecutionStatus)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return executionStatus, op.client.Do(ctx, req, executionStatus)
}

/*List recovery plan jobs.
 * This Operations lists the recovery plan jobs matching the criteria specified in 'request'.
 *
 * @param request pointer to specification of type DSMetadata.
 */
func (op Operations) ListRecoveryPlanJobs(ctx context.Context, request *DSMetadata) (*RecoveryPlanJobListResponse, error) {
	path := "/recovery_plan_jobs/list"

	list := new(RecoveryPlanJobListResponse)

	req, err := op.client.NewRequest(http.MethodPost, path, request)
	if err != nil {
		return nil, err
	}

	return list, op.client.Do(ctx, req, list)
}

/*Get a projection of the attributes of entities of certain type.
 * This operation returns the attributes of the entities that match the
 * filter criteria specified in 'request'.
 *
 * @param request pointer to the specification of type GroupsGetEntitiesRequest
 */
func (op Operations) GroupsGetEntities(ctx context.Context, request *GroupsGetEntitiesRequest,
) (*GroupsGetEntitiesResponse, error) {
	req, err := op.client.NewRequest(http.MethodPost, "/groups", request)
	response := new(GroupsGetEntitiesResponse)

	if err != nil {
		return nil, err
	}

	return response, op.client.Do(ctx, req, response)
}

/*Get information about an availability zone (AZ).
 * This operation gets the information about an AZ identified by 'uuid'.
 *
 * @param uuid UUID of the AZ.
 */
func (op Operations) GetAvailabilityZone(ctx context.Context, uuid string) (*AvailabilityZoneIntentResponse, error) {
	path := fmt.Sprintf("/availability_zones/%s", uuid)
	response := new(AvailabilityZoneIntentResponse)

	req, err := op.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	return response, op.client.Do(ctx, req, response)
}

// GetPrismCentral gets the information about the Prism Central
func (op Operations) GetPrismCentral(ctx context.Context) (*models.PrismCentral, error) {
	path := "/prism_central"
	response := new(models.PrismCentral)

	req, err := op.client.NewRequest(http.MethodGet, path, response)
	if err != nil {
		return nil, err
	}

	return response, op.client.Do(ctx, req, response)
}
