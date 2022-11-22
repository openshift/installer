package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i1d5e0e0253764c43a7170c4dae6d6dada63a8577871ccfb6dbfa89ffcbb240cd "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/multivalueextendedproperties"
    i5134eaa0bf73502180126762f8dba15d49dd36677713ba778154e0a60592a5bd "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/forward"
    i68a4055bc841e0456a2e6639a547657de33eb9b48f925db253682599604a51e3 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/singlevalueextendedproperties"
    i75174cbcb1b898c11858f5f96b2d8847c330a38149bcd9be83e30ad77ed43c22 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/attachments"
    i87cc194ce5c8b6a9311a66fb5a22e3ba337812de9de5c80802d09bdb94b863aa "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/reply"
    ic5c310ab811a9efd80ec215346a10a39c86ebea5d6fadba8de5ac53f120fbdc1 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/extensions"
    ie143dd9268acc57eb633fcf1a20193d5695c0268106b5f598eb2c2f1e7a8a611 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/inreplyto"
    i27611e4379a9836b89c70ec5ef6aa8f3e16e6cc0deccd831272023f099a54b9a "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/multivalueextendedproperties/item"
    i8ec098680e04f7df4b9ae597c1037f94b0a4442db8e1280e22d5710ee0602d81 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/attachments/item"
    ie25007fdca9d41586081045acb368264715f3e0d2b48bebac108a93b0d11781a "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/extensions/item"
    if0f66d4e50c563fdb7e3c99373c63bc365617d543a82f7c19e2cfc84e3dcb3a8 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/threads/item/posts/item/singlevalueextendedproperties/item"
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
func (m *PostItemRequestBuilder) Attachments()(*i75174cbcb1b898c11858f5f96b2d8847c330a38149bcd9be83e30ad77ed43c22.AttachmentsRequestBuilder) {
    return i75174cbcb1b898c11858f5f96b2d8847c330a38149bcd9be83e30ad77ed43c22.NewAttachmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AttachmentsById provides operations to manage the attachments property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) AttachmentsById(id string)(*i8ec098680e04f7df4b9ae597c1037f94b0a4442db8e1280e22d5710ee0602d81.AttachmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["attachment%2Did"] = id
    }
    return i8ec098680e04f7df4b9ae597c1037f94b0a4442db8e1280e22d5710ee0602d81.NewAttachmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewPostItemRequestBuilderInternal instantiates a new PostItemRequestBuilder and sets the default values.
func NewPostItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PostItemRequestBuilder) {
    m := &PostItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/threads/{conversationThread%2Did}/posts/{post%2Did}{?%24select,%24expand}";
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
func (m *PostItemRequestBuilder) Extensions()(*ic5c310ab811a9efd80ec215346a10a39c86ebea5d6fadba8de5ac53f120fbdc1.ExtensionsRequestBuilder) {
    return ic5c310ab811a9efd80ec215346a10a39c86ebea5d6fadba8de5ac53f120fbdc1.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) ExtensionsById(id string)(*ie25007fdca9d41586081045acb368264715f3e0d2b48bebac108a93b0d11781a.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return ie25007fdca9d41586081045acb368264715f3e0d2b48bebac108a93b0d11781a.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Forward provides operations to call the forward method.
func (m *PostItemRequestBuilder) Forward()(*i5134eaa0bf73502180126762f8dba15d49dd36677713ba778154e0a60592a5bd.ForwardRequestBuilder) {
    return i5134eaa0bf73502180126762f8dba15d49dd36677713ba778154e0a60592a5bd.NewForwardRequestBuilderInternal(m.pathParameters, m.requestAdapter);
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
func (m *PostItemRequestBuilder) InReplyTo()(*ie143dd9268acc57eb633fcf1a20193d5695c0268106b5f598eb2c2f1e7a8a611.InReplyToRequestBuilder) {
    return ie143dd9268acc57eb633fcf1a20193d5695c0268106b5f598eb2c2f1e7a8a611.NewInReplyToRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) MultiValueExtendedProperties()(*i1d5e0e0253764c43a7170c4dae6d6dada63a8577871ccfb6dbfa89ffcbb240cd.MultiValueExtendedPropertiesRequestBuilder) {
    return i1d5e0e0253764c43a7170c4dae6d6dada63a8577871ccfb6dbfa89ffcbb240cd.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*i27611e4379a9836b89c70ec5ef6aa8f3e16e6cc0deccd831272023f099a54b9a.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return i27611e4379a9836b89c70ec5ef6aa8f3e16e6cc0deccd831272023f099a54b9a.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Reply provides operations to call the reply method.
func (m *PostItemRequestBuilder) Reply()(*i87cc194ce5c8b6a9311a66fb5a22e3ba337812de9de5c80802d09bdb94b863aa.ReplyRequestBuilder) {
    return i87cc194ce5c8b6a9311a66fb5a22e3ba337812de9de5c80802d09bdb94b863aa.NewReplyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) SingleValueExtendedProperties()(*i68a4055bc841e0456a2e6639a547657de33eb9b48f925db253682599604a51e3.SingleValueExtendedPropertiesRequestBuilder) {
    return i68a4055bc841e0456a2e6639a547657de33eb9b48f925db253682599604a51e3.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.post entity.
func (m *PostItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*if0f66d4e50c563fdb7e3c99373c63bc365617d543a82f7c19e2cfc84e3dcb3a8.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return if0f66d4e50c563fdb7e3c99373c63bc365617d543a82f7c19e2cfc84e3dcb3a8.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
