package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// OpenTypeExtension 
type OpenTypeExtension struct {
    Extension
    // A unique text identifier for an open type data extension. Required.
    extensionName *string
}
// NewOpenTypeExtension instantiates a new OpenTypeExtension and sets the default values.
func NewOpenTypeExtension()(*OpenTypeExtension) {
    m := &OpenTypeExtension{
        Extension: *NewExtension(),
    }
    odataTypeValue := "#microsoft.graph.openTypeExtension";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateOpenTypeExtensionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateOpenTypeExtensionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOpenTypeExtension(), nil
}
// GetExtensionName gets the extensionName property value. A unique text identifier for an open type data extension. Required.
func (m *OpenTypeExtension) GetExtensionName()(*string) {
    return m.extensionName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *OpenTypeExtension) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Extension.GetFieldDeserializers()
    res["extensionName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetExtensionName)
    return res
}
// Serialize serializes information the current object
func (m *OpenTypeExtension) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Extension.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("extensionName", m.GetExtensionName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetExtensionName sets the extensionName property value. A unique text identifier for an open type data extension. Required.
func (m *OpenTypeExtension) SetExtensionName(value *string)() {
    m.extensionName = value
}
