package transitivememberof

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0bc63be75ac303f6b43c28e3a061ea60b28080039abfc589d73deef851e766f1 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/transitivememberof/group"
    i2f498bf25a4d10b00a1bec8d5c4950a3b0b43e362fac0ca28c09772bb3e82b65 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/transitivememberof/application"
    i87c00ed6154c2f13639757f52ad6bacc6a594047248b2f3fd3c1082456701943 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/transitivememberof/device"
    i9745ffa3bef22614df23fff3705b156c1327099a3b1528dd6d0149c1d4ce997f "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/transitivememberof/user"
    iafff1194bef580888235900aac3172d83f360aaf8ec027e4549280d060f7de5d "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/transitivememberof/orgcontact"
    iddd86adf7641cc14f00b8231c8ec278181a76a9871fba9fee14b8bfee85107c2 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/transitivememberof/serviceprincipal"
    ie355b6a2fb2db274cfe967938d840245ee5fc7981beda26a1e698879b15da371 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/transitivememberof/count"
)

// TransitiveMemberOfRequestBuilder provides operations to manage the transitiveMemberOf property of the microsoft.graph.servicePrincipal entity.
type TransitiveMemberOfRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// TransitiveMemberOfRequestBuilderGetQueryParameters get transitiveMemberOf from servicePrincipals
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
func (m *TransitiveMemberOfRequestBuilder) Application()(*i2f498bf25a4d10b00a1bec8d5c4950a3b0b43e362fac0ca28c09772bb3e82b65.ApplicationRequestBuilder) {
    return i2f498bf25a4d10b00a1bec8d5c4950a3b0b43e362fac0ca28c09772bb3e82b65.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewTransitiveMemberOfRequestBuilderInternal instantiates a new TransitiveMemberOfRequestBuilder and sets the default values.
func NewTransitiveMemberOfRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TransitiveMemberOfRequestBuilder) {
    m := &TransitiveMemberOfRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/transitiveMemberOf{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
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
func (m *TransitiveMemberOfRequestBuilder) Count()(*ie355b6a2fb2db274cfe967938d840245ee5fc7981beda26a1e698879b15da371.CountRequestBuilder) {
    return ie355b6a2fb2db274cfe967938d840245ee5fc7981beda26a1e698879b15da371.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation get transitiveMemberOf from servicePrincipals
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
func (m *TransitiveMemberOfRequestBuilder) Device()(*i87c00ed6154c2f13639757f52ad6bacc6a594047248b2f3fd3c1082456701943.DeviceRequestBuilder) {
    return i87c00ed6154c2f13639757f52ad6bacc6a594047248b2f3fd3c1082456701943.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get get transitiveMemberOf from servicePrincipals
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
func (m *TransitiveMemberOfRequestBuilder) Group()(*i0bc63be75ac303f6b43c28e3a061ea60b28080039abfc589d73deef851e766f1.GroupRequestBuilder) {
    return i0bc63be75ac303f6b43c28e3a061ea60b28080039abfc589d73deef851e766f1.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *TransitiveMemberOfRequestBuilder) OrgContact()(*iafff1194bef580888235900aac3172d83f360aaf8ec027e4549280d060f7de5d.OrgContactRequestBuilder) {
    return iafff1194bef580888235900aac3172d83f360aaf8ec027e4549280d060f7de5d.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *TransitiveMemberOfRequestBuilder) ServicePrincipal()(*iddd86adf7641cc14f00b8231c8ec278181a76a9871fba9fee14b8bfee85107c2.ServicePrincipalRequestBuilder) {
    return iddd86adf7641cc14f00b8231c8ec278181a76a9871fba9fee14b8bfee85107c2.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *TransitiveMemberOfRequestBuilder) User()(*i9745ffa3bef22614df23fff3705b156c1327099a3b1528dd6d0149c1d4ce997f.UserRequestBuilder) {
    return i9745ffa3bef22614df23fff3705b156c1327099a3b1528dd6d0149c1d4ce997f.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
