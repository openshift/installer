package inreplyto

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1bb01e9daad0a1d10bc19f61790389f463801a0f48cd11c98f185e03018bfca6 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/inreplyto/multivalueextendedproperties"
    i430f1d4a835b0b99d0b13f11ccfa4df1656409bcbf6bcf41f8ea5681dcb1aa92 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/inreplyto/extensions"
    i4472ef305b18a9aa2e06ba2a54d4719d2be1f8fa6456051cb236477675c9cd03 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/inreplyto/attachments"
    ic06046cc8c7f87f655d91876ca77d67f798d85de5c44c9e4974ce552215dfb8b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/inreplyto/forward"
    ie720d4183606c7d13bf38c1fb20010a0d54dbeb3dca324af9ff5dcc638b2a44c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/inreplyto/singlevalueextendedproperties"
    if1f9571513b1bb74e42cbc5688db7636642423bb6ccc5801422a1c43cfaa0c0c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/inreplyto/reply"
    i0f01ae290b75a0c7094161d59819ebc34d4806cd13012ae529843662aa5da593 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/inreplyto/multivalueextendedproperties/item"
    i5876235ee8bac9db2eddf11773408040864058f168150e6b139e4e22f3fc5a94 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/inreplyto/extensions/item"
    i7f4e85ba5983109d680cfb442209512e035a5aadc8b8e7867a3c3959956d86ba "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/inreplyto/singlevalueextendedproperties/item"
    if6bf231e5b7a5b45b60ecc5a3d445a1512040b77f6170861144cf99bb398767d "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/inreplyto/attachments/item"
)

// InReplyToRequestBuilder provides operations to manage the inReplyTo property of the microsoft.graph.post entity.
type InReplyToRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// InReplyToRequestBuilderGetQueryParameters read-only. Supports $expand.
type InReplyToRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// InReplyToRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type InReplyToRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *InReplyToRequestBuilderGetQueryParameters
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.post entity.
func (m *InReplyToRequestBuilder) Attachments()(*i4472ef305b18a9aa2e06ba2a54d4719d2be1f8fa6456051cb236477675c9cd03.AttachmentsRequestBuilder) {
    return i4472ef305b18a9aa2e06ba2a54d4719d2be1f8fa6456051cb236477675c9cd03.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.post entity.
func (m *InReplyToRequestBuilder) AttachmentsById(id string)(*if6bf231e5b7a5b45b60ecc5a3d445a1512040b77f6170861144cf99bb398767d.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return if6bf231e5b7a5b45b60ecc5a3d445a1512040b77f6170861144cf99bb398767d.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewInReplyToRequestBuilderInternal instantiates a new InReplyToRequestBuilder and sets the default values.
func NewInReplyToRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*InReplyToRequestBuilder) {
    m := &InReplyToRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/threads/{conversationThread%2Did}/posts/{post%2Did}/inReplyTo{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewInReplyToRequestBuilder instantiates a new InReplyToRequestBuilder and sets the default values.
func NewInReplyToRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*InReplyToRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewInReplyToRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation read-only. Supports $expand.
func (m *InReplyToRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *InReplyToRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
func (m *InReplyToRequestBuilder) Extensions()(*i430f1d4a835b0b99d0b13f11ccfa4df1656409bcbf6bcf41f8ea5681dcb1aa92.ExtensionsRequestBuilder) {
    return i430f1d4a835b0b99d0b13f11ccfa4df1656409bcbf6bcf41f8ea5681dcb1aa92.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.post entity.
func (m *InReplyToRequestBuilder) ExtensionsById(id string)(*i5876235ee8bac9db2eddf11773408040864058f168150e6b139e4e22f3fc5a94.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i5876235ee8bac9db2eddf11773408040864058f168150e6b139e4e22f3fc5a94.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *InReplyToRequestBuilder) Forward()(*ic06046cc8c7f87f655d91876ca77d67f798d85de5c44c9e4974ce552215dfb8b.ForwardRequestBuilder) {
    return ic06046cc8c7f87f655d91876ca77d67f798d85de5c44c9e4974ce552215dfb8b.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get read-only. Supports $expand.
func (m *InReplyToRequestBuilder) Get(ctx context.Context, requestConfiguration *InReplyToRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Postable, error) {
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
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.post entity.
func (m *InReplyToRequestBuilder) MultiValueExtendedProperties()(*i1bb01e9daad0a1d10bc19f61790389f463801a0f48cd11c98f185e03018bfca6.MultiValueExtendedPropertiesRequestBuilder) {
    return i1bb01e9daad0a1d10bc19f61790389f463801a0f48cd11c98f185e03018bfca6.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.post entity.
func (m *InReplyToRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i0f01ae290b75a0c7094161d59819ebc34d4806cd13012ae529843662aa5da593.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i0f01ae290b75a0c7094161d59819ebc34d4806cd13012ae529843662aa5da593.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Reply provides operations to call the reply method.
func (m *InReplyToRequestBuilder) Reply()(*if1f9571513b1bb74e42cbc5688db7636642423bb6ccc5801422a1c43cfaa0c0c.ReplyRequestBuilder) {
    return if1f9571513b1bb74e42cbc5688db7636642423bb6ccc5801422a1c43cfaa0c0c.NewReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.post entity.
func (m *InReplyToRequestBuilder) SingleValueExtendedProperties()(*ie720d4183606c7d13bf38c1fb20010a0d54dbeb3dca324af9ff5dcc638b2a44c.SingleValueExtendedPropertiesRequestBuilder) {
    return ie720d4183606c7d13bf38c1fb20010a0d54dbeb3dca324af9ff5dcc638b2a44c.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.post entity.
func (m *InReplyToRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i7f4e85ba5983109d680cfb442209512e035a5aadc8b8e7867a3c3959956d86ba.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i7f4e85ba5983109d680cfb442209512e035a5aadc8b8e7867a3c3959956d86ba.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
