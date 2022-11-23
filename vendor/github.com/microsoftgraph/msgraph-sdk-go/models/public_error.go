package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PublicError 
type PublicError struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Represents the error code.
    code *string
    // Details of the error.
    details []PublicErrorDetailable
    // Details of the inner error.
    innerError PublicInnerErrorable
    // A non-localized message for the developer.
    message *string
    // The OdataType property
    odataType *string
    // The target of the error.
    target *string
}
// NewPublicError instantiates a new publicError and sets the default values.
func NewPublicError()(*PublicError) {
    m := &PublicError{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreatePublicErrorFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePublicErrorFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPublicError(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *PublicError) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCode gets the code property value. Represents the error code.
func (m *PublicError) GetCode()(*string) {
    return m.code
}
// GetDetails gets the details property value. Details of the error.
func (m *PublicError) GetDetails()([]PublicErrorDetailable) {
    return m.details
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PublicError) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["code"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCode)
    res["details"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreatePublicErrorDetailFromDiscriminatorValue , m.SetDetails)
    res["innerError"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreatePublicInnerErrorFromDiscriminatorValue , m.SetInnerError)
    res["message"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetMessage)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["target"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetTarget)
    return res
}
// GetInnerError gets the innerError property value. Details of the inner error.
func (m *PublicError) GetInnerError()(PublicInnerErrorable) {
    return m.innerError
}
// GetMessage gets the message property value. A non-localized message for the developer.
func (m *PublicError) GetMessage()(*string) {
    return m.message
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *PublicError) GetOdataType()(*string) {
    return m.odataType
}
// GetTarget gets the target property value. The target of the error.
func (m *PublicError) GetTarget()(*string) {
    return m.target
}
// Serialize serializes information the current object
func (m *PublicError) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("code", m.GetCode())
        if err != nil {
            return err
        }
    }
    if m.GetDetails() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetDetails())
        err := writer.WriteCollectionOfObjectValues("details", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("innerError", m.GetInnerError())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("message", m.GetMessage())
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
        err := writer.WriteStringValue("target", m.GetTarget())
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
func (m *PublicError) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCode sets the code property value. Represents the error code.
func (m *PublicError) SetCode(value *string)() {
    m.code = value
}
// SetDetails sets the details property value. Details of the error.
func (m *PublicError) SetDetails(value []PublicErrorDetailable)() {
    m.details = value
}
// SetInnerError sets the innerError property value. Details of the inner error.
func (m *PublicError) SetInnerError(value PublicInnerErrorable)() {
    m.innerError = value
}
// SetMessage sets the message property value. A non-localized message for the developer.
func (m *PublicError) SetMessage(value *string)() {
    m.message = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *PublicError) SetOdataType(value *string)() {
    m.odataType = value
}
// SetTarget sets the target property value. The target of the error.
func (m *PublicError) SetTarget(value *string)() {
    m.target = value
}
