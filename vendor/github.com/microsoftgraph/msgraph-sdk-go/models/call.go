package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Call provides operations to manage the cloudCommunications singleton.
type Call struct {
    Entity
    // The audioRoutingGroups property
    audioRoutingGroups []AudioRoutingGroupable
    // The callback URL on which callbacks will be delivered. Must be https.
    callbackUri *string
    // A unique identifier for all the participant calls in a conference or a unique identifier for two participant calls in a P2P call.  This needs to be copied over from Microsoft.Graph.Call.CallChainId.
    callChainId *string
    // Contains the optional features for the call.
    callOptions CallOptionsable
    // The routing information on how the call was retargeted. Read-only.
    callRoutes []CallRouteable
    // The chat information. Required information for joining a meeting.
    chatInfo ChatInfoable
    // The contentSharingSessions property
    contentSharingSessions []ContentSharingSessionable
    // The direction of the call. The possible value are incoming or outgoing. Read-only.
    direction *CallDirection
    // Call context associated with an incoming call.
    incomingContext IncomingContextable
    // The media configuration. Required.
    mediaConfig MediaConfigable
    // Read-only. The call media state.
    mediaState CallMediaStateable
    // The meeting information that's required for joining a meeting.
    meetingInfo MeetingInfoable
    // The myParticipantId property
    myParticipantId *string
    // The operations property
    operations []CommsOperationable
    // The participants property
    participants []Participantable
    // The list of requested modalities. Possible values are: unknown, audio, video, videoBasedScreenSharing, data.
    requestedModalities []Modality
    // The result information. For example can hold termination reason. Read-only.
    resultInfo ResultInfoable
    // The originator of the call.
    source ParticipantInfoable
    // The call state. Possible values are: incoming, establishing, ringing, established, hold, transferring, transferAccepted, redirecting, terminating, terminated. Read-only.
    state *CallState
    // The subject of the conversation.
    subject *string
    // The targets of the call. Required information for creating peer to peer call.
    targets []InvitationParticipantInfoable
    // The tenantId property
    tenantId *string
    // The toneInfo property
    toneInfo ToneInfoable
    // The transcription information for the call. Read-only.
    transcription CallTranscriptionInfoable
}
// NewCall instantiates a new call and sets the default values.
func NewCall()(*Call) {
    m := &Call{
        Entity: *NewEntity(),
    }
    return m
}
// CreateCallFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCallFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCall(), nil
}
// GetAudioRoutingGroups gets the audioRoutingGroups property value. The audioRoutingGroups property
func (m *Call) GetAudioRoutingGroups()([]AudioRoutingGroupable) {
    return m.audioRoutingGroups
}
// GetCallbackUri gets the callbackUri property value. The callback URL on which callbacks will be delivered. Must be https.
func (m *Call) GetCallbackUri()(*string) {
    return m.callbackUri
}
// GetCallChainId gets the callChainId property value. A unique identifier for all the participant calls in a conference or a unique identifier for two participant calls in a P2P call.  This needs to be copied over from Microsoft.Graph.Call.CallChainId.
func (m *Call) GetCallChainId()(*string) {
    return m.callChainId
}
// GetCallOptions gets the callOptions property value. Contains the optional features for the call.
func (m *Call) GetCallOptions()(CallOptionsable) {
    return m.callOptions
}
// GetCallRoutes gets the callRoutes property value. The routing information on how the call was retargeted. Read-only.
func (m *Call) GetCallRoutes()([]CallRouteable) {
    return m.callRoutes
}
// GetChatInfo gets the chatInfo property value. The chat information. Required information for joining a meeting.
func (m *Call) GetChatInfo()(ChatInfoable) {
    return m.chatInfo
}
// GetContentSharingSessions gets the contentSharingSessions property value. The contentSharingSessions property
func (m *Call) GetContentSharingSessions()([]ContentSharingSessionable) {
    return m.contentSharingSessions
}
// GetDirection gets the direction property value. The direction of the call. The possible value are incoming or outgoing. Read-only.
func (m *Call) GetDirection()(*CallDirection) {
    return m.direction
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Call) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["audioRoutingGroups"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAudioRoutingGroupFromDiscriminatorValue , m.SetAudioRoutingGroups)
    res["callbackUri"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCallbackUri)
    res["callChainId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCallChainId)
    res["callOptions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateCallOptionsFromDiscriminatorValue , m.SetCallOptions)
    res["callRoutes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateCallRouteFromDiscriminatorValue , m.SetCallRoutes)
    res["chatInfo"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateChatInfoFromDiscriminatorValue , m.SetChatInfo)
    res["contentSharingSessions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateContentSharingSessionFromDiscriminatorValue , m.SetContentSharingSessions)
    res["direction"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseCallDirection , m.SetDirection)
    res["incomingContext"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIncomingContextFromDiscriminatorValue , m.SetIncomingContext)
    res["mediaConfig"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateMediaConfigFromDiscriminatorValue , m.SetMediaConfig)
    res["mediaState"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateCallMediaStateFromDiscriminatorValue , m.SetMediaState)
    res["meetingInfo"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateMeetingInfoFromDiscriminatorValue , m.SetMeetingInfo)
    res["myParticipantId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetMyParticipantId)
    res["operations"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateCommsOperationFromDiscriminatorValue , m.SetOperations)
    res["participants"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateParticipantFromDiscriminatorValue , m.SetParticipants)
    res["requestedModalities"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfEnumValues(ParseModality , m.SetRequestedModalities)
    res["resultInfo"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateResultInfoFromDiscriminatorValue , m.SetResultInfo)
    res["source"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateParticipantInfoFromDiscriminatorValue , m.SetSource)
    res["state"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseCallState , m.SetState)
    res["subject"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSubject)
    res["targets"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateInvitationParticipantInfoFromDiscriminatorValue , m.SetTargets)
    res["tenantId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTenantId)
    res["toneInfo"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateToneInfoFromDiscriminatorValue , m.SetToneInfo)
    res["transcription"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateCallTranscriptionInfoFromDiscriminatorValue , m.SetTranscription)
    return res
}
// GetIncomingContext gets the incomingContext property value. Call context associated with an incoming call.
func (m *Call) GetIncomingContext()(IncomingContextable) {
    return m.incomingContext
}
// GetMediaConfig gets the mediaConfig property value. The media configuration. Required.
func (m *Call) GetMediaConfig()(MediaConfigable) {
    return m.mediaConfig
}
// GetMediaState gets the mediaState property value. Read-only. The call media state.
func (m *Call) GetMediaState()(CallMediaStateable) {
    return m.mediaState
}
// GetMeetingInfo gets the meetingInfo property value. The meeting information that's required for joining a meeting.
func (m *Call) GetMeetingInfo()(MeetingInfoable) {
    return m.meetingInfo
}
// GetMyParticipantId gets the myParticipantId property value. The myParticipantId property
func (m *Call) GetMyParticipantId()(*string) {
    return m.myParticipantId
}
// GetOperations gets the operations property value. The operations property
func (m *Call) GetOperations()([]CommsOperationable) {
    return m.operations
}
// GetParticipants gets the participants property value. The participants property
func (m *Call) GetParticipants()([]Participantable) {
    return m.participants
}
// GetRequestedModalities gets the requestedModalities property value. The list of requested modalities. Possible values are: unknown, audio, video, videoBasedScreenSharing, data.
func (m *Call) GetRequestedModalities()([]Modality) {
    return m.requestedModalities
}
// GetResultInfo gets the resultInfo property value. The result information. For example can hold termination reason. Read-only.
func (m *Call) GetResultInfo()(ResultInfoable) {
    return m.resultInfo
}
// GetSource gets the source property value. The originator of the call.
func (m *Call) GetSource()(ParticipantInfoable) {
    return m.source
}
// GetState gets the state property value. The call state. Possible values are: incoming, establishing, ringing, established, hold, transferring, transferAccepted, redirecting, terminating, terminated. Read-only.
func (m *Call) GetState()(*CallState) {
    return m.state
}
// GetSubject gets the subject property value. The subject of the conversation.
func (m *Call) GetSubject()(*string) {
    return m.subject
}
// GetTargets gets the targets property value. The targets of the call. Required information for creating peer to peer call.
func (m *Call) GetTargets()([]InvitationParticipantInfoable) {
    return m.targets
}
// GetTenantId gets the tenantId property value. The tenantId property
func (m *Call) GetTenantId()(*string) {
    return m.tenantId
}
// GetToneInfo gets the toneInfo property value. The toneInfo property
func (m *Call) GetToneInfo()(ToneInfoable) {
    return m.toneInfo
}
// GetTranscription gets the transcription property value. The transcription information for the call. Read-only.
func (m *Call) GetTranscription()(CallTranscriptionInfoable) {
    return m.transcription
}
// Serialize serializes information the current object
func (m *Call) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAudioRoutingGroups() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAudioRoutingGroups())
        err = writer.WriteCollectionOfObjectValues("audioRoutingGroups", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("callbackUri", m.GetCallbackUri())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("callChainId", m.GetCallChainId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("callOptions", m.GetCallOptions())
        if err != nil {
            return err
        }
    }
    if m.GetCallRoutes() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetCallRoutes())
        err = writer.WriteCollectionOfObjectValues("callRoutes", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("chatInfo", m.GetChatInfo())
        if err != nil {
            return err
        }
    }
    if m.GetContentSharingSessions() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetContentSharingSessions())
        err = writer.WriteCollectionOfObjectValues("contentSharingSessions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDirection() != nil {
        cast := (*m.GetDirection()).String()
        err = writer.WriteStringValue("direction", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("incomingContext", m.GetIncomingContext())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mediaConfig", m.GetMediaConfig())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mediaState", m.GetMediaState())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("meetingInfo", m.GetMeetingInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("myParticipantId", m.GetMyParticipantId())
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
    if m.GetParticipants() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetParticipants())
        err = writer.WriteCollectionOfObjectValues("participants", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRequestedModalities() != nil {
        err = writer.WriteCollectionOfStringValues("requestedModalities", SerializeModality(m.GetRequestedModalities()))
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("resultInfo", m.GetResultInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("source", m.GetSource())
        if err != nil {
            return err
        }
    }
    if m.GetState() != nil {
        cast := (*m.GetState()).String()
        err = writer.WriteStringValue("state", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("subject", m.GetSubject())
        if err != nil {
            return err
        }
    }
    if m.GetTargets() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetTargets())
        err = writer.WriteCollectionOfObjectValues("targets", cast)
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
    {
        err = writer.WriteObjectValue("toneInfo", m.GetToneInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("transcription", m.GetTranscription())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAudioRoutingGroups sets the audioRoutingGroups property value. The audioRoutingGroups property
func (m *Call) SetAudioRoutingGroups(value []AudioRoutingGroupable)() {
    m.audioRoutingGroups = value
}
// SetCallbackUri sets the callbackUri property value. The callback URL on which callbacks will be delivered. Must be https.
func (m *Call) SetCallbackUri(value *string)() {
    m.callbackUri = value
}
// SetCallChainId sets the callChainId property value. A unique identifier for all the participant calls in a conference or a unique identifier for two participant calls in a P2P call.  This needs to be copied over from Microsoft.Graph.Call.CallChainId.
func (m *Call) SetCallChainId(value *string)() {
    m.callChainId = value
}
// SetCallOptions sets the callOptions property value. Contains the optional features for the call.
func (m *Call) SetCallOptions(value CallOptionsable)() {
    m.callOptions = value
}
// SetCallRoutes sets the callRoutes property value. The routing information on how the call was retargeted. Read-only.
func (m *Call) SetCallRoutes(value []CallRouteable)() {
    m.callRoutes = value
}
// SetChatInfo sets the chatInfo property value. The chat information. Required information for joining a meeting.
func (m *Call) SetChatInfo(value ChatInfoable)() {
    m.chatInfo = value
}
// SetContentSharingSessions sets the contentSharingSessions property value. The contentSharingSessions property
func (m *Call) SetContentSharingSessions(value []ContentSharingSessionable)() {
    m.contentSharingSessions = value
}
// SetDirection sets the direction property value. The direction of the call. The possible value are incoming or outgoing. Read-only.
func (m *Call) SetDirection(value *CallDirection)() {
    m.direction = value
}
// SetIncomingContext sets the incomingContext property value. Call context associated with an incoming call.
func (m *Call) SetIncomingContext(value IncomingContextable)() {
    m.incomingContext = value
}
// SetMediaConfig sets the mediaConfig property value. The media configuration. Required.
func (m *Call) SetMediaConfig(value MediaConfigable)() {
    m.mediaConfig = value
}
// SetMediaState sets the mediaState property value. Read-only. The call media state.
func (m *Call) SetMediaState(value CallMediaStateable)() {
    m.mediaState = value
}
// SetMeetingInfo sets the meetingInfo property value. The meeting information that's required for joining a meeting.
func (m *Call) SetMeetingInfo(value MeetingInfoable)() {
    m.meetingInfo = value
}
// SetMyParticipantId sets the myParticipantId property value. The myParticipantId property
func (m *Call) SetMyParticipantId(value *string)() {
    m.myParticipantId = value
}
// SetOperations sets the operations property value. The operations property
func (m *Call) SetOperations(value []CommsOperationable)() {
    m.operations = value
}
// SetParticipants sets the participants property value. The participants property
func (m *Call) SetParticipants(value []Participantable)() {
    m.participants = value
}
// SetRequestedModalities sets the requestedModalities property value. The list of requested modalities. Possible values are: unknown, audio, video, videoBasedScreenSharing, data.
func (m *Call) SetRequestedModalities(value []Modality)() {
    m.requestedModalities = value
}
// SetResultInfo sets the resultInfo property value. The result information. For example can hold termination reason. Read-only.
func (m *Call) SetResultInfo(value ResultInfoable)() {
    m.resultInfo = value
}
// SetSource sets the source property value. The originator of the call.
func (m *Call) SetSource(value ParticipantInfoable)() {
    m.source = value
}
// SetState sets the state property value. The call state. Possible values are: incoming, establishing, ringing, established, hold, transferring, transferAccepted, redirecting, terminating, terminated. Read-only.
func (m *Call) SetState(value *CallState)() {
    m.state = value
}
// SetSubject sets the subject property value. The subject of the conversation.
func (m *Call) SetSubject(value *string)() {
    m.subject = value
}
// SetTargets sets the targets property value. The targets of the call. Required information for creating peer to peer call.
func (m *Call) SetTargets(value []InvitationParticipantInfoable)() {
    m.targets = value
}
// SetTenantId sets the tenantId property value. The tenantId property
func (m *Call) SetTenantId(value *string)() {
    m.tenantId = value
}
// SetToneInfo sets the toneInfo property value. The toneInfo property
func (m *Call) SetToneInfo(value ToneInfoable)() {
    m.toneInfo = value
}
// SetTranscription sets the transcription property value. The transcription information for the call. Read-only.
func (m *Call) SetTranscription(value CallTranscriptionInfoable)() {
    m.transcription = value
}
