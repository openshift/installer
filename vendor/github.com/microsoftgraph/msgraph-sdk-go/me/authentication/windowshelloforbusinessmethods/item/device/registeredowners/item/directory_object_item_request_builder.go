package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i00967feecc7c4666a9920355ce33d30b6520ef377bfe05a1fa3309c68879d0b2 "github.com/microsoftgraph/msgraph-sdk-go/me/authentication/windowshelloforbusinessmethods/item/device/registeredowners/item/ref"
    i20e984775a5e9f82dd6d9e24f9ffbbae22603db6b07af2a1143cfb24a29e71f9 "github.com/microsoftgraph/msgraph-sdk-go/me/authentication/windowshelloforbusinessmethods/item/device/registeredowners/item/approleassignment"
    i55481771a625884c0b0d1361419ee69566b21c0f51c8274b69ce9a76280d32f9 "github.com/microsoftgraph/msgraph-sdk-go/me/authentication/windowshelloforbusinessmethods/item/device/registeredowners/item/endpoint"
    i64d0e684e9a720456d9863c0a91e227a729ff9f1a997e7c0ab1ae583487a75cd "github.com/microsoftgraph/msgraph-sdk-go/me/authentication/windowshelloforbusinessmethods/item/device/registeredowners/item/user"
    ia5b5d1728b877f29f6389d5019ce78c3745fa24b00cae91253b20120526c1f59 "github.com/microsoftgraph/msgraph-sdk-go/me/authentication/windowshelloforbusinessmethods/item/device/registeredowners/item/serviceprincipal"
)

// DirectoryObjectItemRequestBuilder builds and executes requests for operations under \me\authentication\windowsHelloForBusinessMethods\{windowsHelloForBusinessAuthenticationMethod-id}\device\registeredOwners\{directoryObject-id}
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// AppRoleAssignment casts the previous resource to appRoleAssignment.
func (m *DirectoryObjectItemRequestBuilder) AppRoleAssignment()(*i20e984775a5e9f82dd6d9e24f9ffbbae22603db6b07af2a1143cfb24a29e71f9.AppRoleAssignmentRequestBuilder) {
    return i20e984775a5e9f82dd6d9e24f9ffbbae22603db6b07af2a1143cfb24a29e71f9.NewAppRoleAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/authentication/windowsHelloForBusinessMethods/{windowsHelloForBusinessAuthenticationMethod%2Did}/device/registeredOwners/{directoryObject%2Did}";
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
func (m *DirectoryObjectItemRequestBuilder) Endpoint()(*i55481771a625884c0b0d1361419ee69566b21c0f51c8274b69ce9a76280d32f9.EndpointRequestBuilder) {
    return i55481771a625884c0b0d1361419ee69566b21c0f51c8274b69ce9a76280d32f9.NewEndpointRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of user entities.
func (m *DirectoryObjectItemRequestBuilder) Ref()(*i00967feecc7c4666a9920355ce33d30b6520ef377bfe05a1fa3309c68879d0b2.RefRequestBuilder) {
    return i00967feecc7c4666a9920355ce33d30b6520ef377bfe05a1fa3309c68879d0b2.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*ia5b5d1728b877f29f6389d5019ce78c3745fa24b00cae91253b20120526c1f59.ServicePrincipalRequestBuilder) {
    return ia5b5d1728b877f29f6389d5019ce78c3745fa24b00cae91253b20120526c1f59.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*i64d0e684e9a720456d9863c0a91e227a729ff9f1a997e7c0ab1ae583487a75cd.UserRequestBuilder) {
    return i64d0e684e9a720456d9863c0a91e227a729ff9f1a997e7c0ab1ae583487a75cd.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
