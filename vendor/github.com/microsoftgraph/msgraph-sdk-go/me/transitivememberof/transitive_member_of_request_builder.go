package transitivememberof

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i5f7e61b0c1f178c96c67c89db2e54884580eb9debf5da111645c1b835be3b1a0 "github.com/microsoftgraph/msgraph-sdk-go/me/transitivememberof/application"
    i8f2c143d91124cd4159b7ae4ace9ae2ef48b787d52f6b2388dbda3d565938152 "github.com/microsoftgraph/msgraph-sdk-go/me/transitivememberof/serviceprincipal"
    ia339a388059a4a725d09806d482efb171a74d285f37d4f7efc597b15f5cbb4d0 "github.com/microsoftgraph/msgraph-sdk-go/me/transitivememberof/device"
    iaa783998cb72105d1e7f25d34cd9821ac9ccf1022bb97ec0821c0951413d5834 "github.com/microsoftgraph/msgraph-sdk-go/me/transitivememberof/group"
    icb9716a67683385703764d49f2b2ee7a872082397a36ebefe294533f7c880d2d "github.com/microsoftgraph/msgraph-sdk-go/me/transitivememberof/count"
    icec9d8e0f9ef512c158ef8ecddba7af6a43bd42aebd108dcf3651f9a615652e4 "github.com/microsoftgraph/msgraph-sdk-go/me/transitivememberof/orgcontact"
    if1b8b394295215d53f03821ff96d888ba7244110f9e0d6a033e448745d85790b "github.com/microsoftgraph/msgraph-sdk-go/me/transitivememberof/user"
)

// TransitiveMemberOfRequestBuilder provides operations to manage the transitiveMemberOf property of the microsoft.graph.user entity.
type TransitiveMemberOfRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// TransitiveMemberOfRequestBuilderGetQueryParameters the groups, including nested groups, and directory roles that a user is a member of. Nullable.
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
func (m *TransitiveMemberOfRequestBuilder) Application()(*i5f7e61b0c1f178c96c67c89db2e54884580eb9debf5da111645c1b835be3b1a0.ApplicationRequestBuilder) {
    return i5f7e61b0c1f178c96c67c89db2e54884580eb9debf5da111645c1b835be3b1a0.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewTransitiveMemberOfRequestBuilderInternal instantiates a new TransitiveMemberOfRequestBuilder and sets the default values.
func NewTransitiveMemberOfRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TransitiveMemberOfRequestBuilder) {
    m := &TransitiveMemberOfRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/transitiveMemberOf{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
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
func (m *TransitiveMemberOfRequestBuilder) Count()(*icb9716a67683385703764d49f2b2ee7a872082397a36ebefe294533f7c880d2d.CountRequestBuilder) {
    return icb9716a67683385703764d49f2b2ee7a872082397a36ebefe294533f7c880d2d.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation the groups, including nested groups, and directory roles that a user is a member of. Nullable.
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
func (m *TransitiveMemberOfRequestBuilder) Device()(*ia339a388059a4a725d09806d482efb171a74d285f37d4f7efc597b15f5cbb4d0.DeviceRequestBuilder) {
    return ia339a388059a4a725d09806d482efb171a74d285f37d4f7efc597b15f5cbb4d0.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the groups, including nested groups, and directory roles that a user is a member of. Nullable.
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
func (m *TransitiveMemberOfRequestBuilder) Group()(*iaa783998cb72105d1e7f25d34cd9821ac9ccf1022bb97ec0821c0951413d5834.GroupRequestBuilder) {
    return iaa783998cb72105d1e7f25d34cd9821ac9ccf1022bb97ec0821c0951413d5834.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *TransitiveMemberOfRequestBuilder) OrgContact()(*icec9d8e0f9ef512c158ef8ecddba7af6a43bd42aebd108dcf3651f9a615652e4.OrgContactRequestBuilder) {
    return icec9d8e0f9ef512c158ef8ecddba7af6a43bd42aebd108dcf3651f9a615652e4.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *TransitiveMemberOfRequestBuilder) ServicePrincipal()(*i8f2c143d91124cd4159b7ae4ace9ae2ef48b787d52f6b2388dbda3d565938152.ServicePrincipalRequestBuilder) {
    return i8f2c143d91124cd4159b7ae4ace9ae2ef48b787d52f6b2388dbda3d565938152.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *TransitiveMemberOfRequestBuilder) User()(*if1b8b394295215d53f03821ff96d888ba7244110f9e0d6a033e448745d85790b.UserRequestBuilder) {
    return if1b8b394295215d53f03821ff96d888ba7244110f9e0d6a033e448745d85790b.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
