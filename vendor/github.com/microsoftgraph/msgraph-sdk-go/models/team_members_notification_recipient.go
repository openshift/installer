package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TeamMembersNotificationRecipient 
type TeamMembersNotificationRecipient struct {
    TeamworkNotificationRecipient
    // The unique identifier for the team whose members should receive the notification.
    teamId *string
}
// NewTeamMembersNotificationRecipient instantiates a new TeamMembersNotificationRecipient and sets the default values.
func NewTeamMembersNotificationRecipient()(*TeamMembersNotificationRecipient) {
    m := &TeamMembersNotificationRecipient{
        TeamworkNotificationRecipient: *NewTeamworkNotificationRecipient(),
    }
    odataTypeValue := "#microsoft.graph.teamMembersNotificationRecipient";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateTeamMembersNotificationRecipientFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTeamMembersNotificationRecipientFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamMembersNotificationRecipient(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TeamMembersNotificationRecipient) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.TeamworkNotificationRecipient.GetFieldDeserializers()
    res["teamId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTeamId)
    return res
}
// GetTeamId gets the teamId property value. The unique identifier for the team whose members should receive the notification.
func (m *TeamMembersNotificationRecipient) GetTeamId()(*string) {
    return m.teamId
}
// Serialize serializes information the current object
func (m *TeamMembersNotificationRecipient) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.TeamworkNotificationRecipient.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("teamId", m.GetTeamId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetTeamId sets the teamId property value. The unique identifier for the team whose members should receive the notification.
func (m *TeamMembersNotificationRecipient) SetTeamId(value *string)() {
    m.teamId = value
}
