package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ConditionalAccessRoot 
type ConditionalAccessRoot struct {
    Entity
}
// NewConditionalAccessRoot instantiates a new conditionalAccessRoot and sets the default values.
func NewConditionalAccessRoot()(*ConditionalAccessRoot) {
    m := &ConditionalAccessRoot{
        Entity: *NewEntity(),
    }
    return m
}
// CreateConditionalAccessRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateConditionalAccessRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConditionalAccessRoot(), nil
}
// GetAuthenticationContextClassReferences gets the authenticationContextClassReferences property value. Read-only. Nullable. Returns a collection of the specified authentication context class references.
func (m *ConditionalAccessRoot) GetAuthenticationContextClassReferences()([]AuthenticationContextClassReferenceable) {
    val, err := m.GetBackingStore().Get("authenticationContextClassReferences")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AuthenticationContextClassReferenceable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ConditionalAccessRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["authenticationContextClassReferences"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAuthenticationContextClassReferenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AuthenticationContextClassReferenceable, len(val))
            for i, v := range val {
                res[i] = v.(AuthenticationContextClassReferenceable)
            }
            m.SetAuthenticationContextClassReferences(res)
        }
        return nil
    }
    res["namedLocations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateNamedLocationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]NamedLocationable, len(val))
            for i, v := range val {
                res[i] = v.(NamedLocationable)
            }
            m.SetNamedLocations(res)
        }
        return nil
    }
    res["policies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateConditionalAccessPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ConditionalAccessPolicyable, len(val))
            for i, v := range val {
                res[i] = v.(ConditionalAccessPolicyable)
            }
            m.SetPolicies(res)
        }
        return nil
    }
    res["templates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateConditionalAccessTemplateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ConditionalAccessTemplateable, len(val))
            for i, v := range val {
                res[i] = v.(ConditionalAccessTemplateable)
            }
            m.SetTemplates(res)
        }
        return nil
    }
    return res
}
// GetNamedLocations gets the namedLocations property value. Read-only. Nullable. Returns a collection of the specified named locations.
func (m *ConditionalAccessRoot) GetNamedLocations()([]NamedLocationable) {
    val, err := m.GetBackingStore().Get("namedLocations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]NamedLocationable)
    }
    return nil
}
// GetPolicies gets the policies property value. Read-only. Nullable. Returns a collection of the specified Conditional Access (CA) policies.
func (m *ConditionalAccessRoot) GetPolicies()([]ConditionalAccessPolicyable) {
    val, err := m.GetBackingStore().Get("policies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ConditionalAccessPolicyable)
    }
    return nil
}
// GetTemplates gets the templates property value. Read-only. Nullable. Returns a collection of the specified Conditional Access templates.
func (m *ConditionalAccessRoot) GetTemplates()([]ConditionalAccessTemplateable) {
    val, err := m.GetBackingStore().Get("templates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ConditionalAccessTemplateable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ConditionalAccessRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAuthenticationContextClassReferences() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAuthenticationContextClassReferences()))
        for i, v := range m.GetAuthenticationContextClassReferences() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("authenticationContextClassReferences", cast)
        if err != nil {
            return err
        }
    }
    if m.GetNamedLocations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetNamedLocations()))
        for i, v := range m.GetNamedLocations() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("namedLocations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPolicies()))
        for i, v := range m.GetPolicies() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("policies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTemplates() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTemplates()))
        for i, v := range m.GetTemplates() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("templates", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAuthenticationContextClassReferences sets the authenticationContextClassReferences property value. Read-only. Nullable. Returns a collection of the specified authentication context class references.
func (m *ConditionalAccessRoot) SetAuthenticationContextClassReferences(value []AuthenticationContextClassReferenceable)() {
    err := m.GetBackingStore().Set("authenticationContextClassReferences", value)
    if err != nil {
        panic(err)
    }
}
// SetNamedLocations sets the namedLocations property value. Read-only. Nullable. Returns a collection of the specified named locations.
func (m *ConditionalAccessRoot) SetNamedLocations(value []NamedLocationable)() {
    err := m.GetBackingStore().Set("namedLocations", value)
    if err != nil {
        panic(err)
    }
}
// SetPolicies sets the policies property value. Read-only. Nullable. Returns a collection of the specified Conditional Access (CA) policies.
func (m *ConditionalAccessRoot) SetPolicies(value []ConditionalAccessPolicyable)() {
    err := m.GetBackingStore().Set("policies", value)
    if err != nil {
        panic(err)
    }
}
// SetTemplates sets the templates property value. Read-only. Nullable. Returns a collection of the specified Conditional Access templates.
func (m *ConditionalAccessRoot) SetTemplates(value []ConditionalAccessTemplateable)() {
    err := m.GetBackingStore().Set("templates", value)
    if err != nil {
        panic(err)
    }
}
// ConditionalAccessRootable 
type ConditionalAccessRootable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthenticationContextClassReferences()([]AuthenticationContextClassReferenceable)
    GetNamedLocations()([]NamedLocationable)
    GetPolicies()([]ConditionalAccessPolicyable)
    GetTemplates()([]ConditionalAccessTemplateable)
    SetAuthenticationContextClassReferences(value []AuthenticationContextClassReferenceable)()
    SetNamedLocations(value []NamedLocationable)()
    SetPolicies(value []ConditionalAccessPolicyable)()
    SetTemplates(value []ConditionalAccessTemplateable)()
}
