package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0453d8f2831f761719a2d96702e4c863ff529132a101dbe684210882f2721d5c "github.com/microsoftgraph/msgraph-sdk-go/chats/item/messages"
    i39beb613c864af7bec369a1056bf8c41cb817d04318b7282f041a733e6a359b0 "github.com/microsoftgraph/msgraph-sdk-go/chats/item/hideforuser"
    i3d748ec071b80bb2a3c14b6b3e8b6b9cfd2a64fe6989dcd8d1bff5dd4b81f605 "github.com/microsoftgraph/msgraph-sdk-go/chats/item/unhideforuser"
    i4829a1a874097739700639a97e70426ca9839240f97e9abd24a5c555929e9d7f "github.com/microsoftgraph/msgraph-sdk-go/chats/item/tabs"
    i6e9937a477483db52c918f11ae44d2375f40192b050c02f9aaa34371ac1a12b8 "github.com/microsoftgraph/msgraph-sdk-go/chats/item/installedapps"
    i6feff4d6e3d0767acb48ac48f31e0afd6402b4d1449e0aa23e2bbb0c98addc1a "github.com/microsoftgraph/msgraph-sdk-go/chats/item/sendactivitynotification"
    i75aa69acd6a8f7bab80f3af14760426f17ce1597f5155737243d02bac7efbc9b "github.com/microsoftgraph/msgraph-sdk-go/chats/item/members"
    iba0aed8e34eab1bae9e9a79ddc2bf2fe62cae60975a70e24a088e3b4472fe1d1 "github.com/microsoftgraph/msgraph-sdk-go/chats/item/lastmessagepreview"
    ibc85c1be0dcca59eb279107d952538363f05577124a509e4d5f9d6a90d514b32 "github.com/microsoftgraph/msgraph-sdk-go/chats/item/markchatunreadforuser"
    ic70252b11f26ab5143efdc891ab5fe95b11c40bd8d46d3f7476b1bd89bbcad4f "github.com/microsoftgraph/msgraph-sdk-go/chats/item/markchatreadforuser"
    ie429cd3a59f9c2a7e91ada198f1b1cd17a6eabca6baa6bc1a6b2904876e96ca5 "github.com/microsoftgraph/msgraph-sdk-go/chats/item/pinnedmessages"
    i837ff1257c18874e8d41607d659eb84e1e4dde73a2ec63b89aa3896598c0d0b1 "github.com/microsoftgraph/msgraph-sdk-go/chats/item/tabs/item"
    ia1765451c54b694370a995b8293695056771707eceb99a55ae5f1a319717a9cd "github.com/microsoftgraph/msgraph-sdk-go/chats/item/members/item"
    ie48eb983af14f89a84a71353f7c41fc8847c3001c54421e947a07c7e5017ae2b "github.com/microsoftgraph/msgraph-sdk-go/chats/item/installedapps/item"
    if01cbefdd30a84f3f8ae04a8237bc738cd6d90ab68b9c2ccb4f8f2e8b18e5f34 "github.com/microsoftgraph/msgraph-sdk-go/chats/item/messages/item"
    if9fa69e93dea7c66b64ee19841f68c83182c43780f205df8e472e307c0d41fc6 "github.com/microsoftgraph/msgraph-sdk-go/chats/item/pinnedmessages/item"
)

// ChatItemRequestBuilder provides operations to manage the collection of chat entities.
type ChatItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ChatItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ChatItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ChatItemRequestBuilderGetQueryParameters retrieve a single chat (without its messages). This method supports federation. To access a chat, at least one chat member must belong to the tenant the request initiated from.
type ChatItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ChatItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ChatItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ChatItemRequestBuilderGetQueryParameters
}
// ChatItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ChatItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewChatItemRequestBuilderInternal instantiates a new ChatItemRequestBuilder and sets the default values.
func NewChatItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ChatItemRequestBuilder) {
    m := &ChatItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/chats/{chat%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewChatItemRequestBuilder instantiates a new ChatItemRequestBuilder and sets the default values.
func NewChatItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ChatItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewChatItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete entity from chats
func (m *ChatItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *ChatItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation retrieve a single chat (without its messages). This method supports federation. To access a chat, at least one chat member must belong to the tenant the request initiated from.
func (m *ChatItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ChatItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the properties of a chat object.
func (m *ChatItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Chatable, requestConfiguration *ChatItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete entity from chats
func (m *ChatItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ChatItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get retrieve a single chat (without its messages). This method supports federation. To access a chat, at least one chat member must belong to the tenant the request initiated from.
func (m *ChatItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ChatItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Chatable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateChatFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Chatable), nil
}
// HideForUser provides operations to call the hideForUser method.
func (m *ChatItemRequestBuilder) HideForUser()(*i39beb613c864af7bec369a1056bf8c41cb817d04318b7282f041a733e6a359b0.HideForUserRequestBuilder) {
    return i39beb613c864af7bec369a1056bf8c41cb817d04318b7282f041a733e6a359b0.NewHideForUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstalledApps provides operations to manage the installedApps property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) InstalledApps()(*i6e9937a477483db52c918f11ae44d2375f40192b050c02f9aaa34371ac1a12b8.InstalledAppsRequestBuilder) {
    return i6e9937a477483db52c918f11ae44d2375f40192b050c02f9aaa34371ac1a12b8.NewInstalledAppsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstalledAppsById provides operations to manage the installedApps property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) InstalledAppsById(id string)(*ie48eb983af14f89a84a71353f7c41fc8847c3001c54421e947a07c7e5017ae2b.TeamsAppInstallationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsAppInstallation%2Did"] = id
    }
    return ie48eb983af14f89a84a71353f7c41fc8847c3001c54421e947a07c7e5017ae2b.NewTeamsAppInstallationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// LastMessagePreview provides operations to manage the lastMessagePreview property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) LastMessagePreview()(*iba0aed8e34eab1bae9e9a79ddc2bf2fe62cae60975a70e24a088e3b4472fe1d1.LastMessagePreviewRequestBuilder) {
    return iba0aed8e34eab1bae9e9a79ddc2bf2fe62cae60975a70e24a088e3b4472fe1d1.NewLastMessagePreviewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MarkChatReadForUser provides operations to call the markChatReadForUser method.
func (m *ChatItemRequestBuilder) MarkChatReadForUser()(*ic70252b11f26ab5143efdc891ab5fe95b11c40bd8d46d3f7476b1bd89bbcad4f.MarkChatReadForUserRequestBuilder) {
    return ic70252b11f26ab5143efdc891ab5fe95b11c40bd8d46d3f7476b1bd89bbcad4f.NewMarkChatReadForUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MarkChatUnreadForUser provides operations to call the markChatUnreadForUser method.
func (m *ChatItemRequestBuilder) MarkChatUnreadForUser()(*ibc85c1be0dcca59eb279107d952538363f05577124a509e4d5f9d6a90d514b32.MarkChatUnreadForUserRequestBuilder) {
    return ibc85c1be0dcca59eb279107d952538363f05577124a509e4d5f9d6a90d514b32.NewMarkChatUnreadForUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Members provides operations to manage the members property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) Members()(*i75aa69acd6a8f7bab80f3af14760426f17ce1597f5155737243d02bac7efbc9b.MembersRequestBuilder) {
    return i75aa69acd6a8f7bab80f3af14760426f17ce1597f5155737243d02bac7efbc9b.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById provides operations to manage the members property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) MembersById(id string)(*ia1765451c54b694370a995b8293695056771707eceb99a55ae5f1a319717a9cd.ConversationMemberItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conversationMember%2Did"] = id
    }
    return ia1765451c54b694370a995b8293695056771707eceb99a55ae5f1a319717a9cd.NewConversationMemberItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Messages provides operations to manage the messages property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) Messages()(*i0453d8f2831f761719a2d96702e4c863ff529132a101dbe684210882f2721d5c.MessagesRequestBuilder) {
    return i0453d8f2831f761719a2d96702e4c863ff529132a101dbe684210882f2721d5c.NewMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MessagesById provides operations to manage the messages property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) MessagesById(id string)(*if01cbefdd30a84f3f8ae04a8237bc738cd6d90ab68b9c2ccb4f8f2e8b18e5f34.ChatMessageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["chatMessage%2Did"] = id
    }
    return if01cbefdd30a84f3f8ae04a8237bc738cd6d90ab68b9c2ccb4f8f2e8b18e5f34.NewChatMessageItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the properties of a chat object.
func (m *ChatItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Chatable, requestConfiguration *ChatItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Chatable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateChatFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Chatable), nil
}
// PinnedMessages provides operations to manage the pinnedMessages property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) PinnedMessages()(*ie429cd3a59f9c2a7e91ada198f1b1cd17a6eabca6baa6bc1a6b2904876e96ca5.PinnedMessagesRequestBuilder) {
    return ie429cd3a59f9c2a7e91ada198f1b1cd17a6eabca6baa6bc1a6b2904876e96ca5.NewPinnedMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PinnedMessagesById provides operations to manage the pinnedMessages property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) PinnedMessagesById(id string)(*if9fa69e93dea7c66b64ee19841f68c83182c43780f205df8e472e307c0d41fc6.PinnedChatMessageInfoItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["pinnedChatMessageInfo%2Did"] = id
    }
    return if9fa69e93dea7c66b64ee19841f68c83182c43780f205df8e472e307c0d41fc6.NewPinnedChatMessageInfoItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SendActivityNotification provides operations to call the sendActivityNotification method.
func (m *ChatItemRequestBuilder) SendActivityNotification()(*i6feff4d6e3d0767acb48ac48f31e0afd6402b4d1449e0aa23e2bbb0c98addc1a.SendActivityNotificationRequestBuilder) {
    return i6feff4d6e3d0767acb48ac48f31e0afd6402b4d1449e0aa23e2bbb0c98addc1a.NewSendActivityNotificationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Tabs provides operations to manage the tabs property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) Tabs()(*i4829a1a874097739700639a97e70426ca9839240f97e9abd24a5c555929e9d7f.TabsRequestBuilder) {
    return i4829a1a874097739700639a97e70426ca9839240f97e9abd24a5c555929e9d7f.NewTabsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TabsById provides operations to manage the tabs property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) TabsById(id string)(*i837ff1257c18874e8d41607d659eb84e1e4dde73a2ec63b89aa3896598c0d0b1.TeamsTabItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsTab%2Did"] = id
    }
    return i837ff1257c18874e8d41607d659eb84e1e4dde73a2ec63b89aa3896598c0d0b1.NewTeamsTabItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// UnhideForUser provides operations to call the unhideForUser method.
func (m *ChatItemRequestBuilder) UnhideForUser()(*i3d748ec071b80bb2a3c14b6b3e8b6b9cfd2a64fe6989dcd8d1bff5dd4b81f605.UnhideForUserRequestBuilder) {
    return i3d748ec071b80bb2a3c14b6b3e8b6b9cfd2a64fe6989dcd8d1bff5dd4b81f605.NewUnhideForUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
