package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i35cd693249803017856aee8dfacd1989f6072292094e894ec0a08104a4545506 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberswithlicenseerrors/item/group"
    i60143daff92426f6dcf833351b1bd3ce8ddb92c92ac1149cac1504755b7e4f88 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberswithlicenseerrors/item/serviceprincipal"
    i625d4e85eca110568409d60002a5f45f71cabcd515ab026addca69d34f11311b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberswithlicenseerrors/item/user"
    ia74b2fd46cd25b0da3e840fdd2d90d81716aefed7b27d6614ebe043b4a64bea8 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberswithlicenseerrors/item/application"
    ic8b2ab57f847df16e7f415d4451ad15572b548c5e38c2493b5bedba65993b24c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberswithlicenseerrors/item/device"
    id4609afe24daea13b8f9a558b14db205ce9c111e305b5961b779709c43938ef3 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/memberswithlicenseerrors/item/orgcontact"
)

// DirectoryObjectItemRequestBuilder provides operations to manage the membersWithLicenseErrors property of the microsoft.graph.group entity.
type DirectoryObjectItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// DirectoryObjectItemRequestBuilderGetQueryParameters a list of group members with license errors from this group-based license assignment. Read-only.
type DirectoryObjectItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DirectoryObjectItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DirectoryObjectItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DirectoryObjectItemRequestBuilderGetQueryParameters
}
// Application casts the previous resource to application.
func (m *DirectoryObjectItemRequestBuilder) Application()(*ia74b2fd46cd25b0da3e840fdd2d90d81716aefed7b27d6614ebe043b4a64bea8.ApplicationRequestBuilder) {
    return ia74b2fd46cd25b0da3e840fdd2d90d81716aefed7b27d6614ebe043b4a64bea8.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewDirectoryObjectItemRequestBuilderInternal instantiates a new DirectoryObjectItemRequestBuilder and sets the default values.
func NewDirectoryObjectItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DirectoryObjectItemRequestBuilder) {
    m := &DirectoryObjectItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/membersWithLicenseErrors/{directoryObject%2Did}{?%24select,%24expand}";
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
// CreateGetRequestInformation a list of group members with license errors from this group-based license assignment. Read-only.
func (m *DirectoryObjectItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *DirectoryObjectItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *DirectoryObjectItemRequestBuilder) Device()(*ic8b2ab57f847df16e7f415d4451ad15572b548c5e38c2493b5bedba65993b24c.DeviceRequestBuilder) {
    return ic8b2ab57f847df16e7f415d4451ad15572b548c5e38c2493b5bedba65993b24c.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get a list of group members with license errors from this group-based license assignment. Read-only.
func (m *DirectoryObjectItemRequestBuilder) Get(ctx context.Context, requestConfiguration *DirectoryObjectItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDirectoryObjectFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectable), nil
}
// Group casts the previous resource to group.
func (m *DirectoryObjectItemRequestBuilder) Group()(*i35cd693249803017856aee8dfacd1989f6072292094e894ec0a08104a4545506.GroupRequestBuilder) {
    return i35cd693249803017856aee8dfacd1989f6072292094e894ec0a08104a4545506.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *DirectoryObjectItemRequestBuilder) OrgContact()(*id4609afe24daea13b8f9a558b14db205ce9c111e305b5961b779709c43938ef3.OrgContactRequestBuilder) {
    return id4609afe24daea13b8f9a558b14db205ce9c111e305b5961b779709c43938ef3.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *DirectoryObjectItemRequestBuilder) ServicePrincipal()(*i60143daff92426f6dcf833351b1bd3ce8ddb92c92ac1149cac1504755b7e4f88.ServicePrincipalRequestBuilder) {
    return i60143daff92426f6dcf833351b1bd3ce8ddb92c92ac1149cac1504755b7e4f88.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *DirectoryObjectItemRequestBuilder) User()(*i625d4e85eca110568409d60002a5f45f71cabcd515ab026addca69d34f11311b.UserRequestBuilder) {
    return i625d4e85eca110568409d60002a5f45f71cabcd515ab026addca69d34f11311b.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
