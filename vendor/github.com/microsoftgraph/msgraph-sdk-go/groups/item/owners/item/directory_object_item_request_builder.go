package item

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i01c86d8d244da854467b852e368fa3b56139676cf3bd79f9d0975efe556813f7 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/owners/item/application"
    i0e7d1d0923955ea395c736dadd5714b8aa3b3743f945f031bcab62e88f2c4516 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/owners/item/serviceprincipal"
    i23f494a43dd9e4d3ff3d92dd241d65a443036ed03b5eb21b51f0eb140dbe914f "github.com/microsoftgraph/msgraph-sdk-go/groups/item/owners/item/orgcontact"
    icde18b78835002aae795f439881b2754e9c74d4b07d0590937d1b4ff1f17d9e1 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/owners/item/ref"
    idee3b6c8f7ab29fbc3b37feae5065ea3c62baaaa38ef42a73302a148d75e25da "github.com/microsoftgraph/msgraph-sdk-go/groups/item/owners/item/group"
    iea2a8c1e6b028cc46f70aa1279749fff4540ea19f8900f8687dbd589042f4468 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/owners/item/user"
    ifd576f4c4a0343a3fbc154a1eb8bb119d8a362b790f8ad4ba14f71c0337cb71b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/owners/item/device"
)

// DirectoryObjectItemRequestBuilder builds and executes requests for operations under \groups\{group-id}\owners\{directoryObject-id}
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// Application casts the previous resource to application.
func (m *DirectoryObjectItemRequestBuilder) Application()(*i01c86d8d244da854467b852e368fa3b56139676cf3bd79f9d0975efe556813f7.ApplicationRequestBuilder) {
    return i01c86d8d244da854467b852e368fa3b56139676cf3bd79f9d0975efe556813f7.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/owners/{directoryObject%2Did}";
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
func (m *DirectoryObjectItemRequestBuilder) Device()(*ifd576f4c4a0343a3fbc154a1eb8bb119d8a362b790f8ad4ba14f71c0337cb71b.DeviceRequestBuilder) {
    return ifd576f4c4a0343a3fbc154a1eb8bb119d8a362b790f8ad4ba14f71c0337cb71b.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Group casts the previous resource to group.
func (m *DirectoryObjectItemRequestBuilder) Group()(*idee3b6c8f7ab29fbc3b37feae5065ea3c62baaaa38ef42a73302a148d75e25da.GroupRequestBuilder) {
    return idee3b6c8f7ab29fbc3b37feae5065ea3c62baaaa38ef42a73302a148d75e25da.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *DirectoryObjectItemRequestBuilder) OrgContact()(*i23f494a43dd9e4d3ff3d92dd241d65a443036ed03b5eb21b51f0eb140dbe914f.OrgContactRequestBuilder) {
    return i23f494a43dd9e4d3ff3d92dd241d65a443036ed03b5eb21b51f0eb140dbe914f.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of group entities.
func (m *DirectoryObjectItemRequestBuilder) Ref()(*icde18b78835002aae795f439881b2754e9c74d4b07d0590937d1b4ff1f17d9e1.RefRequestBuilder) {
    return icde18b78835002aae795f439881b2754e9c74d4b07d0590937d1b4ff1f17d9e1.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*i0e7d1d0923955ea395c736dadd5714b8aa3b3743f945f031bcab62e88f2c4516.ServicePrincipalRequestBuilder) {
    return i0e7d1d0923955ea395c736dadd5714b8aa3b3743f945f031bcab62e88f2c4516.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*iea2a8c1e6b028cc46f70aa1279749fff4540ea19f8900f8687dbd589042f4468.UserRequestBuilder) {
    return iea2a8c1e6b028cc46f70aa1279749fff4540ea19f8900f8687dbd589042f4468.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
