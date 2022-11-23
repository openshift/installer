package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CloudAppSecuritySessionControl 
type CloudAppSecuritySessionControl struct {
    ConditionalAccessSessionControl
    // Possible values are: mcasConfigured, monitorOnly, blockDownloads, unknownFutureValue. For more information, see Deploy Conditional Access App Control for featured apps.
    cloudAppSecurityType *CloudAppSecuritySessionControlType
}
// NewCloudAppSecuritySessionControl instantiates a new CloudAppSecuritySessionControl and sets the default values.
func NewCloudAppSecuritySessionControl()(*CloudAppSecuritySessionControl) {
    m := &CloudAppSecuritySessionControl{
        ConditionalAccessSessionControl: *NewConditionalAccessSessionControl(),
    }
    odataTypeValue := "#microsoft.graph.cloudAppSecuritySessionControl";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateCloudAppSecuritySessionControlFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateCloudAppSecuritySessionControlFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCloudAppSecuritySessionControl(), nil
}
// GetCloudAppSecurityType gets the cloudAppSecurityType property value. Possible values are: mcasConfigured, monitorOnly, blockDownloads, unknownFutureValue. For more information, see Deploy Conditional Access App Control for featured apps.
func (m *CloudAppSecuritySessionControl) GetCloudAppSecurityType()(*CloudAppSecuritySessionControlType) {
    return m.cloudAppSecurityType
}
// GetFieldDeserializers the deserialization information for the current model
func (m *CloudAppSecuritySessionControl) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ConditionalAccessSessionControl.GetFieldDeserializers()
    res["cloudAppSecurityType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseCloudAppSecuritySessionControlType , m.SetCloudAppSecurityType)
    return res
}
// Serialize serializes information the current object
func (m *CloudAppSecuritySessionControl) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ConditionalAccessSessionControl.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCloudAppSecurityType() != nil {
        cast := (*m.GetCloudAppSecurityType()).String()
        err = writer.WriteStringValue("cloudAppSecurityType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCloudAppSecurityType sets the cloudAppSecurityType property value. Possible values are: mcasConfigured, monitorOnly, blockDownloads, unknownFutureValue. For more information, see Deploy Conditional Access App Control for featured apps.
func (m *CloudAppSecuritySessionControl) SetCloudAppSecurityType(value *CloudAppSecuritySessionControlType)() {
    m.cloudAppSecurityType = value
}
