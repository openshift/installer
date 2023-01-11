package changescreensharingrole

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// ChangeScreenSharingRolePostRequestBody provides operations to call the changeScreenSharingRole method.
type ChangeScreenSharingRolePostRequestBody struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The role property
    role *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ScreenSharingRole
}
// NewChangeScreenSharingRolePostRequestBody instantiates a new changeScreenSharingRolePostRequestBody and sets the default values.
func NewChangeScreenSharingRolePostRequestBody()(*ChangeScreenSharingRolePostRequestBody) {
    m := &ChangeScreenSharingRolePostRequestBody{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateChangeScreenSharingRolePostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateChangeScreenSharingRolePostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewChangeScreenSharingRolePostRequestBody(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ChangeScreenSharingRolePostRequestBody) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ChangeScreenSharingRolePostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["role"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParseScreenSharingRole , m.SetRole)
    return res
}
// GetRole gets the role property value. The role property
func (m *ChangeScreenSharingRolePostRequestBody) GetRole()(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ScreenSharingRole) {
    return m.role
}
// Serialize serializes information the current object
func (m *ChangeScreenSharingRolePostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetRole() != nil {
        cast := (*m.GetRole()).String()
        err := writer.WriteStringValue("role", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ChangeScreenSharingRolePostRequestBody) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetRole sets the role property value. The role property
func (m *ChangeScreenSharingRolePostRequestBody) SetRole(value *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ScreenSharingRole)() {
    m.role = value
}
