package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i03bf4b34230f15d3e490ba46ab92e65b3c86befa789b18b8d0a25dec32b8ec13 "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/singlevalueextendedproperties"
    i122ff3868c47eff36e63f27e4a5b59d68ad9c0ac6f599d0365b5b2a23664fef8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/attachments"
    i21a12f3061b5040c42eeb7a6e4e887142ebaa674b6c8e349ad4d27e4203799f5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/value"
    i21ac64590aeb4ea4de81e9f5c0bdb9af91b5615ed8258947f1983a17be31a12a "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/createreply"
    i29ffc4815578e16ff0a35471a241efa71e25f9fce417e2ea2e64b556ae6b4084 "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/copy"
    i4382614991ffa5cffd77c6a25c626b3b88a33b415d02b39342eb5f3a52cb881e "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/multivalueextendedproperties"
    i505c2e3f7a924493aa1390d23de4d48964764a251a918b6c13839c1b7c1aaf20 "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/replyall"
    i77a08d89129ea8ff7ca1c841cd1e162b665795b77f30f02ee19f6b11a86fbb8e "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/forward"
    i8ad53531e48d207d730614399c20bde311ba08c3dd03b37e69bef6ba38fd79fa "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/reply"
    i8ba7bb1c08d896d87831ebedb8fe7a65e343758cbd43c2a2dcaea7a802dad0b3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/send"
    i9a54a4fefd859216588c818688d9056c0d0261d20a194bfa8cd910fe376614ea "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/extensions"
    iabf960d8fba3af22e12d1fdca449478d5dda049dece4b4177254bff526c014ae "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/createforward"
    id70c200088dfe1c9223f533b4a73c07c53ce049857b0be9d382686560b3a0a44 "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/createreplyall"
    iee8102d3efbbfcfdd28ff8f8427c68562589a436cd6306dcc0ffdef0f46db71d "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/move"
    i3126e60880d45aa8eb3e50dbb7349d260566b68be6e5204b3bd5ce8027d87ae4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/multivalueextendedproperties/item"
    i4f930801d57aea145b04d038b3787d219e2ba9c79e633cdc6366b0f192b7751d "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/extensions/item"
    i569e8a41910c91a566eb55ee03be3b95dad61f13affae3df81e08c5093475576 "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/attachments/item"
    iab94264dd4909f265fdecd3eca0cdeb9268a7eec19a0c90bc22d0c9cb1ce414b "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item/singlevalueextendedproperties/item"
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
func (m *MessageItemRequestBuilder) Attachments()(*i122ff3868c47eff36e63f27e4a5b59d68ad9c0ac6f599d0365b5b2a23664fef8.AttachmentsRequestBuilder) {
    return i122ff3868c47eff36e63f27e4a5b59d68ad9c0ac6f599d0365b5b2a23664fef8.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) AttachmentsById(id string)(*i569e8a41910c91a566eb55ee03be3b95dad61f13affae3df81e08c5093475576.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i569e8a41910c91a566eb55ee03be3b95dad61f13affae3df81e08c5093475576.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewMessageItemRequestBuilderInternal instantiates a new MessageItemRequestBuilder and sets the default values.
func NewMessageItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MessageItemRequestBuilder) {
    m := &MessageItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/messages/{message%2Did}{?%24select}";
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
func (m *MessageItemRequestBuilder) Content()(*i21a12f3061b5040c42eeb7a6e4e887142ebaa674b6c8e349ad4d27e4203799f5.ContentRequestBuilder) {
    return i21a12f3061b5040c42eeb7a6e4e887142ebaa674b6c8e349ad4d27e4203799f5.NewContentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Copy provides operations to call the copy method.
func (m *MessageItemRequestBuilder) Copy()(*i29ffc4815578e16ff0a35471a241efa71e25f9fce417e2ea2e64b556ae6b4084.CopyRequestBuilder) {
    return i29ffc4815578e16ff0a35471a241efa71e25f9fce417e2ea2e64b556ae6b4084.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property messages for users
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
func (m *MessageItemRequestBuilder) CreateForward()(*iabf960d8fba3af22e12d1fdca449478d5dda049dece4b4177254bff526c014ae.CreateForwardRequestBuilder) {
    return iabf960d8fba3af22e12d1fdca449478d5dda049dece4b4177254bff526c014ae.NewCreateForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
// CreatePatchRequestInformation update the navigation property messages in users
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
func (m *MessageItemRequestBuilder) CreateReply()(*i21ac64590aeb4ea4de81e9f5c0bdb9af91b5615ed8258947f1983a17be31a12a.CreateReplyRequestBuilder) {
    return i21ac64590aeb4ea4de81e9f5c0bdb9af91b5615ed8258947f1983a17be31a12a.NewCreateReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateReplyAll provides operations to call the createReplyAll method.
func (m *MessageItemRequestBuilder) CreateReplyAll()(*id70c200088dfe1c9223f533b4a73c07c53ce049857b0be9d382686560b3a0a44.CreateReplyAllRequestBuilder) {
    return id70c200088dfe1c9223f533b4a73c07c53ce049857b0be9d382686560b3a0a44.NewCreateReplyAllRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Delete delete navigation property messages for users
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
func (m *MessageItemRequestBuilder) Extensions()(*i9a54a4fefd859216588c818688d9056c0d0261d20a194bfa8cd910fe376614ea.ExtensionsRequestBuilder) {
    return i9a54a4fefd859216588c818688d9056c0d0261d20a194bfa8cd910fe376614ea.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) ExtensionsById(id string)(*i4f930801d57aea145b04d038b3787d219e2ba9c79e633cdc6366b0f192b7751d.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i4f930801d57aea145b04d038b3787d219e2ba9c79e633cdc6366b0f192b7751d.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *MessageItemRequestBuilder) Forward()(*i77a08d89129ea8ff7ca1c841cd1e162b665795b77f30f02ee19f6b11a86fbb8e.ForwardRequestBuilder) {
    return i77a08d89129ea8ff7ca1c841cd1e162b665795b77f30f02ee19f6b11a86fbb8e.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) Move()(*iee8102d3efbbfcfdd28ff8f8427c68562589a436cd6306dcc0ffdef0f46db71d.MoveRequestBuilder) {
    return iee8102d3efbbfcfdd28ff8f8427c68562589a436cd6306dcc0ffdef0f46db71d.NewMoveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) MultiValueExtendedProperties()(*i4382614991ffa5cffd77c6a25c626b3b88a33b415d02b39342eb5f3a52cb881e.MultiValueExtendedPropertiesRequestBuilder) {
    return i4382614991ffa5cffd77c6a25c626b3b88a33b415d02b39342eb5f3a52cb881e.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i3126e60880d45aa8eb3e50dbb7349d260566b68be6e5204b3bd5ce8027d87ae4.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i3126e60880d45aa8eb3e50dbb7349d260566b68be6e5204b3bd5ce8027d87ae4.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property messages in users
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
func (m *MessageItemRequestBuilder) Reply()(*i8ad53531e48d207d730614399c20bde311ba08c3dd03b37e69bef6ba38fd79fa.ReplyRequestBuilder) {
    return i8ad53531e48d207d730614399c20bde311ba08c3dd03b37e69bef6ba38fd79fa.NewReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ReplyAll provides operations to call the replyAll method.
func (m *MessageItemRequestBuilder) ReplyAll()(*i505c2e3f7a924493aa1390d23de4d48964764a251a918b6c13839c1b7c1aaf20.ReplyAllRequestBuilder) {
    return i505c2e3f7a924493aa1390d23de4d48964764a251a918b6c13839c1b7c1aaf20.NewReplyAllRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Send provides operations to call the send method.
func (m *MessageItemRequestBuilder) Send()(*i8ba7bb1c08d896d87831ebedb8fe7a65e343758cbd43c2a2dcaea7a802dad0b3.SendRequestBuilder) {
    return i8ba7bb1c08d896d87831ebedb8fe7a65e343758cbd43c2a2dcaea7a802dad0b3.NewSendRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) SingleValueExtendedProperties()(*i03bf4b34230f15d3e490ba46ab92e65b3c86befa789b18b8d0a25dec32b8ec13.SingleValueExtendedPropertiesRequestBuilder) {
    return i03bf4b34230f15d3e490ba46ab92e65b3c86befa789b18b8d0a25dec32b8ec13.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*iab94264dd4909f265fdecd3eca0cdeb9268a7eec19a0c90bc22d0c9cb1ce414b.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return iab94264dd4909f265fdecd3eca0cdeb9268a7eec19a0c90bc22d0c9cb1ce414b.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
