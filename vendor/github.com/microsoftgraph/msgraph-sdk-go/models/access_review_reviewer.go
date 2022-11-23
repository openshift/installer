package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AccessReviewReviewer provides operations to manage the collection of agreement entities.
type AccessReviewReviewer struct {
    Entity
    // The date when the reviewer was added for the access review.
    createdDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // Name of reviewer.
    displayName *string
    // User principal name of the reviewer.
    userPrincipalName *string
}
// NewAccessReviewReviewer instantiates a new accessReviewReviewer and sets the default values.
func NewAccessReviewReviewer()(*AccessReviewReviewer) {
    m := &AccessReviewReviewer{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAccessReviewReviewerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAccessReviewReviewerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewReviewer(), nil
}
// GetCreatedDateTime gets the createdDateTime property value. The date when the reviewer was added for the access review.
func (m *AccessReviewReviewer) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.createdDateTime
}
// GetDisplayName gets the displayName property value. Name of reviewer.
func (m *AccessReviewReviewer) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AccessReviewReviewer) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["createdDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetCreatedDateTime)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["userPrincipalName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserPrincipalName)
    return res
}
// GetUserPrincipalName gets the userPrincipalName property value. User principal name of the reviewer.
func (m *AccessReviewReviewer) GetUserPrincipalName()(*string) {
    return m.userPrincipalName
}
// Serialize serializes information the current object
func (m *AccessReviewReviewer) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userPrincipalName", m.GetUserPrincipalName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCreatedDateTime sets the createdDateTime property value. The date when the reviewer was added for the access review.
func (m *AccessReviewReviewer) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.createdDateTime = value
}
// SetDisplayName sets the displayName property value. Name of reviewer.
func (m *AccessReviewReviewer) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetUserPrincipalName sets the userPrincipalName property value. User principal name of the reviewer.
func (m *AccessReviewReviewer) SetUserPrincipalName(value *string)() {
    m.userPrincipalName = value
}
