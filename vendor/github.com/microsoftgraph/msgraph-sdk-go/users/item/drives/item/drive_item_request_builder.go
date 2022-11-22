package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i26bf19300aed5efb255319d7da9911edd0ae7ca03e6fdfddef4a2f6581682ca5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/recent"
    i3959fc5f718498a64f7b4bcc9ac4b4bfc74564f1973213366e5ba756fa3b6cf1 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/bundles"
    ia43c04139e040bccd3c51a3e13b89c5bd9f4493722f732499431ccd547e8a502 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/sharedwithme"
    iaaf4d433dad634780d17407e36381b8dbd65675f4ec7d040609d686c17e54c3c "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items"
    iaceacb7c57c6e0cab9dfe9f00f869b5f5e1fcd6265e7550db1cab9f5e1e6953b "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/root"
    ic8094ec71dba1b039d40e34e48be10142245be65d5fbc5c57c4f20f6af133ef2 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/following"
    id96c6d741d51026adde933e1dc250eb7ed8ddf087f1f2a9235bcf77ad3f95e68 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/searchwithq"
    if8ffbadcd14d06fa791c22b60b93c57cf0043e5bd611c797b0956a4d1c513322 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/list"
    ifa64565d61311e1ff82b220f60ec0b00f6578cc52d806aed6adb3e97fe080314 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/special"
    i4dccdff9f8dcdf668768e001e5a27186fe4bf619c82ddb04a7646de52b9311fd "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/items/item"
    i619a3607ef091cfc613c8064ad8c26e2c851e6b29b5a8de1c8ae6ad11cfec472 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/following/item"
    i6d900d5819130f59b838504ba76d948433e47ce4b2aa7564aca30ad043db2dfd "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/bundles/item"
    i7707a2d5fd849a1571e51f19174fe037b6b38b3a81bed89f5c7a178732625346 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item/special/item"
)

// DriveItemRequestBuilder provides operations to manage the drives property of the microsoft.graph.user entity.
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
// DriveItemRequestBuilderGetQueryParameters a collection of drives available for this user. Read-only.
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
func (m *DriveItemRequestBuilder) Bundles()(*i3959fc5f718498a64f7b4bcc9ac4b4bfc74564f1973213366e5ba756fa3b6cf1.BundlesRequestBuilder) {
    return i3959fc5f718498a64f7b4bcc9ac4b4bfc74564f1973213366e5ba756fa3b6cf1.NewBundlesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BundlesById provides operations to manage the bundles property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) BundlesById(id string)(*i6d900d5819130f59b838504ba76d948433e47ce4b2aa7564aca30ad043db2dfd.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i6d900d5819130f59b838504ba76d948433e47ce4b2aa7564aca30ad043db2dfd.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewDriveItemRequestBuilderInternal instantiates a new DriveItemRequestBuilder and sets the default values.
func NewDriveItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DriveItemRequestBuilder) {
    m := &DriveItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/drives/{drive%2Did}{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property drives for users
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
// CreateGetRequestInformation a collection of drives available for this user. Read-only.
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
// CreatePatchRequestInformation update the navigation property drives in users
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
// Delete delete navigation property drives for users
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
func (m *DriveItemRequestBuilder) Following()(*ic8094ec71dba1b039d40e34e48be10142245be65d5fbc5c57c4f20f6af133ef2.FollowingRequestBuilder) {
    return ic8094ec71dba1b039d40e34e48be10142245be65d5fbc5c57c4f20f6af133ef2.NewFollowingRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FollowingById provides operations to manage the following property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) FollowingById(id string)(*i619a3607ef091cfc613c8064ad8c26e2c851e6b29b5a8de1c8ae6ad11cfec472.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i619a3607ef091cfc613c8064ad8c26e2c851e6b29b5a8de1c8ae6ad11cfec472.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get a collection of drives available for this user. Read-only.
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
func (m *DriveItemRequestBuilder) Items()(*iaaf4d433dad634780d17407e36381b8dbd65675f4ec7d040609d686c17e54c3c.ItemsRequestBuilder) {
    return iaaf4d433dad634780d17407e36381b8dbd65675f4ec7d040609d686c17e54c3c.NewItemsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ItemsById provides operations to manage the items property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) ItemsById(id string)(*i4dccdff9f8dcdf668768e001e5a27186fe4bf619c82ddb04a7646de52b9311fd.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i4dccdff9f8dcdf668768e001e5a27186fe4bf619c82ddb04a7646de52b9311fd.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// List provides operations to manage the list property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) List()(*if8ffbadcd14d06fa791c22b60b93c57cf0043e5bd611c797b0956a4d1c513322.ListRequestBuilder) {
    return if8ffbadcd14d06fa791c22b60b93c57cf0043e5bd611c797b0956a4d1c513322.NewListRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property drives in users
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
func (m *DriveItemRequestBuilder) Recent()(*i26bf19300aed5efb255319d7da9911edd0ae7ca03e6fdfddef4a2f6581682ca5.RecentRequestBuilder) {
    return i26bf19300aed5efb255319d7da9911edd0ae7ca03e6fdfddef4a2f6581682ca5.NewRecentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Root provides operations to manage the root property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) Root()(*iaceacb7c57c6e0cab9dfe9f00f869b5f5e1fcd6265e7550db1cab9f5e1e6953b.RootRequestBuilder) {
    return iaceacb7c57c6e0cab9dfe9f00f869b5f5e1fcd6265e7550db1cab9f5e1e6953b.NewRootRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchWithQ provides operations to call the search method.
func (m *DriveItemRequestBuilder) SearchWithQ(q *string)(*id96c6d741d51026adde933e1dc250eb7ed8ddf087f1f2a9235bcf77ad3f95e68.SearchWithQRequestBuilder) {
    return id96c6d741d51026adde933e1dc250eb7ed8ddf087f1f2a9235bcf77ad3f95e68.NewSearchWithQRequestBuilderInternal(m.pathParameters, m.requestAdapter, q);
}
// SharedWithMe provides operations to call the sharedWithMe method.
func (m *DriveItemRequestBuilder) SharedWithMe()(*ia43c04139e040bccd3c51a3e13b89c5bd9f4493722f732499431ccd547e8a502.SharedWithMeRequestBuilder) {
    return ia43c04139e040bccd3c51a3e13b89c5bd9f4493722f732499431ccd547e8a502.NewSharedWithMeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Special provides operations to manage the special property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) Special()(*ifa64565d61311e1ff82b220f60ec0b00f6578cc52d806aed6adb3e97fe080314.SpecialRequestBuilder) {
    return ifa64565d61311e1ff82b220f60ec0b00f6578cc52d806aed6adb3e97fe080314.NewSpecialRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SpecialById provides operations to manage the special property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) SpecialById(id string)(*i7707a2d5fd849a1571e51f19174fe037b6b38b3a81bed89f5c7a178732625346.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i7707a2d5fd849a1571e51f19174fe037b6b38b3a81bed89f5c7a178732625346.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
