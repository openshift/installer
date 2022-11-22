package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Win32LobAppRestartSettings contains properties describing restart coordination following an app installation.
type Win32LobAppRestartSettings struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The number of minutes before the restart time to display the countdown dialog for pending restarts.
    countdownDisplayBeforeRestartInMinutes *int32
    // The number of minutes to wait before restarting the device after an app installation.
    gracePeriodInMinutes *int32
    // The OdataType property
    odataType *string
    // The number of minutes to snooze the restart notification dialog when the snooze button is selected.
    restartNotificationSnoozeDurationInMinutes *int32
}
// NewWin32LobAppRestartSettings instantiates a new win32LobAppRestartSettings and sets the default values.
func NewWin32LobAppRestartSettings()(*Win32LobAppRestartSettings) {
    m := &Win32LobAppRestartSettings{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateWin32LobAppRestartSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWin32LobAppRestartSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWin32LobAppRestartSettings(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Win32LobAppRestartSettings) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCountdownDisplayBeforeRestartInMinutes gets the countdownDisplayBeforeRestartInMinutes property value. The number of minutes before the restart time to display the countdown dialog for pending restarts.
func (m *Win32LobAppRestartSettings) GetCountdownDisplayBeforeRestartInMinutes()(*int32) {
    return m.countdownDisplayBeforeRestartInMinutes
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Win32LobAppRestartSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["countdownDisplayBeforeRestartInMinutes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetCountdownDisplayBeforeRestartInMinutes)
    res["gracePeriodInMinutes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetGracePeriodInMinutes)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["restartNotificationSnoozeDurationInMinutes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetRestartNotificationSnoozeDurationInMinutes)
    return res
}
// GetGracePeriodInMinutes gets the gracePeriodInMinutes property value. The number of minutes to wait before restarting the device after an app installation.
func (m *Win32LobAppRestartSettings) GetGracePeriodInMinutes()(*int32) {
    return m.gracePeriodInMinutes
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *Win32LobAppRestartSettings) GetOdataType()(*string) {
    return m.odataType
}
// GetRestartNotificationSnoozeDurationInMinutes gets the restartNotificationSnoozeDurationInMinutes property value. The number of minutes to snooze the restart notification dialog when the snooze button is selected.
func (m *Win32LobAppRestartSettings) GetRestartNotificationSnoozeDurationInMinutes()(*int32) {
    return m.restartNotificationSnoozeDurationInMinutes
}
// Serialize serializes information the current object
func (m *Win32LobAppRestartSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("countdownDisplayBeforeRestartInMinutes", m.GetCountdownDisplayBeforeRestartInMinutes())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("gracePeriodInMinutes", m.GetGracePeriodInMinutes())
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
        err := writer.WriteInt32Value("restartNotificationSnoozeDurationInMinutes", m.GetRestartNotificationSnoozeDurationInMinutes())
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
func (m *Win32LobAppRestartSettings) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCountdownDisplayBeforeRestartInMinutes sets the countdownDisplayBeforeRestartInMinutes property value. The number of minutes before the restart time to display the countdown dialog for pending restarts.
func (m *Win32LobAppRestartSettings) SetCountdownDisplayBeforeRestartInMinutes(value *int32)() {
    m.countdownDisplayBeforeRestartInMinutes = value
}
// SetGracePeriodInMinutes sets the gracePeriodInMinutes property value. The number of minutes to wait before restarting the device after an app installation.
func (m *Win32LobAppRestartSettings) SetGracePeriodInMinutes(value *int32)() {
    m.gracePeriodInMinutes = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Win32LobAppRestartSettings) SetOdataType(value *string)() {
    m.odataType = value
}
// SetRestartNotificationSnoozeDurationInMinutes sets the restartNotificationSnoozeDurationInMinutes property value. The number of minutes to snooze the restart notification dialog when the snooze button is selected.
func (m *Win32LobAppRestartSettings) SetRestartNotificationSnoozeDurationInMinutes(value *int32)() {
    m.restartNotificationSnoozeDurationInMinutes = value
}
