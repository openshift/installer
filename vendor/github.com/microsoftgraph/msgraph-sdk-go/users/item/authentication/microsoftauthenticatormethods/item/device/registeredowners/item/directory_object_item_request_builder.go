package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i01d32e73573aee5bbf169f9d0cf6e793af43962f254609a743be1cdc9d31f138 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/microsoftauthenticatormethods/item/device/registeredowners/item/approleassignment"
    ia23b2b7161f8667225c3234711c851060b0730e7e818c86ad9412cb13c1f6e89 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/microsoftauthenticatormethods/item/device/registeredowners/item/user"
    iadb4063643e081b9ce47f3c6f93732ef9edff3af44c9c41672599ddebf6067ec "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/microsoftauthenticatormethods/item/device/registeredowners/item/endpoint"
    ic87ff92704d94ccee3cd037464a19ca71ae29eb936f2a7849a80cd8a84554654 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/microsoftauthenticatormethods/item/device/registeredowners/item/serviceprincipal"
    ifb39cd9d946877313284177e48266f2a4a94bf9c83373fbf58d03d8373de95cd "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication/microsoftauthenticatormethods/item/device/registeredowners/item/ref"
)

// DirectoryObjectItemRequestBuilder builds and executes requests for operations under \users\{user-id}\authentication\microsoftAuthenticatorMethods\{microsoftAuthenticatorAuthenticationMethod-id}\device\registeredOwners\{directoryObject-id}
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// AppRoleAssignment casts the previous resource to appRoleAssignment.
func (m *DirectoryObjectItemRequestBuilder) AppRoleAssignment()(*i01d32e73573aee5bbf169f9d0cf6e793af43962f254609a743be1cdc9d31f138.AppRoleAssignmentRequestBuilder) {
    return i01d32e73573aee5bbf169f9d0cf6e793af43962f254609a743be1cdc9d31f138.NewAppRoleAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/authentication/microsoftAuthenticatorMethods/{microsoftAuthenticatorAuthenticationMethod%2Did}/device/registeredOwners/{directoryObject%2Did}";
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
func (m *DirectoryObjectItemRequestBuilder) Endpoint()(*iadb4063643e081b9ce47f3c6f93732ef9edff3af44c9c41672599ddebf6067ec.EndpointRequestBuilder) {
    return iadb4063643e081b9ce47f3c6f93732ef9edff3af44c9c41672599ddebf6067ec.NewEndpointRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of user entities.
func (m *DirectoryObjectItemRequestBuilder) Ref()(*ifb39cd9d946877313284177e48266f2a4a94bf9c83373fbf58d03d8373de95cd.RefRequestBuilder) {
    return ifb39cd9d946877313284177e48266f2a4a94bf9c83373fbf58d03d8373de95cd.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*ic87ff92704d94ccee3cd037464a19ca71ae29eb936f2a7849a80cd8a84554654.ServicePrincipalRequestBuilder) {
    return ic87ff92704d94ccee3cd037464a19ca71ae29eb936f2a7849a80cd8a84554654.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*ia23b2b7161f8667225c3234711c851060b0730e7e818c86ad9412cb13c1f6e89.UserRequestBuilder) {
    return ia23b2b7161f8667225c3234711c851060b0730e7e818c86ad9412cb13c1f6e89.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
