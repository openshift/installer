package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ImportedWindowsAutopilotDeviceIdentityUpload 
type ImportedWindowsAutopilotDeviceIdentityUpload struct {
    Entity
    // DateTime when the entity is created.
    createdDateTimeUtc *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Collection of all Autopilot devices as a part of this upload.
    deviceIdentities []ImportedWindowsAutopilotDeviceIdentityable
    // The status property
    status *ImportedWindowsAutopilotDeviceIdentityUploadStatus
}
// NewImportedWindowsAutopilotDeviceIdentityUpload instantiates a new ImportedWindowsAutopilotDeviceIdentityUpload and sets the default values.
func NewImportedWindowsAutopilotDeviceIdentityUpload()(*ImportedWindowsAutopilotDeviceIdentityUpload) {
    m := &ImportedWindowsAutopilotDeviceIdentityUpload{
        Entity: *NewEntity(),
    }
    return m
}
// CreateImportedWindowsAutopilotDeviceIdentityUploadFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateImportedWindowsAutopilotDeviceIdentityUploadFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewImportedWindowsAutopilotDeviceIdentityUpload(), nil
}
// GetCreatedDateTimeUtc gets the createdDateTimeUtc property value. DateTime when the entity is created.
func (m *ImportedWindowsAutopilotDeviceIdentityUpload) GetCreatedDateTimeUtc()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.createdDateTimeUtc
}
// GetDeviceIdentities gets the deviceIdentities property value. Collection of all Autopilot devices as a part of this upload.
func (m *ImportedWindowsAutopilotDeviceIdentityUpload) GetDeviceIdentities()([]ImportedWindowsAutopilotDeviceIdentityable) {
    return m.deviceIdentities
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ImportedWindowsAutopilotDeviceIdentityUpload) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["createdDateTimeUtc"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCreatedDateTimeUtc)
    res["deviceIdentities"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateImportedWindowsAutopilotDeviceIdentityFromDiscriminatorValue , m.SetDeviceIdentities)
    res["status"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseImportedWindowsAutopilotDeviceIdentityUploadStatus , m.SetStatus)
    return res
}
// GetStatus gets the status property value. The status property
func (m *ImportedWindowsAutopilotDeviceIdentityUpload) GetStatus()(*ImportedWindowsAutopilotDeviceIdentityUploadStatus) {
    return m.status
}
// Serialize serializes information the current object
func (m *ImportedWindowsAutopilotDeviceIdentityUpload) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("createdDateTimeUtc", m.GetCreatedDateTimeUtc())
        if err != nil {
            return err
        }
    }
    if m.GetDeviceIdentities() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDeviceIdentities())
        err = writer.WriteCollectionOfObjectValues("deviceIdentities", cast)
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCreatedDateTimeUtc sets the createdDateTimeUtc property value. DateTime when the entity is created.
func (m *ImportedWindowsAutopilotDeviceIdentityUpload) SetCreatedDateTimeUtc(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.createdDateTimeUtc = value
}
// SetDeviceIdentities sets the deviceIdentities property value. Collection of all Autopilot devices as a part of this upload.
func (m *ImportedWindowsAutopilotDeviceIdentityUpload) SetDeviceIdentities(value []ImportedWindowsAutopilotDeviceIdentityable)() {
    m.deviceIdentities = value
}
// SetStatus sets the status property value. The status property
func (m *ImportedWindowsAutopilotDeviceIdentityUpload) SetStatus(value *ImportedWindowsAutopilotDeviceIdentityUploadStatus)() {
    m.status = value
}
