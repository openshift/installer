package transitivememberof

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i18c6cfdf5bb78c2cd84f1ed3a2eaf2d75f610031f0030a212a1b60fb8df92594 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/user"
    i362972e85ac31087cef0b95d1089fe2d7d42e71352bb05120d64c0f28b924248 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/serviceprincipal"
    i364b4877f1c0aed94616371f3fa7a1a78316f9b0f9ebd2347e23125af1dc641a "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/group"
    i45f92b5ebb3c301d1219cd9c21eb9732b1bae2c6f9a3113c0085206ccdf21542 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/orgcontact"
    i5d449587467d81f36dc1ab393ed391230ec7f08d20a2250c93a62391e94e469a "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/device"
    i736d61974be9aedaa8009017ad19f62f35625a45dc3a5a8b56b15b583838711f "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/application"
    ief89c212bb34a0eb54fe4cbb1db9e342fecc50d23976834d4a222dd9959b7c71 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item/transitivememberof/count"
)

// TransitiveMemberOfRequestBuilder provides operations to manage the transitiveMemberOf property of the microsoft.graph.orgContact entity.
type TransitiveMemberOfRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// TransitiveMemberOfRequestBuilderGetQueryParameters get transitiveMemberOf from contacts
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
func (m *TransitiveMemberOfRequestBuilder) Application()(*i736d61974be9aedaa8009017ad19f62f35625a45dc3a5a8b56b15b583838711f.ApplicationRequestBuilder) {
    return i736d61974be9aedaa8009017ad19f62f35625a45dc3a5a8b56b15b583838711f.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewTransitiveMemberOfRequestBuilderInternal instantiates a new TransitiveMemberOfRequestBuilder and sets the default values.
func NewTransitiveMemberOfRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TransitiveMemberOfRequestBuilder) {
    m := &TransitiveMemberOfRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/contacts/{orgContact%2Did}/transitiveMemberOf{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
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
func (m *TransitiveMemberOfRequestBuilder) Count()(*ief89c212bb34a0eb54fe4cbb1db9e342fecc50d23976834d4a222dd9959b7c71.CountRequestBuilder) {
    return ief89c212bb34a0eb54fe4cbb1db9e342fecc50d23976834d4a222dd9959b7c71.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation get transitiveMemberOf from contacts
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
func (m *TransitiveMemberOfRequestBuilder) Device()(*i5d449587467d81f36dc1ab393ed391230ec7f08d20a2250c93a62391e94e469a.DeviceRequestBuilder) {
    return i5d449587467d81f36dc1ab393ed391230ec7f08d20a2250c93a62391e94e469a.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get get transitiveMemberOf from contacts
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
func (m *TransitiveMemberOfRequestBuilder) Group()(*i364b4877f1c0aed94616371f3fa7a1a78316f9b0f9ebd2347e23125af1dc641a.GroupRequestBuilder) {
    return i364b4877f1c0aed94616371f3fa7a1a78316f9b0f9ebd2347e23125af1dc641a.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *TransitiveMemberOfRequestBuilder) OrgContact()(*i45f92b5ebb3c301d1219cd9c21eb9732b1bae2c6f9a3113c0085206ccdf21542.OrgContactRequestBuilder) {
    return i45f92b5ebb3c301d1219cd9c21eb9732b1bae2c6f9a3113c0085206ccdf21542.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *TransitiveMemberOfRequestBuilder) ServicePrincipal()(*i362972e85ac31087cef0b95d1089fe2d7d42e71352bb05120d64c0f28b924248.ServicePrincipalRequestBuilder) {
    return i362972e85ac31087cef0b95d1089fe2d7d42e71352bb05120d64c0f28b924248.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *TransitiveMemberOfRequestBuilder) User()(*i18c6cfdf5bb78c2cd84f1ed3a2eaf2d75f610031f0030a212a1b60fb8df92594.UserRequestBuilder) {
    return i18c6cfdf5bb78c2cd84f1ed3a2eaf2d75f610031f0030a212a1b60fb8df92594.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
