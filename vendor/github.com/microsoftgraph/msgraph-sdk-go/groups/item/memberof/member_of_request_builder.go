package memberof

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i09fdc0ce5de502cfd18ea8450148903fa60b8bcac2ef3671577f3c4a874aa087 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberof/device"
    i25b63a2e147019122ac599b45072d68693da7d9821882892d886b59e2a924bc4 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberof/serviceprincipal"
    i654e755c31785b309d44c74ff33af679221986148a6cb628cabee1ffb915871b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberof/group"
    i7fafc872420566275d04c117039045bf2d646aeea0756dbbfb4767af2fa513be "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberof/count"
    i9b3494fa6fb591c68bb2272a4c1d403a073e057cf410ff93534715438ba4f35c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberof/application"
    ie3917495d4d874e939f0e9018bd6d218be01f46da619449d6d21a5379830f0d5 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberof/orgcontact"
    ie4259814874e8f2389a0386897940f19391cf96421b0e0b392fac8f4924d0252 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberof/user"
)

// MemberOfRequestBuilder provides operations to manage the memberOf property of the microsoft.graph.group entity.
type MemberOfRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// MemberOfRequestBuilderGetQueryParameters groups that this group is a member of. HTTP Methods: GET (supported for all groups). Read-only. Nullable. Supports $expand.
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
func (m *MemberOfRequestBuilder) Application()(*i9b3494fa6fb591c68bb2272a4c1d403a073e057cf410ff93534715438ba4f35c.ApplicationRequestBuilder) {
    return i9b3494fa6fb591c68bb2272a4c1d403a073e057cf410ff93534715438ba4f35c.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewMemberOfRequestBuilderInternal instantiates a new MemberOfRequestBuilder and sets the default values.
func NewMemberOfRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MemberOfRequestBuilder) {
    m := &MemberOfRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/memberOf{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
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
func (m *MemberOfRequestBuilder) Count()(*i7fafc872420566275d04c117039045bf2d646aeea0756dbbfb4767af2fa513be.CountRequestBuilder) {
    return i7fafc872420566275d04c117039045bf2d646aeea0756dbbfb4767af2fa513be.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation groups that this group is a member of. HTTP Methods: GET (supported for all groups). Read-only. Nullable. Supports $expand.
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
func (m *MemberOfRequestBuilder) Device()(*i09fdc0ce5de502cfd18ea8450148903fa60b8bcac2ef3671577f3c4a874aa087.DeviceRequestBuilder) {
    return i09fdc0ce5de502cfd18ea8450148903fa60b8bcac2ef3671577f3c4a874aa087.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get groups that this group is a member of. HTTP Methods: GET (supported for all groups). Read-only. Nullable. Supports $expand.
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
func (m *MemberOfRequestBuilder) Group()(*i654e755c31785b309d44c74ff33af679221986148a6cb628cabee1ffb915871b.GroupRequestBuilder) {
    return i654e755c31785b309d44c74ff33af679221986148a6cb628cabee1ffb915871b.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *MemberOfRequestBuilder) OrgContact()(*ie3917495d4d874e939f0e9018bd6d218be01f46da619449d6d21a5379830f0d5.OrgContactRequestBuilder) {
    return ie3917495d4d874e939f0e9018bd6d218be01f46da619449d6d21a5379830f0d5.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *MemberOfRequestBuilder) ServicePrincipal()(*i25b63a2e147019122ac599b45072d68693da7d9821882892d886b59e2a924bc4.ServicePrincipalRequestBuilder) {
    return i25b63a2e147019122ac599b45072d68693da7d9821882892d886b59e2a924bc4.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *MemberOfRequestBuilder) User()(*ie4259814874e8f2389a0386897940f19391cf96421b0e0b392fac8f4924d0252.UserRequestBuilder) {
    return ie4259814874e8f2389a0386897940f19391cf96421b0e0b392fac8f4924d0252.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
