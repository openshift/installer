package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RoomList 
type RoomList struct {
    Place
    // The email address of the room list.
    emailAddress *string
    // The rooms property
    rooms []Roomable
}
// NewRoomList instantiates a new RoomList and sets the default values.
func NewRoomList()(*RoomList) {
    m := &RoomList{
        Place: *NewPlace(),
    }
    odataTypeValue := "#microsoft.graph.roomList";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateRoomListFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateRoomListFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRoomList(), nil
}
// GetEmailAddress gets the emailAddress property value. The email address of the room list.
func (m *RoomList) GetEmailAddress()(*string) {
    return m.emailAddress
}
// GetFieldDeserializers the deserialization information for the current model
func (m *RoomList) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Place.GetFieldDeserializers()
    res["emailAddress"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetEmailAddress)
    res["rooms"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateRoomFromDiscriminatorValue , m.SetRooms)
    return res
}
// GetRooms gets the rooms property value. The rooms property
func (m *RoomList) GetRooms()([]Roomable) {
    return m.rooms
}
// Serialize serializes information the current object
func (m *RoomList) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Place.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("emailAddress", m.GetEmailAddress())
        if err != nil {
            return err
        }
    }
    if m.GetRooms() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetRooms())
        err = writer.WriteCollectionOfObjectValues("rooms", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEmailAddress sets the emailAddress property value. The email address of the room list.
func (m *RoomList) SetEmailAddress(value *string)() {
    m.emailAddress = value
}
// SetRooms sets the rooms property value. The rooms property
func (m *RoomList) SetRooms(value []Roomable)() {
    m.rooms = value
}
