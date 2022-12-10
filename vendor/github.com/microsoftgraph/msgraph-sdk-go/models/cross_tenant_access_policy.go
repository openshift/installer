package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CrossTenantAccessPolicy 
type CrossTenantAccessPolicy struct {
    PolicyBase
    // The allowedCloudEndpoints property
    allowedCloudEndpoints []string
    // Defines the default configuration for how your organization interacts with external Azure Active Directory organizations.
    default_escaped CrossTenantAccessPolicyConfigurationDefaultable
    // Defines partner-specific configurations for external Azure Active Directory organizations.
    partners []CrossTenantAccessPolicyConfigurationPartnerable
}
// NewCrossTenantAccessPolicy instantiates a new CrossTenantAccessPolicy and sets the default values.
func NewCrossTenantAccessPolicy()(*CrossTenantAccessPolicy) {
    m := &CrossTenantAccessPolicy{
        PolicyBase: *NewPolicyBase(),
    }
    odataTypeValue := "#microsoft.graph.crossTenantAccessPolicy";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateCrossTenantAccessPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCrossTenantAccessPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCrossTenantAccessPolicy(), nil
}
// GetAllowedCloudEndpoints gets the allowedCloudEndpoints property value. The allowedCloudEndpoints property
func (m *CrossTenantAccessPolicy) GetAllowedCloudEndpoints()([]string) {
    return m.allowedCloudEndpoints
}
// GetDefault gets the default property value. Defines the default configuration for how your organization interacts with external Azure Active Directory organizations.
func (m *CrossTenantAccessPolicy) GetDefault()(CrossTenantAccessPolicyConfigurationDefaultable) {
    return m.default_escaped
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CrossTenantAccessPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.PolicyBase.GetFieldDeserializers()
    res["allowedCloudEndpoints"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetAllowedCloudEndpoints)
    res["default"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateCrossTenantAccessPolicyConfigurationDefaultFromDiscriminatorValue , m.SetDefault)
    res["partners"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateCrossTenantAccessPolicyConfigurationPartnerFromDiscriminatorValue , m.SetPartners)
    return res
}
// GetPartners gets the partners property value. Defines partner-specific configurations for external Azure Active Directory organizations.
func (m *CrossTenantAccessPolicy) GetPartners()([]CrossTenantAccessPolicyConfigurationPartnerable) {
    return m.partners
}
// Serialize serializes information the current object
func (m *CrossTenantAccessPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.PolicyBase.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAllowedCloudEndpoints() != nil {
        err = writer.WriteCollectionOfStringValues("allowedCloudEndpoints", m.GetAllowedCloudEndpoints())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("default", m.GetDefault())
        if err != nil {
            return err
        }
    }
    if m.GetPartners() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPartners())
        err = writer.WriteCollectionOfObjectValues("partners", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowedCloudEndpoints sets the allowedCloudEndpoints property value. The allowedCloudEndpoints property
func (m *CrossTenantAccessPolicy) SetAllowedCloudEndpoints(value []string)() {
    m.allowedCloudEndpoints = value
}
// SetDefault sets the default property value. Defines the default configuration for how your organization interacts with external Azure Active Directory organizations.
func (m *CrossTenantAccessPolicy) SetDefault(value CrossTenantAccessPolicyConfigurationDefaultable)() {
    m.default_escaped = value
}
// SetPartners sets the partners property value. Defines partner-specific configurations for external Azure Active Directory organizations.
func (m *CrossTenantAccessPolicy) SetPartners(value []CrossTenantAccessPolicyConfigurationPartnerable)() {
    m.partners = value
}
