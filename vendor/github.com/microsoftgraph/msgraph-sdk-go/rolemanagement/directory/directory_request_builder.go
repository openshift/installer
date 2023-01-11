package directory

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1ecb617cf3506073901d6ebbc926c43ae20c2fe24b4f636c07d76b404da7758e "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleassignmentschedulerequests"
    i57b910acb4194d77e597cd7fc038ff1fcdab95ffb0c6387ecd462abda735c54f "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleassignments"
    i58881c44e386bd8753efbaac86e07be22d26f1b8cdc063b36c6658beb99ed327 "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleeligibilityscheduleinstances"
    i65e8721e08cadca66e59fe29786f285c292916fbfaf6d187c1ae36f76b2d8e58 "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roledefinitions"
    i87846e0996ce71914c22736c34bd7efcde223621a8ded9cd826cc177e595612a "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleassignmentscheduleinstances"
    ia26250e0869da35555e48bccfc062cfedf96ec469b8f446ad12e910062c70fe6 "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleeligibilityschedules"
    ia7a7c59e07028f1682d2f9d60c9932fb4ee80d611ed68b5b5ef11547cd2cba70 "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleeligibilityschedulerequests"
    ibbba73b169d00073d04d75c6dbf204f7d1e7828fda99fa12c50758992179084b "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleassignmentschedules"
    i26564a3dcf6f5b8f339c5029b420a35ad925e9ec8a807e073a522418b894b461 "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleassignments/item"
    i329ceabcf38e19751c19b2a6282c453f471d4acda1f60ae66f8e961ea5ffe9e7 "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleeligibilityschedules/item"
    i38a0099911c322678afe0e91db6bf13f1336c0bd7474a79696eab0396dd25d72 "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleeligibilityscheduleinstances/item"
    i93f0046b6ea93bb22d2872d826c5bb0b32c7bd776d82404d212ed81ae3bfa20c "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roledefinitions/item"
    i97ec61667ae4af0dc73f94ca00501e020bfd2cab6faca9653c9cdf54198c1c7b "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleassignmentschedulerequests/item"
    ied8d8d84ec67ed4fa1c183cd3b843abe9fbf78b547b4078399ffd2ded02c5b9f "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleeligibilityschedulerequests/item"
    if28346ffd76078448f76b8e3a6438c69f288f33de995f80b4470bf3dd729afa7 "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleassignmentscheduleinstances/item"
    ifcfeaaa38c74c27248ac242dc390a943df9fda837089f362cf5a0b616515e16e "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement/directory/roleassignmentschedules/item"
)

// DirectoryRequestBuilder provides operations to manage the directory property of the microsoft.graph.roleManagement entity.
type DirectoryRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// DirectoryRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DirectoryRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// DirectoryRequestBuilderGetQueryParameters get directory from roleManagement
type DirectoryRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DirectoryRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DirectoryRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DirectoryRequestBuilderGetQueryParameters
}
// DirectoryRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DirectoryRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewDirectoryRequestBuilderInternal instantiates a new DirectoryRequestBuilder and sets the default values.
func NewDirectoryRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryRequestBuilder) {
    m := &DirectoryRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/roleManagement/directory{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewDirectoryRequestBuilder instantiates a new DirectoryRequestBuilder and sets the default values.
func NewDirectoryRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDirectoryRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property directory for roleManagement
func (m *DirectoryRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *DirectoryRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE
    if requestConfiguration != nil {
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// CreateGetRequestInformation get directory from roleManagement
func (m *DirectoryRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *DirectoryRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET
    requestInfo.Headers["Accept"] = "application/json"
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// CreatePatchRequestInformation update the navigation property directory in roleManagement
func (m *DirectoryRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RbacApplicationable, requestConfiguration *DirectoryRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH
    requestInfo.Headers["Accept"] = "application/json"
    requestInfo.SetContentFromParsable(ctx, m.requestAdapter, "application/json", body)
    if requestConfiguration != nil {
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// Delete delete navigation property directory for roleManagement
func (m *DirectoryRequestBuilder) Delete(ctx context.Context, requestConfiguration *DirectoryRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.CreateDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.requestAdapter.SendNoContentAsync(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// Get get directory from roleManagement
func (m *DirectoryRequestBuilder) Get(ctx context.Context, requestConfiguration *DirectoryRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RbacApplicationable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateRbacApplicationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RbacApplicationable), nil
}
// Patch update the navigation property directory in roleManagement
func (m *DirectoryRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RbacApplicationable, requestConfiguration *DirectoryRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RbacApplicationable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateRbacApplicationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.RbacApplicationable), nil
}
// RoleAssignments provides operations to manage the roleAssignments property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleAssignments()(*i57b910acb4194d77e597cd7fc038ff1fcdab95ffb0c6387ecd462abda735c54f.RoleAssignmentsRequestBuilder) {
    return i57b910acb4194d77e597cd7fc038ff1fcdab95ffb0c6387ecd462abda735c54f.NewRoleAssignmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RoleAssignmentsById provides operations to manage the roleAssignments property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleAssignmentsById(id string)(*i26564a3dcf6f5b8f339c5029b420a35ad925e9ec8a807e073a522418b894b461.UnifiedRoleAssignmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["unifiedRoleAssignment%2Did"] = id
    }
    return i26564a3dcf6f5b8f339c5029b420a35ad925e9ec8a807e073a522418b894b461.NewUnifiedRoleAssignmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// RoleAssignmentScheduleInstances provides operations to manage the roleAssignmentScheduleInstances property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleAssignmentScheduleInstances()(*i87846e0996ce71914c22736c34bd7efcde223621a8ded9cd826cc177e595612a.RoleAssignmentScheduleInstancesRequestBuilder) {
    return i87846e0996ce71914c22736c34bd7efcde223621a8ded9cd826cc177e595612a.NewRoleAssignmentScheduleInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RoleAssignmentScheduleInstancesById provides operations to manage the roleAssignmentScheduleInstances property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleAssignmentScheduleInstancesById(id string)(*if28346ffd76078448f76b8e3a6438c69f288f33de995f80b4470bf3dd729afa7.UnifiedRoleAssignmentScheduleInstanceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["unifiedRoleAssignmentScheduleInstance%2Did"] = id
    }
    return if28346ffd76078448f76b8e3a6438c69f288f33de995f80b4470bf3dd729afa7.NewUnifiedRoleAssignmentScheduleInstanceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// RoleAssignmentScheduleRequests provides operations to manage the roleAssignmentScheduleRequests property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleAssignmentScheduleRequests()(*i1ecb617cf3506073901d6ebbc926c43ae20c2fe24b4f636c07d76b404da7758e.RoleAssignmentScheduleRequestsRequestBuilder) {
    return i1ecb617cf3506073901d6ebbc926c43ae20c2fe24b4f636c07d76b404da7758e.NewRoleAssignmentScheduleRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RoleAssignmentScheduleRequestsById provides operations to manage the roleAssignmentScheduleRequests property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleAssignmentScheduleRequestsById(id string)(*i97ec61667ae4af0dc73f94ca00501e020bfd2cab6faca9653c9cdf54198c1c7b.UnifiedRoleAssignmentScheduleRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["unifiedRoleAssignmentScheduleRequest%2Did"] = id
    }
    return i97ec61667ae4af0dc73f94ca00501e020bfd2cab6faca9653c9cdf54198c1c7b.NewUnifiedRoleAssignmentScheduleRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// RoleAssignmentSchedules provides operations to manage the roleAssignmentSchedules property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleAssignmentSchedules()(*ibbba73b169d00073d04d75c6dbf204f7d1e7828fda99fa12c50758992179084b.RoleAssignmentSchedulesRequestBuilder) {
    return ibbba73b169d00073d04d75c6dbf204f7d1e7828fda99fa12c50758992179084b.NewRoleAssignmentSchedulesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RoleAssignmentSchedulesById provides operations to manage the roleAssignmentSchedules property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleAssignmentSchedulesById(id string)(*ifcfeaaa38c74c27248ac242dc390a943df9fda837089f362cf5a0b616515e16e.UnifiedRoleAssignmentScheduleItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["unifiedRoleAssignmentSchedule%2Did"] = id
    }
    return ifcfeaaa38c74c27248ac242dc390a943df9fda837089f362cf5a0b616515e16e.NewUnifiedRoleAssignmentScheduleItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// RoleDefinitions provides operations to manage the roleDefinitions property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleDefinitions()(*i65e8721e08cadca66e59fe29786f285c292916fbfaf6d187c1ae36f76b2d8e58.RoleDefinitionsRequestBuilder) {
    return i65e8721e08cadca66e59fe29786f285c292916fbfaf6d187c1ae36f76b2d8e58.NewRoleDefinitionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RoleDefinitionsById provides operations to manage the roleDefinitions property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleDefinitionsById(id string)(*i93f0046b6ea93bb22d2872d826c5bb0b32c7bd776d82404d212ed81ae3bfa20c.UnifiedRoleDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["unifiedRoleDefinition%2Did"] = id
    }
    return i93f0046b6ea93bb22d2872d826c5bb0b32c7bd776d82404d212ed81ae3bfa20c.NewUnifiedRoleDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// RoleEligibilityScheduleInstances provides operations to manage the roleEligibilityScheduleInstances property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleEligibilityScheduleInstances()(*i58881c44e386bd8753efbaac86e07be22d26f1b8cdc063b36c6658beb99ed327.RoleEligibilityScheduleInstancesRequestBuilder) {
    return i58881c44e386bd8753efbaac86e07be22d26f1b8cdc063b36c6658beb99ed327.NewRoleEligibilityScheduleInstancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RoleEligibilityScheduleInstancesById provides operations to manage the roleEligibilityScheduleInstances property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleEligibilityScheduleInstancesById(id string)(*i38a0099911c322678afe0e91db6bf13f1336c0bd7474a79696eab0396dd25d72.UnifiedRoleEligibilityScheduleInstanceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["unifiedRoleEligibilityScheduleInstance%2Did"] = id
    }
    return i38a0099911c322678afe0e91db6bf13f1336c0bd7474a79696eab0396dd25d72.NewUnifiedRoleEligibilityScheduleInstanceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// RoleEligibilityScheduleRequests provides operations to manage the roleEligibilityScheduleRequests property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleEligibilityScheduleRequests()(*ia7a7c59e07028f1682d2f9d60c9932fb4ee80d611ed68b5b5ef11547cd2cba70.RoleEligibilityScheduleRequestsRequestBuilder) {
    return ia7a7c59e07028f1682d2f9d60c9932fb4ee80d611ed68b5b5ef11547cd2cba70.NewRoleEligibilityScheduleRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RoleEligibilityScheduleRequestsById provides operations to manage the roleEligibilityScheduleRequests property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleEligibilityScheduleRequestsById(id string)(*ied8d8d84ec67ed4fa1c183cd3b843abe9fbf78b547b4078399ffd2ded02c5b9f.UnifiedRoleEligibilityScheduleRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["unifiedRoleEligibilityScheduleRequest%2Did"] = id
    }
    return ied8d8d84ec67ed4fa1c183cd3b843abe9fbf78b547b4078399ffd2ded02c5b9f.NewUnifiedRoleEligibilityScheduleRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// RoleEligibilitySchedules provides operations to manage the roleEligibilitySchedules property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleEligibilitySchedules()(*ia26250e0869da35555e48bccfc062cfedf96ec469b8f446ad12e910062c70fe6.RoleEligibilitySchedulesRequestBuilder) {
    return ia26250e0869da35555e48bccfc062cfedf96ec469b8f446ad12e910062c70fe6.NewRoleEligibilitySchedulesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RoleEligibilitySchedulesById provides operations to manage the roleEligibilitySchedules property of the microsoft.graph.rbacApplication entity.
func (m *DirectoryRequestBuilder) RoleEligibilitySchedulesById(id string)(*i329ceabcf38e19751c19b2a6282c453f471d4acda1f60ae66f8e961ea5ffe9e7.UnifiedRoleEligibilityScheduleItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["unifiedRoleEligibilitySchedule%2Did"] = id
    }
    return i329ceabcf38e19751c19b2a6282c453f471d4acda1f60ae66f8e961ea5ffe9e7.NewUnifiedRoleEligibilityScheduleItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
