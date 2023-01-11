package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ApplicationTemplate provides operations to manage the collection of applicationTemplate entities.
type ApplicationTemplate struct {
    Entity
    // The list of categories for the application. Supported values can be: Collaboration, Business Management, Consumer, Content management, CRM, Data services, Developer services, E-commerce, Education, ERP, Finance, Health, Human resources, IT infrastructure, Mail, Management, Marketing, Media, Productivity, Project management, Telecommunications, Tools, Travel, and Web design & hosting.
    categories []string
    // A description of the application.
    description *string
    // The name of the application.
    displayName *string
    // The home page URL of the application.
    homePageUrl *string
    // The URL to get the logo for this application.
    logoUrl *string
    // The name of the publisher for this application.
    publisher *string
    // The list of provisioning modes supported by this application. The only valid value is sync.
    supportedProvisioningTypes []string
    // The list of single sign-on modes supported by this application. The supported values are oidc, password, saml, and notSupported.
    supportedSingleSignOnModes []string
}
// NewApplicationTemplate instantiates a new applicationTemplate and sets the default values.
func NewApplicationTemplate()(*ApplicationTemplate) {
    m := &ApplicationTemplate{
        Entity: *NewEntity(),
    }
    return m
}
// CreateApplicationTemplateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateApplicationTemplateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewApplicationTemplate(), nil
}
// GetCategories gets the categories property value. The list of categories for the application. Supported values can be: Collaboration, Business Management, Consumer, Content management, CRM, Data services, Developer services, E-commerce, Education, ERP, Finance, Health, Human resources, IT infrastructure, Mail, Management, Marketing, Media, Productivity, Project management, Telecommunications, Tools, Travel, and Web design & hosting.
func (m *ApplicationTemplate) GetCategories()([]string) {
    return m.categories
}
// GetDescription gets the description property value. A description of the application.
func (m *ApplicationTemplate) GetDescription()(*string) {
    return m.description
}
// GetDisplayName gets the displayName property value. The name of the application.
func (m *ApplicationTemplate) GetDisplayName()(*string) {
    return m.displayName
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ApplicationTemplate) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["categories"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetCategories)
    res["description"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDescription)
    res["displayName"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetDisplayName)
    res["homePageUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetHomePageUrl)
    res["logoUrl"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetLogoUrl)
    res["publisher"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetStringValue(m.SetPublisher)
    res["supportedProvisioningTypes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetSupportedProvisioningTypes)
    res["supportedSingleSignOnModes"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfPrimitiveValues("string" , m.SetSupportedSingleSignOnModes)
    return res
}
// GetHomePageUrl gets the homePageUrl property value. The home page URL of the application.
func (m *ApplicationTemplate) GetHomePageUrl()(*string) {
    return m.homePageUrl
}
// GetLogoUrl gets the logoUrl property value. The URL to get the logo for this application.
func (m *ApplicationTemplate) GetLogoUrl()(*string) {
    return m.logoUrl
}
// GetPublisher gets the publisher property value. The name of the publisher for this application.
func (m *ApplicationTemplate) GetPublisher()(*string) {
    return m.publisher
}
// GetSupportedProvisioningTypes gets the supportedProvisioningTypes property value. The list of provisioning modes supported by this application. The only valid value is sync.
func (m *ApplicationTemplate) GetSupportedProvisioningTypes()([]string) {
    return m.supportedProvisioningTypes
}
// GetSupportedSingleSignOnModes gets the supportedSingleSignOnModes property value. The list of single sign-on modes supported by this application. The supported values are oidc, password, saml, and notSupported.
func (m *ApplicationTemplate) GetSupportedSingleSignOnModes()([]string) {
    return m.supportedSingleSignOnModes
}
// Serialize serializes information the current object
func (m *ApplicationTemplate) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCategories() != nil {
        err = writer.WriteCollectionOfStringValues("categories", m.GetCategories())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
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
        err = writer.WriteStringValue("homePageUrl", m.GetHomePageUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("logoUrl", m.GetLogoUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("publisher", m.GetPublisher())
        if err != nil {
            return err
        }
    }
    if m.GetSupportedProvisioningTypes() != nil {
        err = writer.WriteCollectionOfStringValues("supportedProvisioningTypes", m.GetSupportedProvisioningTypes())
        if err != nil {
            return err
        }
    }
    if m.GetSupportedSingleSignOnModes() != nil {
        err = writer.WriteCollectionOfStringValues("supportedSingleSignOnModes", m.GetSupportedSingleSignOnModes())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCategories sets the categories property value. The list of categories for the application. Supported values can be: Collaboration, Business Management, Consumer, Content management, CRM, Data services, Developer services, E-commerce, Education, ERP, Finance, Health, Human resources, IT infrastructure, Mail, Management, Marketing, Media, Productivity, Project management, Telecommunications, Tools, Travel, and Web design & hosting.
func (m *ApplicationTemplate) SetCategories(value []string)() {
    m.categories = value
}
// SetDescription sets the description property value. A description of the application.
func (m *ApplicationTemplate) SetDescription(value *string)() {
    m.description = value
}
// SetDisplayName sets the displayName property value. The name of the application.
func (m *ApplicationTemplate) SetDisplayName(value *string)() {
    m.displayName = value
}
// SetHomePageUrl sets the homePageUrl property value. The home page URL of the application.
func (m *ApplicationTemplate) SetHomePageUrl(value *string)() {
    m.homePageUrl = value
}
// SetLogoUrl sets the logoUrl property value. The URL to get the logo for this application.
func (m *ApplicationTemplate) SetLogoUrl(value *string)() {
    m.logoUrl = value
}
// SetPublisher sets the publisher property value. The name of the publisher for this application.
func (m *ApplicationTemplate) SetPublisher(value *string)() {
    m.publisher = value
}
// SetSupportedProvisioningTypes sets the supportedProvisioningTypes property value. The list of provisioning modes supported by this application. The only valid value is sync.
func (m *ApplicationTemplate) SetSupportedProvisioningTypes(value []string)() {
    m.supportedProvisioningTypes = value
}
// SetSupportedSingleSignOnModes sets the supportedSingleSignOnModes property value. The list of single sign-on modes supported by this application. The supported values are oidc, password, saml, and notSupported.
func (m *ApplicationTemplate) SetSupportedSingleSignOnModes(value []string)() {
    m.supportedSingleSignOnModes = value
}
