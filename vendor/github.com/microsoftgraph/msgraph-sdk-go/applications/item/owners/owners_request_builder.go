package owners

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i060d05043ef21a20bca5bba0cf1bf9d4437b1d91d3a0d0ae6c97389d4d00694e "github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners/count"
    i1b94d5383a8ca1e829f2c48b74c036b65d5c7cacb4223e5d4d65ae184d162718 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners/approleassignment"
    i4481ae29272b90e51b803d7ad16e7fa37075d7b4688b9c3d4f259f1c3b71d79b "github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners/endpoint"
    i9d70f273879b63f77ee24edf5ec9817b23b4568da6e6c8a7d96eac4aba1f6139 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners/serviceprincipal"
    iaad2cc02c08e8ffbe0158ad9cce1d33313e353802c2a6aa54e7e7b1c9ee4ab64 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners/user"
    if16d68f8f5e29587f2c3f4904b28d4809aef9ba1ded5e7526612370ed1a35458 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners/ref"
)

// OwnersRequestBuilder provides operations to manage the owners property of the microsoft.graph.application entity.
type OwnersRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// OwnersRequestBuilderGetQueryParameters directory objects that are owners of the application. Read-only. Nullable. Supports $expand and $filter (eq when counting empty collections).
type OwnersRequestBuilderGetQueryParameters struct {
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
// OwnersRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OwnersRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *OwnersRequestBuilderGetQueryParameters
}
// AppRoleAssignment casts the previous resource to appRoleAssignment.
func (m *OwnersRequestBuilder) AppRoleAssignment()(*i1b94d5383a8ca1e829f2c48b74c036b65d5c7cacb4223e5d4d65ae184d162718.AppRoleAssignmentRequestBuilder) {
    return i1b94d5383a8ca1e829f2c48b74c036b65d5c7cacb4223e5d4d65ae184d162718.NewAppRoleAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewOwnersRequestBuilderInternal instantiates a new OwnersRequestBuilder and sets the default values.
func NewOwnersRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OwnersRequestBuilder) {
    m := &OwnersRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/applications/{application%2Did}/owners{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewOwnersRequestBuilder instantiates a new OwnersRequestBuilder and sets the default values.
func NewOwnersRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OwnersRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewOwnersRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *OwnersRequestBuilder) Count()(*i060d05043ef21a20bca5bba0cf1bf9d4437b1d91d3a0d0ae6c97389d4d00694e.CountRequestBuilder) {
    return i060d05043ef21a20bca5bba0cf1bf9d4437b1d91d3a0d0ae6c97389d4d00694e.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation directory objects that are owners of the application. Read-only. Nullable. Supports $expand and $filter (eq when counting empty collections).
func (m *OwnersRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *OwnersRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *OwnersRequestBuilder) Endpoint()(*i4481ae29272b90e51b803d7ad16e7fa37075d7b4688b9c3d4f259f1c3b71d79b.EndpointRequestBuilder) {
    return i4481ae29272b90e51b803d7ad16e7fa37075d7b4688b9c3d4f259f1c3b71d79b.NewEndpointRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get directory objects that are owners of the application. Read-only. Nullable. Supports $expand and $filter (eq when counting empty collections).
func (m *OwnersRequestBuilder) Get(ctx context.Context, requestConfiguration *OwnersRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
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
// Ref provides operations to manage the collection of application entities.
func (m *OwnersRequestBuilder) Ref()(*if16d68f8f5e29587f2c3f4904b28d4809aef9ba1ded5e7526612370ed1a35458.RefRequestBuilder) {
    return if16d68f8f5e29587f2c3f4904b28d4809aef9ba1ded5e7526612370ed1a35458.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *OwnersRequestBuilder) ServicePrincipal()(*i9d70f273879b63f77ee24edf5ec9817b23b4568da6e6c8a7d96eac4aba1f6139.ServicePrincipalRequestBuilder) {
    return i9d70f273879b63f77ee24edf5ec9817b23b4568da6e6c8a7d96eac4aba1f6139.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *OwnersRequestBuilder) User()(*iaad2cc02c08e8ffbe0158ad9cce1d33313e353802c2a6aa54e7e7b1c9ee4ab64.UserRequestBuilder) {
    return iaad2cc02c08e8ffbe0158ad9cce1d33313e353802c2a6aa54e7e7b1c9ee4ab64.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
