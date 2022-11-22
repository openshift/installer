package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i2afd08d9057801444b8fb0855e61f89a9dc9349812fb3358a2a83134e6c9f7cb "github.com/microsoftgraph/msgraph-sdk-go/devices/item/registeredowners/item/ref"
    i3d690bc4924e1d8c94e26a744c97df3aa52501e0d12b37409805433af1a351d8 "github.com/microsoftgraph/msgraph-sdk-go/devices/item/registeredowners/item/user"
    i81db9cb9b5e01d6a9467549b6aeb61786ab61a5e527db590390a72f711931d9f "github.com/microsoftgraph/msgraph-sdk-go/devices/item/registeredowners/item/serviceprincipal"
    ic69c52bc8d21baa76a79334a103205fb49ab2ca4a592044cd18eeee16742e2b9 "github.com/microsoftgraph/msgraph-sdk-go/devices/item/registeredowners/item/approleassignment"
    ifb7b63e7abea4c62bc0552144f0a84df68f6ed2c7ea7fe5bca5637f2380a7dc7 "github.com/microsoftgraph/msgraph-sdk-go/devices/item/registeredowners/item/endpoint"
)

// DirectoryObjectItemRequestBuilder builds and executes requests for operations under \devices\{device-id}\registeredOwners\{directoryObject-id}
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// AppRoleAssignment casts the previous resource to appRoleAssignment.
func (m *DirectoryObjectItemRequestBuilder) AppRoleAssignment()(*ic69c52bc8d21baa76a79334a103205fb49ab2ca4a592044cd18eeee16742e2b9.AppRoleAssignmentRequestBuilder) {
    return ic69c52bc8d21baa76a79334a103205fb49ab2ca4a592044cd18eeee16742e2b9.NewAppRoleAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/devices/{device%2Did}/registeredOwners/{directoryObject%2Did}";
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
func (m *DirectoryObjectItemRequestBuilder) Endpoint()(*ifb7b63e7abea4c62bc0552144f0a84df68f6ed2c7ea7fe5bca5637f2380a7dc7.EndpointRequestBuilder) {
    return ifb7b63e7abea4c62bc0552144f0a84df68f6ed2c7ea7fe5bca5637f2380a7dc7.NewEndpointRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of device entities.
func (m *DirectoryObjectItemRequestBuilder) Ref()(*i2afd08d9057801444b8fb0855e61f89a9dc9349812fb3358a2a83134e6c9f7cb.RefRequestBuilder) {
    return i2afd08d9057801444b8fb0855e61f89a9dc9349812fb3358a2a83134e6c9f7cb.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*i81db9cb9b5e01d6a9467549b6aeb61786ab61a5e527db590390a72f711931d9f.ServicePrincipalRequestBuilder) {
    return i81db9cb9b5e01d6a9467549b6aeb61786ab61a5e527db590390a72f711931d9f.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*i3d690bc4924e1d8c94e26a744c97df3aa52501e0d12b37409805433af1a351d8.UserRequestBuilder) {
    return i3d690bc4924e1d8c94e26a744c97df3aa52501e0d12b37409805433af1a351d8.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
