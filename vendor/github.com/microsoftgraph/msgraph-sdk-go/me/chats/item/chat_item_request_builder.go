package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i05eb5a4625242f589812d7abaff1378dbfaeacd7daaf8058c776a2f6fe5078a2 "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/unhideforuser"
    i2ea8718568dd68e4db5b0c9b13edefe941431dfb1541af1606f4bc0a7d2db21c "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/installedapps"
    i4f7324bbc51d39044294fbc8507574ec35a78d7e6d4521fbb14db70b57abe700 "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/markchatreadforuser"
    i58e23b2939c6432ac8b587c5feeb471fc47d5b55f39473a636fdaaa34d146339 "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/hideforuser"
    i5f59d5a81c4fcc79151c80c703b7058649721f3ebf324ad807f90668dbb2100b "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/pinnedmessages"
    i88a4f2f2a7c4ae16c306dfc70fe5cfa8f9175e29138583b1d60f25decd49e6d9 "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/messages"
    i97c7b1a52b2a1f45a86bbc3b07ceed8a2a6e243e966278d5e86afe0fad5d5beb "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/markchatunreadforuser"
    iaf22e2321aed77ea75fad259764b34bc4b11d87b943ce0c01583b0cf198686bb "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/lastmessagepreview"
    ibfa184694112defc7dad34dc61131fc41730c513148b28b61ceae4b917fc813b "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/members"
    ic48531e217b5a3e24c10e42f0363211f421bdab3ff973fab2c2ccbbb3e67f8c3 "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/sendactivitynotification"
    id387b678bda7815af798dc5b19d04d6775c9e30340b77735c65b245e5bae6867 "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/tabs"
    i25a5bc17ba5983e7d43d8a128bba975fd3706b2f2204b481e7864fe0ae9d5030 "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/tabs/item"
    i2febef273a5d29b65dd84845ed6389ff6349f765563d2c8edcf94919833a9f49 "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/members/item"
    i8d1551e7ddb5563740c721b9a2267c157137532374fa22029cbe1c43d7821a74 "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/installedapps/item"
    i9c260c1538d5b918af33a76a10a514936b0f1f71f65fedc06df9fbe2d4dc5876 "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/messages/item"
    ia85488b3d41cd3644b41191b3d939e66469b14d2c6c8c8399530a704bc606aaf "github.com/microsoftgraph/msgraph-sdk-go/me/chats/item/pinnedmessages/item"
)

// ChatItemRequestBuilder provides operations to manage the chats property of the microsoft.graph.user entity.
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
// ChatItemRequestBuilderGetQueryParameters get chats from me
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
    m.urlTemplate = "{+baseurl}/me/chats/{chat%2Did}{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property chats for me
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
// CreateGetRequestInformation get chats from me
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
// CreatePatchRequestInformation update the navigation property chats in me
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
// Delete delete navigation property chats for me
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
// Get get chats from me
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
func (m *ChatItemRequestBuilder) HideForUser()(*i58e23b2939c6432ac8b587c5feeb471fc47d5b55f39473a636fdaaa34d146339.HideForUserRequestBuilder) {
    return i58e23b2939c6432ac8b587c5feeb471fc47d5b55f39473a636fdaaa34d146339.NewHideForUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstalledApps provides operations to manage the installedApps property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) InstalledApps()(*i2ea8718568dd68e4db5b0c9b13edefe941431dfb1541af1606f4bc0a7d2db21c.InstalledAppsRequestBuilder) {
    return i2ea8718568dd68e4db5b0c9b13edefe941431dfb1541af1606f4bc0a7d2db21c.NewInstalledAppsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstalledAppsById provides operations to manage the installedApps property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) InstalledAppsById(id string)(*i8d1551e7ddb5563740c721b9a2267c157137532374fa22029cbe1c43d7821a74.TeamsAppInstallationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsAppInstallation%2Did"] = id
    }
    return i8d1551e7ddb5563740c721b9a2267c157137532374fa22029cbe1c43d7821a74.NewTeamsAppInstallationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// LastMessagePreview provides operations to manage the lastMessagePreview property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) LastMessagePreview()(*iaf22e2321aed77ea75fad259764b34bc4b11d87b943ce0c01583b0cf198686bb.LastMessagePreviewRequestBuilder) {
    return iaf22e2321aed77ea75fad259764b34bc4b11d87b943ce0c01583b0cf198686bb.NewLastMessagePreviewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MarkChatReadForUser provides operations to call the markChatReadForUser method.
func (m *ChatItemRequestBuilder) MarkChatReadForUser()(*i4f7324bbc51d39044294fbc8507574ec35a78d7e6d4521fbb14db70b57abe700.MarkChatReadForUserRequestBuilder) {
    return i4f7324bbc51d39044294fbc8507574ec35a78d7e6d4521fbb14db70b57abe700.NewMarkChatReadForUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MarkChatUnreadForUser provides operations to call the markChatUnreadForUser method.
func (m *ChatItemRequestBuilder) MarkChatUnreadForUser()(*i97c7b1a52b2a1f45a86bbc3b07ceed8a2a6e243e966278d5e86afe0fad5d5beb.MarkChatUnreadForUserRequestBuilder) {
    return i97c7b1a52b2a1f45a86bbc3b07ceed8a2a6e243e966278d5e86afe0fad5d5beb.NewMarkChatUnreadForUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Members provides operations to manage the members property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) Members()(*ibfa184694112defc7dad34dc61131fc41730c513148b28b61ceae4b917fc813b.MembersRequestBuilder) {
    return ibfa184694112defc7dad34dc61131fc41730c513148b28b61ceae4b917fc813b.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById provides operations to manage the members property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) MembersById(id string)(*i2febef273a5d29b65dd84845ed6389ff6349f765563d2c8edcf94919833a9f49.ConversationMemberItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conversationMember%2Did"] = id
    }
    return i2febef273a5d29b65dd84845ed6389ff6349f765563d2c8edcf94919833a9f49.NewConversationMemberItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Messages provides operations to manage the messages property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) Messages()(*i88a4f2f2a7c4ae16c306dfc70fe5cfa8f9175e29138583b1d60f25decd49e6d9.MessagesRequestBuilder) {
    return i88a4f2f2a7c4ae16c306dfc70fe5cfa8f9175e29138583b1d60f25decd49e6d9.NewMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MessagesById provides operations to manage the messages property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) MessagesById(id string)(*i9c260c1538d5b918af33a76a10a514936b0f1f71f65fedc06df9fbe2d4dc5876.ChatMessageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["chatMessage%2Did"] = id
    }
    return i9c260c1538d5b918af33a76a10a514936b0f1f71f65fedc06df9fbe2d4dc5876.NewChatMessageItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property chats in me
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
func (m *ChatItemRequestBuilder) PinnedMessages()(*i5f59d5a81c4fcc79151c80c703b7058649721f3ebf324ad807f90668dbb2100b.PinnedMessagesRequestBuilder) {
    return i5f59d5a81c4fcc79151c80c703b7058649721f3ebf324ad807f90668dbb2100b.NewPinnedMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PinnedMessagesById provides operations to manage the pinnedMessages property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) PinnedMessagesById(id string)(*ia85488b3d41cd3644b41191b3d939e66469b14d2c6c8c8399530a704bc606aaf.PinnedChatMessageInfoItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["pinnedChatMessageInfo%2Did"] = id
    }
    return ia85488b3d41cd3644b41191b3d939e66469b14d2c6c8c8399530a704bc606aaf.NewPinnedChatMessageInfoItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SendActivityNotification provides operations to call the sendActivityNotification method.
func (m *ChatItemRequestBuilder) SendActivityNotification()(*ic48531e217b5a3e24c10e42f0363211f421bdab3ff973fab2c2ccbbb3e67f8c3.SendActivityNotificationRequestBuilder) {
    return ic48531e217b5a3e24c10e42f0363211f421bdab3ff973fab2c2ccbbb3e67f8c3.NewSendActivityNotificationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Tabs provides operations to manage the tabs property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) Tabs()(*id387b678bda7815af798dc5b19d04d6775c9e30340b77735c65b245e5bae6867.TabsRequestBuilder) {
    return id387b678bda7815af798dc5b19d04d6775c9e30340b77735c65b245e5bae6867.NewTabsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TabsById provides operations to manage the tabs property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) TabsById(id string)(*i25a5bc17ba5983e7d43d8a128bba975fd3706b2f2204b481e7864fe0ae9d5030.TeamsTabItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsTab%2Did"] = id
    }
    return i25a5bc17ba5983e7d43d8a128bba975fd3706b2f2204b481e7864fe0ae9d5030.NewTeamsTabItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// UnhideForUser provides operations to call the unhideForUser method.
func (m *ChatItemRequestBuilder) UnhideForUser()(*i05eb5a4625242f589812d7abaff1378dbfaeacd7daaf8058c776a2f6fe5078a2.UnhideForUserRequestBuilder) {
    return i05eb5a4625242f589812d7abaff1378dbfaeacd7daaf8058c776a2f6fe5078a2.NewUnhideForUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
