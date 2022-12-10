package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i40c8d5da7fc9a6217b6efac58c3740cd58572b63d41abf7ff231e644f35845bf "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items"
    i6048965dbea1b9d4b091a4401571b9e4939fc3dcaa6066a36fadd4ad0d88ebe2 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/following"
    i792e0b31993e4efcc79acff25a9598f5972b2ddf9c703fe62d86508fdf135dfe "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/root"
    i84f47258d3d9d84f919cffd7beaf67d9220276155cd1e9600f024e1c38fc1ac1 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/special"
    ib478a2b7936efcee57f8b9b5defb392493dafd807eac1b17259bc2e4f8aa2c25 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/sharedwithme"
    ib725db0f66452d6a294a8a50ef0739eb1c71f93000b56705a842e7cf9abe54f2 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/searchwithq"
    ic2c397128a0210051e36306fe11b6f791657369b585e62c3d141e8916c6f8fda "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/recent"
    ief3ab11045cc2316a7986e50ca5252ef1bf5247bdd7b782ebc616fd9dedc8c6a "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/list"
    if274b7c54959658e48996ff8d954ee4de4e1313851545174efe44dd6d0573768 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/bundles"
    i1b0986cf2bb3f8cdd5c9626af60c220b8ed74834d8de43da091f289373bb0f7e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/following/item"
    i7549c03bbbc8a5eeb871e3d1e7e255b35d87f970c4a55017c135ddac1bc09a5f "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/items/item"
    i7eaa82d7e156c1f8fe0101ef8557ffce59e5382dc33b38bf5c25f50084ccf3a1 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/special/item"
    i8b396b63c854a201e00ff83a78d93f8292a7b1a3895c9c05cf8b735f3efe3138 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/bundles/item"
)

// DriveItemRequestBuilder provides operations to manage the drives property of the microsoft.graph.group entity.
type DriveItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// DriveItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DriveItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// DriveItemRequestBuilderGetQueryParameters the group's drives. Read-only.
type DriveItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DriveItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DriveItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DriveItemRequestBuilderGetQueryParameters
}
// DriveItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DriveItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Bundles provides operations to manage the bundles property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) Bundles()(*if274b7c54959658e48996ff8d954ee4de4e1313851545174efe44dd6d0573768.BundlesRequestBuilder) {
    return if274b7c54959658e48996ff8d954ee4de4e1313851545174efe44dd6d0573768.NewBundlesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BundlesById provides operations to manage the bundles property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) BundlesById(id string)(*i8b396b63c854a201e00ff83a78d93f8292a7b1a3895c9c05cf8b735f3efe3138.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i8b396b63c854a201e00ff83a78d93f8292a7b1a3895c9c05cf8b735f3efe3138.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewDriveItemRequestBuilderInternal instantiates a new DriveItemRequestBuilder and sets the default values.
func NewDriveItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DriveItemRequestBuilder) {
    m := &DriveItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/drives/{drive%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewDriveItemRequestBuilder instantiates a new DriveItemRequestBuilder and sets the default values.
func NewDriveItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DriveItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDriveItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property drives for groups
func (m *DriveItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *DriveItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE
    if requestConfiguration != nil {
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// CreateGetRequestInformation the group's drives. Read-only.
func (m *DriveItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *DriveItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property drives in groups
func (m *DriveItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Driveable, requestConfiguration *DriveItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH
    requestInfo.Headers["Accept"] = "application/json"
    requestInfo.SetContentFromParsable(ctx, m.requestAdapter, "application/json", body)
    if requestConfiguration != nil {
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// Delete delete navigation property drives for groups
func (m *DriveItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *DriveItemRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.CreateDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.requestAdapter.SendNoContentAsync(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// Following provides operations to manage the following property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) Following()(*i6048965dbea1b9d4b091a4401571b9e4939fc3dcaa6066a36fadd4ad0d88ebe2.FollowingRequestBuilder) {
    return i6048965dbea1b9d4b091a4401571b9e4939fc3dcaa6066a36fadd4ad0d88ebe2.NewFollowingRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FollowingById provides operations to manage the following property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) FollowingById(id string)(*i1b0986cf2bb3f8cdd5c9626af60c220b8ed74834d8de43da091f289373bb0f7e.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i1b0986cf2bb3f8cdd5c9626af60c220b8ed74834d8de43da091f289373bb0f7e.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get the group's drives. Read-only.
func (m *DriveItemRequestBuilder) Get(ctx context.Context, requestConfiguration *DriveItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Driveable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDriveFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Driveable), nil
}
// Items provides operations to manage the items property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) Items()(*i40c8d5da7fc9a6217b6efac58c3740cd58572b63d41abf7ff231e644f35845bf.ItemsRequestBuilder) {
    return i40c8d5da7fc9a6217b6efac58c3740cd58572b63d41abf7ff231e644f35845bf.NewItemsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ItemsById provides operations to manage the items property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) ItemsById(id string)(*i7549c03bbbc8a5eeb871e3d1e7e255b35d87f970c4a55017c135ddac1bc09a5f.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i7549c03bbbc8a5eeb871e3d1e7e255b35d87f970c4a55017c135ddac1bc09a5f.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// List provides operations to manage the list property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) List()(*ief3ab11045cc2316a7986e50ca5252ef1bf5247bdd7b782ebc616fd9dedc8c6a.ListRequestBuilder) {
    return ief3ab11045cc2316a7986e50ca5252ef1bf5247bdd7b782ebc616fd9dedc8c6a.NewListRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property drives in groups
func (m *DriveItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Driveable, requestConfiguration *DriveItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Driveable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDriveFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Driveable), nil
}
// Recent provides operations to call the recent method.
func (m *DriveItemRequestBuilder) Recent()(*ic2c397128a0210051e36306fe11b6f791657369b585e62c3d141e8916c6f8fda.RecentRequestBuilder) {
    return ic2c397128a0210051e36306fe11b6f791657369b585e62c3d141e8916c6f8fda.NewRecentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Root provides operations to manage the root property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) Root()(*i792e0b31993e4efcc79acff25a9598f5972b2ddf9c703fe62d86508fdf135dfe.RootRequestBuilder) {
    return i792e0b31993e4efcc79acff25a9598f5972b2ddf9c703fe62d86508fdf135dfe.NewRootRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchWithQ provides operations to call the search method.
func (m *DriveItemRequestBuilder) SearchWithQ(q *string)(*ib725db0f66452d6a294a8a50ef0739eb1c71f93000b56705a842e7cf9abe54f2.SearchWithQRequestBuilder) {
    return ib725db0f66452d6a294a8a50ef0739eb1c71f93000b56705a842e7cf9abe54f2.NewSearchWithQRequestBuilderInternal(m.pathParameters, m.requestAdapter, q);
}
// SharedWithMe provides operations to call the sharedWithMe method.
func (m *DriveItemRequestBuilder) SharedWithMe()(*ib478a2b7936efcee57f8b9b5defb392493dafd807eac1b17259bc2e4f8aa2c25.SharedWithMeRequestBuilder) {
    return ib478a2b7936efcee57f8b9b5defb392493dafd807eac1b17259bc2e4f8aa2c25.NewSharedWithMeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Special provides operations to manage the special property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) Special()(*i84f47258d3d9d84f919cffd7beaf67d9220276155cd1e9600f024e1c38fc1ac1.SpecialRequestBuilder) {
    return i84f47258d3d9d84f919cffd7beaf67d9220276155cd1e9600f024e1c38fc1ac1.NewSpecialRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SpecialById provides operations to manage the special property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) SpecialById(id string)(*i7eaa82d7e156c1f8fe0101ef8557ffce59e5382dc33b38bf5c25f50084ccf3a1.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i7eaa82d7e156c1f8fe0101ef8557ffce59e5382dc33b38bf5c25f50084ccf3a1.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
