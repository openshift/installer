package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i07fbdadc1f7b6865331a90a151be8842dda6386d7c04f75b3409b63bd49dfb4a "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/replyall"
    i25c6477dfeb7b47e4706252b05b5bf648bfa4496cdddc92bf916eef441283c0c "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/extensions"
    i2c8800384d482fc3961ad31946374258dcb9b3b231bac345fb5cb4a9942a8bc2 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/createreply"
    i39e1a5943feb049778078f6f73510a48767da1b828bbeb45bab1fed0bb2db6a8 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/multivalueextendedproperties"
    i47131f3973f2554586af817bfb58a26418002a3d6a54635d6df69d8fbd1ade98 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/value"
    i553ef424821113114cf1bc00225a9efe230ddc93aac8ccea07cc5de2136cb379 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/copy"
    i64834a2d56614b2ccb2a2dfc49e862606287bf27d77b1b38ab337ddef0b01e97 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/createforward"
    i6d444f26eccd9ca30baa64d941f503616edab9063dae46f04fd9dab7dca6f9c0 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/forward"
    i93ccc038f4b6b4b4fdf5758a718edbfb235d2b09795924a7702496b0bb31228e "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/move"
    i94e02ca6bbe876584f943372d42c52338e0a50742eb8982f9f877828e71aa6c0 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/reply"
    iabf9f4f32d766802607eeef47bf943c667b9d8615b64e3702e5f46d8fa33d716 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/send"
    iac839ac60e5d2e9d11c89356ddc0c42e61a2a4c877b76c029f26c389c9e852fe "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/singlevalueextendedproperties"
    ib51e4147e5d5ab0a840c6a500dcbe72b0941f7612a5c951d54a2dcfe3cbbdccb "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/attachments"
    ic533ba2bbf457c6c6a6868de93bcc711ed265f801e4e5291f5e16e2bea6cc70a "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/createreplyall"
    i05975387f8e334a8ae48d3aed340dcb26535e6d58053099415dfa09207edbb0a "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/singlevalueextendedproperties/item"
    i1d3ed8cbd0f9ab9222cdb9b103733cc535f404c80df689179f38989c88730cdd "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/attachments/item"
    ib1cf63a842d9dbe2fec8426e0ab2baf39197637847eeab9c4a22aea2f6738401 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/multivalueextendedproperties/item"
    ic161cde60b5a9213b856dfb7ac1db801cdfecd0d5586a01f852e410b8ae07e94 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item/messages/item/extensions/item"
)

// MessageItemRequestBuilder provides operations to manage the messages property of the microsoft.graph.mailFolder entity.
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
// MessageItemRequestBuilderGetQueryParameters the collection of messages in the mailFolder.
type MessageItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
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
func (m *MessageItemRequestBuilder) Attachments()(*ib51e4147e5d5ab0a840c6a500dcbe72b0941f7612a5c951d54a2dcfe3cbbdccb.AttachmentsRequestBuilder) {
    return ib51e4147e5d5ab0a840c6a500dcbe72b0941f7612a5c951d54a2dcfe3cbbdccb.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) AttachmentsById(id string)(*i1d3ed8cbd0f9ab9222cdb9b103733cc535f404c80df689179f38989c88730cdd.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i1d3ed8cbd0f9ab9222cdb9b103733cc535f404c80df689179f38989c88730cdd.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewMessageItemRequestBuilderInternal instantiates a new MessageItemRequestBuilder and sets the default values.
func NewMessageItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MessageItemRequestBuilder) {
    m := &MessageItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/mailFolders/{mailFolder%2Did}/childFolders/{mailFolder%2Did1}/messages/{message%2Did}{?%24select,%24expand}";
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
func (m *MessageItemRequestBuilder) Content()(*i47131f3973f2554586af817bfb58a26418002a3d6a54635d6df69d8fbd1ade98.ContentRequestBuilder) {
    return i47131f3973f2554586af817bfb58a26418002a3d6a54635d6df69d8fbd1ade98.NewContentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Copy provides operations to call the copy method.
func (m *MessageItemRequestBuilder) Copy()(*i553ef424821113114cf1bc00225a9efe230ddc93aac8ccea07cc5de2136cb379.CopyRequestBuilder) {
    return i553ef424821113114cf1bc00225a9efe230ddc93aac8ccea07cc5de2136cb379.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) CreateForward()(*i64834a2d56614b2ccb2a2dfc49e862606287bf27d77b1b38ab337ddef0b01e97.CreateForwardRequestBuilder) {
    return i64834a2d56614b2ccb2a2dfc49e862606287bf27d77b1b38ab337ddef0b01e97.NewCreateForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateGetRequestInformation the collection of messages in the mailFolder.
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
func (m *MessageItemRequestBuilder) CreateReply()(*i2c8800384d482fc3961ad31946374258dcb9b3b231bac345fb5cb4a9942a8bc2.CreateReplyRequestBuilder) {
    return i2c8800384d482fc3961ad31946374258dcb9b3b231bac345fb5cb4a9942a8bc2.NewCreateReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateReplyAll provides operations to call the createReplyAll method.
func (m *MessageItemRequestBuilder) CreateReplyAll()(*ic533ba2bbf457c6c6a6868de93bcc711ed265f801e4e5291f5e16e2bea6cc70a.CreateReplyAllRequestBuilder) {
    return ic533ba2bbf457c6c6a6868de93bcc711ed265f801e4e5291f5e16e2bea6cc70a.NewCreateReplyAllRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) Extensions()(*i25c6477dfeb7b47e4706252b05b5bf648bfa4496cdddc92bf916eef441283c0c.ExtensionsRequestBuilder) {
    return i25c6477dfeb7b47e4706252b05b5bf648bfa4496cdddc92bf916eef441283c0c.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) ExtensionsById(id string)(*ic161cde60b5a9213b856dfb7ac1db801cdfecd0d5586a01f852e410b8ae07e94.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ic161cde60b5a9213b856dfb7ac1db801cdfecd0d5586a01f852e410b8ae07e94.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *MessageItemRequestBuilder) Forward()(*i6d444f26eccd9ca30baa64d941f503616edab9063dae46f04fd9dab7dca6f9c0.ForwardRequestBuilder) {
    return i6d444f26eccd9ca30baa64d941f503616edab9063dae46f04fd9dab7dca6f9c0.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get the collection of messages in the mailFolder.
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
func (m *MessageItemRequestBuilder) Move()(*i93ccc038f4b6b4b4fdf5758a718edbfb235d2b09795924a7702496b0bb31228e.MoveRequestBuilder) {
    return i93ccc038f4b6b4b4fdf5758a718edbfb235d2b09795924a7702496b0bb31228e.NewMoveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) MultiValueExtendedProperties()(*i39e1a5943feb049778078f6f73510a48767da1b828bbeb45bab1fed0bb2db6a8.MultiValueExtendedPropertiesRequestBuilder) {
    return i39e1a5943feb049778078f6f73510a48767da1b828bbeb45bab1fed0bb2db6a8.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*ib1cf63a842d9dbe2fec8426e0ab2baf39197637847eeab9c4a22aea2f6738401.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return ib1cf63a842d9dbe2fec8426e0ab2baf39197637847eeab9c4a22aea2f6738401.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) Reply()(*i94e02ca6bbe876584f943372d42c52338e0a50742eb8982f9f877828e71aa6c0.ReplyRequestBuilder) {
    return i94e02ca6bbe876584f943372d42c52338e0a50742eb8982f9f877828e71aa6c0.NewReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ReplyAll provides operations to call the replyAll method.
func (m *MessageItemRequestBuilder) ReplyAll()(*i07fbdadc1f7b6865331a90a151be8842dda6386d7c04f75b3409b63bd49dfb4a.ReplyAllRequestBuilder) {
    return i07fbdadc1f7b6865331a90a151be8842dda6386d7c04f75b3409b63bd49dfb4a.NewReplyAllRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Send provides operations to call the send method.
func (m *MessageItemRequestBuilder) Send()(*iabf9f4f32d766802607eeef47bf943c667b9d8615b64e3702e5f46d8fa33d716.SendRequestBuilder) {
    return iabf9f4f32d766802607eeef47bf943c667b9d8615b64e3702e5f46d8fa33d716.NewSendRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) SingleValueExtendedProperties()(*iac839ac60e5d2e9d11c89356ddc0c42e61a2a4c877b76c029f26c389c9e852fe.SingleValueExtendedPropertiesRequestBuilder) {
    return iac839ac60e5d2e9d11c89356ddc0c42e61a2a4c877b76c029f26c389c9e852fe.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i05975387f8e334a8ae48d3aed340dcb26535e6d58053099415dfa09207edbb0a.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i05975387f8e334a8ae48d3aed340dcb26535e6d58053099415dfa09207edbb0a.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
