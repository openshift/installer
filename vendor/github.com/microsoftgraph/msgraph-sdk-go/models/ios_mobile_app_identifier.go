package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosMobileAppIdentifier 
type IosMobileAppIdentifier struct {
    MobileAppIdentifier
    // The identifier for an app, as specified in the app store.
    bundleId *string
}
// NewIosMobileAppIdentifier instantiates a new IosMobileAppIdentifier and sets the default values.
func NewIosMobileAppIdentifier()(*IosMobileAppIdentifier) {
    m := &IosMobileAppIdentifier{
        MobileAppIdentifier: *NewMobileAppIdentifier(),
    }
    odataTypeValue := "#microsoft.graph.iosMobileAppIdentifier";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateIosMobileAppIdentifierFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateIosMobileAppIdentifierFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosMobileAppIdentifier(), nil
}
// GetBundleId gets the bundleId property value. The identifier for an app, as specified in the app store.
func (m *IosMobileAppIdentifier) GetBundleId()(*string) {
    return m.bundleId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *IosMobileAppIdentifier) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileAppIdentifier.GetFieldDeserializers()
    res["bundleId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetBundleId)
    return res
}
// Serialize serializes information the current object
func (m *IosMobileAppIdentifier) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileAppIdentifier.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("bundleId", m.GetBundleId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBundleId sets the bundleId property value. The identifier for an app, as specified in the app store.
func (m *IosMobileAppIdentifier) SetBundleId(value *string)() {
    m.bundleId = value
}
