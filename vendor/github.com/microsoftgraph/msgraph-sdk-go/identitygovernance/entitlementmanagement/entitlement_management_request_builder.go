package entitlementmanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i069e71eb6a0d6756f4420f4f811017c7470620ad5174a1e3aab2db0d8cb430c2 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/catalogs"
    i3e9bbff4ac7b59c7638e82a3457169acb55d826f7187ac6e78fbe6ec1bb80d7d "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/assignments"
    i441d3614c41d6a7eafb923dbc4975591570abe0ca8005d661c628e49592a07f8 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/settings"
    i458b5d51a7c8ed196efc57bee2af7eeeea44a7093ec6fe9b40b0d87560256cca "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/accesspackages"
    i4a6e5bbedb985fd4e83bec609cf8d21a776d26a1c6807eac6ae0d88110984ab9 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/accesspackageassignmentapprovals"
    i6676dbd3ebb41db58dd801dbe33621dbe43cb8d0142648650dee5f3cf366ddd0 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/assignmentpolicies"
    i890fb4653f4dff2b9ab5a72705a40947b79d9bb64793d05773197622a996fc61 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/connectedorganizations"
    i8d2bfa033b27ef0ce19de47e723f8c3475251f8917d2cbf768fc4e424505f65e "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/assignmentrequests"
    i0b0fdc1abe3e98caf17dadc8636597ba88c5644fc68edc06b6f9a0aeb66660bc "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/assignmentpolicies/item"
    i3a731cbed9792cd58eb16385c435f4b1af4fafbdc6a965b8534ab7ff7df19e39 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/connectedorganizations/item"
    i558fcddcc0d01116d05fce610754153bd32852762580e628e25c34b7d8fffa4a "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/assignments/item"
    i76acc551f9c3cb51943981936cb0b2aa4089fd504de3bf817176dbbe448d3fb7 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/assignmentrequests/item"
    i98d02a4533e64ff168c5023553f526cbe65a72cacee2308e1a90cb06055c2a93 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/accesspackages/item"
    i9f1127e219619b20720c5350652063980d97d3210fea76590720646f27ba45bb "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/catalogs/item"
    ib6ba4a2debf33ace4ef4d7d21e119e5c8a37451c871b2845657ae921beea0814 "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance/entitlementmanagement/accesspackageassignmentapprovals/item"
)

// EntitlementManagementRequestBuilder provides operations to manage the entitlementManagement property of the microsoft.graph.identityGovernance entity.
type EntitlementManagementRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// EntitlementManagementRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// EntitlementManagementRequestBuilderGetQueryParameters get entitlementManagement from identityGovernance
type EntitlementManagementRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// EntitlementManagementRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EntitlementManagementRequestBuilderGetQueryParameters
}
// EntitlementManagementRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AccessPackageAssignmentApprovals provides operations to manage the accessPackageAssignmentApprovals property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) AccessPackageAssignmentApprovals()(*i4a6e5bbedb985fd4e83bec609cf8d21a776d26a1c6807eac6ae0d88110984ab9.AccessPackageAssignmentApprovalsRequestBuilder) {
    return i4a6e5bbedb985fd4e83bec609cf8d21a776d26a1c6807eac6ae0d88110984ab9.NewAccessPackageAssignmentApprovalsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AccessPackageAssignmentApprovalsById provides operations to manage the accessPackageAssignmentApprovals property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) AccessPackageAssignmentApprovalsById(id string)(*ib6ba4a2debf33ace4ef4d7d21e119e5c8a37451c871b2845657ae921beea0814.ApprovalItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["approval%2Did"] = id
    }
    return ib6ba4a2debf33ace4ef4d7d21e119e5c8a37451c871b2845657ae921beea0814.NewApprovalItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AccessPackages provides operations to manage the accessPackages property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) AccessPackages()(*i458b5d51a7c8ed196efc57bee2af7eeeea44a7093ec6fe9b40b0d87560256cca.AccessPackagesRequestBuilder) {
    return i458b5d51a7c8ed196efc57bee2af7eeeea44a7093ec6fe9b40b0d87560256cca.NewAccessPackagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AccessPackagesById provides operations to manage the accessPackages property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) AccessPackagesById(id string)(*i98d02a4533e64ff168c5023553f526cbe65a72cacee2308e1a90cb06055c2a93.AccessPackageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["accessPackage%2Did"] = id
    }
    return i98d02a4533e64ff168c5023553f526cbe65a72cacee2308e1a90cb06055c2a93.NewAccessPackageItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AssignmentPolicies provides operations to manage the assignmentPolicies property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) AssignmentPolicies()(*i6676dbd3ebb41db58dd801dbe33621dbe43cb8d0142648650dee5f3cf366ddd0.AssignmentPoliciesRequestBuilder) {
    return i6676dbd3ebb41db58dd801dbe33621dbe43cb8d0142648650dee5f3cf366ddd0.NewAssignmentPoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AssignmentPoliciesById provides operations to manage the assignmentPolicies property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) AssignmentPoliciesById(id string)(*i0b0fdc1abe3e98caf17dadc8636597ba88c5644fc68edc06b6f9a0aeb66660bc.AccessPackageAssignmentPolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["accessPackageAssignmentPolicy%2Did"] = id
    }
    return i0b0fdc1abe3e98caf17dadc8636597ba88c5644fc68edc06b6f9a0aeb66660bc.NewAccessPackageAssignmentPolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AssignmentRequests provides operations to manage the assignmentRequests property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) AssignmentRequests()(*i8d2bfa033b27ef0ce19de47e723f8c3475251f8917d2cbf768fc4e424505f65e.AssignmentRequestsRequestBuilder) {
    return i8d2bfa033b27ef0ce19de47e723f8c3475251f8917d2cbf768fc4e424505f65e.NewAssignmentRequestsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AssignmentRequestsById provides operations to manage the assignmentRequests property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) AssignmentRequestsById(id string)(*i76acc551f9c3cb51943981936cb0b2aa4089fd504de3bf817176dbbe448d3fb7.AccessPackageAssignmentRequestItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["accessPackageAssignmentRequest%2Did"] = id
    }
    return i76acc551f9c3cb51943981936cb0b2aa4089fd504de3bf817176dbbe448d3fb7.NewAccessPackageAssignmentRequestItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Assignments provides operations to manage the assignments property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) Assignments()(*i3e9bbff4ac7b59c7638e82a3457169acb55d826f7187ac6e78fbe6ec1bb80d7d.AssignmentsRequestBuilder) {
    return i3e9bbff4ac7b59c7638e82a3457169acb55d826f7187ac6e78fbe6ec1bb80d7d.NewAssignmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AssignmentsById provides operations to manage the assignments property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) AssignmentsById(id string)(*i558fcddcc0d01116d05fce610754153bd32852762580e628e25c34b7d8fffa4a.AccessPackageAssignmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["accessPackageAssignment%2Did"] = id
    }
    return i558fcddcc0d01116d05fce610754153bd32852762580e628e25c34b7d8fffa4a.NewAccessPackageAssignmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Catalogs provides operations to manage the catalogs property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) Catalogs()(*i069e71eb6a0d6756f4420f4f811017c7470620ad5174a1e3aab2db0d8cb430c2.CatalogsRequestBuilder) {
    return i069e71eb6a0d6756f4420f4f811017c7470620ad5174a1e3aab2db0d8cb430c2.NewCatalogsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CatalogsById provides operations to manage the catalogs property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) CatalogsById(id string)(*i9f1127e219619b20720c5350652063980d97d3210fea76590720646f27ba45bb.AccessPackageCatalogItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["accessPackageCatalog%2Did"] = id
    }
    return i9f1127e219619b20720c5350652063980d97d3210fea76590720646f27ba45bb.NewAccessPackageCatalogItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ConnectedOrganizations provides operations to manage the connectedOrganizations property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) ConnectedOrganizations()(*i890fb4653f4dff2b9ab5a72705a40947b79d9bb64793d05773197622a996fc61.ConnectedOrganizationsRequestBuilder) {
    return i890fb4653f4dff2b9ab5a72705a40947b79d9bb64793d05773197622a996fc61.NewConnectedOrganizationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ConnectedOrganizationsById provides operations to manage the connectedOrganizations property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) ConnectedOrganizationsById(id string)(*i3a731cbed9792cd58eb16385c435f4b1af4fafbdc6a965b8534ab7ff7df19e39.ConnectedOrganizationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["connectedOrganization%2Did"] = id
    }
    return i3a731cbed9792cd58eb16385c435f4b1af4fafbdc6a965b8534ab7ff7df19e39.NewConnectedOrganizationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewEntitlementManagementRequestBuilderInternal instantiates a new EntitlementManagementRequestBuilder and sets the default values.
func NewEntitlementManagementRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementRequestBuilder) {
    m := &EntitlementManagementRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/identityGovernance/entitlementManagement{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewEntitlementManagementRequestBuilder instantiates a new EntitlementManagementRequestBuilder and sets the default values.
func NewEntitlementManagementRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEntitlementManagementRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property entitlementManagement for identityGovernance
func (m *EntitlementManagementRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *EntitlementManagementRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation get entitlementManagement from identityGovernance
func (m *EntitlementManagementRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *EntitlementManagementRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property entitlementManagement in identityGovernance
func (m *EntitlementManagementRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EntitlementManagementable, requestConfiguration *EntitlementManagementRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property entitlementManagement for identityGovernance
func (m *EntitlementManagementRequestBuilder) Delete(ctx context.Context, requestConfiguration *EntitlementManagementRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get entitlementManagement from identityGovernance
func (m *EntitlementManagementRequestBuilder) Get(ctx context.Context, requestConfiguration *EntitlementManagementRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EntitlementManagementable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEntitlementManagementFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EntitlementManagementable), nil
}
// Patch update the navigation property entitlementManagement in identityGovernance
func (m *EntitlementManagementRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EntitlementManagementable, requestConfiguration *EntitlementManagementRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EntitlementManagementable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEntitlementManagementFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EntitlementManagementable), nil
}
// Settings provides operations to manage the settings property of the microsoft.graph.entitlementManagement entity.
func (m *EntitlementManagementRequestBuilder) Settings()(*i441d3614c41d6a7eafb923dbc4975591570abe0ca8005d661c628e49592a07f8.SettingsRequestBuilder) {
    return i441d3614c41d6a7eafb923dbc4975591570abe0ca8005d661c628e49592a07f8.NewSettingsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
