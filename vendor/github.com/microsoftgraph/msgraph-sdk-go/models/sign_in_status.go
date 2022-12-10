package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SignInStatus 
type SignInStatus struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Provides additional details on the sign-in activity
    additionalDetails *string
    // Provides the 5-6 digit error code that's generated during a sign-in failure. Check out the list of error codes and messages.
    errorCode *int32
    // Provides the error message or the reason for failure for the corresponding sign-in activity. Check out the list of error codes and messages.
    failureReason *string
    // The OdataType property
    odataType *string
}
// NewSignInStatus instantiates a new signInStatus and sets the default values.
func NewSignInStatus()(*SignInStatus) {
    m := &SignInStatus{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateSignInStatusFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSignInStatusFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSignInStatus(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SignInStatus) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAdditionalDetails gets the additionalDetails property value. Provides additional details on the sign-in activity
func (m *SignInStatus) GetAdditionalDetails()(*string) {
    return m.additionalDetails
}
// GetErrorCode gets the errorCode property value. Provides the 5-6 digit error code that's generated during a sign-in failure. Check out the list of error codes and messages.
func (m *SignInStatus) GetErrorCode()(*int32) {
    return m.errorCode
}
// GetFailureReason gets the failureReason property value. Provides the error message or the reason for failure for the corresponding sign-in activity. Check out the list of error codes and messages.
func (m *SignInStatus) GetFailureReason()(*string) {
    return m.failureReason
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SignInStatus) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["additionalDetails"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAdditionalDetails)
    res["errorCode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetInt32Value(m.SetErrorCode)
    res["failureReason"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetFailureReason)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *SignInStatus) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *SignInStatus) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("additionalDetails", m.GetAdditionalDetails())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("errorCode", m.GetErrorCode())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("failureReason", m.GetFailureReason())
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
func (m *SignInStatus) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAdditionalDetails sets the additionalDetails property value. Provides additional details on the sign-in activity
func (m *SignInStatus) SetAdditionalDetails(value *string)() {
    m.additionalDetails = value
}
// SetErrorCode sets the errorCode property value. Provides the 5-6 digit error code that's generated during a sign-in failure. Check out the list of error codes and messages.
func (m *SignInStatus) SetErrorCode(value *int32)() {
    m.errorCode = value
}
// SetFailureReason sets the failureReason property value. Provides the error message or the reason for failure for the corresponding sign-in activity. Check out the list of error codes and messages.
func (m *SignInStatus) SetFailureReason(value *string)() {
    m.failureReason = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SignInStatus) SetOdataType(value *string)() {
    m.odataType = value
}
