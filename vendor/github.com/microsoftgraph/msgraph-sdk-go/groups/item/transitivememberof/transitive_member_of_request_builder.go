package transitivememberof

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0b8df292d0745698570dc5f21fc1137fea3745e60baccff4e0794cab160d6ada "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivememberof/orgcontact"
    i52b285a4b5dfc67eeba0538826a0c75cb8a567fa8b95f79c05235d6bc79a7bd4 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivememberof/user"
    i792c82fbf0df7e5f435214a5dfcc562b38233a1394e54b7bcc3a0e4d1ba58c49 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivememberof/serviceprincipal"
    ia55d293da10d83d4490c67b3b55f4a9c2e1f92b9f194a0b82b7e304bea4eaf82 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivememberof/group"
    ib2fe019740c4ba2e9995a5469a2fc5460e1a45b6dda90a1c11d6808de7ec8150 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivememberof/device"
    ic718df616f470b9baeaf5351ad4d1e0dd4982519c7ddc8d47070ddcb2e8a0a3e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivememberof/application"
    id53e208bb7aff4ca8bed62c924cd8e191407770fed99ba105c23c454af7f7287 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivememberof/count"
)

// TransitiveMemberOfRequestBuilder provides operations to manage the transitiveMemberOf property of the microsoft.graph.group entity.
type TransitiveMemberOfRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// TransitiveMemberOfRequestBuilderGetQueryParameters the groups that a group is a member of, either directly and through nested membership. Nullable.
type TransitiveMemberOfRequestBuilderGetQueryParameters struct {
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
// TransitiveMemberOfRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TransitiveMemberOfRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *TransitiveMemberOfRequestBuilderGetQueryParameters
}
// Application casts the previous resource to application.
func (m *TransitiveMemberOfRequestBuilder) Application()(*ic718df616f470b9baeaf5351ad4d1e0dd4982519c7ddc8d47070ddcb2e8a0a3e.ApplicationRequestBuilder) {
    return ic718df616f470b9baeaf5351ad4d1e0dd4982519c7ddc8d47070ddcb2e8a0a3e.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewTransitiveMemberOfRequestBuilderInternal instantiates a new TransitiveMemberOfRequestBuilder and sets the default values.
func NewTransitiveMemberOfRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TransitiveMemberOfRequestBuilder) {
    m := &TransitiveMemberOfRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/transitiveMemberOf{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewTransitiveMemberOfRequestBuilder instantiates a new TransitiveMemberOfRequestBuilder and sets the default values.
func NewTransitiveMemberOfRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TransitiveMemberOfRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewTransitiveMemberOfRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *TransitiveMemberOfRequestBuilder) Count()(*id53e208bb7aff4ca8bed62c924cd8e191407770fed99ba105c23c454af7f7287.CountRequestBuilder) {
    return id53e208bb7aff4ca8bed62c924cd8e191407770fed99ba105c23c454af7f7287.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation the groups that a group is a member of, either directly and through nested membership. Nullable.
func (m *TransitiveMemberOfRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *TransitiveMemberOfRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *TransitiveMemberOfRequestBuilder) Device()(*ib2fe019740c4ba2e9995a5469a2fc5460e1a45b6dda90a1c11d6808de7ec8150.DeviceRequestBuilder) {
    return ib2fe019740c4ba2e9995a5469a2fc5460e1a45b6dda90a1c11d6808de7ec8150.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the groups that a group is a member of, either directly and through nested membership. Nullable.
func (m *TransitiveMemberOfRequestBuilder) Get(ctx context.Context, requestConfiguration *TransitiveMemberOfRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
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
func (m *TransitiveMemberOfRequestBuilder) Group()(*ia55d293da10d83d4490c67b3b55f4a9c2e1f92b9f194a0b82b7e304bea4eaf82.GroupRequestBuilder) {
    return ia55d293da10d83d4490c67b3b55f4a9c2e1f92b9f194a0b82b7e304bea4eaf82.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *TransitiveMemberOfRequestBuilder) OrgContact()(*i0b8df292d0745698570dc5f21fc1137fea3745e60baccff4e0794cab160d6ada.OrgContactRequestBuilder) {
    return i0b8df292d0745698570dc5f21fc1137fea3745e60baccff4e0794cab160d6ada.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *TransitiveMemberOfRequestBuilder) ServicePrincipal()(*i792c82fbf0df7e5f435214a5dfcc562b38233a1394e54b7bcc3a0e4d1ba58c49.ServicePrincipalRequestBuilder) {
    return i792c82fbf0df7e5f435214a5dfcc562b38233a1394e54b7bcc3a0e4d1ba58c49.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *TransitiveMemberOfRequestBuilder) User()(*i52b285a4b5dfc67eeba0538826a0c75cb8a567fa8b95f79c05235d6bc79a7bd4.UserRequestBuilder) {
    return i52b285a4b5dfc67eeba0538826a0c75cb8a567fa8b95f79c05235d6bc79a7bd4.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
