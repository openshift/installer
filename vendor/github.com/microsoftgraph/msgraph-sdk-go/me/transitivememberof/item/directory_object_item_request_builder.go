package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0541f7430c345354a67654f250a38f52a8e99908a6cf053140b86a973d42a340 "github.com/microsoftgraph/msgraph-sdk-go/me/transitivememberof/item/user"
    i05bd1f37bb79481462e80faaa38171a9f00c55f1cdf782d67a823df980ce60cb "github.com/microsoftgraph/msgraph-sdk-go/me/transitivememberof/item/orgcontact"
    i184d1fdf9b1340963d45be83bf02245925692bb8de937a4f1b75662689effa45 "github.com/microsoftgraph/msgraph-sdk-go/me/transitivememberof/item/device"
    i444360d520b8d9290c6fb7f13eae76d50e8ef2db9bc8f6691e127a2a7384f62e "github.com/microsoftgraph/msgraph-sdk-go/me/transitivememberof/item/group"
    i8b4c0c4d3d5eb767591fdfefafe06b6da5298db9d3530e7f2364d660fe12a091 "github.com/microsoftgraph/msgraph-sdk-go/me/transitivememberof/item/serviceprincipal"
    i9d4a5c081d26c30f5b2071d4b94c0845d6619169e0e9d081a42347267c176ee4 "github.com/microsoftgraph/msgraph-sdk-go/me/transitivememberof/item/application"
)

// DirectoryObjectItemRequestBuilder provides operations to manage the transitiveMemberOf property of the microsoft.graph.user entity.
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// DirectoryObjectItemRequestBuilderGetQueryParameters the groups, including nested groups, and directory roles that a user is a member of. Nullable.
type DirectoryObjectItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DirectoryObjectItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DirectoryObjectItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DirectoryObjectItemRequestBuilderGetQueryParameters
}
// Application casts the previous resource to application.
func (m *DirectoryObjectItemRequestBuilder) Application()(*i9d4a5c081d26c30f5b2071d4b94c0845d6619169e0e9d081a42347267c176ee4.ApplicationRequestBuilder) {
    return i9d4a5c081d26c30f5b2071d4b94c0845d6619169e0e9d081a42347267c176ee4.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/transitiveMemberOf/{directoryObject%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewDirectoryObjectItemRequestBuilder instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDirectoryObjectItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation the groups, including nested groups, and directory roles that a user is a member of. Nullable.
func (m *DirectoryObjectItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *DirectoryObjectItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *DirectoryObjectItemRequestBuilder) Device()(*i184d1fdf9b1340963d45be83bf02245925692bb8de937a4f1b75662689effa45.DeviceRequestBuilder) {
    return i184d1fdf9b1340963d45be83bf02245925692bb8de937a4f1b75662689effa45.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the groups, including nested groups, and directory roles that a user is a member of. Nullable.
func (m *DirectoryObjectItemRequestBuilder) Get(ctx context.Context, requestConfiguration *DirectoryObjectItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryObjectFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable), nil
}
// Group casts the previous resource to group.
func (m *DirectoryObjectItemRequestBuilder) Group()(*i444360d520b8d9290c6fb7f13eae76d50e8ef2db9bc8f6691e127a2a7384f62e.GroupRequestBuilder) {
    return i444360d520b8d9290c6fb7f13eae76d50e8ef2db9bc8f6691e127a2a7384f62e.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *DirectoryObjectItemRequestBuilder) OrgContact()(*i05bd1f37bb79481462e80faaa38171a9f00c55f1cdf782d67a823df980ce60cb.OrgContactRequestBuilder) {
    return i05bd1f37bb79481462e80faaa38171a9f00c55f1cdf782d67a823df980ce60cb.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*i8b4c0c4d3d5eb767591fdfefafe06b6da5298db9d3530e7f2364d660fe12a091.ServicePrincipalRequestBuilder) {
    return i8b4c0c4d3d5eb767591fdfefafe06b6da5298db9d3530e7f2364d660fe12a091.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*i0541f7430c345354a67654f250a38f52a8e99908a6cf053140b86a973d42a340.UserRequestBuilder) {
    return i0541f7430c345354a67654f250a38f52a8e99908a6cf053140b86a973d42a340.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
