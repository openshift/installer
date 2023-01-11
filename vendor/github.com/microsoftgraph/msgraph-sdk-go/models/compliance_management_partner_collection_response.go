package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ComplianceManagementPartnerCollectionResponse provides operations to manage the complianceManagementPartners property of the microsoft.graph.deviceManagement entity.
type ComplianceManagementPartnerCollectionResponse struct {
    BaseCollectionPaginationCountResponse
    // The value property
    value []ComplianceManagementPartnerable
}
// NewComplianceManagementPartnerCollectionResponse instantiates a new ComplianceManagementPartnerCollectionResponse and sets the default values.
func NewComplianceManagementPartnerCollectionResponse()(*ComplianceManagementPartnerCollectionResponse) {
    m := &ComplianceManagementPartnerCollectionResponse{
        BaseCollectionPaginationCountResponse: *NewBaseCollectionPaginationCountResponse(),
    }
    return m
}
// CreateComplianceManagementPartnerCollectionResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateComplianceManagementPartnerCollectionResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewComplianceManagementPartnerCollectionResponse(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ComplianceManagementPartnerCollectionResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseCollectionPaginationCountResponse.GetFieldDeserializers()
    res["value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateComplianceManagementPartnerFromDiscriminatorValue , m.SetValue)
    return res
}
// GetValue gets the value property value. The value property
func (m *ComplianceManagementPartnerCollectionResponse) GetValue()([]ComplianceManagementPartnerable) {
    return m.value
}
// Serialize serializes information the current object
func (m *ComplianceManagementPartnerCollectionResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseCollectionPaginationCountResponse.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetValue() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetValue())
        err = writer.WriteCollectionOfObjectValues("value", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetValue sets the value property value. The value property
func (m *ComplianceManagementPartnerCollectionResponse) SetValue(value []ComplianceManagementPartnerable)() {
    m.value = value
}
