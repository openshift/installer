package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i0d53f661bb96ea28364f561f24c6a15b30ee8383c8b828dac50eca27ee85d3fa "github.com/microsoftgraph/msgraph-sdk-go/teams/item/operations"
    i112483150813fdda19082d152401ea7792fc69e950a26fd2de0f70593b4e9aa3 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/schedule"
    i149d801c662f9e963517ba3a6c49cb4845b36d355c160bcc1729479cd18d6107 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/photo"
    i15af31366dc887244693eaa3d3a5b29b322f057d84108b094ab85c131706bc4d "github.com/microsoftgraph/msgraph-sdk-go/teams/item/sendactivitynotification"
    i3ab10cbdf172c82da95d8c321e7ab0c74610f25b0a5df76707ecc483fe73c949 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/unarchive"
    i6b2db9d1ed1ca06816d98e7327f880f451780be7beb8380462b23278f4c949d8 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/archive"
    i86a4744d5e45722adbed1e0f22b69f1a4b6148d0445ff82385abf3cedd90a29d "github.com/microsoftgraph/msgraph-sdk-go/teams/item/clone"
    i94154a059496bbfe5f4225f44a08db234c17a10c3a7d58124a8fb182751e3bc6 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/installedapps"
    i941dc9f0521a47c6bd3be5dd3a072e94f5381994ffe128d7e54810e76ed8e473 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/template"
    i96200494e669e195393f509bc7a1b076f717becc8f2d0e85d3bdc73e820e2615 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/primarychannel"
    i9c16b9ccd15275976651d11bc1695c103bd12afe48b01714765458594b9f7975 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/channels"
    ic0623992ab864953cddbe4e195f8ecf82c663de82e6f4af4bd9402875dbb3257 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/group"
    ic1586e13f55531e36f79cfc95650dc1050b9aa6f0c57969599b2fe450a05078c "github.com/microsoftgraph/msgraph-sdk-go/teams/item/incomingchannels"
    id10e40d460f5e6fec4c37430edd37ea1b19ebe0818c68acfde864d94b6e465f1 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/tags"
    id19ffdf54b88ab907e8aa643b2969b583758c77cb61e7e304e1e7fbc61c9d1fa "github.com/microsoftgraph/msgraph-sdk-go/teams/item/completemigration"
    id42b130ad2d47c0a14f279b59cb0081ea2cd07584c7a38df737c9dd7c99c1976 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/allchannels"
    ie3ea4358486c3530306ce0ee6275a0285de71b92eff7f63adb5e51ed2ef63235 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/members"
    i167530cb07fd4387f5ffe0c5a728e41dd5608a4963302dd28d98720ea3c093dc "github.com/microsoftgraph/msgraph-sdk-go/teams/item/members/item"
    i3efb7759b468a039683d7e9a7233c6ef567ad6e2bfdb31d83c6c7c4bfb121545 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/operations/item"
    i47ac7f5f25f369795017e987d58959bf433cffedfeee0f41d10f557d0dc64f40 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/incomingchannels/item"
    i6b8094a2be2251280648d1c8aa1f2cb616745b9111719fe2f3f1b754098a2ea7 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/installedapps/item"
    i95bd316e1af2146251d5c5f8f94d9865c17745b134282c19c4fae4c81ba47474 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/channels/item"
    i9fd988942b1914b11d1790841f3af7b85efd256822487ed83679466a66b3fb59 "github.com/microsoftgraph/msgraph-sdk-go/teams/item/allchannels/item"
    id69fe1bc24b951cc84e66e8108d1952ade3a29cca8871917c9be1e3dd0b9449e "github.com/microsoftgraph/msgraph-sdk-go/teams/item/tags/item"
)

// TeamItemRequestBuilder provides operations to manage the collection of team entities.
type TeamItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// TeamItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TeamItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// TeamItemRequestBuilderGetQueryParameters retrieve the properties and relationships of the specified team.
type TeamItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// TeamItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TeamItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *TeamItemRequestBuilderGetQueryParameters
}
// TeamItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TeamItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AllChannels provides operations to manage the allChannels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) AllChannels()(*id42b130ad2d47c0a14f279b59cb0081ea2cd07584c7a38df737c9dd7c99c1976.AllChannelsRequestBuilder) {
    return id42b130ad2d47c0a14f279b59cb0081ea2cd07584c7a38df737c9dd7c99c1976.NewAllChannelsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AllChannelsById provides operations to manage the allChannels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) AllChannelsById(id string)(*i9fd988942b1914b11d1790841f3af7b85efd256822487ed83679466a66b3fb59.ChannelItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["channel%2Did"] = id
    }
    return i9fd988942b1914b11d1790841f3af7b85efd256822487ed83679466a66b3fb59.NewChannelItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Archive provides operations to call the archive method.
func (m *TeamItemRequestBuilder) Archive()(*i6b2db9d1ed1ca06816d98e7327f880f451780be7beb8380462b23278f4c949d8.ArchiveRequestBuilder) {
    return i6b2db9d1ed1ca06816d98e7327f880f451780be7beb8380462b23278f4c949d8.NewArchiveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Channels provides operations to manage the channels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Channels()(*i9c16b9ccd15275976651d11bc1695c103bd12afe48b01714765458594b9f7975.ChannelsRequestBuilder) {
    return i9c16b9ccd15275976651d11bc1695c103bd12afe48b01714765458594b9f7975.NewChannelsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChannelsById provides operations to manage the channels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) ChannelsById(id string)(*i95bd316e1af2146251d5c5f8f94d9865c17745b134282c19c4fae4c81ba47474.ChannelItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["channel%2Did"] = id
    }
    return i95bd316e1af2146251d5c5f8f94d9865c17745b134282c19c4fae4c81ba47474.NewChannelItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Clone provides operations to call the clone method.
func (m *TeamItemRequestBuilder) Clone()(*i86a4744d5e45722adbed1e0f22b69f1a4b6148d0445ff82385abf3cedd90a29d.CloneRequestBuilder) {
    return i86a4744d5e45722adbed1e0f22b69f1a4b6148d0445ff82385abf3cedd90a29d.NewCloneRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CompleteMigration provides operations to call the completeMigration method.
func (m *TeamItemRequestBuilder) CompleteMigration()(*id19ffdf54b88ab907e8aa643b2969b583758c77cb61e7e304e1e7fbc61c9d1fa.CompleteMigrationRequestBuilder) {
    return id19ffdf54b88ab907e8aa643b2969b583758c77cb61e7e304e1e7fbc61c9d1fa.NewCompleteMigrationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewTeamItemRequestBuilderInternal instantiates a new TeamItemRequestBuilder and sets the default values.
func NewTeamItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TeamItemRequestBuilder) {
    m := &TeamItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/teams/{team%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewTeamItemRequestBuilder instantiates a new TeamItemRequestBuilder and sets the default values.
func NewTeamItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TeamItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewTeamItemRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete entity from teams
func (m *TeamItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *TeamItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation retrieve the properties and relationships of the specified team.
func (m *TeamItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *TeamItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update the properties of the specified team.
func (m *TeamItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, requestConfiguration *TeamItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete entity from teams
func (m *TeamItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *TeamItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get retrieve the properties and relationships of the specified team.
func (m *TeamItemRequestBuilder) Get(ctx context.Context, requestConfiguration *TeamItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTeamFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable), nil
}
// Group provides operations to manage the group property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Group()(*ic0623992ab864953cddbe4e195f8ecf82c663de82e6f4af4bd9402875dbb3257.GroupRequestBuilder) {
    return ic0623992ab864953cddbe4e195f8ecf82c663de82e6f4af4bd9402875dbb3257.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// IncomingChannels provides operations to manage the incomingChannels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) IncomingChannels()(*ic1586e13f55531e36f79cfc95650dc1050b9aa6f0c57969599b2fe450a05078c.IncomingChannelsRequestBuilder) {
    return ic1586e13f55531e36f79cfc95650dc1050b9aa6f0c57969599b2fe450a05078c.NewIncomingChannelsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// IncomingChannelsById provides operations to manage the incomingChannels property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) IncomingChannelsById(id string)(*i47ac7f5f25f369795017e987d58959bf433cffedfeee0f41d10f557d0dc64f40.ChannelItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["channel%2Did"] = id
    }
    return i47ac7f5f25f369795017e987d58959bf433cffedfeee0f41d10f557d0dc64f40.NewChannelItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// InstalledApps provides operations to manage the installedApps property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) InstalledApps()(*i94154a059496bbfe5f4225f44a08db234c17a10c3a7d58124a8fb182751e3bc6.InstalledAppsRequestBuilder) {
    return i94154a059496bbfe5f4225f44a08db234c17a10c3a7d58124a8fb182751e3bc6.NewInstalledAppsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstalledAppsById provides operations to manage the installedApps property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) InstalledAppsById(id string)(*i6b8094a2be2251280648d1c8aa1f2cb616745b9111719fe2f3f1b754098a2ea7.TeamsAppInstallationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsAppInstallation%2Did"] = id
    }
    return i6b8094a2be2251280648d1c8aa1f2cb616745b9111719fe2f3f1b754098a2ea7.NewTeamsAppInstallationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Members provides operations to manage the members property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Members()(*ie3ea4358486c3530306ce0ee6275a0285de71b92eff7f63adb5e51ed2ef63235.MembersRequestBuilder) {
    return ie3ea4358486c3530306ce0ee6275a0285de71b92eff7f63adb5e51ed2ef63235.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById provides operations to manage the members property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) MembersById(id string)(*i167530cb07fd4387f5ffe0c5a728e41dd5608a4963302dd28d98720ea3c093dc.ConversationMemberItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conversationMember%2Did"] = id
    }
    return i167530cb07fd4387f5ffe0c5a728e41dd5608a4963302dd28d98720ea3c093dc.NewConversationMemberItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Operations provides operations to manage the operations property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Operations()(*i0d53f661bb96ea28364f561f24c6a15b30ee8383c8b828dac50eca27ee85d3fa.OperationsRequestBuilder) {
    return i0d53f661bb96ea28364f561f24c6a15b30ee8383c8b828dac50eca27ee85d3fa.NewOperationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OperationsById provides operations to manage the operations property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) OperationsById(id string)(*i3efb7759b468a039683d7e9a7233c6ef567ad6e2bfdb31d83c6c7c4bfb121545.TeamsAsyncOperationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsAsyncOperation%2Did"] = id
    }
    return i3efb7759b468a039683d7e9a7233c6ef567ad6e2bfdb31d83c6c7c4bfb121545.NewTeamsAsyncOperationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the properties of the specified team.
func (m *TeamItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, requestConfiguration *TeamItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTeamFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable), nil
}
// Photo provides operations to manage the photo property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Photo()(*i149d801c662f9e963517ba3a6c49cb4845b36d355c160bcc1729479cd18d6107.PhotoRequestBuilder) {
    return i149d801c662f9e963517ba3a6c49cb4845b36d355c160bcc1729479cd18d6107.NewPhotoRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PrimaryChannel provides operations to manage the primaryChannel property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) PrimaryChannel()(*i96200494e669e195393f509bc7a1b076f717becc8f2d0e85d3bdc73e820e2615.PrimaryChannelRequestBuilder) {
    return i96200494e669e195393f509bc7a1b076f717becc8f2d0e85d3bdc73e820e2615.NewPrimaryChannelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Schedule provides operations to manage the schedule property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Schedule()(*i112483150813fdda19082d152401ea7792fc69e950a26fd2de0f70593b4e9aa3.ScheduleRequestBuilder) {
    return i112483150813fdda19082d152401ea7792fc69e950a26fd2de0f70593b4e9aa3.NewScheduleRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SendActivityNotification provides operations to call the sendActivityNotification method.
func (m *TeamItemRequestBuilder) SendActivityNotification()(*i15af31366dc887244693eaa3d3a5b29b322f057d84108b094ab85c131706bc4d.SendActivityNotificationRequestBuilder) {
    return i15af31366dc887244693eaa3d3a5b29b322f057d84108b094ab85c131706bc4d.NewSendActivityNotificationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Tags provides operations to manage the tags property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Tags()(*id10e40d460f5e6fec4c37430edd37ea1b19ebe0818c68acfde864d94b6e465f1.TagsRequestBuilder) {
    return id10e40d460f5e6fec4c37430edd37ea1b19ebe0818c68acfde864d94b6e465f1.NewTagsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TagsById provides operations to manage the tags property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) TagsById(id string)(*id69fe1bc24b951cc84e66e8108d1952ade3a29cca8871917c9be1e3dd0b9449e.TeamworkTagItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamworkTag%2Did"] = id
    }
    return id69fe1bc24b951cc84e66e8108d1952ade3a29cca8871917c9be1e3dd0b9449e.NewTeamworkTagItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Template provides operations to manage the template property of the microsoft.graph.team entity.
func (m *TeamItemRequestBuilder) Template()(*i941dc9f0521a47c6bd3be5dd3a072e94f5381994ffe128d7e54810e76ed8e473.TemplateRequestBuilder) {
    return i941dc9f0521a47c6bd3be5dd3a072e94f5381994ffe128d7e54810e76ed8e473.NewTemplateRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Unarchive provides operations to call the unarchive method.
func (m *TeamItemRequestBuilder) Unarchive()(*i3ab10cbdf172c82da95d8c321e7ab0c74610f25b0a5df76707ecc483fe73c949.UnarchiveRequestBuilder) {
    return i3ab10cbdf172c82da95d8c321e7ab0c74610f25b0a5df76707ecc483fe73c949.NewUnarchiveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
