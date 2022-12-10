package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i179944f90cbf5506a527dd2030c3cef54ae3ad4acbf50eb045881c9bd8549265 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/searchwithq"
    i2e66e6af0433d92d462a30e13b39a28fb7dee1f8ee90a449b4e31ca54ccec7f4 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/sharedwithme"
    i4205fb8ccb30dea17fb8ba72024db369db4bc1d5e093c337fcf552e72e3571c4 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/bundles"
    i53698fe752ade34739ebb810a835811924df4d95e4b93d4e5a47d1449a7c9181 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/recent"
    iae8ee83127418dc4b9186e79ac698ea94d3e30eecc4269743c54c60cf314172b "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items"
    ib2a4039148f5be30653b6dd294c5ab7aeb4ea50c5b8cb75c9ee23c4bdcd6b502 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/root"
    ic658e9fc1679f0fe33489f42aae3b91cbf7d5382b5acf0f9fa109add5f3e4b2b "github.com/microsoftgraph/msgraph-sdk-go/drives/item/special"
    id81d2362f07493ce3357d9624466a24f691c09e5f6d78806ecabf5b6f62f8ab8 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/list"
    iee46aa3df6c96b156c19a7dd468c660bb9bf053a22e578da189bbbaee28b0cdc "github.com/microsoftgraph/msgraph-sdk-go/drives/item/following"
    i1ccf6e5f8701f6397ef714c413639aae46e1c415cb2263eeee0443f0155a743b "github.com/microsoftgraph/msgraph-sdk-go/drives/item/following/item"
    iafa8d5d3e1080d6766c6b71c082fdab35a828b1d0cc84e17edcfefd65ff9b021 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/items/item"
    ibb5bf0dffefc66fdb863261eb4a45d386965bbca07ce1fa075249ee0e414aba9 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/special/item"
    id7a7f8969d72823f194725db5971b674a1dd6dda550bc0b0d5bec20719ad9ee5 "github.com/microsoftgraph/msgraph-sdk-go/drives/item/bundles/item"
)

// DriveItemRequestBuilder provides operations to manage the collection of drive entities.
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
// DriveItemRequestBuilderGetQueryParameters retrieve the properties and relationships of a Drive resource. A Drive is the top-level container for a file system, such as OneDrive or SharePoint document libraries.
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
func (m *DriveItemRequestBuilder) Bundles()(*i4205fb8ccb30dea17fb8ba72024db369db4bc1d5e093c337fcf552e72e3571c4.BundlesRequestBuilder) {
    return i4205fb8ccb30dea17fb8ba72024db369db4bc1d5e093c337fcf552e72e3571c4.NewBundlesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BundlesById provides operations to manage the bundles property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) BundlesById(id string)(*id7a7f8969d72823f194725db5971b674a1dd6dda550bc0b0d5bec20719ad9ee5.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return id7a7f8969d72823f194725db5971b674a1dd6dda550bc0b0d5bec20719ad9ee5.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewDriveItemRequestBuilderInternal instantiates a new DriveItemRequestBuilder and sets the default values.
func NewDriveItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DriveItemRequestBuilder) {
    m := &DriveItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/drives/{drive%2Did}{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete entity from drives
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
// CreateGetRequestInformation retrieve the properties and relationships of a Drive resource. A Drive is the top-level container for a file system, such as OneDrive or SharePoint document libraries.
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
// CreatePatchRequestInformation update entity in drives
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
// Delete delete entity from drives
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
func (m *DriveItemRequestBuilder) Following()(*iee46aa3df6c96b156c19a7dd468c660bb9bf053a22e578da189bbbaee28b0cdc.FollowingRequestBuilder) {
    return iee46aa3df6c96b156c19a7dd468c660bb9bf053a22e578da189bbbaee28b0cdc.NewFollowingRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FollowingById provides operations to manage the following property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) FollowingById(id string)(*i1ccf6e5f8701f6397ef714c413639aae46e1c415cb2263eeee0443f0155a743b.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i1ccf6e5f8701f6397ef714c413639aae46e1c415cb2263eeee0443f0155a743b.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get retrieve the properties and relationships of a Drive resource. A Drive is the top-level container for a file system, such as OneDrive or SharePoint document libraries.
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
func (m *DriveItemRequestBuilder) Items()(*iae8ee83127418dc4b9186e79ac698ea94d3e30eecc4269743c54c60cf314172b.ItemsRequestBuilder) {
    return iae8ee83127418dc4b9186e79ac698ea94d3e30eecc4269743c54c60cf314172b.NewItemsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ItemsById provides operations to manage the items property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) ItemsById(id string)(*iafa8d5d3e1080d6766c6b71c082fdab35a828b1d0cc84e17edcfefd65ff9b021.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return iafa8d5d3e1080d6766c6b71c082fdab35a828b1d0cc84e17edcfefd65ff9b021.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// List provides operations to manage the list property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) List()(*id81d2362f07493ce3357d9624466a24f691c09e5f6d78806ecabf5b6f62f8ab8.ListRequestBuilder) {
    return id81d2362f07493ce3357d9624466a24f691c09e5f6d78806ecabf5b6f62f8ab8.NewListRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update entity in drives
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
func (m *DriveItemRequestBuilder) Recent()(*i53698fe752ade34739ebb810a835811924df4d95e4b93d4e5a47d1449a7c9181.RecentRequestBuilder) {
    return i53698fe752ade34739ebb810a835811924df4d95e4b93d4e5a47d1449a7c9181.NewRecentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Root provides operations to manage the root property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) Root()(*ib2a4039148f5be30653b6dd294c5ab7aeb4ea50c5b8cb75c9ee23c4bdcd6b502.RootRequestBuilder) {
    return ib2a4039148f5be30653b6dd294c5ab7aeb4ea50c5b8cb75c9ee23c4bdcd6b502.NewRootRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SearchWithQ provides operations to call the search method.
func (m *DriveItemRequestBuilder) SearchWithQ(q *string)(*i179944f90cbf5506a527dd2030c3cef54ae3ad4acbf50eb045881c9bd8549265.SearchWithQRequestBuilder) {
    return i179944f90cbf5506a527dd2030c3cef54ae3ad4acbf50eb045881c9bd8549265.NewSearchWithQRequestBuilderInternal(m.pathParameters, m.requestAdapter, q);
}
// SharedWithMe provides operations to call the sharedWithMe method.
func (m *DriveItemRequestBuilder) SharedWithMe()(*i2e66e6af0433d92d462a30e13b39a28fb7dee1f8ee90a449b4e31ca54ccec7f4.SharedWithMeRequestBuilder) {
    return i2e66e6af0433d92d462a30e13b39a28fb7dee1f8ee90a449b4e31ca54ccec7f4.NewSharedWithMeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Special provides operations to manage the special property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) Special()(*ic658e9fc1679f0fe33489f42aae3b91cbf7d5382b5acf0f9fa109add5f3e4b2b.SpecialRequestBuilder) {
    return ic658e9fc1679f0fe33489f42aae3b91cbf7d5382b5acf0f9fa109add5f3e4b2b.NewSpecialRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SpecialById provides operations to manage the special property of the microsoft.graph.drive entity.
func (m *DriveItemRequestBuilder) SpecialById(id string)(*ibb5bf0dffefc66fdb863261eb4a45d386965bbca07ce1fa075249ee0e414aba9.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return ibb5bf0dffefc66fdb863261eb4a45d386965bbca07ce1fa075249ee0e414aba9.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
