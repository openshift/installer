package memberof

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0823b7f9e0620e3c10ed0e03c860d61ba82085fd0cbc86ee37b55c9913603587 "github.com/microsoftgraph/msgraph-sdk-go/me/memberof/orgcontact"
    i28b0877d46e0fef9eb80669bbc452b18b98035cd14fa3ebd9179ffeceb9fb5f8 "github.com/microsoftgraph/msgraph-sdk-go/me/memberof/application"
    i646a4ac5bc99048228608d332efbb2ccdde4b73f5d92de5efa03f625af535113 "github.com/microsoftgraph/msgraph-sdk-go/me/memberof/user"
    i70d1677faec042643c5fd925d8ed7167dc66052b34a9e196d3fe1ffb58757f4a "github.com/microsoftgraph/msgraph-sdk-go/me/memberof/device"
    i8984270e0001edcb14427a5116429fd6f711943c3ae6facd56768c4e791daed8 "github.com/microsoftgraph/msgraph-sdk-go/me/memberof/group"
    i95775dfc90a74159759a9ac3f2ca610131c365a6945f7ae481a376a79ebf9573 "github.com/microsoftgraph/msgraph-sdk-go/me/memberof/count"
    ic31f98d1fc9879e7abb6cf5d636d24a5626dcf2f368d97151d4b8d470d24e50d "github.com/microsoftgraph/msgraph-sdk-go/me/memberof/serviceprincipal"
)

// MemberOfRequestBuilder provides operations to manage the memberOf property of the microsoft.graph.user entity.
type MemberOfRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// MemberOfRequestBuilderGetQueryParameters the groups and directory roles that the user is a member of. Read-only. Nullable. Supports $expand.
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
func (m *MemberOfRequestBuilder) Application()(*i28b0877d46e0fef9eb80669bbc452b18b98035cd14fa3ebd9179ffeceb9fb5f8.ApplicationRequestBuilder) {
    return i28b0877d46e0fef9eb80669bbc452b18b98035cd14fa3ebd9179ffeceb9fb5f8.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewMemberOfRequestBuilderInternal instantiates a new MemberOfRequestBuilder and sets the default values.
func NewMemberOfRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MemberOfRequestBuilder) {
    m := &MemberOfRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/memberOf{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
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
func (m *MemberOfRequestBuilder) Count()(*i95775dfc90a74159759a9ac3f2ca610131c365a6945f7ae481a376a79ebf9573.CountRequestBuilder) {
    return i95775dfc90a74159759a9ac3f2ca610131c365a6945f7ae481a376a79ebf9573.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation the groups and directory roles that the user is a member of. Read-only. Nullable. Supports $expand.
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
func (m *MemberOfRequestBuilder) Device()(*i70d1677faec042643c5fd925d8ed7167dc66052b34a9e196d3fe1ffb58757f4a.DeviceRequestBuilder) {
    return i70d1677faec042643c5fd925d8ed7167dc66052b34a9e196d3fe1ffb58757f4a.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the groups and directory roles that the user is a member of. Read-only. Nullable. Supports $expand.
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
func (m *MemberOfRequestBuilder) Group()(*i8984270e0001edcb14427a5116429fd6f711943c3ae6facd56768c4e791daed8.GroupRequestBuilder) {
    return i8984270e0001edcb14427a5116429fd6f711943c3ae6facd56768c4e791daed8.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *MemberOfRequestBuilder) OrgContact()(*i0823b7f9e0620e3c10ed0e03c860d61ba82085fd0cbc86ee37b55c9913603587.OrgContactRequestBuilder) {
    return i0823b7f9e0620e3c10ed0e03c860d61ba82085fd0cbc86ee37b55c9913603587.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *MemberOfRequestBuilder) ServicePrincipal()(*ic31f98d1fc9879e7abb6cf5d636d24a5626dcf2f368d97151d4b8d470d24e50d.ServicePrincipalRequestBuilder) {
    return ic31f98d1fc9879e7abb6cf5d636d24a5626dcf2f368d97151d4b8d470d24e50d.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *MemberOfRequestBuilder) User()(*i646a4ac5bc99048228608d332efbb2ccdde4b73f5d92de5efa03f625af535113.UserRequestBuilder) {
    return i646a4ac5bc99048228608d332efbb2ccdde4b73f5d92de5efa03f625af535113.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
