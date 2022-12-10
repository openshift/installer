package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceComplianceActionItem scheduled Action Configuration
type DeviceComplianceActionItem struct {
    Entity
    // Scheduled Action Type Enum
    actionType *DeviceComplianceActionType
    // Number of hours to wait till the action will be enforced. Valid values 0 to 8760
    gracePeriodHours *int32
    // A list of group IDs to speicify who to CC this notification message to.
    notificationMessageCCList []string
    // What notification Message template to use
    notificationTemplateId *string
}
// NewDeviceComplianceActionItem instantiates a new deviceComplianceActionItem and sets the default values.
func NewDeviceComplianceActionItem()(*DeviceComplianceActionItem) {
    m := &DeviceComplianceActionItem{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceComplianceActionItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDeviceComplianceActionItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceComplianceActionItem(), nil
}
// GetActionType gets the actionType property value. Scheduled Action Type Enum
func (m *DeviceComplianceActionItem) GetActionType()(*DeviceComplianceActionType) {
    return m.actionType
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DeviceComplianceActionItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["actionType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseDeviceComplianceActionType , m.SetActionType)
    res["gracePeriodHours"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetGracePeriodHours)
    res["notificationMessageCCList"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetNotificationMessageCCList)
    res["notificationTemplateId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetNotificationTemplateId)
    return res
}
// GetGracePeriodHours gets the gracePeriodHours property value. Number of hours to wait till the action will be enforced. Valid values 0 to 8760
func (m *DeviceComplianceActionItem) GetGracePeriodHours()(*int32) {
    return m.gracePeriodHours
}
// GetNotificationMessageCCList gets the notificationMessageCCList property value. A list of group IDs to speicify who to CC this notification message to.
func (m *DeviceComplianceActionItem) GetNotificationMessageCCList()([]string) {
    return m.notificationMessageCCList
}
// GetNotificationTemplateId gets the notificationTemplateId property value. What notification Message template to use
func (m *DeviceComplianceActionItem) GetNotificationTemplateId()(*string) {
    return m.notificationTemplateId
}
// Serialize serializes information the current object
func (m *DeviceComplianceActionItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetActionType() != nil {
        cast := (*m.GetActionType()).String()
        err = writer.WriteStringValue("actionType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("gracePeriodHours", m.GetGracePeriodHours())
        if err != nil {
            return err
        }
    }
    if m.GetNotificationMessageCCList() != nil {
        err = writer.WriteCollectionOfStringValues("notificationMessageCCList", m.GetNotificationMessageCCList())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("notificationTemplateId", m.GetNotificationTemplateId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActionType sets the actionType property value. Scheduled Action Type Enum
func (m *DeviceComplianceActionItem) SetActionType(value *DeviceComplianceActionType)() {
    m.actionType = value
}
// SetGracePeriodHours sets the gracePeriodHours property value. Number of hours to wait till the action will be enforced. Valid values 0 to 8760
func (m *DeviceComplianceActionItem) SetGracePeriodHours(value *int32)() {
    m.gracePeriodHours = value
}
// SetNotificationMessageCCList sets the notificationMessageCCList property value. A list of group IDs to speicify who to CC this notification message to.
func (m *DeviceComplianceActionItem) SetNotificationMessageCCList(value []string)() {
    m.notificationMessageCCList = value
}
// SetNotificationTemplateId sets the notificationTemplateId property value. What notification Message template to use
func (m *DeviceComplianceActionItem) SetNotificationTemplateId(value *string)() {
    m.notificationTemplateId = value
}
