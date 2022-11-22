package memberof

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i029f0c0199ac2f16a34b891aaa66007dfe970b188ffe3a660548e72edd2d54d2 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/memberof/user"
    i330026bb7097de85f21954ec1291006646b985660e599c8561c0374478e57b22 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/memberof/application"
    i3cbf60cf25a8a7d98e0770bddb788896ffb9ca9bb8a3868eae317d1c7a617957 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/memberof/group"
    i48933bc6898e7be908e9e1388eb77ac36efb2000246042ab9193911b8f5480d6 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/memberof/orgcontact"
    i48951b7b70cfd00dfafe4b76c663ba1f507d6517af65ce64e8f929e50e1dd1b1 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/memberof/device"
    i66f3dc3f4431afaf863e47991418d423a8053e9a717f18beaf3d9c6b4d1358b7 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/memberof/serviceprincipal"
    i773ad93098e0a6741e933b07b53fbe3fa3f43bcf03aae20c9debde1b3fe54b51 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/memberof/count"
)

// MemberOfRequestBuilder provides operations to manage the memberOf property of the microsoft.graph.orgContact entity.
type MemberOfRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// MemberOfRequestBuilderGetQueryParameters get memberOf from contacts
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
func (m *MemberOfRequestBuilder) Application()(*i330026bb7097de85f21954ec1291006646b985660e599c8561c0374478e57b22.ApplicationRequestBuilder) {
    return i330026bb7097de85f21954ec1291006646b985660e599c8561c0374478e57b22.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewMemberOfRequestBuilderInternal instantiates a new MemberOfRequestBuilder and sets the default values.
func NewMemberOfRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MemberOfRequestBuilder) {
    m := &MemberOfRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/contacts/{orgContact%2Did}/memberOf{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
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
func (m *MemberOfRequestBuilder) Count()(*i773ad93098e0a6741e933b07b53fbe3fa3f43bcf03aae20c9debde1b3fe54b51.CountRequestBuilder) {
    return i773ad93098e0a6741e933b07b53fbe3fa3f43bcf03aae20c9debde1b3fe54b51.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation get memberOf from contacts
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
func (m *MemberOfRequestBuilder) Device()(*i48951b7b70cfd00dfafe4b76c663ba1f507d6517af65ce64e8f929e50e1dd1b1.DeviceRequestBuilder) {
    return i48951b7b70cfd00dfafe4b76c663ba1f507d6517af65ce64e8f929e50e1dd1b1.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get get memberOf from contacts
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
func (m *MemberOfRequestBuilder) Group()(*i3cbf60cf25a8a7d98e0770bddb788896ffb9ca9bb8a3868eae317d1c7a617957.GroupRequestBuilder) {
    return i3cbf60cf25a8a7d98e0770bddb788896ffb9ca9bb8a3868eae317d1c7a617957.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *MemberOfRequestBuilder) OrgContact()(*i48933bc6898e7be908e9e1388eb77ac36efb2000246042ab9193911b8f5480d6.OrgContactRequestBuilder) {
    return i48933bc6898e7be908e9e1388eb77ac36efb2000246042ab9193911b8f5480d6.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *MemberOfRequestBuilder) ServicePrincipal()(*i66f3dc3f4431afaf863e47991418d423a8053e9a717f18beaf3d9c6b4d1358b7.ServicePrincipalRequestBuilder) {
    return i66f3dc3f4431afaf863e47991418d423a8053e9a717f18beaf3d9c6b4d1358b7.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *MemberOfRequestBuilder) User()(*i029f0c0199ac2f16a34b891aaa66007dfe970b188ffe3a660548e72edd2d54d2.UserRequestBuilder) {
    return i029f0c0199ac2f16a34b891aaa66007dfe970b188ffe3a660548e72edd2d54d2.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
