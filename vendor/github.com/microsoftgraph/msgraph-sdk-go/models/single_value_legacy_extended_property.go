package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SingleValueLegacyExtendedProperty provides operations to manage the collection of agreement entities.
type SingleValueLegacyExtendedProperty struct {
    Entity
    // A property value.
    value *string
}
// NewSingleValueLegacyExtendedProperty instantiates a new singleValueLegacyExtendedProperty and sets the default values.
func NewSingleValueLegacyExtendedProperty()(*SingleValueLegacyExtendedProperty) {
    m := &SingleValueLegacyExtendedProperty{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSingleValueLegacyExtendedPropertyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSingleValueLegacyExtendedPropertyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSingleValueLegacyExtendedProperty(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SingleValueLegacyExtendedProperty) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetValue)
    return res
}
// GetValue gets the value property value. A property value.
func (m *SingleValueLegacyExtendedProperty) GetValue()(*string) {
    return m.value
}
// Serialize serializes information the current object
func (m *SingleValueLegacyExtendedProperty) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("value", m.GetValue())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetValue sets the value property value. A property value.
func (m *SingleValueLegacyExtendedProperty) SetValue(value *string)() {
    m.value = value
}
