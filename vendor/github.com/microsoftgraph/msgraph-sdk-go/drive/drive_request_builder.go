package drive

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    i33074db2b8fa9be0b6a9d94cf471e252f008c27ab9bb5a8e0ca53ac0f0210cf7 "github.com/microsoftgraph/msgraph-sdk-go/drive/items"
    i56f7e69bbb7d6f10a3721658326d4d734b2ed1edb084522652d8117a07ccaaff "github.com/microsoftgraph/msgraph-sdk-go/drive/sharedwithme"
    i62e01495b26c3ef220d18dd4a2ec05f79ddb93094e7d5a328fdfc46f27cef7b6 "github.com/microsoftgraph/msgraph-sdk-go/drive/bundles"
    i6540397555b3c095b7c6c5d4fcae5bef89c2baee0dcf1f8bba0233611b3d3f95 "github.com/microsoftgraph/msgraph-sdk-go/drive/list"
    i9316c3bd03b4192d3d91a742a5f05f6406c00c15b3b5de33a8e843dac77d53ad "github.com/microsoftgraph/msgraph-sdk-go/drive/following"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    ia82c9355f6f8d356282e506ef46bfccb28fc443731e9a05eab283791a6de8220 "github.com/microsoftgraph/msgraph-sdk-go/drive/special"
    ie13063c1cadad44d23aa1c8d0fa86d8d1520f918579ef541dec1ccdb6ed9a220 "github.com/microsoftgraph/msgraph-sdk-go/drive/searchwithq"
    ie72d87938961cb52922de029a87b3e51e1f879639a8ade8b937c65900a85220c "github.com/microsoftgraph/msgraph-sdk-go/drive/recent"
    ie85bb22170659be58e3c60d2b8b86f9a5d711197d7702af8a50ef22c9883d76a "github.com/microsoftgraph/msgraph-sdk-go/drive/root"
    i196a6eae213268ad21678594954acf3ccd7f039d9e67578273af389beab82ac2 "github.com/microsoftgraph/msgraph-sdk-go/drive/bundles/item"
    i3087ffea8c433c599359b8347be323c02e65364584d9de678689a61cc31d9d68 "github.com/microsoftgraph/msgraph-sdk-go/drive/special/item"
    i442efb2ad5fc2ef8bc831b0f78ab7125b2330176f4baf9ac638a842e8b683dc9 "github.com/microsoftgraph/msgraph-sdk-go/drive/following/item"
    ibd7a7a1ce01856c78689c5524ff307585d4847757b7dc445f97bff4e3fdb5a1b "github.com/microsoftgraph/msgraph-sdk-go/drive/items/item"
)

// DriveRequestBuilder provides operations to manage the drive singleton.
type DriveRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// DriveRequestBuilderGetQueryParameters get drive
type DriveRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DriveRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DriveRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DriveRequestBuilderGetQueryParameters
}
// DriveRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DriveRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Bundles provides operations to manage the bundles property of the microsoft.graph.drive entity.
func (m *DriveRequestBuilder) Bundles()(*i62e01495b26c3ef220d18dd4a2ec05f79ddb93094e7d5a328fdfc46f27cef7b6.BundlesRequestBuilder) {
    return i62e01495b26c3ef220d18dd4a2ec05f79ddb93094e7d5a328fdfc46f27cef7b6.NewBundlesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BundlesById provides operations to manage the bundles property of the microsoft.graph.drive entity.
func (m *DriveRequestBuilder) BundlesById(id string)(*i196a6eae213268ad21678594954acf3ccd7f039d9e67578273af389beab82ac2.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i196a6eae213268ad21678594954acf3ccd7f039d9e67578273af389beab82ac2.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewDriveRequestBuilderInternal instantiates a new DriveRequestBuilder and sets the default values.
func NewDriveRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DriveRequestBuilder) {
    m := &DriveRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/drive{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewDriveRequestBuilder instantiates a new DriveRequestBuilder and sets the default values.
func NewDriveRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DriveRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDriveRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get drive
func (m *DriveRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *DriveRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update drive
func (m *DriveRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Driveable, requestConfiguration *DriveRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Following provides operations to manage the following property of the microsoft.graph.drive entity.
func (m *DriveRequestBuilder) Following()(*i9316c3bd03b4192d3d91a742a5f05f6406c00c15b3b5de33a8e843dac77d53ad.FollowingRequestBuilder) {
    return i9316c3bd03b4192d3d91a742a5f05f6406c00c15b3b5de33a8e843dac77d53ad.NewFollowingRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FollowingById provides operations to manage the following property of the microsoft.graph.drive entity.
func (m *DriveRequestBuilder) FollowingById(id string)(*i442efb2ad5fc2ef8bc831b0f78ab7125b2330176f4baf9ac638a842e8b683dc9.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i442efb2ad5fc2ef8bc831b0f78ab7125b2330176f4baf9ac638a842e8b683dc9.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get get drive
func (m *DriveRequestBuilder) Get(ctx context.Context, requestConfiguration *DriveRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Driveable, error) {
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
func (m *DriveRequestBuilder) Items()(*i33074db2b8fa9be0b6a9d94cf471e252f008c27ab9bb5a8e0ca53ac0f0210cf7.ItemsRequestBuilder) {
    return i33074db2b8fa9be0b6a9d94cf471e252f008c27ab9bb5a8e0ca53ac0f0210cf7.NewItemsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ItemsById provides operations to manage the items property of the microsoft.graph.drive entity.
func (m *DriveRequestBuilder) ItemsById(id string)(*ibd7a7a1ce01856c78689c5524ff307585d4847757b7dc445f97bff4e3fdb5a1b.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return ibd7a7a1ce01856c78689c5524ff307585d4847757b7dc445f97bff4e3fdb5a1b.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// List provides operations to manage the list property of the microsoft.graph.drive entity.
func (m *DriveRequestBuilder) List()(*i6540397555b3c095b7c6c5d4fcae5bef89c2baee0dcf1f8bba0233611b3d3f95.ListRequestBuilder) {
    return i6540397555b3c095b7c6c5d4fcae5bef89c2baee0dcf1f8bba0233611b3d3f95.NewListRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update drive
func (m *DriveRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Driveable, requestConfiguration *DriveRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Driveable, error) {
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
func (m *DriveRequestBuilder) Recent()(*ie72d87938961cb52922de029a87b3e51e1f879639a8ade8b937c65900a85220c.RecentRequestBuilder) {
    return ie72d87938961cb52922de029a87b3e51e1f879639a8ade8b937c65900a85220c.NewRecentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Root provides operations to manage the root property of the microsoft.graph.drive entity.
func (m *DriveRequestBuilder) Root()(*ie85bb22170659be58e3c60d2b8b86f9a5d711197d7702af8a50ef22c9883d76a.RootRequestBuilder) {
    return ie85bb22170659be58e3c60d2b8b86f9a5d711197d7702af8a50ef22c9883d76a.NewRootRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchWithQ provides operations to call the search method.
func (m *DriveRequestBuilder) SearchWithQ(q *string)(*ie13063c1cadad44d23aa1c8d0fa86d8d1520f918579ef541dec1ccdb6ed9a220.SearchWithQRequestBuilder) {
    return ie13063c1cadad44d23aa1c8d0fa86d8d1520f918579ef541dec1ccdb6ed9a220.NewSearchWithQRequestBuilderInternal(m.pathParameters, m.requestAdapter, q);
}
// SharedWithMe provides operations to call the sharedWithMe method.
func (m *DriveRequestBuilder) SharedWithMe()(*i56f7e69bbb7d6f10a3721658326d4d734b2ed1edb084522652d8117a07ccaaff.SharedWithMeRequestBuilder) {
    return i56f7e69bbb7d6f10a3721658326d4d734b2ed1edb084522652d8117a07ccaaff.NewSharedWithMeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Special provides operations to manage the special property of the microsoft.graph.drive entity.
func (m *DriveRequestBuilder) Special()(*ia82c9355f6f8d356282e506ef46bfccb28fc443731e9a05eab283791a6de8220.SpecialRequestBuilder) {
    return ia82c9355f6f8d356282e506ef46bfccb28fc443731e9a05eab283791a6de8220.NewSpecialRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SpecialById provides operations to manage the special property of the microsoft.graph.drive entity.
func (m *DriveRequestBuilder) SpecialById(id string)(*i3087ffea8c433c599359b8347be323c02e65364584d9de678689a61cc31d9d68.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i3087ffea8c433c599359b8347be323c02e65364584d9de678689a61cc31d9d68.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
