package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TabUpdatedEventMessageDetail 
type TabUpdatedEventMessageDetail struct {
    EventMessageDetail
    // Initiator of the event.
    initiator IdentitySetable
    // Unique identifier of the tab.
    tabId *string
}
// NewTabUpdatedEventMessageDetail instantiates a new TabUpdatedEventMessageDetail and sets the default values.
func NewTabUpdatedEventMessageDetail()(*TabUpdatedEventMessageDetail) {
    m := &TabUpdatedEventMessageDetail{
        EventMessageDetail: *NewEventMessageDetail(),
    }
    odataTypeValue := "#microsoft.graph.tabUpdatedEventMessageDetail";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateTabUpdatedEventMessageDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateTabUpdatedEventMessageDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTabUpdatedEventMessageDetail(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *TabUpdatedEventMessageDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EventMessageDetail.GetFieldDeserializers()
    res["initiator"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentitySetFromDiscriminatorValue , m.SetInitiator)
    res["tabId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTabId)
    return res
}
// GetInitiator gets the initiator property value. Initiator of the event.
func (m *TabUpdatedEventMessageDetail) GetInitiator()(IdentitySetable) {
    return m.initiator
}
// GetTabId gets the tabId property value. Unique identifier of the tab.
func (m *TabUpdatedEventMessageDetail) GetTabId()(*string) {
    return m.tabId
}
// Serialize serializes information the current object
func (m *TabUpdatedEventMessageDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EventMessageDetail.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("initiator", m.GetInitiator())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("tabId", m.GetTabId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetInitiator sets the initiator property value. Initiator of the event.
func (m *TabUpdatedEventMessageDetail) SetInitiator(value IdentitySetable)() {
    m.initiator = value
}
// SetTabId sets the tabId property value. Unique identifier of the tab.
func (m *TabUpdatedEventMessageDetail) SetTabId(value *string)() {
    m.tabId = value
}
