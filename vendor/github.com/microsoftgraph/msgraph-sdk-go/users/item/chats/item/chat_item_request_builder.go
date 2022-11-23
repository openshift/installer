package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i161ad19b0c7a99cae8a80283db44d458d311ea7ea0fda7b7bee6231e5a71bf34 "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/installedapps"
    i29250e65d3047fac7086c3da168f1dc97b716e57040772b108c92eb7e41e5082 "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/members"
    i359872a537fd70c716e7e5ba029dc803c22b94686718a3e10361f77d83422f79 "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/pinnedmessages"
    i4faa2569de9f59350ee6b8738c1c2f9f7d3be53a464b6fcf482f642b99a3ea1d "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/markchatunreadforuser"
    i747f2653d0182726c4364050b88c4c017361f81ae437c821457319c34de8bbe4 "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/messages"
    i7dd8f04ee639104b986c6c5422c345187262c5c456d2627297ccf5846978d15b "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/lastmessagepreview"
    i8bf0adcfacbb5b5e805bc0b691c0068ae5b5c5f80fce37bf131df5325929232d "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/sendactivitynotification"
    ibdf68b59493fcbfe5243f120ca2b97f1b87a4f0eca5bbed6558b780e96d3beec "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/hideforuser"
    ic3b0551739f3e070055bc2d75ba06c264a887bf783cdb8d1d265e7e955d9f288 "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/unhideforuser"
    ic797b970aa8f13def47c8a9daa30e882f6b3ffba8f376b46454ecdca4bbfff2a "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/tabs"
    ied1f558ecd84e0ae9eb5ca8fc926f3fa1232c1f5894efa8cae1410b1cd659c38 "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/markchatreadforuser"
    i03a5d9d6f43c34754e384f5d06895085102259e8c0fd8d72b735c4f7e27e19e3 "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/members/item"
    i77281ca541a7a54c245fdd2770ea294bd7c00f30c47f9c268ae075c47417ee91 "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/pinnedmessages/item"
    i89204df365cb0b4df35e680ee9e843167bce6597538b2e059b21aa0890ed91ae "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/tabs/item"
    ia5fe7898b11b622f06cabbbf98b672db477f4cd5aa1bdc1f1091441663a7ac94 "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/messages/item"
    ice5d918650d80a2fd00a96c09c41f928c4c75fc531425a6ae4f4cc396b70c4f2 "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item/installedapps/item"
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
// ChatItemRequestBuilderGetQueryParameters get chats from users
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
    m.urlTemplate = "{+baseurl}/users/{user%2Did}/chats/{chat%2Did}{?%24select,%24expand}";
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
// CreateDeleteRequestInformation delete navigation property chats for users
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
// CreateGetRequestInformation get chats from users
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
// CreatePatchRequestInformation update the navigation property chats in users
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
// Delete delete navigation property chats for users
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
// Get get chats from users
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
func (m *ChatItemRequestBuilder) HideForUser()(*ibdf68b59493fcbfe5243f120ca2b97f1b87a4f0eca5bbed6558b780e96d3beec.HideForUserRequestBuilder) {
    return ibdf68b59493fcbfe5243f120ca2b97f1b87a4f0eca5bbed6558b780e96d3beec.NewHideForUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstalledApps provides operations to manage the installedApps property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) InstalledApps()(*i161ad19b0c7a99cae8a80283db44d458d311ea7ea0fda7b7bee6231e5a71bf34.InstalledAppsRequestBuilder) {
    return i161ad19b0c7a99cae8a80283db44d458d311ea7ea0fda7b7bee6231e5a71bf34.NewInstalledAppsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstalledAppsById provides operations to manage the installedApps property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) InstalledAppsById(id string)(*ice5d918650d80a2fd00a96c09c41f928c4c75fc531425a6ae4f4cc396b70c4f2.TeamsAppInstallationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsAppInstallation%2Did"] = id
    }
    return ice5d918650d80a2fd00a96c09c41f928c4c75fc531425a6ae4f4cc396b70c4f2.NewTeamsAppInstallationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// LastMessagePreview provides operations to manage the lastMessagePreview property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) LastMessagePreview()(*i7dd8f04ee639104b986c6c5422c345187262c5c456d2627297ccf5846978d15b.LastMessagePreviewRequestBuilder) {
    return i7dd8f04ee639104b986c6c5422c345187262c5c456d2627297ccf5846978d15b.NewLastMessagePreviewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MarkChatReadForUser provides operations to call the markChatReadForUser method.
func (m *ChatItemRequestBuilder) MarkChatReadForUser()(*ied1f558ecd84e0ae9eb5ca8fc926f3fa1232c1f5894efa8cae1410b1cd659c38.MarkChatReadForUserRequestBuilder) {
    return ied1f558ecd84e0ae9eb5ca8fc926f3fa1232c1f5894efa8cae1410b1cd659c38.NewMarkChatReadForUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MarkChatUnreadForUser provides operations to call the markChatUnreadForUser method.
func (m *ChatItemRequestBuilder) MarkChatUnreadForUser()(*i4faa2569de9f59350ee6b8738c1c2f9f7d3be53a464b6fcf482f642b99a3ea1d.MarkChatUnreadForUserRequestBuilder) {
    return i4faa2569de9f59350ee6b8738c1c2f9f7d3be53a464b6fcf482f642b99a3ea1d.NewMarkChatUnreadForUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Members provides operations to manage the members property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) Members()(*i29250e65d3047fac7086c3da168f1dc97b716e57040772b108c92eb7e41e5082.MembersRequestBuilder) {
    return i29250e65d3047fac7086c3da168f1dc97b716e57040772b108c92eb7e41e5082.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById provides operations to manage the members property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) MembersById(id string)(*i03a5d9d6f43c34754e384f5d06895085102259e8c0fd8d72b735c4f7e27e19e3.ConversationMemberItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conversationMember%2Did"] = id
    }
    return i03a5d9d6f43c34754e384f5d06895085102259e8c0fd8d72b735c4f7e27e19e3.NewConversationMemberItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Messages provides operations to manage the messages property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) Messages()(*i747f2653d0182726c4364050b88c4c017361f81ae437c821457319c34de8bbe4.MessagesRequestBuilder) {
    return i747f2653d0182726c4364050b88c4c017361f81ae437c821457319c34de8bbe4.NewMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MessagesById provides operations to manage the messages property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) MessagesById(id string)(*ia5fe7898b11b622f06cabbbf98b672db477f4cd5aa1bdc1f1091441663a7ac94.ChatMessageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["chatMessage%2Did"] = id
    }
    return ia5fe7898b11b622f06cabbbf98b672db477f4cd5aa1bdc1f1091441663a7ac94.NewChatMessageItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the navigation property chats in users
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
func (m *ChatItemRequestBuilder) PinnedMessages()(*i359872a537fd70c716e7e5ba029dc803c22b94686718a3e10361f77d83422f79.PinnedMessagesRequestBuilder) {
    return i359872a537fd70c716e7e5ba029dc803c22b94686718a3e10361f77d83422f79.NewPinnedMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PinnedMessagesById provides operations to manage the pinnedMessages property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) PinnedMessagesById(id string)(*i77281ca541a7a54c245fdd2770ea294bd7c00f30c47f9c268ae075c47417ee91.PinnedChatMessageInfoItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["pinnedChatMessageInfo%2Did"] = id
    }
    return i77281ca541a7a54c245fdd2770ea294bd7c00f30c47f9c268ae075c47417ee91.NewPinnedChatMessageInfoItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SendActivityNotification provides operations to call the sendActivityNotification method.
func (m *ChatItemRequestBuilder) SendActivityNotification()(*i8bf0adcfacbb5b5e805bc0b691c0068ae5b5c5f80fce37bf131df5325929232d.SendActivityNotificationRequestBuilder) {
    return i8bf0adcfacbb5b5e805bc0b691c0068ae5b5c5f80fce37bf131df5325929232d.NewSendActivityNotificationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Tabs provides operations to manage the tabs property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) Tabs()(*ic797b970aa8f13def47c8a9daa30e882f6b3ffba8f376b46454ecdca4bbfff2a.TabsRequestBuilder) {
    return ic797b970aa8f13def47c8a9daa30e882f6b3ffba8f376b46454ecdca4bbfff2a.NewTabsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TabsById provides operations to manage the tabs property of the microsoft.graph.chat entity.
func (m *ChatItemRequestBuilder) TabsById(id string)(*i89204df365cb0b4df35e680ee9e843167bce6597538b2e059b21aa0890ed91ae.TeamsTabItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsTab%2Did"] = id
    }
    return i89204df365cb0b4df35e680ee9e843167bce6597538b2e059b21aa0890ed91ae.NewTeamsTabItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// UnhideForUser provides operations to call the unhideForUser method.
func (m *ChatItemRequestBuilder) UnhideForUser()(*ic3b0551739f3e070055bc2d75ba06c264a887bf783cdb8d1d265e7e955d9f288.UnhideForUserRequestBuilder) {
    return ic3b0551739f3e070055bc2d75ba06c264a887bf783cdb8d1d265e7e955d9f288.NewUnhideForUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
