package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// NotificationMessageTemplate notification messages are messages that are sent to end users who are determined to be not-compliant with the compliance policies defined by the administrator. Administrators choose notifications and configure them in the Intune Admin Console using the compliance policy creation page under the “Actions for non-compliance” section. Use the notificationMessageTemplate object to create your own custom notifications for administrators to choose while configuring actions for non-compliance.
type NotificationMessageTemplate struct {
    Entity
    // Branding Options for the Message Template. Branding is defined in the Intune Admin Console.
    brandingOptions *NotificationTemplateBrandingOptions
    // The default locale to fallback onto when the requested locale is not available.
    defaultLocale *string
    // Display name for the Notification Message Template.
    displayName *string
    // DateTime the object was last modified.
    lastModifiedDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
    // The list of localized messages for this Notification Message Template.
    localizedNotificationMessages []LocalizedNotificationMessageable
}
// NewNotificationMessageTemplate instantiates a new notificationMessageTemplate and sets the default values.
func NewNotificationMessageTemplate()(*NotificationMessageTemplate) {
    m := &NotificationMessageTemplate{
        Entity: *NewEntity(),
    }
    return m
}
// CreateNotificationMessageTemplateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateNotificationMessageTemplateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewNotificationMessageTemplate(), nil
}
// GetBrandingOptions gets the brandingOptions property value. Branding Options for the Message Template. Branding is defined in the Intune Admin Console.
func (m *NotificationMessageTemplate) GetBrandingOptions()(*NotificationTemplateBrandingOptions) {
    return m.brandingOptions
}
// GetDefaultLocale gets the defaultLocale property value. The default locale to fallback onto when the requested locale is not available.
func (m *NotificationMessageTemplate) GetDefaultLocale()(*string) {
    return m.defaultLocale
}
// GetDisplayName gets the displayName property value. Display name for the Notification Message Template.
func (m *NotificationMessageTemplate) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *NotificationMessageTemplate) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["brandingOptions"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetEnumValue(ParseNotificationTemplateBrandingOptions , m.SetBrandingOptions)
    res["defaultLocale"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDefaultLocale)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["lastModifiedDateTime"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetTimeValue(m.SetLastModifiedDateTime)
    res["localizedNotificationMessages"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateLocalizedNotificationMessageFromDiscriminatorValue , m.SetLocalizedNotificationMessages)
    return res
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. DateTime the object was last modified.
func (m *NotificationMessageTemplate) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    return m.lastModifiedDateTime
}
// GetLocalizedNotificationMessages gets the localizedNotificationMessages property value. The list of localized messages for this Notification Message Template.
func (m *NotificationMessageTemplate) GetLocalizedNotificationMessages()([]LocalizedNotificationMessageable) {
    return m.localizedNotificationMessages
}
// Serialize serializes information the current object
func (m *NotificationMessageTemplate) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetBrandingOptions() != nil {
        cast := (*m.GetBrandingOptions()).String()
        err = writer.WriteStringValue("brandingOptions", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("defaultLocale", m.GetDefaultLocale())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetLocalizedNotificationMessages() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetLocalizedNotificationMessages())
        err = writer.WriteCollectionOfObjectValues("localizedNotificationMessages", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBrandingOptions sets the brandingOptions property value. Branding Options for the Message Template. Branding is defined in the Intune Admin Console.
func (m *NotificationMessageTemplate) SetBrandingOptions(value *NotificationTemplateBrandingOptions)() {
    m.brandingOptions = value
}
// SetDefaultLocale sets the defaultLocale property value. The default locale to fallback onto when the requested locale is not available.
func (m *NotificationMessageTemplate) SetDefaultLocale(value *string)() {
    m.defaultLocale = value
}
// SetDisplayName sets the displayName property value. Display name for the Notification Message Template.
func (m *NotificationMessageTemplate) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. DateTime the object was last modified.
func (m *NotificationMessageTemplate) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    m.lastModifiedDateTime = value
}
// SetLocalizedNotificationMessages sets the localizedNotificationMessages property value. The list of localized messages for this Notification Message Template.
func (m *NotificationMessageTemplate) SetLocalizedNotificationMessages(value []LocalizedNotificationMessageable)() {
    m.localizedNotificationMessages = value
}
