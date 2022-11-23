package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RiskUserActivity 
type RiskUserActivity struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Details of the detected risk. Possible values are: none, adminGeneratedTemporaryPassword, userPerformedSecuredPasswordChange, userPerformedSecuredPasswordReset, adminConfirmedSigninSafe, aiConfirmedSigninSafe, userPassedMFADrivenByRiskBasedPolicy, adminDismissedAllRiskForUser, adminConfirmedSigninCompromised, hidden, adminConfirmedUserCompromised, unknownFutureValue.
    detail *RiskDetail
    // The OdataType property
    odataType *string
    // The type of risk event detected.
    riskEventTypes []string
}
// NewRiskUserActivity instantiates a new riskUserActivity and sets the default values.
func NewRiskUserActivity()(*RiskUserActivity) {
    m := &RiskUserActivity{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateRiskUserActivityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateRiskUserActivityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRiskUserActivity(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *RiskUserActivity) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetDetail gets the detail property value. Details of the detected risk. Possible values are: none, adminGeneratedTemporaryPassword, userPerformedSecuredPasswordChange, userPerformedSecuredPasswordReset, adminConfirmedSigninSafe, aiConfirmedSigninSafe, userPassedMFADrivenByRiskBasedPolicy, adminDismissedAllRiskForUser, adminConfirmedSigninCompromised, hidden, adminConfirmedUserCompromised, unknownFutureValue.
func (m *RiskUserActivity) GetDetail()(*RiskDetail) {
    return m.detail
}
// GetFieldDeserializers the deserialization information for the current model
func (m *RiskUserActivity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["detail"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseRiskDetail , m.SetDetail)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["riskEventTypes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetRiskEventTypes)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *RiskUserActivity) GetOdataType()(*string) {
    return m.odataType
}
// GetRiskEventTypes gets the riskEventTypes property value. The type of risk event detected.
func (m *RiskUserActivity) GetRiskEventTypes()([]string) {
    return m.riskEventTypes
}
// Serialize serializes information the current object
func (m *RiskUserActivity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetDetail() != nil {
        cast := (*m.GetDetail()).String()
        err := writer.WriteStringValue("detail", &cast)
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
    if m.GetRiskEventTypes() != nil {
        err := writer.WriteCollectionOfStringValues("riskEventTypes", m.GetRiskEventTypes())
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
func (m *RiskUserActivity) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetDetail sets the detail property value. Details of the detected risk. Possible values are: none, adminGeneratedTemporaryPassword, userPerformedSecuredPasswordChange, userPerformedSecuredPasswordReset, adminConfirmedSigninSafe, aiConfirmedSigninSafe, userPassedMFADrivenByRiskBasedPolicy, adminDismissedAllRiskForUser, adminConfirmedSigninCompromised, hidden, adminConfirmedUserCompromised, unknownFutureValue.
func (m *RiskUserActivity) SetDetail(value *RiskDetail)() {
    m.detail = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *RiskUserActivity) SetOdataType(value *string)() {
    m.odataType = value
}
// SetRiskEventTypes sets the riskEventTypes property value. The type of risk event detected.
func (m *RiskUserActivity) SetRiskEventTypes(value []string)() {
    m.riskEventTypes = value
}
