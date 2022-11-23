package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i1597afe854aa773a29f976b59e316338efbe03c7e646ded0ae8fc8ceca5db37d "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/item/orgcontact"
    i1e74e6374c62a9a81da83a465994d4ab4c60d2c098c02fec71b0f6c3ef409fe3 "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/item/serviceprincipal"
    i55626291fe3763517907fa1e8e48ed2e02fcb0c36a6d953c46eb656dc95bc14e "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/item/user"
    i83f7743adfd8c81de6ad254bb26959c0550ce9a311619d12ae6b587c5b9498f0 "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/item/group"
    i9c51595b2591fc467f81c1fadd948928c238e6b8413d3027cb1da891d7800f12 "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/item/device"
    ia2a0200c01baa8b81b2a0f264081d0eeb64074a4058bc543e4bd9b3c1697c4a4 "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/item/application"
    ibff336f36c3e7565bbfc28782a6891d97d35883a605a9647190505947733ea02 "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/item/ref"
)

// DirectoryObjectItemRequestBuilder builds and executes requests for operations under \directory\administrativeUnits\{administrativeUnit-id}\members\{directoryObject-id}
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// Application casts the previous resource to application.
func (m *DirectoryObjectItemRequestBuilder) Application()(*ia2a0200c01baa8b81b2a0f264081d0eeb64074a4058bc543e4bd9b3c1697c4a4.ApplicationRequestBuilder) {
    return ia2a0200c01baa8b81b2a0f264081d0eeb64074a4058bc543e4bd9b3c1697c4a4.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/directory/administrativeUnits/{administrativeUnit%2Did}/members/{directoryObject%2Did}";
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
func (m *DirectoryObjectItemRequestBuilder) Device()(*i9c51595b2591fc467f81c1fadd948928c238e6b8413d3027cb1da891d7800f12.DeviceRequestBuilder) {
    return i9c51595b2591fc467f81c1fadd948928c238e6b8413d3027cb1da891d7800f12.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Group casts the previous resource to group.
func (m *DirectoryObjectItemRequestBuilder) Group()(*i83f7743adfd8c81de6ad254bb26959c0550ce9a311619d12ae6b587c5b9498f0.GroupRequestBuilder) {
    return i83f7743adfd8c81de6ad254bb26959c0550ce9a311619d12ae6b587c5b9498f0.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *DirectoryObjectItemRequestBuilder) OrgContact()(*i1597afe854aa773a29f976b59e316338efbe03c7e646ded0ae8fc8ceca5db37d.OrgContactRequestBuilder) {
    return i1597afe854aa773a29f976b59e316338efbe03c7e646ded0ae8fc8ceca5db37d.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of directory entities.
func (m *DirectoryObjectItemRequestBuilder) Ref()(*ibff336f36c3e7565bbfc28782a6891d97d35883a605a9647190505947733ea02.RefRequestBuilder) {
    return ibff336f36c3e7565bbfc28782a6891d97d35883a605a9647190505947733ea02.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*i1e74e6374c62a9a81da83a465994d4ab4c60d2c098c02fec71b0f6c3ef409fe3.ServicePrincipalRequestBuilder) {
    return i1e74e6374c62a9a81da83a465994d4ab4c60d2c098c02fec71b0f6c3ef409fe3.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*i55626291fe3763517907fa1e8e48ed2e02fcb0c36a6d953c46eb656dc95bc14e.UserRequestBuilder) {
    return i55626291fe3763517907fa1e8e48ed2e02fcb0c36a6d953c46eb656dc95bc14e.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
