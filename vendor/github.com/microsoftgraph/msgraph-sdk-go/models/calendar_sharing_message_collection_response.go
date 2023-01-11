package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CalendarSharingMessageCollectionResponse 
type CalendarSharingMessageCollectionResponse struct {
    BaseCollectionPaginationCountResponse
    // The value property
    value []CalendarSharingMessageable
}
// NewCalendarSharingMessageCollectionResponse instantiates a new CalendarSharingMessageCollectionResponse and sets the default values.
func NewCalendarSharingMessageCollectionResponse()(*CalendarSharingMessageCollectionResponse) {
    m := &CalendarSharingMessageCollectionResponse{
        BaseCollectionPaginationCountResponse: *NewBaseCollectionPaginationCountResponse(),
    }
    return m
}
// CreateCalendarSharingMessageCollectionResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCalendarSharingMessageCollectionResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCalendarSharingMessageCollectionResponse(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CalendarSharingMessageCollectionResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseCollectionPaginationCountResponse.GetFieldDeserializers()
    res["value"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateCalendarSharingMessageFromDiscriminatorValue , m.SetValue)
    return res
}
// GetValue gets the value property value. The value property
func (m *CalendarSharingMessageCollectionResponse) GetValue()([]CalendarSharingMessageable) {
    return m.value
}
// Serialize serializes information the current object
func (m *CalendarSharingMessageCollectionResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
func (m *CalendarSharingMessageCollectionResponse) SetValue(value []CalendarSharingMessageable)() {
    m.value = value
}
