package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i3208fb270ef1573d8b02648df64e0496368c11b1b0232a27f57128ca46c84c92 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/transitivememberof/item/orgcontact"
    i406b9936e444ade32d7efacb15a1360f81c47b94c1a15ea9fe45f5eb570c9cb9 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/transitivememberof/item/application"
    ia169eb3e8a7d81653c08d76ee45873313fc53cd24e0946e75096742095058179 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/transitivememberof/item/user"
    id25ef0b3c256771e35062363296bde72395bd20c7d18467cd5c79051f66a6ffc "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/transitivememberof/item/group"
    idfcc5aea052eeb9ad5789e1d96f6fbacb50afb27dcaaff6a2167943b9cc32c3f "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/transitivememberof/item/device"
    iea6fa7dc608380a76ae23c98db406f86b93c82ed24b88d2bdfadbc5ef1487db3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/transitivememberof/item/serviceprincipal"
)

// DirectoryObjectItemRequestBuilder provides operations to manage the transitiveMemberOf property of the microsoft.graph.device entity.
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// DirectoryObjectItemRequestBuilderGetQueryParameters groups and administrative units that the device is a member of. This operation is transitive. Supports $expand.
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
func (m *DirectoryObjectItemRequestBuilder) Application()(*i406b9936e444ade32d7efacb15a1360f81c47b94c1a15ea9fe45f5eb570c9cb9.ApplicationRequestBuilder) {
    return i406b9936e444ade32d7efacb15a1360f81c47b94c1a15ea9fe45f5eb570c9cb9.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/authentication/windowsHelloForBusinessMethods/{windowsHelloForBusinessAuthenticationMethod%2Did}/device/transitiveMemberOf/{directoryObject%2Did}{?%24select,%24expand}";
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
// CreateGetRequestInformation groups and administrative units that the device is a member of. This operation is transitive. Supports $expand.
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
func (m *DirectoryObjectItemRequestBuilder) Device()(*idfcc5aea052eeb9ad5789e1d96f6fbacb50afb27dcaaff6a2167943b9cc32c3f.DeviceRequestBuilder) {
    return idfcc5aea052eeb9ad5789e1d96f6fbacb50afb27dcaaff6a2167943b9cc32c3f.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get groups and administrative units that the device is a member of. This operation is transitive. Supports $expand.
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
func (m *DirectoryObjectItemRequestBuilder) Group()(*id25ef0b3c256771e35062363296bde72395bd20c7d18467cd5c79051f66a6ffc.GroupRequestBuilder) {
    return id25ef0b3c256771e35062363296bde72395bd20c7d18467cd5c79051f66a6ffc.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *DirectoryObjectItemRequestBuilder) OrgContact()(*i3208fb270ef1573d8b02648df64e0496368c11b1b0232a27f57128ca46c84c92.OrgContactRequestBuilder) {
    return i3208fb270ef1573d8b02648df64e0496368c11b1b0232a27f57128ca46c84c92.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*iea6fa7dc608380a76ae23c98db406f86b93c82ed24b88d2bdfadbc5ef1487db3.ServicePrincipalRequestBuilder) {
    return iea6fa7dc608380a76ae23c98db406f86b93c82ed24b88d2bdfadbc5ef1487db3.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*ia169eb3e8a7d81653c08d76ee45873313fc53cd24e0946e75096742095058179.UserRequestBuilder) {
    return ia169eb3e8a7d81653c08d76ee45873313fc53cd24e0946e75096742095058179.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
