package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// KeyCredentialCollectionResponse 
type KeyCredentialCollectionResponse struct {
    BaseCollectionPaginationCountResponse
    // The value property
    value []KeyCredentialable
}
// NewKeyCredentialCollectionResponse instantiates a new KeyCredentialCollectionResponse and sets the default values.
func NewKeyCredentialCollectionResponse()(*KeyCredentialCollectionResponse) {
    m := &KeyCredentialCollectionResponse{
        BaseCollectionPaginationCountResponse: *NewBaseCollectionPaginationCountResponse(),
    }
    return m
}
// CreateKeyCredentialCollectionResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateKeyCredentialCollectionResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewKeyCredentialCollectionResponse(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *KeyCredentialCollectionResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseCollectionPaginationCountResponse.GetFieldDeserializers()
    res["value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateKeyCredentialFromDiscriminatorValue , m.SetValue)
    return res
}
// GetValue gets the value property value. The value property
func (m *KeyCredentialCollectionResponse) GetValue()([]KeyCredentialable) {
    return m.value
}
// Serialize serializes information the current object
func (m *KeyCredentialCollectionResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
func (m *KeyCredentialCollectionResponse) SetValue(value []KeyCredentialable)() {
    m.value = value
}
