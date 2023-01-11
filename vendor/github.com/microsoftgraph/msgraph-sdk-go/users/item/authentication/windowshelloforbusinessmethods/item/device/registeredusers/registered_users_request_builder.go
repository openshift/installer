package registeredusers

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i01f23acda21fbab4878fd68794c77fc4b40af367d7a5b2c449d68d0637f03232 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/registeredusers/serviceprincipal"
    i204d6f500221aef044863536e9b0d0deec0def51d40e66f7e006a2e3de09ceac "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/registeredusers/user"
    i2794a1473ba42345204449cadbc6b596c062643a2a52fb8ac73a58bfb9f56ac1 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/registeredusers/count"
    i4e2441dd82ccc677601124c522f836c09ed805f5850b03b16c5162c7d7d9abef "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/registeredusers/endpoint"
    ic345b5f6c874e9acc90b7d0083940bfea879132083dd717008f0e797b76637c2 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/registeredusers/approleassignment"
)

// RegisteredUsersRequestBuilder provides operations to manage the registeredUsers property of the microsoft.graph.device entity.
type RegisteredUsersRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// RegisteredUsersRequestBuilderGetQueryParameters collection of registered users of the device. For cloud joined devices and registered personal devices, registered users are set to the same value as registered owners at the time of registration. Read-only. Nullable. Supports $expand.
type RegisteredUsersRequestBuilderGetQueryParameters struct {
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
// RegisteredUsersRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type RegisteredUsersRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *RegisteredUsersRequestBuilderGetQueryParameters
}
// AppRoleAssignment casts the previous resource to appRoleAssignment.
func (m *RegisteredUsersRequestBuilder) AppRoleAssignment()(*ic345b5f6c874e9acc90b7d0083940bfea879132083dd717008f0e797b76637c2.AppRoleAssignmentRequestBuilder) {
    return ic345b5f6c874e9acc90b7d0083940bfea879132083dd717008f0e797b76637c2.NewAppRoleAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewRegisteredUsersRequestBuilderInternal instantiates a new RegisteredUsersRequestBuilder and sets the default values.
func NewRegisteredUsersRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RegisteredUsersRequestBuilder) {
    m := &RegisteredUsersRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/authentication/windowsHelloForBusinessMethods/{windowsHelloForBusinessAuthenticationMethod%2Did}/device/registeredUsers{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewRegisteredUsersRequestBuilder instantiates a new RegisteredUsersRequestBuilder and sets the default values.
func NewRegisteredUsersRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RegisteredUsersRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewRegisteredUsersRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *RegisteredUsersRequestBuilder) Count()(*i2794a1473ba42345204449cadbc6b596c062643a2a52fb8ac73a58bfb9f56ac1.CountRequestBuilder) {
    return i2794a1473ba42345204449cadbc6b596c062643a2a52fb8ac73a58bfb9f56ac1.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation collection of registered users of the device. For cloud joined devices and registered personal devices, registered users are set to the same value as registered owners at the time of registration. Read-only. Nullable. Supports $expand.
func (m *RegisteredUsersRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *RegisteredUsersRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Endpoint casts the previous resource to endpoint.
func (m *RegisteredUsersRequestBuilder) Endpoint()(*i4e2441dd82ccc677601124c522f836c09ed805f5850b03b16c5162c7d7d9abef.EndpointRequestBuilder) {
    return i4e2441dd82ccc677601124c522f836c09ed805f5850b03b16c5162c7d7d9abef.NewEndpointRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get collection of registered users of the device. For cloud joined devices and registered personal devices, registered users are set to the same value as registered owners at the time of registration. Read-only. Nullable. Supports $expand.
func (m *RegisteredUsersRequestBuilder) Get(ctx context.Context, requestConfiguration *RegisteredUsersRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
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
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *RegisteredUsersRequestBuilder) ServicePrincipal()(*i01f23acda21fbab4878fd68794c77fc4b40af367d7a5b2c449d68d0637f03232.ServicePrincipalRequestBuilder) {
    return i01f23acda21fbab4878fd68794c77fc4b40af367d7a5b2c449d68d0637f03232.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *RegisteredUsersRequestBuilder) User()(*i204d6f500221aef044863536e9b0d0deec0def51d40e66f7e006a2e3de09ceac.UserRequestBuilder) {
    return i204d6f500221aef044863536e9b0d0deec0def51d40e66f7e006a2e3de09ceac.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
