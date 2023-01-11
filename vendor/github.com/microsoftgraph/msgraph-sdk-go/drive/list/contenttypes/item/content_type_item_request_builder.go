package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0b98bee81fa5a4fb2b5fa646275edac4f17fd1135e34963ffdb081964f70fa29 "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/base"
    i0e4a4f160e936ed0b701995d280a973c677a1c132f35feaca1aae10c863ea1b4 "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/ispublished"
    i2bab590cdf67b22b3d8348c3477a9ef52e970dff41681d07f85c1a7e2f18726d "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/publish"
    i350e95307b635e7cbd9b80327f06108e3c2d26f1b80f2ffe02ab9e26040919c0 "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/columnlinks"
    i37ed59fef414ed542b467a743d67ed5f7b35e828361fcee35f794973680f0273 "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/columns"
    i45c3970936281df10318dee6941b2dd766057e61fc3af1a2061b427dcc3f44b0 "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/columnpositions"
    i8a8731025cbe5699c843f367738d2c9ee2193b6d64796d2224611f93056402c4 "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/copytodefaultcontentlocation"
    iaf680a4c4cb484ea82e42a31267fda11f665a0e1c57cc2bc8d76dacd06dbc5f6 "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/unpublish"
    ibbfe9e07ff81e3a1ac4b3ed908b90bd28a86fd235c5c1fcf568ba59eabeca6ec "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/basetypes"
    idb984764a96aa1929d9736abf97514ee454b1e8782c03102f624a72dec321d85 "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/associatewithhubsites"
    i12c62140c5b78ceed5111ee524200560ab5e66ec6923f649782e9d66e3cd955e "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/columnpositions/item"
    i15ae992f8c06ad8a7ba1ec98d7e1e47ad2aadb64d711a178a018a0753d67fbcf "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/columns/item"
    i27f7157489366635a4d889a21bcf438a31aaa8b9e8f1fc16c1fe83d3b58b85aa "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/basetypes/item"
    ib8f5fc0530361dc8f60f00c441136b066b0ef8c5599697f14d1092db89ae8857 "github.com/microsoftgraph/msgraph-sdk-go/drive/list/contenttypes/item/columnlinks/item"
)

// ContentTypeItemRequestBuilder provides operations to manage the contentTypes property of the microsoft.graph.list entity.
type ContentTypeItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ContentTypeItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ContentTypeItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ContentTypeItemRequestBuilderGetQueryParameters the collection of content types present in this list.
type ContentTypeItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ContentTypeItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ContentTypeItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ContentTypeItemRequestBuilderGetQueryParameters
}
// ContentTypeItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ContentTypeItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AssociateWithHubSites provides operations to call the associateWithHubSites method.
func (m *ContentTypeItemRequestBuilder) AssociateWithHubSites()(*idb984764a96aa1929d9736abf97514ee454b1e8782c03102f624a72dec321d85.AssociateWithHubSitesRequestBuilder) {
    return idb984764a96aa1929d9736abf97514ee454b1e8782c03102f624a72dec321d85.NewAssociateWithHubSitesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Base provides operations to manage the base property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) Base()(*i0b98bee81fa5a4fb2b5fa646275edac4f17fd1135e34963ffdb081964f70fa29.BaseRequestBuilder) {
    return i0b98bee81fa5a4fb2b5fa646275edac4f17fd1135e34963ffdb081964f70fa29.NewBaseRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BaseTypes provides operations to manage the baseTypes property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) BaseTypes()(*ibbfe9e07ff81e3a1ac4b3ed908b90bd28a86fd235c5c1fcf568ba59eabeca6ec.BaseTypesRequestBuilder) {
    return ibbfe9e07ff81e3a1ac4b3ed908b90bd28a86fd235c5c1fcf568ba59eabeca6ec.NewBaseTypesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BaseTypesById provides operations to manage the baseTypes property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) BaseTypesById(id string)(*i27f7157489366635a4d889a21bcf438a31aaa8b9e8f1fc16c1fe83d3b58b85aa.ContentTypeItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["contentType%2Did1"] = id
    }
    return i27f7157489366635a4d889a21bcf438a31aaa8b9e8f1fc16c1fe83d3b58b85aa.NewContentTypeItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ColumnLinks provides operations to manage the columnLinks property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnLinks()(*i350e95307b635e7cbd9b80327f06108e3c2d26f1b80f2ffe02ab9e26040919c0.ColumnLinksRequestBuilder) {
    return i350e95307b635e7cbd9b80327f06108e3c2d26f1b80f2ffe02ab9e26040919c0.NewColumnLinksRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnLinksById provides operations to manage the columnLinks property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnLinksById(id string)(*ib8f5fc0530361dc8f60f00c441136b066b0ef8c5599697f14d1092db89ae8857.ColumnLinkItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnLink%2Did"] = id
    }
    return ib8f5fc0530361dc8f60f00c441136b066b0ef8c5599697f14d1092db89ae8857.NewColumnLinkItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ColumnPositions provides operations to manage the columnPositions property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnPositions()(*i45c3970936281df10318dee6941b2dd766057e61fc3af1a2061b427dcc3f44b0.ColumnPositionsRequestBuilder) {
    return i45c3970936281df10318dee6941b2dd766057e61fc3af1a2061b427dcc3f44b0.NewColumnPositionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnPositionsById provides operations to manage the columnPositions property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnPositionsById(id string)(*i12c62140c5b78ceed5111ee524200560ab5e66ec6923f649782e9d66e3cd955e.ColumnDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnDefinition%2Did"] = id
    }
    return i12c62140c5b78ceed5111ee524200560ab5e66ec6923f649782e9d66e3cd955e.NewColumnDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Columns provides operations to manage the columns property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) Columns()(*i37ed59fef414ed542b467a743d67ed5f7b35e828361fcee35f794973680f0273.ColumnsRequestBuilder) {
    return i37ed59fef414ed542b467a743d67ed5f7b35e828361fcee35f794973680f0273.NewColumnsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnsById provides operations to manage the columns property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnsById(id string)(*i15ae992f8c06ad8a7ba1ec98d7e1e47ad2aadb64d711a178a018a0753d67fbcf.ColumnDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnDefinition%2Did"] = id
    }
    return i15ae992f8c06ad8a7ba1ec98d7e1e47ad2aadb64d711a178a018a0753d67fbcf.NewColumnDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewContentTypeItemRequestBuilderInternal instantiates a new ContentTypeItemRequestBuilder and sets the default values.
func NewContentTypeItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ContentTypeItemRequestBuilder) {
    m := &ContentTypeItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/drive/list/contentTypes/{contentType%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewContentTypeItemRequestBuilder instantiates a new ContentTypeItemRequestBuilder and sets the default values.
func NewContentTypeItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ContentTypeItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewContentTypeItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CopyToDefaultContentLocation provides operations to call the copyToDefaultContentLocation method.
func (m *ContentTypeItemRequestBuilder) CopyToDefaultContentLocation()(*i8a8731025cbe5699c843f367738d2c9ee2193b6d64796d2224611f93056402c4.CopyToDefaultContentLocationRequestBuilder) {
    return i8a8731025cbe5699c843f367738d2c9ee2193b6d64796d2224611f93056402c4.NewCopyToDefaultContentLocationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property contentTypes for drive
func (m *ContentTypeItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *ContentTypeItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation the collection of content types present in this list.
func (m *ContentTypeItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ContentTypeItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property contentTypes in drive
func (m *ContentTypeItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ContentTypeable, requestConfiguration *ContentTypeItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property contentTypes for drive
func (m *ContentTypeItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ContentTypeItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get the collection of content types present in this list.
func (m *ContentTypeItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ContentTypeItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ContentTypeable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateContentTypeFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ContentTypeable), nil
}
// IsPublished provides operations to call the isPublished method.
func (m *ContentTypeItemRequestBuilder) IsPublished()(*i0e4a4f160e936ed0b701995d280a973c677a1c132f35feaca1aae10c863ea1b4.IsPublishedRequestBuilder) {
    return i0e4a4f160e936ed0b701995d280a973c677a1c132f35feaca1aae10c863ea1b4.NewIsPublishedRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property contentTypes in drive
func (m *ContentTypeItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ContentTypeable, requestConfiguration *ContentTypeItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ContentTypeable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateContentTypeFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ContentTypeable), nil
}
// Publish provides operations to call the publish method.
func (m *ContentTypeItemRequestBuilder) Publish()(*i2bab590cdf67b22b3d8348c3477a9ef52e970dff41681d07f85c1a7e2f18726d.PublishRequestBuilder) {
    return i2bab590cdf67b22b3d8348c3477a9ef52e970dff41681d07f85c1a7e2f18726d.NewPublishRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Unpublish provides operations to call the unpublish method.
func (m *ContentTypeItemRequestBuilder) Unpublish()(*iaf680a4c4cb484ea82e42a31267fda11f665a0e1c57cc2bc8d76dacd06dbc5f6.UnpublishRequestBuilder) {
    return iaf680a4c4cb484ea82e42a31267fda11f665a0e1c57cc2bc8d76dacd06dbc5f6.NewUnpublishRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
