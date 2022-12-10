package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Invitation 
type Invitation struct {
    Entity
    // The user created as part of the invitation creation. Read-Only
    invitedUser Userable
    // The display name of the user being invited.
    invitedUserDisplayName *string
    // The email address of the user being invited. Required. The following special characters are not permitted in the email address:Tilde (~)Exclamation point (!)Number sign (#)Dollar sign ($)Percent (%)Circumflex (^)Ampersand (&)Asterisk (*)Parentheses (( ))Plus sign (+)Equal sign (=)Brackets ([ ])Braces ({ })Backslash (/)Slash mark (/)Pipe (/|)Semicolon (;)Colon (:)Quotation marks (')Angle brackets (< >)Question mark (?)Comma (,)However, the following exceptions apply:A period (.) or a hyphen (-) is permitted anywhere in the user name, except at the beginning or end of the name.An underscore (_) is permitted anywhere in the user name. This includes at the beginning or end of the name.
    invitedUserEmailAddress *string
    // Additional configuration for the message being sent to the invited user, including customizing message text, language and cc recipient list.
    invitedUserMessageInfo InvitedUserMessageInfoable
    // The userType of the user being invited. By default, this is Guest. You can invite as Member if you are a company administrator.
    invitedUserType *string
    // The URL the user can use to redeem their invitation. Read-only.
    inviteRedeemUrl *string
    // The URL the user should be redirected to once the invitation is redeemed. Required.
    inviteRedirectUrl *string
    // The resetRedemption property
    resetRedemption *bool
    // Indicates whether an email should be sent to the user being invited. The default is false.
    sendInvitationMessage *bool
    // The status of the invitation. Possible values are: PendingAcceptance, Completed, InProgress, and Error.
    status *string
}
// NewInvitation instantiates a new Invitation and sets the default values.
func NewInvitation()(*Invitation) {
    m := &Invitation{
        Entity: *NewEntity(),
    }
    return m
}
// CreateInvitationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateInvitationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewInvitation(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *Invitation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["invitedUser"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateUserFromDiscriminatorValue , m.SetInvitedUser)
    res["invitedUserDisplayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetInvitedUserDisplayName)
    res["invitedUserEmailAddress"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetInvitedUserEmailAddress)
    res["invitedUserMessageInfo"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetObjectValue(CreateInvitedUserMessageInfoFromDiscriminatorValue , m.SetInvitedUserMessageInfo)
    res["invitedUserType"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetInvitedUserType)
    res["inviteRedeemUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetInviteRedeemUrl)
    res["inviteRedirectUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetInviteRedirectUrl)
    res["resetRedemption"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetResetRedemption)
    res["sendInvitationMessage"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetBoolValue(m.SetSendInvitationMessage)
    res["status"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetStatus)
    return res
}
// GetInvitedUser gets the invitedUser property value. The user created as part of the invitation creation. Read-Only
func (m *Invitation) GetInvitedUser()(Userable) {
    return m.invitedUser
}
// GetInvitedUserDisplayName gets the invitedUserDisplayName property value. The display name of the user being invited.
func (m *Invitation) GetInvitedUserDisplayName()(*string) {
    return m.invitedUserDisplayName
}
// GetInvitedUserEmailAddress gets the invitedUserEmailAddress property value. The email address of the user being invited. Required. The following special characters are not permitted in the email address:Tilde (~)Exclamation point (!)Number sign (#)Dollar sign ($)Percent (%)Circumflex (^)Ampersand (&)Asterisk (*)Parentheses (( ))Plus sign (+)Equal sign (=)Brackets ([ ])Braces ({ })Backslash (/)Slash mark (/)Pipe (/|)Semicolon (;)Colon (:)Quotation marks (')Angle brackets (< >)Question mark (?)Comma (,)However, the following exceptions apply:A period (.) or a hyphen (-) is permitted anywhere in the user name, except at the beginning or end of the name.An underscore (_) is permitted anywhere in the user name. This includes at the beginning or end of the name.
func (m *Invitation) GetInvitedUserEmailAddress()(*string) {
    return m.invitedUserEmailAddress
}
// GetInvitedUserMessageInfo gets the invitedUserMessageInfo property value. Additional configuration for the message being sent to the invited user, including customizing message text, language and cc recipient list.
func (m *Invitation) GetInvitedUserMessageInfo()(InvitedUserMessageInfoable) {
    return m.invitedUserMessageInfo
}
// GetInvitedUserType gets the invitedUserType property value. The userType of the user being invited. By default, this is Guest. You can invite as Member if you are a company administrator.
func (m *Invitation) GetInvitedUserType()(*string) {
    return m.invitedUserType
}
// GetInviteRedeemUrl gets the inviteRedeemUrl property value. The URL the user can use to redeem their invitation. Read-only.
func (m *Invitation) GetInviteRedeemUrl()(*string) {
    return m.inviteRedeemUrl
}
// GetInviteRedirectUrl gets the inviteRedirectUrl property value. The URL the user should be redirected to once the invitation is redeemed. Required.
func (m *Invitation) GetInviteRedirectUrl()(*string) {
    return m.inviteRedirectUrl
}
// GetResetRedemption gets the resetRedemption property value. The resetRedemption property
func (m *Invitation) GetResetRedemption()(*bool) {
    return m.resetRedemption
}
// GetSendInvitationMessage gets the sendInvitationMessage property value. Indicates whether an email should be sent to the user being invited. The default is false.
func (m *Invitation) GetSendInvitationMessage()(*bool) {
    return m.sendInvitationMessage
}
// GetStatus gets the status property value. The status of the invitation. Possible values are: PendingAcceptance, Completed, InProgress, and Error.
func (m *Invitation) GetStatus()(*string) {
    return m.status
}
// Serialize serializes information the current object
func (m *Invitation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("invitedUser", m.GetInvitedUser())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("invitedUserDisplayName", m.GetInvitedUserDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("invitedUserEmailAddress", m.GetInvitedUserEmailAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("invitedUserMessageInfo", m.GetInvitedUserMessageInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("invitedUserType", m.GetInvitedUserType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("inviteRedeemUrl", m.GetInviteRedeemUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("inviteRedirectUrl", m.GetInviteRedirectUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("resetRedemption", m.GetResetRedemption())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("sendInvitationMessage", m.GetSendInvitationMessage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("status", m.GetStatus())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetInvitedUser sets the invitedUser property value. The user created as part of the invitation creation. Read-Only
func (m *Invitation) SetInvitedUser(value Userable)() {
    m.invitedUser = value
}
// SetInvitedUserDisplayName sets the invitedUserDisplayName property value. The display name of the user being invited.
func (m *Invitation) SetInvitedUserDisplayName(value *string)() {
    m.invitedUserDisplayName = value
}
// SetInvitedUserEmailAddress sets the invitedUserEmailAddress property value. The email address of the user being invited. Required. The following special characters are not permitted in the email address:Tilde (~)Exclamation point (!)Number sign (#)Dollar sign ($)Percent (%)Circumflex (^)Ampersand (&)Asterisk (*)Parentheses (( ))Plus sign (+)Equal sign (=)Brackets ([ ])Braces ({ })Backslash (/)Slash mark (/)Pipe (/|)Semicolon (;)Colon (:)Quotation marks (')Angle brackets (< >)Question mark (?)Comma (,)However, the following exceptions apply:A period (.) or a hyphen (-) is permitted anywhere in the user name, except at the beginning or end of the name.An underscore (_) is permitted anywhere in the user name. This includes at the beginning or end of the name.
func (m *Invitation) SetInvitedUserEmailAddress(value *string)() {
    m.invitedUserEmailAddress = value
}
// SetInvitedUserMessageInfo sets the invitedUserMessageInfo property value. Additional configuration for the message being sent to the invited user, including customizing message text, language and cc recipient list.
func (m *Invitation) SetInvitedUserMessageInfo(value InvitedUserMessageInfoable)() {
    m.invitedUserMessageInfo = value
}
// SetInvitedUserType sets the invitedUserType property value. The userType of the user being invited. By default, this is Guest. You can invite as Member if you are a company administrator.
func (m *Invitation) SetInvitedUserType(value *string)() {
    m.invitedUserType = value
}
// SetInviteRedeemUrl sets the inviteRedeemUrl property value. The URL the user can use to redeem their invitation. Read-only.
func (m *Invitation) SetInviteRedeemUrl(value *string)() {
    m.inviteRedeemUrl = value
}
// SetInviteRedirectUrl sets the inviteRedirectUrl property value. The URL the user should be redirected to once the invitation is redeemed. Required.
func (m *Invitation) SetInviteRedirectUrl(value *string)() {
    m.inviteRedirectUrl = value
}
// SetResetRedemption sets the resetRedemption property value. The resetRedemption property
func (m *Invitation) SetResetRedemption(value *bool)() {
    m.resetRedemption = value
}
// SetSendInvitationMessage sets the sendInvitationMessage property value. Indicates whether an email should be sent to the user being invited. The default is false.
func (m *Invitation) SetSendInvitationMessage(value *bool)() {
    m.sendInvitationMessage = value
}
// SetStatus sets the status property value. The status of the invitation. Possible values are: PendingAcceptance, Completed, InProgress, and Error.
func (m *Invitation) SetStatus(value *string)() {
    m.status = value
}
