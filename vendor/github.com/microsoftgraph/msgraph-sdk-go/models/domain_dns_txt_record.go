package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DomainDnsTxtRecord 
type DomainDnsTxtRecord struct {
    DomainDnsRecord
    // Value used when configuring the text property at the DNS host.
    text *string
}
// NewDomainDnsTxtRecord instantiates a new DomainDnsTxtRecord and sets the default values.
func NewDomainDnsTxtRecord()(*DomainDnsTxtRecord) {
    m := &DomainDnsTxtRecord{
        DomainDnsRecord: *NewDomainDnsRecord(),
    }
    return m
}
// CreateDomainDnsTxtRecordFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateDomainDnsTxtRecordFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDomainDnsTxtRecord(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *DomainDnsTxtRecord) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DomainDnsRecord.GetFieldDeserializers()
    res["text"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetText)
    return res
}
// GetText gets the text property value. Value used when configuring the text property at the DNS host.
func (m *DomainDnsTxtRecord) GetText()(*string) {
    return m.text
}
// Serialize serializes information the current object
func (m *DomainDnsTxtRecord) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DomainDnsRecord.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("text", m.GetText())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetText sets the text property value. Value used when configuring the text property at the DNS host.
func (m *DomainDnsTxtRecord) SetText(value *string)() {
    m.text = value
}
