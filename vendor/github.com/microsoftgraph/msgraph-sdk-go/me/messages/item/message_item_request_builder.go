package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i09c7e9344c4531879e9a4292bc9a67eb53529e726c1878c6c2e9e8e53916ba47 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/send"
    i2632294208677b2133bf271ae434a941bb72968036431b6303130fd27a80ddb4 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/createforward"
    i263d542cf24ae9e9ef78bc9a276ebed9e7e6590c4a4ae274bd71b3012abf44d2 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/reply"
    i2bb36c6899b162cafc6e656fe32fc77e6a563cef2a267d7fb5f954e663b9e808 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/value"
    i3933a4d952f201657cf6f03d324ef4e33d0780a4ea49f70704cf76a2c191f976 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/copy"
    i4ccf4f74cd764907f4b93d428388f1e772673c2f57b25c776ed0040e37810d8e "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/replyall"
    i5d0b9e0e9542b209896282893fed24a1ce4587679e91e205e15df0f1aeb037ff "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/attachments"
    i638c3159dd70df11db383b956f2700e3cdf9726dede7a3126289ca8deebbba33 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/multivalueextendedproperties"
    i83f515fce5f469eba14ae8086d23694bda995238a9cf6aeda9e4a77f60d46e51 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/createreplyall"
    i9a10feb98b69df8b4be73cf64661437fb8eba6c417528ea15e3fa7f97690cf45 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/createreply"
    ia12ddc34b64b94dfc125c48cb2ce271ab4c81b4e234623928b468956f0ac7218 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/singlevalueextendedproperties"
    id541c444118e749467445e6ddae5eccd7ce1bcbde5376a875eb7a6ac14639c51 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/extensions"
    iee30ea951ebbeb52a3eefc3847404284a2701aa3eee3ab7b930e1466e6871b06 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/forward"
    ifb4c655ccae06417a7b44976ed7bbc4f45bb8f773bda157b976061d66736f4c7 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/move"
    i66622497d9843fd844c973e3655196d0fcde41543efedcdb9bc3c54e70023033 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/extensions/item"
    i6e6c8882e6f07057df4c3ec98b0922cb303131fe0d4d661ef64ab6392a2c2221 "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/attachments/item"
    ib81b123557158ed5006240abf09bd6ea17d88c3dbd6ea5625f18100a9f024c8a "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/multivalueextendedproperties/item"
    ic3747645bd85467205b599ea11e617092b39c2dcfaef23e7852445e46b0d20fa "github.com/microsoftgraph/msgraph-sdk-go/me/messages/item/singlevalueextendedproperties/item"
)

// MessageItemRequestBuilder provides operations to manage the messages property of the microsoft.graph.user entity.
type MessageItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// MessageItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MessageItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// MessageItemRequestBuilderGetQueryParameters the messages in a mailbox or folder. Read-only. Nullable.
type MessageItemRequestBuilderGetQueryParameters struct {
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// MessageItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MessageItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MessageItemRequestBuilderGetQueryParameters
}
// MessageItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MessageItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) Attachments()(*i5d0b9e0e9542b209896282893fed24a1ce4587679e91e205e15df0f1aeb037ff.AttachmentsRequestBuilder) {
    return i5d0b9e0e9542b209896282893fed24a1ce4587679e91e205e15df0f1aeb037ff.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) AttachmentsById(id string)(*i6e6c8882e6f07057df4c3ec98b0922cb303131fe0d4d661ef64ab6392a2c2221.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i6e6c8882e6f07057df4c3ec98b0922cb303131fe0d4d661ef64ab6392a2c2221.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewMessageItemRequestBuilderInternal instantiates a new MessageItemRequestBuilder and sets the default values.
func NewMessageItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MessageItemRequestBuilder) {
    m := &MessageItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/messages/{message%2Did}{?%24select}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewMessageItemRequestBuilder instantiates a new MessageItemRequestBuilder and sets the default values.
func NewMessageItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MessageItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMessageItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Content provides operations to manage the media for the user entity.
func (m *MessageItemRequestBuilder) Content()(*i2bb36c6899b162cafc6e656fe32fc77e6a563cef2a267d7fb5f954e663b9e808.ContentRequestBuilder) {
    return i2bb36c6899b162cafc6e656fe32fc77e6a563cef2a267d7fb5f954e663b9e808.NewContentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Copy provides operations to call the copy method.
func (m *MessageItemRequestBuilder) Copy()(*i3933a4d952f201657cf6f03d324ef4e33d0780a4ea49f70704cf76a2c191f976.CopyRequestBuilder) {
    return i3933a4d952f201657cf6f03d324ef4e33d0780a4ea49f70704cf76a2c191f976.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property messages for me
func (m *MessageItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *MessageItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateForward provides operations to call the createForward method.
func (m *MessageItemRequestBuilder) CreateForward()(*i2632294208677b2133bf271ae434a941bb72968036431b6303130fd27a80ddb4.CreateForwardRequestBuilder) {
    return i2632294208677b2133bf271ae434a941bb72968036431b6303130fd27a80ddb4.NewCreateForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation the messages in a mailbox or folder. Read-only. Nullable.
func (m *MessageItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *MessageItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property messages in me
func (m *MessageItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Messageable, requestConfiguration *MessageItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateReply provides operations to call the createReply method.
func (m *MessageItemRequestBuilder) CreateReply()(*i9a10feb98b69df8b4be73cf64661437fb8eba6c417528ea15e3fa7f97690cf45.CreateReplyRequestBuilder) {
    return i9a10feb98b69df8b4be73cf64661437fb8eba6c417528ea15e3fa7f97690cf45.NewCreateReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateReplyAll provides operations to call the createReplyAll method.
func (m *MessageItemRequestBuilder) CreateReplyAll()(*i83f515fce5f469eba14ae8086d23694bda995238a9cf6aeda9e4a77f60d46e51.CreateReplyAllRequestBuilder) {
    return i83f515fce5f469eba14ae8086d23694bda995238a9cf6aeda9e4a77f60d46e51.NewCreateReplyAllRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Delete delete navigation property messages for me
func (m *MessageItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *MessageItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Extensions provides operations to manage the extensions property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) Extensions()(*id541c444118e749467445e6ddae5eccd7ce1bcbde5376a875eb7a6ac14639c51.ExtensionsRequestBuilder) {
    return id541c444118e749467445e6ddae5eccd7ce1bcbde5376a875eb7a6ac14639c51.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) ExtensionsById(id string)(*i66622497d9843fd844c973e3655196d0fcde41543efedcdb9bc3c54e70023033.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i66622497d9843fd844c973e3655196d0fcde41543efedcdb9bc3c54e70023033.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *MessageItemRequestBuilder) Forward()(*iee30ea951ebbeb52a3eefc3847404284a2701aa3eee3ab7b930e1466e6871b06.ForwardRequestBuilder) {
    return iee30ea951ebbeb52a3eefc3847404284a2701aa3eee3ab7b930e1466e6871b06.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the messages in a mailbox or folder. Read-only. Nullable.
func (m *MessageItemRequestBuilder) Get(ctx context.Context, requestConfiguration *MessageItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Messageable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMessageFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Messageable), nil
}
// Move provides operations to call the move method.
func (m *MessageItemRequestBuilder) Move()(*ifb4c655ccae06417a7b44976ed7bbc4f45bb8f773bda157b976061d66736f4c7.MoveRequestBuilder) {
    return ifb4c655ccae06417a7b44976ed7bbc4f45bb8f773bda157b976061d66736f4c7.NewMoveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) MultiValueExtendedProperties()(*i638c3159dd70df11db383b956f2700e3cdf9726dede7a3126289ca8deebbba33.MultiValueExtendedPropertiesRequestBuilder) {
    return i638c3159dd70df11db383b956f2700e3cdf9726dede7a3126289ca8deebbba33.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*ib81b123557158ed5006240abf09bd6ea17d88c3dbd6ea5625f18100a9f024c8a.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return ib81b123557158ed5006240abf09bd6ea17d88c3dbd6ea5625f18100a9f024c8a.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property messages in me
func (m *MessageItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Messageable, requestConfiguration *MessageItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Messageable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMessageFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Messageable), nil
}
// Reply provides operations to call the reply method.
func (m *MessageItemRequestBuilder) Reply()(*i263d542cf24ae9e9ef78bc9a276ebed9e7e6590c4a4ae274bd71b3012abf44d2.ReplyRequestBuilder) {
    return i263d542cf24ae9e9ef78bc9a276ebed9e7e6590c4a4ae274bd71b3012abf44d2.NewReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ReplyAll provides operations to call the replyAll method.
func (m *MessageItemRequestBuilder) ReplyAll()(*i4ccf4f74cd764907f4b93d428388f1e772673c2f57b25c776ed0040e37810d8e.ReplyAllRequestBuilder) {
    return i4ccf4f74cd764907f4b93d428388f1e772673c2f57b25c776ed0040e37810d8e.NewReplyAllRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Send provides operations to call the send method.
func (m *MessageItemRequestBuilder) Send()(*i09c7e9344c4531879e9a4292bc9a67eb53529e726c1878c6c2e9e8e53916ba47.SendRequestBuilder) {
    return i09c7e9344c4531879e9a4292bc9a67eb53529e726c1878c6c2e9e8e53916ba47.NewSendRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) SingleValueExtendedProperties()(*ia12ddc34b64b94dfc125c48cb2ce271ab4c81b4e234623928b468956f0ac7218.SingleValueExtendedPropertiesRequestBuilder) {
    return ia12ddc34b64b94dfc125c48cb2ce271ab4c81b4e234623928b468956f0ac7218.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*ic3747645bd85467205b599ea11e617092b39c2dcfaef23e7852445e46b0d20fa.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return ic3747645bd85467205b599ea11e617092b39c2dcfaef23e7852445e46b0d20fa.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
