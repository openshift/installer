package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// InvitationParticipantInfo 
type InvitationParticipantInfo struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Optional. Whether to hide the participant from the roster.
    hidden *bool
    // The identity property
    identity IdentitySetable
    // The OdataType property
    odataType *string
    // Optional. The ID of the target participant.
    participantId *string
    // Optional. Whether to remove them from the main mixer.
    removeFromDefaultAudioRoutingGroup *bool
    // Optional. The call which the target identity is currently a part of. For peer-to-peer case, the call will be dropped once the participant is added successfully.
    replacesCallId *string
}
// NewInvitationParticipantInfo instantiates a new invitationParticipantInfo and sets the default values.
func NewInvitationParticipantInfo()(*InvitationParticipantInfo) {
    m := &InvitationParticipantInfo{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateInvitationParticipantInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateInvitationParticipantInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewInvitationParticipantInfo(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *InvitationParticipantInfo) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *InvitationParticipantInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["hidden"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetHidden)
    res["identity"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentitySetFromDiscriminatorValue , m.SetIdentity)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["participantId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetParticipantId)
    res["removeFromDefaultAudioRoutingGroup"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetRemoveFromDefaultAudioRoutingGroup)
    res["replacesCallId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetReplacesCallId)
    return res
}
// GetHidden gets the hidden property value. Optional. Whether to hide the participant from the roster.
func (m *InvitationParticipantInfo) GetHidden()(*bool) {
    return m.hidden
}
// GetIdentity gets the identity property value. The identity property
func (m *InvitationParticipantInfo) GetIdentity()(IdentitySetable) {
    return m.identity
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *InvitationParticipantInfo) GetOdataType()(*string) {
    return m.odataType
}
// GetParticipantId gets the participantId property value. Optional. The ID of the target participant.
func (m *InvitationParticipantInfo) GetParticipantId()(*string) {
    return m.participantId
}
// GetRemoveFromDefaultAudioRoutingGroup gets the removeFromDefaultAudioRoutingGroup property value. Optional. Whether to remove them from the main mixer.
func (m *InvitationParticipantInfo) GetRemoveFromDefaultAudioRoutingGroup()(*bool) {
    return m.removeFromDefaultAudioRoutingGroup
}
// GetReplacesCallId gets the replacesCallId property value. Optional. The call which the target identity is currently a part of. For peer-to-peer case, the call will be dropped once the participant is added successfully.
func (m *InvitationParticipantInfo) GetReplacesCallId()(*string) {
    return m.replacesCallId
}
// Serialize serializes information the current object
func (m *InvitationParticipantInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("hidden", m.GetHidden())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("identity", m.GetIdentity())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("participantId", m.GetParticipantId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("removeFromDefaultAudioRoutingGroup", m.GetRemoveFromDefaultAudioRoutingGroup())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("replacesCallId", m.GetReplacesCallId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *InvitationParticipantInfo) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetHidden sets the hidden property value. Optional. Whether to hide the participant from the roster.
func (m *InvitationParticipantInfo) SetHidden(value *bool)() {
    m.hidden = value
}
// SetIdentity sets the identity property value. The identity property
func (m *InvitationParticipantInfo) SetIdentity(value IdentitySetable)() {
    m.identity = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *InvitationParticipantInfo) SetOdataType(value *string)() {
    m.odataType = value
}
// SetParticipantId sets the participantId property value. Optional. The ID of the target participant.
func (m *InvitationParticipantInfo) SetParticipantId(value *string)() {
    m.participantId = value
}
// SetRemoveFromDefaultAudioRoutingGroup sets the removeFromDefaultAudioRoutingGroup property value. Optional. Whether to remove them from the main mixer.
func (m *InvitationParticipantInfo) SetRemoveFromDefaultAudioRoutingGroup(value *bool)() {
    m.removeFromDefaultAudioRoutingGroup = value
}
// SetReplacesCallId sets the replacesCallId property value. Optional. The call which the target identity is currently a part of. For peer-to-peer case, the call will be dropped once the participant is added successfully.
func (m *InvitationParticipantInfo) SetReplacesCallId(value *string)() {
    m.replacesCallId = value
}
