package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// BookingQuestionAssignment 
type BookingQuestionAssignment struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The ID of the custom question.
    isRequired *bool
    // The OdataType property
    odataType *string
    // Indicates whether it is mandatory to answer the custom question.
    questionId *string
}
// NewBookingQuestionAssignment instantiates a new bookingQuestionAssignment and sets the default values.
func NewBookingQuestionAssignment()(*BookingQuestionAssignment) {
    m := &BookingQuestionAssignment{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateBookingQuestionAssignmentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateBookingQuestionAssignmentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBookingQuestionAssignment(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *BookingQuestionAssignment) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *BookingQuestionAssignment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["isRequired"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetIsRequired)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["questionId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetQuestionId)
    return res
}
// GetIsRequired gets the isRequired property value. The ID of the custom question.
func (m *BookingQuestionAssignment) GetIsRequired()(*bool) {
    return m.isRequired
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *BookingQuestionAssignment) GetOdataType()(*string) {
    return m.odataType
}
// GetQuestionId gets the questionId property value. Indicates whether it is mandatory to answer the custom question.
func (m *BookingQuestionAssignment) GetQuestionId()(*string) {
    return m.questionId
}
// Serialize serializes information the current object
func (m *BookingQuestionAssignment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("isRequired", m.GetIsRequired())
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
        err := writer.WriteStringValue("questionId", m.GetQuestionId())
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
func (m *BookingQuestionAssignment) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetIsRequired sets the isRequired property value. The ID of the custom question.
func (m *BookingQuestionAssignment) SetIsRequired(value *bool)() {
    m.isRequired = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *BookingQuestionAssignment) SetOdataType(value *string)() {
    m.odataType = value
}
// SetQuestionId sets the questionId property value. Indicates whether it is mandatory to answer the custom question.
func (m *BookingQuestionAssignment) SetQuestionId(value *string)() {
    m.questionId = value
}
