package memberof

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i079cd090a0b1a43d285bfb42b24e0a45643b2cf978b487da585428327d3b1dd8 "github.com/microsoftgraph/msgraph-sdk-go/devices/item/memberof/application"
    i3e5bdea644de6960184732c44e7819ee6a85e03fe0fca854b731d2e032b07d70 "github.com/microsoftgraph/msgraph-sdk-go/devices/item/memberof/orgcontact"
    i63d02bb0e0653edb04f13a3a4522177d1238029e03a60504b29ac1abd83d1026 "github.com/microsoftgraph/msgraph-sdk-go/devices/item/memberof/count"
    i934ca7ab724fb534dfdf54df44ee091153306d5034de28e2fdabb83bda3eaf7c "github.com/microsoftgraph/msgraph-sdk-go/devices/item/memberof/device"
    ia820a6462df43d2ac8305df6039ff657b595340a11a45b3f77d852e557c9af6e "github.com/microsoftgraph/msgraph-sdk-go/devices/item/memberof/user"
    ie4078c687f41cedea0888972e52441970ebc8cf2939a19d7e1cc868183a4839a "github.com/microsoftgraph/msgraph-sdk-go/devices/item/memberof/group"
    ie88d98c78b2b454609b8767e72e6f67784c0a29be4f588094f591cc72a70d0b8 "github.com/microsoftgraph/msgraph-sdk-go/devices/item/memberof/serviceprincipal"
)

// MemberOfRequestBuilder provides operations to manage the memberOf property of the microsoft.graph.device entity.
type MemberOfRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// MemberOfRequestBuilderGetQueryParameters groups and administrative units that this device is a member of. Read-only. Nullable. Supports $expand.
type MemberOfRequestBuilderGetQueryParameters struct {
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
// MemberOfRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MemberOfRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MemberOfRequestBuilderGetQueryParameters
}
// Application casts the previous resource to application.
func (m *MemberOfRequestBuilder) Application()(*i079cd090a0b1a43d285bfb42b24e0a45643b2cf978b487da585428327d3b1dd8.ApplicationRequestBuilder) {
    return i079cd090a0b1a43d285bfb42b24e0a45643b2cf978b487da585428327d3b1dd8.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewMemberOfRequestBuilderInternal instantiates a new MemberOfRequestBuilder and sets the default values.
func NewMemberOfRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MemberOfRequestBuilder) {
    m := &MemberOfRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/devices/{device%2Did}/memberOf{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewMemberOfRequestBuilder instantiates a new MemberOfRequestBuilder and sets the default values.
func NewMemberOfRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MemberOfRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMemberOfRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *MemberOfRequestBuilder) Count()(*i63d02bb0e0653edb04f13a3a4522177d1238029e03a60504b29ac1abd83d1026.CountRequestBuilder) {
    return i63d02bb0e0653edb04f13a3a4522177d1238029e03a60504b29ac1abd83d1026.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation groups and administrative units that this device is a member of. Read-only. Nullable. Supports $expand.
func (m *MemberOfRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *MemberOfRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *MemberOfRequestBuilder) Device()(*i934ca7ab724fb534dfdf54df44ee091153306d5034de28e2fdabb83bda3eaf7c.DeviceRequestBuilder) {
    return i934ca7ab724fb534dfdf54df44ee091153306d5034de28e2fdabb83bda3eaf7c.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get groups and administrative units that this device is a member of. Read-only. Nullable. Supports $expand.
func (m *MemberOfRequestBuilder) Get(ctx context.Context, requestConfiguration *MemberOfRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
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
func (m *MemberOfRequestBuilder) Group()(*ie4078c687f41cedea0888972e52441970ebc8cf2939a19d7e1cc868183a4839a.GroupRequestBuilder) {
    return ie4078c687f41cedea0888972e52441970ebc8cf2939a19d7e1cc868183a4839a.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *MemberOfRequestBuilder) OrgContact()(*i3e5bdea644de6960184732c44e7819ee6a85e03fe0fca854b731d2e032b07d70.OrgContactRequestBuilder) {
    return i3e5bdea644de6960184732c44e7819ee6a85e03fe0fca854b731d2e032b07d70.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *MemberOfRequestBuilder) ServicePrincipal()(*ie88d98c78b2b454609b8767e72e6f67784c0a29be4f588094f591cc72a70d0b8.ServicePrincipalRequestBuilder) {
    return ie88d98c78b2b454609b8767e72e6f67784c0a29be4f588094f591cc72a70d0b8.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *MemberOfRequestBuilder) User()(*ia820a6462df43d2ac8305df6039ff657b595340a11a45b3f77d852e557c9af6e.UserRequestBuilder) {
    return ia820a6462df43d2ac8305df6039ff657b595340a11a45b3f77d852e557c9af6e.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
