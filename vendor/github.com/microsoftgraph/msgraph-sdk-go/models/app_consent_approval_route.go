package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AppConsentApprovalRoute 
type AppConsentApprovalRoute struct {
    Entity
    // A collection of userConsentRequest objects for a specific application.
    appConsentRequests []AppConsentRequestable
}
// NewAppConsentApprovalRoute instantiates a new AppConsentApprovalRoute and sets the default values.
func NewAppConsentApprovalRoute()(*AppConsentApprovalRoute) {
    m := &AppConsentApprovalRoute{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAppConsentApprovalRouteFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAppConsentApprovalRouteFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppConsentApprovalRoute(), nil
}
// GetAppConsentRequests gets the appConsentRequests property value. A collection of userConsentRequest objects for a specific application.
func (m *AppConsentApprovalRoute) GetAppConsentRequests()([]AppConsentRequestable) {
    return m.appConsentRequests
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AppConsentApprovalRoute) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["appConsentRequests"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAppConsentRequestFromDiscriminatorValue , m.SetAppConsentRequests)
    return res
}
// Serialize serializes information the current object
func (m *AppConsentApprovalRoute) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAppConsentRequests() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetAppConsentRequests())
        err = writer.WriteCollectionOfObjectValues("appConsentRequests", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppConsentRequests sets the appConsentRequests property value. A collection of userConsentRequest objects for a specific application.
func (m *AppConsentApprovalRoute) SetAppConsentRequests(value []AppConsentRequestable)() {
    m.appConsentRequests = value
}
