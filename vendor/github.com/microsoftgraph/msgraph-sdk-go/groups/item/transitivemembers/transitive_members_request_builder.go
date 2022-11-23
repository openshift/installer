package transitivemembers

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i12a3b468c9c0a3c8da68e094edc520b93e18d703211ee0b47096782f28bad499 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivemembers/orgcontact"
    i26cce04991ff66b1d72fc86177122d68b128b1a98365febf8e825e10d25130d6 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivemembers/serviceprincipal"
    i2ff12bb0273f0398126710f0aaff76c96e0476da5e02f0aec51fe60c03355b90 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivemembers/application"
    i66eaff71b4e75de2efa3edfb0945eda946c319d6ff2614b77f8e2655b5441840 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivemembers/user"
    i781807388bd7de0f338a2206afb969a3d310830b4eeccd08cb6380a0bd5013a6 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivemembers/count"
    ib4252bc1b5260565f25f5386665bc3f6473de29f04739bb1907f79500765ade2 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivemembers/device"
    icac3c120cfb8f1f90d877bdc23c985f5b254788eed73b61fd9c5ac62ceab091d "github.com/microsoftgraph/msgraph-sdk-go/groups/item/transitivemembers/group"
)

// TransitiveMembersRequestBuilder provides operations to manage the transitiveMembers property of the microsoft.graph.group entity.
type TransitiveMembersRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// TransitiveMembersRequestBuilderGetQueryParameters the direct and transitive members of a group. Nullable.
type TransitiveMembersRequestBuilderGetQueryParameters struct {
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
// TransitiveMembersRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TransitiveMembersRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *TransitiveMembersRequestBuilderGetQueryParameters
}
// Application casts the previous resource to application.
func (m *TransitiveMembersRequestBuilder) Application()(*i2ff12bb0273f0398126710f0aaff76c96e0476da5e02f0aec51fe60c03355b90.ApplicationRequestBuilder) {
    return i2ff12bb0273f0398126710f0aaff76c96e0476da5e02f0aec51fe60c03355b90.NewApplicationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewTransitiveMembersRequestBuilderInternal instantiates a new TransitiveMembersRequestBuilder and sets the default values.
func NewTransitiveMembersRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TransitiveMembersRequestBuilder) {
    m := &TransitiveMembersRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/transitiveMembers{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewTransitiveMembersRequestBuilder instantiates a new TransitiveMembersRequestBuilder and sets the default values.
func NewTransitiveMembersRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TransitiveMembersRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewTransitiveMembersRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *TransitiveMembersRequestBuilder) Count()(*i781807388bd7de0f338a2206afb969a3d310830b4eeccd08cb6380a0bd5013a6.CountRequestBuilder) {
    return i781807388bd7de0f338a2206afb969a3d310830b4eeccd08cb6380a0bd5013a6.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation the direct and transitive members of a group. Nullable.
func (m *TransitiveMembersRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *TransitiveMembersRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *TransitiveMembersRequestBuilder) Device()(*ib4252bc1b5260565f25f5386665bc3f6473de29f04739bb1907f79500765ade2.DeviceRequestBuilder) {
    return ib4252bc1b5260565f25f5386665bc3f6473de29f04739bb1907f79500765ade2.NewDeviceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the direct and transitive members of a group. Nullable.
func (m *TransitiveMembersRequestBuilder) Get(ctx context.Context, requestConfiguration *TransitiveMembersRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DirectoryObjectCollectionResponseable, error) {
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
func (m *TransitiveMembersRequestBuilder) Group()(*icac3c120cfb8f1f90d877bdc23c985f5b254788eed73b61fd9c5ac62ceab091d.GroupRequestBuilder) {
    return icac3c120cfb8f1f90d877bdc23c985f5b254788eed73b61fd9c5ac62ceab091d.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrgContact casts the previous resource to orgContact.
func (m *TransitiveMembersRequestBuilder) OrgContact()(*i12a3b468c9c0a3c8da68e094edc520b93e18d703211ee0b47096782f28bad499.OrgContactRequestBuilder) {
    return i12a3b468c9c0a3c8da68e094edc520b93e18d703211ee0b47096782f28bad499.NewOrgContactRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipal casts the previous resource to servicePrincipal.
func (m *TransitiveMembersRequestBuilder) ServicePrincipal()(*i26cce04991ff66b1d72fc86177122d68b128b1a98365febf8e825e10d25130d6.ServicePrincipalRequestBuilder) {
    return i26cce04991ff66b1d72fc86177122d68b128b1a98365febf8e825e10d25130d6.NewServicePrincipalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// User casts the previous resource to user.
func (m *TransitiveMembersRequestBuilder) User()(*i66eaff71b4e75de2efa3edfb0945eda946c319d6ff2614b77f8e2655b5441840.UserRequestBuilder) {
    return i66eaff71b4e75de2efa3edfb0945eda946c319d6ff2614b77f8e2655b5441840.NewUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
