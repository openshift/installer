package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// OrganizerMeetingInfo 
type OrganizerMeetingInfo struct {
    MeetingInfo
    // The organizer property
    organizer IdentitySetable
}
// NewOrganizerMeetingInfo instantiates a new OrganizerMeetingInfo and sets the default values.
func NewOrganizerMeetingInfo()(*OrganizerMeetingInfo) {
    m := &OrganizerMeetingInfo{
        MeetingInfo: *NewMeetingInfo(),
    }
    odataTypeValue := "#microsoft.graph.organizerMeetingInfo";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateOrganizerMeetingInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateOrganizerMeetingInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOrganizerMeetingInfo(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *OrganizerMeetingInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MeetingInfo.GetFieldDeserializers()
    res["organizer"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentitySetFromDiscriminatorValue , m.SetOrganizer)
    return res
}
// GetOrganizer gets the organizer property value. The organizer property
func (m *OrganizerMeetingInfo) GetOrganizer()(IdentitySetable) {
    return m.organizer
}
// Serialize serializes information the current object
func (m *OrganizerMeetingInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MeetingInfo.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("organizer", m.GetOrganizer())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOrganizer sets the organizer property value. The organizer property
func (m *OrganizerMeetingInfo) SetOrganizer(value IdentitySetable)() {
    m.organizer = value
}
