package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0872ea9e09fbbda9c62afcf7b96dbdbde1218b51a8d644cf53cad94477c4ede8 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/item/orgcontact"
    i16a651ddeec2b255863aed03760a6f1a52c220825a3e94f45f008d1e527ee112 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/item/serviceprincipal"
    i1b8b4c7b9ebd7e3570aba7eeca8de15502d0618bf40d50f2bd4cae1cccd53955 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/item/user"
    i601f4fc0c95eb86d53701c236b29e0f1b3f099536cd1ee724b22b193a4b0dbe2 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/item/device"
    if9ca4653e52194b3981433c08f3eb14c8a741544ae5f94a20172c8af95ce3eba "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/item/application"
    ifa71b059aa3e6e16e0c4dcbbf4c3e875691ccf9b05c3b2e67bb1da0f4139a45e "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/item/group"
)

// DirectoryObjectItemRequestBuilder provides operations to manage the transitiveMemberOf property of the microsoft.graph.orgContact entity.
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// DirectoryObjectItemRequestBuilderGetQueryParameters get transitiveMemberOf from contacts
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
func (m *DirectoryObjectItemRequestBuilder) Application()(*if9ca4653e52194b3981433c08f3eb14c8a741544ae5f94a20172c8af95ce3eba.ApplicationRequestBuilder) {
    return if9ca4653e52194b3981433c08f3eb14c8a741544ae5f94a20172c8af95ce3eba.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/contacts/{orgContact%2Did}/transitiveMemberOf/{directoryObject%2Did}{?%24select,%24expand}";
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
// CreateGetRequestInformation get transitiveMemberOf from contacts
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
func (m *DirectoryObjectItemRequestBuilder) Device()(*i601f4fc0c95eb86d53701c236b29e0f1b3f099536cd1ee724b22b193a4b0dbe2.DeviceRequestBuilder) {
    return i601f4fc0c95eb86d53701c236b29e0f1b3f099536cd1ee724b22b193a4b0dbe2.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get get transitiveMemberOf from contacts
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
func (m *DirectoryObjectItemRequestBuilder) Group()(*ifa71b059aa3e6e16e0c4dcbbf4c3e875691ccf9b05c3b2e67bb1da0f4139a45e.GroupRequestBuilder) {
    return ifa71b059aa3e6e16e0c4dcbbf4c3e875691ccf9b05c3b2e67bb1da0f4139a45e.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *DirectoryObjectItemRequestBuilder) OrgContact()(*i0872ea9e09fbbda9c62afcf7b96dbdbde1218b51a8d644cf53cad94477c4ede8.OrgContactRequestBuilder) {
    return i0872ea9e09fbbda9c62afcf7b96dbdbde1218b51a8d644cf53cad94477c4ede8.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*i16a651ddeec2b255863aed03760a6f1a52c220825a3e94f45f008d1e527ee112.ServicePrincipalRequestBuilder) {
    return i16a651ddeec2b255863aed03760a6f1a52c220825a3e94f45f008d1e527ee112.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*i1b8b4c7b9ebd7e3570aba7eeca8de15502d0618bf40d50f2bd4cae1cccd53955.UserRequestBuilder) {
    return i1b8b4c7b9ebd7e3570aba7eeca8de15502d0618bf40d50f2bd4cae1cccd53955.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
