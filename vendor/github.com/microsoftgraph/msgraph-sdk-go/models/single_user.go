package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// SingleUser 
type SingleUser struct {
    SubjectSet
    // The name of the user in Azure AD. Read only.
    description *string
    // The ID of the user in Azure AD.
    userId *string
}
// NewSingleUser instantiates a new SingleUser and sets the default values.
func NewSingleUser()(*SingleUser) {
    m := &SingleUser{
        SubjectSet: *NewSubjectSet(),
    }
    odataTypeValue := "#microsoft.graph.singleUser";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateSingleUserFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateSingleUserFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSingleUser(), nil
}
// GetDescription gets the description property value. The name of the user in Azure AD. Read only.
func (m *SingleUser) GetDescription()(*string) {
    return m.description
}
// GetFieldDeserializers the deserialization information for the current model
func (m *SingleUser) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.SubjectSet.GetFieldDeserializers()
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDescription)
    res["userId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetUserId)
    return res
}
// GetUserId gets the userId property value. The ID of the user in Azure AD.
func (m *SingleUser) GetUserId()(*string) {
    return m.userId
}
// Serialize serializes information the current object
func (m *SingleUser) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.SubjectSet.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userId", m.GetUserId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDescription sets the description property value. The name of the user in Azure AD. Read only.
func (m *SingleUser) SetDescription(value *string)() {
    m.description = value
}
// SetUserId sets the userId property value. The ID of the user in Azure AD.
func (m *SingleUser) SetUserId(value *string)() {
    m.userId = value
}
