package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i08455bff0d3a507a966599516654d5d8197a83a8e12cbe332bb57d758e84a1e4 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivememberof/item/orgcontact"
    i9d23f424a62d000bb302a0fb1e0b7a0edc5d4cc85eaaac4be253b4b826e4c7c5 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivememberof/item/serviceprincipal"
    ic4c58f02936e92cbe39bb8bcf8716daac685f9af6344908e910f3d66c83a157a "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivememberof/item/device"
    idb62748373c010ffa1d69fb54e0d1a1fac98da927f2539931b876248d2e908fe "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivememberof/item/group"
    if7a7c47ae89f8e51e9f52344cd749f14c46dbc405e95cb638e35b4573967a810 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivememberof/item/application"
    ifa847279b4905d792fd6782676bf7969090332a2a77fdde2cb1a5526f0c95d17 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivememberof/item/user"
)

// DirectoryObjectItemRequestBuilder provides operations to manage the transitiveMemberOf property of the microsoft.graph.group entity.
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// DirectoryObjectItemRequestBuilderGetQueryParameters the groups that a group is a member of, either directly and through nested membership. Nullable.
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
func (m *DirectoryObjectItemRequestBuilder) Application()(*if7a7c47ae89f8e51e9f52344cd749f14c46dbc405e95cb638e35b4573967a810.ApplicationRequestBuilder) {
    return if7a7c47ae89f8e51e9f52344cd749f14c46dbc405e95cb638e35b4573967a810.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/transitiveMemberOf/{directoryObject%2Did}{?%24select,%24expand}";
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
// CreateGetRequestInformation the groups that a group is a member of, either directly and through nested membership. Nullable.
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
func (m *DirectoryObjectItemRequestBuilder) Device()(*ic4c58f02936e92cbe39bb8bcf8716daac685f9af6344908e910f3d66c83a157a.DeviceRequestBuilder) {
    return ic4c58f02936e92cbe39bb8bcf8716daac685f9af6344908e910f3d66c83a157a.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the groups that a group is a member of, either directly and through nested membership. Nullable.
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
func (m *DirectoryObjectItemRequestBuilder) Group()(*idb62748373c010ffa1d69fb54e0d1a1fac98da927f2539931b876248d2e908fe.GroupRequestBuilder) {
    return idb62748373c010ffa1d69fb54e0d1a1fac98da927f2539931b876248d2e908fe.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *DirectoryObjectItemRequestBuilder) OrgContact()(*i08455bff0d3a507a966599516654d5d8197a83a8e12cbe332bb57d758e84a1e4.OrgContactRequestBuilder) {
    return i08455bff0d3a507a966599516654d5d8197a83a8e12cbe332bb57d758e84a1e4.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*i9d23f424a62d000bb302a0fb1e0b7a0edc5d4cc85eaaac4be253b4b826e4c7c5.ServicePrincipalRequestBuilder) {
    return i9d23f424a62d000bb302a0fb1e0b7a0edc5d4cc85eaaac4be253b4b826e4c7c5.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*ifa847279b4905d792fd6782676bf7969090332a2a77fdde2cb1a5526f0c95d17.UserRequestBuilder) {
    return ifa847279b4905d792fd6782676bf7969090332a2a77fdde2cb1a5526f0c95d17.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
