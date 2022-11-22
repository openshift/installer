package members

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i38118184414eceea71fc31f6be7fe2e203577cb78e16d0a0da1de0f7a7d9401c "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/group"
    i581742d9ba330e0969e1840b1f914a8bfcfd0f4f2c24d1e6282e1c491f718d3d "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/device"
    i7352778687f57b1dfa4fe6104377e6746ade4e2394f2c09fa95d0311bb329c51 "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/serviceprincipal"
    i74d3c02f44d332e8ca392295116b7d5cba5f372c1417890c6920f8a1d33cd2ea "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/orgcontact"
    i880bb46770a047cd20259c2ceb994c8031f774bcd4a705e353f11f9ba9a18001 "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/count"
    ib1a07bcd313a696054e8f98522bf52c4d946fdda956eda86dec0afb5df4ef19b "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/ref"
    if3d2be8e725c7d5a1a3d17d6ebd815c1e12924ac442933289a45b39c4fc975bc "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/application"
    ifa3a9d1a41a94314756760c02164cb049f8994e712fb73e0431b24bd7c65f277 "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/user"
)

// MembersRequestBuilder provides operations to manage the members property of the microsoft.graph.directoryRole entity.
type MembersRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// MembersRequestBuilderGetQueryParameters users that are members of this directory role. HTTP Methods: GET, POST, DELETE. Read-only. Nullable. Supports $expand.
type MembersRequestBuilderGetQueryParameters struct {
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
// MembersRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MembersRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MembersRequestBuilderGetQueryParameters
}
// Application casts the previous resource to application.
func (m *MembersRequestBuilder) Application()(*if3d2be8e725c7d5a1a3d17d6ebd815c1e12924ac442933289a45b39c4fc975bc.ApplicationRequestBuilder) {
    return if3d2be8e725c7d5a1a3d17d6ebd815c1e12924ac442933289a45b39c4fc975bc.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewMembersRequestBuilderInternal instantiates a new MembersRequestBuilder and sets the default values.
func NewMembersRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MembersRequestBuilder) {
    m := &MembersRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/directoryRoles/{directoryRole%2Did}/members{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewMembersRequestBuilder instantiates a new MembersRequestBuilder and sets the default values.
func NewMembersRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MembersRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMembersRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *MembersRequestBuilder) Count()(*i880bb46770a047cd20259c2ceb994c8031f774bcd4a705e353f11f9ba9a18001.CountRequestBuilder) {
    return i880bb46770a047cd20259c2ceb994c8031f774bcd4a705e353f11f9ba9a18001.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation users that are members of this directory role. HTTP Methods: GET, POST, DELETE. Read-only. Nullable. Supports $expand.
func (m *MembersRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *MembersRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *MembersRequestBuilder) Device()(*i581742d9ba330e0969e1840b1f914a8bfcfd0f4f2c24d1e6282e1c491f718d3d.DeviceRequestBuilder) {
    return i581742d9ba330e0969e1840b1f914a8bfcfd0f4f2c24d1e6282e1c491f718d3d.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get users that are members of this directory role. HTTP Methods: GET, POST, DELETE. Read-only. Nullable. Supports $expand.
func (m *MembersRequestBuilder) Get(ctx context.Context, requestConfiguration *MembersRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
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
func (m *MembersRequestBuilder) Group()(*i38118184414eceea71fc31f6be7fe2e203577cb78e16d0a0da1de0f7a7d9401c.GroupRequestBuilder) {
    return i38118184414eceea71fc31f6be7fe2e203577cb78e16d0a0da1de0f7a7d9401c.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *MembersRequestBuilder) OrgContact()(*i74d3c02f44d332e8ca392295116b7d5cba5f372c1417890c6920f8a1d33cd2ea.OrgContactRequestBuilder) {
    return i74d3c02f44d332e8ca392295116b7d5cba5f372c1417890c6920f8a1d33cd2ea.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of directoryRole entities.
func (m *MembersRequestBuilder) Ref()(*ib1a07bcd313a696054e8f98522bf52c4d946fdda956eda86dec0afb5df4ef19b.RefRequestBuilder) {
    return ib1a07bcd313a696054e8f98522bf52c4d946fdda956eda86dec0afb5df4ef19b.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *MembersRequestBuilder) ServicePrincipal()(*i7352778687f57b1dfa4fe6104377e6746ade4e2394f2c09fa95d0311bb329c51.ServicePrincipalRequestBuilder) {
    return i7352778687f57b1dfa4fe6104377e6746ade4e2394f2c09fa95d0311bb329c51.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *MembersRequestBuilder) User()(*ifa3a9d1a41a94314756760c02164cb049f8994e712fb73e0431b24bd7c65f277.UserRequestBuilder) {
    return ifa3a9d1a41a94314756760c02164cb049f8994e712fb73e0431b24bd7c65f277.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
