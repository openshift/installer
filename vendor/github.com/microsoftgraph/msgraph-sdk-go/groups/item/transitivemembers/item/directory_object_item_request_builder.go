package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i2ae8f05afc4ab7755eda7ac570865f043c6a61aeb68d9792631c6e7297790f17 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivemembers/item/group"
    i350aaf1b458d5143e8c2c0377a329c7309406d5c18dbea162e2e59039870e589 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivemembers/item/orgcontact"
    i709d4a15bc3a3bf78e9c6210bf45793f6eaf5531c0fcb372ecfb169e1d16f617 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivemembers/item/serviceprincipal"
    ic3bf19ec033cea06ca57ed6e10c7103f65b6eac7c428254155d46e494e4c37ac "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivemembers/item/application"
    ie031435fe6106bcf58e7e7c567cfef730ae796ab68ef27fcafba68fcfcb4b029 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivemembers/item/device"
    iec67ab186e30acc0171e6c47d2dda98a3542cda89f6d46a52eebbe1628ddf9c6 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivemembers/item/user"
)

// DirectoryObjectItemRequestBuilder provides operations to manage the transitiveMembers property of the microsoft.graph.group entity.
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// DirectoryObjectItemRequestBuilderGetQueryParameters the direct and transitive members of a group. Nullable.
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
func (m *DirectoryObjectItemRequestBuilder) Application()(*ic3bf19ec033cea06ca57ed6e10c7103f65b6eac7c428254155d46e494e4c37ac.ApplicationRequestBuilder) {
    return ic3bf19ec033cea06ca57ed6e10c7103f65b6eac7c428254155d46e494e4c37ac.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/transitiveMembers/{directoryObject%2Did}{?%24select,%24expand}";
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
// CreateGetRequestInformation the direct and transitive members of a group. Nullable.
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
func (m *DirectoryObjectItemRequestBuilder) Device()(*ie031435fe6106bcf58e7e7c567cfef730ae796ab68ef27fcafba68fcfcb4b029.DeviceRequestBuilder) {
    return ie031435fe6106bcf58e7e7c567cfef730ae796ab68ef27fcafba68fcfcb4b029.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the direct and transitive members of a group. Nullable.
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
func (m *DirectoryObjectItemRequestBuilder) Group()(*i2ae8f05afc4ab7755eda7ac570865f043c6a61aeb68d9792631c6e7297790f17.GroupRequestBuilder) {
    return i2ae8f05afc4ab7755eda7ac570865f043c6a61aeb68d9792631c6e7297790f17.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *DirectoryObjectItemRequestBuilder) OrgContact()(*i350aaf1b458d5143e8c2c0377a329c7309406d5c18dbea162e2e59039870e589.OrgContactRequestBuilder) {
    return i350aaf1b458d5143e8c2c0377a329c7309406d5c18dbea162e2e59039870e589.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*i709d4a15bc3a3bf78e9c6210bf45793f6eaf5531c0fcb372ecfb169e1d16f617.ServicePrincipalRequestBuilder) {
    return i709d4a15bc3a3bf78e9c6210bf45793f6eaf5531c0fcb372ecfb169e1d16f617.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*iec67ab186e30acc0171e6c47d2dda98a3542cda89f6d46a52eebbe1628ddf9c6.UserRequestBuilder) {
    return iec67ab186e30acc0171e6c47d2dda98a3542cda89f6d46a52eebbe1628ddf9c6.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
