package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i68434da45e4bb225cbc9b29081433b990c981fdb21a3332a2482c471251d9006 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners/item/serviceprincipal"
    i8452d4525438568644b8d3f8a41e4b9d4e29c9c85f07acbfc5b8451d0e984854 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners/item/endpoint"
    i943e8e991ddfe45932e81d03b32a50c6031bb5999252762ebbc389d545fb1cf5 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners/item/user"
    ibff776cb0947ce892f7e2234c70ce708a7a2e2570ea82e2d70bb902fb1dda1ae "github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners/item/ref"
    if871d8e45bc3009a6cdd0e483570f92d82c94cbbeb63a506e645c048fd4d9080 "github.com/microsoftgraph/msgraph-sdk-go/applications/item/owners/item/approleassignment"
)

// DirectoryObjectItemRequestBuilder builds and executes requests for operations under \applications\{application-id}\owners\{directoryObject-id}
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// AppRoleAssignment casts the previous resource to appRoleAssignment.
func (m *DirectoryObjectItemRequestBuilder) AppRoleAssignment()(*if871d8e45bc3009a6cdd0e483570f92d82c94cbbeb63a506e645c048fd4d9080.AppRoleAssignmentRequestBuilder) {
    return if871d8e45bc3009a6cdd0e483570f92d82c94cbbeb63a506e645c048fd4d9080.NewAppRoleAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/applications/{application%2Did}/owners/{directoryObject%2Did}";
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
func (m *DirectoryObjectItemRequestBuilder) Endpoint()(*i8452d4525438568644b8d3f8a41e4b9d4e29c9c85f07acbfc5b8451d0e984854.EndpointRequestBuilder) {
    return i8452d4525438568644b8d3f8a41e4b9d4e29c9c85f07acbfc5b8451d0e984854.NewEndpointRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of application entities.
func (m *DirectoryObjectItemRequestBuilder) Ref()(*ibff776cb0947ce892f7e2234c70ce708a7a2e2570ea82e2d70bb902fb1dda1ae.RefRequestBuilder) {
    return ibff776cb0947ce892f7e2234c70ce708a7a2e2570ea82e2d70bb902fb1dda1ae.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*i68434da45e4bb225cbc9b29081433b990c981fdb21a3332a2482c471251d9006.ServicePrincipalRequestBuilder) {
    return i68434da45e4bb225cbc9b29081433b990c981fdb21a3332a2482c471251d9006.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*i943e8e991ddfe45932e81d03b32a50c6031bb5999252762ebbc389d545fb1cf5.UserRequestBuilder) {
    return i943e8e991ddfe45932e81d03b32a50c6031bb5999252762ebbc389d545fb1cf5.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
