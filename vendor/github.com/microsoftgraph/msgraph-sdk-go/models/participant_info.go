package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ParticipantInfo 
type ParticipantInfo struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // The ISO 3166-1 Alpha-2 country code of the participant's best estimated physical location at the start of the call. Read-only.
    countryCode *string
    // The type of endpoint the participant is using. Possible values are: default, skypeForBusiness, or skypeForBusinessVoipPhone. Read-only.
    endpointType *EndpointType
    // The identity property
    identity IdentitySetable
    // The language culture string. Read-only.
    languageId *string
    // The OdataType property
    odataType *string
    // The participant ID of the participant. Read-only.
    participantId *string
    // The home region of the participant. This can be a country, a continent, or a larger geographic region. This does not change based on the participant's current physical location. Read-only.
    region *string
}
// NewParticipantInfo instantiates a new participantInfo and sets the default values.
func NewParticipantInfo()(*ParticipantInfo) {
    m := &ParticipantInfo{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateParticipantInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateParticipantInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewParticipantInfo(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ParticipantInfo) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCountryCode gets the countryCode property value. The ISO 3166-1 Alpha-2 country code of the participant's best estimated physical location at the start of the call. Read-only.
func (m *ParticipantInfo) GetCountryCode()(*string) {
    return m.countryCode
}
// GetEndpointType gets the endpointType property value. The type of endpoint the participant is using. Possible values are: default, skypeForBusiness, or skypeForBusinessVoipPhone. Read-only.
func (m *ParticipantInfo) GetEndpointType()(*EndpointType) {
    return m.endpointType
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ParticipantInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["countryCode"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCountryCode)
    res["endpointType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseEndpointType , m.SetEndpointType)
    res["identity"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateIdentitySetFromDiscriminatorValue , m.SetIdentity)
    res["languageId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetLanguageId)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    res["participantId"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetParticipantId)
    res["region"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetRegion)
    return res
}
// GetIdentity gets the identity property value. The identity property
func (m *ParticipantInfo) GetIdentity()(IdentitySetable) {
    return m.identity
}
// GetLanguageId gets the languageId property value. The language culture string. Read-only.
func (m *ParticipantInfo) GetLanguageId()(*string) {
    return m.languageId
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ParticipantInfo) GetOdataType()(*string) {
    return m.odataType
}
// GetParticipantId gets the participantId property value. The participant ID of the participant. Read-only.
func (m *ParticipantInfo) GetParticipantId()(*string) {
    return m.participantId
}
// GetRegion gets the region property value. The home region of the participant. This can be a country, a continent, or a larger geographic region. This does not change based on the participant's current physical location. Read-only.
func (m *ParticipantInfo) GetRegion()(*string) {
    return m.region
}
// Serialize serializes information the current object
func (m *ParticipantInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("countryCode", m.GetCountryCode())
        if err != nil {
            return err
        }
    }
    if m.GetEndpointType() != nil {
        cast := (*m.GetEndpointType()).String()
        err := writer.WriteStringValue("endpointType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("identity", m.GetIdentity())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("languageId", m.GetLanguageId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("participantId", m.GetParticipantId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("region", m.GetRegion())
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
func (m *ParticipantInfo) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCountryCode sets the countryCode property value. The ISO 3166-1 Alpha-2 country code of the participant's best estimated physical location at the start of the call. Read-only.
func (m *ParticipantInfo) SetCountryCode(value *string)() {
    m.countryCode = value
}
// SetEndpointType sets the endpointType property value. The type of endpoint the participant is using. Possible values are: default, skypeForBusiness, or skypeForBusinessVoipPhone. Read-only.
func (m *ParticipantInfo) SetEndpointType(value *EndpointType)() {
    m.endpointType = value
}
// SetIdentity sets the identity property value. The identity property
func (m *ParticipantInfo) SetIdentity(value IdentitySetable)() {
    m.identity = value
}
// SetLanguageId sets the languageId property value. The language culture string. Read-only.
func (m *ParticipantInfo) SetLanguageId(value *string)() {
    m.languageId = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ParticipantInfo) SetOdataType(value *string)() {
    m.odataType = value
}
// SetParticipantId sets the participantId property value. The participant ID of the participant. Read-only.
func (m *ParticipantInfo) SetParticipantId(value *string)() {
    m.participantId = value
}
// SetRegion sets the region property value. The home region of the participant. This can be a country, a continent, or a larger geographic region. This does not change based on the participant's current physical location. Read-only.
func (m *ParticipantInfo) SetRegion(value *string)() {
    m.region = value
}
