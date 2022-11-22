package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i091165ee47a53365087a9b57bf31e9ca719ff7771179173713ec1e124a4a753d "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/following"
    i0dae3c20cb0fc75243bf5f3fe7a800224db7512d907490fbb718cb2eecffc1f6 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items"
    i14f7d808b16988af7c461c10ab368dc33cf12f6dd059a02532983cbe76ecb3cf "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/list"
    i71a78fc710509695f355c29c03de4fc7f02c5982fcd79825f756f25d45ea9021 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/root"
    iaf44c83cc64b40a51ac635afee8863205e793ee7c58a4741990f4aed67d74b06 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/sharedwithme"
    ib3a7e73bc5ddfc97d824cd78d5150cb9dfa843ede58050738ee6654be3931ce2 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/special"
    ib9ba40d643265a3cd641fc287bf1cdfba67d4ecd3b7e85f186bdc1e2719a8508 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/bundles"
    ica2191d69f833bfcc9010b770b2c97d0bd50d83ae2a22555252325d7ee3dd360 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/recent"
    id426c78d772b4fa2958336706b8ba9c798bde3fda24c6e8b99f3022f31b43fa5 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/searchwithq"
    i02963047ed83bfaaa492b881f396af2b691df3601a94c1338b7f279004d8950c "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/special/item"
    i904578fd49676fdd39eec5e2a3b32acf6b122da24f860ba553c0f41b6d894de0 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/items/item"
    ie35b5069eeea0359b3dc0a8264d4aa578446aa78e86d3bf1873094d7bfbafc39 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/bundles/item"
    ie5906441abbc73187c4eeab57ee5b669993d8fe5f2b96ae57c4a1f3ea9cbbfd0 "github.com/microsoftgraph/msgraph-sdk-go/me/drives/item/following/item"
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
func (m *DriveItemRequestBuilder) Bundles()(*ib9ba40d643265a3cd641fc287bf1cdfba67d4ecd3b7e85f186bdc1e2719a8508.BundlesRequestBuilder) {
    return ib9ba40d643265a3cd641fc287bf1cdfba67d4ecd3b7e85f186bdc1e2719a8508.NewBundlesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BundlesById provides operations to manage the bundles property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) BundlesById(id string)(*ie35b5069eeea0359b3dc0a8264d4aa578446aa78e86d3bf1873094d7bfbafc39.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return ie35b5069eeea0359b3dc0a8264d4aa578446aa78e86d3bf1873094d7bfbafc39.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewDriveItemRequestBuilderInternal instantiates a new DriveItemRequestBuilder and sets the default values.
func NewDriveItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DriveItemRequestBuilder) {
    m := &DriveItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/drives/{drive%2Did}{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property drives for me
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
// CreatePatchRequestInformation update the navigation property drives in me
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
// Delete delete navigation property drives for me
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
func (m *DriveItemRequestBuilder) Following()(*i091165ee47a53365087a9b57bf31e9ca719ff7771179173713ec1e124a4a753d.FollowingRequestBuilder) {
    return i091165ee47a53365087a9b57bf31e9ca719ff7771179173713ec1e124a4a753d.NewFollowingRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FollowingById provides operations to manage the following property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) FollowingById(id string)(*ie5906441abbc73187c4eeab57ee5b669993d8fe5f2b96ae57c4a1f3ea9cbbfd0.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return ie5906441abbc73187c4eeab57ee5b669993d8fe5f2b96ae57c4a1f3ea9cbbfd0.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
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
func (m *DriveItemRequestBuilder) Items()(*i0dae3c20cb0fc75243bf5f3fe7a800224db7512d907490fbb718cb2eecffc1f6.ItemsRequestBuilder) {
    return i0dae3c20cb0fc75243bf5f3fe7a800224db7512d907490fbb718cb2eecffc1f6.NewItemsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ItemsById provides operations to manage the items property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) ItemsById(id string)(*i904578fd49676fdd39eec5e2a3b32acf6b122da24f860ba553c0f41b6d894de0.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i904578fd49676fdd39eec5e2a3b32acf6b122da24f860ba553c0f41b6d894de0.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// List provides operations to manage the list property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) List()(*i14f7d808b16988af7c461c10ab368dc33cf12f6dd059a02532983cbe76ecb3cf.ListRequestBuilder) {
    return i14f7d808b16988af7c461c10ab368dc33cf12f6dd059a02532983cbe76ecb3cf.NewListRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property drives in me
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
func (m *DriveItemRequestBuilder) Recent()(*ica2191d69f833bfcc9010b770b2c97d0bd50d83ae2a22555252325d7ee3dd360.RecentRequestBuilder) {
    return ica2191d69f833bfcc9010b770b2c97d0bd50d83ae2a22555252325d7ee3dd360.NewRecentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Root provides operations to manage the root property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) Root()(*i71a78fc710509695f355c29c03de4fc7f02c5982fcd79825f756f25d45ea9021.RootRequestBuilder) {
    return i71a78fc710509695f355c29c03de4fc7f02c5982fcd79825f756f25d45ea9021.NewRootRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchWithQ provides operations to call the search method.
func (m *DriveItemRequestBuilder) SearchWithQ(q *string)(*id426c78d772b4fa2958336706b8ba9c798bde3fda24c6e8b99f3022f31b43fa5.SearchWithQRequestBuilder) {
    return id426c78d772b4fa2958336706b8ba9c798bde3fda24c6e8b99f3022f31b43fa5.NewSearchWithQRequestBuilderInternal(m.pathParameters, m.requestAdapter, q);
}
// SharedWithMe provides operations to call the sharedWithMe method.
func (m *DriveItemRequestBuilder) SharedWithMe()(*iaf44c83cc64b40a51ac635afee8863205e793ee7c58a4741990f4aed67d74b06.SharedWithMeRequestBuilder) {
    return iaf44c83cc64b40a51ac635afee8863205e793ee7c58a4741990f4aed67d74b06.NewSharedWithMeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Special provides operations to manage the special property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) Special()(*ib3a7e73bc5ddfc97d824cd78d5150cb9dfa843ede58050738ee6654be3931ce2.SpecialRequestBuilder) {
    return ib3a7e73bc5ddfc97d824cd78d5150cb9dfa843ede58050738ee6654be3931ce2.NewSpecialRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SpecialById provides operations to manage the special property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) SpecialById(id string)(*i02963047ed83bfaaa492b881f396af2b691df3601a94c1338b7f279004d8950c.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i02963047ed83bfaaa492b881f396af2b691df3601a94c1338b7f279004d8950c.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
