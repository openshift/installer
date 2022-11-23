package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i3024d5fb5e98d595fa699297ee7859a5760e8b0d247c03a1de15be6e9cd4b47e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/item/serviceprincipal"
    i4511ccf77c999c33a307743400000d916b8ef6f73b6ce5b8618995d9057a9b20 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/item/ref"
    i5537b56aece055af280278a8530d48680880ab267a0d8bd8c079bf40aefe7579 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/item/device"
    i886d422bd918de0b2928aca7f6f80bdcbbc0707096cbfb6c36af837178c4f2bb "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/item/orgcontact"
    i95c21ccf661796cd8e94a65125227983436c27288e57f178d0a9ea7e201cb62d "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/item/user"
    iefa34d46d801a95d0fed5bdd1aa528e120037e065245037862a34a4c39e5b6ca "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/item/application"
    if26fe6278b9799c8a3493f5e61e5d6b3bc3a629544bcfc96a440c979cf445852 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/item/group"
)

// DirectoryObjectItemRequestBuilder builds and executes requests for operations under \groups\{group-id}\members\{directoryObject-id}
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// Application casts the previous resource to application.
func (m *DirectoryObjectItemRequestBuilder) Application()(*iefa34d46d801a95d0fed5bdd1aa528e120037e065245037862a34a4c39e5b6ca.ApplicationRequestBuilder) {
    return iefa34d46d801a95d0fed5bdd1aa528e120037e065245037862a34a4c39e5b6ca.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/members/{directoryObject%2Did}";
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
func (m *DirectoryObjectItemRequestBuilder) Device()(*i5537b56aece055af280278a8530d48680880ab267a0d8bd8c079bf40aefe7579.DeviceRequestBuilder) {
    return i5537b56aece055af280278a8530d48680880ab267a0d8bd8c079bf40aefe7579.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Group casts the previous resource to group.
func (m *DirectoryObjectItemRequestBuilder) Group()(*if26fe6278b9799c8a3493f5e61e5d6b3bc3a629544bcfc96a440c979cf445852.GroupRequestBuilder) {
    return if26fe6278b9799c8a3493f5e61e5d6b3bc3a629544bcfc96a440c979cf445852.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *DirectoryObjectItemRequestBuilder) OrgContact()(*i886d422bd918de0b2928aca7f6f80bdcbbc0707096cbfb6c36af837178c4f2bb.OrgContactRequestBuilder) {
    return i886d422bd918de0b2928aca7f6f80bdcbbc0707096cbfb6c36af837178c4f2bb.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of group entities.
func (m *DirectoryObjectItemRequestBuilder) Ref()(*i4511ccf77c999c33a307743400000d916b8ef6f73b6ce5b8618995d9057a9b20.RefRequestBuilder) {
    return i4511ccf77c999c33a307743400000d916b8ef6f73b6ce5b8618995d9057a9b20.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*i3024d5fb5e98d595fa699297ee7859a5760e8b0d247c03a1de15be6e9cd4b47e.ServicePrincipalRequestBuilder) {
    return i3024d5fb5e98d595fa699297ee7859a5760e8b0d247c03a1de15be6e9cd4b47e.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*i95c21ccf661796cd8e94a65125227983436c27288e57f178d0a9ea7e201cb62d.UserRequestBuilder) {
    return i95c21ccf661796cd8e94a65125227983436c27288e57f178d0a9ea7e201cb62d.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
