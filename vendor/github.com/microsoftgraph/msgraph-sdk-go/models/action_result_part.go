package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ActionResultPart 
type ActionResultPart struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The error that occurred, if any, during the course of the bulk operation.
    error PublicErrorable
    // The OdataType property
    odataType *string
}
// NewActionResultPart instantiates a new actionResultPart and sets the default values.
func NewActionResultPart()(*ActionResultPart) {
    m := &ActionResultPart{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateActionResultPartFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateActionResultPartFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.aadUserConversationMemberResult":
                        return NewAadUserConversationMemberResult(), nil
                }
            }
        }
    }
    return NewActionResultPart(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ActionResultPart) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetError gets the error property value. The error that occurred, if any, during the course of the bulk operation.
func (m *ActionResultPart) GetError()(PublicErrorable) {
    return m.error
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ActionResultPart) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["error"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePublicErrorFromDiscriminatorValue , m.SetError)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ActionResultPart) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *ActionResultPart) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("error", m.GetError())
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
func (m *ActionResultPart) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetError sets the error property value. The error that occurred, if any, during the course of the bulk operation.
func (m *ActionResultPart) SetError(value PublicErrorable)() {
    m.error = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ActionResultPart) SetOdataType(value *string)() {
    m.odataType = value
}
