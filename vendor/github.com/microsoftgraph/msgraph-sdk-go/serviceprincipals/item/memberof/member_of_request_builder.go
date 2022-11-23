package memberof

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0b27242c7dcbdd771786e0f446664593f09940a382381f2622fc6a1c3e458e07 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/memberof/serviceprincipal"
    i4454307873c5851e347787ba6ab2d3940ed9e7a55cf56ac268b1aac06400e6d7 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/memberof/device"
    i5e6f3053da49caf881450f09581c96acb41911545fcc8509c552944f2e192a67 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/memberof/orgcontact"
    iabf57bc1560ca5354a1c2b8ce106ad4b265aae6ce3ba5d26216db3b613dd81db "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/memberof/group"
    ib70fe1873991f79b66e17a20a4d22ae11e48b9c6decf378d8edf9a0a0969a919 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/memberof/user"
    id7da27b83068feaec5bb92a93f3dccf50f1695decc5eff8e15ef6f2653caaed0 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/memberof/application"
    id92750a7b3cc4ed2656a6e56bb755836906fdc5131ade6ac4b3ffd52bab76f3a "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/memberof/count"
)

// MemberOfRequestBuilder provides operations to manage the memberOf property of the microsoft.graph.servicePrincipal entity.
type MemberOfRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// MemberOfRequestBuilderGetQueryParameters roles that this service principal is a member of. HTTP Methods: GET Read-only. Nullable. Supports $expand.
type MemberOfRequestBuilderGetQueryParameters struct {
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
// MemberOfRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MemberOfRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MemberOfRequestBuilderGetQueryParameters
}
// Application casts the previous resource to application.
func (m *MemberOfRequestBuilder) Application()(*id7da27b83068feaec5bb92a93f3dccf50f1695decc5eff8e15ef6f2653caaed0.ApplicationRequestBuilder) {
    return id7da27b83068feaec5bb92a93f3dccf50f1695decc5eff8e15ef6f2653caaed0.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewMemberOfRequestBuilderInternal instantiates a new MemberOfRequestBuilder and sets the default values.
func NewMemberOfRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MemberOfRequestBuilder) {
    m := &MemberOfRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/memberOf{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewMemberOfRequestBuilder instantiates a new MemberOfRequestBuilder and sets the default values.
func NewMemberOfRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MemberOfRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMemberOfRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *MemberOfRequestBuilder) Count()(*id92750a7b3cc4ed2656a6e56bb755836906fdc5131ade6ac4b3ffd52bab76f3a.CountRequestBuilder) {
    return id92750a7b3cc4ed2656a6e56bb755836906fdc5131ade6ac4b3ffd52bab76f3a.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation roles that this service principal is a member of. HTTP Methods: GET Read-only. Nullable. Supports $expand.
func (m *MemberOfRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *MemberOfRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Device casts the previous resource to device.
func (m *MemberOfRequestBuilder) Device()(*i4454307873c5851e347787ba6ab2d3940ed9e7a55cf56ac268b1aac06400e6d7.DeviceRequestBuilder) {
    return i4454307873c5851e347787ba6ab2d3940ed9e7a55cf56ac268b1aac06400e6d7.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get roles that this service principal is a member of. HTTP Methods: GET Read-only. Nullable. Supports $expand.
func (m *MemberOfRequestBuilder) Get(ctx context.Context, requestConfiguration *MemberOfRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
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
func (m *MemberOfRequestBuilder) Group()(*iabf57bc1560ca5354a1c2b8ce106ad4b265aae6ce3ba5d26216db3b613dd81db.GroupRequestBuilder) {
    return iabf57bc1560ca5354a1c2b8ce106ad4b265aae6ce3ba5d26216db3b613dd81db.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *MemberOfRequestBuilder) OrgContact()(*i5e6f3053da49caf881450f09581c96acb41911545fcc8509c552944f2e192a67.OrgContactRequestBuilder) {
    return i5e6f3053da49caf881450f09581c96acb41911545fcc8509c552944f2e192a67.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *MemberOfRequestBuilder) ServicePrincipal()(*i0b27242c7dcbdd771786e0f446664593f09940a382381f2622fc6a1c3e458e07.ServicePrincipalRequestBuilder) {
    return i0b27242c7dcbdd771786e0f446664593f09940a382381f2622fc6a1c3e458e07.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *MemberOfRequestBuilder) User()(*ib70fe1873991f79b66e17a20a4d22ae11e48b9c6decf378d8edf9a0a0969a919.UserRequestBuilder) {
    return ib70fe1873991f79b66e17a20a4d22ae11e48b9c6decf378d8edf9a0a0969a919.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
