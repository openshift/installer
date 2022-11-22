package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i34747938a4f644018c7eb75dba6350962d8e7044c80c35711e4b90e95a736e43 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/owners/item/endpoint"
    i9f41b73e55924f4d6fbb7b19a34cc5ab46e3cebe58bbaa7d1ee72ef411a66e22 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/owners/item/serviceprincipal"
    ic5d2b48d40d8e6c63f1b3a73e8fdfbfa83736635640c04f7f87f6ab7a2006a06 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/owners/item/user"
    id521778c91b8c914afefedfddeb983090fa72b8986d2cecf42eafea5829df7ce "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/owners/item/ref"
    if03c06d05c56ccfd9745adc9a233ef7a0089bd1626b1d761d4d68fbea3dd4ff8 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item/owners/item/approleassignment"
)

// DirectoryObjectItemRequestBuilder builds and executes requests for operations under \servicePrincipals\{servicePrincipal-id}\owners\{directoryObject-id}
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// AppRoleAssignment casts the previous resource to appRoleAssignment.
func (m *DirectoryObjectItemRequestBuilder) AppRoleAssignment()(*if03c06d05c56ccfd9745adc9a233ef7a0089bd1626b1d761d4d68fbea3dd4ff8.AppRoleAssignmentRequestBuilder) {
    return if03c06d05c56ccfd9745adc9a233ef7a0089bd1626b1d761d4d68fbea3dd4ff8.NewAppRoleAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/owners/{directoryObject%2Did}";
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
func (m *DirectoryObjectItemRequestBuilder) Endpoint()(*i34747938a4f644018c7eb75dba6350962d8e7044c80c35711e4b90e95a736e43.EndpointRequestBuilder) {
    return i34747938a4f644018c7eb75dba6350962d8e7044c80c35711e4b90e95a736e43.NewEndpointRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of servicePrincipal entities.
func (m *DirectoryObjectItemRequestBuilder) Ref()(*id521778c91b8c914afefedfddeb983090fa72b8986d2cecf42eafea5829df7ce.RefRequestBuilder) {
    return id521778c91b8c914afefedfddeb983090fa72b8986d2cecf42eafea5829df7ce.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*i9f41b73e55924f4d6fbb7b19a34cc5ab46e3cebe58bbaa7d1ee72ef411a66e22.ServicePrincipalRequestBuilder) {
    return i9f41b73e55924f4d6fbb7b19a34cc5ab46e3cebe58bbaa7d1ee72ef411a66e22.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*ic5d2b48d40d8e6c63f1b3a73e8fdfbfa83736635640c04f7f87f6ab7a2006a06.UserRequestBuilder) {
    return ic5d2b48d40d8e6c63f1b3a73e8fdfbfa83736635640c04f7f87f6ab7a2006a06.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
