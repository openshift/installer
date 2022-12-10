package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EmployeeOrgData 
type EmployeeOrgData struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The cost center associated with the user. Returned only on $select. Supports $filter.
    costCenter *string
    // The name of the division in which the user works. Returned only on $select. Supports $filter.
    division *string
    // The OdataType property
    odataType *string
}
// NewEmployeeOrgData instantiates a new employeeOrgData and sets the default values.
func NewEmployeeOrgData()(*EmployeeOrgData) {
    m := &EmployeeOrgData{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateEmployeeOrgDataFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateEmployeeOrgDataFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEmployeeOrgData(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *EmployeeOrgData) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCostCenter gets the costCenter property value. The cost center associated with the user. Returned only on $select. Supports $filter.
func (m *EmployeeOrgData) GetCostCenter()(*string) {
    return m.costCenter
}
// GetDivision gets the division property value. The name of the division in which the user works. Returned only on $select. Supports $filter.
func (m *EmployeeOrgData) GetDivision()(*string) {
    return m.division
}
// GetFieldDeserializers the deserialization information for the current model
func (m *EmployeeOrgData) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["costCenter"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCostCenter)
    res["division"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDivision)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *EmployeeOrgData) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *EmployeeOrgData) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("costCenter", m.GetCostCenter())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("division", m.GetDivision())
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
func (m *EmployeeOrgData) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCostCenter sets the costCenter property value. The cost center associated with the user. Returned only on $select. Supports $filter.
func (m *EmployeeOrgData) SetCostCenter(value *string)() {
    m.costCenter = value
}
// SetDivision sets the division property value. The name of the division in which the user works. Returned only on $select. Supports $filter.
func (m *EmployeeOrgData) SetDivision(value *string)() {
    m.division = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *EmployeeOrgData) SetOdataType(value *string)() {
    m.odataType = value
}
