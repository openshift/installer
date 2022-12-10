package members

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i11e20919f832341e262584ad928d6b68fc6761ffe45ba837fbba743b93049c6a "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/orgcontact"
    i22b76b06873ef22d081b87cd39242b00c9f3d874b53044609e6f4d69799c2217 "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/device"
    i538cd9e6fe479bb94a7b01859eb2ff3243e9b18ab461b9efe2acd70664138be5 "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/group"
    i91760463bf94ce2631f955f311292734ae90c0ded3031721f5840a86ee64d9d2 "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/serviceprincipal"
    ia005616f40555875a5c23d1cb98bd7cf53800b59083688a9bfd85f38c4296f65 "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/application"
    ia06c236a7f5077be7d0e8b7befc60817495891cbc6824fe0d96315aaf59a3799 "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/count"
    icb61932ea5b80cfcfea2b517015daa1cb0d27b7dfe85d2a67daa77e365e38bc6 "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/user"
    ifc82e34929252e43b26961c4aa1341cdd29b763a4e341c8bd7002b8d7770e3ec "github.com/microsoftgraph/msgraph-sdk-go/directory/administrativeunits/item/members/ref"
)

// MembersRequestBuilder provides operations to manage the members property of the microsoft.graph.administrativeUnit entity.
type MembersRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// MembersRequestBuilderGetQueryParameters users and groups that are members of this administrative unit. Supports $expand.
type MembersRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// MembersRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MembersRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MembersRequestBuilderGetQueryParameters
}
// Application casts the previous resource to application.
func (m *MembersRequestBuilder) Application()(*ia005616f40555875a5c23d1cb98bd7cf53800b59083688a9bfd85f38c4296f65.ApplicationRequestBuilder) {
    return ia005616f40555875a5c23d1cb98bd7cf53800b59083688a9bfd85f38c4296f65.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewMembersRequestBuilderInternal instantiates a new MembersRequestBuilder and sets the default values.
func NewMembersRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MembersRequestBuilder) {
    m := &MembersRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/directory/administrativeUnits/{administrativeUnit%2Did}/members{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewMembersRequestBuilder instantiates a new MembersRequestBuilder and sets the default values.
func NewMembersRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MembersRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMembersRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *MembersRequestBuilder) Count()(*ia06c236a7f5077be7d0e8b7befc60817495891cbc6824fe0d96315aaf59a3799.CountRequestBuilder) {
    return ia06c236a7f5077be7d0e8b7befc60817495891cbc6824fe0d96315aaf59a3799.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation users and groups that are members of this administrative unit. Supports $expand.
func (m *MembersRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *MembersRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET
    requestInfo.Headers["Accept"] = "application/json"
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// Device casts the previous resource to device.
func (m *MembersRequestBuilder) Device()(*i22b76b06873ef22d081b87cd39242b00c9f3d874b53044609e6f4d69799c2217.DeviceRequestBuilder) {
    return i22b76b06873ef22d081b87cd39242b00c9f3d874b53044609e6f4d69799c2217.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get users and groups that are members of this administrative unit. Supports $expand.
func (m *MembersRequestBuilder) Get(ctx context.Context, requestConfiguration *MembersRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable), nil
}
// Group casts the previous resource to group.
func (m *MembersRequestBuilder) Group()(*i538cd9e6fe479bb94a7b01859eb2ff3243e9b18ab461b9efe2acd70664138be5.GroupRequestBuilder) {
    return i538cd9e6fe479bb94a7b01859eb2ff3243e9b18ab461b9efe2acd70664138be5.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *MembersRequestBuilder) OrgContact()(*i11e20919f832341e262584ad928d6b68fc6761ffe45ba837fbba743b93049c6a.OrgContactRequestBuilder) {
    return i11e20919f832341e262584ad928d6b68fc6761ffe45ba837fbba743b93049c6a.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of directory entities.
func (m *MembersRequestBuilder) Ref()(*ifc82e34929252e43b26961c4aa1341cdd29b763a4e341c8bd7002b8d7770e3ec.RefRequestBuilder) {
    return ifc82e34929252e43b26961c4aa1341cdd29b763a4e341c8bd7002b8d7770e3ec.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *MembersRequestBuilder) ServicePrincipal()(*i91760463bf94ce2631f955f311292734ae90c0ded3031721f5840a86ee64d9d2.ServicePrincipalRequestBuilder) {
    return i91760463bf94ce2631f955f311292734ae90c0ded3031721f5840a86ee64d9d2.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *MembersRequestBuilder) User()(*icb61932ea5b80cfcfea2b517015daa1cb0d27b7dfe85d2a67daa77e365e38bc6.UserRequestBuilder) {
    return icb61932ea5b80cfcfea2b517015daa1cb0d27b7dfe85d2a67daa77e365e38bc6.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
