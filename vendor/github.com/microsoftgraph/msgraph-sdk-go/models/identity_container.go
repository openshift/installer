package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IdentityContainer 
type IdentityContainer struct {
    Entity
    // Represents entry point for API connectors.
    apiConnectors []IdentityApiConnectorable
    // Represents entry point for B2X/self-service sign-up identity userflows.
    b2xUserFlows []B2xIdentityUserFlowable
    // the entry point for the Conditional Access (CA) object model.
    conditionalAccess ConditionalAccessRootable
    // The identityProviders property
    identityProviders []IdentityProviderBaseable
    // Represents entry point for identity userflow attributes.
    userFlowAttributes []IdentityUserFlowAttributeable
}
// NewIdentityContainer instantiates a new IdentityContainer and sets the default values.
func NewIdentityContainer()(*IdentityContainer) {
    m := &IdentityContainer{
        Entity: *NewEntity(),
    }
    return m
}
// CreateIdentityContainerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateIdentityContainerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIdentityContainer(), nil
}
// GetApiConnectors gets the apiConnectors property value. Represents entry point for API connectors.
func (m *IdentityContainer) GetApiConnectors()([]IdentityApiConnectorable) {
    return m.apiConnectors
}
// GetB2xUserFlows gets the b2xUserFlows property value. Represents entry point for B2X/self-service sign-up identity userflows.
func (m *IdentityContainer) GetB2xUserFlows()([]B2xIdentityUserFlowable) {
    return m.b2xUserFlows
}
// GetConditionalAccess gets the conditionalAccess property value. the entry point for the Conditional Access (CA) object model.
func (m *IdentityContainer) GetConditionalAccess()(ConditionalAccessRootable) {
    return m.conditionalAccess
}
// GetFieldDeserializers the deserialization information for the current model
func (m *IdentityContainer) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["apiConnectors"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateIdentityApiConnectorFromDiscriminatorValue , m.SetApiConnectors)
    res["b2xUserFlows"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateB2xIdentityUserFlowFromDiscriminatorValue , m.SetB2xUserFlows)
    res["conditionalAccess"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateConditionalAccessRootFromDiscriminatorValue , m.SetConditionalAccess)
    res["identityProviders"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateIdentityProviderBaseFromDiscriminatorValue , m.SetIdentityProviders)
    res["userFlowAttributes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateIdentityUserFlowAttributeFromDiscriminatorValue , m.SetUserFlowAttributes)
    return res
}
// GetIdentityProviders gets the identityProviders property value. The identityProviders property
func (m *IdentityContainer) GetIdentityProviders()([]IdentityProviderBaseable) {
    return m.identityProviders
}
// GetUserFlowAttributes gets the userFlowAttributes property value. Represents entry point for identity userflow attributes.
func (m *IdentityContainer) GetUserFlowAttributes()([]IdentityUserFlowAttributeable) {
    return m.userFlowAttributes
}
// Serialize serializes information the current object
func (m *IdentityContainer) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetApiConnectors() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetApiConnectors())
        err = writer.WriteCollectionOfObjectValues("apiConnectors", cast)
        if err != nil {
            return err
        }
    }
    if m.GetB2xUserFlows() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetB2xUserFlows())
        err = writer.WriteCollectionOfObjectValues("b2xUserFlows", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("conditionalAccess", m.GetConditionalAccess())
        if err != nil {
            return err
        }
    }
    if m.GetIdentityProviders() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetIdentityProviders())
        err = writer.WriteCollectionOfObjectValues("identityProviders", cast)
        if err != nil {
            return err
        }
    }
    if m.GetUserFlowAttributes() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetUserFlowAttributes())
        err = writer.WriteCollectionOfObjectValues("userFlowAttributes", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApiConnectors sets the apiConnectors property value. Represents entry point for API connectors.
func (m *IdentityContainer) SetApiConnectors(value []IdentityApiConnectorable)() {
    m.apiConnectors = value
}
// SetB2xUserFlows sets the b2xUserFlows property value. Represents entry point for B2X/self-service sign-up identity userflows.
func (m *IdentityContainer) SetB2xUserFlows(value []B2xIdentityUserFlowable)() {
    m.b2xUserFlows = value
}
// SetConditionalAccess sets the conditionalAccess property value. the entry point for the Conditional Access (CA) object model.
func (m *IdentityContainer) SetConditionalAccess(value ConditionalAccessRootable)() {
    m.conditionalAccess = value
}
// SetIdentityProviders sets the identityProviders property value. The identityProviders property
func (m *IdentityContainer) SetIdentityProviders(value []IdentityProviderBaseable)() {
    m.identityProviders = value
}
// SetUserFlowAttributes sets the userFlowAttributes property value. Represents entry point for identity userflow attributes.
func (m *IdentityContainer) SetUserFlowAttributes(value []IdentityUserFlowAttributeable)() {
    m.userFlowAttributes = value
}
