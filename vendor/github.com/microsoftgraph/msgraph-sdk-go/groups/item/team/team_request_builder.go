package team

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i04e49204322d308913c8af78daf5ae0e14fb7df81490288b98bd1c408480bbb4 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/group"
    i19cce4ed716ea1e1791100abc23b02546b679f6617320780eac07feae2047b70 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/allchannels"
    i34294a7456473fba252b3ca961ed3381f55b7a770d74ee4720d752d9a45ad739 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/primarychannel"
    i36b62d6f777744a15bac3ff84127f195718092c76813b8607f20150831a58b5a "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/incomingchannels"
    i40ee01a629bf206681620d9e9951dc387487230caf141be819f95c589a8f0779 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/sendactivitynotification"
    i477873d925c8fe5f007dfb4043f82ed6374d80069f96d4931626a23aeb0f2a87 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/template"
    i495e85248bfc2c7bf36f36ede17adf8bc672f00e69df645f352a27b2b7d30e61 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/schedule"
    i4d1bfd5f8aba268d7200e167f158f47cd0fbcd973717f0a5d357c3ba868d2609 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/clone"
    i7aee2305413b3b465064d02a6aa79c0f32df27bc874d76350b572847a880a22f "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/installedapps"
    i8f663f97c6eef0cee280b998b250f55e4d63fc289392b2a09af9c95210ca5d5d "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels"
    i9086fbb99e21665fd895070750ee43da3017e1397eaa8655faeb9701ad70d234 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/completemigration"
    i946cc8c196bfeaf31caf19315aac4a6c73de740196dc6305506cfd9d4f440ddb "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/photo"
    i97dd0d985d4eaf1cbaeb3d771ef48c101ec691107d9a79dcd3f38213ccae863b "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/members"
    i9f6a3c1c4300552f4792a53f0d0fe603636d1b865d408ff28e1cdd722683ebe3 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/operations"
    iac2547e415c1013b90e41a271000fafc1a87da2e4529b1d8b5ba6f0aa631ab8c "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/archive"
    idf8489132b9d5484e716c4faff83b41054a8bc066341cf00a201254e31131154 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/tags"
    ie49085238af0f76cc56d306db6cbd7ea879a083402f3fffae10acf7f0e81362d "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/unarchive"
    i6e500f69af937bd22bbfc2b6ee2698da77bfc652c449bb8ce6a2ee75470b9de9 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/installedapps/item"
    ia28d226e35a69ac659de8c513208f3f357b7559e89b7a71709ae6d2bc8229124 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/incomingchannels/item"
    id5fe3547558f37b7931bf00eb734fbe581be897e9aea405e91e3ce2dbeec74b3 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/operations/item"
    idbd66e2d729b9850c37a45204c98fab6700e673bdde41eb4e751af024dc0a897 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/tags/item"
    ie59099b1e42dc6f55628ccb467e95d4e7416a9bd33399bd6beff110ff67ea34e "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/members/item"
    ie594d6624cc5e7f45c840beec42560c6f347a2da3e424188e2fbe60f77fbc850 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/channels/item"
    ifea2e46130f0932465c5758491c2645c019fd218aa02d4f0dc6adb86351531d3 "github.com/microsoftgraph/msgraph-sdk-go/groups/item/team/allchannels/item"
)

// TeamRequestBuilder provides operations to manage the team property of the microsoft.graph.group entity.
type TeamRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// TeamRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TeamRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// TeamRequestBuilderGetQueryParameters the team associated with this group.
type TeamRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// TeamRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TeamRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *TeamRequestBuilderGetQueryParameters
}
// TeamRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TeamRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AllChannels provides operations to manage the allChannels property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) AllChannels()(*i19cce4ed716ea1e1791100abc23b02546b679f6617320780eac07feae2047b70.AllChannelsRequestBuilder) {
    return i19cce4ed716ea1e1791100abc23b02546b679f6617320780eac07feae2047b70.NewAllChannelsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AllChannelsById provides operations to manage the allChannels property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) AllChannelsById(id string)(*ifea2e46130f0932465c5758491c2645c019fd218aa02d4f0dc6adb86351531d3.ChannelItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["channel%2Did"] = id
    }
    return ifea2e46130f0932465c5758491c2645c019fd218aa02d4f0dc6adb86351531d3.NewChannelItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Archive provides operations to call the archive method.
func (m *TeamRequestBuilder) Archive()(*iac2547e415c1013b90e41a271000fafc1a87da2e4529b1d8b5ba6f0aa631ab8c.ArchiveRequestBuilder) {
    return iac2547e415c1013b90e41a271000fafc1a87da2e4529b1d8b5ba6f0aa631ab8c.NewArchiveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Channels provides operations to manage the channels property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) Channels()(*i8f663f97c6eef0cee280b998b250f55e4d63fc289392b2a09af9c95210ca5d5d.ChannelsRequestBuilder) {
    return i8f663f97c6eef0cee280b998b250f55e4d63fc289392b2a09af9c95210ca5d5d.NewChannelsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChannelsById provides operations to manage the channels property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) ChannelsById(id string)(*ie594d6624cc5e7f45c840beec42560c6f347a2da3e424188e2fbe60f77fbc850.ChannelItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["channel%2Did"] = id
    }
    return ie594d6624cc5e7f45c840beec42560c6f347a2da3e424188e2fbe60f77fbc850.NewChannelItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Clone provides operations to call the clone method.
func (m *TeamRequestBuilder) Clone()(*i4d1bfd5f8aba268d7200e167f158f47cd0fbcd973717f0a5d357c3ba868d2609.CloneRequestBuilder) {
    return i4d1bfd5f8aba268d7200e167f158f47cd0fbcd973717f0a5d357c3ba868d2609.NewCloneRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CompleteMigration provides operations to call the completeMigration method.
func (m *TeamRequestBuilder) CompleteMigration()(*i9086fbb99e21665fd895070750ee43da3017e1397eaa8655faeb9701ad70d234.CompleteMigrationRequestBuilder) {
    return i9086fbb99e21665fd895070750ee43da3017e1397eaa8655faeb9701ad70d234.NewCompleteMigrationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewTeamRequestBuilderInternal instantiates a new TeamRequestBuilder and sets the default values.
func NewTeamRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TeamRequestBuilder) {
    m := &TeamRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/team{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewTeamRequestBuilder instantiates a new TeamRequestBuilder and sets the default values.
func NewTeamRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TeamRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewTeamRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateDeleteRequestInformation delete navigation property team for groups
func (m *TeamRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *TeamRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreateGetRequestInformation the team associated with this group.
func (m *TeamRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *TeamRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation create a new team under a group. In order to create a team, the group must have a least one owner. If the group was created less than 15 minutes ago, it's possible for the Create team call to fail with a 404 error code due to replication delays. The recommended pattern is to retry the Create team call three times, with a 10 second delay between calls.
func (m *TeamRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, requestConfiguration *TeamRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Delete delete navigation property team for groups
func (m *TeamRequestBuilder) Delete(ctx context.Context, requestConfiguration *TeamRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get the team associated with this group.
func (m *TeamRequestBuilder) Get(ctx context.Context, requestConfiguration *TeamRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, error) {
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
func (m *TeamRequestBuilder) Group()(*i04e49204322d308913c8af78daf5ae0e14fb7df81490288b98bd1c408480bbb4.GroupRequestBuilder) {
    return i04e49204322d308913c8af78daf5ae0e14fb7df81490288b98bd1c408480bbb4.NewGroupRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// IncomingChannels provides operations to manage the incomingChannels property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) IncomingChannels()(*i36b62d6f777744a15bac3ff84127f195718092c76813b8607f20150831a58b5a.IncomingChannelsRequestBuilder) {
    return i36b62d6f777744a15bac3ff84127f195718092c76813b8607f20150831a58b5a.NewIncomingChannelsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// IncomingChannelsById provides operations to manage the incomingChannels property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) IncomingChannelsById(id string)(*ia28d226e35a69ac659de8c513208f3f357b7559e89b7a71709ae6d2bc8229124.ChannelItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["channel%2Did"] = id
    }
    return ia28d226e35a69ac659de8c513208f3f357b7559e89b7a71709ae6d2bc8229124.NewChannelItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// InstalledApps provides operations to manage the installedApps property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) InstalledApps()(*i7aee2305413b3b465064d02a6aa79c0f32df27bc874d76350b572847a880a22f.InstalledAppsRequestBuilder) {
    return i7aee2305413b3b465064d02a6aa79c0f32df27bc874d76350b572847a880a22f.NewInstalledAppsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InstalledAppsById provides operations to manage the installedApps property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) InstalledAppsById(id string)(*i6e500f69af937bd22bbfc2b6ee2698da77bfc652c449bb8ce6a2ee75470b9de9.TeamsAppInstallationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsAppInstallation%2Did"] = id
    }
    return i6e500f69af937bd22bbfc2b6ee2698da77bfc652c449bb8ce6a2ee75470b9de9.NewTeamsAppInstallationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Members provides operations to manage the members property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) Members()(*i97dd0d985d4eaf1cbaeb3d771ef48c101ec691107d9a79dcd3f38213ccae863b.MembersRequestBuilder) {
    return i97dd0d985d4eaf1cbaeb3d771ef48c101ec691107d9a79dcd3f38213ccae863b.NewMembersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MembersById provides operations to manage the members property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) MembersById(id string)(*ie59099b1e42dc6f55628ccb467e95d4e7416a9bd33399bd6beff110ff67ea34e.ConversationMemberItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["conversationMember%2Did"] = id
    }
    return ie59099b1e42dc6f55628ccb467e95d4e7416a9bd33399bd6beff110ff67ea34e.NewConversationMemberItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Operations provides operations to manage the operations property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) Operations()(*i9f6a3c1c4300552f4792a53f0d0fe603636d1b865d408ff28e1cdd722683ebe3.OperationsRequestBuilder) {
    return i9f6a3c1c4300552f4792a53f0d0fe603636d1b865d408ff28e1cdd722683ebe3.NewOperationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OperationsById provides operations to manage the operations property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) OperationsById(id string)(*id5fe3547558f37b7931bf00eb734fbe581be897e9aea405e91e3ce2dbeec74b3.TeamsAsyncOperationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsAsyncOperation%2Did"] = id
    }
    return id5fe3547558f37b7931bf00eb734fbe581be897e9aea405e91e3ce2dbeec74b3.NewTeamsAsyncOperationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch create a new team under a group. In order to create a team, the group must have a least one owner. If the group was created less than 15 minutes ago, it's possible for the Create team call to fail with a 404 error code due to replication delays. The recommended pattern is to retry the Create team call three times, with a 10 second delay between calls.
func (m *TeamRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, requestConfiguration *TeamRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, error) {
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
func (m *TeamRequestBuilder) Photo()(*i946cc8c196bfeaf31caf19315aac4a6c73de740196dc6305506cfd9d4f440ddb.PhotoRequestBuilder) {
    return i946cc8c196bfeaf31caf19315aac4a6c73de740196dc6305506cfd9d4f440ddb.NewPhotoRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PrimaryChannel provides operations to manage the primaryChannel property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) PrimaryChannel()(*i34294a7456473fba252b3ca961ed3381f55b7a770d74ee4720d752d9a45ad739.PrimaryChannelRequestBuilder) {
    return i34294a7456473fba252b3ca961ed3381f55b7a770d74ee4720d752d9a45ad739.NewPrimaryChannelRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Schedule provides operations to manage the schedule property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) Schedule()(*i495e85248bfc2c7bf36f36ede17adf8bc672f00e69df645f352a27b2b7d30e61.ScheduleRequestBuilder) {
    return i495e85248bfc2c7bf36f36ede17adf8bc672f00e69df645f352a27b2b7d30e61.NewScheduleRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SendActivityNotification provides operations to call the sendActivityNotification method.
func (m *TeamRequestBuilder) SendActivityNotification()(*i40ee01a629bf206681620d9e9951dc387487230caf141be819f95c589a8f0779.SendActivityNotificationRequestBuilder) {
    return i40ee01a629bf206681620d9e9951dc387487230caf141be819f95c589a8f0779.NewSendActivityNotificationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Tags provides operations to manage the tags property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) Tags()(*idf8489132b9d5484e716c4faff83b41054a8bc066341cf00a201254e31131154.TagsRequestBuilder) {
    return idf8489132b9d5484e716c4faff83b41054a8bc066341cf00a201254e31131154.NewTagsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TagsById provides operations to manage the tags property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) TagsById(id string)(*idbd66e2d729b9850c37a45204c98fab6700e673bdde41eb4e751af024dc0a897.TeamworkTagItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamworkTag%2Did"] = id
    }
    return idbd66e2d729b9850c37a45204c98fab6700e673bdde41eb4e751af024dc0a897.NewTeamworkTagItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Template provides operations to manage the template property of the microsoft.graph.team entity.
func (m *TeamRequestBuilder) Template()(*i477873d925c8fe5f007dfb4043f82ed6374d80069f96d4931626a23aeb0f2a87.TemplateRequestBuilder) {
    return i477873d925c8fe5f007dfb4043f82ed6374d80069f96d4931626a23aeb0f2a87.NewTemplateRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Unarchive provides operations to call the unarchive method.
func (m *TeamRequestBuilder) Unarchive()(*ie49085238af0f76cc56d306db6cbd7ea879a083402f3fffae10acf7f0e81362d.UnarchiveRequestBuilder) {
    return ie49085238af0f76cc56d306db6cbd7ea879a083402f3fffae10acf7f0e81362d.NewUnarchiveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
