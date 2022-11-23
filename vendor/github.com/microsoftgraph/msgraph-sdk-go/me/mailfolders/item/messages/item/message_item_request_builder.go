package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1e31ad36dd6f6af890b8fbc16fdd724e1b79f825180e2146209442903e1536c6 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/send"
    i23d1871eaa3b34b9da4dc3013d78975116ba3c998a34b6b2f7deffdacf14481b "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/move"
    i4d76a5f9ebc4fd08800dbc2923dc874961ba25497bb7fa44bc7745d855311d6d "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/createreply"
    i525f8876c52a4815a231b5ca35b01865617a535594600f8665da18bdb784de49 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/value"
    i53c1e297b8c3f78ec3582370d703e4eb456f41af03120d845e2806968caea8a8 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/reply"
    i57c84c8fec460faeb8d03d0e3a5e76bff8ceceef2d1a73f98c319b311312185f "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/extensions"
    i7a919721ec6ea9b78407f2ee5aa295d1812c6812aeef0242dec74082f64d5fa2 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/attachments"
    i7aade3e497faeed2338d7b27d6d8b2216ae41fa1ae3717530eaf821301271040 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/replyall"
    i7e844353f50f778cff7d0fb3536d9e8900c8f5f4958bac6ec2331dabe2b1670c "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/createforward"
    i95c67b80ede95c8639d833166ec28903fde06719540ffe7742eec48fc08cd740 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/createreplyall"
    ia421f89e87068ec5411ac4214056c2997f61057486c146323e7ff0854e755cca "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/copy"
    iaf5c26a284e4a05c5a5a6cfa25d3c65b91541124829eb429f6ffc67343739406 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/singlevalueextendedproperties"
    idc279dd0f6da07bc29bf2a0f7b5d7528ee60dea48bfd2284a8f21c10e3e53e61 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/forward"
    ie70f42a3bdeda2a4c05e634ad649af47641c62c1ebdb49d12f0653327fce4e07 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/multivalueextendedproperties"
    i39886b4e20e3dea3e0ba21523a984da9b2ca189aff6561f3bc956ebf0d1c17ee "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/singlevalueextendedproperties/item"
    i93ef364c0f89b4997f731376d9fd8c3497ed63375a7efc0ee24e1a2122c20a64 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/attachments/item"
    id44a9e5325d96aaff187d75056e8b221eae0fe1190baf4c4c0318b6dfc483672 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/multivalueextendedproperties/item"
    idec392630c701cf349a9fadb8a05ca2f81715c4de9cfc77853598541c08f6920 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item/extensions/item"
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
func (m *MessageItemRequestBuilder) Attachments()(*i7a919721ec6ea9b78407f2ee5aa295d1812c6812aeef0242dec74082f64d5fa2.AttachmentsRequestBuilder) {
    return i7a919721ec6ea9b78407f2ee5aa295d1812c6812aeef0242dec74082f64d5fa2.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) AttachmentsById(id string)(*i93ef364c0f89b4997f731376d9fd8c3497ed63375a7efc0ee24e1a2122c20a64.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i93ef364c0f89b4997f731376d9fd8c3497ed63375a7efc0ee24e1a2122c20a64.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewMessageItemRequestBuilderInternal instantiates a new MessageItemRequestBuilder and sets the default values.
func NewMessageItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MessageItemRequestBuilder) {
    m := &MessageItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/mailFolders/{mailFolder%2Did}/messages/{message%2Did}{?%24select,%24expand}";
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
func (m *MessageItemRequestBuilder) Content()(*i525f8876c52a4815a231b5ca35b01865617a535594600f8665da18bdb784de49.ContentRequestBuilder) {
    return i525f8876c52a4815a231b5ca35b01865617a535594600f8665da18bdb784de49.NewContentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Copy provides operations to call the copy method.
func (m *MessageItemRequestBuilder) Copy()(*ia421f89e87068ec5411ac4214056c2997f61057486c146323e7ff0854e755cca.CopyRequestBuilder) {
    return ia421f89e87068ec5411ac4214056c2997f61057486c146323e7ff0854e755cca.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) CreateForward()(*i7e844353f50f778cff7d0fb3536d9e8900c8f5f4958bac6ec2331dabe2b1670c.CreateForwardRequestBuilder) {
    return i7e844353f50f778cff7d0fb3536d9e8900c8f5f4958bac6ec2331dabe2b1670c.NewCreateForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) CreateReply()(*i4d76a5f9ebc4fd08800dbc2923dc874961ba25497bb7fa44bc7745d855311d6d.CreateReplyRequestBuilder) {
    return i4d76a5f9ebc4fd08800dbc2923dc874961ba25497bb7fa44bc7745d855311d6d.NewCreateReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateReplyAll provides operations to call the createReplyAll method.
func (m *MessageItemRequestBuilder) CreateReplyAll()(*i95c67b80ede95c8639d833166ec28903fde06719540ffe7742eec48fc08cd740.CreateReplyAllRequestBuilder) {
    return i95c67b80ede95c8639d833166ec28903fde06719540ffe7742eec48fc08cd740.NewCreateReplyAllRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) Extensions()(*i57c84c8fec460faeb8d03d0e3a5e76bff8ceceef2d1a73f98c319b311312185f.ExtensionsRequestBuilder) {
    return i57c84c8fec460faeb8d03d0e3a5e76bff8ceceef2d1a73f98c319b311312185f.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) ExtensionsById(id string)(*idec392630c701cf349a9fadb8a05ca2f81715c4de9cfc77853598541c08f6920.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return idec392630c701cf349a9fadb8a05ca2f81715c4de9cfc77853598541c08f6920.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *MessageItemRequestBuilder) Forward()(*idc279dd0f6da07bc29bf2a0f7b5d7528ee60dea48bfd2284a8f21c10e3e53e61.ForwardRequestBuilder) {
    return idc279dd0f6da07bc29bf2a0f7b5d7528ee60dea48bfd2284a8f21c10e3e53e61.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) Move()(*i23d1871eaa3b34b9da4dc3013d78975116ba3c998a34b6b2f7deffdacf14481b.MoveRequestBuilder) {
    return i23d1871eaa3b34b9da4dc3013d78975116ba3c998a34b6b2f7deffdacf14481b.NewMoveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) MultiValueExtendedProperties()(*ie70f42a3bdeda2a4c05e634ad649af47641c62c1ebdb49d12f0653327fce4e07.MultiValueExtendedPropertiesRequestBuilder) {
    return ie70f42a3bdeda2a4c05e634ad649af47641c62c1ebdb49d12f0653327fce4e07.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*id44a9e5325d96aaff187d75056e8b221eae0fe1190baf4c4c0318b6dfc483672.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return id44a9e5325d96aaff187d75056e8b221eae0fe1190baf4c4c0318b6dfc483672.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) Reply()(*i53c1e297b8c3f78ec3582370d703e4eb456f41af03120d845e2806968caea8a8.ReplyRequestBuilder) {
    return i53c1e297b8c3f78ec3582370d703e4eb456f41af03120d845e2806968caea8a8.NewReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ReplyAll provides operations to call the replyAll method.
func (m *MessageItemRequestBuilder) ReplyAll()(*i7aade3e497faeed2338d7b27d6d8b2216ae41fa1ae3717530eaf821301271040.ReplyAllRequestBuilder) {
    return i7aade3e497faeed2338d7b27d6d8b2216ae41fa1ae3717530eaf821301271040.NewReplyAllRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Send provides operations to call the send method.
func (m *MessageItemRequestBuilder) Send()(*i1e31ad36dd6f6af890b8fbc16fdd724e1b79f825180e2146209442903e1536c6.SendRequestBuilder) {
    return i1e31ad36dd6f6af890b8fbc16fdd724e1b79f825180e2146209442903e1536c6.NewSendRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) SingleValueExtendedProperties()(*iaf5c26a284e4a05c5a5a6cfa25d3c65b91541124829eb429f6ffc67343739406.SingleValueExtendedPropertiesRequestBuilder) {
    return iaf5c26a284e4a05c5a5a6cfa25d3c65b91541124829eb429f6ffc67343739406.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i39886b4e20e3dea3e0ba21523a984da9b2ca189aff6561f3bc956ebf0d1c17ee.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i39886b4e20e3dea3e0ba21523a984da9b2ca189aff6561f3bc956ebf0d1c17ee.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
