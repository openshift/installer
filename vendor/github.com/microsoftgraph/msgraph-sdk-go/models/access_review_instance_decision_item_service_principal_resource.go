package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AccessReviewInstanceDecisionItemServicePrincipalResource 
type AccessReviewInstanceDecisionItemServicePrincipalResource struct {
    AccessReviewInstanceDecisionItemResource
    // The appId property
    appId *string
}
// NewAccessReviewInstanceDecisionItemServicePrincipalResource instantiates a new AccessReviewInstanceDecisionItemServicePrincipalResource and sets the default values.
func NewAccessReviewInstanceDecisionItemServicePrincipalResource()(*AccessReviewInstanceDecisionItemServicePrincipalResource) {
    m := &AccessReviewInstanceDecisionItemServicePrincipalResource{
        AccessReviewInstanceDecisionItemResource: *NewAccessReviewInstanceDecisionItemResource(),
    }
    odataTypeValue := "#microsoft.graph.accessReviewInstanceDecisionItemServicePrincipalResource";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateAccessReviewInstanceDecisionItemServicePrincipalResourceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAccessReviewInstanceDecisionItemServicePrincipalResourceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewInstanceDecisionItemServicePrincipalResource(), nil
}
// GetAppId gets the appId property value. The appId property
func (m *AccessReviewInstanceDecisionItemServicePrincipalResource) GetAppId()(*string) {
    return m.appId
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AccessReviewInstanceDecisionItemServicePrincipalResource) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessReviewInstanceDecisionItemResource.GetFieldDeserializers()
    res["appId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetAppId)
    return res
}
// Serialize serializes information the current object
func (m *AccessReviewInstanceDecisionItemServicePrincipalResource) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessReviewInstanceDecisionItemResource.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("appId", m.GetAppId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppId sets the appId property value. The appId property
func (m *AccessReviewInstanceDecisionItemServicePrincipalResource) SetAppId(value *string)() {
    m.appId = value
}
