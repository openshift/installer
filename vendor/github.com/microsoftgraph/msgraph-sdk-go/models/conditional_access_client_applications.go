package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ConditionalAccessClientApplications 
type ConditionalAccessClientApplications struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Service principal IDs excluded from the policy scope.
    excludeServicePrincipals []string
    // Service principal IDs included in the policy scope, or ServicePrincipalsInMyTenant.
    includeServicePrincipals []string
    // The OdataType property
    odataType *string
}
// NewConditionalAccessClientApplications instantiates a new conditionalAccessClientApplications and sets the default values.
func NewConditionalAccessClientApplications()(*ConditionalAccessClientApplications) {
    m := &ConditionalAccessClientApplications{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateConditionalAccessClientApplicationsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateConditionalAccessClientApplicationsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConditionalAccessClientApplications(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ConditionalAccessClientApplications) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetExcludeServicePrincipals gets the excludeServicePrincipals property value. Service principal IDs excluded from the policy scope.
func (m *ConditionalAccessClientApplications) GetExcludeServicePrincipals()([]string) {
    return m.excludeServicePrincipals
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ConditionalAccessClientApplications) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["excludeServicePrincipals"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetExcludeServicePrincipals)
    res["includeServicePrincipals"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetIncludeServicePrincipals)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetIncludeServicePrincipals gets the includeServicePrincipals property value. Service principal IDs included in the policy scope, or ServicePrincipalsInMyTenant.
func (m *ConditionalAccessClientApplications) GetIncludeServicePrincipals()([]string) {
    return m.includeServicePrincipals
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ConditionalAccessClientApplications) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *ConditionalAccessClientApplications) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetExcludeServicePrincipals() != nil {
        err := writer.WriteCollectionOfStringValues("excludeServicePrincipals", m.GetExcludeServicePrincipals())
        if err != nil {
            return err
        }
    }
    if m.GetIncludeServicePrincipals() != nil {
        err := writer.WriteCollectionOfStringValues("includeServicePrincipals", m.GetIncludeServicePrincipals())
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
func (m *ConditionalAccessClientApplications) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetExcludeServicePrincipals sets the excludeServicePrincipals property value. Service principal IDs excluded from the policy scope.
func (m *ConditionalAccessClientApplications) SetExcludeServicePrincipals(value []string)() {
    m.excludeServicePrincipals = value
}
// SetIncludeServicePrincipals sets the includeServicePrincipals property value. Service principal IDs included in the policy scope, or ServicePrincipalsInMyTenant.
func (m *ConditionalAccessClientApplications) SetIncludeServicePrincipals(value []string)() {
    m.includeServicePrincipals = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ConditionalAccessClientApplications) SetOdataType(value *string)() {
    m.odataType = value
}
