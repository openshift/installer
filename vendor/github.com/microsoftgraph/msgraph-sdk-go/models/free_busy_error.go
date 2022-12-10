package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// FreeBusyError 
type FreeBusyError struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Describes the error.
    message *string
    // The OdataType property
    odataType *string
    // The response code from querying for the availability of the user, distribution list, or resource.
    responseCode *string
}
// NewFreeBusyError instantiates a new freeBusyError and sets the default values.
func NewFreeBusyError()(*FreeBusyError) {
    m := &FreeBusyError{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateFreeBusyErrorFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateFreeBusyErrorFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFreeBusyError(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *FreeBusyError) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *FreeBusyError) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["message"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetMessage)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["responseCode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetResponseCode)
    return res
}
// GetMessage gets the message property value. Describes the error.
func (m *FreeBusyError) GetMessage()(*string) {
    return m.message
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *FreeBusyError) GetOdataType()(*string) {
    return m.odataType
}
// GetResponseCode gets the responseCode property value. The response code from querying for the availability of the user, distribution list, or resource.
func (m *FreeBusyError) GetResponseCode()(*string) {
    return m.responseCode
}
// Serialize serializes information the current object
func (m *FreeBusyError) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err := writer.WriteStringValue("responseCode", m.GetResponseCode())
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
func (m *FreeBusyError) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetMessage sets the message property value. Describes the error.
func (m *FreeBusyError) SetMessage(value *string)() {
    m.message = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *FreeBusyError) SetOdataType(value *string)() {
    m.odataType = value
}
// SetResponseCode sets the responseCode property value. The response code from querying for the availability of the user, distribution list, or resource.
func (m *FreeBusyError) SetResponseCode(value *string)() {
    m.responseCode = value
}
