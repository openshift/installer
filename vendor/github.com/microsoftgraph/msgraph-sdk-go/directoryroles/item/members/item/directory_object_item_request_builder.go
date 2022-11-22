package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i48ba5b711c92bd162655921d6d58e2d12bdd2ee8bd7a154e89f1b364bc5fe390 "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/item/group"
    i96fa7400fc436ca3b0224a0fe4903b59431dd1ad30f062971852ecc2653b2acb "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/item/orgcontact"
    iab6fee5f045e1f6aca945cfb7a1db41d3d7a39f9b26d1a9b40802c849fb3c65a "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/item/ref"
    id308d7b0e6f6b421378887115c265f2e955f684ce9659e0fcf892ad06446657e "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/item/serviceprincipal"
    id67784ffeb10bf166001104716bb18030b9adb87d6959064ed974b7b6e5263c7 "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/item/user"
    idbd71bd2e4ad0b1a442ad82378e107ec7bf1dcfe240bc3f55c831d04643e06cc "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/item/application"
    ie580b0330186e76aaeab0f4f9e301bfe9adff56cd9d17adde0a8cae9dddd9d8c "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item/members/item/device"
)

// DirectoryObjectItemRequestBuilder builds and executes requests for operations under \directoryRoles\{directoryRole-id}\members\{directoryObject-id}
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// Application casts the previous resource to application.
func (m *DirectoryObjectItemRequestBuilder) Application()(*idbd71bd2e4ad0b1a442ad82378e107ec7bf1dcfe240bc3f55c831d04643e06cc.ApplicationRequestBuilder) {
    return idbd71bd2e4ad0b1a442ad82378e107ec7bf1dcfe240bc3f55c831d04643e06cc.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/directoryRoles/{directoryRole%2Did}/members/{directoryObject%2Did}";
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
// Device casts the previous resource to device.
func (m *DirectoryObjectItemRequestBuilder) Device()(*ie580b0330186e76aaeab0f4f9e301bfe9adff56cd9d17adde0a8cae9dddd9d8c.DeviceRequestBuilder) {
    return ie580b0330186e76aaeab0f4f9e301bfe9adff56cd9d17adde0a8cae9dddd9d8c.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Group casts the previous resource to group.
func (m *DirectoryObjectItemRequestBuilder) Group()(*i48ba5b711c92bd162655921d6d58e2d12bdd2ee8bd7a154e89f1b364bc5fe390.GroupRequestBuilder) {
    return i48ba5b711c92bd162655921d6d58e2d12bdd2ee8bd7a154e89f1b364bc5fe390.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *DirectoryObjectItemRequestBuilder) OrgContact()(*i96fa7400fc436ca3b0224a0fe4903b59431dd1ad30f062971852ecc2653b2acb.OrgContactRequestBuilder) {
    return i96fa7400fc436ca3b0224a0fe4903b59431dd1ad30f062971852ecc2653b2acb.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of directoryRole entities.
func (m *DirectoryObjectItemRequestBuilder) Ref()(*iab6fee5f045e1f6aca945cfb7a1db41d3d7a39f9b26d1a9b40802c849fb3c65a.RefRequestBuilder) {
    return iab6fee5f045e1f6aca945cfb7a1db41d3d7a39f9b26d1a9b40802c849fb3c65a.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*id308d7b0e6f6b421378887115c265f2e955f684ce9659e0fcf892ad06446657e.ServicePrincipalRequestBuilder) {
    return id308d7b0e6f6b421378887115c265f2e955f684ce9659e0fcf892ad06446657e.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*id67784ffeb10bf166001104716bb18030b9adb87d6959064ed974b7b6e5263c7.UserRequestBuilder) {
    return id67784ffeb10bf166001104716bb18030b9adb87d6959064ed974b7b6e5263c7.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
