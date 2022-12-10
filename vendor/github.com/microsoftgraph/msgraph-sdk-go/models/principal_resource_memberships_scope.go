package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PrincipalResourceMembershipsScope 
type PrincipalResourceMembershipsScope struct {
    AccessReviewScope
    // Defines the scopes of the principals whose access to resources are reviewed in the access review.
    principalScopes []AccessReviewScopeable
    // Defines the scopes of the resources for which access is reviewed.
    resourceScopes []AccessReviewScopeable
}
// NewPrincipalResourceMembershipsScope instantiates a new PrincipalResourceMembershipsScope and sets the default values.
func NewPrincipalResourceMembershipsScope()(*PrincipalResourceMembershipsScope) {
    m := &PrincipalResourceMembershipsScope{
        AccessReviewScope: *NewAccessReviewScope(),
    }
    odataTypeValue := "#microsoft.graph.principalResourceMembershipsScope";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreatePrincipalResourceMembershipsScopeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreatePrincipalResourceMembershipsScopeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrincipalResourceMembershipsScope(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *PrincipalResourceMembershipsScope) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessReviewScope.GetFieldDeserializers()
    res["principalScopes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAccessReviewScopeFromDiscriminatorValue , m.SetPrincipalScopes)
    res["resourceScopes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAccessReviewScopeFromDiscriminatorValue , m.SetResourceScopes)
    return res
}
// GetPrincipalScopes gets the principalScopes property value. Defines the scopes of the principals whose access to resources are reviewed in the access review.
func (m *PrincipalResourceMembershipsScope) GetPrincipalScopes()([]AccessReviewScopeable) {
    return m.principalScopes
}
// GetResourceScopes gets the resourceScopes property value. Defines the scopes of the resources for which access is reviewed.
func (m *PrincipalResourceMembershipsScope) GetResourceScopes()([]AccessReviewScopeable) {
    return m.resourceScopes
}
// Serialize serializes information the current object
func (m *PrincipalResourceMembershipsScope) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessReviewScope.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetPrincipalScopes() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPrincipalScopes())
        err = writer.WriteCollectionOfObjectValues("principalScopes", cast)
        if err != nil {
            return err
        }
    }
    if m.GetResourceScopes() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetResourceScopes())
        err = writer.WriteCollectionOfObjectValues("resourceScopes", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetPrincipalScopes sets the principalScopes property value. Defines the scopes of the principals whose access to resources are reviewed in the access review.
func (m *PrincipalResourceMembershipsScope) SetPrincipalScopes(value []AccessReviewScopeable)() {
    m.principalScopes = value
}
// SetResourceScopes sets the resourceScopes property value. Defines the scopes of the resources for which access is reviewed.
func (m *PrincipalResourceMembershipsScope) SetResourceScopes(value []AccessReviewScopeable)() {
    m.resourceScopes = value
}
