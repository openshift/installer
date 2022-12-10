package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// OnPremisesProvisioningError 
type OnPremisesProvisioningError struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Category of the provisioning error. Note: Currently, there is only one possible value. Possible value: PropertyConflict - indicates a property value is not unique. Other objects contain the same value for the property.
    category *string
    // The date and time at which the error occurred.
    occurredDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The OdataType property
    odataType *string
    // Name of the directory property causing the error. Current possible values: UserPrincipalName or ProxyAddress
    propertyCausingError *string
    // Value of the property causing the error.
    value *string
}
// NewOnPremisesProvisioningError instantiates a new onPremisesProvisioningError and sets the default values.
func NewOnPremisesProvisioningError()(*OnPremisesProvisioningError) {
    m := &OnPremisesProvisioningError{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateOnPremisesProvisioningErrorFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateOnPremisesProvisioningErrorFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnPremisesProvisioningError(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *OnPremisesProvisioningError) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCategory gets the category property value. Category of the provisioning error. Note: Currently, there is only one possible value. Possible value: PropertyConflict - indicates a property value is not unique. Other objects contain the same value for the property.
func (m *OnPremisesProvisioningError) GetCategory()(*string) {
    return m.category
}
// GetFieldDeserializers the deserialization information for the current model
func (m *OnPremisesProvisioningError) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["category"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCategory)
    res["occurredDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetOccurredDateTime)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["propertyCausingError"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPropertyCausingError)
    res["value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetValue)
    return res
}
// GetOccurredDateTime gets the occurredDateTime property value. The date and time at which the error occurred.
func (m *OnPremisesProvisioningError) GetOccurredDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.occurredDateTime
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *OnPremisesProvisioningError) GetOdataType()(*string) {
    return m.odataType
}
// GetPropertyCausingError gets the propertyCausingError property value. Name of the directory property causing the error. Current possible values: UserPrincipalName or ProxyAddress
func (m *OnPremisesProvisioningError) GetPropertyCausingError()(*string) {
    return m.propertyCausingError
}
// GetValue gets the value property value. Value of the property causing the error.
func (m *OnPremisesProvisioningError) GetValue()(*string) {
    return m.value
}
// Serialize serializes information the current object
func (m *OnPremisesProvisioningError) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("category", m.GetCategory())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("occurredDateTime", m.GetOccurredDateTime())
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
        err := writer.WriteStringValue("propertyCausingError", m.GetPropertyCausingError())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("value", m.GetValue())
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
func (m *OnPremisesProvisioningError) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCategory sets the category property value. Category of the provisioning error. Note: Currently, there is only one possible value. Possible value: PropertyConflict - indicates a property value is not unique. Other objects contain the same value for the property.
func (m *OnPremisesProvisioningError) SetCategory(value *string)() {
    m.category = value
}
// SetOccurredDateTime sets the occurredDateTime property value. The date and time at which the error occurred.
func (m *OnPremisesProvisioningError) SetOccurredDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.occurredDateTime = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *OnPremisesProvisioningError) SetOdataType(value *string)() {
    m.odataType = value
}
// SetPropertyCausingError sets the propertyCausingError property value. Name of the directory property causing the error. Current possible values: UserPrincipalName or ProxyAddress
func (m *OnPremisesProvisioningError) SetPropertyCausingError(value *string)() {
    m.propertyCausingError = value
}
// SetValue sets the value property value. Value of the property causing the error.
func (m *OnPremisesProvisioningError) SetValue(value *string)() {
    m.value = value
}
