package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Team provides operations to manage the collection of agreement entities.
type Team struct {
    Entity
    // List of channels either hosted in or shared with the team (incoming channels).
    allChannels []Channelable
    // The collection of channels and messages associated with the team.
    channels []Channelable
    // An optional label. Typically describes the data or business sensitivity of the team. Must match one of a pre-configured set in the tenant's directory.
    classification *string
    // Timestamp at which the team was created.
    createdDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // An optional description for the team. Maximum length: 1024 characters.
    description *string
    // The name of the team.
    displayName *string
    // Settings to configure use of Giphy, memes, and stickers in the team.
    funSettings TeamFunSettingsable
    // The group property
    group Groupable
    // Settings to configure whether guests can create, update, or delete channels in the team.
    guestSettings TeamGuestSettingsable
    // List of channels shared with the team.
    incomingChannels []Channelable
    // The apps installed in this team.
    installedApps []TeamsAppInstallationable
    // A unique ID for the team that has been used in a few places such as the audit log/Office 365 Management Activity API.
    internalId *string
    // Whether this team is in read-only mode.
    isArchived *bool
    // Members and owners of the team.
    members []ConversationMemberable
    // Settings to configure whether members can perform certain actions, for example, create channels and add bots, in the team.
    memberSettings TeamMemberSettingsable
    // Settings to configure messaging and mentions in the team.
    messagingSettings TeamMessagingSettingsable
    // The async operations that ran or are running on this team.
    operations []TeamsAsyncOperationable
    // The profile photo for the team.
    photo ProfilePhotoable
    // The general channel for the team.
    primaryChannel Channelable
    // The schedule of shifts for this team.
    schedule Scheduleable
    // Optional. Indicates whether the team is intended for a particular use case.  Each team specialization has access to unique behaviors and experiences targeted to its use case.
    specialization *TeamSpecialization
    // The summary property
    summary TeamSummaryable
    // The tags associated with the team.
    tags []TeamworkTagable
    // The template this team was created from. See available templates.
    template TeamsTemplateable
    // The ID of the Azure Active Directory tenant.
    tenantId *string
    // The visibility of the group and team. Defaults to Public.
    visibility *TeamVisibilityType
    // A hyperlink that will go to the team in the Microsoft Teams client. This is the URL that you get when you right-click a team in the Microsoft Teams client and select Get link to team. This URL should be treated as an opaque blob, and not parsed.
    webUrl *string
}
// NewTeam instantiates a new team and sets the default values.
func NewTeam()(*Team) {
    m := &Team{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTeamFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeam(), nil
}
// GetAllChannels gets the allChannels property value. List of channels either hosted in or shared with the team (incoming channels).
func (m *Team) GetAllChannels()([]Channelable) {
    return m.allChannels
}
// GetChannels gets the channels property value. The collection of channels and messages associated with the team.
func (m *Team) GetChannels()([]Channelable) {
    return m.channels
}
// GetClassification gets the classification property value. An optional label. Typically describes the data or business sensitivity of the team. Must match one of a pre-configured set in the tenant's directory.
func (m *Team) GetClassification()(*string) {
    return m.classification
}
// GetCreatedDateTime gets the createdDateTime property value. Timestamp at which the team was created.
func (m *Team) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.createdDateTime
}
// GetDescription gets the description property value. An optional description for the team. Maximum length: 1024 characters.
func (m *Team) GetDescription()(*string) {
    return m.description
}
// GetDisplayName gets the displayName property value. The name of the team.
func (m *Team) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Team) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["allChannels"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateChannelFromDiscriminatorValue , m.SetAllChannels)
    res["channels"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateChannelFromDiscriminatorValue , m.SetChannels)
    res["classification"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetClassification)
    res["createdDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCreatedDateTime)
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDescription)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["funSettings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTeamFunSettingsFromDiscriminatorValue , m.SetFunSettings)
    res["group"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateGroupFromDiscriminatorValue , m.SetGroup)
    res["guestSettings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTeamGuestSettingsFromDiscriminatorValue , m.SetGuestSettings)
    res["incomingChannels"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateChannelFromDiscriminatorValue , m.SetIncomingChannels)
    res["installedApps"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTeamsAppInstallationFromDiscriminatorValue , m.SetInstalledApps)
    res["internalId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetInternalId)
    res["isArchived"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsArchived)
    res["members"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateConversationMemberFromDiscriminatorValue , m.SetMembers)
    res["memberSettings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTeamMemberSettingsFromDiscriminatorValue , m.SetMemberSettings)
    res["messagingSettings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTeamMessagingSettingsFromDiscriminatorValue , m.SetMessagingSettings)
    res["operations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTeamsAsyncOperationFromDiscriminatorValue , m.SetOperations)
    res["photo"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateProfilePhotoFromDiscriminatorValue , m.SetPhoto)
    res["primaryChannel"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateChannelFromDiscriminatorValue , m.SetPrimaryChannel)
    res["schedule"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateScheduleFromDiscriminatorValue , m.SetSchedule)
    res["specialization"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseTeamSpecialization , m.SetSpecialization)
    res["summary"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTeamSummaryFromDiscriminatorValue , m.SetSummary)
    res["tags"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateTeamworkTagFromDiscriminatorValue , m.SetTags)
    res["template"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateTeamsTemplateFromDiscriminatorValue , m.SetTemplate)
    res["tenantId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTenantId)
    res["visibility"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseTeamVisibilityType , m.SetVisibility)
    res["webUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetWebUrl)
    return res
}
// GetFunSettings gets the funSettings property value. Settings to configure use of Giphy, memes, and stickers in the team.
func (m *Team) GetFunSettings()(TeamFunSettingsable) {
    return m.funSettings
}
// GetGroup gets the group property value. The group property
func (m *Team) GetGroup()(Groupable) {
    return m.group
}
// GetGuestSettings gets the guestSettings property value. Settings to configure whether guests can create, update, or delete channels in the team.
func (m *Team) GetGuestSettings()(TeamGuestSettingsable) {
    return m.guestSettings
}
// GetIncomingChannels gets the incomingChannels property value. List of channels shared with the team.
func (m *Team) GetIncomingChannels()([]Channelable) {
    return m.incomingChannels
}
// GetInstalledApps gets the installedApps property value. The apps installed in this team.
func (m *Team) GetInstalledApps()([]TeamsAppInstallationable) {
    return m.installedApps
}
// GetInternalId gets the internalId property value. A unique ID for the team that has been used in a few places such as the audit log/Office 365 Management Activity API.
func (m *Team) GetInternalId()(*string) {
    return m.internalId
}
// GetIsArchived gets the isArchived property value. Whether this team is in read-only mode.
func (m *Team) GetIsArchived()(*bool) {
    return m.isArchived
}
// GetMembers gets the members property value. Members and owners of the team.
func (m *Team) GetMembers()([]ConversationMemberable) {
    return m.members
}
// GetMemberSettings gets the memberSettings property value. Settings to configure whether members can perform certain actions, for example, create channels and add bots, in the team.
func (m *Team) GetMemberSettings()(TeamMemberSettingsable) {
    return m.memberSettings
}
// GetMessagingSettings gets the messagingSettings property value. Settings to configure messaging and mentions in the team.
func (m *Team) GetMessagingSettings()(TeamMessagingSettingsable) {
    return m.messagingSettings
}
// GetOperations gets the operations property value. The async operations that ran or are running on this team.
func (m *Team) GetOperations()([]TeamsAsyncOperationable) {
    return m.operations
}
// GetPhoto gets the photo property value. The profile photo for the team.
func (m *Team) GetPhoto()(ProfilePhotoable) {
    return m.photo
}
// GetPrimaryChannel gets the primaryChannel property value. The general channel for the team.
func (m *Team) GetPrimaryChannel()(Channelable) {
    return m.primaryChannel
}
// GetSchedule gets the schedule property value. The schedule of shifts for this team.
func (m *Team) GetSchedule()(Scheduleable) {
    return m.schedule
}
// GetSpecialization gets the specialization property value. Optional. Indicates whether the team is intended for a particular use case.  Each team specialization has access to unique behaviors and experiences targeted to its use case.
func (m *Team) GetSpecialization()(*TeamSpecialization) {
    return m.specialization
}
// GetSummary gets the summary property value. The summary property
func (m *Team) GetSummary()(TeamSummaryable) {
    return m.summary
}
// GetTags gets the tags property value. The tags associated with the team.
func (m *Team) GetTags()([]TeamworkTagable) {
    return m.tags
}
// GetTemplate gets the template property value. The template this team was created from. See available templates.
func (m *Team) GetTemplate()(TeamsTemplateable) {
    return m.template
}
// GetTenantId gets the tenantId property value. The ID of the Azure Active Directory tenant.
func (m *Team) GetTenantId()(*string) {
    return m.tenantId
}
// GetVisibility gets the visibility property value. The visibility of the group and team. Defaults to Public.
func (m *Team) GetVisibility()(*TeamVisibilityType) {
    return m.visibility
}
// GetWebUrl gets the webUrl property value. A hyperlink that will go to the team in the Microsoft Teams client. This is the URL that you get when you right-click a team in the Microsoft Teams client and select Get link to team. This URL should be treated as an opaque blob, and not parsed.
func (m *Team) GetWebUrl()(*string) {
    return m.webUrl
}
// Serialize serializes information the current object
func (m *Team) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAllChannels() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAllChannels())
        err = writer.WriteCollectionOfObjectValues("allChannels", cast)
        if err != nil {
            return err
        }
    }
    if m.GetChannels() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetChannels())
        err = writer.WriteCollectionOfObjectValues("channels", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("classification", m.GetClassification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("funSettings", m.GetFunSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("group", m.GetGroup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("guestSettings", m.GetGuestSettings())
        if err != nil {
            return err
        }
    }
    if m.GetIncomingChannels() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetIncomingChannels())
        err = writer.WriteCollectionOfObjectValues("incomingChannels", cast)
        if err != nil {
            return err
        }
    }
    if m.GetInstalledApps() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetInstalledApps())
        err = writer.WriteCollectionOfObjectValues("installedApps", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("internalId", m.GetInternalId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isArchived", m.GetIsArchived())
        if err != nil {
            return err
        }
    }
    if m.GetMembers() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetMembers())
        err = writer.WriteCollectionOfObjectValues("members", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("memberSettings", m.GetMemberSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("messagingSettings", m.GetMessagingSettings())
        if err != nil {
            return err
        }
    }
    if m.GetOperations() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetOperations())
        err = writer.WriteCollectionOfObjectValues("operations", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("photo", m.GetPhoto())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("primaryChannel", m.GetPrimaryChannel())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("schedule", m.GetSchedule())
        if err != nil {
            return err
        }
    }
    if m.GetSpecialization() != nil {
        cast := (*m.GetSpecialization()).String()
        err = writer.WriteStringValue("specialization", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("summary", m.GetSummary())
        if err != nil {
            return err
        }
    }
    if m.GetTags() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTags())
        err = writer.WriteCollectionOfObjectValues("tags", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("template", m.GetTemplate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("tenantId", m.GetTenantId())
        if err != nil {
            return err
        }
    }
    if m.GetVisibility() != nil {
        cast := (*m.GetVisibility()).String()
        err = writer.WriteStringValue("visibility", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("webUrl", m.GetWebUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllChannels sets the allChannels property value. List of channels either hosted in or shared with the team (incoming channels).
func (m *Team) SetAllChannels(value []Channelable)() {
    m.allChannels = value
}
// SetChannels sets the channels property value. The collection of channels and messages associated with the team.
func (m *Team) SetChannels(value []Channelable)() {
    m.channels = value
}
// SetClassification sets the classification property value. An optional label. Typically describes the data or business sensitivity of the team. Must match one of a pre-configured set in the tenant's directory.
func (m *Team) SetClassification(value *string)() {
    m.classification = value
}
// SetCreatedDateTime sets the createdDateTime property value. Timestamp at which the team was created.
func (m *Team) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.createdDateTime = value
}
// SetDescription sets the description property value. An optional description for the team. Maximum length: 1024 characters.
func (m *Team) SetDescription(value *string)() {
    m.description = value
}
// SetDisplayName sets the displayName property value. The name of the team.
func (m *Team) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetFunSettings sets the funSettings property value. Settings to configure use of Giphy, memes, and stickers in the team.
func (m *Team) SetFunSettings(value TeamFunSettingsable)() {
    m.funSettings = value
}
// SetGroup sets the group property value. The group property
func (m *Team) SetGroup(value Groupable)() {
    m.group = value
}
// SetGuestSettings sets the guestSettings property value. Settings to configure whether guests can create, update, or delete channels in the team.
func (m *Team) SetGuestSettings(value TeamGuestSettingsable)() {
    m.guestSettings = value
}
// SetIncomingChannels sets the incomingChannels property value. List of channels shared with the team.
func (m *Team) SetIncomingChannels(value []Channelable)() {
    m.incomingChannels = value
}
// SetInstalledApps sets the installedApps property value. The apps installed in this team.
func (m *Team) SetInstalledApps(value []TeamsAppInstallationable)() {
    m.installedApps = value
}
// SetInternalId sets the internalId property value. A unique ID for the team that has been used in a few places such as the audit log/Office 365 Management Activity API.
func (m *Team) SetInternalId(value *string)() {
    m.internalId = value
}
// SetIsArchived sets the isArchived property value. Whether this team is in read-only mode.
func (m *Team) SetIsArchived(value *bool)() {
    m.isArchived = value
}
// SetMembers sets the members property value. Members and owners of the team.
func (m *Team) SetMembers(value []ConversationMemberable)() {
    m.members = value
}
// SetMemberSettings sets the memberSettings property value. Settings to configure whether members can perform certain actions, for example, create channels and add bots, in the team.
func (m *Team) SetMemberSettings(value TeamMemberSettingsable)() {
    m.memberSettings = value
}
// SetMessagingSettings sets the messagingSettings property value. Settings to configure messaging and mentions in the team.
func (m *Team) SetMessagingSettings(value TeamMessagingSettingsable)() {
    m.messagingSettings = value
}
// SetOperations sets the operations property value. The async operations that ran or are running on this team.
func (m *Team) SetOperations(value []TeamsAsyncOperationable)() {
    m.operations = value
}
// SetPhoto sets the photo property value. The profile photo for the team.
func (m *Team) SetPhoto(value ProfilePhotoable)() {
    m.photo = value
}
// SetPrimaryChannel sets the primaryChannel property value. The general channel for the team.
func (m *Team) SetPrimaryChannel(value Channelable)() {
    m.primaryChannel = value
}
// SetSchedule sets the schedule property value. The schedule of shifts for this team.
func (m *Team) SetSchedule(value Scheduleable)() {
    m.schedule = value
}
// SetSpecialization sets the specialization property value. Optional. Indicates whether the team is intended for a particular use case.  Each team specialization has access to unique behaviors and experiences targeted to its use case.
func (m *Team) SetSpecialization(value *TeamSpecialization)() {
    m.specialization = value
}
// SetSummary sets the summary property value. The summary property
func (m *Team) SetSummary(value TeamSummaryable)() {
    m.summary = value
}
// SetTags sets the tags property value. The tags associated with the team.
func (m *Team) SetTags(value []TeamworkTagable)() {
    m.tags = value
}
// SetTemplate sets the template property value. The template this team was created from. See available templates.
func (m *Team) SetTemplate(value TeamsTemplateable)() {
    m.template = value
}
// SetTenantId sets the tenantId property value. The ID of the Azure Active Directory tenant.
func (m *Team) SetTenantId(value *string)() {
    m.tenantId = value
}
// SetVisibility sets the visibility property value. The visibility of the group and team. Defaults to Public.
func (m *Team) SetVisibility(value *TeamVisibilityType)() {
    m.visibility = value
}
// SetWebUrl sets the webUrl property value. A hyperlink that will go to the team in the Microsoft Teams client. This is the URL that you get when you right-click a team in the Microsoft Teams client and select Get link to team. This URL should be treated as an opaque blob, and not parsed.
func (m *Team) SetWebUrl(value *string)() {
    m.webUrl = value
}
