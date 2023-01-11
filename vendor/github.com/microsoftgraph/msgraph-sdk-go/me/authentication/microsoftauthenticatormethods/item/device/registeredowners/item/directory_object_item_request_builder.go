package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i1a9a2b210784a76c8fd74e819a489398ebb0382206781ced5c461aa04bc8db8c "github.com/microsoftgraph/msgraph-sdk-go/me/authentication/microsoftauthenticatormethods/item/device/registeredowners/item/serviceprincipal"
    i1df9bfba66ca4e458925f3578495d75c2eeed86f2237c766bbf95dd4ae83d760 "github.com/microsoftgraph/msgraph-sdk-go/me/authentication/microsoftauthenticatormethods/item/device/registeredowners/item/approleassignment"
    i24d2c08413bf06ae1f4e4108ace657d96ac048ed9b59928bfcf793bf44023179 "github.com/microsoftgraph/msgraph-sdk-go/me/authentication/microsoftauthenticatormethods/item/device/registeredowners/item/user"
    i2dff487c7016556062845ff9b22b1d863e77bd10e1c36ce76e1a875939c974b0 "github.com/microsoftgraph/msgraph-sdk-go/me/authentication/microsoftauthenticatormethods/item/device/registeredowners/item/endpoint"
    ic287c7cfa2f68fe178c491138c156ac2611fb8a71ba0b0b61d17fe2e82910eeb "github.com/microsoftgraph/msgraph-sdk-go/me/authentication/microsoftauthenticatormethods/item/device/registeredowners/item/ref"
)

// DirectoryObjectItemRequestBuilder builds and executes requests for operations under \me\authentication\microsoftAuthenticatorMethods\{microsoftAuthenticatorAuthenticationMethod-id}\device\registeredOwners\{directoryObject-id}
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// AppRoleAssignment casts the previous resource to appRoleAssignment.
func (m *DirectoryObjectItemRequestBuilder) AppRoleAssignment()(*i1df9bfba66ca4e458925f3578495d75c2eeed86f2237c766bbf95dd4ae83d760.AppRoleAssignmentRequestBuilder) {
    return i1df9bfba66ca4e458925f3578495d75c2eeed86f2237c766bbf95dd4ae83d760.NewAppRoleAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/authentication/microsoftAuthenticatorMethods/{microsoftAuthenticatorAuthenticationMethod%2Did}/device/registeredOwners/{directoryObject%2Did}";
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
// Endpoint casts the previous resource to endpoint.
func (m *DirectoryObjectItemRequestBuilder) Endpoint()(*i2dff487c7016556062845ff9b22b1d863e77bd10e1c36ce76e1a875939c974b0.EndpointRequestBuilder) {
    return i2dff487c7016556062845ff9b22b1d863e77bd10e1c36ce76e1a875939c974b0.NewEndpointRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of user entities.
func (m *DirectoryObjectItemRequestBuilder) Ref()(*ic287c7cfa2f68fe178c491138c156ac2611fb8a71ba0b0b61d17fe2e82910eeb.RefRequestBuilder) {
    return ic287c7cfa2f68fe178c491138c156ac2611fb8a71ba0b0b61d17fe2e82910eeb.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*i1a9a2b210784a76c8fd74e819a489398ebb0382206781ced5c461aa04bc8db8c.ServicePrincipalRequestBuilder) {
    return i1a9a2b210784a76c8fd74e819a489398ebb0382206781ced5c461aa04bc8db8c.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*i24d2c08413bf06ae1f4e4108ace657d96ac048ed9b59928bfcf793bf44023179.UserRequestBuilder) {
    return i24d2c08413bf06ae1f4e4108ace657d96ac048ed9b59928bfcf793bf44023179.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
