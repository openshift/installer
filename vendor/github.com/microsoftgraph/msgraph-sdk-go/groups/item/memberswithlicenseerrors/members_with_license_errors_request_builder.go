package memberswithlicenseerrors

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0dab4ab5d03cdb1e2bc2146c8f8f393369639cf2951df17099c75880b3405169 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberswithlicenseerrors/group"
    i38d5613daf6a74842f466cb5f7e1d5e1b49776cadcc1ccffe3f985c86669b1ba "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberswithlicenseerrors/user"
    i44b96843de9edb7831feeb7753f66dad8edcd1bbfaeea1ff3b97835638702ed8 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberswithlicenseerrors/device"
    i9eee8ca196107e21e220b940811218fe3f0c4bfce89d18b327b903b1f1fa451c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberswithlicenseerrors/application"
    ia0ddb948c0200f3efec2048b85c5e412b832549e67625f1bde2d7930f12a7216 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberswithlicenseerrors/orgcontact"
    ia742a37f75ffa85a15c52758b54a8cfbacdcae12597cc38b77fd0f412aaa1e32 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberswithlicenseerrors/count"
    iddc54b8e85093999a061b74455c42139479196ec70c7ff2077c4332c3635da87 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberswithlicenseerrors/serviceprincipal"
)

// MembersWithLicenseErrorsRequestBuilder provides operations to manage the membersWithLicenseErrors property of the microsoft.graph.group entity.
type MembersWithLicenseErrorsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// MembersWithLicenseErrorsRequestBuilderGetQueryParameters a list of group members with license errors from this group-based license assignment. Read-only.
type MembersWithLicenseErrorsRequestBuilderGetQueryParameters struct {
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
// MembersWithLicenseErrorsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MembersWithLicenseErrorsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MembersWithLicenseErrorsRequestBuilderGetQueryParameters
}
// Application casts the previous resource to application.
func (m *MembersWithLicenseErrorsRequestBuilder) Application()(*i9eee8ca196107e21e220b940811218fe3f0c4bfce89d18b327b903b1f1fa451c.ApplicationRequestBuilder) {
    return i9eee8ca196107e21e220b940811218fe3f0c4bfce89d18b327b903b1f1fa451c.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewMembersWithLicenseErrorsRequestBuilderInternal instantiates a new MembersWithLicenseErrorsRequestBuilder and sets the default values.
func NewMembersWithLicenseErrorsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MembersWithLicenseErrorsRequestBuilder) {
    m := &MembersWithLicenseErrorsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/membersWithLicenseErrors{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewMembersWithLicenseErrorsRequestBuilder instantiates a new MembersWithLicenseErrorsRequestBuilder and sets the default values.
func NewMembersWithLicenseErrorsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MembersWithLicenseErrorsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMembersWithLicenseErrorsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *MembersWithLicenseErrorsRequestBuilder) Count()(*ia742a37f75ffa85a15c52758b54a8cfbacdcae12597cc38b77fd0f412aaa1e32.CountRequestBuilder) {
    return ia742a37f75ffa85a15c52758b54a8cfbacdcae12597cc38b77fd0f412aaa1e32.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation a list of group members with license errors from this group-based license assignment. Read-only.
func (m *MembersWithLicenseErrorsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *MembersWithLicenseErrorsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *MembersWithLicenseErrorsRequestBuilder) Device()(*i44b96843de9edb7831feeb7753f66dad8edcd1bbfaeea1ff3b97835638702ed8.DeviceRequestBuilder) {
    return i44b96843de9edb7831feeb7753f66dad8edcd1bbfaeea1ff3b97835638702ed8.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get a list of group members with license errors from this group-based license assignment. Read-only.
func (m *MembersWithLicenseErrorsRequestBuilder) Get(ctx context.Context, requestConfiguration *MembersWithLicenseErrorsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
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
func (m *MembersWithLicenseErrorsRequestBuilder) Group()(*i0dab4ab5d03cdb1e2bc2146c8f8f393369639cf2951df17099c75880b3405169.GroupRequestBuilder) {
    return i0dab4ab5d03cdb1e2bc2146c8f8f393369639cf2951df17099c75880b3405169.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *MembersWithLicenseErrorsRequestBuilder) OrgContact()(*ia0ddb948c0200f3efec2048b85c5e412b832549e67625f1bde2d7930f12a7216.OrgContactRequestBuilder) {
    return ia0ddb948c0200f3efec2048b85c5e412b832549e67625f1bde2d7930f12a7216.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *MembersWithLicenseErrorsRequestBuilder) ServicePrincipal()(*iddc54b8e85093999a061b74455c42139479196ec70c7ff2077c4332c3635da87.ServicePrincipalRequestBuilder) {
    return iddc54b8e85093999a061b74455c42139479196ec70c7ff2077c4332c3635da87.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *MembersWithLicenseErrorsRequestBuilder) User()(*i38d5613daf6a74842f466cb5f7e1d5e1b49776cadcc1ccffe3f985c86669b1ba.UserRequestBuilder) {
    return i38d5613daf6a74842f466cb5f7e1d5e1b49776cadcc1ccffe3f985c86669b1ba.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
