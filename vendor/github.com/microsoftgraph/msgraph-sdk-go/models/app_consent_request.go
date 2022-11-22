package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AppConsentRequest provides operations to manage the collection of agreement entities.
type AppConsentRequest struct {
    Entity
    // The display name of the app for which consent is requested. Required. Supports $filter (eq only) and $orderby.
    appDisplayName *string
    // The identifier of the application. Required. Supports $filter (eq only) and $orderby.
    appId *string
    // A list of pending scopes waiting for approval. Required.
    pendingScopes []AppConsentRequestScopeable
    // A list of pending user consent requests. Supports $filter (eq).
    userConsentRequests []UserConsentRequestable
}
// NewAppConsentRequest instantiates a new appConsentRequest and sets the default values.
func NewAppConsentRequest()(*AppConsentRequest) {
    m := &AppConsentRequest{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAppConsentRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAppConsentRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppConsentRequest(), nil
}
// GetAppDisplayName gets the appDisplayName property value. The display name of the app for which consent is requested. Required. Supports $filter (eq only) and $orderby.
func (m *AppConsentRequest) GetAppDisplayName()(*string) {
    return m.appDisplayName
}
// GetAppId gets the appId property value. The identifier of the application. Required. Supports $filter (eq only) and $orderby.
func (m *AppConsentRequest) GetAppId()(*string) {
    return m.appId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AppConsentRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["appDisplayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAppDisplayName)
    res["appId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAppId)
    res["pendingScopes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateAppConsentRequestScopeFromDiscriminatorValue , m.SetPendingScopes)
    res["userConsentRequests"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateUserConsentRequestFromDiscriminatorValue , m.SetUserConsentRequests)
    return res
}
// GetPendingScopes gets the pendingScopes property value. A list of pending scopes waiting for approval. Required.
func (m *AppConsentRequest) GetPendingScopes()([]AppConsentRequestScopeable) {
    return m.pendingScopes
}
// GetUserConsentRequests gets the userConsentRequests property value. A list of pending user consent requests. Supports $filter (eq).
func (m *AppConsentRequest) GetUserConsentRequests()([]UserConsentRequestable) {
    return m.userConsentRequests
}
// Serialize serializes information the current object
func (m *AppConsentRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("appDisplayName", m.GetAppDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appId", m.GetAppId())
        if err != nil {
            return err
        }
    }
    if m.GetPendingScopes() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetPendingScopes())
        err = writer.WriteCollectionOfObjectValues("pendingScopes", cast)
        if err != nil {
            return err
        }
    }
    if m.GetUserConsentRequests() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetUserConsentRequests())
        err = writer.WriteCollectionOfObjectValues("userConsentRequests", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppDisplayName sets the appDisplayName property value. The display name of the app for which consent is requested. Required. Supports $filter (eq only) and $orderby.
func (m *AppConsentRequest) SetAppDisplayName(value *string)() {
    m.appDisplayName = value
}
// SetAppId sets the appId property value. The identifier of the application. Required. Supports $filter (eq only) and $orderby.
func (m *AppConsentRequest) SetAppId(value *string)() {
    m.appId = value
}
// SetPendingScopes sets the pendingScopes property value. A list of pending scopes waiting for approval. Required.
func (m *AppConsentRequest) SetPendingScopes(value []AppConsentRequestScopeable)() {
    m.pendingScopes = value
}
// SetUserConsentRequests sets the userConsentRequests property value. A list of pending user consent requests. Supports $filter (eq).
func (m *AppConsentRequest) SetUserConsentRequests(value []UserConsentRequestable)() {
    m.userConsentRequests = value
}
