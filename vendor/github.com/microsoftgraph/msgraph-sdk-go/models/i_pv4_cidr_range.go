package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IPv4CidrRange 
type IPv4CidrRange struct {
    IpRange
    // IPv4 address in CIDR notation. Not nullable.
    cidrAddress *string
}
// NewIPv4CidrRange instantiates a new IPv4CidrRange and sets the default values.
func NewIPv4CidrRange()(*IPv4CidrRange) {
    m := &IPv4CidrRange{
        IpRange: *NewIpRange(),
    }
    odataTypeValue := "#microsoft.graph.iPv4CidrRange";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateIPv4CidrRangeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateIPv4CidrRangeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIPv4CidrRange(), nil
}
// GetCidrAddress gets the cidrAddress property value. IPv4 address in CIDR notation. Not nullable.
func (m *IPv4CidrRange) GetCidrAddress()(*string) {
    return m.cidrAddress
}
// GetFieldDeserializers the deserialization information for the current model
func (m *IPv4CidrRange) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.IpRange.GetFieldDeserializers()
    res["cidrAddress"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCidrAddress)
    return res
}
// Serialize serializes information the current object
func (m *IPv4CidrRange) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.IpRange.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("cidrAddress", m.GetCidrAddress())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCidrAddress sets the cidrAddress property value. IPv4 address in CIDR notation. Not nullable.
func (m *IPv4CidrRange) SetCidrAddress(value *string)() {
    m.cidrAddress = value
}
