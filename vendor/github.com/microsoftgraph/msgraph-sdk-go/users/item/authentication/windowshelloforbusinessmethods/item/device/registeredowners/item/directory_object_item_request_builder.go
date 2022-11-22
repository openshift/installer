package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i4add89608b6a10347c3f4c00f51e651d51d73668a55257a6c83f7955ef936d62 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/registeredowners/item/serviceprincipal"
    i82efd56eb2c7b6598fff20c6bd4d55025484ba7763ce7a3546161e1b091afe05 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/registeredowners/item/endpoint"
    i9517e1b91927c70542e0afee13ae6d1dc538e920891a127823afc0dd51d03266 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/registeredowners/item/user"
    iaeec103b7deba7ea9f2d05f91eda2cc58739a5a0a5db3fc47ee229a39994e649 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/registeredowners/item/approleassignment"
    id6dd6b5e96a42f68e5808fbc66b9e46deb31d7cbaabe60d450825b94a0f08b52 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/windowshelloforbusinessmethods/item/device/registeredowners/item/ref"
)

// DirectoryObjectItemRequestBuilder builds and executes requests for operations under \users\{user-id}\authentication\windowsHelloForBusinessMethods\{windowsHelloForBusinessAuthenticationMethod-id}\device\registeredOwners\{directoryObject-id}
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// AppRoleAssignment casts the previous resource to appRoleAssignment.
func (m *DirectoryObjectItemRequestBuilder) AppRoleAssignment()(*iaeec103b7deba7ea9f2d05f91eda2cc58739a5a0a5db3fc47ee229a39994e649.AppRoleAssignmentRequestBuilder) {
    return iaeec103b7deba7ea9f2d05f91eda2cc58739a5a0a5db3fc47ee229a39994e649.NewAppRoleAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/authentication/windowsHelloForBusinessMethods/{windowsHelloForBusinessAuthenticationMethod%2Did}/device/registeredOwners/{directoryObject%2Did}";
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
func (m *DirectoryObjectItemRequestBuilder) Endpoint()(*i82efd56eb2c7b6598fff20c6bd4d55025484ba7763ce7a3546161e1b091afe05.EndpointRequestBuilder) {
    return i82efd56eb2c7b6598fff20c6bd4d55025484ba7763ce7a3546161e1b091afe05.NewEndpointRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of user entities.
func (m *DirectoryObjectItemRequestBuilder) Ref()(*id6dd6b5e96a42f68e5808fbc66b9e46deb31d7cbaabe60d450825b94a0f08b52.RefRequestBuilder) {
    return id6dd6b5e96a42f68e5808fbc66b9e46deb31d7cbaabe60d450825b94a0f08b52.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*i4add89608b6a10347c3f4c00f51e651d51d73668a55257a6c83f7955ef936d62.ServicePrincipalRequestBuilder) {
    return i4add89608b6a10347c3f4c00f51e651d51d73668a55257a6c83f7955ef936d62.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*i9517e1b91927c70542e0afee13ae6d1dc538e920891a127823afc0dd51d03266.UserRequestBuilder) {
    return i9517e1b91927c70542e0afee13ae6d1dc538e920891a127823afc0dd51d03266.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
