package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0df9b813acca9f52cb58191b981da2631e2ce7ec1e0f39baea115430a4d91ec7 "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/columnlinks"
    i1ec18c0a4fbb8f904f34e72c387a8703eb3af7790657ffef68bd11ec2a9905ef "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/columns"
    i49bc8900bfdd8b3b1d3f1b52a9d8d36be89f0d252009782f35e12ba5c6fae675 "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/associatewithhubsites"
    i4a1ce58e21fb551499684168220c8c6eee92bb5bc6a4b99639b513c5cfe71a1b "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/basetypes"
    i553035798d555384c1d78519561934573cbdb4b0f2bd9da038ce9a71c4991d05 "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/columnpositions"
    i7f1730a5f677932cd121090203a3bc4480b7f148f3eae6dd7de05441a796ed7c "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/unpublish"
    i88052e6ec6afde3eb4c8a4e647edd6f1b4d0d4c30b1be1a17f71fdd90f644c90 "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/ispublished"
    i934cd417c891824167d128199e27f353bdac8d75faa96747cb1f2f67634a829e "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/publish"
    ibda66b96fbe7db9cdba0b0f6fb218ff35d517fc7b3ac53ebf9bdc4c45957d5c8 "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/copytodefaultcontentlocation"
    if20ce4a9c8987ec45c17d11c97d9255ead0b63c6bb3cb2a948e79d243584e4de "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/base"
    i5a14dba7acc9a67d9b71d060c0e0de130dbeaf1185def884644520254fd0a2bf "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/columns/item"
    i684c0ce66c376c68b7eda568c4a0e38ff6f9905b99e380e5c7d6735c0ce5813b "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/basetypes/item"
    i798aa8d3eb2beecb739a9f0d6422dbcd303b5e2935628722cd84d0e6fbea4bb1 "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/columnpositions/item"
    i8459e0e9fdbef8a8ad4c7a315789b0cca6cc0c5d53c4d9f685c083f3c7d8381d "github.com/microsoftgraph/msgraph-sdk-go/shares/item/list/contenttypes/item/columnlinks/item"
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
func (m *ContentTypeItemRequestBuilder) AssociateWithHubSites()(*i49bc8900bfdd8b3b1d3f1b52a9d8d36be89f0d252009782f35e12ba5c6fae675.AssociateWithHubSitesRequestBuilder) {
    return i49bc8900bfdd8b3b1d3f1b52a9d8d36be89f0d252009782f35e12ba5c6fae675.NewAssociateWithHubSitesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Base provides operations to manage the base property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) Base()(*if20ce4a9c8987ec45c17d11c97d9255ead0b63c6bb3cb2a948e79d243584e4de.BaseRequestBuilder) {
    return if20ce4a9c8987ec45c17d11c97d9255ead0b63c6bb3cb2a948e79d243584e4de.NewBaseRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BaseTypes provides operations to manage the baseTypes property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) BaseTypes()(*i4a1ce58e21fb551499684168220c8c6eee92bb5bc6a4b99639b513c5cfe71a1b.BaseTypesRequestBuilder) {
    return i4a1ce58e21fb551499684168220c8c6eee92bb5bc6a4b99639b513c5cfe71a1b.NewBaseTypesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// BaseTypesById provides operations to manage the baseTypes property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) BaseTypesById(id string)(*i684c0ce66c376c68b7eda568c4a0e38ff6f9905b99e380e5c7d6735c0ce5813b.ContentTypeItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["contentType%2Did1"] = id
    }
    return i684c0ce66c376c68b7eda568c4a0e38ff6f9905b99e380e5c7d6735c0ce5813b.NewContentTypeItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ColumnLinks provides operations to manage the columnLinks property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnLinks()(*i0df9b813acca9f52cb58191b981da2631e2ce7ec1e0f39baea115430a4d91ec7.ColumnLinksRequestBuilder) {
    return i0df9b813acca9f52cb58191b981da2631e2ce7ec1e0f39baea115430a4d91ec7.NewColumnLinksRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnLinksById provides operations to manage the columnLinks property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnLinksById(id string)(*i8459e0e9fdbef8a8ad4c7a315789b0cca6cc0c5d53c4d9f685c083f3c7d8381d.ColumnLinkItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnLink%2Did"] = id
    }
    return i8459e0e9fdbef8a8ad4c7a315789b0cca6cc0c5d53c4d9f685c083f3c7d8381d.NewColumnLinkItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ColumnPositions provides operations to manage the columnPositions property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnPositions()(*i553035798d555384c1d78519561934573cbdb4b0f2bd9da038ce9a71c4991d05.ColumnPositionsRequestBuilder) {
    return i553035798d555384c1d78519561934573cbdb4b0f2bd9da038ce9a71c4991d05.NewColumnPositionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnPositionsById provides operations to manage the columnPositions property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnPositionsById(id string)(*i798aa8d3eb2beecb739a9f0d6422dbcd303b5e2935628722cd84d0e6fbea4bb1.ColumnDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnDefinition%2Did"] = id
    }
    return i798aa8d3eb2beecb739a9f0d6422dbcd303b5e2935628722cd84d0e6fbea4bb1.NewColumnDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Columns provides operations to manage the columns property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) Columns()(*i1ec18c0a4fbb8f904f34e72c387a8703eb3af7790657ffef68bd11ec2a9905ef.ColumnsRequestBuilder) {
    return i1ec18c0a4fbb8f904f34e72c387a8703eb3af7790657ffef68bd11ec2a9905ef.NewColumnsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnsById provides operations to manage the columns property of the microsoft.graph.contentType entity.
func (m *ContentTypeItemRequestBuilder) ColumnsById(id string)(*i5a14dba7acc9a67d9b71d060c0e0de130dbeaf1185def884644520254fd0a2bf.ColumnDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnDefinition%2Did"] = id
    }
    return i5a14dba7acc9a67d9b71d060c0e0de130dbeaf1185def884644520254fd0a2bf.NewColumnDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewContentTypeItemRequestBuilderInternal instantiates a new ContentTypeItemRequestBuilder and sets the default values.
func NewContentTypeItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ContentTypeItemRequestBuilder) {
    m := &ContentTypeItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/shares/{sharedDriveItem%2Did}/list/contentTypes/{contentType%2Did}{?%24select,%24expand}";
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
func (m *ContentTypeItemRequestBuilder) CopyToDefaultContentLocation()(*ibda66b96fbe7db9cdba0b0f6fb218ff35d517fc7b3ac53ebf9bdc4c45957d5c8.CopyToDefaultContentLocationRequestBuilder) {
    return ibda66b96fbe7db9cdba0b0f6fb218ff35d517fc7b3ac53ebf9bdc4c45957d5c8.NewCopyToDefaultContentLocationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property contentTypes for shares
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
// CreatePatchRequestInformation update the navigation property contentTypes in shares
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
// Delete delete navigation property contentTypes for shares
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
func (m *ContentTypeItemRequestBuilder) IsPublished()(*i88052e6ec6afde3eb4c8a4e647edd6f1b4d0d4c30b1be1a17f71fdd90f644c90.IsPublishedRequestBuilder) {
    return i88052e6ec6afde3eb4c8a4e647edd6f1b4d0d4c30b1be1a17f71fdd90f644c90.NewIsPublishedRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Patch update the navigation property contentTypes in shares
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
func (m *ContentTypeItemRequestBuilder) Publish()(*i934cd417c891824167d128199e27f353bdac8d75faa96747cb1f2f67634a829e.PublishRequestBuilder) {
    return i934cd417c891824167d128199e27f353bdac8d75faa96747cb1f2f67634a829e.NewPublishRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Unpublish provides operations to call the unpublish method.
func (m *ContentTypeItemRequestBuilder) Unpublish()(*i7f1730a5f677932cd121090203a3bc4480b7f148f3eae6dd7de05441a796ed7c.UnpublishRequestBuilder) {
    return i7f1730a5f677932cd121090203a3bc4480b7f148f3eae6dd7de05441a796ed7c.NewUnpublishRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
