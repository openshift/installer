package ownedobjects

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i17b02e26889fb0f080604a4de66ce6fd380d546f8c196280db674e7b7cd97f48 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/ownedobjects/endpoint"
    i32bc732ebd661320e929f285389b8d78cd33a237e87cd32d59d98a4ce77be380 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/ownedobjects/serviceprincipal"
    i3bd45a7989145bacecf06a82a4bf204177f08a7d528fa9fa7c7793c38d76c6f6 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/ownedobjects/count"
    ic3da8a993087f20f550c11ecd9ef2a6e63ba2afd8289604927f1186658f3fd97 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/ownedobjects/application"
    ie13636b5ae2b4b848c8e6999da44a02f263a358d521250880031e785e6a673a2 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/ownedobjects/approleassignment"
    if4d65e73dedc21559f062dc84d7e5b84a8fd8a7e20959d667f30a36c4d3273fd "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/ownedobjects/group"
)

// OwnedObjectsRequestBuilder provides operations to manage the ownedObjects property of the microsoft.graph.servicePrincipal entity.
type OwnedObjectsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// OwnedObjectsRequestBuilderGetQueryParameters directory objects that are owned by this service principal. Read-only. Nullable. Supports $expand.
type OwnedObjectsRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// OwnedObjectsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OwnedObjectsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *OwnedObjectsRequestBuilderGetQueryParameters
}
// Application casts the previous resource to application.
func (m *OwnedObjectsRequestBuilder) Application()(*ic3da8a993087f20f550c11ecd9ef2a6e63ba2afd8289604927f1186658f3fd97.ApplicationRequestBuilder) {
    return ic3da8a993087f20f550c11ecd9ef2a6e63ba2afd8289604927f1186658f3fd97.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AppRoleAssignment casts the previous resource to appRoleAssignment.
func (m *OwnedObjectsRequestBuilder) AppRoleAssignment()(*ie13636b5ae2b4b848c8e6999da44a02f263a358d521250880031e785e6a673a2.AppRoleAssignmentRequestBuilder) {
    return ie13636b5ae2b4b848c8e6999da44a02f263a358d521250880031e785e6a673a2.NewAppRoleAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewOwnedObjectsRequestBuilderInternal instantiates a new OwnedObjectsRequestBuilder and sets the default values.
func NewOwnedObjectsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OwnedObjectsRequestBuilder) {
    m := &OwnedObjectsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/ownedObjects{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewOwnedObjectsRequestBuilder instantiates a new OwnedObjectsRequestBuilder and sets the default values.
func NewOwnedObjectsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OwnedObjectsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewOwnedObjectsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *OwnedObjectsRequestBuilder) Count()(*i3bd45a7989145bacecf06a82a4bf204177f08a7d528fa9fa7c7793c38d76c6f6.CountRequestBuilder) {
    return i3bd45a7989145bacecf06a82a4bf204177f08a7d528fa9fa7c7793c38d76c6f6.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation directory objects that are owned by this service principal. Read-only. Nullable. Supports $expand.
func (m *OwnedObjectsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *OwnedObjectsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Endpoint casts the previous resource to endpoint.
func (m *OwnedObjectsRequestBuilder) Endpoint()(*i17b02e26889fb0f080604a4de66ce6fd380d546f8c196280db674e7b7cd97f48.EndpointRequestBuilder) {
    return i17b02e26889fb0f080604a4de66ce6fd380d546f8c196280db674e7b7cd97f48.NewEndpointRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get directory objects that are owned by this service principal. Read-only. Nullable. Supports $expand.
func (m *OwnedObjectsRequestBuilder) Get(ctx context.Context, requestConfiguration *OwnedObjectsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable), nil
}
// Group casts the previous resource to group.
func (m *OwnedObjectsRequestBuilder) Group()(*if4d65e73dedc21559f062dc84d7e5b84a8fd8a7e20959d667f30a36c4d3273fd.GroupRequestBuilder) {
    return if4d65e73dedc21559f062dc84d7e5b84a8fd8a7e20959d667f30a36c4d3273fd.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *OwnedObjectsRequestBuilder) ServicePrincipal()(*i32bc732ebd661320e929f285389b8d78cd33a237e87cd32d59d98a4ce77be380.ServicePrincipalRequestBuilder) {
    return i32bc732ebd661320e929f285389b8d78cd33a237e87cd32d59d98a4ce77be380.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
