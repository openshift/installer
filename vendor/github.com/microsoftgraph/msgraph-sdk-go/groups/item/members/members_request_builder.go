package members

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i00889486dac3934978a99d8748dcd819aded9de67792784727a21eb08855d421 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/orgcontact"
    i016e4d6b6da76cf3891ad04b07a3aefb156c8c1a05ddf1dc0cc71bdff83020ca "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/application"
    i12e900b7863a530abb96344b2b364579136e814fd983663c5e21b9ef49ebc17c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/ref"
    i425cbef62da0ca6bd55549fd715a8e2d5a002b68479cb44a959d46ff4b53195b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/user"
    i7645bde648239f0a83e41332f0d46978f518babb2fcccb4f37aad18bd120f5ae "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/serviceprincipal"
    i819cbd5f3c989b4705fbafa07eff5a10397df1f7a3143fc77561710ad50430a3 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/count"
    icf485d88350fd6a8f4a8781f23e95e361c9b59f52510406c32121387a3edb467 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/group"
    iea34c95e06fc88b5aa637c42561062dba26018c8f549b2630e930189661808c8 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/members/device"
)

// MembersRequestBuilder provides operations to manage the members property of the microsoft.graph.group entity.
type MembersRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// MembersRequestBuilderGetQueryParameters the members of this group, who can be users, devices, other groups, or service principals. Supports the List members, Add member, and Remove member operations. Nullable. Supports $expand including nested $select. For example, /groups?$filter=startsWith(displayName,'Role')&$select=id,displayName&$expand=members($select=id,userPrincipalName,displayName).
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
func (m *MembersRequestBuilder) Application()(*i016e4d6b6da76cf3891ad04b07a3aefb156c8c1a05ddf1dc0cc71bdff83020ca.ApplicationRequestBuilder) {
    return i016e4d6b6da76cf3891ad04b07a3aefb156c8c1a05ddf1dc0cc71bdff83020ca.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewMembersRequestBuilderInternal instantiates a new MembersRequestBuilder and sets the default values.
func NewMembersRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MembersRequestBuilder) {
    m := &MembersRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/members{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
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
func (m *MembersRequestBuilder) Count()(*i819cbd5f3c989b4705fbafa07eff5a10397df1f7a3143fc77561710ad50430a3.CountRequestBuilder) {
    return i819cbd5f3c989b4705fbafa07eff5a10397df1f7a3143fc77561710ad50430a3.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation the members of this group, who can be users, devices, other groups, or service principals. Supports the List members, Add member, and Remove member operations. Nullable. Supports $expand including nested $select. For example, /groups?$filter=startsWith(displayName,'Role')&$select=id,displayName&$expand=members($select=id,userPrincipalName,displayName).
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
func (m *MembersRequestBuilder) Device()(*iea34c95e06fc88b5aa637c42561062dba26018c8f549b2630e930189661808c8.DeviceRequestBuilder) {
    return iea34c95e06fc88b5aa637c42561062dba26018c8f549b2630e930189661808c8.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the members of this group, who can be users, devices, other groups, or service principals. Supports the List members, Add member, and Remove member operations. Nullable. Supports $expand including nested $select. For example, /groups?$filter=startsWith(displayName,'Role')&$select=id,displayName&$expand=members($select=id,userPrincipalName,displayName).
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
func (m *MembersRequestBuilder) Group()(*icf485d88350fd6a8f4a8781f23e95e361c9b59f52510406c32121387a3edb467.GroupRequestBuilder) {
    return icf485d88350fd6a8f4a8781f23e95e361c9b59f52510406c32121387a3edb467.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *MembersRequestBuilder) OrgContact()(*i00889486dac3934978a99d8748dcd819aded9de67792784727a21eb08855d421.OrgContactRequestBuilder) {
    return i00889486dac3934978a99d8748dcd819aded9de67792784727a21eb08855d421.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Ref provides operations to manage the collection of group entities.
func (m *MembersRequestBuilder) Ref()(*i12e900b7863a530abb96344b2b364579136e814fd983663c5e21b9ef49ebc17c.RefRequestBuilder) {
    return i12e900b7863a530abb96344b2b364579136e814fd983663c5e21b9ef49ebc17c.NewRefRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *MembersRequestBuilder) ServicePrincipal()(*i7645bde648239f0a83e41332f0d46978f518babb2fcccb4f37aad18bd120f5ae.ServicePrincipalRequestBuilder) {
    return i7645bde648239f0a83e41332f0d46978f518babb2fcccb4f37aad18bd120f5ae.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *MembersRequestBuilder) User()(*i425cbef62da0ca6bd55549fd715a8e2d5a002b68479cb44a959d46ff4b53195b.UserRequestBuilder) {
    return i425cbef62da0ca6bd55549fd715a8e2d5a002b68479cb44a959d46ff4b53195b.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
