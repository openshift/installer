package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i23a948d948931da4c2473f25bf2771887399012591e170bf2c38ffb2a3ba66c0 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages"
    i41238d4f3e944195a28128ba7c7e506c78ae7944f63ee5b99f6ccf43fc959a1f "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders"
    ia3a6564ed24f76e12dff691849f92f63255d225c1ca3736c45221460427057ed "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/move"
    iad30d3cef76cd9f7b9285b76f6f4477abd7036d9f2c35bc5ed1522f0af02c78a "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/singlevalueextendedproperties"
    id3d1b1d14374855b32699cf86c45aa2624debfc4187f8f34b6f03ceb6bcdc5ac "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messagerules"
    ief5060c34a3b1d4aad14d434f772d2eb07c6f99ebfee6ebac51aeec75a2ea7fb "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/multivalueextendedproperties"
    if16b5f206adf92e3678b287eee1edc09a825e313dddd8774b3a749501630578c "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/copy"
    i355444ee85a0200730d90e442a49ca5aee393c423237c41182037176c7a79d49 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messagerules/item"
    i90c07370d4369a4b6aee4e8f38d2071aa7e2def7f82f01fe6ce1274c5065e166 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/singlevalueextendedproperties/item"
    ia104b6ddf37a2b233498082fb22da19da091795fbd96aa1548ddd53087f11025 "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/childfolders/item"
    ieb796547151eb3ed23dc629e3abe93c6b6c3f96fa2ab5ff99d5ec88c1e8ff28a "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/multivalueextendedproperties/item"
    ife199debe7c0cb4dd460b2e4c2e997cb1496b1b2299d0849f63f7a2d76121d8e "github.com/microsoftgraph/msgraph-sdk-go/me/mailfolders/item/messages/item"
)

// MailFolderItemRequestBuilder provides operations to manage the mailFolders property of the microsoft.graph.user entity.
type MailFolderItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// MailFolderItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MailFolderItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// MailFolderItemRequestBuilderGetQueryParameters the user's mail folders. Read-only. Nullable.
type MailFolderItemRequestBuilderGetQueryParameters struct {
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// MailFolderItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MailFolderItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MailFolderItemRequestBuilderGetQueryParameters
}
// MailFolderItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MailFolderItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ChildFolders provides operations to manage the childFolders property of the microsoft.graph.mailFolder entity.
func (m *MailFolderItemRequestBuilder) ChildFolders()(*i41238d4f3e944195a28128ba7c7e506c78ae7944f63ee5b99f6ccf43fc959a1f.ChildFoldersRequestBuilder) {
    return i41238d4f3e944195a28128ba7c7e506c78ae7944f63ee5b99f6ccf43fc959a1f.NewChildFoldersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChildFoldersById provides operations to manage the childFolders property of the microsoft.graph.mailFolder entity.
func (m *MailFolderItemRequestBuilder) ChildFoldersById(id string)(*ia104b6ddf37a2b233498082fb22da19da091795fbd96aa1548ddd53087f11025.MailFolderItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["mailFolder%2Did1"] = id
    }
    return ia104b6ddf37a2b233498082fb22da19da091795fbd96aa1548ddd53087f11025.NewMailFolderItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewMailFolderItemRequestBuilderInternal instantiates a new MailFolderItemRequestBuilder and sets the default values.
func NewMailFolderItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MailFolderItemRequestBuilder) {
    m := &MailFolderItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/me/mailFolders/{mailFolder%2Did}{?%24select}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewMailFolderItemRequestBuilder instantiates a new MailFolderItemRequestBuilder and sets the default values.
func NewMailFolderItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MailFolderItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMailFolderItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Copy provides operations to call the copy method.
func (m *MailFolderItemRequestBuilder) Copy()(*if16b5f206adf92e3678b287eee1edc09a825e313dddd8774b3a749501630578c.CopyRequestBuilder) {
    return if16b5f206adf92e3678b287eee1edc09a825e313dddd8774b3a749501630578c.NewCopyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreateDeleteRequestInformation delete navigation property mailFolders for me
func (m *MailFolderItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *MailFolderItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation the user's mail folders. Read-only. Nullable.
func (m *MailFolderItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *MailFolderItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the navigation property mailFolders in me
func (m *MailFolderItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MailFolderable, requestConfiguration *MailFolderItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property mailFolders for me
func (m *MailFolderItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *MailFolderItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get the user's mail folders. Read-only. Nullable.
func (m *MailFolderItemRequestBuilder) Get(ctx context.Context, requestConfiguration *MailFolderItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MailFolderable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMailFolderFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MailFolderable), nil
}
// MessageRules provides operations to manage the messageRules property of the microsoft.graph.mailFolder entity.
func (m *MailFolderItemRequestBuilder) MessageRules()(*id3d1b1d14374855b32699cf86c45aa2624debfc4187f8f34b6f03ceb6bcdc5ac.MessageRulesRequestBuilder) {
    return id3d1b1d14374855b32699cf86c45aa2624debfc4187f8f34b6f03ceb6bcdc5ac.NewMessageRulesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MessageRulesById provides operations to manage the messageRules property of the microsoft.graph.mailFolder entity.
func (m *MailFolderItemRequestBuilder) MessageRulesById(id string)(*i355444ee85a0200730d90e442a49ca5aee393c423237c41182037176c7a79d49.MessageRuleItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["messageRule%2Did"] = id
    }
    return i355444ee85a0200730d90e442a49ca5aee393c423237c41182037176c7a79d49.NewMessageRuleItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Messages provides operations to manage the messages property of the microsoft.graph.mailFolder entity.
func (m *MailFolderItemRequestBuilder) Messages()(*i23a948d948931da4c2473f25bf2771887399012591e170bf2c38ffb2a3ba66c0.MessagesRequestBuilder) {
    return i23a948d948931da4c2473f25bf2771887399012591e170bf2c38ffb2a3ba66c0.NewMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MessagesById provides operations to manage the messages property of the microsoft.graph.mailFolder entity.
func (m *MailFolderItemRequestBuilder) MessagesById(id string)(*ife199debe7c0cb4dd460b2e4c2e997cb1496b1b2299d0849f63f7a2d76121d8e.MessageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["message%2Did"] = id
    }
    return ife199debe7c0cb4dd460b2e4c2e997cb1496b1b2299d0849f63f7a2d76121d8e.NewMessageItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Move provides operations to call the move method.
func (m *MailFolderItemRequestBuilder) Move()(*ia3a6564ed24f76e12dff691849f92f63255d225c1ca3736c45221460427057ed.MoveRequestBuilder) {
    return ia3a6564ed24f76e12dff691849f92f63255d225c1ca3736c45221460427057ed.NewMoveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedProperties provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.mailFolder entity.
func (m *MailFolderItemRequestBuilder) MultiValueExtendedProperties()(*ief5060c34a3b1d4aad14d434f772d2eb07c6f99ebfee6ebac51aeec75a2ea7fb.MultiValueExtendedPropertiesRequestBuilder) {
    return ief5060c34a3b1d4aad14d434f772d2eb07c6f99ebfee6ebac51aeec75a2ea7fb.NewMultiValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MultiValueExtendedPropertiesById provides operations to manage the multiValueExtendedProperties property of the microsoft.graph.mailFolder entity.
func (m *MailFolderItemRequestBuilder) MultiValueExtendedPropertiesById(id string)(*ieb796547151eb3ed23dc629e3abe93c6b6c3f96fa2ab5ff99d5ec88c1e8ff28a.MultiValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["multiValueLegacyExtendedProperty%2Did"] = id
    }
    return ieb796547151eb3ed23dc629e3abe93c6b6c3f96fa2ab5ff99d5ec88c1e8ff28a.NewMultiValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property mailFolders in me
func (m *MailFolderItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MailFolderable, requestConfiguration *MailFolderItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MailFolderable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMailFolderFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MailFolderable), nil
}
// SingleValueExtendedProperties provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.mailFolder entity.
func (m *MailFolderItemRequestBuilder) SingleValueExtendedProperties()(*iad30d3cef76cd9f7b9285b76f6f4477abd7036d9f2c35bc5ed1522f0af02c78a.SingleValueExtendedPropertiesRequestBuilder) {
    return iad30d3cef76cd9f7b9285b76f6f4477abd7036d9f2c35bc5ed1522f0af02c78a.NewSingleValueExtendedPropertiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SingleValueExtendedPropertiesById provides operations to manage the singleValueExtendedProperties property of the microsoft.graph.mailFolder entity.
func (m *MailFolderItemRequestBuilder) SingleValueExtendedPropertiesById(id string)(*i90c07370d4369a4b6aee4e8f38d2071aa7e2def7f82f01fe6ce1274c5065e166.SingleValueLegacyExtendedPropertyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["singleValueLegacyExtendedProperty%2Did"] = id
    }
    return i90c07370d4369a4b6aee4e8f38d2071aa7e2def7f82f01fe6ce1274c5065e166.NewSingleValueLegacyExtendedPropertyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
