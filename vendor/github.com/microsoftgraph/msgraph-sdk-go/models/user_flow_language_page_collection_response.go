package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserFlowLanguagePageCollectionResponse 
type UserFlowLanguagePageCollectionResponse struct {
    BaseCollectionPaginationCountResponse
    // The value property
    value []UserFlowLanguagePageable
}
// NewUserFlowLanguagePageCollectionResponse instantiates a new UserFlowLanguagePageCollectionResponse and sets the default values.
func NewUserFlowLanguagePageCollectionResponse()(*UserFlowLanguagePageCollectionResponse) {
    m := &UserFlowLanguagePageCollectionResponse{
        BaseCollectionPaginationCountResponse: *NewBaseCollectionPaginationCountResponse(),
    }
    return m
}
// CreateUserFlowLanguagePageCollectionResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateUserFlowLanguagePageCollectionResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserFlowLanguagePageCollectionResponse(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UserFlowLanguagePageCollectionResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseCollectionPaginationCountResponse.GetFieldDeserializers()
    res["value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateUserFlowLanguagePageFromDiscriminatorValue , m.SetValue)
    return res
}
// GetValue gets the value property value. The value property
func (m *UserFlowLanguagePageCollectionResponse) GetValue()([]UserFlowLanguagePageable) {
    return m.value
}
// Serialize serializes information the current object
func (m *UserFlowLanguagePageCollectionResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseCollectionPaginationCountResponse.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetValue() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetValue())
        err = writer.WriteCollectionOfObjectValues("value", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetValue sets the value property value. The value property
func (m *UserFlowLanguagePageCollectionResponse) SetValue(value []UserFlowLanguagePageable)() {
    m.value = value
}
