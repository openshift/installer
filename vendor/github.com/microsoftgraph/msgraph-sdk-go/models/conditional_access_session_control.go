package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ConditionalAccessSessionControl 
type ConditionalAccessSessionControl struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Specifies whether the session control is enabled.
    isEnabled *bool
    // The OdataType property
    odataType *string
}
// NewConditionalAccessSessionControl instantiates a new conditionalAccessSessionControl and sets the default values.
func NewConditionalAccessSessionControl()(*ConditionalAccessSessionControl) {
    m := &ConditionalAccessSessionControl{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateConditionalAccessSessionControlFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateConditionalAccessSessionControlFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.applicationEnforcedRestrictionsSessionControl":
                        return NewApplicationEnforcedRestrictionsSessionControl(), nil
                    case "#microsoft.graph.cloudAppSecuritySessionControl":
                        return NewCloudAppSecuritySessionControl(), nil
                    case "#microsoft.graph.persistentBrowserSessionControl":
                        return NewPersistentBrowserSessionControl(), nil
                    case "#microsoft.graph.signInFrequencySessionControl":
                        return NewSignInFrequencySessionControl(), nil
                }
            }
        }
    }
    return NewConditionalAccessSessionControl(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ConditionalAccessSessionControl) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ConditionalAccessSessionControl) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["isEnabled"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsEnabled)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetIsEnabled gets the isEnabled property value. Specifies whether the session control is enabled.
func (m *ConditionalAccessSessionControl) GetIsEnabled()(*bool) {
    return m.isEnabled
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ConditionalAccessSessionControl) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *ConditionalAccessSessionControl) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("isEnabled", m.GetIsEnabled())
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
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ConditionalAccessSessionControl) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetIsEnabled sets the isEnabled property value. Specifies whether the session control is enabled.
func (m *ConditionalAccessSessionControl) SetIsEnabled(value *bool)() {
    m.isEnabled = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ConditionalAccessSessionControl) SetOdataType(value *string)() {
    m.odataType = value
}
