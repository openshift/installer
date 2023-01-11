package registeredowners

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1ab9540be26799d5e525377eae39b656dcfd671ed49f12f5ba1572b0c1c5ea15 "github.com/microsoftgraph/msgraph-sdk-go/devices/item/registeredowners/user"
    i58d81abe2a94e2e873f495b3f6bdebbf13ae06b03f1b58c6521100d5acada25b "github.com/microsoftgraph/msgraph-sdk-go/devices/item/registeredowners/approleassignment"
    i764856d5a0d4d32224ab2b48aa026d2cdef858237c89fdf7e542463f4cfcfaf3 "github.com/microsoftgraph/msgraph-sdk-go/devices/item/registeredowners/count"
    ib5cc46b500e9b24fa270a613e8e49d845012dfb5efeb884206b23274ffedaa91 "github.com/microsoftgraph/msgraph-sdk-go/devices/item/registeredowners/ref"
    icf24949f036260cdde78a7a1b7a3bd43d492d348b222c0088fc7b628e85f2f76 "github.com/microsoftgraph/msgraph-sdk-go/devices/item/registeredowners/serviceprincipal"
    ifb6aab027c17b4a3f323306d77c23d7c4148efc04152f28b877b9b61ab62ddbd "github.com/microsoftgraph/msgraph-sdk-go/devices/item/registeredowners/endpoint"
)

// RegisteredOwnersRequestBuilder provides operations to manage the registeredOwners property of the microsoft.graph.device entity.
type RegisteredOwnersRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// RegisteredOwnersRequestBuilderGetQueryParameters the user that cloud joined the device or registered their personal device. The registered owner is set at the time of registration. Currently, there can be only one owner. Read-only. Nullable. Supports $expand.
type RegisteredOwnersRequestBuilderGetQueryParameters struct {
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
// RegisteredOwnersRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type RegisteredOwnersRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *RegisteredOwnersRequestBuilderGetQueryParameters
}
// AppRoleAssignment casts the previous resource to appRoleAssignment.
func (m *RegisteredOwnersRequestBuilder) AppRoleAssignment()(*i58d81abe2a94e2e873f495b3f6bdebbf13ae06b03f1b58c6521100d5acada25b.AppRoleAssignmentRequestBuilder) {
    return i58d81abe2a94e2e873f495b3f6bdebbf13ae06b03f1b58c6521100d5acada25b.NewAppRoleAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewRegisteredOwnersRequestBuilderInternal instantiates a new RegisteredOwnersRequestBuilder and sets the default values.
func NewRegisteredOwnersRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RegisteredOwnersRequestBuilder) {
    m := &RegisteredOwnersRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/devices/{device%2Did}/registeredOwners{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewRegisteredOwnersRequestBuilder instantiates a new RegisteredOwnersRequestBuilder and sets the default values.
func NewRegisteredOwnersRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RegisteredOwnersRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewRegisteredOwnersRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *RegisteredOwnersRequestBuilder) Count()(*i764856d5a0d4d32224ab2b48aa026d2cdef858237c89fdf7e542463f4cfcfaf3.CountRequestBuilder) {
    return i764856d5a0d4d32224ab2b48aa026d2cdef858237c89fdf7e542463f4cfcfaf3.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation the user that cloud joined the device or registered their personal device. The registered owner is set at the time of registration. Currently, there can be only one owner. Read-only. Nullable. Supports $expand.
func (m *RegisteredOwnersRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *RegisteredOwnersRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *RegisteredOwnersRequestBuilder) Endpoint()(*ifb6aab027c17b4a3f323306d77c23d7c4148efc04152f28b877b9b61ab62ddbd.EndpointRequestBuilder) {
    return ifb6aab027c17b4a3f323306d77c23d7c4148efc04152f28b877b9b61ab62ddbd.NewEndpointRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the user that cloud joined the device or registered their personal device. The registered owner is set at the time of registration. Currently, there can be only one owner. Read-only. Nullable. Supports $expand.
func (m *RegisteredOwnersRequestBuilder) Get(ctx context.Context, requestConfiguration *RegisteredOwnersRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
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
// Ref provides operations to manage the collection of device entities.
func (m *RegisteredOwnersRequestBuilder) Ref()(*ib5cc46b500e9b24fa270a613e8e49d845012dfb5efeb884206b23274ffedaa91.RefRequestBuilder) {
    return ib5cc46b500e9b24fa270a613e8e49d845012dfb5efeb884206b23274ffedaa91.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *RegisteredOwnersRequestBuilder) ServicePrincipal()(*icf24949f036260cdde78a7a1b7a3bd43d492d348b222c0088fc7b628e85f2f76.ServicePrincipalRequestBuilder) {
    return icf24949f036260cdde78a7a1b7a3bd43d492d348b222c0088fc7b628e85f2f76.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *RegisteredOwnersRequestBuilder) User()(*i1ab9540be26799d5e525377eae39b656dcfd671ed49f12f5ba1572b0c1c5ea15.UserRequestBuilder) {
    return i1ab9540be26799d5e525377eae39b656dcfd671ed49f12f5ba1572b0c1c5ea15.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
