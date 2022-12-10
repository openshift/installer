package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// BookingCurrency 
type BookingCurrency struct {
    Entity
    // The currency symbol. For example, the currency symbol for the US dollar and for the Australian dollar is $.
    symbol *string
}
// NewBookingCurrency instantiates a new BookingCurrency and sets the default values.
func NewBookingCurrency()(*BookingCurrency) {
    m := &BookingCurrency{
        Entity: *NewEntity(),
    }
    return m
}
// CreateBookingCurrencyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateBookingCurrencyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBookingCurrency(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *BookingCurrency) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["symbol"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetSymbol)
    return res
}
// GetSymbol gets the symbol property value. The currency symbol. For example, the currency symbol for the US dollar and for the Australian dollar is $.
func (m *BookingCurrency) GetSymbol()(*string) {
    return m.symbol
}
// Serialize serializes information the current object
func (m *BookingCurrency) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("symbol", m.GetSymbol())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetSymbol sets the symbol property value. The currency symbol. For example, the currency symbol for the US dollar and for the Australian dollar is $.
func (m *BookingCurrency) SetSymbol(value *string)() {
    m.symbol = value
}
