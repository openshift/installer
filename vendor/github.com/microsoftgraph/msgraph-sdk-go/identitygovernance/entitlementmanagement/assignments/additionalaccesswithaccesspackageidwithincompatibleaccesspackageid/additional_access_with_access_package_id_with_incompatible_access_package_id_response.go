package additionalaccesswithaccesspackageidwithincompatibleaccesspackageid

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// AdditionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponse provides operations to call the additionalAccess method.
type AdditionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponse struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.BaseCollectionPaginationCountResponse
    // The value property
    value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageAssignmentable
}
// NewAdditionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponse instantiates a new additionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponse and sets the default values.
func NewAdditionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponse()(*AdditionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponse) {
    m := &AdditionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponse{
        BaseCollectionPaginationCountResponse: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewBaseCollectionPaginationCountResponse(),
    }
    return m
}
// CreateAdditionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateAdditionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAdditionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponse(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *AdditionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseCollectionPaginationCountResponse.GetFieldDeserializers()
    res["value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAccessPackageAssignmentFromDiscriminatorValue , m.SetValue)
    return res
}
// GetValue gets the value property value. The value property
func (m *AdditionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponse) GetValue()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageAssignmentable) {
    return m.value
}
// Serialize serializes information the current object
func (m *AdditionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
func (m *AdditionalAccessWithAccessPackageIdWithIncompatibleAccessPackageIdResponse) SetValue(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AccessPackageAssignmentable)() {
    m.value = value
}
