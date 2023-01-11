package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedMobileApp the identifier for the deployment an app.
type ManagedMobileApp struct {
    Entity
    // The identifier for an app with it's operating system type.
    mobileAppIdentifier MobileAppIdentifierable
    // Version of the entity.
    version *string
}
// NewManagedMobileApp instantiates a new managedMobileApp and sets the default values.
func NewManagedMobileApp()(*ManagedMobileApp) {
    m := &ManagedMobileApp{
        Entity: *NewEntity(),
    }
    return m
}
// CreateManagedMobileAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateManagedMobileAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewManagedMobileApp(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ManagedMobileApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["mobileAppIdentifier"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateMobileAppIdentifierFromDiscriminatorValue , m.SetMobileAppIdentifier)
    res["version"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetVersion)
    return res
}
// GetMobileAppIdentifier gets the mobileAppIdentifier property value. The identifier for an app with it's operating system type.
func (m *ManagedMobileApp) GetMobileAppIdentifier()(MobileAppIdentifierable) {
    return m.mobileAppIdentifier
}
// GetVersion gets the version property value. Version of the entity.
func (m *ManagedMobileApp) GetVersion()(*string) {
    return m.version
}
// Serialize serializes information the current object
func (m *ManagedMobileApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("mobileAppIdentifier", m.GetMobileAppIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMobileAppIdentifier sets the mobileAppIdentifier property value. The identifier for an app with it's operating system type.
func (m *ManagedMobileApp) SetMobileAppIdentifier(value MobileAppIdentifierable)() {
    m.mobileAppIdentifier = value
}
// SetVersion sets the version property value. Version of the entity.
func (m *ManagedMobileApp) SetVersion(value *string)() {
    m.version = value
}
