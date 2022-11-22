package messages

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0a1ac23d2ad4f7421bfe2ad68e1ae9861244745d3a0eb9b2339dfa4bf574093b "github.com/microsoftgraph/msgraph-sdk-go/admin/serviceannouncement/messages/favorite"
    i184768f41f49a81af236b2a9bc61c9295feb546cfe69781d725d312259ea39b5 "github.com/microsoftgraph/msgraph-sdk-go/admin/serviceannouncement/messages/markunread"
    i34e6b5bfcdf3235a943d804a8018471597527eea943054880b2e367767db738f "github.com/microsoftgraph/msgraph-sdk-go/admin/serviceannouncement/messages/markread"
    i5b2e93574c1603b0456a604836b4dcf3f5b0f24d49a3a2258788e9fa0c9b58f1 "github.com/microsoftgraph/msgraph-sdk-go/admin/serviceannouncement/messages/count"
    i6fe59a721c674e36b7fd2811e2df3f93a4dd9bf30a3958c160ecfa9b1bde6ad7 "github.com/microsoftgraph/msgraph-sdk-go/admin/serviceannouncement/messages/unfavorite"
    i7c8002d50d3d01307f13e1746331c6c44c1d8adb28af73e57a6058f41851d2c0 "github.com/microsoftgraph/msgraph-sdk-go/admin/serviceannouncement/messages/archive"
    if4711e0cd2dcd5bcd3b2fe5ff1a931184507edf584a1b39b0489f9d289e9acff "github.com/microsoftgraph/msgraph-sdk-go/admin/serviceannouncement/messages/unarchive"
)

// MessagesRequestBuilder provides operations to manage the messages property of the microsoft.graph.serviceAnnouncement entity.
type MessagesRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// MessagesRequestBuilderGetQueryParameters retrieve the serviceUpdateMessage resources from the **messages** navigation property. This operation retrieves all service update messages that exist for the tenant.
type MessagesRequestBuilderGetQueryParameters struct {
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
// MessagesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MessagesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MessagesRequestBuilderGetQueryParameters
}
// MessagesRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MessagesRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Archive provides operations to call the archive method.
func (m *MessagesRequestBuilder) Archive()(*i7c8002d50d3d01307f13e1746331c6c44c1d8adb28af73e57a6058f41851d2c0.ArchiveRequestBuilder) {
    return i7c8002d50d3d01307f13e1746331c6c44c1d8adb28af73e57a6058f41851d2c0.NewArchiveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewMessagesRequestBuilderInternal instantiates a new MessagesRequestBuilder and sets the default values.
func NewMessagesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MessagesRequestBuilder) {
    m := &MessagesRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/admin/serviceAnnouncement/messages{?%24top,%24skip,%24search,%24filter,%24count,%24orderby,%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewMessagesRequestBuilder instantiates a new MessagesRequestBuilder and sets the default values.
func NewMessagesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MessagesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMessagesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
func (m *MessagesRequestBuilder) Count()(*i5b2e93574c1603b0456a604836b4dcf3f5b0f24d49a3a2258788e9fa0c9b58f1.CountRequestBuilder) {
    return i5b2e93574c1603b0456a604836b4dcf3f5b0f24d49a3a2258788e9fa0c9b58f1.NewCountRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation retrieve the serviceUpdateMessage resources from the **messages** navigation property. This operation retrieves all service update messages that exist for the tenant.
func (m *MessagesRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *MessagesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePostRequestInformation create new navigation property to messages for admin
func (m *MessagesRequestBuilder) CreatePostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceUpdateMessageable, requestConfiguration *MessagesRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST
    requestInfo.Headers["Accept"] = "application/json"
    requestInfo.SetContentFromParsable(ctx, m.requestAdapter, "application/json", body)
    if requestConfiguration != nil {
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// Favorite provides operations to call the favorite method.
func (m *MessagesRequestBuilder) Favorite()(*i0a1ac23d2ad4f7421bfe2ad68e1ae9861244745d3a0eb9b2339dfa4bf574093b.FavoriteRequestBuilder) {
    return i0a1ac23d2ad4f7421bfe2ad68e1ae9861244745d3a0eb9b2339dfa4bf574093b.NewFavoriteRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get retrieve the serviceUpdateMessage resources from the **messages** navigation property. This operation retrieves all service update messages that exist for the tenant.
func (m *MessagesRequestBuilder) Get(ctx context.Context, requestConfiguration *MessagesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceUpdateMessageCollectionResponseable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateServiceUpdateMessageCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceUpdateMessageCollectionResponseable), nil
}
// MarkRead provides operations to call the markRead method.
func (m *MessagesRequestBuilder) MarkRead()(*i34e6b5bfcdf3235a943d804a8018471597527eea943054880b2e367767db738f.MarkReadRequestBuilder) {
    return i34e6b5bfcdf3235a943d804a8018471597527eea943054880b2e367767db738f.NewMarkReadRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MarkUnread provides operations to call the markUnread method.
func (m *MessagesRequestBuilder) MarkUnread()(*i184768f41f49a81af236b2a9bc61c9295feb546cfe69781d725d312259ea39b5.MarkUnreadRequestBuilder) {
    return i184768f41f49a81af236b2a9bc61c9295feb546cfe69781d725d312259ea39b5.NewMarkUnreadRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Post create new navigation property to messages for admin
func (m *MessagesRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceUpdateMessageable, requestConfiguration *MessagesRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceUpdateMessageable, error) {
    requestInfo, err := m.CreatePostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateServiceUpdateMessageFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceUpdateMessageable), nil
}
// Unarchive provides operations to call the unarchive method.
func (m *MessagesRequestBuilder) Unarchive()(*if4711e0cd2dcd5bcd3b2fe5ff1a931184507edf584a1b39b0489f9d289e9acff.UnarchiveRequestBuilder) {
    return if4711e0cd2dcd5bcd3b2fe5ff1a931184507edf584a1b39b0489f9d289e9acff.NewUnarchiveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Unfavorite provides operations to call the unfavorite method.
func (m *MessagesRequestBuilder) Unfavorite()(*i6fe59a721c674e36b7fd2811e2df3f93a4dd9bf30a3958c160ecfa9b1bde6ad7.UnfavoriteRequestBuilder) {
    return i6fe59a721c674e36b7fd2811e2df3f93a4dd9bf30a3958c160ecfa9b1bde6ad7.NewUnfavoriteRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
