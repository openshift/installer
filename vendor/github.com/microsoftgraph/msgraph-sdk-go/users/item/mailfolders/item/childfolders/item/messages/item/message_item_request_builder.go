package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i078dc43cfacb492d7c14c400159d76effb3ddc82fd0975278340abf380ca4ac0 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/createforward"
    i1e2e330cf63776672ee46cb29c0931a99d88b111c325d0eddcba469f590284ef "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/value"
    i283292ee0714f6d62978117a50ae89004640243f72c98fb77b3d3ad5d895724b "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/multivalueextendedproperties"
    i46afd8dc1f68253c5ec58fffe851f31d4606c9896eeca162270731a770163159 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/move"
    i4d24d04db91a658e0a0039257207563b182ae368b5564bc2651941a177fcec1f "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/send"
    i5a8c865ef21d4e26d59de283ac0bf75d8fa86b7ef5610dbf01669ac1a6a6426c "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/singlevalueextendedproperties"
    i601006ca5c8d4a2bd32a308480c958aa565af8b7eef89ca92de125fddcd8df37 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/createreply"
    i8bf596c0443d82bd098088cf313820c77518a4903e0996592bf7aa9bb0fc6204 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/forward"
    ia2b7d74de2d2ae40f811c9f5ca231ae8bcd6383bd423cbdadabf72fbdbfb6e62 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/replyall"
    ia7dce7bbbad8834f61a47bf8d2d11a325ca6cb60f92e7e712275f7e8c25fc693 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/reply"
    ib2ca11a51f3340d7f385f534ea5827066380a80aa53f7fb73586d468aad2cf61 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/createreplyall"
    id3a837a6f0f6a0a991f73cd274bed4374b20b7ab621b89f0c8b54728f1438971 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/copy"
    id871d310ec2178d7aee89ca6f65864dd837e1389228be08e4040bb3486b367d4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/extensions"
    if72e21ea0481f36a049b51f38179e467041e8723702278110ed61ad1e65dde50 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/attachments"
    i5b5457ef2553e64079400434de96e0e9776d9caa96dc2683cdee41584b29dcf5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/multivalueextendedproperties/item"
    ic067f029b1404c524dee30760b61043b1957ec2a9b5b5f2858424190465ad31a "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/singlevalueextendedproperties/item"
    ie43dfbb21dca01399bcbaa4461736852466b5a2f62ecd895499c91d040c14a01 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/attachments/item"
    ief0529f669908aba2ad97bcc88a9d64281988927f10a7814f2b993f88439be1b "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item/childfolders/item/messages/item/extensions/item"
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
func (m *MessageItemRequestBuilder) Attachments()(*if72e21ea0481f36a049b51f38179e467041e8723702278110ed61ad1e65dde50.AttachmentsRequestBuilder) {
    return if72e21ea0481f36a049b51f38179e467041e8723702278110ed61ad1e65dde50.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) AttachmentsById(id string)(*ie43dfbb21dca01399bcbaa4461736852466b5a2f62ecd895499c91d040c14a01.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return ie43dfbb21dca01399bcbaa4461736852466b5a2f62ecd895499c91d040c14a01.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewMessageItemRequestBuilderInternal instantiates a new MessageItemRequestBuilder and sets the default values.
func NewMessageItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MessageItemRequestBuilder) {
    m := &MessageItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/mailFolders/{mailFolder%2Did}/childFolders/{mailFolder%2Did1}/messages/{message%2Did}{?%24select,%24expand}";
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
func (m *MessageItemRequestBuilder) Content()(*i1e2e330cf63776672ee46cb29c0931a99d88b111c325d0eddcba469f590284ef.ContentRequestBuilder) {
    return i1e2e330cf63776672ee46cb29c0931a99d88b111c325d0eddcba469f590284ef.NewContentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Copy provides operations to call the copy method.
func (m *MessageItemRequestBuilder) Copy()(*id3a837a6f0f6a0a991f73cd274bed4374b20b7ab621b89f0c8b54728f1438971.CopyRequestBuilder) {
    return id3a837a6f0f6a0a991f73cd274bed4374b20b7ab621b89f0c8b54728f1438971.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) CreateForward()(*i078dc43cfacb492d7c14c400159d76effb3ddc82fd0975278340abf380ca4ac0.CreateForwardRequestBuilder) {
    return i078dc43cfacb492d7c14c400159d76effb3ddc82fd0975278340abf380ca4ac0.NewCreateForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) CreateReply()(*i601006ca5c8d4a2bd32a308480c958aa565af8b7eef89ca92de125fddcd8df37.CreateReplyRequestBuilder) {
    return i601006ca5c8d4a2bd32a308480c958aa565af8b7eef89ca92de125fddcd8df37.NewCreateReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateReplyAll provides operations to call the createReplyAll method.
func (m *MessageItemRequestBuilder) CreateReplyAll()(*ib2ca11a51f3340d7f385f534ea5827066380a80aa53f7fb73586d468aad2cf61.CreateReplyAllRequestBuilder) {
    return ib2ca11a51f3340d7f385f534ea5827066380a80aa53f7fb73586d468aad2cf61.NewCreateReplyAllRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) Extensions()(*id871d310ec2178d7aee89ca6f65864dd837e1389228be08e4040bb3486b367d4.ExtensionsRequestBuilder) {
    return id871d310ec2178d7aee89ca6f65864dd837e1389228be08e4040bb3486b367d4.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) ExtensionsById(id string)(*ief0529f669908aba2ad97bcc88a9d64281988927f10a7814f2b993f88439be1b.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ief0529f669908aba2ad97bcc88a9d64281988927f10a7814f2b993f88439be1b.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *MessageItemRequestBuilder) Forward()(*i8bf596c0443d82bd098088cf313820c77518a4903e0996592bf7aa9bb0fc6204.ForwardRequestBuilder) {
    return i8bf596c0443d82bd098088cf313820c77518a4903e0996592bf7aa9bb0fc6204.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) Move()(*i46afd8dc1f68253c5ec58fffe851f31d4606c9896eeca162270731a770163159.MoveRequestBuilder) {
    return i46afd8dc1f68253c5ec58fffe851f31d4606c9896eeca162270731a770163159.NewMoveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) MultiValueExtendedProperties()(*i283292ee0714f6d62978117a50ae89004640243f72c98fb77b3d3ad5d895724b.MultiValueExtendedPropertiesRequestBuilder) {
    return i283292ee0714f6d62978117a50ae89004640243f72c98fb77b3d3ad5d895724b.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i5b5457ef2553e64079400434de96e0e9776d9caa96dc2683cdee41584b29dcf5.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i5b5457ef2553e64079400434de96e0e9776d9caa96dc2683cdee41584b29dcf5.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
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
func (m *MessageItemRequestBuilder) Reply()(*ia7dce7bbbad8834f61a47bf8d2d11a325ca6cb60f92e7e712275f7e8c25fc693.ReplyRequestBuilder) {
    return ia7dce7bbbad8834f61a47bf8d2d11a325ca6cb60f92e7e712275f7e8c25fc693.NewReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ReplyAll provides operations to call the replyAll method.
func (m *MessageItemRequestBuilder) ReplyAll()(*ia2b7d74de2d2ae40f811c9f5ca231ae8bcd6383bd423cbdadabf72fbdbfb6e62.ReplyAllRequestBuilder) {
    return ia2b7d74de2d2ae40f811c9f5ca231ae8bcd6383bd423cbdadabf72fbdbfb6e62.NewReplyAllRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Send provides operations to call the send method.
func (m *MessageItemRequestBuilder) Send()(*i4d24d04db91a658e0a0039257207563b182ae368b5564bc2651941a177fcec1f.SendRequestBuilder) {
    return i4d24d04db91a658e0a0039257207563b182ae368b5564bc2651941a177fcec1f.NewSendRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) SingleValueExtendedProperties()(*i5a8c865ef21d4e26d59de283ac0bf75d8fa86b7ef5610dbf01669ac1a6a6426c.SingleValueExtendedPropertiesRequestBuilder) {
    return i5a8c865ef21d4e26d59de283ac0bf75d8fa86b7ef5610dbf01669ac1a6a6426c.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.message entity.
func (m *MessageItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*ic067f029b1404c524dee30760b61043b1957ec2a9b5b5f2858424190465ad31a.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return ic067f029b1404c524dee30760b61043b1957ec2a9b5b5f2858424190465ad31a.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
