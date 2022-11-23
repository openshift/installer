package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i524bf6311b24b6f141d885728d606289b8ee24cbd67c3599ccb2ec6084cb17bd "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/reply"
    i78e7a8cdddabb1f91bbd12d096fb8e9c132c2ec07d32fd988ce7552cfd89352b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/forward"
    i9884a55ef28596c5da82f51922aa27afafadf6bc4977af30511de7bda6744049 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/singlevalueextendedproperties"
    ia669f29093557a86be3b54973c4aba4e78b09891dc40e9d967d02b3938cbb08d "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/inreplyto"
    ie5a656c04186903bcf4da8c7768bc7acdefb45b9b92af375c844ab6dde924c3b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/multivalueextendedproperties"
    if00c8974308ec134f8b54a4b85c9b080451e578b97ae45f3a96f2983e98b235b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/attachments"
    if33d4aff747bbefe22c9d4c0581a32241f79f8c00b65004202a691770673b2d8 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/extensions"
    i81cf5d0e294523ab35bef33098c9a475ca60e40c2f945ab2e768942447194086 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/multivalueextendedproperties/item"
    i85387a3ef1c4b59c2db207f407864b3a4fca529975d344ab83cd21fd615de2e6 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/singlevalueextendedproperties/item"
    i98e0b63002f7d7bb02acdac8815a30c2058502950f307d6652a0b2d766b7f02b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/attachments/item"
    i9deed3605dbdfd0f6b5413628ad1be773f48a2ee7f6266b85266a4248af81477 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/extensions/item"
)

// PostItemRequestBuilder provides operations to manage the posts property of the microsoft.graph.conversationThread entity.
type PostItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// PostItemRequestBuilderGetQueryParameters get posts from groups
type PostItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// PostItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PostItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *PostItemRequestBuilderGetQueryParameters
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) Attachments()(*if00c8974308ec134f8b54a4b85c9b080451e578b97ae45f3a96f2983e98b235b.AttachmentsRequestBuilder) {
    return if00c8974308ec134f8b54a4b85c9b080451e578b97ae45f3a96f2983e98b235b.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) AttachmentsById(id string)(*i98e0b63002f7d7bb02acdac8815a30c2058502950f307d6652a0b2d766b7f02b.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i98e0b63002f7d7bb02acdac8815a30c2058502950f307d6652a0b2d766b7f02b.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewPostItemRequestBuilderInternal instantiates a new PostItemRequestBuilder and sets the default values.
func NewPostItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PostItemRequestBuilder) {
    m := &PostItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/conversations/{conversation%2Did}/threads/{conversationThread%2Did}/posts/{post%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewPostItemRequestBuilder instantiates a new PostItemRequestBuilder and sets the default values.
func NewPostItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PostItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPostItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get posts from groups
func (m *PostItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *PostItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Extensions provides operations to manage the extensions property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) Extensions()(*if33d4aff747bbefe22c9d4c0581a32241f79f8c00b65004202a691770673b2d8.ExtensionsRequestBuilder) {
    return if33d4aff747bbefe22c9d4c0581a32241f79f8c00b65004202a691770673b2d8.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) ExtensionsById(id string)(*i9deed3605dbdfd0f6b5413628ad1be773f48a2ee7f6266b85266a4248af81477.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i9deed3605dbdfd0f6b5413628ad1be773f48a2ee7f6266b85266a4248af81477.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *PostItemRequestBuilder) Forward()(*i78e7a8cdddabb1f91bbd12d096fb8e9c132c2ec07d32fd988ce7552cfd89352b.ForwardRequestBuilder) {
    return i78e7a8cdddabb1f91bbd12d096fb8e9c132c2ec07d32fd988ce7552cfd89352b.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get get posts from groups
func (m *PostItemRequestBuilder) Get(ctx context.Context, requestConfiguration *PostItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Postable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePostFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Postable), nil
}
// InReplyTo provides operations to manage the inReplyTo property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) InReplyTo()(*ia669f29093557a86be3b54973c4aba4e78b09891dc40e9d967d02b3938cbb08d.InReplyToRequestBuilder) {
    return ia669f29093557a86be3b54973c4aba4e78b09891dc40e9d967d02b3938cbb08d.NewInReplyToRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) MultiValueExtendedProperties()(*ie5a656c04186903bcf4da8c7768bc7acdefb45b9b92af375c844ab6dde924c3b.MultiValueExtendedPropertiesRequestBuilder) {
    return ie5a656c04186903bcf4da8c7768bc7acdefb45b9b92af375c844ab6dde924c3b.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i81cf5d0e294523ab35bef33098c9a475ca60e40c2f945ab2e768942447194086.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i81cf5d0e294523ab35bef33098c9a475ca60e40c2f945ab2e768942447194086.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Reply provides operations to call the reply method.
func (m *PostItemRequestBuilder) Reply()(*i524bf6311b24b6f141d885728d606289b8ee24cbd67c3599ccb2ec6084cb17bd.ReplyRequestBuilder) {
    return i524bf6311b24b6f141d885728d606289b8ee24cbd67c3599ccb2ec6084cb17bd.NewReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) SingleValueExtendedProperties()(*i9884a55ef28596c5da82f51922aa27afafadf6bc4977af30511de7bda6744049.SingleValueExtendedPropertiesRequestBuilder) {
    return i9884a55ef28596c5da82f51922aa27afafadf6bc4977af30511de7bda6744049.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i85387a3ef1c4b59c2db207f407864b3a4fca529975d344ab83cd21fd615de2e6.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i85387a3ef1c4b59c2db207f407864b3a4fca529975d344ab83cd21fd615de2e6.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
