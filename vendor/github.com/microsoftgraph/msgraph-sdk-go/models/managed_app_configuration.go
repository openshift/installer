package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedAppConfiguration 
type ManagedAppConfiguration struct {
    ManagedAppPolicy
    // A set of string key and string value pairs to be sent to apps for users to whom the configuration is scoped, unalterned by this service
    customSettings []KeyValuePairable
}
// NewManagedAppConfiguration instantiates a new ManagedAppConfiguration and sets the default values.
func NewManagedAppConfiguration()(*ManagedAppConfiguration) {
    m := &ManagedAppConfiguration{
        ManagedAppPolicy: *NewManagedAppPolicy(),
    }
    odataTypeValue := "#microsoft.graph.managedAppConfiguration";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateManagedAppConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateManagedAppConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.targetedManagedAppConfiguration":
                        return NewTargetedManagedAppConfiguration(), nil
                }
            }
        }
    }
    return NewManagedAppConfiguration(), nil
}
// GetCustomSettings gets the customSettings property value. A set of string key and string value pairs to be sent to apps for users to whom the configuration is scoped, unalterned by this service
func (m *ManagedAppConfiguration) GetCustomSettings()([]KeyValuePairable) {
    return m.customSettings
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ManagedAppConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ManagedAppPolicy.GetFieldDeserializers()
    res["customSettings"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateKeyValuePairFromDiscriminatorValue , m.SetCustomSettings)
    return res
}
// Serialize serializes information the current object
func (m *ManagedAppConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ManagedAppPolicy.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCustomSettings() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetCustomSettings())
        err = writer.WriteCollectionOfObjectValues("customSettings", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCustomSettings sets the customSettings property value. A set of string key and string value pairs to be sent to apps for users to whom the configuration is scoped, unalterned by this service
func (m *ManagedAppConfiguration) SetCustomSettings(value []KeyValuePairable)() {
    m.customSettings = value
}
