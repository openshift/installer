package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ProvisioningErrorInfo 
type ProvisioningErrorInfo struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Additional details in case of error.
    additionalDetails *string
    // Categorizes the error code. Possible values are failure, nonServiceFailure, success, unknownFutureValue
    errorCategory *ProvisioningStatusErrorCategory
    // Unique error code if any occurred. Learn more
    errorCode *string
    // The OdataType property
    odataType *string
    // Summarizes the status and describes why the status happened.
    reason *string
    // Provides the resolution for the corresponding error.
    recommendedAction *string
}
// NewProvisioningErrorInfo instantiates a new provisioningErrorInfo and sets the default values.
func NewProvisioningErrorInfo()(*ProvisioningErrorInfo) {
    m := &ProvisioningErrorInfo{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateProvisioningErrorInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateProvisioningErrorInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewProvisioningErrorInfo(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ProvisioningErrorInfo) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetAdditionalDetails gets the additionalDetails property value. Additional details in case of error.
func (m *ProvisioningErrorInfo) GetAdditionalDetails()(*string) {
    return m.additionalDetails
}
// GetErrorCategory gets the errorCategory property value. Categorizes the error code. Possible values are failure, nonServiceFailure, success, unknownFutureValue
func (m *ProvisioningErrorInfo) GetErrorCategory()(*ProvisioningStatusErrorCategory) {
    return m.errorCategory
}
// GetErrorCode gets the errorCode property value. Unique error code if any occurred. Learn more
func (m *ProvisioningErrorInfo) GetErrorCode()(*string) {
    return m.errorCode
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ProvisioningErrorInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["additionalDetails"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAdditionalDetails)
    res["errorCategory"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseProvisioningStatusErrorCategory , m.SetErrorCategory)
    res["errorCode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetErrorCode)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["reason"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetReason)
    res["recommendedAction"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetRecommendedAction)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ProvisioningErrorInfo) GetOdataType()(*string) {
    return m.odataType
}
// GetReason gets the reason property value. Summarizes the status and describes why the status happened.
func (m *ProvisioningErrorInfo) GetReason()(*string) {
    return m.reason
}
// GetRecommendedAction gets the recommendedAction property value. Provides the resolution for the corresponding error.
func (m *ProvisioningErrorInfo) GetRecommendedAction()(*string) {
    return m.recommendedAction
}
// Serialize serializes information the current object
func (m *ProvisioningErrorInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("additionalDetails", m.GetAdditionalDetails())
        if err != nil {
            return err
        }
    }
    if m.GetErrorCategory() != nil {
        cast := (*m.GetErrorCategory()).String()
        err := writer.WriteStringValue("errorCategory", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("errorCode", m.GetErrorCode())
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
        err := writer.WriteStringValue("reason", m.GetReason())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("recommendedAction", m.GetRecommendedAction())
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
func (m *ProvisioningErrorInfo) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetAdditionalDetails sets the additionalDetails property value. Additional details in case of error.
func (m *ProvisioningErrorInfo) SetAdditionalDetails(value *string)() {
    m.additionalDetails = value
}
// SetErrorCategory sets the errorCategory property value. Categorizes the error code. Possible values are failure, nonServiceFailure, success, unknownFutureValue
func (m *ProvisioningErrorInfo) SetErrorCategory(value *ProvisioningStatusErrorCategory)() {
    m.errorCategory = value
}
// SetErrorCode sets the errorCode property value. Unique error code if any occurred. Learn more
func (m *ProvisioningErrorInfo) SetErrorCode(value *string)() {
    m.errorCode = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ProvisioningErrorInfo) SetOdataType(value *string)() {
    m.odataType = value
}
// SetReason sets the reason property value. Summarizes the status and describes why the status happened.
func (m *ProvisioningErrorInfo) SetReason(value *string)() {
    m.reason = value
}
// SetRecommendedAction sets the recommendedAction property value. Provides the resolution for the corresponding error.
func (m *ProvisioningErrorInfo) SetRecommendedAction(value *string)() {
    m.recommendedAction = value
}
