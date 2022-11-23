package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// InvitedUserMessageInfo 
type InvitedUserMessageInfo struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]interface{}
    // Additional recipients the invitation message should be sent to. Currently only 1 additional recipient is supported.
    ccRecipients []Recipientable
    // Customized message body you want to send if you don't want the default message.
    customizedMessageBody *string
    // The language you want to send the default message in. If the customizedMessageBody is specified, this property is ignored, and the message is sent using the customizedMessageBody. The language format should be in ISO 639. The default is en-US.
    messageLanguage *string
    // The OdataType property
    odataType *string
}
// NewInvitedUserMessageInfo instantiates a new invitedUserMessageInfo and sets the default values.
func NewInvitedUserMessageInfo()(*InvitedUserMessageInfo) {
    m := &InvitedUserMessageInfo{
    }
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateInvitedUserMessageInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateInvitedUserMessageInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewInvitedUserMessageInfo(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *InvitedUserMessageInfo) GetAdditionalData()(map[string]interface{}) {
    return m.additionalData
}
// GetCcRecipients gets the ccRecipients property value. Additional recipients the invitation message should be sent to. Currently only 1 additional recipient is supported.
func (m *InvitedUserMessageInfo) GetCcRecipients()([]Recipientable) {
    return m.ccRecipients
}
// GetCustomizedMessageBody gets the customizedMessageBody property value. Customized message body you want to send if you don't want the default message.
func (m *InvitedUserMessageInfo) GetCustomizedMessageBody()(*string) {
    return m.customizedMessageBody
}
// GetFieldDeserializers the deserialization information for the current model
func (m *InvitedUserMessageInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["ccRecipients"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateRecipientFromDiscriminatorValue , m.SetCcRecipients)
    res["customizedMessageBody"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetCustomizedMessageBody)
    res["messageLanguage"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetMessageLanguage)
    res["@odata.type"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetOdataType)
    return res
}
// GetMessageLanguage gets the messageLanguage property value. The language you want to send the default message in. If the customizedMessageBody is specified, this property is ignored, and the message is sent using the customizedMessageBody. The language format should be in ISO 639. The default is en-US.
func (m *InvitedUserMessageInfo) GetMessageLanguage()(*string) {
    return m.messageLanguage
}
// GetOdataType gets the @odata.type property value. The OdataType property
func (m *InvitedUserMessageInfo) GetOdataType()(*string) {
    return m.odataType
}
// Serialize serializes information the current object
func (m *InvitedUserMessageInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetCcRecipients() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetCcRecipients())
        err := writer.WriteCollectionOfObjectValues("ccRecipients", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("customizedMessageBody", m.GetCustomizedMessageBody())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("messageLanguage", m.GetMessageLanguage())
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
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *InvitedUserMessageInfo) SetAdditionalData(value map[string]interface{})() {
    m.additionalData = value
}
// SetCcRecipients sets the ccRecipients property value. Additional recipients the invitation message should be sent to. Currently only 1 additional recipient is supported.
func (m *InvitedUserMessageInfo) SetCcRecipients(value []Recipientable)() {
    m.ccRecipients = value
}
// SetCustomizedMessageBody sets the customizedMessageBody property value. Customized message body you want to send if you don't want the default message.
func (m *InvitedUserMessageInfo) SetCustomizedMessageBody(value *string)() {
    m.customizedMessageBody = value
}
// SetMessageLanguage sets the messageLanguage property value. The language you want to send the default message in. If the customizedMessageBody is specified, this property is ignored, and the message is sent using the customizedMessageBody. The language format should be in ISO 639. The default is en-US.
func (m *InvitedUserMessageInfo) SetMessageLanguage(value *string)() {
    m.messageLanguage = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *InvitedUserMessageInfo) SetOdataType(value *string)() {
    m.odataType = value
}
