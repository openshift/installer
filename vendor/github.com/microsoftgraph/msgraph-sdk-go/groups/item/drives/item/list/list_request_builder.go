package list

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0c5fe127e9c5262e799b334489225b8b473f658d635bf8a2f9ea988592e40c73 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/list/operations"
    i3578d9428d0707ea2bec7d6551526c78f173c0439f2efaf932ae98769c033a32 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/list/contenttypes"
    i4d1a7c58922a3780a3960a688fcba637b3958ac4b810c7ed3cef37800edcf59b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/list/items"
    i9d9153a76ee2e7e0ad524d49537e6d7c83594661c33143bd1d26f908da0c026c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/list/columns"
    i9fe3aa331fee615703c3de231e4be421e9d0696764a7336de7b4cc10ee09f6ea "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/list/subscriptions"
    ic2276ed9745d93fb6be340f7c5acc5d401f7d9fd94179287d38c665c6139a676 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/list/drive"
    i0fa31d16e9df3001f88c2f2d477d86d01da8e058f9cc697cd639901a4a2cb7d9 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/list/operations/item"
    i8adb960b5e027d83cfd4d2a17b9f81e8ad35943e5ddca93e6088cf52d91563a4 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/list/items/item"
    i9e9296ac6f47bb3e815ee7414f0731cf71c700e247dd839e7d4ba5a567255fd8 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/list/contenttypes/item"
    i9f1484543b641257061220b1573e626a4861a5febd78f2eedb10dac0a33620c3 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/list/columns/item"
    ib0884420617b51c22b37f0f0db204956f442f838d747b3752ca7df43f4e9618c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/drives/item/list/subscriptions/item"
)

// ListRequestBuilder provides operations to manage the list property of the microsoft.graph.drive entity.
type ListRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ListRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ListRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ListRequestBuilderGetQueryParameters for drives in SharePoint, the underlying document library list. Read-only. Nullable.
type ListRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ListRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ListRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ListRequestBuilderGetQueryParameters
}
// ListRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ListRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Columns provides operations to manage the columns property of the microsoft.graph.list entity.
func (m *ListRequestBuilder) Columns()(*i9d9153a76ee2e7e0ad524d49537e6d7c83594661c33143bd1d26f908da0c026c.ColumnsRequestBuilder) {
    return i9d9153a76ee2e7e0ad524d49537e6d7c83594661c33143bd1d26f908da0c026c.NewColumnsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ColumnsById provides operations to manage the columns property of the microsoft.graph.list entity.
func (m *ListRequestBuilder) ColumnsById(id string)(*i9f1484543b641257061220b1573e626a4861a5febd78f2eedb10dac0a33620c3.ColumnDefinitionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["columnDefinition%2Did"] = id
    }
    return i9f1484543b641257061220b1573e626a4861a5febd78f2eedb10dac0a33620c3.NewColumnDefinitionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewListRequestBuilderInternal instantiates a new ListRequestBuilder and sets the default values.
func NewListRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ListRequestBuilder) {
    m := &ListRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/drives/{drive%2Did}/list{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewListRequestBuilder instantiates a new ListRequestBuilder and sets the default values.
func NewListRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ListRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewListRequestBuilderInternal(urlParams, requestAdapter)
}
// ContentTypes provides operations to manage the contentTypes property of the microsoft.graph.list entity.
func (m *ListRequestBuilder) ContentTypes()(*i3578d9428d0707ea2bec7d6551526c78f173c0439f2efaf932ae98769c033a32.ContentTypesRequestBuilder) {
    return i3578d9428d0707ea2bec7d6551526c78f173c0439f2efaf932ae98769c033a32.NewContentTypesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ContentTypesById provides operations to manage the contentTypes property of the microsoft.graph.list entity.
func (m *ListRequestBuilder) ContentTypesById(id string)(*i9e9296ac6f47bb3e815ee7414f0731cf71c700e247dd839e7d4ba5a567255fd8.ContentTypeItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["contentType%2Did"] = id
    }
    return i9e9296ac6f47bb3e815ee7414f0731cf71c700e247dd839e7d4ba5a567255fd8.NewContentTypeItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property list for groups
func (m *ListRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *ListRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation for drives in SharePoint, the underlying document library list. Read-only. Nullable.
func (m *ListRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ListRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property list in groups
func (m *ListRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Listable, requestConfiguration *ListRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property list for groups
func (m *ListRequestBuilder) Delete(ctx context.Context, requestConfiguration *ListRequestBuilderDeleteRequestConfiguration)(error) {
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
// Drive provides operations to manage the drive property of the microsoft.graph.list entity.
func (m *ListRequestBuilder) Drive()(*ic2276ed9745d93fb6be340f7c5acc5d401f7d9fd94179287d38c665c6139a676.DriveRequestBuilder) {
    return ic2276ed9745d93fb6be340f7c5acc5d401f7d9fd94179287d38c665c6139a676.NewDriveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get for drives in SharePoint, the underlying document library list. Read-only. Nullable.
func (m *ListRequestBuilder) Get(ctx context.Context, requestConfiguration *ListRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Listable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateListFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Listable), nil
}
// Items provides operations to manage the items property of the microsoft.graph.list entity.
func (m *ListRequestBuilder) Items()(*i4d1a7c58922a3780a3960a688fcba637b3958ac4b810c7ed3cef37800edcf59b.ItemsRequestBuilder) {
    return i4d1a7c58922a3780a3960a688fcba637b3958ac4b810c7ed3cef37800edcf59b.NewItemsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ItemsById provides operations to manage the items property of the microsoft.graph.list entity.
func (m *ListRequestBuilder) ItemsById(id string)(*i8adb960b5e027d83cfd4d2a17b9f81e8ad35943e5ddca93e6088cf52d91563a4.ListItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["listItem%2Did"] = id
    }
    return i8adb960b5e027d83cfd4d2a17b9f81e8ad35943e5ddca93e6088cf52d91563a4.NewListItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Operations provides operations to manage the operations property of the microsoft.graph.list entity.
func (m *ListRequestBuilder) Operations()(*i0c5fe127e9c5262e799b334489225b8b473f658d635bf8a2f9ea988592e40c73.OperationsRequestBuilder) {
    return i0c5fe127e9c5262e799b334489225b8b473f658d635bf8a2f9ea988592e40c73.NewOperationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OperationsById provides operations to manage the operations property of the microsoft.graph.list entity.
func (m *ListRequestBuilder) OperationsById(id string)(*i0fa31d16e9df3001f88c2f2d477d86d01da8e058f9cc697cd639901a4a2cb7d9.RichLongRunningOperationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["richLongRunningOperation%2Did"] = id
    }
    return i0fa31d16e9df3001f88c2f2d477d86d01da8e058f9cc697cd639901a4a2cb7d9.NewRichLongRunningOperationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property list in groups
func (m *ListRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Listable, requestConfiguration *ListRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Listable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateListFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Listable), nil
}
// Subscriptions provides operations to manage the subscriptions property of the microsoft.graph.list entity.
func (m *ListRequestBuilder) Subscriptions()(*i9fe3aa331fee615703c3de231e4be421e9d0696764a7336de7b4cc10ee09f6ea.SubscriptionsRequestBuilder) {
    return i9fe3aa331fee615703c3de231e4be421e9d0696764a7336de7b4cc10ee09f6ea.NewSubscriptionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubscriptionsById provides operations to manage the subscriptions property of the microsoft.graph.list entity.
func (m *ListRequestBuilder) SubscriptionsById(id string)(*ib0884420617b51c22b37f0f0db204956f442f838d747b3752ca7df43f4e9618c.SubscriptionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["subscription%2Did"] = id
    }
    return ib0884420617b51c22b37f0f0db204956f442f838d747b3752ca7df43f4e9618c.NewSubscriptionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
