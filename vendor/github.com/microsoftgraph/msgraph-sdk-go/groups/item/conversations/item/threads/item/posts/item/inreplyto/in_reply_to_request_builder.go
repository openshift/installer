package inreplyto

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0b7197b35866560bb75d5d0f4efa896921a6b20436b029349143d8acf9ad3267 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/inreplyto/extensions"
    i5a96cf81c3b745fcb623e502cb75bbd77276c4b7ec8789db5e8e4f2a6f262330 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/inreplyto/attachments"
    i92495245446099b92ab7c7f0065f500a4d14975a7d0b9c96573f137956ca320b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/inreplyto/multivalueextendedproperties"
    i9d94bd11d6d6bd3a2e1f8fa9a4b9aadfb3092f358c65db04ee394627108d0ce5 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/inreplyto/forward"
    ib88622d9430ce7882a6ecfd828e52de9a25082dfb3670f150a4e6d6e20463b04 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/inreplyto/singlevalueextendedproperties"
    if19c96c1013ba36d2ec558a175696439aec92252a53d9eabf93ad72118c4a530 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/inreplyto/reply"
    i0c96d2db214e9fe53246ea7d2b3be413b98c5662576cab468ddd391583915a33 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/inreplyto/singlevalueextendedproperties/item"
    i3e8cb7e4bf2b4a512a1cbb6155635c676f8c2f106acb67a77867e8fe307c6e2c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/inreplyto/attachments/item"
    ic430108688a3dd1740084b1d6f5636f80554b3bca84a2ff93f74d8950f0ff2fd "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/inreplyto/multivalueextendedproperties/item"
    ifbfba05084a93efb147200ec38527fedfa451327d8d0168fc33f1244aa12772f "github.com/microsoftgraph/msgraph-sdk-go/groups/item/conversations/item/threads/item/posts/item/inreplyto/extensions/item"
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
func (m *InReplyToRequestBuilder) Attachments()(*i5a96cf81c3b745fcb623e502cb75bbd77276c4b7ec8789db5e8e4f2a6f262330.AttachmentsRequestBuilder) {
    return i5a96cf81c3b745fcb623e502cb75bbd77276c4b7ec8789db5e8e4f2a6f262330.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.post entity.
func (m *InReplyToRequestBuilder) AttachmentsById(id string)(*i3e8cb7e4bf2b4a512a1cbb6155635c676f8c2f106acb67a77867e8fe307c6e2c.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i3e8cb7e4bf2b4a512a1cbb6155635c676f8c2f106acb67a77867e8fe307c6e2c.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewInReplyToRequestBuilderInternal instantiates a new InReplyToRequestBuilder and sets the default values.
func NewInReplyToRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*InReplyToRequestBuilder) {
    m := &InReplyToRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/conversations/{conversation%2Did}/threads/{conversationThread%2Did}/posts/{post%2Did}/inReplyTo{?%24select,%24expand}";
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
func (m *InReplyToRequestBuilder) Extensions()(*i0b7197b35866560bb75d5d0f4efa896921a6b20436b029349143d8acf9ad3267.ExtensionsRequestBuilder) {
    return i0b7197b35866560bb75d5d0f4efa896921a6b20436b029349143d8acf9ad3267.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.post entity.
func (m *InReplyToRequestBuilder) ExtensionsById(id string)(*ifbfba05084a93efb147200ec38527fedfa451327d8d0168fc33f1244aa12772f.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ifbfba05084a93efb147200ec38527fedfa451327d8d0168fc33f1244aa12772f.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *InReplyToRequestBuilder) Forward()(*i9d94bd11d6d6bd3a2e1f8fa9a4b9aadfb3092f358c65db04ee394627108d0ce5.ForwardRequestBuilder) {
    return i9d94bd11d6d6bd3a2e1f8fa9a4b9aadfb3092f358c65db04ee394627108d0ce5.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *InReplyToRequestBuilder) MultiValueExtendedProperties()(*i92495245446099b92ab7c7f0065f500a4d14975a7d0b9c96573f137956ca320b.MultiValueExtendedPropertiesRequestBuilder) {
    return i92495245446099b92ab7c7f0065f500a4d14975a7d0b9c96573f137956ca320b.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.post entity.
func (m *InReplyToRequestBuilder) MultiValueExtendedPropertiesById(id string)(*ic430108688a3dd1740084b1d6f5636f80554b3bca84a2ff93f74d8950f0ff2fd.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return ic430108688a3dd1740084b1d6f5636f80554b3bca84a2ff93f74d8950f0ff2fd.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Reply provides operations to call the reply method.
func (m *InReplyToRequestBuilder) Reply()(*if19c96c1013ba36d2ec558a175696439aec92252a53d9eabf93ad72118c4a530.ReplyRequestBuilder) {
    return if19c96c1013ba36d2ec558a175696439aec92252a53d9eabf93ad72118c4a530.NewReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.post entity.
func (m *InReplyToRequestBuilder) SingleValueExtendedProperties()(*ib88622d9430ce7882a6ecfd828e52de9a25082dfb3670f150a4e6d6e20463b04.SingleValueExtendedPropertiesRequestBuilder) {
    return ib88622d9430ce7882a6ecfd828e52de9a25082dfb3670f150a4e6d6e20463b04.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.post entity.
func (m *InReplyToRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i0c96d2db214e9fe53246ea7d2b3be413b98c5662576cab468ddd391583915a33.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i0c96d2db214e9fe53246ea7d2b3be413b98c5662576cab468ddd391583915a33.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
