package onenote

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i392b78ec87e6862ef0bb0965913ffdca97eaf69084704fd05b09f2a82f718e11 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/onenote/sectiongroups"
    i3e5ab660d6ac3db9593f97ea896aa48566cf72944732b076c4a9d061912a6c04 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/onenote/notebooks"
    i56afdd94701540796ba18ff7485293b0edb3d0ff46c482e936c1a66059939a7d "github.com/microsoftgraph/msgraph-sdk-go/groups/item/onenote/resources"
    ic4302b238aa48da7d88cf2ada9f17f6eed828e36e8f54efabf0292996911ef1a "github.com/microsoftgraph/msgraph-sdk-go/groups/item/onenote/operations"
    id5409a02c0495e9c5e6b38f89b83f82f93291dc5e79dd29c377d082d3a67b973 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/onenote/sections"
    id7a0200daa8eeafd045771e775fa3934eb9a171334b734bc94cdc8222f05a916 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/onenote/pages"
    i032333af94d0562d937990b412828a992cf6dcd9300441d921dd966122896df0 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/onenote/sectiongroups/item"
    i0a0eb5562b0dd916c680d0a554cb25b167a2ebe18945dc353bdaa39a7ec4fa14 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/onenote/pages/item"
    i15f518181fda87a8cdf2c10fc20b18c70bacacca5e08a1dc09a45eb124da89ad "github.com/microsoftgraph/msgraph-sdk-go/groups/item/onenote/operations/item"
    i18fb32726d606974074adcc29ef533aa0b798a82757b5717d4487c0bec214a39 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/onenote/sections/item"
    ic982193c51308e3c875bd915544c104ad98988f51026864216f9494619169091 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/onenote/resources/item"
    icf88544b86cbd77f33fc028168410031e602c247e6524b0ec76c9ec5301c1b8c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/onenote/notebooks/item"
)

// OnenoteRequestBuilder provides operations to manage the onenote property of the microsoft.graph.group entity.
type OnenoteRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// OnenoteRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OnenoteRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// OnenoteRequestBuilderGetQueryParameters get onenote from groups
type OnenoteRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// OnenoteRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OnenoteRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *OnenoteRequestBuilderGetQueryParameters
}
// OnenoteRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OnenoteRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewOnenoteRequestBuilderInternal instantiates a new OnenoteRequestBuilder and sets the default values.
func NewOnenoteRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OnenoteRequestBuilder) {
    m := &OnenoteRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/onenote{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewOnenoteRequestBuilder instantiates a new OnenoteRequestBuilder and sets the default values.
func NewOnenoteRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OnenoteRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewOnenoteRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property onenote for groups
func (m *OnenoteRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *OnenoteRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation get onenote from groups
func (m *OnenoteRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *OnenoteRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property onenote in groups
func (m *OnenoteRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Onenoteable, requestConfiguration *OnenoteRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property onenote for groups
func (m *OnenoteRequestBuilder) Delete(ctx context.Context, requestConfiguration *OnenoteRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get onenote from groups
func (m *OnenoteRequestBuilder) Get(ctx context.Context, requestConfiguration *OnenoteRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Onenoteable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOnenoteFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Onenoteable), nil
}
// Notebooks provides operations to manage the notebooks property of the microsoft.graph.onenote entity.
func (m *OnenoteRequestBuilder) Notebooks()(*i3e5ab660d6ac3db9593f97ea896aa48566cf72944732b076c4a9d061912a6c04.NotebooksRequestBuilder) {
    return i3e5ab660d6ac3db9593f97ea896aa48566cf72944732b076c4a9d061912a6c04.NewNotebooksRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NotebooksById provides operations to manage the notebooks property of the microsoft.graph.onenote entity.
func (m *OnenoteRequestBuilder) NotebooksById(id string)(*icf88544b86cbd77f33fc028168410031e602c247e6524b0ec76c9ec5301c1b8c.NotebookItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["notebook%2Did"] = id
    }
    return icf88544b86cbd77f33fc028168410031e602c247e6524b0ec76c9ec5301c1b8c.NewNotebookItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Operations provides operations to manage the operations property of the microsoft.graph.onenote entity.
func (m *OnenoteRequestBuilder) Operations()(*ic4302b238aa48da7d88cf2ada9f17f6eed828e36e8f54efabf0292996911ef1a.OperationsRequestBuilder) {
    return ic4302b238aa48da7d88cf2ada9f17f6eed828e36e8f54efabf0292996911ef1a.NewOperationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OperationsById provides operations to manage the operations property of the microsoft.graph.onenote entity.
func (m *OnenoteRequestBuilder) OperationsById(id string)(*i15f518181fda87a8cdf2c10fc20b18c70bacacca5e08a1dc09a45eb124da89ad.OnenoteOperationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["onenoteOperation%2Did"] = id
    }
    return i15f518181fda87a8cdf2c10fc20b18c70bacacca5e08a1dc09a45eb124da89ad.NewOnenoteOperationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Pages provides operations to manage the pages property of the microsoft.graph.onenote entity.
func (m *OnenoteRequestBuilder) Pages()(*id7a0200daa8eeafd045771e775fa3934eb9a171334b734bc94cdc8222f05a916.PagesRequestBuilder) {
    return id7a0200daa8eeafd045771e775fa3934eb9a171334b734bc94cdc8222f05a916.NewPagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PagesById provides operations to manage the pages property of the microsoft.graph.onenote entity.
func (m *OnenoteRequestBuilder) PagesById(id string)(*i0a0eb5562b0dd916c680d0a554cb25b167a2ebe18945dc353bdaa39a7ec4fa14.OnenotePageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["onenotePage%2Did"] = id
    }
    return i0a0eb5562b0dd916c680d0a554cb25b167a2ebe18945dc353bdaa39a7ec4fa14.NewOnenotePageItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property onenote in groups
func (m *OnenoteRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Onenoteable, requestConfiguration *OnenoteRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Onenoteable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOnenoteFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Onenoteable), nil
}
// Resources provides operations to manage the resources property of the microsoft.graph.onenote entity.
func (m *OnenoteRequestBuilder) Resources()(*i56afdd94701540796ba18ff7485293b0edb3d0ff46c482e936c1a66059939a7d.ResourcesRequestBuilder) {
    return i56afdd94701540796ba18ff7485293b0edb3d0ff46c482e936c1a66059939a7d.NewResourcesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ResourcesById provides operations to manage the resources property of the microsoft.graph.onenote entity.
func (m *OnenoteRequestBuilder) ResourcesById(id string)(*ic982193c51308e3c875bd915544c104ad98988f51026864216f9494619169091.OnenoteResourceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["onenoteResource%2Did"] = id
    }
    return ic982193c51308e3c875bd915544c104ad98988f51026864216f9494619169091.NewOnenoteResourceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SectionGroups provides operations to manage the sectionGroups property of the microsoft.graph.onenote entity.
func (m *OnenoteRequestBuilder) SectionGroups()(*i392b78ec87e6862ef0bb0965913ffdca97eaf69084704fd05b09f2a82f718e11.SectionGroupsRequestBuilder) {
    return i392b78ec87e6862ef0bb0965913ffdca97eaf69084704fd05b09f2a82f718e11.NewSectionGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SectionGroupsById provides operations to manage the sectionGroups property of the microsoft.graph.onenote entity.
func (m *OnenoteRequestBuilder) SectionGroupsById(id string)(*i032333af94d0562d937990b412828a992cf6dcd9300441d921dd966122896df0.SectionGroupItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["sectionGroup%2Did"] = id
    }
    return i032333af94d0562d937990b412828a992cf6dcd9300441d921dd966122896df0.NewSectionGroupItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Sections provides operations to manage the sections property of the microsoft.graph.onenote entity.
func (m *OnenoteRequestBuilder) Sections()(*id5409a02c0495e9c5e6b38f89b83f82f93291dc5e79dd29c377d082d3a67b973.SectionsRequestBuilder) {
    return id5409a02c0495e9c5e6b38f89b83f82f93291dc5e79dd29c377d082d3a67b973.NewSectionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SectionsById provides operations to manage the sections property of the microsoft.graph.onenote entity.
func (m *OnenoteRequestBuilder) SectionsById(id string)(*i18fb32726d606974074adcc29ef533aa0b798a82757b5717d4487c0bec214a39.OnenoteSectionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["onenoteSection%2Did"] = id
    }
    return i18fb32726d606974074adcc29ef533aa0b798a82757b5717d4487c0bec214a39.NewOnenoteSectionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
