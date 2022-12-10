package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceManagementReports 
type DeviceManagementReports struct {
    Entity
    // Entity representing a job to export a report
    exportJobs []DeviceManagementExportJobable
}
// NewDeviceManagementReports instantiates a new deviceManagementReports and sets the default values.
func NewDeviceManagementReports()(*DeviceManagementReports) {
    m := &DeviceManagementReports{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceManagementReportsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDeviceManagementReportsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceManagementReports(), nil
}
// GetExportJobs gets the exportJobs property value. Entity representing a job to export a report
func (m *DeviceManagementReports) GetExportJobs()([]DeviceManagementExportJobable) {
    return m.exportJobs
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DeviceManagementReports) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["exportJobs"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateDeviceManagementExportJobFromDiscriminatorValue , m.SetExportJobs)
    return res
}
// Serialize serializes information the current object
func (m *DeviceManagementReports) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetExportJobs() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetExportJobs())
        err = writer.WriteCollectionOfObjectValues("exportJobs", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetExportJobs sets the exportJobs property value. Entity representing a job to export a report
func (m *DeviceManagementReports) SetExportJobs(value []DeviceManagementExportJobable)() {
    m.exportJobs = value
}
